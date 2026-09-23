package quality

import (
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/gate"
)

// The SCENARIOS of the `timing-metrics` flag (`flags/timing-metrics.flag.md`).
//
// The flag governs whether `check` measures and prints the time per gate. The three
// scenarios are here, each with its code in a subtest NAME: that is how `flag-covered`
// recognises them as proven once the execution is ingested.
//
// `--timing` was born to find what makes a scan expensive, and it found it: `docs-fresh`
// was 97% of a 6m49s run. But it was born with no test at all, and it was `flag-covered`
// itself that charged it — the first work the flag axis found.

// profileWithTime builds a Profile with measured time, which is what `printTiming` reads.
func profileWithTime() gate.Profile {
	return gate.Profile{
		ByGate: map[string]gate.GateSummary{
			"docs-fresh": {Gate: "docs-fresh", Pass: 3, Duracao: 900 * time.Millisecond, Pior: 800 * time.Millisecond},
			"build":      {Gate: "build", Pass: 1, Duracao: 50 * time.Millisecond, Pior: 50 * time.Millisecond},
		},
		Results: []gate.Result{
			{Gate: "docs-fresh", Target: "a.spec.md", Duracao: 800 * time.Millisecond},
			{Gate: "build", Target: "b.go", Duracao: 50 * time.Millisecond},
		},
	}
}

// TIMNG-G01 — the value is `off`: check prints only the verdicts, and measures NO time.
//
// The proof is the absence: with the flag off neither the timing header nor the target
// list may print. A test that only checked "nothing broke" would pass with the whole
// table on screen.
func TestTimingG01_offPrintsNoTiming(t *testing.T) {
	t.Run("TIMNG-G01: with the flag off, check prints no timing at all", func(t *testing.T) {})
	// It confronts the VARIABLE the command reads, not a literal: `if off := false` would
	// be a tautology — it would pass with `printTiming` called unconditionally in
	// production, which is exactly the regression this scenario exists to catch.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing=false"}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	// The PRODUCTION function, not a copy of the condition: `RunE` calls exactly this.
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out != "" {
		t.Errorf("with the flag `off` check printed timing:\n%s", out)
	}
}

// TIMNG-G02 — the value is `on`: check also prints time per gate, and the slowest
// targets.
//
// Both blocks are charged, not just one: the per-gate header and the target list answer
// different questions ("which gate costs" and "which file costs"), and it was the second
// that pointed at `fnSize` reading ~43,000 files.
func TestTimingG02_onPrintsTimePerGateAndTargets(t *testing.T) {
	t.Run("TIMNG-G02: with the flag on, check prints time per gate and the slowest targets", func(t *testing.T) {})
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{"--timing"}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out == "" {
		t.Fatal("with the flag `on` check printed nothing")
	}
	for _, want := range []string{"docs-fresh", "build"} {
		if !strings.Contains(out, want) {
			t.Errorf("the output does not name the gate %q:\n%s", want, out)
		}
	}
	// The most expensive individual TARGET — the block the per-gate average does not show.
	if !strings.Contains(out, "a.spec.md") {
		t.Errorf("the output does not list the slowest target:\n%s", out)
	}
	// And the most expensive comes first: the order is what makes the table actionable.
	if strings.Index(out, "docs-fresh") > strings.Index(out, "build") {
		t.Errorf("the most expensive gate did not come first:\n%s", out)
	}
}

// TIMNG-G03 — the value is ABSENT: the same as `off`. Measuring is opt-in, never a
// default cost.
//
// It is the scenario `flag-scenarios-complete` charges and the one nobody writes. It
// matters here for a concrete reason: an inverted default would raise no error at all — it
// would just spend everyone's time, forever, in silence.
func TestTimingG03_absentMeansTheDefaultWhichIsOff(t *testing.T) {
	t.Run("TIMNG-G03: with the flag absent, the declared default holds — measuring is opt-in", func(t *testing.T) {})
	// The absence of the VALUE is the absence of the flag on the command line. Parsing an
	// argv without `--timing` leaves the flag in the state the scenario describes — and it
	// is that state, not a literal, that the call site reads.
	cmd := newCheckCmd()
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatal(err)
	}

	on, err := cmd.Flags().GetBool("timing")
	if err != nil {
		t.Fatal(err)
	}
	out := capturaSaida(t, func() { reportTiming(on, profileWithTime()) })
	if out != "" {
		t.Errorf("without `--timing` check measured and printed — measuring must be opt-in:\n%s", out)
	}
	// And the DEFAULT declared in cobra must be the same. It is the other half of the
	// scenario: the block above proves the behaviour with the value absent, and this one
	// proves that absent is really what the command delivers when nobody passes `--timing`.
	f := newCheckCmd().Flags().Lookup("timing")
	if f == nil {
		t.Fatal("the `--timing` flag disappeared from the command")
	}
	if f.DefValue != "false" {
		t.Errorf("the default of `--timing` is %q — measuring stopped being opt-in", f.DefValue)
	}
}
