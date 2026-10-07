// @anchors
//   code: NVAPN
//   ref: DCNAV

package doct

import (
	"path"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// NavScreen is a screen of the app: a spec of the `screen` layer.
type NavScreen struct {
	Spec
	Name  string // the spec's file stem — `GoalDetailScreen`
	Entry bool   // a route the app opens on (`navigation.entry`)
	// Reached says some entry leads to it; always true when no entry is declared, where
	// there is nothing to reach from.
	Reached bool
}

// NavEdge is one navigation the code declares: from a screen to another, by a rule.
type NavEdge struct {
	From, To NavScreen
	Rule     string
}

// NavMap is the app's navigation: its screens and every edge the code's `@navigates:`
// flags declare.
type NavMap struct {
	Screens []NavScreen
	Edges   []NavEdge
}

// fnNavigation draws the navigation map from the map's `navigates-to` edges — the flags on
// the code's navigation calls —, each end taken to the screen whose spec specifies its
// file. The flags are what the code does; the gates keep them equal to the Out tables.
func (c *Compiler) fnNavigation() NavMap {
	var m NavMap
	bySpec := map[string]int{}
	for _, sp := range c.specs {
		if sp.Layer != "screen" {
			continue
		}
		bySpec[sp.Path] = len(m.Screens)
		m.Screens = append(m.Screens, NavScreen{Spec: sp, Name: strings.TrimSuffix(path.Base(sp.Path), ".spec.md")})
	}
	if c.Graph == nil {
		return m
	}
	screenOf := map[string]int{} // a code file → the screen whose spec specifies it
	for _, e := range c.Graph.Edges {
		if e.Type == mapx.EdgeSpecifies {
			if i, ok := bySpec[e.From]; ok {
				screenOf[e.To] = i
			}
		}
	}
	seen := map[string]bool{}
	type pair struct{ from, to int }
	var pairs []pair
	var rules []string
	for _, e := range c.Graph.Edges {
		if e.Type != mapx.EdgeNavigatesTo {
			continue
		}
		// A flag names the screen by its code — its spec's —, so an end is the spec itself,
		// or a file the spec specifies.
		end := func(id string) (int, bool) {
			if i, ok := bySpec[id]; ok {
				return i, true
			}
			i, ok := screenOf[id]
			return i, ok
		}
		from, ok1 := end(e.From)
		to, ok2 := end(e.To)
		if !ok1 || !ok2 || from == to {
			continue
		}
		k := e.From + "\x00" + e.To + "\x00" + e.Method
		if seen[k] {
			continue
		}
		seen[k] = true
		pairs = append(pairs, pair{from, to})
		rules = append(rules, e.Method)
	}
	var entries []string
	if c.Config != nil && c.Config.Navigation != nil {
		entries = c.Config.Navigation.Entry
	}
	for i := range m.Screens {
		for _, en := range entries {
			if screenMatches(m.Screens[i].Name, en) {
				m.Screens[i].Entry = true
			}
		}
	}
	// Reachability: from each entry, along the edges; with no entry, nothing is unreached.
	next := map[int][]int{}
	for _, p := range pairs {
		next[p.from] = append(next[p.from], p.to)
	}
	var queue []int
	for i, s := range m.Screens {
		if s.Entry || len(entries) == 0 {
			m.Screens[i].Reached = true
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		for _, j := range next[i] {
			if !m.Screens[j].Reached {
				m.Screens[j].Reached = true
				queue = append(queue, j)
			}
		}
	}
	for k, p := range pairs {
		m.Edges = append(m.Edges, NavEdge{From: m.Screens[p.from], To: m.Screens[p.to], Rule: rules[k]})
	}
	sort.SliceStable(m.Edges, func(i, j int) bool {
		a, b := m.Edges[i], m.Edges[j]
		if a.From.Name != b.From.Name {
			return a.From.Name < b.From.Name
		}
		if a.To.Name != b.To.Name {
			return a.To.Name < b.To.Name
		}
		return a.Rule < b.Rule
	})
	sort.SliceStable(m.Screens, func(i, j int) bool { return m.Screens[i].Name < m.Screens[j].Name })
	return m
}

// screenMatches says whether an entry names a screen: its stem, or the stem without the
// `Screen` suffix, in any case.
func screenMatches(name, entry string) bool {
	e := strings.ToLower(strings.Trim(entry, "/ "))
	n := strings.ToLower(name)
	return e != "" && (e == n || e == strings.TrimSuffix(n, "screen"))
}

// Mermaid is the map as a Mermaid flowchart: an entry drawn as a stadium, a screen no entry
// reaches dashed, each edge labelled with its rule.
func (m NavMap) Mermaid() string {
	var b strings.Builder
	b.WriteString("flowchart LR\n")
	for _, s := range m.Screens {
		id := mermaidID(s)
		if s.Entry {
			b.WriteString("  " + id + "([" + s.Name + "])\n")
		} else {
			b.WriteString("  " + id + "[" + s.Name + "]\n")
		}
	}
	for _, e := range m.Edges {
		arrow := " --> "
		if e.Rule != "" {
			arrow = " -->|" + e.Rule + "| "
		}
		b.WriteString("  " + mermaidID(e.From) + arrow + mermaidID(e.To) + "\n")
	}
	var lost []string
	for _, s := range m.Screens {
		if !s.Reached {
			lost = append(lost, mermaidID(s))
		}
	}
	if len(lost) > 0 {
		b.WriteString("  classDef unreached stroke-dasharray: 5 5\n  class " + strings.Join(lost, ",") + " unreached\n")
	}
	return b.String()
}

// mermaidID is a screen's node id: its code, or its name when it has none.
func mermaidID(s NavScreen) string {
	if s.Code != "" {
		return s.Code
	}
	return s.Name
}

// hasScreens says whether any spec is of the `screen` layer.
func (c *Compiler) hasScreens() bool {
	for _, sp := range c.specs {
		if sp.Layer == "screen" {
			return true
		}
	}
	return false
}

// ScaffoldNavigation is the template of the app's navigation page, in the project's
// language: the flowchart, then one row per screen with where it leads.
func ScaffoldNavigation(lang string) Scaffold {
	title, cols, entry, unreached, why := "# Navigation", "| Screen | Code | Leads to |", "entry", "reached by no entry",
		"The app's navigation map, drawn from the `@navigates:` flags on the code's navigation calls: every screen, where it leads and by which rule."
	switch lang {
	case "pt", "pt-BR", "pt-br":
		title, cols, entry, unreached, why = "# Navegação", "| Tela | Código | Leva a |", "entrada", "nenhuma entrada alcança",
			"O mapa de navegação do app, desenhado das flags `@navigates:` nas chamadas de navegação do código: cada tela, aonde leva e por qual regra."
	case "es":
		title, cols, entry, unreached, why = "# Navegación", "| Pantalla | Código | Lleva a |", "entrada", "ninguna entrada la alcanza",
			"El mapa de navegación de la app, dibujado de las flags `@navigates:` en las llamadas de navegación del código: cada pantalla, adónde lleva y por cuál regla."
	}
	body := title + "\n\n{{with navigation}}{{$m := .}}```mermaid\n{{.Mermaid}}```\n\n" + cols + "\n| --- | --- | --- |\n" +
		"{{range $s := .Screens}}| {{$s.Name}}{{if $s.Entry}} (" + entry + "){{end}}{{if not $s.Reached}} — " + unreached + "{{end}} | `{{$s.Code}}` | " +
		"{{$n := 0}}{{range $m.Edges}}{{if eq .From.Path $s.Path}}{{if $n}}, {{end}}{{$n = 1}}{{.To.Name}}{{if .Rule}} (`{{.Rule}}`){{end}}{{end}}{{end}} |\n{{end}}{{end}}"
	return Scaffold{Nome: "navigation.md.tmpl", Corpo: body, Porque: why}
}
