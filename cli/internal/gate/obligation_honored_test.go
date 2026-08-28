package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func obligCfg() *config.Config {
	return &config.Config{Obligations: []config.Obligation{{
		Name:         "pii-purgavel",
		When:         "carries: pii",
		MustAppearIn: []string{"purge.ts"},
		IdentifiedBy: "screaming-snake",
		Because:      "LGPD — dado pessoal precisa ser apagável",
	}}}
}

// tmpRoot cria uma raiz com o arquivo-destino no conteúdo dado.
func tmpRoot(t *testing.T, purge string) string {
	t.Helper()
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "purge.ts"), []byte(purge), 0o644); err != nil {
		t.Fatal(err)
	}
	return d
}

func obligNode() mapx.Node { return mapx.Node{ID: "models/MetadataEntry.spec.md", Kind: mapx.KindSpec} }

func TestObligation_gatilhoSemDeverReprova(t *testing.T) {
	root := tmpRoot(t, "const tables = [USER_PROFILE_TABLE_NAME]\n")
	content := "<!-- @anchors\n  code: ABCDX\n  carries: pii\n-->\n"
	v, msg := checkObligationHonored(content, obligNode(), root, nil, obligCfg())
	if v != Fail {
		t.Fatalf("esperava Fail, got %v", v)
	}
	if !strings.Contains(msg, "LGPD") {
		t.Errorf("a mensagem deve dizer POR QUE (o `because`): %q", msg)
	}
}

func TestObligation_deverCumpridoPassa(t *testing.T) {
	// screaming-snake: MetadataEntry → METADATA_ENTRY
	root := tmpRoot(t, "const t = [METADATA_ENTRY_TABLE]\n")
	content := "<!-- @anchors\n  carries: pii\n-->\n"
	if v, msg := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Errorf("dever cumprido deveria passar: %v (%s)", v, msg)
	}
}

func TestObligation_semGatilhoPassa(t *testing.T) {
	root := tmpRoot(t, "vazio\n")
	content := "<!-- @anchors\n  code: ABCDX\n-->\n" // não declara carries: pii
	if v, _ := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Errorf("nó sem o atributo-gatilho não contrai a obrigação")
	}
}

func TestObligation_waiverComMotivoEximeMasSemMotivoNao(t *testing.T) {
	root := tmpRoot(t, "vazio\n")
	comMotivo := "<!-- @anchors\n  carries: pii\n  obligation_waived: pii-purgavel — grupo compartilhado; apagar destruiria dado de terceiros\n-->\n"
	if v, _ := checkObligationHonored(comMotivo, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Errorf("waiver COM motivo deveria eximir")
	}
	semMotivo := "<!-- @anchors\n  carries: pii\n  obligation_waived: pii-purgavel\n-->\n"
	if v, _ := checkObligationHonored(semMotivo, obligNode(), root, nil, obligCfg()); v != Fail {
		t.Errorf("waiver SEM motivo não vale — é o que separa exceção honesta de silêncio")
	}
}

func TestObligation_semObrigacaoDeclaradaPula(t *testing.T) {
	root := tmpRoot(t, "")
	if v, _ := checkObligationHonored("carries: pii", obligNode(), root, nil, &config.Config{}); v != Skip {
		t.Errorf("projeto sem obrigações declaradas deveria dar Skip")
	}
}

func TestApplyIdentifierForm(t *testing.T) {
	cases := map[string]string{
		"as-is":                    "MetadataEntry",
		"screaming-snake":          "METADATA_ENTRY",
		"snake":                    "metadata_entry",
		"kebab":                    "metadata-entry",
		"{{SCREAMING}}_TABLE_NAME": "METADATA_ENTRY_TABLE_NAME",
	}
	for form, want := range cases {
		if got := applyIdentifierForm("MetadataEntry", form); got != want {
			t.Errorf("applyIdentifierForm(%q) = %q, quer %q", form, got, want)
		}
	}
}

// O TERCEIRO ESTADO: a dívida assumida.
//
// O gate só oferecia "cumpra" ou "dispense", e nenhuma servia ao caso mais comum: a
// obrigação é REALX e será cumprida noutra fase. Um agente relatou a armadilha com
// precisão — dispensar seria mentira (o dever não deixou de existir), e deixar vermelho
// confunde dívida assumida com esquecimento, que é justamente a distinção que o pilar
// existe para preservar.
func TestObligation_dividaAssumida(t *testing.T) {
	root := tmpRoot(t, "// nada aqui\n")
	base := "<!-- @anchors\n  code: MTENX\n  carries: pii\n%s-->\n# x\n"

	// sem declaração: falha, e a mensagem oferece as TRÊS saídas
	v, d := checkObligationHonored(fmt.Sprintf(base, ""), obligNode(), root, nil, obligCfg())
	if v != Fail {
		t.Fatalf("obrigação descumprida sem declaração deveria falhar, foi %s (%s)", v, d)
	}
	for _, saida := range []string{"CUMPRA", "obligation_pending", "obligation_waived"} {
		if !strings.Contains(d, saida) {
			t.Errorf("a mensagem não oferece a saída %q: %s", saida, d)
		}
	}

	// dívida assumida COM o quando: Pendente (registro visível), nunca Pass nem Fail
	comQuando := fmt.Sprintf(base, "  obligation_pending: pii-purgavel — o handler nasce na fase 2 do plano\n")
	v, d = checkObligationHonored(comQuando, obligNode(), root, nil, obligCfg())
	if v != Pending {
		t.Fatalf("dívida assumida deveria ser Pendente, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "DÍVIDA ASSUMIDA") || !strings.Contains(d, "fase 2") {
		t.Errorf("o registro não mostra o compromisso: %s", d)
	}

	// marcador NU não assume dívida nenhuma — continua falhando
	nu := fmt.Sprintf(base, "  obligation_pending: pii-purgavel\n")
	if v, _ := checkObligationHonored(nu, obligNode(), root, nil, obligCfg()); v != Fail {
		t.Fatalf("dívida sem QUANDO deveria continuar falhando, foi %s", v)
	}
}

// Dispensa e dívida continuam distintas: dispensar diz "o dever não se aplica";
// assumir diz "o dever vale e será pago". Só a primeira passa como resolvida.
func TestObligation_dispensaEDividaSaoDistintas(t *testing.T) {
	root := tmpRoot(t, "// nada aqui\n")
	base := "<!-- @anchors\n  code: MTENX\n  carries: pii\n%s-->\n# x\n"

	dispensa := fmt.Sprintf(base, "  obligation_waived: pii-purgavel — este nó não guarda dado do titular\n")
	if v, d := checkObligationHonored(dispensa, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Fatalf("dispensa COM motivo passa (o dever não se aplica), foi %s (%s)", v, d)
	}
	divida := fmt.Sprintf(base, "  obligation_pending: pii-purgavel — fase 2\n")
	if v, _ := checkObligationHonored(divida, obligNode(), root, nil, obligCfg()); v == Pass {
		t.Fatal("dívida assumida NÃO pode passar como cumprida — ela ainda é devida")
	}
}
