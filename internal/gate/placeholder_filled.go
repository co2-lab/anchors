package gate

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// placeholder-preenchido: o esqueleto que o `anchors new` emite tem de ser PREENCHIDO.
//
// O `anchors work spec` promete, textualmente: "placeholder não preenchido reprova". Não
// reprovava. Medido: uma spec recém-gerada, com `layer: TODO`, `updated_at: TODO`, o
// título `# X — TODO propósito em uma frase` e a regra `### CRUXX-B01 — TODO regra`,
// atravessava TODOS os gates bloqueantes com "✓ pode promover" — inclusive
// `header-valid ✓3`, que leu `layer: TODO` e aprovou.
//
// A razão é estrutural: os gates de header validam a FORMA (o campo existe? tem o formato
// certo?) e nunca perguntam se o valor SIGNIFICA alguma coisa. `TODO` é um valor
// bem-formado. E `spec-completa` conta seções, que o esqueleto tem todas — vazias.
//
// O custo desse silêncio é o pior tipo: um artefato que ninguém escreveu passa por
// escrito. Os gates relacionais seguintes acham a spec, acham o código, confrontam duas
// coisas que se referenciam, e todo o pipeline certifica trabalho que não existe.
//
// Medido antes de ligar, contra o repositório real: 0 achados em 590 specs. Nenhuma spec
// viva carrega placeholder do gerador — o que confirma que quem escreve, preenche, e que
// o gate cobra apenas o que ficou pelo caminho.
func checkPlaceholderFilled(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec && n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.placeholder_filled.skip_not_applicable")
	}
	achados := openPlaceholders(content, cfg)
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	return Fail, i18n.T("gate.placeholder.unfilled", len(achados), strings.Join(achados, "; "))
}

// placeholderRegexps are the three positions where a marker means "nobody filled this in",
// built from the project's vocabulary (`placeholder_markers`, default: the templates' word):
//
//	field  a header field whose VALUE is the marker (`layer: TODO`). The gravest case —
//	       `layer: TODO` is no layer at all, and the header gate approved it.
//	cell   a table cell holding only the marker (`| MTVRX-B01 | TODO |`): the rule exists
//	       as a code and says nothing.
//	title  a title or body line that opens with the marker (`# X — TODO purpose`,
//	       `TODO: what the unit does`).
//
// `<…>` is a placeholder in any vocabulary: it is the shape, not a word.
type placeholderRegexps struct{ field, cell, title *regexp.Regexp }

var (
	placeholderMu    sync.Mutex
	placeholderCache = map[string]placeholderRegexps{}
)

func placeholderREs(markers []string) placeholderRegexps {
	key := strings.Join(markers, "\x00")
	placeholderMu.Lock()
	defer placeholderMu.Unlock()
	if r, ok := placeholderCache[key]; ok {
		return r
	}
	quoted := make([]string, len(markers))
	for i, m := range markers {
		quoted[i] = regexp.QuoteMeta(m)
	}
	words := strings.Join(quoted, "|")
	r := placeholderRegexps{
		field: regexp.MustCompile(`(?mi)^\s*(?://|#|<!--|\*)?\s*([a-z_]+):\s*(` + words + `|<[^>]+>)\s*$`),
		cell:  regexp.MustCompile(`(?m)^\s*\|[^|\n]*\|[^|\n]*\b(?:` + words + `)\b[^|\n]*\|`),
		title: regexp.MustCompile(`(?m)^(?:#{1,6}\s+.*—\s*(?:` + words + `)\b.*|(?:` + words + `)[: ].*)$`),
	}
	placeholderCache[key] = r
	return r
}

// openPlaceholders acha os marcadores que o GERADOR deixou, e só eles.
//
// A distinção que evita o falso positivo: uma seção `## TODOs` (lista de pendências que o
// autor escreveu de propósito) é legítima e comum — medido, 77 specs de um projeto real a
// têm. O gate não pode confundir "o autor listou o que falta" com "o autor não escreveu
// nada". Por isso só conta o marcador em POSIÇÃO DE VALOR: campo de header, célula de
// tabela, título de regra.
func openPlaceholders(content string, cfg *config.Config) []string {
	var out []string
	visto := map[string]bool{}
	add := func(s string) {
		s = strings.TrimSpace(s)
		if len(s) > 60 {
			s = s[:60] + "…"
		}
		if s != "" && !visto[s] {
			visto[s] = true
			out = append(out, "«"+s+"»")
		}
	}
	re := placeholderREs(cfg.Placeholders())
	for _, m := range re.field.FindAllStringSubmatch(content, -1) {
		add(m[1] + ": " + m[2])
	}
	for _, m := range re.cell.FindAllString(content, -1) {
		add(m)
	}
	for _, m := range re.title.FindAllString(content, -1) {
		add(m)
	}
	return out
}
