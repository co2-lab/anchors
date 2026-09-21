package flagx

import (
	"os"
	"path/filepath"
	"testing"
)

const exemplo = "<!-- @anchors\n  code: CHKUT\n-->\n" + `# Flag: novo-checkout

## Cenários

| Cenário | Quando o valor | Então |
| --- | --- | --- |
| ` + "`CHKUT-G01`" + ` | ` + "`= \"off\"`" + ` | o checkout antigo responde |
| ` + "`CHKUT-G02`" + ` | ` + "`= \"on\"`" + ` | o checkout novo responde |
| ` + "`CHKUT-G03`" + ` | ` + "`>= 50`" + ` | o rollout decide por usuário |
| ` + "`CHKUT-G04`" + ` | ausente | vale o default do cabeçalho |
`

func TestParse_leOsCenarios(t *testing.T) {
	f := parse(exemplo, "flags/novo-checkout.flag.md")
	if len(f.Scenarios) != 4 {
		t.Fatalf("len(Scenarios) = %d, esperado 4", len(f.Scenarios))
	}
	if f.Code != "CHKUT" {
		t.Errorf("Code = %q, esperado CHKUT", f.Code)
	}
	if f.Name != "novo-checkout" {
		t.Errorf("Name = %q, esperado novo-checkout", f.Name)
	}
	if got := f.Scenarios[0]; got.Cond.Op != OpEq || got.Cond.Operand != "off" {
		t.Errorf("G01 = %q/%q, esperado eq/off", got.Cond.Op, got.Cond.Operand)
	}
	if got := f.Scenarios[2]; got.Cond.Op != OpGte || got.Cond.Operand != "50" {
		t.Errorf("G03 = %q/%q, esperado gte/50", got.Cond.Op, got.Cond.Operand)
	}
	if got := f.Scenarios[1].Then; got != "o checkout novo responde" {
		t.Errorf("G02.Then = %q", got)
	}
}

// A linha é registrada para o veredito poder apontar onde.
func TestParse_registraALinha(t *testing.T) {
	f := parse(exemplo, "flags/novo-checkout.flag.md")
	if f.Scenarios[0].Line == 0 {
		t.Error("Line não foi registrada — o veredito não teria onde apontar")
	}
	if f.Scenarios[1].Line <= f.Scenarios[0].Line {
		t.Error("as linhas não são crescentes")
	}
}

// O caso AUSENTE é a pergunta que a flag mais esquece.
func TestAbsent(t *testing.T) {
	f := parse(exemplo, "x.flag.md")
	if !f.Absent() {
		t.Error("o exemplo declara `ausente` e Absent() disse que não")
	}

	semAusente := "| `CHKUT-G01` | `= \"on\"` | liga |\n"
	if parse(semAusente, "x.flag.md").Absent() {
		t.Error("não há cenário ausente e Absent() disse que sim")
	}
}

// Uma condição que a gramática recusa vira ACHADO, não desaparece. Linha que some é
// exatamente o silêncio que os gates existem para acabar.
func TestParse_condicaoInvalidaViraAchadoNaoSumico(t *testing.T) {
	src := "| `CHKUT-G01` | quando o usuário for beta | liga |\n"
	f := parse(src, "x.flag.md")
	if len(f.Scenarios) != 1 {
		t.Fatalf("a linha sumiu: len = %d, esperado 1", len(f.Scenarios))
	}
	if f.Scenarios[0].Err == nil {
		t.Error("a condição é prosa e foi aceita — a gramática fixa não recusou nada")
	}
}

// Só a letra `G` é cenário de flag. Uma linha de outra letra na mesma tabela não entra.
func TestParse_soALetraG(t *testing.T) {
	src := "| `CHKUT-B01` | `= \"on\"` | não é cenário de flag |\n" +
		"| `CHKUT-G01` | `= \"on\"` | este é |\n"
	f := parse(src, "x.flag.md")
	if len(f.Scenarios) != 1 {
		t.Fatalf("len = %d, esperado 1 — só a letra G", len(f.Scenarios))
	}
	if f.Scenarios[0].Code != "CHKUT-G01" {
		t.Errorf("veio %q", f.Scenarios[0].Code)
	}
}

// Projeto sem `flags/` não é erro: é projeto que ainda não declarou flag.
func TestLoad_semPastaNaoEErro(t *testing.T) {
	fs, err := Load(t.TempDir())
	if err != nil || fs != nil {
		t.Errorf("Load sem flags/ = (%v, %v), esperado (nil, nil)", fs, err)
	}
}

func TestLoad_ordemEstavel(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"zeta", "alfa", "meio"} {
		src := "| `" + map[string]string{"zeta": "ZETAA", "alfa": "ALFAA", "meio": "MEIOO"}[n] +
			"-G01` | `= \"on\"` | liga |\n"
		if err := os.WriteFile(filepath.Join(dir, n+Suffix), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fs, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 3 {
		t.Fatalf("len = %d", len(fs))
	}
	if fs[0].Name != "alfa" || fs[1].Name != "meio" || fs[2].Name != "zeta" {
		t.Errorf("ordem instável: %q, %q, %q", fs[0].Name, fs[1].Name, fs[2].Name)
	}
}

func TestByCode(t *testing.T) {
	idx := ByCode([]Flag{parse(exemplo, "x.flag.md")})
	if s, ok := idx["CHKUT-G03"]; !ok || s.Cond.Op != OpGte {
		t.Errorf("ByCode não resolveu CHKUT-G03: %v/%v", s.Cond.Op, ok)
	}
	if _, ok := idx["CHKUT-G99"]; ok {
		t.Error("ByCode resolveu um código que não existe")
	}
}
