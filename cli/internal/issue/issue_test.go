package issue

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func viol() Issue {
	return Issue{
		Kind: Violation, Target: "features/x/A.spec.md", Gate: "spec-sections",
		Detail: "falta a seção Regras", Date: "2026-08-07",
	}
}

func TestKeyIsStableAcrossDates(t *testing.T) {
	a := viol()
	b := viol()
	b.Date = "2027-01-01" // outra data
	if a.Key() != b.Key() {
		t.Fatalf("a Key deve ser estável no tempo: %q vs %q", a.Key(), b.Key())
	}
	if a.ID() == b.ID() {
		t.Fatal("o ID (nome de arquivo) deve variar com a data")
	}
	// gate diferente → Key diferente (mesmo alvo pode violar dois gates)
	c := viol()
	c.Gate = "spec-tem-codigo"
	if a.Key() == c.Key() {
		t.Fatal("gates distintos deveriam gerar Keys distintas")
	}
}

func TestIDIsLegibleAndSanitized(t *testing.T) {
	id := viol().ID()
	if !strings.HasPrefix(id, "2026-08-07--violation--") {
		t.Fatalf("ID inesperado: %s", id)
	}
	if strings.Contains(id, "/") {
		t.Fatalf("ID não deve conter barra: %s", id)
	}
}

func TestOpenCreatesInTodo(t *testing.T) {
	root := t.TempDir()
	created, at, err := Open(root, viol())
	if err != nil || !created || at != Todo {
		t.Fatalf("open: created=%v at=%v err=%v", created, at, err)
	}
	ids, _ := List(root, Todo)
	if len(ids) != 1 {
		t.Fatalf("esperava 1 issue em todo/, veio %v", ids)
	}
	body, _ := os.ReadFile(filepath.Join(root, Dir, string(Todo), ids[0]))
	s := string(body)
	for _, want := range []string{"VIOLATION", "features/x/A.spec.md", "spec-sections", "falta a seção Regras"} {
		if !strings.Contains(s, want) {
			t.Errorf("corpo não contém %q:\n%s", want, s)
		}
	}
}

func TestOpenIsIdempotentAcrossDates(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, viol())
	// reconfrontar em OUTRO dia não duplica (mesma Key)
	later := viol()
	later.Date = "2026-09-15"
	created, at, _ := Open(root, later)
	if created {
		t.Fatal("não deveria recriar issue já existente (mesma Key, outra data)")
	}
	if at != Todo {
		t.Fatalf("estado atual deveria ser todo, veio %v", at)
	}
	if ids, _ := List(root, Todo); len(ids) != 1 {
		t.Fatalf("esperava 1 issue, veio %d", len(ids))
	}
}

func TestResolveMovesToDone(t *testing.T) {
	root := t.TempDir()
	_, _, _ = Open(root, viol())
	// o confronto voltou a passar → resolve
	ok, err := Resolve(root, viol().Key())
	if err != nil || !ok {
		t.Fatalf("resolve: ok=%v err=%v", ok, err)
	}
	if ids, _ := List(root, Todo); len(ids) != 0 {
		t.Fatalf("todo/ deveria esvaziar após resolve, veio %v", ids)
	}
	if ids, _ := List(root, Done); len(ids) != 1 {
		t.Fatalf("done/ deveria ter 1, veio %v", ids)
	}
	// resolver de novo é no-op
	ok2, _ := Resolve(root, viol().Key())
	if ok2 {
		t.Fatal("resolver uma issue já resolvida deveria ser no-op")
	}
}

func TestResolveFromDoing(t *testing.T) {
	root := t.TempDir()
	// simula issue em doing/ (alguém pegou para resolver)
	i := viol()
	doingDir := filepath.Join(root, Dir, string(Doing))
	_ = os.MkdirAll(doingDir, 0o755)
	_ = os.WriteFile(filepath.Join(doingDir, i.ID()), []byte(i.Body()), 0o644)
	// o check passa → resolve mesmo estando em doing/
	ok, _ := Resolve(root, i.Key())
	if !ok {
		t.Fatal("deveria resolver issue que estava em doing/")
	}
	if ids, _ := List(root, Done); len(ids) != 1 {
		t.Fatalf("done/ deveria ter 1, veio %v", ids)
	}
}

func TestOpenDoesNotResurrectResolvedIssue(t *testing.T) {
	root := t.TempDir()
	i := viol()
	// já resolvida (em done/)
	doneDir := filepath.Join(root, Dir, string(Done))
	_ = os.MkdirAll(doneDir, 0o755)
	_ = os.WriteFile(filepath.Join(doneDir, i.ID()), []byte("resolvida"), 0o644)
	created, at, _ := Open(root, i)
	if created {
		t.Fatal("não deveria reabrir issue já em done/")
	}
	if at != Done {
		t.Fatalf("deveria reportar done, veio %v", at)
	}
	if todos, _ := List(root, Todo); len(todos) != 0 {
		t.Fatalf("todo/ deveria ficar vazio, veio %v", todos)
	}
}

// TestDividaNasceEmFutureEFechaAoSerPaga prova o ciclo de vida da dívida assumida.
//
// Antes, `obligation_pending` era uma linha no cabeçalho de um arquivo: visível só para
// quem o abrisse, sem estado, sem como ser paga, sem como vencer. O gate afirmava que quem
// declara diz três coisas — conhece o dever, ele vale, e QUANDO será pago — e o "quando"
// era prosa livre que nada confrontava.
func TestDividaNasceEmFutureEFechaAoSerPaga(t *testing.T) {
	root := t.TempDir()
	iss := Issue{
		Kind: Violation, Target: "infra/models/X.spec.md", Gate: "obligation-honored",
		Detail: "não alcançada pelo purge", Date: "2026-08-13",
		Prazo: "`lgpd-eliminacao` — na etapa de código da Fase 1",
	}

	created, at, err := OpenAt(root, iss, Future)
	if err != nil || !created {
		t.Fatalf("dívida deveria nascer: created=%v err=%v", created, err)
	}
	if at != Future {
		t.Errorf("dívida nasce em `future/`, não em %q — quem lê `todo/` pergunta "+
			"'o que faço AGORA'", at)
	}

	// Reconfrontar não duplica nem promove para `todo/`: continua adiada.
	if created, at, _ := OpenAt(root, iss, Future); created || at != Future {
		t.Errorf("reconfronto não pode duplicar; veio created=%v at=%q", created, at)
	}

	// Pagar a dívida (o gate volta a passar) fecha a issue — o mesmo ciclo das demais.
	ok, err := Resolve(root, iss.Key())
	if err != nil || !ok {
		t.Fatalf("dívida paga deve fechar: ok=%v err=%v", ok, err)
	}
	if st, existe := Exists(root, iss.Key()); !existe || st != Done {
		t.Errorf("depois de paga a dívida vive em `done/`; veio %q (existe=%v)", st, existe)
	}
}

// O prazo tem de aparecer no corpo: uma dívida sem vencimento à vista é um TODO com nome
// melhor.
func TestCorpoDaDividaMostraOVencimento(t *testing.T) {
	iss := Issue{
		Kind: Violation, Target: "a.ts", Gate: "obligation-honored",
		Detail: "laudo", Date: "2026-08-13", Prazo: "`x` — na Fase 2",
	}
	corpo := iss.Body()
	for _, quer := range []string{"Quando será paga", "na Fase 2", "Dívida ASSUMIDA"} {
		if !strings.Contains(corpo, quer) {
			t.Errorf("o corpo da dívida precisa conter %q", quer)
		}
	}
}

// TestAchadoNovoNaoSomeEmSilencio guarda o defeito medido num E2E: um segundo `anchors
// judge --verdict fail` sobre a MESMA unidade só imprimia "issue já registrada" e
// descartava o `--reason`. Dois laudos distintos foram perdidos assim, e a issue seguiu
// com o texto antigo.
//
// Idempotência é a política certa para o MESMO problema detectado duas vezes; não é para
// um problema DIFERENTE no mesmo lugar.
func TestAchadoNovoNaoSomeEmSilencio(t *testing.T) {
	root := t.TempDir()
	base := Issue{Kind: Violation, Target: "a.ts", Gate: "review", Date: "2026-08-13"}

	primeiro := base
	primeiro.Detail = "## Laudo A\nperda de dado na leitura"
	if created, _, err := Open(root, primeiro); err != nil || !created {
		t.Fatalf("primeira issue deveria nascer: %v", err)
	}

	// achado DIFERENTE no mesmo alvo: acrescenta, preservando o anterior
	segundo := base
	segundo.Detail = "## Laudo B\ncontradição entre duas regras"
	reaberta, err := Reabrir(root, segundo)
	if err != nil || !reaberta {
		t.Fatalf("achado novo deveria reabrir: reaberta=%v err=%v", reaberta, err)
	}
	_, name, _ := byKey(root, base.Key())
	corpo, _ := os.ReadFile(pathFor(root, Todo, name))
	for _, quer := range []string{"Laudo A", "Laudo B", "perda de dado", "contradição"} {
		if !strings.Contains(string(corpo), quer) {
			t.Errorf("o corpo precisa preservar %q — laudo antigo E novo", quer)
		}
	}

	// MESMO achado de novo: idempotente, não duplica
	if reaberta, _ := Reabrir(root, segundo); reaberta {
		t.Error("o mesmo achado não pode ser acrescentado duas vezes")
	}
	corpo2, _ := os.ReadFile(pathFor(root, Todo, name))
	if strings.Count(string(corpo2), "Laudo B") != 1 {
		t.Errorf("`Laudo B` deve aparecer 1x; apareceu %d", strings.Count(string(corpo2), "Laudo B"))
	}
}
