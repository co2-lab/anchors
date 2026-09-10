package doct

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func projetoComContainers(t *testing.T) *Compiler {
	t.Helper()
	specs := map[string]string{
		"app/Tela.spec.md":   "---\ncode: TELAX\nlayer: screen\n---\n\n# Tela — a tela\n",
		"api/Rota.spec.md":   "---\ncode: ROTAX\nlayer: lambdas\n---\n\n# Rota — a rota\n",
		"api/Util.spec.md":   "---\ncode: UTILX\nlayer: shared\n---\n\n# Util — o util\n",
		"solto/Orfa.spec.md": "---\ncode: ORFAX\nlayer: orfa\n---\n\n# Orfa — sem contêiner\n",
	}
	root, g := projetoDeTeste(t, specs)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	c.Config = &config.Config{ContainersDecl: []config.Container{
		{Name: "app", Description: "a interface", Layers: []string{"screen"},
			Talks: []config.Talk{{To: "api", Protocol: "HTTPS/JSON"}}},
		{Name: "api", Description: "as rotas", Layers: []string{"lambdas", "shared"},
			Talks: []config.Talk{{To: "banco", Protocol: "SQL/TLS"}}},
		// O BANCO é contêiner — contêiner é o que executa OU ARMAZENA dado, e não
		// "processo que escrevemos". Externo porque é de terceiro.
		{Name: "banco", Description: "o que persiste", External: true},
	}}
	return c
}

// O NÍVEL 3 É POR CONTÊINER, e o externo não tem um.
//
// É a regra central do C4: cada nível amplia UMA caixa do anterior. Um diagrama de
// componentes misturando app, API e banco não é nível 3 de coisa nenhuma — é a falha que
// o modelo existe para evitar. E o banco não ganha nível 3 porque não temos componentes
// lá dentro; desenhá-los afirmaria um conhecimento que não temos.
func TestInternalContainers_oExternoNaoGanhaNivel3(t *testing.T) {
	c := projetoComContainers(t)

	if todos := c.fnContainers(); len(todos) != 3 {
		t.Fatalf("contêineres = %d, queria 3 (o banco É contêiner)", len(todos))
	}
	internos := c.fnInternalContainers()
	if len(internos) != 2 {
		t.Fatalf("internos = %d, queria 2 — o externo não tem nível 3", len(internos))
	}
	for _, k := range internos {
		if k.Name == "banco" {
			t.Error("o banco ganhou nível 3 — não temos componentes dentro dele")
		}
	}
}

// As UNIDADES de um contêiner são as das camadas que ele declara. É a ponte entre os dois
// vocabulários: camada é agrupamento de código, componente é peça dentro de um contêiner.
func TestContainers_unidadesVemDasCamadasDeclaradas(t *testing.T) {
	c := projetoComContainers(t)
	por := map[string][]string{}
	for _, k := range c.fnContainers() {
		for _, u := range k.Units {
			por[k.Name] = append(por[k.Name], u.Code)
		}
	}
	if len(por["app"]) != 1 || por["app"][0] != "TELAX" {
		t.Errorf("app = %v, queria [TELAX]", por["app"])
	}
	if len(por["api"]) != 2 {
		t.Errorf("api = %v, queria ROTAX e UTILX (lambdas + shared)", por["api"])
	}
	if len(por["banco"]) != 0 {
		t.Errorf("banco = %v — ele não tem unidade nossa", por["banco"])
	}
}

// A camada FORA de todo contêiner é dita, não escondida: um diagrama que a omite em
// silêncio afirma, por ausência, que ela não existe.
func TestOrphanLayers_saoDitas(t *testing.T) {
	if orfas := projetoComContainers(t).fnOrphanLayers(); len(orfas) != 1 || orfas[0] != "orfa" {
		t.Errorf("órfãs = %v, queria [orfa]", orfas)
	}
}

// O nível 2 traz o PROTOCOLO em cada seta, o banco entre as caixas, e há um nível 3 por
// contêiner interno — nenhum para o externo.
func TestBuild_oC4SegueOModelo(t *testing.T) {
	c := projetoComContainers(t)
	if _, _, err := c.InitScaffolds(false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, err := readFile(c.Root, "arquitetura.md")
	if err != nil {
		t.Fatal(err)
	}

	// O protocolo: uma seta sem ele diz que os dois se falam e não diz o que acontece
	// quando a conversa falha — a única coisa que um nível 2 tem a dizer sobre risco.
	for _, p := range []string{"HTTPS/JSON", "SQL/TLS"} {
		if !strings.Contains(b, p) {
			t.Errorf("o protocolo %q não entrou no diagrama", p)
		}
	}
	if !strings.Contains(b, "o que persiste") {
		t.Error("o banco não aparece no nível 2 — contêiner é o que executa OU ARMAZENA")
	}
	if !strings.Contains(b, "### app") || !strings.Contains(b, "### api") {
		t.Errorf("falta o nível 3 de um contêiner interno:\n%s", b)
	}
	if strings.Contains(b, "### banco") {
		t.Error("o banco ganhou seção de nível 3")
	}
}

// Sem contêiner declarado, o documento DIZ isso em vez de sair com diagramas vazios.
func TestBuild_semContainerDeclaradoAvisa(t *testing.T) {
	root, g := projetoDeTeste(t, map[string]string{
		"a/X.spec.md": "---\ncode: XXXXX\nlayer: infra\n---\n\n# X — x\n"})
	c, _ := New(root, g)
	c.Config = &config.Config{}
	c.InitScaffolds(false)
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, _ := readFile(c.Root, "arquitetura.md")
	if !strings.Contains(b, "Nenhum contêiner declarado") {
		t.Errorf("o documento não diz que falta declarar:\n%s", b)
	}
}
