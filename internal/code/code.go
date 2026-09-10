// Package code gera e valida o código de identidade de uma unidade (TRACEABILITY
// §3). Implementa a doutrina universal do SPEC_GUIDE do app de referência:
//
//	Camada 2 — compressão do nome em N chars (iniciais / consoantes 2+2 / consoantes
//	           + vogais / X-padding). É a régua agnóstica.
//	Camada 3 — resolução de colisão determinística contra os códigos já tomados.
//
// A Camada 1 (prefixos de módulo tipo AU/IR/FM) é DIALETO de projeto — não entra
// aqui; um projeto que a queira passa um prefixo explícito. A nota do guide manda:
// "o importante é a unicidade no namespace global, não a estética".
package code

import (
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// Slots é o comprimento do código que o `anchors code` GERA. Default 5 (o canônico do
// framework), reconfigurado por `SetSlots` a partir do `code_lengths` do projeto.
//
// Era `const 5`, e a constante estava errada de duas formas ao mesmo tempo: sugeria código
// de 5 num projeto que declara 4 (o `anchors code Login` respondia `LGNOI` para um mundo de
// `LOGI`), e tornava impossível checar identidade contra a régua — comparar um código
// declarado de 4 com um gerado de 5 acusa TODOS como errados.
//
// Medido no app de referência: dos 8 primeiros nomes de tela, zero batiam com 5 slots; com 4, três
// batiam exatamente. Os cinco que ainda divergem são escolha humana legítima (`ARSC` para
// AlertsScreen é mais legível que `LRTS`), e é por isso que o gerado é SUGESTÃO — a
// doutrina diz "o importante é a unicidade no namespace global, não a estética".
//
// Continua distinto de `config.CodeLengths`: aquele é o que o engine RECONHECE ao ler (pode
// aceitar 4 e 5 durante uma migração), este é o comprimento ÚNICO que a geração emite.
var Slots = 5

// SetSlots ajusta o comprimento gerado ao que o projeto declara.
//
// Recebe a lista de `code_lengths` e usa o MENOR: durante uma migração o projeto aceita 4 e
// 5 ao mesmo tempo, e gerar no menor mantém o código novo compatível com o vocabulário
// antigo — o contrário (gerar 5 num projeto que ainda lê 4) criaria código que metade das
// ferramentas do próprio projeto não reconhece.
func SetSlots(lengths []int) {
	menor := 0
	for _, l := range lengths {
		if l >= 2 && (menor == 0 || l < menor) {
			menor = l
		}
	}
	if menor > 0 {
		Slots = menor
	}
}

// sufixos genéricos que não agregam identidade — removidos antes de comprimir.
var genericSuffixes = []string{"Screen", "Component", "Modal", "Sheet", "Layout", "Overlay", "View", "Page"}

func isVowel(c byte) bool     { return strings.IndexByte("AEIOUaeiou", c) >= 0 }
func isLetter(c byte) bool    { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isConsonant(c byte) bool { return isLetter(c) && !isVowel(c) }
func isDigit(c byte) bool     { return c >= '0' && c <= '9' }

// StripGeneric remove um sufixo genérico do nome (LoginScreen → Login), preservando
// o nome se ele FOR só o genérico.
func StripGeneric(name string) string {
	for _, s := range genericSuffixes {
		if len(name) > len(s) && strings.HasSuffix(name, s) {
			return name[:len(name)-len(s)]
		}
	}
	return name
}

// tokenize quebra um nome em palavras: separadores, camelCase e ACRÔNIMOSeguido.
func tokenize(name string) []string {
	var b strings.Builder
	runes := []byte(name)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if c == '-' || c == '_' || c == ' ' {
			b.WriteByte(' ')
			continue
		}
		// fronteira camelCase: minúscula/dígito seguida de Maiúscula
		if i > 0 && c >= 'A' && c <= 'Z' {
			prev := runes[i-1]
			if (prev >= 'a' && prev <= 'z') || isDigit(prev) {
				b.WriteByte(' ')
			} else if prev >= 'A' && prev <= 'Z' && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				// ACR + Onym → ACR | Onym
				b.WriteByte(' ')
			}
		}
		b.WriteByte(c)
	}
	return strings.Fields(b.String())
}

// Generate produz o código canônico de um nome, SEM resolução de colisão (Camada 2).
func Generate(name string) string { return GenerateWithPrefix(name, "") }

// GenerateWithPrefix aplica a Camada 1 quando `prefix` não é vazio: o código começa
// pelo prefixo do MÓDULO (2 chars, vindo da Estrutura) + chars distintivos do nome
// preenchendo o resto dos slots. Sem prefixo, é a Camada 2 pura (compressão do nome).
// Espelha o SPEC_GUIDE: ≥4 palavras → iniciais; 2-3 → consoantes 2+2; 1 → consoantes
// e depois vogais. Sempre padronizado para Slots com X-padding.
func GenerateWithPrefix(name, prefix string) string {
	if prefix != "" {
		p := strings.ToUpper(prefix)
		if len(p) > Slots {
			p = p[:Slots]
		}
		// distintivo: comprime o nome nos slots restantes
		rest := Slots - len(p)
		distinct := lettersFrom(StripGeneric(name), rest, true)
		return pad(p + string(distinct))
	}
	base := StripGeneric(name)
	words := tokenize(base)
	if len(words) == 0 {
		words = tokenize(name) // o nome era só genérico — usa-o inteiro
		if len(words) == 0 {
			return pad("")
		}
	}

	var picked []byte
	switch {
	case len(words) >= Slots:
		// iniciais das primeiras Slots palavras
		for _, w := range words {
			if len(picked) >= Slots {
				break
			}
			picked = append(picked, upper(w[0]))
		}
	case len(words) >= 2:
		// distribui os slots entre as palavras: consoantes de cada, 2+2 típico
		per := Slots / len(words)
		if per < 1 {
			per = 1
		}
		for _, w := range words {
			picked = append(picked, lettersFrom(w, per, true)...)
		}
		// se faltou (poucas consoantes), completa com a estratégia de 1-palavra
		if len(picked) < Slots {
			picked = append(picked, lettersFrom(strings.Join(words, ""), Slots-len(picked), false)...)
		}
	default:
		// 1 palavra: consoantes primeiro, depois vogais
		picked = lettersFrom(words[0], Slots, false)
	}
	return pad(string(picked))
}

// lettersFrom colhe até n letras de w. Se consonantsOnly, mira consoantes primeiro
// (mantendo a 1ª letra mesmo se vogal, como o guide: "aLRTS"→ mantém a inicial);
// senão colhe consoantes e depois vogais na ordem do nome.
func lettersFrom(w string, n int, preferInitial bool) []byte {
	var out []byte
	if n <= 0 || w == "" {
		return out
	}
	// sempre tenta a 1ª letra como âncora (inicial da palavra)
	start := 0
	if preferInitial {
		out = append(out, upper(w[0]))
		start = 1
	}
	// consoantes
	for i := start; i < len(w) && len(out) < n; i++ {
		if isConsonant(w[i]) {
			out = append(out, upper(w[i]))
		}
	}
	// vogais/dígitos, se ainda faltar
	for i := 0; i < len(w) && len(out) < n; i++ {
		if isVowel(w[i]) || isDigit(w[i]) {
			c := upper(w[i])
			if !contains(out, c) || isDigit(w[i]) {
				out = append(out, c)
			}
		}
	}
	if len(out) > n {
		out = out[:n]
	}
	return out
}

// GenerateUnique resolve colisão (Camada 3): parte do código canônico e, se já
// tomado, varia deterministicamente até achar um livre. `taken` são os códigos já em
// uso no namespace. Determinístico dado (name, taken).
func GenerateUnique(name string, taken map[string]bool) string {
	return GenerateUniqueWithPrefix(name, "", taken)
}

// GenerateUniqueWithPrefix é GenerateUnique honrando o prefixo de módulo (Camada 1).
func GenerateUniqueWithPrefix(name, prefix string, taken map[string]bool) string {
	code := GenerateWithPrefix(name, prefix)
	if !taken[code] {
		return code
	}
	// colisão: troca a última letra por A..Z, depois a penúltima — busca
	// determinística por um código livre (Camada 3 do SPEC_GUIDE, forma agnóstica).
	// Preserva o prefixo (as primeiras len(prefix) posições não mudam).
	fixed := len(prefix)
	if fixed > Slots {
		fixed = Slots
	}
	runes := []byte(code)
	for pos := len(runes) - 1; pos >= fixed; pos-- { // não mexe nas posições do prefixo
		orig := runes[pos]
		for c := byte('A'); c <= 'Z'; c++ {
			if c == orig {
				continue
			}
			runes[pos] = c
			if cand := string(runes); !taken[cand] {
				return cand
			}
		}
		runes[pos] = orig
	}
	return code // desistiu (namespace saturado) — devolve o canônico
}

// ModulePrefix deriva um prefixo de 2 chars de um nome de módulo (auth→AU,
// family→FM, inventory→IV...) — o mesmo espírito do gerador, comprimido a 2 slots.
// Usado pelo init para DEDUZIR os prefixos de módulo a partir da Estrutura, em vez
// de uma tabela hardcoded.
func ModulePrefix(module string) string {
	letters := lettersFrom(module, 2, true)
	return pad2(string(letters))
}

func pad2(s string) string {
	s = strings.ToUpper(s)
	for len(s) < 2 {
		s += "X"
	}
	return s[:2]
}

func pad(s string) string {
	s = strings.ToUpper(s)
	for len(s) < Slots {
		s += "X"
	}
	return s[:Slots]
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}

func contains(bs []byte, c byte) bool {
	for _, b := range bs {
		if b == c {
			return true
		}
	}
	return false
}

// Liga o gerador ao `code_lengths` do projeto na carga da config.
//
// `init` no pacote do GERADOR, e não no `config`, pela mesma razão que o
// SetCanonicalGateResolver mora no `initx`: quem tem a dependência é quem se registra, e
// assim `config` continua sem importar ninguém — é a base da pilha.
func init() { config.SetSlotsHook(SetSlots) }

// Pad completa um código existente até `Slots` com X — a MESMA regra do gerador.
//
// Exportado para o `anchors code list --check`: quando um código está com o comprimento
// errado, o conserto é completá-lo, não gerar outro. `MTVR` → `MTVRX` preserva a escolha de
// quem nomeou; `MTVR` → `LRTS` (regenerar do nome) a descartaria — e testado no app de referência, o
// gerado divergia do declarado em 5 de 5 casos sem nenhum deles estar errado.
func Pad(s string) string { return pad(s) }
