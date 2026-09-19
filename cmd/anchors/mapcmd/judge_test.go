package mapcmd

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O VALOR ANTIGO `dispensado` continua aceito, e vira `waived`.
//
// Ele já está em mapas commitados (`verdict: dispensado`) e em scripts. Recusá-lo faria
// o `check` reperguntar julgamentos que alguém já respondeu — com o carimbo ali, visível
// no arquivo, sem nada os ligando.
func TestVeredito_dispensadoAntigoAindaVale(t *testing.T) {
	if err := ValidateVerdict("dispensado", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor antigo foi recusado: %v", err)
	}
	if err := ValidateVerdict("waived", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor novo foi recusado: %v", err)
	}
}

// E o motivo continua obrigatório nos dois — a compatibilidade não pode afrouxar a régua.
func TestVeredito_dispensadoAntigoTambemExigeMotivo(t *testing.T) {
	if err := ValidateVerdict("dispensado", ""); err == nil {
		t.Error("o valor antigo sem motivo foi aceito — a compatibilidade afrouxou a régua")
	}
}

func TestJudge_aceitaDispensadoEExigeMotivo(t *testing.T) {
	// A validação é a do comando; o que se prova aqui é o CONTRATO dos três vereditos.
	for _, c := range []struct {
		verdict string
		reason  string
		querErr bool
		porque  string
	}{
		{"pass", "", false, "pass sem motivo é aceito (aprovação não precisa de laudo)"},
		{"fail", "", true, "fail sem motivo não é acionável"},
		{"waived", "", true, "dispensado sem a ausência nomeada é indistinguível de gate desligado"},
		{"waived", "a spec declara @TBD: code e MTHRN não existe", false, "dispensado com motivo é aceito"},
		{"inventado", "x", true, "veredito fora dos três é recusado"},
	} {
		err := ValidateVerdict(c.verdict, c.reason)
		if (err != nil) != c.querErr {
			t.Errorf("verdict=%q reason=%q: err=%v, queria erro=%v (%s)",
				c.verdict, c.reason, err, c.querErr, c.porque)
		}
	}
}

// O NOME DO GATE é EXATO, e o nome legado em português não é mais aceito.
//
// Havia uma tabela de alias que convertia `mock-detect-cobre-o-dialeto` para o canônico, e
// este teste provava que ela funcionava. Ela saiu: resolvia e não fechava — o arquivo nunca
// se consertava, e o mapa acumulava carimbos nos dois formatos (medido no blue-eyes: 40
// julgamentos como `regra-cumprida` convivendo com 2 como `rule-fulfilled`).
//
// A conversão virou o passo de migração `1→2`, que roda uma vez. Um projeto que ainda use
// o nome antigo está no formato 1, e é RECUSADO com a mensagem que manda migrar — não lido
// pela metade.
func TestFindJudgmentGate_soOCanonico(t *testing.T) {
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "mock-detect-covers-dialect", Measures: "judgment"},
	}}
	if _, ok := findJudgmentGate(cfg, "mock-detect-covers-dialect"); !ok {
		t.Error("o nome canônico tem de ser encontrado")
	}
	// O legado NÃO resolve mais. Se resolvesse, a tabela teria voltado por algum caminho —
	// e com ela a assimetria que duplicava carimbo.
	if _, ok := findJudgmentGate(cfg, "mock-detect-cobre-o-dialeto"); ok {
		t.Error("o nome legado não deveria resolver: ele é convertido na MIGRAÇÃO, não na leitura")
	}
	if _, ok := findJudgmentGate(cfg, "gate-que-nao-existe"); ok {
		t.Error("um gate não declarado não pode ser encontrado")
	}
}
