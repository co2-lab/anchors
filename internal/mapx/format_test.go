package mapx

import (
	"strings"
	"testing"
)

// THE CONTRACT: the map's `version:` says in which FORMAT the file is written, and the binary
// refuses what it cannot read.
//
// Without it, a binary that meets an unknown format interprets what it recognises, ignores the
// rest, and the next save writes only what survived — losing data with nothing reported. The
// concrete example is the renamed `judgments` key: the AI judgment stamps evaporate and `check`
// asks again what someone already answered.

func TestCheckFormat_acceptsOnlyTheRangeItReads(t *testing.T) {
	t.Run("MPFRM-B01: The written format and the oldest readable format are both accepted", func(t *testing.T) {})
	if err := ConfereFormato("m.yaml", FormatoAtual); err != nil {
		t.Errorf("the format this binary writes must be readable: %v", err)
	}
	if err := ConfereFormato("m.yaml", FormatoMinimoLegivel); err != nil {
		t.Errorf("the oldest readable format must be accepted: %v", err)
	}
}

// Every format from below the range to above it gets the answer the range dictates — the
// endpoints alone would leave an off-by-one inside or outside unnoticed.
func TestCheckFormat_everyFormatAgainstTheRange(t *testing.T) {
	t.Run("MPFRM-I01: Exactly format 6 is readable", func(t *testing.T) {})
	if FormatoAtual != 6 || FormatoMinimoLegivel != 6 {
		t.Fatalf("the spec states the binary writes 6 and reads from 6; got %d and %d", FormatoAtual, FormatoMinimoLegivel)
	}
	for v := -1; v <= FormatoAtual+2; v++ {
		err := ConfereFormato("m.yaml", v)
		readable := v == 6
		if readable && err != nil {
			t.Errorf("format %d is inside the range and was refused: %v", v, err)
		}
		if !readable && err == nil {
			t.Errorf("format %d is outside the range and was accepted", v)
		}
	}
}

// A map from the FUTURE asks for a newer binary, and the message must state the RISK —
// otherwise the cheap way out is to delete the file and rebuild it, which loses the judgments.
func TestCheckFormat_futureFormatAsksForUpgrade(t *testing.T) {
	t.Run("MPFRM-B02: A map from a newer binary is refused with the upgrade message", func(t *testing.T) {})
	err := ConfereFormato("m.yaml", FormatoAtual+1)
	if err == nil {
		t.Fatal("a format above the one the binary writes must be refused")
	}
	msg := err.Error()
	for _, want := range []string{"NEWER", "in silence", "judgment stamps", "upgrade"} {
		if !strings.Contains(msg, want) {
			t.Errorf("the message should contain %q; got:\n%s", want, msg)
		}
	}
}

// A map from the PAST asks for MIGRATION, and the message names the command. An error that
// says "run X" when X does not exist is worse than none: the reader tries, fails, and distrusts
// the next message.
func TestCheckFormat_oldFormatAsksForMigration(t *testing.T) {
	t.Run("MPFRM-B03: A map older than the readable range asks for migration", func(t *testing.T) {})
	t.Run("MPFRM-X01: Format 1 is migrated, not read", func(t *testing.T) {})
	err := ConfereFormato("m.yaml", 1)
	if err == nil {
		t.Fatal("format 1 must be refused — it is migrated, not read")
	}
	if !strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("the message should name the command that fixes it; got:\n%s", err)
	}
}

// A map WITHOUT `version:` predates the field — format 1, and migratable. Treating it as "0"
// and refusing with the future message would send the person to upgrade a binary that is
// already the newest.
func TestCheckFormat_noVersionIsFormatOne(t *testing.T) {
	t.Run("MPFRM-B04: A map with no version is format 1 and asks for migration", func(t *testing.T) {})
	err := ConfereFormato("m.yaml", 0)
	if err == nil {
		t.Fatal("with no `version:` the map is format 1, and needs migrating")
	}
	if !strings.Contains(err.Error(), "migrate") {
		t.Errorf("it should ask for MIGRATION, not an upgrade; got:\n%s", err)
	}
	// The message must say FORMAT 1, not "format 0".
	//
	// Without the normalisation, 0 still falls in the "below the minimum" range and produces
	// the right message for the wrong reason — and whoever reads "is in format 0" looks for a
	// corrupted file, not an old one. The mutation that removes the normalisation passed the
	// first version of this test.
	if !strings.Contains(err.Error(), "format 1") {
		t.Errorf("the message should say FORMAT 1 (the file is old, not corrupted); got:\n%s", err)
	}
	if strings.Contains(err.Error(), "format 0") {
		t.Error("`format 0` does not exist: a map with no `version:` is format 1")
	}
}

// The two messages are DISTINCT, and the distinction is what sends the person the right way:
// the future asks for a new binary, the past asks for migration. One "error loading the map"
// for both would send them looking for corruption where there is only a version.
func TestCheckFormat_theTwoMessagesDoNotMix(t *testing.T) {
	t.Run("MPFRM-B05: The newer-map refusal never names the migration command", func(t *testing.T) {})
	future := ConfereFormato("m.yaml", FormatoAtual+1).Error()
	past := ConfereFormato("m.yaml", 1).Error()

	if strings.Contains(future, "anchors migrate") {
		t.Error("a map from the future is not fixed by migrating — it asks for a new binary")
	}
	if strings.Contains(past, "NEWER") {
		t.Error("the old map was not written by a newer binary")
	}
}
