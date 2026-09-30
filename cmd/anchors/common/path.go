package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
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

// FilesAnnotation marks a command whose arguments are a list of files read by FileArgs.
const FilesAnnotation = "anchors.files"

// FileArgs reads the arguments of a command that takes files: each argument is a file, or
// several separated by commas — the same list a `--changed` flag takes. Blanks are
// dropped, a file named twice is kept once, and the order is kept.
func FileArgs(args []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, a := range args {
		for _, f := range strings.Split(a, ",") {
			if f = strings.TrimSpace(f); f != "" && !seen[f] {
				seen[f] = true
				out = append(out, f)
			}
		}
	}
	return out
}

// TakesFiles makes the command read its arguments through FileArgs, and marks it: a list
// of files is taken the same way by every command, and the root's test finds a command
// that takes several files without it.
func TakesFiles(cmd *cobra.Command) *cobra.Command {
	if run := cmd.RunE; run != nil {
		cmd.RunE = func(c *cobra.Command, args []string) error { return run(c, FileArgs(args)) }
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[FilesAnnotation] = "true"
	return cmd
}
