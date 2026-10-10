// @anchors
//   code: CPMPC
//   ref: VRCPT

package mapx

import (
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// captureEdges ties each visual-regression test to what it captures: the main code file of
// the unit whose `{CODE}-VR` it names, and that unit's baseline images.
//
// A capture flow does not import the screen — it reaches it through the running app — so
// nothing in the map linked them, and the evidence of a capture never went stale: the
// screen changed and its VR stayed green over a look it no longer had. A baseline swapped
// without a new capture went unseen the same way.
//
// ONE LEVEL, by decision: the edge reaches the unit's own file and images, and the
// evidence closure does not descend past them. A component used by a screen has states and
// a VR of its own; when it changes, its own capture goes stale, not every screen that uses
// it.
//
// A CONTRACT test (`{CODE}-CT`) is tied the same way to what it validates: the API unit's
// code, its spec — the OpenAPI is compiled from it — and the project's OpenAPI document when
// the map knows it. A handler changed, or a contract regenerated, stales the contract test.
func captureEdges(files []scan.File) []Edge {
	vrRE := regexp.MustCompile(`\b([A-Z0-9]` + config.CodeLengthPattern() + `)-(VR|CT)(?:-(S\d{2}))?\b`)
	var openapiDocs []string
	for _, f := range files {
		if isOpenAPIPath(f.Path) {
			openapiDocs = append(openapiDocs, f.Path)
		}
	}
	sort.Strings(openapiDocs)
	codeByStem := map[string]string{} // unit stem → main code file
	for _, f := range files {
		if f.Kind == string(KindCode) {
			stem := strings.TrimSuffix(f.Path, path.Ext(f.Path))
			if _, ok := codeByStem[stem]; !ok {
				codeByStem[stem] = f.Path
			}
		}
	}
	stemByCode := map[string]string{} // unit code → stem, from the specs' declared identity
	for _, f := range files {
		if f.Kind == string(KindSpec) && f.HeaderCode != "" && strings.HasSuffix(f.Path, ".spec.md") {
			stemByCode[f.HeaderCode] = strings.TrimSuffix(f.Path, ".spec.md")
		}
	}
	images := map[string][]string{} // unit code → its baseline images
	for _, f := range files {
		if !isImagePath(f.Path) {
			continue
		}
		for code, stem := range stemByCode {
			if strings.HasPrefix(f.Path, stem+"."+code+"-VR") {
				images[code] = append(images[code], f.Path)
			}
		}
	}
	var edges []Edge
	for _, f := range files {
		if f.Kind != string(KindTest) || isImagePath(f.Path) {
			continue
		}
		// What the test compares against: each unit code, VR or CT, and the states it names
		// ("" = every state). A code only a comment mentions captures nothing; a test of one
		// state is no test of the others (reported from MIF: a Home baseline staled the
		// Alerts flow, whose comment named HOMEH-VR, and every test of another Home state).
		codes := map[string]string{}
		states := map[string]map[string]bool{}
		note := func(code, kind, state string) {
			codes[code] = kind
			if _, seen := states[code]; !seen {
				states[code] = map[string]bool{}
			}
			states[code][state] = true
		}
		for _, m := range vrRE.FindAllStringSubmatch(f.Path, -1) {
			note(m[1], m[2], m[3])
		}
		for _, c := range f.CaptureCodes {
			for _, kind := range []string{"VR", "CT"} {
				if i := strings.Index(c, "-"+kind); i > 0 {
					note(c[:i], kind, strings.TrimPrefix(c[i+3:], "-"))
				}
			}
		}
		for _, code := range sortedKeysOf(codes) {
			stem, ok := stemByCode[code]
			if !ok {
				continue
			}
			var targets []string
			if src, ok := codeByStem[stem]; ok {
				targets = append(targets, src)
			}
			if codes[code] == "CT" {
				targets = append(targets, stem+".spec.md")
				targets = append(targets, openapiDocs...)
			}
			imgs := imagesOfStates(images[code], stem+"."+code+"-VR", states[code])
			sort.Strings(imgs)
			targets = append(targets, imgs...)
			for _, t := range targets {
				if t != f.Path {
					edges = append(edges, Edge{From: f.Path, To: t, Type: EdgeCaptures, Origin: OriginConvention})
				}
			}
		}
	}
	return edges
}

// imagesOfStates are the baseline images of the states named — all of them when a state
// is left unnamed ("") —: `Home.HOMEH-VR-S02-loaded.png` is the image of S02.
func imagesOfStates(imgs []string, prefix string, states map[string]bool) []string {
	if states[""] {
		return append([]string(nil), imgs...)
	}
	var out []string
	for _, img := range imgs {
		rest := strings.TrimPrefix(img, prefix)
		for st := range states {
			if strings.HasPrefix(rest, "-"+st) {
				after := rest[len(st)+1:]
				if after == "" || after[0] == '-' || after[0] == '.' {
					out = append(out, img)
					break
				}
			}
		}
	}
	return out
}

func isImagePath(p string) bool {
	switch strings.ToLower(path.Ext(p)) {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".svg":
		return true
	}
	return false
}

// isOpenAPIPath says whether a file is an OpenAPI document by its name: `openapi.yaml`,
// `openapi.json`, `api.openapi.yml`.
func isOpenAPIPath(p string) bool {
	b := strings.ToLower(path.Base(p))
	return strings.Contains(b, "openapi") && (strings.HasSuffix(b, ".yaml") || strings.HasSuffix(b, ".yml") || strings.HasSuffix(b, ".json"))
}

func sortedKeysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// composesEdges ties a spec to the code of each unit its Parts Used section names: the code
// file whose name, without extension, is the name written (`BottomSheet` →
// `…/BottomSheet.tsx`). A name no code file carries ties nothing; a name two carry, the first
// by path.
func composesEdges(files []scan.File) []Edge {
	byName := map[string][]string{}
	for _, f := range files {
		if f.Kind == string(KindCode) {
			stem := strings.TrimSuffix(path.Base(f.Path), path.Ext(f.Path))
			byName[stem] = append(byName[stem], f.Path)
		}
	}
	var edges []Edge
	for _, f := range files {
		if f.Kind != string(KindSpec) {
			continue
		}
		seen := map[string]bool{}
		for _, name := range f.Composes {
			cands := byName[name]
			if len(cands) == 0 {
				continue
			}
			sort.Strings(cands)
			if seen[cands[0]] {
				continue
			}
			seen[cands[0]] = true
			edges = append(edges, Edge{From: f.Path, To: cands[0], Type: EdgeComposes, Origin: OriginDeclared})
		}
	}
	return edges
}
