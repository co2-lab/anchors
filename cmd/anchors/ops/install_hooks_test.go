// @anchors
//   code: IHTNS
//   ref: INHKN

package ops

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

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
	t.Run("INHKN-E01: A directory without anchors.yaml gets no hook", func(t *testing.T) {})
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
	t.Run("INHKN-E02: Outside a repository the error explains what needs git", func(t *testing.T) {})
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
	t.Run("INHKN-B02: A fresh install writes the three managed hooks, executable", func(t *testing.T) {})
	t.Run("INHKN-B11: Both merge drivers are registered in git config and .gitattributes", func(t *testing.T) {})
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
		// Windows has no executable bit; Git for Windows runs a hook without one.
		if fi, _ := os.Stat(p); runtime.GOOS != "windows" && fi.Mode()&0o111 == 0 {
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
	t.Run("INHKN-I01: Reinstalling never duplicates an attribute line nor damages the user's lines", func(t *testing.T) {})
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
	t.Run("INHKN-B03: Foreign hooks are respected unless forced", func(t *testing.T) {})
	t.Run("INHKN-X01: A hook the user wrote is never replaced without --force", func(t *testing.T) {})
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
	t.Run("INHKN-B01: The hooks go where git looks for them", func(t *testing.T) {})
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
	t.Run("INHKN-B04: The hooks anchors wrote, old or new, are updated on reinstall", func(t *testing.T) {})
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

// --- the installed scripts, run by a real git against a fake `anchors` ---

// fakeAnchors puts a fake `anchors` first on the PATH. `commit-msg` exits with
// $FAKE_MSG_STATUS, `verify` with $FAKE_VERIFY_STATUS, `--version` prints $FAKE_VERSION.
func fakeAnchors(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	script := `#!/bin/bash
case "$1" in
commit-msg) exit "${FAKE_MSG_STATUS:-0}" ;;
verify) exit "${FAKE_VERIFY_STATUS:-0}" ;;
--version) echo "anchors version ${FAKE_VERSION:-1.0.0}" ;;
esac
exit 0
`
	testkit.FakeBin(t, dir, "anchors", script)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// hookedRepo is a repository pushed to a bare remote whose main holds `remoteCfg`, with
// the hooks installed and a fake `anchors` on the PATH.
func hookedRepo(t *testing.T, remoteCfg string) (root string) {
	t.Helper()
	root, _ = frozenRepo(t, remoteCfg)
	fakeAnchors(t)
	if err, out := installHooks(t, root); err != nil {
		t.Fatalf("install-hooks: %v\n%s", err, out)
	}
	return root
}

// commitFile stages a new file and commits it through the hooks.
func commitFile(t *testing.T, root, name, subject string) (string, error) {
	t.Helper()
	writeFile(t, root, name, "x\n")
	if out, err := runGit(root, "add", name); err != nil {
		t.Fatal(out)
	}
	return runGit(root, "commit", "-q", "-m", subject)
}

// A frozen remote stops the commit and the push, and says why.
func TestInstalledHooksRefuseWorkOnAFrozenRemote(t *testing.T) {
	t.Run("INHKN-B05: The pre-commit refuses a commit while the remote is frozen", func(t *testing.T) {})
	t.Run("INHKN-B08: The pre-push refuses a push while the remote is frozen", func(t *testing.T) {})
	root := hookedRepo(t, "enabled: false\nfreeze_reason: \"rotate the key\"\nversion: 1\n")

	out, err := commitFile(t, root, "a.txt", "feat: a")
	if err == nil || !strings.Contains(out, "COMMIT REFUSED — the project is FROZEN") ||
		!strings.Contains(out, "rotate the key") {
		t.Fatalf("the pre-commit let a commit through on a frozen remote (%v):\n%s", err, out)
	}

	// Committed past the hooks, the push is still stopped.
	if out, err := runGit(root, "commit", "-q", "--no-verify", "-m", "feat: a"); err != nil {
		t.Fatal(out)
	}
	out, err = runGit(root, "push", "-q", "origin", "main")
	if err == nil || !strings.Contains(out, "PUSH REFUSED — the project is FROZEN") {
		t.Fatalf("the pre-push let a push through on a frozen remote (%v):\n%s", err, out)
	}
}

// The gates' verdict: "not governed" (exit 3) passes; a failure is deferred by the
// pre-commit to the commit-msg, which has the message in hand and blocks.
func TestInstalledHooksDeferTheVerdictToTheCommitMsg(t *testing.T) {
	t.Run("INHKN-B06: A staged set with nothing governed passes the pre-commit", func(t *testing.T) {})
	t.Run("INHKN-B07: A gate failure is deferred to the commit-msg, which blocks it", func(t *testing.T) {})
	t.Run("INHKN-B09: The commit-msg refuses a subject the message check refuses", func(t *testing.T) {})
	root := hookedRepo(t, "version: 1\n")

	t.Setenv("FAKE_VERIFY_STATUS", "3")
	out, err := commitFile(t, root, "package.json", "chore: deps")
	if err != nil {
		t.Fatalf("a commit with nothing governed was refused (%v):\n%s", err, out)
	}
	if !strings.Contains(out, "no staged file is governed by the Structure") || strings.Contains(out, "gates failed") {
		t.Errorf("exit 3 must read as \"nothing to confront\", not as a failure:\n%s", out)
	}

	t.Setenv("FAKE_VERIFY_STATUS", "1")
	out, err = commitFile(t, root, "b.txt", "feat: b")
	if err == nil {
		t.Fatalf("a failing gate let the commit through:\n%s", out)
	}
	if !strings.Contains(out, "gates failed. If it is deliberate") ||
		!strings.Contains(out, "commit BLOCKED by the anchors gates.\n  If the failure is deliberate") {
		t.Errorf("the pre-commit must defer and the commit-msg must block:\n%s", out)
	}

	t.Setenv("FAKE_VERIFY_STATUS", "0")
	t.Setenv("FAKE_MSG_STATUS", "1")
	if out, err := commitFile(t, root, "c.txt", "whatever"); err == nil {
		t.Fatalf("the commit-msg let a refused subject through:\n%s", out)
	}
}

// A remote that requires a newer binary stops the push.
func TestInstalledPrePushRefusesAnOlderBinary(t *testing.T) {
	t.Run("INHKN-B10: The pre-push refuses a binary older than the remote's minimum version", func(t *testing.T) {})
	root := hookedRepo(t, "version: 1\nmin_version: 2.0.0\n")
	t.Setenv("FAKE_VERSION", "1.9.0")
	if out, err := runGit(root, "commit", "-q", "--allow-empty", "-m", "feat: x"); err != nil {
		t.Fatal(out)
	}
	out, err := runGit(root, "push", "-q", "origin", "main")
	if err == nil || !strings.Contains(out, "requires 'anchors' 2.0.0 or newer") {
		t.Fatalf("an older binary pushed (%v):\n%s", err, out)
	}
	t.Setenv("FAKE_VERSION", "2.0.0")
	if out, err := runGit(root, "push", "-q", "origin", "main"); err != nil {
		t.Fatalf("the minimum version itself was refused (%v):\n%s", err, out)
	}
}

// The map's writer is `generated_by:` since the format moved to English; the pre-push
// grepped the old `gerado_por:` and the warning never fired. It reads both: a remote
// still on the old format has to be understood too.
func TestInstalledPrePushWarnsWhenTheMapWriterDiffers(t *testing.T) {
	t.Run("INHKN-B12: The pre-push warns when the remote map was written by another version", func(t *testing.T) {})
	for _, key := range []string{"generated_by", "gerado_por"} {
		t.Run(key, func(t *testing.T) {
			root := hookedRepo(t, "version: 1\n")
			t.Setenv("FAKE_VERSION", "1.0.0")
			writeFile(t, root, "anchors.graph.yaml", key+": 0.9.0\nnodes: {}\n")
			for _, args := range [][]string{
				{"add", "anchors.graph.yaml"},
				{"commit", "-q", "--no-verify", "-m", "chore: map"},
				{"push", "-q", "--no-verify", "origin", "main"},
				{"commit", "-q", "--no-verify", "--allow-empty", "-m", "feat: x"},
			} {
				if out, err := runGit(root, args...); err != nil {
					t.Fatalf("git %v: %s", args, out)
				}
			}
			out, err := runGit(root, "push", "-q", "origin", "main")
			if err != nil {
				t.Fatalf("a version divergence must warn, not refuse (%v):\n%s", err, out)
			}
			if !strings.Contains(out, "was written by 'anchors 0.9.0'") {
				t.Errorf("no warning for a map written by 0.9.0 (key %s):\n%s", key, out)
			}
			// The same version is silent.
			t.Setenv("FAKE_VERSION", "0.9.0")
			if out, _ := runGit(root, "commit", "-q", "--no-verify", "--allow-empty", "-m", "feat: y"); out != "" {
				t.Fatal(out)
			}
			if out, _ := runGit(root, "push", "-q", "origin", "main"); strings.Contains(out, "was written by") {
				t.Errorf("the same version warned:\n%s", out)
			}
		})
	}
}

// The hooks are bash, and on Windows they run through Git for Windows' own: nothing but a
// real commit shows they do, with the gates' `sh` behind them.
func TestHooksRunOnARealCommit(t *testing.T) {
	t.Run("INHKN-B13: The installed hooks run on a real commit, on every system", func(t *testing.T) {})
	if testing.Short() {
		t.Skip("builds the anchors binary")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	bin := t.TempDir()
	exe := filepath.Join(bin, "anchors")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", exe, "github.com/co2-lab/anchors/cmd/anchors").CombinedOutput(); err != nil {
		t.Fatalf("build anchors: %v %s", err, out)
	}
	testkit.OnPath(t, bin)
	root := t.TempDir()
	run := func(name string, args ...string) (string, error) {
		c := exec.Command(name, args...)
		c.Dir = root
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := c.CombinedOutput()
		return string(out), err
	}
	must := func(name string, args ...string) {
		t.Helper()
		if out, err := run(name, args...); err != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, err, out)
		}
	}
	gates := func(runLine string) string {
		return "version: 6\nlayers:\n  notes:\n    pattern: \"*.txt\"\n    kind: code\n" +
			"gates:\n  - name: must-pass\n    on: [code]\n    blocking: true\n    run: \"" + runLine + "\"\n"
	}
	must("git", "init", "-q")
	writeFile(t, root, "anchors.yaml", gates("exit 0"))
	writeFile(t, root, "a.txt", "one\n")
	must(exe, "map", "build")
	must("git", "add", ".")
	must("git", "commit", "-qm", "chore: base", "--no-verify")
	if err := runInstallHooks(root, false); err != nil {
		t.Fatal(err)
	}

	writeFile(t, root, "anchors.yaml", gates("echo refused-by-the-gate; exit 1"))
	writeFile(t, root, "a.txt", "two\n")
	must("git", "add", "anchors.yaml", "a.txt")
	if out, err := run("git", "commit", "-qm", "chore: refused"); err == nil || !strings.Contains(out, "must-pass") {
		t.Fatalf("a failing blocking gate refuses the commit through the hook, got %v:\n%s", err, out)
	}

	writeFile(t, root, "anchors.yaml", gates("exit 0"))
	must("git", "add", "anchors.yaml")
	if out, err := run("git", "commit", "-qm", "chore: lands"); err != nil {
		t.Fatalf("a passing gate lets the commit through, got %v:\n%s", err, out)
	}
	if out, _ := run("git", "log", "--oneline"); !strings.Contains(out, "chore: lands") {
		t.Errorf("the commit landed, log:\n%s", out)
	}
}
