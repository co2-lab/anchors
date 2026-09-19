package gate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// region-pair-honored: toda região aberta fecha, e fecha com o próprio código.
//
// A região (`// #region [CODEX-A03]` … `// #endregion [CODEX-A03]`) é o que dá INTERVALO à
// identidade no fonte — sem ela, sabe-se que o arquivo realiza um requisito, não onde
// (TRACEABILITY §3). O par é a única parte frágil do mecanismo, e a fragilidade é do tipo
// que não dá erro: um fecho no lugar errado produz um intervalo VÁLIDO e ERRADO.
//
// Daí o gate. São três defeitos, e nenhum deles é detectável lendo o arquivo de cima a
// baixo num arquivo grande:
//
//   - sem-fecho     — abriu e nunca fechou; o intervalo iria até o fim do arquivo
//   - fecho-orfao   — fechou sem ter aberto; sobra de um recorte/colagem
//   - fecho-trocado — fechou com OUTRO código; o aninhamento está invertido
//
// O terceiro é o que justifica o código no fecho. Com `#endregion` anônimo a contagem de
// abre/fecha bateria, o gate ficaria verde, e o carimbo de frescor mediria o intervalo do
// vizinho — o requisito A pareceria mudar quando mudou o B. Erro silencioso e cruzado é o
// pior tipo: some na medição e reaparece como conclusão errada semanas depois.
//
// FAIL e não PENDING: diferente de uma letra não declarada (que é decisão de domínio), um
// par mal fechado não tem duas leituras legítimas — é defeito de escrita, com conserto
// único e local. E a ausência de região NÃO é defeito: a delimitação é opcional, e onde
// ela não existe vale o rev do arquivo.
func checkRegionPairHonored(content string, n mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode && n.Kind != mapx.KindTest {
		return Skip, ""
	}
	regioes, erros := scan.Regioes(content)
	if len(regioes) == 0 && len(erros) == 0 {
		return Skip, i18n.T("gate.region_pair_honored.skip_no_regions")
	}
	if len(erros) == 0 {
		return Pass, i18n.T("gate.region_pair_honored.pass_well_formed", len(regioes))
	}

	// Ordena por linha: o leitor conserta de cima para baixo no arquivo, e uma lista fora
	// de ordem faz procurar duas vezes.
	sort.Slice(erros, func(i, j int) bool { return erros[i].Linha < erros[j].Linha })

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", i18n.T("gate.region_pair_honored.fail_pairing_defects", len(erros)))
	for _, e := range erros {
		switch e.Kind {
		case "sem-fecho":
			fmt.Fprintf(&b, "  %s\n", i18n.T("gate.region_pair_honored.err_unclosed", e.Linha, e.Code, e.Code))
		case "fecho-orfao":
			id := e.Achou
			if id == "" {
				id = i18n.T("gate.region_pair_honored.no_code")
			}
			fmt.Fprintf(&b, "  %s\n", i18n.T("gate.region_pair_honored.err_orphan_close", e.Linha, id))
		case "fecho-trocado":
			fmt.Fprintf(&b, "  %s\n", i18n.T("gate.region_pair_honored.err_mismatched_close", e.Linha, e.Achou, e.Code))
		}
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}
