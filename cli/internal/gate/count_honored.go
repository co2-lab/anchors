package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// count-honored: um NÚMERO que a spec afirma sobre o código precisa bater com o código.
//
// A "âncora que mente" tem uma variante numérica, e ela é a que envelhece sozinha: a spec
// diz "os 50 modelos do produto", alguém adiciona o 51º, e a frase vira mentira sem que
// ninguém a tenha editado. Medido num projeto real: um único arquivo de spec continha 7
// afirmações numéricas, e adicionar UM modelo tornava 10 frases obsoletas de uma vez.
//
// Diferente das outras mentiras de âncora, esta não depende de descuido — depende só do
// tempo. Toda contagem escrita em prosa é uma dívida com juros.
//
// O gate NÃO adivinha o que contar: "51 modelos", "51 cláusulas de autorização" e "46
// modelos indexados" são três perguntas diferentes sobre os mesmos arquivos, e nenhuma
// heurística acerta as três. A spec DECLARA como conferir, no mesmo padrão da Tabela de
// Dependências (o que está entre crases vira contrato):
//
//	<!-- @anchors-count: 51 = amplify/data/models/*.ts -->
//	<!-- @anchors-count: 51 cláusulas = amplify/data/models/*.ts /allow\.(owner|group)/ -->
//
// Sem padrão declarado, conta ARQUIVOS que casam o glob; com padrão, conta OCORRÊNCIAS
// dele. Uma contagem não declarada é ignorada — o gate cobra o que foi prometido, não
// todo número que aparece no texto (um "90 dias" de retenção não é contagem de código).
func checkCountHonored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a afirmação numérica é da spec — é ela que descreve o código"
	}

	decls := contagensDeclaradas(content)
	if len(decls) == 0 {
		// Ausência não é falha: nem toda spec afirma número, e exigir a marcação de toda
		// spec seria ritual. O gate cobra quem DECLAROU como conferir — e quem escreve uma
		// contagem sem declarar assume que ela envelheça.
		return Skip, "a spec não declara nenhuma contagem confrontável " +
			"(`<!-- @anchors-count: N = <glob> [/regex/] -->`)"
	}

	var erros []string
	for _, d := range decls {
		real, err := contar(root, d)
		if err != nil {
			erros = append(erros, fmt.Sprintf("`%s`: %v", d.Glob, err))
			continue
		}
		// A PROSA ao lado do marcador também afirma o número, e é ela que o leitor lê.
		// Achado real: o marcador dizia 51 (certo), e uma frase três linhas abaixo dizia
		// "50 cláusulas de autorização" — o gate ficava ✓ e a spec mentia para quem a
		// abrisse. Conferir só a marcação é conferir a metade que ninguém lê.
		if real == d.Esperado && d.Rotulo != "" {
			if divergentes := prosaDivergente(content, d, real); len(divergentes) > 0 {
				erros = append(erros, fmt.Sprintf(
					"o marcador de `%s` bate (%d), mas a PROSA diz %s — é a frase que o leitor "+
						"lê, e ela está errada", d.Rotulo, real, strings.Join(divergentes, " e ")))
			}
			continue
		}
		if real != d.Esperado {
			rotulo := d.Rotulo
			if rotulo == "" {
				rotulo = "itens"
			}
			erros = append(erros, fmt.Sprintf(
				"a spec afirma **%d %s**, o código tem **%d** (`%s`%s)",
				d.Esperado, rotulo, real, d.Glob, sufixoPadrao(d.Padrao)))
		}
	}
	if len(erros) == 0 {
		return Pass, ""
	}
	sort.Strings(erros)
	return Fail, fmt.Sprintf("%d contagem(ns) desatualizada(s): %s. "+
		"Um número em prosa envelhece sozinho — ninguém precisa errar para a frase virar "+
		"mentira, basta o código crescer. Atualize o número (e as frases que dependem dele) "+
		"ou corrija o que a contagem confronta",
		len(erros), strings.Join(erros, "; "))
}

// contagem declarada na spec.
type contagem struct {
	Esperado int
	Rotulo   string // "modelos", "cláusulas" — só para a mensagem
	Glob     string
	Padrao   string // regex opcional: conta OCORRÊNCIAS em vez de arquivos
}

// countRE casa `@anchors-count: 51 modelos = <glob> /regex/`. O rótulo e o regex são
// opcionais; o número e o glob, não.
// O glob CONTÉM barras (`models/*.ts`), então não pode ser delimitado por `[^\s/]`. A
// separação glob↔regex é feita pelo espaço antes da `/` de abertura: `<glob> /<regex>/`.
var countRE = regexp.MustCompile(`@anchors-count:\s*(\d+)\s*([^=\n]*?)\s*=\s*(\S+)(?:\s+/(.+?)/)?\s*(?:-->)?\s*$`)

func contagensDeclaradas(content string) []contagem {
	var out []contagem
	for _, linha := range strings.Split(content, "\n") {
		m := countRE.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		out = append(out, contagem{
			Esperado: n,
			Rotulo:   strings.TrimSpace(m[2]),
			Glob:     strings.TrimSpace(m[3]),
			Padrao:   m[4],
		})
	}
	return out
}

// contar resolve a contagem contra o disco: arquivos que casam o glob, ou ocorrências do
// padrão dentro deles.
func contar(root string, d contagem) (int, error) {
	arquivos, err := doublestar.Glob(os.DirFS(root), d.Glob)
	if err != nil {
		return 0, fmt.Errorf("glob inválido: %w", err)
	}
	if len(arquivos) == 0 {
		// Zero arquivos quase nunca é "o código tem zero": é o glob apontando para o lugar
		// errado. Dizer isso evita a conclusão errada — e a correção é o glob, não o número.
		return 0, fmt.Errorf("nenhum arquivo casa este glob — confira o caminho antes de mexer no número")
	}
	if d.Padrao == "" {
		return len(arquivos), nil
	}
	re, err := regexp.Compile(d.Padrao)
	if err != nil {
		return 0, fmt.Errorf("regex inválido: %w", err)
	}
	total := 0
	for _, f := range arquivos {
		b, rerr := os.ReadFile(filepath.Join(root, f))
		if rerr != nil {
			continue
		}
		total += len(re.FindAllIndex(b, -1))
	}
	return total, nil
}

func sufixoPadrao(p string) string {
	if p == "" {
		return ""
	}
	return ", contando `/" + p + "/`"
}

// prosaDivergente acha, no texto, afirmações "<N> <rótulo>" cujo número não bate com o
// real. Só olha o MESMO rótulo que a marcação declara — é o que evita acusar o "90 dias"
// de retenção que nada tem a ver com contagem de código.
func prosaDivergente(content string, d contagem, real int) []string {
	rot := regexp.QuoteMeta(d.Rotulo)
	// `\*{0,2}` aceita a forma em negrito (`**50 modelos**`), que é como o número
	// costuma aparecer quando o autor quer destacá-lo.
	// O QUALIFICADOR é o que separa a afirmação sobre o TOTAL de uma sobre um
	// SUBCONJUNTO: "os 50 modelos" fala do todo; "nos 46 modelos INDEXADOS" e "em 16
	// modelos DE DADO FINANCEIRO" falam de recortes, e acusá-los seria o ruído que faz
	// desligar o gate (aconteceu: 2 falsos positivos na primeira execução).
	//
	// A regra: se vem palavra depois do rótulo antes da pontuação, é recorte — a menos
	// que seja ligação ("de autorização", "do produto"), que não restringe nada.
	re := regexp.MustCompile(`\*{0,2}(\d+)\*{0,2}\s+` + rot + `\b([^\n.,;:)]*)`)
	// A própria linha de marcação contém "<N> <rótulo>" e não é prosa — lê-la faria o
	// gate acusar a si mesmo sempre que o número estivesse desatualizado (e o teste pegou).
	semMarcacao := countRE.ReplaceAllString(content, "")
	vistos := map[string]bool{}
	var out []string
	for _, m := range re.FindAllStringSubmatch(semMarcacao, -1) {
		n, err := strconv.Atoi(m[1])
		if err != nil || n == real || vistos[m[1]] {
			continue
		}
		if qualificado(m[2]) {
			continue // fala de um recorte, não do total
		}
		vistos[m[1]] = true
		out = append(out, "**"+m[1]+"**")
	}
	sort.Strings(out)
	return out
}

// qualificado diz se a frase fala de um SUBCONJUNTO em vez do total.
//
// Distinguir isso é SEMÂNTICA, e o gate só resolve o caso inequívoco. Medido no texto
// real: "50 cláusulas DE AUTORIZAÇÃO" fala do total, "16 modelos DE DADO FINANCEIRO" fala
// de um recorte — as duas começam igual, e nenhum regex separa as duas sem entender o que
// as palavras significam.
//
// Então a régua é conservadora: QUALQUER palavra depois do rótulo torna a frase suspeita
// de recorte, e o gate se cala. Ele confronta apenas a forma nua ("os 50 modelos", "51
// cláusulas."), que é onde a afirmação é inequivocamente sobre o total.
//
// O custo é perder achados; o custo de errar para o outro lado é acusar quem está certo,
// e um gate que acusa a maioria é desligado no primeiro dia — levando junto os que
// funcionam. (QUALITY: "antes de escrever um gate, meça o falso-positivo contra o
// repositório real".)
func qualificado(resto string) bool {
	// o negrito pode fechar DEPOIS do complemento (`**3 cláusulas de autorização**`),
	// então os asteriscos não fazem parte do complemento.
	r := strings.TrimSpace(strings.Trim(strings.TrimSpace(resto), "*"))
	if r == "" {
		return false // forma nua: fala do total
	}
	// A lista de complementos NÃO-RESTRITIVOS é fechada e curta: são os que ligam o
	// rótulo ao seu domínio ("cláusulas de autorização", "modelos do produto") sem
	// recortar nada. Qualquer outra coisa depois do rótulo — adjetivo, oração — é
	// tratada como recorte, e o gate se cala.
	//
	// Ser fechada é o que a torna honesta: em vez de adivinhar semântica ("de dado
	// financeiro" restringe, "de autorização" não), o projeto vê a lista e sabe o que o
	// gate confronta. Errar para o lado do silêncio custa achado; errar para o outro
	// acusa quem está certo, e um gate assim é desligado.
	for _, lig := range complementosNaoRestritivos {
		if strings.EqualFold(r, lig) {
			return false
		}
	}
	return true
}

// complementosNaoRestritivos: o que segue o rótulo sem recortar o conjunto.
var complementosNaoRestritivos = []string{
	"de autorização", "de autorizacao", "do produto", "do projeto", "no total",
	"of authorization", "in total", "declarados", "no schema",
}
