// @anchors
//   code: LTNGC
//   ref: NGSTI

package mapcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/testsig"
)

// --- a late ingestion: a report stamped with what its run measured (DESIGN-late-ingestion.md) ---
//
// A report ingested after its run was stamped with the tree of the ingestion, and read as if
// it had run against it: a screen changed after the run kept a fresh proof. A manual
// ingestion now places the run in time — `--ran-at`, the report's own timestamp, or the
// report file's modification time —, takes the commit of that time, and confronts each file
// the proof stands on with its content there, by the evidence the map uses. Codes renamed
// since are read through `--alias`, in the report and in that content alike.

// LateRun is how a manual ingestion places its run: when it ran (a time, or the commit it
// measured) and the code aliases (`OLD=NEW`) the report and the files are read through.
type LateRun struct {
	RanAt     string
	Aliases   []string
	AliasFile string
}

// codeAlias renames one code: a unit's (`PHA1`, and every scenario of it) or a scenario's.
type codeAlias struct{ from, to string }

// aliases are the run's aliases, from the flags and the file, the longest first so a scenario
// alias wins over its unit's.
func (l LateRun) aliases() ([]codeAlias, error) {
	lines := append([]string{}, l.Aliases...)
	if l.AliasFile != "" {
		b, err := os.ReadFile(l.AliasFile)
		if err != nil {
			return nil, fmt.Errorf("--alias-file: %w", err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
				lines = append(lines, line)
			}
		}
	}
	var out []codeAlias
	seen := map[string]bool{}
	for _, line := range lines {
		from, to, ok := strings.Cut(line, "=")
		from, to = strings.TrimSpace(from), strings.TrimSpace(to)
		if !ok || from == "" || to == "" || from == to || !isCode(from) || !isCode(to) {
			return nil, fmt.Errorf("--alias %q is not OLD=NEW with two different codes", line)
		}
		if seen[from] {
			return nil, fmt.Errorf("--alias renames %s twice", from)
		}
		seen[from] = true
		out = append(out, codeAlias{from, to})
	}
	sort.SliceStable(out, func(i, j int) bool { return len(out[i].from) > len(out[j].from) })
	return out, nil
}

func isCode(s string) bool {
	for _, r := range s {
		if !(r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return true
}

// remap renames the codes of a text in one pass: each alias's old code where it stands whole —
// no letter or digit glued before or after it —, the longest alias first, so `A=B` and `B=C`
// do not chain.
func remap(b []byte, aliases []codeAlias) []byte {
	if len(aliases) == 0 {
		return b
	}
	word := func(c byte) bool { return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' }
	out := make([]byte, 0, len(b))
	for i := 0; i < len(b); {
		matched := false
		if i == 0 || !word(b[i-1]) {
			for _, a := range aliases {
				end := i + len(a.from)
				if end <= len(b) && string(b[i:end]) == a.from && (end == len(b) || !word(b[end])) {
					out = append(out, a.to...)
					i, matched = end, true
					break
				}
			}
		}
		if !matched {
			out = append(out, b[i])
			i++
		}
	}
	return out
}

// readReport reads a JUnit report through the aliases: its case names and file paths take the
// new codes before anything is matched.
func readReport(path string, aliases []codeAlias) (*testsig.ExecReport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return testsig.ParseJUnitBytes(remap(b, aliases))
}

// runTime is when the run happened, and the commit it measured: `--ran-at` (a time, or a
// commit), else the report's timestamp, else the report file's modification time.
func (l LateRun) runTime(root, report string, rep *testsig.ExecReport) (time.Time, string, error) {
	parse := func(s string) (time.Time, bool) {
		for _, f := range []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
			if t, err := time.ParseInLocation(f, s, time.Local); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}
	if l.RanAt != "" {
		if t, ok := parse(l.RanAt); ok {
			sha, _ := gitmeta.CommitAt(root, t)
			return t, sha, nil
		}
		sha, ok := gitmeta.Commit(root, l.RanAt)
		if !ok {
			return time.Time{}, "", fmt.Errorf("--ran-at %q is neither a time nor a commit", l.RanAt)
		}
		t, _ := gitmeta.CommitTime(root, sha)
		return t, sha, nil
	}
	if rep != nil && rep.RanAt != "" {
		if t, ok := parse(rep.RanAt); ok {
			sha, _ := gitmeta.CommitAt(root, t)
			return t, sha, nil
		}
	}
	fi, err := os.Stat(report)
	if err != nil {
		return time.Time{}, "", err
	}
	sha, _ := gitmeta.CommitAt(root, fi.ModTime())
	return fi.ModTime(), sha, nil
}

// unchangedSince answers, for a file of the map, whether its evidence is the one the run
// measured: a file not modified after the run is; another is confronted with its content at
// the run's commit, read through the aliases. A file the commit lacks, or with no commit to
// read, changed.
func unchangedSince(root string, g *mapx.Graph, cfg *config.Config, at time.Time, sha string, aliases []codeAlias) func(string) bool {
	memo := map[string]bool{}
	return func(id string) bool {
		if v, ok := memo[id]; ok {
			return v
		}
		v := func() bool {
			n := g.Node(id)
			if n == nil {
				return false
			}
			fi, err := os.Stat(filepath.Join(root, filepath.FromSlash(id)))
			if err != nil {
				return false
			}
			if !fi.ModTime().After(at) {
				return true
			}
			if sha == "" {
				return false
			}
			then, ok := gitmeta.AtCommit(root, sha, id)
			if !ok {
				return false
			}
			now := n.EvidenceRev
			if now == "" {
				now = n.Rev
			}
			return scan.EvidenceOf(string(n.Kind), id, remap(then, aliases), cfg).Rev == now
		}()
		memo[id] = v
		return v
	}
}

// aliasNames are the aliases as the signal records them.
func aliasNames(aliases []codeAlias) []string {
	var out []string
	for _, a := range aliases {
		out = append(out, a.from+"="+a.to)
	}
	sort.Strings(out)
	return out
}
