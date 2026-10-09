// @anchors
//   code: RNCFG
//   ref: MNTRS

package runs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/testsig"
)

// Settings are the monitor's settings as the project declares them over the defaults.
type Settings struct {
	Timing  Timing
	Runners []Runner
	Reports []string
	Keep    int
	MaxAge  time.Duration
}

// Configure reads the project's `monitor:` block over the defaults: its timing, its runners
// — read before the built-in ones, a runner of a built-in's name replacing it —, the reports
// its suites declare and the ones it adds, and how many ended runs to keep. A duration or a
// pattern that does not read is an error, naming it.
func Configure(cfg *config.Config) (Settings, error) {
	s := Settings{Timing: DefaultTiming(), Runners: DefaultRunners(), Keep: 50, MaxAge: 7 * 24 * time.Hour}
	if cfg == nil {
		return s, nil
	}
	for _, suites := range [][]config.Suite{cfg.Tests, cfg.Mutation} {
		for _, su := range suites {
			if su.JUnit != "" {
				s.Reports = append(s.Reports, su.JUnit)
			}
		}
	}
	m := cfg.Monitor
	if m == nil {
		return s, nil
	}
	for _, d := range []struct {
		name string
		val  string
		to   *time.Duration
	}{{"every", m.Every, &s.Timing.Every}, {"heartbeat", m.Heartbeat, &s.Timing.Heartbeat},
		{"progress", m.Progress, &s.Timing.Progress}, {"stall", m.Stall, &s.Timing.Stall}} {
		if err := setDuration(d.val, d.to); err != nil {
			return s, fmt.Errorf("monitor.%s: %w", d.name, err)
		}
	}
	var own []Runner
	replaced := map[string]bool{}
	for _, r := range m.Runners {
		rn := Runner{Name: r.Name, Kind: r.Kind}
		if rn.Kind == "" {
			rn.Kind = "other"
		}
		for _, p := range []struct {
			name string
			val  string
			to   **regexp.Regexp
		}{{"match", r.Match, &rn.Match}, {"pass", r.Pass, &rn.Pass}, {"fail", r.Fail, &rn.Fail}, {"progress", r.Progress, &rn.Progress}} {
			if p.val == "" {
				continue
			}
			x, err := regexp.Compile(p.val)
			if err != nil {
				return s, fmt.Errorf("monitor.runners[%s].%s: %w", r.Name, p.name, err)
			}
			*p.to = x
		}
		if rn.Match == nil {
			return s, fmt.Errorf("monitor.runners[%s]: match is required", r.Name)
		}
		if err := setDuration(r.Stall, &rn.Stall); err != nil {
			return s, fmt.Errorf("monitor.runners[%s].stall: %w", r.Name, err)
		}
		own = append(own, rn)
		replaced[r.Name] = true
	}
	for _, d := range s.Runners {
		if !replaced[d.Name] {
			own = append(own, d)
		}
	}
	s.Runners = own
	s.Reports = append(s.Reports, m.Reports...)
	if m.Keep > 0 {
		s.Keep = m.Keep
	}
	if m.KeepDays > 0 {
		s.MaxAge = time.Duration(m.KeepDays) * 24 * time.Hour
	}
	return s, nil
}

func setDuration(val string, to *time.Duration) error {
	if val == "" {
		return nil
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	if d <= 0 {
		return fmt.Errorf("%q is not a positive duration", val)
	}
	*to = d
	return nil
}

// ReadReports reads the reports the globs find under root: each one modified after what
// `known` holds for it is parsed for its tests and failures; the others are given with their
// time alone, which is all the monitor compares.
func ReadReports(root string, globs []string, known map[string]time.Time) []Report {
	var out []Report
	seen := map[string]bool{}
	for _, g := range globs {
		pattern := g
		if !filepath.IsAbs(pattern) {
			pattern = filepath.Join(root, filepath.FromSlash(g))
		}
		matches, _ := doublestar.FilepathGlob(pattern)
		for _, p := range matches {
			if seen[p] {
				continue
			}
			seen[p] = true
			st, err := os.Stat(p)
			if err != nil || st.IsDir() {
				continue
			}
			rel, err := filepath.Rel(root, p)
			if err != nil {
				rel = p
			}
			rel = filepath.ToSlash(rel)
			rep := Report{Path: rel, Modified: st.ModTime()}
			if last, ok := known[rel]; !ok || st.ModTime().After(last) {
				if ex, err := testsig.ParseJUnit(p); err == nil {
					for _, c := range ex.Cases {
						rep.Tests++
						if c.Failed {
							rep.Failures++
							if rep.First == "" {
								rep.First = c.Name
							}
						}
					}
				}
			}
			out = append(out, rep)
		}
	}
	return out
}
