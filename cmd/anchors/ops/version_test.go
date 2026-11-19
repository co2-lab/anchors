// @anchors
//   code: VRSNK
//   ref: VRSNI

package ops

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/shell"
)

const appVersioning = `version: 8
versioning:
  app:
    scheme: semver
    source:
      - file: app/version.json
        path: version
    tag: "app-v{version}"
    next: semantic
    branches:
      "hotfix/*":
        next: { strategy: increment, part: patch }
      "feat/*":
        next: none
`

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(shell.WithoutRepoEnv(os.Environ()), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func versionRun(t *testing.T, args ...string) (string, error) {
	t.Helper()
	c := newVersionCmd()
	var out, errb bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&errb)
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestVersion_shows(t *testing.T) {
	root := historyRepo(t, "file:anchors.yaml="+appVersioning, "file:app/version.json={\"version\": \"1.2.0\"}\n",
		"feat: zero", "tag:app-v1.1.0", "feat: one", "feat: two")
	t.Run("VRSNI-B01: The command shows the sources, the release, the branch's rule and the next version", func(t *testing.T) {
		out, err := versionRun(t, "--root", root, "--branch", "develop")
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"app (semver)", "1.2.0      app/version.json", "1.1.0      app-v1.1.0", "develop → semantic (default)", "next       1.2.0      semantic: 2 feat since 1.1.0", "2 since app-v1.1.0"} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in\n%s", want, out)
			}
		}
		out, err = versionRun(t, "--root", root, "--branch", "hotfix/login")
		if err != nil || !strings.Contains(out, "hotfix/login → increment (hotfix/*)") {
			t.Errorf("hotfix: %v\n%s", err, out)
		}
	})
	t.Run("VRSNI-B07: A single product needs no name, and the current branch decides", func(t *testing.T) {
		gitIn(t, root, "checkout", "-q", "-b", "hotfix/login")
		out, err := versionRun(t, "--root", root)
		if err != nil || !strings.Contains(out, "app (semver)") || !strings.Contains(out, "hotfix/login → increment (hotfix/*)") {
			t.Errorf("%v\n%s", err, out)
		}
	})
	t.Run("VRSNI-E02: A detached HEAD needs a branch named", func(t *testing.T) {
		gitIn(t, root, "checkout", "-q", "--detach")
		if _, err := versionRun(t, "--root", root); err == nil || !strings.Contains(err.Error(), "--branch") {
			t.Errorf("got %v", err)
		}
	})
}

func TestVersion_release(t *testing.T) {
	t.Run("VRSNI-B02: The latest release is the highest version of the product, not the last name", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+appVersioning, "file:app/version.json={\"version\": \"1.10.0\"}",
			"feat: a", "tag:app-v1.9.0", "tag:app-v1.10.0", "tag:api-v9.0.0")
		out, err := versionRun(t, "--root", root, "--branch", "main")
		if err != nil || !strings.Contains(out, "released   1.10.0     app-v1.10.0") {
			t.Errorf("%v\n%s", err, out)
		}
	})
	t.Run("VRSNI-B03: JSON prints the state with the commits classified", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+appVersioning, "file:app/version.json={\"version\": \"2.0.0\"}",
			"chore: start", "tag:app-v1.1.0", "feat!: a new store", "fix: a crash\n\nBug: seen in production")
		out, err := versionRun(t, "--root", root, "--branch", "main", "--json")
		if err != nil {
			t.Fatal(err)
		}
		var st struct {
			Product, Scheme, Current, Released, Next, Reason string
			Commits                                          []struct {
				Subject, Type string
				Breaking      bool
				Footers       map[string]string
			}
		}
		if err := json.Unmarshal([]byte(out), &st); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
		if st.Product != "app" || st.Current != "2.0.0" || st.Released != "1.1.0" || st.Next != "2.0.0" || len(st.Commits) != 2 {
			t.Fatalf("got %+v", st)
		}
		if c := st.Commits[0]; c.Type != "feat" || !c.Breaking {
			t.Errorf("breaking feat: %+v", c)
		}
		if c := st.Commits[1]; c.Type != "fix" || c.Breaking || c.Footers["Bug"] != "seen in production" {
			t.Errorf("fix: %+v", c)
		}
	})
}

const twoSources = `version: 8
versioning:
  app:
    scheme: semver
    source:
      - file: app/version.json
        path: version
      - file: app/build.txt
        regex: 'VERSION=(\S+)'
    tag: "app-v{version}"
    next: semantic
`

func TestVersion_bump(t *testing.T) {
	t.Run("VRSNI-B04: Bump writes the strategy's version into every source, and only the version", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+twoSources, "file:app/version.json={\n  \"name\": \"x\",\n  \"version\": \"1.1.0\"\n}\n",
			"file:app/build.txt=NAME=x\nVERSION=1.1.0\n", "chore: start", "tag:app-v1.1.0", "feat: one")
		out, err := versionRun(t, "--root", root, "--branch", "main", "bump")
		if err != nil {
			t.Fatal(err)
		}
		if got := readFile(t, root, "app/version.json"); got != "{\n  \"name\": \"x\",\n  \"version\": \"1.2.0\"\n}\n" {
			t.Errorf("json: %q", got)
		}
		if got := readFile(t, root, "app/build.txt"); got != "NAME=x\nVERSION=1.2.0\n" {
			t.Errorf("txt: %q", got)
		}
		if !strings.Contains(out, "app 1.1.0 → 1.2.0") || !strings.Contains(out, "app/version.json") || !strings.Contains(out, "app/build.txt") {
			t.Errorf("output: %s", out)
		}
	})
	t.Run("VRSNI-B05: Bump raises a part, or writes the version given", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+appVersioning, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start", "tag:app-v1.1.0")
		if _, err := versionRun(t, "--root", root, "--branch", "main", "bump", "major"); err != nil || readFile(t, root, "app/version.json") != `{"version": "2.0.0"}` {
			t.Errorf("major: %v %s", err, readFile(t, root, "app/version.json"))
		}
		if _, err := versionRun(t, "--root", root, "--branch", "main", "bump", "3.0.1"); err != nil || readFile(t, root, "app/version.json") != `{"version": "3.0.1"}` {
			t.Errorf("exact: %v %s", err, readFile(t, root, "app/version.json"))
		}
	})
	t.Run("VRSNI-B06: Bump to the version the sources hold writes nothing", func(t *testing.T) {
		manual := strings.Replace(appVersioning, "next: semantic", "next: manual", 1)
		root := historyRepo(t, "file:anchors.yaml="+manual, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start")
		before, _ := os.Stat(filepath.Join(root, "app/version.json"))
		out, err := versionRun(t, "--root", root, "--branch", "main", "bump")
		after, _ := os.Stat(filepath.Join(root, "app/version.json"))
		if err != nil || !strings.Contains(out, "app is already at 1.2.0") || !after.ModTime().Equal(before.ModTime()) {
			t.Errorf("%v\n%s", err, out)
		}
	})
	t.Run("VRSNI-E04: Bump refuses a version that does not move forward, and writes nothing", func(t *testing.T) {
		cfg := strings.Replace(appVersioning, "next: semantic", "next: { strategy: increment, part: patch }", 1) +
			"  api:\n    scheme: semver\n    tag: \"api-v{version}\"\n    next: semantic\n"
		root := historyRepo(t, "file:anchors.yaml="+cfg, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start", "tag:app-v1.3.0", "fix: x")
		for _, c := range []struct {
			args []string
			says string
		}{
			{[]string{"--product", "app", "--branch", "main", "bump", "1.1.0"}, "below 1.2.0"},
			{[]string{"--product", "app", "--branch", "main", "bump", "1.2.x"}, `"1.2.x"`},
			{[]string{"--product", "app", "--branch", "main", "bump", "1.2.5"}, "not above 1.3.0"},
			{[]string{"--product", "app", "--branch", "feat/x", "bump"}, "never versions"},
			{[]string{"--product", "api", "--branch", "main", "bump"}, "no source"},
		} {
			if _, err := versionRun(t, append([]string{"--root", root}, c.args...)...); err == nil || !strings.Contains(err.Error(), c.says) {
				t.Errorf("%v: %v, want it to say %q", c.args, err, c.says)
			}
		}
		if got := readFile(t, root, "app/version.json"); got != `{"version": "1.2.0"}` {
			t.Errorf("a refused bump wrote: %s", got)
		}
	})
}

func TestVersion_refused(t *testing.T) {
	t.Run("VRSNI-E01: No block, an unknown product, or several products without a name", func(t *testing.T) {
		bare := historyRepo(t, "file:anchors.yaml=version: 8\n", "chore: start")
		if _, err := versionRun(t, "--root", bare, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "no `versioning:`") {
			t.Errorf("no block: %v", err)
		}
		two := appVersioning + "  api:\n    scheme: number\n    next: none\n"
		root := historyRepo(t, "file:anchors.yaml="+two, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start")
		if _, err := versionRun(t, "--root", root, "--branch", "main", "--product", "web"); err == nil || !strings.Contains(err.Error(), "api, app") {
			t.Errorf("unknown: %v", err)
		}
		if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "--product") || !strings.Contains(err.Error(), "api, app") {
			t.Errorf("several: %v", err)
		}
	})
	t.Run("VRSNI-E03: Sources that disagree or a tag the scheme cannot read stop the command", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+twoSources, "file:app/version.json={\"version\": \"1.2.0\"}", "file:app/build.txt=VERSION=1.3.0", "chore: start")
		if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "app/version.json 1.2.0") || !strings.Contains(err.Error(), "app/build.txt 1.3.0") {
			t.Errorf("disagree: %v", err)
		}
		root = historyRepo(t, "file:anchors.yaml="+appVersioning, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start", "tag:app-vnext")
		if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "app-vnext") {
			t.Errorf("bad tag: %v", err)
		}
	})
}

func TestVersion_range(t *testing.T) {
	t.Run("VRSNI-B08: A branch's range shows with its rule, and bump refuses a version outside it", func(t *testing.T) {
		cfg := strings.Replace(appVersioning, "        next: { strategy: increment, part: patch }\n", "        next: { strategy: increment, part: patch }\n        range: from_base\n", 1)
		root := historyRepo(t, "file:anchors.yaml="+cfg, "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start", "tag:app-v1.2.0", "fix: x")
		out, err := versionRun(t, "--root", root, "--branch", "hotfix/login")
		if err != nil || !strings.Contains(out, "range      from_base  hotfix/*") {
			t.Errorf("%v\n%s", err, out)
		}
		if _, err := versionRun(t, "--root", root, "--branch", "hotfix/login", "bump", "minor"); err == nil || !strings.Contains(err.Error(), "1.3.0 is outside 1.2.x (from_base: 1.2.0)") {
			t.Errorf("minor on a hotfix: %v", err)
		}
		if _, err := versionRun(t, "--root", root, "--branch", "hotfix/login", "bump"); err != nil || readFile(t, root, "app/version.json") != `{"version": "1.2.1"}` {
			t.Errorf("patch: %v %s", err, readFile(t, root, "app/version.json"))
		}
	})
}

func TestVersion_intents(t *testing.T) {
	t.Run("VRSNI-B09: With the intents strategy, the next version is the latest release raised by the largest bump of the product's intents", func(t *testing.T) {
		cfg := strings.Replace(appVersioning, "next: semantic", "next: intents", 1)
		root := historyRepo(t, "file:anchors.yaml="+cfg, "file:app/version.json={\"version\": \"1.1.0\"}",
			"file:.anchors/changes/a.md=---\nproduct: app\nbump: patch\ntype: fix\n---\nA crash.\n",
			"file:.anchors/changes/b.md=---\nproduct: app\nbump: minor\ntype: feat\n---\nSharing.\n",
			"file:.anchors/changes/c.md=---\nproduct: api\nbump: major\ntype: feat\n---\nOther product.\n",
			"chore: start", "tag:app-v1.0.0")
		out, err := versionRun(t, "--root", root, "--branch", "main")
		if err != nil || !strings.Contains(out, "next       1.1.0      intents: minor, the largest bump of 2 intent(s) since 1.0.0") {
			t.Errorf("%v\n%s", err, out)
		}
		writeIn(t, root, ".anchors/changes/d.md", "no front matter")
		if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "d.md") {
			t.Errorf("an intent that does not parse: %v", err)
		}
	})
}

func writeIn(t *testing.T, root, rel, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, rel)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestVersion_promoteAndScript(t *testing.T) {
	t.Run("VRSNI-B10: promote takes the version the branch carries at its head, or its release with no source, and fails naming a branch it cannot read", func(t *testing.T) {
		cfg := strings.Replace(appVersioning, "next: semantic", "next: { strategy: promote, from: review }", 1)
		root := historyRepo(t, "file:anchors.yaml="+cfg, "file:app/version.json={\"version\": \"1.0.0\"}", "chore: start", "tag:app-v1.0.0")
		gitIn(t, root, "add", "-A")
		gitIn(t, root, "commit", "-q", "-m", "chore: files")
		gitIn(t, root, "checkout", "-q", "-b", "review")
		writeIn(t, root, "app/version.json", `{"version": "1.1.0"}`)
		gitIn(t, root, "commit", "-q", "-am", "chore: 1.1.0")
		gitIn(t, root, "checkout", "-q", "main")
		out, err := versionRun(t, "--root", root, "--branch", "main")
		if err != nil || !strings.Contains(out, "next       1.1.0      promote: the version review carries (") {
			t.Errorf("%v\n%s", err, out)
		}
		tagOnly := strings.Replace(cfg, "    source:\n      - file: app/version.json\n        path: version\n", "", 1)
		writeIn(t, root, "anchors.yaml", tagOnly)
		gitIn(t, root, "tag", "app-v1.2.0", "review")
		if out, err := versionRun(t, "--root", root, "--branch", "main"); err != nil || !strings.Contains(out, "promote: the release of review, app-v1.2.0") {
			t.Errorf("tag only: %v\n%s", err, out)
		}
		writeIn(t, root, "anchors.yaml", strings.Replace(cfg, "from: review", "from: qa", 1))
		if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), "promote: no branch qa") {
			t.Errorf("no branch: %v", err)
		}
	})
	t.Run("VRSNI-B11: script prints the version, reading the product's state on its stdin; a script that fails or prints nothing decides nothing", func(t *testing.T) {
		script := `next: { strategy: script, run: "cat > ctx.json; echo 3.0.0" }`
		root := historyRepo(t, "file:anchors.yaml="+strings.Replace(appVersioning, "next: semantic", script, 1), "file:app/version.json={\"version\": \"1.2.0\"}", "chore: start", "tag:app-v1.1.0", "feat: one")
		out, err := versionRun(t, "--root", root, "--branch", "main")
		if err != nil || !strings.Contains(out, "next       3.0.0      script: `cat > ctx.json; echo 3.0.0`") {
			t.Errorf("%v\n%s", err, out)
		}
		if ctx := readFile(t, root, "ctx.json"); !strings.Contains(ctx, `"released":"1.1.0"`) || !strings.Contains(ctx, `"subject":"feat: one"`) {
			t.Errorf("the state on stdin: %s", ctx)
		}
		for run, says := range map[string]string{"exit 3": "`exit 3` failed", "true": "`true` printed no version"} {
			writeIn(t, root, "anchors.yaml", strings.Replace(appVersioning, "next: semantic", `next: { strategy: script, run: "`+run+`" }`, 1))
			if _, err := versionRun(t, "--root", root, "--branch", "main"); err == nil || !strings.Contains(err.Error(), says) {
				t.Errorf("%s: %v", run, err)
			}
		}
	})
}

const labelledProduct = `version: 8
versioning:
  app:
    scheme: semver
    source:
      - file: app/version.json
        path: version
    tag: "app-v{version}"
    next: semantic
    after_bump: 'printf "%s %s\n" "$ANCHORS_PRODUCT" "$ANCHORS_VERSION" > bumped.txt'
    branches:
      review:
        label: { name: rc, in: source }
      qa:
        label: { name: rc }
`

func TestVersion_label(t *testing.T) {
	t.Run("VRSNI-B12: A labelled branch shows its next prerelease, one above the tags; bump writes it where the label goes; a prerelease is not the latest release", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+labelledProduct, "file:app/version.json={\"version\": \"1.0.0\"}", "chore: start", "tag:app-v1.0.0", "feat: one", "tag:app-v1.1.0-rc.1")
		out, err := versionRun(t, "--root", root, "--branch", "review")
		if err != nil || !strings.Contains(out, "released   1.0.0") || !strings.Contains(out, "label      1.1.0-rc.2 in source") || !strings.Contains(out, "next       1.1.0") {
			t.Errorf("%v\n%s", err, out)
		}
		if _, err := versionRun(t, "--root", root, "--branch", "review", "bump"); err != nil || readFile(t, root, "app/version.json") != `{"version": "1.1.0-rc.2"}` {
			t.Errorf("in source: %v %s", err, readFile(t, root, "app/version.json"))
		}
		writeIn(t, root, "app/version.json", `{"version": "1.0.0"}`)
		if _, err := versionRun(t, "--root", root, "--branch", "qa", "bump"); err != nil || readFile(t, root, "app/version.json") != `{"version": "1.1.0"}` {
			t.Errorf("in tag: %v %s", err, readFile(t, root, "app/version.json"))
		}
	})
	t.Run("VRSNI-B13: after_bump runs in the root once every source was written, with the product and the version; when it fails the command says the sources were written", func(t *testing.T) {
		root := historyRepo(t, "file:anchors.yaml="+labelledProduct, "file:app/version.json={\"version\": \"1.0.0\"}", "chore: start", "tag:app-v1.0.0", "fix: one")
		if _, err := versionRun(t, "--root", root, "--branch", "main", "bump"); err != nil || readFile(t, root, "bumped.txt") != "app 1.0.1\n" {
			t.Errorf("%v", err)
		}
		writeIn(t, root, "anchors.yaml", strings.Replace(labelledProduct, `after_bump: 'printf "%s %s\n" "$ANCHORS_PRODUCT" "$ANCHORS_VERSION" > bumped.txt'`, "after_bump: exit 4", 1))
		if _, err := versionRun(t, "--root", root, "--branch", "main", "bump", "1.0.2"); err == nil || !strings.Contains(err.Error(), "after_bump `exit 4` failed") || !strings.Contains(err.Error(), "already hold 1.0.2") {
			t.Errorf("%v", err)
		}
	})
}
