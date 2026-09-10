// Package pack carrega CONJUNTOS DE OBRIGAÇÕES distribuíveis — o mecanismo que torna a
// conformidade plugável em vez de reescrita por projeto.
//
// O problema que resolve: o mecanismo de `obligations` já era agnóstico (o engine nunca
// soube o que é LGPD — só que existem deveres transversais). Mas cada projeto declarava
// as suas do zero, inline no anchors.yaml. Consequências medidas num projeto real:
//
//   - 3 obrigações escritas à mão, com vocabulário local (`identifica: pessoa`,
//     `titular: usuario` — em português, específicos daquele projeto);
//   - o conteúdo regulatório vazou para a PROSA de 594 specs (47 avisos de revisão e
//     explicações de LGPD dentro de cada modelo);
//   - nenhuma forma de responder "quais deveres se aplicam a mim e onde estou".
//
// Trocar de jurisdição, ali, significaria reescrever prosa — não trocar configuração.
//
// UM PACK É DECLARATIVO. Não é plugin, não executa nada. Código de terceiro rodando
// dentro do gate seria superfície de ataque e quebraria o determinismo — a mesma razão
// pela qual gates externos vivem atrás de `run:` explícito. Um pack é YAML: obrigações,
// com metadados que apontam para a NORMA que as origina.
//
// O QUE É DO PACK E O QUE É DO PROJETO:
//
//	do PACK    — o dever (o que a norma exige), o gatilho canônico, o artigo/cláusula
//	do PROJETO — ONDE isso vive nele (o caminho do handler de exclusão, do log de
//	             auditoria), via placeholders que o anchors.yaml resolve
//
// Sem essa divisão, um pack teria de adivinhar a estrutura de todo projeto — e voltaríamos
// à complexidade do projeto subindo para o framework (CONCEPT §5.2).
package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Dir é onde os packs do projeto vivem, relativo à raiz. Semeado pelo `anchors init`.
const Dir = "packs"

// Pack — um conjunto de deveres de uma norma.
type Pack struct {
	// Name identifica o pack (`lgpd`, `gdpr`, `pci-dss`, `wcag`).
	Name string `yaml:"name"`
	// Domain agrupa packs da mesma natureza (`privacy`, `payment`, `health`,
	// `accessibility`). Serve ao relatório: "quais deveres de privacidade eu tenho".
	Domain string `yaml:"domain"`
	// Jurisdiction é onde a norma vale (`br`, `eu`, `us-ca`, `global`). Um projeto
	// declara em quais opera, e só os packs correspondentes se aplicam — é o que permite
	// o mesmo código servir a mercados diferentes sem duplicar declaração.
	Jurisdiction string `yaml:"jurisdiction,omitempty"`
	// Authority é a norma em si (`Lei 13.709/2018`, `Regulation (EU) 2016/679`,
	// `PCI DSS v4.0`). É o que torna o relatório apresentável a um auditor.
	Authority string `yaml:"authority,omitempty"`
	// Requires são os placeholders que o PROJETO precisa resolver para este pack
	// funcionar (ex.: `purge_handler`). Declará-los permite falhar com uma mensagem útil
	// em vez de silenciosamente não confrontar nada.
	Requires []string `yaml:"requires,omitempty"`
	// Obligations são os deveres. O formato é o mesmo de `obligations:` no anchors.yaml —
	// um pack não inventa mecanismo, só distribui declaração.
	Obligations []PackObligation `yaml:"obligations"`
}

// PackObligation espelha config.Obligation e acrescenta a rastreabilidade normativa.
// Não importa config para evitar ciclo: config carrega packs.
type PackObligation struct {
	Name             string   `yaml:"name"`
	When             string   `yaml:"when"`
	MustAppearIn     []string `yaml:"must_appear_in"`
	IdentifiedBy     string   `yaml:"identified_by,omitempty"`
	IdentifiedAsForm string   `yaml:"identified_as_form,omitempty"`
	Because          string   `yaml:"because,omitempty"`
	// Article é a cláusula da norma que origina o dever (`Art. 18, VI`, `Art. 17`,
	// `Req. 3.4`). É o que permite responder "o que a norma exige e onde estou" — a
	// pergunta que um auditor faz e que nenhum gate respondia.
	Article string `yaml:"article,omitempty"`
}

// Load lê um pack de um arquivo.
func Load(path string) (*Pack, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Pack
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	if p.Name == "" {
		return nil, fmt.Errorf("%s: pack sem `name`", filepath.Base(path))
	}
	if len(p.Obligations) == 0 {
		return nil, fmt.Errorf("%s: pack sem obrigações — um pack vazio não confronta nada", p.Name)
	}
	return &p, nil
}

// LoadAll carrega os packs pedidos, resolvendo os placeholders com os valores do projeto.
//
// `refs` são as referências declaradas no anchors.yaml: um nome (`privacy/lgpd`, resolvido
// dentro de `packs/`) ou um caminho explícito (`./packs/interno.yaml`).
//
// `jurisdictions` filtra: um pack de jurisdição que o projeto não declara é IGNORADO com
// aviso, não carregado em silêncio — declarar `gdpr` num app que só opera no Brasil é
// provavelmente engano, e o silêncio esconderia isso.
func LoadAll(root string, refs []string, values map[string]string, jurisdictions []string) ([]*Pack, []string, error) {
	var out []*Pack
	var avisos []string
	jur := map[string]bool{}
	for _, j := range jurisdictions {
		jur[strings.ToLower(strings.TrimSpace(j))] = true
	}

	for _, ref := range refs {
		path := resolveRef(root, ref)
		p, err := Load(path)
		if err != nil {
			return nil, nil, err
		}
		if j := strings.ToLower(p.Jurisdiction); j != "" && j != "global" && len(jur) > 0 && !jur[j] {
			avisos = append(avisos, fmt.Sprintf(
				"pack `%s` é da jurisdição `%s`, que o projeto não declara em `jurisdictions:` — "+
					"não carregado. Se o app opera lá, declare; senão, remova o pack.", p.Name, p.Jurisdiction))
			continue
		}
		if faltando := unresolved(p, values); len(faltando) > 0 {
			return nil, nil, fmt.Errorf(
				"o pack `%s` precisa que o projeto declare %s em `pack_values:` — sem isso ele "+
					"não sabe ONDE confrontar, e um dever que não aponta para lugar nenhum passa "+
					"verde sem verificar nada",
				p.Name, "`"+strings.Join(faltando, "`, `")+"`")
		}
		resolve(p, values)
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, avisos, nil
}

// resolveRef converte a referência em caminho. Nome curto vira `packs/<nome>.yaml`.
func resolveRef(root, ref string) string {
	if strings.HasSuffix(ref, ".yaml") || strings.HasSuffix(ref, ".yml") {
		if filepath.IsAbs(ref) {
			return ref
		}
		return filepath.Join(root, ref)
	}
	return filepath.Join(root, Dir, ref+".yaml")
}

var placeholderRE = regexp.MustCompile(`\{\{\s*([a-z0-9_]+)\s*\}\}`)

// unresolved lista os placeholders que o pack usa e o projeto não declarou.
func unresolved(p *Pack, values map[string]string) []string {
	falta := map[string]bool{}
	for _, r := range p.Requires {
		if _, ok := values[r]; !ok {
			falta[r] = true
		}
	}
	for _, ob := range p.Obligations {
		for _, alvo := range ob.MustAppearIn {
			for _, m := range placeholderRE.FindAllStringSubmatch(alvo, -1) {
				if _, ok := values[m[1]]; !ok {
					falta[m[1]] = true
				}
			}
		}
	}
	var out []string
	for k := range falta {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// resolve substitui os placeholders pelos caminhos do projeto.
func resolve(p *Pack, values map[string]string) {
	for i := range p.Obligations {
		alvos := p.Obligations[i].MustAppearIn
		for j, alvo := range alvos {
			alvos[j] = placeholderRE.ReplaceAllStringFunc(alvo, func(m string) string {
				k := placeholderRE.FindStringSubmatch(m)[1]
				if v, ok := values[k]; ok {
					return v
				}
				return m
			})
		}
	}
}
