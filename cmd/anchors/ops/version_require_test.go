// @anchors
//   code: VRQRC
//   ref: VRQRA

package ops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

const releasedProject = `version: 8
branch_flow:
  merges:
    review: [develop, "hotfix/*"]
    production: [review, "hotfix/*"]
versioning:
  app:
    scheme: semver
    source: [{ file: app/version.json, path: version }]
    tag: "app-v{version}"
    released: { file: policy.json, keys: versions, where: { target: production } }
    next: semantic
    branches:
      "hotfix/*":
        next: { strategy: increment, part: patch }
        range: from_base
      "feat/*":
        next: none
      "maint/*":
        next: manual
pipeline:
  stages:
    review-release:
      branches: [review]
      requires:
        - contains: production
        - version: { fresh: true, listed: { where: { target: review } } }
      verify: x
    production-release:
      branches: [production]
      requires:
        - version: { same_as_passed: review-release }
      verify: x
`

// releasedRepo is a project whose production released app 1.0.0 — tagged, and listed with
// target production —, cloned from a bare origin that has its production branch.
func releasedRepo(t *testing.T) string {
	t.Helper()
	origin := filepath.Join(t.TempDir(), "origin.git")
	gitIn(t, filepath.Dir(origin), "init", "-q", "--bare", "-b", "production", origin)
	root := t.TempDir()
	gitIn(t, root, "init", "-q", "-b", "production")
	gitIn(t, root, "remote", "add", "origin", origin)
	commitFiles(t, root, "chore: release 1.0.0", map[string]string{
		"anchors.yaml":     releasedProject,
		"app/version.json": `{"version": "1.0.0"}`,
		"policy.json":      `{"versions": {"1.0.0": {"target": "production"}}}`,
	})
	gitIn(t, root, "tag", "app-v1.0.0")
	gitIn(t, root, "push", "-q", "origin", "production")
	gitIn(t, root, "checkout", "-q", "-b", "review")
	return root
}

// commitFiles writes the files and commits them, returning the commit.
func commitFiles(t *testing.T, root, msg string, files map[string]string) string {
	t.Helper()
	for f, body := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, f)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, f), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "--allow-empty", "-m", msg)
	return headOf(t, root)
}

func headOf(t *testing.T, root string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func configOf(t *testing.T, root string) *config.Config {
	t.Helper()
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func appVersion(v string) map[string]string {
	return map[string]string{"app/version.json": `{"version": "` + v + `"}`}
}

func holdsSays(t *testing.T, label string, ok bool, why string, wantOK bool, says ...string) {
	t.Helper()
	if ok != wantOK {
		t.Errorf("%s: holds=%v, want %v: %s", label, ok, wantOK, why)
	}
	for _, s := range says {
		if !strings.Contains(why, s) {
			t.Errorf("%s: %q does not say %q", label, why, s)
		}
	}
}

func TestVersionRequirement_fresh(t *testing.T) {
	t.Run("VRQRA-B01: fresh holds for a version above every release — the tags and what released: lists — and refuses one released or below", func(t *testing.T) {
		root := releasedRepo(t)
		fresh := &config.VersionRequirement{Fresh: true}
		sha := commitFiles(t, root, "fix: a crash", nil)
		ok, why := versionHolds(root, configOf(t, root), fresh, "review", sha, nil)
		holdsSays(t, "kept 1.0.0", ok, why, false, "1.0.0 is released (tag app-v1.0.0)", "bump to a greater version")
		sha = commitFiles(t, root, "chore: 1.1.0", map[string]string{"app/version.json": `{"version": "1.1.0"}`,
			"policy.json": `{"versions": {"1.0.0": {"target": "production"}, "1.1.0": {"target": "production"}}}`})
		ok, why = versionHolds(root, configOf(t, root), fresh, "review", sha, nil)
		holdsSays(t, "listed as released, with no tag", ok, why, false, "1.1.0 is released (policy.json versions, where target = production)")
		sha = commitFiles(t, root, "chore: 1.1.0 in review", map[string]string{"policy.json": `{"versions": {"1.0.0": {"target": "production"}, "1.1.0": {"target": "review"}}}`})
		ok, why = versionHolds(root, configOf(t, root), fresh, "review", sha, nil)
		holdsSays(t, "in review, resubmitted", ok, why, true, "1.1.0 is above every release (1.0.0, tag app-v1.0.0)")
		below := commitFiles(t, root, "chore: back", appVersion("0.9.0"))
		ok, why = versionHolds(root, configOf(t, root), fresh, "review", below, nil)
		holdsSays(t, "below", ok, why, false, "0.9.0 is below 1.0.0")
	})
	t.Run("VRQRA-B07: The version is read at the sha, not in the working tree", func(t *testing.T) {
		root := releasedRepo(t)
		sha := commitFiles(t, root, "chore: 1.1.0", appVersion("1.1.0"))
		if err := os.WriteFile(filepath.Join(root, "app/version.json"), []byte(`{"version": "1.0.0"}`), 0o644); err != nil {
			t.Fatal(err)
		}
		ok, why := versionHolds(root, configOf(t, root), &config.VersionRequirement{Fresh: true}, "review", sha, nil)
		holdsSays(t, "at the sha", ok, why, true, "1.1.0")
	})
}

func TestVersionRequirement_listed(t *testing.T) {
	t.Run("VRQRA-B02: listed holds when the file at the sha has an entry for the version matching where, and names the listing when it has not", func(t *testing.T) {
		root := releasedRepo(t)
		review := &config.VersionRequirement{Listed: &config.Listing{Where: map[string]string{"target": "review"}}}
		sha := commitFiles(t, root, "chore: 1.0.1", appVersion("1.0.1"))
		ok, why := versionHolds(root, configOf(t, root), review, "review", sha, nil)
		holdsSays(t, "not listed", ok, why, false, "1.0.1 is not in policy.json versions, where target = review")
		sha = commitFiles(t, root, "chore: list 1.0.1", map[string]string{"policy.json": `{"versions": {"1.0.0": {"target": "production"}, "1.0.1": {"target": "review"}}}`})
		ok, why = versionHolds(root, configOf(t, root), review, "review", sha, nil)
		holdsSays(t, "listed", ok, why, true, "1.0.1 is in policy.json versions, where target = review")
		other := &config.VersionRequirement{Listed: &config.Listing{File: "envs.yaml", Where: map[string]string{"env": "qa"}}}
		sha = commitFiles(t, root, "chore: envs", map[string]string{"envs.yaml": "1.0.1:\n  env: qa\n"})
		ok, why = versionHolds(root, configOf(t, root), other, "review", sha, nil)
		holdsSays(t, "a file of its own", ok, why, true, "1.0.1 is in envs.yaml, where env = qa")
	})
}

func TestVersionRequirement_sameAsPassed(t *testing.T) {
	t.Run("VRQRA-B03: same_as_passed holds when the sha carries the version the commit that passed the stage carried", func(t *testing.T) {
		root := releasedRepo(t)
		passedWith := commitFiles(t, root, "chore: 1.0.1", appVersion("1.0.1"))
		same := commitFiles(t, root, "docs: a note", nil)
		other := commitFiles(t, root, "chore: 1.0.2", appVersion("1.0.2"))
		req := &config.VersionRequirement{SameAsPassed: "review-release"}
		by := func(string) (string, string, string) { return passedWith, "the head of #9", "" }
		ok, why := versionHolds(root, configOf(t, root), req, "production", same, by)
		holdsSays(t, "same", ok, why, true, "1.0.1, the version the head of #9 passed review-release with")
		ok, why = versionHolds(root, configOf(t, root), req, "production", other, by)
		holdsSays(t, "other", ok, why, false, "this commit carries 1.0.2, and the head of #9 passed review-release with 1.0.1")
		none := func(string) (string, string, string) {
			return "", "", "this commit has not passed review-release (no record)"
		}
		ok, why = versionHolds(root, configOf(t, root), req, "production", same, none)
		holdsSays(t, "nothing passed", ok, why, false, "has not passed review-release")
	})
}

func TestVersionRequirement_expected(t *testing.T) {
	t.Run("VRQRA-B04: expected holds when the version is the one the branch's strategy computes; manual asks it be fresh, none refuses", func(t *testing.T) {
		root := releasedRepo(t)
		req := &config.VersionRequirement{Expected: true}
		commitFiles(t, root, "feat: a new screen", nil)
		sha := commitFiles(t, root, "chore: 1.0.1", appVersion("1.0.1"))
		ok, why := versionHolds(root, configOf(t, root), req, "review", sha, nil)
		holdsSays(t, "semantic", ok, why, false, "review: semantic asks for 1.1.0 (semantic: ", "the source says 1.0.1")
		sha = commitFiles(t, root, "chore: 1.1.0", appVersion("1.1.0"))
		ok, why = versionHolds(root, configOf(t, root), req, "review", sha, nil)
		holdsSays(t, "as asked", ok, why, true, "1.1.0, as semantic asks")
		ok, why = versionHolds(root, configOf(t, root), req, "feat/x", sha, nil)
		holdsSays(t, "none", ok, why, false, "feat/x never versions (feat/*: next: none)")
		ok, why = versionHolds(root, configOf(t, root), req, "maint/1", sha, nil)
		holdsSays(t, "manual, fresh", ok, why, true, "1.1.0 is above every release")
		kept := commitFiles(t, root, "chore: back to 1.0.0", appVersion("1.0.0"))
		ok, why = versionHolds(root, configOf(t, root), req, "maint/1", kept, nil)
		holdsSays(t, "manual, released", ok, why, false, "1.0.0 is released")
	})
	t.Run("VRQRA-B05: A branch's range bounds the version whatever is asked: from_base keeps the release the branch was cut from", func(t *testing.T) {
		root := releasedRepo(t)
		gitIn(t, root, "checkout", "-q", "-b", "hotfix/login", "production")
		req := &config.VersionRequirement{Expected: true}
		sha := commitFiles(t, root, "fix: login", appVersion("1.1.0"))
		ok, why := versionHolds(root, configOf(t, root), req, "hotfix/login", sha, nil)
		holdsSays(t, "outside", ok, why, false, "1.1.0 is outside 1.0.x (from_base: 1.0.0) (hotfix/*)")
		sha = commitFiles(t, root, "chore: 1.0.1", appVersion("1.0.1"))
		ok, why = versionHolds(root, configOf(t, root), req, "hotfix/login", sha, nil)
		holdsSays(t, "inside", ok, why, true, "1.0.1 is inside 1.0.x", "1.0.1, as increment asks")
	})
}

func TestContainsRequirement(t *testing.T) {
	t.Run("VRQRA-B06: contains holds when the head of the branch on the remote is an ancestor of the sha, and names the commits missing when it is not", func(t *testing.T) {
		root := releasedRepo(t)
		gitIn(t, root, "checkout", "-q", "production")
		commitFiles(t, root, "fix: a crash on production", nil)
		gitIn(t, root, "push", "-q", "origin", "production")
		gitIn(t, root, "checkout", "-q", "review")
		sha := commitFiles(t, root, "feat: review goes on", nil)
		ok, why := containsHolds(root, "production", sha)
		holdsSays(t, "missing", ok, why, false, "1 commit(s) of origin/production are not in this commit", "fix: a crash on production")
		gitIn(t, root, "merge", "-q", "--no-edit", "production")
		ok, why = containsHolds(root, "production", headOf(t, root))
		holdsSays(t, "back", ok, why, true, "origin/production (")
	})
	t.Run("VRQRA-E01: A branch the remote does not have, or a version the sha cannot give, fails the requirement naming it", func(t *testing.T) {
		root := releasedRepo(t)
		ok, why := containsHolds(root, "nope", headOf(t, root))
		holdsSays(t, "no branch", ok, why, false, "no origin/nope")
		sha := commitFiles(t, root, "chore: broken", map[string]string{"app/version.json": `{"name": "x"}`})
		ok, why = versionHolds(root, configOf(t, root), &config.VersionRequirement{Fresh: true}, "review", sha, nil)
		holdsSays(t, "no version", ok, why, false, "app/version.json")
	})
}

func TestVersionRequirement_label(t *testing.T) {
	t.Run("VRQRA-B08: expected takes a prerelease of the computed version on a branch whose label goes into the sources", func(t *testing.T) {
		root := releasedRepo(t)
		cfg := strings.Replace(releasedProject, "      \"maint/*\":\n        next: manual\n", "      \"maint/*\":\n        next: manual\n      review:\n        label: { name: rc, in: both }\n", 1)
		commitFiles(t, root, "feat: a screen", map[string]string{"anchors.yaml": cfg})
		sha := commitFiles(t, root, "chore: rc", appVersion("1.1.0-rc.1"))
		ok, why := versionHolds(root, configOf(t, root), &config.VersionRequirement{Expected: true}, "review", sha, nil)
		holdsSays(t, "labelled", ok, why, true, "1.1.0-rc.1")
		sha = commitFiles(t, root, "chore: other", appVersion("1.2.0-rc.1"))
		ok, why = versionHolds(root, configOf(t, root), &config.VersionRequirement{Expected: true}, "review", sha, nil)
		holdsSays(t, "another version", ok, why, false, "semantic asks for 1.1.0")
	})
}
