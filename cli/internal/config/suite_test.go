package config

import (
	"strings"
	"testing"
)

// O cenário é o monorepo real que motivou os dois eixos: `unit` existe em DOIS
// workspaces, e `integration` só num deles.
func suites() []Suite {
	return []Suite{
		{Workspace: "backend", Layer: "unit", Run: "b-unit"},
		{Workspace: "backend", Layer: "integration", Run: "b-it"},
		{Workspace: "mobile", Layer: "unit", Run: "m-unit"},
		{Workspace: "mobile", Layer: "e2e", Run: "m-e2e"},
	}
}

func comandos(sel []Suite) string {
	out := make([]string, 0, len(sel))
	for _, s := range sel {
		out = append(out, s.Run)
	}
	return strings.Join(out, ",")
}

// TestSemFiltroRodaTudo — o default de um comando declarativo é "o que foi declarado".
// Não pode ser "nada": um comando que só age com flag ensina a sempre passar flag, e a
// declaração no arquivo perde a função.
func TestSemFiltroRodaTudo(t *testing.T) {
	sel, ausentes := SelecionaSuites(suites(), nil, nil, nil)
	if len(sel) != 4 || len(ausentes) != 0 {
		t.Fatalf("sem filtro devia rodar as 4; veio %d (ausentes %v)", len(sel), ausentes)
	}
}

// TestCamadaAtravessaWorkspaces é a razão de os eixos serem separados: `anchors test
// unit` significa "a unit de TODO MUNDO". Com um campo só, isto era impossível de
// expressar — era preciso inventar uma camada por workspace.
func TestCamadaAtravessaWorkspaces(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, nil, nil)
	if got := comandos(sel); got != "b-unit,m-unit" {
		t.Errorf("unit devia pegar os dois workspaces; veio %q", got)
	}
}

// TestWorkspaceAtravessaCamadas — o eixo oposto: "tudo do backend", sem nomear camada.
func TestWorkspaceAtravessaCamadas(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), nil, []string{"backend"}, nil)
	if got := comandos(sel); got != "b-unit,b-it" {
		t.Errorf("backend devia pegar as duas camadas dele; veio %q", got)
	}
}

// TestOsDoisEixosSeCombinam — a interseção é o caso do dia a dia ("só a unit do
// backend, que é rápida").
func TestOsDoisEixosSeCombinam(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, []string{"backend"}, nil)
	if got := comandos(sel); got != "b-unit" {
		t.Errorf("interseção errada; veio %q", got)
	}
}

func TestVariosValoresEmCadaEixo(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), []string{"unit", "e2e"}, []string{"mobile"}, nil)
	if got := comandos(sel); got != "m-unit,m-e2e" {
		t.Errorf("queria as duas camadas do mobile; veio %q", got)
	}
}

// TestOrdemEhADoArquivoNaoADoUsuario — quem escreveu o anchors.yaml decidiu que unit
// vem antes de e2e, e isso costuma ser dependência real. Respeitar a digitação faria o
// mesmo pedido se comportar diferente conforme quem digita.
func TestOrdemEhADoArquivoNaoADoUsuario(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), []string{"e2e", "unit"}, nil, nil)
	if got := comandos(sel); got != "b-unit,m-unit,m-e2e" {
		t.Errorf("a ordem é a da declaração; veio %q", got)
	}
}

// TestNomeInexistenteVemRotuladoEJunto — errar o nome é o engano comum, e o usuário
// precisa saber em QUAL eixo errou: "smoke" e "web" falham por motivos diferentes.
func TestNomeInexistenteVemRotuladoEJunto(t *testing.T) {
	_, ausentes := SelecionaSuites(suites(), []string{"unit", "smoke"}, []string{"web"}, nil)
	if len(ausentes) != 2 {
		t.Fatalf("as duas ausências vêm juntas; veio %v", ausentes)
	}
	if !strings.Contains(ausentes[0], "camada") || !strings.Contains(ausentes[0], "smoke") {
		t.Errorf("a primeira devia ser a camada smoke; veio %q", ausentes[0])
	}
	if !strings.Contains(ausentes[1], "workspace") || !strings.Contains(ausentes[1], "web") {
		t.Errorf("a segunda devia ser o workspace web; veio %q", ausentes[1])
	}
}

// TestCombinacaoVaziaNaoEhNomeErrado guarda a distinção que mais confunde: `integration`
// e `mobile` existem, mas não JUNTOS. Reportar "camada integration não existe" mandaria
// o usuário corrigir um nome que está certo — o que falta é a combinação.
func TestCombinacaoVaziaNaoEhNomeErrado(t *testing.T) {
	sel, ausentes := SelecionaSuites(suites(), []string{"integration"}, []string{"mobile"}, nil)
	if len(ausentes) != 0 {
		t.Errorf("os dois nomes existem; nada a reportar como ausente. veio %v", ausentes)
	}
	if len(sel) != 0 {
		t.Errorf("a combinação não existe, então nada roda; veio %q", comandos(sel))
	}
}

func TestIgnoraCaixaEEspaco(t *testing.T) {
	sel, ausentes := SelecionaSuites(suites(), []string{" Unit "}, []string{"BackEnd"}, nil)
	if len(sel) != 1 || len(ausentes) != 0 {
		t.Errorf("caixa e espaço não deviam separar; veio %q (ausentes %v)", comandos(sel), ausentes)
	}
}

// TestWorkspaceEhOpcional — projeto de pacote único não tem esse eixo. Exigir o campo
// obrigaria metade dos projetos a inventar um nome para "o projeto".
func TestWorkspaceEhOpcional(t *testing.T) {
	simples := []Suite{{Layer: "unit", Run: "u"}, {Layer: "e2e", Run: "e"}}
	sel, ausentes := SelecionaSuites(simples, []string{"unit"}, nil, nil)
	if len(sel) != 1 || sel[0].Run != "u" || len(ausentes) != 0 {
		t.Errorf("suíte sem workspace devia funcionar; veio %q / %v", comandos(sel), ausentes)
	}
	if ws := WorkspacesDeclarados(simples); len(ws) != 0 {
		t.Errorf("sem workspace declarado, a lista é vazia; veio %v", ws)
	}
}

// TestVocabularioEhDoProjeto guarda a agnosticidade: não há lista fixa de camadas nem
// de workspaces. Validar contra unit/integration/e2e seria o framework decidindo o
// vocabulário de quem o usa.
func TestVocabularioEhDoProjeto(t *testing.T) {
	proprias := []Suite{
		{Workspace: "cobranca", Layer: "contrato", Run: "c"},
		{Workspace: "cobranca", Layer: "carga", Run: "k"},
	}
	sel, ausentes := SelecionaSuites(proprias, []string{"carga"}, []string{"cobranca"}, nil)
	if len(sel) != 1 || sel[0].Run != "k" || len(ausentes) != 0 {
		t.Errorf("vocabulário próprio devia ser aceito; veio %q / %v", comandos(sel), ausentes)
	}
}

// TestListasNaoRepetemNomes — as listas são o que se mostra a quem errou; `unit`
// aparecendo duas vezes (uma por workspace) faria a mensagem parecer um dump.
func TestListasNaoRepetemNomes(t *testing.T) {
	if got := strings.Join(CamadasDeclaradas(suites()), ","); got != "unit,integration,e2e" {
		t.Errorf("camadas repetidas ou fora de ordem: %q", got)
	}
	if got := strings.Join(WorkspacesDeclarados(suites()), ","); got != "backend,mobile" {
		t.Errorf("workspaces repetidos ou fora de ordem: %q", got)
	}
}

// ── terceiro eixo: escopo (só na mutação) ──────────────────────────────────────

// O escopo é a MESMA unidade medida com abrangências diferentes: contra o teste dela
// (`isolated`) e contra o de quem a importa (`full`). Não é camada nem workspace, e o
// gate julga pelo isolado quando os dois existem — por isso precisa ser endereçável.
func suitesMutacao() []Suite {
	return []Suite{
		{Workspace: "backend", Layer: "unit", Scope: "full", Run: "b-full"},
		{Workspace: "backend", Layer: "unit", Scope: "isolated", Run: "b-iso"},
		{Workspace: "mobile", Layer: "unit", Scope: "full", Run: "m-full"},
	}
}

// TestSemEscopoRodaOsDois é o default que importa: é o PAR que produz a leitura de
// acoplamento do gate. Exigir dois comandos para obtê-lo devolveria o passo manual que
// estes comandos vieram eliminar.
func TestSemEscopoRodaOsDois(t *testing.T) {
	sel, _ := SelecionaSuites(suitesMutacao(), nil, []string{"backend"}, nil)
	if got := comandos(sel); got != "b-full,b-iso" {
		t.Errorf("sem --scope, roda os dois escopos do backend; veio %q", got)
	}
}

func TestEscopoFiltraUmSo(t *testing.T) {
	sel, _ := SelecionaSuites(suitesMutacao(), nil, nil, []string{"isolated"})
	if got := comandos(sel); got != "b-iso" {
		t.Errorf("só o isolated; veio %q", got)
	}
}

// TestEscopoCombinaComOsOutrosEixos — os três se cruzam, como camada e workspace já
// faziam.
func TestEscopoCombinaComOsOutrosEixos(t *testing.T) {
	sel, _ := SelecionaSuites(suitesMutacao(), []string{"unit"}, []string{"backend"}, []string{"full"})
	if got := comandos(sel); got != "b-full" {
		t.Errorf("interseção dos três eixos; veio %q", got)
	}
}

// TestMesmaCamadaEWorkspaceComEscoposDiferentesNaoColidem guarda o que motivou o eixo:
// as duas suítes do backend têm (workspace, layer) IDÊNTICOS e só diferem no escopo.
// Sem o terceiro eixo elas seriam indistinguíveis, e o isolado ficaria inalcançável —
// que é exatamente o buraco que este eixo fecha.
func TestMesmaCamadaEWorkspaceComEscoposDiferentesNaoColidem(t *testing.T) {
	sel, ausentes := SelecionaSuites(suitesMutacao(), []string{"unit"}, []string{"backend"}, []string{"isolated"})
	if len(ausentes) != 0 {
		t.Fatalf("nada ausente aqui; veio %v", ausentes)
	}
	if len(sel) != 1 || sel[0].Run != "b-iso" {
		t.Errorf("o escopo tem de desempatar as duas; veio %q", comandos(sel))
	}
}

// TestEscopoInexistenteVemRotulado — o usuário precisa saber em QUAL eixo errou.
func TestEscopoInexistenteVemRotulado(t *testing.T) {
	_, ausentes := SelecionaSuites(suitesMutacao(), nil, nil, []string{"parcial"})
	if len(ausentes) != 1 || !strings.Contains(ausentes[0], "escopo") {
		t.Errorf("ausência devia ser rotulada como escopo; veio %v", ausentes)
	}
}

// TestSuiteSemEscopoNaoSomeQuandoNinguemFiltra — `tests:` não declara escopo, e o eixo
// não pode excluir quem não participa dele.
func TestSuiteSemEscopoNaoSomeQuandoNinguemFiltra(t *testing.T) {
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, nil, nil)
	if len(sel) != 2 {
		t.Errorf("suítes sem escopo continuam valendo; veio %q", comandos(sel))
	}
}
