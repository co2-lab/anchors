package main

import (
	"bytes"
	"github.com/co2-lab/anchors/cmd/anchors/common"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/cmd/anchors/governance"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/spf13/cobra"
)

// registeredNames collects the name of every command under the root.
func registeredNames() map[string]bool {
	registered := map[string]bool{}
	var collect func(*cobra.Command)
	collect = func(c *cobra.Command) {
		registered[c.Name()] = true
		for _, f := range c.Commands() {
			collect(f)
		}
	}
	collect(newRootCmd())
	return registered
}

func TestGuideCitesOnlyCommandsThatExist(t *testing.T) {
	t.Run("CLRTC-I01: Every command the work guide and the pipelines teach is registered", func(t *testing.T) {})
	registered := registeredNames()
	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	seen := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(governance.WorkGuide, -1) {
		name := m[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		if name == "guide" || registered[name] {
			continue
		}
		t.Errorf("the guide says to run `anchors %s`, and that command does not exist — whoever "+
			"follows it gets `unknown command` and stops", name)
	}
	if len(seen) == 0 {
		t.Fatal("no `anchors <command>` in the guide — the regex broke and the test would pass empty")
	}
}

func TestPipelineTeachesOnlyCommandsThatExist(t *testing.T) {
	t.Run("CLRTC-I01: Every command the work guide and the pipelines teach is registered", func(t *testing.T) {})
	registered := registeredNames()
	dir := filepath.Join("..", "..", "internal", "initx", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read the workflows: %v", err)
	}
	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	total := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("%s: %v", e.Name(), err)
		}
		seen := map[string]bool{}
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			name := m[1]
			if seen[name] {
				continue
			}
			seen[name] = true
			total++
			if registered[name] {
				continue
			}
			t.Errorf("%s cites `anchors %s`, and that command does not exist -- whoever follows "+
				"the instruction gets `unknown command`", e.Name(), name)
		}
	}
	if total == 0 {
		t.Fatal("no `anchors <command>` in the workflows -- the regex broke and the test would pass empty")
	}
}

// --- the freeze and the language, applied before every command ---

// keepLang restores the process language after a test that changes it.
func keepLang(t *testing.T) {
	t.Helper()
	t.Setenv("ANCHORS_TELEMETRY", "off") // no notice, no event from a test
	prev := i18n.Current()
	t.Cleanup(func() { _ = i18n.Set(prev) })
}

func project(t *testing.T, cfg string) string {
	t.Helper()
	dir := t.TempDir()
	if cfg != "" {
		if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(cfg), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

const frozenCfg = "enabled: false\nfreeze_reason: \"rotate the key\"\nversion: 1\n"

// find returns the command at path under a fresh root, with --root set when it has one.
func find(t *testing.T, root string, path ...string) *cobra.Command {
	t.Helper()
	r := newRootCmd()
	r.InitDefaultHelpCmd() // cobra adds help and completion only when it executes
	r.InitDefaultCompletionCmd()
	c, _, err := r.Find(path)
	if err != nil || c == nil {
		t.Fatalf("find %v: %v", path, err)
	}
	if f := c.Flags().Lookup("root"); f != nil {
		if err := c.Flags().Set("root", root); err != nil {
			t.Fatal(err)
		}
	}
	return c
}

// A frozen project refuses every command that produces project state, and names the
// command and the reason.
func TestFrozenProjectRefusesTheCommand(t *testing.T) {
	t.Run("CLRTC-B01: A frozen project refuses the command, naming it and the reason", func(t *testing.T) {})
	t.Run("CLRTC-B05: The root flag decides which project is read", func(t *testing.T) {})
	keepLang(t)
	root := project(t, frozenCfg)
	err := refuseIfFrozen(find(t, root, "generated-paths"))
	if err == nil || !strings.Contains(err.Error(), "`generated-paths` will not run") ||
		!strings.Contains(err.Error(), "rotate the key") {
		t.Fatalf("want the freeze refusal naming the command and the reason, got %v", err)
	}
	// The same command pointed at a project that is not frozen runs.
	if err := refuseIfFrozen(find(t, project(t, "version: 1\n"), "generated-paths")); err != nil {
		t.Errorf("a project that is not frozen was refused: %v", err)
	}
}

// What only reads keeps working while frozen — including a SUBcommand of an allowed one,
// because the whole parent chain is checked.
func TestFrozenProjectKeepsTheReadingCommands(t *testing.T) {
	t.Run("CLRTC-B02: The commands that only read, and thaw, run while frozen", func(t *testing.T) {})
	t.Run("CLRTC-B03: A subcommand of an allowed command runs while frozen", func(t *testing.T) {})
	keepLang(t)
	root := project(t, frozenCfg)
	t.Chdir(root) // commands without --root read the working directory
	for _, path := range [][]string{{"thaw"}, {"freeze"}, {"status"}, {"doctor"}, {"guide"},
		{"guide", "work"}, {"help"}, {"completion"}, {"coverage"}, {"impact"}} {
		if err := refuseIfFrozen(find(t, root, path...)); err != nil {
			t.Errorf("`%s` was refused while frozen: %v", strings.Join(path, " "), err)
		}
	}
}

// No config, or a config that does not load, is not "frozen": the command proceeds, and
// whoever investigates is not sent the wrong way.
func TestMissingOrBrokenConfigIsNotFrozen(t *testing.T) {
	t.Run("CLRTC-B04: A missing or broken config is not frozen", func(t *testing.T) {})
	keepLang(t)
	for name, cfg := range map[string]string{"missing": "", "broken": "enabled: [\n", "unknown key": "enabled: false\nnope: 1\n"} {
		if err := refuseIfFrozen(find(t, project(t, cfg), "generated-paths")); err != nil {
			t.Errorf("%s config: refused with %v", name, err)
		}
	}
}

// The project's top-level `lang:` is applied before any output — read loosely, so a
// broken YAML never stops a command and an unknown language is simply ignored.
func TestProjectLangIsAppliedBeforeTheCommand(t *testing.T) {
	t.Run("CLRTC-B06: The project's top-level lang is applied before the command runs", func(t *testing.T) {})
	keepLang(t)
	apply := func(cfg string) string {
		t.Helper()
		_ = i18n.Set("en")
		applyProjectLang(find(t, project(t, cfg), "generated-paths"))
		return i18n.Current()
	}
	for cfg, want := range map[string]string{
		"version: 1\nlang: pt-BR\n":  "pt-BR",
		"lang: \"es\"\nlayers: [\n":  "es", // invalid YAML does not block
		"workflow:\n  lang: pt-BR\n": "en", // nested: not the project's language
		"lang: xx\n":                 "en", // unknown: ignored
		"":                           "en", // no config
	} {
		if got := apply(cfg); got != want {
			t.Errorf("config %q: language = %s, want %s", cfg, got, want)
		}
	}
}

// Through the real root: the language comes FIRST, so even the freeze refusal speaks the
// project's language, and the root prints neither the error nor the usage itself.
func TestTheFreezeRefusalSpeaksTheProjectLanguage(t *testing.T) {
	t.Run("CLRTC-B07: The root prints neither the error nor the usage", func(t *testing.T) {})
	t.Run("CLRTC-X01: The freeze refusal is written in the project language", func(t *testing.T) {})
	keepLang(t)
	_ = i18n.Set("en")
	root := project(t, "lang: pt-BR\n"+frozenCfg)
	c := newRootCmd()
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs([]string{"generated-paths", "--root", root})
	err := c.Execute()
	if err == nil || !strings.Contains(err.Error(), "o projeto está CONGELADO") {
		t.Fatalf("want the refusal in pt-BR, got %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("the root printed on its own (the error must come out once, from main):\n%s", out.String())
	}
}

func TestLangNoYAML_crlf(t *testing.T) {
	t.Run("CLRTC-B08: The language is read from a CRLF file", func(t *testing.T) {})
	for _, src := range []string{"version: 6\r\nlang: pt-BR\r\n", "lang: \"es\"  \r\n"} {
		m := langNoYAML.FindSubmatch([]byte(src))
		if m == nil || (string(m[1]) != "pt-BR" && string(m[1]) != "es") {
			t.Errorf("%q: the language is read, got %q", src, m)
		}
	}
}

func TestEveryCommandTakesFilesTheSameWay(t *testing.T) {
	t.Run("CLRTC-I02: Every command that takes several files reads them the same way", func(t *testing.T) {})
	// several files: a usage that names files and repeats them — `<file>...`, `[spec files...]`;
	// `[layers...]` or `<card>...` are other lists.
	takesFiles := regexp.MustCompile(`(?i)file[^.\]>]*[\]>]?\.\.\.|files?[^\]>]*\.\.\.`)
	var multi int
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		if takesFiles.MatchString(c.Use) {
			multi++
			if c.Annotations[common.FilesAnnotation] != "true" {
				t.Errorf("`anchors %s` takes several files and does not read them through common.TakesFiles: "+
					"a list that works on another command fails on this one", c.CommandPath())
			}
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(newRootCmd())
	if multi < 4 {
		t.Fatalf("found %d command(s) taking several files — keep-evidence, touch, stamp and renumber at least; the walk broke", multi)
	}
}
