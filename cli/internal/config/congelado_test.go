package config

import "testing"

// AUSENTE significa HABILITADO — e este é o teste que impede o pior modo de falha.
//
// `Enabled` é ponteiro justamente por isto: fosse `bool`, o zero-value `false` faria todo
// projeto que nunca declarou o campo nascer CONGELADO. Quem estivesse adotando o Anchors
// veria todos os comandos recusarem, sem ter pedido nada.
func TestCongelado_ausenteNaoCongela(t *testing.T) {
	if (&Config{}).Congelado() {
		t.Fatal("config sem `enabled` está congelada — todo projeto que nunca declarou o " +
			"campo nasceria parado")
	}
}

func TestCongelado_soOFalseExplicitoCongela(t *testing.T) {
	sim, nao := false, true
	if !(&Config{Enabled: &sim}).Congelado() {
		t.Error("`enabled: false` não congelou")
	}
	if (&Config{Enabled: &nao}).Congelado() {
		t.Error("`enabled: true` congelou")
	}
}

// Nil-safe: uma config que não carregou não congela nada.
//
// O comando que a recebeu vazia tem outro problema, e responder "congelado" ali mandaria
// quem investiga para o lado errado — ele procuraria um freeze que ninguém declarou.
func TestCongelado_configNilNaoCongela(t *testing.T) {
	var c *Config
	if c.Congelado() {
		t.Fatal("config nil respondeu congelado")
	}
	if c.MotivoDoCongelamento() != "" {
		t.Fatal("config nil devolveu motivo")
	}
}

// O motivo AUSENTE não vira silêncio: a mensagem cobra quem congelou.
//
// Um congelamento sem razão escrita é indistinguível de configuração quebrada, e quem
// esbarra nele tenta contornar em vez de ler.
func TestCongelado_semMotivoCobraQuemCongelou(t *testing.T) {
	sim := false
	m := (&Config{Enabled: &sim}).MotivoDoCongelamento()
	if m == "" {
		t.Fatal("congelado sem `freeze_reason` devolveu motivo vazio — a recusa ficaria " +
			"sem explicação")
	}
	if !contemFreezeReason(m) {
		t.Errorf("a mensagem não diz onde escrever o motivo: %q", m)
	}
}

func TestCongelado_motivoDeclaradoEhOQueAparece(t *testing.T) {
	sim := false
	c := &Config{Enabled: &sim, FreezeReason: "  o plano 0002 aponta para spec inexistente  "}
	// Espaço em volta não é conteúdo: o texto vai para uma mensagem de terminal.
	if got := c.MotivoDoCongelamento(); got != "o plano 0002 aponta para spec inexistente" {
		t.Errorf("motivo: %q", got)
	}
}

// Motivo só em branco cai no mesmo caso do ausente — senão a recusa sairia com uma linha
// vazia onde deveria estar a explicação.
func TestCongelado_motivoEmBrancoContaComoAusente(t *testing.T) {
	sim := false
	if m := (&Config{Enabled: &sim, FreezeReason: "   \n  "}).MotivoDoCongelamento(); !contemFreezeReason(m) {
		t.Errorf("motivo em branco não caiu no aviso: %q", m)
	}
}

func contemFreezeReason(s string) bool {
	for i := 0; i+13 <= len(s); i++ {
		if s[i:i+13] == "freeze_reason" {
			return true
		}
	}
	return false
}
