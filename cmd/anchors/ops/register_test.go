// @anchors
//   ref: OPRGP

package ops

import (
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Register hands the root the fifteen operation commands, each exactly once.
func TestRegisterAddsEachOpsCommandOnce(t *testing.T) {
	t.Run("OPRGP-B01: The root receives the fifteen operation commands", func(t *testing.T) {})
	t.Run("OPRGP-I01: Each operation command is registered exactly once", func(t *testing.T) {})
	t.Run("OPRGP-X01: The registration adds commands and nothing else", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)
	var names []string
	for _, c := range root.Commands() {
		names = append(names, c.Name())
	}
	sort.Strings(names)
	want := "board,changelog,code,commit-msg,docs,freeze,generated-paths,init,install-hooks,migrate,new,settings,suggest,synthesize,thaw"
	if got := strings.Join(names, ","); got != want {
		t.Errorf("registered = %s\nwant       = %s", got, want)
	}
	if root.HasAvailableFlags() || root.PersistentPreRunE != nil || root.RunE != nil {
		t.Error("the registration changed the root beyond adding commands")
	}
}
