package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// regra-implementada: a spec cataloga regras; o código tem de mostrar que as realizou.
//
// É o confronto que faltava, e ele é o inverso do `code-reference-valid`: aquele confere
// que os códigos CITADOS pelo código existem na spec; este confere que a spec não ficou
// falando sozinha. Sem ele, uma spec pode declarar cinco regras novas e o código não
// ganhar linha nenhuma — com todos os gates verdes, porque a spec existe, o código existe,
// e os dois se referenciam pelo header.
//
// Medido num E2E real: uma spec de handler ganhou `MRHMX-B04..B08` (5 regras, 98 linhas) e
// o handler correspondente tinha ZERO ocorrências de `metadata|migrat`. A metade
// "recorrência" da feature era código morto declarado como pronto, e nenhum dos 26 gates
// perguntou. O defeito só apareceu quando um revisor leu spec e código na mesma passada.
//
// # A régua: cada regra, contra o que a SPEC declarou sobre ela
//
// Exigir todas as regras marcadas seria falso por construção — medido contra 592 unidades,
// daria 3.121 achados, e nem as bem-feitas passariam: das 8 que marcam o código, NENHUMA
// marca 100% (76% a 93%). A razão é boa, não descuido: restrição (`X01` — "a unidade NÃO
// faz Y") é satisfeita pela AUSÊNCIA de código, e ausência não tem onde receber comentário.
//
// Mas "ao menos uma" também não serve: separa quem implementou de quem não implementou, e
// não diz nada sobre as outras quinze regras da spec. É heurística, e heurística sobre
// intenção sempre erra o caso que não está sendo olhado.
//
// A saída é não adivinhar: quem escreve a spec DECLARA, regra a regra, se ela tem código.
// O padrão é o do `@no-scenario` (CONCEPT §5.1) — a dispensa vai na linha do que ela
// dispensa, com razão escrita:
//
//	| `MTVRX-B01` | resolve a versão vigente do ciclo |
//	| `MTVRX-X01` | NÃO faz I/O @no-code: satisfeita pela ausência de import |
//
// O gate deixa de estimar e passa a confrontar uma afirmação. E o `anchors work` diz o que
// fazer quando a afirmação não se sustenta na hora de implementar: **abrir issue**, porque
// aí a spec está errada sobre a própria unidade — não é detalhe de formatação, e apagar a
// regra ou inventar um comentário para calar o gate esconde o defeito em vez de tratá-lo.
func checkRegraImplementada(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "quem cataloga a regra é a spec — é dela que o confronto parte"
	}
	unidade := codigoDaSpec(content)
	if unidade == "" {
		return Skip, "a spec não declara `code:` no header — sem identidade não há o que confrontar"
	}
	regras := regrasDeclaradas(content, unidade)
	if len(regras) == 0 {
		return Skip, "a spec não cataloga nenhuma regra com código"
	}
	alvo, achou := alvoDaSpecNoDisco(root, n.ID)
	if !achou {
		// Sem código, quem acusa é o `trinca-completa` — este gate confronta o código que
		// existe, e duplicar a cobrança produziria dois gates apontando o mesmo dedo.
		return Skip, "o código desta unidade ainda não existe — a ausência é do gate `trinca-completa`"
	}
	cod, err := os.ReadFile(filepath.Join(root, alvo))
	if err != nil {
		return Pending, "não foi possível ler `" + alvo + "`"
	}
	texto := string(cod)

	// Regra a regra: cada uma que a spec NÃO dispensou precisa aparecer no código.
	//
	// É aqui que a heurística vira confronto. Antes o gate perguntava "citou alguma?" —
	// o que separa "não implementou" de "implementou", mas não diz NADA sobre as outras
	// 15 regras da spec. Com a dispensa declarada, a pergunta passa a ser respondível uma
	// a uma: quem escreveu a spec sabe se a regra tem código, e diz.
	var faltando []string
	for _, r := range regras {
		if strings.Contains(texto, r) {
			continue
		}
		if dispensada(content, r) {
			continue
		}
		faltando = append(faltando, r)
	}
	if len(faltando) == 0 {
		return Pass, ""
	}
	// DÍVIDA DE MIGRAÇÃO: a unidade não declarou NENHUMA regra — nem marcando o código,
	// nem dispensando. Isso é a marca de quem nasceu antes da prática existir, não de
	// quem errou.
	//
	// Medido no repositório que originou o gate: 3.114 regras em 590 unidades nessa
	// situação. Acusá-las como defeito faria o gate reprovar 98% do projeto no primeiro
	// dia — e um gate assim é desligado, levando junto os que funcionam.
	//
	// Pendente, e não Pass: o veredito NOMEIA a dívida (aparece no `check` como
	// pendência, entra no relatório) sem fingir que está tudo certo. Conforme cada
	// unidade for tocada e declarar suas regras, ela sai deste ramo e passa a ser
	// confrontada de verdade.
	// A pendência de migração só vale enquanto o projeto NÃO declarou que a prática é
	// obrigatória. Declarado `rule_marking: required`, a migração acabou: a unidade que
	// não marcou nada é cobrada como qualquer outra.
	//
	// Sem esta saída a pendência nunca vence — e pendência que não vence é
	// indistinguível de "ninguém olhou". Medido: uma spec descrevia uma unidade que o
	// código não implementava, o gate viu a regra ausente, caiu aqui e devolveu
	// pendência. O defeito atravessou os 44 gates.
	if len(faltando) == len(regras) && !temAlgumaDeclaracao(content, texto, regras) &&
		exigeMarcacao(cfg) {
		sort.Strings(faltando)
		return Fail, fmt.Sprintf(
			"nenhuma das %d regra(s) da spec aparece no código, e este projeto exige a "+
				"marcação (`derived.rule_marking: required`).\n\nMarque no código o trecho "+
				"que realiza cada regra (`// %s-B01: …`) ou dispense na linha dela "+
				"(`@no-code: <razão>`). Se NENHUMA regra tem código, verifique antes se a "+
				"spec descreve mesmo esta unidade — é o sintoma de spec que fala de outra coisa",
			len(regras), unidade)
	}

	if len(faltando) == len(regras) && !temAlgumaDeclaracao(content, texto, regras) {
		return Pending, fmt.Sprintf("%d regra(s) catalogada(s) e nenhuma declarada — esta "+
			"unidade é anterior à prática de ligar regra↔código. Ao tocá-la, marque no código "+
			"a regra que cada trecho realiza (`// %s-B01: …`) ou dispense na linha dela "+
			"(`@no-code: <razão>`); aí o gate passa a confrontar de verdade",
			len(regras), unidade)
	}
	sort.Strings(faltando)
	mostra := faltando
	if len(mostra) > 6 {
		mostra = mostra[:6]
	}
	return Fail, fmt.Sprintf("%d regra(s) que a spec declara e o código `%s` não realiza: %s%s. "+
		"Uma regra catalogada sem implementador atravessa o pipeline inteiro — a spec existe, "+
		"o código existe, os dois se referenciam pelo header, e todos os gates ficam verdes "+
		"sobre trabalho que não foi feito.\n\n"+
		"Duas saídas, e as duas são honestas: marque no código o trecho que realiza a regra "+
		"(`// %s-B01: …`), ou declare na linha dela `@no-code: <razão>` — para o que é "+
		"satisfeito pela AUSÊNCIA de código (restrição, limite de escopo).\n\n"+
		"Se ao implementar a razão da dispensa não se sustentar, isso não é detalhe de "+
		"formatação: é a spec estar errada sobre a própria unidade. Abra issue "+
		"(`anchors judge %s --gate review --verdict fail --reason \"…\"`) em vez de "+
		"apagar a regra ou inventar um comentário para calar o gate",
		len(faltando), alvo, strings.Join(mostra, ", "),
		sufixoResto(len(faltando)-len(mostra)), unidade, n.ID)
}

// noCodeRE: a dispensa DECLARADA de que uma regra apareça no código. Vai na linha da
// regra, e exige razão escrita — marcador nu não dispensa nada.
//
//	| `MTVRX-X01` | NÃO faz I/O @no-code: satisfeita pela ausência de import |
//
// O padrão é o do `@no-scenario` (CONCEPT §5.1): a dispensa fica À VISTA, na linha do que
// ela dispensa, versionada pelo git e legível para quem abrir a spec depois. Escondê-la
// num Skip do gate seria a mesma decisão sem a mesma prestação de contas.
// A razão tem de ser TEXTO, não o fechamento da célula: `@no-code: |` num markdown de
// tabela satisfaria um `\S+` ingênuo e dispensaria a regra sem dizer nada. Marcador nu não
// dispensa — é o que separa prestar contas de calar o gate.
var noCodeRE = regexp.MustCompile(`@no-code[^\S\n]*:[^\S\n]*[^\s|]\S*`)

// dispensada diz se a LINHA que declara a regra carrega a dispensa.
//
// A dispensa é da linha, não do arquivo: um `@no-code` solto no topo da spec dispensaria
// tudo de uma vez, que é o oposto de prestar contas.
func dispensada(content, regra string) bool {
	for _, linha := range strings.Split(content, "\n") {
		if strings.Contains(linha, regra) && noCodeRE.MatchString(linha) {
			return true
		}
	}
	return false
}

// codigoRE lê a identidade declarada no header da spec.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func codigoRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

func codigoDaSpec(content string) string {
	if m := codigoRE().FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// regrasDeclaradas devolve as regras que a spec CATALOGA — as que abrem linha, em tabela,
// lista ou título. Uma citação no meio de um parágrafo é referência, não declaração, e
// contá-la faria o gate cobrar do código regras que pertencem a outra unidade.
func regrasDeclaradas(content, unidade string) []string {
	re := regexp.MustCompile(`(?m)^\s*(?:\|\s*|-\s*|#{2,4}\s+)` + "`?" + `(` + unidade + `-[A-Z]\d{2})` + "`?")
	visto := map[string]bool{}
	var out []string
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if !visto[m[1]] {
			visto[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

// alvoDaSpecNoDisco acha o arquivo de código que a spec descreve, pelo nome.
func alvoDaSpecNoDisco(root, specID string) (string, bool) {
	base := strings.TrimSuffix(specID, ".spec.md")
	if base == specID {
		return "", false
	}
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
		if _, err := os.Stat(filepath.Join(root, base+ext)); err == nil {
			return base + ext, true
		}
	}
	return "", false
}

func sufixoResto(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf(" e mais %d", n)
}

// temAlgumaDeclaracao diz se a unidade JÁ ENTROU na prática — se alguma regra foi marcada
// no código ou dispensada na spec. É o que separa a dívida de migração (ninguém declarou
// nada, porque a prática não existia) do defeito (declarou umas, esqueceu outras).
func temAlgumaDeclaracao(spec, codigo string, regras []string) bool {
	for _, r := range regras {
		if strings.Contains(codigo, r) || dispensada(spec, r) {
			return true
		}
	}
	return false
}

// exigeMarcacao diz se o projeto declarou que a marcação regra↔código é obrigatória.
// Vazio = migração em curso (a pendência vale); "required" = a migração acabou.
func exigeMarcacao(cfg *config.Config) bool {
	return cfg != nil && cfg.Derived != nil &&
		strings.EqualFold(strings.TrimSpace(cfg.Derived.RuleMarking), "required")
}
