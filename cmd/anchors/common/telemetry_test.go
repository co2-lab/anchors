package common

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// These tests build the emitter and never emit: NoticeTelemetry only CREATES it, and no
// test here calls Emit. The endpoint is still pointed at a closed local port, so that a
// future change that starts sending on creation fails locally instead of reaching the
// network.
func isolateTelemetry(t *testing.T) {
	t.Helper()
	t.Setenv("ANCHORS_TELEMETRY_ENDPOINT", "http://127.0.0.1:1/never")
	t.Setenv("ANCHORS_TELEMETRY_KEY", "")
	t.Setenv("ANCHORS_TELEMETRY", "")
	old := Emitter
	Emitter = nil
	t.Cleanup(func() { Emitter = old })
}

func captureStderr(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = old }()
	f()
	w.Close()
	b, _ := io.ReadAll(r)
	return string(b)
}

func cmdWithRoot(root string) *cobra.Command {
	cmd := &cobra.Command{Use: "probe"}
	cmd.Flags().String("root", ".", "project root")
	if root != "" {
		_ = cmd.Flags().Set("root", root)
	}
	return cmd
}

// First run on a machine: the notice is shown, with the ways to turn it off, the emitter
// is turned on — and the second run is silent.
func TestNoticeTelemetry_noticeOnceThenQuiet(t *testing.T) {
	isolateTelemetry(t)
	root := t.TempDir()
	cmd := cmdWithRoot(root)

	first := captureStderr(t, func() { NoticeTelemetry(cmd) })
	if !strings.Contains(first, "ANCHORS_TELEMETRY=off") {
		t.Errorf("the notice does not say how to turn it off:\n%s", first)
	}
	if Emitter == nil {
		t.Fatal("telemetry on by default, and the emitter was not created")
	}
	if _, err := os.Stat(filepath.Join(root, ".anchors", "telemetry-noticed")); err != nil {
		t.Errorf("the notice was not marked as shown in the project: %v", err)
	}

	second := captureStderr(t, func() { NoticeTelemetry(cmd) })
	if second != "" {
		t.Errorf("the notice was repeated on the second run:\n%s", second)
	}
	// Nothing in flight: the flush returns without waiting on anyone.
	FlushTelemetry()
}

// The environment turns it off with no file edited: no notice, no emitter.
func TestNoticeTelemetry_envOffCreatesNothing(t *testing.T) {
	isolateTelemetry(t)
	t.Setenv("ANCHORS_TELEMETRY", "off")
	root := t.TempDir()

	out := captureStderr(t, func() { NoticeTelemetry(cmdWithRoot(root)) })
	if out != "" || Emitter != nil {
		t.Errorf("telemetry off, and still: notice=%q emitter=%v", out, Emitter)
	}
	if _, err := os.Stat(filepath.Join(root, ".anchors")); err == nil {
		t.Error("telemetry off, and the project got a marker written")
	}
	FlushTelemetry() // nil emitter: must not panic
}

// `telemetry: off` in anchors.yaml turns it off for the project (run from its root).
func TestNoticeTelemetry_projectFileOff(t *testing.T) {
	isolateTelemetry(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "anchors.yaml"), "telemetry: off\n")
	t.Chdir(root)

	out := captureStderr(t, func() { NoticeTelemetry(cmdWithRoot("")) })
	if out != "" || Emitter != nil {
		t.Errorf("the project declared `telemetry: off`, and still: notice=%q emitter=%v", out, Emitter)
	}
}

// The project's `telemetry: off` holds from a subdirectory and with `--root` too: the
// config is read from the project root, not from the cwd.
func TestNoticeTelemetry_projectFileOffFromElsewhere(t *testing.T) {
	isolateTelemetry(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "anchors.yaml"), "telemetry: off\n")
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Chdir(filepath.Join(root, "sub"))
	if out := captureStderr(t, func() { NoticeTelemetry(cmdWithRoot("")) }); out != "" || Emitter != nil {
		t.Errorf("from a subdirectory the project's opt-out was ignored: notice=%q emitter=%v", out, Emitter)
	}

	t.Chdir(t.TempDir())
	if out := captureStderr(t, func() { NoticeTelemetry(cmdWithRoot(root)) }); out != "" || Emitter != nil {
		t.Errorf("with --root the project's opt-out was ignored: notice=%q emitter=%v", out, Emitter)
	}
}

func TestTelemetryHeaders_onlyFromTheEnvironment(t *testing.T) {
	t.Setenv("ANCHORS_TELEMETRY_KEY", "")
	if h := telemetryHeaders(); h != nil {
		t.Errorf("no key, no headers: %v", h)
	}
	t.Setenv("ANCHORS_TELEMETRY_KEY", "secret")
	if h := telemetryHeaders(); h["x-honeycomb-team"] != "secret" || len(h) != 1 {
		t.Errorf("the key did not become the auth header: %v", h)
	}
}

// An explicit --root is taken as given; with none, the root is found walking up to the
// directory that has anchors.yaml.
func TestProjectRoot(t *testing.T) {
	root := t.TempDir()
	if got := ProjectRoot(cmdWithRoot(root)); got != root {
		t.Errorf("explicit --root: got %q, want %q", got, root)
	}

	proj := t.TempDir()
	writeFile(t, filepath.Join(proj, "anchors.yaml"), "version: 2\n")
	sub := filepath.Join(proj, "pkg", "a")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(sub)
	want, _ := filepath.EvalSymlinks(proj)
	got, _ := filepath.EvalSymlinks(ProjectRoot(&cobra.Command{Use: "noflag"}))
	if got != want {
		t.Errorf("no --root, from a subdirectory: got %q, want %q", got, want)
	}
}
