package gate

import (
	"os"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// vr-baseline: o cenário de regressão VISUAL prometido tem imagem de referência.
//
// O VR é a única superfície de prova que nenhum gate alcançava. Um cenário marcado
// `@vr-level` declara que aquela tela é provada por CAPTURA — não por asserção em teste
// unitário, que é onde os demais gates olham. Sem baseline, o cenário existe, o gate de
// feature o conta como coberto, e não há imagem contra a qual comparar nada: a prova
// prometida não acontece, e nada acusa.
//
// Medido no repositório que originou o gate: 105 baselines no disco, 4 features declarando
// cenário VR, e 3 cenários declarados SEM imagem — dois deles de uma feature recém-escrita,
// onde o autor copiou a convenção do vizinho sem gerar a captura.
//
// # O que este gate NÃO faz, e por quê
//
// A pergunta seguinte — "o baseline está DESATUALIZADO?" — parece o complemento óbvio e
// não vira gate. Medido pelas duas fontes possíveis:
//
//   - `mtime` do disco: 105 de 105 acusados. Não é sinal — um `git clone` reescreve o
//     timestamp de tudo, e o número diria o mesmo num repositório impecável.
//   - data de COMMIT do git: 103 de 105 acusados. É sinal verdadeiro (o projeto não
//     regrava baseline desde junho), e ainda assim inútil como gate: acusar 98% do
//     repositório é ser desligado no primeiro dia, levando junto os gates que funcionam.
//
// A dívida é real e fica NOMEADA no veredito Pendente, sem virar acusação. Quando o
// projeto passar a regravar baseline junto com a tela, o número cai e o gate endurece.
func checkVRBaseline(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.vr_baseline.skip_not_feature")
	}
	cenarios := vrScenarios(content, cfg)
	if len(cenarios) == 0 {
		return Skip, i18n.T("gate.vr_baseline.skip_no_vr_scenarios")
	}

	base := strings.TrimSuffix(n.ID, ".feature")
	var semImagem []string
	for _, c := range cenarios {
		// A convenção do baseline é `<Unidade>.<CODEX-VR-variante>.png`, ao lado da
		// unidade. O glob cobre a variante opcional: `TCDTX-VR` casa
		// `TCDTX-VR-loaded.png` tanto quanto `TCDTX-VR.png`.
		achou, _ := doublestar.Glob(os.DirFS(root), base+"."+c+"*.png")
		if len(achou) == 0 {
			semImagem = append(semImagem, c)
		}
	}
	if len(semImagem) == 0 {
		return Pass, ""
	}
	sort.Strings(semImagem)
	return Fail, i18n.T("gate.vr.missing", len(semImagem), strings.Join(semImagem, ", "))
}

// vrScenarios devolve os códigos de cenário marcados como regressão visual na feature.
//
// O que é "regressão visual" vem do PROJETO (`derived.regimes` no anchors.yaml diz qual
// tag nomeia esse regime), com `vr-level` como default.
//
// O default é em inglês pelo mesmo motivo que o resto do produto: ele INDUZ. Um projeto
// que não declara `regimes` herda esta grafia e a escreve em cada cenário; um default em
// português faria todo projeto novo nascer com vocabulário misto — e a tag é
// identificador, então consertar depois custa migração, não tradução. Medido: um projeto
// real declarou `nivel-unit` porque foi o exemplo que o Anchors lhe deu, e hoje são 800
// cenários. O Anchors não impõe a nomenclatura — mas o que ele sugere vira a do projeto.
func vrScenarios(content string, cfg *config.Config) []string {
	tag := visualRegimeTag(cfg)
	var out []string
	visto := map[string]bool{}
	for _, linha := range strings.Split(content, "\n") {
		if !strings.Contains(linha, "@"+tag) {
			continue
		}
		for _, m := range featScenarioCodeRE.FindAllStringSubmatch(linha, -1) {
			// só o código que É de VR — a tag de regime marca a linha, mas o cenário
			// pode co-etiquetar outros códigos que não são visuais.
			if strings.Contains(m[1], "-VR") && !visto[m[1]] {
				visto[m[1]] = true
				out = append(out, m[1])
			}
		}
	}
	return out
}

// visualRegimeTag lê do projeto qual tag nomeia o regime de captura visual.
func visualRegimeTag(cfg *config.Config) string {
	// O de-para é TAG → REGIME (`vr-level: vr`): a chave é o que aparece na feature, o
	// valor é o nome do regime. Ler invertido devolvia `vr` como tag e o gate não
	// encontrava cenário nenhum — silenciosamente, porque "não declara cenário visual" é
	// um Skip legítimo para a maioria das features.
	if cfg != nil && cfg.Derived != nil {
		for tag, regime := range cfg.Derived.Regimes {
			r := strings.ToLower(regime)
			if r == "vr" || strings.Contains(r, "visual") {
				return strings.TrimPrefix(tag, "@")
			}
		}
	}
	return "vr-level"
}
