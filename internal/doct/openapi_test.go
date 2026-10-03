// @anchors
//   ref: OPNAP

package doct

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const apiSpec = `<!-- @anchors
  code: GQCHG
-->
# Generate a QR Code charge

## Overview

Creates a PIX QR Code charge for the marketplace's wallet.

## Endpoint
| Method | Path | Operation | Deprecated |
| --- | --- | --- | --- |
| POST | /v1/qrcodes | generateQrCode | no |

## Parameters
| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| x-marketplace-id | header | yes | string (uuid) | the marketplace that calls |

## Request Body
| Content type | Contract | Required |
| --- | --- | --- |
| application/json | ` + "`QRBDY`" + ` | yes |

## Responses
| Status | When | Content type | Contract |
| --- | --- | --- | --- |
| 201 | the QR Code was created | application/json | ` + "`QRCDM`" + ` |
| 4xx | the request was refused | application/json | ` + "`QRCDM`" + ` |
| 204 | nothing to do | — | — |

## Error Responses
| Rule | When | Status | Error code | Message |
| --- | --- | --- | --- | --- |
| ` + "`GQCHG-E01`" + ` | the wallet does not exist | 404 | WALLET_NOT_FOUND | "Wallet not found" |

## Security
| Scheme | Type | Where | Scopes |
| --- | --- | --- | --- |
| marketplaceKey | apiKey | header x-api-key | — |

## Limits
| Limit | Value | Why |
| --- | --- | --- |
| rate | 100/min | the provider's quota |
`

// A contract written in Portuguese: the section and its columns are read in any language.
const bodySpec = `<!-- @anchors
  code: QRBDY
-->
# Corpo da cobrança

## Domínio

| Entrada | Tipo | Obrigatório | Aceita | Fora do domínio | Quem garante |
| --- | --- | --- | --- | --- | --- |
| ` + "`walletId`" + ` | string (uuid) | sim | um UUID canônico | o resto | esta unidade |
| ` + "`amount`" + ` | integer | sim | centavos, > 0 | zero | esta unidade |
| ` + "`tags`" + ` | array of string | não | rótulos livres | — | esta unidade |
`

const qrSpec = `<!-- @anchors
  code: QRCDM
-->
# QR Code

## Domain

| Input | Type | Required | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- | --- | --- |
| ` + "`id`" + ` | string | yes | the charge id | — | payments |
`

func buildOpenAPI(t *testing.T, specs map[string]string) (string, error) {
	t.Helper()
	root, g := projetoDeTeste(t, specs)
	escreveTemplate(t, root, "openapi.yaml.tmpl", `{{ openapi "BaaS Proxy" "1.2.0" "https://api.example.com" }}`)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "openapi.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b), nil
}

func TestOpenAPI_anOperationFromAnAPISpec(t *testing.T) {
	t.Run("OPNAP-B01: Each Endpoint row of a spec is an operation, with its parameters, body, responses, errors, security and limits", func(t *testing.T) {})
	t.Run("OPNAP-B02: A cited contract is a shared schema built from its Domain, in any language", func(t *testing.T) {})
	t.Run("OPNAP-B03: The document is written in OpenAPI's order, two spaces deep", func(t *testing.T) {})
	t.Run("OPNAP-B04: The compiled document carries the generated marker as a YAML comment and is valid YAML", func(t *testing.T) {})
	out, err := buildOpenAPI(t, map[string]string{"src/generateQrCode.spec.md": apiSpec,
		"src/models/qrBody.spec.md": bodySpec, "src/models/qrCode.spec.md": qrSpec})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "# anchors:generated from doct/openapi.yaml.tmpl") {
		t.Errorf("the marker is a YAML comment:\n%.120s", out)
	}
	var doc map[string]any
	if err := yaml.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("the document is YAML: %v\n%s", err, out)
	}
	for _, want := range []string{
		"openapi: 3.1.0", "title: BaaS Proxy", "version: 1.2.0", "url: https://api.example.com",
		"/v1/qrcodes:", "post:", "operationId: generateQrCode", "summary: Generate a QR Code charge",
		"Creates a PIX QR Code charge", "name: x-marketplace-id", "in: header", "format: uuid",
		"$ref: '#/components/schemas/QrBody'", "\"201\":", "4XX:", "\"204\":", "\"404\":",
		"code: WALLET_NOT_FOUND", "message: Wallet not found", "marketplaceKey:", "type: apiKey", "name: x-api-key",
		"x-limits:", "value: 100/min", "x-anchors-code: GQCHG",
		"QrBody:", "walletId:", "amount:", "type: integer", "- amount", "- walletId", "items:",
		"QrCode:",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the OpenAPI lacks %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "\nopenapi: 3.1.0\ninfo:\n  title:") || strings.Index(out, "\npaths:") > strings.Index(out, "\ncomponents:") {
		t.Errorf("the keys come in OpenAPI's order, two spaces deep:\n%.300s", out)
	}
	if i := strings.Index(out, "\"204\":"); i < 0 || strings.Contains(out[i:i+80], "content") {
		t.Errorf("a body-less response has no content:\n%s", out)
	}
	if strings.Contains(out, "- tags") {
		t.Error("a field the contract does not require is not required")
	}
}

func TestOpenAPI_aContractThatIsNoSpecFailsTheBuild(t *testing.T) {
	t.Run("OPNAP-E01: A contract cited and not found fails the build, naming it", func(t *testing.T) {})
	t.Run("OPNAP-E02: An Endpoint row with no path fails the build", func(t *testing.T) {})
	_, err := buildOpenAPI(t, map[string]string{"src/generateQrCode.spec.md": apiSpec, "src/models/qrCode.spec.md": qrSpec})
	if err == nil || !strings.Contains(err.Error(), "QRBDY") {
		t.Errorf("the missing contract fails the build by name, got %v", err)
	}
	noPath := strings.Replace(apiSpec, "| POST | /v1/qrcodes |", "| POST |  |", 1)
	if _, err := buildOpenAPI(t, map[string]string{"src/generateQrCode.spec.md": noPath}); err == nil || !strings.Contains(err.Error(), "generateQrCode.spec.md") {
		t.Errorf("an Endpoint row with no path fails the build naming the spec, got %v", err)
	}
}

func TestOpenAPI_aCompiledYAMLIsKeptFresh(t *testing.T) {
	t.Run("OPNAP-B05: A compiled YAML is recognised as generated and its freshness checked", func(t *testing.T) {})
	root, g := projetoDeTeste(t, map[string]string{"src/generateQrCode.spec.md": apiSpec,
		"src/models/qrBody.spec.md": bodySpec, "src/models/qrCode.spec.md": qrSpec})
	escreveTemplate(t, root, "openapi.yaml.tmpl", `{{ openapi "BaaS Proxy" "1.2.0" }}`)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	if stale, err := c.Stale(); err != nil || len(stale) != 0 {
		t.Fatalf("fresh after the build: %v %v", stale, err)
	}
	if err := os.WriteFile(filepath.Join(root, "src/models/qrCode.spec.md"), []byte(strings.Replace(qrSpec, "the charge id", "the charge's id", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	c2, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if stale, err := c2.Stale(); err != nil || len(stale) != 1 || stale[0] != "openapi.yaml" {
		t.Errorf("a contract changed makes the OpenAPI stale: %v %v", stale, err)
	}
	if res, _ := c2.Build(false); len(res.Skipped) != 0 {
		t.Errorf("a generated YAML is overwritten, never skipped: %+v", res)
	}
}
