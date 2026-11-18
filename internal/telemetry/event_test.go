// @anchors
//   code: EVTSA
//   ref: TLEVT

package telemetry

import (
	"testing"
	"time"
)

func TestNew_keepsNameAndAttributes(t *testing.T) {
	t.Run("TLEVT-B01: A new event keeps its name and attributes", func(t *testing.T) {})
	ev := New(ClaimServed, map[string]any{"candidates": 14, "skipped": 2}, fixedClock)
	if ev.Name != "claim.served" {
		t.Errorf("name = %q", ev.Name)
	}
	if len(ev.Attrs) != 2 || ev.Attrs["candidates"] != 14 || ev.Attrs["skipped"] != 2 {
		t.Errorf("attributes changed: %v", ev.Attrs)
	}
}

func TestNew_withoutAttributesCarriesAnEmptySet(t *testing.T) {
	t.Run("TLEVT-B02: An event without attributes carries an empty set", func(t *testing.T) {})
	ev := New(CheckFinished, nil, fixedClock)
	if ev.Attrs == nil {
		t.Fatal("the attribute set must exist, even when empty")
	}
	if len(ev.Attrs) != 0 {
		t.Errorf("expected an empty set, got %v", ev.Attrs)
	}
}

func TestNew_theInstantComesFromTheCallersClock(t *testing.T) {
	t.Run("TLEVT-B03: The instant comes from the caller's clock", func(t *testing.T) {})
	ev := New(TurnEnded, nil, fixedClock)
	if !ev.At.Equal(time.Date(2026, 9, 14, 18, 30, 0, 0, time.UTC)) {
		t.Errorf("instant = %v", ev.At)
	}
}

func TestNames_theVocabularyIsClosed(t *testing.T) {
	t.Run("TLEVT-I01: The vocabulary is the five declared names", func(t *testing.T) {})
	want := map[Name]string{
		ClaimServed:     "claim.served",
		ClaimEmpty:      "claim.empty",
		CheckFinished:   "check.finished",
		EscalatedToUser: "escalate.raised",
		TurnEnded:       "turn.ended",
	}
	if len(want) != 5 {
		t.Fatalf("two names share a wire value: %v", want)
	}
	for n, s := range want {
		if string(n) != s {
			t.Errorf("%q should be %q", n, s)
		}
	}
}

func TestNew_neverReadsTheSystemClock(t *testing.T) {
	t.Run("TLEVT-X01: The system clock is never read", func(t *testing.T) {})
	past := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	ev := New(ClaimEmpty, nil, func() time.Time { return past })
	if !ev.At.Equal(past) {
		t.Errorf("instant = %v, expected the injected %v", ev.At, past)
	}
}
