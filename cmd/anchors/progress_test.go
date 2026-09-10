package main

import (
	"io"
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
	if !scan.IsProgressFile("plans/0001-x" + progressSuffix) {
		t.Fatalf("o sufixo do comando (%q) não é reconhecido pelo scan — o `new` criaria "+
			"um arquivo que o mapa NÃO exclui, e os gates voltariam a confrontá-lo",
			progressSuffix)
	}
}

func TestProgresso_caminhoFicaAoLadoDoPlano(t *testing.T) {
	got := progressPath("plans/0017-mutacao.md")
	want := "plans/0017-mutacao" + progressSuffix
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

	destino, err := writeInitialProgress(p, plano, "MTUAO")
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
	if err := os.WriteFile(progressPath(p), []byte(jaFeito), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := writeInitialProgress(p, "# Plano\n\n### ABCDE-F01 — fase\n", "ABCDE"); err == nil {
		t.Fatal("sobrescreveu o progresso existente — o `[x]` de quem trabalhou seria apagado")
	}

	b, _ := os.ReadFile(progressPath(p))
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
	destino, err := writeInitialProgress(p, plano, "YYYYY")
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

	destino, err := writeInitialProgress(p, plano, "ABC")
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
		{"waived", "", true, "dispensado sem a ausência nomeada é indistinguível de gate desligado"},
		{"waived", "a spec declara @TBD: code e MTHRN não existe", false, "dispensado com motivo é aceito"},
		{"inventado", "x", true, "veredito fora dos três é recusado"},
	} {
		err := validateVerdict(c.verdict, c.reason)
		if (err != nil) != c.querErr {
			t.Errorf("verdict=%q reason=%q: err=%v, queria erro=%v (%s)",
				c.verdict, c.reason, err, c.querErr, c.porque)
		}
	}
}

// O `anchors new progress --for <plano>` existe para os planos que nasceram ANTES do
// mecanismo. O `anchors new plan` cria o companheiro junto — mas só ele, e um projeto que
// adotou o Anchors antes desta versão fica com todos os planos sem companheiro, para
// sempre, sem nada acusar.
//
// Medido no blue-eyes: 17 planos, ZERO com `-progress.md`, e 17 com checkbox DENTRO do
// plano — que é exatamente o que este arquivo existe para tirar de lá. O caminho feliz
// (plano novo) funcionava; era a ADOÇÃO que não tinha caminho.
func TestNewProgress_criaParaPlanoExistente(t *testing.T) {
	root := t.TempDir()
	plano := filepath.Join(root, "plans", "0002-plataforma.md")
	if err := os.MkdirAll(filepath.Dir(plano), 0o755); err != nil {
		t.Fatal(err)
	}
	const conteudo = `<!-- @anchors
  code: PLTFR
  layer: plan
-->
# Plataforma

### PLTFR-F01 — o contrato

### PLTFR-F02 — o acesso às fontes
`
	if err := os.WriteFile(plano, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := newProgressCmd()
	cmd.SetArgs([]string{"--root", root, "--for", "plans/0002-plataforma.md"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("criar o progresso: %v", err)
	}

	prog := filepath.Join(root, "plans", "0002-plataforma-progress.md")
	b, err := os.ReadFile(prog)
	if err != nil {
		t.Fatalf("o arquivo não foi criado: %v", err)
	}
	got := string(b)
	// a IDENTIDADE vem do plano, não de um argumento: um progresso com código
	// divergente do plano que ele acompanha deixaria de ser localizável.
	if !strings.Contains(got, "# Progresso — PLTFR") {
		t.Errorf("o progresso não herdou o código do plano:\n%s", got)
	}
	for _, fase := range []string{"## PLTFR-F01", "## PLTFR-F02"} {
		if !strings.Contains(got, fase) {
			t.Errorf("falta a seção %q — as fases vêm dos cabeçalhos do plano", fase)
		}
	}

	// NÃO SOBRESCREVE: o arquivo guarda estado, e regravá-lo apagaria o que já foi
	// registrado. É o que torna seguro rodar o comando sobre um projeto inteiro.
	cmd2 := newProgressCmd()
	cmd2.SetArgs([]string{"--root", root, "--for", "plans/0002-plataforma.md"})
	cmd2.SetOut(io.Discard)
	cmd2.SetErr(io.Discard)
	if err := cmd2.Execute(); err == nil {
		t.Error("rodar de novo sobrescreveu o estado — devia recusar")
	}
}

// Sem `code:` no header o progresso nasceria sem identidade, e o par plano/progresso
// deixaria de ser localizável por código.
func TestNewProgress_recusaPlanoSemCodigo(t *testing.T) {
	root := t.TempDir()
	plano := filepath.Join(root, "p.md")
	if err := os.WriteFile(plano, []byte("# Plano sem header\n\n### ABCDE-F01 — fase\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := newProgressCmd()
	cmd.SetArgs([]string{"--root", root, "--for", "p.md"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	if err := cmd.Execute(); err == nil {
		t.Error("aceitou plano sem `code:` — o progresso nasceria sem identidade")
	}
}
