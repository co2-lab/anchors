package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// As duas constantes têm de ser a MESMA string.
//
// O `scan` mantém o arquivo fora do mapa; o comando o cria. Se divergirem, o `new` gera
// um arquivo com um sufixo e o mapa exclui outro — e o progresso volta a ser confrontado
// pelos gates, reintroduzindo em silêncio o defeito que a separação removeu. Nenhum teste
// de comportamento pegaria isso: os dois lados funcionariam, cada um com a sua régua.
func TestProgresso_sufixoBateComOScan(t *testing.T) {
	if !scan.EhArquivoDeProgresso("plans/0001-x" + sufixoProgresso) {
		t.Fatalf("o sufixo do comando (%q) não é reconhecido pelo scan — o `new` criaria "+
			"um arquivo que o mapa NÃO exclui, e os gates voltariam a confrontá-lo",
			sufixoProgresso)
	}
}

func TestProgresso_caminhoFicaAoLadoDoPlano(t *testing.T) {
	got := caminhoDeProgresso("plans/0017-mutacao.md")
	want := "plans/0017-mutacao" + sufixoProgresso
	if got != want {
		t.Fatalf("caminho: %q, queria %q", got, want)
	}
}

// O progresso nasce com uma seção por FASE DECLARADA, lida do cabeçalho.
//
// A fonte é a mesma que os gates `fase-existe` e `fase-ordenada` usam. Manter uma segunda
// lista faria o progresso falar de fases que não existem — e a divergência só apareceria
// quando alguém renomeasse uma fase.
func TestProgresso_umaSecaoPorFaseDoPlano(t *testing.T) {
	plano := `# Plano 0017

## Fases

### MTUAO-F01 — a ferramenta e o relatório

- ` + "`packages/shared/MutationHarness.spec.md`" + `

### MTUAO-F02 — o CI ingere o sinal (depende de MTUAO-F01)
`
	dir := t.TempDir()
	p := filepath.Join(dir, "0017-mutacao.md")
	if err := os.WriteFile(p, []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}

	destino, err := escreveProgressoInicial(p, plano, "MTUAO")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(destino)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)

	for _, fase := range []string{"MTUAO-F01", "MTUAO-F02"} {
		if !strings.Contains(got, "## "+fase) {
			t.Errorf("falta a seção da fase %s:\n%s", fase, got)
		}
	}
	// O TÍTULO da fase vem junto: sem ele o arquivo é uma lista de códigos, e quem o abre
	// tem de voltar ao plano para saber do que cada um trata.
	if !strings.Contains(got, "a ferramenta e o relatório") {
		t.Errorf("o título da fase não foi copiado:\n%s", got)
	}
	// O checkbox mora AQUI, e é o ponto inteiro da separação.
	if !strings.Contains(got, "- [ ]") {
		t.Errorf("o progresso nasce sem checkbox — é onde o `[x]` deve ser marcado:\n%s", got)
	}
}

// NÃO SOBRESCREVE: o arquivo guarda estado, e regravá-lo apagaria o trabalho registrado.
func TestProgresso_naoSobrescreveEstadoExistente(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "0001-x.md")
	if err := os.WriteFile(p, []byte("# Plano\n\n### ABCDE-F01 — fase\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jaFeito := "# Progresso — ABCDE\n\n## ABCDE-F01\n\n- [x] feito\n"
	if err := os.WriteFile(caminhoDeProgresso(p), []byte(jaFeito), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := escreveProgressoInicial(p, "# Plano\n\n### ABCDE-F01 — fase\n", "ABCDE"); err == nil {
		t.Fatal("sobrescreveu o progresso existente — o `[x]` de quem trabalhou seria apagado")
	}

	b, _ := os.ReadFile(caminhoDeProgresso(p))
	if string(b) != jaFeito {
		t.Fatalf("o conteúdo mudou:\n%s", b)
	}
}

// Um plano SEM fases declaradas ainda ganha o arquivo, com o que fazer escrito.
func TestProgresso_planoSemFasesDizOQueFazer(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "0002-y.md")
	plano := "# Plano\n\n## Objetivo\n\nnada de fases ainda\n"
	if err := os.WriteFile(p, []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	destino, err := escreveProgressoInicial(p, plano, "YYYYY")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(destino)
	if !strings.Contains(string(b), "TODO") {
		t.Errorf("plano sem fases devia dizer o que fazer:\n%s", b)
	}
}

// O TAMANHO DO CÓDIGO vem da CONFIG, não de um literal.
//
// `code_lengths` é configurável por projeto. Enquanto o regex fixava `{4,5}`, um projeto
// com código de 3 letras tinha as fases reconhecidas pelos gates (que usam
// `config.CodeLengthPattern()`) e IGNORADAS por este comando — o progresso nascia vazio,
// sem nada acusar. O comentário do código afirmava paridade com os gates; o código só a
// tinha para quem estivesse no default.
//
// Achado no review do próprio PR que introduziu o arquivo.
func TestProgresso_respeitaCodeLengthsDoProjeto(t *testing.T) {
	original := config.CodeLengths
	t.Cleanup(func() { config.CodeLengths = original })
	config.CodeLengths = []int{3}

	plano := "# Plano\n\n### ABC-F01 — fase de código curto\n"
	dir := t.TempDir()
	p := filepath.Join(dir, "0001-curto.md")
	if err := os.WriteFile(p, []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}

	destino, err := escreveProgressoInicial(p, plano, "ABC")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(destino)
	if !strings.Contains(string(b), "## ABC-F01") {
		t.Fatalf("a fase de um projeto com code_lengths=[3] não foi reconhecida — os gates "+
			"a veem e este comando não:\n%s", b)
	}
}

// O `dispensado` é ACEITO, e é distinto de `pass` no mapa.
//
// A instrução do `@TBD` (ver `internal/initx/instrucao_tbd.go`) manda responder
// DISPENSADO quando o alvo da pergunta não existe. Se o comando recusasse esse veredito,
// a instrução mandaria o agente para uma parede — e a saída seria voltar ao `pass`
// mentiroso que ela existe para evitar.
//
// Achado ao implementar a instrução: o `judge` aceitava só `pass` e `fail`.
func TestJudge_aceitaDispensadoEExigeMotivo(t *testing.T) {
	// A validação é a do comando; o que se prova aqui é o CONTRATO dos três vereditos.
	for _, c := range []struct {
		verdict string
		reason  string
		querErr bool
		porque  string
	}{
		{"pass", "", false, "pass sem motivo é aceito (aprovação não precisa de laudo)"},
		{"fail", "", true, "fail sem motivo não é acionável"},
		{"dispensado", "", true, "dispensado sem a ausência nomeada é indistinguível de gate desligado"},
		{"dispensado", "a spec declara @TBD: code e MTHRN não existe", false, "dispensado com motivo é aceito"},
		{"inventado", "x", true, "veredito fora dos três é recusado"},
	} {
		err := validaVeredito(c.verdict, c.reason)
		if (err != nil) != c.querErr {
			t.Errorf("verdict=%q reason=%q: err=%v, queria erro=%v (%s)",
				c.verdict, c.reason, err, c.querErr, c.porque)
		}
	}
}
