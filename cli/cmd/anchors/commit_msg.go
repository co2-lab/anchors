package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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

func assuntoRE() *regexp.Regexp {
	return regexp.MustCompile(`^(` + strings.Join(tiposConvencionais, "|") +
		`)(\([^)]+\))?(!)?: .+`)
}

// LimiteDoAssunto é onde o assunto para de caber onde as pessoas o leem.
//
// 100 é o número do commitlint, e o motivo de adotá-lo não é autoridade: é que as
// ferramentas de changelog e as interfaces de git (git log --oneline, a lista de commits
// do GitHub) já são desenhadas em torno dele. Escolher 72 seria mais "correto" pelo
// costume do git e cortaria assuntos que hoje passam.
const LimiteDoAssunto = 100

// cabecalhoRE separa as partes do assunto para confrontá-las uma a uma. Um `^...$` único
// diria só "não casou", e quem foi barrado precisa saber QUAL parte está errada.
func cabecalhoRE() *regexp.Regexp {
	return regexp.MustCompile(`^([A-Za-z]+)(\(([^)]*)\))?(!)?: *(.*)$`)
}

// problemaNoAssunto devolve o defeito e como consertá-lo, ou "" se o assunto está bom.
//
// A ordem das conferências é do mais estrutural ao mais cosmético: quem errou o tipo não
// precisa ouvir sobre maiúscula na mesma volta.
func problemaNoAssunto(assunto string) string {
	// O commitlint tem `header-trim`, e aqui ele não faz falta: o assunto já chega aparado
	// (ver primeiraLinhaUtil) e o próprio git apara o cabeçalho ao gravar — o espaço não
	// alcança o changelog. Conferir seria acusar um defeito que não existe.
	m := cabecalhoRE().FindStringSubmatch(assunto)
	if m == nil {
		return "não está no formato `tipo(escopo): o que mudou`"
	}
	tipo, temEscopo, escopo, texto := m[1], m[2] != "", m[3], m[5]

	// TIPO em minúsculas: `Feat` e `feat` viram grupos SEPARADOS no changelog, e ninguém
	// percebe até ver a mesma seção duas vezes na mesma versão.
	if tipo != strings.ToLower(tipo) {
		return fmt.Sprintf("o tipo `%s` tem maiúscula — use `%s`, senão o changelog cria "+
			"dois grupos para a mesma coisa", tipo, strings.ToLower(tipo))
	}
	if !tipoConhecido(tipo) {
		return fmt.Sprintf("`%s` não é um tipo conhecido. Tipos: %s",
			tipo, strings.Join(tiposConvencionais, ", "))
	}
	// ESCOPO VAZIO (`feat(): x`) é pior que nenhum: parece que alguém ia dizer algo e
	// parou. Aqui a régua é MAIS estrita que o commitlint, que aceita — e ser mais estrita
	// é seguro: o que passa aqui passa lá, então um projeto que migre para o commitlint
	// não descobre um histórico que a nova ferramenta reprova.
	if temEscopo && strings.TrimSpace(escopo) == "" {
		return "o escopo está vazio — escreva `feat(board):` ou simplesmente `feat:`"
	}
	if strings.TrimSpace(texto) == "" {
		return "o tipo está certo, mas não há assunto depois dos dois-pontos"
	}
	if len(assunto) > LimiteDoAssunto {
		return fmt.Sprintf("o assunto tem %d caracteres e o limite é %d — ele é cortado na "+
			"lista de commits e no changelog. O detalhe vai no CORPO da mensagem, que não "+
			"tem limite", len(assunto), LimiteDoAssunto)
	}
	// PONTO FINAL: o changelog junta o assunto a marcadores e links, e a frase termina
	// com dois pontos.
	if strings.HasSuffix(texto, ".") {
		return "o assunto termina em ponto — o changelog o emenda a outros textos e o ponto sobra"
	}
	return ""
}

func tipoConhecido(t string) bool {
	for _, v := range tiposConvencionais {
		if v == t {
			return true
		}
	}
	return false
}

// mensagensQueOGitGera não são escritas por ninguém — barrá-las quebraria operações
// normais do git em vez de melhorar o histórico.
func geradaPeloGit(assunto string) bool {
	for _, p := range []string{"Merge ", "Revert ", "fixup! ", "squash! ", "Reapply "} {
		if strings.HasPrefix(assunto, p) {
			return true
		}
	}
	return false
}

func newCommitMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-msg <arquivo>",
		Short: "Confronta a mensagem de commit com o formato que o changelog vai ler",
		Long: `Lê o arquivo de mensagem (o que o git passa ao hook commit-msg) e confronta
o assunto com o Conventional Commits.

O changelog nasce dos commits. Um assunto fora do formato não some do histórico —
some do CHANGELOG, e só se descobre quando alguém vai gerar a primeira versão e o
que faltou já está a centenas de commits de distância.

Mensagens geradas pelo git (merge, revert, fixup) passam: ninguém as escreveu.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("ler a mensagem: %w", err)
			}
			assunto := primeiraLinhaUtil(string(b))
			if assunto == "" {
				return nil // mensagem vazia: quem barra é o git, e a mensagem dele é melhor
			}
			if geradaPeloGit(assunto) {
				return nil
			}
			problema := problemaNoAssunto(assunto)
			if problema == "" {
				return nil
			}
			cmd.SilenceUsage = true
			return fmt.Errorf("%s:\n"+
				"    %s\n\n"+
				"  Use `tipo(escopo): o que mudou` — o escopo é opcional, o `!` marca quebra:\n"+
				"    feat(board): a contagem do refresh aparece no botão\n"+
				"    fix: o gate acusava quem não mudou\n"+
				"    feat!: `--changed` passa a exigir caminho relativo\n\n"+
				"  Tipos: %s\n\n"+
				"  Isto não é estética: o changelog nasce dos commits, e o que está fora do\n"+
				"  formato não aparece nele. Commit já feito não se conserta.",
				problema, assunto, strings.Join(tiposConvencionais, ", "))
		},
	}
	return cmd
}

// primeiraLinhaUtil devolve o assunto, pulando os comentários que o git põe no arquivo.
func primeiraLinhaUtil(s string) string {
	for _, l := range strings.Split(s, "\n") {
		t := strings.TrimSpace(l)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		return t
	}
	return ""
}
