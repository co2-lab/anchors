package initx

import (
	"fmt"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// THE CONTRIBUTING GUIDE IS THE PROJECT'S OWN CONFIGURATION, READ ALOUD.
//
// A contributor — a person or an agent — arrives at a project under Anchors and does not
// know that the spec comes first, which command gives the next task, or which gate bars a
// commit. All of that is already decided in anchors.yaml, in a form a reader does not
// read. The seeded guide says it in prose, from the configuration init just wrote, so it
// names this project's layers and this project's blocking gates, not an example's.
//
// It proposes no structure: the layers it lists are the ones the project declared, and the
// layer kinds it names are illustration, never a layout to move files into.

// ContributingFile is the guide's name at the project root.
const ContributingFile = "CONTRIBUTING.md"

// RenderContributing is the whole guide for a project that has none.
func RenderContributing(cfg *config.Config, guideDir string) string {
	var b strings.Builder
	b.WriteString("# Contributing\n\n")
	b.WriteString("This project is kept coherent by [Anchors](https://github.com/co2-lab/anchors). The\n")
	b.WriteString("section below says how work flows here; it was seeded by `anchors init` from\n")
	b.WriteString("`anchors.yaml`, and the configuration stays the source of truth.\n\n")
	b.WriteString(ContributingSection(cfg, guideDir))
	return b.String()
}

// ContributingSection is the part about working with Anchors: what init would add to a
// guide the project already has.
func ContributingSection(cfg *config.Config, guideDir string) string {
	var b strings.Builder
	b.WriteString("## Working with Anchors\n\n")

	b.WriteString("### The order of the work\n\n")
	b.WriteString("The spec is the anchor: a change starts in the spec (`*.spec.md`), which states the\n")
	b.WriteString("rules; the feature (`*.feature`) turns each rule into scenarios; the tests prove the\n")
	b.WriteString("scenarios; the code realises them. A bug fixed without a rule and a scenario is a bug\n")
	b.WriteString("the next change can bring back — write the rule, then the scenario, then the test.\n\n")
	if a := artifactLines(cfg); len(a) > 0 {
		b.WriteString("Where each piece lives here:\n\n")
		for _, l := range a {
			b.WriteString(l + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("### Layers\n\n")
	b.WriteString("Anchors does not care how the folders are arranged, only that the layers are kept\n")
	b.WriteString("apart and that each file respects the layer it is in. Layers are the project's own\n")
	b.WriteString("decision — for example: entry points, use cases, domain, repositories,\n")
	b.WriteString("infrastructure, presentation. Those names are illustration, not a layout to fit.\n\n")
	if ls := codeLayerLines(cfg); len(ls) > 0 {
		b.WriteString("The code layers this project declares:\n\n")
		for _, l := range ls {
			b.WriteString(l + "\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("No code layer is declared yet: declare each one in `anchors.yaml` once the code exists.\n\n")
	}

	b.WriteString("### Day to day\n\n")
	b.WriteString("- `anchors status` — where the project stands, and the next step.\n")
	if cfg != nil && cfg.Workflow != nil && cfg.Workflow.Mode == config.ModeGitHub {
		fmt.Fprintf(&b, "- `anchors next` — claim the next task; the queue lives in the issues of `%s`.\n", cfg.Workflow.Repo)
	} else {
		b.WriteString("- `anchors next` — claim the next task from the local queue (`.anchors/tasks/`).\n")
	}
	b.WriteString("- `anchors check --changed <file>` — run the gates over what a change reaches.\n")
	b.WriteString("- `anchors test` — run the test suites and ingest their results into the map.\n")
	b.WriteString("- `anchors touch` — bump `updated_at` in the headers of the files that changed.\n")
	b.WriteString("- `anchors install-hooks` — once per clone: the pre-commit runs the gates over what is staged.\n\n")

	b.WriteString("### What blocks and what informs\n\n")
	blocking, informing := gateSplit(cfg)
	if len(blocking) > 0 {
		b.WriteString("A failure of these gates bars the commit: " + quoteAll(blocking) + ".\n")
	} else {
		b.WriteString("No gate bars a commit yet: every gate only informs.\n")
	}
	if len(informing) > 0 {
		b.WriteString("These only inform: " + quoteAll(informing) + ".\n")
	}
	b.WriteString("A divergence (⚠) and a pending item (?) inform unless `severity` in `anchors.yaml`\n")
	b.WriteString("says otherwise; `anchors check --show-drift` lists them.\n\n")

	b.WriteString("### Guides\n\n")
	b.WriteString("`anchors guide` prints the guides for writing each piece.")
	if guideDir != "" {
		fmt.Fprintf(&b, " The project's own guides are in `%s/`.", guideDir)
	}
	b.WriteString("\n")
	return b.String()
}

// artifactLines are the unit's patterns, one line per declared artifact layer, and the
// colocation templates when the project declares colocation.
func artifactLines(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}
	var out []string
	for _, kind := range []string{"spec", "feature", "test"} {
		for _, name := range sortedLayerNames(cfg) {
			if l := cfg.Layers[name]; l.Kind == kind {
				out = append(out, fmt.Sprintf("- %s: `%s`", kind, l.Pattern))
			}
		}
	}
	if d := cfg.Derived; d != nil && len(d.Files) > 0 {
		kinds := make([]string, 0, len(d.Files))
		for k := range d.Files {
			kinds = append(kinds, k)
		}
		sort.Strings(kinds)
		var parts []string
		for _, k := range kinds {
			parts = append(parts, fmt.Sprintf("the %s at `%s`", k, strings.Join(d.Files[k], "`, `")))
		}
		out = append(out, "- beside its spec: "+strings.Join(parts, "; "))
	}
	return out
}

// codeLayerLines are the declared code layers, one line each, by name.
func codeLayerLines(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}
	var out []string
	for _, name := range sortedLayerNames(cfg) {
		if l := cfg.Layers[name]; l.Kind == "code" {
			out = append(out, fmt.Sprintf("- `%s` — `%s`", name, l.Pattern))
		}
	}
	return out
}

// gateSplit names the gates whose failure blocks, and the rest.
func gateSplit(cfg *config.Config) (blocking, informing []string) {
	if cfg == nil {
		return nil, nil
	}
	for _, g := range cfg.Gates {
		if g.ActionFor(config.LevelFail) == config.ActionBlock {
			blocking = append(blocking, g.Name)
		} else {
			informing = append(informing, g.Name)
		}
	}
	return blocking, informing
}

func sortedLayerNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Layers))
	for n := range cfg.Layers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func quoteAll(names []string) string {
	q := make([]string, len(names))
	for i, n := range names {
		q[i] = "`" + n + "`"
	}
	return strings.Join(q, ", ")
}
