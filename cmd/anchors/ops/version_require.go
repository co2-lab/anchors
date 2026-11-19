// @anchors
//   code: VRQRV
//   ref: VRQRA

package ops

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/version"
)

// A stage asks two things of the commit it runs on that the commit gate cannot see. The
// version: commit time does not know what was released on another branch, nor which file
// sends each version where (reported from a project: a review build whose version was not
// listed fell to the file's default and talked to production). And the fixes: a branch that
// goes on without what was fixed on production releases the bug again. Both are read at the
// sha, so the stage confronts the commit it builds, not the working tree.

// passedBy finds the commit that passed a stage, and says which it was; empty with the reason
// when none did.
type passedBy func(stage string) (sha, of, why string)

// versionHolds answers a `version:` requirement on sha, for the branch the stage runs on.
// passed is nil where no stage has run yet — a pull request —, and same_as_passed is then
// not asked.
func versionHolds(root string, v config.Versioning, r *config.VersionRequirement, branch, sha string, passed passedBy) (bool, string) {
	st, err := stateAt(root, v, r.Product, "product:", branch, sha)
	if err != nil {
		return false, err.Error()
	}
	cur := st.Current
	if cur == "" {
		cur = st.Next
	}
	if cur == "" {
		return false, st.Product + " has no version here: " + st.Reason
	}
	var says []string
	if st.Range != "" {
		ok, why, err := st.inRange(cur)
		if err != nil {
			return false, err.Error()
		}
		if !ok {
			return false, why
		}
		says = append(says, why)
	}
	if r.Fresh || (r.Expected && st.Strategy == version.Manual) {
		ok, why := isFresh(root, st, cur, sha)
		if !ok {
			return false, why
		}
		says = append(says, why)
	}
	if r.Listed != nil {
		l := st.product.ListingOf(*r.Listed)
		vs, err := l.Versions(readAt(root, sha))
		if err != nil {
			return false, err.Error()
		}
		if !contains(vs, cur) {
			return false, fmt.Sprintf("%s is not in %s", cur, l)
		}
		says = append(says, fmt.Sprintf("%s is in %s", cur, l))
	}
	if r.SameAsPassed != "" && passed != nil {
		at, of, why := passed(r.SameAsPassed)
		if at == "" {
			return false, why
		}
		then, err := versionAt(root, st.product, at)
		if err != nil {
			return false, fmt.Sprintf("the version %s passed %s with: %v", of, r.SameAsPassed, err)
		}
		if then != cur {
			return false, fmt.Sprintf("this commit carries %s, and %s passed %s with %s", cur, of, r.SameAsPassed, then)
		}
		says = append(says, fmt.Sprintf("%s, the version %s passed %s with", cur, of, r.SameAsPassed))
	}
	if r.Expected {
		switch st.Strategy {
		case version.None:
			return false, fmt.Sprintf("%s never versions (%s: next: none)", branch, st.Rule)
		case version.Manual:
		default:
			if st.Next != cur {
				return false, fmt.Sprintf("%s: %s asks for %s (%s), the source says %s", branch, st.Strategy, st.Next, st.Reason, cur)
			}
			says = append(says, fmt.Sprintf("%s, as %s asks (%s)", cur, st.Strategy, st.Reason))
		}
	}
	return true, strings.Join(says, "; ")
}

// isFresh says whether a version is above every released one: the product's tags, wherever
// they are reachable from, and what `released:` lists at the sha. A version still in review
// is not released, so resubmitting it passes.
func isFresh(root string, st *versionStateT, cur, sha string) (bool, string) {
	v, err := st.scheme.Parse(cur)
	if err != nil {
		return false, err.Error()
	}
	type release struct{ raw, where string }
	var rels []release
	out, err := exec.Command("git", "-C", root, "tag", "--list").Output()
	if err != nil {
		return false, "read the tags: " + err.Error()
	}
	for _, t := range strings.Fields(string(out)) {
		if raw, ok := st.product.VersionOfTag(t); ok {
			rels = append(rels, release{raw, "tag " + t})
		}
	}
	if l := st.product.Released; l != nil {
		vs, err := l.Versions(readAt(root, sha))
		if err != nil {
			return false, err.Error()
		}
		for _, raw := range vs {
			rels = append(rels, release{raw, l.String()})
		}
	}
	var top *release
	var topV version.V
	for i, r := range rels {
		rv, err := st.scheme.Parse(r.raw)
		if err != nil {
			return false, fmt.Sprintf("a release, %s (%s): %v", r.raw, r.where, err)
		}
		if top == nil || version.Compare(rv, topV) > 0 {
			top, topV = &rels[i], rv
		}
	}
	switch {
	case top == nil:
		return true, cur + " is fresh: nothing is released yet"
	case version.Compare(v, topV) == 0:
		return false, fmt.Sprintf("%s is released (%s); bump to a greater version", cur, top.where)
	case version.Compare(v, topV) < 0:
		return false, fmt.Sprintf("%s is below %s, released (%s); bump to a greater version", cur, top.raw, top.where)
	}
	return true, fmt.Sprintf("%s is above every release (%s, %s)", cur, top.raw, top.where)
}

// versionAt is the version a product's sources hold at a commit, fetched when the checkout
// does not have it — the head of a squashed pull request is no ancestor of its merge.
func versionAt(root string, p *config.Product, sha string) (string, error) {
	if _, ok := gitmeta.Commit(root, sha); !ok {
		_ = exec.Command("git", "-C", root, "fetch", "-q", "--no-tags", "origin", sha).Run()
	}
	svs, err := p.ReadSources(readAt(root, sha))
	if err != nil {
		return "", err
	}
	if len(svs) == 0 {
		return "", fmt.Errorf("the product has no source: its version is the tag alone")
	}
	return config.Agreed(svs)
}

// containsHolds answers `contains: <branch>`: the head of that branch on the remote is an
// ancestor of sha. The branch is fetched first, so a stale or missing remote ref does not
// answer for it.
func containsHolds(root, branch, sha string) (bool, string) {
	ref := "refs/remotes/origin/" + branch
	fetchErr := exec.Command("git", "-C", root, "fetch", "-q", "--no-tags", "origin", "+refs/heads/"+branch+":"+ref).Run()
	head, ok := gitmeta.Commit(root, ref)
	if !ok {
		why := "no origin/" + branch
		if fetchErr != nil {
			why += " (fetch: " + fetchErr.Error() + ")"
		}
		return false, why
	}
	err := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", ref, sha).Run()
	if err == nil {
		return true, fmt.Sprintf("origin/%s (%s) is in this commit", branch, short(head))
	}
	if ee, isExit := err.(*exec.ExitError); !isExit || ee.ExitCode() != 1 {
		return false, fmt.Sprintf("compare origin/%s with %s: %v", branch, short(sha), err)
	}
	out, _ := exec.Command("git", "-C", root, "log", "--format=%h %s", sha+".."+ref).Output()
	missing := strings.Split(strings.TrimSpace(string(out)), "\n")
	shown := missing
	if len(shown) > 5 {
		shown = shown[:5]
	}
	why := fmt.Sprintf("%d commit(s) of origin/%s are not in this commit: %s", len(missing), branch, strings.Join(shown, "; "))
	if len(missing) > len(shown) {
		why += fmt.Sprintf("; and %d more", len(missing)-len(shown))
	}
	return false, why
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
