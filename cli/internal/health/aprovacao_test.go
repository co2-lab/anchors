package health

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
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
	if fs := checkAprovacaoAlcancavel(cfg); len(fs) != 0 {
		t.Errorf("zero exigido não pode gerar achado: %+v", fs)
	}

	// Sem `gh` no PATH o doctor já reclama noutro achado — não duplicar.
	if fs := checkAprovacaoAlcancavel(nil); len(fs) != 0 {
		t.Errorf("config nula não deveria gerar achado: %+v", fs)
	}
}

// A mensagem precisa NOMEAR as duas saídas. Um aviso que diz "está travado" sem dizer o
// que fazer transfere o problema em vez de resolvê-lo.
func TestMensagemNomeiaAsDuasSaidas(t *testing.T) {
	um := 1
	cfg := &config.Config{Workflow: &config.Workflow{
		Mode: config.ModeGitHub, Repo: "repo/inexistente-de-proposito", RequiredApprovals: &um,
	}}
	fs := checkAprovacaoAlcancavel(cfg)
	if len(fs) == 0 {
		t.Skip("sem `gh` autenticado neste ambiente — o achado não é produzido")
	}
	d := fs[0].Detail
	for _, esperado := range []string{"conta de serviço", "required_approvals: 0", "doctor --fix"} {
		if !strings.Contains(d, esperado) {
			t.Errorf("a mensagem deveria citar %q: %s", esperado, d)
		}
	}
}
