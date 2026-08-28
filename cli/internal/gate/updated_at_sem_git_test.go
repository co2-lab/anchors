package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
)

// A pior mensagem que o silêncio do gitmeta produzia: sem REPOSITÓRIO, o gate caía no
// ramo "commitado", `LastCommitDate` falhava, e o veredito saía
// "arquivo sem commit no git (novo/untracked)" — uma afirmação falsa e ESPECÍFICA
// sobre o arquivo, que mandava o autor investigar exatamente onde o problema não está.
func TestUpdatedAtSemRepoNaoCulpaOArquivo(t *testing.T) {
	dir := t.TempDir()
	if gitmeta.Verifica(dir) == gitmeta.Disponível {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	conteudo := "// @anchors\n// updated_at: 2026-01-01\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}

	v, d := checkUpdatedAt(conteudo, mapx.Node{ID: "a.go"}, dir)

	if v != Skip {
		t.Fatalf("sem repositório não há veredito a dar sobre a data, foi %s (%s)", v, d)
	}
	if strings.Contains(d, "novo/untracked") {
		t.Errorf("culpa o ARQUIVO por uma falta que é do REPOSITÓRIO: %s", d)
	}
	if !strings.Contains(d, "repositório") {
		t.Errorf("a mensagem tem de nomear a causa real: %s", d)
	}
}

// A contrapartida: num repo de verdade, um arquivo novo e não commitado com a data de
// hoje continua passando — a regra do gate não mudou.
func TestUpdatedAtComRepoSegueConferindoAData(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("preparo: %s", out)
	}
	hoje := gitmeta.Today()
	conteudo := "// @anchors\n// updated_at: " + hoje + "\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}

	if v, d := checkUpdatedAt(conteudo, mapx.Node{ID: "a.go"}, dir); v != Pass {
		t.Fatalf("arquivo em edição com a data de hoje deveria passar, foi %s (%s)", v, d)
	}

	// E a data errada continua reprovando.
	errado := "// @anchors\n// updated_at: 2020-01-01\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(errado), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkUpdatedAt(errado, mapx.Node{ID: "a.go"}, dir); v != Fail {
		t.Fatalf("data errada em arquivo modificado deveria reprovar, foi %s", v)
	}
}
