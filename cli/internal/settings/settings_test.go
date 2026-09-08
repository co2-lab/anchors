package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AUSÊNCIA NÃO É ERRO, e é o caso mais comum: um projeto recém-clonado não tem o arquivo,
// e é exatamente quando o agente precisa perguntar.
func TestLoad_semArquivoNaoEhErro(t *testing.T) {
	s, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("projeto sem `%s` deu erro: %v", File, err)
	}
	if s.Decided() {
		t.Error("sem arquivo, a decisão não pode contar como tomada")
	}
	if s.HandlesUserIssues() {
		t.Error("o padrão é FECHADO — quem pode decidir declara que pode")
	}
}

// TRÊS ESTADOS, e a distinção é o mecanismo: `nil` é a única situação em que o agente
// pergunta. Sem o ponteiro, "não declarado" e "declarado como não" seriam o mesmo valor, e
// o agente perguntaria a cada sessão a quem já disse não.
func TestSettings_tresEstados(t *testing.T) {
	casos := []struct {
		nome     string
		valor    *bool
		decidido bool
		atua     bool
	}{
		{"nunca perguntei", nil, false, false},
		{"disse que não", Bool(false), true, false},
		{"disse que sim", Bool(true), true, true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			s := Settings{UserIssues: c.valor}
			if s.Decided() != c.decidido {
				t.Errorf("Decided() = %v, queria %v", s.Decided(), c.decidido)
			}
			if s.HandlesUserIssues() != c.atua {
				t.Errorf("HandlesUserIssues() = %v, queria %v", s.HandlesUserIssues(), c.atua)
			}
		})
	}
}

// O que foi gravado é o que se lê de volta.
func TestSaveLoad_daVoltaOQueFoiGravado(t *testing.T) {
	root := t.TempDir()
	if err := Save(root, Settings{
		UserIssues: Bool(true), Agent: "maquina/sessao", DecidedAt: "2026-09-08",
	}); err != nil {
		t.Fatal(err)
	}

	s, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !s.HandlesUserIssues() || s.Agent != "maquina/sessao" || s.DecidedAt != "2026-09-08" {
		t.Errorf("voltou %+v", s)
	}
}

// O arquivo EXPLICA o que é, porque quem o encontra depois não estava na conversa.
func TestSave_oArquivoDizOQueEle(t *testing.T) {
	root := t.TempDir()
	if err := Save(root, Settings{UserIssues: Bool(false)}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(Path(root))
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	for _, esperado := range []string{"não vai para o git", "gitignore", "escalonados"} {
		if !strings.Contains(texto, esperado) {
			t.Errorf("o cabeçalho não menciona %q:\n%s", esperado, texto)
		}
	}
}

// O arquivo vive em `.anchors/`, que já está no `.gitignore` dos projetos — é o que o
// mantém local, e é a razão de ele não ficar ao lado do `anchors.yaml`.
func TestPath_ficaEmAnchors(t *testing.T) {
	if got := Path("/proj"); got != filepath.Join("/proj", ".anchors", "settings.yaml") {
		t.Errorf("Path = %q", got)
	}
}

// A resposta é lida nas DUAS línguas e nas formas curtas: a pergunta é no terminal, e quem
// responde digita o que lhe vem primeiro.
func TestParseAnswer(t *testing.T) {
	sim := []string{"s", "sim", "y", "yes", "SIM", " Sim \n"}
	nao := []string{"n", "nao", "não", "no", "NÃO", " nao \n"}
	naoEntendi := []string{"", "talvez", "1", "ok", "quem sabe"}

	for _, r := range sim {
		if v := ParseAnswer(r); v == nil || !*v {
			t.Errorf("ParseAnswer(%q) não foi lido como sim", r)
		}
	}
	for _, r := range nao {
		if v := ParseAnswer(r); v == nil || *v {
			t.Errorf("ParseAnswer(%q) não foi lido como não", r)
		}
	}
	// O que não se entende devolve `nil` — e quem chama pergunta de novo. Assumir "não"
	// aqui seria conveniente e errado: quem digitou algo estava respondendo, e descartar a
	// resposta em silêncio faz o agente decidir por conta.
	for _, r := range naoEntendi {
		if v := ParseAnswer(r); v != nil {
			t.Errorf("ParseAnswer(%q) = %v, queria nil", r, *v)
		}
	}
}
