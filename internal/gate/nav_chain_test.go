// @anchors
//   code: NCTNV
//   ref: NCGNV

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func screenSpecText(code, route, in, out string) string {
	return "<!-- @anchors\n  code: " + code + "\n  layer: screen\n-->\n# S\n\n> **Route**: `" + route + "`\n\n## Navigation\n\n### In\n\n" +
		"| Origin | Action |\n| --- | --- |\n" + in + "\n### Out\n\n| Rule | Destination | Action |\n| --- | --- | --- |\n" + out + "\n## Rules\n"
}

// navProject is an app: Home leads to GoalDetail, which leads to GoalEdit, which goes back
// — GoalDetail's In forgets GoalEdit —, and Lonely is reached by nothing.
func navProject(t *testing.T) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"s/HomeScreen.spec.md":       screenSpecText("HOMEH", "Home", "", "| `HOMEH-A01` | GoalDetailScreen | tap a goal |\n"),
		"s/GoalDetailScreen.spec.md": screenSpecText("GLDTG", "GoalDetail", "| HomeScreen | tap |\n", "| `GLDTG-A01` | GoalEditScreen | edit |\n"),
		"s/GoalEditScreen.spec.md":   screenSpecText("GLETG", "GoalEdit", "| GoalDetailScreen | edit |\n", "| `GLETG-A02` | GoalDetailScreen | save |\n"),
		"s/LonelyScreen.spec.md":     screenSpecText("LNLYS", "Lonely", "", ""),
		"s/HomeScreen.tsx":           "onPress={() => navigation.navigate('GoalDetail')}\n",
		"s/GoalDetailScreen.tsx":     "navigation.navigate('GoalEdit') // @navigates: GLETG [GLDTG-A01]\n// @no-nav: closes the sheet of this screen\nnavigation.goBack()\n",
		"s/GoalEditScreen.tsx":       "navigation.navigate('GoalDetail') // @navigates: HOMEH\n",
		"s/LonelyScreen.tsx":         "export const L = 1\n",
	}
	for rel, body := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{}
	for _, s := range []struct{ name, code, file string }{{"Home", "HOMEH", "HMSCR"}, {"GoalDetail", "GLDTG", "GDSCR"}, {"GoalEdit", "GLETG", "GESCR"}, {"Lonely", "LNLYS", "LNSCR"}} {
		spec := "s/" + s.name + "Screen.spec.md"
		code := "s/" + s.name + "Screen.tsx"
		g.Nodes = append(g.Nodes, mapx.Node{ID: spec, Kind: mapx.KindSpec, Code: s.code, Layer: "spec"},
			mapx.Node{ID: code, Kind: mapx.KindCode, Code: s.code, FileCode: s.file})
		g.Edges = append(g.Edges, mapx.Edge{From: spec, To: code, Type: mapx.EdgeSpecifies})
	}
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts"}, Navigation: &config.Navigation{Entry: []string{"Home"}}}
	screensMu.Lock()
	delete(screensCache, root)
	screensMu.Unlock()
	return root, g, cfg
}

func navNode(g *mapx.Graph, id string) mapx.Node { return *g.Node(id) }

func TestNavChain_annotated(t *testing.T) {
	t.Run("NCGNV-B01: nav-annotated names each navigation call with no flag, and each flag naming another screen than its route", func(t *testing.T) {})
	root, g, cfg := navProject(t)
	v, msg := checkNavAnnotated(read(t, root, "s/HomeScreen.tsx"), navNode(g, "s/HomeScreen.tsx"), root, g, cfg)
	if v != Fail || !strings.Contains(msg, "line 1") || !strings.Contains(msg, "GoalDetail") || !strings.Contains(msg, "anchors guide navigation") {
		t.Errorf("the unflagged call is named: %v %s", v, msg)
	}
	if v, msg := checkNavAnnotated(read(t, root, "s/GoalDetailScreen.tsx"), navNode(g, "s/GoalDetailScreen.tsx"), root, g, cfg); v != Pass {
		t.Errorf("a right flag and a waived back navigation pass: %v %s", v, msg)
	}
	v, msg = checkNavAnnotated(read(t, root, "s/GoalEditScreen.tsx"), navNode(g, "s/GoalEditScreen.tsx"), root, g, cfg)
	if v != Fail || !strings.Contains(msg, "HOMEH") || !strings.Contains(msg, "GLDTG") {
		t.Errorf("a flag naming another screen than its route is named: %v %s", v, msg)
	}
}

func TestNavChain_matchesSpecAndSymmetric(t *testing.T) {
	t.Run("NCGNV-B02: nav-matches-spec names each Out row no flag answers, and each flag the Out table does not declare", func(t *testing.T) {})
	t.Run("NCGNV-B03: nav-symmetric names each Out with no matching In, and each In with no matching Out", func(t *testing.T) {})
	root, g, cfg := navProject(t)
	if v, msg := checkNavMatchesSpec("", navNode(g, "s/GoalDetailScreen.spec.md"), root, g, cfg); v != Pass {
		t.Errorf("GoalDetail's Out and its flag agree: %v %s", v, msg)
	}
	v, msg := checkNavMatchesSpec("", navNode(g, "s/GoalEditScreen.spec.md"), root, g, cfg)
	if v != Fail || !strings.Contains(msg, "GoalDetailScreen") || !strings.Contains(msg, "HomeScreen") {
		t.Errorf("GoalEdit's Out says GoalDetail, its flag says Home: %v %s", v, msg)
	}
	v, msg = checkNavSymmetric("", navNode(g, "s/GoalEditScreen.spec.md"), root, g, cfg)
	if v != Fail || !strings.Contains(msg, "GoalDetailScreen") {
		t.Errorf("GoalEdit leads to GoalDetail, whose In forgets it: %v %s", v, msg)
	}
	if v, msg := checkNavSymmetric("", navNode(g, "s/HomeScreen.spec.md"), root, g, cfg); v != Pass {
		t.Errorf("Home → GoalDetail is answered by GoalDetail's In: %v %s", v, msg)
	}
	if v, _ := checkNavMatchesSpec("", mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}, root, g, cfg); v != Skip {
		t.Error("a spec that is no screen is skipped")
	}
}

func TestNavChain_reachableAndTheFix(t *testing.T) {
	t.Run("NCGNV-B04: nav-reachable fails a screen no entry route reaches, and is pending with no entry declared", func(t *testing.T) {})
	t.Run("NCGNV-B05: The fixer flags a navigation call whose route names one screen, and leaves a back navigation to the author", func(t *testing.T) {})
	root, g, cfg := navProject(t)
	if v, msg := checkNavReachable("", navNode(g, "s/GoalEditScreen.spec.md"), root, g, cfg); v != Pass {
		t.Errorf("GoalEdit is reached from Home through GoalDetail: %v %s", v, msg)
	}
	if v, msg := checkNavReachable("", navNode(g, "s/LonelyScreen.spec.md"), root, g, cfg); v != Fail || !strings.Contains(msg, "LonelyScreen") {
		t.Errorf("Lonely is reached by nothing: %v %s", v, msg)
	}
	if v, _ := checkNavReachable("", navNode(g, "s/LonelyScreen.spec.md"), root, g, &config.Config{}); v != Pending {
		t.Error("with no entry declared, there is nothing to reach from")
	}
	gates := []config.Gate{{Name: "nav-annotated", Check: "nav-annotated", On: []string{"code"}}}
	FixWithConfig(gates, g.Nodes, root, g, cfg)
	if home := read(t, root, "s/HomeScreen.tsx"); !strings.Contains(home, "navigation.navigate('GoalDetail')} // @navigates: GLDTG") {
		t.Errorf("the call to GoalDetail is flagged with its screen:\n%s", home)
	}
	if detail := read(t, root, "s/GoalDetailScreen.tsx"); strings.Count(detail, "@navigates") != 1 {
		t.Errorf("the waived back navigation is left as it is:\n%s", detail)
	}
}

func TestNavChain_theFixerNeverWritesIntoJSXText(t *testing.T) {
	t.Run("NCGNV-B06: On a line ending in a JSX tag the fixer writes the flag as a block comment right after the call, where it renders nothing, and the flag is read there", func(t *testing.T) {})
	root, g, cfg := navProject(t)
	jsx := "<Text style={link} onPress={() => navigation.navigate('GoalDetail')}>\n" +
		"<Pressable onPress={() => go(navigation.navigate('GoalDetail', { id: f(x) }))}>\n"
	if err := os.WriteFile(filepath.Join(root, "s/HomeScreen.tsx"), []byte(jsx), 0o644); err != nil {
		t.Fatal(err)
	}
	gates := []config.Gate{{Name: "nav-annotated", Check: "nav-annotated", On: []string{"code"}}}
	FixWithConfig(gates, g.Nodes, root, g, cfg)
	home := read(t, root, "s/HomeScreen.tsx")
	for _, want := range []string{
		"navigation.navigate('GoalDetail') /* @navigates: GLDTG */}>",
		"navigation.navigate('GoalDetail', { id: f(x) }) /* @navigates: GLDTG */)}>",
	} {
		if !strings.Contains(home, want) {
			t.Errorf("the flag sits after the call, inside the expression — want %q in:\n%s", want, home)
		}
	}
	if strings.Contains(home, "> //") {
		t.Errorf("nothing is written after the tag:\n%s", home)
	}
	if v, msg := checkNavAnnotated(home, navNode(g, "s/HomeScreen.tsx"), root, g, cfg); v != Pass {
		t.Errorf("the block flags are read: %v %s", v, msg)
	}
}
