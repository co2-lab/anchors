package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// vr-baseline: o cenário de regressão VISUAL prometido tem imagem de referência.
//
// O VR é a única superfície de prova que nenhum gate alcançava. Um cenário marcado
// `@nivel-vr` declara que aquela tela é provada por CAPTURA — não por asserção em teste
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
		return Skip, "o cenário de regressão visual é declarado na feature — é dela que o confronto parte"
	}
	cenarios := cenariosVR(content, cfg)
	if len(cenarios) == 0 {
		return Skip, "a feature não declara cenário de regressão visual"
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
	return Fail, fmt.Sprintf("%d cenário(s) de regressão visual sem imagem de referência: %s. "+
		"O `@nivel-vr` declara que esta tela é provada por CAPTURA, não por asserção — e sem "+
		"baseline não há contra o que comparar. O cenário existe, o gate de feature o conta "+
		"como coberto, e a prova prometida não acontece.\n\n"+
		"Gere a captura (`%s.%s.png`, ao lado da unidade) ou tire a tag `@nivel-vr` do "+
		"cenário — um cenário que ninguém prova é pior que um cenário ausente, porque parece "+
		"cobertura",
		len(semImagem), strings.Join(semImagem, ", "), filepath.Base(base), semImagem[0])
}

// cenariosVR devolve os códigos de cenário marcados como regressão visual na feature.
//
// O que é "regressão visual" vem do PROJETO (`derived.regimes` no anchors.yaml diz qual
// tag nomeia esse regime), com `nivel-vr` como default — o Anchors não impõe a
// nomenclatura, do mesmo modo que não impõe idioma nem nome de vendor.
func cenariosVR(content string, cfg *config.Config) []string {
	tag := tagDeRegimeVisual(cfg)
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

// tagDeRegimeVisual lê do projeto qual tag nomeia o regime de captura visual.
func tagDeRegimeVisual(cfg *config.Config) string {
	// O de-para é TAG → REGIME (`nivel-vr: vr`): a chave é o que aparece na feature, o
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
	return "nivel-vr"
}
