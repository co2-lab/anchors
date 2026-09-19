package initx

import (
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
)

// O guide que o `anchors init` SEMEIA tem de passar no gate que o próprio init declara —
// em qualquer idioma.
//
// O acoplamento: `RenderHeaderGuide` escreve a seção de conformidade traduzida pelo
// `lang:` do projeto, e o `guide-checklist` (internal/gate) a procura por título. Enquanto
// a regex do gate só conhecia a forma portuguesa, um projeto `lang: en` nascia reprovando
// um arquivo que o init acabara de escrever — o framework cobrando o que ele mesmo não
// gera, e a pior primeira impressão possível.
//
// Este teste duplica de propósito a régua do gate (a regex abaixo espelha a
// `checklistHeadingRE`): pôr o confronto AQUI, do lado de quem PRODUZ, é o que faz a
// quebra aparecer para quem mexer no guia — e o teste irmão em internal/gate faz o mesmo
// para quem mexer no gate. As duas pontas do acoplamento, cada uma com seu alarme.
func TestGuiaSemeadoTemAChecklistEmTodoIdioma(t *testing.T) {
	// Espelha checklistHeadingRE: a seção é reconhecida em qualquer idioma do catálogo.
	titulos := i18n.AllTranslations("section.title.compliance_points")
	if len(titulos) == 0 {
		t.Fatal("sem títulos no catálogo para `section.title.compliance_points` — o gate não teria o que casar")
	}
	escapados := make([]string, 0, len(titulos))
	for _, x := range titulos {
		escapados = append(escapados, regexp.QuoteMeta(x))
	}
	secaoRE := regexp.MustCompile(`(?mi)^##+\s+(?:` + strings.Join(escapados, "|") + `)\b`)
	itemRE := regexp.MustCompile(`(?m)\bCK\d+\b`)

	original := i18n.Current()
	t.Cleanup(func() { _ = i18n.Set(original) })

	for _, lang := range i18n.SupportedLangs {
		if err := i18n.Set(lang); err != nil {
			t.Fatalf("idioma %s: %v", lang, err)
		}
		g := RenderHeaderGuide(Preset{}, nil)
		if !secaoRE.MatchString(g) {
			t.Errorf("lang=%s: o guide semeado não traz a seção de conformidade que o `guide-checklist` cobra", lang)
		}
		// A seção sozinha não basta: sem ao menos um ponto CK, o gate de julgamento
		// recai em heurística — que é justamente o que a checklist existe para impedir.
		if !itemRE.MatchString(g) {
			t.Errorf("lang=%s: a seção existe mas não tem ponto CK — nada a confrontar item a item", lang)
		}
	}
}
