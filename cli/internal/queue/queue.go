// Package queue é a fila de tarefas do Anchors: o mecanismo que desacopla "algo
// mudou" de "alguém trabalha nisso". O watcher ENFILEIRA (Enqueue) quando classifica
// uma mudança; um worker PUXA (Claim) o próximo item, executa um passo, e marca DONE.
//
// A fila é materializada em .anchors/tasks/ — arquivos YAML inspecionáveis, um por
// task. O claim é ATÔMICO via rename (atômico no POSIX), então dois workers em
// terminais diferentes nunca pegam a mesma task. Isso é o que permite paralelismo
// manual em QUALQUER cliente de IA (dois terminais, dois `anchors next`), sem
// depender de subagente em background.
package queue

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Dir é a subpasta (sob .anchors/) onde as tasks VIVAS vivem (pending + claimed).
// DoneDir é para onde vão as concluídas — o histórico, fora da fila de trabalho.
const (
	Dir     = ".anchors/tasks"
	DoneDir = ".anchors/done"
)

// State é o ciclo de vida de uma task. Refletido no NOME do arquivo (não só no
// conteúdo) para que o claim seja um rename atômico e a listagem seja barata.
type State string

const (
	Pending State = "pending" // enfileirada, ninguém pegou
	Claimed State = "claimed" // um worker reivindicou; está trabalhando
	Done    State = "done"    // concluída
)

// Task é uma unidade de trabalho: uma mudança que pede um passo do ciclo de vida.
type Task struct {
	ID string `yaml:"id"` // estável: <seq>-<kind>-<slug>

	// O fato observado.
	Changed string `yaml:"changed"` // arquivo (relativo) que mudou
	Kind    string `yaml:"kind"`    // kind do arquivo (spec|feature|code|test|plan|...)
	Origin  string `yaml:"origin"`  // "watch" | "manual" — quem enfileirou

	// A SUGESTÃO (não ordem): o watcher deriva o próximo passo do kind/estrutura,
	// mas o worker (a IA) tem autonomia para divergir. Ver o playbook.
	// SuggestedNext é o VERBO da próxima etapa, e ele tem de ser um artefato que o
	// `anchors work` aceita — senão quem puxa a task não consegue compor o prompt.
	//
	// Medido num E2E real: a fila emitia `specify`/`implement`/`verify`, que o `work`
	// recusa ("artefato desconhecido"). O orquestrador teve de traduzir os três verbos
	// por conta própria, todas as vezes — o único ponto onde o ciclo não se auto-roteia.
	// Dois vocabulários para a mesma coisa, no mesmo binário.
	SuggestedNext string `yaml:"suggested_next"` // ex: "spec", "code", "test", "review"
	Reason        string `yaml:"reason"`         // por que este é o próximo passo

	State     State  `yaml:"state"`
	CreatedAt string `yaml:"created_at"`           // RFC3339; carimbado por quem enfileira
	ClaimedBy string `yaml:"claimed_by,omitempty"` // identificador do worker
	// ClaimedAt é QUANDO a task foi reivindicada. É o que permite ao `reclaim` respeitar
	// quem está trabalhando: o worker do Anchors não é um processo observável (o `anchors
	// next` morre ao imprimir a task; quem trabalha depois é o agente), então a única
	// evidência disponível é há quanto tempo alguém a pegou.
	ClaimedAt string `yaml:"claimed_at,omitempty"`
}

// ArtefatosDeTrabalho são os verbos que o `anchors work` sabe compor. A fila só pode
// sugerir um destes: quem puxa a task usa o verbo para montar o prompt, e um verbo que o
// `work` recusa deixa o ciclo sem roteamento — quem estiver executando tem de adivinhar a
// tradução, todas as vezes.
//
// A lista vive AQUI, e não só no comando, porque é aqui que ela é honrada: `work` valida
// contra ela e um teste confronta as duas pontas. Já divergiram em silêncio uma vez.
var ArtefatosDeTrabalho = []string{
	"spec", "code", "feature", "test", "review", "review-plan", "review-plan-draft",
}

// ArtefatoDeTrabalhoValido diz se o verbo é composível pelo `anchors work`.
func ArtefatoDeTrabalhoValido(v string) bool {
	for _, a := range ArtefatosDeTrabalho {
		if a == v {
			return true
		}
	}
	return false
}

// SuggestNext mapeia o kind de um arquivo que mudou para o próximo passo do ciclo de
// vida. É uma SUGESTÃO determinística (spec→implementar, feature→testar, …); a IA
// pode divergir. Deriva só do kind — a estrutura fina (via impact) fica com o worker.
func SuggestNext(kind string) (next, reason string) {
	switch kind {
	case "plan-draft":
		// Um plano em REVISÃO (`plans/review/`) ainda não é roteiro de ninguém: vai para o
		// revisor, não para o executor.
		//
		// O ESTADO É A PASTA — a mesma mecânica das issues (`todo/doing/done`) e da fila
		// (`pending__`/`claimed__`): mover é mudar de estado, e o estado é visível sem
		// abrir o arquivo.
		//
		// Por que o estado precisa existir: o watcher vê "um arquivo mudou" e nada mais.
		// Sem separar os dois, ele não distingue "o plano nasceu" de "o plano foi editado
		// durante a execução" — e marcar um item como feito dispararia review de novo a
		// cada vez, até o passo virar ruído e alguém desligá-lo.
		//
		// Promover (mover de `plans/review/` para `plans/`) é o gesto que diz "revisei e
		// aprovo". Um gesto, não um campo a manter.
		return "review-plan-draft", "um plano em RASCUNHO mudou — REVISE-O antes de promovê-lo " +
			"(`anchors work review-plan-draft --for <plano>`): item cuja EXISTÊNCIA depende de " +
			"decisão não tomada faz quem executa decidir sozinho ou parar o fluxo. Aprovado, " +
			"renomeie tirando `.draft` — aí ele semeia trabalho"
	case "plan":
		// Plano PROMOVIDO (já revisado) semeia trabalho. O custo de não revisar antes é
		// assimétrico: um defeito na spec afeta uma unidade; um no plano se espalha por
		// todas as que ele semeia — medido, um item "nasce SE …" fez o executor decidir
		// sozinho, e uma fase inteira nasceu sobre premissa que ninguém aprovou.
		return "spec", "um plano foi semeado — gere as specs que ele lista, uma por alvo " +
			"(`anchors work spec --for <alvo>`)"
	case "spec":
		return "code", "uma spec mudou — crie/atualize o código e a feature " +
			"(`anchors work code --for <alvo>`)"
	case "feature":
		return "test", "uma feature mudou — escreva/atualize os testes"
	case "code":
		// Depois do código vem a FEATURE: é ela que descreve o comportamento em cenários, e
		// é dela que os testes nascem (`feature` → `test` logo acima). Confrontar com
		// `anchors check` é parte de toda etapa, não uma etapa — sugerir "check" aqui
		// mandava o executor para um comando que o `work` não sabe compor.
		return "feature", "o código mudou — descreva o comportamento em cenários " +
			"(`anchors work feature --for <alvo>`), que é de onde os testes nascem"
	case "test":
		// O teste é a ÚLTIMA peça da trinca a nascer: quando ele muda, a unidade está
		// completa e é exatamente aí que o trabalho PARECE pronto. Por isso a cadeia não
		// termina em `verify-tests` — ela chama o REVIEW.
		//
		// Medido em três rodadas de um E2E real: 7 defeitos graves (perda silenciosa de
		// dado do usuário, tabela inexistente na infra, regra sem teste que a prove,
		// contradição entre duas regras da mesma spec) passaram com TODOS os gates verdes.
		// Nenhum foi achado por gate; os 7 saíram de revisão adversarial. O gate confronta
		// o que é DECLARÁVEL; o que sobra precisa de alguém atacando de fora.
		return "review", "a trinca fechou — rode a suíte, `anchors ingest` os sinais, e então " +
			"REVISE a unidade inteira (`anchors work review --for <alvo>`): os gates verdes " +
			"não provam que está certo"
	case "guide":
		return "review-governed", "uma régua mudou — revise os artefatos que ela rege"
	default:
		return "triage", "mudança de kind não mapeado — decida o próximo passo"
	}
}

// --- persistência ---

func dirFor(root string) string     { return filepath.Join(root, Dir) }
func doneDirFor(root string) string { return filepath.Join(root, DoneDir) }

// fileName codifica o estado no nome: <state>__<id>.yaml. Claim = rename entre
// prefixos de estado (atômico). done idem.
func fileName(state State, id string) string {
	return fmt.Sprintf("%s__%s.yaml", state, id)
}

func parseFileName(name string) (state State, id string, ok bool) {
	if !strings.HasSuffix(name, ".yaml") {
		return "", "", false
	}
	base := strings.TrimSuffix(name, ".yaml")
	st, id, found := strings.Cut(base, "__")
	if !found {
		return "", "", false
	}
	return State(st), id, true
}

// Enqueue grava uma nova task no estado pending. Idempotente por ID: se já existe
// uma task (em qualquer estado) para o mesmo (changed, suggested_next) ainda não
// concluída, NÃO duplica — evita enxurrada quando um arquivo é salvo várias vezes.
func Enqueue(root string, t Task) (created bool, err error) {
	d := dirFor(root)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return false, err
	}
	// dedup: já há task viva (pending|claimed) para o mesmo alvo+passo?
	existing, err := List(root)
	if err != nil {
		return false, err
	}
	for _, e := range existing {
		if e.State != Done && e.Changed == t.Changed && e.SuggestedNext == t.SuggestedNext {
			return false, nil // já enfileirada; não duplica
		}
	}
	t.State = Pending
	data, err := yaml.Marshal(t)
	if err != nil {
		return false, err
	}
	path := filepath.Join(d, fileName(Pending, t.ID))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// List devolve todas as tasks (qualquer estado), ordenadas por ID (estável).
func List(root string) ([]Task, error) {
	d := dirFor(root)
	entries, err := os.ReadDir(d)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Task
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		state, _, ok := parseFileName(e.Name())
		if !ok {
			continue
		}
		data, err := os.ReadFile(filepath.Join(d, e.Name()))
		if err != nil {
			continue
		}
		var t Task
		if yaml.Unmarshal(data, &t) != nil {
			continue
		}
		t.State = state // o nome do arquivo é a fonte da verdade do estado
		// Task cujo ALVO não existe mais é ruído: o arquivo foi apagado (uma sonda de
		// revisor, um rascunho, um rename) e a task ficou. O watcher enfileira na
		// mudança e nunca desenfileira na deleção.
		//
		// Medido em quatro execuções seguidas: sempre sobravam tasks-fantasma que o
		// orquestrador tinha de descartar à mão — e uma fila que exige triagem manual
		// deixa de ser fila. Remover aqui (e não só filtrar) evita que ela reapareça no
		// `list` seguinte.
		if t.Changed != "" && !existeNaRaiz(root, t.Changed) {
			_ = os.Remove(filepath.Join(d, e.Name()))
			continue
		}
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Claim reivindica a próxima task pending, atomicamente. Percorre as pending em ordem e
// tenta CRIAR claimed__X com O_EXCL; a criação exclusiva é o que elege um só dono, e o
// pending__X só é removido depois. Devolve a task reivindicada, ou (nil, nil) se a fila
// está vazia.
//
// Antes isto era um rename pending__X → claimed__X, apoiado em "o rename é atômico, só
// um vence". É verdade no POSIX e FALSO no Windows: medido com 8 goroutines disputando o
// mesmo arquivo, 167 de 200 rodadas tiveram mais de um vencedor — uma delas, os 8. Dois
// terminais rodando `anchors next` pegavam a MESMA task e trabalhavam em cima um do
// outro. O_EXCL (CREATE_NEW no Windows) é atômico nos dois sistemas.
//
// A ordem — criar o claimed, depois apagar o pending — deixa como pior caso um pending
// resíduo se o processo morrer no meio, nunca uma task reivindicada duas vezes. Esse
// resíduo se limpa sozinho: quem o encontrar vai achar o claimed já criado e o apaga.
func Claim(root, worker, now string) (*Task, error) {
	d := dirFor(root)
	tasks, err := List(root)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.State != Pending {
			continue
		}
		from := filepath.Join(d, fileName(Pending, t.ID))
		to := filepath.Join(d, fileName(Claimed, t.ID))
		// A criação exclusiva é a disputa: quem cria o claimed__X é o dono.
		f, err := os.OpenFile(to, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			// Já existe claimed__X: outro worker pegou. O pending que sobrou é resíduo
			// da morte de alguém entre criar e apagar — apagar aqui é o que o dono faria.
			if os.IsExist(err) {
				_ = os.Remove(from)
			}
			continue
		}
		t.State = Claimed
		t.ClaimedBy, t.ClaimedAt = worker, now
		data, _ := yaml.Marshal(t)
		_, werr := f.Write(data)
		cerr := f.Close()
		if werr != nil || cerr != nil {
			// Sem registrar a posse a task ficaria presa em claimed vazio, invisível para
			// a fila e sem dono. Desfaz e deixa o pending para o próximo.
			_ = os.Remove(to)
			continue
		}
		_ = os.Remove(from)
		return &t, nil
	}
	return nil, nil // fila vazia
}

// MarkDone tira uma task da fila viva (.anchors/tasks/) e a move para o histórico
// (.anchors/done/). Assim `queue`/`next` só enxergam trabalho vivo; o done fica
// arquivado, inspecionável, fora do caminho. Procura em claimed e (raro) pending.
func MarkDone(root, id string) error {
	d := dirFor(root)
	done := doneDirFor(root)
	if err := os.MkdirAll(done, 0o755); err != nil {
		return err
	}
	for _, from := range []string{fileName(Claimed, id), fileName(Pending, id)} {
		src := filepath.Join(d, from)
		if _, err := os.Stat(src); err == nil {
			return os.Rename(src, filepath.Join(done, fileName(Done, id)))
		}
	}
	return fmt.Errorf("task não encontrada (pending/claimed): %s", id)
}

// Drop descarta uma task da fila viva SEM concluí-la — remove o arquivo (pending ou
// claimed). Para lixo: planos-doc que viraram triage, tasks obsoletas, duplicatas.
// Diferente de MarkDone (que arquiva em done/): Drop apaga, não é histórico de
// trabalho feito.
func Drop(root, id string) error {
	d := dirFor(root)
	for _, name := range []string{fileName(Pending, id), fileName(Claimed, id)} {
		src := filepath.Join(d, name)
		if _, err := os.Stat(src); err == nil {
			return os.Remove(src)
		}
	}
	return fmt.Errorf("task não encontrada (pending/claimed): %s", id)
}

// Reclaim devolve à fila (claimed → pending) as tasks presas em claimed — tipicamente
// órfãs de um worker que morreu sem chamar done. Devolve quantas recuperou. O usuário
// roda isto após um crash para o trabalho não ficar preso indefinidamente.
// Reclaim devolve à fila as tasks presas em `claimed` por um worker que MORREU.
//
// A checagem de vida não é detalhe: sem ela, `reclaim` devolve tudo — inclusive o que um
// worker ATIVO está fazendo agora — e dois agentes passam a trabalhar no mesmo arquivo
// sem saber. Aconteceu, medido: um worker parecia travado (45 min sem escrever, porque
// um subagente longo estava rodando), alguém deu `reclaim`, e o trabalho foi duplicado.
// Deu certo por sorte.
//
// O ID do worker é `pid@hostname`. Só se decide sobre worker da MESMA máquina: de outra,
// não há como saber, e roubar por suposição é pior que deixar preso (o `--force` existe
// para quem sabe o que está fazendo).
func Reclaim(root string) (int, error) { return reclaim(root, false) }

// ReclaimForce devolve à fila TODA task claimed, viva ou não. É a saída para quando o
// worker está noutra máquina, ou travado de verdade — e por isso é explícita.
func ReclaimForce(root string) (int, error) { return reclaim(root, true) }

func reclaim(root string, force bool) (int, error) {
	d := dirFor(root)
	tasks, err := List(root)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, t := range tasks {
		if t.State != Claimed {
			continue
		}
		if !force && t.ClaimedBy != "" && !claimAntigo(t) {
			continue // reivindicado há pouco: provavelmente alguém está nisto agora
		}
		from := filepath.Join(d, fileName(Claimed, t.ID))
		to := filepath.Join(d, fileName(Pending, t.ID))
		if os.Rename(from, to) == nil {
			// limpa o claimed_by no conteúdo (voltou a ser de ninguém)
			t.State, t.ClaimedBy, t.ClaimedAt = Pending, "", ""
			data, _ := yaml.Marshal(t)
			_ = os.WriteFile(to, data, 0o644)
			n++
		}
	}
	return n, nil
}

// PendingCount conta as tasks ainda não concluídas (pending + claimed).
func PendingCount(root string) (int, error) {
	tasks, err := List(root)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, t := range tasks {
		if t.State != Done {
			n++
		}
	}
	return n, nil
}

// claimAntigo diz se o claim já passou da janela em que se presume que alguém trabalha.
//
// Por que TEMPO e não "o processo está vivo": o worker do Anchors não é um processo. O
// `anchors next` imprime a task e MORRE; quem trabalha depois é o agente (ou a pessoa),
// que o CLI não tem como observar. Checar o PID do `next` é checar um processo que já
// terminou — mediria sempre "morto", e o reclaim voltaria a roubar trabalho ativo.
//
// A janela é generosa de propósito. Medido: um subagente rodou 90 minutos numa única
// etapa (e pareceu travado a quem observava de fora, o que levou alguém a dar reclaim e
// duplicar o trabalho). Devolver cedo demais é o erro caro; devolver tarde só custa
// espera, e o `--force` resolve para quem tem certeza.
const JanelaDeTrabalho = 4 * time.Hour

func claimAntigo(t Task) bool {
	if t.ClaimedAt == "" {
		return true // sem carimbo: não há o que respeitar
	}
	quando, err := time.Parse(time.RFC3339, t.ClaimedAt)
	if err != nil {
		return true
	}
	return time.Since(quando) > JanelaDeTrabalho
}

// existeNaRaiz diz se o alvo de uma task ainda está no disco.
func existeNaRaiz(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, rel))
	return err == nil
}
