// @anchors
//   code: TLTST
//   ref: TLSTT

package common

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"

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
	return testkit.CaptureStderr(t, f)
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
	t.Run("TLSTT-B03: The notice is shown once, then telemetry runs quietly", func(t *testing.T) {})
	t.Run("TLSTT-I01: No emitter without the notice", func(t *testing.T) {})
	isolateTelemetry(t)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "anchors.yaml"), "version: 4\n")
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
	t.Run("TLSTT-I01: No emitter without the notice", func(t *testing.T) {})
	t.Run("TLSTT-B02: The environment opt-out builds nothing", func(t *testing.T) {})
	t.Run("TLSTT-B06: Waiting at exit with no emitter returns at once", func(t *testing.T) {})
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
	t.Run("TLSTT-B01: The project opt-out is read from the project root", func(t *testing.T) {})
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
	t.Run("TLSTT-B01: The project opt-out is read from the project root", func(t *testing.T) {})
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

// Starting telemetry builds the emitter and sends nothing: the endpoint is a local server
// that counts requests, and after the start and the exit wait it has received none.
func TestNoticeTelemetry_startingSendsNothing(t *testing.T) {
	t.Run("TLSTT-X01: Starting telemetry sends nothing", func(t *testing.T) {})
	isolateTelemetry(t)
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()
	t.Setenv("ANCHORS_TELEMETRY_ENDPOINT", srv.URL)

	proj := t.TempDir()
	writeFile(t, filepath.Join(proj, "anchors.yaml"), "version: 4\n")
	captureStderr(t, func() { NoticeTelemetry(cmdWithRoot(proj)) })
	if Emitter == nil {
		t.Fatal("telemetry on, and the emitter was not built")
	}
	FlushTelemetry()
	srv.Close() // waits for any request still being served
	if hits != 0 {
		t.Errorf("starting telemetry sent %d request(s)", hits)
	}
}

func TestTelemetryHeaders_onlyFromTheEnvironment(t *testing.T) {
	t.Run("TLSTT-B04: The authentication header comes only from the environment", func(t *testing.T) {})
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
	t.Run("TLSTT-B05: The project root is the explicit root or the nearest configured directory above", func(t *testing.T) {})
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

// A configuration that does not load (an unknown key, a format to migrate) was skipped, and
// the project's `telemetry: off` with it: the notice showed and the emitter was built for
// a project that had opted out.
func TestNoticeTelemetry_optOutHoldsWhenTheConfigDoesNotLoad(t *testing.T) {
	t.Run("TLSTT-B07: A declared opt-out holds when the configuration does not load", func(t *testing.T) {
		isolateTelemetry(t)
		root := t.TempDir()
		writeFile(t, filepath.Join(root, "anchors.yaml"), "telemetry: off\nnot_a_key_anchors_knows: 1\n")
		out := captureStderr(t, func() { NoticeTelemetry(cmdWithRoot(root)) })
		if out != "" || Emitter != nil {
			t.Errorf("the declared opt-out was ignored: notice=%q emitter=%v", out, Emitter)
		}
	})
}

// Outside a project the root fell back to the working directory, and the notice's marker
// was written there — a `.anchors/` left in whatever directory the command ran from.
func TestNoticeTelemetry_outsideAProjectWritesNothing(t *testing.T) {
	t.Run("TLSTT-B08: Outside a project there is no notice, no emitter and no mark", func(t *testing.T) {
		isolateTelemetry(t)
		dir := t.TempDir()
		t.Chdir(dir)
		if got := ProjectRoot(&cobra.Command{Use: "noflag"}); got != "" {
			t.Errorf("outside a project the root is empty, got %q", got)
		}
		for _, cmd := range []*cobra.Command{cmdWithRoot(""), cmdWithRoot(dir)} {
			out := captureStderr(t, func() { NoticeTelemetry(cmd) })
			if out != "" || Emitter != nil {
				t.Errorf("outside a project: notice=%q emitter=%v", out, Emitter)
			}
		}
		if _, err := os.Stat(filepath.Join(dir, ".anchors")); err == nil {
			t.Error("outside a project a `.anchors/` mark was written in the working directory")
		}
	})
}
