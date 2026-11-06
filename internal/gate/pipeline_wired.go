// @anchors
//   code: PPWRP
//   ref: PWSPP

package gate

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// pipeline-wired: the stages the configuration declares and the stages the workflows call are
// the same. A workflow calling a stage the configuration lacks is refused at run time, and a
// declared stage no workflow calls is a guard that guards nothing while the configuration says
// it does; a workflow that calls no stage runs unconfronted. Each is read here, before anything
// runs (requested by MIF: "toda etapa que chama o Anchors no pipeline tem de estar
// configurada").

var pipelineCallRE = regexp.MustCompile(`anchors\s+pipeline\s+(can-run|done|run)\s+([A-Za-z0-9_.-]+)`)

func checkPipelineWired(_ config.Gate, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	p := cfg.Pipeline
	if p == nil || len(p.Stages) == 0 {
		return Pending, i18n.T("gate.pipeline_wired.no_pipeline")
	}
	var files []string
	for _, glob := range p.WorkflowGlobs() {
		found, _ := doublestar.Glob(os.DirFS(root), glob, doublestar.WithFilesOnly())
		files = append(files, found...)
	}
	sort.Strings(files)
	if len(files) == 0 {
		return Pending, i18n.T("gate.pipeline_wired.no_workflows", strings.Join(p.WorkflowGlobs(), ", "))
	}
	called := map[string]map[string]bool{} // stage → verbs
	var gaps []string
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		calls := pipelineCallRE.FindAllStringSubmatch(string(b), -1)
		if len(calls) == 0 {
			if !unguarded(p, f) {
				gaps = append(gaps, i18n.T("gate.pipeline_wired.unguarded_workflow", f))
			}
			continue
		}
		for _, c := range calls {
			verb, stage := c[1], c[2]
			if _, ok := p.Stages[stage]; !ok {
				gaps = append(gaps, i18n.T("gate.pipeline_wired.unknown_stage", f, stage))
				continue
			}
			if called[stage] == nil {
				called[stage] = map[string]bool{}
			}
			called[stage][verb] = true
		}
	}
	for _, s := range p.StageNames() {
		switch verbs := called[s]; {
		case verbs == nil:
			gaps = append(gaps, i18n.T("gate.pipeline_wired.uncalled_stage", s))
		case verbs["can-run"] && !verbs["done"] && !verbs["run"]:
			gaps = append(gaps, i18n.T("gate.pipeline_wired.no_record", s))
		}
	}
	if len(gaps) > 0 {
		return Fail, strings.Join(gaps, "; ")
	}
	return Pass, i18n.T("gate.pipeline_wired.pass", len(files), len(p.Stages))
}

// unguarded says whether the configuration declares a workflow as calling no stage, by its
// path or its file name.
func unguarded(p *config.Pipeline, file string) bool {
	if _, ok := p.Unguarded[file]; ok {
		return true
	}
	_, ok := p.Unguarded[path.Base(file)]
	return ok
}
