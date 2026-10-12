// @anchors
//   code: RGTSR
//   ref: FLRGF

package flow

import (
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func childNames(c *cobra.Command) []string {
	var names []string
	for _, sub := range c.Commands() {
		names = append(names, sub.Name())
	}
	sort.Strings(names)
	return names
}

func TestRegister_attachesTheFlowCommandsOnce(t *testing.T) {
	t.Run("FLRGF-B01: The root holds exactly the nineteen flow commands after registration", func(t *testing.T) {})
	t.Run("FLRGF-I01: No flow command is attached twice", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)

	// Commands() already sorts and a duplicate would show up as a repeated name.
	got := childNames(root)
	want := []string{"backfill-labels", "decided", "deliver", "discard", "done", "drop", "escalate",
		"merge-progress", "monitor", "next", "pr-body", "progress", "queue", "reclaim", "report-bug", "task-status", "unblock", "watch", "work"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("root children:\n got  %v\n want %v", got, want)
	}
	seen := map[string]int{}
	for _, n := range got {
		seen[n]++
		if seen[n] > 1 {
			t.Errorf("%q is attached more than once", n)
		}
	}
}

func TestRegister_watchControlsLiveUnderWatch(t *testing.T) {
	t.Run("FLRGF-B02: The watcher's controls live under watch", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)
	watch, _, err := root.Find([]string{"watch"})
	if err != nil || watch.Name() != "watch" {
		t.Fatalf("watch not found: %v", err)
	}
	want := "logs,pause,resume,run,start,status,stop"
	if got := strings.Join(childNames(watch), ","); got != want {
		t.Errorf("watch children: got %s, want %s", got, want)
	}
	for _, n := range strings.Split(want, ",") {
		for _, top := range childNames(root) {
			if top == n {
				t.Errorf("%q must not be a root command", n)
			}
		}
	}
}
