package initx

import (
	"sort"

	"github.com/co2-lab/anchors/internal/config"
)

// Este arquivo contém as funções PURAS de decisão do init: recebem dados (a
// proposta + as escolhas do usuário) e devolvem dados (o Config ajustado). Nenhuma
// delas toca em TUI — por isso são testáveis com `go test` normal. A casca de
// prompt (cmd/anchors/init.go) só coleta respostas e chama estas funções.

// CodeLayerNames devolve os nomes das camadas de código da proposta (as que o
// usuário vai escolher manter). Ordenado, determinístico.
func CodeLayerNames(cfg *config.Config) []string {
	var names []string
	for name, l := range cfg.Layers {
		if l.Kind == "code" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// PruneCodeLayers remove as camadas de código que NÃO estão em `keep`. As camadas
// de artefato (spec/feature/test/guide/doc) nunca são removidas aqui. Devolve o
// próprio cfg mutado, para encadear.
func PruneCodeLayers(cfg *config.Config, keep map[string]bool) *config.Config {
	for name, l := range cfg.Layers {
		if l.Kind == "code" && !keep[name] {
			delete(cfg.Layers, name)
		}
	}
	return cfg
}

// Tags coleta todas as tags declaradas nas camadas do cfg (dedup, ordenado). É a
// lista de tags candidatas que um guide pode reger no governs.
func Tags(cfg *config.Config) []string {
	seen := map[string]bool{}
	for _, l := range cfg.Layers {
		for _, t := range l.Tags {
			seen[t] = true
		}
	}
	var tags []string
	for t := range seen {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}

// BuildGovernRules monta as regras de governs a partir das respostas do usuário:
// answers[guidePath] = tag (ou "" / "(nenhuma)" para pular). Puro e testável.
func BuildGovernRules(answers map[string]string) []config.GovernRule {
	// ordena por guide para saída determinística
	var guides []string
	for g := range answers {
		guides = append(guides, g)
	}
	sort.Strings(guides)

	var rules []config.GovernRule
	for _, g := range guides {
		tag := answers[g]
		if tag == "" || tag == NoneTag {
			continue
		}
		rules = append(rules, config.GovernRule{From: g, Governs: tag})
	}
	return rules
}

// NoneTag é a opção "nenhuma" apresentada ao usuário no governs.
const NoneTag = "(nenhuma)"
