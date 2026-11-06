// @anchors
//   code: BRFLB
//   ref: BFSBR

package config

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
)

// BranchFlow is the project's branch flow: which branches a branch may receive a merge from,
// and which branches take no commit of their own. A release flow run backwards — commits on
// the production branch, the review branch aligned from it — went unnoticed: nothing said
// which way the work goes (reported from MIF, 2026-10-23).
type BranchFlow struct {
	// Merges names, for each base branch, the head branches a pull request may merge into it.
	// A base no entry names takes a merge from anywhere.
	Merges map[string][]string `yaml:"merges,omitempty"`
	// Protected are the branches that take no commit of their own, only merges.
	Protected []string `yaml:"protected,omitempty"`
}

// BranchMatches says whether a branch name matches a pattern: a glob whose `*` does not
// cross `/` (`hotfix/*`), or a regular expression after `re:`, matched whole.
func BranchMatches(pattern, branch string) bool {
	if re, ok := strings.CutPrefix(pattern, "re:"); ok {
		r, err := regexp.Compile("^(?:" + re + ")$")
		return err == nil && r.MatchString(branch)
	}
	ok, err := path.Match(pattern, branch)
	return err == nil && ok
}

// MergeAllowed says whether a pull request may merge head into base, and the heads base
// takes. The entry whose pattern names base exactly wins over a pattern that matches it; a
// base no entry matches takes any head, unless it is protected — a protected branch takes
// only what an entry allows.
func (f *BranchFlow) MergeAllowed(head, base string) (bool, []string) {
	if f == nil {
		return true, nil
	}
	heads, found := f.Merges[base]
	if !found {
		keys := make([]string, 0, len(f.Merges))
		for k := range f.Merges {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if BranchMatches(k, base) {
				heads, found = f.Merges[k], true
				break
			}
		}
	}
	if !found {
		return !f.IsProtected(base), nil
	}
	for _, h := range heads {
		if BranchMatches(h, head) {
			return true, heads
		}
	}
	return false, heads
}

// IsProtected says whether the branch takes no commit of its own.
func (f *BranchFlow) IsProtected(branch string) bool {
	if f == nil || branch == "" {
		return false
	}
	for _, p := range f.Protected {
		if BranchMatches(p, branch) {
			return true
		}
	}
	return false
}

// validateBranchFlow refuses a pattern that is neither a glob nor a regular expression: it
// would match nothing, and the flow it guards would pass unguarded.
func validateBranchFlow(f *BranchFlow) error {
	if f == nil {
		return nil
	}
	var pats []string
	for base, heads := range f.Merges {
		pats = append(pats, base)
		pats = append(pats, heads...)
	}
	pats = append(pats, f.Protected...)
	sort.Strings(pats)
	for _, p := range pats {
		if re, ok := strings.CutPrefix(p, "re:"); ok {
			if _, err := regexp.Compile(re); err != nil {
				return fmt.Errorf("branch_flow: %q is not a regular expression: %v", p, err)
			}
		} else if _, err := path.Match(p, ""); err != nil {
			return fmt.Errorf("branch_flow: %q is not a branch pattern: %v", p, err)
		}
	}
	return nil
}
