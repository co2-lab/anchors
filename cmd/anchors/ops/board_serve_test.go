package ops

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// O BOARD PUBLICADO E UMA FOTO, e a foto envelhece.
//
// O pipeline roda, gera o `board.json`, publica. Entre uma execucao e outra o que se ve e'
// passado -- e num projeto com nove agentes entregando, minutos bastam para a foto mentir.
//
// O `board serve` sobe o MESMO html apontando para dado vivo. O "mesmo" nao e' detalhe:
// duas renderizacoes do board divergiriam, e a local passaria a mentir de outro jeito.
func TestBoardServeExiste(t *testing.T) {
	c := newBoardCmd()
	var temServe bool
	for _, s := range c.Commands() {
		if s.Name() == "serve" {
			temServe = true
		}
	}
	if !temServe {
		t.Error("`anchors board` existe mas nao tem `serve` — sem ele o board " +
			"local nao sobe")
	}
}

// A ESCUTA e' o desenho B: o custo acompanha a MUDANCA, nao o tempo.
//
// Um poll de 30s custa 1800 chamadas por hora, mudando algo ou nao -- e o rate limit do
// GitHub ja estourou neste projeto com seis agentes. A escuta incremental pergunta "o que
// mudou desde X?", e o custo cai para o que de fato mudou.
//
// O FALLBACK e' o desenho A, e existe porque a escuta pode nao estar disponivel (sem
// webhook, sem permissao, API recusando). Sem ele o comando falharia onde o poll
// funcionaria -- e um board que nao sobe e' pior que um board com alguns segundos de
// atraso.
func TestBoardServeEscutaComFallbackParaPoll(t *testing.T) {
	c := newBoardCmd()
	var longo string
	var temIntervalo bool
	for _, sub := range c.Commands() {
		if sub.Name() != "serve" {
			continue
		}
		longo = sub.Long
		if sub.Flags().Lookup("interval") != nil {
			temIntervalo = true
		}
	}
	if longo == "" {
		t.Fatal("`board serve` nao existe ou nao tem doc")
	}

	// O intervalo do fallback e' configuravel: projeto grande e projeto pequeno nao
	// toleram a mesma frequencia.
	if !temIntervalo {
		t.Error("falta a flag `--interval` — o fallback precisa ser calibravel pelo projeto")
	}

	// A doc precisa NOMEAR o fallback. Um comando que silenciosamente troca de estrategia
	// e' um comando cujo comportamento ninguem consegue prever.
	if !strings.Contains(strings.ToLower(longo), "fallback") &&
		!strings.Contains(strings.ToLower(longo), "recua") {
		t.Error("a doc nao diz que ha fallback — trocar de estrategia em silencio faz o " +
			"comportamento ficar imprevisivel para quem depende dele")
	}
}

// O INCREMENTAL PRECISA PERGUNTAR PELA MAIS RECENTE, e nao pela de maior numero.
//
// MEDIDO: o `gh issue list --limit 1` devolve a issue de maior NUMERO, nao a tocada por
// ultimo. Onze cards foram fechados e o board NAO releu, porque o topo da lista (`#683`,
// 19:22) era mais velho que o card de fato tocado (`#630`, 19:29).
//
// O defeito e' silencioso do pior jeito: o board segue servindo, com cara de fresco, e
// mostra passado. E' a foto de novo -- agora local, e sem ninguem saber.
//
// Duas exigencias, e as duas vem de medicao:
//
//	sort=updated   sem isso a lista vem por numero, e a pergunta responde errado
//	gh api (REST)  o `gh issue list` usa GraphQL, que tem limite SECUNDARIO invisivel
//	               nos contadores -- derruba a consulta enquanto `rate_limit` diz cheio
func TestIncrementalPerguntaPelaMaisRecente(t *testing.T) {
	b, err := os.ReadFile("board_serve.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)

	i := strings.Index(src, "func (f *boardSource) algoMudou()")
	if i < 0 {
		t.Fatal("nao achei o `algoMudou` — o incremental sumiu")
	}
	corpo := src[i:]
	if j := strings.Index(corpo, "\n}"); j > 0 {
		corpo = corpo[:j]
	}

	if !strings.Contains(corpo, "sort=updated") {
		t.Error("o incremental nao ordena por `updated` — a consulta devolve a issue de " +
			"maior NUMERO, e uma issue antiga tocada agora fica invisivel. O board segue " +
			"servindo com cara de fresco, mostrando passado")
	}
	if strings.Contains(corpo, `"issue", "list"`) {
		t.Error("o incremental usa `gh issue list`, que vai por GraphQL — e o limite " +
			"SECUNDARIO dele derruba a consulta enquanto `gh api rate_limit` ainda " +
			"reporta cota cheia. O REST nao tem esse problema")
	}
}

// The live board says it is live, and stamps `takenAt` like the published one.
//
// The page's agent filter counts "the last 30 minutes" from `takenAt` for the published
// board and from NOW when `live` is true. Without `live`, `board serve` would count from
// the stamp of a read the floor keeps handing back, and active agents would drop off the
// list while still working. Without `takenAt`, the header read "Invalid Date".
func TestServePayloadIsLiveAndStamped(t *testing.T) {
	agora := time.Date(2026, 9, 24, 12, 0, 0, 0, time.FixedZone("BRT", -3*3600))
	var b struct {
		TakenAt string           `json:"takenAt"`
		Live    bool             `json:"live"`
		Items   []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(montaBoard(`[{"number":1}]`, agora), &b); err != nil {
		t.Fatalf("the payload is not JSON: %v", err)
	}
	if !b.Live {
		t.Error("the `board serve` payload does not say `live: true` — the page would count " +
			"agent activity from the stamp of a cached read, not from now")
	}
	if b.TakenAt != "2026-09-24T15:00:00Z" {
		t.Errorf("takenAt = %q, want the read time in UTC (2026-09-24T15:00:00Z)", b.TakenAt)
	}
	if len(b.Items) != 1 {
		t.Errorf("items = %v, want the one item passed in", b.Items)
	}
}

// --- the live collection, against a fake `gh` (nothing reaches GitHub) ---

// ghIssuesPages is what `gh api --paginate` hands back: the pages glued together, which
// is not valid JSON on its own. Issue 3 is a pull request and must not become a card.
const ghIssuesPages = `[{"number":1,"title":"A","url":"u1","body":"","labels":[{"name":"anchors"},{"name":"anchors:doing"}],"user":{"login":"x"},"updated_at":"2026-09-20T10:00:00Z","created_at":"2026-09-01T00:00:00Z","closed_at":null}]` +
	`[{"number":2,"title":"B","url":"u2","body":"","labels":[{"name":"anchors"}],"user":{"login":"y"},"updated_at":"2026-09-21T10:00:00Z","created_at":"2026-09-01T00:00:00Z","closed_at":null},` +
	`{"number":3,"pull_request":{},"title":"PR","labels":[],"user":{"login":"z"},"updated_at":"2026-09-22T00:00:00Z"}]`

const ghOwnerComments = `{"data":{"repository":{"issues":{"pageInfo":{"hasNextPage":false,"endCursor":""},"nodes":[{"number":1,"comments":{"nodes":[{"body":"anchors-owner: alice"}]}}]}}}}`

// fakeBoardGH answers the three calls of the collection: the comments (GraphQL), the
// issues (REST, paginated), and the "anything newer?" question, answered with `newest`.
func fakeBoardGH(t *testing.T, newest string) string {
	t.Helper()
	return fakeGH(t, `case "$*" in
*graphql*) echo '`+ghOwnerComments+`' ;;
*--paginate*) printf '%s' '`+ghIssuesPages+`' ;;
*sort=updated*) echo '`+newest+`' ;;
*) exit 1 ;;
esac`)
}

type boardDoc struct {
	Live  bool `json:"live"`
	Items []struct {
		Number  int    `json:"number"`
		Updated string `json:"updated"`
		Owner   string `json:"owner"`
	} `json:"items"`
}

func decodeBoard(t *testing.T, b []byte) boardDoc {
	t.Helper()
	var d boardDoc
	if err := json.Unmarshal(b, &d); err != nil {
		t.Fatalf("the board is not JSON: %v\n%s", err, b)
	}
	return d
}

// The full sweep goes through REST, stitches the pages, drops the pull requests, and
// takes the owner from the `anchors-owner:` comment fetched through GraphQL.
func TestFullCollectionStitchesPagesAndKeepsTheOwner(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not installed")
	}
	log := fakeBoardGH(t, "")
	b, err := coletaCompleta("acme/app")
	if err != nil {
		t.Fatalf("coletaCompleta: %v", err)
	}
	d := decodeBoard(t, b)
	if !d.Live || len(d.Items) != 2 {
		t.Fatalf("want the two issues of both pages (not the PR), live; got %+v", d)
	}
	if d.Items[0].Number != 1 || d.Items[0].Owner != "alice" || d.Items[1].Owner != "" {
		t.Errorf("items = %+v; want #1 owned by alice and #2 without owner", d.Items)
	}
	if calls := readLog(t, log); !strings.Contains(calls, "repos/acme/app/issues") || strings.Contains(calls, "issue list") {
		t.Errorf("the sweep must go through REST, not `gh issue list`:\n%s", calls)
	}
}

// A refusal from GitHub is an error that carries gh's own message.
func TestFullCollectionReportsTheRefusal(t *testing.T) {
	fakeGH(t, `echo "API rate limit exceeded" >&2; exit 1`)
	_, err := coletaCompleta("acme/app")
	if err == nil || !strings.Contains(err.Error(), "collect from GitHub") ||
		!strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Fatalf("want the collection error with gh's stderr, got %v", err)
	}
}

// The read floor, the incremental question, and the fallback — through `leia`.
func TestBoardSourceReadsOnlyWhenSomethingChanged(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq is not installed")
	}
	log := fakeBoardGH(t, "2026-09-21T10:00:00Z") // the newest issue is the one already seen
	sweeps := func() int { return strings.Count(readLog(t, log), "--paginate") }

	f := &boardSource{repo: "acme/app", piso: time.Hour}
	first, err := f.leia()
	if err != nil {
		t.Fatalf("first read: %v", err)
	}
	if f.desdeAt != "2026-09-21T10:00:00Z" {
		t.Errorf("desdeAt = %q, want the newest `updated` of the board", f.desdeAt)
	}
	// Inside the floor: the same bytes, no call at all.
	before := readLog(t, log)
	again, _ := f.leia()
	if string(again) != string(first) || readLog(t, log) != before {
		t.Error("a read inside the floor called GitHub again")
	}

	// Past the floor, nothing newer: ONE question, and no sweep.
	f.piso = 0
	if _, err := f.leia(); err != nil {
		t.Fatal(err)
	}
	if n := sweeps(); n != 1 {
		t.Errorf("with nothing changed the board swept again (%d sweeps)", n)
	}

	// Something newer: the incremental says so, and the board sweeps.
	f.desdeAt = "2026-09-01T00:00:00Z"
	if _, err := f.leia(); err != nil {
		t.Fatal(err)
	}
	if n := sweeps(); n != 2 {
		t.Errorf("a newer issue did not trigger the sweep (%d sweeps)", n)
	}
}

// When the sweep fails after a good read, the board keeps serving the previous read; with
// no previous read, the failure is the answer.
func TestBoardSourceFallsBackToTheLastGoodRead(t *testing.T) {
	fakeGH(t, `exit 1`)
	fresh := &boardSource{repo: "acme/app"}
	if _, err := fresh.leia(); err == nil {
		t.Error("the first read failing must be an error — there is nothing to serve")
	}
	old := []byte(`{"items":[]}`)
	f := &boardSource{repo: "acme/app", ultimo: old, desdeAt: "2026-01-01T00:00:00Z"}
	got, err := f.leia()
	if err != nil || string(got) != string(old) {
		t.Errorf("leia = %q, %v; want the previous read and no error", got, err)
	}
}

func TestAlgoMudouComparesWithTheNewestSeen(t *testing.T) {
	log := fakeGH(t, `echo "2026-09-21T10:00:00Z"`)
	f := &boardSource{repo: "acme/app", desdeAt: "2026-09-20T00:00:00Z"}
	if changed, err := f.algoMudou(); err != nil || !changed {
		t.Errorf("newer issue: algoMudou = %v, %v; want true", changed, err)
	}
	f.desdeAt = "2026-09-21T10:00:00Z"
	if changed, _ := f.algoMudou(); changed {
		t.Error("the same timestamp is not a change")
	}
	if !strings.Contains(readLog(t, log), "/repos/acme/app/issues?") {
		t.Errorf("the question did not go to the repository's issues:\n%s", readLog(t, log))
	}
}

func TestMaisRecenteTakesTheNewestUpdated(t *testing.T) {
	got := maisRecente([]byte(`{"items":[{"updated":"2026-01-02"},{"updated":"2026-03-01"},{"updated":"2026-02-01"}]}`))
	if got != "2026-03-01" {
		t.Errorf("maisRecente = %q, want 2026-03-01", got)
	}
	if maisRecente([]byte("not json")) != "" {
		t.Error("an unreadable board has no newest date")
	}
}

func TestRepoAtualAsksTheClone(t *testing.T) {
	fakeGH(t, `echo "acme/app"`)
	if r, err := repoAtual(); err != nil || r != "acme/app" {
		t.Errorf("repoAtual = %q, %v", r, err)
	}
	fakeGH(t, `echo ""`)
	if _, err := repoAtual(); err == nil {
		t.Error("an empty answer from gh must be an error, not an empty repository")
	}
}

// The comments come in pages; a failure keeps what was read, and a malformed repository
// name reads nothing.
func TestCommentsOfOpenCardsFollowsTheCursor(t *testing.T) {
	log := fakeGH(t, `case "$*" in
*cursor=C1*) echo '{"data":{"repository":{"issues":{"pageInfo":{"hasNextPage":false},"nodes":[{"number":7,"comments":{"nodes":[{"body":"b7"}]}}]}}}}' ;;
*) echo '{"data":{"repository":{"issues":{"pageInfo":{"hasNextPage":true,"endCursor":"C1"},"nodes":[{"number":5,"comments":{"nodes":[{"body":"b5"}]}}]}}}}' ;;
esac`)
	var got map[string][]map[string]string
	if err := json.Unmarshal([]byte(commentsOfOpenCards("acme/app")), &got); err != nil {
		t.Fatal(err)
	}
	if got["5"][0]["body"] != "b5" || got["7"][0]["body"] != "b7" {
		t.Errorf("both pages should be read, got %v", got)
	}
	if n := strings.Count(readLog(t, log), "graphql"); n != 2 {
		t.Errorf("want two GraphQL calls (one per page), got %d", n)
	}
	if commentsOfOpenCards("no-slash") != "{}" {
		t.Error("a repository without owner/name must read nothing")
	}

	fakeGH(t, `exit 1`)
	if commentsOfOpenCards("acme/app") != "{}" {
		t.Error("a refused query must yield the empty map, not break the board")
	}
	path, cleanup := commentsToFile("acme/app")
	cleanup()
	if path != "" {
		t.Errorf("with no comments there is no file to pass, got %q", path)
	}
}

func TestCommentsToFileWritesTheMapAndCleansUp(t *testing.T) {
	fakeGH(t, `echo '`+ghOwnerComments+`'`)
	path, cleanup := commentsToFile("acme/app")
	if path == "" {
		t.Fatal("the comments were not written")
	}
	b, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(b), "anchors-owner: alice") {
		t.Errorf("file = %q, %v", b, err)
	}
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("the cleanup did not remove the temporary file")
	}
}

// `board serve` without --repo asks the clone; when that fails the error says how to
// declare it. With a repository, it prints where it serves and returns the listen error.
func TestBoardServeCommandDiscoversTheRepoOrFails(t *testing.T) {
	fakeGH(t, `exit 1`)
	c := newBoardServeCmd()
	c.SetArgs([]string{})
	c.SetOut(io.Discard)
	c.SetErr(io.Discard)
	err := c.Execute()
	if err == nil || !strings.Contains(err.Error(), "--repo <owner/name>") {
		t.Fatalf("want the discovery error pointing at --repo, got %v", err)
	}

	c = newBoardServeCmd()
	c.SetArgs([]string{"--repo", "acme/app", "--port", "-1"})
	c.SetOut(io.Discard)
	c.SetErr(io.Discard)
	out := captureStdout(t, func() { err = c.Execute() })
	if err == nil {
		t.Fatal("an invalid port must fail to listen")
	}
	if !strings.Contains(out, "repository: acme/app") || !strings.Contains(out, "localhost:-1") {
		t.Errorf("the banner does not say what is served:\n%s", out)
	}
}
