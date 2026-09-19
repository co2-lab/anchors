package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// RelTo converte caminho em relativo à raiz com barras normais (/).
func RelTo(root, arg string) string {
	if !filepath.IsAbs(arg) {
		if _, err := os.Stat(filepath.Join(root, arg)); err == nil {
			return filepath.ToSlash(filepath.Clean(arg))
		}
	}
	abs, err := filepath.Abs(arg)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	return filepath.ToSlash(rel)
}

// NodeExists confere se um nó com aquele id existe no mapa.
func NodeExists(m *mapx.Graph, path string) bool {
	if m == nil {
		return false
	}
	for _, n := range m.Nodes {
		if n.ID == path {
			return true
		}
	}
	return false
}

// RelSlug reduz um caminho a algo usável em ID de task.
func RelSlug(p string) string {
	return strings.TrimSuffix(p, filepath.Ext(p))
}
