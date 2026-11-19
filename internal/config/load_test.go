// @anchors
//   code: LDCFT
//   ref: LDCFS

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func loadWorkload(t *testing.T, body string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\nload:\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestWorkload_loads(t *testing.T) {
	t.Run("LDCFS-B01: load: declares levels over the defaults, tasks by scope in the order written, wait limits and CI", func(t *testing.T) {
		c, err := loadWorkload(t, `  levels:
    heavy: { cpu: 6, memory: 8GB }
    huge: { cpu: 8, memory: 12GB }
  wait_limit: 10m
  ci: true
  tasks:
    mutation:
      default: medium
      wait_limit: 20m
      scopes:
        "packages/backend/**": huge
        "packages/**": heavy
    test:
      default: light
`)
		if err != nil {
			t.Fatal(err)
		}
		w := c.Load
		if cpu, mem, ok := w.LevelCost("heavy"); !ok || cpu != 6 || mem != 8<<30 {
			t.Errorf("heavy: %v %v %v", cpu, mem, ok)
		}
		if cpu, mem, ok := w.LevelCost("light"); !ok || cpu != 1 || mem != 512<<20 {
			t.Errorf("light keeps its default: %v %v %v", cpu, mem, ok)
		}
		for _, c := range []struct {
			task  string
			files []string
			want  string
		}{
			{"mutation", []string{"packages/backend/a.ts"}, "huge"},
			{"mutation", []string{"apps/x.ts", "packages/ui/b.ts"}, "heavy"},
			{"mutation", []string{"apps/x.ts"}, "medium"},
			{"test", nil, "light"},
			{"gate", []string{"a.ts"}, "medium"},
		} {
			if got := w.LevelFor(c.task, c.files); got != c.want {
				t.Errorf("%s %v: %s, want %s", c.task, c.files, got, c.want)
			}
		}
		if w.WaitLimitFor("mutation") != 20*time.Minute || w.WaitLimitFor("test") != 10*time.Minute || !w.CI {
			t.Errorf("wait limits: %v %v", w.WaitLimitFor("mutation"), w.WaitLimitFor("test"))
		}
		if (&Workload{}).WaitLimitFor("test") != DefaultWaitLimit {
			t.Errorf("default wait limit")
		}
		if b, err := ParseBytes("1.5GB"); err != nil || b != 3<<29 {
			t.Errorf("1.5GB: %v %v", b, err)
		}
	})
	t.Run("LDCFS-E01: A level with no CPU or an unreadable size, a task the queue does not know, an undeclared level, a bad pattern or duration refuse the configuration", func(t *testing.T) {
		for body, says := range map[string]string{
			"  levels: { big: { cpu: 0, memory: 1GB } }\n":      "cpu is a positive number",
			"  levels: { big: { cpu: 2, memory: lots } }\n":     "a size is a number and a unit",
			"  tasks: { lint: { default: light } }\n":           "the queue admits test, mutation, gate",
			"  tasks: { test: { default: enormous } }\n":        `level "enormous" is not declared`,
			"  tasks: { test: { scopes: { \"[\": light } } }\n": "is not a file pattern",
			"  wait_limit: soon\n":                              `"soon" is not a duration`,
			"  tasks: { test: { wait_limit: -1m } }\n":          "is not a duration",
			"  tasks: { test: { level: light } }\n":             "level",
		} {
			if _, err := loadWorkload(t, body); err == nil || !strings.Contains(err.Error(), says) {
				t.Errorf("%s: %v, want %q", body, err, says)
			}
		}
	})
}
