package initx

import (
	"fmt"
	"strings"
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
	b.WriteString("# Guia de cabeçalho — " + presetTitleOr(preset) + "\n\n")
	b.WriteString("> O bloco de marcações no topo de CADA arquivo deste projeto. Semeado por\n")
	b.WriteString("> `anchors init`; é a régua embutida (`anchors guide header`) instanciada para a\n")
	b.WriteString("> stack aqui. Mandatório: um arquivo sem este cabeçalho é invisível ao que o\n")
	b.WriteString("> Anchors sabe fazer melhor.\n\n")

	b.WriteString("## O bloco, no dialeto desta stack\n\n")
	b.WriteString("No CÓDIGO/teste/feature (referencia a unidade da spec):\n\n```\n")
	fmt.Fprintf(&b, "%s @anchors\n", c)
	fmt.Fprintf(&b, "%s   ref: LGNN             # referencia a unidade dona (a spec); NÃO é posse\n", c)
	fmt.Fprintf(&b, "%s   updated_at: 2026-08-08 # dia da última alteração (o gate confere vs. git)\n", c)
	fmt.Fprintf(&b, "%s   layer: screen         # camada da Estrutura (normalmente deduzida do caminho)\n", c)
	fmt.Fprintf(&b, "%s   @feature: %s\n", c, feat)
	b.WriteString("```\n\n")
	b.WriteString("Na SPEC (a DONA da identidade):\n\n```\n")
	if cs.Open != "" {
		fmt.Fprintf(&b, "%s @anchors\n  code: LGNN            # a spec POSSUI o código\n  updated_at: 2026-08-08\n  layer: screen\n%s\n", cs.Open, cs.Close)
	} else {
		fmt.Fprintf(&b, "%s @anchors\n%s   code: LGNN            # a spec POSSUI o código\n%s   updated_at: 2026-08-08\n", c, c, c)
	}
	b.WriteString("```\n\n")

	b.WriteString("## As marcações\n\n")
	b.WriteString("- `code:` — POSSE da identidade (a SPEC é a dona). `ref:` — REFERÊNCIA (code/\n")
	b.WriteString("  feature/test apontam a unidade da spec; pode ser múltiplo: `ref: A, B`). Todo\n")
	b.WriteString("  arquivo precisa de um dos dois. Gere o código com `anchors code <nome>`.\n")
	b.WriteString("- `updated_at:` — o dia da última alteração. Quem altera atualiza; o gate\n")
	b.WriteString("  `updated-at-atual` confere contra o git (só ano-mês-dia) e `anchors check --fix`\n")
	b.WriteString("  corrige. NÃO invente a data — deixe bater com o commit.\n")
	b.WriteString("- `layer:` — a camada; normalmente deduzida do caminho, declare só p/ sobrepor.\n")
	b.WriteString("- `@feature: <nome>` — o módulo/feature vertical. ")
	if len(moduleNames) > 0 {
		b.WriteString("Neste projeto: " + strings.Join(moduleNames, ", ") + ".\n")
	} else {
		b.WriteString("\n")
	}
	b.WriteString("- `@noPropagation`, `@anchors-shared-code` — opt-outs honestos (sempre com o porquê ao lado).\n\n")

	b.WriteString("## Regras\n\n")
	b.WriteString("- Sempre no TOPO do arquivo.\n")
	b.WriteString("- `code` é o mínimo obrigatório (gate `header-conforme`).\n")
	b.WriteString("- `updated_at` bate com o dia do último commit (gate `updated-at-atual`; `--fix` conserta).\n")
	b.WriteString("- Opt-out sempre com um porquê ao lado.\n\n")
	b.WriteString("_(Régua completa e universal: `anchors guide header`.)_\n")
	return b.String()
}

func presetTitleOr(p Preset) string {
	if p.Title != "" {
		return p.Title
	}
	return "projeto"
}
