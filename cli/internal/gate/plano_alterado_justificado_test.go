package gate

import (
	"github.com/co2-lab/anchors/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O CASO QUE MOTIVOU O GATE: alguém achou uma inconsistência no plano e corrigiu em
// silêncio. O arquivo fica VÁLIDO — a inconsistência foi removida —, e é por isso que
// nenhum gate de estado o pega. O que denuncia é a mudança sem justificativa.
func TestAlterado_semRevisaoReprova(t *testing.T) {
	v, d := checkPlanoAlteradoJustificado(
		"# Plano 0001 — Fundação\n\n## Objetivo\n\nTexto corrigido sem dizer nada.\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Fail {
		t.Fatalf("plano alterado sem revisão deve reprovar, veio %v: %s", v, d)
	}
	// O laudo tem de ENSINAR as duas saídas, senão o agente escolhe a errada por não
	// saber que a outra existe.
	if !strings.Contains(d, "FNDTN-R0001") {
		t.Errorf("o laudo deve mostrar o formato com o código do arquivo; veio: %s", d)
	}
	if !strings.Contains(d, "anchors escalate") {
		t.Errorf("o laudo deve dar a SAÍDA para quando a correção MUDA A DIREÇÃO — nomear a "+
			"label não basta, o comando é o que torna escalar tão barato quanto corrigir; "+
			"veio: %s", d)
	}
}

func TestAlterado_comRevisaoPassa(t *testing.T) {
	v, d := checkPlanoAlteradoJustificado(
		"# Plano 0001\n\n> **FNDTN-R0001:** o exemplo citava um pacote que não existe.\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Pass {
		t.Fatalf("revisão declarada deveria passar, veio %v: %s", v, d)
	}
}

// A revisão de OUTRO documento não justifica a mudança deste. Sem isto, bastaria citar a
// revisão alheia ao explicar o contexto para o gate calar.
func TestAlterado_revisaoDeOutroNaoConta(t *testing.T) {
	v, _ := checkPlanoAlteradoJustificado(
		"# Plano 0017\n\nComo explica a `MTUAO-R0001`, o CI mudou.\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Fail {
		t.Fatalf("revisão de outro código não justifica a mudança deste, veio %v", v)
	}
}

// A NUMERAÇÃO é o que uma marca solta não daria: quantas vezes o documento mudou. Com
// buraco ela para de responder isso.
func TestAlterado_numeracaoComBuracoReprova(t *testing.T) {
	v, d := checkPlanoAlteradoJustificado(
		"> **FNDTN-R0001:** primeira.\n\n> **FNDTN-R0007:** pulou do 1 para o 7.\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Fail {
		t.Fatalf("numeração com buraco deve reprovar, veio %v: %s", v, d)
	}
}

func TestAlterado_variasRevisoesSequenciaisPassam(t *testing.T) {
	v, d := checkPlanoAlteradoJustificado(
		"> **FNDTN-R0001:** primeira.\n\n> **FNDTN-R0002:** segunda.\n\n> **FNDTN-R0003:** terceira.\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Pass {
		t.Fatalf("três revisões sequenciais deveriam passar, veio %v: %s", v, d)
	}
	// O laudo cita a ÚLTIMA — é a que explica a mudança que está sendo confrontada.
	if !strings.Contains(d, "R0003") {
		t.Errorf("o laudo deve citar a revisão mais recente; veio: %s", d)
	}
}

// Sem código não há como saber de quem é a revisão. Skip, e não Fail: quem cobra o código
// é outro gate, e reprovar aqui daria dois achados para um defeito só.
func TestAlterado_semCodigoPula(t *testing.T) {
	v, _ := checkPlanoAlteradoJustificado("# Plano\n", mapx.Node{ID: "plans/0001.md"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
	if v != Skip {
		t.Fatalf("sem código o gate deve pular, veio %v", v)
	}
}

// O marcador tem de sobreviver às formas que as pessoas escrevem markdown: dentro de
// alerta, em negrito, como item de lista.
func TestAlterado_formasDeEscreverAMarca(t *testing.T) {
	for _, forma := range []string{
		"> **FNDTN-R0001:** dentro de citação e negrito",
		"FNDTN-R0001: cru, no começo da linha",
		"> [!NOTE]\n> FNDTN-R0001: dentro de um alerta do GitHub",
		"**FNDTN-R0001**: negrito antes dos dois-pontos",
	} {
		if v, d := checkPlanoAlteradoJustificado(forma, mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md")); v != Pass {
			t.Errorf("deveria aceitar %q, veio %v: %s", forma, v, d)
		}
	}
}

// cfgAlterado monta o config com a lista do que MUDOU nesta execução.
func cfgAlterado(paths ...string) *config.Config {
	return &config.Config{Alterados: paths}
}

// O CASO QUE O GATE ACUSAVA ERRADO, e que os outros testes não pegavam por passarem só o
// nó alterado: `--changed X` entrega o RAIO DE IMPACTO de X, e as specs que X semeia
// entram no escopo sem terem mudado.
//
// Medido no blue-eyes: alterar UM plano acusava 8 arquivos, 7 intocados. Um gate
// bloqueante que acusa inocente é pior que gate nenhum — a saída barata vira desligá-lo.
func TestAlterado_noRaioDeImpactoMasIntocadoPula(t *testing.T) {
	v, d := checkPlanoAlteradoJustificado(
		"# Spec que ninguém tocou\n",
		mapx.Node{ID: "packages/shared/Workspace.spec.md", Code: "WKSPC"},
		semGit(t), nil, cfgAlterado("plans/0001-fundacao.md"))
	if v != Skip {
		t.Fatalf("nó alcançado só pelo raio de impacto não tem o que justificar, veio %v: %s", v, d)
	}
}

// Sem a lista (chamada que não passou config) o gate se cala: acusar sem saber o que mudou
// é exatamente o defeito que ele acabou de corrigir.
func TestAlterado_semListaDeAlteradosPula(t *testing.T) {
	v, _ := checkPlanoAlteradoJustificado("# Plano\n",
		mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, nil)
	if v != Skip {
		t.Fatalf("sem saber o que mudou, o gate deve se calar, veio %v", v)
	}
}

// O DEFEITO QUE O EXERCÍCIO REVELOU: os planos do PR de mutação já diziam por que mudaram
// — um com `revises:`, o outro com `@revised-by` —, e o gate reprovava os dois por não
// usarem a notação `-R0001`.
//
// Exigir a mesma informação em duas notações ensina a satisfazer o gate em vez de
// comunicar, que é o oposto do que ele existe para fazer.
func TestAlterado_jaExplicadoPorRevisaoDePlanoPassa(t *testing.T) {
	for _, conteudo := range []string{
		"<!-- @anchors\n  revises: plans/0001-fundacao.md\n-->\n# Plano 0017\n",
		"# Plano 0001\n\n> [!IMPORTANT]\n> `@revised-by: plans/0017-mutacao.md`\n",
		"# Plano 0001\n\n> [!WARNING]\n> `@amended-by: plans/0017-mutacao.md`\n",
	} {
		v, d := checkPlanoAlteradoJustificado(conteudo,
			mapx.Node{ID: "plans/0001.md", Code: "FNDTN"}, semGit(t), nil, cfgAlterado("plans/0001.md"))
		if v != Pass {
			t.Errorf("já explicado pelo mecanismo de revisão deveria passar, veio %v: %s\n%s",
				v, d, conteudo)
		}
	}
}

// `--changed X` é uma AFIRMAÇÃO de quem chama, não uma medição — e o gate não pode
// acreditar nela sozinho.
//
// Medido: rodei `check --changed Workspace.spec.md` para testar outra coisa, com o
// arquivo INTOCADO, e o gate cobrou justificativa. O mesmo vale no pre-commit, que passa
// todo arquivo staged — inclusive os que entraram por rebase ou merge sem ninguém os ter
// editado. Cobrar de quem não mexeu é o pior defeito num gate bloqueante: barra trabalho
// correto, e a saída barata vira desligá-lo.
func TestAlterado_gitDesmenteAListaRecebida(t *testing.T) {
	// Diretório sem git: o gate NÃO pode se calar por não conseguir medir — ali ele
	// devolve o benefício da dúvida a quem chamou.
	if !gitDizQueMudou(t.TempDir(), "qualquer.md") {
		t.Error("sem git não há como desmentir a lista: o gate deve confiar em quem chamou, " +
			"senão se silencia justamente onde não consegue medir")
	}
}

// semGit devolve um diretório FORA de repositório git.
//
// Sem isto os testes rodavam com `root: ""` — o diretório do próprio Anchors, um repo de
// verdade onde `plans/0001.md` não existe e o git diz, com razão, que nada mudou. O teste
// media o repositório errado.
func semGit(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// ARQUIVO NOVO não tem o que justificar — ele não existia, então nada foi ALTERADO.
//
// Medido ao escrever as duas primeiras specs de uma fase: o gate cobrou `-R0001` de
// arquivos que nasciam naquele commit. Num arquivo novo o texto INTEIRO é a decisão;
// exigir uma revisão de si mesmo no primeiro commit é ruído puro, e ruído em gate
// bloqueante é o que faz alguém desligá-lo.
func TestAlterado_arquivoNovoNaoTemOQueJustificar(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Skipf("sem git: %v — %s", err, out)
	}
	novo := filepath.Join(dir, "nova.spec.md")
	if err := os.WriteFile(novo, []byte("# Spec nova\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, d := checkPlanoAlteradoJustificado("# Spec nova\n",
		mapx.Node{ID: "nova.spec.md", Code: "NOVAA"}, dir, nil, cfgAlterado("nova.spec.md"))
	if v != Skip {
		t.Fatalf("arquivo que nunca foi commitado não tem alteração a justificar, veio %v: %s", v, d)
	}
}
