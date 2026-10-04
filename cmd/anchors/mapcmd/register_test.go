// @anchors
//   code: RGTSA
//   ref: MPRGM

package mapcmd

import (
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// The map domain hangs exactly its eight commands on the root, and each one is reachable
// by its name.
func TestRegister_addsTheMapDomainCommands(t *testing.T) {
	t.Run("MPRGM-B01: Registering the map domain makes each of its commands reachable from the root", func(t *testing.T) {})
	t.Run("MPRGM-X01: Registering the map domain adds no command outside it", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)

	want := []string{"failures", "flow", "impact", "ingest", "judge", "map", "recode", "renumber", "review"}
	for _, name := range want {
		c, _, err := root.Find([]string{name})
		if err != nil || c == nil || c.Name() != name {
			t.Errorf("%s is not reachable from the root: %v", name, err)
		}
	}
	var got []string
	for _, c := range root.Commands() {
		got = append(got, c.Name())
	}
	sort.Strings(got)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("registered %v, want exactly %v", got, want)
	}
}
