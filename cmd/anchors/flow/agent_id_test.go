package flow

import (
	"io"
	"os"
	"strings"
	"testing"
)

func stderrOf(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = orig
	b, _ := io.ReadAll(r)
	return string(b)
}

// Without ANCHORS_SESSION the identity falls back to the OS user, and two agents of that
// user collide (blue-eyes #650). The fallback is announced; a declared session is not.
func TestAgentID_fallbackIsAnnounced(t *testing.T) {
	t.Setenv("ANCHORS_SESSION", "")
	t.Setenv("USER", "alice")
	var id string
	out := stderrOf(t, func() { id = agentID() })
	if !strings.HasSuffix(id, "/alice") {
		t.Fatalf("fallback identity = %q, want the OS user", id)
	}
	if !strings.Contains(out, "ANCHORS_SESSION is not set") {
		t.Errorf("the fallback identity was used silently:\n%q", out)
	}

	t.Setenv("ANCHORS_SESSION", "devA")
	out = stderrOf(t, func() { id = agentID() })
	if !strings.HasSuffix(id, "/devA") || out != "" {
		t.Errorf("a declared session: id=%q, stderr=%q (want no warning)", id, out)
	}
}
