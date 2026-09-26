package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
)

// --- anchors touch: bump `updated_at` on the files that changed ---
//
// The `updated-at-current` gate charges every touched file whose header date is not the
// day of the change. Keeping it meant rewriting `updated_at:` by hand, or a one-off script:
// measured in MIF, one `anchors stamp --refresh` rewrote the stamps of 111 tests and each
// header had to be bumped after. `check --fix` bumps too, but only inside the full check,
// with no scope and no filter.
//
// It does not lie about what changed. A file whose only difference from HEAD is its own
// `updated_at` is not bumped again, and a file with no difference at all is not touched —
// the date says when the CONTENT changed, and bumping it for nothing would make it say
// something false.

// updatedAtLineRE finds the date of an `@anchors` header, in any comment dialect
// (`//`, `#`, `<!-- -->`): the field is plain text inside the block.
var updatedAtLineRE = regexp.MustCompile(`updated_at:\s*(\d{4}-\d{2}-\d{2})`)

// touchSkip is why a changed file was not bumped.
type touchSkip string

const (
	skipNoHeader   touchSkip = "no @anchors header with updated_at"
	skipDateOnly   touchSkip = "only its updated_at differs from HEAD"
	skipUnchanged  touchSkip = "no change from HEAD"
	skipAlready    touchSkip = "already dated"
	skipExcluded   touchSkip = "excluded"
	skipUnstaged   touchSkip = "has changes outside the index — stage them or use --changed"
	skipUnreadable touchSkip = "could not be read"
)

// touchDecision is what `anchors touch` does with one file.
type touchDecision struct {
	File    string
	Bump    bool
	From    string // the date it had
	Skip    touchSkip
	Content string // the new content, when Bump
}

// decideTouch decides one file: `current` is its content now (worktree, or index with
// --staged), `base` its content at HEAD (`hasBase` false for a new file), `date` the day
// to write.
func decideTouch(file, current, base string, hasBase bool, date string) touchDecision {
	d := touchDecision{File: file}
	start, end, ok := headerSpan(current)
	var m []int
	if ok {
		if mm := updatedAtLineRE.FindStringSubmatchIndex(current[start:end]); mm != nil {
			m = make([]int, len(mm))
			for i, v := range mm {
				m[i] = v + start
			}
		}
	}
	if m == nil {
		d.Skip = skipNoHeader
		return d
	}
	d.From = current[m[2]:m[3]]
	if hasBase {
		if current == base {
			d.Skip = skipUnchanged
			return d
		}
		// Blank the date on both sides: equal means the header date is all that changed.
		if blankDate(current) == blankDate(base) {
			d.Skip = skipDateOnly
			return d
		}
	}
	if d.From == date {
		d.Skip = skipAlready
		return d
	}
	d.Bump = true
	d.Content = current[:m[2]] + date + current[m[3]:]
	return d
}

func blankDate(s string) string {
	return updatedAtLineRE.ReplaceAllString(s, "updated_at: ")
}

func newTouchCmd() *cobra.Command {
	var root, date string
	var staged, dryRun bool
	var exclude []string
	cmd := &cobra.Command{
		Use:   "touch",
		Short: "Bump `updated_at` in the @anchors header of the files that changed",
		Long: `Writes the day's date into the ` + "`updated_at`" + ` of the @anchors header of every
file that changed — the date the ` + "`updated-at-current`" + ` gate charges.

  anchors touch                 files changed in the worktree vs HEAD (new files too)
  anchors touch --staged        files in the index vs HEAD, re-staged after (for a pre-commit)
  anchors touch --dry-run       says what it would bump
  anchors touch --exclude 'docs/**' --date 2026-09-25

It does not lie: a file whose only difference from HEAD is its own updated_at, or that
did not change, is not bumped. Only files with an @anchors header carrying updated_at are
touched, in any comment dialect. The project's ` + "`touch.exclude`" + ` globs (anchors.yaml) add to
--exclude. With --staged, a file that also has changes outside the index is skipped:
re-staging it would put those changes in the commit.

The installed pre-commit already runs it on the staged files before the gates
(` + "`touch.pre_commit`" + ` in anchors.yaml, on by default; ` + "`false`" + ` turns it off).

After a ` + "`stamp --refresh`" + `, order does not matter: a ` + "`@contract`" + ` stamp never covers the header.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if date == "" {
				date = gitmeta.Today()
			}
			bumped, skipped, err := touchRun(absRoot, staged, dryRun, date, exclude)
			if err != nil {
				return err
			}
			verb := "bumped"
			if dryRun {
				verb = "would bump"
			}
			for _, d := range bumped {
				fmt.Printf("  %s %s  (%s → %s)\n", verb, d.File, d.From, date)
			}
			for _, why := range []touchSkip{skipDateOnly, skipUnchanged, skipAlready, skipExcluded, skipUnstaged, skipUnreadable} {
				for _, f := range skipped[why] {
					fmt.Printf("  skipped %s — %s\n", f, why)
				}
			}
			fmt.Printf("\n%s %d file(s) to %s.\n", verb, len(bumped), date)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&staged, "staged", false, "the files in the index (vs HEAD), re-staged after rewriting")
	cmd.Flags().StringVar(&date, "date", "", "the date to write (YYYY-MM-DD; default: today)")
	cmd.Flags().StringArrayVar(&exclude, "exclude", nil, "glob of files never to bump (repeatable; adds to `touch.exclude`)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "say what would be bumped, without writing")
	cmd.Flags().Bool("changed", true, "the files changed in the worktree vs HEAD (the default)")
	return cmd
}

// touchCandidates lists the files to consider: changed vs HEAD in the worktree (plus the
// untracked ones), or in the index with --staged. Deleted files are not candidates.
func touchCandidates(root string, staged bool) ([]string, error) {
	var lists [][]string
	if staged {
		out, err := exec.Command("git", "-C", root, "diff", "--cached", "--name-only", "--diff-filter=ACMR").Output()
		if err != nil {
			return nil, fmt.Errorf("git diff --cached: %w", err)
		}
		lists = append(lists, strings.Fields(string(out)))
	} else {
		out, err := exec.Command("git", "-C", root, "diff", "HEAD", "--name-only", "--diff-filter=ACMR").Output()
		if err != nil {
			return nil, fmt.Errorf("git diff HEAD: %w", err)
		}
		lists = append(lists, strings.Fields(string(out)))
		out, err = exec.Command("git", "-C", root, "ls-files", "--others", "--exclude-standard").Output()
		if err != nil {
			return nil, fmt.Errorf("git ls-files: %w", err)
		}
		lists = append(lists, strings.Fields(string(out)))
	}
	seen := map[string]bool{}
	var files []string
	for _, l := range lists {
		for _, f := range l {
			if !seen[f] {
				seen[f] = true
				files = append(files, f)
			}
		}
	}
	sort.Strings(files)
	return files, nil
}

// touchCurrent is the content to date: the worktree file, or the staged blob.
func touchCurrent(root, f string, staged bool) (string, bool) {
	if staged {
		return gitShow(root, ":"+f)
	}
	b, err := os.ReadFile(filepath.Join(root, f))
	if err != nil {
		return "", false
	}
	return string(b), true
}

func gitShow(root, spec string) (string, bool) {
	out, err := exec.Command("git", "-C", root, "show", spec).Output()
	if err != nil {
		return "", false
	}
	return string(out), true
}

// hasUnstaged says whether the worktree copy differs from the index copy.
func hasUnstaged(root, f string) bool {
	return exec.Command("git", "-C", root, "diff", "--quiet", "--", f).Run() != nil
}

func excludedBy(f string, globs []string) bool {
	for _, g := range globs {
		if ok, _ := doublestar.Match(g, f); ok {
			return true
		}
	}
	return false
}

// touchRun bumps the changed files (worktree, or index with staged) and says what it did.
// The project's `touch.exclude` adds to exclude. Shared by `anchors touch` and the
// pre-commit phase of `anchors verify`.
func touchRun(absRoot string, staged, dryRun bool, date string, exclude []string) ([]touchDecision, map[touchSkip][]string, error) {
	if cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile)); err == nil && cfg.Touch != nil {
		exclude = append(append([]string{}, exclude...), cfg.Touch.Exclude...)
	}
	files, err := touchCandidates(absRoot, staged)
	if err != nil {
		return nil, nil, err
	}
	var bumped []touchDecision
	skipped := map[touchSkip][]string{}
	for _, f := range files {
		if excludedBy(f, exclude) {
			skipped[skipExcluded] = append(skipped[skipExcluded], f)
			continue
		}
		current, ok := touchCurrent(absRoot, f, staged)
		if !ok {
			skipped[skipUnreadable] = append(skipped[skipUnreadable], f)
			continue
		}
		if staged && hasUnstaged(absRoot, f) {
			if updatedAtLineRE.MatchString(current) {
				skipped[skipUnstaged] = append(skipped[skipUnstaged], f)
			}
			continue
		}
		base, hasBase := gitShow(absRoot, "HEAD:"+f)
		d := decideTouch(f, current, base, hasBase, date)
		if !d.Bump {
			if d.Skip != skipNoHeader {
				skipped[d.Skip] = append(skipped[d.Skip], f)
			}
			continue
		}
		bumped = append(bumped, d)
		if dryRun {
			continue
		}
		if err := os.WriteFile(filepath.Join(absRoot, f), []byte(d.Content), 0o644); err != nil {
			return bumped, skipped, err
		}
		if staged {
			if out, err := exec.Command("git", "-C", absRoot, "add", "--", f).CombinedOutput(); err != nil {
				return bumped, skipped, fmt.Errorf("re-stage %s: %v %s", f, err, out)
			}
			if err := restageRealIndex(absRoot, f); err != nil {
				return bumped, skipped, err
			}
		}
	}
	return bumped, skipped, nil
}

// restageRealIndex keeps the REAL index in step during a partial commit.
//
// `git commit -- <paths>` runs the hook against a temporary index (`next-index-<pid>.lock`,
// in GIT_INDEX_FILE), and prepares the real one separately in `index.lock` of the same
// directory, from the worktree as it was BEFORE the hook. The hook's `git add` reached
// only the temporary index: the commit came out dated, and the real index kept the old
// date — `MM` in `git status`, and a later `git add` could revert the date (reported from
// MIF, where every commit names its paths so another session's changes never mix in).
// Staging the file in that `index.lock` too makes the real index end equal to HEAD.
//
// Outside a partial commit it does nothing: a plain commit's hook stages in the real index
// already, and `commit -a` stages in `index.lock`, which becomes the index.
func restageRealIndex(root, f string) error {
	idx := os.Getenv("GIT_INDEX_FILE")
	if !strings.HasPrefix(filepath.Base(idx), "next-index-") {
		return nil
	}
	real := filepath.Join(filepath.Dir(idx), "index.lock")
	if _, err := os.Stat(real); err != nil {
		fmt.Printf("· touch: partial commit, and the real index was not found — after the commit, "+
			"`git restore --staged %s` puts the index back in step with HEAD\n", f)
		return nil
	}
	c := exec.Command("git", "-C", root, "add", "--", f)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+real)
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("re-stage %s in the real index: %v %s", f, err, out)
	}
	return nil
}

// touchOnPreCommit says whether the pre-commit phase bumps the staged files first: on
// unless the project declares `touch.pre_commit: false`.
func touchOnPreCommit(cfg *config.Config) bool {
	return cfg == nil || cfg.Touch == nil || cfg.Touch.PreCommit == nil || *cfg.Touch.PreCommit
}

// headerSearchLines bounds where the @anchors header may start: it sits at the top.
const touchHeaderSearchLines = 10

// headerSpan is the byte range of the @anchors header block, when the file has one at its
// top: a run of line comments (`//`, `#`, `*`, `--`) that includes `@anchors`, or an HTML
// block `<!-- @anchors … -->`. Only a date INSIDE it is the header's.
//
// Searching the whole file dated text that only looks like a header. Measured on the
// Anchors repository itself: a test fixture — the string "// @anchors … updated_at:
// 2026-09-01" in a Go constant — was rewritten by the pre-commit, and the commit went in
// with the test broken.
func headerSpan(content string) (start, end int, ok bool) {
	lines := strings.SplitAfter(content, "\n")
	offset := 0
	for i := 0; i < len(lines) && i < touchHeaderSearchLines; i++ {
		t := strings.TrimSpace(lines[i])
		if !strings.Contains(t, "@anchors") || !isHeaderComment(t) {
			offset += len(lines[i])
			continue
		}
		start = offset
		end = offset + len(lines[i])
		if strings.HasPrefix(t, "<!--") {
			// HTML block: up to the line that closes it.
			for j := i; j < len(lines); j++ {
				if j > i {
					end += len(lines[j])
				}
				if strings.Contains(lines[j], "-->") {
					return start, end, true
				}
			}
			return 0, 0, false // never closed: not a header
		}
		for j := i + 1; j < len(lines) && isHeaderComment(strings.TrimSpace(lines[j])); j++ {
			end += len(lines[j])
		}
		return start, end, true
	}
	return 0, 0, false
}

func isHeaderComment(t string) bool {
	for _, p := range []string{"//", "#", "/*", "*", "--", "<!--"} {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}
