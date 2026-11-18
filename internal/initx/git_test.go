// @anchors
//   code: GTTSB
//   ref: GTSTG

package initx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
)

// isolatedRoot returns a folder that is NOT inside any git repository. Without it the test
// would be a false positive: t.TempDir() on macOS lands in /var/folders, but running the
// suite from inside the Anchors repository with a relative path would find the `.git`
// above and every state would become GitPronto.
func isolatedRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if _, inside := repoAcima(dir); inside {
		t.Skipf("the temporary folder %s is inside a git repository", dir)
	}
	return dir
}

// The distinction behind the four states: "no git" on the machine is NOT the same as
// "no git in this project", because Anchors acts differently — without the program there
// is no `git init` to offer.
func TestGitNotInstalledIsNotTheSameAsNotInitialised(t *testing.T) {
	t.Run("GTSTG-B01: With git not installed the state is not installed", func(t *testing.T) {})
	t.Run("GTSTG-B02: With git installed and no repository anywhere the state is not initialised", func(t *testing.T) {})
	t.Run("GTSTG-I01: Not installed and not initialised stay two states with different offers", func(t *testing.T) {})
	t.Run("GTSTG-X02: Whether git is installed comes from the caller", func(t *testing.T) {})
	dir := isolatedRoot(t)

	withoutProgram := DetectGit(dir, false)
	if withoutProgram != GitNaoInstalado {
		t.Fatalf("without the program the state must be GitNaoInstalado, got %v", withoutProgram)
	}
	withProgram := DetectGit(dir, true)
	if withProgram != GitNaoIniciado {
		t.Fatalf("with the program and no repository the state must be GitNaoIniciado, got %v", withProgram)
	}

	// And the action differs: only one of the two has something to offer.
	if OfferAction(GitNaoInstalado) {
		t.Error("there is no `git init` to offer on a machine without git — asking would promise what would fail")
	}
	if !OfferAction(GitNaoIniciado) {
		t.Error("with git installed and no repository, offering to initialise is exactly the step")
	}
}

// Each state teaches something different. A message telling someone who already has git
// to "install" it sends them looking for the problem where it is not.
func TestGitWarningSaysWhatToDoInEachState(t *testing.T) {
	t.Run("GTSTG-B09: Each unready state has its own warning and ready has none", func(t *testing.T) {})
	// In the project's language: the warnings were hard-coded in Portuguese.
	t.Cleanup(func() { _ = i18n.Set(i18n.Default) })
	for lang, cases := range map[string][]struct {
		state    GitState
		contains string
	}{
		"en":    {{GitNaoInstalado, "git is not installed"}, {GitNaoIniciado, "not under git"}, {GitSemCommit, "no commit yet"}},
		"pt-BR": {{GitNaoInstalado, "não está instalado"}, {GitNaoIniciado, "não está sob git"}, {GitSemCommit, "nenhum commit"}},
	} {
		if err := i18n.Set(lang); err != nil {
			t.Fatal(err)
		}
		for _, c := range cases {
			w := AvisoGit(c.state)
			if !strings.Contains(w, c.contains) {
				t.Errorf("%s: the warning for %v should contain %q, got: %s", lang, c.state, c.contains, w)
			}
		}
	}
	if AvisoGit(GitPronto) != "" {
		t.Error("a ready repository has nothing to warn about")
	}
}

// Offering is decided per state, exhaustively over the four.
func TestGitOfferOnlyWhereThereIsAStepToTake(t *testing.T) {
	t.Run("GTSTG-B08: Only the not-initialised and no-commit states have an action to offer", func(t *testing.T) {})
	want := map[GitState]bool{GitNaoInstalado: false, GitNaoIniciado: true, GitSemCommit: true, GitPronto: false}
	for state, offer := range want {
		if got := OfferAction(state); got != offer {
			t.Errorf("OfferAction(%v) = %v, want %v", state, got, offer)
		}
	}
}

// A repository created and never committed is a state of its OWN: calling it "has git"
// would claim `git log`/`git diff` work, and they do not without HEAD.
func TestGitRepositoryWithoutCommitIsNotReady(t *testing.T) {
	t.Run("GTSTG-B05: A repository with no reference at all has no commit yet", func(t *testing.T) {})
	dir := isolatedRoot(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}

	if e := DetectGit(dir, true); e != GitSemCommit {
		t.Fatalf("a repository with no reference must be GitSemCommit, got %v", e)
	}
	// An EMPTY packed list is no history either.
	if err := os.WriteFile(filepath.Join(dir, ".git", "packed-refs"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if e := DetectGit(dir, true); e != GitSemCommit {
		t.Fatalf("an empty packed-refs is no commit, got %v", e)
	}
	if !OfferAction(GitSemCommit) {
		t.Error("the commit is what is missing — and that is what Anchors has to offer here")
	}
}

func TestGitRepositoryWithCommitIsReady(t *testing.T) {
	t.Run("GTSTG-B03: A repository with a branch reference is ready", func(t *testing.T) {})
	dir := isolatedRoot(t)
	heads := filepath.Join(dir, ".git", "refs", "heads")
	if err := os.MkdirAll(heads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(heads, "main"), []byte("abc123\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectGit(dir, true); e != GitPronto {
		t.Fatalf("a repository with a reference in heads/ is ready, got %v", e)
	}
	if OfferAction(GitPronto) {
		t.Error("nothing to offer on a ready repository")
	}
}

// `git gc` packs the references: heads/ becomes empty and the history lives in
// packed-refs. Reading heads/ alone would classify a repository with years of history as
// "no commit".
func TestGitRepositoryWithPackedRefsIsReady(t *testing.T) {
	t.Run("GTSTG-B04: A repository with packed references is ready", func(t *testing.T) {})
	dir := isolatedRoot(t)
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}
	packed := "# pack-refs with: peeled fully-peeled sorted \nabc123 refs/heads/main\n"
	if err := os.WriteFile(filepath.Join(dir, ".git", "packed-refs"), []byte(packed), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectGit(dir, true); e != GitPronto {
		t.Fatalf("packed references count as a commit, got %v", e)
	}
}

// Running `init` in a subfolder of an existing repository must NOT offer a nested
// repository: the subproject's history would vanish from the repository above, and that
// is not undone with a revert.
func TestGitSubfolderOfExistingRepositoryOffersNoNestedRepository(t *testing.T) {
	t.Run("GTSTG-B06: A subfolder of an existing repository is ready", func(t *testing.T) {})
	dir := isolatedRoot(t)
	heads := filepath.Join(dir, ".git", "refs", "heads")
	if err := os.MkdirAll(heads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(heads, "main"), []byte("abc123\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "pacotes", "app")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if e := DetectGit(sub, true); e != GitPronto {
		t.Fatalf("a subfolder of an existing repository is already versioned, got %v", e)
	}
}

// Worktrees and submodules have `.git` as a FILE (a `gitdir: …` pointer). Treating it as
// "no repository" would offer to initialise over a valid worktree.
func TestGitMarkerAsFileIsReady(t *testing.T) {
	t.Run("GTSTG-B07: A repository marker that is a file is ready", func(t *testing.T) {})
	dir := isolatedRoot(t)
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /outro/lugar/.git/worktrees/x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if e := DetectGit(dir, true); e != GitPronto {
		t.Fatalf("`.git` as a file is a worktree or submodule — the real repository exists, got %v", e)
	}
}

// The seeded .gitignore covers what Anchors ITSELF generates. It does not guess the
// stack: at this point of init the language has not been read yet.
func TestGitignoreCoversWhatAnchorsGenerates(t *testing.T) {
	t.Run("GTSTG-B10: The seeded ignore list covers what Anchors generates", func(t *testing.T) {})
	t.Run("GTSTG-X01: The seeded ignore list does not guess the stack", func(t *testing.T) {})
	for _, want := range []string{".DS_Store", ".anchors/"} {
		if !strings.Contains(GitignorePadrão, want) {
			t.Errorf("the default .gitignore should ignore %q", want)
		}
	}
	// NO EXCEPTION inside `.anchors/`. A negation there fails silently: `.anchors/`
	// excludes the FOLDER, git never descends into it, and the file "saved" by the
	// negation vanishes with no error. What must be versioned is born OUTSIDE — as was
	// done with sbom.json, which is now generated at the root.
	if strings.Contains(GitignorePadrão, "!.anchors/") {
		t.Error("an exception inside `.anchors/` does not work (the folder is excluded first) " +
			"and vanishes silently: the versionable artifact must be born outside the working area")
	}
	for _, stack := range []string{"node_modules", "target/", "vendor/"} {
		if strings.Contains(GitignorePadrão, stack) {
			t.Errorf(".gitignore must not guess the stack (%q): the language has not been read yet", stack)
		}
	}
}

// O QUE DESCREVE O PROJETO NÃO NASCE EM `.anchors/`.
//
// A distinção não é "quem gerou" — é o que o arquivo DESCREVE. O `.anchors/` guarda
// execução (fila de julgamento, espelho do check, cache) e está inteiro no `.gitignore`.
// O que descreve o PROJETO tem de nascer fora, ou é engolido.
//
// Aconteceu com o SBOM: o gate dizia `Measures: "... é gerado e versionável"` enquanto o
// `Run` o escrevia em `.anchors/sbom.json` — a contradição estava a duas linhas de
// distância e sobreviveu ao commit inicial. A saída seria uma exceção no `.gitignore`, e
// exceção ali é onde o silêncio mora: `.anchors/` com barra exclui o DIRETÓRIO, o git nem
// desce nele para avaliar a negação, e o arquivo some sem erro.
func TestArtefatoVersionavelNaoNasceNaAreaDeTrabalho(t *testing.T) {
	for _, g := range DefaultGates(map[string]bool{"code": true, "spec": true}, true) {
		if g.Run == "" || !strings.Contains(g.Run, ".anchors/") {
			continue
		}
		// Um gate que grava em `.anchors/` e se diz versionável está se contradizendo:
		// aquilo nunca vai chegar ao repositório.
		if strings.Contains(g.Measures, "versionável") || strings.Contains(g.Measures, "versionavel") {
			t.Errorf("gate %q grava em `.anchors/` (ignorado) e promete algo VERSIONÁVEL — "+
				"o arquivo nunca chegaria ao repositório. Grave na raiz.\n  run: %s\n  measures: %s",
				g.ID, g.Run, g.Measures)
		}
	}
}
