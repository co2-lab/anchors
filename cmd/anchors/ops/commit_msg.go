// @anchors
//   code: CMCCM
//   ref: CMMSC

package ops

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/spf13/cobra"
)

// --- a mensagem de commit é matéria-prima do changelog ---
//
// O changelog nasce dos COMMITS, e por isso o padrão tem de valer ANTES de ele existir:
// commit já feito não se conserta, e um histórico onde metade das mensagens não casa o
// formato produz um changelog com buracos que ninguém consegue preencher depois.
//
// Medido no projeto de referência: o squash do PR do plano de mutação entrou como
// "[MTUAO] Plano 0017 — mutação", porque o GitHub usa o TÍTULO DO PR como mensagem do
// squash e o título estava no formato do card. Esse commit — o que introduziu o plano —
// simplesmente não apareceria no changelog.
//
// A régua é o Conventional Commits, e não uma invenção nossa: é o que as ferramentas de
// changelog já sabem ler. Inventar formato aqui daria trabalho duas vezes.

// O QUE FOI APRENDIDO DO COMMITLINT, e o que foi deixado de fora.
//
// O commitlint é a ferramenta madura desta régua, e confrontar a implementação com a dele
// mostrou quatro conferências que faltavam aqui: tipo em minúsculas, escopo não vazio,
// tamanho do assunto e ponto final. Todas protegem o changelog de um jeito concreto —
// `Feat` e `feat` viram grupos separados, o assunto longo é cortado na lista de commits.
//
// UMA regra dele fica DE FORA: `subject-case`, que barra assunto começando com maiúscula.
// Medido: ela reprova `feat: SBOM sai da pasta ignorada` — não distingue sigla legítima
// de "frase capitalizada". Num projeto que fala de SBOM, CI e PR isso é atrito sem ganho,
// e a saída barata para quem é barrado sem razão é desligar o hook.
//
// A régua embutida existe para o projeto que NÃO tem Node: o commitlint exige um runtime
// inteiro para validar uma linha de texto, e o Anchors se propõe a funcionar em qualquer
// stack. Onde o Node já está, usar o commitlint é a escolha melhor — mais regras, mais
// configurável, e é o que as ferramentas de changelog esperam encontrar.

// tiposConvencionais são os tipos aceitos, e o que cada um significa PARA O CHANGELOG.
//
// `feat` e `fix` são os que aparecem para o usuário; o resto é histórico interno, que a
// maioria das ferramentas agrupa ou omite. A lista é fechada de propósito: tipo livre
// vira sinônimo ("bugfix", "hotfix", "correção") e o agrupamento se desfaz.
var tiposConvencionais = []string{
	"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert",
}

// LimiteDoAssunto é onde o assunto para de caber onde as pessoas o leem.
//
// 100 é o número do commitlint, e o motivo de adotá-lo não é autoridade: é que as
// ferramentas de changelog e as interfaces de git (git log --oneline, a lista de commits
// do GitHub) já são desenhadas em torno dele. Escolher 72 seria mais "correto" pelo
// costume do git e cortaria assuntos que hoje passam.
const LimiteDoAssunto = 100

// headerRE separa as partes do assunto para confrontá-las uma a uma. Um `^...$` único
// diria só "não casou", e quem foi barrado precisa saber QUAL parte está errada.
func headerRE() *regexp.Regexp {
	return regexp.MustCompile(`^([A-Za-z]+)(\(([^)]*)\))?(!)?:( *)(.*)$`)
}

// subjectProblem devolve o defeito e como consertá-lo, ou "" se o assunto está bom.
//
// A ordem das conferências é do mais estrutural ao mais cosmético: quem errou o tipo não
// precisa ouvir sobre maiúscula na mesma volta.
func subjectProblem(assunto string) string {
	// O commitlint tem `header-trim`, e aqui ele não faz falta: o assunto já chega aparado
	// (ver primeiraLinhaUtil) e o próprio git apara o cabeçalho ao gravar — o espaço não
	// alcança o changelog. Conferir seria acusar um defeito que não existe.
	m := headerRE().FindStringSubmatch(assunto)
	if m == nil {
		return "is not in the `type(scope): what changed` format"
	}
	tipo, temEscopo, escopo, espaco, texto := m[1], m[2] != "", m[3], m[5], m[6]

	// TIPO em minúsculas: `Feat` e `feat` viram grupos SEPARADOS no changelog, e ninguém
	// percebe até ver a mesma seção duas vezes na mesma versão.
	if tipo != strings.ToLower(tipo) {
		return fmt.Sprintf("the type `%s` has an uppercase letter — use `%s`, otherwise the changelog creates "+
			"two groups for the same thing", tipo, strings.ToLower(tipo))
	}
	if !knownType(tipo) {
		return fmt.Sprintf("`%s` is not a known type. Types: %s",
			tipo, strings.Join(tiposConvencionais, ", "))
	}
	// ESCOPO VAZIO (`feat(): x`) é pior que nenhum: parece que alguém ia dizer algo e
	// parou. Aqui a régua é MAIS estrita que o commitlint, que aceita — e ser mais estrita
	// é seguro: o que passa aqui passa lá, então um projeto que migre para o commitlint
	// não descobre um histórico que a nova ferramenta reprova.
	if temEscopo && strings.TrimSpace(escopo) == "" {
		return "the scope is empty — write `feat(board):` or simply `feat:`"
	}
	if strings.TrimSpace(texto) == "" {
		return "the type is right, but there is no subject after the colon"
	}
	// SEM ESPAÇO depois dos dois-pontos (`feat:x`): o parser do Conventional Commits não
	// separa o tipo do assunto, e o commitlint recusa. Só a regex morta dos testes pegava
	// isto; a régua real deixava passar — divergindo para o lado frouxo (CMMSC-I01).
	if espaco == "" {
		return "there is no space after the colon — write `feat: what changed`, otherwise the changelog " +
			"cannot tell the type from the subject"
	}
	// CHARACTERS, not bytes: `len` counted `ç` and `ã` twice, and an accented subject of 60
	// letters was refused as if it had 120.
	if n := utf8.RuneCountInString(assunto); n > LimiteDoAssunto {
		return fmt.Sprintf("the subject has %d characters and the limit is %d — it is cut off in the "+
			"commit list and in the changelog. The detail goes in the BODY of the message, which has "+
			"no limit", n, LimiteDoAssunto)
	}
	// PONTO FINAL: o changelog junta o assunto a marcadores e links, e a frase termina
	// com dois pontos.
	if strings.HasSuffix(texto, ".") {
		return "the subject ends in a period — the changelog splices it to other texts and the period is left over"
	}
	return ""
}

func knownType(t string) bool {
	for _, v := range tiposConvencionais {
		if v == t {
			return true
		}
	}
	return false
}

// mensagensQueOGitGera não são escritas por ninguém — barrá-las quebraria operações
// normais do git em vez de melhorar o histórico.
func generatedByGit(assunto string) bool {
	for _, p := range []string{"Merge ", "Revert ", "fixup! ", "squash! ", "Reapply "} {
		if strings.HasPrefix(assunto, p) {
			return true
		}
	}
	return false
}

func newCommitMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-msg <file>",
		Short: "Confront the commit message with the format the changelog will read",
		Long: `Reads the message file (the one git passes to the commit-msg hook) and confronts
the subject with Conventional Commits.

The changelog is born from the commits. A subject outside the format does not vanish from the history —
it vanishes from the CHANGELOG, and it is only discovered when someone goes to generate the first version and
what is missing is already hundreds of commits away.

Messages generated by git (merge, revert, fixup) pass: nobody wrote them.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read the message: %w", err)
			}
			assunto := firstUsefulLine(string(b))
			if assunto == "" {
				return nil // mensagem vazia: quem barra é o git, e a mensagem dele é melhor
			}
			if generatedByGit(assunto) {
				return nil
			}
			problema := subjectProblem(assunto)
			if problema == "" {
				if p := bugFooterProblem(string(b), assunto); p != "" {
					cmd.SilenceUsage = true
					return fmt.Errorf("%s\n\n"+
						"  A bug is a defect that SHIPPED — reached a release or production, where someone could\n"+
						"  hit it. Its fix is a `fix` commit with a footer saying where it was seen:\n"+
						"    fix(board): the refresh count no longer doubles\n\n"+
						"    <what failed, where, how it was measured>\n\n"+
						"    Bug: v0.1.204 — card #812\n\n"+
						"  A correction of work that never shipped is a `fix` WITHOUT the footer: it is not a bug,\n"+
						"  and the changelog's bug list would name something no one ever saw.", p)
				}
				return nil
			}
			cmd.SilenceUsage = true
			return fmt.Errorf("%s:\n"+
				"    %s\n\n"+
				"  Use `type(scope): what changed` — the scope is optional, the `!` marks a break:\n"+
				"    feat(board): the refresh count appears on the button\n"+
				"    fix: the gate accused whoever did not change\n"+
				"    feat!: `--changed` now requires a relative path\n\n"+
				"  Types: %s\n\n"+
				"  This is not aesthetics: the changelog is born from the commits, and what is outside the\n"+
				"  format does not appear in it. A commit already made cannot be fixed.",
				problema, assunto, strings.Join(tiposConvencionais, ", "))
		},
	}
	return cmd
}

// firstUsefulLine devolve o assunto, pulando os comentários que o git põe no arquivo.
func firstUsefulLine(s string) string {
	for _, l := range strings.Split(s, "\n") {
		t := strings.TrimSpace(l)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		return t
	}
	return ""
}

// bugFooterRE is a `Bug:` footer line, in any casing, so a misspelling is caught too.
var bugFooterRE = regexp.MustCompile(`(?i)^(bug)\s*:(.*)$`)

// bugFooterProblem confronts the `Bug:` footer, when the message has one.
//
// A BUG is a defect that shipped — reached a release or production. A FIX is any
// correction, including one of something broken in the same work that never reached
// anyone. The type cannot tell them apart (both are `fix`), and a changelog listing every
// `fix` as a bug fixed would name corrections no user ever saw. The `Bug:` footer is the
// marker: on a `fix`, it says the defect shipped and where it was seen, and it is what a
// search of the history and a changelog's bug list read.
//
// Only the format is checked. Whether the defect shipped the hook cannot know, so a `fix`
// without the footer passes. The footer is read as the changelog reads it: the trailing
// paragraphs made only of `Key: value` lines, so a `Co-Authored-By:` block after the
// `Bug:` does not hide it; a paragraph of prose that begins with "Bug:" is not a footer.
func bugFooterProblem(msg, subject string) string {
	// The footer is read as the changelog reads it (`changelog.FooterLines`), from the
	// body only: the subject's paragraph — up to the first blank line — is never a footer.
	var body []string
	inSubject, started := true, false
	for _, l := range strings.Split(strings.ReplaceAll(msg, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "#") {
			continue
		}
		if inSubject {
			started = started || t != ""
			inSubject = !(started && t == "")
			continue
		}
		body = append(body, l)
	}
	for _, l := range changelog.FooterLines(strings.Join(body, "\n")) {
		m := bugFooterRE.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		if m[1] != "Bug" {
			return fmt.Sprintf("the footer is written `%s:` — write `Bug:`; the changelog and a search of the "+
				"history find it by that exact token", m[1])
		}
		if strings.TrimSpace(m[2]) == "" {
			return "the `Bug:` footer is empty — say where the bug was seen: the release or version, the card, " +
				"the report"
		}
		if h := headerRE().FindStringSubmatch(subject); h == nil || h[1] != "fix" {
			tipo := subject
			if h != nil {
				tipo = h[1]
			}
			return fmt.Sprintf("a `Bug:` footer on a `%s` commit — only a `fix` fixes a bug, and the footer "+
				"would list this commit among the bugs fixed", tipo)
		}
	}
	return ""
}
