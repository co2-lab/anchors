// @anchors
//   code: BFTBR
//   ref: BFSBR

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mifFlow is the flow of a project with a work branch, a review branch, a production branch
// and hotfixes taken from production.
func mifFlow() *BranchFlow {
	return &BranchFlow{
		Merges: map[string][]string{
			"review":     {"develop", "hotfix/*"},
			"production": {"review", "hotfix/*"},
			"develop":    {"review", "production", "hotfix/*"},
		},
		Protected: []string{"review", "production"},
	}
}

func TestBranchFlow_aMergeIsAllowedByItsPair(t *testing.T) {
	t.Run("BFSBR-B01: A merge is allowed when the base's entry names its head", func(t *testing.T) {})
	t.Run("BFSBR-B02: A base no entry names takes any head, unless it is protected", func(t *testing.T) {})
	f := mifFlow()
	for _, c := range []struct {
		head, base string
		ok         bool
	}{
		{"develop", "review", true}, {"hotfix/login", "review", true}, {"review", "production", true},
		{"hotfix/login", "production", true}, {"production", "develop", true},
		{"develop", "production", false}, {"production", "review", false}, {"feat/x", "review", false},
		{"hotfix/a/b", "production", false},
		{"feat/x", "feat/y", true},
	} {
		if ok, _ := f.MergeAllowed(c.head, c.base); ok != c.ok {
			t.Errorf("%s → %s: allowed %v, want %v", c.head, c.base, ok, c.ok)
		}
	}
	if _, heads := f.MergeAllowed("develop", "production"); strings.Join(heads, ",") != "review,hotfix/*" {
		t.Errorf("a refusal names what the base takes: %v", heads)
	}
	guarded := &BranchFlow{Protected: []string{"production"}}
	if ok, _ := guarded.MergeAllowed("review", "production"); ok {
		t.Error("a protected base no entry names takes nothing")
	}
	var none *BranchFlow
	if ok, _ := none.MergeAllowed("develop", "production"); !ok || none.IsProtected("production") {
		t.Error("no flow declared guards nothing")
	}
}

func TestBranchFlow_patterns(t *testing.T) {
	t.Run("BFSBR-B03: A branch pattern is a glob whose * does not cross /, or a regular expression after re:", func(t *testing.T) {})
	for _, c := range []struct {
		pat, branch string
		ok          bool
	}{
		{"hotfix/*", "hotfix/login", true}, {"hotfix/*", "hotfix/a/b", false}, {"review", "review", true},
		{"review", "review-2", false}, {"re:release/v[0-9]+", "release/v12", true}, {"re:release/v[0-9]+", "xrelease/v1", false},
	} {
		if BranchMatches(c.pat, c.branch) != c.ok {
			t.Errorf("%q ~ %q: want %v", c.pat, c.branch, c.ok)
		}
	}
	f := &BranchFlow{Merges: map[string][]string{"release/*": {"develop"}, "release/hot": {"hotfix/*"}}}
	if ok, _ := f.MergeAllowed("hotfix/x", "release/hot"); !ok {
		t.Error("the entry naming the base exactly wins over a pattern matching it")
	}
	if ok, _ := f.MergeAllowed("develop", "release/1"); !ok {
		t.Error("a pattern entry rules the bases it matches")
	}
	if !mifFlow().IsProtected("production") || mifFlow().IsProtected("develop") {
		t.Error("the protected branches are those named")
	}
}

func TestBranchFlow_aBadPatternDoesNotLoad(t *testing.T) {
	t.Run("BFSBR-E01: A branch pattern that is neither a glob nor a regular expression refuses the configuration", func(t *testing.T) {})
	dir := t.TempDir()
	for _, bad := range []string{"  merges:\n    review: [\"re:(\"]\n", "  protected: [\"[\"]\n"} {
		p := filepath.Join(dir, "anchors.yaml")
		if err := os.WriteFile(p, []byte("version: 8\nbranch_flow:\n"+bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(p); err == nil || !strings.Contains(err.Error(), "branch_flow") {
			t.Errorf("a bad pattern refuses the configuration, naming branch_flow: %v", err)
		}
	}
	p := filepath.Join(dir, "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\nbranch_flow:\n  merges:\n    review: [develop]\n  protected: [review]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(p)
	if err != nil || c.BranchFlow == nil || !c.BranchFlow.IsProtected("review") {
		t.Errorf("the flow loads: %v %+v", err, c)
	}
}
