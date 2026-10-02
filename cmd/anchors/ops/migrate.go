// @anchors
//   ref: MGCMM

package ops

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/flowx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
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
			for _, alvo := range alvos {
				if f, err := migra.FormatOf(alvo); err == nil && f < from {
					from = f
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
