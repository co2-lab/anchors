// @anchors
//   code: PSRCF
//   ref: PSRSA

package config

import (
	"fmt"
	"path"
	"strings"
)

// PushReview asks each push to carry its evidence before it leaves the machine: a review of
// the commits it sends, written to a file the repository keeps, and — apart — an answer from
// every judgment gate those commits reach. The pre-push hook confronts both; the agent
// pushing reviews, runs the judgment gates while it does, and records the verdict.
type PushReview struct {
	// Required makes the pre-push hook refuse commits no review covers, or whose latest
	// review failed.
	Required bool `yaml:"required,omitempty"`
	// File is where the reviews are written, relative to the root (default REVIEWS.md).
	File string `yaml:"file,omitempty" path:"write"`
	// Retention is how long a review stays in the file, in days (`30d`): writing a review
	// drops those older. Empty keeps every review.
	Retention string `yaml:"retention,omitempty"`
	// RequireJudgment makes the pre-push hook refuse commits that reach a judgment gate no
	// verdict answers in the map the push sends.
	RequireJudgment bool `yaml:"require_judgment,omitempty"`
}

// DefaultReviewFile is where the reviews are written when the project names no file.
const DefaultReviewFile = "REVIEWS.md"

// FileOrDefault is the declared file, or REVIEWS.md.
func (p *PushReview) FileOrDefault() string {
	if p != nil && strings.TrimSpace(p.File) != "" {
		return path.Clean(strings.TrimSpace(p.File))
	}
	return DefaultReviewFile
}

// RetentionDays is Retention in days, 0 when every review is kept.
func (p *PushReview) RetentionDays() int {
	if p == nil {
		return 0
	}
	return ExternalSystem{MaxAge: p.Retention}.MaxAgeDays()
}

// Enforced says whether the pre-push hook confronts anything: a review, or the judgments.
func (p *PushReview) Enforced() bool {
	return p != nil && (p.Required || p.RequireJudgment)
}

// validatePushReview refuses a retention that is not a number of days, and a file outside
// the project.
func validatePushReview(p *PushReview) error {
	if p == nil {
		return nil
	}
	if r := strings.TrimSpace(p.Retention); r != "" && !maxAgeRE.MatchString(r) {
		return fmt.Errorf("`push_review.retention` must be a number of days (30d), got %q", p.Retention)
	}
	if f := strings.TrimSpace(p.File); f != "" {
		clean := path.Clean(strings.ReplaceAll(f, "\\", "/"))
		if path.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(f, ":") {
			return fmt.Errorf("`push_review.file` must be a path inside the project, got %q", p.File)
		}
	}
	return nil
}
