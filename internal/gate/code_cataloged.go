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
// O RUÍDO NÃO É ARGUMENTO PARA NÃO CONSTRUIR — e esta linha já esteve errada aqui.
//
// A versão anterior dizia: "cobrar spec de todo símbolo produziria centenas de achados
// legítimos-porém-inúteis, e um gate que acusa tudo é desligado". O medo estava certo;
// a conclusão, não. Compare os dois desfechos:
//
//	gate granular, desligado pelo projeto → não protege, e o projeto SABE
//	gate grosso demais, ligado            → não protege, e reporta VERDE
//
// O resultado é o mesmo; o que muda é a honestidade. O gate desligado DECLARA que não
// cobre. O gate grosso FINGE que cobre — e sendo `blocking`, carimba aprovação sobre o
// que não conferiu. É a mesma família do sinal que afirma prova inexistente.
//
// E há uma assimetria decisiva: um gate ruidoso é CALIBRÁVEL por quem o usa — desliga,
// torna informativo, escopa, dispensa caso a caso; o `anchors.yaml` já declara cada gate
// com seu próprio `blocking`, e omitir a entrada o desliga. Um gate grosso demais NÃO é
// afiável pelo projeto: a decisão foi tomada aqui dentro e ele não tem como recuperá-la.
//
// Na dúvida entre granular-com-ruído e grosso-com-silêncio, o padrão é GRANULAR. A
// cobertura é responsabilidade do projeto; o dever do Anchors é dar a informação para
// ele decidir, não decidir por ele.
//
// A dispensa continua existindo — é o que torna o ruído administrável sem mentir:
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
func checkCodeCataloged(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a spec é o catálogo — é dela que o confronto parte"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	alvo, texto, ok := specTarget(n, root, g)
	if !ok {
		return Skip, "spec sem código ligado (`specifies`) — a ausência é do gate trinca-completa"
	}

	// SEM SABER LER, O GATE CALA — nunca aprova.
	//
	// Antes o padrão TS/JS vinha embutido, e num projeto Go ou Python ele casava zero
	// símbolos e devolvia Pass: verde sobre o que não conferiu, com `blocking: true`.
	exportRE := exportDetectDe(cfg)
	if exportRE == nil {
		return Skip, "o projeto não declarou `derived.export_detect` — o padrão que " +
			"reconhece um símbolo público NESTA linguagem.\n\nSem ele o gate não sabe o " +
			"que ler, e aprovar seria carimbar o que não foi conferido.\n\nDeclare em " +
			"`anchors.yaml`, com UM grupo de captura (o nome do símbolo):\n" +
			"    derived:\n" +
			"      export_detect: \"" + exportedREPadraoTS + "\"   # TS/JS\n" +
			"      # export_detect: \"^func\\\\s+([A-Z]\\\\w*)\"          # Go\n" +
			"      # export_detect: \"^(?:def|class)\\\\s+([a-zA-Z]\\\\w*)\"  # Python"
	}

	// Os símbolos que a spec já nomeia ou que o código dispensa saem da conta.
	var orfaos []string
	for _, s := range symbolsWithLine(texto, exportRE) {
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
		len(orfaos), alvo, firstOnes(orfaos, 5))
}

// noRuleRE — a dispensa por SÍMBOLO, com razão obrigatória. Mesmo padrão do
// `@no-code`/`@no-scenario` (CONCEPT §5.1).
var noRuleRE = regexp.MustCompile(`@no-rule[^\S\n]*:[^\S\n]*\S+`)

// exportedREPadraoTS é o padrão de TypeScript/JavaScript, e só vale como SUGESTÃO ao
// projeto que ainda não declarou o seu — nunca como default silencioso.
//
// Antes ele era embutido: num projeto Go, Python ou Ruby casava ZERO símbolos e o gate
// reportava VERDE, sendo `blocking: true`. Carimbava aprovação sobre o que não tinha
// conferido, que é a pior falha possível num medidor — e com viés de ecossistema
// escondido no silêncio.
//
// Quem decide é `derived.export_detect`, pelo mesmo motivo do `mock_detect`: reconhecer
// o que é público depende da linguagem (`export const X` em TS, maiúscula inicial em Go,
// `__all__` em Python), e o Anchors não presume.
const exportedREPadraoTS = `(?m)^\s*export\s+(?:async\s+)?(?:function|const|let|var|class|interface|type|enum)\s+([A-Za-z_$][\w$]*)`

// exportDetectDe devolve o regex que ESTE projeto declarou, ou nil se não declarou.
// nil não é erro: é o gate admitindo que não sabe ler, para pular em vez de aprovar.
func exportDetectDe(cfg *config.Config) *regexp.Regexp {
	if cfg == nil || cfg.Derived == nil || strings.TrimSpace(cfg.Derived.ExportDetect) == "" {
		return nil
	}
	re, err := regexp.Compile(cfg.Derived.ExportDetect)
	if err != nil || re.NumSubexp() < 1 {
		// Regex inválido ou sem grupo de captura: também não sabemos ler. O gate pula
		// e a mensagem cobra a correção — melhor que aprovar por engano.
		return nil
	}
	return re
}

type exportedSymbol struct {
	nome  string
	linha string // a linha da declaração + a anterior (onde o comentário de dispensa vive)
	num   int    // número da linha da declaração (1-based) — o endereço para quem vai corrigir
}

// symbolsWithLine devolve cada símbolo exportado junto do contexto onde a dispensa
// poderia estar escrita — a própria linha ou a de cima, que é onde o comentário fica.
func symbolsWithLine(codigo string, re *regexp.Regexp) []exportedSymbol {
	linhas := strings.Split(codigo, "\n")
	var out []exportedSymbol
	for i, l := range linhas {
		m := re.FindStringSubmatch(l)
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
		ctx := symbolContext(linhas, i)
		out = append(out, exportedSymbol{nome: m[1], linha: ctx, num: i + 1})
	}
	return out
}

// specTarget lê o arquivo que a spec descreve (aresta `specifies`).
func specTarget(n mapx.Node, root string, g *mapx.Graph) (alvo, texto string, ok bool) {
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

// symbolContext devolve a linha do símbolo mais o bloco de comentário acima dela.
//
// Sobe enquanto encontrar comentário (`//`, `*`, `/*`, `*/`) ou linha em branco DENTRO do
// bloco — a branco separa parágrafos de um mesmo comentário, e parar nela cortaria a
// explicação ao meio. A primeira linha de código encerra a subida.
func symbolContext(linhas []string, i int) string {
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
			if j == 0 || !isCommentAbove(linhas, j) {
				break
			}
		}
		inicio = j
	}
	return strings.Join(linhas[inicio:i+1], "\n")
}

func isCommentAbove(linhas []string, j int) bool {
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
