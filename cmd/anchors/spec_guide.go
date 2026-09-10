package main

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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
// Ninguém escreve markdown conforme a partir de descrição — copia-se de um exemplo.
func renderSpecGuide(cfg *config.Config, exemploCode string) string {
	if exemploCode == "" {
		exemploCode = "LOGI"
	}
	var b strings.Builder

	b.WriteString("# Guia de spec — como escrever um `.spec.md` neste projeto\n\n")
	b.WriteString("> Semeado por `anchors init`. É a régua embutida (`anchors guide spec`)\n")
	b.WriteString("> instanciada com o dialeto DESTE projeto. Leia antes de escrever qualquer spec.\n\n")

	b.WriteString("## Comece pelo comando, não pelo texto\n\n")
	b.WriteString("Não escreva a spec do zero. O CLI emite o esqueleto já conforme:\n\n")
	b.WriteString("```sh\n")
	b.WriteString("anchors new spec <Nome> --out <caminho>/<Nome>.spec.md   # gera o esqueleto\n")
	b.WriteString("anchors new spec --list-sections                          # as seções e QUANDO usar cada\n")
	b.WriteString("```\n\n")
	b.WriteString("O comando resolve o que é mais fácil errar: gera o código de identidade, escreve o\n")
	b.WriteString("cabeçalho `@anchors` no dialeto certo, e usa o formato exato de regra catalogada.\n")
	b.WriteString("Depois preencha e confronte com `anchors check --changed <arquivo>`.\n\n")

	b.WriteString("## O formato que o gate exige\n\n")
	b.WriteString("Uma regra é **catalogada** quando tem código E lugar estruturado. Três formas\n")
	b.WriteString("valem, e menção solta em prosa NÃO conta:\n\n")
	fmt.Fprintf(&b, "```md\n### %s-B01 — descrição da regra          <- cabeçalho (preferido)\n", exemploCode)
	fmt.Fprintf(&b, "| `%s-B02` | descrição |                     <- linha de tabela\n", exemploCode)
	fmt.Fprintf(&b, "- **%s-B03** descrição                       <- bullet-negrito\n```\n\n", exemploCode)

	b.WriteString("## Exemplo completo (copie e adapte)\n\n")
	b.WriteString("```md\n")
	fmt.Fprintf(&b, "<!-- @anchors\n  code: %s\n  updated_at: 2026-01-15\n  layer: screen\n-->\n", exemploCode)
	fmt.Fprintf(&b, "# Login — autentica o usuário e o leva ao app\n\n> **Código**: `%s`\n\n", exemploCode)
	b.WriteString("## Visão Geral\n\nTela de entrada: recebe e-mail e senha, autentica, e navega para a Home.\n\n")
	b.WriteString("## Regras\n\n")
	fmt.Fprintf(&b, "### %s-S01 — Estado inicial\nCampos vazios, botão de entrar desabilitado.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-A01 — Entrar com credenciais válidas\nAutentica e navega para a Home.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-V01 — E-mail inválido\nO campo mostra a mensagem e o submit não dispara.\n\n", exemploCode)
	fmt.Fprintf(&b, "### %s-R01 — Só anônimo acessa\nSessão ativa é redirecionada para a Home.\n\n", exemploCode)
	b.WriteString("## Decisões em aberto\n\n| Pergunta | Quem decide | Vira |\n| --- | --- | --- |\n\nnenhuma\n")
	b.WriteString("```\n\n")

	b.WriteString("## A letra do código diz a NATUREZA da regra\n\n")
	if len(cfg.RuleTypes) > 0 {
		b.WriteString("Este projeto declara em `rule_types`:\n\n| letra | natureza |\n| --- | --- |\n")
		for _, rt := range cfg.RuleTypes {
			fmt.Fprintf(&b, "| `%s` | %s |\n", rt.Letter, rt.Term)
		}
		b.WriteString("\n")
	} else {
		b.WriteString("Este projeto não declara `rule_types`, então valem as letras canônicas do\n")
		b.WriteString("framework: `S` estado, `R` permissão, `V` validação, `A` ação, `X` restrição,\n")
		b.WriteString("`B` comportamento, `N` navegação, `M` mensagem, `D` dado. Declarar as suas em\n")
		b.WriteString("`rule_types` faz o vocabulário do time valer no lugar do genérico.\n\n")
	}

	if len(cfg.CodeLengths) > 0 {
		fmt.Fprintf(&b, "O código de identidade tem %s caractere(s) neste projeto (`code_lengths`).\n\n",
			joinInts(cfg.CodeLengths))
	}

	b.WriteString("## As seções\n\n")
	b.WriteString("Três são obrigatórias — cabeçalho, visão geral e regras — mais as decisões em\n")
	b.WriteString("aberto. As demais entram com `--with <chave>` quando a unidade pede. Rode\n")
	b.WriteString("`anchors new spec --list-sections` para a lista com o critério de escolha de cada\n")
	b.WriteString("uma; ela inclui as ALTERNATIVAS mutuamente exclusivas (`contract` ou `signature`,\n")
	b.WriteString("`rules` ou `effects`), que é onde a escolha errada custa reescrita.\n\n")

	b.WriteString("## O que NÃO fazer\n\n")
	b.WriteString("- **Regra sem código.** Sem identidade, a feature e o teste não têm o que citar —\n")
	b.WriteString("  a trinca não fecha e os gates relacionais ficam sem alvo.\n")
	b.WriteString("- **Descrever implementação.** A spec diz o COMPORTAMENTO; o nome da função e a\n")
	b.WriteString("  biblioteca mudam sem a regra mudar.\n")
	b.WriteString("- **Repetir a copy.** O texto ao usuário mora uma vez (na seção de mensagens); as\n")
	b.WriteString("  outras seções referenciam o código dela.\n")
	b.WriteString("- **Chutar o que está ambíguo.** Vira linha em *Decisões em aberto* — o gate\n")
	b.WriteString("  `open-questions-resolved` cobra que alguém decida, e é isso que se quer.\n\n")

	b.WriteString("## Especialize este arquivo\n\n")
	b.WriteString("Ele nasce genérico. Conforme o projeto firma convenções (perfis de spec por tipo\n")
	b.WriteString("de unidade, exemplos reais, casos CORRETO/ERRADO tirados do próprio repo), edite\n")
	b.WriteString("aqui — é a régua DESTE projeto, e o `governs` do `anchors.yaml` liga este guide aos\n")
	b.WriteString("alvos que ele rege.\n")

	return b.String()
}

func joinInts(xs []int) string {
	partes := make([]string, len(xs))
	for i, x := range xs {
		partes[i] = fmt.Sprint(x)
	}
	return strings.Join(partes, " ou ")
}
