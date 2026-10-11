// @anchors
//   code: ATCFF
//   ref: AFTSF

package quality

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/testlist"
	"github.com/co2-lab/anchors/internal/testsig"
	"github.com/spf13/cobra"
)

// --- which tests to run for what changed: the short cycle without the whole suite ---
//
// The gates and the CI run everything; whoever is working wants the tests a change can break,
// now. A change is read at the finest grain the project declares: a spec's rule, a code
// region `#region [CODE]` — and then the tests that carry that code in their title —, or a
// whole file — and then every test whose evidence closure holds it. The list is printed the
// way the project's runner takes it.

func newAffectedTestsCmd() *cobra.Command {
	var root, mapPath, format string
	cmd := &cobra.Command{
		Use:   "affected-tests [file...]",
		Short: "List the tests a change can break — by rule and region where the project declares them",
		Long: `Lists the tests to run for what changed: the given files, or, with none, every file that
differs from the last commit.

  a spec's rule, or a code region '#region [CODE]', changed  → the tests whose title carries CODE
  a file changed outside its regions                          → every test whose evidence closure holds it
  a test file changed                                         → all of its tests

The formats:
  lines   one test title per line (default)
  regex   one anchored alternation of the titles: go test -run "$(anchors affected-tests --format regex)"
  files   the test files, one per line
  json    the tests with file, line and title

With no 'dialect.tests' declared the titles cannot be read, and the test files are printed.
The gate run stays whole for the CI; this serves the short cycle of whoever is working.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w%s", config.DefaultFile, err, config.InitHint(err))
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}
			files := args
			if len(files) == 0 {
				dirty, known := gitmeta.DirtyFiles(absRoot)
				if !known {
					return fmt.Errorf("no file given and git cannot say what changed — name the files")
				}
				for f := range dirty {
					files = append(files, f)
				}
			}
			for i, f := range files {
				files[i] = common.RelTo(absRoot, f)
			}
			sort.Strings(files)
			tests, declared, err := gate.ProjectTests(absRoot, g, cfg)
			if err != nil {
				return fmt.Errorf("read the tests: %w", err)
			}
			a := affectedBy(absRoot, g, cfg, files, func(rel string) (string, bool) { return gitmeta.AtHead(absRoot, rel) })
			return printAffected(cmd.OutOrStdout(), cmd.ErrOrStderr(), a, tests, declared, format)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&format, "format", "lines", "lines | regex | files | json")
	return common.TakesFiles(cmd)
}

// affected is what a change reaches: the test files run whole, and the codes whose tests run.
type affected struct {
	whole map[string]bool // test files: all of their tests
	codes map[string]bool // scenario codes: the tests whose title carries one
}

// affectedBy reads each changed file against its last committed version (atHead) and says what
// it reaches. A file a freeze covers reaches nothing: its tests are not run until it is thawed.
func affectedBy(root string, g *mapx.Graph, cfg *config.Config, files []string, atHead func(string) (string, bool)) affected {
	a := affected{whole: map[string]bool{}, codes: map[string]bool{}}
	var wholeFiles []string
	for _, f := range files {
		n := g.Node(f)
		if n != nil && n.Frozen != "" {
			continue
		}
		if n != nil && n.Kind == mapx.KindTest {
			a.whole[f] = true
			continue
		}
		now, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
		before, had := atHead(f)
		if err != nil || !had {
			wholeFiles = append(wholeFiles, f) // new, or gone: whole
			continue
		}
		if n != nil && n.Kind == mapx.KindSpec {
			for _, c := range changedRuleCodes(f, before, string(now), cfg) {
				a.codes[c] = true
			}
			continue
		}
		codes, outside := changedRegions(before, string(now))
		for _, c := range codes {
			a.codes[c] = true
		}
		if outside {
			wholeFiles = append(wholeFiles, f)
		}
	}
	if len(wholeFiles) > 0 {
		want := map[string]bool{}
		for _, f := range wholeFiles {
			want[f] = true
		}
		for _, n := range g.Nodes {
			if n.Kind != mapx.KindTest || n.Frozen != "" || a.whole[n.ID] {
				continue
			}
			for id := range g.EvidenceClosure(n.ID) {
				if want[id] {
					a.whole[n.ID] = true
					break
				}
			}
		}
	}
	return a
}

// changedRuleCodes are the rules of a spec whose definition differs between two versions,
// added and removed ones included.
func changedRuleCodes(rel, before, now string, cfg *config.Config) []string {
	b := scan.EvidenceOf("spec", rel, []byte(before), cfg).Rules
	n := scan.EvidenceOf("spec", rel, []byte(now), cfg).Rules
	var out []string
	for c, r := range n {
		if b[c] != r {
			out = append(out, c)
		}
	}
	for c := range b {
		if _, ok := n[c]; !ok {
			out = append(out, c)
		}
	}
	sort.Strings(out)
	return out
}

// changedRegions are the codes of the regions that differ between two versions of a file —
// added and removed ones included —, and whether anything changed outside every region.
func changedRegions(before, now string) (codes []string, outside bool) {
	rb, _ := scan.Regioes(before)
	rn, _ := scan.Regioes(now)
	revB, revN := map[string]string{}, map[string]string{}
	for _, r := range rb {
		revB[r.Code] += r.Rev
	}
	for _, r := range rn {
		revN[r.Code] += r.Rev
	}
	for c, r := range revN {
		if revB[c] != r {
			codes = append(codes, c)
		}
	}
	for c := range revB {
		if _, ok := revN[c]; !ok {
			codes = append(codes, c)
		}
	}
	sort.Strings(codes)
	return codes, outsideRegions(before, rb) != outsideRegions(now, rn)
}

// outsideRegions is a file's text with every region taken out.
func outsideRegions(content string, regions []scan.Regiao) string {
	lines := strings.Split(content, "\n")
	in := make([]bool, len(lines)+1)
	for _, r := range regions {
		for l := r.Start; l <= r.End && l <= len(lines); l++ {
			in[l] = true
		}
	}
	var b strings.Builder
	for i, l := range lines {
		if !in[i+1] {
			b.WriteString(l + "\n")
		}
	}
	return b.String()
}

// selectAffected are the tests a change reaches, by file and line: every test of a file run
// whole, and every test whose title carries a changed code.
func selectAffected(a affected, tests []testlist.Test) []testlist.Test {
	var out []testlist.Test
	for _, t := range tests {
		take := a.whole[t.File]
		if !take {
			for _, c := range testsig.CodesInCase(t.Title) {
				if a.codes[c] {
					take = true
					break
				}
			}
		}
		if take {
			out = append(out, t)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		return out[i].Line < out[j].Line
	})
	return out
}

// printAffected prints the tests in the runner's format; with no way to read the titles, the
// test files run whole are what can be said, and stderr says why.
func printAffected(w, errw io.Writer, a affected, tests []testlist.Test, declared bool, format string) error {
	if !declared {
		var files []string
		for f := range a.whole {
			files = append(files, f)
		}
		sort.Strings(files)
		fmt.Fprintln(errw, "no `dialect.tests` declared: the test titles cannot be read, so the tests by code are left out and the test files are printed")
		for _, f := range files {
			fmt.Fprintln(w, f)
		}
		return nil
	}
	sel := selectAffected(a, tests)
	switch format {
	case "lines":
		for _, t := range sel {
			fmt.Fprintln(w, t.Title)
		}
	case "files":
		seen := map[string]bool{}
		for _, t := range sel {
			if !seen[t.File] {
				seen[t.File] = true
				fmt.Fprintln(w, t.File)
			}
		}
	case "regex":
		if len(sel) == 0 {
			return nil
		}
		var alts []string
		seen := map[string]bool{}
		for _, t := range sel {
			if !seen[t.Title] {
				seen[t.Title] = true
				alts = append(alts, regexp.QuoteMeta(t.Title))
			}
		}
		fmt.Fprintf(w, "^(%s)$\n", strings.Join(alts, "|"))
	case "json":
		if sel == nil {
			sel = []testlist.Test{}
		}
		b, err := json.MarshalIndent(sel, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(b))
	default:
		return fmt.Errorf("unknown format %q — lines, regex, files or json", format)
	}
	return nil
}
