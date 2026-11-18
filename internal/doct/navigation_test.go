// @anchors
//   code: NVTSN
//   ref: DCNAV

package doct

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func screenSpec(code string) string {
	return "<!-- @anchors\n  code: " + code + "\n  layer: screen\n-->\n# S\n"
}

// navApp is an app of four screens: Home leads to GoalDetail, which leads to GoalEdit, which
// goes back; Lonely is reached by nothing. Each screen's spec specifies its code file, and
// the flags on the calls are the map's navigates-to edges.
func navApp(t *testing.T) *Compiler {
	t.Helper()
	root, g := projetoDeTeste(t, map[string]string{
		"s/HomeScreen.spec.md":       screenSpec("HOMEH"),
		"s/GoalDetailScreen.spec.md": screenSpec("GLDTG"),
		"s/GoalEditScreen.spec.md":   screenSpec("GLETG"),
		"s/LonelyScreen.spec.md":     screenSpec("LNLYS"),
		"s/other.spec.md":            "<!-- @anchors\n  code: OTHRS\n-->\n# Other\n",
	})
	for _, s := range []string{"Home", "GoalDetail", "GoalEdit", "Lonely"} {
		code := "s/" + s + "Screen.tsx"
		g.Nodes = append(g.Nodes, mapx.Node{ID: code, Kind: mapx.KindCode, FileCode: strings.ToUpper(s[:4]) + "C"})
		g.Edges = append(g.Edges, mapx.Edge{From: "s/" + s + "Screen.spec.md", To: code, Type: mapx.EdgeSpecifies})
	}
	g.Edges = append(g.Edges,
		// The flag names the screen's code, its spec's: the edge lands on the spec.
		mapx.Edge{From: "s/HomeScreen.tsx", To: "s/GoalDetailScreen.spec.md", Type: mapx.EdgeNavigatesTo, Method: "HOMEH-A01"},
		mapx.Edge{From: "s/GoalDetailScreen.tsx", To: "s/GoalEditScreen.spec.md", Type: mapx.EdgeNavigatesTo, Method: "GLDTG-A01"},
		mapx.Edge{From: "s/GoalEditScreen.tsx", To: "s/GoalDetailScreen.tsx", Type: mapx.EdgeNavigatesTo})
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	c.Config = &config.Config{Navigation: &config.Navigation{Entry: []string{"Home"}}}
	return c
}

func TestNavigation_theMapOfTheApp(t *testing.T) {
	t.Run("DCNAV-B01: navigation lists the screen specs and the edges the code's flags declare, each end taken to its screen", func(t *testing.T) {})
	t.Run("DCNAV-B02: An entry is marked, and a screen no entry reaches is marked unreached; with no entry, none is", func(t *testing.T) {})
	c := navApp(t)
	m := c.fnNavigation()
	var names []string
	for _, s := range m.Screens {
		names = append(names, s.Name)
	}
	if strings.Join(names, ",") != "GoalDetailScreen,GoalEditScreen,HomeScreen,LonelyScreen" {
		t.Fatalf("screens = %v", names)
	}
	var edges []string
	for _, e := range m.Edges {
		edges = append(edges, e.From.Code+">"+e.To.Code+":"+e.Rule)
	}
	if strings.Join(edges, " ") != "GLDTG>GLETG:GLDTG-A01 GLETG>GLDTG: HOMEH>GLDTG:HOMEH-A01" {
		t.Errorf("edges = %v", edges)
	}
	if !m.Screens[2].Entry || m.Screens[3].Reached || !m.Screens[1].Reached {
		t.Errorf("Home is the entry, Lonely unreached, GoalEdit reached: %+v", m.Screens)
	}
	mer := m.Mermaid()
	for _, want := range []string{"HOMEH([HomeScreen])", "HOMEH -->|HOMEH-A01| GLDTG", "GLETG --> GLDTG", "class LNLYS unreached"} {
		if !strings.Contains(mer, want) {
			t.Errorf("the flowchart lacks %q:\n%s", want, mer)
		}
	}
	c.Config = &config.Config{}
	for _, s := range c.fnNavigation().Screens {
		if !s.Reached {
			t.Errorf("with no entry, nothing is unreached: %+v", s)
		}
	}
}

func TestNavigation_thePagesAreSeededAndCompiled(t *testing.T) {
	t.Run("DCNAV-B03: ScaffoldNavigation and ScaffoldDependencies are seeded when the app has screens and dependency flags, and compile into the pages", func(t *testing.T) {})
	c := navApp(t)
	c.Graph.Edges = append(c.Graph.Edges, mapx.Edge{From: "s/HomeScreen.tsx", To: "s/GoalDetailScreen.tsx", Type: mapx.EdgeDependsOn, Origin: mapx.OriginDeclared, Method: "GoalCard, default"})
	written, _, err := c.InitScaffolds(false)
	if err != nil || !strings.Contains(strings.Join(written, ","), "navigation.md.tmpl") || !strings.Contains(strings.Join(written, ","), "dependencies.md.tmpl") {
		t.Fatalf("both templates are seeded: %v %v", written, err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}
	nav, err := os.ReadFile(filepath.Join(c.Root, OutDir, "navigation.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"```mermaid\nflowchart LR", "| HomeScreen (entry) | `HOMEH` | GoalDetailScreen (`HOMEH-A01`) |", "| LonelyScreen — reached by no entry | `LNLYS` |  |", "| GoalEditScreen | `GLETG` | GoalDetailScreen |"} {
		if !strings.Contains(string(nav), want) {
			t.Errorf("the navigation page lacks %q:\n%s", want, nav)
		}
	}
	deps, err := os.ReadFile(filepath.Join(c.Root, OutDir, "dependencies.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"| `s/HomeScreen.tsx` | `HOMEC` | `GOALC` (GoalCard, default) |  |", "| `s/GoalDetailScreen.tsx` | `GOALC` |  | `HOMEC` |"} {
		if !strings.Contains(string(deps), want) {
			t.Errorf("the dependencies page lacks %q:\n%s", want, deps)
		}
	}
}

func TestNavigation_aComponentNavigatesForItsScreens(t *testing.T) {
	t.Run("DCNAV-B04: A navigation flagged in a component is drawn from every screen that renders it along the dependency chain", func(t *testing.T) {})
	c := navApp(t)
	c.Graph.Nodes = append(c.Graph.Nodes, mapx.Node{ID: "s/Card.tsx", Kind: mapx.KindCode, FileCode: "CARDC"})
	c.Graph.Edges = append(c.Graph.Edges,
		mapx.Edge{From: "s/HomeScreen.tsx", To: "s/Card.tsx", Type: mapx.EdgeDependsOn},
		mapx.Edge{From: "s/Card.tsx", To: "s/LonelyScreen.spec.md", Type: mapx.EdgeNavigatesTo, Method: "HOMEH-A02"})
	m := c.fnNavigation()
	var edges []string
	for _, e := range m.Edges {
		edges = append(edges, e.From.Code+">"+e.To.Code+":"+e.Rule)
	}
	if !strings.Contains(strings.Join(edges, " "), "HOMEH>LNLYS:HOMEH-A02") {
		t.Errorf("the card's navigation is drawn from Home: %v", edges)
	}
	for _, s := range m.Screens {
		if s.Code == "LNLYS" && !s.Reached {
			t.Error("Lonely is reached through the card")
		}
	}
}
