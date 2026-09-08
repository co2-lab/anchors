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
	// A identidade de um artefato DERIVADO vem da âncora irmã, não do texto dele. Montado
	// antes do laço porque a âncora pode aparecer depois na lista.
	ancoraDeDerivado := anchorCodeByDerived(files, cfg)
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
			Layer:         f.Layer,
			Parent:        f.Parent,
			Revises:       f.Revises,
			Code:          nodeCode(f, ancoraDeDerivado),
			CodeDeclarado: declaredCode(f),
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
				// CAMINHO DECLARADO É CAMINHO — não se resolve por nome.
				//
				// A busca por nome existe para a CITAÇÃO em prosa, onde o autor escreve só
				// o arquivo. Aplicá-la a um caminho inteiro faz o mapa apontar para outro
				// diretório: medido no blue-eyes, o plano 0010 semeia
				// `packages/lambdas/redis/InstanceList.spec.md`, o alvo não existia ainda,
				// e a aresta foi para `packages/lambdas/database/InstanceList.spec.md` —
				// a spec do plano 0009.
				//
				// O sintoma visível foi o `anchors next` respondendo "2 de 4" e depois
				// "3 de 3" para o mesmo plano cujas quatro sementes faltam. O dano real é
				// maior: o mapa afirma que um plano semeia a spec de outro, e todo gate
				// relacional passa a confrontar o par errado.
				if strings.ContainsRune(alvo, '/') || strings.ContainsRune(alvo, filepath.Separator) {
					continue
				}
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
		name, ext := StemOfAnchor(f.Path)
		module := filepath.Base(dir) // {{module}} — o dir-pai (ex.: run-audits em .../run-audits/handler.ts)
		// Templates efetivos = default sobrescrito pelos overrides cuja `when` casa a
		// camada da âncora (STRUCTURE.md §2.2: padrão de localização por camada quando
		// não co-localizado). Só as camadas presentes no override sobrescrevem.
		tmpls := map[string]config.Padroes{}
		for layer, tmpl := range cfg.Derived.PadroesDe() {
			tmpls[layer] = tmpl
		}
		// PRECEDÊNCIA: camada primeiro, código depois. O override por CÓDIGO é o mais
		// específico e precisa vencer — duas specs da mesma camada podem governar
		// conjuntos diferentes de arquivos (a de `tsconfig` e a de `package.json`), e sem
		// essa granularidade cada uma governaria os arquivos da outra.
		for _, ov := range cfg.Derived.Overrides {
			// A camada da UNIDADE, e não a do arquivo. Uma spec casa `**/*.spec.md` e o
			// `f.Layer` dela é `spec` — um `when: screen` nunca casaria, e o override por
			// camada seria inútil justamente para a âncora, que é quem o consulta.
			//
			// Medido: as camadas de UI do projeto de referência exigem `.tsx` no pattern,
			// o `derived.files` global diz `.ts`, e o override que reconciliava os dois
			// não era aplicado. O `triad-complete` respondia "falta o código" com o
			// arquivo no disco.
			if ov.Code != "" || ov.When != layerOfUnit(f) {
				continue
			}
			for layer, tmpl := range ov.PadroesDe() {
				tmpls[layer] = tmpl
			}
		}
		for _, ov := range cfg.Derived.Overrides {
			// nil: este laço só vê a ÂNCORA, e ela declara a identidade no header — não
			// há derivado a resolver aqui.
			if ov.Code == "" || ov.Code != nodeCode(f, nil) {
				continue
			}
			// O override por código SUBSTITUI a camada inteira, e não a completa: uma
			// spec de configuração não tem o `{{name}}.ts` da co-location, e herdá-lo
			// faria o mapa procurar um arquivo que ninguém vai escrever.
			tmpls = map[string]config.Padroes{}
			for layer, tmpl := range ov.PadroesDe() {
				tmpls[layer] = tmpl
			}
			break
		}
		// Resolve o caminho esperado de cada derivado e liga se existir.
		//
		// Uma camada pode ter VÁRIOS padrões: a spec de configuração governa vários
		// arquivos (`TypeScriptConfig` descreve seis `tsconfig.json`), e ligar só o
		// primeiro deixaria os outros órfãos no mapa.
		derived := map[string]string{} // camada → PRIMEIRO caminho existente (a trinca)
		todos := map[string][]string{} // camada → TODOS os que existem (as arestas)
		for layer, padroes := range tmpls {
			for _, tmpl := range padroes {
				for _, want := range expandePadrao(resolveTemplateM(tmpl, dir, name, ext, module), byPath) {
					if derived[layer] == "" {
						derived[layer] = want
					}
					todos[layer] = append(todos[layer], want)
				}
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
			// TODOS os arquivos que a spec governa, e não só o primeiro: uma spec de
			// configuração descreve vários, e ligar um deixaria os outros órfãos — o
			// `doctor` os reportaria como código sem spec, que é falso.
			alvos := todos["code"]
			if cfg.Derived.Anchor != "spec" {
				alvos = []string{codigo}
			}
			for _, alvo := range alvos {
				edges = append(edges, Edge{From: spec, To: alvo, Type: EdgeSpecifies, Origin: OriginConvention})
			}
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

	// SÓ O CÓDIGO DA PRÓPRIA UNIDADE liga. Citar não é declarar.
	//
	// Uma spec cita códigos de irmãs o tempo todo — "o `ARSTS-B05` decidiu que métrica
	// ausente não é métrica boa" — e o teste da irmã prova aquele código. Sem esta guarda,
	// o mapa conclui que um testa o outro.
	//
	// Medido no projeto de referência: 40 arestas `tested-by` para um único
	// `AreaStatus.test.ts`, vindas de specs que só o mencionavam em prosa. O `RateLimiting`
	// aparecia testado por ele e NÃO pelo próprio `RateLimiting.test.ts` — e o
	// `triad-complete` ficava indeterminado, que é o pior resultado: nem passa nem acusa.
	//
	// É a terceira porta do mesmo defeito que a v0.1.57 fechou nas outras duas (o dado de
	// teste lido como declaração, e o derivado sem âncora).
	codigoDe := func(f scan.File) string {
		if f.HeaderCode != "" {
			return f.HeaderCode
		}
		// Sem header, o primeiro código é a inferência que o `nodeCode` já usa.
		if len(f.Codes) > 0 {
			return RuleRoot(f.Codes[0])
		}
		return ""
	}

	type ref struct{ path, kind string }
	byCode := map[string][]ref{}
	for _, f := range files {
		dono := codigoDe(f)
		for _, c := range f.Codes {
			// A citação de uma irmã não entra no índice: só o código cuja RAIZ é a
			// identidade deste arquivo.
			if dono == "" || RuleRoot(c) != dono {
				continue
			}
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
// StemOfAnchor extrai o nome da unidade a partir do caminho da ÂNCORA.
//
// EXPORTADA porque havia uma segunda implementação. O `anchors work` cortava o nome com
// `strings.TrimSuffix(base, filepath.Ext(rel))` — e `filepath.Ext("X.spec.md")` é `.md`,
// então o `{{name}}` saía `X.spec` e o prompt mandava criar `X.spec.ts` e
// `X.spec.test.ts`.
//
// Medido no blue-eyes: `anchors work code --for packages/infra/GoLiveChecklist.spec.md`
// dizia para escrever `GoLiveChecklist.spec.feature` e `GoLiveChecklist.spec.test.ts`,
// quando o padrão do projeto — e o que o MAPA usa — é `GoLiveChecklist.ts` e
// `GoLiveChecklist.test.ts`. O prompt de trabalho e o mapa discordavam sobre onde a peça
// nasce, e quem seguisse o prompt criaria arquivo que nenhum gate encontra.
//
// A âncora é a spec (`Login.spec.md`), e dela o nome da unidade é `Login` — não
// `Login.spec`, que é o que um corte na última extensão devolveria. Os sufixos de
// artefato são compostos, e cortá-los é o que faz `{{name}}` valer para os derivados.
//
// O `{{ext}}` sai VAZIO de propósito: a extensão do código não está na spec, e inventá-la
// (`.ts`? `.tsx`? `.go`?) escolheria por um projeto que não declarou. Um template que
// precise dela declara a extensão literalmente — `{{dir}}/{{name}}.test.ts` —, e assim a
// decisão fica escrita onde se lê.
func StemOfAnchor(path string) (name, ext string) {
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
// não declara — e ela erra quando o arquivo CITA outra unidade antes de definir a sua
// (medido: uma spec de modelo que abria referenciando `DTAXX-B11` entrava no mapa como
// dona de `DTAX`, e os gates relacionais passavam a confrontar a unidade errada).
//
// O `anchors` recebe a lista de arquivos para poder derivar a identidade de um ARTEFATO
// DERIVADO da âncora irmã — ver `codeFromSibling`.
func nodeCode(f scan.File, anchors map[string]string) string {
	if f.HeaderCode != "" {
		return f.HeaderCode
	}
	// A ÂNCORA IRMÃ vence a inferência pelo texto, para artefato derivado.
	//
	// Um teste ou uma feature não são donos de identidade: eles PROVAM uma unidade, e a
	// unidade é a spec ao lado. Inferir do texto ali é ler o dado de teste como
	// declaração.
	//
	// Medido no blue-eyes: `GoLiveChecklist.test.ts` citava `ELKAD-B01` numa string —
	// o `estadoAtual` de uma dívida fictícia, "o ELKAD-B01 preparou o caminho" — e o
	// arquivo entrou no mapa com `code: ELKAD`. O `scenario-coverage` passou a cobrar
	// vinte e um cenários de outras specs, e os três invariantes que o teste PROVAVA
	// apareciam como não provados.
	//
	// O arquivo vizinho não tem esse problema (`AreaStatus.test.ts` funciona sem header)
	// porque todas as menções dele são da própria unidade. É o caso fácil, e ele esconde
	// o defeito até alguém citar outra spec.
	if c, ok := anchors[f.Path]; ok && c != "" {
		return c
	}
	return primaryCode(f.Codes)
}

// anchorCodeByDerived mapeia cada artefato DERIVADO para o código da âncora dele.
//
// A âncora é a spec (`derived.anchor`), e o vínculo é o STEM: `X.spec.md` é a âncora de
// `X.ts`, `X.feature` e `X.test.ts` no mesmo diretório. Só o que a âncora DECLARA no
// header conta — se ela mesma não tem identidade declarada, não há o que propagar.
func anchorCodeByDerived(files []scan.File, cfg *config.Config) map[string]string {
	if cfg == nil || cfg.Derived == nil {
		return nil
	}
	// stem da âncora → código declarado
	porStem := map[string]string{}
	for _, f := range files {
		if f.Kind != cfg.Derived.Anchor || f.HeaderCode == "" {
			continue
		}
		name, _ := StemOfAnchor(f.Path)
		porStem[filepath.ToSlash(filepath.Dir(f.Path))+"/"+name] = f.HeaderCode
	}
	if len(porStem) == 0 {
		return nil
	}
	out := map[string]string{}
	for _, f := range files {
		if f.Kind == cfg.Derived.Anchor || f.HeaderCode != "" {
			continue // a âncora tem a sua; quem declarou não precisa de fallback
		}
		chave := filepath.ToSlash(filepath.Dir(f.Path)) + "/" + stemOfDerived(f.Path)
		if c, ok := porStem[chave]; ok {
			out[f.Path] = c
		}
	}
	return out
}

// declaredCode diz se a identidade foi DECLARADA (header `code:`) ou apenas inferida do
// texto. Ver Node.CodeDeclarado para o porquê da distinção.
func declaredCode(f scan.File) bool { return f.HeaderCode != "" }

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

// PreserveStamps transfere os carimbos de validação de um grafo ANTERIOR para o recém
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
func PreserveStamps(novo, antigo *Graph) {
	preservarSinais(novo, antigo)
	if novo == nil || antigo == nil {
		return
	}
	carimbos := make(map[string]*Stamp, len(antigo.Edges))
	// Os julgamentos de IA seguem o mesmo caminho do carimbo: reconstruir o mapa não
	// pode apagar quem já leu. Sem isto, um `map build` entre o `judge` e o `check`
	// desfazia o julgamento — e é exatamente essa a sequência que o `check --all` roda.
	julgamentos := make(map[string][]Judgment, len(antigo.Edges))
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

// expandePadrao devolve os caminhos do mapa que casam o padrão.
//
// Sem curinga, é uma busca exata — o caminho existe ou não. Com `*`, casa vários: a spec
// de configuração descreve `packages/*/tsconfig.json`, e enumerar os pacotes na config
// obrigaria a editá-la a cada pacote novo, que é o tipo de manutenção que se esquece.
func expandePadrao(padrao string, byPath map[string]scan.File) []string {
	if !strings.ContainsAny(padrao, "*?[") {
		if _, ok := byPath[padrao]; ok {
			return []string{padrao}
		}
		return nil
	}
	var out []string
	for p := range byPath {
		if ok, err := filepath.Match(padrao, p); err == nil && ok {
			out = append(out, p)
		}
	}
	sort.Strings(out) // ordem estável: o mapa não pode mudar entre execuções
	return out
}

// stemOfDerived corta os sufixos de um artefato DERIVADO até o nome da unidade.
//
// O `StemOfAnchor` corta os sufixos da ÂNCORA (`.spec.md`, `.feature`); aqui o alvo é o
// outro lado — `X.test.ts`, `X.ts`, `X.tsx`, `X.test.tsx`. Cortar só a última extensão
// deixaria `X.test`, que não casa o stem da spec.
//
// Não vem do `derived.files` de propósito: os templates são caminhos com `{{name}}`, e
// invertê-los para extrair o nome exigiria casá-los como regex — que falha no primeiro
// projeto que use um template com mais de uma variável. Cortar sufixo é o que os dois
// lados já fazem.
func stemOfDerived(path string) string {
	base := filepath.Base(path)
	// Do mais específico para o menos: `.test.ts` antes de `.ts`, senão `X.test.ts`
	// viraria `X.test`.
	for _, suf := range []string{
		".spec.md", ".feature",
		".test.ts", ".test.tsx", ".test.js", ".test.jsx",
		".spec.ts", ".spec.tsx",
		".test.go", "_test.go",
	} {
		if b, ok := strings.CutSuffix(base, suf); ok {
			return b
		}
	}
	if i := strings.LastIndex(base, "."); i > 0 {
		return base[:i]
	}
	return base
}

// RuleRoot devolve a UNIDADE dona de um código de regra ou cenário.
//
// `GLCGL-B01` → `GLCGL`; `GLCGL-B01#02` → `GLCGL`. É o prefixo antes do primeiro hífen, e
// não uma regex do formato completo: o vocabulário de letras é extensível por projeto
// (`rule_types`), e uma regex aqui teria de acompanhar cada extensão — divergindo em
// silêncio na primeira que alguém declarasse.
func RuleRoot(codigo string) string {
	if i := strings.Index(codigo, "-"); i > 0 {
		return codigo[:i]
	}
	return codigo
}

// layerOfUnit devolve a camada da UNIDADE de um arquivo.
//
// O header vence porque é onde o autor a declara; sem ele, o `Layer` do arquivo já é a
// camada da unidade para tudo o que não é spec.
func layerOfUnit(f scan.File) string {
	if f.HeaderLayer != "" {
		return f.HeaderLayer
	}
	return f.Layer
}
