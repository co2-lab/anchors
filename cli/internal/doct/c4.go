package doct

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- o C4 como o C4 é ---
//
// A regra central do modelo é que cada nível AMPLIA UMA CAIXA do anterior. O nível 3 é o
// zoom de UM contêiner — não do sistema. Um diagrama de componentes misturando o app, a
// API e a infraestrutura não é nível 3 de coisa nenhuma: é a falha que o C4 existe para
// evitar, um desenho só com tudo dentro.
//
// A primeira versão deste template errava exatamente nisso — desenhava as oito camadas do
// projeto num diagrama e o chamava de nível 3.
//
// E CAMADA NÃO É COMPONENTE. `screen`, `lambdas`, `shared` são agrupamentos de código no
// repositório; o componente do C4 é a peça com responsabilidade dentro de um contêiner. A
// declaração `containers.layers` é a ponte entre os dois vocabulários: ela diz que camadas
// rodam em que contêiner, e é o que permite o nível 3 mostrar as UNIDADES daquele
// contêiner em vez da lista de camadas do repositório.
//
// O BANCO é contêiner, e essa é a parte do C4 que mais se erra. Contêiner é o que executa
// ou armazena dado — aplicação, serviço, banco, fila, sistema de arquivos — e não "processo
// que escrevemos". Ele aparece no nível 2 com o protocolo da conversa; sendo de terceiro,
// não ganha nível 3, porque não temos componentes lá dentro e desenhá-los seria afirmar um
// conhecimento que não temos.

// Container é o que o template vê de um contêiner.
type Container struct {
	Name        string
	Description string
	External    bool
	Layers      []string
	Talks       []config.Talk
	// Units são as unidades que rodam nele — o material do nível 3.
	Units []Spec
}

// ID é o identificador seguro para o Mermaid.
func (k Container) ID() string { return MermaidID(k.Name) }

// fnContainers devolve os contêineres declarados, com as unidades de cada um.
func (c *Compiler) fnContainers() []Container {
	if c.Config == nil {
		return nil
	}
	var out []Container
	for _, k := range c.Config.Containers() {
		box := Container{
			Name: k.Name, Description: k.Description,
			External: k.External, Layers: k.Layers, Talks: k.Talks,
		}
		for _, l := range k.Layers {
			for _, s := range c.specs {
				if strings.EqualFold(s.Layer, strings.TrimSpace(l)) {
					box.Units = append(box.Units, s)
				}
			}
		}
		sort.Slice(box.Units, func(i, j int) bool {
			if box.Units[i].Layer != box.Units[j].Layer {
				return box.Units[i].Layer < box.Units[j].Layer
			}
			return box.Units[i].Code < box.Units[j].Code
		})
		out = append(out, box)
	}
	return out
}

// fnInternalContainers são os que ganham diagrama de nível 3.
//
// O externo fica de fora por definição do modelo, não por economia: não temos componentes
// dentro dele. Um nível 3 do banco seria uma caixa vazia afirmando que olhamos lá dentro.
func (c *Compiler) fnInternalContainers() []Container {
	var out []Container
	for _, k := range c.fnContainers() {
		if !k.External {
			out = append(out, k)
		}
	}
	return out
}

// fnOrphanLayers são as camadas que existem no projeto e nenhum contêiner declara.
//
// Existe para o template DIZER isso, e não para escondê-las. Uma camada fora de todo
// contêiner não aparece em nível 3 nenhum — e um diagrama que a omite em silêncio afirma,
// por ausência, que ela não existe.
func (c *Compiler) fnOrphanLayers() []string {
	if c.Config == nil {
		return nil
	}
	return c.Config.OrphanLayers(c.fnLayers())
}

// readFile lê um documento compilado — auxiliar dos testes de integração do compilador.
func readFile(root, nome string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, OutDir, filepath.FromSlash(nome)))
	return string(b), err
}
