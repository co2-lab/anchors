package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// NENHUM GATE PODE DEPENDER DA REDAÇÃO DE UM TEXTO EM PORTUGUÊS.
//
// O projeto vai ser traduzido, e a tradução é substituição mecânica — o que sobrevive a
// ela é o VOCABULÁRIO ESTÁVEL (`@no-test`, `TODO`, `FNDTN-F01`, `[decisao-em-aberto]`);
// o que não sobrevive é prosa.
//
// O custo de errar isto é silencioso: o gate não falha, ele para de disparar. Já mordeu
// três vezes — o `plano-revisado` casava "revisado por", o `open-questions` lia o próprio
// laudo ("que a spec ainda NÃO tomou") para decidir se BARRAVA, e o `fase-ordenada`
// procurava a palavra "fase" (num plano em inglês, nunca dispararia).
func TestGateNaoCasaProsaEmPortugues(t *testing.T) {
	// Palavra com acento, ou sequência de palavras comuns: é prosa, não vocabulário.
	prosa := regexp.MustCompile(`"[^"]*(?:[àáâãéêíóôõúç]|\b(?:que|não|para|com|uma|dos|das|pela|ainda|aqui|descrever)\b)[^"]*"`)

	arquivos, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range arquivos {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for i, linha := range strings.Split(string(b), "\n") {
			// Só as linhas que CONFRONTAM texto. Uma mensagem de laudo pode (e deve) ser
			// prosa — o problema é lê-la de volta para decidir.
			if !strings.Contains(linha, "strings.Contains(") &&
				!strings.Contains(linha, "strings.HasPrefix(") &&
				!strings.Contains(linha, "strings.HasSuffix(") {
				continue
			}
			if m := prosa.FindString(linha); m != "" {
				t.Errorf("%s:%d confronta PROSA (%s) — traduza o projeto e o gate para de "+
					"disparar, sem erro nenhum. Use vocabulário estável: um marcador "+
					"(`@algo`, `[algo]`), um código, ou a ESTRUTURA do documento.",
					f, i+1, m)
			}
		}
	}
}
