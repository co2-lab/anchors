package governance

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
)

// renderSpecGuide gera o SPEC_GUIDE.md do PROJETO — a régua de como escrever uma spec,
// semeada no repositório e não só embutida no binário.
//
// Por que o projeto precisa do arquivo, e não basta `anchors guide spec`:
//
//  1. O guide embutido é do FRAMEWORK; cada projeto tem dialeto (títulos de seção via
//     `section_titles`, vocabulário de letras via `rule_types`, comprimento de código via
//     `code_lengths`). Um agente que lê só o genérico não sabe o que ESTE projeto exige.
//  2. O Anchors cobra que todo artefato tenha guide no repositório (`guide-sem-governo` no
//     doctor). Não distribuir o de spec era o framework exigindo o que ele mesmo não dava.
//  3. Um arquivo no repo é encontrável: o agente lista `guides/`, o humano abre no editor,
//     o `governs` liga ao alvo. Um texto atrás de um subcomando só é achado por quem já
//     sabe que ele existe — e o custo disso foi medido: uma sessão inteira travada
//     "batendo cabeça" porque não sabia o formato da spec, com a informação existindo em
//     `anchors new spec --list-sections`.
//
// A diferença central em relação ao guide embutido: aqui vai um EXEMPLO COMPLETO. O guide
// do framework tem 94 linhas de regra e zero exemplos, e chega a instruir a descobrir o
// formato ERRANDO ("rode `anchors check`; se reprovar, a mensagem do gate diz o formato").
var renderSpecGuide = RenderSpecGuide

// RenderSpecGuide gera o SPEC_GUIDE.md do PROJETO — a régua de como escrever uma spec,
// instanciada com os prefixos e metadados reais do anchors.yaml gerado.
func RenderSpecGuide(cfg *config.Config, exemploCode string) string {
	if exemploCode == "" {
		exemploCode = "LOGI"
	}
	var b strings.Builder

	b.WriteString("# Spec guide — how to write a `.spec.md` in this project\n\n")
	b.WriteString("> Seeded by `anchors init`. It is the built-in ruler (`anchors guide spec`)\n")
	b.WriteString("> instantiated with THIS project's dialect. Read it before writing any spec.\n\n")

	b.WriteString("## Start from the command, not from the text\n\n")
	b.WriteString("Do not write the spec from scratch. The CLI emits a skeleton that already conforms:\n\n")
	b.WriteString("```sh\n")
	b.WriteString("anchors new spec <Name> --out <path>/<Name>.spec.md       # generates the skeleton\n")
	b.WriteString("anchors new spec --list-sections                          # the sections and WHEN to use each\n")
	b.WriteString("```\n\n")
	b.WriteString("The command solves what is easiest to get wrong: it generates the identity code, writes the\n")
	b.WriteString("`@anchors` header in the right dialect, and uses the exact catalogued-rule format.\n")
	b.WriteString("Then fill it in and confront it with `anchors check --changed <file>`.\n\n")

	b.WriteString("## The format the gate requires\n\n")
	b.WriteString("A rule is **catalogued** when it has a code AND a structured place. Three forms\n")
	b.WriteString("are valid, and a loose mention in prose does NOT count:\n\n")
	fmt.Fprintf(&b, "```md\n### %s-B01 — rule description             <- heading (preferred)\n", exemploCode)
	fmt.Fprintf(&b, "| `%s-B02` | description |                    <- table row\n", exemploCode)
	fmt.Fprintf(&b, "- **%s-B03** description                      <- bold bullet\n```\n\n", exemploCode)

	b.WriteString("## Complete example (copy and adapt)\n\n")
	b.WriteString("```md\n")
	fmt.Fprintf(&b, "<!-- @anchors\n  code: %s\n  updated_at: 2026-01-15\n  layer: screen\n-->\n", exemploCode)
	// O EXEMPLO tem de casar com o que o `anchors new spec` gera NESTE projeto: os
	// títulos de seção vêm do mesmo catálogo i18n que o gerador usa (`section.title.*`,
	// resolvido pelo `lang:`), e o rótulo do código idem. Cravá-los aqui mostraria ao
	// leitor um exemplo que o `spec-sections` reprova — o guia contradizendo a ferramenta
	// que ele existe para ensinar, que é o defeito que a régua chama de "duas fontes da
	// própria régua discordando".
	fmt.Fprintf(&b, "# Login — authenticates the user and takes them to the app\n\n> **%s**: `%s`\n\n",
		i18n.T("spec_guide.code_label"), exemploCode)
	fmt.Fprintf(&b, "## %s\n\nEntry screen: takes e-mail and password, authenticates, and navigates to the Home.\n\n",
		sectionTitle("overview"))
	fmt.Fprintf(&b, "## %s\n\n", sectionTitle("rules"))
	fmt.Fprintf(&b, "### %s-S01 — Initial state\nEmpty fields, sign-in button disabled.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-A01 — Sign in with valid credentials\nAuthenticates and navigates to the Home.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-V01 — Invalid e-mail\nThe field shows the message and submit does not fire.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-R01 — Only anonymous may access\nAn active session is redirected to the Home.\n\n", exemploCode)
	fmt.Fprintf(&b, "## %s\n\n| %s | %s | %s |\n| --- | --- | --- |\n\n%s\n",
		sectionTitle("open"), i18n.T("spec_guide.col_question"), i18n.T("spec_guide.col_who"),
		i18n.T("spec_guide.col_becomes"), i18n.T("spec_guide.none"))
	b.WriteString("```\n\n")

	b.WriteString("## The code letter states the rule's NATURE\n\n")
	if len(cfg.RuleTypes) > 0 {
		b.WriteString("This project declares in `rule_types`:\n\n| letter | nature |\n| --- | --- |\n")
		for _, rt := range cfg.RuleTypes {
			fmt.Fprintf(&b, "| `%s` | %s |\n", rt.Letter, rt.Term)
		}
		b.WriteString("\n")
	} else {
		b.WriteString("This project does not declare `rule_types`, so the framework's canonical letters\n")
		b.WriteString("apply: `S` state, `R` permission, `V` validation, `A` action, `X` restriction,\n")
		b.WriteString("`B` behaviour, `N` navigation, `M` message, `D` data. Declaring your own in\n")
		b.WriteString("`rule_types` makes the team's vocabulary apply in place of the generic one.\n\n")
	}

	if len(cfg.CodeLengths) > 0 {
		fmt.Fprintf(&b, "The identity code has %s character(s) in this project (`code_lengths`).\n\n",
			joinInts(cfg.CodeLengths))
	}

	b.WriteString("## The sections\n\n")
	b.WriteString("Three are mandatory — header, overview and rules — plus the open\n")
	b.WriteString("decisions. The rest come in with `--with <key>` when the unit calls for them. Run\n")
	b.WriteString("`anchors new spec --list-sections` for the list with the selection criterion of each\n")
	b.WriteString("one; it includes the mutually exclusive ALTERNATIVES (`contract` or `signature`,\n")
	b.WriteString("`rules` or `effects`), which is where the wrong choice costs a rewrite.\n\n")

	b.WriteString("## What NOT to do\n\n")
	b.WriteString("- **A rule without a code.** Without identity, the feature and the test have nothing to cite —\n")
	b.WriteString("  the triad does not close and the relational gates are left with no target.\n")
	b.WriteString("- **Describing implementation.** The spec states the BEHAVIOUR; the function name and the\n")
	b.WriteString("  library change without the rule changing.\n")
	b.WriteString("- **Repeating the copy.** The user-facing text lives once (in the messages section); the\n")
	b.WriteString("  other sections reference its code.\n")
	fmt.Fprintf(&b, "- **Guessing what is ambiguous.** It becomes a row in *%s* — the gate\n", sectionTitle("open"))
	b.WriteString("  `open-questions-resolved` demands that someone decide, and that is what is wanted.\n\n")

	b.WriteString("## Specialize this file\n\n")
	b.WriteString("It is born generic. As the project settles conventions (spec profiles per kind\n")
	b.WriteString("of unit, real examples, RIGHT/WRONG cases taken from the repo itself), edit\n")
	b.WriteString("here — it is THIS project's ruler, and the `governs` in `anchors.yaml` links this guide to the\n")
	b.WriteString("targets it governs.\n")

	return b.String()
}

// sectionTitle resolve o título de uma seção pelo MESMO catálogo que o gerador de
// artefatos usa, para o exemplo do guia não divergir do que o `anchors new spec` escreve.
func sectionTitle(chave string) string {
	if t := i18n.TIn(i18n.Current(), "section.title."+chave); t != "" {
		return t
	}
	return chave
}

func joinInts(xs []int) string {
	partes := make([]string, len(xs))
	for i, x := range xs {
		partes[i] = fmt.Sprint(x)
	}
	return strings.Join(partes, " or ")
}
