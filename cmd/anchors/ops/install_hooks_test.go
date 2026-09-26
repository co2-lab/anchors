package ops

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O HOOK PRECISA FUNCIONAR EM WORKTREE, e `$ROOT/.git` não serve.
//
// Num worktree, `.git` é um ARQUIVO que aponta para `…/.git/worktrees/<nome>` — escrever
// dentro dele falha com `Not a directory`.
//
// MEDIDO: um agente trabalhando em worktree levou
//
//	.git/hooks/pre-commit: line 48: …/be-rev1/.git/anchors-freeze-cache: Not a directory
//
// e commitou com `--no-verify`. O hook INTEIRO deixou de rodar por causa do cache — e um
// hook que falha assim é pior que hook nenhum: ele treina quem o usa a contorná-lo.
//
// `git rev-parse --git-dir` devolve o caminho certo nos dois casos.
func TestHookAchaOGitDirEmWorktree(t *testing.T) {
	if strings.Contains(preCommitScript, `"$ROOT/.git/`) {
		t.Error("o hook monta um caminho dentro de `$ROOT/.git` — num worktree isso é um " +
			"ARQUIVO, e a escrita falha com `Not a directory`, derrubando o hook inteiro")
	}
	if !strings.Contains(preCommitScript, "rev-parse --git-dir") {
		t.Error("o hook não usa `git rev-parse --git-dir` — é o único jeito de achar o " +
			"diretório do git que funciona em clone comum E em worktree")
	}
}

// installHooks runs `anchors install-hooks` over root and returns the error and output.
func installHooks(t *testing.T, root string, flags ...string) (error, string) {
	t.Helper()
	cmd := newInstallHooksCmd()
	cmd.SetArgs(append([]string{"--root", root}, flags...))
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return err, out
}

func hookProject(t *testing.T) string {
	t.Helper()
	root := newGitRepo(t)
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	return root
}

// Without anchors.yaml there is nothing to govern: the command refuses and says why.
func TestInstallHooksRequiresAnAnchorsProject(t *testing.T) {
	root := newGitRepo(t)
	err, _ := installHooks(t, root)
	if err == nil || !strings.Contains(err.Error(), "anchors init") {
		t.Fatalf("want a refusal pointing at `anchors init`, got %v", err)
	}
	if _, serr := os.Stat(filepath.Join(root, ".git", "hooks", "pre-commit")); serr == nil {
		t.Error("a hook was installed in a project without anchors.yaml")
	}
}

// Outside a repository the hook has nowhere to live, and the error says that — not a raw
// git message about rev-parse.
func TestInstallHooksOutsideGitExplainsWhy(t *testing.T) {
	isolateGit(t)
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	err, _ := installHooks(t, root)
	if err == nil || !strings.Contains(err.Error(), "install the pre-commit") {
		t.Fatalf("want an explanation naming what needs git, got %v", err)
	}
}

// A fresh install writes the three hooks (executable, with the marker) and registers both
// merge drivers: the attribute line and the git config that says how to call them.
func TestInstallHooksWritesTheHooksAndTheMergeDrivers(t *testing.T) {
	root := hookProject(t)
	err, out := installHooks(t, root)
	if err != nil {
		t.Fatalf("install-hooks: %v\n%s", err, out)
	}
	hooks := filepath.Join(root, ".git", "hooks")
	for name, want := range map[string]string{
		"pre-commit": preCommitScript, "commit-msg": commitMsgScript, "pre-push": prePushScript,
	} {
		p := filepath.Join(hooks, name)
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			t.Errorf("%s not written: %v", name, rerr)
			continue
		}
		if string(b) != want {
			t.Errorf("%s does not hold the managed script", name)
		}
		if fi, _ := os.Stat(p); fi.Mode()&0o111 == 0 {
			t.Errorf("%s is not executable (%v)", name, fi.Mode())
		}
	}
	attrs, _ := os.ReadFile(filepath.Join(root, ".gitattributes"))
	for _, line := range []string{"anchors.graph.yaml merge=anchors-map", "*-progress.md merge=anchors-progress"} {
		if strings.Count(string(attrs), line) != 1 {
			t.Errorf(".gitattributes should hold %q once:\n%s", line, attrs)
		}
	}
	for key, want := range map[string]string{
		"merge.anchors-map.driver":      "anchors map merge %O %A %B",
		"merge.anchors-progress.driver": "anchors merge-progress %O %A %B",
	} {
		if got, _ := runGit(root, "config", "--get", key); got != want {
			t.Errorf("git config %s = %q, want %q", key, got, want)
		}
	}
	if !strings.Contains(out, "map merge driver installed") {
		t.Errorf("the output does not report the driver:\n%s", out)
	}
}

// Reinstalling is safe: the managed hooks are replaced, and the attribute lines are not
// appended a second time.
func TestInstallHooksIsIdempotent(t *testing.T) {
	root := hookProject(t)
	writeFile(t, root, ".gitattributes", "*.png binary") // no trailing newline
	if err, out := installHooks(t, root); err != nil {
		t.Fatalf("first install: %v\n%s", err, out)
	}
	err, out := installHooks(t, root)
	if err != nil {
		t.Fatalf("second install: %v\n%s", err, out)
	}
	attrs, _ := os.ReadFile(filepath.Join(root, ".gitattributes"))
	if !strings.HasPrefix(string(attrs), "*.png binary\n") {
		t.Errorf("the user's line was damaged:\n%s", attrs)
	}
	if n := strings.Count(string(attrs), "merge=anchors-map"); n != 1 {
		t.Errorf("the map driver line appears %d times after reinstalling:\n%s", n, attrs)
	}
	if !strings.Contains(out, "already configured") {
		t.Errorf("the second run should say the drivers were already there:\n%s", out)
	}
}

// A hook the user wrote is theirs. A foreign pre-commit refuses the whole install without
// --force; a foreign commit-msg or pre-push is left alone with a warning. --force replaces
// them.
func TestInstallHooksRespectsForeignHooksUnlessForced(t *testing.T) {
	root := hookProject(t)
	hooks := filepath.Join(root, ".git", "hooks")
	const mine = "#!/bin/sh\necho mine\n"
	for _, h := range []string{"pre-commit", "commit-msg", "pre-push"} {
		writeFile(t, hooks, h, mine)
	}

	err, _ := installHooks(t, root)
	if err == nil || !strings.Contains(err.Error(), "--force") {
		t.Fatalf("a foreign pre-commit must be refused with a pointer to --force, got %v", err)
	}
	for _, h := range []string{"pre-commit", "commit-msg", "pre-push"} {
		if b, _ := os.ReadFile(filepath.Join(hooks, h)); string(b) != mine {
			t.Errorf("the refused install still overwrote the user's %s", h)
		}
	}

	// Without the foreign pre-commit the install proceeds, around the other two.
	if err := os.Remove(filepath.Join(hooks, "pre-commit")); err != nil {
		t.Fatal(err)
	}
	err, out := installHooks(t, root)
	if err != nil {
		t.Fatalf("install-hooks: %v\n%s", err, out)
	}
	for _, h := range []string{"commit-msg", "pre-push"} {
		if b, _ := os.ReadFile(filepath.Join(hooks, h)); string(b) != mine {
			t.Errorf("the user's %s was overwritten without --force", h)
		}
		if !strings.Contains(out, h+" exists and was not written by anchors") {
			t.Errorf("no warning about the foreign %s:\n%s", h, out)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(hooks, "pre-commit")); string(b) != preCommitScript {
		t.Error("the pre-commit was not installed next to the foreign hooks")
	}

	if err, out := installHooks(t, root, "--force"); err != nil {
		t.Fatalf("--force: %v\n%s", err, out)
	}
	for h, want := range map[string]string{"commit-msg": commitMsgScript, "pre-push": prePushScript} {
		if b, _ := os.ReadFile(filepath.Join(hooks, h)); string(b) != want {
			t.Errorf("--force did not replace %s", h)
		}
	}
}

// core.hooksPath wins over .git/hooks — relative to the root, or absolute.
func TestGitHooksDirHonorsCoreHooksPath(t *testing.T) {
	root := newGitRepo(t)
	def, err := gitHooksDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(def) != "hooks" || filepath.Base(filepath.Dir(def)) != ".git" {
		t.Errorf("default hooks dir = %s, want <root>/.git/hooks", def)
	}
	if out, err := runGit(root, "config", "core.hooksPath", ".githooks"); err != nil {
		t.Fatal(out)
	}
	if got, _ := gitHooksDir(root); got != filepath.Join(root, ".githooks") {
		t.Errorf("relative core.hooksPath: got %s", got)
	}
	abs := t.TempDir()
	if out, err := runGit(root, "config", "core.hooksPath", abs); err != nil {
		t.Fatal(out)
	}
	if got, _ := gitHooksDir(root); got != abs {
		t.Errorf("absolute core.hooksPath: got %s, want %s", got, abs)
	}
}

// Every script anchors writes carries the marker it checks on reinstall. The commit-msg
// did not (its header was `# anchors:hook — instalado por…`), so a reinstall without
// --force took the hook anchors itself wrote for the user's, and never updated it. A
// commit-msg with that old header is anchors' own too, and is replaced.
func TestInstallHooksUpdatesItsOwnCommitMsg(t *testing.T) {
	for name, script := range map[string]string{
		"pre-commit": preCommitScript, "commit-msg": commitMsgScript, "pre-push": prePushScript,
	} {
		if !strings.Contains(script, hookMarker) {
			t.Errorf("the %s script does not carry %q: a reinstall would take it for the user's", name, hookMarker)
		}
	}

	root := hookProject(t)
	hooks := filepath.Join(root, ".git", "hooks")
	const oldCommitMsg = "#!/usr/bin/env bash\n# anchors:hook — instalado por 'anchors install-hooks'\n" +
		"set -euo pipefail\nanchors commit-msg \"$1\"\n"
	writeFile(t, hooks, "commit-msg", oldCommitMsg)

	err, out := installHooks(t, root)
	if err != nil {
		t.Fatalf("install-hooks: %v\n%s", err, out)
	}
	if b, _ := os.ReadFile(filepath.Join(hooks, "commit-msg")); string(b) != commitMsgScript {
		t.Errorf("the commit-msg an older anchors wrote was not updated:\n%s", b)
	}
	if strings.Contains(out, "commit-msg exists and was not written by anchors") {
		t.Errorf("anchors' own commit-msg was reported as foreign:\n%s", out)
	}

	// And the next reinstall, over the hook just written, updates it again.
	writeFile(t, hooks, "commit-msg", strings.Replace(commitMsgScript, "set -euo pipefail", "set -eu", 1))
	if err, out := installHooks(t, root); err != nil {
		t.Fatalf("second install: %v\n%s", err, out)
	}
	if b, _ := os.ReadFile(filepath.Join(hooks, "commit-msg")); string(b) != commitMsgScript {
		t.Errorf("the commit-msg anchors wrote was not updated on reinstall:\n%s", b)
	}
}
