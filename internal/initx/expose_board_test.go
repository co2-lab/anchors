package initx

import (
	"embed"
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
