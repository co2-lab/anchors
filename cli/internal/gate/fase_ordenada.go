package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a FASE do plano tem código, e a spec declara de qual depende ---
//
// O `needs:` já existia e resolvia ORDEM ENTRE PLANOS. Dentro de um plano, porém, a ordem
// vivia em prosa:
//
//	### Fase 2 — a régua mecânica (depende da Fase 1)
//
// Prosa não é confrontável. Quando o pipeline de identificação cria um card por spec, as
// cinco nascem em `to-do` indistinguíveis, e o pipeline de claim entrega qualquer uma.
//
// Medido no primeiro uso real: um agente recebeu a spec de test harness (Fase 3) com as
// Fases 1 e 2 abertas — sem `package.json` sequer, não havia onde configurar nada. Só não
// virou trabalho perdido porque a pessoa foi LER o plano; um agente que confiasse no card
// teria tentado.
//
// A fase passa a ser um item CATALOGADO, com código derivado do plano: `FNDTN-F01`. E a
// spec semeada declara `needs: FNDTN-F02` no header — a mesma palavra que o plano já usa
// para ordem, agora um nível abaixo.
//
// A identidade é o que faz a diferença: `F01` não muda quando alguém reescreve o título
// da fase, e é confrontável — "esta spec pode ser trabalhada agora?" vira uma pergunta com
// resposta, em vez de uma leitura.

// faseRE casa o cabeçalho de uma fase catalogada: `### FNDTN-F01 — a árvore e o gerenciador`.
func faseRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^#{2,4}\s+([A-Z0-9]` + config.CodeLengthPattern() + `-F\d{2})\b`)
}

// FasesDoPlano devolve os códigos de fase catalogados no plano, na ordem em que aparecem.
func FasesDoPlano(content string) []string {
	var out []string
	for _, m := range faseRE().FindAllStringSubmatch(content, -1) {
		out = append(out, m[1])
	}
	return out
}

// checkFaseOrdenada confronta a ORDEM declarada: uma fase não pode depender de outra que
// vem depois dela, nem de si mesma.
//
// Roda no PLANO, e não na spec: é o plano que declara as fases, e um ciclo entre elas é
// erro de quem escreveu o plano — a spec só aponta para uma fase que já existe.
func checkFaseOrdenada(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "só o plano declara fases"
	}
	fases := FasesDoPlano(content)
	if len(fases) == 0 {
		// Um plano sem fases catalogadas não está errado: plano pequeno não precisa de
		// fase. Pendência, e não falha — é dívida de quem quiser a ordem confrontável.
		if strings.Contains(strings.ToLower(content), "fase") {
			return Pending, "o plano fala em FASE na prosa mas não cataloga nenhuma com " +
				"código (`### " + strings.ToUpper(n.Code) + "-F01 — …`). Sem código, a ordem " +
				"existe para quem lê e não para quem confronta: as specs semeadas nascem " +
				"todas disponíveis, e o agente pega a da fase 3 com a fase 1 em aberto"
		}
		return Skip, "plano sem fases declaradas"
	}

	// A posição de cada fase, para saber o que vem antes.
	pos := map[string]int{}
	for i, f := range fases {
		if _, repetida := pos[f]; repetida {
			return Fail, fmt.Sprintf("a fase `%s` está catalogada duas vezes — "+
				"o código é a identidade dela, e duas fases com o mesmo código tornam "+
				"impossível dizer de qual uma spec depende", f)
		}
		pos[f] = i
	}

	// `depende de` na mesma seção da fase: `### FNDTN-F02 — … (depende de FNDTN-F01)`.
	dependeRE := regexp.MustCompile(`(?i)depende\s+d[eao]\s+` + "`?" + `([A-Z0-9]` +
		config.CodeLengthPattern() + `-F\d{2})`)
	var erros []string
	secoes := regexp.MustCompile(`(?m)^#{2,4}\s+`).Split(content, -1)
	for _, sec := range secoes {
		m := faseRE().FindStringSubmatch("### " + sec)
		if m == nil {
			continue
		}
		esta := m[1]
		for _, d := range dependeRE.FindAllStringSubmatch(sec, -1) {
			alvo := d[1]
			p, existe := pos[alvo]
			if !existe {
				erros = append(erros, fmt.Sprintf("`%s` depende de `%s`, que não está "+
					"catalogada neste plano", esta, alvo))
				continue
			}
			if alvo == esta {
				erros = append(erros, fmt.Sprintf("`%s` depende de si mesma", esta))
				continue
			}
			if p > pos[esta] {
				erros = append(erros, fmt.Sprintf("`%s` depende de `%s`, que vem DEPOIS "+
					"dela — a ordem declarada é impossível de cumprir", esta, alvo))
			}
		}
	}
	if len(erros) > 0 {
		return Fail, "ordem de fases inconsistente: " + strings.Join(erros, "; ")
	}
	return Pass, ""
}

// checkFaseExiste confronta o `needs:` de uma SPEC contra as fases catalogadas nos planos.
//
// É o outro lado do par: `fase-ordenada` cobra a coerência DENTRO do plano, e este cobra
// que a spec aponte para uma fase que existe. Um `needs: FNDTN-F09` num plano de quatro
// fases é uma dependência que nunca fecha — a spec ficaria bloqueada para sempre, e nada
// diria por quê.
func checkFaseExiste(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "só a spec declara de qual fase depende"
	}
	if len(n.Needs) == 0 {
		return Skip, "a spec não declara `needs:` — não há ordem a confrontar"
	}
	if g == nil {
		return Pending, "sem mapa não dá para achar os planos"
	}
	// Todas as fases catalogadas em todos os planos do mapa.
	conhecidas := map[string]string{} // código da fase → plano que a declara
	for _, p := range g.Nodes {
		if p.Kind != mapx.KindPlan {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, p.ID))
		if err != nil {
			continue
		}
		for _, f := range FasesDoPlano(string(b)) {
			conhecidas[f] = p.ID
		}
	}
	var ausentes []string
	for _, need := range n.Needs {
		if _, ok := conhecidas[need]; !ok {
			ausentes = append(ausentes, need)
		}
	}
	if len(ausentes) > 0 {
		return Fail, fmt.Sprintf("depende da(s) fase(s) %s, que nenhum plano cataloga. "+
			"Uma dependência que não existe nunca fecha: esta spec ficaria bloqueada para "+
			"sempre, e nada diria por quê. Confira o código no plano (`### CODE-F0N — …`)",
			strings.Join(ausentes, ", "))
	}
	return Pass, ""
}

// checkParentValido confronta o `parent:` contra o que existe.
//
// Um `parent` que aponta para o nada não produz erro em lugar nenhum — produz SILÊNCIO: o
// item some da árvore (ninguém o contém) ou fica pendurado numa raiz que não deveria
// existir. É a mesma classe de defeito do `needs` quebrado, e recebe o mesmo tratamento.
//
// Aceita como pai o CÓDIGO de um artefato (`FNDTN`) ou de uma fase (`FNDTN-F01`) — a
// hierarquia é livre porque a forma de organizar trabalho varia entre projetos.
func checkParentValido(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Parent == "" {
		return Skip, "sem `parent:` declarado — nada a confrontar"
	}
	if g == nil {
		return Pending, "sem mapa não dá para achar o pai"
	}
	if n.Parent == n.Code {
		return Fail, fmt.Sprintf("`%s` declara a si mesmo como pai — um item contido em si "+
			"não tem lugar na árvore, e quem a monta entra em laço", n.Code)
	}

	// Os códigos que existem: artefatos do mapa e fases catalogadas nos planos.
	existe := map[string]bool{}
	for _, o := range g.Nodes {
		if o.Code != "" {
			existe[o.Code] = true
		}
		if o.Kind != mapx.KindPlan {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, o.ID))
		if err != nil {
			continue
		}
		for _, f := range FasesDoPlano(string(b)) {
			existe[f] = true
		}
	}
	if !existe[n.Parent] {
		return Fail, fmt.Sprintf("declara `parent: %s`, que não existe — nenhum artefato "+
			"tem esse código e nenhum plano cataloga essa fase. Um pai inexistente não "+
			"falha em lugar nenhum: o item simplesmente SOME da árvore, ou fica pendurado "+
			"numa raiz que não deveria existir", n.Parent)
	}

	// CICLO: A contém B contém A. Percorre a cadeia até a raiz.
	visto := map[string]bool{n.Code: true}
	atual := n.Parent
	for i := 0; i < 64 && atual != ""; i++ {
		if visto[atual] {
			return Fail, fmt.Sprintf("a cadeia de `parent` volta a `%s` — a hierarquia tem "+
				"um ciclo, e quem a percorre nunca chega à raiz", atual)
		}
		visto[atual] = true
		atual = paiDe(g, atual)
	}
	return Pass, ""
}

// paiDe devolve o `parent` do nó com este código, ou vazio.
func paiDe(g *mapx.Graph, code string) string {
	for _, n := range g.Nodes {
		if n.Code == code {
			return n.Parent
		}
	}
	return "" // fase não é nó do mapa: a cadeia termina nela
}
