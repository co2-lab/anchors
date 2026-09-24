package initx

import (
	"encoding/json"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// THE AGENT FILTER answers "who is working right now?": one chip per agent that owns an
// open card touched in the last 30 minutes. These tests run the page's JS in `node`, like
// the blocked-cards strip, because the rule lives in the page — both producers (the
// pipeline and `board serve`) feed it the same fields.

type agenteAtivo struct {
	Nome  string `json:"nome"`
	Cards int    `json:"cards"`
}

func rodaAgentesAtivos(t *testing.T, montaDado string) []agenteAtivo {
	t.Helper()
	out := rodaJSDaPagina(t, `
const d = (() => { `+montaDado+` })();
console.log(JSON.stringify(agentesAtivos(d, referenciaDaAtividade(d))));
`)
	var r []agenteAtivo
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("node output is not the agent list: %v\n%s", err, out)
	}
	return r
}

// The window counts from the SNAPSHOT (`takenAt`), not from the reader's clock: the
// published board can be hours old. The data below is from January; counting from now
// would leave the list empty.
//
// Each card pins one edge of the rule:
//
//	alice  20 min before the snapshot  → active; her OLD second card still counts
//	erin   exactly 30 min before        → active (the edge is inclusive)
//	bob    31 min before                → out of the window
//	(liberado) card touched 1 min ago   → owner is "" — a released card has no agent
//	dave   closed 1 min ago             → a closed card is not work in progress
func TestActiveAgents_windowCountsFromTheSnapshot(t *testing.T) {
	got := rodaAgentesAtivos(t, `return {
	  takenAt: "2026-01-01T12:00:00Z",
	  items: [
	    { number: 1, owner: "alice", updated: "2026-01-01T11:40:00Z" },
	    { number: 2, owner: "alice", updated: "2025-12-01T00:00:00Z" },
	    { number: 3, owner: "erin",  updated: "2026-01-01T11:30:00Z" },
	    { number: 4, owner: "bob",   updated: "2026-01-01T11:29:00Z" },
	    { number: 5, owner: "",      updated: "2026-01-01T11:59:00Z",
	      ownership: ["alice", "(liberado)"] },
	    { number: 6, owner: "dave",  updated: "2026-01-01T11:59:00Z", closed: "2026-01-01T11:59:00Z" },
	  ],
	};`)
	want := []agenteAtivo{{"alice", 2}, {"erin", 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("active agents = %+v, want %+v\n"+
			"(bob is 31 min old, the released card has no owner, dave's card is closed)", got, want)
	}
}

// `board serve` is LIVE: its stamp is the time of a read the floor keeps handing back, so
// the window counts from now. The stamp here is a year old on purpose — counting from it
// would put an agent idle for a year on the list.
func TestActiveAgents_liveBoardCountsFromNow(t *testing.T) {
	got := rodaAgentesAtivos(t, `
	const ha = (min) => new Date(Date.now() - min * 60000).toISOString();
	return {
	  live: true,
	  takenAt: new Date(Date.now() - 365 * 864e5).toISOString(),
	  items: [
	    { number: 1, owner: "alice", updated: ha(5) },
	    { number: 2, owner: "bob",   updated: ha(45) },
	    { number: 3, owner: "idle",  updated: new Date(Date.now() - 365 * 864e5 - 60000).toISOString() },
	  ],
	};`)
	want := []agenteAtivo{{"alice", 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("live board: active agents = %+v, want %+v", got, want)
	}
}

// Clicking a chip shows only that agent's cards; the chip row shows the active agents with
// their counts and an "all" chip, or a quiet line when nobody is active.
func TestAgentChips_filterAndRow(t *testing.T) {
	out := rodaJSDaPagina(t, `
const d = {
  takenAt: "2026-01-01T12:00:00Z",
  items: [
    { number: 1, owner: "alice", updated: "2026-01-01T11:50:00Z", labels: ["bug"] },
    { number: 2, owner: "bob",   updated: "2026-01-01T11:55:00Z", labels: ["bug"] },
    { number: 3, owner: "carol", updated: "2025-01-01T00:00:00Z", labels: [] },
  ],
};
const linha = mk("agentes-chips");
rotulasAgentes(d);
const comAgentes = linha.innerHTML;
agenteAtivo = "alice";
const passa = d.items.map((i) => passaNoFiltro(i));
agenteAtivo = "";
const semFiltro = d.items.map((i) => passaNoFiltro(i));
rotulasAgentes({ takenAt: "2026-01-01T12:00:00Z", items: [] });
console.log(JSON.stringify({ comAgentes, passa, semFiltro, vazio: linha.innerHTML }));
`)
	var r struct {
		ComAgentes string `json:"comAgentes"`
		Passa      []bool `json:"passa"`
		SemFiltro  []bool `json:"semFiltro"`
		Vazio      string `json:"vazio"`
	}
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("node output: %v\n%s", err, out)
	}
	if !reflect.DeepEqual(r.Passa, []bool{true, false, false}) {
		t.Errorf("with alice selected, passaNoFiltro = %v, want only alice's card", r.Passa)
	}
	if !reflect.DeepEqual(r.SemFiltro, []bool{true, true, true}) {
		t.Errorf("with no agent selected, passaNoFiltro = %v, want every card", r.SemFiltro)
	}
	for _, s := range []string{`data-agente="alice"`, `data-agente="bob"`, `data-agente=""`} {
		if !strings.Contains(r.ComAgentes, s) {
			t.Errorf("the chip row lacks %s:\n%s", s, r.ComAgentes)
		}
	}
	if strings.Contains(r.ComAgentes, "carol") {
		t.Error("carol's only card is a year old, and she got a chip")
	}
	if strings.Index(r.ComAgentes, `"alice"`) > strings.Index(r.ComAgentes, `"bob"`) {
		t.Error("the chips are not sorted by name")
	}
	if !strings.Contains(r.Vazio, "nenhum agente ativo") {
		t.Errorf("with nobody active the row should say so quietly, got:\n%s", r.Vazio)
	}
}

// THE RELEASED OWNER IS NOT AN AGENT, and the rule is the producer's: the collect jq —
// the same one `board serve` runs — turns a last `anchors-owner: (liberado)` into an empty
// owner, and the page skips empty owners. Run the jq for real on `gh issue list` shaped
// input.
func TestBoardCollectJQReleasedOwnerIsNoAgent(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq missing: the collect expression cannot be run on this machine")
	}
	expr, err := BoardCollectJQ()
	if err != nil {
		t.Fatal(err)
	}
	in := `[
	  {"number":1,"title":"[AAAAA] a","url":"u","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}],
	   "updatedAt":"2026-01-01T11:59:00Z","createdAt":"2026-01-01T00:00:00Z","closedAt":null,
	   "author":{"login":"x"},"body":"",
	   "comments":[{"body":"anchors-owner: alice","createdAt":"2026-01-01T10:00:00Z"},
	               {"body":"anchors-owner: (liberado) — parado há 2h","createdAt":"2026-01-01T11:59:00Z"}]},
	  {"number":2,"title":"[BBBBB] b","url":"u","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}],
	   "updatedAt":"2026-01-01T11:59:00Z","createdAt":"2026-01-01T00:00:00Z","closedAt":null,
	   "author":{"login":"x"},"body":"",
	   "comments":[{"body":"anchors-owner: bob\nrelatório longo","createdAt":"2026-01-01T11:00:00Z"}]}
	]`
	cmd := exec.Command("jq", "-c", expr)
	cmd.Stdin = strings.NewReader(in)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("jq: %v\n%s", err, out)
	}
	var itens []struct {
		Number  int    `json:"number"`
		Owner   string `json:"owner"`
		Updated string `json:"updated"`
	}
	if err := json.Unmarshal(out, &itens); err != nil {
		t.Fatalf("jq output: %v\n%s", err, out)
	}
	if len(itens) != 2 {
		t.Fatalf("got %d items, want 2:\n%s", len(itens), out)
	}
	if itens[0].Owner != "" {
		t.Errorf("a card whose last `anchors-owner:` is `(liberado)` has owner %q — the "+
			"agent filter would list whoever just released it as working", itens[0].Owner)
	}
	if itens[1].Owner != "bob" {
		t.Errorf("owner = %q, want bob (first line of the comment only)", itens[1].Owner)
	}
	// `updated` is the activity signal the page reads; it must reach the JSON.
	if itens[0].Updated != "2026-01-01T11:59:00Z" {
		t.Errorf("updated = %q, want the issue's updatedAt", itens[0].Updated)
	}
}
