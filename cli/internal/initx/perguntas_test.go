package initx

import (
	"strings"
	"testing"
)

func qsDeTeste() []Pergunta {
	return Perguntas(&Proposal{Config: nil}, []string{"go", "nextjs"})
}

// O contrato existe para um agente responder sem ver a TUI. Cada pergunta precisa trazer
// o que ele precisa para DECIDIR: o que é aceito, o que o Anchors inferiu, e o que a
// resposta muda no projeto.
func TestCadaPerguntaTrazOQueOAgentePrecisaParaDecidir(t *testing.T) {
	for _, q := range qsDeTeste() {
		if q.ID == "" || q.Texto == "" || q.Tipo == "" {
			t.Errorf("pergunta incompleta: %+v", q)
		}
		// Sem o "por que", o agente escolhe pelo NOME da opção — que é como se escolhe
		// errado. É a diferença entre responder e adivinhar.
		if q.PorQue == "" {
			t.Errorf("%s não diz o que a resposta muda no projeto", q.ID)
		}
		if q.Tipo == "select" && len(q.Opcoes) == 0 {
			t.Errorf("%s é select e não oferece opções", q.ID)
		}
	}
}

// As sete decisões da TUI têm de estar todas aqui: uma pergunta que ficasse de fora
// seria decidida em silêncio pelo default, num arquivo que o usuário acha que decidiu.
func TestTodasAsDecisoesDaTUIEstaoNoContrato(t *testing.T) {
	esperadas := []string{"preset", "header", "artifacts", "gates", "colocation", "layers",
		"governs", "workflow", "repo", "labels"}
	achadas := map[string]bool{}
	for _, q := range qsDeTeste() {
		achadas[q.ID] = true
	}
	for _, e := range esperadas {
		if !achadas[e] {
			t.Errorf("falta a pergunta %q — ela seria decidida em silêncio", e)
		}
	}
}

// Uma resposta inválida recusa o CONJUNTO. Escrever as válidas produziria um
// anchors.yaml que ninguém decidiu por completo — e um arquivo assim carrega sem erro,
// governa errado, e não acusa a causa.
func TestRespostaInvalidaRecusaTudo(t *testing.T) {
	qs := qsDeTeste()
	ruim := "preset-que-nao-existe"

	st := ValidaRespostas(qs, Respostas{Preset: &ruim})

	if TudoAceito(st) {
		t.Fatal("preset inexistente deveria recusar o conjunto")
	}
	// E o status tem de dizer QUAL falhou e por quê — não basta "algo deu errado".
	var achou bool
	for _, s := range st {
		if s.ID == "preset" {
			achou = true
			if s.Aceita {
				t.Error("o preset inválido foi aceito")
			}
			if !strings.Contains(s.Detalhe, "aceitos:") {
				t.Errorf("a mensagem deveria listar os valores válidos: %s", s.Detalhe)
			}
		}
	}
	if !achou {
		t.Error("o status do preset não foi reportado")
	}
}

// O status é de TODAS as respostas, não só das inválidas. É o que permite ao agente
// conferir que o Anchors entendeu o que ele quis dizer — uma flag escrita errada, sem
// isso, seria indistinguível de uma resposta aceita.
func TestStatusReportaTodasAsRespostas(t *testing.T) {
	qs := qsDeTeste()
	sim := true

	st := ValidaRespostas(qs, Respostas{Header: &sim})

	if len(st) != len(qs) {
		t.Fatalf("esperava %d status (um por pergunta), veio %d", len(qs), len(st))
	}
	for _, s := range st {
		if s.ID == "header" && s.UsouPada {
			t.Error("header foi respondido explicitamente — não é default")
		}
		if s.ID == "gates" && !s.UsouPada {
			t.Error("gates não foi respondido — tinha de constar como default")
		}
	}
}

// A distinção que motiva os ponteiros: "não respondi" (vale o default inferido do disco)
// e "respondi vazio" (`--artifacts=""`, nenhum artefato) são decisões OPOSTAS, e um bool
// zero-value não as separa.
func TestNaoRespondidoNaoEhOMesmoQueRespondidoVazio(t *testing.T) {
	qs := Perguntas(&Proposal{}, nil)
	vazio := []string{}

	semResposta := ValidaRespostas(qs, Respostas{})
	comVazio := ValidaRespostas(qs, Respostas{Artifacts: &vazio})

	for _, s := range semResposta {
		if s.ID == "artifacts" && !s.UsouPada {
			t.Error("sem a flag, artifacts tem de cair no default")
		}
	}
	for _, s := range comVazio {
		if s.ID == "artifacts" && s.UsouPada {
			t.Error("`--artifacts=` vazio é uma escolha deliberada, não ausência de resposta")
		}
	}
}

// O modo de trabalho é decisão HUMANA (onde a fila mora) e o `init` não a fazia: todo
// projeto nascia `local` por omissão, e quem queria `github` tinha de descobrir o campo e
// editar o YAML à mão.
func TestModoDeTrabalhoEhPerguntado(t *testing.T) {
	var achou bool
	for _, q := range qsDeTeste() {
		if q.ID == "workflow" {
			achou = true
			if !contem(q.Opcoes, "local") || !contem(q.Opcoes, "github") {
				t.Errorf("os dois modos têm de ser oferecidos: %v", q.Opcoes)
			}
			if q.Default != "local" {
				t.Errorf("o default é local (bloco ausente = local), veio %v", q.Default)
			}
		}
	}
	if !achou {
		t.Fatal("o modo de trabalho não é perguntado — o projeto nasce local sem ninguém decidir")
	}
}

// No modo `github`, `repo` e `labels` são obrigatórios (WORKFLOW.md §2). Sem repo, a
// escrita cairia no lugar errado; sem label, o fluxo pegaria issue de produto.
func TestGitHubExigeRepoELabels(t *testing.T) {
	qs := qsDeTeste()
	github := "github"

	st := ValidaRespostas(qs, Respostas{Workflow: &github})

	if TudoAceito(st) {
		t.Fatal("modo github sem repo nem labels deveria recusar")
	}
	faltando := map[string]bool{}
	for _, s := range st {
		if !s.Aceita {
			faltando[s.ID] = true
		}
	}
	for _, id := range []string{"repo", "labels"} {
		if !faltando[id] {
			t.Errorf("%s é obrigatório no modo github e passou", id)
		}
	}
}

// No modo `local`, declarar `repo` não é campo inofensivo: quem lê o arquivo conclui que
// a integração está ativa (WORKFLOW.md §2).
func TestLocalRecusaCamposDoGitHub(t *testing.T) {
	qs := qsDeTeste()
	local, repo := "local", "owner/nome"

	st := ValidaRespostas(qs, Respostas{Workflow: &local, Repo: &repo})

	if TudoAceito(st) {
		t.Error("`repo` no modo local faz o arquivo mentir sobre a integração estar ativa")
	}
}

// O caminho feliz: github com as duas exigências satisfeitas passa.
func TestGitHubCompletoEhAceito(t *testing.T) {
	qs := qsDeTeste()
	github, repo := "github", "acme/exemplo"
	labels := []string{"anchors"}

	st := ValidaRespostas(qs, Respostas{Workflow: &github, Repo: &repo, Labels: &labels})

	if !TudoAceito(st) {
		for _, s := range st {
			if !s.Aceita {
				t.Errorf("%s recusado: %s", s.ID, s.Detalhe)
			}
		}
	}
}
