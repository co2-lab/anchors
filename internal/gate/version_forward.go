// @anchors
//   code: VRFRV
//   ref: VRFRA

package gate

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/version"
)

// version-forward: a change to a product's version moves it forward. The version the
// sources hold now is read against the version they held at HEAD: it may stay — a branch
// keeps one version over many commits —, and when it changes it must be readable and
// above the one before; two or more sources must hold the same one. Whether the version
// was already released elsewhere is not a question of the commit, which cannot see the
// other branches: the pipeline's stage asks it.

func checkVersionForward(_ config.Gate, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	names := make([]string, 0, len(cfg.Versioning))
	for n, p := range cfg.Versioning {
		if len(p.Source) > 0 {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return Skip, i18n.T("gate.version_forward.no_source")
	}
	var fails []string
	for _, name := range names {
		if msg := versionMovedForward(name, cfg.Versioning[name], root); msg != "" {
			fails = append(fails, msg)
		}
	}
	if len(fails) > 0 {
		return Fail, strings.Join(fails, "; ")
	}
	return Pass, ""
}

// versionMovedForward confronts one product's version now with the one at HEAD; empty when
// it holds.
func versionMovedForward(name string, p *config.Product, root string) string {
	now, err := p.ReadSources(func(f string) ([]byte, error) {
		return os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
	})
	if err != nil {
		return i18n.T("gate.version_forward.unreadable", name, err.Error())
	}
	cur, err := config.Agreed(now)
	if err != nil {
		return i18n.T("gate.version_forward.unreadable", name, err.Error())
	}
	// At HEAD: a source that did not exist, or held no version, is a version born now.
	before, err := p.ReadSources(func(f string) ([]byte, error) {
		s, ok := gitmeta.AtHead(root, f)
		if !ok {
			return nil, os.ErrNotExist
		}
		return []byte(s), nil
	})
	if err != nil {
		return ""
	}
	prev, err := config.Agreed(before)
	if err != nil || prev == cur {
		return ""
	}
	s, _ := p.SchemeOf()
	v, err := s.Parse(cur)
	if err != nil {
		return i18n.T("gate.version_forward.unreadable", name, err.Error())
	}
	pv, err := s.Parse(prev)
	if err != nil {
		return ""
	}
	if version.Compare(v, pv) < 0 {
		return i18n.T("gate.version_forward.backwards", name, prev, cur)
	}
	if version.Compare(v, pv) == 0 {
		return i18n.T("gate.version_forward.sideways", name, prev, cur)
	}
	return ""
}
