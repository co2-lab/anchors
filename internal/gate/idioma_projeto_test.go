package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// NENHUM identificador do projeto está em português — e este teste é o que impede o
// trabalho de migração de se desfazer.
//
// Ele varre o código de PRODUÇÃO (os `_test.go` ficam de fora: um teste pode nomear um
// caso em português para descrever o que testa) e falha nomeando cada acusação com a
// palavra que a motivou.
//
// A régua está no `idioma_do_codigo.go`: identificadores em inglês, comentários no idioma
// do time, mensagens pelo i18n.
func TestNenhumIdentificadorEmPortugues(t *testing.T) {
	raiz := "../.."
	achados := map[string]string{}

	err := filepath.Walk(raiz, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || !strings.HasSuffix(p, ".go") {
			return nil
		}
		if strings.Contains(p, "vendor") || strings.Contains(p, "testdata") ||
			strings.HasSuffix(p, "_test.go") {
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return nil
		}
		for id, palavra := range IdentificadoresPT(string(b)) {
			achados[id] = palavra + "  " + p
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(achados) == 0 {
		return
	}
	t.Errorf("%d identificador(es) em português — o código é em inglês, "+
		"e os comentários é que ficam no idioma do time:", len(achados))
	for id, onde := range achados {
		t.Errorf("    %-40s (%s)", id, onde)
	}
}
