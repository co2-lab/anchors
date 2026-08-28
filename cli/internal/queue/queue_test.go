package queue

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func task(id, changed, kind, next string) Task {
	return Task{ID: id, Changed: changed, Kind: kind, Origin: "watch", SuggestedNext: next, CreatedAt: "2026-08-07T00:00:00Z"}
}

// comAlvo cria o arquivo-alvo no disco. Desde que `List` descarta task cujo alvo não
// existe (a task-fantasma que sobrava em toda execução real), o alvo tem de ser material.
func comAlvo(t *testing.T, root, rel string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestEnqueueAndList(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "AddItem.spec.md")
	created, err := Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	if err != nil || !created {
		t.Fatalf("enqueue: created=%v err=%v", created, err)
	}
	tasks, err := List(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].State != Pending || tasks[0].Changed != "AddItem.spec.md" {
		t.Fatalf("esperava 1 task pending; veio %+v", tasks)
	}
}

func TestEnqueueDedup(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	// mesma (changed, suggested_next) com ID diferente → NÃO duplica
	created, err := Enqueue(root, task("2-spec-add", "AddItem.spec.md", "spec", "implement"))
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("não deveria duplicar task viva para o mesmo alvo+passo")
	}
	if n, _ := PendingCount(root); n != 1 {
		t.Fatalf("esperava 1 pendente, veio %d", n)
	}
}

func TestClaimEmpty(t *testing.T) {
	root := t.TempDir()
	got, err := Claim(root, "w1", "2026-08-13T00:00:00-03:00")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("fila vazia deveria devolver nil, veio %+v", got)
	}
}

func TestClaimMovesToClaimed(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	got, err := Claim(root, "worker-A", "2026-08-13T00:00:00-03:00")
	if err != nil || got == nil {
		t.Fatalf("claim: %+v err=%v", got, err)
	}
	if got.State != Claimed || got.ClaimedBy != "worker-A" {
		t.Fatalf("esperava claimed por worker-A, veio %+v", got)
	}
	// o arquivo pending sumiu, o claimed existe
	d := dirFor(root)
	if _, err := os.Stat(filepath.Join(d, fileName(Pending, "1-spec-add"))); !os.IsNotExist(err) {
		t.Fatal("arquivo pending deveria ter sumido")
	}
	if _, err := os.Stat(filepath.Join(d, fileName(Claimed, "1-spec-add"))); err != nil {
		t.Fatal("arquivo claimed deveria existir")
	}
}

// TestClaimAtomicNoDoubleClaim: N workers concorrentes sobre M tasks nunca pegam a
// mesma task. É a garantia que permite dois terminais rodando `anchors next`.
func TestClaimAtomicNoDoubleClaim(t *testing.T) {
	root := t.TempDir()
	const M = 20
	for i := 0; i < M; i++ {
		id := filepath.Base(t.Name()) + "-" + string(rune('a'+i))
		alvo := "f" + string(rune('a'+i)) + ".spec.md"
		comAlvo(t, root, alvo)
		_, _ = Enqueue(root, task(id, alvo, "spec", "implement"))
	}
	var mu sync.Mutex
	claimed := map[string]int{}
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for {
				got, err := Claim(root, "w", "2026-08-13T00:00:00-03:00")
				if err != nil || got == nil {
					return
				}
				mu.Lock()
				claimed[got.ID]++
				mu.Unlock()
			}
		}(w)
	}
	wg.Wait()
	if len(claimed) != M {
		t.Fatalf("esperava %d tasks reivindicadas, veio %d", M, len(claimed))
	}
	for id, n := range claimed {
		if n != 1 {
			t.Fatalf("task %s reivindicada %d vezes (double-claim!)", id, n)
		}
	}
}

func TestMarkDoneMovesToHistory(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "AddItem.spec.md")
	_, _ = Enqueue(root, task("1-spec-add", "AddItem.spec.md", "spec", "implement"))
	_, _ = Claim(root, "w1", "2026-08-13T00:00:00-03:00")
	if err := MarkDone(root, "1-spec-add"); err != nil {
		t.Fatal(err)
	}
	// sumiu da fila viva
	if n, _ := PendingCount(root); n != 0 {
		t.Fatalf("esperava fila vazia após done, veio %d", n)
	}
	// apareceu no histórico
	donePath := filepath.Join(root, DoneDir, fileName(Done, "1-spec-add"))
	if _, err := os.Stat(donePath); err != nil {
		t.Fatalf("task concluída deveria estar em %s: %v", DoneDir, err)
	}
	// e uma nova task para o mesmo alvo agora É permitida (a anterior saiu da fila)
	created, _ := Enqueue(root, task("2-spec-add", "AddItem.spec.md", "spec", "implement"))
	if !created {
		t.Fatal("após done, novo enqueue do mesmo alvo deveria ser permitido")
	}
}

func TestDropRemovesFromQueue(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "plans/x.md")
	_, _ = Enqueue(root, task("1-doc-x", "plans/x.md", "doc", "triage"))
	if err := Drop(root, "1-doc-x"); err != nil {
		t.Fatal(err)
	}
	if n, _ := PendingCount(root); n != 0 {
		t.Fatalf("após drop a fila deveria esvaziar, veio %d", n)
	}
	// drop não cria histórico (diferente de done)
	if _, err := os.Stat(filepath.Join(root, DoneDir)); !os.IsNotExist(err) {
		t.Fatal("drop não deveria criar .anchors/done/")
	}
	// drop de task inexistente → erro
	if err := Drop(root, "nao-existe"); err == nil {
		t.Fatal("drop de task inexistente deveria falhar")
	}
}

func TestReclaimReturnsClaimedToPending(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "A.spec.md")
	comAlvo(t, root, "B.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))
	_, _ = Enqueue(root, task("2-spec-b", "B.spec.md", "spec", "implement"))
	// Dois workers pegam as duas — HÁ MUITO TEMPO (fora da janela de trabalho), que é o
	// caso do worker que morreu sem fechar.
	antigo := time.Now().Add(-2 * JanelaDeTrabalho).Format(time.RFC3339)
	_, _ = Claim(root, "worker-morto-1", antigo)
	_, _ = Claim(root, "worker-morto-2", antigo)
	n, err := Reclaim(root)
	if err != nil || n != 2 {
		t.Fatalf("reclaim: n=%d err=%v (esperava 2)", n, err)
	}
	// ambas voltaram a pending e sem claimed_by
	tasks, _ := List(root)
	for _, tk := range tasks {
		if tk.State != Pending {
			t.Errorf("task %s deveria estar pending, veio %s", tk.ID, tk.State)
		}
		if tk.ClaimedBy != "" {
			t.Errorf("task %s deveria ter claimed_by limpo, veio %q", tk.ID, tk.ClaimedBy)
		}
	}
	// e são puxáveis de novo
	got, _ := Claim(root, "worker-novo", "2026-08-13T00:00:00-03:00")
	if got == nil {
		t.Fatal("task recuperada deveria ser puxável por Claim")
	}
}

func TestSuggestNext(t *testing.T) {
	cases := map[string]string{
		// O plano PROMOVIDO (já revisado, em `plans/`) semeia trabalho; o RASCUNHO
		// (`plans/review/`) vai para o review antes. São dois kinds porque o estado é a
		// PASTA — sem isso, editar o plano durante a execução dispararia review de novo a
		// cada vez, até o passo virar ruído.
		// Os verbos são os ARTEFATOS do `anchors work` — ver ArtefatosDeTrabalho. Já
		// foram `specify`/`implement`/`verify`, que o `work` recusa: quem puxava a task
		// não conseguia compor o prompt e traduzia por conta própria, todas as vezes.
		"plan":       "spec",
		"plan-draft": "review-plan-draft", "spec": "code", "feature": "test",
		// `test` fecha a trinca — e é aí que o trabalho PARECE pronto. A cadeia não
		// termina em verificar: chama o REVIEW. Medido em três rodadas de um E2E real, 7
		// defeitos graves passaram com todos os gates verdes; nenhum foi achado por gate.
		"code": "feature", "test": "review", "guide": "review-governed",
		"mistério": "triage",
	}
	for kind, want := range cases {
		if got, _ := SuggestNext(kind); got != want {
			t.Errorf("SuggestNext(%q) = %q, quer %q", kind, got, want)
		}
	}
}

// O `reclaim` RESPEITA quem pegou a task há pouco. Sem isso ele devolve tudo — inclusive
// o que um worker ATIVO está fazendo — e dois agentes passam a escrever no mesmo arquivo
// sem saber. Aconteceu, medido: um subagente rodou 90 minutos numa etapa, pareceu travado
// a quem observava de fora, alguém deu reclaim, e o trabalho foi duplicado.
func TestReclaimRespeitaClaimRecente(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "A.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))
	agora := time.Now().Format(time.RFC3339)
	if _, err := Claim(root, "worker-ativo", agora); err != nil {
		t.Fatal(err)
	}

	n, err := Reclaim(root)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("reclaim devolveu %d task(s) de um worker ATIVO — é assim que dois agentes "+
			"acabam no mesmo arquivo", n)
	}

	// `--force` é a saída explícita de quem SABE que o worker parou.
	if n, err := ReclaimForce(root); err != nil || n != 1 {
		t.Fatalf("reclaim --force: n=%d err=%v (esperava 1)", n, err)
	}
}

// Task sem carimbo de quando foi reivindicada não tem o que respeitar — devolve.
func TestReclaimSemCarimboDevolve(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "A.spec.md")
	_, _ = Enqueue(root, task("1-spec-a", "A.spec.md", "spec", "implement"))
	if _, err := Claim(root, "worker-sem-carimbo", ""); err != nil {
		t.Fatal(err)
	}
	if n, _ := Reclaim(root); n != 1 {
		t.Fatalf("claim sem carimbo deveria ser devolvido, veio %d", n)
	}
}

// TestSugestaoDaFilaEhComponivelPeloWork trava a divergência que quebrou o roteamento num
// E2E real: a fila sugeria `specify`/`implement`/`verify`, e o `anchors work` recusa os
// três ("artefato desconhecido"). Quem puxava a task não conseguia compor o prompt e tinha
// de traduzir os verbos por conta própria — o único ponto onde o ciclo não se auto-roteia.
//
// Dois vocabulários para a mesma coisa no mesmo binário é a forma mais barata do modo de
// falha mais caro deste framework: duas fontes da régua discordando.
func TestSugestaoDaFilaEhComponivelPeloWork(t *testing.T) {
	// todo kind que a fila reconhece precisa sugerir um verbo que o `work` aceita
	kinds := []string{"plan", "spec", "feature", "code", "test"}
	for _, k := range kinds {
		verbo, porque := SuggestNext(k)
		if verbo == "" {
			continue // kind sem próxima etapa é legítimo
		}
		if !ArtefatoDeTrabalhoValido(verbo) {
			t.Errorf("kind %q sugere %q, que o `anchors work` recusa — quem puxa a task não "+
				"consegue compor o prompt", k, verbo)
		}
		if porque == "" {
			t.Errorf("kind %q sugere %q sem dizer por quê", k, verbo)
		}
	}
}

// TestTaskDeAlvoInexistenteEhDescartada guarda o ruído medido em quatro execuções reais:
// o watcher enfileira na MUDANÇA e nunca desenfileira na DELEÇÃO, então a sonda de um
// revisor (`__probe.test.tsx`), apagada logo depois, deixava uma task viva para sempre.
//
// O orquestrador tinha de descartá-la à mão em toda rodada — e uma fila que exige triagem
// manual deixa de ser fila.
func TestTaskDeAlvoInexistenteEhDescartada(t *testing.T) {
	root := t.TempDir()
	comAlvo(t, root, "vive.ts")
	if _, err := Enqueue(root, task("1-code-vive", "vive.ts", "code", "feature")); err != nil {
		t.Fatal(err)
	}
	if _, err := Enqueue(root, task("2-code-sonda", "sonda.test.ts", "code", "review")); err != nil {
		t.Fatal(err)
	}

	tasks, err := List(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].Changed != "vive.ts" {
		t.Fatalf("a task do alvo apagado devia sumir; veio %+v", tasks)
	}
	// E some do DISCO, não só da listagem: filtrar sem remover faria a fantasma
	// reaparecer no `list` seguinte.
	restantes, _ := os.ReadDir(filepath.Join(root, ".anchors", "tasks"))
	for _, e := range restantes {
		if strings.Contains(e.Name(), "sonda") {
			t.Error("a task-fantasma continua no disco — reapareceria na próxima listagem")
		}
	}
}
