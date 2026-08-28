package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// scenario-asserts: o passo de RESULTADO de um cenário precisa afirmar um resultado.
//
// A falha que este gate pega é a **tautologia** — "dado X, então X":
//
//	@XXXXX-B05 @nivel-unit
//	Cenário: cada chave versiona sozinha
//	  Dado ...
//	  Quando ...
//	  Então o efeito XXXXX-B05 se verifica    ← não afirma NADA
//
// O passo final repete o código da regra em vez de dizer o que deveria acontecer. Dois
// cenários de regras diferentes ficam indistinguíveis, e a definição do que "se verifica"
// migra inteiramente para o teste — que é justamente o artefato que a feature deveria
// governar. A inversão é silenciosa: todos os gates passam, porque o código de cenário
// está lá e `feature-test-match` confronta código e descrição, não força de asserção.
//
// Por que importa além do estilo: um revisor externo identificou a tautologia como a razão
// ESTRUTURAL pela qual uma regra ficou sem caso discriminante. Ninguém, entre a spec e o
// teste, foi obrigado a escrever o resultado esperado — então o teste que "prova" a regra
// foi escrito sem saber qual era o desfecho, e uma mutação que apagava a regra manteve a
// suíte verde. Medido no mesmo projeto: 359 cenários com esta forma.
//
// O gate NÃO julga a qualidade da prosa (isso é julgamento, e há gate de IA para isso).
// Ele pega uma forma específica e mecânica: o passo de resultado cujo conteúdo é
// essencialmente o CÓDIGO DA REGRA e nada mais.
func checkScenarioAsserts(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, "o passo de resultado é da feature — é ela que descreve o cenário"
	}

	d := cfg.DialectFor()
	_, kw := d.GherkinFor()
	// O passo de resultado é o `Então`/`Then` do idioma do projeto. Reconhece também os
	// demais idiomas: um repositório pode ter features herdadas noutra língua, e ficar
	// cego aqui significaria reportar verde sobre o que não se enxerga.
	thens := map[string]bool{kw.Then: true}
	for _, alt := range config.GherkinThenAlternatives() {
		thens[alt] = true
	}

	var vazios []string
	for _, linha := range strings.Split(content, "\n") {
		passo := strings.TrimSpace(linha)
		if passo == "" || strings.HasPrefix(passo, "#") {
			continue
		}
		primeira, resto := primeiraPalavra(passo)
		if !thens[primeira] {
			continue
		}
		if code, tautologico := ehTautologia(resto); tautologico {
			vazios = append(vazios, code)
		}
	}
	if len(vazios) == 0 {
		return Pass, ""
	}

	sort.Strings(vazios)
	vazios = dedup(vazios)
	return Fail, fmt.Sprintf("%d cenário(s) com passo de resultado que não afirma resultado "+
		"(%s): o `%s` repete o código da regra em vez de dizer o que deveria acontecer. "+
		"Escreva o resultado OBSERVÁVEL, com o valor esperado — sem isso, a definição do que "+
		"\"se verifica\" migra para o teste, e o teste passa a ser escrito sem saber qual era "+
		"o desfecho",
		len(vazios), strings.Join(vazios, ", "), kw.Then)
}

// ehTautologia: o resto do passo é só o código da regra, cercado de palavras de ligação?
// A régua é mecânica de propósito — julgar prosa é trabalho de outro gate.
func ehTautologia(resto string) (string, bool) {
	m := tautoCodeRE().FindStringSubmatch(resto)
	if m == nil {
		return "", false
	}
	// Remove o código e as palavras de ligação; o que sobra é o conteúdo real do passo.
	semCode := tautoCodeRE().ReplaceAllString(resto, " ")
	conteudo := ligacaoRE.ReplaceAllString(strings.ToLower(semCode), " ")
	conteudo = strings.TrimSpace(regexp.MustCompile(`[^\p{L}\p{N}]+`).ReplaceAllString(conteudo, " "))
	// Até duas palavras residuais ainda é tautologia ("o efeito X se verifica" →
	// "efeito verifica"). Acima disso, o autor escreveu algo de próprio.
	if len(strings.Fields(conteudo)) <= 2 {
		return m[1], true
	}
	return "", false
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func tautoCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`\b([A-Z0-9]` + config.CodeLengthPattern() + `-[A-Z]\d{2})\b`)
}

// ligacaoRE são as palavras que só ligam — sem elas, o passo não perde afirmação. Cobre os
// idiomas do Gherkin que o Anchors conhece.
var ligacaoRE = regexp.MustCompile(`\b(o|a|os|as|um|uma|de|do|da|se|e|que|the|a|an|of|is|are|el|la|los|las|del|se|le|les|du|de|der|die|das|` +
	`efeito|efeitos|regra|regras|comportamento|cenário|cenario|requisito|` +
	`verifica|verificado|verificada|aplica|aplicado|aplicada|vale|válido|valido|ocorre|acontece|` +
	`effect|rule|behavior|behaviour|verified|applies|holds|is met|met|satisfied|` +
	`efecto|regla|verifica|cumple|` +
	`effet|règle|vérifie|` +
	`effekt|regel|gilt|erfüllt)\b`)

func primeiraPalavra(s string) (string, string) {
	f := strings.Fields(s)
	if len(f) == 0 {
		return "", ""
	}
	return f[0], strings.TrimSpace(strings.TrimPrefix(s, f[0]))
}

func dedup(xs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
