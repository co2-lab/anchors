package common

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func cmdWithFlags() (*cobra.Command, *string, *bool) {
	var root string
	var all bool
	cmd := &cobra.Command{Use: "probe", RunE: func(*cobra.Command, []string) error { return nil }}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&all, "all", false, "every node")
	return cmd, &root, &all
}

// The OLD name keeps working: a script that still passes it reaches the new flag.
func TestAliasDeFlag_oldNameReachesTheNewFlag(t *testing.T) {
	cmd, root, all := cmdWithFlags()
	AliasDeFlag(cmd, "root", "raiz")
	AliasDeFlag(cmd, "all", "todos")

	if err := cmd.ParseFlags([]string{"--raiz", "/tmp/project", "--todos"}); err != nil {
		t.Fatalf("the old names were not accepted: %v", err)
	}
	if err := ResolveAliases(cmd, map[string]string{"root": "raiz", "all": "todos"}); err != nil {
		t.Fatal(err)
	}
	if *root != "/tmp/project" {
		t.Errorf("--raiz did not reach --root: %q", *root)
	}
	if !*all {
		t.Error("--todos did not reach --all")
	}
	// The alias is hidden and deprecated: it keeps old scripts alive without being taught.
	f := cmd.Flags().Lookup("raiz")
	if !f.Hidden || !strings.Contains(f.Deprecated, "--root") {
		t.Errorf("the alias should be hidden and point at --root: hidden=%v deprecated=%q", f.Hidden, f.Deprecated)
	}
	if cmd.Flags().Lookup("todos").Value.Type() != "bool" {
		t.Error("the alias of a bool flag must be a bool, or `--todos` alone would demand a value")
	}
}

// When BOTH are passed the new name wins: copying the old one over it would ignore the
// value the person typed with the current name.
func TestResolveAliases_newNameWins(t *testing.T) {
	cmd, root, _ := cmdWithFlags()
	AliasDeFlag(cmd, "root", "raiz")
	if err := cmd.ParseFlags([]string{"--raiz", "/old", "--root", "/new"}); err != nil {
		t.Fatal(err)
	}
	if err := ResolveAliases(cmd, map[string]string{"root": "raiz", "missing": "gone"}); err != nil {
		t.Fatal(err)
	}
	if *root != "/new" {
		t.Errorf("the old name overwrote the new one: %q", *root)
	}
}

// A value the new flag cannot hold is an error that names the OLD flag — the one typed.
func TestResolveAliases_badValueNamesTheOldFlag(t *testing.T) {
	var n int
	cmd := &cobra.Command{Use: "probe"}
	cmd.Flags().IntVar(&n, "limit", 0, "limit")
	cmd.Flags().String("limite", "", "old limit")
	if err := cmd.ParseFlags([]string{"--limite", "many"}); err != nil {
		t.Fatal(err)
	}
	err := ResolveAliases(cmd, map[string]string{"limit": "limite"})
	if err == nil || !strings.Contains(err.Error(), "--limite") {
		t.Errorf("expected an error naming --limite, got %v", err)
	}
}

// Declaring an alias of a flag that does not exist is a programming error, caught at start.
func TestAliasDeFlag_panicsOnUnknownFlag(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || !strings.Contains(r.(string), `"nope"`) {
			t.Errorf("expected a panic naming the flag, got %v", r)
		}
	}()
	cmd, _, _ := cmdWithFlags()
	AliasDeFlag(cmd, "nope", "nao")
}
