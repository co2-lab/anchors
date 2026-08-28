package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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

func checkPlanSeedsValid(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "não é um plano — o gate confronta o que um plano SEMEIA"
	}
	if cfg == nil {
		return Pending, "sem config carregada — o gate precisa da Estrutura"
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
	if len(seeds) == 0 {
		return Skip, "o plano não semeia nenhuma spec (nenhum caminho `*.spec.md` citado)"
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
			declarativa = append(declarativa, fmt.Sprintf("%s (camada `%s`)", s, layer))
		}
	}

	var parts []string
	if len(declarativa) > 0 {
		sort.Strings(declarativa)
		parts = append(parts, "semeia spec em camada RECONHECIDA (declarativa), que não tem spec: "+
			strings.Join(declarativa, ", ")+" — a decisão pertence à camada que DECIDE, não a esta")
	}
	if len(semCamada) > 0 {
		sort.Strings(semCamada)
		parts = append(parts, "semeia spec em caminho que não pertence a nenhuma camada declarada: "+
			strings.Join(semCamada, ", ")+" — confira o caminho, ou declare a camada em `layers:`")
	}
	if len(parts) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(parts, "; ") + ". Corrija o PLANO: executá-lo produziria " +
		"artefato que a Estrutura proíbe (ou fará quem executa gastar uma rodada descobrindo isso)."
}

// rootDirExists diz se o primeiro segmento do caminho é um diretório REAL na raiz do
// projeto. Distingue caminho de verdade ("packages/backend/...") de abreviação em prosa
// ("features/subscription/..." sem o prefixo do app).
func rootDirExists(root, rel string) bool {
	first := strings.SplitN(rel, "/", 2)[0]
	if first == "" || first == rel {
		return false
	}
	fi, err := os.Stat(filepath.Join(root, first))
	return err == nil && fi.IsDir()
}
