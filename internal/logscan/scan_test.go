// @anchors
//   code: SCTSB
//   ref: LGSCL

package logscan

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func projectWithLogs(t *testing.T, files map[string]string) (string, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(root, "logs", name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, &config.Config{Logs: &config.LogsConfig{
		Paths:     []string{"logs/*"},
		Timestamp: `(\d{4}-\d{2}-\d{2}[T ][\d:]+)`,
	}}
}

// THE FORMAT DOES NOT MATTER, and that is what makes scanning possible without dictating
// anything. The three lines below say the same thing in JSON, plain text and syslog — and
// the scanner reads them alike, because it looks for the CODE, not the format.
func TestScan_theFormatDoesNotMatter(t *testing.T) {
	t.Run("LGSCL-B02: The format does not matter", func(t *testing.T) {})
	t.Run("LGSCL-B07: Each occurrence records its time window", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{
		"a.jsonl":  `{"level":"error","msg":"#[CRED-E01] failed","ts":"2026-09-05T10:00:00Z"}` + "\n",
		"b.txt":    "2026-09-19 11:00:00 ERROR #[CRED-E01] refused\n",
		"c.syslog": "<134>1 2026-09-01T08:00:00Z app - - - #[CRED-E01] refused\n",
	})
	r, err := Scan(root, cfg, specDeclaring(t, root, "CRED-E01"))
	if err != nil || r == nil {
		t.Fatalf("scan: %v", err)
	}
	if len(r.Occurrences) != 1 || r.Occurrences[0].Rule != "CRED-E01" {
		t.Fatalf("expected one rule across the three formats, got %+v", r.Occurrences)
	}
	if r.Occurrences[0].Count != 3 {
		t.Errorf("the three formats must count alike, got %d", r.Occurrences[0].Count)
	}
	// The window spans the three files: a failure that happened a lot and STOPPED says
	// something different from one happening now, and `last` is what tells them apart.
	if !strings.HasPrefix(r.Occurrences[0].First, "2026-09-01") || !strings.HasPrefix(r.Occurrences[0].Last, "2026-09-19") {
		t.Errorf("the window must span the earliest and the latest line: %+v", r.Occurrences[0])
	}
	if r.Files != 3 || r.Lines != 3 {
		t.Errorf("read %d files and %d lines, want 3 and 3", r.Files, r.Lines)
	}
}

// The most valuable finding of this layer, and the one no observability tool gives: the
// log shows an error the spec never foresaw. Without this separation it would be
// discarded in silence — which is exactly the case nobody is watching.
func TestScan_separatesWhatNoSpecDeclares(t *testing.T) {
	t.Run("LGSCL-B04: An undeclared failure is reported apart", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{
		"a.log": "ERROR #[CRED-E01] known\nERROR #[ORPHA-E99] nobody declared this\n",
	})
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, "src", "x.spec.md"),
		[]byte("| `CRED-E01` | known | refuses |\n"), 0o644)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "src/x.spec.md", Kind: mapx.KindSpec}}}

	r, _ := Scan(root, cfg, g)
	if len(r.Occurrences) != 1 || r.Occurrences[0].Rule != "CRED-E01" {
		t.Errorf("only the declared rule becomes an occurrence: %+v", r.Occurrences)
	}
	if r.Unknown["ORPHA-E99"] != 1 {
		t.Errorf("the undeclared code must be reported, not discarded: %+v", r.Unknown)
	}
}

// Anchors never guesses where the logs are: with no declared path it scans nothing.
// Hunting for files that look like logs would read what it should not — a log usually
// carries exactly what must not leak.
func TestScan_withoutDeclaredPathScansNothing(t *testing.T) {
	t.Run("LGSCL-B01: Without a declared log path nothing is scanned", func(t *testing.T) {})
	root, _ := projectWithLogs(t, map[string]string{"a.log": "ERROR #[CRED-E01] x\n"})
	r, err := Scan(root, &config.Config{}, nil)
	if err != nil || r != nil {
		t.Errorf("expected no scan and no error, got %+v / %v", r, err)
	}
}

// The alias is the LEGACY bridge — the log that already exists and nobody will rewrite to
// carry the code. It comes after the code on purpose: a line that already says `CRED-E01`
// needs no guess, and letting the guess override the identity would trade certainty for
// approximation.
func TestScan_aliasIsOnlyForWhatCarriesNoCode(t *testing.T) {
	t.Run("LGSCL-B06: The alias only serves lines that carry no code", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{
		"a.log": "ERROR INSUFFICIENT_BALANCE legacy line\nERROR #[CRED-E02] carries the code\n",
	})
	cfg.Logs.Aliases = map[string]string{"CRED-E01": "INSUFFICIENT_BALANCE"}
	r, _ := Scan(root, cfg, specDeclaring(t, root, "CRED-E01", "CRED-E02"))
	got := map[string]int{}
	for _, o := range r.Occurrences {
		got[o.Rule] = o.Count
	}
	if got["CRED-E01"] != 1 {
		t.Errorf("the alias must catch the legacy line: %+v", got)
	}
	if got["CRED-E02"] != 1 {
		t.Errorf("the line carrying the code must not need the alias: %+v", got)
	}
}

// specDeclaring writes a spec cataloguing the given failures and returns a map with it.
func specDeclaring(t *testing.T, root string, codes ...string) *mapx.Graph {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	body := ""
	for _, c := range codes {
		body += "| `" + c + "` | condition | failure |\n"
	}
	if err := os.WriteFile(filepath.Join(root, "src", "x.spec.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return &mapx.Graph{Nodes: []mapx.Node{{ID: "src/x.spec.md", Kind: mapx.KindSpec}}}
}

// The pattern is built from the codes the project DECLARES, never from the generic shape
// of a code — and this is the test that keeps it that way.
//
// Measured with a bare `\b[A-Z0-9]{3,6}-E[0-9]{2}\b`: five captures in a five-line log,
// one of them a failure. The other four were a build id, a cache key, a SKU and a trace id
// inside a URL. Eighty per cent false positive — worse than the language heuristic this
// project already discarded for erring nine times in ten.
func TestScan_doesNotCaptureWhatMerelyLooksLikeACode(t *testing.T) {
	t.Run("LGSCL-B03: Look-alikes of a code are not captured", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{
		"a.log": "INFO build BUILD-E01 finished\n" +
			"INFO cache key USER-E42 evicted\n" +
			"INFO shipping SKU ITEM-E07 dispatched\n" +
			"INFO url https://x/y?trace=ABCD-E99\n" +
			"ERROR #[CRED-E01] real failure\n",
	})
	r, _ := Scan(root, cfg, specDeclaring(t, root, "CRED-E01"))
	if len(r.Occurrences) != 1 || r.Occurrences[0].Rule != "CRED-E01" {
		t.Fatalf("only the declared code is an occurrence, got %+v", r.Occurrences)
	}
	// None of the four look-alikes may be reported as unknown either: they sit on INFO
	// lines, and a build id on an INFO line is not a failure.
	if len(r.Unknown) != 0 {
		t.Errorf("look-alikes on non-failure lines must not be reported: %+v", r.Unknown)
	}
}

// The second pass — the undeclared code — is where the measured risk lives, so it demands
// that the line SAY it is reporting a failure. Without that, an INFO line brings the four
// false positives back.
func TestScan_theUndeclaredCodeNeedsAFailureLine(t *testing.T) {
	root, cfg := projectWithLogs(t, map[string]string{
		"a.log": "INFO cache USER-E42 evicted\nERROR #[ORPHA-E88] undeclared failure\nINFO #[ORPHA-E98] delimited\n",
	})
	r, _ := Scan(root, cfg, specDeclaring(t, root, "CRED-E01"))
	if r.Unknown["ORPHA-E88"] != 1 || r.Unknown["ORPHA-E98"] != 1 {
		t.Errorf("a delimited undeclared code must be reported, with or without a level word: %+v", r.Unknown)
	}
	if _, ok := r.Unknown["USER-E42"]; ok {
		t.Errorf("a look-alike on an INFO line must not be reported: %+v", r.Unknown)
	}

	// With no delimiter the level word is the only defence left: the bare code is believed
	// only on a failure line.
	bare, bcfg := projectWithLogs(t, map[string]string{
		"a.log": "ERROR ORPHA-E77 bare failure\nINFO ORPHA-E76 bare id\n",
	})
	bcfg.Logs.Delimiters = []string{"", ""}
	rb, _ := Scan(bare, bcfg, specDeclaring(t, bare, "CRED-E01"))
	if rb.Unknown["ORPHA-E77"] != 1 {
		t.Errorf("a bare undeclared code on a failure line must be reported: %+v", rb.Unknown)
	}
	if _, ok := rb.Unknown["ORPHA-E76"]; ok {
		t.Errorf("a bare undeclared code on an INFO line must not be reported: %+v", rb.Unknown)
	}
}

// The DELIMITER is what separates the failure from what merely LOOKS like one, and the
// choice is measured against real noise:
//
//	bare code     5 false positives   build id, cache key, SKU, trace id
//	[CRED-E01]    2 false positives   markdown in prose, array index (`logs[...]`)
//	#[CRED-E01]   0 false positives
//
// All three catch the real failures alike — what changes is what they capture BESIDES. And
// writing `#[` instead of `[` costs whoever logs exactly the same.
func TestScan_theDelimiterIsWhatRemovesTheAmbiguity(t *testing.T) {
	t.Run("LGSCL-B05: The delimiter removes the ambiguity", func(t *testing.T) {})
	noise := "INFO markdown [CRED-E01] cited in prose\n" +
		"INFO array logs[CRED-E01] index\n" +
		`INFO json {"tags":["CRED-E01"]}` + "\n"
	root, cfg := projectWithLogs(t, map[string]string{
		"a.log": noise + "ERROR #[CRED-E01] the real one\n",
	})
	g := specDeclaring(t, root, "CRED-E01")
	r, _ := Scan(root, cfg, g)
	if len(r.Occurrences) != 1 || r.Occurrences[0].Count != 1 {
		t.Fatalf("only the delimited line is an occurrence, got %+v", r.Occurrences)
	}

	// With an empty pair the project goes back to the bare code — the legacy log nobody
	// will rewrite — and pays the measured cost: the three noise lines come back.
	cfg.Logs.Delimiters = []string{"", ""}
	r2, _ := Scan(root, cfg, g)
	if r2.Occurrences[0].Count == 1 {
		t.Error("with no delimiter the bare code captures the noise too — that is the cost the default avoids")
	}
}

func TestScan_occurrencesSortedByRule(t *testing.T) {
	t.Run("LGSCL-B08: Occurrences are sorted by rule", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{"a.log": "ERROR #[CRED-E02] b\nERROR #[CRED-E01] a\n"})
	r, _ := Scan(root, cfg, specDeclaring(t, root, "CRED-E01", "CRED-E02"))
	if len(r.Occurrences) != 2 || r.Occurrences[0].Rule != "CRED-E01" || r.Occurrences[1].Rule != "CRED-E02" {
		t.Errorf("occurrences = %+v, want CRED-E01 then CRED-E02", r.Occurrences)
	}
}

func TestSpecFailureCodes_everyForm(t *testing.T) {
	t.Run("LGSCL-B09: A spec's failure rules are listed in every form", func(t *testing.T) {})
	spec := "### ABCDE-E01 — a heading\n\n| `ABCDE-E02` | a row |\n\n- **ABCDE-E03** a bullet\n\n| `ABCDE-B01` | a behaviour |\n"
	if got := SpecFailureCodes(spec); !reflect.DeepEqual(got, []string{"ABCDE-E01", "ABCDE-E02", "ABCDE-E03"}) {
		t.Errorf("SpecFailureCodes = %v, want the three -E rules", got)
	}
}

func TestScan_invalidPatternsFail(t *testing.T) {
	t.Run("LGSCL-E01: An invalid alias or timestamp pattern fails the scan", func(t *testing.T) {})
	root, cfg := projectWithLogs(t, map[string]string{"a.log": "ERROR #[CRED-E01] x\n"})
	cfg.Logs.Aliases = map[string]string{"CRED-E01": "("}
	if r, err := Scan(root, cfg, nil); err == nil || r != nil {
		t.Errorf("an invalid alias = %+v, %v; want the error and no result", r, err)
	}
	cfg.Logs.Aliases = nil
	cfg.Logs.Timestamp = "("
	if r, err := Scan(root, cfg, nil); err == nil || r != nil {
		t.Errorf("an invalid timestamp = %+v, %v; want the error and no result", r, err)
	}
}
