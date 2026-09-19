package issue

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// O `Reassign` tem de achar a âncora em issue ANTIGA — a que traz o rótulo em português.
//
// Por que este caso é justamente o que importa: o ramo de fallback do Reassign SÓ roda em
// issue anterior ao campo `owner`, e issue anterior ao campo é, por construção, issue
// escrita pelo binário antigo — com `- **alvo`. Quando os rótulos foram traduzidos, a
// âncora passou a ser `- **target`, que nenhuma issue antiga tem.
//
// E a falha é SILENCIOSA: o Replace não casa, a linha do dono não entra, o Reassign
// termina sem erro e a issue segue sem dono. Quem a reivindicar depois recebe o default,
// que é o agente — exatamente a pessoa que o reassign existia para tirar do caminho.
func TestReassignAchaAAncoraEmIssueComRotuloAntigo(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, Dir, string(Todo))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Uma issue como o binário ANTIGO a gravava: rótulo em português e SEM campo de dono.
	antiga := "# algo quebrou\n\n- **alvo (regido):** `src/a.ts`\n- **detectada em:** 2026-01-01\n"
	nome := "0001-algo.md"
	if err := os.WriteFile(filepath.Join(dir, nome), []byte(antiga), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Reassign(root, Todo, nome, DonoUsuário, "só uma pessoa decide isto"); err != nil {
		t.Fatalf("reassign falhou: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(dir, nome))
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "- **owner:**") {
		t.Fatalf("a linha do dono não foi inserida numa issue com rótulo antigo — o reassign "+
			"passou em silêncio e a issue ficou sem dono:\n%s", texto)
	}
	// E o dono lido de volta tem de ser o que se pediu: escrever a linha e não conseguir
	// relê-la seria o mesmo defeito um passo adiante.
	if got := Owner(strings.TrimSpace(string(ownerRE.FindSubmatch(b)[1]))); got != DonoUsuário {
		t.Errorf("dono relido = %q, queria %q", got, DonoUsuário)
	}
}
