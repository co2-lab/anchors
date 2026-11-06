// @anchors
//   code: PPCNP
//   ref: PCSPP

package config

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Pipeline is what the project's pipeline asks the Anchors before each stage runs: which
// stages exist, on which branches, what each requires, and the command that validates it. A
// release step ran on a sha no one had reviewed, because nothing a workflow consulted said it
// had to (reported from MIF: the production release, run from commits made on production).
type Pipeline struct {
	// Workflows are the globs of the pipeline's definitions, read by `pipeline-wired`.
	// Empty: `.github/workflows/*.yml` and `*.yaml`.
	Workflows []string `yaml:"workflows,omitempty"`
	// Unguarded names the workflows that call no stage, each with the reason.
	Unguarded map[string]string `yaml:"unguarded,omitempty"`
	// Stages are the stages a workflow may ask about, by name.
	Stages map[string]Stage `yaml:"stages,omitempty"`
}

// Stage is one stage of the pipeline: where it runs, what it requires, and what runs and
// validates it.
type Stage struct {
	// Branches the stage runs on. Empty: any.
	Branches []string `yaml:"branches,omitempty"`
	// Requires are the conditions the stage's sha must meet, all of them.
	Requires []Requirement `yaml:"requires,omitempty"`
	// Run is the command that executes the stage, when the Anchors runs it (`pipeline run`).
	Run string `yaml:"run,omitempty"`
	// Verify is the command that validates the stage once it ran (`pipeline done`).
	Verify string `yaml:"verify,omitempty"`
}

// Requirement is one condition of a stage: `flow` — the sha came by the branch flow —, a
// phase whose gates pass (`gates: ci`), a stage the sha passed (`passed: review-release`), the
// head the sha was merged from (`from: "hotfix/*"`), or one of several (`any:`).
type Requirement struct {
	Flow   bool
	Gates  string
	Passed string
	From   string
	Any    []Requirement
}

// UnmarshalYAML reads a requirement written as the word `flow` or as a one-key map.
func (r *Requirement) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		if n.Value != "flow" {
			return fmt.Errorf("line %d: requirement %q: write flow, or one of gates:, passed:, from:, any:", n.Line, n.Value)
		}
		r.Flow = true
		return nil
	}
	var m map[string]yaml.Node
	if err := n.Decode(&m); err != nil || len(m) != 1 {
		return fmt.Errorf("line %d: a requirement is flow, or a map of one key: gates:, passed:, from: or any:", n.Line)
	}
	for k, v := range m {
		switch k {
		case "gates":
			return v.Decode(&r.Gates)
		case "passed":
			return v.Decode(&r.Passed)
		case "from":
			return v.Decode(&r.From)
		case "any":
			return v.Decode(&r.Any)
		}
		return fmt.Errorf("line %d: requirement %q: write flow, or one of gates:, passed:, from:, any:", n.Line, k)
	}
	return nil
}

// String writes the requirement as it is declared.
func (r Requirement) String() string {
	switch {
	case r.Flow:
		return "flow"
	case r.Gates != "":
		return "gates: " + r.Gates
	case r.Passed != "":
		return "passed: " + r.Passed
	case r.From != "":
		return "from: " + r.From
	}
	parts := make([]string, len(r.Any))
	for i, a := range r.Any {
		parts[i] = a.String()
	}
	return "any: [" + strings.Join(parts, ", ") + "]"
}

// WorkflowGlobs are the globs the pipeline's definitions are read from.
func (p *Pipeline) WorkflowGlobs() []string {
	if p == nil || len(p.Workflows) == 0 {
		return []string{".github/workflows/*.yml", ".github/workflows/*.yaml"}
	}
	return p.Workflows
}

// StageNames are the declared stages, sorted.
func (p *Pipeline) StageNames() []string {
	if p == nil {
		return nil
	}
	out := make([]string, 0, len(p.Stages))
	for n := range p.Stages {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// validatePipeline refuses a stage nothing executes or validates, a stage requiring a stage
// that does not exist, and an empty requirement: each would let a stage run unconfronted.
func validatePipeline(p *Pipeline) error {
	for _, name := range p.StageNames() {
		s := p.Stages[name]
		if strings.TrimSpace(s.Run) == "" && strings.TrimSpace(s.Verify) == "" {
			return fmt.Errorf("pipeline.stages.%s: a stage needs the command that runs it (run:) or validates it (verify:)", name)
		}
		for _, b := range s.Branches {
			if err := validBranchPattern(b); err != nil {
				return fmt.Errorf("pipeline.stages.%s: %w", name, err)
			}
		}
		if err := validateRequirements(p, name, s.Requires); err != nil {
			return err
		}
	}
	return nil
}

func validateRequirements(p *Pipeline, stage string, reqs []Requirement) error {
	for _, r := range reqs {
		switch {
		case r.Flow, r.Gates != "":
		case r.Passed != "":
			if _, ok := p.Stages[r.Passed]; !ok {
				return fmt.Errorf("pipeline.stages.%s: passed: %s is no stage of the pipeline", stage, r.Passed)
			}
		case r.From != "":
			if err := validBranchPattern(r.From); err != nil {
				return fmt.Errorf("pipeline.stages.%s: %w", stage, err)
			}
		case len(r.Any) > 0:
			if err := validateRequirements(p, stage, r.Any); err != nil {
				return err
			}
		default:
			return fmt.Errorf("pipeline.stages.%s: an empty requirement", stage)
		}
	}
	return nil
}
