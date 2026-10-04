// @anchors
//   code: MGTSA
//   ref: MGFLM

package migra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func escreve(t *testing.T, dir, nome, conteudo string) string {
	t.Helper()
	p := filepath.Join(dir, nome)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// withSteps swaps the registry for the test and restores it afterwards.
func withSteps(t *testing.T, s ...Step) {
	t.Helper()
	original := steps
	t.Cleanup(func() { steps = original })
	steps = nil
	for _, p := range s {
		Register(p)
	}
}

func TestFormatOf_readsTheTopLevelVersion(t *testing.T) {
	t.Run("MGFLM-B01: The format is the top-level version line, and its absence means format 1", func(t *testing.T) {})
	dir := t.TempDir()
	none := escreve(t, dir, "a/anchors.graph.yaml", "# header\n# derived\nnodes: []\n")
	three := escreve(t, dir, "b/anchors.graph.yaml", "# header\nversion: 3\nnodes: []\n")
	if f, err := FormatOf(none); err != nil || f != 1 {
		t.Errorf("without `version:` the format is 1; got %d, %v", f, err)
	}
	if f, err := FormatOf(three); err != nil || f != 3 {
		t.Errorf("`version: 3` is format 3; got %d, %v", f, err)
	}
}

func TestMigrateFile_alreadyAtTheTargetIsUntouched(t *testing.T) {
	t.Run("MGFLM-B02: A file already at the target is left untouched", func(t *testing.T) {})
	original := "version: 2\ngerado_por: dev\n"
	p := escreve(t, t.TempDir(), "anchors.graph.yaml", original)
	r, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if r.Changed || len(r.Replaced) != 0 {
		t.Errorf("a file at the target should report no change; got %+v", r)
	}
	if string(mustRead(t, p)) != original {
		t.Error("a file at the target was rewritten")
	}
	// Above the target is the same answer, not an error.
	above := escreve(t, t.TempDir(), "anchors.graph.yaml", "version: 3\ngerado_por: dev\n")
	if r, err := MigrateFile(above, 2, false); err != nil || r.Changed {
		t.Errorf("a file above the target should be left alone; got %+v, %v", r, err)
	}
}

// THE CASE THAT MOTIVATES IT ALL: `julgamentos` became `judgments`, and losing the key would
// not raise an error — it would be SILENCE. The AI judgment stamps would vanish and `check`
// would redo them all, charging again what someone already answered.
func TestMigrateFile_renamesTheMapKeys(t *testing.T) {
	t.Run("MGFLM-B03: Old map keys are renamed where they are keys and keep their values", func(t *testing.T) {})
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", `# header
version: 1
gerado_por: 0.1.83
nodes:
    - id: a.spec.md
      kind: spec
      code_declarado: true
edges:
    - from: a.spec.md
      to: a.ts
      julgamentos:
        - gate: doc-self-contained
          verdict: pass
`)
	r, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Changed {
		t.Fatal("the file should have been changed")
	}

	text := string(mustRead(t, p))
	for _, novo := range []string{"generated_by:", "code_declared:", "judgments:", "version: 2"} {
		if !strings.Contains(text, novo) {
			t.Errorf("expected %q in the migrated file", novo)
		}
	}
	for _, velho := range []string{"gerado_por:", "code_declarado:", "julgamentos:"} {
		if strings.Contains(text, velho) {
			t.Errorf("the old key %q is still in the file", velho)
		}
	}
	// The VALUE has to survive: a migration that renames and loses the content is worse than
	// none, because it looks like it worked.
	if !strings.Contains(text, "0.1.83") || !strings.Contains(text, "doc-self-contained") {
		t.Error("the migration lost a value while renaming the key")
	}
}

// Anchoring on the key is what separates migration from blind rewriting: `julgamentos`
// appears in COMMENTS and prose inside Anchors' files, and a substring replacement would
// corrupt text that is no key at all.
func TestMigrateFile_doesNotTouchWhatIsNotAKey(t *testing.T) {
	t.Run("MGFLM-B03: Old map keys are renamed where they are keys and keep their values", func(t *testing.T) {})
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", `# os julgamentos ficam aqui; gerado_por diz quem gravou
version: 1
nodes:
    - id: a.spec.md
      # este nó tem julgamentos pendentes
      title: "sobre julgamentos e code_declarado"
`)
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	if !strings.Contains(text, "# os julgamentos ficam aqui; gerado_por diz quem gravou") {
		t.Error("the comment was rewritten — the replacement is not anchored on the key")
	}
	if !strings.Contains(text, `"sobre julgamentos e code_declarado"`) {
		t.Error("a VALUE holding the key's text was rewritten")
	}
}

// `trinca_opcional` lives in `anchors.yaml`, and `julgamentos` in the map. Renaming the key in
// the wrong file would rewrite something that happens to have the same name.
func TestMigrateFile_eachKeyInItsOwnFile(t *testing.T) {
	t.Run("MGFLM-B04: Each rename applies only to its own file", func(t *testing.T) {})
	dir := t.TempDir()
	cfg := escreve(t, dir, "anchors.yaml", `version: 1
layers:
    spec:
        trinca_opcional: [covered-by, tested-by]
`)
	if _, err := MigrateFile(cfg, 2, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, cfg))
	if !strings.Contains(text, "triad_optional:") {
		t.Error("the `trinca_opcional` of anchors.yaml was not migrated")
	}
	// And the VALUE of the waiver, which is what cannot vanish: without it the
	// `triad-complete` gate charges again what the team decided not to have.
	if !strings.Contains(text, "[covered-by, tested-by]") {
		t.Error("the migration lost the declared waiver")
	}

	// THE OTHER HALF of the rule: a key OF THE MAP cannot be renamed in `anchors.yaml`.
	//
	// The first version of this test only checked that the right key migrated, and survived
	// the mutation that removes the per-file filter — because a common `anchors.yaml` has no
	// `julgamentos` to rename by mistake. Here it has one, on purpose: it is a GATE name, and
	// the gate is called that in any file.
	cfg2 := escreve(t, dir, "anchors.yaml", `version: 1
gates:
    - id: julgamentos
      when: [pre-commit]
layers:
    spec:
        gerado_por: quem-sabe
`)
	if _, err := MigrateFile(cfg2, 2, false); err != nil {
		t.Fatal(err)
	}
	t2 := string(mustRead(t, cfg2))
	if strings.Contains(t2, "judgments") {
		t.Error("a key OF THE MAP was renamed in anchors.yaml")
	}
	if strings.Contains(t2, "generated_by") {
		t.Error("`gerado_por` is a key OF THE MAP and should not be touched here")
	}
}

func TestMigrateFile_valueRenamedOnlyUnderItsKeyAndWhole(t *testing.T) {
	t.Run("MGFLM-B05: A value is renamed only under its key and when it matches whole", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.yaml", `version: 1
gates:
    - name: regra-cumprida
      when: regra-cumprida
    - name: regra-cumprida-extra
`)
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	for _, want := range []string{"- name: rule-fulfilled\n", "when: regra-cumprida\n", "- name: regra-cumprida-extra\n"} {
		if !strings.Contains(text, want) {
			t.Errorf("expected %q in:\n%s", want, text)
		}
	}
}

func TestMigrateFile_stepsAppliedInOrderOnThePreviousResult(t *testing.T) {
	t.Run("MGFLM-B06: Steps are applied in order, each on the result of the previous", func(t *testing.T) {})
	withSteps(t,
		Step{To: 3, Why: "b to c", RenameKeys: map[string]map[string]string{"anchors.yaml": {"b": "c"}}},
		Step{To: 2, Why: "a to b", RenameKeys: map[string]map[string]string{"anchors.yaml": {"a": "b"}}},
	)
	p := escreve(t, t.TempDir(), "anchors.yaml", "version: 1\na: 1\n")
	if _, err := MigrateFile(p, 3, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	if !strings.Contains(text, "c: 1") || strings.Contains(text, "a:") || strings.Contains(text, "b:") {
		t.Errorf("expected `c: 1` only, got:\n%s", text)
	}
}

func TestMigrateFile_versionRaisedEvenWithNothingElse(t *testing.T) {
	t.Run("MGFLM-B07: The version is raised even when nothing else changes", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.graph.yaml", "version: 1\nnodes: []\n")
	r, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Changed || string(mustRead(t, p)) != "version: 2\nnodes: []\n" {
		t.Errorf("expected only the version raised; changed=%v text=%q", r.Changed, mustRead(t, p))
	}
}

// A map from before the `version:` field existed: it is format 1, and migrating it requires
// INSERTING the line — not just replacing it.
func TestMigrateFile_mapWithoutVersionGetsTheField(t *testing.T) {
	t.Run("MGFLM-B08: A map or a configuration without a version gets one after its comment header", func(t *testing.T) {})
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", "# anchors.graph.yaml — o mapa\n# derivado\nnodes: []\n")

	if de, _ := FormatOf(p); de != 1 {
		t.Fatalf("without `version:` the format should be 1, got %d", de)
	}
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	if !strings.Contains(text, "version: 2") {
		t.Error("the `version:` field was not inserted")
	}
	// After the header, not before: `Save` writes the comments first, and a `version:` above
	// them would move on the next save.
	if strings.HasPrefix(text, "version:") {
		t.Error("`version:` went in BEFORE the comment header")
	}
	// The configuration too: it was migrated but never stamped, so it stayed format 1 and
	// every run found it old again (`anchors migrate` said "already in format N").
	cfg := escreve(t, dir, "anchors.yaml", "# config\nproject: x\n")
	r, err := MigrateFile(cfg, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(mustRead(t, cfg)); !r.Changed || got != "# config\nversion: 2\nproject: x\n" {
		t.Errorf("the configuration gets `version: 2` after its header; changed=%v text=%q", r.Changed, got)
	}
	if f, _ := FormatOf(cfg); f != 2 {
		t.Errorf("after the migration the configuration is format 2, got %d", f)
	}
}

func TestMigrateFile_countsEachRenameByItsOldForm(t *testing.T) {
	t.Run("MGFLM-B09: The result counts each rename by its old form", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.yaml", `version: 1
layers:
    spec:
        trinca_opcional: [tested-by]
    feature:
        trinca_opcional: [covered-by]
gates:
    - name: regra-cumprida
`)
	r, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if r.Replaced["trinca_opcional"] != 2 || r.Replaced["name: regra-cumprida"] != 1 || !r.Changed {
		t.Errorf("unexpected counts: %+v", r)
	}
}

// `dryRun` MEASURES without writing — it is what lets the `doctor` warn without touching the
// repository of someone who did not ask.
func TestMigrateFile_dryRunDoesNotWrite(t *testing.T) {
	t.Run("MGFLM-B10: A dry run reports without writing", func(t *testing.T) {})
	dir := t.TempDir()
	original := "version: 1\ngerado_por: dev\nnodes: []\n"
	p := escreve(t, dir, "anchors.graph.yaml", original)

	r, err := MigrateFile(p, 2, true)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Changed {
		t.Error("the dry run should REPORT there is something to change")
	}
	if r.Replaced["gerado_por"] != 1 {
		t.Errorf("the dry run should count the key to rename; got %v", r.Replaced)
	}
	if string(mustRead(t, p)) != original {
		t.Error("the dry run WROTE to the file")
	}
}

// Migrating twice changes nothing the second time: `version:` goes up EVEN with no key
// renamed, and without it the migration would run again on every command.
func TestMigrateFile_isIdempotent(t *testing.T) {
	t.Run("MGFLM-I01: Migrating twice changes nothing the second time", func(t *testing.T) {})
	dir := t.TempDir()
	p := escreve(t, dir, "anchors.graph.yaml", "version: 1\nnodes: []\nedges: []\n")

	r1, err := MigrateFile(p, 2, false)
	if err != nil || !r1.Changed {
		t.Fatalf("the first migration should change the file; err=%v", err)
	}
	first := mustRead(t, p)

	r2, err := MigrateFile(p, 2, false)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Changed {
		t.Error("the second migration should change nothing")
	}
	if string(mustRead(t, p)) != string(first) {
		t.Error("migrating again changed the file")
	}
}

func TestMigrateFile_textTheParserRefusesStillMigrates(t *testing.T) {
	t.Run("MGFLM-X01: A file the YAML parser would refuse is still migrated", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.graph.yaml", "version: 1\ngerado_por: dev\nbroken: [a, b\n")
	if _, err := MigrateFile(p, 2, false); err != nil {
		t.Fatalf("an unparsable file must still migrate: %v", err)
	}
	text := string(mustRead(t, p))
	for _, want := range []string{"generated_by: dev", "version: 2", "broken: [a, b"} {
		if !strings.Contains(text, want) {
			t.Errorf("expected %q in:\n%s", want, text)
		}
	}
}

func TestMigrateFile_missingFileIsAnError(t *testing.T) {
	t.Run("MGFLM-E01: A missing file is an error and nothing is written", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "anchors.graph.yaml")
	if _, err := MigrateFile(p, 2, false); err == nil {
		t.Fatal("a missing file must be an error")
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("migrating a missing file created it")
	}
}

func TestFormatOf_hugeVersionIsAnError(t *testing.T) {
	t.Run("MGFLM-E02: A version that is not a number is an error and nothing is written", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.yaml", "version: 99999999999999999999999\n")
	_, err := FormatOf(p)
	if err == nil || !strings.Contains(err.Error(), "not a number") {
		t.Errorf("expected a `not a number` error, got %v", err)
	}
}

func TestMigrateFile_holeInTheChainLeavesTheFileUntouched(t *testing.T) {
	t.Run("MGFLM-E03: A hole in the chain leaves the file untouched", func(t *testing.T) {})
	withSteps(t, Step{To: 2, Why: "two", RenameKeys: map[string]map[string]string{
		"anchors.graph.yaml": {"gerado_por": "generated_by"}}})
	original := "version: 1\ngerado_por: dev\n"
	p := escreve(t, t.TempDir(), "anchors.graph.yaml", original)
	_, err := MigrateFile(p, 3, false)
	if err == nil || !strings.Contains(err.Error(), "format 3") {
		t.Errorf("expected an error naming format 3, got %v", err)
	}
	if string(mustRead(t, p)) != original {
		t.Error("a failed migration wrote to the file")
	}
}

func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// `version: 3  # why` is valid YAML and read as format 3; the anchored pattern missed it,
// took the file for format 1, and then inserted a second `version:` line into a map.
func TestVersionLineWithTrailingComment(t *testing.T) {
	t.Run("MGFLM-B11: A version line with a trailing comment is read and raised keeping the comment", func(t *testing.T) {
		p := escreve(t, t.TempDir(), "anchors.graph.yaml", "version: 1  # the first\nnodes: []\n")
		if f, err := FormatOf(p); err != nil || f != 1 {
			t.Fatalf("want format 1, got %d, %v", f, err)
		}
		if _, err := MigrateFile(p, 2, false); err != nil {
			t.Fatal(err)
		}
		if got := string(mustRead(t, p)); got != "version: 2  # the first\nnodes: []\n" {
			t.Errorf("want the line raised in place with its comment, got %q", got)
		}
		three := escreve(t, t.TempDir(), "anchors.yaml", "version: 3 # c\n")
		if f, err := FormatOf(three); err != nil || f != 3 {
			t.Errorf("want format 3, got %d, %v", f, err)
		}
	})
}

// `version: abc` was not matched at all, so the file read as format 1 and a map got a
// second `version:` line inserted above the bad one.
func TestVersionThatIsNotANumber(t *testing.T) {
	t.Run("MGFLM-E02: A version that is not a number is an error and nothing is written", func(t *testing.T) {
		original := "# h\nversion: abc\nnodes: []\n"
		p := escreve(t, t.TempDir(), "anchors.graph.yaml", original)
		if _, err := FormatOf(p); err == nil || !strings.Contains(err.Error(), "not a number") {
			t.Errorf("want a `not a number` error, got %v", err)
		}
		if _, err := MigrateFile(p, 2, false); err == nil {
			t.Error("migrating a file whose version is not a number must fail")
		}
		if got := string(mustRead(t, p)); got != original {
			t.Errorf("the file was written: %q", got)
		}
	})
}

func TestMigrateFile_crlf(t *testing.T) {
	t.Run("MGFLM-B12: A CRLF file is read and migrated as it is", func(t *testing.T) {})
	dir := t.TempDir()
	current := escreve(t, dir, "a/anchors.graph.yaml", "# header\r\nversion: 2\r\ngerado_por: dev\r\n")
	if f, err := FormatOf(current); err != nil || f != 2 {
		t.Fatalf("a CRLF version line is read, got %d %v", f, err)
	}
	if r, err := MigrateFile(current, 2, false); err != nil || r.Changed {
		t.Errorf("a current CRLF file is untouched, got %+v %v", r, err)
	}
	older := escreve(t, dir, "b/anchors.graph.yaml", "version: 1 # old\r\nnodes: []\r\n")
	if _, err := MigrateFile(older, 2, false); err != nil {
		t.Fatal(err)
	}
	if got := string(mustRead(t, older)); strings.Count(got, "version:") != 1 || !strings.Contains(got, "version: 2 # old\r\n") {
		t.Errorf("the version is stamped in place, keeping CRLF, got %q", got)
	}
	none := escreve(t, dir, "c/anchors.graph.yaml", "# header\r\nnodes: []\r\n")
	if _, err := MigrateFile(none, 2, false); err != nil {
		t.Fatal(err)
	}
	if got := string(mustRead(t, none)); strings.Count(got, "version:") != 1 || !strings.Contains(got, "version: 2\r\n") {
		t.Errorf("an inserted version ends in CRLF, got %q", got)
	}
}
