// Helpers de testID: reconhecer o handle nas duas pontas que o ESCREVEM — o código
// (JSX/atributo) e a spec (tabela de inventário).
//
// Este arquivo já foi `testid_declared.go`, sede do gate homônimo. O gate foi
// aposentado por `testid-coerente`, que confronta a trinca inteira em vez de duas
// pontas; os reconhecedores continuam aqui porque são dele que todos partem.
package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// reExposto/reTemplate são construídos a partir do ATRIBUTO que o projeto declara
// (`derived.test_handle`) — `testID`, `data-testid`, `contentDescription`. O nome é do
// ecossistema, não do framework; cravá-lo aqui tornaria o gate específico de React
// Native e inútil (ou pior, silenciosamente verde) em qualquer outro projeto.
//
// Duas formas por atributo:
//   - LITERAL   `handle="foo"` / `handle={"foo"}`
//   - TEMPLATE  handle={`foo-${id}`} — captura só a parte ESTÁTICA. O sufixo varia em
//     runtime; exigir que a spec declare `foo-42` seria exigir que ela declare o DADO.
//     Entra no inventário como `foo-*`, que é como a spec o declara.
//
// O sufixo `[a-zA-Z]*` antes do atributo cobre as props DERIVADAS que carregam um
// handle sem se chamarem exatamente como ele (`backTestID`, `optionTestIDPrefix`,
// `confirmTestID`). O valor delas termina num handle real dentro do componente
// filho — ignorá-las fazia o gate declarar inexistente um id que a spec documenta
// corretamente e que os flows usam.
// O `(?:[a-zA-Z]*(?:Prefix)?)?` NÃO entra aqui: a prop terminada em `Prefix` é a
// CABEÇA de um id montado no filho, não um id — ela é tratada à parte, virando
// `cabeça-*`. Deixá-la casar o regex literal registraria `otp-input` (que nunca
// aparece na tela) como se fosse um handle exposto.
var sufixoPrefix = regexp.MustCompile(`(?i)prefix$`)

// O `(?i)` é essencial e não é frouxidão: a prop derivada CAPITALIZA o atributo ao
// compô-lo (`testID` → `backTestID`, `confirmTestID`). Sem ignorar a caixa, o regex
// casaria só a forma exata e perderia todas as derivadas — que é como o gate acabou
// declarando inexistentes ids que a spec documenta e os flows usam.
func regexDeHandle(attr string) (literal, template *regexp.Regexp) {
	a := `(?i)([a-z]*` + regexp.QuoteMeta(attr) + `[a-z]*)`
	return regexp.MustCompile(a + `\s*[=:]\s*\{?["'](:?[a-zA-Z][a-zA-Z0-9._-]*)["']`),
		regexp.MustCompile(a + "\\s*[=:]\\s*\\{?`(:?[a-zA-Z][a-zA-Z0-9._-]*)\\$?\\{?")
}

// testIDsExpostos devolve os handles que a unidade oferece ao mundo. O template
// entra na forma `prefixo-*`, que é como a spec o declara.
func testIDsExpostos(src string, attr string) []string {
	testIDExpostoRE, testIDTemplateRE := regexDeHandle(attr)
	visto := map[string]bool{}
	var out []string
	add := func(s string) {
		if !visto[s] {
			visto[s] = true
			out = append(out, s)
		}
	}
	// m[1] = nome da prop, m[2] = valor. A prop terminada em `Prefix` é a CABEÇA de um
	// id que o FILHO completa (`testIDPrefix=":otp-input"` → `otp-input-0`…`-5`):
	// entra como `cabeça-*`, não como id. Registrá-la literalmente cobraria da spec um
	// `otp-input` que nunca aparece na tela, e acusaria de inexistentes os
	// `otp-input-N` que a spec documenta e que 6 flows usam.
	for _, m := range testIDExpostoRE.FindAllStringSubmatch(src, -1) {
		if sufixoPrefix.MatchString(m[1]) {
			add(strings.TrimSuffix(m[2], "-") + "-*")
			continue
		}
		add(m[2])
	}
	for _, m := range testIDTemplateRE.FindAllStringSubmatch(src, -1) {
		add(strings.TrimSuffix(m[2], "-") + "-*")
	}
	// LITERAL DENTRO DE EXPRESSÃO: `testID={i === 0 ? ':goal-new-icon-option-0' :
	// undefined}`. Os regexes acima exigem o valor colado ao `=`, então param no `{`
	// e perdem o id — que existe, é consultado pelo teste e aparece no flow. A
	// expressão é delimitada por `{...}` sem aninhamento (é o formato destes
	// ternários), e dentro dela todo literal kebab-case é candidato a handle.
	for _, m := range regexp.MustCompile(
		`(?i)[a-z]*`+regexp.QuoteMeta(attr)+`[a-z]*\s*=\s*\{([^{}]*\?[^{}]*)\}`,
	).FindAllStringSubmatch(src, -1) {
		// Só o que vem DEPOIS do primeiro `?`. A CONDIÇÃO também tem literais, e eles
		// são estado de domínio (`target === 'email'`, `c2.k === 'push'`) — colhê-los
		// registraria `email` e `push` como handles expostos e o gate cobraria da spec
		// a declaração de um id que não existe.
		i := strings.Index(m[1], "?")
		if i < 0 {
			continue
		}
		for _, lit := range regexp.MustCompile(`["'](:?[a-z][a-z0-9-]*)["']`).
			FindAllStringSubmatch(m[1][i+1:], -1) {
			add(lit[1])
		}
	}
	return out
}

// handleDeTeste — o atributo com que ESTE projeto marca elementos para alcance
// externo. Sem declaração devolve vazio, e os gates de inventário pulam: o Anchors
// não adivinha o ecossistema, do mesmo modo que não adivinha idioma nem vendor.
func handleDeTeste(cfg *config.Config) string {
	if cfg == nil || cfg.Derived == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Derived.TestHandle)
}

// testIDDeclaradoRE captura o id citado em crase na spec. Aceita as três formas
// catalogadas que o projeto já usa (SPEC_TYPES §5): linha de tabela, bullet e
// cabeçalho — o que importa é o id estar em crase, não a moldura em volta.
var testIDDeclaradoRE = regexp.MustCompile("`(:?[a-zA-Z][a-zA-Z0-9._-]*(?:-\\*)?)`")

// secaoSuperficieRE delimita a seção de inventário. Sem ela, qualquer crase na prosa
// da spec (nome de classe CSS, de campo, de arquivo) contaria como declaração — e o
// gate passaria a aprovar por acidente, que é pior que reprovar por engano.
//
// O título aceita QUALIFICADOR depois do nome (`## Test IDs (Maestro)`): 42 das 159
// specs do app de referência nomeiam ali a superfície que consome os ids. Exigir fim-de-linha logo
// após o título fazia o gate não enxergar a seção dessas specs e acusá-las de não
// declarar nada — reprovando justamente as que documentam melhor.
var secaoSuperficieRE = regexp.MustCompile(`(?im)^#{1,6}\s*(?:superf[íi]cie de teste|test\s*ids?)\b.*$`)

// testIDsDeclarados lê o inventário DENTRO da seção de Superfície de Teste.
func testIDsDeclarados(spec, attr string) []string {
	loc := secaoSuperficieRE.FindStringIndex(spec)
	if loc == nil {
		return nil
	}
	corpo := spec[loc[1]:]
	// A seção termina no próximo cabeçalho de mesmo nível ou acima.
	if fim := regexp.MustCompile(`(?m)^#{1,6}\s`).FindStringIndex(corpo); fim != nil {
		corpo = corpo[:fim[0]]
	}
	visto := map[string]bool{}
	var out []string
	for _, linha := range strings.Split(corpo, "\n") {
		// Numa linha de TABELA, só a PRIMEIRA célula é o id — as demais descrevem o
		// elemento e citam os cenários que o usam, e essas colunas também levam crase.
		// Ler a linha inteira fazia o gate colher `TouchableOpacity` (coluna Elemento)
		// e `ATLNX-VR` (coluna "Usado em") como se fossem testIDs declarados, e depois
		// acusá-los de órfãos — inventando dívida a partir da própria documentação.
		alvo := linha
		if celulas := celulasDaLinha(linha); celulas != nil {
			alvo = celulas[0]
		}
		for _, m := range testIDDeclaradoRE.FindAllStringSubmatch(alvo, -1) {
			// O PRÓPRIO nome do atributo não é um id. O átomo genérico (Button,
			// Avatar, ActionLink) recebe o handle de fora e a spec o documenta como
			// `testID` (prop) — declarar isso é dizer "aceito um handle", não "exponho
			// este". Cobrá-lo faria o gate exigir do átomo um id fixo, que é o oposto
			// do que um componente reusável deve ter.
			if strings.EqualFold(strings.TrimPrefix(m[1], ":"), attr) {
				continue
			}
			if !visto[m[1]] {
				visto[m[1]] = true
				out = append(out, m[1])
			}
		}
	}
	return out
}

// celulasDaLinha devolve as células de uma linha de tabela Markdown, ou nil se a
// linha não for tabela. Separador (`| --- |`) não tem conteúdo e devolve nil.
func celulasDaLinha(linha string) []string {
	t := strings.TrimSpace(linha)
	if !strings.HasPrefix(t, "|") {
		return nil
	}
	if strings.Trim(t, "| -:") == "" {
		return nil // linha separadora
	}
	partes := strings.Split(strings.Trim(t, "|"), "|")
	for i := range partes {
		partes[i] = strings.TrimSpace(partes[i])
	}
	return partes
}

// diferenca devolve o que está em `a` e não em `b`, comparando sem o `:` de marcação
// — a spec pode citar `bdgt-screen` ou `:bdgt-screen` e as duas dizem a mesma coisa.
// O curinga casa dos DOIS lados: exposto `otp-input-*` cobre o declarado
// `otp-input-0`, e declarado `abcd-item-*` cobre o exposto `abcd-item-3`. Sem isso o
// gate acusaria em ambos os sentidos um id que o outro lado descreve corretamente —
// só que na forma genérica em vez da concreta, ou vice-versa.
func cobre(padrao, id string) bool {
	p := strings.TrimPrefix(padrao, ":")
	i := strings.TrimPrefix(id, ":")
	if p == i {
		return true
	}
	if strings.HasSuffix(p, "-*") {
		return strings.HasPrefix(i, strings.TrimSuffix(p, "*"))
	}
	return false
}

func diferenca(a, b []string) []string {
	var out []string
	for _, s := range a {
		coberto := false
		for _, t := range b {
			// Nos dois sentidos: o curinga pode estar em qualquer lado, porque a forma
			// genérica tanto é exposta pelo código (template) quanto declarada na spec.
			if cobre(t, s) || cobre(s, t) {
				coberto = true
				break
			}
		}
		if !coberto {
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// listar formata a lista de ids, truncando para a mensagem não virar despejo.
func listar(ids []string) string {
	const max = 8
	if len(ids) <= max {
		return "`" + strings.Join(ids, "`, `") + "`"
	}
	return "`" + strings.Join(ids[:max], "`, `") + "` (+" +
		fmt.Sprint(len(ids)-max) + ")"
}
