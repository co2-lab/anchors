// @anchors
//   ref: DCENV

package doct

import (
	"sort"
	"strings"
)

// EnvVar is an environment variable of the project, as its specs declare it, with the units
// that read it.
type EnvVar struct {
	Name, Type, Required, Default, Values, Deprecated, Description string
	Units                                                          []Spec
}

// fnEnvVars gathers the environment variables every spec declares in its Environment
// Variables section into one list, by name: the page whoever deploys reads. A variable two
// units read is listed once, with both; its description is the first spec's that gives one.
func (c *Compiler) fnEnvVars() []EnvVar {
	byName := map[string]*EnvVar{}
	var used []Spec
	for _, sp := range c.specs {
		rows := sectionTable(sp, "Environment Variables")
		if len(rows) == 0 {
			continue
		}
		used = append(used, sp)
		for _, r := range rows {
			name := strings.Trim(cell(r, "variable"), "` ")
			if name == "" || strings.HasPrefix(strings.ToUpper(name), "TODO") {
				continue
			}
			v := byName[name]
			if v == nil {
				v = &EnvVar{Name: name, Type: cell(r, "type"), Required: cell(r, "required"), Default: cell(r, "default"),
					Values: cell(r, "values"), Deprecated: cell(r, "deprecated"), Description: cell(r, "description")}
				byName[name] = v
			}
			if v.Description == "" {
				v.Description = cell(r, "description")
			}
			v.Units = append(v.Units, sp)
		}
	}
	c.markConsumed(used)
	out := make([]EnvVar, 0, len(byName))
	for _, v := range byName {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// hasEnvSpecs says whether any spec declares an environment variable.
func (c *Compiler) hasEnvSpecs() bool {
	for _, sp := range c.specs {
		if len(sectionTable(sp, "Environment Variables")) > 0 {
			return true
		}
	}
	return false
}

// ScaffoldEnvironment is the template of the project's environment-variables page, in the
// project's language: one table of every variable, with the units that read it.
func ScaffoldEnvironment(lang string) Scaffold {
	title, cols, why := "# Environment variables", "| Variable | Type | Required | Default | Values | Deprecated | Description | Read by |",
		"Every environment variable the project reads, compiled from the Environment Variables section of its specs: the page whoever deploys reads."
	switch lang {
	case "pt", "pt-BR", "pt-br":
		title, cols, why = "# Variáveis de ambiente", "| Variável | Tipo | Obrigatória | Default | Valores | Deprecated | Descrição | Lida por |",
			"Toda variável de ambiente que o projeto lê, compilada da seção Variáveis de Ambiente das specs: a página de quem implanta."
	case "es":
		title, cols, why = "# Variables de entorno", "| Variable | Tipo | Obligatoria | Default | Valores | Deprecated | Descripción | Leída por |",
			"Toda variable de entorno que el proyecto lee, compilada de la sección Variables de Entorno de las specs: la página de quien despliega."
	}
	body := title + "\n\n" + cols + "\n| --- | --- | --- | --- | --- | --- | --- | --- |\n" +
		"{{range envVars}}| `{{.Name}}` | {{.Type}} | {{.Required}} | {{.Default}} | {{.Values}} | {{.Deprecated}} | {{.Description}} | {{range $i, $u := .Units}}{{if $i}}, {{end}}`{{$u.Code}}`{{end}} |\n{{end}}"
	return Scaffold{Nome: "environment.md.tmpl", Corpo: body, Porque: why}
}
