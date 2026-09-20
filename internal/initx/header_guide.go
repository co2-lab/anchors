package initx

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
)

// CommentStyle é o dialeto de comentário de uma stack, para os exemplos do
// HEADER_GUIDE semeado. A maioria usa `//`; scripts/config usam `#`; markdown `<!--`.
type CommentStyle struct {
	Line  string // prefixo de comentário de linha (ex: "//", "#")
	Open  string // abertura de bloco, se houver (ex: "<!--")
	Close string // fechamento de bloco (ex: "-->")
}

// commentStyleFor deduz o dialeto principal de um preset pela stack. Conservador:
// `//` para as linguagens C-like (JS/TS/Go/Java/C#/Rust/Dart/PHP/C++), `#` para
// Python/Ruby/Elixir/shell. É só para o EXEMPLO no guide — o Anchors lê qualquer
// comentário em runtime, independente disto.
func commentStyleFor(preset Preset) CommentStyle {
	switch preset.Name {
	case "django", "fastapi", "python-lib", "rails", "phoenix":
		return CommentStyle{Line: "#"}
	default: // node-ts, nextjs, angular, nuxt, expo-rn, spring, dotnet-clean, laravel, go, rust, flutter, cpp…
		return CommentStyle{Line: "//"}
	}
}

// RenderHeaderGuide gera o HEADER_GUIDE.md concreto para o projeto — a régua embutida
// (a doutrina) instanciada com o dialeto de comentário da stack e as features/módulos
// reais detectados. `moduleNames` são as features do projeto (para o exemplo de
// @feature). Vazio → exemplo genérico.
func RenderHeaderGuide(preset Preset, moduleNames []string) string {
	cs := commentStyleFor(preset)
	c := cs.Line

	feat := "auth"
	if len(moduleNames) > 0 {
		feat = moduleNames[0]
	}

	var b strings.Builder
	b.WriteString("# Header guide — " + presetTitleOr(preset) + "\n\n")
	b.WriteString("> The annotation block at the top of EVERY file in this project. Seeded by\n")
	b.WriteString("> `anchors init`; it is the built-in ruler (`anchors guide header`) instantiated for\n")
	b.WriteString("> the stack here. Mandatory: a file without this header is invisible to what\n")
	b.WriteString("> Anchors does best.\n\n")

	b.WriteString("## The block, in this stack's dialect\n\n")
	b.WriteString("In CODE/test/feature (references the spec unit):\n\n```\n")
	fmt.Fprintf(&b, "%s @anchors\n", c)
	fmt.Fprintf(&b, "%s   ref: LGNN             # references the owning unit (the spec); NOT ownership\n", c)
	fmt.Fprintf(&b, "%s   updated_at: 2026-08-08 # day of the last change (the gate checks vs. git)\n", c)
	fmt.Fprintf(&b, "%s   layer: screen         # Structure layer (normally inferred from the path)\n", c)
	fmt.Fprintf(&b, "%s   @feature: %s\n", c, feat)
	b.WriteString("```\n\n")
	b.WriteString("In the SPEC (the OWNER of the identity):\n\n```\n")
	if cs.Open != "" {
		fmt.Fprintf(&b, "%s @anchors\n  code: LGNN            # the spec OWNS the code\n  updated_at: 2026-08-08\n  layer: screen\n%s\n", cs.Open, cs.Close)
	} else {
		fmt.Fprintf(&b, "%s @anchors\n%s   code: LGNN            # the spec OWNS the code\n%s   updated_at: 2026-08-08\n", c, c, c)
	}
	b.WriteString("```\n\n")

	b.WriteString("## The annotations\n\n")
	b.WriteString("- `code:` — OWNERSHIP of the identity (the SPEC is the owner). `ref:` — REFERENCE\n")
	b.WriteString("  (code/feature/test point to the spec unit; it may be multiple: `ref: A, B`). Every\n")
	b.WriteString("  file needs one of the two. Generate the code with `anchors code <name>`.\n")
	b.WriteString("- `updated_at:` — the day of the last change. Whoever changes it updates it; the gate\n")
	b.WriteString("  `updated-at-atual` checks against git (year-month-day only) and `anchors check --fix`\n")
	b.WriteString("  fixes it. Do NOT make up the date — let it match the commit.\n")
	b.WriteString("- `layer:` — the layer; normally inferred from the path, declare it only to override.\n")
	// A tag AGRUPA e nao e' confrontada; a doutrina de produto DECIDE e e' confrontada.
	//
	// Dizer so' "o modulo vertical" deixava as duas parecendo o mesmo eixo, uma delas nao
	// implementada — e desde que `product/` existe, essa leitura custa caro: quem quer
	// que uma regra transversal tenha um lugar escreveria a tag e esperaria um gate que
	// nunca vem.
	b.WriteString("- `@feature: <name>` — a free GROUPING label for the vertical module. ")
	if len(moduleNames) > 0 {
		b.WriteString("In this project: " + strings.Join(moduleNames, ", ") + ".")
	}
	b.WriteString(" Nothing reads it: no gate, no edge.\n")
	b.WriteString("  The vertical axis that IS confronted is product doctrine — `product/<name>.doctrine.md`,\n")
	b.WriteString("  which the spec points at with `@realizes` (see `anchors guide product`).\n")
	b.WriteString("  The tag groups; the doctrine decides. The mechanism grows toward FILTER and VIEW,\n")
	b.WriteString("  never toward a gate that demands the label.\n")
	b.WriteString("- `@noPropagation`, `@anchors-shared-code` — honest opt-outs (always with the why alongside).\n\n")

	b.WriteString("## Rules\n\n")
	b.WriteString("- Always at the TOP of the file.\n")
	b.WriteString("- `code` is the mandatory minimum (gate `header-valid`).\n")
	b.WriteString("- `updated_at` matches the day of the last commit (gate `updated-at-atual`; `--fix` repairs it).\n")
	b.WriteString("- Opt-out always with a why alongside.\n\n")
	// A seção de conformidade não é ornamento: o gate `guide-checklist` a exige, e um
	// guide semeado pelo init que reprova o próprio gate do init é a pior primeira
	// impressão possível — medido num projeto real, foi o primeiro achado bloqueante.
	//
	// Ela existe por um motivo mais fundo que o gate: um guide em prosa é lido e
	// interpretado; um guide com pontos CK é CONFRONTÁVEL. É o que separa "siga o
	// padrão" de "estes cinco itens são verificáveis um a um".
	fmt.Fprintf(&b, "## %s\n\n", i18n.TIn(i18n.Current(), "section.title.compliance_points"))
	b.WriteString("Each item is verifiable in a single file, in isolation. It is what a judgment\n")
	b.WriteString("gate confronts — and what keeps judgment from turning into heuristics.\n\n")
	b.WriteString("- **CK1** — the block is at the TOP of the file, before any code.\n")
	b.WriteString("- **CK2** — there is a `code:` line with a single identity code.\n")
	b.WriteString("- **CK3** — the comment dialect is that of the file's language.\n")
	b.WriteString("- **CK4** — `updated_at`, when present, is the day of the last commit that touched the file.\n")
	b.WriteString("- **CK5** — every opt-out (`@no-…`, `@allow-…`) carries the why on the same line.\n\n")
	b.WriteString("_(Complete, universal ruler: `anchors guide header`.)_\n")
	return b.String()
}

func presetTitleOr(p Preset) string {
	if p.Title != "" {
		return p.Title
	}
	return "project"
}
