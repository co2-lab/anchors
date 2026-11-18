// @anchors
//   code: PRVLG
//   ref: PRVLS

// Package pushreview keeps the reviews a push carries: the file the repository holds, one
// entry per review, newest first — when, by whom, the verdict, the commits it covers and what
// was validated or failed —, and the commits a push still owes one.
package pushreview

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// The verdicts of a review.
const (
	OK   = "ok"
	Fail = "fail"
)

// Commit is a commit a review covers: its hash, its patch id — the same change after a rebase
// or an amend of its message —, and its subject.
type Commit struct {
	Hash, PatchID, Subject string
}

// Entry is one review as the file holds it: what the machine reads, and the text as written.
type Entry struct {
	At      time.Time
	Verdict string
	By      string
	Commits []Commit
	Raw     string
}

// Review is what a reviewer records: the verdict, who, what was validated (ok) or what failed
// (fail), the observations, and the judgment gates' answers.
type Review struct {
	At       time.Time
	Verdict  string
	By       string
	Range    string
	Commits  []Commit
	Text     string
	Notes    []string
	Judgment []string
}

const header = "# Push reviews\n\n" +
	"<!-- Written by `anchors push-review`, newest first. The pre-push hook reads it: each pushed\n" +
	"     commit needs a review whose verdict is ok. Reviews older than the retention are dropped\n" +
	"     when one is written. -->\n"

var (
	metaRE   = regexp.MustCompile(`^<!-- push-review (.*) -->\s*$`)
	attrRE   = regexp.MustCompile(`(\w+)="([^"]*)"`)
	commitRE = regexp.MustCompile("^- `([0-9a-f]{7,40})`(?: `p:([0-9a-f]+)`)? (.*)$")
)

// Entry renders a review as the file holds it.
func (r Review) Entry() Entry {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s · %s · %s\n", r.At.UTC().Format("2006-01-02 15:04 UTC"), r.Verdict, r.Range)
	fmt.Fprintf(&b, "<!-- push-review at=%q verdict=%q by=%q -->\n\n", r.At.UTC().Format(time.RFC3339), r.Verdict, quoteless(r.By))
	b.WriteString("**Commits**\n\n")
	for _, c := range r.Commits {
		fmt.Fprintf(&b, "- `%s`", c.Hash)
		if c.PatchID != "" {
			fmt.Fprintf(&b, " `p:%s`", c.PatchID)
		}
		fmt.Fprintf(&b, " %s\n", c.Subject)
	}
	if len(r.Judgment) > 0 {
		b.WriteString("\n**Judgment gates**\n\n")
		for _, j := range r.Judgment {
			fmt.Fprintf(&b, "- %s\n", j)
		}
	}
	title := "Validated"
	if r.Verdict == Fail {
		title = "Failed"
	}
	fmt.Fprintf(&b, "\n**%s**\n\n%s\n", title, strings.TrimSpace(r.Text))
	if len(r.Notes) > 0 {
		b.WriteString("\n**Notes**\n\n")
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(n))
		}
	}
	return Entry{At: r.At.UTC().Truncate(time.Second), Verdict: r.Verdict, By: quoteless(r.By), Commits: r.Commits, Raw: b.String()}
}

func quoteless(s string) string { return strings.ReplaceAll(strings.TrimSpace(s), `"`, "'") }

// Parse reads the file's entries in its order, newest first. A section the machine line does
// not open is no review, and is left out.
func Parse(content string) []Entry {
	var out []Entry
	var cur []string
	flush := func() {
		if len(cur) < 2 {
			cur = nil
			return
		}
		m := metaRE.FindStringSubmatch(cur[1])
		if m == nil {
			cur = nil
			return
		}
		e := Entry{Raw: strings.TrimRight(strings.Join(cur, "\n"), "\n") + "\n"}
		for _, a := range attrRE.FindAllStringSubmatch(m[1], -1) {
			switch a[1] {
			case "at":
				e.At, _ = time.Parse(time.RFC3339, a[2])
			case "verdict":
				e.Verdict = a[2]
			case "by":
				e.By = a[2]
			}
		}
		for _, l := range cur[2:] {
			if c := commitRE.FindStringSubmatch(l); c != nil {
				e.Commits = append(e.Commits, Commit{Hash: c[1], PatchID: c[2], Subject: c[3]})
			}
		}
		out = append(out, e)
		cur = nil
	}
	for _, l := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(l, "## ") {
			flush()
			cur = []string{l}
			continue
		}
		if cur != nil {
			cur = append(cur, l)
		}
	}
	flush()
	return out
}

// Render writes the file: its header, then the entries in the order given.
func Render(entries []Entry) string {
	var b strings.Builder
	b.WriteString(header)
	for _, e := range entries {
		b.WriteString("\n")
		b.WriteString(e.Raw)
	}
	return b.String()
}

// Keep drops the entries older than days before now; 0 keeps them all. An entry whose date
// nobody can read is kept: dropping it would lose a review by accident.
func Keep(entries []Entry, days int, now time.Time) []Entry {
	if days <= 0 {
		return entries
	}
	limit := now.AddDate(0, 0, -days)
	var out []Entry
	for _, e := range entries {
		if e.At.IsZero() || !e.At.Before(limit) {
			out = append(out, e)
		}
	}
	return out
}

// Covering is the newest entry that names the commit — by its hash, or by its patch id when
// a rebase gave it another hash; false when none does.
func Covering(entries []Entry, c Commit) (Entry, bool) {
	for _, e := range entries {
		for _, x := range e.Commits {
			if sameHash(x.Hash, c.Hash) || (x.PatchID != "" && c.PatchID != "" && x.PatchID == c.PatchID) {
				return e, true
			}
		}
	}
	return Entry{}, false
}

func sameHash(a, b string) bool {
	if len(a) < 7 || len(b) < 7 {
		return false
	}
	if len(a) > len(b) {
		a, b = b, a
	}
	return strings.HasPrefix(b, a)
}

// Exempt says whether a commit owes no review: it changes nothing but the reviews' file and
// the map — the evidence of a review, which a review of itself would never end.
func Exempt(files []string, reviewFile, mapFile string) bool {
	if len(files) == 0 {
		return false
	}
	for _, f := range files {
		if f != reviewFile && f != mapFile {
			return false
		}
	}
	return true
}

// Commits are the commits rev-list gives for the arguments (`a..b`, or `b --not --remotes`),
// oldest first, each with its patch id and subject.
func Commits(root string, revs ...string) ([]Commit, error) {
	args := append([]string{"-C", root, "rev-list", "--reverse", "--no-merges"}, revs...)
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("git rev-list %s: %w", strings.Join(revs, " "), err)
	}
	var cs []Commit
	for _, h := range strings.Fields(string(out)) {
		subj, _ := exec.Command("git", "-C", root, "show", "-s", "--format=%s", h).Output()
		cs = append(cs, Commit{Hash: h, PatchID: patchID(root, h), Subject: strings.TrimSpace(string(subj))})
	}
	return cs, nil
}

// Files are the paths a commit changes, relative to the root.
func Files(root, hash string) []string {
	out, err := exec.Command("git", "-C", root, "diff-tree", "--no-commit-id", "--name-only", "-r", "--root", hash).Output()
	if err != nil {
		return nil
	}
	return strings.Fields(string(out))
}

// patchID is the commit's stable patch id, short; empty when git gives none (an empty commit).
func patchID(root, hash string) string {
	show := exec.Command("git", "-C", root, "show", "--format=", hash)
	pid := exec.Command("git", "-C", root, "patch-id", "--stable")
	pipe, err := show.StdoutPipe()
	if err != nil {
		return ""
	}
	pid.Stdin = pipe
	if err := show.Start(); err != nil {
		return ""
	}
	out, err := pid.Output()
	_ = show.Wait()
	if err != nil {
		return ""
	}
	f := strings.Fields(string(out))
	if len(f) == 0 || len(f[0]) < 12 {
		return ""
	}
	return f[0][:12]
}
