package quality

import (
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var qualityCommandNames = []string{"check", "coverage", "doctor", "keep-evidence", "mutation", "report", "stale", "stamp", "status", "test", "touch", "verify"}

func TestRegisterAddsExactlyTheQualityCommands(t *testing.T) {
	t.Run("QLCMQ-B01: The quality domain registers exactly its twelve commands", func(t *testing.T) {})
	t.Run("QLCMQ-I01: No two quality commands share a name", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)
	var got []string
	seen := map[string]bool{}
	for _, c := range root.Commands() {
		if seen[c.Name()] {
			t.Errorf("command %q registered twice", c.Name())
		}
		seen[c.Name()] = true
		got = append(got, c.Name())
	}
	sort.Strings(got)
	if strings.Join(got, ",") != strings.Join(qualityCommandNames, ",") {
		t.Errorf("registered %v, want %v", got, qualityCommandNames)
	}
}

func TestRegisterMakesEachCommandReachable(t *testing.T) {
	t.Run("QLCMQ-B02: Each quality command is reachable by its name", func(t *testing.T) {})
	root := &cobra.Command{Use: "anchors"}
	Register(root)
	for _, name := range qualityCommandNames {
		c, _, err := root.Find([]string{name})
		if err != nil || c == nil || c.Name() != name {
			t.Errorf("%q is not reachable from the root: cmd=%v err=%v", name, c, err)
		}
	}
}

func TestRegisterPrintsNothing(t *testing.T) {
	t.Run("QLCMQ-X01: Registering the quality commands prints nothing", func(t *testing.T) {})
	out := captureStdout(t, func() { Register(&cobra.Command{Use: "anchors"}) })
	if out != "" {
		t.Errorf("registering printed:\n%s", out)
	}
}
