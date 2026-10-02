// @anchors
//   ref: USSTS

package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ABSENCE IS NOT AN ERROR, and it is the most common case: a freshly cloned project has no file,
// and that is exactly when the agent must ask.
func TestLoad_noFileIsNotAnError(t *testing.T) {
	t.Run("USSTS-B02: A missing settings file is not an error", func(t *testing.T) {})
	s, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("a project without `%s` gave an error: %v", File, err)
	}
	if s.Decided() {
		t.Error("without a file the decision cannot count as taken")
	}
	if s.HandlesUserIssues() {
		t.Error("the default is CLOSED — whoever can decide declares they can")
	}
}

// THREE STATES, and the distinction is the mechanism: `nil` is the only case where the agent
// asks. Without the pointer, "not declared" and "declared no" would be the same value.
func TestSettings_threeStates(t *testing.T) {
	t.Run("USSTS-B03: The legacy field has three states", func(t *testing.T) {})
	cases := []struct {
		name    string
		value   *bool
		decided bool
		handles bool
	}{
		{"never asked", nil, false, false},
		{"said no", Bool(false), true, false},
		{"said yes", Bool(true), true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Settings{UserIssues: c.value}
			if s.Decided() != c.decided {
				t.Errorf("Decided() = %v, want %v", s.Decided(), c.decided)
			}
			if s.HandlesUserIssues() != c.handles {
				t.Errorf("HandlesUserIssues() = %v, want %v", s.HandlesUserIssues(), c.handles)
			}
		})
	}
}

// THE OLD FIELD still holds for whoever already declared: ignoring it would make whoever declared
// `user_issues: true` before roles lose the capability with nothing warning, mid-work.
func TestSettings_compatibilityWithTheOldField(t *testing.T) {
	t.Run("USSTS-B04: The legacy yes grants only the product decision", func(t *testing.T) {})
	t.Run("USSTS-B05: A declared role wins over the legacy field", func(t *testing.T) {})
	old := Settings{UserIssues: Bool(true)}
	if !old.Can(CapDecideProduct) {
		t.Error("`user_issues: true` without a role must keep deciding the product")
	}
	if !old.Decided() {
		t.Error("whoever declared through the old field already decided — must not be asked again")
	}
	// But it answers ONE question only: it was a boolean, and had no way to say more.
	if old.Can(CapWritePlan) || old.Can(CapReview) {
		t.Error("the old field cannot unlock a capability it did not declare")
	}

	// And the role WINS when both exist: it is the new source, and keeping both would make the
	// answer come from whichever someone forgot to change.
	both := Settings{Role: RoleDev, UserIssues: Bool(true)}
	if both.Can(CapDecideProduct) {
		t.Error("the declared role must win over the old field")
	}
}

// What was saved is what is read back.
func TestSaveLoad_givesBackWhatWasSaved(t *testing.T) {
	t.Run("USSTS-B06: What is saved is what is loaded", func(t *testing.T) {})
	root := t.TempDir()
	if err := Save(root, Settings{
		UserIssues: Bool(true), Agent: "maquina/sessao", DecidedAt: "2026-09-08",
	}); err != nil {
		t.Fatal(err)
	}

	s, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !s.HandlesUserIssues() || s.Agent != "maquina/sessao" || s.DecidedAt != "2026-09-08" {
		t.Errorf("got back %+v", s)
	}
}

// The file EXPLAINS what it is, because whoever finds it later was not in the conversation.
func TestSave_theFileSaysWhatItIs(t *testing.T) {
	t.Run("USSTS-B07: The saved file explains itself", func(t *testing.T) {})
	root := t.TempDir()
	if err := Save(root, Settings{UserIssues: Bool(false)}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(Path(root))
	if err != nil {
		t.Fatal(err)
	}
	// What the header must say, and why:
	//   · that it is local — or someone tries to version it;
	//   · who decides the product — the difference between roles that shows most;
	//   · how to declare — whoever finds the file was not in the conversation.
	text := strings.ToLower(string(b))
	for _, want := range []string{
		"não vai para o git", "gitignore", "escalonados",
		"product-owner", "architect", "anchors settings role",
	} {
		if !strings.Contains(text, strings.ToLower(want)) {
			t.Errorf("the header does not mention %q:\n%s", want, b)
		}
	}
}

// The file lives in `.anchors/`, already in the projects' `.gitignore` — which keeps it local.
func TestPath_livesInAnchors(t *testing.T) {
	t.Run("USSTS-B01: The settings live in the local state folder", func(t *testing.T) {})
	if got := Path("/proj"); got != filepath.Join("/proj", ".anchors", "settings.yaml") {
		t.Errorf("Path = %q", got)
	}
}

// The answer is read in BOTH languages and in short forms: the question is asked in the terminal,
// and whoever answers types what comes first.
func TestParseAnswer(t *testing.T) {
	t.Run("USSTS-B08: Typed answers are read in both languages", func(t *testing.T) {})
	yes := []string{"s", "sim", "y", "yes", "SIM", " Sim \n"}
	no := []string{"n", "nao", "não", "no", "NÃO", " nao \n"}
	notUnderstood := []string{"", "talvez", "1", "ok", "quem sabe"}

	for _, r := range yes {
		if v := ParseAnswer(r); v == nil || !*v {
			t.Errorf("ParseAnswer(%q) was not read as yes", r)
		}
	}
	for _, r := range no {
		if v := ParseAnswer(r); v == nil || *v {
			t.Errorf("ParseAnswer(%q) was not read as no", r)
		}
	}
	// What is not understood returns `nil` — and the caller asks again. Assuming "no" here would
	// be convenient and wrong: whoever typed something was answering.
	for _, r := range notUnderstood {
		if v := ParseAnswer(r); v != nil {
			t.Errorf("ParseAnswer(%q) = %v, want nil", r, *v)
		}
	}
}

func TestDescribe_namesTheRoleOrAsksForOne(t *testing.T) {
	t.Run("USSTS-B09: The description names the role or asks for one", func(t *testing.T) {})
	if got := (Settings{Role: RolePO, Agent: "maq/1"}).Describe(); got != "Product Owner · maq/1" {
		t.Errorf("Describe with a role = %q, want %q", got, "Product Owner · maq/1")
	}
	for _, legacy := range []Settings{{UserIssues: Bool(true)}, {UserIssues: Bool(false)}} {
		if got := legacy.Describe(); !strings.Contains(got, "anchors settings role") {
			t.Errorf("Describe of legacy settings = %q, want it to point to `anchors settings role`", got)
		}
	}
	if (Settings{}).Describe() == (Settings{UserIssues: Bool(true)}).Describe() {
		t.Error("undecided and legacy settings must not read the same")
	}
}

func TestLoad_invalidYAMLNamesTheFile(t *testing.T) {
	t.Run("USSTS-E01: A settings file that is not YAML fails naming it", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root), []byte("role: [unclosed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err == nil || !strings.Contains(err.Error(), Path(root)) {
		t.Fatalf("Load of invalid YAML = %v, want an error naming %s", err, Path(root))
	}
}
