package config

import "strings"

// --- as DOCUMENTAÇÕES que o projeto deve ter ---
//
// O Anchors sempre cobrou o que está DENTRO da unidade — a spec, a feature, o teste, o
// código. A documentação ficava como uma frase no entregável do card ("a documentação
// evolui junto"), e uma frase assim não diz QUAL documentação.
//
// O buraco é real e depende do tipo de projeto. Uma API tem um contrato que vive fora do
// código — o OpenAPI — e um endpoint novo que não entra nele é invisível para quem
// consome. Um design system tem o catálogo de componentes. Um projeto com banco tem o
// esquema. Nada disso o gate de trinca alcança, porque não é arquivo da unidade: é o
// artefato AGREGADO que várias unidades alimentam.
//
// E há um que todo projeto tem: a ARQUITETURA. O Anchors já regula a estrutura — as
// camadas, o `needs:` entre planos, as fronteiras que o `layer-boundary` cobra — mas essa
// regulação vive espalhada no `anchors.yaml` e no mapa, em formato de máquina. Quem chega
// no projeto precisa da mesma informação em formato de gente, e o C4 é a notação que
// separa isso em níveis (contexto → contêineres → componentes → código), o que impede a
// falha clássica do diagrama de arquitetura: um desenho só, com tudo, ilegível.
//
// A DECLARAÇÃO é o ato. Um projeto que não declara `docs:` não é cobrado — cobrar OpenAPI
// de quem não tem API seria ruído. Declarar é dizer "esta documentação é contrato, cobre-a
// de mim".

// Docs declara as documentações agregadas do projeto.
type Docs struct {
	// Required são as documentações obrigatórias. Cada uma nomeia o arquivo, o que ela
	// existe para responder, e QUANDO precisa ser tocada.
	Required []DocArtifact `yaml:"required,omitempty"`
}

// DocArtifact é uma documentação agregada — um artefato que várias unidades alimentam.
type DocArtifact struct {
	// Kind é o tipo, e é o que dá ao agente a instrução certa: `openapi` pede o contrato
	// dos endpoints, `c4` pede os quatro níveis, `schema` pede as tabelas.
	//
	// Tipo desconhecido não é erro: o projeto pode ter uma documentação que o Anchors não
	// conhece, e recusá-la faria o mecanismo servir só ao que já foi previsto. O que o
	// Anchors não conhece ele não instrui — mas cobra a existência.
	Kind string `yaml:"kind"`
	// Path é o arquivo. É o que o gate confere e o que o `anchors next` nomeia.
	Path string `yaml:"path"`
	// Trigger são as camadas cuja alteração OBRIGA tocar esta doc.
	//
	// É o campo que transforma "documente" em algo verificável: mexer numa lambda de
	// rota muda o contrato da API, e o OpenAPI tem de acompanhar. Mexer num utilitário
	// interno não muda contrato nenhum, e cobrar ali seria cobrar por cobrar — o agente
	// aprenderia a atualizar o arquivo sem pensar, que é pior que não atualizar.
	//
	// Vazio = a doc não tem gatilho por camada (é o caso do C4, que muda quando a
	// ESTRUTURA muda, não quando uma unidade muda).
	Trigger []string `yaml:"trigger,omitempty"`
	// Why é por que esta documentação existe. Vai para o texto do card: um agente que
	// sabe o que a doc responde escreve melhor do que um que só sabe o caminho dela.
	Why string `yaml:"why,omitempty"`
}

// KindOpenAPI e companhia são os tipos que o Anchors sabe instruir.
const (
	KindOpenAPI   = "openapi"
	KindC4        = "c4"
	KindSchema    = "schema"
	KindComponent = "components"
	KindADR       = "adr"
	KindRunbook   = "runbook"
)

// TriggeredBy diz se alterar uma unidade desta camada obriga tocar esta documentação.
func (d DocArtifact) TriggeredBy(layer string) bool {
	for _, l := range d.Trigger {
		if strings.EqualFold(strings.TrimSpace(l), layer) {
			return true
		}
	}
	return false
}

// RequiredFor devolve as documentações que a alteração de uma camada obriga tocar.
func (c *Config) RequiredFor(layer string) []DocArtifact {
	if c == nil || c.Docs == nil {
		return nil
	}
	var out []DocArtifact
	for _, d := range c.Docs.Required {
		if d.TriggeredBy(layer) {
			out = append(out, d)
		}
	}
	return out
}

// AllRequiredDocs devolve todas as documentações declaradas.
func (c *Config) AllRequiredDocs() []DocArtifact {
	if c == nil || c.Docs == nil {
		return nil
	}
	return c.Docs.Required
}
