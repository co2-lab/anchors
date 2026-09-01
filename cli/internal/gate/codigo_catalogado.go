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

// codigo-catalogado: o que o código EXPORTA precisa estar na spec — ou dispensado nele.
//
// É o inverso do `regra-implementada`. Aquele parte da spec e pergunta "esta regra tem
// código?"; este parte do código e pergunta "este símbolo tem regra?". Sem os dois, a
// divergência escapa por um dos lados: medido num projeto real, uma spec catalogava 2
// regras para 7 funções exportadas, e nenhum gate perguntou pelas 5 restantes.
//
// O risco de um gate assim é o ruído: nem toda função exportada MERECE regra. Um helper
// de formatação, um tipo, uma constante de configuração — cobrar spec deles produziria
// centenas de achados legítimos-porém-inúteis, e um gate que acusa tudo é desligado.
//
// A saída é a mesma do resto do vocabulário (`@no-test`, `@no-code`, `@no-scenario`):
// quem escreve o código DECLARA, na linha do símbolo, que ele não carrega regra —
// com razão escrita.
//
//	// @no-rule: formatação pura, sem decisão de negócio
//	export function formatarMoeda(v: number) { … }
//
// A razão é obrigatória pelo mesmo motivo de sempre: um marcador nu vira um jeito
// silencioso de calar o gate, e some o rastro de que houve decisão.
func checkCodigoCatalogado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a spec é o catálogo — é dela que o confronto parte"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	alvo, texto, ok := alvoDaSpec(n, root, g)
	if !ok {
		return Skip, "spec sem código ligado (`specifies`) — a ausência é do gate trinca-completa"
	}

	// Os símbolos que a spec já nomeia ou que o código dispensa saem da conta.
	var orfaos []string
	for _, s := range simbolosComLinha(texto) {
		if strings.Contains(content, s.nome) {
			continue
		}
		if noRuleRE.MatchString(s.linha) {
			continue
		}
		// O nome sozinho obriga quem lê a caçar o símbolo no arquivo; com a linha,
		// o endereço está completo.
		orfaos = append(orfaos, fmt.Sprintf("%s (linha %d)", s.nome, s.num))
	}
	if len(orfaos) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(
		"%d símbolo(s) exportado(s) por `%s` que a spec não cataloga: %s.\n\nUm símbolo "+
			"fora do catálogo atravessa o pipeline sem regra a confrontar — ninguém sabe o "+
			"que ele deveria fazer.\n\nDuas saídas: catalogue-o na spec (basta o NOME "+
			"aparecer no texto de uma regra), ou declare que ele não carrega regra — na "+
			"linha do símbolo, ou no comentário logo acima dela:\n"+
			"    // @no-rule: <por que este símbolo não tem regra>\n"+
			"    export function algo() { … }",
		len(orfaos), alvo, primeiros(orfaos, 5))
}

// noRuleRE — a dispensa por SÍMBOLO, com razão obrigatória. Mesmo padrão do
// `@no-code`/`@no-scenario` (CONCEPT §5.1).
var noRuleRE = regexp.MustCompile(`@no-rule[^\S\n]*:[^\S\n]*\S+`)

// exportadoRE casa o nome de um símbolo exportado.
//
// ⚠️ É sintaxe de TypeScript/JavaScript. Num projeto Python, Go ou Ruby ele casa ZERO
// símbolos e o gate passa em silêncio — verde sobre o que não conferiu. A generalização
// (o projeto declarar o padrão, como em `mock_detect`) está pendente; até lá este gate
// só tem efeito real em projetos JS/TS.
var exportadoRE = regexp.MustCompile(
	`(?m)^\s*export\s+(?:async\s+)?(?:function|const|let|var|class|interface|type|enum)\s+([A-Za-z_$][\w$]*)`)

type simboloExportado struct {
	nome  string
	linha string // a linha da declaração + a anterior (onde o comentário de dispensa vive)
	num   int    // número da linha da declaração (1-based) — o endereço para quem vai corrigir
}

// simbolosComLinha devolve cada símbolo exportado junto do contexto onde a dispensa
// poderia estar escrita — a própria linha ou a de cima, que é onde o comentário fica.
func simbolosComLinha(codigo string) []simboloExportado {
	linhas := strings.Split(codigo, "\n")
	var out []simboloExportado
	for i, l := range linhas {
		m := exportadoRE.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		// O CONTEXTO é o símbolo mais o BLOCO DE COMENTÁRIO imediatamente acima, e não
		// só a linha anterior.
		//
		// A declaração `@no-rule` é documentação: quem a escreve a põe junto com a
		// explicação, e a explicação raramente cabe numa linha. Olhar só `i-1` fazia o
		// gate ignorar a declaração e seguir acusando — medido num arquivo onde ela
		// estava na segunda linha de um comentário de duas.
		//
		// Pior: o mesmo arquivo passava em OUTRO símbolo por acaso, porque o nome dele
		// aparecia no texto da spec. A declaração nunca era lida, e ninguém notava.
		ctx := contextoDoSimbolo(linhas, i)
		out = append(out, simboloExportado{nome: m[1], linha: ctx, num: i + 1})
	}
	return out
}

// alvoDaSpec lê o arquivo que a spec descreve (aresta `specifies`).
func alvoDaSpec(n mapx.Node, root string, g *mapx.Graph) (alvo, texto string, ok bool) {
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		return e.To, string(b), true
	}
	return "", "", false
}

// contextoDoSimbolo devolve a linha do símbolo mais o bloco de comentário acima dela.
//
// Sobe enquanto encontrar comentário (`//`, `*`, `/*`, `*/`) ou linha em branco DENTRO do
// bloco — a branco separa parágrafos de um mesmo comentário, e parar nela cortaria a
// explicação ao meio. A primeira linha de código encerra a subida.
func contextoDoSimbolo(linhas []string, i int) string {
	inicio := i
	for j := i - 1; j >= 0; j-- {
		t := strings.TrimSpace(linhas[j])
		ehComentario := strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") ||
			strings.HasPrefix(t, "/*") || strings.HasSuffix(t, "*/") || strings.HasPrefix(t, "#")
		if !ehComentario && t != "" {
			break
		}
		// Linha em branco só continua a subida se ainda houver comentário acima — senão
		// o "bloco" engoliria o arquivo inteiro até o topo.
		if t == "" {
			if j == 0 || !ehComentarioAcima(linhas, j) {
				break
			}
		}
		inicio = j
	}
	return strings.Join(linhas[inicio:i+1], "\n")
}

func ehComentarioAcima(linhas []string, j int) bool {
	for k := j - 1; k >= 0; k-- {
		t := strings.TrimSpace(linhas[k])
		if t == "" {
			continue
		}
		return strings.HasPrefix(t, "//") || strings.HasPrefix(t, "*") ||
			strings.HasPrefix(t, "/*") || strings.HasSuffix(t, "*/") || strings.HasPrefix(t, "#")
	}
	return false
}
