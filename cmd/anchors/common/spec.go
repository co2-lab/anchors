// @anchors
//   code: SPCMS
//   ref: UNCDN

package common

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testsig"
)

func SpecHeaderCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)` + config.HeaderLinePrefix + `code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

// CodeDoHeaderSpec extrai o código de identidade do header de uma spec.
func CodeDoHeaderSpec(content string) string {
	if m := SpecHeaderCodeRE().FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// CodeOfUnit descobre o código da unidade consultando o mapa.
func CodeOfUnit(root, unit string) string {
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		return ""
	}
	unit = filepath.ToSlash(unit)
	stem, _ := mapx.StemOfAnchor(unit)
	for _, n := range g.Nodes {
		if n.ID == unit && n.Code != "" {
			return n.Code
		}
	}
	for _, n := range g.Nodes {
		if s, _ := mapx.StemOfAnchor(n.ID); s == stem && n.Code != "" {
			return n.Code
		}
	}
	return ""
}

// CodesInFileOfUnit extrai os códigos que o arquivo DECLARA, descartando os que ele apenas CITA.
func CodesInFileOfUnit(path, unit string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := string(b)
	seen := map[string]bool{}
	var out []string
	for _, c := range testsig.CodesInCase(data) {
		if unit != "" && !strings.HasPrefix(c, unit+"-") {
			continue
		}
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	// A data state written bare — `DS-seen-no`, as a spec's data-state table names it — is
	// the unit's own: the test names it with the prefix (`ARSCA-DS-seen-no: …`), and without
	// it here the passing test was dropped at ingestion, and the data state never proven
	// (reported from MIF: 129 data states, every one with a passing test).
	if unit != "" {
		for _, m := range bareDataStateRE.FindAllStringSubmatch(data, -1) {
			if c := unit + "-DS-" + m[1]; !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return out, nil
}

// bareDataStateRE is a data state with no unit prefix: `DS-` not glued to a code before it.
var bareDataStateRE = regexp.MustCompile(`(?:^|[^A-Za-z0-9_-])DS-([A-Za-z0-9][A-Za-z0-9-]*)`)
