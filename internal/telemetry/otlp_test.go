package telemetry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func fixedClock() time.Time { return time.Date(2026, 9, 14, 18, 30, 0, 0, time.UTC) }

// otlpBody is the part of the OTLP/JSON logs body the tests read.
type otlpBody struct {
	ResourceLogs []struct {
		Resource struct {
			Attributes []struct {
				Key   string         `json:"key"`
				Value map[string]any `json:"value"`
			} `json:"attributes"`
		} `json:"resource"`
		ScopeLogs []struct {
			LogRecords []struct {
				TimeUnixNano string         `json:"timeUnixNano"`
				Body         map[string]any `json:"body"`
				Attributes   []struct {
					Key   string         `json:"key"`
					Value map[string]any `json:"value"`
				} `json:"attributes"`
			} `json:"logRecords"`
		} `json:"scopeLogs"`
	} `json:"resourceLogs"`
}

// The BODY must be valid OTLP — and the only way to know is to receive it and check the shape
// the protocol demands. A malformed body is accepted with 200 by some collectors and dropped
// in silence: the defect would only show when nobody found the data.
func TestEmit_producesValidOTLP(t *testing.T) {
	t.Run("TLEMT-B03: An event is sent as one log record whose body is the event name", func(t *testing.T) {})
	var mu sync.Mutex
	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		received = b
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(ClaimServed, map[string]any{"candidates": 14}, fixedClock))
	e.Flush()

	mu.Lock()
	defer mu.Unlock()
	var p otlpBody
	if err := json.Unmarshal(received, &p); err != nil {
		t.Fatalf("the body is not JSON: %v (%q)", err, received)
	}
	if len(p.ResourceLogs) != 1 || len(p.ResourceLogs[0].ScopeLogs) != 1 ||
		len(p.ResourceLogs[0].ScopeLogs[0].LogRecords) != 1 {
		t.Fatalf("expected one resource with one log record, got %s", received)
	}
	rec := p.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	if rec.Body["stringValue"] != "claim.served" {
		t.Errorf("the record's body should be the event name, got %v", rec.Body)
	}
	if want := "1789410600000000000"; rec.TimeUnixNano != want {
		t.Errorf("time = %q, want %q", rec.TimeUnixNano, want)
	}
	if len(rec.Attributes) != 1 || rec.Attributes[0].Key != "candidates" ||
		rec.Attributes[0].Value["intValue"] != "14" {
		t.Errorf("attributes = %+v", rec.Attributes)
	}
}

// The RESOURCE identifies the PRODUCT, never who uses it. Without this, the first attribute
// someone added for convenience — hostname, user, repository — would enter with nothing to
// accuse it. The set of keys is closed, not checked against a list of forbidden names.
func TestEmit_theResourceDoesNotIdentifyTheUser(t *testing.T) {
	t.Run("TLEMT-X01: The resource carries only the product name and version", func(t *testing.T) {})
	e := NewEmitter(Config{Enabled: true}, "0.1.99")
	body, _ := e.build(New(CheckFinished, nil, fixedClock))

	var p otlpBody
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatal(err)
	}
	got := map[string]any{}
	for _, a := range p.ResourceLogs[0].Resource.Attributes {
		got[a.Key] = a.Value["stringValue"]
	}
	if len(got) != 2 || got["service.name"] != "anchors" || got["service.version"] != "0.1.99" {
		t.Errorf("the resource must be exactly the product and its version, got %v", got)
	}
}

// TELEMETRY NEVER BLOCKS. A collector that is down cannot hold anyone's `anchors check` —
// losing an event is infinitely better.
func TestEmit_neitherBlocksNorFailsWithADeadServer(t *testing.T) {
	t.Run("TLEMT-B04: Emitting to a dead collector neither blocks nor fails", func(t *testing.T) {})
	e := NewEmitter(Config{Enabled: true, Endpoint: "http://127.0.0.1:1"}, "0.1.99")
	done := make(chan struct{})
	go func() {
		e.Emit(New(TurnEnded, map[string]any{"state": "in-progress"}, fixedClock))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("`Emit` blocked — telemetry cannot delay anyone's work")
	}
}

// DISABLED returns nil, and the nil accepts `Emit` doing nothing. It is what spares the caller
// an `if` around each call — and a forgotten `if` would be a panic.
func TestNewEmitter_disabledReturnsNilThatAcceptsEmit(t *testing.T) {
	t.Run("TLEMT-B01: A disabled configuration yields no emitter that still accepts calls", func(t *testing.T) {})
	e := NewEmitter(Config{Enabled: false}, "0.1.99")
	if e != nil {
		t.Fatal("disabled, the emitter should be nil")
	}
	e.Emit(New(ClaimEmpty, nil, fixedClock)) // must not panic
}

func TestNewEmitter_blankEndpointIsHoneycomb(t *testing.T) {
	t.Run("TLEMT-B02: A blank endpoint means Honeycomb's logs endpoint", func(t *testing.T) {})
	e := NewEmitter(Config{Enabled: true}, "0.1.99")
	if e.cfg.Endpoint != "https://api.honeycomb.io/v1/logs" {
		t.Errorf("endpoint = %q", e.cfg.Endpoint)
	}
	if e := NewEmitter(Config{Enabled: true, Endpoint: "http://collector"}, "x"); e.cfg.Endpoint != "http://collector" {
		t.Errorf("a declared endpoint must be kept, got %q", e.cfg.Endpoint)
	}
}

func TestEmit_isAJSONPostWithTheConfiguredHeaders(t *testing.T) {
	t.Run("TLEMT-B07: Each request is a JSON POST carrying the configured headers", func(t *testing.T) {})
	var mu sync.Mutex
	var method, ctype, team string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		method, ctype, team = r.Method, r.Header.Get("Content-Type"), r.Header.Get("x-honeycomb-team")
		mu.Unlock()
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL,
		Headers: map[string]string{"x-honeycomb-team": "k"}}, "0.1.99")
	e.Emit(New(CheckFinished, nil, fixedClock))
	e.Flush()

	mu.Lock()
	defer mu.Unlock()
	if method != http.MethodPost || ctype != "application/json" || team != "k" {
		t.Errorf("method=%q content-type=%q header=%q", method, ctype, team)
	}
}

func TestEmit_anErrorStatusIsDiscarded(t *testing.T) {
	t.Run("TLEMT-E02: An error status from the collector is discarded", func(t *testing.T) {})
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits++
		mu.Unlock()
		w.WriteHeader(500)
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(CheckFinished, nil, fixedClock))
	e.Flush()

	mu.Lock()
	defer mu.Unlock()
	if hits != 1 {
		t.Errorf("the collector should see the event exactly once (no retry), saw %d", hits)
	}
}

// The VALUE only accepts number, boolean and vocabulary string. Anything else becomes the
// TYPE DESCRIPTION — the guard that keeps an `error` with a file path, or a `map` with
// content, from leaking through a future call written without care.
func TestOTLPValue_whatIsNotVocabularyBecomesTheType(t *testing.T) {
	t.Run("TLEMT-I01: A value that is not vocabulary becomes its type description", func(t *testing.T) {})
	v := otlpValue(map[string]string{"path": "/home/someone/project/spec.md"})
	s, _ := v["stringValue"].(string)
	if strings.Contains(s, "/home/someone") {
		t.Errorf("an unexpected value leaked content: %q", s)
	}
	if s != "(map[string]string)" {
		t.Errorf("expected the type description, got %q", s)
	}
	// And the expected ones pass.
	if otlpValue(14)["intValue"] != "14" {
		t.Error("a number should become intValue")
	}
	if otlpValue(int64(7))["intValue"] != "7" {
		t.Error("an int64 should become intValue")
	}
	if otlpValue(true)["boolValue"] != true {
		t.Error("a boolean should become boolValue")
	}
	if otlpValue("pr-aberto")["stringValue"] != "pr-aberto" {
		t.Error("product vocabulary should pass")
	}
}

// FLUSH is what makes the event ARRIVE in a CLI.
//
// `Emit` fires the POST in a goroutine so as not to delay the command — and in a command-line
// program that means the process dies before the send goes out.
//
// The defect is INVISIBLE in a unit test: there the process stays alive after the call, and
// the POST ends up going out. It only showed when running the real binary against a local
// collector — which received NOTHING until this mechanism existed.
func TestFlush_waitsForTheSendBeforeTheProcessExits(t *testing.T) {
	t.Run("TLEMT-B05: Flushing waits for the send in flight", func(t *testing.T) {})
	arrived := make(chan struct{}, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The delay simulates what really happens: the POST is not instantaneous, and
		// without the flush the process would die in this interval.
		time.Sleep(50 * time.Millisecond)
		arrived <- struct{}{}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(TurnEnded, map[string]any{"state": "in-progress"}, fixedClock))
	e.Flush()

	select {
	case <-arrived:
	default:
		t.Fatal("`Flush` returned before the send arrived — in a CLI the event would be lost")
	}
}

// FLUSH CANNOT BLOCK EITHER. A server that does not answer cannot delay a command's exit — if
// it is slow, the event is lost, and that is the right outcome.
func TestFlush_doesNotWaitForever(t *testing.T) {
	// The handler waits for a signal that only arrives at the end of the test. `time.Sleep(10s)`
	// would be the same from the client's point of view, and would make the TEST last ten
	// seconds — `srv.Close()` waits for open connections.
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer func() { close(release); srv.Close() }()

	e := NewEmitter(Config{Enabled: true, Endpoint: srv.URL}, "0.1.99")
	e.Emit(New(CheckFinished, nil, fixedClock))

	start := time.Now()
	e.Flush()
	if d := time.Since(start); d > sendDeadline+time.Second {
		t.Errorf("`Flush` waited %v — telemetry cannot delay a command's exit", d)
	}
}

// The FLUSH DEADLINE is a defence OF ITS OWN, not the HTTP client's timeout.
//
// There are two, and the distinction only showed in the mutation run: removing the
// `time.After` from the `select` left the test above green, because the `http.Client` has a
// timeout and returns control on its own. The outer defence hid the inner one.
//
// Why the inner one matters: a send that does not go through the HTTP client — a goroutine
// stuck anywhere else, today or in a future change — would hold the process exit forever.
// `Flush` cannot depend on another layer having a deadline.
func TestFlush_hasItsOwnDeadlineIndependentOfTheHTTPClient(t *testing.T) {
	t.Run("TLEMT-B06: Flushing has its own deadline independent of the HTTP client", func(t *testing.T) {})
	e := NewEmitter(Config{Enabled: true, Endpoint: "http://example.invalid"}, "0.1.99")

	// A send that NEVER finishes and does not go through the HTTP client: it is what tells
	// the `Flush` defence apart from the `http.Client` one.
	e.emVoo.Add(1)
	defer e.emVoo.Done() // release at the end of the test, not to leak the goroutine

	// `Flush` runs in a goroutine with a limit of the TEST's own: without the inner deadline
	// it waits forever, and a hanging test does not fail — it exhausts the whole package's
	// time and hides which assertion broke.
	returned := make(chan time.Duration, 1)
	go func() {
		start := time.Now()
		e.Flush()
		returned <- time.Since(start)
	}()

	select {
	case d := <-returned:
		if d > sendDeadline+500*time.Millisecond {
			t.Errorf("`Flush` waited %v without the HTTP client in the way", d)
		}
	case <-time.After(sendDeadline + 2*time.Second):
		t.Fatal("`Flush` did not return — it needs its OWN deadline, or a stuck goroutine " +
			"holds the process exit forever")
	}
}

// Nil accepts `Flush` — it lets `main` call it without knowing whether telemetry is on, and a
// forgotten `if` there would be a panic at process exit.
func TestFlush_nilDoesNotPanic(t *testing.T) {
	var e *Emitter
	e.Flush()
}
