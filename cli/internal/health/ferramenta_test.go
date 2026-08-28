package health

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// Gate pulado por ferramenta ausente é indistinguível de gate que passou, para quem só
// olha o verde do `check`. O doctor é o lugar onde essa diferença tem de aparecer — sem
// este aviso o projeto ficaria descoberto exatamente na medida em que ninguém notasse.
func TestDoctorAvisaFerramentaAusenteComoWarn(t *testing.T) {
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "secret-nao-vazado", NeedsTool: "binario-inexistente-xyz", InstallHint: "brew install foo"},
		{Name: "gate-sem-exigencia"},
		{Name: "gate-com-ferramenta", NeedsTool: "sh"},
	}}

	fs := checkFerramentasAusentes(cfg)
	if len(fs) != 1 {
		t.Fatalf("só o gate com ferramenta AUSENTE devia render achado; vieram %d: %+v", len(fs), fs)
	}
	if fs[0].Severity != Warn {
		t.Errorf("a cobertura declarada não é a real — isso é Warn, veio %q", fs[0].Severity)
	}
	if fs[0].Subject != "secret-nao-vazado" {
		t.Errorf("o achado deve apontar o GATE desabilitado, veio %q", fs[0].Subject)
	}
	// O aviso tem de terminar em ação: sem o hint o leitor sabe do problema e não do
	// conserto, e o achado vira só mais uma linha vermelha para conviver.
	if !strings.Contains(fs[0].Detail, "brew install foo") {
		t.Errorf("o aviso deve carregar o install_hint; veio: %q", fs[0].Detail)
	}
}
