package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mapaComCodigos escreve um mapa mínimo com nós de identidade em workspaces diferentes.
func mapaComCodigos(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.graph.yaml")
	const y = `version: 1
nodes:
    - id: apps/mobile/src/features/auth/LoginScreen.spec.md
      kind: spec
      code: LOGI
    - id: apps/mobile/src/features/auth/LoginScreen.tsx
      kind: code
      code: LOGI
    - id: packages/backend/services/security.spec.md
      kind: spec
      code: SGSB
    - id: apps/mobile/src/components/atoms/Button.spec.md
      kind: spec
      code: BTTN
edges: []
`
	if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func rodarList(t *testing.T, mapPath string, args ...string) string {
	t.Helper()
	cmd := newCodeListCmd()
	var out strings.Builder
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(append([]string{"--map", mapPath}, args...))
	// O comando imprime em os.Stdout; captura por pipe.
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := cmd.Execute()
	w.Close()
	os.Stdout = orig
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	return string(buf[:n])
}

func TestCodeListEnumeraDoMapa(t *testing.T) {
	// A razão de o comando existir: a alternativa é grep por padrão de código, que casa
	// menção em prosa, comentário e nome de arquivo — e depende do comprimento do código,
	// que varia por projeto (`code_lengths`). Aqui a fonte é o campo estruturado do mapa.
	out := rodarList(t, mapaComCodigos(t))
	for _, c := range []string{"LOGI", "SGSB", "BTTN"} {
		if !strings.Contains(out, c) {
			t.Errorf("o código %s devia aparecer, veio:\n%s", c, out)
		}
	}
	// Ordenado: quem lê procura um código, e lista fora de ordem obriga a procurar duas vezes.
	if i, j := strings.Index(out, "BTTN"), strings.Index(out, "LOGI"); i > j {
		t.Errorf("a saída deve estar ordenada por código, veio:\n%s", out)
	}
}

func TestCodeListUmaLinhaPorCodigoMesmoComVariosArquivos(t *testing.T) {
	// LOGI está na spec E no .tsx — é UMA unidade, não duas. Repetir a linha faria a
	// contagem mentir sobre quantas identidades existem.
	out := rodarList(t, mapaComCodigos(t))
	if n := strings.Count(out, "LOGI"); n != 1 {
		t.Errorf("LOGI devia aparecer em UMA linha (a pasta é a mesma), apareceu %d× em:\n%s", n, out)
	}
}

func TestCodeListFiltraPorWorkspace(t *testing.T) {
	// A pergunta real de quem trabalha num monorepo: "os códigos DESTE workspace".
	out := rodarList(t, mapaComCodigos(t), "--in", "packages/backend")
	if !strings.Contains(out, "SGSB") {
		t.Errorf("SGSB está sob packages/backend e devia aparecer, veio:\n%s", out)
	}
	for _, fora := range []string{"LOGI", "BTTN"} {
		if strings.Contains(out, fora) {
			t.Errorf("%s está fora do filtro e não devia aparecer, veio:\n%s", fora, out)
		}
	}
}

func TestCodeListFiltroSemResultadoNaoMenteVazio(t *testing.T) {
	// "nenhum código sob X" é diferente de "o projeto não tem código" — a segunda mensagem
	// mandaria rodar `map build` sem necessidade.
	out := rodarList(t, mapaComCodigos(t), "--in", "apps/web")
	if !strings.Contains(out, "apps/web") {
		t.Errorf("a mensagem deve citar o filtro que não casou, veio:\n%s", out)
	}
}

// O título do artefato serve para NOMEAR o trabalho num card, e por isso a parte que
// repete o tipo e o número ("Plano 0001 — ") sai: quem consome já tem o kind e o código,
// e a repetição só faz o texto crescer.
func TestTituloDoArquivoTiraOPrefixoRedundante(t *testing.T) {
	dir := t.TempDir()
	casos := map[string]string{
		"# Login\n":                       "Login",
		"# Plano 0001 — Fundação\n":       "Fundação",
		"# Spec 0042 — Recuperação\n":     "Recuperação",
		"sem titulo\n":                    "",
		"<!-- @anchors -->\n\n# Depois\n": "Depois",
	}
	for conteudo, esperado := range casos {
		if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := tituloDoArquivo(dir, "a.md"); got != esperado {
			t.Errorf("%q → %q, queria %q", conteudo, got, esperado)
		}
	}
	// Arquivo que não é markdown não tem título a extrair.
	if got := tituloDoArquivo(dir, "x.go"); got != "" {
		t.Errorf("arquivo não-markdown devolveu %q", got)
	}
}
