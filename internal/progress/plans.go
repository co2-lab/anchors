// @anchors
//   code: PRMGP
//   ref: DMLDD

package progress

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// --- a plan as a demand, and the per-plan progress files it replaces ---

// phaseRE is a phase a plan declares in a header: `## FNDTN-W01 — Foundation`.
func phaseRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^#{2,4}[^\S\n]+([A-Z0-9]` + config.CodeLengthPattern() + `-` +
		config.PhaseLetter + `\d{2})\b[^\S\n]*—?[^\S\n]*(.*)$`)
}

// titleRE is a markdown document's title.
var titleRE = regexp.MustCompile(`(?m)^#[^\S\n]+(.+)$`)

// cleanText makes a text from a plan fit the format: one line, no separator.
func cleanText(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.ReplaceAll(s, " · ", " - ")
}

// PlanDemand adds a plan as a demand: titled as the plan (or by the title given), with one
// stage per phase the plan declares in its headers. A plan already tracked is refused.
func (f *File) PlanDemand(plan, content, title string) (string, error) {
	plan = filepath.ToSlash(plan)
	for _, d := range f.Demands {
		if d.Plan == plan {
			return "", fmt.Errorf("%s is already tracked by %s", plan, d.ID)
		}
	}
	if title == "" {
		if m := titleRE.FindStringSubmatch(content); m != nil {
			title = cleanText(m[1])
		} else {
			title = strings.TrimSuffix(path.Base(plan), path.Ext(plan))
		}
	}
	id, err := f.AddDemand(title, plan)
	if err != nil {
		return "", err
	}
	for _, m := range phaseRE().FindAllStringSubmatch(content, -1) {
		stage := m[1]
		if t := cleanText(m[2]); t != "" {
			stage += " — " + t
		}
		if _, err := f.AddStage(id, stage); err != nil {
			return "", err
		}
	}
	return id, nil
}

// LegacySuffix is the suffix of the per-plan progress files anchors.progress.md replaces.
const LegacySuffix = "-progress.md"

var (
	legacyStageRE = regexp.MustCompile(`^#{2,4}\s+(.+?)\s*$`)
	legacyItemRE  = regexp.MustCompile(`^\s*- \[([ xX])\]\s+(.+?)\s*$`)
)

// Migrated is what the migration of one per-plan progress file did.
type Migrated struct {
	File   string // the per-plan file, deleted
	Plan   string
	Demand string // the demand it became; empty when the plan was already tracked
	Tasks  int
	Done   int
}

// LegacyFiles are the per-plan progress files of a project: beside a plan of the same name.
// The directories Anchors never walks are skipped.
func LegacyFiles(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if p != root && (strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "_build") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(name, LegacySuffix) {
			return nil
		}
		plan := strings.TrimSuffix(p, LegacySuffix) + ".md"
		if _, err := os.Stat(plan); err != nil {
			return nil // a progress file of no plan is not Anchors'
		}
		rel, _ := filepath.Rel(root, p)
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	return out, err
}

// MigrateLegacy moves one per-plan progress file into f: its plan becomes a demand with one
// stage per phase section, and each item a task — done when it was checked, weighing 1 as an
// estimate to revise. The placeholder items of an unfilled template are left behind. A plan
// already tracked keeps its demand.
func (f *File) MigrateLegacy(rel, content, who, date string) (Migrated, error) {
	plan := strings.TrimSuffix(rel, LegacySuffix) + ".md"
	m := Migrated{File: rel, Plan: plan}
	for _, d := range f.Demands {
		if d.Plan == plan {
			return m, nil
		}
	}
	title := strings.TrimSuffix(path.Base(plan), ".md")
	if t := titleRE.FindStringSubmatch(content); t != nil {
		title = cleanText(strings.TrimPrefix(t[1], "Progress — "))
	}
	id, err := f.AddDemand(title, plan)
	if err != nil {
		return m, err
	}
	m.Demand = id
	parent := id
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "# ") {
			continue
		}
		if s := legacyStageRE.FindStringSubmatch(line); s != nil {
			parent, err = f.AddStage(id, cleanText(s[1]))
			if err != nil {
				return m, err
			}
			continue
		}
		it := legacyItemRE.FindStringSubmatch(line)
		if it == nil || placeholderRE.MatchString(it[2]) {
			continue
		}
		tid, err := f.AddTask(parent, cleanText(it[2]), 1, true)
		if err != nil {
			return m, err
		}
		m.Tasks++
		if strings.ToLower(it[1]) == "x" {
			if err := f.Done(tid, who, date); err != nil {
				return m, err
			}
			m.Done++
		}
	}
	return m, nil
}

// placeholderRE is an item an unfilled template left: it cites no work yet.
var placeholderRE = regexp.MustCompile(`^(TODO|FIXME|XXX|\.\.\.)\b`)

// MigrateProject moves every per-plan progress file of the project into anchors.progress.md,
// each demand titled as its plan, under its lock, and deletes them. With dryRun nothing is written, and what would be done is
// returned.
func MigrateProject(root, unit, date string, dryRun bool) ([]Migrated, error) {
	files, err := LegacyFiles(root)
	if err != nil || len(files) == 0 {
		return nil, err
	}
	var out []Migrated
	migrate := func(f *File) error {
		for _, rel := range files {
			b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				return err
			}
			m, err := f.MigrateLegacy(rel, string(b), "migrated", date)
			if err != nil {
				return fmt.Errorf("%s: %w", rel, err)
			}
			// The template titled the file by the plan's code; the demand takes the plan's title.
			if d := f.Demand(m.Demand); d != nil {
				if b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(m.Plan))); err == nil {
					if t := titleRE.FindStringSubmatch(string(b)); t != nil {
						d.Title = cleanText(t[1])
					}
				}
			}
			out = append(out, m)
		}
		return nil
	}
	if dryRun {
		f, err := Load(root)
		if err != nil {
			return nil, err
		}
		return out, migrate(f)
	}
	if err := Update(root, unit, migrate); err != nil {
		return nil, err
	}
	for _, rel := range files {
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			return out, err
		}
	}
	return out, nil
}
