package health

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
)

func cfgGitHub() *config.Config {
	return &config.Config{Workflow: &config.Workflow{
		Mode:   config.ModeGitHub,
		Repo:   "acme/exemplo",
		Labels: []string{"anchors"},
	}}
}

// Cobrar board e pipelines de um projeto que declarou `mode: local` seria ruído garantido
// — e ruído recorrente treina a equipe a ignorar o doctor, que é o oposto do que ele quer.
func TestAmbienteNaoCobraNadaNoModoLocal(t *testing.T) {
	dir := t.TempDir()

	if fs := checkAmbienteGitHub(&config.Config{}, dir); len(fs) != 0 {
		t.Errorf("modo local não usa board nem pipelines: %+v", fs)
	}
}

// Um pipeline ausente não gera erro em lugar nenhum — gera SILÊNCIO, e os artefatos ficam
// sem card para sempre. É por isso que o doctor tem de avisar: é a única coisa que
// aparece.
func TestAmbienteAvisaPipelineAusente(t *testing.T) {
	dir := t.TempDir()

	fs := checkPipelines(dir)

	if len(fs) != len(initx.WorkflowsDoFluxo) {
		t.Fatalf("esperava %d achados, veio %d", len(initx.WorkflowsDoFluxo), len(fs))
	}
	for _, f := range fs {
		if f.Check != "pipeline-ausente" || f.Severity != Warn {
			t.Errorf("achado errado: %+v", f)
		}
		// A mensagem tem de dizer o que fica SEM ACONTECER — "falta um arquivo" não
		// explica por que isso importa.
		if !strings.Contains(f.Detail, "não acontece") {
			t.Errorf("a mensagem deveria nomear o que deixa de acontecer: %s", f.Detail)
		}
		if !strings.Contains(f.Detail, "--fix") {
			t.Errorf("a mensagem deveria apontar o conserto: %s", f.Detail)
		}
	}
}

// O pior caso: o pipeline ESTÁ lá, então parece configurado — mas sem `concurrency` ele
// roda em paralelo consigo mesmo e devolve a corrida que existia para eliminar. Achado
// próprio, e nunca contado como "ausente".
func TestAmbientePegaPipelineSemSerializacao(t *testing.T) {
	dir := t.TempDir()
	if _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}
	// Estraga um: tira a serialização, mantém o arquivo.
	alvo := filepath.Join(dir, initx.DirWorkflows, "anchors-claim.yml")
	if err := os.WriteFile(alvo, []byte("name: claim\non:\n  workflow_dispatch:\njobs:\n  x:\n    runs-on: ubuntu-latest\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs := checkPipelines(dir)

	if len(fs) != 1 {
		t.Fatalf("esperava 1 achado (o pipeline quebrado), veio %d: %+v", len(fs), fs)
	}
	if fs[0].Check != "pipeline-sem-serializacao" {
		t.Errorf("o arquivo existe — não pode ser reportado como ausente: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "mesmo card a dois agentes") {
		t.Errorf("a mensagem deveria dizer o que quebra na prática: %s", fs[0].Detail)
	}
}

// Depois do `--fix` não pode sobrar achado de pipeline: o que o doctor cobra e o que o
// `--fix` cria saem da MESMA lista, e discordarem faria o doctor pedir eternamente algo
// que ele mesmo acabou de criar.
func TestFixSatisfazOQueODoctorCobra(t *testing.T) {
	dir := t.TempDir()
	if fs := checkPipelines(dir); len(fs) == 0 {
		t.Fatal("preparo: o projeto vazio deveria ter achados")
	}

	if _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}

	if fs := checkPipelines(dir); len(fs) != 0 {
		t.Errorf("o --fix criou os pipelines e o doctor ainda reclama: %+v", fs)
	}
}

// O board NÃO é conferido, e é deliberado: o estado do trabalho é uma label, e o
// Project é espelho opcional (BOOTSTRAP.md §7.13). A decisão anterior — estado na coluna
// — exigia um PAT com escopo `project` de todo adotante, e isso era atrito de adoção por
// uma escolha de visualização.
//
// Este teste é o que impede a volta: cobrar board faria o doctor pedir uma peça que o
// fluxo não usa.
func TestDoctorNaoCobraBoard(t *testing.T) {
	dir := t.TempDir()
	if _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}

	fs := checkAmbienteGitHub(cfgGitHub(), dir)

	for _, f := range fs {
		if strings.Contains(strings.ToLower(f.Check), "board") ||
			strings.Contains(strings.ToLower(f.Detail), "project") {
			t.Errorf("o doctor cobra board, que virou opcional: %+v", f)
		}
	}
}
