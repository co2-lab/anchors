package health

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
)

// O GitHub RECUSA que o autor aprove o próprio PR — regra de plataforma, sem toggle. Como
// agentes na mesma máquina compartilham a conta, `required_approvals: 1` trava o fluxo
// inteiro: o card chega a `ready-to-review`, o revisor confronta e aprova… e não consegue.
//
// O doctor precisa DIZER isso antes de alguém descobrir no meio de um merge.
func TestAvisaQuandoAprovacaoEhInalcancavel(t *testing.T) {
	// Com ZERO exigido não há o que ficar inalcançável.
	zero := 0
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode: config.ModeGitHub, Repo: "acme/x", RequiredApprovals: &zero,
	}}
	if fs := checkApprovalReachable(cfg); len(fs) != 0 {
		t.Errorf("zero exigido não pode gerar achado: %+v", fs)
	}

	// Sem `gh` no PATH o doctor já reclama noutro achado — não duplicar.
	if fs := checkApprovalReachable(nil); len(fs) != 0 {
		t.Errorf("config nula não deveria gerar achado: %+v", fs)
	}
}

// A mensagem precisa NOMEAR as duas saídas. Um aviso que diz "está travado" sem dizer o
// que fazer transfere o problema em vez de resolvê-lo.
func TestMensagemNomeiaAsDuasSaidas(t *testing.T) {
	um := 1
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode: config.ModeGitHub, Repo: "acme/x", RequiredApprovals: &um,
	}}
	// The fake `gh`: a writer with no escape, so the finding is produced — offline. The
	// test used to ask the real GitHub about a repository that does not exist.
	fakeGH(t, ghUser, ghWriter)
	fs := checkApprovalReachable(cfg)
	if len(fs) == 0 {
		t.Fatal("a writer with no escape must produce the finding")
	}
	d := fs[0].Detail
	for _, esperado := range []string{"required_approvals: 0", "doctor --fix"} {
		if !strings.Contains(d, esperado) {
			t.Errorf("a mensagem deveria citar %q: %s", esperado, d)
		}
	}
	if !strings.Contains(d, "conta de serviço") && !strings.Contains(d, "service account") {
		t.Errorf("a mensagem deveria citar conta de serviço / service account: %s", d)
	}
}

// ghAnswer is one scripted answer of the fake `gh`: when the joined arguments match the
// shell glob `match`, it prints `out` and exits with `code`.
type ghAnswer struct {
	match, out string
	code       int
}

// fakeGH puts on the PATH a `gh` that answers by rule, records every call's arguments
// and whatever arrives on stdin. Nothing reaches the network: a call no rule matches
// exits 1, which every caller here reads as "could not ask".
func fakeGH(t *testing.T, answers ...ghAnswer) (calls, stdin func() string) {
	t.Helper()
	dir := t.TempDir()
	log, in := filepath.Join(dir, "calls.txt"), filepath.Join(dir, "stdin.txt")
	s := "#!/bin/sh\necho \"$*\" >> '" + log + "'\ncase \"$*\" in\n"
	for _, a := range answers {
		s += a.match + ")\n"
		if a.out != "" {
			s += "printf '%s\\n' '" + a.out + "'\n"
		}
		s += "[ -t 0 ] || /bin/cat >> '" + in + "'\nexit " + strconv.Itoa(a.code) + " ;;\n"
	}
	s += "esac\necho 'no rule' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(s), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	read := func(p string) func() string {
		return func() string { b, _ := os.ReadFile(p); return string(b) }
	}
	return read(log), read(in)
}

var (
	ghUser       = ghAnswer{match: "'api user --jq .login'", out: "bot"}
	ghAdmin      = ghAnswer{match: "'api repos/acme/x/collaborators/bot/permission --jq .permission'", out: "admin"}
	ghWriter     = ghAnswer{match: "'api repos/acme/x/collaborators/bot/permission --jq .permission'", out: "write"}
	ghNoEnforce  = ghAnswer{match: "'api repos/acme/x/branches/main/protection --jq .enforce_admins.enabled'", out: "false"}
	ghEnforce    = ghAnswer{match: "'api repos/acme/x/branches/main/protection --jq .enforce_admins.enabled'", out: "true"}
	ghProtFailed = ghAnswer{match: "'api repos/acme/x/branches/main/protection --jq .enforce_admins.enabled'", code: 1}
)

func TestCanBypassProtection(t *testing.T) {
	for name, tc := range map[string]struct {
		answers []ghAnswer
		ok      bool
		reason  string
	}{
		"admin, admins not enforced": {[]ghAnswer{ghUser, ghAdmin, ghNoEnforce}, true, ""},
		"admin, admins enforced":     {[]ghAnswer{ghUser, ghAdmin, ghEnforce}, false, "health.approval_enforce_admins"},
		"not admin":                  {[]ghAnswer{ghUser, ghWriter}, false, "health.approval_not_admin"},
		"permission unreadable":      {[]ghAnswer{ghUser}, false, "health.approval_not_admin"},
		"user unreadable":            {nil, false, "health.approval_user_failed"},
		"protection unreadable":      {[]ghAnswer{ghUser, ghAdmin, ghProtFailed}, false, "health.approval_read_failed"},
	} {
		t.Run(name, func(t *testing.T) {
			fakeGH(t, tc.answers...)
			ok, reason := CanBypassProtection("acme/x", "main")
			want := ""
			if tc.reason != "" {
				want = i18n.T(tc.reason)
			}
			if ok != tc.ok || reason != want {
				t.Fatalf("CanBypassProtection = %v, %q; want %v, %q", ok, reason, tc.ok, want)
			}
		})
	}

	t.Run("gh missing", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		ok, reason := CanBypassProtection("acme/x", "main")
		if ok || reason != i18n.T("health.approval_gh_missing") {
			t.Fatalf("CanBypassProtection = %v, %q; want the gh-missing reason", ok, reason)
		}
	})
}

func TestCheckApprovalReachable_withFakeGH(t *testing.T) {
	um := 1
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode: config.ModeGitHub, Repo: "acme/x", RequiredApprovals: &um,
	}}

	// Admin with an escape: the flow works, nothing to report.
	fakeGH(t, ghUser, ghAdmin, ghNoEnforce)
	if fs := checkApprovalReachable(cfg); len(fs) != 0 {
		t.Fatalf("an admin with an escape must not be warned: %+v", fs)
	}

	// No escape: one warning, on the repository, naming the requirement.
	fakeGH(t, ghUser, ghWriter)
	fs := checkApprovalReachable(cfg)
	want := Finding{"aprovacao-inalcancavel", Warn, "acme/x", i18n.T("health.approval_unreachable", 1)}
	if len(fs) != 1 || fs[0] != want {
		t.Fatalf("checkApprovalReachable = %+v, want [%+v]", fs, want)
	}

	// Without gh at all the doctor already complains elsewhere.
	t.Setenv("PATH", t.TempDir())
	if fs := checkApprovalReachable(cfg); len(fs) != 0 {
		t.Fatalf("without gh there must be no duplicate finding: %+v", fs)
	}
}

func TestDisableApprovalRequirement(t *testing.T) {
	calls, stdin := fakeGH(t, ghAnswer{match: "'api --method PUT repos/acme/x/branches/main/protection --input -'"})
	if err := DisableApprovalRequirement("acme/x", "main"); err != nil {
		t.Fatal(err)
	}
	if got := calls(); got != "api --method PUT repos/acme/x/branches/main/protection --input -\n" {
		t.Fatalf("gh calls = %q", got)
	}
	var body struct {
		EnforceAdmins bool `json:"enforce_admins"`
		Reviews       struct {
			Count int `json:"required_approving_review_count"`
		} `json:"required_pull_request_reviews"`
	}
	if err := json.Unmarshal([]byte(stdin()), &body); err != nil {
		t.Fatalf("the protection body must be JSON on stdin: %v (%q)", err, stdin())
	}
	if body.EnforceAdmins || body.Reviews.Count != 0 {
		t.Fatalf("the body must zero the approvals and not enforce admins: %+v", body)
	}

	fakeGH(t) // every call fails with "no rule"
	err := DisableApprovalRequirement("acme/x", "main")
	if err == nil || err.Error() != "main: no rule" {
		t.Fatalf("a failed PUT must name the branch and carry gh's output, got %v", err)
	}
}

func TestCurrentApproval(t *testing.T) {
	q := "'api repos/acme/x/branches/main/protection --jq .required_pull_request_reviews.required_approving_review_count'"
	for name, tc := range map[string]struct {
		answers []ghAnswer
		want    int
	}{
		"two required":    {[]ghAnswer{{match: q, out: "2"}}, 2},
		"zero required":   {[]ghAnswer{{match: q, out: "0"}}, 0},
		"not a number":    {[]ghAnswer{{match: q, out: "null-ish"}}, -1},
		"protection gone": {nil, -1},
	} {
		t.Run(name, func(t *testing.T) {
			fakeGH(t, tc.answers...)
			if got := currentApproval("acme/x", "main"); got != tc.want {
				t.Fatalf("currentApproval = %d, want %d", got, tc.want)
			}
		})
	}
}
