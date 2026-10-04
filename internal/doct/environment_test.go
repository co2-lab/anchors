// @anchors
//   ref: DCENV

package doct

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const envSpecEN = `<!-- @anchors
  code: PAYMT
-->
# Pay

## Environment Variables
| Variable | Type | Required | Default | Values | Deprecated | Description |
| --- | --- | --- | --- | --- | --- | --- |
| PAYMENTS_URL | string (uri) | yes | — | — | no | the payments base URL |
| TODO: VARIABLE_NAME | TODO | yes | — | — | no | TODO |
`

const envSpecPT = `<!-- @anchors
  code: LOGGR
-->
# Logger

## Variáveis de Ambiente
| Variável | Tipo | Obrigatória | Default | Valores | Deprecated | Descrição |
| --- | --- | --- | --- | --- | --- | --- |
| LOG_LEVEL | string | não | info | debug, info, warn | não | quanto se registra |
| PAYMENTS_URL | string | sim | — | — | não | |
`

func TestEnvVars_listsEveryDeclaredVariable(t *testing.T) {
	t.Run("DCENV-B01: Every declared variable is listed by name, in any language", func(t *testing.T) {})
	t.Run("DCENV-B02: A variable two units read is listed once, with both", func(t *testing.T) {})
	t.Run("DCENV-E01: A placeholder row is no variable", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pay/pay.spec.md": envSpecEN, "log/logger.spec.md": envSpecPT})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	vars := c.fnEnvVars()
	var names []string
	for _, v := range vars {
		names = append(names, v.Name)
	}
	if strings.Join(names, ",") != "LOG_LEVEL,PAYMENTS_URL" {
		t.Fatalf("variables = %v", names)
	}
	if v := vars[0]; v.Default != "info" || v.Values != "debug, info, warn" || v.Required != "não" {
		t.Errorf("LOG_LEVEL read from the Portuguese table: %+v", v)
	}
	if v := vars[1]; len(v.Units) != 2 || v.Description != "the payments base URL" {
		t.Errorf("PAYMENTS_URL once, with both units and the first description: %+v", v)
	}
}

func TestScaffoldEnvironment_thePageIsSeededAndCompiled(t *testing.T) {
	t.Run("DCENV-B03: The environment page is seeded and compiled", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"pay/pay.spec.md": envSpecEN, "log/logger.spec.md": envSpecPT})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	written, _, err := c.InitScaffolds(false)
	if err != nil || !strings.Contains(strings.Join(written, ","), "environment.md.tmpl") {
		t.Fatalf("the environment template is seeded: %v %v", written, err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "environment.md"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(b)
	for _, want := range []string{"| `LOG_LEVEL` | string | não | info |", "| `PAYMENTS_URL` |", "`LOGGR`, `PAYMT`"} {
		if !strings.Contains(page, want) {
			t.Errorf("the page lacks %q:\n%s", want, page)
		}
	}
}
