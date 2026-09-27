package mapcmd

import (
	"io"
	"os"
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
