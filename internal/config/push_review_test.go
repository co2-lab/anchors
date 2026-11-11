// @anchors
//   code: PSRTA
//   ref: PSRSA

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadPushReview(t *testing.T, section string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), DefaultFile)
	if err := os.WriteFile(p, []byte("version: 1\nlayers: {}\npush_review:\n"+section), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestPushReview_declared(t *testing.T) {
	t.Run("PSRSA-B01: The section declares the review, the judgments, the file and the retention", func(t *testing.T) {})
	c, err := loadPushReview(t, "  required: true\n  require_judgment: true\n  file: docs/reviews.md\n  retention: 30d\n")
	if err != nil {
		t.Fatal(err)
	}
	p := c.PushReview
	if !p.Required || !p.RequireJudgment || p.FileOrDefault() != "docs/reviews.md" || p.RetentionDays() != 30 || !p.Enforced() {
		t.Errorf("the section: %+v", p)
	}
	c, err = loadPushReview(t, "  required: true\n")
	if err != nil {
		t.Fatal(err)
	}
	if p := c.PushReview; p.FileOrDefault() != "REVIEWS.md" || p.RetentionDays() != 0 || p.RequireJudgment || !p.Enforced() {
		t.Errorf("only the review: %+v", p)
	}
	var none *PushReview
	if none.Enforced() || none.FileOrDefault() != "REVIEWS.md" || (&PushReview{File: "x.md"}).Enforced() {
		t.Error("without a demand the push confronts nothing")
	}
	c, err = loadPushReview(t, "  require_gates: true\n")
	if err != nil || !c.PushReview.RequireGates || c.PushReview.Required || !c.PushReview.Enforced() {
		t.Errorf("only the gates: %+v %v", c.PushReview, err)
	}
}

func TestPushReview_refused(t *testing.T) {
	t.Run("PSRSA-E01: A retention or a file that cannot hold refuses the configuration", func(t *testing.T) {})
	for want, section := range map[string]string{
		"push_review.retention": "  retention: 30\n",
		"push_review.file":      "  file: /tmp/reviews.md\n",
	} {
		if _, err := loadPushReview(t, section); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: want an error naming %s, got %v", section, want, err)
		}
	}
	if _, err := loadPushReview(t, "  file: ../reviews.md\n"); err == nil || !strings.Contains(err.Error(), "push_review.file") {
		t.Errorf("a file outside the project: %v", err)
	}
}
