package initx

import (
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func testQuestions() []Question {
	return Questions(&Proposal{Config: nil}, []string{"go", "nextjs"})
}

// verdictOf returns the verdict of one question, failing the test when it is missing.
func verdictOf(t *testing.T, st []StatusResposta, id string) StatusResposta {
	t.Helper()
	for _, s := range st {
		if s.ID == id {
			return s
		}
	}
	t.Fatalf("no verdict for %q", id)
	return StatusResposta{}
}

// The contract exists so an agent can answer without seeing the TUI. Every question must
// carry what it needs to DECIDE: what is accepted, what Anchors inferred, and what the
// answer changes in the project.
func TestQuestionsCarryWhatTheAgentNeedsToDecide(t *testing.T) {
	t.Run("INQSN-B02: Every question carries what the agent needs to decide", func(t *testing.T) {})
	for _, q := range testQuestions() {
		if q.ID == "" || q.Texto == "" || q.Tipo == "" {
			t.Errorf("incomplete question: %+v", q)
		}
		// Without the "why", the agent chooses by the option's NAME — which is how one
		// chooses wrong. It is the difference between answering and guessing.
		if q.PorQue == "" {
			t.Errorf("%s does not say what the answer changes in the project", q.ID)
		}
		if q.Tipo == "select" && len(q.Opcoes) == 0 {
			t.Errorf("%s is a select and offers no options", q.ID)
		}
	}
}

// Every decision of the TUI must be here, in the TUI's order: a question left out would be
// decided in silence by the default, in a file the user thinks they decided.
func TestQuestionsFollowTheTUIOrder(t *testing.T) {
	t.Run("INQSN-B01: The questions come in the order of the terminal UI", func(t *testing.T) {})
	want := []string{"preset", "header", "artifacts", "gates", "colocation", "layers",
		"workflow", "repo", "labels", "governs"}
	var got []string
	for _, q := range testQuestions() {
		got = append(got, q.ID)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}
}

// The defaults are the reading of the real project, never a guess.
func TestQuestionsDefaultsComeFromTheInference(t *testing.T) {
	t.Run("INQSN-B03: The defaults come from the inference", func(t *testing.T) {})
	p := &Proposal{
		Colocated: true,
		HasSpecMD: true,
		Config: &config.Config{Layers: map[string]config.Layer{
			"web-code": {Kind: "code"},
			"api-code": {Kind: "code"},
			"spec":     {Kind: "spec"},
		}},
	}
	byID := map[string]Question{}
	for _, q := range Questions(p, nil) {
		byID[q.ID] = q
	}
	if byID["colocation"].Default != true {
		t.Errorf("colocation default = %v, want true", byID["colocation"].Default)
	}
	layers := []string{"api-code", "web-code"}
	if !reflect.DeepEqual(byID["layers"].Opcoes, layers) || !reflect.DeepEqual(byID["layers"].Default, layers) {
		t.Errorf("layers = %v / %v, want %v", byID["layers"].Opcoes, byID["layers"].Default, layers)
	}
	if d, _ := byID["artifacts"].Default.([]string); !contem(d, "spec") {
		t.Errorf("artifacts default = %v, want it to include spec", byID["artifacts"].Default)
	}

	// No proposal at all: empty defaults, and no panic.
	for _, q := range Questions(nil, nil) {
		switch q.ID {
		case "colocation":
			if q.Default != false {
				t.Errorf("with no proposal colocation defaults to false, got %v", q.Default)
			}
		case "layers":
			if len(q.Opcoes) != 0 {
				t.Errorf("with no proposal there are no layers, got %v", q.Opcoes)
			}
		}
	}
}

// The work-queue mode is a HUMAN decision (where the queue lives) and `init` did not make
// it: every project was born `local` by omission, and whoever wanted `github` had to find
// the field and edit the YAML by hand.
func TestQuestionsPresetAndWorkflowChoices(t *testing.T) {
	t.Run("INQSN-B04: The preset and the work-queue mode have their choices and defaults", func(t *testing.T) {})
	byID := map[string]Question{}
	for _, q := range testQuestions() {
		byID[q.ID] = q
	}
	if p := byID["preset"]; !reflect.DeepEqual(p.Opcoes, []string{"nenhum", "go", "nextjs"}) || p.Default != "nenhum" {
		t.Errorf("preset = %v (default %v)", p.Opcoes, p.Default)
	}
	if w := byID["workflow"]; !reflect.DeepEqual(w.Opcoes, []string{"local", "manual", "github"}) || w.Default != "local" {
		t.Errorf("workflow = %v (default %v)", w.Opcoes, w.Default)
	}
}

// One invalid answer refuses the SET. Writing the valid ones would produce an anchors.yaml
// nobody decided in full — and such a file loads without error, governs wrong, and does
// not point at the cause.
func TestQuestionsInvalidAnswerRefusesEverything(t *testing.T) {
	t.Run("INQSN-B06: An answer outside the options is refused with the accepted values", func(t *testing.T) {})
	t.Run("INQSN-B09: One refused answer refuses the whole set", func(t *testing.T) {})
	t.Run("INQSN-X01: An invalid answer is not corrected", func(t *testing.T) {})
	qs := testQuestions()
	bad := "preset-que-nao-existe"

	st := ValidateAnswers(qs, Respostas{Preset: &bad})

	if TudoAceito(st) {
		t.Fatal("a nonexistent preset should refuse the set")
	}
	// And the verdict must say WHICH failed and why — "something went wrong" is not enough.
	s := verdictOf(t, st, "preset")
	if s.Aceita {
		t.Error("the invalid preset was accepted")
	}
	if s.Valor != bad {
		t.Errorf("the invalid value must be kept as given, got %v", s.Valor)
	}
	if !strings.Contains(s.Detalhe, "aceitos:") {
		t.Errorf("the message should list the valid values: %s", s.Detalhe)
	}

	// A list answer outside the options is refused the same way.
	code := []string{"code"}
	if a := verdictOf(t, ValidateAnswers(qs, Respostas{Artifacts: &code}), "artifacts"); a.Aceita || !strings.Contains(a.Detalhe, "aceitos:") {
		t.Errorf("an artifact outside the options must be refused with the accepted values: %+v", a)
	}
	// A multiple choice with no declared options (the labels) accepts any value.
	free := []string{"whatever"}
	if l := verdictOf(t, ValidateAnswers(qs, Respostas{Labels: &free}), "labels"); !l.Aceita {
		t.Errorf("a free label should be accepted: %+v", l)
	}
}

// The verdict covers ALL the answers, not only the invalid ones. It is what lets the agent
// check that Anchors understood what it meant — a misspelt flag, without it, would be
// indistinguishable from an accepted answer.
func TestQuestionsVerdictCoversEveryAnswer(t *testing.T) {
	t.Run("INQSN-B05: Every question gets a verdict and unanswered ones take the default", func(t *testing.T) {})
	t.Run("INQSN-I01: No answer goes missing from the verdict", func(t *testing.T) {})
	qs := testQuestions()
	yes := true

	st := ValidateAnswers(qs, Respostas{Header: &yes})

	if len(st) != len(qs) {
		t.Fatalf("expected %d verdicts (one per question), got %d", len(qs), len(st))
	}
	for i, s := range st {
		if s.ID != qs[i].ID {
			t.Errorf("verdict %d is %q, want %q", i, s.ID, qs[i].ID)
		}
		switch s.ID {
		case "header":
			if s.UsouPada {
				t.Error("header was answered explicitly — it is not the default")
			}
		case "repo":
			// Outside github, repo is answered by its default "" and marked as such.
			if !s.UsouPada {
				t.Error("repo was not answered — it must show as the default")
			}
		case "governs":
			if !s.UsouPada {
				t.Error("governs was not answered — it must show as the default")
			}
		default:
			if !s.UsouPada {
				t.Errorf("%s was not answered — it must show as the default", s.ID)
			}
			if !reflect.DeepEqual(s.Valor, qs[i].Default) {
				t.Errorf("%s: unanswered value %v, want the default %v", s.ID, s.Valor, qs[i].Default)
			}
		}
	}
}

// The distinction behind the pointers: "not answered" (the default inferred from disk
// applies) and "answered empty" (`--artifacts=""`, no artifact) are OPPOSITE decisions,
// and a zero-value bool does not tell them apart.
func TestQuestionsNotAnsweredIsNotAnsweredEmpty(t *testing.T) {
	t.Run("INQSN-B10: An answer given empty is not the default", func(t *testing.T) {})
	qs := Questions(&Proposal{}, nil)
	empty := []string{}

	unanswered := verdictOf(t, ValidateAnswers(qs, Respostas{}), "artifacts")
	answeredEmpty := verdictOf(t, ValidateAnswers(qs, Respostas{Artifacts: &empty}), "artifacts")

	if !unanswered.UsouPada {
		t.Error("without the flag, artifacts must fall back to the default")
	}
	if answeredEmpty.UsouPada {
		t.Error("an empty `--artifacts=` is a deliberate choice, not a missing answer")
	}
	if v, ok := answeredEmpty.Valor.([]string); !ok || len(v) != 0 {
		t.Errorf("the empty answer must be kept as given, got %#v", answeredEmpty.Valor)
	}
}

// In `github` mode, `repo` and `labels` are required (WORKFLOW.md §2). Without the repo, the
// writing would land in the wrong place; without the label, the flow would take a product
// issue.
func TestQuestionsGitHubRequiresRepoAndLabels(t *testing.T) {
	t.Run("INQSN-B07: The github mode requires repository and labels", func(t *testing.T) {})
	qs := testQuestions()
	github := "github"

	st := ValidateAnswers(qs, Respostas{Workflow: &github})

	if TudoAceito(st) {
		t.Fatal("github mode with no repo nor labels should refuse")
	}
	missing := map[string]bool{}
	for _, s := range st {
		if !s.Aceita {
			missing[s.ID] = true
		}
	}
	for _, id := range []string{"repo", "labels"} {
		if !missing[id] {
			t.Errorf("%s is required in github mode and passed", id)
		}
	}

	// The happy path: github with both requirements met passes.
	repo := "acme/exemplo"
	labels := []string{"anchors"}
	st = ValidateAnswers(qs, Respostas{Workflow: &github, Repo: &repo, Labels: &labels})
	if !TudoAceito(st) {
		for _, s := range st {
			if !s.Aceita {
				t.Errorf("%s refused: %s", s.ID, s.Detalhe)
			}
		}
	}
}

// In `local` mode, declaring `repo` is not a harmless field: whoever reads the file
// concludes the integration is active (WORKFLOW.md §2).
func TestQuestionsLocalRefusesGitHubFields(t *testing.T) {
	t.Run("INQSN-B08: A repository outside the github mode is refused", func(t *testing.T) {})
	qs := testQuestions()
	local, repo := "local", "owner/nome"

	st := ValidateAnswers(qs, Respostas{Workflow: &local, Repo: &repo})

	if verdictOf(t, st, "repo").Aceita {
		t.Error("`repo` in local mode makes the file lie about the integration being active")
	}
}
