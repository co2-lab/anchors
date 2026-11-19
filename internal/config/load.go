// @anchors
//   code: LDCFV
//   ref: LDCFS

package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// Workload is what each run of the project costs, so heavy runs wait for one another instead of
// fighting over one machine (reported from a project: two worktrees' suites at a load of 30
// to 40 on 10 cores, and timing tests breaking). The costs are the project's, in real
// resources, the same on every machine; each machine's capacity is its own.
type Workload struct {
	// Levels are the costs a task may be declared at. Unset levels keep their defaults.
	Levels map[string]LoadCost `yaml:"levels,omitempty"`
	// Tasks are the costs of each kind of run: `test`, `mutation`, `gate`.
	Tasks map[string]LoadTask `yaml:"tasks,omitempty"`
	// WaitLimit is how long a run may wait before no new run is admitted until it fits.
	// Empty: 5m.
	WaitLimit string `yaml:"wait_limit,omitempty"`
	// CI admits the runs through the queue in CI too; by default CI runs everything at once.
	CI bool `yaml:"ci,omitempty"`
}

// LoadCost is what a run holds while it runs: CPUs and memory (`512MB`, `6GB`).
type LoadCost struct {
	CPU    float64 `yaml:"cpu"`
	Memory string  `yaml:"memory"`
}

// LoadTask is the cost of one kind of run: the level of the first scope a run's files match,
// else its default, and how long it may wait.
type LoadTask struct {
	Default   string     `yaml:"default,omitempty"`
	WaitLimit string     `yaml:"wait_limit,omitempty"`
	Scopes    LoadScopes `yaml:"scopes,omitempty"`
}

// LoadScope says that runs over files a pattern matches cost a level.
type LoadScope struct {
	Pattern string `yaml:"-"`
	Level   string `yaml:"-"`
}

// LoadScopes keep the order they are written in: the first pattern a run's files match wins.
type LoadScopes []LoadScope

// UnmarshalYAML reads `scopes:` as a map of pattern → level, in the order it is written.
func (s *LoadScopes) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: scopes: a map of file pattern → level", node.Line)
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		*s = append(*s, LoadScope{Pattern: node.Content[i].Value, Level: node.Content[i+1].Value})
	}
	return nil
}

// LoadTasks are the kinds of run the queue admits.
var LoadTasks = []string{"test", "mutation", "gate"}

// DefaultLoadLevels are the costs of the levels a project does not set: a guess until runs
// are measured.
var DefaultLoadLevels = map[string]LoadCost{
	"light":  {CPU: 1, Memory: "512MB"},
	"medium": {CPU: 2, Memory: "2GB"},
	"heavy":  {CPU: 4, Memory: "6GB"},
}

// DefaultWaitLimit is how long a run waits before the queue admits no new run until it fits.
const DefaultWaitLimit = 5 * time.Minute

// LevelCost is a level's CPUs and bytes, the project's or the default; false for a level
// neither declares.
func (l *Workload) LevelCost(level string) (cpu float64, mem int64, ok bool) {
	c, found := DefaultLoadLevels[level]
	if l != nil {
		if own, set := l.Levels[level]; set {
			c, found = own, true
		}
	}
	if !found {
		return 0, 0, false
	}
	b, err := ParseBytes(c.Memory)
	if err != nil {
		return 0, 0, false
	}
	return c.CPU, b, true
}

// LevelFor is the level a run of a task over files costs: the first scope one of the files
// matches, else the task's default, else medium.
func (l *Workload) LevelFor(task string, files []string) string {
	t := l.Tasks[task]
	for _, s := range t.Scopes {
		for _, f := range files {
			if ok, _ := doublestar.Match(s.Pattern, f); ok {
				return s.Level
			}
		}
	}
	if t.Default != "" {
		return t.Default
	}
	return "medium"
}

// WaitLimitFor is how long a run of the task may wait: its own limit, else the project's,
// else DefaultWaitLimit.
func (l *Workload) WaitLimitFor(task string) time.Duration {
	for _, v := range []string{l.Tasks[task].WaitLimit, l.WaitLimit} {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return DefaultWaitLimit
}

// ParseBytes reads a size: a number with B, KB, MB, GB or TB (powers of 1024).
func ParseBytes(s string) (int64, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	units := []struct {
		suffix string
		mult   int64
	}{{"TB", 1 << 40}, {"GB", 1 << 30}, {"MB", 1 << 20}, {"KB", 1 << 10}, {"B", 1}}
	for _, u := range units {
		if strings.HasSuffix(t, u.suffix) {
			n, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(t, u.suffix)), 64)
			if err != nil || n < 0 {
				return 0, fmt.Errorf("%q: a size is a number and a unit (512MB, 6GB)", s)
			}
			return int64(n * float64(u.mult)), nil
		}
	}
	return 0, fmt.Errorf("%q: a size is a number and a unit (512MB, 6GB)", s)
}

// validateLoad refuses a level with no cost or an unreadable size, a task the queue does not
// know, a level no one declares, a pattern that is no glob, and a duration that is not one.
func validateLoad(l *Workload) error {
	if l == nil {
		return nil
	}
	names := make([]string, 0, len(l.Levels))
	for n := range l.Levels {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		c := l.Levels[n]
		if c.CPU <= 0 {
			return fmt.Errorf("load.levels.%s: cpu is a positive number", n)
		}
		if _, err := ParseBytes(c.Memory); err != nil {
			return fmt.Errorf("load.levels.%s: memory %w", n, err)
		}
	}
	if err := validDuration("load.wait_limit", l.WaitLimit); err != nil {
		return err
	}
	tasks := make([]string, 0, len(l.Tasks))
	for t := range l.Tasks {
		tasks = append(tasks, t)
	}
	sort.Strings(tasks)
	for _, name := range tasks {
		known := false
		for _, k := range LoadTasks {
			known = known || k == name
		}
		if !known {
			return fmt.Errorf("load.tasks.%s: the queue admits %s", name, strings.Join(LoadTasks, ", "))
		}
		t := l.Tasks[name]
		levels := []string{t.Default}
		for _, s := range t.Scopes {
			if strings.TrimSpace(s.Pattern) == "" || !doublestar.ValidatePattern(s.Pattern) {
				return fmt.Errorf("load.tasks.%s.scopes: %q is not a file pattern", name, s.Pattern)
			}
			levels = append(levels, s.Level)
		}
		for _, lv := range levels {
			if lv == "" {
				continue
			}
			if _, _, ok := l.LevelCost(lv); !ok {
				return fmt.Errorf("load.tasks.%s: level %q is not declared (light, medium, heavy, or load.levels)", name, lv)
			}
		}
		if err := validDuration("load.tasks."+name+".wait_limit", t.WaitLimit); err != nil {
			return err
		}
	}
	return nil
}

func validDuration(field, v string) error {
	if v == "" {
		return nil
	}
	if d, err := time.ParseDuration(v); err != nil || d <= 0 {
		return fmt.Errorf("%s: %q is not a duration (30s, 5m, 1h)", field, v)
	}
	return nil
}
