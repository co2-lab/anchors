// @anchors
//   ref: RVCMR

package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// reviewProject is a project whose rule-fulfilled gate is judged and reviewed.
func reviewProject(t *testing.T, workflow string) string {
	t.Helper()
	useEnglish(t)
	root := t.TempDir()
	for p, c := range map[string]string{
		"anchors.yaml": "version: 6\n" + workflow + "layers:\n  spec: {pattern: \"src/*.spec.md\", kind: spec}\n" +
			"gates:\n  - name: rule-fulfilled\n    measures: judgment\n    on: [spec]\n    ask: \"does the code do the rule?\"\n    review:\n      ask: \"look again at what the agent decided\"\n",
		"src/pay.spec.md": "<!-- @anchors\n  code: PAYMT\n-->\n# Pay\n\nPAYMT-B01 — pays.\n",
	} {
		writeProjectFile(t, root, p, c)
	}
	runCmd(t, newMapCmd(), "build", "--root", root)
	return root
}

func reviewsOf(t *testing.T, root string) []mapx.Review {
	t.Helper()
	return node(t, root, "src/pay.spec.md").Reviews
}

func TestReview_pendingListsWhatIsDue(t *testing.T) {
	t.Run("RVCMR-B01: Pending lists what is to review with its question", func(t *testing.T) {})
	root := reviewProject(t, "")
	out := runCmd(t, newReviewCmd(), "--root", root, "--pending")
	for _, want := range []string{"rule-fulfilled — look again at what the agent decided", "○ src/pay.spec.md", "1 target(s) to review"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	runCmd(t, newReviewCmd(), "--root", root, "src/pay.spec.md", "--gate", "rule-fulfilled", "--by", "human:ana")
	if out := runCmd(t, newReviewCmd(), "--root", root, "--pending"); !strings.Contains(out, "nothing to review") {
		t.Errorf("reviewed, nothing is due:\n%s", out)
	}
}

func TestReview_withNoFindingsRecordsWhoLooked(t *testing.T) {
	t.Run("RVCMR-B02: A review with no findings records who looked and opens no issue", func(t *testing.T) {})
	root := reviewProject(t, "")
	runCmd(t, newReviewCmd(), "--root", root, "src/pay.spec.md", "--gate", "rule-fulfilled", "--by", "human:ana")
	rs := reviewsOf(t, root)
	if len(rs) != 1 || rs[0].By != "human:ana" || rs[0].Findings || rs[0].Rev == "" {
		t.Errorf("one clean review by human:ana at the spec's revision: %+v", rs)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "issues", "todo")); len(entries) != 0 {
		t.Errorf("a clean review opens no issue: %v", entries)
	}
}

func TestReview_withFindingsOpensTheIssue(t *testing.T) {
	t.Run("RVCMR-B03: A review with findings records them and opens the issue", func(t *testing.T) {})
	root := reviewProject(t, "")
	runCmd(t, newReviewCmd(), "--root", root, "src/pay.spec.md", "--gate", "rule-fulfilled", "--by", "agent:openai/gpt",
		"--findings", "### 1. B01 pays twice (src/pay.go:3)")
	if rs := reviewsOf(t, root); len(rs) != 1 || !rs[0].Findings {
		t.Errorf("the review is marked as having found something: %+v", rs)
	}
	entries, _ := os.ReadDir(filepath.Join(root, "issues", "todo"))
	if len(entries) != 1 {
		t.Fatalf("one issue for the findings, got %v", entries)
	}
	b, _ := os.ReadFile(filepath.Join(root, "issues", "todo", entries[0].Name()))
	if !strings.Contains(string(b), "B01 pays twice") || !strings.Contains(entries[0].Name(), "rule-fulfilled") {
		t.Errorf("the issue holds the report for the gate:\n%s", b)
	}
}

func TestReview_manualModeWritesTheIssueOnlyWhenAsked(t *testing.T) {
	t.Run("RVCMR-B04: In manual mode the findings write no issue unless asked", func(t *testing.T) {})
	root := reviewProject(t, "workflow:\n  mode: manual\n")
	out := runCmd(t, newReviewCmd(), "--root", root, "src/pay.spec.md", "--gate", "rule-fulfilled", "--by", "x", "--findings", "### 1. finding")
	if !strings.Contains(out, "### 1. finding") {
		t.Errorf("manual mode prints the report:\n%s", out)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "issues", "todo")); len(entries) != 0 {
		t.Errorf("manual mode writes no issue: %v", entries)
	}
	runCmd(t, newReviewCmd(), "--root", root, "src/pay.spec.md", "--gate", "rule-fulfilled", "--by", "x", "--findings", "### 1. finding", "--record-issues")
	if entries, _ := os.ReadDir(filepath.Join(root, "issues", "todo")); len(entries) != 1 {
		t.Errorf("--record-issues writes it: %v", entries)
	}
}

func TestReview_refusals(t *testing.T) {
	t.Run("RVCMR-I01: A review record always names who reviewed", func(t *testing.T) {})
	t.Run("RVCMR-E01: A wrong target, gate or no reviewer refuses the record", func(t *testing.T) {})
	t.Run("RVCMR-X01: There is no waived", func(t *testing.T) {})
	root := reviewProject(t, "")
	for _, args := range [][]string{
		{"src/pay.spec.md", "--gate", "rule-fulfilled"},
		{"src/nope.spec.md", "--gate", "rule-fulfilled", "--by", "x"},
		{"src/pay.spec.md", "--gate", "other", "--by", "x"},
	} {
		if _, err := runCmdErr(newReviewCmd(), t, append([]string{"--root", root}, args...)...); err == nil {
			t.Errorf("%v should be refused", args)
		}
	}
	if rs := reviewsOf(t, root); len(rs) != 0 {
		t.Errorf("a refused record writes nothing: %+v", rs)
	}
	if newReviewCmd().Flags().Lookup("waived") != nil || newReviewCmd().Flags().Lookup("verdict") != nil {
		t.Error("the review command offers no waived and no verdict")
	}
}
