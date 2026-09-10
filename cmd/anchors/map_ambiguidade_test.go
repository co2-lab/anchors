package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// A `scan.Ambiguities` existia — completa e testada — e NINGUÉM a chamava.
//
// O comentário do campo `priority` promete: "Declare quando a heurística errar; o `check`
// avisa onde ela decidiu sozinha." Não avisava.
//
// Custo medido no blue-eyes (blue-eyes#100): `**/*.test.*` e `packages/shared/**/*.ts`
// casavam o mesmo `AreaStatus.test.ts`, o desempate por comprimento escolheu `shared`, e
// o projeto ficou com ZERO nós `kind: test` tendo 70 testes verdes. Em cascata QUATRO
// gates ficaram cegos, incluindo o `test-traceable`, que é bloqueante.
//
// Só se descobriu contando os nós do mapa à mão.
func TestAmbiguidade_oParDeCamadasApareceNoAviso(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"test":   {Pattern: "**/*.test.ts", Kind: "test"},
		"shared": {Pattern: "packages/shared/**/*.ts", Kind: "code"},
	}}
	files := []scan.File{
		{Path: "packages/shared/AreaStatus.test.ts"},
		{Path: "packages/shared/QueryScope.test.ts"},
	}

	amb := scan.Ambiguities(files, cfg)
	if len(amb) != 2 {
		t.Fatalf("Ambiguities devolveu %d, queria 2 — os dois arquivos casam as duas camadas", len(amb))
	}
	// `packages/shared/**/*.ts` (23) é mais longo que `**/*.test.ts` (12), então o
	// desempate por comprimento entrega o TESTE à camada de código — o defeito do #100.
	if amb[0].Vencedora != "shared" {
		t.Errorf("vencedora = %q; o comprimento do pattern favorece `shared`", amb[0].Vencedora)
	}

	// declarar `priority` na camada certa CALA o aviso: o projeto já decidiu, e não há
	// adivinhação a denunciar.
	cfg.Layers["test"] = config.Layer{Pattern: "**/*.test.ts", Kind: "test", Priority: 100}
	if amb := scan.Ambiguities(files, cfg); len(amb) != 0 {
		t.Errorf("com `priority` declarada o aviso continuou: %+v", amb)
	}
}

// O aviso agrupa por PAR de camadas, não por arquivo: num monorepo o mesmo par produz
// centenas de linhas idênticas, e despejá-las esconde o padrão que precisa ser visto.
func TestAmbiguidade_avisoAgrupaPorParDeCamadas(t *testing.T) {
	amb := []scan.LayerAmbiguity{
		{Arquivo: "a.test.ts", Vencedora: "shared", Perdedoras: []string{"test"}},
		{Arquivo: "b.test.ts", Vencedora: "shared", Perdedoras: []string{"test"}},
		{Arquivo: "c.test.ts", Vencedora: "lambdas", Perdedoras: []string{"test"}},
	}
	saida := capturaStdout(t, func() { printLayerAmbiguities(amb) })

	if !strings.Contains(saida, "3 arquivo(s)") {
		t.Errorf("o total não aparece:\n%s", saida)
	}
	if !strings.Contains(saida, "shared venceu test") || !strings.Contains(saida, "(2 arquivo(s)") {
		t.Errorf("o par `shared venceu test` com contagem 2 não aparece:\n%s", saida)
	}
	if !strings.Contains(saida, "lambdas venceu test") {
		t.Errorf("o segundo par não aparece:\n%s", saida)
	}
	// a consequência tem de estar dita: sem ela o aviso parece cosmético
	if !strings.Contains(saida, "TODO gate") {
		t.Errorf("o aviso não diz o que se perde:\n%s", saida)
	}
}

// silêncio quando não há ambiguidade: um aviso que sempre aparece deixa de ser lido.
func TestAmbiguidade_semAmbiguidadeNaoImprimeNada(t *testing.T) {
	if s := capturaStdout(t, func() { printLayerAmbiguities(nil) }); s != "" {
		t.Errorf("imprimiu sem ambiguidade: %q", s)
	}
}

// capturaStdout coleta o que a função escreve em `os.Stdout`.
//
// Os comandos imprimem com `fmt.Println` direto (e não pelo `cmd.OutOrStdout()`), então
// testar a SAÍDA exige interceptar o descritor. É o que permite afirmar sobre o texto que
// o usuário lê, em vez de só sobre o valor de retorno.
func capturaStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	antigo := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = antigo }()

	f()
	w.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// A LIGAÇÃO, e não só a função.
//
// Este é o teste que faltava no defeito original: a `scan.Ambiguities` estava completa e
// testada, e o `map build` não a chamava. Um teste que exercita só a função passa com o
// mecanismo desligado — foi assim que ele ficou desligado por tanto tempo.
//
// Rodar o comando de verdade sobre um projeto ambíguo é o que confronta a ligação.
func TestMapBuild_avisaDaAmbiguidadeDeCamada(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "packages", "shared"), 0o755); err != nil {
		t.Fatal(err)
	}
	// as duas camadas casam o mesmo arquivo, e nenhuma declara `priority`
	yaml := `version: 1
layers:
  test:
    pattern: "**/*.test.ts"
    kind: test
  shared:
    pattern: "packages/shared/**/*.ts"
    kind: code
`
	if err := os.WriteFile(filepath.Join(root, config.DefaultFile), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	alvo := filepath.Join(root, "packages", "shared", "AreaStatus.test.ts")
	if err := os.WriteFile(alvo, []byte("export const x = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	saida := capturaStdout(t, func() {
		cmd := newMapCmd()
		cmd.SetArgs([]string{"build", "--root", root})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("map build: %v", err)
		}
	})

	if !strings.Contains(saida, "HEURÍSTICA") {
		t.Errorf("o `map build` NÃO avisou da ambiguidade — a `scan.Ambiguities` não está "+
			"ligada, que é exatamente o defeito de blue-eyes#100:\n%s", saida)
	}
	if !strings.Contains(saida, "shared venceu test") {
		t.Errorf("o par de camadas não aparece na saída:\n%s", saida)
	}
}
