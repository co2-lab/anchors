// @anchors
//   code: MNTMP
//   ref: MNCMP

package mapcmd

import (
	"strings"
	"testing"
)

func TestNavLines(t *testing.T) {
	t.Run("MNCMP-B01: The navigation is printed screen by screen, each with where it leads, sorted", func(t *testing.T) {})
	t.Run("MNCMP-B02: One screen shows where it comes from and where it leads; a name that is no screen is refused", func(t *testing.T) {})
	edges := map[string]map[string]bool{
		"s/HomeScreen.spec.md":       {"s/GoalDetailScreen.spec.md": true},
		"s/GoalDetailScreen.spec.md": {"s/GoalEditScreen.spec.md": true, "s/HomeScreen.spec.md": true},
	}
	all, ok := NavLines(edges, "")
	if !ok || strings.Join(all, "\n") != "GoalDetailScreen → GoalEditScreen, HomeScreen\nHomeScreen → GoalDetailScreen" {
		t.Errorf("all:\n%s", strings.Join(all, "\n"))
	}
	one, ok := NavLines(edges, "GoalDetail")
	if !ok || strings.Join(one, "\n") != "← HomeScreen\n→ GoalEditScreen, HomeScreen" {
		t.Errorf("one:\n%s", strings.Join(one, "\n"))
	}
	if _, ok := NavLines(edges, "Nowhere"); ok {
		t.Error("a name that is no screen is refused")
	}
}
