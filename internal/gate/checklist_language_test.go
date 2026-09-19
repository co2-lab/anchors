package gate

import "testing"

// O `guide-checklist` reconhece a seção de conformidade em QUALQUER idioma do catálogo.
//
// O acoplamento que este teste tranca: o `anchors init` semeia o HEADER_GUIDE.md com o
// título traduzido pelo `lang:` do projeto, e a regex antiga só casava a forma
// portuguesa. Um projeto `lang: en` nasceria reprovando o gate por um guide que o
// próprio init acabara de escrever — o framework cobrando o que ele mesmo não gera.
//
// Sem este teste, a regressão é silenciosa dos dois lados: quem traduzir o guia sem
// tocar no gate quebra o init; quem estreitar o gate de volta quebra os projetos
// traduzidos. Nenhum dos dois aparece no build.
func TestChecklistHeadingReconheceTodosOsIdiomas(t *testing.T) {
	for _, titulo := range []string{
		"Pontos de conformidade", // pt-BR
		"Compliance points",      // en
		"Puntos de conformidad",  // es
	} {
		if !checklistHeadingRE.MatchString("## " + titulo + "\n\n- CK1 item\n") {
			t.Errorf("o gate não reconheceu a seção %q — um projeto nesse idioma reprova por um guide que o init escreveu", titulo)
		}
	}
	// E segue recusando o que NÃO é a seção: aceitar todos os idiomas não pode virar
	// aceitar qualquer título.
	if checklistHeadingRE.MatchString("## Outra coisa\n\n- CK1 item\n") {
		t.Error("o gate aceitou um título que não é a seção de conformidade")
	}
}
