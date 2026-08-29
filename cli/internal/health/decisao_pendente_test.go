package health

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O doctor listava sinal de mutação ausente como ponta de atenção enquanto declarava
// "0 pontas" com uma regra pendente de decisão. A hierarquia estava invertida: o que falta
// MEDIR aparecia, o que falta DECIDIR não — e é a decisão que mais precisa sobreviver à
// sessão, porque depende de humano e pode levar semanas.
func TestDoctorReportaDecisaoPendente(t *testing.T) {
	dir := t.TempDir()
	spec := "# Unidade\n\n## Decisões em aberto\n\n| Código | Pergunta | Quem decide |\n| --- | --- | --- |\n" +
		"| `PARCX-Q01` | Fuso do vencimento? | Produto |\n| `PARCX-Q02` | Retenção dos logs? | Segurança |\n"
	if err := os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "u.spec.md", Kind: mapx.KindSpec}}}

	fs := checkDecisoesPendentes(g, dir, nil)

	if len(fs) != 1 {
		t.Fatalf("esperava 1 achado (a spec com pendências), veio %d", len(fs))
	}
	if fs[0].Check != "decisao-pendente" || fs[0].Severity != Warn {
		t.Errorf("achado errado: %+v", fs[0])
	}
	if !strings.Contains(fs[0].Detail, "2 decisão") {
		t.Errorf("deveria contar as duas perguntas: %s", fs[0].Detail)
	}

	// Spec que FECHOU a seção com "nenhuma" não é pendência — é o opt-out honesto.
	fechada := "# Unidade\n\n## Decisões em aberto\n\nnenhuma\n"
	if err := os.WriteFile(filepath.Join(dir, "u.spec.md"), []byte(fechada), 0o644); err != nil {
		t.Fatal(err)
	}
	if fs := checkDecisoesPendentes(g, dir, nil); len(fs) != 0 {
		t.Errorf("seção fechada com `nenhuma` não é pendência: %+v", fs)
	}
}
