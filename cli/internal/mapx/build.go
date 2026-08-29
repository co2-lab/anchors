package mapx

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// Build monta o grafo a partir dos arquivos escaneados E da Estrutura declarada
// (a config). O CLI não sabe mais nada sobre TS/TSX ou spec.md — tudo vem da
// config (STRUCTURE §2.1: a Estrutura é o grafo virtual; o mapa é sua projeção).
//
// Arestas por origem:
//   - convention: co-location dos derivados (config.Derived) — spec→código,
//     spec→feature, feature→teste, seguindo os templates de caminho.
//   - inferred (por código de cenário): liga encarnações de um requisito onde a
//     co-location NÃO alcança (cross-target). Tipo = a relação real (tested-by),
//     não references.
//   - declared (governs): as arestas verticais de alto grau, da config.Governs.
//
// Build monta o grafo. `updatedAt` mapeia caminho→data do último commit (o carimbo
// de alteração, TRACEABILITY §rev/updated_at) — vem do git, resolvido pelo chamador
// (o mapx não invoca git para não acoplar). Pode ser nil (updated_at fica vazio).
func Build(files []scan.File, cfg *config.Config, updatedAt map[string]string) *Graph {
	g := &Graph{Version: 1}
	for _, f := range files {
		var tags []string
		var regime string
		if l, ok := cfg.Layers[f.Layer]; ok {
			tags = l.Tags
			regime = l.Regime
		}
		g.Nodes = append(g.Nodes, Node{
			ID:            f.Path,
			Kind:          Kind(f.Kind),
			Rev:           f.Rev,
			UpdatedAt:     updatedAt[f.Path],
			Code:          nodeCode(f),
			CodeDeclarado: codeDeclarado(f),
			Tags:          tags,
			Regime:        regime,
			NoPropagation: f.NoPropagation,
			SharedCode:    f.SharedCode,
			Needs:         f.Needs,
		})
	}

	colo := colocationEdges(files, cfg)
	g.Edges = append(g.Edges, colo...)
	g.Edges = append(g.Edges, scenarioEdges(files, colo)...)
	g.Edges = append(g.Edges, governsEdges(files, cfg)...)
	g.Edges = append(g.Edges, dependsOnEdges(files)...)
	g.Edges = append(g.Edges, seedEdges(files)...)

	sortGraph(g)
	return g
}

// dependsOnEdges constrói as arestas de REÚSO entre camadas (SPEC_TYPES §5): de cada
// spec consumidora que declarou uma Tabela de Dependências, uma aresta `depends-on`
// para o ARQUIVO de cada linha, carregando Method e o código local (DEPn) como
// metadados. Origem `declared` (a spec afirma a dependência; não é co-location nem
// inferência). Só liga a um alvo que EXISTE no grafo — uma linha que aponta arquivo
// inexistente vira aresta morta que o anti-drift do mapa (TRACEABILITY §"não mentir")
// acusa depois; aqui a omitimos para não poluir o grafo com nós fantasma.
func dependsOnEdges(files []scan.File) []Edge {
	exists := map[string]bool{}
	for _, f := range files {
		exists[f.Path] = true
	}
	var edges []Edge
	for _, f := range files {
		for _, d := range f.Deps {
			if d.File == "" || !exists[d.File] {
				continue
			}
			edges = append(edges, Edge{
				From:   f.Path,
				To:     d.File,
				Type:   EdgeDependsOn,
				Origin: OriginDeclared,
				Method: d.Method,
				Dep:    d.Code,
			})
		}
	}
	return edges
}

// seedEdges liga cada PLANO às specs que ele semeia.
//
// Sem elas o plano é órfão do grafo: `anchors impact` respondia "nenhum filho depende
// dele" sobre um plano que nomeia 11 specs, e mudar o plano não propagava para nada — o
// documento que ORIGINA o trabalho ficava fora do mapa do trabalho.
//
// Só liga o que já EXISTE. Um plano cita specs que ainda vão nascer, e uma aresta para o
// vazio seria a âncora que mente: apontaria para um arquivo ausente com a mesma cara de
// uma ligação real. O gate `plan-seeds-valid` é quem confronta a lista inteira.
func seedEdges(files []scan.File) []Edge {
	exists := map[string]bool{}
	for _, f := range files {
		exists[f.Path] = true
	}
	// Índice por NOME do arquivo: um plano cita a spec de duas formas legítimas — pelo
	// caminho (`apps/.../X.spec.md`) ou só pelo nome (`X.spec.md`), que é como se
	// escreve em prosa. Medido num repositório real: 10 das 26 citações são por nome.
	//
	// Tratar "citado por nome" como "não existe" fazia o plano parecer não-cumprido para
	// sempre — e a partida a frio da fila semeava fases antigas, já concluídas.
	porNome := map[string][]string{}
	for _, f := range files {
		base := filepath.Base(f.Path)
		porNome[base] = append(porNome[base], f.Path)
	}

	var edges []Edge
	for _, f := range files {
		for _, alvo := range f.Seeds {
			destino := alvo
			if !exists[destino] {
				// Só resolve por nome quando ele é ÚNICO no repositório. Dois arquivos com
				// o mesmo nome tornam a citação ambígua, e escolher um seria inventar uma
				// aresta que o autor não declarou.
				if cands := porNome[filepath.Base(alvo)]; len(cands) == 1 {
					destino = cands[0]
				} else {
					continue
				}
			}
			edges = append(edges, Edge{From: f.Path, To: destino, Type: EdgeSeeds, Origin: OriginDeclared})
		}
		// `needs:` — a ordem de trabalho entre planos. A aresta aponta do dependente para
		// o pré-requisito, e um alvo que não existe NÃO vira aresta: o doctor a reporta
		// como `needs-quebrado`, que é mais útil do que uma aresta morta.
		for _, alvo := range f.Needs {
			if exists[alvo] {
				edges = append(edges, Edge{From: f.Path, To: alvo, Type: EdgeNeeds, Origin: OriginDeclared})
			}
		}
	}
	return edges
}

// colocationEdges liga os derivados co-localizados usando os templates de
// config.Derived. Agrupa por stem (o {{dir}}/{{name}} da âncora) e liga a trinca.
func colocationEdges(files []scan.File, cfg *config.Config) []Edge {
	if cfg.Derived == nil {
		return nil
	}
	// índice: caminho → arquivo (para achar os derivados esperados)
	byPath := map[string]scan.File{}
	for _, f := range files {
		byPath[f.Path] = f
	}

	var edges []Edge
	for _, f := range files {
		// a âncora casa por KIND (ex.: "code") — pode haver várias layers de código
		// (screen, component, mobile-code…), todas âncoras da co-location.
		if f.Kind != cfg.Derived.Anchor {
			continue
		}
		// ToSlash porque o dir alimenta {{dir}} nos templates dos derivados, e o que sai
		// dali é comparado contra ids do mapa (barra normal). No Windows o filepath.Dir
		// passa por Clean e devolve "functions\run-audits": o caminho montado não casa
		// id nenhum e a co-location inteira deixa de ligar spec/feature/teste.
		dir := filepath.ToSlash(filepath.Dir(f.Path))
		name, ext := stemDaAncora(f.Path)
		module := filepath.Base(dir) // {{module}} — o dir-pai (ex.: run-audits em .../run-audits/handler.ts)
		// Templates efetivos = default sobrescrito pelos overrides cuja `when` casa a
		// camada da âncora (STRUCTURE.md §2.2: padrão de localização por camada quando
		// não co-localizado). Só as camadas presentes no override sobrescrevem.
		tmpls := map[string]string{}
		for layer, tmpl := range cfg.Derived.Files {
			tmpls[layer] = tmpl
		}
		for _, ov := range cfg.Derived.Overrides {
			if ov.When != f.Layer {
				continue
			}
			for layer, tmpl := range ov.Files {
				tmpls[layer] = tmpl
			}
		}
		// resolve o caminho esperado de cada derivado e liga se existir
		derived := map[string]string{} // camada → caminho existente
		for layer, tmpl := range tmpls {
			want := resolveTemplateM(tmpl, dir, name, ext, module)
			if _, ok := byPath[want]; ok {
				derived[layer] = want
			}
		}
		// A ÂNCORA é a spec, e os derivados saem dela. Até a v0.1 a âncora era o código e
		// a spec vinha em `derived["spec"]` — direção invertida em relação à doutrina
		// ("da spec nascem o código, a feature e o teste").
		//
		// O código segue aceito como âncora, e aí a spec volta a ser derivada: um projeto
		// que declarou `anchor: code` não pode perder as arestas por causa desta mudança.
		spec, codigo := f.Path, derived["code"]
		if cfg.Derived.Anchor != "spec" {
			spec, codigo = derived["spec"], f.Path
		}
		feat := derived["feature"]
		test := derived["test"]
		if spec != "" && codigo != "" {
			edges = append(edges, Edge{From: spec, To: codigo, Type: EdgeSpecifies, Origin: OriginConvention})
		}
		if spec != "" && feat != "" {
			edges = append(edges, Edge{From: spec, To: feat, Type: EdgeCoveredBy, Origin: OriginConvention})
		}
		if feat != "" && test != "" {
			edges = append(edges, Edge{From: feat, To: test, Type: EdgeTestedBy, Origin: OriginConvention})
		}
		// Sem feature declarada, o teste ainda deriva da spec. Um projeto que não usa
		// feature (o próprio Anchors é um) perderia a ligação spec→teste de outro modo.
		if feat == "" && spec != "" && test != "" {
			edges = append(edges, Edge{From: spec, To: test, Type: EdgeTestedBy, Origin: OriginConvention})
		}
	}
	return edges
}

// scenarioEdges liga encarnações de um mesmo código de cenário que a co-location
// NÃO ligou (cross-target). Corrige os furos #1/#2 da v1: não duplica arestas de
// co-location, e usa o tipo da RELAÇÃO (tested-by), não references.
func scenarioEdges(files []scan.File, colo []Edge) []Edge {
	// pares já ligados por co-location (para não duplicar)
	linked := map[string]bool{}
	for _, e := range colo {
		linked[e.From+"\x00"+e.To] = true
	}

	type ref struct{ path, kind string }
	byCode := map[string][]ref{}
	for _, f := range files {
		for _, c := range f.Codes {
			byCode[c] = append(byCode[c], ref{f.Path, f.Kind})
		}
	}

	seen := map[string]bool{}
	var edges []Edge
	for _, refs := range byCode {
		var srcs []string // spec/feature
		var tests []string
		for _, r := range refs {
			switch Kind(r.kind) {
			case KindSpec, KindFeature:
				srcs = append(srcs, r.path)
			case KindTest:
				tests = append(tests, r.path)
			}
		}
		for _, s := range srcs {
			for _, t := range tests {
				// pula se co-localizados (mesmo diretório) ou já ligados
				if filepath.Dir(s) == filepath.Dir(t) {
					continue
				}
				key := s + "\x00" + t
				if linked[key] || seen[key] {
					continue
				}
				seen[key] = true
				edges = append(edges, Edge{From: s, To: t, Type: EdgeTestedBy, Origin: OriginInferred})
			}
		}
	}
	return edges
}

// governsEdges cria as arestas verticais declaradas. A régua `from` (um guide)
// rege os nós de todas as layers que carregam a tag `governs` da regra. Sempre por
// TAG — o escopo vem dos patterns dessas layers (DRY), sem glob duplicado. Isso é
// o que evita a explosão cartesiana: cada guide rege só as layers de sua tag.
func governsEdges(files []scan.File, cfg *config.Config) []Edge {
	// tag → conjunto de nomes de layer que a carregam
	layersByTag := map[string]map[string]bool{}
	for name, l := range cfg.Layers {
		for _, tag := range l.Tags {
			if layersByTag[tag] == nil {
				layersByTag[tag] = map[string]bool{}
			}
			layersByTag[tag][name] = true
		}
	}
	// nós por layer
	byLayer := map[string][]string{}
	for _, f := range files {
		byLayer[f.Layer] = append(byLayer[f.Layer], f.Path)
	}

	var edges []Edge
	for _, rule := range cfg.Governs {
		targetLayers := layersByTag[rule.Governs]
		if targetLayers == nil {
			continue // tag não declarada em nenhuma layer — nada a reger
		}
		for layerName := range targetLayers {
			for _, to := range byLayer[layerName] {
				if rule.From == to {
					continue
				}
				edges = append(edges, Edge{From: rule.From, To: to, Type: EdgeGoverns, Origin: OriginDeclared})
			}
		}
	}
	return edges
}

// stemName devolve o nome-base sem extensão e a extensão (sem ponto).
// dir/Login.tsx → ("Login", "tsx").
// stemDaAncora extrai o nome da unidade a partir do caminho da ÂNCORA.
//
// A âncora é a spec (`Login.spec.md`), e dela o nome da unidade é `Login` — não
// `Login.spec`, que é o que um corte na última extensão devolveria. Os sufixos de
// artefato são compostos, e cortá-los é o que faz `{{name}}` valer para os derivados.
//
// O `{{ext}}` sai VAZIO de propósito: a extensão do código não está na spec, e inventá-la
// (`.ts`? `.tsx`? `.go`?) escolheria por um projeto que não declarou. Um template que
// precise dela declara a extensão literalmente — `{{dir}}/{{name}}.test.ts` —, e assim a
// decisão fica escrita onde se lê.
func stemDaAncora(path string) (name, ext string) {
	base := filepath.Base(path)
	for _, suf := range []string{".spec.md", ".feature"} {
		if b, ok := strings.CutSuffix(base, suf); ok {
			return b, ""
		}
	}
	// Não é um artefato de sufixo composto: cai no comportamento antigo, que serve a
	// qualquer âncora que um projeto venha a declarar.
	return stemName(path)
}

func stemName(path string) (name, ext string) {
	base := filepath.Base(path)
	if i := strings.LastIndex(base, "."); i >= 0 {
		return base[:i], base[i+1:]
	}
	return base, ""
}

// resolveTemplate expande {{dir}}, {{name}}, {{ext}} num template de caminho.
func resolveTemplate(tmpl, dir, name, ext string) string {
	return resolveTemplateM(tmpl, dir, name, ext, filepath.Base(dir))
}

// resolveTemplateM expande {{dir}}, {{name}}, {{ext}} e {{module}} (o dir-pai) — usado
// pelos padrões de localização por camada (STRUCTURE.md §2.2), onde a peça derivada
// mora numa região própria e o módulo (não o basename) compõe o caminho esperado.
func resolveTemplateM(tmpl, dir, name, ext, module string) string {
	r := strings.NewReplacer("{{dir}}", dir, "{{name}}", name, "{{ext}}", ext, "{{module}}", module)
	return r.Replace(tmpl)
}

// nodeCode decide a identidade do nó. O header DECLARADO vence: é onde o autor diz de
// quem é o arquivo. A inferência pelo primeiro código do texto é só o fallback para quem
// não declara — e ela erra quando a spec CITA outra unidade antes de definir a sua
// (medido: uma spec de modelo que abria referenciando `DTAXX-B11` entrava no mapa como
// dona de `DTAX`, e os gates relacionais passavam a confrontar a unidade errada).
func nodeCode(f scan.File) string {
	if f.HeaderCode != "" {
		return f.HeaderCode
	}
	return primaryCode(f.Codes)
}

// codeDeclarado diz se a identidade foi DECLARADA (header `code:`) ou apenas inferida do
// texto. Ver Node.CodeDeclarado para o porquê da distinção.
func codeDeclarado(f scan.File) bool { return f.HeaderCode != "" }

func primaryCode(codes []string) string {
	if len(codes) == 0 {
		return ""
	}
	return strings.SplitN(codes[0], "-", 2)[0]
}

func sortGraph(g *Graph) {
	sort.Slice(g.Nodes, func(i, j int) bool { return g.Nodes[i].ID < g.Nodes[j].ID })
	sort.Slice(g.Edges, func(i, j int) bool {
		if g.Edges[i].From != g.Edges[j].From {
			return g.Edges[i].From < g.Edges[j].From
		}
		if g.Edges[i].To != g.Edges[j].To {
			return g.Edges[i].To < g.Edges[j].To
		}
		return g.Edges[i].Type < g.Edges[j].Type
	})
}

// PreservarCarimbos transfere os carimbos de validação de um grafo ANTERIOR para o recém
// construído, para as arestas que sobreviveram ao rebuild.
//
// Sem isto, o carimbo é memória de uma execução só: o `Build` cria o grafo do zero, e o
// `anchors work` prescreve `anchors map build` ANTES de todo `check` — então cada etapa
// apagava a validação da etapa anterior. Medido: `anchors stale` reportava 9.975 de 9.975
// arestas "nunca validada" num repositório com 590 specs confrontadas dezenas de vezes.
// O comando parecia quebrado e o que faltava era o rebuild não jogar fora o que já se
// sabia.
//
// A aresta é identificada por (from, to, tipo) — o mesmo par ligado pelo mesmo motivo. Se
// qualquer ponta mudou de rev, o `StaleEdges` continua acusando: preservar o carimbo não
// é fingir que o confronto é atual, é lembrar QUANDO ele aconteceu.
func PreservarCarimbos(novo, antigo *Graph) {
	preservarSinais(novo, antigo)
	if novo == nil || antigo == nil {
		return
	}
	carimbos := make(map[string]*Stamp, len(antigo.Edges))
	// Os julgamentos de IA seguem o mesmo caminho do carimbo: reconstruir o mapa não
	// pode apagar quem já leu. Sem isto, um `map build` entre o `judge` e o `check`
	// desfazia o julgamento — e é exatamente essa a sequência que o `check --all` roda.
	julgamentos := make(map[string][]Julgamento, len(antigo.Edges))
	for i := range antigo.Edges {
		e := &antigo.Edges[i]
		k := string(e.Type) + "\x00" + e.From + "\x00" + e.To
		if e.Stamp != nil {
			carimbos[k] = e.Stamp
		}
		if len(e.Julgamentos) > 0 {
			julgamentos[k] = e.Julgamentos
		}
	}
	for i := range novo.Edges {
		e := &novo.Edges[i]
		k := string(e.Type) + "\x00" + e.From + "\x00" + e.To
		if st, ok := carimbos[k]; ok {
			e.Stamp = st
		}
		if js, ok := julgamentos[k]; ok {
			e.Julgamentos = js
		}
	}
}

// preservarSinais transfere os SINAIS DE EXECUÇÃO (teste passou? mutante sobreviveu?) do
// grafo anterior para o novo, quando o nó não mudou de rev.
//
// Mesmo defeito dos carimbos, na outra metade: o `Build` cria o grafo do zero, e o
// `anchors work` manda rodar `map build` antes de todo `check` — então cada etapa apagava
// o sinal ingerido pela anterior. Medido: depois de um `ingest --junit` casar 352 arquivos,
// um `map build` deixava o gate `testes-passam` com ⚠494 ✓0, como se nunca ninguém tivesse
// rodado a suíte.
//
// O corte é a REV: sinal de um arquivo que mudou não vale mais, e apagá-lo é o certo — é
// isso que faz `SignalStale()` acusar "reingira". Preservar só o que não mudou mantém a
// memória sem mentir sobre atualidade.
func preservarSinais(novo, antigo *Graph) {
	if novo == nil || antigo == nil {
		return
	}
	sinais := make(map[string]*TestSignal, len(antigo.Nodes))
	revs := make(map[string]string, len(antigo.Nodes))
	for i := range antigo.Nodes {
		n := &antigo.Nodes[i]
		if n.Signal != nil {
			sinais[n.ID] = n.Signal
			revs[n.ID] = n.Rev
		}
	}
	for i := range novo.Nodes {
		n := &novo.Nodes[i]
		if s, ok := sinais[n.ID]; ok && revs[n.ID] == n.Rev {
			n.Signal = s
		}
	}
}
