package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func carregar(t *testing.T, yaml string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestWorkflowAusenteEhLocal(t *testing.T) {
	// Projeto que nunca declarou nada continua funcionando como sempre: a fila é local.
	// Exigir a declaração faria toda config existente virar erro de um dia para o outro.
	c, err := carregar(t, "version: 1\nlayers: {}\n")
	if err != nil {
		t.Fatalf("config sem bloco workflow tem de carregar: %v", err)
	}
	if c.ModoGitHub() {
		t.Error("sem declaração, o modo é local")
	}
}

func TestWorkflowGitHubExigeRepo(t *testing.T) {
	// Inferir do remote do git faria a escrita cair no repositório errado quando alguém
	// trabalha num fork — e escrita em lugar errado não se desfaz com revert.
	_, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  labels: [anchors]\n")
	if err == nil || !strings.Contains(err.Error(), "repo: owner/nome") {
		t.Fatalf("esperava erro exigindo repo, veio %v", err)
	}
}

func TestWorkflowGitHubExigeLabels(t *testing.T) {
	// Sem label, `anchors next` puxaria issue de produto, que não tem a forma do ciclo.
	_, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: acme/exemplo\n")
	if err == nil || !strings.Contains(err.Error(), "labels") {
		t.Fatalf("esperava erro exigindo labels, veio %v", err)
	}
}

func TestWorkflowGitHubValido(t *testing.T) {
	c, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: acme/exemplo\n  labels: [anchors]\n")
	if err != nil {
		t.Fatalf("config completa tem de carregar: %v", err)
	}
	if !c.ModoGitHub() {
		t.Error("mode: github tem de ligar o ModoGitHub()")
	}
}

func TestWorkflowRepoSemBarra(t *testing.T) {
	_, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: exemplo\n  labels: [anchors]\n")
	if err == nil || !strings.Contains(err.Error(), "owner/nome") {
		t.Fatalf("repo sem barra tem de falhar, veio %v", err)
	}
}

func TestWorkflowLocalNaoAceitaCamposDoGitHub(t *testing.T) {
	// Declarar repo/labels no modo local faz o ARQUIVO afirmar uma integração que não
	// existe — quem lê conclui que está ativa. É erro de declaração, não campo inofensivo.
	_, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: local\n  repo: acme/exemplo\n")
	if err == nil || !strings.Contains(err.Error(), "só valem em `mode: github`") {
		t.Fatalf("esperava erro de campo indevido, veio %v", err)
	}
}

func TestWorkflowModoDesconhecidoNaoCaiEmFallback(t *testing.T) {
	// A decisão de projeto: modo é DECLARADO, não adivinhado. Um modo com typo tem de
	// falhar, e não silenciosamente virar local — senão o agente trabalha achando que a
	// fila é do GitHub e reivindica no lugar errado.
	_, err := carregar(t, "version: 1\nlayers: {}\nworkflow:\n  mode: guithub\n")
	if err == nil {
		t.Fatal("modo desconhecido tem de falhar, não cair em fallback")
	}
	if !strings.Contains(err.Error(), "não há fallback") {
		t.Errorf("a mensagem tem de dizer que não há fallback, veio: %v", err)
	}
}

func TestChaveDesconhecidaFalhaEmVezDeSerIgnorada(t *testing.T) {
	// O atrito que motivou o KnownFields, medido ao pôr o Anchors no próprio Anchors:
	// escrevi `governs` com as chaves erradas (`guide`/`tags` em vez de `from`/`governs`),
	// o `map build` respondeu "222 nós, 0 arestas" e eu segui achando que o projeto não
	// tinha relações. A declaração inteira havia sido descartada em silêncio.
	_, err := carregar(t, "version: 1\nlayers: {}\ngovernz: []\n")
	if err == nil {
		t.Fatal("chave desconhecida tem de falhar — ignorá-la faz a declaração do usuário virar no-op invisível")
	}
	if !strings.Contains(err.Error(), "governz") {
		t.Errorf("o erro tem de NOMEAR a chave errada (senão não ajuda a consertar), veio: %v", err)
	}
}

func TestChaveErradaDentroDeGovernsTambemFalha(t *testing.T) {
	// O caso exato que me pegou: as chaves de topo estão certas, o erro está DENTRO do
	// item da lista. É o mais traiçoeiro, porque o arquivo parece certo à primeira vista.
	_, err := carregar(t, "version: 1\nlayers: {}\ngoverns:\n  - guide: X.md\n    tags: [y]\n")
	if err == nil {
		t.Fatal("chave errada dentro de governs[] tem de falhar")
	}
	if !strings.Contains(err.Error(), "guide") && !strings.Contains(err.Error(), "tags") {
		t.Errorf("o erro tem de nomear a chave, veio: %v", err)
	}
}

func TestConfigValidaContinuaCarregando(t *testing.T) {
	// Contraprova: o KnownFields não pode transformar config legítima em erro. Se este
	// teste quebrar ao adicionar campo novo, o campo faltou no struct — não relaxe o modo.
	c, err := carregar(t, `version: 1
layers:
  spec:
    pattern: "**/*.spec.md"
    kind: spec
    tags: [spec]
governs:
  - from: GUIDE.md
    governs: spec
gates:
  - name: spec-completa
    blocking: false
`)
	if err != nil {
		t.Fatalf("config válida tem de carregar: %v", err)
	}
	if len(c.Layers) != 1 || len(c.Governs) != 1 || len(c.Gates) != 1 {
		t.Errorf("config carregou incompleta: %d layers, %d governs, %d gates",
			len(c.Layers), len(c.Governs), len(c.Gates))
	}
}

func TestScopeInvalidoFalha(t *testing.T) {
	// O caso real: declarei `scope: repo` (o certo é `project`), o YAML carregou sem
	// reclamar e três gates entraram como "sem nada para medir". Enum errado é gate
	// DESLIGADO em silêncio — o arquivo afirma uma proteção que não existe.
	_, err := carregar(t, "version: 1\nlayers: {}\ngates:\n  - name: build\n    run: go build ./...\n    scope: repo\n")
	if err == nil {
		t.Fatal("scope inválido tem de falhar — cair no default por nó em silêncio desliga o gate")
	}
	if !strings.Contains(err.Error(), "repo") || !strings.Contains(err.Error(), "project") {
		t.Errorf("o erro tem de citar o valor errado E os válidos, veio: %v", err)
	}
}

func TestCostInvalidoFalha(t *testing.T) {
	// `cost` errado faz o gate ser tratado como rápido e entrar num loop apertado que ele
	// deveria evitar — o oposto do que o autor declarou.
	_, err := carregar(t, "version: 1\nlayers: {}\ngates:\n  - name: x\n    cost: lento\n")
	if err == nil || !strings.Contains(err.Error(), "lento") {
		t.Fatalf("cost inválido tem de falhar citando o valor, veio: %v", err)
	}
}

func TestWhenComFaseInvalidaFalha(t *testing.T) {
	// `when` com typo tira o gate da fase que o autor queria cobrar, sem aviso.
	_, err := carregar(t, "version: 1\nlayers: {}\ngates:\n  - name: x\n    when: [precommit]\n")
	if err == nil || !strings.Contains(err.Error(), "precommit") {
		t.Fatalf("fase inválida tem de falhar, veio: %v", err)
	}
}

func TestEnumsValidosCarregam(t *testing.T) {
	// Contraprova: os valores legítimos das três enumerações passam.
	c, err := carregar(t, `version: 1
layers: {}
gates:
  - name: a
    scope: project
    cost: slow
    when: [ci]
  - name: b
    scope: batch
    cost: fast
    when: [pre-commit, pre-push, manual]
  - name: c
    scope: node
`)
	if err != nil {
		t.Fatalf("enums válidos têm de carregar: %v", err)
	}
	if len(c.Gates) != 3 {
		t.Errorf("esperava 3 gates, veio %d", len(c.Gates))
	}
}
