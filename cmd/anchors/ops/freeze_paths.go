// @anchors
//   code: FZPTH
//   ref: FRZEX

package ops

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a folder frozen: a freeze's scope by path, in anchors.yaml (FRUNF-W04) ---
//
// The freeze's record stays the doctrine's entry; `frozen_paths` says where it applies by
// glob, and the map marks every file a glob matches. The file is the team's, written by hand
// and commented: an entry is added as text at the end of the block, and removed by the lines
// the YAML parser places it on — the rest of the file is never rewritten.

const frozenPathsKey = "frozen_paths"

// checkFrozenGlobs refuses a glob that is no glob, or that matches no file of the map: a
// folder frozen by a typo freezes nothing, and freeze-consistent would say so after the fact.
func checkFrozenGlobs(g *mapx.Graph, globs []string) error {
	for _, glob := range globs {
		if !doublestar.ValidatePattern(glob) {
			return fmt.Errorf("--path %q is not a valid glob", glob)
		}
		found := false
		for _, n := range g.Nodes {
			if ok, _ := doublestar.Match(glob, n.ID); ok {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("--path %q matches no file of the project", glob)
		}
	}
	return nil
}

// addFrozenPaths writes an entry of the freeze with the globs at the end of the block,
// creating the block at the end of the file when it has none.
func addFrozenPaths(cfgPath, value string, globs []string) error {
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	item := []string{"  - freeze: " + strconv.Quote(value), "    paths:"}
	for _, glob := range globs {
		item = append(item, "      - "+strconv.Quote(glob))
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, frozenPathsKey+":") {
			start = i
			break
		}
	}
	if start < 0 {
		lines = append(lines, "", frozenPathsKey+":")
		lines = append(lines, item...)
		return os.WriteFile(cfgPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
	}
	if strings.TrimSpace(strings.TrimPrefix(lines[start], frozenPathsKey+":")) != "" {
		return fmt.Errorf("%s is written inline in %s — add the paths by hand", frozenPathsKey, cfgPath)
	}
	end := blockEnd(lines, start)
	at := end
	for at > start+1 && strings.TrimSpace(lines[at-1]) == "" {
		at--
	}
	lines = append(lines[:at], append(item, lines[at:]...)...)
	return os.WriteFile(cfgPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// blockEnd is the index after the last line of the top-level block that starts at start: the
// next line written at the first column, a comment included.
func blockEnd(lines []string, start int) int {
	for i := start + 1; i < len(lines); i++ {
		if l := lines[i]; l != "" && l[0] != ' ' && l[0] != '\t' {
			return i
		}
	}
	return len(lines)
}

// removeFrozenPaths removes, from the block, the entries of a freeze (value) — or, with globs
// and no value, those globs from whichever entry holds them, and an entry left with none —,
// and the block once empty. It says how many globs it removed.
func removeFrozenPaths(cfgPath, value string, globs []string) (int, error) {
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return 0, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return 0, err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return 0, nil
	}
	top := doc.Content[0]
	var seq *yaml.Node
	keyLine := 0
	for i := 0; i+1 < len(top.Content); i += 2 {
		if top.Content[i].Value == frozenPathsKey {
			seq, keyLine = top.Content[i+1], top.Content[i].Line
		}
	}
	if seq == nil || seq.Kind != yaml.SequenceNode {
		return 0, nil
	}
	lines := strings.Split(string(b), "\n")
	blockStop := blockEnd(lines, keyLine-1) // 0-based, exclusive
	// The blank lines that close the block separate it from what follows, and stay.
	for blockStop > keyLine && strings.TrimSpace(lines[blockStop-1]) == "" {
		blockStop--
	}
	drop := map[int]bool{} // 0-based lines to remove
	removed, kept := 0, 0
	want := map[string]bool{}
	for _, g := range globs {
		want[g] = true
	}
	for idx, item := range seq.Content {
		itemStop := blockStop
		if idx+1 < len(seq.Content) {
			itemStop = seq.Content[idx+1].Line - 1
		}
		freeze, paths := frozenItem(item)
		if value != "" && freeze == value {
			for l := item.Line - 1; l < itemStop; l++ {
				drop[l] = true
			}
			removed += len(paths.Content)
			continue
		}
		left := 0
		for _, p := range paths.Content {
			if value == "" && want[p.Value] {
				if paths.Style == yaml.FlowStyle {
					return 0, fmt.Errorf("the paths of %s are written inline in %s — remove %q by hand", freeze, cfgPath, p.Value)
				}
				drop[p.Line-1] = true
				removed++
				continue
			}
			left++
		}
		if left == 0 && len(paths.Content) > 0 {
			for l := item.Line - 1; l < itemStop; l++ {
				drop[l] = true
			}
		}
		kept += left
	}
	if removed == 0 {
		return 0, nil
	}
	if kept == 0 {
		for l := keyLine - 1; l < blockStop; l++ {
			drop[l] = true
		}
	}
	var out []string
	for i, l := range lines {
		if !drop[i] {
			out = append(out, l)
		}
	}
	text := strings.Join(out, "\n")
	if kept == 0 {
		text = strings.TrimRight(text, "\n") + "\n"
	}
	return removed, os.WriteFile(cfgPath, []byte(text), 0o644)
}

// frozenItem is an entry's freeze and its paths node.
func frozenItem(item *yaml.Node) (string, *yaml.Node) {
	freeze, paths := "", &yaml.Node{}
	for i := 0; i+1 < len(item.Content); i += 2 {
		switch item.Content[i].Value {
		case "freeze":
			freeze = strings.TrimSpace(item.Content[i+1].Value)
		case "paths":
			paths = item.Content[i+1]
		}
	}
	return freeze, paths
}

// frozenGlobsByValue are the globs anchors.yaml freezes, by freeze, each list sorted.
func frozenGlobsByValue(cfg *config.Config) map[string][]string {
	out := map[string][]string{}
	if cfg == nil {
		return out
	}
	for _, f := range cfg.FrozenPaths {
		v := strings.TrimSpace(f.Freeze)
		out[v] = append(out[v], f.Paths...)
	}
	for v := range out {
		sort.Strings(out[v])
	}
	return out
}
