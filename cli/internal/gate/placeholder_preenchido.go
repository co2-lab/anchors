package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// placeholder-preenchido: o esqueleto que o `anchors new` emite tem de ser PREENCHIDO.
//
// O `anchors work spec` promete, textualmente: "placeholder não preenchido reprova". Não
// reprovava. Medido: uma spec recém-gerada, com `layer: TODO`, `updated_at: TODO`, o
// título `# X — TODO propósito em uma frase` e a regra `### CRUXX-B01 — TODO regra`,
// atravessava TODOS os gates bloqueantes com "✓ pode promover" — inclusive
// `header-conforme ✓3`, que leu `layer: TODO` e aprovou.
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
func checkPlaceholderPreenchido(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec && n.Kind != mapx.KindFeature {
		return Skip, "o esqueleto com placeholder é o que o `anchors new` emite — spec e feature"
	}
	achados := placeholdersAbertos(content, cfg)
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	return Fail, fmt.Sprintf("%d placeholder(s) do gerador não preenchido(s): %s. "+
		"O esqueleto do `anchors new` não é a spec — é o formulário dela. Enquanto houver "+
		"marcador aqui, os gates de forma passam (o campo existe, o formato está certo) e "+
		"os relacionais confrontam duas peças que se referenciam sem que ninguém tenha "+
		"escrito o que a unidade faz",
		len(achados), strings.Join(achados, "; "))
}

// placeholderCampoRE: campo de header cujo VALOR é o marcador (`layer: TODO`). É o caso
// mais grave — `layer: TODO` não é camada nenhuma, e o gate de header aprovava.
var placeholderCampoRE = regexp.MustCompile(`(?mi)^\s*(?://|#|<!--|\*)?\s*([a-z_]+):\s*(TODO|FIXME|XXX|<[^>]+>)\s*$`)

// placeholderCelulaRE: célula de tabela que só tem o marcador (`| MTVRX-B01 | TODO |`) —
// a regra existe como código e não diz nada.
var placeholderCelulaRE = regexp.MustCompile(`(?m)^\s*\|[^|\n]*\|[^|\n]*\bTODO\b[^|\n]*\|`)

// placeholderTituloRE: título ou linha de corpo que abre com o marcador
// (`# X — TODO propósito`, `TODO: o que a unidade faz`).
var placeholderTituloRE = regexp.MustCompile(`(?m)^(?:#{1,6}\s+.*—\s*TODO\b.*|TODO[: ].*)$`)

// placeholdersAbertos acha os marcadores que o GERADOR deixou, e só eles.
//
// A distinção que evita o falso positivo: uma seção `## TODOs` (lista de pendências que o
// autor escreveu de propósito) é legítima e comum — medido, 77 specs de um projeto real a
// têm. O gate não pode confundir "o autor listou o que falta" com "o autor não escreveu
// nada". Por isso só conta o marcador em POSIÇÃO DE VALOR: campo de header, célula de
// tabela, título de regra.
func placeholdersAbertos(content string, cfg *config.Config) []string {
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
	for _, m := range placeholderCampoRE.FindAllStringSubmatch(content, -1) {
		add(m[1] + ": " + m[2])
	}
	for _, m := range placeholderCelulaRE.FindAllString(content, -1) {
		add(m)
	}
	for _, m := range placeholderTituloRE.FindAllString(content, -1) {
		add(m)
	}
	return out
}
