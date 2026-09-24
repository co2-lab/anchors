package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- renumbering the revisions a branch added, when the base took the same number ---
//
// A revision is numbered by the file it revises: `VLANV-R0003` is the third change of
// `VLANV`. Two open pull requests that revise the same spec each take the next free number,
// and both take the SAME one. The second to merge carries a `R0003` that already means
// something else on the base: two revisions answering to one code, and every citation of it
// ambiguous. `plan-change-justified` sees the count and the maximum disagree.
//
// The number stays (it is what says "changed three times"), and the collision is resolved
// where it happens — on the branch, at rebase. This file is the engine behind
// `anchors renumber`: pure text in, text out. The command does the git.
//
// ONLY WHAT THE BRANCH ADDED MOVES. A revision the base already has is history other people
// have read and cited; renumbering it would make every one of those citations point at the
// wrong change. The branch's revision is the one nobody outside the branch has seen yet.

// RevisionFile is a revision-bearing file (a spec, a plan) in the three versions the
// renumbering needs.
type RevisionFile struct {
	Path      string
	Working   string // the branch's content now
	MergeBase string // the content where the branch forked from the base ("" when absent)
	Base      string // the content on the base ref ("" when absent)
}

// Renumbering is one revision the branch added that moves to a free number.
type Renumbering struct {
	File string
	Old  string // `VLANV-R0003`
	New  string // `VLANV-R0004`
}

// revisionLine is one revision as it appears in a version of the file: the whole line is
// its identity, because two revisions with the same number differ in what they say.
type revisionLine struct {
	code string
	num  int
	line string
}

func revisionLines(content string) []revisionLine {
	var out []revisionLine
	for _, m := range revisionRE().FindAllStringSubmatch(content, -1) {
		n, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		out = append(out, revisionLine{code: m[1], num: n, line: m[0]})
	}
	return out
}

// lineSet is the set of lines of a version: a line present in the merge base existed before
// the branch, and a line absent from it is the branch's.
func lineSet(content string) map[string]bool {
	s := map[string]bool{}
	for _, l := range strings.Split(content, "\n") {
		s[l] = true
	}
	return s
}

// PlanRenumber finds, in each file, the revisions the branch ADDED whose number the base
// already uses, and assigns them the next free numbers.
//
// A revision is the branch's when its line did not exist at the merge base — with one
// exception: a revision the merge base already numbered and that appears ONCE on the branch
// is the same revision with its explanation edited, not a new one. After a rebase the base's
// revision and the branch's sit in the same file with one number; the base's line is in the
// merge base, and the branch's is not.
//
// When one added revision collides, EVERY added revision of that code is renumbered, in the
// order of its old number, starting after the highest number the base and the branch's kept
// revisions use. Renumbering only the colliding one would leave the sequence out of order
// (`R0005` above `R0004`).
func PlanRenumber(files []RevisionFile) ([]Renumbering, error) {
	var out []Renumbering
	for _, f := range files {
		mbLines := lineSet(f.MergeBase)
		mbNums := map[string]map[int]bool{}
		for _, r := range revisionLines(f.MergeBase) {
			if mbNums[r.code] == nil {
				mbNums[r.code] = map[int]bool{}
			}
			mbNums[r.code][r.num] = true
		}
		baseNums := map[string]map[int]bool{}
		for _, r := range revisionLines(f.Base) {
			if baseNums[r.code] == nil {
				baseNums[r.code] = map[int]bool{}
			}
			baseNums[r.code][r.num] = true
		}

		work := revisionLines(f.Working)
		count := map[string]int{}
		for _, r := range work {
			count[fmt.Sprintf("%s-%d", r.code, r.num)]++
		}

		added := map[string][]revisionLine{}
		taken := map[string]map[int]bool{} // kept on the branch + used on the base
		var codes []string
		for _, r := range work {
			if taken[r.code] == nil {
				taken[r.code] = map[int]bool{}
				codes = append(codes, r.code)
			}
			edited := mbNums[r.code][r.num] && count[fmt.Sprintf("%s-%d", r.code, r.num)] == 1
			if !mbLines[r.line] && !edited {
				added[r.code] = append(added[r.code], r)
				continue
			}
			taken[r.code][r.num] = true
		}

		for _, code := range codes {
			adds := added[code]
			if len(adds) == 0 {
				continue
			}
			for n := range baseNums[code] {
				taken[code][n] = true
			}
			seen := map[int]bool{}
			collides := false
			for _, a := range adds {
				if seen[a.num] {
					return nil, fmt.Errorf("%s: the branch adds `%s-R%04d` twice — "+
						"a citation of it cannot be told apart; number them by hand first", f.Path, code, a.num)
				}
				seen[a.num] = true
				if taken[code][a.num] {
					collides = true
				}
			}
			if !collides {
				continue
			}
			next := 0
			for n := range taken[code] {
				if n > next {
					next = n
				}
			}
			sort.SliceStable(adds, func(i, j int) bool { return adds[i].num < adds[j].num })
			for _, a := range adds {
				next++
				if next == a.num {
					continue
				}
				out = append(out, Renumbering{File: f.Path,
					Old: fmt.Sprintf("%s-R%04d", code, a.num), New: fmt.Sprintf("%s-R%04d", code, next)})
			}
		}
	}
	return out, nil
}

// revisionCodeRE matches a citation of a revision anywhere in a line: `VLANV-R0003`.
func revisionCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`\b[A-Z0-9]` + config.CodeLengthPattern() + `-R\d{4}\b`)
}

// RewriteRevisionCitations rewrites, in the lines the BRANCH added, every revision code the
// map renames, and returns the new content with the number of citations rewritten.
//
// Only the branch's lines: a line that existed at the merge base cites what the base means
// by that code — after a rebase, the other pull request's revision — and rewriting it would
// point a correct citation at the wrong change.
//
// One pass, reading the map once per match: `R0003→R0004` and `R0004→R0005` in the same
// file would otherwise chain, and the first revision would end as `R0005`.
func RewriteRevisionCitations(working, mergeBase string, renames map[string]string) (string, int) {
	if len(renames) == 0 {
		return working, 0
	}
	mb := lineSet(mergeBase)
	re := revisionCodeRE()
	n := 0
	lines := strings.Split(working, "\n")
	for i, l := range lines {
		if mb[l] {
			continue
		}
		lines[i] = re.ReplaceAllStringFunc(l, func(c string) string {
			if to, ok := renames[c]; ok {
				n++
				return to
			}
			return c
		})
	}
	return strings.Join(lines, "\n"), n
}
