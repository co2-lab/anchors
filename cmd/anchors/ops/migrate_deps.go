// @anchors
//   code: MDCMG
//   ref: MGCMM

package ops

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// The spec declares no dependency (DESIGN-dependencies-out-of-the-spec.md): it precedes the
// code, and the files a unit imports are declared where the import is, by `@dep:`. This step
// takes the Dependencies tables out of the specs and the `DEPn` out of what the rules use.
// It is no format change — an old binary reads the result —, so it runs whenever a table is
// left, and does nothing once none is: idempotent like the rest of the migration.

// SpecDepsReport is what the step did to one spec, and what it leaves to the author.
type SpecDepsReport struct {
	File string
	// Sections are the titles of the tables removed; External, the rows of them that named
	// no file of the project — an API, a table, a queue —, to be flagged in the code where
	// it is called (`@dep[<kind>]:`).
	Sections []string
	External []string
	// Dropped is how many `DEPn` left the rules' uses; Emptied, the rules whose uses held
	// nothing else, which now say nothing they read.
	Dropped int
	Emptied []string
	// Mentions are the lines that still cite a `DEPn` — an origin in a data contract, a
	// sentence —, which the step does not rewrite: what the datum comes from is the author's.
	Mentions []int
}

var (
	// depTableRowRE is a row of a Dependencies table: it opens with its `DEPn`.
	depTableRowRE = regexp.MustCompile("^\\s*\\|\\s*`?DEP\\d+`?\\s*\\|")
	// depItemRE is a use that names a `DEPn` — alone, or with a member (`DEP1.getProfile`).
	depItemRE      = regexp.MustCompile("^`?DEP\\d+(?:[.:][^`]*)?`?$")
	depMentionRE   = regexp.MustCompile(`\bDEP\d+\b`)
	headingLineRE  = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
	ruleCodeCellRE = regexp.MustCompile("^`?[A-Z0-9]{3,}-[A-Z]\\d{2}")
)

// migrateSpecDependencies runs the step over every spec of the project.
func migrateSpecDependencies(absRoot string, cfg *config.Config, dryRun bool) ([]SpecDepsReport, error) {
	files, err := versionedFiles(absRoot)
	if err != nil {
		return nil, err
	}
	var out []SpecDepsReport
	type rewrite struct{ rel, to string }
	var written []rewrite
	defer func() {
		// Taking the table out proves nothing new: each spec keeps what held at its revision
		// before — its scenarios' proofs, its stamps —, as a repair by `check --fix` does.
		if len(written) == 0 {
			return
		}
		mapPath := filepath.Join(absRoot, mapx.DefaultPath)
		if _, err := os.Stat(mapPath); err != nil {
			return
		}
		today := gitmeta.Today()
		_ = mapx.Update(mapPath, func(g *mapx.Graph) error {
			for _, w := range written {
				g.KeepFixedEvidence(w.rel, w.to, "the spec's Dependencies table moved to the code's @dep flags", today, false)
			}
			return nil
		})
	}()
	for _, f := range files {
		if !strings.HasSuffix(f, ".spec.md") {
			continue
		}
		p := filepath.Join(absRoot, filepath.FromSlash(f))
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		next, r := withoutSpecDependencies(string(b), absRoot, path.Dir(f), cfg)
		if next == string(b) {
			continue
		}
		r.File = f
		out = append(out, r)
		if !dryRun {
			if err := writeKeepingMode(p, next); err != nil {
				return out, err
			}
			written = append(written, rewrite{f, scan.ShortHash([]byte(next))})
		}
	}
	return out, nil
}

// withoutSpecDependencies is a spec's text without its Dependencies tables and without the
// `DEPn` its rules use, with what was done and what is left.
func withoutSpecDependencies(content, absRoot, specDir string, cfg *config.Config) (string, SpecDepsReport) {
	var r SpecDepsReport
	lines := strings.Split(content, "\n")

	// 1. The sections whose table rows open with a `DEPn` — under any title.
	type span struct{ start, end int }
	var cut []span
	for i := 0; i < len(lines); i++ {
		m := headingLineRE.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		level := len(m[1])
		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			if h := headingLineRE.FindStringSubmatch(lines[j]); h != nil && len(h[1]) <= level {
				end = j
				break
			}
		}
		rows := 0
		for j := i + 1; j < end; j++ {
			if headingLineRE.MatchString(lines[j]) {
				break // a subsection is its own section
			}
			if depTableRowRE.MatchString(lines[j]) {
				rows++
				if !namesProjectFile(lines[j], absRoot, specDir) {
					r.External = append(r.External, strings.TrimSpace(lines[j]))
				}
			}
		}
		if rows > 0 {
			cut = append(cut, span{i, end})
			r.Sections = append(r.Sections, m[2])
			i = end - 1
		}
	}
	for k := len(cut) - 1; k >= 0; k-- {
		at := cut[k].start
		lines = append(lines[:at], lines[cut[k].end:]...)
		// The cut leaves no blank line doubled where the section was; the rest of the spec
		// keeps its own spacing.
		for at > 0 && at < len(lines) && strings.TrimSpace(lines[at-1]) == "" && strings.TrimSpace(lines[at]) == "" {
			lines = append(lines[:at], lines[at+1:]...)
		}
	}

	// 2. The `DEPn` of what each rule uses.
	uses := useSectionTitles(cfg)
	in := false
	for i, l := range lines {
		if m := headingLineRE.FindStringSubmatch(l); m != nil {
			in = uses[strings.ToLower(strings.TrimSpace(noteRE.ReplaceAllString(m[2], "")))] || uses[strings.ToLower(m[2])]
			continue
		}
		t := strings.TrimSpace(l)
		if !in || !strings.HasPrefix(t, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(t, "|"), "|")
		if len(cells) < 2 || !ruleCodeCellRE.MatchString(strings.TrimSpace(cells[0])) {
			continue
		}
		changed := false
		for c := 1; c < len(cells); c++ {
			items := strings.Split(cells[c], ",")
			var keep []string
			cellChanged := false
			for _, it := range items {
				if depItemRE.MatchString(strings.TrimSpace(it)) {
					r.Dropped++
					cellChanged = true
					continue
				}
				keep = append(keep, strings.TrimSpace(it))
			}
			if cellChanged {
				cells[c] = " " + strings.Join(keep, ", ") + " "
				changed = true
			}
		}
		if !changed {
			continue
		}
		empty := true
		for _, c := range cells[1:] {
			if strings.TrimSpace(c) != "" {
				empty = false
			}
		}
		if empty {
			r.Emptied = append(r.Emptied, strings.Trim(strings.TrimSpace(cells[0]), "`"))
		}
		lines[i] = "|" + strings.Join(cells, "|") + "|"
	}

	// 3. What still cites a `DEPn`.
	for i, l := range lines {
		if depMentionRE.MatchString(l) {
			r.Mentions = append(r.Mentions, i+1)
		}
	}
	if len(r.Sections) == 0 && r.Dropped == 0 {
		return content, SpecDepsReport{} // nothing taken out: the spec is left exactly as it was
	}
	return strings.Join(lines, "\n"), r
}

// noteRE is a title's note in parentheses.
var noteRE = regexp.MustCompile(`\s*\([^)]*\)\s*$`)

// useSectionTitles are the titles of the sections where a rule says what it uses: Rule uses,
// Validations and Presentation validations, in every language and as the project names them.
func useSectionTitles(cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	for _, k := range []string{"rule_uses", "validations", "presentation_validations"} {
		for _, t := range i18n.AllTranslations("section.title." + k) {
			out[strings.ToLower(t)] = true
		}
		for _, t := range cfg.SectionTitlesFor(strings.ReplaceAll(k, "_", "-")) {
			out[strings.ToLower(t)] = true
		}
	}
	return out
}

// namesProjectFile says whether a Dependencies row names a file of the project: a cell after
// its code holds a path — or several, separated by commas — and one of them exists, from the
// root or from the spec's folder.
func namesProjectFile(row, absRoot, specDir string) bool {
	cells := strings.Split(strings.Trim(strings.TrimSpace(row), "|"), "|")
	for _, c := range cells[1:] {
		for _, f := range strings.Split(c, ",") {
			f = strings.Trim(strings.TrimSpace(f), "`")
			if f == "" || strings.Contains(f, " ") {
				continue
			}
			if fileExists(absRoot, f) || fileExists(absRoot, path.Join(specDir, f)) {
				return true
			}
		}
	}
	return false
}

func fileExists(absRoot, rel string) bool {
	_, err := os.Stat(filepath.Join(absRoot, filepath.FromSlash(rel)))
	return err == nil
}

// printSpecDepsReport says what the step did, and what each spec leaves to its author.
func printSpecDepsReport(reports []SpecDepsReport, dryRun bool) {
	if len(reports) == 0 {
		return
	}
	verb := "removed"
	if dryRun {
		verb = "would remove"
	}
	tables, dropped := 0, 0
	for _, r := range reports {
		tables += len(r.Sections)
		dropped += r.Dropped
	}
	fmt.Printf("✓ specs: %s %d Dependencies table(s) from %d spec(s), and %d `DEPn` from the rules' uses\n", verb, tables, len(reports), dropped)
	fmt.Println("  a spec precedes the code: the files a unit imports are flagged in the code (`@dep:`)")
	fmt.Println("  — `anchors check --fix` writes the flags on the imports.")
	sort.Slice(reports, func(i, j int) bool { return reports[i].File < reports[j].File })
	for _, r := range reports {
		if len(r.External) == 0 && len(r.Emptied) == 0 && len(r.Mentions) == 0 {
			continue
		}
		fmt.Printf("  %s:\n", r.File)
		for _, e := range r.External {
			fmt.Printf("    names no file of the project — flag it where the code calls it (`@dep[<kind>]: <name>`): %s\n", e)
		}
		if len(r.Emptied) > 0 {
			fmt.Printf("    say nothing they read now — name the fields: %s\n", strings.Join(r.Emptied, ", "))
		}
		if len(r.Mentions) > 0 {
			fmt.Printf("    still cite a `DEPn` (a data origin, a sentence) on line(s) %s — say what the datum comes from\n", joinLines(r.Mentions))
		}
	}
}

func joinLines(ns []int) string {
	s := make([]string, len(ns))
	for i, n := range ns {
		s[i] = fmt.Sprint(n)
	}
	return strings.Join(s, ", ")
}
