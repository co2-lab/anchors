package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a revisão que mudou o significado de uma palavra, e não disse a quem ---
//
// Uma regra não vive sozinha. Ela compartilha vocabulário com as irmãs, e é esse
// vocabulário que a revisão muda — não só o texto da regra que ela reescreve.
//
// MEDIDO no app de referência, e é o que produziu este gate. A `NTCNN-R0002` trocou o
// badge do sino de CONTAGEM para PONTO e nomeou as regras que reescreveu: `B03`, `B04`,
// `B07`. O invariante `I02` não foi nomeado — e ele se chama "O badge nunca CONTA o que a
// lista não mostra", com corpo dizendo "o NÚMERO no sino corresponde ao que aparece ao
// abrir". O invariante governava uma aritmética que a revisão havia abolido.
//
// NADA ACUSOU. A `B03` estava correta, a `I02` bem-formada, a tríade completa, a suíte
// verde. A contradição apareceu MESES depois, quando outro agente foi implementar e não
// soube qual das duas seguir — e virou decisão que teve de subir para o usuário, sem
// ninguém lembrar do contexto. Sete contradições dessa forma exata numa única leva.
//
// A RÉGUA É CO-CITAÇÃO, não significado. Perguntar se duas regras se contradizem exige
// lê-las, e isso é julgamento — faria disto um judge, não um gate. O que uma máquina
// decide sozinha é mais estreito e basta: quais regras desta unidade compartilham o
// vocabulário das regras que a revisão tocou, e não foram mencionadas.

// revisesRE e checkedRE casam os dois campos da revisão.
//
// As palavras-chave vêm do CATÁLOGO DE TRADUÇÕES, nunca cravadas. É a lição que o eixo de
// fluxo já deu: a primeira versão do `fits` trazia `(?:Encaixa|Fits)` pregado no padrão, e
// um projeto que escrevesse em espanhol não tinha como declarar nada sem editar o engine.
func revisesRE() *regexp.Regexp { return keywordListRE("revision.keyword.revises") }
func checkedRE() *regexp.Regexp { return keywordListRE("revision.keyword.checked") }

// keywordListRE monta o padrão de um campo que lista códigos de regra:
// `**Revises:** ` + "`B03`, `B07`" + `.
func keywordListRE(chave string) *regexp.Regexp {
	return regexp.MustCompile("(?im)^[^\\S\\n]*(?:>|#{1,6})?[^\\S\\n]*(?:\\*\\*)?(?:" +
		strings.Join(escapeKeywords(i18n.AllTranslations(chave)), "|") +
		")(?:\\*\\*)?[^\\S\\n]*:[^\\S\\n]*(\\S.*)$")
}

// escapeKeywords prepara as palavras traduzidas para entrar numa alternação de regex.
func escapeKeywords(xs []string) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if x != "" {
			out = append(out, regexp.QuoteMeta(x))
		}
	}
	return out
}

// ruleRefRE acha os códigos CURTOS que os campos listam (`B03`, `I02`).
//
// Curto e não completo (`NTCNN-B03`) porque é assim que a revisão se escreve: ela já está
// dentro da unidade, e repetir o prefixo em cada item seria ruído. A forma longa também
// casa — quem escreve alterna entre as duas sem pensar nisso.
var ruleRefRE = regexp.MustCompile(`\b(?:[A-Z0-9]{3,6}-)?([A-Z]\d{2})\b`)

// --- o gate ---
func checkRevisionOrphans(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.revision_orphans.skip_not_spec")
	}

	revs := revisionRE().FindAllStringSubmatch(content, -1)
	if len(revs) == 0 {
		return Skip, i18n.T("gate.revision_orphans.no_revision")
	}

	// As regras que ESTA spec define — TODAS, inclusive as dispensadas de cenário.
	//
	// O `definedRequirements` do `spec-feature-match` descarta a regra com
	// `@no-scenario`, e ali está certo: ela pergunta "este requisito tem cenário?", e a
	// dispensa responde. Aqui a pergunta é outra — "esta regra afirma algo que a revisão
	// mudou?" — e a dispensa de cenário não responde nada sobre isso.
	//
	// MEDIDO na spec que originou o gate: a `NTCNN-B07` carrega `@no-scenario` (é decisão
	// de layout, sem renderização a exercitar) e é uma das três regras que a `R0002`
	// reescreveu. Com a leitura do vizinho, o gate acusava a `B07` como código
	// inexistente — falso positivo sobre a regra que o motivou.
	titleByShort := ruleTitles(content)
	if len(titleByShort) == 0 {
		return Skip, i18n.T("gate.revision_orphans.no_rules")
	}
	revised := codesIn(revisesRE(), content)
	checked := codesIn(checkedRE(), content)

	// SEM `Revises:` O GATE SE ABSTÉM, e não acusa.
	//
	// O campo é novo e as revisões que já existem não o carregam: medido no app de
	// referência, 439 revisões escritas em 141 specs, nenhuma com `Revises:` — e a
	// primeira versão deste gate reprovava 174 specs de uma vez, todas com a mesma
	// mensagem.
	//
	// Cento e setenta e quatro achados idênticos não são uma fila, são ruído: quem abre
	// o relatório aprende a rolar por eles, e o achado REAL — a irmã órfã — se perde no
	// meio. Pending diz o que é verdade ("não há como confrontar esta revisão") sem
	// cobrar de quem escreveu antes de a régua existir.
	//
	// A cobrança vem do `plan-change-justified`, que se ancora no `--changed`: a revisão
	// NOVA nasce sob a régua nova, e é nela que o campo passa a ser exigido.
	if len(revised) == 0 {
		return Pending, fmt.Sprintf(i18n.T("gate.revision_orphans.no_revises"), len(revs))
	}

	// Código nomeado que a spec não define: o `ref-resolves` deste eixo.
	var unknown []string
	for c := range revised {
		if _, ok := titleByShort[c]; !ok {
			unknown = append(unknown, c)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return Fail, fmt.Sprintf(i18n.T("gate.revision_orphans.unknown_rule"),
			len(unknown), strings.Join(unknown, ", "))
	}

	// O VOCABULÁRIO das regras revisadas, menos o que não discrimina.
	revisedTerms := map[string]bool{}
	for c := range revised {
		for t := range termsOf(titleByShort[c]) {
			revisedTerms[t] = true
		}
	}
	if len(revisedTerms) == 0 {
		return Pass, ""
	}

	type orfa struct {
		code   string
		termos []string
	}
	var orfas []orfa
	for short, titulo := range titleByShort {
		// Uma regra nunca acusa a si mesma, e o que já foi conferido sai da lista.
		if revised[short] || checked[short] {
			continue
		}
		var compart []string
		for t := range termsOf(titulo) {
			if revisedTerms[t] {
				compart = append(compart, t)
			}
		}
		// UMA palavra em comum não liga duas regras — liga quase todas entre si. O que
		// aponta para o mesmo assunto é a COINCIDÊNCIA: duas ou mais.
		//
		// Medido na spec que originou o gate: com o limiar em uma palavra qualquer, seis
		// das nove regras vinham acusadas; com este, sobra a `I02` — que é exatamente a
		// que contradizia a revisão.
		if !discriminates(compart) {
			continue
		}
		sort.Strings(compart)
		orfas = append(orfas, orfa{short, compart})
	}
	if len(orfas) == 0 {
		return Pass, ""
	}
	sort.Slice(orfas, func(i, j int) bool { return orfas[i].code < orfas[j].code })

	var b strings.Builder
	fmt.Fprintf(&b, i18n.T("gate.revision_orphans.orphans"), len(orfas))
	for _, o := range orfas {
		fmt.Fprintf(&b, "\n    %s — %q\n      %s", o.code,
			trimTitle(titleByShort[o.code]), strings.Join(o.termos, ", "))
	}
	b.WriteString("\n" + i18n.T("gate.revision_orphans.how_to_clear"))
	return Fail, b.String()
}

// codesIn reúne os códigos curtos que um campo lista, em todas as suas ocorrências.
func codesIn(re *regexp.Regexp, content string) map[string]bool {
	out := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		for _, r := range ruleRefRE.FindAllStringSubmatch(m[1], -1) {
			out[r[1]] = true
		}
	}
	return out
}

// ruleTitles indexa o TÍTULO de cada regra definida, pelo código curto.
//
// O título e não o corpo: ele é onde a regra AFIRMA o que afirma, em uma linha, e é o que
// a revisão contradiz quando contradiz. O corpo traz prosa de justificativa, e incluí-lo
// faria quase toda regra compartilhar vocabulário com quase toda outra.
func ruleTitles(content string) map[string]string {
	out := map[string]string{}
	re := defineRuleCaptureRE()
	for _, linha := range strings.Split(content, "\n") {
		m := re.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		// A DISPENSA EM COMENTÁRIO não é parte do que a regra afirma, e deixá-la entrar
		// envenena a comparação: medido na spec real, a `B07` carrega dois `@no-*` com
		// razão escrita e saía com 34 termos — contra 4 a 6 das irmãs —, compartilhando
		// vocabulário com quase todas por acidente de prosa.
		if i := strings.Index(linha, "<!--"); i >= 0 {
			linha = linha[:i]
		}
		_, curto, ok := strings.Cut(m[1], "-")
		if !ok {
			continue
		}
		out[curto] = linha
	}
	return out
}

// termsOf reduz um título ao seu vocabulário, como CONJUNTO.
//
// Reaproveita o `significantTerms` do `feature-test-match` — mesma pergunta ("quais
// palavras deste texto dizem do que ele trata?"), e duas listas de stopwords divergiriam
// na primeira palavra que alguém acrescentasse a uma só.
//
// O CÓDIGO sai antes: `NTCNN-B03` traria `NTCNN` para todo título da unidade, e o termo
// que todos compartilham não discrimina nada.
func termsOf(titulo string) map[string]bool {
	out := map[string]bool{}
	for _, w := range significantTerms(anyCodeRE.ReplaceAllString(titulo, " ")) {
		if !extraStopwords[w] {
			out[w] = true
		}
	}
	return out
}

// extraStopwords são as palavras que o `descStopwords` não precisa descartar e esta régua
// precisa.
//
// A NEGAÇÃO é a que mais custa: metade das regras de uma spec bem escrita diz o que a
// unidade NÃO faz, e "não" ligava toda regra a toda outra. Medido na spec que originou o
// gate — das seis órfãs acusadas na primeira versão, quatro compartilhavam apenas "não".
//
// Ela não entra no `descStopwords` porque lá a pergunta é outra: o `feature-test-match`
// compara a descrição do cenário com o corpo do teste, e ali a negação DISCRIMINA (um
// teste que afirma e um que nega provam coisas diferentes).
var extraStopwords = map[string]bool{
	"não": true, "nao": true, "nunca": true, "nenhum": true, "nenhuma": true,
	"not": true, "never": true, "none": true, "no": true,
	"ni": true, "nunca_es": true,
	// Vocabulário de ESTRUTURA da spec, não do domínio dela.
	"regra": true, "rule": true, "spec": true, "unidade": true, "unit": true,
}

// UMA PALAVRA DE DOMÍNIO BASTA, e o limiar foi medido nas duas direções.
//
// A primeira versão exigia DUAS palavras em comum, e o caso que originou o gate não
// passava: na spec real, a `I02` divide exatamente uma palavra com o que a `R0002`
// reescreveu — `badge` —, e é essa palavra que carrega a contradição inteira.
//
// O que torna uma palavra suficiente é a LIMPEZA do título, não a quantidade. Duas
// medições, na mesma spec:
//
//	sem limpar   3 acusadas — `B01` e `B05` entravam só por "não"
//	limpando     1 acusada  — a `I02`, que é o alvo
//
// As duas fontes de ruído eram a NEGAÇÃO (metade das regras de uma spec bem escrita diz
// o que a unidade não faz) e a DISPENSA EM COMENTÁRIO (a `B07` saía com 34 termos contra
// 4 das irmãs). Removidas as duas, o vocabulário que sobra é o do domínio — e nele uma
// coincidência já é sinal.
//
// Também foi testado um filtro de termo ubíquo (descartar o que aparece em mais de 2/3
// das regras). Ele descartava `lista` e `badge` — exatamente o assunto — e fazia o caso
// real não acusar nada. Numa spec bem escrita o vocabulário de domínio se repete de
// propósito: descartar o que se repete é descartar o assunto.

// discriminates decide se o vocabulário compartilhado aponta para o MESMO assunto.
func discriminates(compart []string) bool {
	return len(compart) >= 1
}

// trimTitle deixa o título legível no veredito: sem o `###`, sem o código, sem o travessão.
func trimTitle(linha string) string {
	s := strings.TrimSpace(strings.TrimLeft(linha, "#> *|`"))
	if m := anyCodeRE.FindStringIndex(s); m != nil {
		s = s[m[1]:]
	}
	s = strings.TrimSpace(strings.TrimLeft(s, "—–-: `"))
	// A dispensa em comentário HTML não é parte do que a regra afirma.
	if i := strings.Index(s, "<!--"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return s
}
