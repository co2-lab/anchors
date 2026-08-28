// Package health é o validador de saúde do ecossistema (QUALITY §5.2) — a visão
// GLOBAL. Diferente dos gates (que confrontam nó contra critério, incremental), o
// health varre o mapa + a config + o disco e caça as pontas SISTÊMICAS: integridade
// do mapa, órfãos, camadas frouxas, cobertura de gates. Ele APRESENTA e REGISTRA,
// mas NÃO bloqueia (roda sob demanda, fora do caminho de merge; "detecta e
// apresenta, não arbitra", CONCEPT §2).
package health

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// Severity de um achado.
type Severity string

const (
	Warn Severity = "warn" // ponta aberta / débito — precisa de atenção
	Info Severity = "info" // observação — não é problema
)

// Finding é um achado sistêmico do doctor.
type Finding struct {
	Check    string // qual verificação (ex.: "aresta-morta", "orfao-identidade")
	Severity Severity
	Subject  string // o nó/aresta/camada afetado
	Detail   string
}

// Report é o resultado da varredura.
type Report struct {
	Findings []Finding
	Nodes    int
	Edges    int
	Layers   int
}

// Warnings devolve só os achados de atenção.
func (r Report) Warnings() []Finding {
	var w []Finding
	for _, f := range r.Findings {
		if f.Severity == Warn {
			w = append(w, f)
		}
	}
	return w
}

// Diagnose roda todas as verificações sistêmicas. root é a raiz (para checar o
// disco). Puro em relação a (graph, cfg); toca o disco só para read/stat.
func Diagnose(g *mapx.Graph, cfg *config.Config, root string) Report {
	r := Report{Nodes: len(g.Nodes), Edges: len(g.Edges), Layers: len(cfg.Layers)}
	r.Findings = append(r.Findings, checkMapFidelity(g, root)...)
	r.Findings = append(r.Findings, checkOrphans(g)...)
	r.Findings = append(r.Findings, checkDuplicateCodes(g)...)
	r.Findings = append(r.Findings, checkLooseLayers(g, cfg)...)
	r.Findings = append(r.Findings, checkGateCoverage(g, cfg)...)
	r.Findings = append(r.Findings, checkSkipOnValido(cfg)...)
	r.Findings = append(r.Findings, checkFerramentasAusentes(cfg)...)
	r.Findings = append(r.Findings, checkGitAusente(cfg, root)...)
	r.Findings = append(r.Findings, checkAmbienteGitHub(cfg, root)...)
	r.Findings = append(r.Findings, checkNeedsDosPlanos(g)...)
	r.Findings = append(r.Findings, checkSinaisAusentes(g)...)
	sortFindings(r.Findings)
	return r
}

// --- identidade DUPLICADA (colisão de código de cenário, TRACEABILITY §3) ---
//
// O código de cenário deve ser ÚNICO por unidade. Quando duas unidades em pastas
// diferentes carregam o MESMO código, é uma colisão: a rastreabilidade cruza fios e
// a propagação gera arestas espúrias (uma spec "propaga" para o teste de outra
// unidade só porque compartilham o código). Este check pega a colisão na raiz — é a
// contrapartida de "identidade-ausente".
func checkDuplicateCodes(g *mapx.Graph) []Finding {
	// Só kinds que são DONOS de identidade contam. doc/guide/plan apenas REFERENCIAM
	// códigos (documentam, planejam) — não são donos, então não colidem.
	owns := func(k mapx.Kind) bool {
		switch k {
		case mapx.KindSpec, mapx.KindFeature, mapx.KindTest, mapx.KindCode:
			return true
		}
		return false
	}
	// O dono de um código é a UNIDADE, não o diretório. Unidade = o stem co-locado
	// (dir/nome sem os sufixos .spec.md/.feature/.test.*). Assim a mesma unidade
	// espalhada em lib/+screens/ (irmãs, stems iguais) não conta como colisão; só
	// stems REALMENTE distintos compartilhando um código é que colidem.
	unitsByCode := map[string]map[string]bool{}
	for _, n := range g.Nodes {
		if n.Code == "" || !owns(n.Kind) {
			continue
		}
		// opt-out honesto: o arquivo declarou @anchors-shared-code — seus códigos
		// pertencem a outra unidade de propósito, não são posse. Não conta.
		if n.SharedCode {
			continue
		}
		if unitsByCode[n.Code] == nil {
			unitsByCode[n.Code] = map[string]bool{}
		}
		unitsByCode[n.Code][unitStem(n.ID)] = true
	}
	var out []Finding
	for code, units := range unitsByCode {
		if len(units) < 2 {
			continue // um só dono (uma unidade, mesmo em pastas irmãs) — ok
		}
		owners := make([]string, 0, len(units))
		for u := range units {
			owners = append(owners, u)
		}
		sort.Strings(owners)
		// Só a colisão CROSS-DOMAIN é reportada — dois donos em features/domínios
		// DIFERENTES cruzam a rastreabilidade e geram propagação espúria (Warn).
		// Unidades da MESMA feature (a tela + sua lib de domínio) compartilhando
		// código é DELIBERADO e legítimo: não é achado — um "alerta" que nunca exige
		// ação é ruído, então não o emitimos.
		if !crossDomain(owners) {
			continue
		}
		out = append(out, Finding{
			"identidade-duplicada", Warn, code,
			"código usado por unidades distintas: " + strings.Join(owners, ", "),
		})
	}
	return out
}

// crossDomain: os donos pertencem a features/domínios DIFERENTES? Usa o segmento de
// feature (features/<x>) como domínio; se todos compartilham o mesmo, é intra-feature.
func crossDomain(owners []string) bool {
	domain := func(p string) string {
		parts := strings.Split(p, "/")
		for i, s := range parts {
			if s == "features" && i+1 < len(parts) {
				return "features/" + parts[i+1]
			}
		}
		return parts[0] // fora de features/: o topo (apps, packages, …) já separa
	}
	first := domain(owners[0])
	for _, o := range owners[1:] {
		if domain(o) != first {
			return true
		}
	}
	return false
}

// unitStem devolve a identidade da UNIDADE de um arquivo: o caminho sem os sufixos de
// artefato (.spec.md, .feature, .test.*) e sem a extensão de código. A trinca
// Login.spec.md / Login.feature / Login.test.tsx / Login.tsx colapsa em ".../Login".
func unitStem(id string) string {
	// Junta sempre com barra normal: o retorno é identidade de unidade, comparada contra
	// ids do mapa e fatiada por "/" adiante (crossDomain). No Windows o Join devolveria
	// "features\b\lib\x", o domínio deixaria de ser reconhecido, e um
	// compartilhamento legítimo intra-feature apareceria como colisão cross-domain.
	junta := func(dir, base string) string { return filepath.ToSlash(filepath.Join(dir, base)) }

	base := filepath.Base(id)
	for _, suf := range []string{".spec.md", ".feature"} {
		if b, ok := strings.CutSuffix(base, suf); ok {
			return junta(filepath.Dir(id), b)
		}
	}
	// .test.<ext> → tira .test e a extensão
	if i := strings.Index(base, ".test."); i >= 0 {
		base = base[:i]
		return junta(filepath.Dir(id), base)
	}
	// código comum: tira só a extensão
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return junta(filepath.Dir(id), base)
}

// --- integridade do mapa: arestas mortas / nós fantasma (TRACEABILITY §6) ---

func checkMapFidelity(g *mapx.Graph, root string) []Finding {
	var out []Finding
	exists := func(id string) bool {
		_, err := os.Stat(filepath.Join(root, id))
		return err == nil
	}
	// nós cujo arquivo sumiu
	nodeIDs := map[string]bool{}
	for _, n := range g.Nodes {
		nodeIDs[n.ID] = true
		if !exists(n.ID) {
			out = append(out, Finding{"no-fantasma", Warn, n.ID, "nó no mapa sem arquivo em disco"})
		}
	}
	// arestas apontando para nó inexistente (aresta morta)
	for _, e := range g.Edges {
		if !nodeIDs[e.From] || !nodeIDs[e.To] {
			out = append(out, Finding{"aresta-morta", Warn, e.From + " → " + e.To, "aresta aponta para nó ausente do mapa"})
		}
	}
	return out
}

// --- órfãos (TRACEABILITY §6) ---

func checkOrphans(g *mapx.Graph) []Finding {
	var out []Finding

	// índice: um nó é `from` de alguma aresta? (usado para detectar spec sem realização)
	hasOutgoing := map[string]bool{}
	for _, e := range g.Edges {
		hasOutgoing[e.From] = true
	}

	for _, n := range g.Nodes {
		switch n.Kind {
		// NB: código SEM spec NÃO é órfão — é o caso normal. Utils, constantes, hooks,
		// lib pura não precisam de spec; specs são para as unidades que as merecem
		// (telas, componentes). Marcar todo código sem spec como débito seria ruído.
		// A direção que importa é a INVERSA (abaixo): uma SPEC sem código é que é um
		// requisito sem encarnação.
		case mapx.KindSpec:
			// spec sem código realizado — requisito sem encarnação
			if !hasOutgoing[n.ID] {
				out = append(out, Finding{"spec-sem-realizacao", Warn, n.ID, "spec não liga a nenhum código/feature/teste"})
			}
			// identidade ausente (o órfão invisível): spec sem código de cenário
			if n.Code == "" {
				out = append(out, Finding{"identidade-ausente", Warn, n.ID, "spec sem código de cenário (órfão invisível)"})
			}
		}
	}
	return out
}

// --- pilares frouxos: camada declarada sem arquivos; guide que não rege nada ---

func checkLooseLayers(g *mapx.Graph, cfg *config.Config) []Finding {
	var out []Finding
	// contagem de nós por kind e por "é regido"
	kindCount := map[mapx.Kind]int{}
	for _, n := range g.Nodes {
		kindCount[n.Kind]++
	}
	// camada declarada na config mas sem nenhum nó do seu kind
	for name, l := range cfg.Layers {
		if kindCount[mapx.Kind(l.Kind)] == 0 {
			out = append(out, Finding{"camada-vazia", Warn, name, "camada declarada na Estrutura mas sem nenhum arquivo"})
		}
	}
	// guide que não aparece como `from` de nenhuma aresta governs (não rege nada)
	governsFrom := map[string]bool{}
	for _, e := range g.Edges {
		if e.Type == mapx.EdgeGoverns {
			governsFrom[e.From] = true
		}
	}
	for _, n := range g.Nodes {
		if n.Kind == mapx.KindGuide && !governsFrom[n.ID] {
			// um guide que não rege ninguém NÃO é guide — é doc parado. É débito de
			// Estrutura (Warn): ou falta a regra `governs` (o guide rege algo que não
			// foi declarado), ou o arquivo deveria ser kind `doc`.
			out = append(out, Finding{"guide-sem-governo", Warn, n.ID,
				"guide não rege ninguém — declare uma regra `governs` (a tag que ele rege) ou reclassifique como doc"})
		}
	}
	return out
}

// --- cobertura de gates: kinds sem nenhum gate declarado (meta-gate, §5.1) ---

func checkGateCoverage(g *mapx.Graph, cfg *config.Config) []Finding {
	var out []Finding
	// kinds presentes no projeto
	present := map[string]bool{}
	for _, n := range g.Nodes {
		present[string(n.Kind)] = true
	}
	// kinds cobertos por algum gate
	covered := map[string]bool{}
	for _, gt := range cfg.Gates {
		for _, k := range gt.On {
			covered[k] = true
		}
	}
	var kinds []string
	for k := range present {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		// código e doc podem legitimamente não ter gate próprio; specs/features/testes sim
		if !covered[k] && (k == "spec" || k == "feature" || k == "test") {
			out = append(out, Finding{"kind-sem-gate", Warn, k, "nenhum gate declarado confronta este kind (buraco de cobertura)"})
		}
	}
	return out
}

// checkSinaisAusentes acusa os sinais de verificação que o projeto NÃO ingere.
//
// A distinção que este check materializa: o Anchors pode EXIGIR que a spec decida e que
// a peça exista — isso é dele. Não pode exigir que o projeto tenha uma ferramenta de
// mutação ou um reporter JUnit: depende de stack, de linguagem, de tempo de CI. Exigir
// seria o framework decidindo pelo projeto.
//
// Mas não exigir não é o mesmo que ficar calado. Sem esses sinais, uma família inteira
// de gates fica Pending — passa por verde no relatório e certifica trabalho que ninguém
// mediu. O doctor é o lugar certo para o alerta: apresenta o risco, nomeia o que se
// perde, e deixa a decisão com o projeto.
//
// A escala de severidade segue o que se perde:
//   - execução (JUnit): Warn — sem ele nem se sabe se o teste passou.
//   - cobertura (lcov): Warn — sem ela não se sabe se a linha rodou.
//   - mutação: Info — o mais caro de adotar; a ausência é aceitável, a ignorância não.
func checkSinaisAusentes(g *mapx.Graph) []Finding {
	var códigos, testes int
	var comExec, comCov, comMut int
	for _, n := range g.Nodes {
		switch n.Kind {
		case mapx.KindCode:
			códigos++
			if n.Signal != nil && n.Signal.TotalLines > 0 {
				comCov++
			}
			if n.Signal != nil && (n.Signal.MutantsKilled > 0 || n.Signal.MutantsSurvived > 0) {
				comMut++
			}
		case mapx.KindTest:
			testes++
			if n.Signal != nil && (n.Signal.Passed > 0 || n.Signal.Failed > 0 || n.Signal.Skipped > 0) {
				comExec++
			}
		}
	}

	var out []Finding
	if testes > 0 && comExec == 0 {
		out = append(out, Finding{"sinal-ausente", Warn, "execução (JUnit)",
			fmt.Sprintf("%d teste(s) no mapa e NENHUM resultado ingerido — os gates de "+
				"execução ficam Pending, então o pipeline não sabe se algum teste passa. "+
				"Rode `anchors ingest --junit <arquivo>` (todo runner sério emite JUnit XML)", testes)})
	}
	if códigos > 0 && comCov == 0 {
		out = append(out, Finding{"sinal-ausente", Warn, "cobertura (lcov)",
			fmt.Sprintf("%d arquivo(s) de código e NENHUMA cobertura ingerida — sem ela não "+
				"se distingue código exercitado de código nunca executado. "+
				"Rode `anchors ingest --lcov <arquivo>` (lcov é formato universal)", códigos)})
	}
	if códigos > 0 && comMut == 0 {
		out = append(out, Finding{"sinal-ausente", Info, "mutação",
			"nenhum sinal de mutação ingerido. É o único sinal que responde se o teste " +
				"PROVA a linha — cobertura só diz que ela EXECUTOU. O risco concreto: uma " +
				"suíte 100% verde e 100% coberta pode não derrubar teste algum quando você " +
				"apaga uma guarda do código. Se o stack tiver ferramenta (Stryker/PIT/" +
				"Infection/mutmut), rode e `anchors ingest --mutation <relatório>`"})
	} else if comMut > 0 && comMut < códigos {
		out = append(out, Finding{"sinal-ausente", Info, "mutação (parcial)",
			fmt.Sprintf("%d de %d arquivo(s) com sinal de mutação — o resto do código não "+
				"tem essa prova", comMut, códigos)})
	}
	return out
}

func sortFindings(fs []Finding) {
	sort.Slice(fs, func(i, j int) bool {
		if fs[i].Check != fs[j].Check {
			return fs[i].Check < fs[j].Check
		}
		return fs[i].Subject < fs[j].Subject
	})
}

// checkSkipOnValido acusa perspectiva desconhecida em `skip_on`.
//
// Vale porque o campo é uma lista de EXCLUSÃO por nome: um typo (`chnage`) não casa
// perspectiva nenhuma, o gate segue rodando nas duas, e o autor acredita tê-lo
// desligado numa delas. O erro é silencioso justamente do lado perigoso — quem escreveu
// `skip_on` queria menos execução e recebeu mais, sem nenhum sinal.
func checkSkipOnValido(cfg *config.Config) []Finding {
	if cfg == nil {
		return nil
	}
	validas := map[string]bool{
		config.PerspectiveChange: true,
		config.PerspectiveAll:    true,
	}
	var out []Finding
	for _, gt := range cfg.Gates {
		for _, p := range gt.SkipOn {
			if !validas[p] {
				out = append(out, Finding{"skip-on-invalido", Warn, gt.Name,
					"perspectiva desconhecida em `skip_on`: " + p +
						" (use `change` ou `all`) — como está, o gate não é pulado em lugar nenhum"})
			}
		}
	}
	return out
}

// --- ferramenta AUSENTE (gate declarado que não tem como rodar) ---
//
// Um gate com `needs_tool` cuja ferramenta não está no PATH é PULADO pelo motor
// (Skip, nunca Fail — ver gate.ferramentaAusente). Isso protege o veredito, mas abre
// um risco novo: gate que não roda é indistinguível de gate que passou, para quem só
// olha o verde. O projeto ficaria descoberto exatamente na medida em que ninguém
// notasse.
//
// É aqui que a ausência tem de aparecer, e não a cada varredura: o doctor é lido uma
// vez e resolvido, enquanto um aviso no `check` viraria ruído a cada commit — e ruído
// recorrente treina a equipe a ignorar, que é o oposto do que este aviso quer.
//
// Warn, não Info: a cobertura que o projeto DECLARA ter não é a que ele tem.
func checkFerramentasAusentes(cfg *config.Config) []Finding {
	if cfg == nil {
		return nil
	}
	var out []Finding
	for _, gt := range cfg.Gates {
		if gt.NeedsTool == "" {
			continue
		}
		if _, err := exec.LookPath(gt.NeedsTool); err == nil {
			continue
		}
		msg := "gate DESABILITADO: `" + gt.NeedsTool + "` não está no PATH"
		if gt.InstallHint != "" {
			msg += " — instale com: " + gt.InstallHint
		}
		out = append(out, Finding{"ferramenta-ausente", Warn, gt.Name, msg})
	}
	return out
}

// --- git AUSENTE (o substrato que não avisa quando falta) ---
//
// Git não é conforto do Anchors: é de onde vêm o carimbo de alteração das âncoras
// (gitmeta), a cobertura de diff, o pre-commit e — no modo `github` — o repositório da
// fila de trabalho. Nada disso falha ruidosamente quando o git falta; simplesmente não
// acontece — e o que o usuário perde não é sinalizado por nenhum erro.
//
// É o mesmo motivo de `checkFerramentasAusentes`: o doctor existe para avisar COM
// ANTECEDÊNCIA, antes de o usuário descobrir o buraco no meio de um trabalho e
// investigar a camada errada.
//
// Os dois casos são achados DIFERENTES porque o conserto é diferente: sem o binário,
// instalar; sem o repositório, iniciar. Uma mensagem que colapse os dois manda o
// usuário procurar o problema onde ele não está.
func checkGitAusente(cfg *config.Config, root string) []Finding {
	return gitAusente(cfg, root, gitNoPath())
}

// gitNoPath é uma variável, e não uma chamada direta, para que o teste alcance o caso
// "git não instalado" numa máquina que TEM git. Sem isso a distinção que separa
// "instalar" de "iniciar" ficaria sem cobertura justamente na metade mais rara — a que
// ninguém reproduz por acidente.
var gitNoPath = func() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func gitAusente(cfg *config.Config, root string, instalado bool) []Finding {
	if !instalado {
		return []Finding{{"git-ausente", Warn, "git",
			"o git não está no PATH — o carimbo de alteração (updated_at), " +
				"`coverage --diff` e o pre-commit ficam desligados, em silêncio"}}
	}
	if temRepo(root) {
		return nil
	}
	det := "este projeto não está sob git — o carimbo de alteração (updated_at), " +
		"`coverage --diff` e `install-hooks` não têm como funcionar; rode `git init`"
	// No modo `github` isso deixa de ser débito e vira impedimento: a fila de trabalho
	// mora nas issues de um repositório, e sem repo não há de onde puxar.
	if cfg.ModoGitHub() {
		det = "o `workflow.mode: github` exige repositório, e este projeto não está sob " +
			"git — a fila de trabalho não tem de onde ser puxada; rode `git init`"
	}
	return []Finding{{"git-ausente", Warn, "repositório", det}}
}

// temRepo sobe a árvore procurando `.git` (arquivo ou diretório: worktree e submódulo
// usam um ponteiro `gitdir:`). Não roda `git` — a resposta é sobre o DISCO, e serve
// igual quando o binário não existe.
func temRepo(root string) bool {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return true
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return false
		}
		dir = pai
	}
}

// --- `needs` quebrado e CICLO entre planos (BOOTSTRAP.md §7) ---
//
// O `needs:` de um plano declara ORDEM DE TRABALHO: "não comece antes daquele terminar".
// É o que permite escrever a sequência inteira de planos de uma vez e liberar aos poucos
// — o arquivo existe, mas o card só nasce quando o pré-requisito fecha.
//
// Duas formas de quebrar, e as duas são silenciosas sem este check:
//
//   - apontar para um plano que NÃO EXISTE: a dependência simplesmente não é respeitada,
//     e o trabalho começa fora de ordem sem nada acusar;
//   - CICLO (A precisa de B, B precisa de A): nenhum dos dois pode começar, nunca, e o
//     board fica com dois cards que ninguém pega sem que se saiba por quê.
func checkNeedsDosPlanos(g *mapx.Graph) []Finding {
	planos := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Kind == mapx.KindPlan {
			planos[n.ID] = true
		}
	}
	if len(planos) == 0 {
		return nil
	}

	// As arestas `needs` que o mapa conseguiu ligar.
	depende := map[string][]string{}
	for _, e := range g.Edges {
		if e.Type == mapx.EdgeNeeds {
			depende[e.From] = append(depende[e.From], e.To)
		}
	}

	var out []Finding

	// `needs` para plano inexistente: o mapa DESCARTA essas arestas (não liga a nada),
	// então a ausência é invisível no grafo — e o efeito é a dependência simplesmente
	// não valer, com o trabalho começando fora de ordem sem nada acusar.
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindPlan {
			continue
		}
		for _, alvo := range n.Needs {
			if !planos[alvo] {
				out = append(out, Finding{"needs-quebrado", Warn, n.ID,
					"`needs: " + alvo + "` aponta para um plano que não existe — a ordem de " +
						"trabalho declarada não vale, e o card nasce antes da base"})
			}
		}
	}

	// Ciclo: uma busca em profundidade por nó, marcando o caminho em curso.
	estado := map[string]int{} // 0=novo, 1=no caminho, 2=fechado
	var caminho []string
	var visita func(string) bool
	visita = func(p string) bool {
		estado[p] = 1
		caminho = append(caminho, p)
		for _, alvo := range depende[p] {
			if estado[alvo] == 1 {
				out = append(out, Finding{"needs-ciclo", Warn, p,
					"ciclo de `needs` entre planos (" + strings.Join(caminho, " → ") + " → " + alvo +
						") — nenhum deles pode começar, e o board fica com cards que ninguém pega"})
				return true
			}
			if estado[alvo] == 0 && visita(alvo) {
				return true
			}
		}
		caminho = caminho[:len(caminho)-1]
		estado[p] = 2
		return false
	}
	for p := range planos {
		if estado[p] == 0 {
			caminho = nil
			if visita(p) {
				break // um ciclo relatado basta: o segundo seria o mesmo problema
			}
		}
	}
	return out
}
