// @anchors
//   ref: VTRST

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const loginHead = "<!-- @anchors\n  code: LOGIN\n-->\n# Login\n\n## States\n\n### LOGIN-S01: Empty\n\n### LOGIN-S02: Filled\n\n### LOGIN-S03: Refused\n\n"

func loginUnit(t *testing.T, spec string) string {
	t.Helper()
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "ui"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, "ui", "Login.spec.md"), []byte(spec), 0o644))
	return root
}

var loginCode = mapx.Node{ID: "ui/Login.tsx", Kind: mapx.KindCode}

func TestStateTransitions_nothingToConfront(t *testing.T) {
	t.Run("VTRST-B01: What has nothing to confront leaves both gates without a verdict", func(t *testing.T) {})
	t.Run("VTRST-E01: A code file with no spec beside it leaves without a verdict", func(t *testing.T) {})
	root := loginUnit(t, "# Login, no code\n")
	if v, _ := checkValidationTransitions("", mapx.Node{ID: "ui/Login.spec.md", Kind: mapx.KindSpec}, root, nil, nil); v != Skip {
		t.Errorf("a spec node: %v", v)
	}
	if v, _ := checkErrorMessageDeclared("", loginCode, root, nil, nil); v != Skip {
		t.Errorf("a spec with no code: %v", v)
	}
	if v, d := checkValidationTransitions("", mapx.Node{ID: "ui/Other.tsx", Kind: mapx.KindCode}, root, nil, nil); v != Skip || !strings.Contains(d, "spec") {
		t.Errorf("no spec beside it: %v %s", v, d)
	}
}

func TestValidationTransitions(t *testing.T) {
	t.Run("VTRST-B02: Every validation is the trigger of a transition", func(t *testing.T) {})
	t.Run("VTRST-B03: A transition from or to an unknown state does not count", func(t *testing.T) {})
	t.Run("VTRST-B04: A validation exempted with a reason is not asked", func(t *testing.T) {})
	validations := "## Validations\n\n| Rule | Field | Condition | Behavior |\n| --- | --- | --- | --- |\n" +
		"| `LOGIN-V01` | `email` | invalid | refuses |\n| `LOGIN-V02` | `email` | spaces | trims |\n\n" +
		"## Presentation validations\n\n| Rule | Prop/State | Condition | Appearance |\n| --- | --- | --- | --- |\n| `LOGIN-P01` | `loading` | true | spinner |\n\n"
	flow := "## State Flow\n\n| From | Trigger | To |\n| --- | --- | --- |\n| `LOGIN-S02` | `LOGIN-V01` | `LOGIN-S03` |\n\n"
	v, d := checkValidationTransitions("", loginCode, loginUnit(t, loginHead+validations+flow), nil, nil)
	if v != Fail || !strings.Contains(d, "LOGIN-P01, LOGIN-V02") || strings.Contains(d, "LOGIN-V01") {
		t.Errorf("V02 and P01 lead nowhere: %v %s", v, d)
	}
	bad := "## State Flow\n\n| From | Trigger | To |\n| --- | --- | --- |\n| `LOGIN-S01` | `LOGIN-V01` | `LOGIN-S09` |\n\n"
	if v, d := checkValidationTransitions("", loginCode, loginUnit(t, loginHead+validations+bad), nil, nil); v != Fail || !strings.Contains(d, "LOGIN-V01 (LOGIN-S01 → LOGIN-S09)") {
		t.Errorf("an unknown state does not count: %v %s", v, d)
	}
	exempted := strings.Replace(strings.Replace(validations, "| trims |", "| trims @no-state: only trims spaces |", 1), "| spinner |", "| spinner @no-state |", 1)
	v, d = checkValidationTransitions("", loginCode, loginUnit(t, loginHead+exempted+flow), nil, nil)
	if v != Fail || strings.Contains(d, "LOGIN-V02") || !strings.Contains(d, "with no reason: LOGIN-P01") {
		t.Errorf("V02 is exempt, P01's exemption has no reason: %v %s", v, d)
	}
}

func TestErrorMessageDeclared(t *testing.T) {
	t.Run("VTRST-B05: Every error names the message it shows", func(t *testing.T) {})
	spec := loginHead + "## User Messages\n| Rule | Condition | Message |\n| --- | --- | --- |\n| `LOGIN-M01` | refused | \"Wrong password\" |\n\n" +
		"## Errors / Failures\n| Rule | Condition | Failure |\n| --- | --- | --- |\n" +
		"| `LOGIN-E01` | wrong password | shows `LOGIN-M01` |\n| `LOGIN-E02` | locked | shows `LOGIN-M09` |\n" +
		"| `LOGIN-E03` | offline | retries |\n| `LOGIN-E04` | analytics down | @no-message: logged only |\n"
	v, d := checkErrorMessageDeclared("", loginCode, loginUnit(t, spec), nil, nil)
	if v != Fail || !strings.Contains(d, "show: LOGIN-E03") || !strings.Contains(d, "LOGIN-E02 (LOGIN-M09)") ||
		strings.Contains(d, "LOGIN-E01") || strings.Contains(d, "LOGIN-E04") {
		t.Errorf("E03 shows nothing named, E02 a message that is none: %v %s", v, d)
	}
}

func TestStateTransitions_nestedSections(t *testing.T) {
	t.Run("VTRST-B06: Sections nested under another are found", func(t *testing.T) {})
	spec := loginHead + "## Rules (Regras de Negócio)\n\n### Validações\n\n| Regra | Campo | Condição | Comportamento |\n| --- | --- | --- | --- |\n" +
		"| `LOGIN-V01` | `email` | inválido | recusa |\n\n### Erros / Falhas\n\n| Regra | Condição | Falha |\n| --- | --- | --- |\n| `LOGIN-E01` | offline | \"Sem conexão\" |\n\n" +
		"## Fluxo de Estados\n\n| De | Gatilho | Para |\n| --- | --- | --- |\n| `LOGIN-S02` | `LOGIN-V01` | `LOGIN-S03` |\n"
	root := loginUnit(t, spec)
	if v, d := checkValidationTransitions("", loginCode, root, nil, nil); v != Pass {
		t.Errorf("the nested validation is read and linked: %v %s", v, d)
	}
	if v, d := checkErrorMessageDeclared("", loginCode, root, nil, nil); v != Fail || !strings.Contains(d, "LOGIN-E01") {
		t.Errorf("the nested error is read: %v %s", v, d)
	}
}
