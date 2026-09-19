package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ghFalso põe no PATH um `gh` que responde o JSON dado — sem rede e sem repositório.
func ghFalso(t *testing.T, saida string) {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\ncat <<'FIM'\n" + saida + "\nFIM\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// ESCALAR UM ALVO QUE JÁ TEM CARD precisa avisar.
//
// Dois agentes entregaram o MESMO trabalho no mesmo dia no projeto de referência: um pegou
// o card do gate para `MetricCard.spec.md` às 11:38; outro, trabalhando noutro card,
// encontrou o mesmo problema às 12:02 e abriu um card NOVO. Os dois PRs acrescentaram a
// mesma seção ao mesmo documento.
//
// O `claim` impede dois agentes de pegarem o mesmo card — não impedia um agente de CRIAR um
// card para trabalho já em andamento noutro.
func TestAvisaQuandoOAlvoJaTemCardAberto(t *testing.T) {
	ghFalso(t, `[{"number":405,"title":"[doc-required] Violação @ apps/mobile/src/components/MetricCard.spec.md","body":"corpo"}]`)
	got := openCardsAbout("apps/mobile/src/components/MetricCard.spec.md", "anchors")
	if len(got) != 1 {
		t.Fatalf("o card aberto do alvo não foi achado: %v", got)
	}
	// O CONTEÚDO do aviso, e não só a contagem. Quem o lê precisa saber QUAL card ir
	// ver — um aviso que diz "há um card" sem dizer qual manda procurar na fila inteira,
	// e é aí que a pessoa desiste e cria o card novo mesmo assim.
	if !strings.Contains(got[0], "#405") {
		t.Errorf("o aviso não traz o NÚMERO do card: %q", got[0])
	}
	if !strings.Contains(got[0], "MetricCard.spec.md") {
		t.Errorf("o aviso não traz o TÍTULO, que é o que diz se é o mesmo trabalho: %q", got[0])
	}
}

// A BUSCA DO GITHUB É APROXIMADA, e sem confirmação ela casa o alvo errado.
//
// É a mesma defesa que o `internal/issue` tem com o marcador: `MetricCard.spec.md` casaria
// `MetricCardList.spec.md` numa busca por texto, e o aviso apontaria trabalho que não é o
// mesmo — ensinando a ignorar o aviso.
func TestNaoConfundeAlvoComOutroQueOContem(t *testing.T) {
	ghFalso(t, `[{"number":999,"title":"[doc-required] Violação @ apps/mobile/src/components/MetricCardList.spec.md","body":"outro alvo"}]`)
	if v := openCardsAbout("apps/mobile/src/components/MetricCard.spec.md", "anchors"); len(v) != 0 {
		t.Errorf("casou um alvo DIFERENTE que apenas contém o nome: %v", v)
	}
}

// SEM ALVO não há o que perguntar — e chamar o `gh` à toa atrasaria todo `escalate` sem
// `--about`.
func TestSemAlvoNaoConsulta(t *testing.T) {
	ghFalso(t, `[{"number":1,"title":"qualquer","body":"x"}]`)
	if v := openCardsAbout("", "anchors"); len(v) != 0 {
		t.Errorf("consultou sem alvo: %v", v)
	}
	if v := openCardsAbout("X.spec.md", ""); len(v) != 0 {
		t.Errorf("consultou sem label: %v", v)
	}
}

// O `gh` QUE FALHA não pode derrubar o `escalate`: a conferência é auxiliar, e impedir o
// registro por causa dela seria pior que a duplicata que ela evita.
func TestFalhaDaConsultaNaoDerrubaOEscalate(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if v := openCardsAbout("X.spec.md", "anchors"); v != nil {
		t.Errorf("devolveu %v quando o `gh` falhou — deveria seguir sem aviso", v)
	}
}
