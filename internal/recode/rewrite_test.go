package recode

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestRewrite_headerDialects(t *testing.T) {
	cases := map[string]string{
		"md":      "<!-- @anchors\n  code: TCDTX\n-->\n",
		"ts":      "// @anchors\n//   ref: TCDTX\n",
		"feature": "# @anchors\n#   ref: TCDTX\n",
	}
	for name, in := range cases {
		got, n := Rewrite(in, "TCDTX", "TCTXX")
		if strings.Contains(got, "TCDTX") {
			t.Errorf("%s: sobrou TCDTX: %q", name, got)
		}
		if !strings.Contains(got, "TCTXX") || n == 0 {
			t.Errorf("%s: não trocou (n=%d): %q", name, n, got)
		}
	}
}

func TestRewrite_refListPreservesOthers(t *testing.T) {
	// ref com lista: só o TCDTX muda, o MNDTX fica.
	in := "//   ref: TCDTX, MNDTX\n"
	got, _ := Rewrite(in, "TCDTX", "TCTXX")
	if !strings.Contains(got, "TCTXX") || !strings.Contains(got, "MNDTX") {
		t.Errorf("lista não preservada: %q", got)
	}
	if strings.Contains(got, "TCDTX") {
		t.Errorf("sobrou TCDTX na lista: %q", got)
	}
}

func TestRewrite_scenarioCodeVariety(t *testing.T) {
	in := "TCDTX-B01 TCDTX-S04 TCDTX-A04 TCDTX-FP01b TCDTX-DS-receipt-none TCDTX-VR TCDTX-RA03"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if strings.Contains(got, "TCDTX") {
		t.Errorf("sobrou algum TCDTX-*: %q", got)
	}
	for _, want := range []string{"TCTXX-B01", "TCTXX-FP01b", "TCTXX-DS-receipt-none", "TCTXX-VR", "TCTXX-RA03"} {
		if !strings.Contains(got, want) {
			t.Errorf("faltou %s em: %q", want, got)
		}
	}
	if n != 7 {
		t.Errorf("esperava 7 trocas, veio %d", n)
	}
}

func TestRewrite_bareCrossRef(t *testing.T) {
	in := "Ver a regra em TCDTX (a tela de detalhe)."
	got, _ := Rewrite(in, "TCDTX", "TCTXX")
	if got != "Ver a regra em TCTXX (a tela de detalhe)." {
		t.Errorf("ref nua não trocada: %q", got)
	}
}

func TestRewrite_wordBoundary_doesNotTouchLongerCode(t *testing.T) {
	// O alvo NÃO pode casar dentro de um código maior nem de um vizinho que o contenha.
	// Os exemplos aqui eram de 4 chars e a conversão para 5 os fez colidir com o alvo —
	// `TCDTX-B01` passou a SER o alvo, então o "não deveria tocar" virou "deveria".
	in := "TCDTXX-B01 e XTCDTX e TCDTXABCD"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if got != in || n != 0 {
		t.Errorf("tocou código maior/vizinho (n=%d): %q", n, got)
	}
}

func TestRewrite_idempotentOnNewAbsent(t *testing.T) {
	in := "// @anchors\n//   ref: WXYZX\n\nWXYZ-B01"
	got, n := Rewrite(in, "TCDTX", "TCTXX")
	if got != in || n != 0 {
		t.Errorf("mexeu onde não havia TCDTX (n=%d)", n)
	}
}

func TestFind_classifies(t *testing.T) {
	in := "// @anchors\n//   ref: TCDTX\n\nit('TCDTX-S02: x')\n// ver TCDTX adiante\n"
	occ := Find(in, "TCDTX")
	kinds := map[string]int{}
	for _, o := range occ {
		kinds[o.Kind]++
	}
	if kinds["scenario-code"] != 1 {
		t.Errorf("esperava 1 scenario-code, veio %d (%v)", kinds["scenario-code"], occ)
	}
	if kinds["header"] != 1 {
		t.Errorf("esperava 1 header, veio %d", kinds["header"])
	}
	if kinds["bare-ref"] != 1 {
		t.Errorf("esperava 1 bare-ref, veio %d", kinds["bare-ref"])
	}
}

func TestValidCode(t *testing.T) {
	// O comprimento válido é do PROJETO (`code_lengths`), não do engine — então o teste
	// declara o que exercita, como um projeto faria. Sem isto, o caso depende do default e
	// quebra quando o default muda (foi o que aconteceu: `MN01`, de 4 chars, deixou de
	// valer quando o canônico passou a ser 5).
	config.SetCodeLengths([]int{4, 5})
	t.Cleanup(func() { config.SetCodeLengths([]int{5}) })
	for _, ok := range []string{"TCDTX", "MN01", "ABCDX", "ABCDE"} {
		if !ValidCode(ok) {
			t.Errorf("%s deveria ser válido", ok)
		}
	}
	// Continuam inválidos: fora da faixa de comprimento, minúscula, e caractere que não
	// é alfanumérico maiúsculo.
	for _, bad := range []string{"TCD", "ABCDEF", "tcdt", "TC-T"} {
		if ValidCode(bad) {
			t.Errorf("%s NÃO deveria ser válido", bad)
		}
	}
}
