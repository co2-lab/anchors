package telemetry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The NOTICE appears ONCE. Repeating it would be worse than not showing it: whoever reads the
// same thing every time stops reading, and the notice stops doing what justifies the opt-out.
func TestNotice_oncePerProjectRoot(t *testing.T) {
	t.Run("TLNTT-B01: The notice is shown once per project root", func(t *testing.T) {})
	dir := t.TempDir()
	var buf strings.Builder

	Notice(&buf, dir)
	if buf.Len() == 0 {
		t.Fatal("the first run must show the notice")
	}

	buf.Reset()
	Notice(&buf, dir)
	if buf.Len() != 0 {
		t.Errorf("the second run should not show the notice; it wrote %d bytes", buf.Len())
	}
}

// The TEXT says three things, and the third makes the opt-out honest: whoever reads "we
// collect data" and does not find HOW TO TURN IT OFF on the same screen assumes the worst.
func TestNotice_saysWhatItSendsWhatItDoesNotAndHowToTurnItOff(t *testing.T) {
	t.Run("TLNTT-B02: The text says what is sent, what is not, and how to turn it off", func(t *testing.T) {})
	var buf strings.Builder
	Notice(&buf, t.TempDir())
	// The comparison normalizes case: the notice writes DECISÃO in upper case for emphasis,
	// and tying the assertion to the exact form would test the typography, not the information.
	text := strings.ToLower(buf.String())

	for _, want := range []string{
		"decisão",               // what it collects
		"não envia",             // what it does not collect
		"conteúdo de arquivo",   // the most likely fear of whoever reads it
		"anchors_telemetry=off", // how to turn it off, without editing a file
		"telemetry: off",
	} {
		if !strings.Contains(text, strings.ToLower(want)) {
			t.Errorf("the notice should contain %q", want)
		}
	}
}

func TestNotice_theMarkerStaysOutOfGit(t *testing.T) {
	t.Run("TLNTT-B03: The marker is written under the project's unversioned anchors directory", func(t *testing.T) {})
	if !strings.HasPrefix(noticeFile, ".anchors/") {
		t.Errorf("the marker must live in `.anchors/` (not versioned); it is at %q", noticeFile)
	}
	dir := t.TempDir()
	Notice(&strings.Builder{}, dir)
	if _, err := os.Stat(filepath.Join(dir, ".anchors", "telemetry-noticed")); err != nil {
		t.Errorf("the marker should exist under the root's .anchors/: %v", err)
	}
}

func TestNotice_writtenMeansAlreadyNoticed(t *testing.T) {
	t.Run("TLNTT-I01: After the notice is written it counts as already shown", func(t *testing.T) {})
	dir := t.TempDir()
	if AlreadyNoticed(dir) {
		t.Fatal("a fresh root has not seen the notice")
	}
	Notice(&strings.Builder{}, dir)
	if !AlreadyNoticed(dir) {
		t.Error("after the notice was written, the root should count as noticed")
	}
}

func TestNotice_doesNotConsultTheOptOut(t *testing.T) {
	t.Run("TLNTT-X01: The notice does not consult the opt-out itself", func(t *testing.T) {})
	t.Setenv(EnvVar, "off")
	var buf strings.Builder
	Notice(&buf, t.TempDir())
	if buf.Len() == 0 {
		t.Error("the notice announces; whether to show it is the caller's decision")
	}
}

func TestNotice_unwritableMarkerNeverBlocks(t *testing.T) {
	t.Run("TLNTT-E01: An unwritable marker never blocks and the notice shows again", func(t *testing.T) {})
	dir := t.TempDir()
	// `.anchors` is a FILE, so the directory for the marker cannot be created.
	if err := os.WriteFile(filepath.Join(dir, ".anchors"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		var buf strings.Builder
		Notice(&buf, dir)
		if buf.Len() == 0 {
			t.Fatalf("call %d: the notice should be written even without a marker", i+1)
		}
	}
}
