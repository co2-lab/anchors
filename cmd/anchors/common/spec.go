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
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
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
	return out, nil
}
