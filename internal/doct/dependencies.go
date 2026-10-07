// @anchors
//   code: DPAPD
//   ref: DCNAV

package doct

import (
	"sort"

	"github.com/co2-lab/anchors/internal/mapx"
)

// DepLink is one end of a dependency between code files: the other file, its code, and
// the symbols the import brings.
type DepLink struct {
	Path, Code, Symbols string
}

// DepFile is a code file of the dependency chain: what it uses and who uses it.
type DepFile struct {
	Path, Code   string
	Uses, UsedBy []DepLink
}

// fnDependencies lists the dependency chain from the map's `depends-on` edges between code
// files — the `@dep:` flags on the import lines —, one entry per file that takes part, with
// what it uses and, in reverse, who uses it. A file outside the chain is left out.
func (c *Compiler) fnDependencies() []DepFile {
	if c.Graph == nil {
		return nil
	}
	byPath := map[string]*DepFile{}
	at := func(n *mapx.Node) *DepFile {
		f := byPath[n.ID]
		if f == nil {
			f = &DepFile{Path: n.ID, Code: n.FileCode}
			byPath[n.ID] = f
		}
		return f
	}
	for _, e := range c.Graph.Edges {
		if e.Type != mapx.EdgeDependsOn {
			continue
		}
		from, to := c.Graph.Node(e.From), c.Graph.Node(e.To)
		if from == nil || to == nil || from.Kind != mapx.KindCode || to.Kind != mapx.KindCode || from.FileCode == "" || to.FileCode == "" {
			continue
		}
		at(from).Uses = append(at(from).Uses, DepLink{Path: to.ID, Code: to.FileCode, Symbols: e.Method})
		at(to).UsedBy = append(at(to).UsedBy, DepLink{Path: from.ID, Code: from.FileCode, Symbols: e.Method})
	}
	out := make([]DepFile, 0, len(byPath))
	for _, f := range byPath {
		byCode := func(l []DepLink) func(i, j int) bool { return func(i, j int) bool { return l[i].Code < l[j].Code } }
		sort.Slice(f.Uses, byCode(f.Uses))
		sort.Slice(f.UsedBy, byCode(f.UsedBy))
		out = append(out, *f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// hasDependencyFlags says whether any code file declares a dependency on another.
func (c *Compiler) hasDependencyFlags() bool {
	return len(c.fnDependencies()) > 0
}

// ScaffoldDependencies is the template of the project's dependency page, in the project's
// language: one row per file of the chain, with what it uses and who uses it.
func ScaffoldDependencies(lang string) Scaffold {
	title, cols, why := "# Dependencies", "| File | Code | Uses | Used by |",
		"Every file of the dependency chain, drawn from the `@dep:` flags on the import lines: what each file uses, with the symbols, and who uses it — `anchors map deps <CODE>` walks it as a tree."
	switch lang {
	case "pt", "pt-BR", "pt-br":
		title, cols, why = "# Dependências", "| Arquivo | Código | Usa | Usado por |",
			"Cada arquivo da cadeia de dependências, desenhada das flags `@dep:` nas linhas de import: o que cada arquivo usa, com os símbolos, e quem o usa — `anchors map deps <CODE>` percorre como árvore."
	case "es":
		title, cols, why = "# Dependencias", "| Archivo | Código | Usa | Usado por |",
			"Cada archivo de la cadena de dependencias, dibujada de las flags `@dep:` en las líneas de import: lo que cada archivo usa, con los símbolos, y quién lo usa — `anchors map deps <CODE>` la recorre como árbol."
	}
	body := title + "\n\n" + cols + "\n| --- | --- | --- | --- |\n" +
		"{{range dependencies}}| `{{.Path}}` | `{{.Code}}` | {{range $i, $u := .Uses}}{{if $i}}, {{end}}`{{$u.Code}}`{{if $u.Symbols}} ({{$u.Symbols}}){{end}}{{end}} | " +
		"{{range $i, $u := .UsedBy}}{{if $i}}, {{end}}`{{$u.Code}}`{{end}} |\n{{end}}"
	return Scaffold{Nome: "dependencies.md.tmpl", Corpo: body, Porque: why}
}
