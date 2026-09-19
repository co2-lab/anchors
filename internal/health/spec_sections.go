package health

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// As SEÇÕES DE CRUZAMENTO e o gate que cada uma alimenta.
//
// POR QUE ESTE CHECK EXISTE. A seção é a unidade de cruzamento do Anchors — gates
// procuram seção NOMEADA (`domain-declared` procura `## Domínio`, `route-declared`
// procura a rota) e `rule-types` liga a LETRA do código ao TÍTULO da seção. Conteúdo
// fora de seção não é confrontável: a decisão pode estar escrita, correta e completa, e
// ainda assim nenhum gate a alcança.
//
// O buraco que ele fecha: `Default: true` no catálogo do `anchors new spec` só age no
// MOMENTO DA CRIAÇÃO. Uma spec já escrita não passa a reprovar porque o catálogo mudou —
// e os gates que cobram seção só cobram as seções que a spec já tem. Projeto que nasceu
// com um catálogo pobre fica sem nenhum caminho de correção: a única via seria alguém
// comparar à mão contra o guide.
//
// Medido no projeto que originou o achado (ver FINDINGS-spec-sections.md): 85 specs
// colapsaram para SEIS títulos de seção; 25 declaravam `layer: screen` e 24 não tinham
// rota. O `route-declared` existe, com doutrina escrita, e ficou cego em 96% dos alvos
// sem que uma linha de saída reclamasse.
//
// AGREGA POR CAMADA, NUNCA POR ARQUIVO. 37 specs × 6 seções ausentes são 200 linhas que
// ninguém lê — o erro que o próprio anchors.yaml deste repo documenta ("ligar tudo de uma
// vez produz centenas de issues no primeiro dia e ninguém lê a lista"). Uma linha por
// seção, com o número, é o que permite corrigir POR LEVA.
type specSection struct {
	// Key é a chave da seção no catálogo do `anchors new spec` (`--with <key>`).
	Key string
	// TitleKeys são as chaves i18n dos títulos que CONTAM como esta seção presente. São
	// chaves, e não literais, porque o título da seção é TEXTO TRADUZIDO: uma spec escrita
	// num projeto `lang: en` traz "Navigation", e procurar só "Navegação" acusaria ausência
	// onde a seção existe. O confronto roda contra os títulos de TODOS os idiomas (ver
	// sectionTitles) — a spec pode ter nascido sob outro idioma que o atual.
	TitleKeys []string
	// ExtraTitles são títulos que valem além dos traduzidos: o léxico próprio de projetos
	// (`section_titles` no anchors.yaml) e grafias sem acento.
	ExtraTitles []string
	// Gate é o gate que esta seção alimenta. Vazio = a seção vale por si.
	Gate string
	// Layers restringe a quais camadas a seção se aplica. Vazio = todas.
	// Cobrar rota de um hook ou de uma regra pura é falso-positivo — o vício do
	// validador legado que só conhecia tela.
	Layers []string
	// WhyKey é a chave i18n da frase que explica o CUSTO de omitir — não o que a seção
	// contém.
	WhyKey string
}

var specSections = []specSection{
	{
		Key:         "navigation",
		TitleKeys:   []string{"section.title.navigation"},
		ExtraTitles: []string{"Navegacao"},
		Gate:        "route-declared",
		Layers:      []string{"screen"},
		WhyKey:      "health.spec_section.why.navigation",
	},
	{
		Key:       "data-contract",
		TitleKeys: []string{"section.title.data_contract"},
		Gate:      "dependency-honored",
		Layers:    []string{"screen", "component"},
		WhyKey:    "health.spec_section.why.data_contract",
	},
	{
		Key:       "data-states",
		TitleKeys: []string{"section.title.data_states"},
		Gate:      "scenario-asserts",
		Layers:    []string{"screen", "component"},
		WhyKey:    "health.spec_section.why.data_states",
	},
	{
		Key:         "testids",
		TitleKeys:   []string{"section.title.testids"},
		ExtraTitles: []string{"Test IDs", "TestIDs"},
		Gate:        "testid-consistent",
		Layers:      []string{"screen", "component"},
		WhyKey:      "health.spec_section.why.testids",
	},
	{
		Key:         "domain",
		TitleKeys:   []string{"section.title.domain"},
		ExtraTitles: []string{"Dominio", "Modelo de Dado"},
		Gate:        "domain-declared",
		WhyKey:      "health.spec_section.why.domain",
	},
	{
		Key:       "states",
		TitleKeys: []string{"section.title.states"},
		Layers:    []string{"screen", "component"},
		WhyKey:    "health.spec_section.why.states",
	},
}

// specInfo é o que se lê de uma spec UMA vez: a camada e o conjunto de títulos de seção.
type specInfo struct {
	layer  string
	titles map[string]bool
}

// checkSpecSections confronta as specs do mapa contra as seções de cruzamento, agregando
// por camada. Emite dois checks distintos porque a AÇÃO CORRETIVA é diferente:
//
//	secao-ausente     → o gate está DECLARADO e cego. Ação: editar as specs.
//	secao-recomendada → o gate nem está declarado. Ação: adotar a seção E o gate.
func checkSpecSections(g *mapx.Graph, cfg *config.Config, root string) []Finding {
	if g == nil || cfg == nil {
		return nil
	}

	declaredGate := map[string]bool{}
	for _, gt := range cfg.Gates {
		declaredGate[gt.Name] = true
	}

	specs := readSpecSections(g, root)
	if len(specs) == 0 {
		return nil
	}

	var out []Finding
	for _, sec := range specSections {
		accepted := sectionTitles(sec)
		targets, missing := 0, 0
		for _, sp := range specs {
			if len(sec.Layers) > 0 && !hasLayer(sec.Layers, sp.layer) {
				continue
			}
			targets++
			if !hasSection(sp.titles, accepted) {
				missing++
			}
		}
		// Só reporta o que é PADRÃO, não o caso isolado: uma spec sem a seção é decisão
		// de quem escreveu; a maioria sem ela é o catálogo (ou o hábito) que não a pede.
		if targets == 0 || missing == 0 || missing*2 <= targets {
			continue
		}

		why := i18n.T(sec.WhyKey)
		check, sev := "secao-recomendada", Info
		detail := i18n.T("health.spec_section.recommended", missing, targets, why, sec.Key)
		switch {
		case sec.Gate != "" && declaredGate[sec.Gate]:
			// O caso grave: há um gate que JÁ PODERIA estar medindo e está cego.
			check, sev = "secao-ausente", Warn
			detail = i18n.T("health.spec_section.blind_gate", missing, targets, sec.Gate, why)
		case sec.Gate != "":
			detail += i18n.T("health.spec_section.and_declare_gate", sec.Gate)
		}

		out = append(out, Finding{
			Check:    check,
			Severity: sev,
			Subject:  sec.Key,
			Detail:   detail,
		})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Check < out[j].Check })
	return out
}

// readSpecSections lê cada spec do mapa UMA vez e guarda seus títulos de seção — em vez
// de reler o arquivo a cada seção procurada.
func readSpecSections(g *mapx.Graph, root string) []specInfo {
	var out []specInfo
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil {
			continue
		}
		// A camada vem do HEADER da spec (`layer: screen`), com o nó do mapa como
		// fallback. O header é quem diz a camada da UNIDADE que a spec descreve; o nó
		// costuma dizer `spec` (a camada do ARQUIVO, que casou o pattern `**/*.spec.md`).
		// O mapx já resolve isso na construção, mas ler o header aqui torna o check
		// independente de mapa recém-construído — e mapa velho é o caso comum justamente
		// nos projetos que este check existe para diagnosticar.
		info := specInfo{layer: n.Layer, titles: map[string]bool{}}
		for line := range strings.SplitSeq(string(b), "\n") {
			t := strings.TrimSpace(line)
			if l, ok := strings.CutPrefix(t, "layer:"); ok {
				if l = strings.TrimSpace(l); l != "" {
					info.layer = l
				}
				continue
			}
			if !strings.HasPrefix(t, "#") {
				continue
			}
			// Normaliza "### {CODE}-B01 — regra" e "## Domínio" ao texto do título.
			info.titles[strings.ToLower(strings.TrimSpace(strings.TrimLeft(t, "#")))] = true
		}
		out = append(out, info)
	}
	return out
}

// sectionTitles devolve os títulos aceitos para a seção, em TODOS os idiomas conhecidos.
// Varrer todos (e não só o atual) porque a spec pode ter nascido sob outro idioma — ou o
// projeto ter trocado de `lang` depois de escrever metade do acervo.
func sectionTitles(sec specSection) []string {
	out := append([]string{}, sec.ExtraTitles...)
	for _, k := range sec.TitleKeys {
		out = append(out, i18n.AllTranslations(k)...)
	}
	return out
}

func hasLayer(layers []string, l string) bool {
	for _, c := range layers {
		if strings.EqualFold(c, l) {
			return true
		}
	}
	return false
}

func hasSection(titles map[string]bool, accepted []string) bool {
	for _, a := range accepted {
		if titles[strings.ToLower(a)] {
			return true
		}
	}
	return false
}
