package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// outsideAnyRepo gives a temporary folder with no repository above it (or skips).
func outsideAnyRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if Check(dir) == Disponível {
		t.Skipf("the temporary folder %s is inside a git repository", dir)
	}
	return dir
}

// The reason the diagnosis exists: `exit status 128` does not say whether to INSTALL or to
// INITIALISE, and the fixes are opposite. Each cause must come with its own.
func TestExplainNamesTheFixOfEachCause(t *testing.T) {
	t.Run("GTAVG-B04: The missing binary is explained with the install fix", func(t *testing.T) {})
	t.Run("GTAVG-B05: The missing repository is explained with the init fix", func(t *testing.T) {})
	noBin := Explain(SemBinário, "list the staged files")
	if !strings.Contains(noBin, "não está instalado") || !strings.Contains(noBin, "list the staged files") {
		t.Errorf("no binary must name the action and send to install: %s", noBin)
	}
	if strings.Contains(noBin, "git init") {
		t.Errorf("`git init` does not run without the binary — it sends to the wrong place: %s", noBin)
	}

	noRepo := Explain(SemRepo, "list the staged files")
	if !strings.Contains(noRepo, "git init") {
		t.Errorf("no repo must send to init: %s", noRepo)
	}
	if strings.Contains(noRepo, "instalado") {
		t.Errorf("git IS installed in this case: %s", noRepo)
	}
	// The command's action is in the sentence: it links the symptom to the cause in one line.
	if !strings.Contains(noRepo, "list the staged files") {
		t.Errorf("the message must name the action left incomplete: %s", noRepo)
	}
}

// With git available the failure is ANOTHER one — inventing a git cause would hide the real one.
func TestExplainIsSilentWhenGitIsAvailable(t *testing.T) {
	t.Run("GTAVG-B06: An available git is not explained", func(t *testing.T) {})
	if msg := Explain(Disponível, "anything"); msg != "" {
		t.Errorf("an available git must produce no explanation: %q", msg)
	}
}

func TestCheckTellsRepoFromNoRepo(t *testing.T) {
	t.Run("GTAVG-B03: A folder outside any repository is told apart from one inside", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := outsideAnyRepo(t)
	if d := Check(dir); d != SemRepo {
		t.Fatalf("folder without a repo = %v, want SemRepo", d)
	}

	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if d := Check(dir); d != Disponível {
		t.Fatalf("with .git = %v, want Disponível", d)
	}
}

// A subfolder of a repo is versioned: the `.git` is above, and complaining there would send
// the user to create a nested repo.
func TestCheckSeesARepoAbove(t *testing.T) {
	t.Run("GTAVG-B02: A repository above the root is seen", func(t *testing.T) {})
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := outsideAnyRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "pacotes", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if d := Check(sub); d != Disponível {
		t.Fatalf("subfolder of a repo = %v, want Disponível", d)
	}
}

func TestCheck_noGitOnThePath(t *testing.T) {
	t.Run("GTAVG-B01: Without git on the PATH the answer is no binary", func(t *testing.T) {})
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", t.TempDir()) // an empty folder: no git to find
	if d := Check(dir); d != SemBinário {
		t.Fatalf("without git on the PATH = %v, want SemBinário", d)
	}
}
