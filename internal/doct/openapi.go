// @anchors
//   ref: OPNAP

package doct

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v3"
)

// fnOpenAPI compiles the project's OpenAPI document from the specs of its API units: every
// spec with an `Endpoint` section is an operation, and the contracts its body and responses
// cite are the shared schemas.
//
// The spec is the source and the OpenAPI is compiled, like every page of `docs/`: written by
// hand next to the specs, it was a second copy of the same contract, and the two drifted in
// silence. Compiled, it carries the generated marker, `docs-fresh` says when it is out of
// date, and the contract tests that validate the API against it run against what the specs
// say.
//
// The API spec says WHICH contract a body is, citing the model's spec by its code; the
// contract's spec says the fields, in its `Domain` table. A contract cited and not found
// fails the build: an OpenAPI with a hole would compile green.
func (c *Compiler) fnOpenAPI(title, version string, servers ...string) (string, error) {
	doc := oaDoc{OpenAPI: "3.1.0", Info: oaInfo{Title: title, Version: version}, Paths: map[string]map[string]*oaOperation{}}
	for _, s := range servers {
		doc.Servers = append(doc.Servers, oaServer{URL: s})
	}
	schemas := map[string]*oaSchema{}
	security := map[string]*oaSecurityScheme{}
	var used []Spec
	for _, sp := range c.specs {
		endpoints := sectionTable(sp, "Endpoint")
		if len(endpoints) == 0 {
			continue
		}
		used = append(used, sp)
		for _, ep := range endpoints {
			method := strings.ToLower(strings.TrimSpace(cell(ep, "method")))
			route := strings.Trim(cell(ep, "path"), "` ")
			if method == "" || route == "" {
				return "", fmt.Errorf("%s: an Endpoint row needs a method and a path", sp.Path)
			}
			op, err := c.operation(sp, ep, schemas, security)
			if err != nil {
				return "", err
			}
			if doc.Paths[route] == nil {
				doc.Paths[route] = map[string]*oaOperation{}
			}
			doc.Paths[route][method] = op
		}
	}
	if len(schemas) > 0 || len(security) > 0 {
		doc.Components = &oaComponents{Schemas: schemas, SecuritySchemes: security}
	}
	c.markConsumed(used)
	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(doc); err != nil {
		return "", err
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// operation builds one operation from an API spec and one of its Endpoint rows.
func (c *Compiler) operation(sp Spec, ep map[string]string, schemas map[string]*oaSchema, security map[string]*oaSecurityScheme) (*oaOperation, error) {
	op := &oaOperation{
		OperationID: strings.Trim(cell(ep, "operation"), "` "),
		Summary:     sp.Titulo,
		Description: strings.TrimSpace(fnSection(sp, "Overview")),
		Deprecated:  yes(cell(ep, "deprecated")),
		Responses:   map[string]*oaResponse{},
		XSpec:       sp.Path,
		XCode:       sp.Code,
	}
	for _, p := range sectionTable(sp, "Parameters") {
		op.Parameters = append(op.Parameters, oaParameter{
			Name:        strings.Trim(cell(p, "name"), "` "),
			In:          strings.ToLower(strings.TrimSpace(cell(p, "in"))),
			Required:    yes(cell(p, "required")),
			Description: cell(p, "description"),
			Schema:      typeSchema(cell(p, "type")),
		})
	}
	if rows := sectionTable(sp, "Request Body"); len(rows) > 0 {
		body := &oaRequestBody{Content: map[string]oaMedia{}}
		for _, r := range rows {
			ref, err := c.contractRef(sp, cell(r, "contract"), schemas)
			if err != nil {
				return nil, err
			}
			body.Required = body.Required || yes(cell(r, "required"))
			body.Content[contentType(cell(r, "content type"))] = oaMedia{Schema: ref}
		}
		op.RequestBody = body
	}
	for _, r := range sectionTable(sp, "Responses") {
		status := statusKey(cell(r, "status"))
		if status == "" {
			continue
		}
		resp := &oaResponse{Description: cell(r, "when")}
		if ref, err := c.contractRef(sp, cell(r, "contract"), schemas); err != nil {
			return nil, err
		} else if ref != nil {
			resp.Content = map[string]oaMedia{contentType(cell(r, "content type")): {Schema: ref}}
		}
		op.Responses[status] = resp
	}
	// Each refusal goes under its own status — declared in Responses or not —, with the
	// error code the client handles and the message it shows.
	for _, e := range sectionTable(sp, "Error Responses") {
		status := statusKey(cell(e, "status"))
		if status == "" {
			continue
		}
		resp := op.Responses[status]
		if resp == nil {
			resp = &oaResponse{Description: cell(e, "when")}
			if rng := op.Responses[status[:1]+"XX"]; rng != nil {
				resp.Content = rng.Content
			}
			op.Responses[status] = resp
		}
		resp.XErrors = append(resp.XErrors, oaError{
			Rule: strings.Trim(cell(e, "rule"), "` "), Code: strings.Trim(cell(e, "error code"), "` "),
			Message: strings.Trim(cell(e, "message"), `"`), When: cell(e, "when"),
		})
	}
	for _, s := range sectionTable(sp, "Security") {
		name := strings.Trim(cell(s, "scheme"), "` ")
		if name == "" {
			continue
		}
		security[name] = securityScheme(cell(s, "type"), cell(s, "where"))
		var scopes []string
		for _, sc := range strings.Split(cell(s, "scopes"), ",") {
			if sc = strings.Trim(strings.TrimSpace(sc), "`"); sc != "" && sc != "—" && sc != "-" {
				scopes = append(scopes, sc)
			}
		}
		if scopes == nil {
			scopes = []string{}
		}
		op.Security = append(op.Security, map[string][]string{name: scopes})
	}
	for _, l := range sectionTable(sp, "Limits") {
		op.XLimits = append(op.XLimits, oaLimit{Limit: cell(l, "limit"), Value: cell(l, "value"), Why: cell(l, "why")})
	}
	return op, nil
}

var codeInCellRE = regexp.MustCompile("`([A-Z0-9][A-Z0-9-]*)`|^([A-Z0-9]{3,})$")

// contractRef resolves a Contract cell — the cited code of the contract's spec — to a
// reference to its schema, building the schema from the contract's Domain once. An empty
// cell, or `—`, is a body-less response.
func (c *Compiler) contractRef(sp Spec, raw string, schemas map[string]*oaSchema) (*oaSchema, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "—" || raw == "-" {
		return nil, nil
	}
	m := codeInCellRE.FindStringSubmatch(raw)
	if m == nil {
		return nil, fmt.Errorf("%s: the contract %q is not a spec's code — cite it as `CODE`", sp.Path, raw)
	}
	code := m[1] + m[2]
	contract, err := c.fnSpecByCode(code)
	if err != nil {
		return nil, fmt.Errorf("%s: the contract `%s` is no spec of the project", sp.Path, code)
	}
	name := schemaName(contract.Path)
	if _, ok := schemas[name]; !ok {
		schemas[name] = domainSchema(*contract)
		c.markConsumed([]Spec{*contract})
	}
	return &oaSchema{Ref: "#/components/schemas/" + name}, nil
}

// domainSchema is the object a contract's Domain describes: one property per input, with
// the `Type` and `Required` columns when the table has them, and what it accepts as the
// description.
func domainSchema(sp Spec) *oaSchema {
	s := &oaSchema{Type: "object", Properties: map[string]*oaSchema{}, Description: sp.Titulo, XSpec: sp.Path}
	for _, r := range sectionTable(sp, "Domain") {
		name := strings.Trim(cell(r, "input"), "` ")
		if name == "" {
			continue
		}
		p := typeSchema(cell(r, "type"))
		if p == nil {
			p = &oaSchema{}
		}
		p.Description = cell(r, "accepts")
		s.Properties[name] = p
		if yes(cell(r, "required")) {
			s.Required = append(s.Required, name)
		}
	}
	sort.Strings(s.Required)
	return s
}

var typeFormatRE = regexp.MustCompile(`^([a-zA-Z]+)\s*(?:\(([^)]*)\))?$`)

// typeSchema reads a Type cell — `string (uuid)`, `integer`, `array of string` — into a
// schema; an empty cell is no type at all.
func typeSchema(raw string) *oaSchema {
	raw = strings.Trim(strings.TrimSpace(raw), "`")
	if raw == "" || raw == "—" {
		return nil
	}
	low := strings.ToLower(raw)
	for _, pre := range []string{"array of ", "list of ", "lista de ", "arreglo de "} {
		if strings.HasPrefix(low, pre) {
			return &oaSchema{Type: "array", Items: typeSchema(raw[len(pre):])}
		}
	}
	if strings.HasSuffix(raw, "[]") {
		return &oaSchema{Type: "array", Items: typeSchema(strings.TrimSuffix(raw, "[]"))}
	}
	m := typeFormatRE.FindStringSubmatch(raw)
	if m == nil {
		return &oaSchema{Description: raw}
	}
	t := map[string]string{"string": "string", "text": "string", "int": "integer", "integer": "integer", "long": "integer",
		"number": "number", "decimal": "number", "float": "number", "boolean": "boolean", "bool": "boolean", "object": "object"}[strings.ToLower(m[1])]
	if t == "" {
		return &oaSchema{Description: raw}
	}
	return &oaSchema{Type: t, Format: strings.TrimSpace(m[2])}
}

// securityScheme reads a Security row: `apiKey` with where it travels (`header x-api-key`),
// `http bearer` / `http basic`, or the type as written.
func securityScheme(typ, where string) *oaSecurityScheme {
	f := strings.Fields(strings.ToLower(strings.Trim(typ, "` ")))
	if len(f) == 0 {
		return &oaSecurityScheme{}
	}
	switch f[0] {
	case "apikey":
		s := &oaSecurityScheme{Type: "apiKey"}
		if w := strings.Fields(strings.Trim(where, "` ")); len(w) >= 2 {
			s.In, s.Name = strings.ToLower(w[0]), w[1]
		}
		return s
	case "http", "bearer", "basic":
		scheme := "bearer"
		if f[0] == "basic" || (len(f) > 1 && f[1] == "basic") {
			scheme = "basic"
		}
		return &oaSecurityScheme{Type: "http", Scheme: scheme}
	}
	return &oaSecurityScheme{Type: f[0], Description: strings.TrimSpace(where)}
}

func schemaName(specPath string) string {
	stem := strings.TrimSuffix(path.Base(specPath), ".spec.md")
	if stem == "" {
		return stem
	}
	r := []rune(stem)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func statusKey(raw string) string {
	raw = strings.Trim(strings.TrimSpace(raw), "`*")
	if len(raw) != 3 {
		return ""
	}
	return strings.ToUpper(raw)
}

func contentType(raw string) string {
	if raw = strings.Trim(strings.TrimSpace(raw), "`"); raw == "" || raw == "—" {
		return "application/json"
	}
	return raw
}

func yes(raw string) bool {
	switch strings.ToLower(strings.Trim(strings.TrimSpace(raw), "`*")) {
	case "yes", "sim", "sí", "si", "true", "x", "✓", "required", "obrigatório", "obligatorio":
		return true
	}
	return false
}

// columnAliases name each column the API sections use, in the catalog's languages.
var columnAliases = map[string][]string{
	"method": {"method", "metodo"}, "path": {"path", "caminho", "ruta"}, "operation": {"operation", "operacao", "operacion"},
	"deprecated": {"deprecated"}, "name": {"name", "nome", "nombre"}, "in": {"in", "em", "en"},
	"required": {"required", "obrigatorio", "obrigatoria", "obligatorio", "obligatoria"}, "type": {"type", "tipo"},
	"description": {"description", "descricao", "descripcion"}, "content type": {"content type"},
	"contract": {"contract", "contrato"}, "status": {"status"}, "when": {"when", "quando", "cuando"},
	"rule": {"rule", "regra", "regla"}, "error code": {"error code", "codigo de erro", "codigo de error"},
	"message": {"message", "mensagem", "mensaje"}, "scheme": {"scheme", "esquema"}, "where": {"where", "onde", "donde"},
	"scopes": {"scopes", "escopos"}, "limit": {"limit", "limite"}, "value": {"value", "valor"},
	"why": {"why", "por que", "por que?"}, "input": {"input", "entrada"}, "accepts": {"accepts", "aceita", "acepta"},
	"variable": {"variable", "variavel"}, "default": {"default", "padrao", "predeterminado"}, "values": {"values", "valores"},
}

// sectionTable reads the first table of a section — found under its title in any
// language — as rows keyed by the canonical column names.
func sectionTable(sp Spec, title string) []map[string]string {
	text := fnSection(sp, title)
	if text == "" {
		return nil
	}
	var header []string
	var rows []map[string]string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			if header != nil && len(rows) > 0 {
				break
			}
			continue
		}
		cells := splitRow(line)
		if header == nil {
			for _, h := range cells {
				header = append(header, canonicalColumn(h))
			}
			continue
		}
		if isSeparator(cells) {
			continue
		}
		row := map[string]string{}
		for i, v := range cells {
			if i < len(header) && header[i] != "" {
				row[header[i]] = v
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func splitRow(line string) []string {
	line = strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isSeparator(cells []string) bool {
	for _, c := range cells {
		if strings.Trim(c, ":- ") != "" {
			return false
		}
	}
	return true
}

func canonicalColumn(h string) string {
	h = foldAccents(strings.ToLower(strings.Trim(strings.TrimSpace(h), "`*")))
	for canon, aliases := range columnAliases {
		for _, a := range aliases {
			if h == a {
				return canon
			}
		}
	}
	return ""
}

func foldAccents(s string) string {
	var b strings.Builder
	for _, r := range norm.NFD.String(s) {
		if !unicode.Is(unicode.Mn, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func cell(row map[string]string, col string) string { return strings.TrimSpace(row[col]) }

// The document's shape, in the order OpenAPI writes it.
type oaDoc struct {
	OpenAPI    string                             `yaml:"openapi"`
	Info       oaInfo                             `yaml:"info"`
	Servers    []oaServer                         `yaml:"servers,omitempty"`
	Paths      map[string]map[string]*oaOperation `yaml:"paths"`
	Components *oaComponents                      `yaml:"components,omitempty"`
}

type oaInfo struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

type oaServer struct {
	URL string `yaml:"url"`
}

type oaOperation struct {
	OperationID string                 `yaml:"operationId,omitempty"`
	Summary     string                 `yaml:"summary,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Deprecated  bool                   `yaml:"deprecated,omitempty"`
	Parameters  []oaParameter          `yaml:"parameters,omitempty"`
	RequestBody *oaRequestBody         `yaml:"requestBody,omitempty"`
	Responses   map[string]*oaResponse `yaml:"responses"`
	Security    []map[string][]string  `yaml:"security,omitempty"`
	XLimits     []oaLimit              `yaml:"x-limits,omitempty"`
	XSpec       string                 `yaml:"x-anchors-spec,omitempty"`
	XCode       string                 `yaml:"x-anchors-code,omitempty"`
}

type oaParameter struct {
	Name        string    `yaml:"name"`
	In          string    `yaml:"in"`
	Required    bool      `yaml:"required,omitempty"`
	Description string    `yaml:"description,omitempty"`
	Schema      *oaSchema `yaml:"schema,omitempty"`
}

type oaRequestBody struct {
	Required bool               `yaml:"required,omitempty"`
	Content  map[string]oaMedia `yaml:"content"`
}

type oaMedia struct {
	Schema *oaSchema `yaml:"schema,omitempty"`
}

type oaResponse struct {
	Description string             `yaml:"description"`
	Content     map[string]oaMedia `yaml:"content,omitempty"`
	XErrors     []oaError          `yaml:"x-error-codes,omitempty"`
}

type oaError struct {
	Rule    string `yaml:"rule,omitempty"`
	Code    string `yaml:"code,omitempty"`
	Message string `yaml:"message,omitempty"`
	When    string `yaml:"when,omitempty"`
}

type oaLimit struct {
	Limit string `yaml:"limit"`
	Value string `yaml:"value,omitempty"`
	Why   string `yaml:"why,omitempty"`
}

type oaSchema struct {
	Ref         string               `yaml:"$ref,omitempty"`
	Type        string               `yaml:"type,omitempty"`
	Format      string               `yaml:"format,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Items       *oaSchema            `yaml:"items,omitempty"`
	Properties  map[string]*oaSchema `yaml:"properties,omitempty"`
	Required    []string             `yaml:"required,omitempty"`
	XSpec       string               `yaml:"x-anchors-spec,omitempty"`
}

type oaComponents struct {
	Schemas         map[string]*oaSchema         `yaml:"schemas,omitempty"`
	SecuritySchemes map[string]*oaSecurityScheme `yaml:"securitySchemes,omitempty"`
}

type oaSecurityScheme struct {
	Type        string `yaml:"type,omitempty"`
	Scheme      string `yaml:"scheme,omitempty"`
	In          string `yaml:"in,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}
