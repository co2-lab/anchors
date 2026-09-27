package gate

import (
	"reflect"
	"testing"
	"time"
)

func TestAggregate_countsVerdictsPerGate(t *testing.T) {
	t.Run("PRFLO-B01: Each gate gets a summary counting its verdicts", func(t *testing.T) {})
	t.Run("PRFLO-B03: Every failed result is listed as a failure", func(t *testing.T) {})
	t.Run("PRFLO-B06: Results awaiting judgement are listed as awaiting judgement", func(t *testing.T) {})
	results := []Result{
		{Gate: "a", Target: "x", Verdict: Pass, Blocking: true, Duracao: 2 * time.Millisecond},
		{Gate: "a", Target: "y", Verdict: Pass, Blocking: true, Duracao: 5 * time.Millisecond},
		{Gate: "a", Target: "z", Verdict: Fail, Blocking: true, Duracao: 3 * time.Millisecond},
		{Gate: "b", Target: "x", Verdict: Skip},
		{Gate: "b", Target: "y", Verdict: Skip},
		{Gate: "b", Target: "z", Verdict: Pending},
		{Gate: "b", Target: "w", Verdict: Pending},
		{Gate: "c", Target: "x", Verdict: Judge},
		{Gate: "c", Target: "y", Verdict: Judge},
	}
	p := Aggregate(results)

	want := map[string]GateSummary{
		"a": {Gate: "a", Blocking: true, Pass: 2, Fail: 1, Duracao: 10 * time.Millisecond, Pior: 5 * time.Millisecond},
		"b": {Gate: "b", Skip: 2, Pending: 2},
		"c": {Gate: "c", Judge: 2},
	}
	if !reflect.DeepEqual(p.ByGate, want) {
		t.Errorf("ByGate:\n got %+v\nwant %+v", p.ByGate, want)
	}
	if len(p.Failures) != 1 || p.Failures[0].Target != "z" {
		t.Errorf("every fail is a failure: %+v", p.Failures)
	}
	if len(p.Judged) != 2 {
		t.Errorf("both judge verdicts await judgement: %+v", p.Judged)
	}
	if p.Passed || len(p.Blocked) != 1 {
		t.Errorf("a blocking fail blocks promotion: passed=%v blocked=%+v", p.Passed, p.Blocked)
	}
	if got := p.GateNames(); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("GateNames is sorted: %v", got)
	}
}

func TestAggregate_worstIsTheMostExpensiveSingleRun(t *testing.T) {
	t.Run("PRFLO-B02: The summary carries the total time and the most expensive run", func(t *testing.T) {})
	p := Aggregate([]Result{
		{Gate: "a", Verdict: Pass, Duracao: 7 * time.Millisecond},
		{Gate: "a", Verdict: Pass, Duracao: 7 * time.Millisecond},
		{Gate: "a", Verdict: Pass, Duracao: 1 * time.Millisecond},
	})
	if s := p.ByGate["a"]; s.Pior != 7*time.Millisecond || s.Duracao != 15*time.Millisecond {
		t.Errorf("Pior=%v Duracao=%v", s.Pior, s.Duracao)
	}
}

func TestAggregate_whatBlocksPromotion(t *testing.T) {
	t.Run("PRFLO-B04: Only a failure of a blocking gate blocks promotion", func(t *testing.T) {})
	t.Run("PRFLO-B05: A pending verdict blocks only when its gate blocks and it impedes", func(t *testing.T) {})
	t.Run("PRFLO-I01: Promotion is refused exactly when something blocks", func(t *testing.T) {})
	t.Run("PRFLO-X01: The profile does not decide whether a pending impedes", func(t *testing.T) {})
	cases := []struct {
		name   string
		r      Result
		passed bool
	}{
		{"a non-blocking fail does not block", Result{Gate: "a", Verdict: Fail}, true},
		{"a blocking fail blocks", Result{Gate: "a", Verdict: Fail, Blocking: true}, false},
		{"a blocking pending that does not impede does not block", Result{Gate: "a", Verdict: Pending, Blocking: true}, true},
		{"a pending that impedes on a non-blocking gate does not block", Result{Gate: "a", Verdict: Pending, Impede: true}, true},
		{"a blocking pending that impedes blocks", Result{Gate: "a", Verdict: Pending, Blocking: true, Impede: true}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := Aggregate([]Result{c.r})
			if p.Passed != c.passed || (len(p.Blocked) == 0) != c.passed {
				t.Errorf("passed=%v blocked=%+v, want passed=%v", p.Passed, p.Blocked, c.passed)
			}
		})
	}
}

func TestProfile_NodeVerdicts(t *testing.T) {
	t.Run("PRFLO-B07: Per-node verdicts include only confronted nodes, sorted", func(t *testing.T) {})
	t.Run("PRFLO-B08: A node is failed when a blocking gate failed on it or left it an impeding pending", func(t *testing.T) {})
	p := Aggregate([]Result{
		{Gate: "a", Target: "c.go", Verdict: Pass, Blocking: true},
		{Gate: "b", Target: "c.go", Verdict: Fail, Blocking: true},
		{Gate: "a", Target: "a.go", Verdict: Fail}, // non-blocking fail: confronted, not failed
		{Gate: "a", Target: "b.go", Verdict: Pass},
		{Gate: "a", Target: "skip.go", Verdict: Skip},
		{Gate: "a", Target: "judge.go", Verdict: Judge},
		{Gate: "a", Target: "pend.go", Verdict: Pending, Blocking: true},
		// A pending that impedes blocks promotion (PRFLO-B05), so the node is failed too.
		{Gate: "a", Target: "impede.go", Verdict: Pending, Blocking: true, Impede: true},
		{Gate: "a", Target: "soft.go", Verdict: Pending, Impede: true}, // non-blocking gate
	})
	got := p.NodeVerdicts()
	want := []NodeVerdict{
		{ID: "a.go"},
		{ID: "b.go"},
		{ID: "c.go", Failed: true},
		{ID: "impede.go", Failed: true},
		{ID: "pend.go"},
		{ID: "soft.go"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NodeVerdicts:\n got %+v\nwant %+v", got, want)
	}
}
