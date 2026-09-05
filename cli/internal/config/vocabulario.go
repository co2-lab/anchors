package config

import "strings"

// --- o VOCABULÁRIO em inglês, e os nomes antigos que ainda funcionam ---
//
// Os nomes de gate são IDENTIFICADORES: eles vão para o `anchors.yaml` de cada projeto,
// para tutoriais, para respostas de fórum. Por isso são fixos em inglês e NÃO se
// traduzem — um `anchors.yaml` escrito por um time brasileiro tem de funcionar num time
// espanhol sem tradução nenhuma.
//
// Este arquivo existe porque eles NASCERAM em português, e renomear sem alias quebraria
// todo projeto que já os declarou. O nome antigo continua sendo aceito: o `config.Load`
// o converte, e o `anchors doctor` avisa que ele está obsoleto — o projeto migra quando
// quiser, não quando o Anchors decidir.

// nomesAntigos mapeia o nome em português para o nome canônico em inglês.
var nomesAntigos = map[string]string{
	"codigo-catalogado":           "code-cataloged",
	"dependencia-vulneravel":      "dependency-vulnerable",
	"fase-existe":                 "phase-exists",
	"fase-ordenada":               "phase-ordered",
	"feature-nao-vazia":           "feature-not-empty",
	"mock-carimbado":              "mock-stamped",
	"mock-detect-cobre-o-dialeto": "mock-detect-covers-dialect",
	"mock-tipado":                 "mock-typed",
	"no-test-prova-real":          "no-test-proof-real",
	"parent-valido":               "parent-valid",
	"plano-alterado-justificado":  "plan-change-justified",
	"plano-revisado":              "plan-revised",
	"prova-cruza-fronteira":       "proof-crosses-boundary",
	"regra-cumprida":              "rule-fulfilled",
	"sbom-gerado":                 "sbom-generated",
	"secret-nao-vazado":           "no-secret-leaked",
	"sem-duplicacao":              "no-duplication",
	"spec-completa":               "spec-complete",
	"spec-tem-codigo":             "spec-has-code",
	"teste-rastreavel":            "test-traceable",
	"trinca-completa":             "triad-complete",
}

// checksAntigos mapeia o `check:` em português para o canônico.
//
// Separado dos gates porque um projeto pode ter renomeado o GATE e mantido o `check:`
// (o nome do gate é livre; o do check é o verificador interno).
var checksAntigos = map[string]string{
	"codigo-catalogado":          "code-cataloged",
	"fase-existe":                "phase-exists",
	"fase-ordenada":              "phase-ordered",
	"mock-carimbado":             "mock-stamped",
	"mock-tipado":                "mock-typed",
	"parent-valido":              "parent-valid",
	"plano-alterado-justificado": "plan-change-justified",
	"plano-revisado":             "plan-revised",
	"prova-cruza-fronteira":      "proof-crosses-boundary",
	"teste-rastreavel":           "test-traceable",
	"trinca-completa":            "triad-complete",
}

// CanonicalizaNome devolve o nome canônico e se houve conversão.
//
// O segundo retorno é o que permite ao `doctor` avisar sem que o `Load` precise imprimir:
// carregar a configuração não é lugar de escrever na tela.
func CanonicalizaNome(n string) (string, bool) {
	if c, ok := nomesAntigos[n]; ok {
		return c, true
	}
	return n, false
}

// CanonicalizaCheck faz o mesmo para o campo `check:`.
func CanonicalizaCheck(c string) (string, bool) {
	if n, ok := checksAntigos[c]; ok {
		return n, true
	}
	return c, false
}

// NomesAntigos devolve o de-para inteiro, para o `doctor` listar o que migrar.
func NomesAntigos() map[string]string {
	out := make(map[string]string, len(nomesAntigos))
	for k, v := range nomesAntigos {
		out[k] = v
	}
	return out
}

// canonicalizaVocabulario converte os nomes antigos para os canônicos, in-place.
//
// Roda na CARGA para que todo o resto do Anchors veja um vocabulário só. A alternativa —
// cada consumidor consultar o de-para — espalharia a conversão por dezenas de lugares, e
// o primeiro que esquecesse produziria um gate que "não existe" num projeto que o
// declarou.
func (c *Config) canonicalizaVocabulario() {
	if c == nil {
		return
	}
	for i := range c.Gates {
		if n, mudou := CanonicalizaNome(c.Gates[i].Name); mudou {
			c.Gates[i].Name = n
			// O ID acompanha o nome quando eram iguais: o ID é a identidade estável do
			// gate, e um projeto que não o declarou explicitamente tinha os dois iguais.
			// Deixar o ID no nome antigo faria o `check` reportar um e o mapa gravar
			// outro.
			if c.Gates[i].ID == "" {
				c.Gates[i].ID = n
			}
		}
		if id, mudou := CanonicalizaNome(c.Gates[i].ID); mudou {
			c.Gates[i].ID = id
		}
		if ch, mudou := CanonicalizaCheck(c.Gates[i].Check); mudou {
			c.Gates[i].Check = ch
		}
	}
}

// VocabularioObsoleto devolve os nomes antigos que ESTE projeto ainda usa.
//
// É o que o `doctor` lista. Lê o arquivo cru em vez do Config já carregado, porque a
// carga já converteu — perguntar ao Config depois disso devolveria sempre vazio.
func VocabularioObsoleto(conteudoYAML string) map[string]string {
	achados := map[string]string{}
	for antigo, novo := range nomesAntigos {
		if strings.Contains(conteudoYAML, antigo) {
			achados[antigo] = novo
		}
	}
	for antigo, novo := range checksAntigos {
		if strings.Contains(conteudoYAML, antigo) {
			achados[antigo] = novo
		}
	}
	return achados
}

// --- a ponte para o teste do de-para ---
//
// O teste que confronta "todo destino aponta para um gate que EXISTE" precisa da lista de
// gates default, que vive no `initx`. E `config` não pode importar `initx` — seria ciclo,
// já que o `initx` monta gates a partir de tipos daqui.
//
// A injeção resolve: o `initx` registra a lista no seu `init()`, e o teste a consulta. Sem
// isso o teste não teria contra o que confrontar, e um de-para apontando para o vazio
// passaria — o projeto carregaria, o gate viraria um nome que nenhum verificador conhece,
// e o `check` reportaria "gate sem nada para medir" sem dizer que a causa foi a conversão.

var defaultGateNames func() []string

// RegistraNomesDeGate liga a lista de gates default ao pacote config.
func RegistraNomesDeGate(f func() []string) { defaultGateNames = f }

// DefaultGatesParaTeste devolve os nomes registrados, ou nil se ninguém registrou.
func DefaultGatesParaTeste() []string {
	if defaultGateNames == nil {
		return nil
	}
	return defaultGateNames()
}
