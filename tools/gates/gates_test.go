package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSummarizeOSV(t *testing.T) {
	report := `{"results":[{"packages":[
		{"package":{"name":"golang.org/x/net"},"vulnerabilities":[{"database_specific":{"severity":"HIGH"}},{}]},
		{"package":{"name":"golang.org/x/text"},"vulnerabilities":[{"database_specific":{"severity":"HIGH"}}]}]}]}`
	got := summarizeOSV([]byte(report))
	for _, want := range []string{"3 vulnerabilities in 2 packages", "HIGH         2  (golang.org/x/net, golang.org/x/text)", "UNKNOWN      1  (golang.org/x/net)"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if got := summarizeOSV([]byte(`{"results":[]}`)); !strings.Contains(got, "no known vulnerability") {
		t.Errorf("no vulnerability: %s", got)
	}
	if got := summarizeOSV([]byte("not json")); !strings.Contains(got, "SKIPPED") {
		t.Errorf("an unreadable report is a skip, never a pass: %s", got)
	}
}

func TestClassifyLicenses(t *testing.T) {
	csv := "github.com/co2-lab/anchors/internal/x,u,Unknown\n" +
		"github.com/a/mit,u,MIT\n" +
		"github.com/b/mpl,u,MPL-2.0\n" +
		"github.com/c/none,u,Unknown\n" +
		"github.com/d/gplx,u,GPL-2.0-with-classpath-exception\n"
	text, bad := classifyLicenses([]byte(csv), "github.com/co2-lab/anchors")
	if bad {
		t.Errorf("a GPL with an explicit exception is not strong copyleft:\n%s", text)
	}
	for _, want := range []string{"4 third-party dependencies", "weak copyleft (1)", "licence NOT declared (1)", "no strong copyleft licence"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in:\n%s", want, text)
		}
	}
	text, bad = classifyLicenses([]byte(csv+"github.com/e/agpl,u,AGPL-3.0\n"), "github.com/co2-lab/anchors")
	if !bad || !strings.Contains(text, "STRONG copyleft (1)") {
		t.Errorf("an AGPL dependency fails the gate:\n%s", text)
	}
}

func TestSummarizeSemgrep(t *testing.T) {
	report := `{"results":[
		{"check_id":"go.lang.security.use-tls","path":"a.go","start":{"line":3}},
		{"check_id":"go.lang.security.use-tls","path":"b.go","start":{"line":9}},
		{"check_id":"go.lang.security.use-of-sha1","path":"b.go","start":{"line":1}}]}`
	got := summarizeSemgrep([]byte(report))
	for _, want := range []string{"3 gosec finding(s) in 2 file(s)", "2x use-tls  (e.g. a.go:3)", "1x use-of-sha1  (e.g. b.go:1)"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if !strings.Contains(summarizeSemgrep([]byte(`{"results":[]}`)), "no known insecure pattern") {
		t.Error("no finding should say so")
	}
	if !strings.Contains(summarizeSemgrep([]byte("x")), "SKIPPED") {
		t.Error("an unreadable report is a skip")
	}
}

func TestSummarizeSBOM(t *testing.T) {
	got := summarizeSBOM([]byte(`{"components":[{"type":"library"},{"type":"library"},{"type":"file"}]}`), "r.json")
	for _, want := range []string{"3 components", "library        2", "file           1", "→ r.json"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestMutationEntriesAndExcludes(t *testing.T) {
	got := mutationEntries([]string{"internal/migra/format_6.go", "internal/gate"})
	want := []mutationEntry{{pkg: "internal/migra", only: "format_6.go"}, {pkg: "internal/gate"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("entries = %+v, want %+v", got, want)
	}
	ex := excludesFor(want[0], []string{"format_5.go", "format_6.go", "step.go"})
	if !reflect.DeepEqual(ex, []string{"-E", "/", "-E", `^format_5\.go$`, "-E", `^step\.go$`}) {
		t.Errorf("a file's run excludes its siblings by exact name and the subfolders, got %v", ex)
	}
	if ex := excludesFor(want[1], []string{"a.go"}); !reflect.DeepEqual(ex, []string{"-E", "/"}) {
		t.Errorf("a package's run excludes only the subfolders, got %v", ex)
	}
}

func TestMutationReportsArePrefixedAndMerged(t *testing.T) {
	a := prefixMutationReport(map[string]any{"go_module": "m", "files": []any{map[string]any{"file_name": "queue.go"}}}, "internal/queue")
	b := prefixMutationReport(emptyMutationReport("m", mutationEntry{pkg: "internal/migra", only: "format_6.go"}, []string{"format_5.go", "format_6.go"}), "internal/migra")
	merged := mergeMutationReports([]map[string]any{a, b})
	files := merged["files"].([]any)
	if merged["go_module"] != "m" || len(files) != 2 {
		t.Fatalf("merged = %+v", merged)
	}
	if files[0].(map[string]any)["file_name"] != "internal/queue/queue.go" || files[1].(map[string]any)["file_name"] != "internal/migra/format_6.go" {
		t.Errorf("each file is named from the root, and a file with nothing to mutate is listed alone: %+v", files)
	}
	if m := files[1].(map[string]any)["mutations"].([]any); len(m) != 0 {
		t.Errorf("nothing to mutate lists no mutation: %v", m)
	}
}

func TestReportFormat(t *testing.T) {
	if code := reportFormat(""); code != 0 {
		t.Errorf("nothing listed passes, got %d", code)
	}
	if code := reportFormat("internal/gate/a_test.go\n"); code != 1 {
		t.Errorf("a file out of form fails, got %d", code)
	}
}
