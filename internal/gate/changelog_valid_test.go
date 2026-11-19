// @anchors
//   code: CHVLC
//   ref: CHVLA

package gate

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/shell"
)

// changelogRepo is a project whose app released 1.0.0 and 1.1.0, each tagged on a commit of
// its own, with the given files committed at HEAD.
func changelogRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", root}, args...)...)
		c.Env = append(shell.WithoutRepoEnv(os.Environ()), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := c.CombinedOutput(); err != nil {
			t.Skipf("git: %v %s", err, out)
		}
	}
	git("init", "-q")
	git("commit", "-q", "--allow-empty", "-m", "feat: one")
	git("tag", "app-v1.0.0")
	git("commit", "-q", "--allow-empty", "-m", "feat: two")
	git("tag", "app-v1.1.0")
	for f, body := range files {
		writeFile(t, root, f, body)
	}
	git("add", ".")
	git("commit", "-q", "--allow-empty", "-m", "docs: changelog")
	return root
}

var tagDay = time.Now().Format("2006-01-02")

func block(v, body string) string {
	if v == "" {
		return changelog.Marker("") + "\n## Unreleased\n" + body + "\n"
	}
	return changelog.Marker(v) + "\n## " + v + " — " + tagDay + "\n" + body + "\n"
}

func appWithOutput(o config.ChangelogOutput) *config.Config {
	return &config.Config{Versioning: config.Versioning{"app": {Scheme: "semver", Tag: "app-v{version}", Next: &config.Next{Strategy: "none"},
		Changelog: &config.ProductChangelog{Outputs: []config.ChangelogOutput{o}}}}}
}

func incremental(f *config.ChangelogFormat) config.ChangelogOutput {
	return config.ChangelogOutput{Changelog: config.Changelog{Path: "CHANGELOG.md"}, Format: f}
}

const goodChangelog = "# Changelog\n\n"

func goodFile() string {
	return goodChangelog + block("", "### Features\n- three") + "\n" + block("1.1.0", "### Features\n- two") + "\n" + block("1.0.0", "### Features\n- one")
}

func failsSaying(t *testing.T, label string, v Verdict, msg string, says ...string) {
	t.Helper()
	if v != Fail {
		t.Errorf("%s: %s %s", label, v, msg)
		return
	}
	for _, s := range says {
		if !strings.Contains(msg, s) {
			t.Errorf("%s: %q does not say %q", label, msg, s)
		}
	}
}

func TestChangelogValid_passes(t *testing.T) {
	t.Run("CHVLA-B01: An output whose blocks read as versions, newest first, one per tag, in its format, passes", func(t *testing.T) {
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": goodFile()})
		f := &config.ChangelogFormat{Heading: "^## {version} — {date}$", RequiredSections: []string{"Features"}, MaxChars: 200}
		if v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(f))); v != Pass {
			t.Errorf("%s %s", v, msg)
		}
	})
	t.Run("CHVLA-B08: No product with an output skips the gate", func(t *testing.T) {
		cfg := &config.Config{Versioning: config.Versioning{"app": {Scheme: "semver"}}}
		if v, msg := checkChangelogValid(config.Gate{}, t.TempDir(), nil, cfg); v != Skip || !strings.Contains(msg, "changelog.outputs") {
			t.Errorf("%s %s", v, msg)
		}
	})
}

func TestChangelogValid_structure(t *testing.T) {
	t.Run("CHVLA-B02: Blocks out of order, twice, an unreleased one below another, or a marker no version, fail", func(t *testing.T) {
		file := goodChangelog + block("1.0.0", "- one") + "\n" + block("1.1.0", "- two") + "\n" + block("1.1.0", "- again") + "\n" + block("", "- next") + "\n" + block("1.x", "- bad")
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": file})
		v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(nil)))
		failsSaying(t, "structure", v, msg, "1.1.0 comes after 1.0.0", "1.1.0 has two blocks", "unreleased block is not on top", `CHANGELOG.md 1.x: "1.x"`)
	})
	t.Run("CHVLA-B03: A tag with no block, a block with no tag, a heading off the format or a date that is not the tag's fail", func(t *testing.T) {
		file := goodChangelog + changelog.Marker("1.2.0") + "\n## 1.2.0 — " + tagDay + "\n- three\n\n" + changelog.Marker("1.1.0") + "\n## 1.1.0 — 2001-01-01\n- two\n\n" + changelog.Marker("1.0.5") + "\n# 1.0.5\n- x\n"
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": file})
		v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(&config.ChangelogFormat{Heading: "^## {version} — {date}$"})))
		failsSaying(t, "agreement", v, msg, "no block for the release(s) 1.0.0", "CHANGELOG.md 1.2.0: no tag", "1.1.0: the heading's date 2001-01-01 is not the tag's", `1.0.5: the heading "# 1.0.5"`)
		root = changelogRepo(t, nil)
		v, msg = checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(nil)))
		failsSaying(t, "no file", v, msg, "does not exist, and 2 release tag(s)")
	})
}

func TestChangelogValid_content(t *testing.T) {
	t.Run("CHVLA-B04: A required section missing, a section left empty, or a block above max_chars fails", func(t *testing.T) {
		file := goodChangelog + block("1.1.0", "### Fixes\n\n### Features\n- two") + "\n" + block("1.0.0", "### Fixes\n- "+strings.Repeat("x", 80))
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": file})
		f := &config.ChangelogFormat{RequiredSections: []string{"Features"}, MaxChars: 60}
		v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(f)))
		failsSaying(t, "content", v, msg, "1.1.0: the section Fixes is empty", "1.0.0: no section Features", "1.0.0: 114 characters, above format.max_chars 60")
	})
	t.Run("CHVLA-B05: The block of a release that changed against HEAD fails; the unreleased one may change", func(t *testing.T) {
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": goodFile()})
		writeFile(t, root, "CHANGELOG.md", strings.Replace(strings.Replace(goodFile(), "- two", "- two, rewritten", 1), "- three", "- three and four", 1))
		v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(incremental(nil)))
		failsSaying(t, "rewritten", v, msg, "CHANGELOG.md 1.1.0: the block of a release changed")
		if strings.Contains(msg, "unreleased") {
			t.Errorf("the unreleased block grows until the tag: %s", msg)
		}
	})
	t.Run("CHVLA-B06: required: changed fails an output unchanged since the latest release while the sources hold a version not released", func(t *testing.T) {
		root := changelogRepo(t, map[string]string{"CHANGELOG.md": goodFile(), "v.txt": "1.2.0"})
		exec.Command("git", "-C", root, "tag", "-f", "app-v1.1.0").Run()
		cfg := appWithOutput(incremental(&config.ChangelogFormat{Required: "changed"}))
		cfg.Versioning["app"].Source = []config.VersionSource{{File: "v.txt", Regex: `(.+)`}}
		v, msg := checkChangelogValid(config.Gate{}, root, nil, cfg)
		failsSaying(t, "unchanged", v, msg, "the sources hold 1.2.0, not released yet, and the output has not changed since app-v1.1.0")
		writeFile(t, root, "CHANGELOG.md", strings.Replace(goodFile(), "- three", "- three and four", 1))
		if v, msg := checkChangelogValid(config.Gate{}, root, nil, cfg); v != Pass {
			t.Errorf("changed: %s %s", v, msg)
		}
		writeFile(t, root, "v.txt", "1.1.0")
		writeFile(t, root, "CHANGELOG.md", goodFile())
		if v, msg := checkChangelogValid(config.Gate{}, root, nil, cfg); v != Pass {
			t.Errorf("a released version asks nothing: %s %s", v, msg)
		}
	})
}

func TestChangelogValid_perVersion(t *testing.T) {
	t.Run("CHVLA-B07: A per_version output has a file per tag named by format.file, plus an unreleased one at most", func(t *testing.T) {
		o := config.ChangelogOutput{Changelog: config.Changelog{Path: "notes", Mode: "per_version"}, Format: &config.ChangelogFormat{File: "v{version}.txt"}}
		root := changelogRepo(t, map[string]string{"notes/v1.1.0.txt": "two\n", "notes/unreleased.md": "next\n", "notes/README": "x\n", "notes/v1.3.0.txt": "?\n"})
		v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(o))
		failsSaying(t, "per_version", v, msg, "notes/README does not match format.file v{version}.txt", "no block for the release(s) 1.0.0", "notes/v1.3.0.txt 1.3.0: no tag")
		root = changelogRepo(t, map[string]string{"notes/v1.1.0.txt": "two\n", "notes/v1.0.0.txt": "one\n"})
		if v, msg := checkChangelogValid(config.Gate{}, root, nil, appWithOutput(o)); v != Pass {
			t.Errorf("complete: %s %s", v, msg)
		}
	})
}
