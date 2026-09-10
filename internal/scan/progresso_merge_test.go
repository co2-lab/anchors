package scan

import (
	"strings"
	"testing"
)

// O CASO REAL, reproduzido dos três conflitos medidos no blue-eyes (#217, #221 e o
// seguinte): duas branches marcam checkboxes VIZINHOS do mesmo plano, o git não sabe
// resolver, e a resolução à mão é sempre idêntica.
func TestMergeProgress_uneOsDoisLados(t *testing.T) {
	nosso := "- [x] `a.spec.md` — a primeira\n- [ ] `b.spec.md` — a segunda\n"
	deles := "- [ ] `a.spec.md` — a primeira\n- [x] `b.spec.md` — a segunda\n"

	got := MergeProgress(nosso, deles)

	if n := ProgressDone(got); n != 2 {
		t.Errorf("união perdeu marca: %d concluído(s), esperava 2\n%s", n, got)
	}
}

// `[x]` VENCE `[ ]` — marcar concluído é um fato (a spec existe, o check carimbou, o PR
// mergeou). Desmarcar por merge apagaria o fato, e o `progress-honest` passaria a acusar
// um arquivo que existe como não feito.
func TestMergeProgress_marcaNaoEDesfeita(t *testing.T) {
	got := MergeProgress(
		"- [x] `a.spec.md` — feita\n",
		"- [ ] `a.spec.md` — feita\n",
	)
	if !strings.Contains(got, "- [x] `a.spec.md`") {
		t.Errorf("o merge DESMARCOU um item concluído:\n%s", got)
	}
}

// Item que existe só do outro lado não pode sumir: é a spec semeada por uma branch
// enquanto a outra marcava outra coisa.
func TestMergeProgress_naoPerdeItemDeUmLadoSo(t *testing.T) {
	got := MergeProgress(
		"- [x] `a.spec.md` — nossa\n",
		"- [x] `a.spec.md` — nossa\n- [x] `c.spec.md` — só deles\n",
	)
	if !strings.Contains(got, "c.spec.md") {
		t.Errorf("item que só existia do outro lado sumiu:\n%s", got)
	}
	if n := ProgressDone(got); n != 2 {
		t.Errorf("esperava 2 concluídos, veio %d:\n%s", n, got)
	}
}

// Tudo que NÃO é item sobrevive intacto — cabeçalho, prosa, seções. O arquivo é lido por
// gente, e um merge que reescreve a prosa é pior que o conflito que ele evita.
func TestMergeProgress_preservaOQueNaoEItem(t *testing.T) {
	nosso := "# Progresso — plano 0007\n\n> nota que explica algo\n\n## Fase 1\n\n- [ ] `a.spec.md` — a\n"
	deles := "# Progresso — plano 0007\n\n> nota que explica algo\n\n## Fase 1\n\n- [x] `a.spec.md` — a\n"

	got := MergeProgress(nosso, deles)

	for _, esperado := range []string{"# Progresso — plano 0007", "> nota que explica algo", "## Fase 1"} {
		if !strings.Contains(got, esperado) {
			t.Errorf("o merge comeu %q:\n%s", esperado, got)
		}
	}
	if !strings.Contains(got, "- [x] `a.spec.md`") {
		t.Errorf("a marca do outro lado não foi promovida:\n%s", got)
	}
}

// A INDENTAÇÃO do nosso lado é preservada ao promover a marca: itens aninhados existem
// nos planos, e reescrevê-los sem o recuo mudaria a estrutura do documento.
func TestMergeProgress_preservaIndentacao(t *testing.T) {
	got := MergeProgress(
		"  - [ ] `a.spec.md` — aninhada\n",
		"  - [x] `a.spec.md` — aninhada\n",
	)
	if !strings.Contains(got, "  - [x] `a.spec.md`") {
		t.Errorf("a indentação se perdeu na promoção:\n%q", got)
	}
}

// Lados idênticos são idempotentes: o driver roda em todo merge, inclusive nos que não
// tinham conflito nenhum, e não pode introduzir diferença onde não havia.
func TestMergeProgress_ladosIguaisNaoMudamNada(t *testing.T) {
	s := "## Fase 1\n\n- [x] `a.spec.md` — a\n- [ ] `b.spec.md` — b\n"
	if got := MergeProgress(s, s); got != strings.TrimSuffix(s, "\n") && got != s {
		t.Errorf("lados iguais produziram saída diferente:\n%q\n%q", s, got)
	}
	if n := ProgressDone(MergeProgress(s, s)); n != 1 {
		t.Errorf("contagem mudou em merge idempotente: %d", n)
	}
}
