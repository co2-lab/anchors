// @anchors
//   code: EBTXP
//   ref: BREXB

package initx

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// THE `board.json` CONTRACT IS SINGLE, and `board serve` reads the PIPELINE's expression
// instead of reimplementing it.
//
// Reimplementing would create a second version that diverges at the first new field: the
// `owner` enters the pipeline and the local board does not show it — or worse, shows it
// differently, and nobody knows which of the two is right.
func TestBoardCollectJQComesFromThePipeline(t *testing.T) {
	t.Run("BREXB-B02: The collect expression comes out of the pipeline", func(t *testing.T) {})
	t.Run("BREXB-I01: The extracted expression never carries the redirection", func(t *testing.T) {})
	jq, err := BoardCollectJQ()
	if err != nil {
		t.Fatalf("the expression was not extracted: %v", err)
	}
	// The fields the board needs — if the regex matched the wrong stretch, one goes missing.
	for _, field := range []string{"number:", "title:", "state:", "owner:", "ownership:"} {
		if !strings.Contains(jq, field) {
			t.Errorf("the extracted expression lacks `%s` — the regex matched the wrong stretch, "+
				"and the local board would serve an incomplete contract", field)
		}
	}
	// The discarded card leaves at the SOURCE, and that is a board rule: if it left the
	// board, it left the data. `serve` has to inherit that.
	if !strings.Contains(jq, "anchors:discarded") {
		t.Error("the expression does not filter `anchors:discarded` — the local board would show " +
			"a discarded card the published one does not")
	}
	if jq != strings.TrimSpace(jq) {
		t.Error("the expression should come back without surrounding whitespace")
	}

	// THE EXPRESSION MUST BE VALID JQ, and the previous ruler did not measure that.
	//
	// The first version of the regex closed on a non-greedy `'` and swallowed the
	// `> _board/board.json` that follows in the shell. All the fields were there — the
	// assertions above passed — and `jq` died with `unexpected token "'"` only when the
	// command really ran.
	//
	// It is the trap of checking the SHAPE and not what the shape is for. The `>` alone is
	// NOT a signal: it appears legitimately in jq's named capture (`(?<c>...)`). What gives
	// away the wrong cut is the REDIRECTION, which exists only in the shell.
	if strings.Contains(jq, "_board/board.json") {
		t.Errorf("the extracted expression carries the shell redirection — jq will "+
			"receive `>` as syntax and die:\n%s", lastChars(jq, 80))
	}
	if !strings.HasPrefix(jq, "[") || !strings.HasSuffix(jq, "]") {
		t.Errorf("the expression does not open and close the array — the regex cut in the wrong place:\n%s",
			lastChars(jq, 80))
	}
}

// Both shapes of the collect step are read, and the cut ends before the redirection. A
// step of any other shape yields nothing — never a built-in expression in its place.
func TestBoardCollectJQShapes(t *testing.T) {
	t.Run("BREXB-B03: Both collect shapes are accepted and cut before the redirection", func(t *testing.T) {})
	t.Run("BREXB-X01: A collect step of another shape yields no expression", func(t *testing.T) {})
	t.Run("BREXB-E02: A pipeline that changed shape fails naming the change", func(t *testing.T) {})
	for _, step := range []string{
		"gh issue list --jq '[.[] | {n: .number}]' > _board/board.json\n",
		"jq -s '[.[] | {n: .number}]' _board/raw.jsonl > _board/board.json\n",
	} {
		m := jqDaColeta.FindSubmatch([]byte(step))
		if m == nil || string(m[1]) != "[.[] | {n: .number}]" {
			t.Errorf("%q: expected the array expression, got %q", step, m)
		}
	}
	if m := jqDaColeta.FindSubmatch([]byte("gh issue list --jq '.[] | .x' > out\n")); m != nil {
		t.Errorf("a collect step of another shape must yield nothing, got %q", m[1])
	}
}

// The HTML is the SAME — not a copy that ages in parallel.
func TestBoardHTMLIsThePipelinePage(t *testing.T) {
	t.Run("BREXB-B01: The local board page is the published page", func(t *testing.T) {})
	h, err := BoardHTML()
	if err != nil {
		t.Fatal(err)
	}
	carried, err := boardFS.ReadFile("board/anchors-board.html")
	if err != nil || h != string(carried) {
		t.Fatalf("the page is not the carried board page (err %v)", err)
	}
	if !strings.Contains(h, "board.json") {
		t.Error("the HTML does not fetch `board.json` — `serve` would serve a page that reads nothing")
	}
	// The marks below belong to the board's STRUCTURE, not to a specific fix: if they go
	// away, `serve` started reading another file.
	for _, mark := range []string{"id=\"colunas\"", "id=\"detalhe\"", "id=\"v-roadmap\""} {
		if !strings.Contains(h, mark) {
			t.Errorf("the HTML lacks `%s` — `serve` is reading another file", mark)
		}
	}
}

// lastChars returns the end of a string, so the error message shows WHERE it cut.
func lastChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "…" + s[len(s)-n:]
}

// Without the embedded files, both readers fail loudly instead of serving an empty board.
func TestBoardReadersFailWithoutTheEmbeddedFiles(t *testing.T) {
	t.Run("BREXB-E01: A binary without the board files fails loudly", func(t *testing.T) {})
	savedBoard, savedWorkflows := boardFS, workflowsFS
	t.Cleanup(func() { boardFS, workflowsFS = savedBoard, savedWorkflows })
	boardFS, workflowsFS = embed.FS{}, embed.FS{}

	if h, err := BoardHTML(); err == nil || h != "" || !strings.Contains(err.Error(), "the board HTML is not embedded") {
		t.Errorf("BoardHTML = %q, %v; want the not-embedded error", h, err)
	}
	if jq, err := BoardCollectJQ(); err == nil || jq != "" || !strings.Contains(err.Error(), "the board pipeline is not embedded") {
		t.Errorf("BoardCollectJQ = %q, %v; want the not-embedded error", jq, err)
	}
}

// ── the page's own behaviour ─────────────────────────────────────────────────────────
//
// THE BLOCKED-CARDS STRIP is the one piece of the board page that decides something: the
// columns list what exists, and this strip judges what is stuck, who blocks how much, and in
// which order to answer.
//
// It was born because `anchors:needs-user` has no column: the card was not merely
// unhighlighted, it VANISHED from the page. It stayed open on GitHub blocking what depended
// on it, and the board showed a project flowing.
//
// The JS of an embedded page has no suite, and the consequence is known — four mutations of
// this behaviour would pass unnoticed. These tests run the real JS in `node` and confront
// what it produces. Without `node` on the machine they skip: the ruler belongs to CI and to
// whoever has the tool, and failing `go test` for whoever does not would be worse than the
// missing test.

//go:embed board
var boardPageForTest embed.FS

// domDouble is the minimum for the page to load outside a browser. A PROXY, and not an
// enumerated object: the page touches dozens of DOM methods while initializing, and patching
// them one by one at each error costs more than answering all of them. What the tests observe
// is the `hidden` and `innerHTML` of the element they target.
const domDouble = `
const els = {};
function mk(id) {
  if (els[id]) return els[id];
  const alvo = { id, hidden: false, innerHTML: '', textContent: '', dataset: {}, style: {} };
  return els[id] = new Proxy(alvo, {
    get(o, k) {
      if (k in o) return o[k];
      // 'querySelectorAll' returns a LIST. A generic method that returns another node makes
      // the page's '[...]' spread blow up, and the error points at the page's code instead of
      // the double — the symptom is "exit status 1" pointing at a line that is correct.
      if (k === 'querySelectorAll') return () => [];
      return typeof k === 'string' ? () => mk(id + ':' + k) : undefined;
    },
    set(o, k, v) { o[k] = v; return true; },
  });
}
const qualquer = new Proxy(function(){}, {
  get(_, k) { return k === 'matches' ? false : qualquer; },
  set() { return true; }, apply() { return qualquer; }, construct() { return qualquer; },
});
globalThis.document = new Proxy({ getElementById: mk }, { get(o, k) { return k in o ? o[k] : qualquer; } });
globalThis.window = qualquer;
globalThis.localStorage = { getItem: () => null, setItem() {} };
globalThis.fetch = () => new Promise(() => {});
globalThis.setInterval = () => 0;
globalThis.setTimeout = () => 0;
globalThis.matchMedia = () => ({ matches: false, addEventListener() {} });
`

var pageScriptRE = regexp.MustCompile(`(?s)<script>(.*?)</script>`)

// runDrawBlocked runs `desenhaBloqueio` over the given items and returns the state of the
// strip element.
func runDrawBlocked(t *testing.T, items []map[string]any) (hidden bool, html string) {
	return runOnPage(t, "bloqueio", "desenhaBloqueio", items)
}

// runOnPage runs one of the page's functions over the given items and returns the element it
// wrote. Generalized from `desenhaBloqueio` when the card's mark in the COLUMN needed the same
// apparatus — the DOM double and the script extraction are the same.
func runOnPage(t *testing.T, element, fn string, items []map[string]any) (hidden bool, html string) {
	t.Helper()

	data, err := json.Marshal(map[string]any{"items": items})
	if err != nil {
		t.Fatalf("data: %v", err)
	}
	out := runPageJS(t, `
const alvo = mk('`+element+`');
`+fn+`(`+string(data)+`);
console.log(JSON.stringify({hidden: alvo.hidden, html: alvo.innerHTML}));
`)
	var r struct {
		Hidden bool   `json:"hidden"`
		HTML   string `json:"html"`
	}
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("node output is not the expected JSON: %v\n%s", err, out)
	}
	return r.Hidden, r.HTML
}

// runPageJS loads the page's JS over the DOM double, runs `code` after it and returns the
// LAST line node printed. Split from `runOnPage` when the agent filter got a PURE function to
// confront — it returns a value instead of writing an element, and the apparatus that loads
// the page is the same.
func runPageJS(t *testing.T, code string) string {
	t.Helper()

	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node missing: the page's JS cannot be confronted on this machine")
	}

	page, err := fs.ReadFile(boardPageForTest, "board/anchors-board.html")
	if err != nil {
		t.Fatalf("board page: %v", err)
	}

	// The page has MORE THAN ONE `<script>` block, and a greedy regex would take the
	// `</script>` of the last one — the first mistake this test made.
	var js strings.Builder
	for _, m := range pageScriptRE.FindAllStringSubmatch(string(page), -1) {
		js.WriteString(m[1])
		js.WriteString("\n")
	}
	if js.Len() == 0 {
		t.Fatal("no <script> block in the page: the test has nothing to confront")
	}

	prog := domDouble + js.String() + "\n" + code
	file := filepath.Join(t.TempDir(), "board.mjs")
	if err := os.WriteFile(file, []byte(prog), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	// `CombinedOutput` and not `Output`: `Output` swallows stderr, and an error inside the JS
	// becomes "exit status 1" without saying where — ten minutes looking for what the message
	// would already have said.
	out, err := exec.Command("node", file).CombinedOutput()
	if err != nil {
		t.Fatalf("node: %v\n%s", err, out)
	}
	// The last line: the page may log before.
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines[len(lines)-1]
}

func TestBlockedStrip_hiddenWithNoEscalatedCard(t *testing.T) {
	t.Run("BREXB-B08: With no card waiting for a person the blocked strip is hidden and empty", func(t *testing.T) {})
	hidden, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:to-do", "number": 1, "code": "A", "title": "x", "url": "u"},
	})
	if !hidden {
		t.Error("with no `needs-user` card the strip should be hidden")
	}
	// Hiding without CLEARING would leave the old content ready to reappear the next time the
	// strip is shown — with decisions already resolved.
	if html != "" {
		t.Errorf("the hidden strip should be empty, and it has %q", html)
	}
}

func TestBlockedStrip_ordersByHowManyCardsEachBlocks(t *testing.T) {
	t.Run("BREXB-B09: The strip orders the waiting cards by how many cards each blocks", func(t *testing.T) {})
	// The LOW number blocks little, the HIGH one blocks a lot: the two possible orderings give
	// opposite results here. With data where they agree, the mutation that orders by number
	// survives — which is what happened in the first version of this test.
	//
	// And `DEC2` depends on `DEC1`: it proves the count does NOT count as blocked a card that
	// is itself stopped waiting for a person.
	hidden, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:needs-user", "number": 12, "code": "DEC1", "title": "[DEC1] qual vocabulário?", "url": "https://x/12"},
		{"state": "anchors:needs-user", "number": 99, "code": "DEC2", "title": "[DEC2] onde fica o limite?", "needs": []string{"DEC1"}, "url": "https://x/99"},
		{"state": "anchors:to-do", "number": 5, "code": "Z1", "title": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"state": "anchors:to-do", "number": 6, "code": "Z2", "title": "z", "needs": []string{"DEC2"}, "url": "u"},
		{"state": "anchors:to-do", "number": 7, "code": "Z3", "title": "z", "needs": []string{"DEC1"}, "url": "u"},
	})
	if hidden {
		t.Fatal("with two `needs-user` cards the strip should show")
	}

	i99, i12 := strings.Index(html, "#99"), strings.Index(html, "#12")
	if i99 < 0 || i12 < 0 {
		t.Fatalf("the strip should show BOTH numbers; html=%q", html)
	}
	if i99 > i12 {
		t.Error("the card that blocks more should come first: #99 blocks 2 and #12 blocks 1")
	}
	if !strings.Contains(html, "trava 2 cards") {
		t.Error("#99 blocks two cards, and the strip should say how many")
	}
	// `trava 1 card<` and not `trava 1 card`: without the delimiter, "trava 1 card" would match
	// inside "trava 1 cards" — which the plural rule should not produce.
	if !strings.Contains(html, "trava 1 card<") {
		t.Error("#12 blocks only Z3 (DEC2 is stopped too, and does not count as blocked)")
	}
	if !strings.Contains(html, "Esperando você · 2") {
		t.Error("the strip should say HOW MANY decisions are stopped")
	}
	// The title's `[CODE]` is noise here: the code already shows in the board's column, and in
	// the strip what matters is the question.
	if strings.Contains(html, "[DEC2]") {
		t.Error("the `[CODE]` prefix should leave the title in the strip")
	}
}

// `needs-user` ADDS UP with the work state, and the strip has to see both.
//
// A card in review that needs a decision carries BOTH labels, and the JSON's `state` is the
// FIRST of them. The first version filtered by `state === "anchors:needs-user"` and the card
// VANISHED from the strip — it stayed only in the IN REVIEW column, indistinguishable from a
// card being reviewed normally.
//
// Measured: #311 of the reference project was waiting for a decision of the user and did not
// show here. The user saw it on the board with nothing saying nobody would touch it.
func TestBlockedStrip_escalatedCardInAWorkStateAppears(t *testing.T) {
	t.Run("BREXB-B10: An escalated card in a work state appears in the strip with the state it stopped in", func(t *testing.T) {})
	hidden, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"state": "anchors:to-do", "escalated": true, "number": 443, "code": "DEC",
			"title": "[DEC] a revisão achou defeito crítico", "url": "https://x/443"},
	})
	if hidden {
		t.Fatal("two escalated cards and the strip did not show")
	}
	if !strings.Contains(html, "#311") {
		t.Error("the `in-review` + `needs-user` card has to show — it is the case that vanished")
	}
	if !strings.Contains(html, "#443") {
		t.Error("the `to-do` + `needs-user` card too")
	}
	if !strings.Contains(html, "Esperando você · 2") {
		t.Error("the count has to include both")
	}
	// And say WHERE the card stopped: `in-review` waiting for a decision differs from `to-do`
	// waiting for a decision — the first has work half done, the second does not.
	if !strings.Contains(html, "parado em in-review") {
		t.Error("the strip should say in which state the card stopped")
	}
	// `to-do` is not "stopped in": it is the normal state of a card not yet taken.
	if strings.Contains(html, "parado em to-do") {
		t.Error("`to-do` is not an interrupted work state")
	}
}

// THE CARD THAT WAITS FOR YOU differs from the one that waits for ALREADY DECIDED WORK.
//
// When the decision comes out and creates work, the change's card is born with
// `anchors:desbloqueia-<n>` and the blocked one keeps `needs-user` — because removing the
// label before the delivery would make the claim hand the card back, and the agent would hit
// the same impasse.
//
// The two cases look alike on screen unless the strip tells them apart: one asks the user's
// attention NOW, the other asks nothing — it only waits for a queue to move.
func TestBlockedStrip_tellsWaitingForYouFromWaitingForWork(t *testing.T) {
	t.Run("BREXB-B11: The strip tells a card waiting for a person from one waiting for decided work", func(t *testing.T) {})
	_, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "https://x/311"},
		{"state": "anchors:to-do", "escalated": true, "number": 500, "code": "DEC",
			"title": "[DEC] qual vocabulário?", "url": "https://x/500"},
		// The change's card: it is not escalated, and it UNBLOCKS 311.
		{"state": "anchors:to-do", "number": 444, "code": "FIX", "unblocks": "311",
			"title": "[destrava #311] try/catch por token", "url": "https://x/444"},
	})

	if !strings.Contains(html, "espera a entrega do #444") {
		t.Error("the card whose decision already came out should say WHICH card it waits for")
	}
	// #500 still waits for the PERSON — it cannot get the same sentence.
	i500 := strings.Index(html, "#500")
	i444 := strings.Index(html, "espera a entrega do #444")
	if i500 >= 0 && i444 >= 0 && i444 > i500 {
		t.Error("the waiting sentence stuck to the wrong card — it is #311's, not #500's")
	}
	// And the CHANGE's card does not enter the strip: it waits for no one, it is ordinary work.
	if strings.Contains(html, "#444 ") || strings.Contains(html, ">#444<") {
		t.Error("the unblocking card is not an escalated card — it should not show here")
	}
}

func TestBlockedStrip_escapesAHostileTitle(t *testing.T) {
	t.Run("BREXB-B12: A hostile title is escaped in the strip", func(t *testing.T) {})
	// The title comes from an issue, and anyone with access to the repository writes issues.
	// The page is served on the project's Pages — a title with markup would run in the browser
	// of whoever opens the board.
	_, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:needs-user", "number": 7, "code": "X",
			"title": "<img src=x onerror=alert(1)>", "url": "https://x/7"},
	})
	if strings.Contains(html, "<img src=x") {
		t.Errorf("the title should be escaped; html=%q", html)
	}
}

// THE ESCALATED CARD needs a mark IN THE COLUMN, not only in the strip.
//
// The strip lists what waits for a person; the column is where one looks to know what is
// happening with a specific card. Without a mark there, a card stopped waiting for a decision
// is visually IDENTICAL to the card next to it being worked on.
//
// Measured: the user looked for #311 in the column, found it, and asked "and the highlight on
// the card?" — the strip showed it, and the card said nothing.
func TestEscalatedCard_isMarkedInItsColumn(t *testing.T) {
	t.Run("BREXB-B13: An escalated card is marked in its column", func(t *testing.T) {})
	_, html := runOnPage(t, "colunas", "desenha", []map[string]any{
		{"state": "anchors:in-review", "escalated": true, "number": 311, "code": "NTDSN",
			"title": "[NTDSN] o que sai da VPC", "url": "u"},
		{"state": "anchors:in-review", "number": 312, "code": "OUTRO",
			"title": "[OUTRO] revisão normal", "url": "u"},
	})

	if !strings.Contains(html, `class="card escalado"`) {
		t.Error("the escalated card needs a class of its own — the border is what sets it " +
			"apart when scanning the column")
	}
	if !strings.Contains(html, "esperando você") {
		t.Error("the badge has to SAY what the card waits for: a colored border alone " +
			"requires whoever looks to already know what the color means")
	}
	// And the ordinary card must NOT get either — a mark that shows on everything marks
	// nothing.
	if n := strings.Count(html, `class="card escalado"`); n != 1 {
		t.Errorf("only the escalated card should have the class; it showed %d times", n)
	}
	if n := strings.Count(html, "esperando você"); n != 1 {
		t.Errorf("only the escalated card should have the badge; it showed %d times", n)
	}
}

// THE ROADMAP did not open the details, and the inconsistency was invisible.
//
// In the Board and Tree tabs a click on a card opens the modal — owner, ownership history,
// what the card asks. In the Roadmap it opened NOTHING: the label on the left had no
// `data-number` nor listener, and the bar was an `<a href>` that LEFT for GitHub.
//
// The cost: whoever wanted to see whose a card is had no way to, precisely in the tab that
// shows the work over time — where the question "who has this?" comes up the most.
//
// Reported by the user: "I am clicking the board's items and the modal does not open, I
// wanted to see the owners of the cards in the in-review column".
func TestRoadmapOpensTheDetails(t *testing.T) {
	t.Run("BREXB-B14: The roadmap opens a card's details", func(t *testing.T) {})
	b, err := fs.ReadFile(boardFS, "board/anchors-board.html")
	if err != nil {
		t.Fatal(err)
	}
	html := string(b)

	// The roadmap's container enters the click delegation.
	if !strings.Contains(html, `"colunas", "arvore", "v-roadmap"`) {
		t.Error("`v-roadmap` is not in the click delegation — the click on the roadmap does " +
			"not reach `abre()`, and the tab has no details")
	}

	// The selector recognizes the roadmap's label.
	if !strings.Contains(html, `.g-rot[data-number]`) {
		t.Error("the selector does not match the roadmap's label (`.g-rot[data-number]`) — " +
			"the listener exists and finds no target")
	}

	// And the label carries the number, otherwise there is nothing to match.
	if !strings.Contains(html, `data-number="${esc(i.number)}"`) {
		t.Error("the roadmap's label does not emit `data-number` — the selector has nothing to read")
	}
}

// THE BOARD SAYS WHICH CARD A CARD WAITS FOR, and not only that it waits.
//
// `needs-user` stops the card — and does not say WHO holds it. The board showed "waiting for
// you" with no link, and whoever answers the decision does not know what they just released:
// each card has to be found again by hand, and what is not found stays stopped after the
// decision already came out.
//
// MEASURED in the reference project: 23 open decisions were SEVEN questions, and unblocking
// the dependants required rereading card by card to find out who waited for what.
func TestBoardSaysWhichCardBlocksACard(t *testing.T) {
	t.Run("BREXB-B15: A blocked card says which card blocks it", func(t *testing.T) {})
	html, err := BoardHTML()
	if err != nil {
		t.Fatal(err)
	}

	// The collected FIELD has to reach the HTML: without it the badge never shows.
	if !strings.Contains(html, "blockedBy") {
		t.Error("the board does not read `blockedBy` — the link exists in the label and is not drawn")
	}

	// The NUMBER in the badge, and not only the word "blocked": the number is what answers
	// "waiting for what?", and without it the badge repeats what `needs-user` already said.
	if !strings.Contains(html, "bloqueado por #") {
		t.Error("the blocking badge does not show the NUMBER of the card that holds it — " +
			"without it the reader knows the card stopped and not where the answer goes")
	}

	// The CLICK goes to the holder. Without `data-vai-para` the click falls on the enclosing
	// card and opens the blocked one itself, which is what the reader is already looking at.
	if !strings.Contains(html, "data-vai-para") {
		t.Error("the badge does not lead to the card that holds it — the question of whoever " +
			"reads `blocked by #123` is what #123 asks")
	}

	// The COLOR has to exist in the THREE palettes: the board is read in the light theme, dark
	// by `prefers-color-scheme` and dark by `data-theme`. A color defined in only one block
	// becomes `inherit` in the others, and the badge loses precisely what sets it apart from
	// the escalated one.
	if n := strings.Count(html, "--bloqueio:"); n < 3 {
		t.Errorf("the blocking color is defined %d time(s), and the palettes are 3 "+
			"(light, prefers-color-scheme, data-theme) — in the missing themes the badge "+
			"loses the color that sets it apart from a card that is only escalated", n)
	}
}

// THE USER'S TWO QUEUES ARE APART ON THE BOARD, because their weights differ.
//
//	needs-user     "decide between A and B" — there is more than one defensible answer, and
//	               choosing between them changes what the product does
//	needs-framing  "check whether this is yours" — whoever escalated could not tell whether
//	               it changes the direction, and preferred declaring the doubt to asserting
//	               an impact they did not measure
//
// THE SECOND HAS A CHEAP WAY OUT: if it has no impact, swapping the label for `to-do` returns
// the card to the queue and nobody decides the merit. Mixing them makes the cheap one cost as
// much as the expensive one — whoever opens the list spends the effort of deciding before
// finding out they only had to return the card.
//
// THE CAUSE was the product's: the `escalate` help said "use `--for-user` (...) or if you are
// unsure whether it has an impact", inviting escalation for safety. The cost of that caution
// is invisible to whoever escalates.
//
// MEASURED in the reference project: 26 open decisions, and a triage concluded they were
// SEVEN real decisions.
func TestBoardSeparatesDecidingFromFraming(t *testing.T) {
	t.Run("BREXB-B16: Decisions and framing requests are listed apart", func(t *testing.T) {})
	html, err := BoardHTML()
	if err != nil {
		t.Fatal(err)
	}

	// The collected FIELD has to reach the HTML: without it the two queues become one.
	if !strings.Contains(html, "framing") {
		t.Fatal("the board does not read `framing` — the user's two queues show as one, and " +
			"the one that asks only framing costs as much as the one that asks a decision")
	}

	// THE TWO LISTS HAVE TO COME FROM THE SAME SET, each with one side of the field.
	//
	// The first version of this assertion looked for `p.framing` in the HTML and SURVIVED the
	// mutation that empties the framing list (`const enquadrar = []`) — because the string
	// still shows in the other list's filter. It measured the presence of the text, not the
	// separation.
	//
	// The two halves of the predicate are what cannot be lost: one list with the field true
	// and another with it false. Either one alone leaves a queue empty.
	for _, half := range []string{"pendentes.filter((p) => !p.framing)", "pendentes.filter((p) => p.framing)"} {
		if !strings.Contains(html, half) {
			t.Errorf("the board lost `%s` — without both halves of the filter one of the "+
				"queues is born empty, and its cards vanish from view", half)
		}
	}

	// And THE CHEAP WAY OUT said in the text: without it whoever reads does not know they can
	// return the card instead of deciding, which is the only reason the queue is apart.
	if !strings.Contains(html, "anchors:to-do") {
		t.Error("the framing queue does not say HOW to return the card to the queue — without " +
			"it whoever reads decides the merit for lack of a visible alternative")
	}
}

// THE RENDER AND THE HANDLER HAVE TO SPEAK THE SAME ATTRIBUTE.
//
// board.json became an English contract, the handler went with it — and the two card renders
// kept writing `data-numero`. The measured effect: 88 cards in the columns tab, ZERO matching
// `closest(".card, ...[data-number]")`, and the click opened nothing. The roadmap tab, which
// already wrote `data-number`, was the only live one.
//
// Nothing flagged it because the embedded HTML had no ruler at all over its own attributes:
// the render and the handler live in the same file and diverged inside it.
//
// This ruler compares both sides: every `data-*` the handler looks for has to exist in the
// render, and every `data-*` the render writes has to be looked for by someone — an attribute
// written and never read is dead work, and is usually the wrong side of a bridge.
func TestBoardDataAttributesMatchBetweenRenderAndHandler(t *testing.T) {
	t.Run("BREXB-I02: The page's data attributes are written and read in pairs", func(t *testing.T) {})
	b, err := fs.ReadFile(boardPageForTest, "board/anchors-board.html")
	if err != nil {
		t.Fatal(err)
	}
	page := string(b)

	// `data-vai-para` in HTML is `dataset.vaiPara` in JS: the platform turns hyphens into
	// camelCase, and comparing both raw spellings would report a bridge that does not exist.
	// Every key becomes the hyphenless lower-case form before entering the sets.
	key := func(s string) string { return strings.ToLower(strings.ReplaceAll(s, "-", "")) }

	// Writing has THREE forms, and missing one becomes a false positive: the attribute in the
	// template, `setAttribute`, and the assignment `dataset.x = ...`.
	written := map[string]bool{}
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`data-([a-z-]+)=`),
		regexp.MustCompile(`setAttribute\("data-([a-z-]+)"`),
		regexp.MustCompile(`dataset\.([a-zA-Z]+)\s*=[^=]`),
	} {
		for _, m := range re.FindAllStringSubmatch(page, -1) {
			written[key(m[1])] = true
		}
	}

	// Reading: the selector `[data-x]`, the selector with a value `[data-x="v"]` (CSS
	// included), and `dataset.x` outside an assignment.
	read := map[string]bool{}
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`\[data-([a-z-]+)[\]=]`),
		regexp.MustCompile(`dataset\.([a-zA-Z]+)`),
	} {
		for _, m := range re.FindAllStringSubmatch(page, -1) {
			read[key(m[1])] = true
		}
	}

	for attr := range read {
		if !written[attr] {
			t.Errorf("the handler looks for `data-%s` (normalized spelling), and NO render "+
				"writes that attribute — the click finds no target and the modal does not open", attr)
		}
	}
	for attr := range written {
		if !read[attr] {
			t.Errorf("the render writes `data-%s` and nobody reads it — an attribute written and "+
				"never looked up is usually the renamed side of a broken bridge", attr)
		}
	}
}

// A BUG IS NOT A DECISION, and the board keeps them apart: "Esperando você" counts only what
// asks a choice, and `anchors:bug` gets its own strip. Before the label, six bugs opened as
// decisions made the strip ask for choices that did not exist.
func TestBlockedStrip_separatesBugsFromDecisions(t *testing.T) {
	t.Run("BREXB-B17: Bugs are listed apart from decisions", func(t *testing.T) {})
	hidden, html := runDrawBlocked(t, []map[string]any{
		{"state": "anchors:to-do", "escalated": true, "number": 10, "code": "DEC", "title": "[decision] A ou B", "url": "u"},
		{"state": "anchors:to-do", "bug": true, "number": 20, "code": "bug", "title": "[bug] o job pega o card errado", "url": "u"},
	})
	if hidden {
		t.Fatal("with a decision and a bug the strip must show")
	}
	if !strings.Contains(html, "Esperando você · 1") {
		t.Errorf("the bug must not count as waiting for you:\n%s", html)
	}
	if !strings.Contains(html, "Bugs · 1") || !strings.Contains(html, "o job pega o card errado") {
		t.Errorf("the bug must be listed in its own strip:\n%s", html)
	}

	hidden, html = runDrawBlocked(t, []map[string]any{
		{"state": "anchors:to-do", "bug": true, "number": 20, "code": "bug", "title": "[bug] x", "url": "u"},
	})
	if hidden || strings.Contains(html, "Esperando você") || !strings.Contains(html, "Bugs · 1") {
		t.Errorf("only a bug: the strip shows the bug and no 'waiting for you' (hidden=%v):\n%s", hidden, html)
	}
}

// THE AGENT FILTER answers "who is working right now?": one chip per agent that owns an
// open card touched in the last 30 minutes. These tests run the page's JS in `node`, like
// the blocked-cards strip, because the rule lives in the page — both producers (the
// pipeline and `board serve`) feed it the same fields.

type activeAgent struct {
	Name  string `json:"nome"`
	Cards int    `json:"cards"`
}

func runActiveAgents(t *testing.T, buildData string) []activeAgent {
	t.Helper()
	out := runPageJS(t, `
const d = (() => { `+buildData+` })();
console.log(JSON.stringify(agentesAtivos(d, referenciaDaAtividade(d))));
`)
	var r []activeAgent
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
	t.Run("BREXB-B04: The active agents are counted from the snapshot", func(t *testing.T) {})
	got := runActiveAgents(t, `return {
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
	want := []activeAgent{{"alice", 2}, {"erin", 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("active agents = %+v, want %+v\n"+
			"(bob is 31 min old, the released card has no owner, dave's card is closed)", got, want)
	}
}

// `board serve` is LIVE: its stamp is the time of a read the floor keeps handing back, so
// the window counts from now. The stamp here is a year old on purpose — counting from it
// would put an agent idle for a year on the list.
func TestActiveAgents_liveBoardCountsFromNow(t *testing.T) {
	t.Run("BREXB-B05: A live board counts the active agents from now", func(t *testing.T) {})
	got := runActiveAgents(t, `
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
	want := []activeAgent{{"alice", 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("live board: active agents = %+v, want %+v", got, want)
	}
}

// Clicking a chip shows only that agent's cards; the chip row shows the active agents with
// their counts and an "all" chip, or a quiet line when nobody is active.
func TestAgentChips_filterAndRow(t *testing.T) {
	t.Run("BREXB-B06: The agent chips filter the cards by agent", func(t *testing.T) {})
	out := runPageJS(t, `
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
		ChipRow  string `json:"comAgentes"`
		Passes   []bool `json:"passa"`
		NoFilter []bool `json:"semFiltro"`
		Empty    string `json:"vazio"`
	}
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("node output: %v\n%s", err, out)
	}
	if !reflect.DeepEqual(r.Passes, []bool{true, false, false}) {
		t.Errorf("with alice selected, passaNoFiltro = %v, want only alice's card", r.Passes)
	}
	if !reflect.DeepEqual(r.NoFilter, []bool{true, true, true}) {
		t.Errorf("with no agent selected, passaNoFiltro = %v, want every card", r.NoFilter)
	}
	for _, s := range []string{`data-agente="alice"`, `data-agente="bob"`, `data-agente=""`} {
		if !strings.Contains(r.ChipRow, s) {
			t.Errorf("the chip row lacks %s:\n%s", s, r.ChipRow)
		}
	}
	if strings.Contains(r.ChipRow, "carol") {
		t.Error("carol's only card is a year old, and she got a chip")
	}
	if strings.Index(r.ChipRow, `"alice"`) > strings.Index(r.ChipRow, `"bob"`) {
		t.Error("the chips are not sorted by name")
	}
	if !strings.Contains(r.Empty, "nenhum agente ativo") {
		t.Errorf("with nobody active the row should say so quietly, got:\n%s", r.Empty)
	}
}

// THE RELEASED OWNER IS NOT AN AGENT, and the rule is the producer's: the collect jq —
// the same one `board serve` runs — turns a last `anchors-owner: (liberado)` into an empty
// owner, and the page skips empty owners. Run the jq for real on `gh issue list` shaped
// input.
func TestBoardCollectJQReleasedOwnerIsNoAgent(t *testing.T) {
	t.Run("BREXB-B07: A released card is collected with no owner", func(t *testing.T) {})
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
	var items []struct {
		Number  int    `json:"number"`
		Owner   string `json:"owner"`
		Updated string `json:"updated"`
	}
	if err := json.Unmarshal(out, &items); err != nil {
		t.Fatalf("jq output: %v\n%s", err, out)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2:\n%s", len(items), out)
	}
	if items[0].Owner != "" {
		t.Errorf("a card whose last `anchors-owner:` is `(liberado)` has owner %q — the "+
			"agent filter would list whoever just released it as working", items[0].Owner)
	}
	if items[1].Owner != "bob" {
		t.Errorf("owner = %q, want bob (first line of the comment only)", items[1].Owner)
	}
	// `updated` is the activity signal the page reads; it must reach the JSON.
	if items[0].Updated != "2026-01-01T11:59:00Z" {
		t.Errorf("updated = %q, want the issue's updatedAt", items[0].Updated)
	}
}
