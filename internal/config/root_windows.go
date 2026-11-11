// @anchors
//   code: RTWND
//   ref: PRRPR

package config

import (
	"os"
	"path/filepath"
	"strings"
)

// diskPath is the path with each part as its folder lists it: Windows ignores case, and a
// working directory typed in another case (`cd c:\users\x\mif`) is kept as typed. A part no
// folder lists in another case — a short `RUNNER~1` name — stays as given.
func diskPath(path string) (string, bool) {
	vol := filepath.VolumeName(path)
	rest := strings.TrimPrefix(path[len(vol):], string(filepath.Separator))
	cur := vol + string(filepath.Separator)
	for _, part := range strings.Split(rest, string(filepath.Separator)) {
		if part == "" {
			continue
		}
		cur = filepath.Join(cur, diskName(cur, part))
	}
	return cur, true
}

// diskName is the name the folder lists for `part`: the part itself when listed so, else the
// one entry that differs only in case; the part as given when there is none, or several.
func diskName(dir, part string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return part
	}
	match := ""
	for _, e := range entries {
		if e.Name() == part {
			return part
		}
		if strings.EqualFold(e.Name(), part) {
			if match != "" {
				return part
			}
			match = e.Name()
		}
	}
	if match == "" {
		return part
	}
	return match
}
