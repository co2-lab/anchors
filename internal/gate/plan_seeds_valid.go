package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// plan-seeds-valid: um PLANO semeia artefatos (as specs que vão nascer). Este gate
// confronta cada caminho semeado contra a Estrutura, ANTES de alguém executar o plano.
//
// O caso que o motivou, observado três rodadas seguidas: um plano listava
// `packages/backend/models/metadata.spec.md — **nasce**`, mas `models/` é camada
// RECONHECIDA (`regime: declarativo`) e não tem spec. Toda vez alguém gastou uma rodada
// redescobrindo a contradição — e a defesa (gate, `new`, `work`) só age na hora de
// executar. A origem do erro é o plano, e o plano é uma camada declarada: é checável.
//
// Confronta duas coisas:
//  1. spec semeada em camada DECLARATIVA (não tem spec por definição);
//  2. caminho semeado que não pertence a NENHUMA camada (typo, ou camada a declarar).
//
// Deliberadamente NÃO cobra: se a spec já existe, se o plano está atualizado, ou o
// conteúdo dos itens — só a validade estrutural do que ele promete criar.
func filepathBase(p string) string { return filepath.Base(p) }

var planSeedRE = regexp.MustCompile("`([^`]+\\.spec\\.md)`")

// doctrineSeedRE — a DOUTRINA DE PRODUTO semeada por um plano.
//
// Reconhecida a' parte, e nao junto das specs no regex acima, porque a validacao daqui
// nao se aplica a ela: o gate confere a CAMADA DO ALVO que a spec descreveria
// (`Foo.spec.md` -> `Foo.ts`), e doutrina nao tem alvo — ela E' o artefato. Cobrar dela
// uma camada de codigo acusaria toda doutrina semeada de "camada nao declarada".
//
// O scan a registra como semente (aresta `seeds` no grafo, que e' o que faz o plano
// deixar de ficar orfao); aqui ela so' precisa nao ser confundida com spec.
var doctrineSeedRE = regexp.MustCompile("`([^`]+\\.doctrine\\.md)`")

func checkPlanSeedsValid(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, i18n.T("gate.plan_seeds_valid.skip_not_plan")
	}
	if cfg == nil {
		return Pending, i18n.T("gate.plan_seeds_valid.pending_no_config")
	}

	seeds := map[string]bool{}
	for _, m := range planSeedRE.FindAllStringSubmatch(content, -1) {
		s := m[1]
		// MOLDES (`_TEMPLATE_*.spec.md`) não são specs semeadas — são o gabarito que a
		// spec copia. O anchors.yaml já os exclui da camada `spec`; excluí-los aqui
		// evita acusar um plano por citar o template que ele manda usar.
		if strings.HasPrefix(filepathBase(s), "_TEMPLATE") {
			continue
		}
		seeds[s] = true
	}
	// SEEDED DOCTRINE: the only rule that holds for it is living under `product/`.
	//
	// No target layer is demanded (it has no target), nor existence (the doctrine is
	// about to BE BORN — that is what seeding means). The PLACE is demanded, because the
	// `product` kind comes from the path: a doctrine seeded outside `product/` is born as
	// a plain doc, and the spec realizing it points at a file the map does not recognise
	// as doctrine.
	var outsideProduct []string
	for _, m := range doctrineSeedRE.FindAllStringSubmatch(content, -1) {
		d := m[1]
		if strings.HasPrefix(filepathBase(d), "_TEMPLATE") {
			continue
		}
		if !strings.Contains(d, "/") || !rootDirExists(root, d) {
			continue // citacao em prosa, nao semeadura — mesma regra das specs
		}
		if !strings.HasPrefix(d, "product/") {
			outsideProduct = append(outsideProduct, d)
		}
	}

	// The "no spec seeded" guard comes AFTER the doctrine check, and the order is the
	// decision: a plan seeding only product doctrine (no spec at all) is legitimate, and
	// with the guard first it came out as Skip — no verdict, with a doctrine outside
	// `product/` passing in silence. Found by probe.
	if len(seeds) == 0 && len(outsideProduct) == 0 {
		return Skip, i18n.T("gate.plan_seeds_valid.skip_no_seeds")
	}

	var declarativa, semCamada []string
	for s := range seeds {
		// o ALVO que a spec descreveria — a camada é dele, não do .spec.md.
		base := strings.TrimSuffix(s, ".spec.md")
		layer := ""
		for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
			if lay, _ := scan.Classify(base+ext, cfg); lay != "" {
				layer = lay
				break
			}
		}
		if layer == "" {
			// Um plano CITA nomes de spec em prosa ("ver SubscriptionScreen.spec.md",
			// "igual ao _TEMPLATE_COMPONENT.spec.md") — isso não é semeadura, é
			// referência. Só tratamos como SEMEADO o que tem cara de caminho real
			// (contém `/`); um nome solto não é promessa de criar arquivo.
			// Um fragmento de caminho ("features/x/Y.spec.md", sem o prefixo do app)
			// é abreviação em prosa, não semeadura. Só acusamos o que começa numa RAIZ
			// REAL do repositório — aí o caminho é uma promessa concreta, e não casar
			// nenhuma camada é erro de verdade (typo ou camada a declarar).
			if strings.Contains(s, "/") && rootDirExists(root, s) {
				semCamada = append(semCamada, s)
			}
			continue
		}
		if l, ok := cfg.Layers[layer]; ok && l.Regime == "declarativo" {
			declarativa = append(declarativa, i18n.T("gate.plan_seeds_valid.declarative_item", s, layer))
		}
	}

	var parts []string
	if len(outsideProduct) > 0 {
		sort.Strings(outsideProduct)
		parts = append(parts, i18n.T("gate.plan_seeds_valid.part_doctrine_outside", strings.Join(outsideProduct, ", ")))
	}
	if len(declarativa) > 0 {
		sort.Strings(declarativa)
		parts = append(parts, i18n.T("gate.plan_seeds_valid.part_declarative", strings.Join(declarativa, ", ")))
	}
	if len(semCamada) > 0 {
		sort.Strings(semCamada)
		parts = append(parts, i18n.T("gate.plan_seeds_valid.part_undeclared", strings.Join(semCamada, ", ")))
	}
	if len(parts) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.plan_seeds_valid.fail_invalid_seeds", strings.Join(parts, "; "))
}

// rootDirExists diz se o primeiro segmento do caminho é um diretório REAL na raiz do
// projeto. Distingue caminho de verdade ("packages/backend/...") de abreviação em prosa
// ("features/subscription/..." sem o prefixo do app).
func rootDirExists(root, rel string) bool {
	// O `first == rel` que havia aqui era INALCANCAVEL: o unico chamador ja guarda com
	// `strings.Contains(s, "/")`, e sem barra o SplitN devolve a string inteira. Achado por
	// mutacao — trocar a condicao deixava a suite verde, e a razao nao era falta de teste:
	// era codigo defensivo que a via de producao nao alcanca. Teste que o cobrisse teria de
	// chamar o helper direto, pinando um ramo que o gate nunca percorre.
	first := strings.SplitN(rel, "/", 2)[0]
	if first == "" {
		return false
	}
	fi, err := os.Stat(filepath.Join(root, first))
	return err == nil && fi.IsDir()
}
