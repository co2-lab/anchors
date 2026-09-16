package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// value-anchored: cada valor de um conjunto fechado aponta a regra que o justifica — e a
// âncora carrega o valor, para que a referência seja verificável.
//
//	export const JANELAS = [
//	  // @code-reference-[WINDOW-001]-[15m]
//	  '15m',
//	  // @code-reference-[WINDOW-002]-[1h]
//	  '1h',
//	] as const;
//
// POR QUE POR VALOR, e não pelo símbolo. O `code-cataloged` cobra o símbolo: um `@no-rule`
// sobre `export const JANELAS` libera a lista inteira de uma vez. Mas cada valor de um
// conjunto fechado é uma decisão de domínio separada — quem acrescenta `6h` está decidindo
// algo, e quem remove `24h` também.
//
// MEDIDO: o `QueryScope` de um projeto real declara `['15m','1h','6h','24h']` com um
// `@no-rule` na linha do `export`. Os quatro valores nunca foram confrontados
// individualmente, e duas telas passaram a oferecer `5m`, `30m` e `1d` — que o contrato
// não aceita. O `resolveEscopo` cai no padrão `1h` quando o valor não é aceito, então a
// tela mostrava `5m` exibindo dados de uma hora. Todos os gates verdes.
//
// O SEGUNDO COLCHETE é o que separa citar de provar. Uma âncora que só aponta a regra
// apodrece em silêncio: alguém troca `'15m'` por `'5m'` e o comentário continua parecendo
// correto. Com o valor escrito nela, o gate confronta o que a âncora AFIRMA contra o que a
// linha DIZ — e isso não depende de linguagem, porque compara duas partes do próprio
// comentário com a linha que ele anota.
//
// O QUE ESTE GATE NÃO FAZ: julgar se `15m` é um valor bom, nem se a regra `WINDOW-001`
// diz o que deveria. Ele cobra que a decisão TENHA endereço e que o endereço NÃO MINTA.
func checkValueAnchored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "o conjunto fechado vive no código; a spec é o ponto de partida do confronto"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	anchorRE := valueAnchorDe(cfg)
	if anchorRE == nil {
		return Skip, "o projeto não declarou `derived.value_anchor` — o padrão da âncora " +
			"que liga um VALOR à regra que o justifica.\n\nSem ele o gate não sabe o que " +
			"ler, e aprovar seria carimbar o que não foi conferido.\n\nDeclare em " +
			"`anchors.yaml`, com DOIS grupos de captura (a chave da regra e o valor):\n" +
			"    derived:\n" +
			"      value_anchor: \"@code-reference-\\\\[([^\\\\]]+)\\\\]-\\\\[([^\\\\]]+)\\\\]\""
	}
	exportRE := exportDetectDe(cfg)
	if exportRE == nil {
		return Skip, "o projeto não declarou `derived.export_detect` — sem ele o gate não " +
			"sabe onde um conjunto fechado começa"
	}

	alvo, texto, ok := specTarget(n, root, g)
	if !ok {
		return Skip, "spec sem código ligado (`specifies`) — a ausência é do gate trinca-completa"
	}

	semAncora, mentirosas := confrontValues(texto, exportRE, anchorRE)
	if len(semAncora) == 0 && len(mentirosas) == 0 {
		return Pass, ""
	}

	var b strings.Builder
	if len(mentirosas) > 0 {
		// A ÂNCORA QUE MENTE vem primeiro: é pior que a ausente. A ausente se vê; esta
		// parece rastreabilidade e aponta para o lugar errado.
		fmt.Fprintf(&b, "%d âncora(s) que NÃO batem com o valor ao lado, em `%s`:\n", len(mentirosas), alvo)
		for _, m := range mentirosas {
			fmt.Fprintf(&b, "    linha %d — a âncora afirma `%s`, a linha diz `%s`\n", m.linha, m.afirmado, m.real)
		}
		b.WriteString("\nUma âncora que não confere apodreceu sem ninguém notar: ela " +
			"continua parecendo rastreabilidade e aponta para o valor errado.\n\n")
	}
	if len(semAncora) > 0 {
		fmt.Fprintf(&b, "%d valor(es) de conjunto fechado sem âncora, em `%s`: %s.\n\n",
			len(semAncora), alvo, firstOnes(semAncora, 5))
		b.WriteString("Cada valor de um conjunto fechado é uma decisão — acrescentar ou " +
			"remover um muda o que o sistema aceita. Sem âncora, ninguém sabe qual regra o " +
			"justifica, e quem consome o conjunto redigita os valores em vez de exigi-los.\n\n" +
			"Ancore cada valor na linha acima dele:\n" +
			"    // @code-reference-[WINDOW-001]-[15m]\n    '15m',")
	}
	return Fail, b.String()
}

// valueAnchorDe devolve o regex declarado pelo projeto, ou nil se ele não declarou.
// Exige DOIS grupos: sem o segundo, a âncora não é verificável — e uma âncora que não
// se confere é só um comentário que envelhece.
func valueAnchorDe(cfg *config.Config) *regexp.Regexp {
	if cfg == nil || cfg.Derived == nil || strings.TrimSpace(cfg.Derived.ValueAnchor) == "" {
		return nil
	}
	re, err := regexp.Compile(cfg.Derived.ValueAnchor)
	if err != nil || re.NumSubexp() < 2 {
		return nil
	}
	return re
}

type lyingAnchor struct {
	linha    int
	afirmado string
	real     string
}

// literalRE extrai o valor literal de uma linha de conjunto — entre aspas simples,
// duplas ou crase. Deliberadamente simples: o gate confronta o que ACHA, e o que não
// casa não é acusado. Falso-negativo aqui é melhor que falso-positivo em massa.
var literalRE = regexp.MustCompile(`['"` + "`" + `]([^'"` + "`" + `]+)['"` + "`" + `]`)

// confrontValues percorre os conjuntos fechados e devolve os valores sem âncora e as
// âncoras que não batem com o valor ao lado.
func confrontValues(codigo string, exportRE, anchorRE *regexp.Regexp) ([]string, []lyingAnchor) {
	linhas := strings.Split(codigo, "\n")
	var semAncora []string
	var mentirosas []lyingAnchor

	dentro := false
	for i, l := range linhas {
		if exportRE.MatchString(l) {
			// Só interessa o conjunto: a declaração precisa ABRIR uma lista na mesma
			// linha. Um `export const X = 1` não é conjunto e não é cobrado.
			dentro = strings.Contains(l, "[")
			continue
		}
		if !dentro {
			continue
		}

		// O FECHAMENTO É CONFERIDO ANTES DO LITERAL, mas só num `]` que fecha DE FATO.
		//
		// A âncora `@code-reference-[WINDOW-001]-[15m]` carrega dois `]`, e tratá-los
		// como fim de lista encerrava o conjunto na primeira âncora — o gate saía do
		// bloco e não confrontava valor nenhum. Foi o que a régua pegou: um caso que
		// deveria reprovar passava, porque o confronto nunca acontecia.
		//
		// Uma linha que carrega âncora nunca é fechamento de lista.
		if anchorRE.MatchString(l) {
			// segue: é comentário de âncora, não fim do conjunto
		} else if strings.Contains(l, "]") {
			dentro = false
			continue
		}

		lit := literalRE.FindStringSubmatch(l)
		if lit == nil {
			continue // linha de comentário, branco, ou o que o gate não sabe ler
		}
		valor := lit[1]

		// A âncora vive na linha acima; aceita também na própria linha.
		var m []string
		if i > 0 {
			m = anchorRE.FindStringSubmatch(linhas[i-1])
		}
		if m == nil {
			m = anchorRE.FindStringSubmatch(l)
		}
		if m == nil {
			semAncora = append(semAncora, fmt.Sprintf("`%s` (linha %d)", valor, i+1))
			continue
		}
		if m[2] != valor {
			mentirosas = append(mentirosas, lyingAnchor{linha: i + 1, afirmado: m[2], real: valor})
		}
	}
	return semAncora, mentirosas
}
