package config

import "strings"

// --- os CONTÊINERES: o que roda separado ---
//
// A Estrutura declara CAMADAS — o que cada peça é. O que ela não diz é o que roda
// SEPARADO: o app, a API, a infraestrutura, o banco. E essa é justamente a pergunta do
// nível 2 do C4, e a que decide quantos diagramas o nível 3 tem.
//
// A distinção importa e não é derivável do layout de pastas. `apps/mobile/` e
// `packages/lambdas/` parecem dois contêineres porque estão em pastas diferentes — mas
// `packages/shared/` roda DENTRO da API, e `packages/infra/` não roda em lugar nenhum (é
// declaração de infraestrutura). Inferir do caminho acertaria por acaso neste projeto e
// erraria no próximo, sem nada avisando.
//
// O BANCO é o caso que torna a declaração necessária. Ele é contêiner — aparece no nível
// 2, com o protocolo da conversa — e ao mesmo tempo é EXTERNO: não temos componentes
// dentro dele, e um diagrama de nível 3 do banco seria uma caixa vazia afirmando que
// olhamos lá dentro. `external: true` é o que diz "existe, conversa conosco, e não é
// nosso por dentro".

// Container é o que roda separado — a unidade do nível 2 do C4.
type Container struct {
	// Name é como ele aparece no diagrama.
	Name string `yaml:"name"`
	// Description é a linha sob o nome. Uma caixa chamada "api" não diz nada; "as rotas
	// que o app consulta" diz.
	Description string `yaml:"description,omitempty"`
	// Layers são as camadas que rodam DENTRO dele. É o que dá o nível 3.
	//
	// Vazio num contêiner interno é declaração de que ele não tem camada nossa — e o
	// nível 3 dele sai dizendo isso, em vez de um diagrama vazio.
	Layers []string `yaml:"layers,omitempty"`
	// External marca o que conversa conosco e não é nosso por dentro: o banco, o gateway
	// de pagamento, o provedor de identidade.
	//
	// Um contêiner externo aparece no nível 2 (a conversa existe, e o protocolo dela
	// importa) e NÃO ganha nível 3: não temos componentes lá dentro, e desenhar um seria
	// afirmar um conhecimento que não temos.
	External bool `yaml:"external,omitempty"`
	// Talks são os contêineres com que ele conversa, e COMO.
	Talks []Talk `yaml:"talks,omitempty"`
}

// Talk é uma conversa entre contêineres, com o protocolo.
//
// O PROTOCOLO é obrigatório por desenho: uma seta sem ele diz que os dois se falam e não
// diz o que acontece quando a conversa falha — que é a única coisa que um diagrama de
// contêineres tem a dizer sobre risco.
type Talk struct {
	To       string `yaml:"to"`
	Protocol string `yaml:"protocol"`
	// Why é o que passa por ali. Opcional, e útil quando o protocolo não basta:
	// "HTTPS/JSON" não diz se são consultas ou comandos.
	Why string `yaml:"why,omitempty"`
}

// Containers devolve os contêineres declarados.
func (c *Config) Containers() []Container {
	if c == nil {
		return nil
	}
	return c.ContainersDecl
}

// InternalContainers são os que ganham diagrama de nível 3.
func (c *Config) InternalContainers() []Container {
	var out []Container
	for _, k := range c.Containers() {
		if !k.External {
			out = append(out, k)
		}
	}
	return out
}

// ContainerOfLayer diz em que contêiner uma camada roda.
//
// Devolve vazio quando nenhum a declara — e isso é informação, não erro: uma camada que
// não está em contêiner nenhum não aparece em diagrama de nível 3, e o template pode
// dizê-lo em vez de escondê-la.
func (c *Config) ContainerOfLayer(layer string) string {
	for _, k := range c.Containers() {
		for _, l := range k.Layers {
			if strings.EqualFold(strings.TrimSpace(l), layer) {
				return k.Name
			}
		}
	}
	return ""
}

// OrphanLayers são as camadas que existem e não estão em contêiner nenhum.
//
// Serve ao template e ao `doctor`: uma camada fora de todo contêiner é ou uma declaração
// esquecida, ou uma peça que não roda em lugar nenhum — e as duas merecem ser ditas.
func (c *Config) OrphanLayers(existentes []string) []string {
	var out []string
	for _, l := range existentes {
		if c.ContainerOfLayer(l) == "" {
			out = append(out, l)
		}
	}
	return out
}
