// @anchors
//   code: MGCMA
//   ref: MGCMM

package ops

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/flowx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
	"github.com/co2-lab/anchors/internal/recode"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// O comando que a mensagem de erro do formato PROMETE.
//
// Um erro que diz "rode `anchors migrate`" e não tem o comando é pior que nenhum erro:
// quem lê tenta, falha, e passa a desconfiar da próxima mensagem.
func newMigrateCmd() *cobra.Command {
	var root string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Bring the Anchors files up to the current format",
		Long: `Rewrites the ` + "`anchors.graph.yaml`" + ` and the ` + "`anchors.yaml`" + ` into the format
that this binary reads.

The format goes up when a change makes the file UNREADABLE for the previous
version — a renamed key, a value that changed shape. It does not go up for a new
optional field: an old binary simply ignores it.

The migration runs ONCE and leaves the project ready to commit. It is idempotent:
running it again changes nothing.

    anchors migrate              # migrates
    anchors migrate --dry-run    # says what it would do, without writing`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true

			alvos := []string{
				filepath.Join(absRoot, mapx.DefaultPath),
				filepath.Join(absRoot, config.DefaultFile),
			}
			// The project's format BEFORE this run: the lowest of the two files. Code letters
			// live in the project's own files, not in these two, and are rewritten once, when
			// the project crosses the step that renamed them.
			from := mapx.FormatoAtual
			formatOf := map[string]int{}
			for _, alvo := range alvos {
				if f, err := migra.FormatOf(alvo); err == nil {
					formatOf[alvo] = f
					if f < from {
						from = f
					}
				}
			}

			var mudou bool
			for _, alvo := range alvos {
				r, err := migra.MigrateFile(alvo, mapx.FormatoAtual, dryRun)
				if err != nil {
					// Arquivo ausente não é erro: um projeto pode não ter mapa ainda, e
					// recusar a migração inteira por causa disso deixaria o outro arquivo
					// no formato velho.
					fmt.Printf("· %s: %v\n", filepath.Base(alvo), err)
					continue
				}
				if !r.Changed {
					fmt.Printf("· %s is already in format %d\n", filepath.Base(alvo), r.To)
					continue
				}
				mudou = true
				verbo := "migrated"
				if dryRun {
					verbo = "would be migrated"
				}
				fmt.Printf("✓ %s %s: format %d → %d\n", filepath.Base(alvo), verbo, r.De, r.To)

				// AS CHAVES, uma a uma e com a contagem. Um "migrado" seco não diz o que
				// mudou, e quem revisa o diff precisa saber o que esperar antes de abri-lo.
				chaves := make([]string, 0, len(r.Replaced))
				for k := range r.Replaced {
					chaves = append(chaves, k)
				}
				sort.Strings(chaves)
				for _, k := range chaves {
					fmt.Printf("    %s  (%d occurrence(s))\n", k, r.Replaced[k])
				}
			}

			if renames := migra.LetterRenamesBetween(from, mapx.FormatoAtual); len(renames) > 0 {
				changed, err := migrateCodeLetters(absRoot, renames, dryRun)
				if err != nil {
					return err
				}
				mudou = mudou || changed
				// The pipelines Anchors installed read the phase letter in their scripts. The
				// ones nobody edited are brought up by the doctor, which already knows how to
				// tell them from the customized ones; the migration does not rewrite a
				// pipeline — that is a copy of a template, not a project file.
				fmt.Println("· pipelines: the installed ones read the old letters — `anchors doctor --fix` updates the")
				fmt.Println("  ones nobody edited; a customized one needs its `-F[0-9]` patterns changed to `-W[0-9]` by hand")
			}

			// Format 7: every file a code of its own, of five characters.
			if from < 7 && mapx.FormatoAtual >= 7 {
				changed, err := migrateToFileCodes(absRoot, dryRun)
				if err != nil {
					// The files go back to the format they had, so the next run crosses the
					// step again — and codes only what this one left.
					for alvo, f := range formatOf {
						if f < 7 {
							_ = rewriteText(alvo, func(s string) string {
								return regexp.MustCompile(`(?m)^version:\s*\d+\s*$`).ReplaceAllString(s, fmt.Sprintf("version: %d", f))
							})
						}
					}
					return err
				}
				mudou = mudou || changed
			}

			if !mudou {
				return nil
			}
			fmt.Println()
			if dryRun {
				fmt.Println("  nothing was written — run without `--dry-run` to apply.")
				return nil
			}
			// O COMMIT é de quem rodou, e dizer isso importa: o mapa é versionado, e uma
			// migração que fica só na máquina faz o próximo agente reencontrar o formato
			// velho — e migrar de novo, gerando o mesmo diff outra vez.
			fmt.Println("  the files are VERSIONED: commit the migration so that the")
			fmt.Println("  other agents do not redo it.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "say what it would do, without writing")
	return cmd
}

// migrateCodeLetters rewrites, in the project's versioned files, the codes whose letter a
// migration step renamed. The units are found by what they are: a plan is a file of a
// `kind: plan` layer, a flow a `.flow.md`, an action a `.action.md` — each by the code its
// header declares. Only those units' codes change: the same letter in a spec means
// something else. Every versioned text file is read, because a phase is cited from specs'
// `needs:`, a result from flows, and both from docs.
func migrateCodeLetters(absRoot string, renames []migra.LetterRename, dryRun bool) (bool, error) {
	files, err := versionedFiles(absRoot)
	if err != nil {
		return false, err
	}
	kindOf := map[string]string{}
	var plans map[string]bool
	if cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile)); err == nil {
		if walked, err := scan.Walk(absRoot, cfg); err == nil {
			plans = map[string]bool{}
			for _, f := range walked {
				if f.Kind == string(mapx.KindPlan) {
					plans[f.Path] = true
				}
			}
		}
	}
	contents := map[string]string{}
	for _, rel := range files {
		b, err := os.ReadFile(filepath.Join(absRoot, rel))
		if err != nil || bytes.IndexByte(b, 0) >= 0 {
			continue // unreadable or binary: nothing to rewrite
		}
		contents[rel] = string(b)
		kind := ""
		switch {
		case strings.HasSuffix(rel, flowx.FlowSuffix):
			kind = "flow"
		case strings.HasSuffix(rel, flowx.ActionSuffix):
			kind = "action"
		case plans[rel]:
			kind = "plan"
		}
		if code := scan.HeaderCodeOf(contents[rel]); kind != "" && code != "" {
			kindOf[code] = kind
		}
	}
	if plans == nil {
		fmt.Println("· plans: the configuration did not load, so plan phases were not looked for")
	}
	if len(kindOf) == 0 {
		return false, nil
	}
	changed := false
	for _, rel := range files {
		text, ok := contents[rel]
		if !ok {
			continue
		}
		out, counts := migra.RewriteCodeLetters(text, kindOf, renames)
		if len(counts) == 0 {
			continue
		}
		changed = true
		verb := "rewritten"
		if dryRun {
			verb = "would be rewritten"
		}
		fmt.Printf("✓ %s %s:\n", rel, verb)
		for _, k := range migra.SortedCounts(counts) {
			fmt.Printf("    %s  (%d)\n", k, counts[k])
		}
		if dryRun {
			continue
		}
		info, err := os.Stat(filepath.Join(absRoot, rel))
		if err != nil {
			return changed, err
		}
		if err := os.WriteFile(filepath.Join(absRoot, rel), []byte(out), info.Mode()); err != nil {
			return changed, err
		}
	}
	return changed, nil
}

// versionedFiles lists the project's files as git knows them — tracked and untracked but
// not ignored —, relative to the root. Outside a repository, every file under the root
// except hidden directories.
func versionedFiles(absRoot string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--cached", "--others", "--exclude-standard", "-z")
	cmd.Dir = absRoot
	if out, err := cmd.Output(); err == nil {
		var files []string
		for _, f := range strings.Split(string(out), "\x00") {
			if f != "" {
				files = append(files, filepath.ToSlash(f))
			}
		}
		return files, nil
	}
	var files []string
	err := filepath.WalkDir(absRoot, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && p != absRoot && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(absRoot, p)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	return files, err
}

// RenamesFile is where the migration records each code it renamed: outside the repository —
// the issues, the pull requests, the old commit messages — the old code stays, and whoever
// searches the history needs the way from it to the new one.
const RenamesFile = "anchors.renames.yaml"

// migrateToFileCodes brings a project to format 7: its four-character codes widened to five,
// and a code of its own on every governed file that can carry one. What every file was
// measured at is carried to its new revision, so the migration proves nothing new and loses
// nothing proven.
func migrateToFileCodes(absRoot string, dryRun bool) (bool, error) {
	// Without a configuration there is no governed file to give a code to; a configuration
	// this binary cannot read is said, and the rest of the migration still goes on.
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		fmt.Printf("· file codes: %v\n", firstLine(err.Error()))
		return false, nil
	}
	// Both lengths are read while the project crosses from one to the other.
	prevLengths := config.CodeLengths
	config.SetCodeLengths([]int{4, 5})
	defer config.SetCodeLengths(prevLengths)

	before, err := scan.Walk(absRoot, cfg)
	if err != nil {
		return false, err
	}
	oldRev := map[string]string{}
	taken := map[string]bool{}
	for _, f := range before {
		oldRev[f.Path] = f.Rev
		if f.HeaderCode != "" {
			taken[f.HeaderCode] = true
		}
		for _, c := range f.Codes {
			if i := strings.Index(c, "-"); i > 0 {
				taken[c[:i]] = true
			}
		}
	}
	changed := false
	touched := map[string]bool{}
	renamedPath := map[string]string{}

	// 1. The four-character codes, widened.
	type pair struct{ old, new string }
	var pairs []pair
	seen := map[string]bool{}
	for _, f := range before {
		c := f.HeaderCode
		if len(c) != 4 || len(f.HeaderRefs) > 0 || seen[c] {
			continue
		}
		seen[c] = true
		n := migra.WidenedCode(c, unitName(f.Path), taken)
		taken[n] = true
		pairs = append(pairs, pair{c, n})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].old < pairs[j].old })
	rcfg := *cfg
	rc := config.Recode{}
	if cfg.Recode != nil {
		rc = *cfg.Recode
	}
	// Every file whose NAME carries the code is renamed with it — a baseline, a capture
	// flow —, declared in `recode:` or not: left behind, it would no longer be found.
	rc.FilePatterns = append(append([]string(nil), rc.FilePatterns...), "**/*{{code}}*")
	rcfg.Recode = &rc
	for _, p := range pairs {
		plan, err := recode.BuildPlan(absRoot, &rcfg, p.old, p.new)
		if err != nil {
			fmt.Printf("· %s → %s: %v\n", p.old, p.new, err)
			continue
		}
		changed = true
		fmt.Printf("✓ code %s → %s: %d occurrence(s) in %d file(s), %d file(s) renamed\n",
			p.old, p.new, plan.Total+plan.TestIDs, len(plan.Files), len(plan.Renames))
		for _, fc := range plan.Files {
			touched[fc.Path] = true
		}
		for _, r := range plan.Renames {
			renamedPath[r.From] = r.To
		}
		if !dryRun {
			if _, err := plan.Apply(absRoot); err != nil {
				return changed, err
			}
		}
	}

	// 2. The configuration and the map speak the new codes, and the renamed paths.
	if !dryRun {
		if err := rewriteText(filepath.Join(absRoot, config.DefaultFile), func(s string) string {
			for _, p := range pairs {
				s, _ = recode.Rewrite(s, p.old, p.new)
			}
			return regexp.MustCompile(`(?m)^code_lengths:\s*\[[^\]]*\]`).ReplaceAllString(s, "code_lengths: [5]")
		}); err != nil {
			return changed, err
		}
		if err := rewriteText(filepath.Join(absRoot, mapx.DefaultPath), func(s string) string {
			for from, to := range renamedPath {
				s = strings.ReplaceAll(s, from, to)
			}
			for _, p := range pairs {
				s, _ = recode.Rewrite(s, p.old, p.new)
			}
			return s
		}); err != nil {
			return changed, err
		}
	}
	config.SetCodeLengths([]int{5})

	// 3. A code of its own on every governed file that can carry one.
	var after []scan.File
	if dryRun {
		after = before
	} else if after, err = scan.Walk(absRoot, cfg); err != nil {
		return changed, err
	}
	// A file that carried its spec's code — written before the code was the file's — keeps
	// naming the unit by a `ref:`, and gets a code of its own like any other.
	owners := map[string]int{}
	specOwns := map[string]bool{}
	for _, f := range after {
		if f.HeaderCode != "" && len(f.HeaderRefs) == 0 {
			owners[f.HeaderCode]++
			if f.Kind == "spec" {
				specOwns[f.HeaderCode] = true
			}
		}
	}
	given := 0
	for _, f := range after {
		if f.HeaderCode != "" && len(f.HeaderRefs) == 0 && f.Kind != "spec" &&
			owners[f.HeaderCode] > 1 && specOwns[f.HeaderCode] && !f.Upstream {
			abs := filepath.Join(absRoot, filepath.FromSlash(f.Path))
			b, err := os.ReadFile(abs)
			if err != nil {
				continue
			}
			c := migra.FileCode(f, taken)
			out := migra.AsRefWithOwnCode(string(b), f.HeaderCode, c)
			if out == string(b) {
				continue
			}
			taken[c] = true
			given++
			touched[f.Path] = true
			if dryRun {
				continue
			}
			info, err := os.Stat(abs)
			if err != nil {
				return true, err
			}
			if err := os.WriteFile(abs, []byte(out), info.Mode()); err != nil {
				return true, err
			}
			continue
		}
		if f.HeaderCode != "" || f.Upstream {
			continue
		}
		abs := filepath.Join(absRoot, filepath.FromSlash(f.Path))
		b, err := os.ReadFile(abs)
		if err != nil || !migra.CanCarryCode(f.Path, b) {
			continue
		}
		c := migra.FileCode(f, taken)
		taken[c] = true
		given++
		touched[f.Path] = true
		if dryRun {
			continue
		}
		info, err := os.Stat(abs)
		if err != nil {
			return true, err
		}
		if err := os.WriteFile(abs, []byte(migra.WithHeaderCode(string(b), f.Path, c)), info.Mode()); err != nil {
			return true, err
		}
	}
	if given > 0 {
		changed = true
		verb := "get"
		if dryRun {
			verb = "would get"
		}
		fmt.Printf("✓ %d file(s) %s a code of its own (`code:` in the header)\n", given, verb)
	}
	if dryRun || !changed {
		return changed, nil
	}

	// 4. What each touched file was measured at goes with it to its new revision — when the
	// map had it measured at the content the migration found.
	g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
	if err == nil {
		var rels []string
		for p := range touched {
			if to, ok := renamedPath[p]; ok {
				p = to
			}
			rels = append(rels, p)
		}
		sort.Strings(rels)
		now, err := scan.ScanPaths(absRoot, cfg, rels)
		if err == nil {
			oldOf := map[string]string{}
			for p, r := range oldRev {
				if to, ok := renamedPath[p]; ok {
					p = to
				}
				oldOf[p] = r
			}
			carried := 0
			for _, f := range now {
				n := g.Node(f.Path)
				if n == nil || n.Rev != oldOf[f.Path] || n.Rev == f.Rev {
					continue
				}
				g.RebaseRev(f.Path, n.Rev, f.Rev)
				carried++
			}
			if err := mapx.Save(g, filepath.Join(absRoot, mapx.DefaultPath)); err != nil {
				return true, err
			}
			fmt.Printf("✓ the measurements of %d file(s) carried to their new revision\n", carried)
		}
	}

	// 5. The way from each old code to its new one.
	if len(pairs) > 0 {
		var b strings.Builder
		b.WriteString("# Codes renamed by `anchors migrate` (format 7): old → new. Outside the repository —\n")
		b.WriteString("# issues, pull requests, old commit messages — the old code stays; this is the way to the new.\n")
		fmt.Fprintf(&b, "%s:\n", time.Now().Format("2006-01-02"))
		for _, p := range pairs {
			fmt.Fprintf(&b, "  %s: %s\n", p.old, p.new)
		}
		prev, _ := os.ReadFile(filepath.Join(absRoot, RenamesFile))
		if err := os.WriteFile(filepath.Join(absRoot, RenamesFile), append(prev, []byte(b.String())...), 0o644); err != nil {
			return true, err
		}
		fmt.Printf("✓ %d code rename(s) recorded in %s\n", len(pairs), RenamesFile)
	}
	return changed, nil
}

// rewriteText rewrites a file in place; a missing file is left alone.
func rewriteText(p string, f func(string) string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	out := f(string(b))
	if out == string(b) {
		return nil
	}
	info, err := os.Stat(p)
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(out), info.Mode())
}
