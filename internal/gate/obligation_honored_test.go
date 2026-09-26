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
	t.Run("OBHNB-B01: A node that carries the trigger and is absent from the demanded file fails", func(t *testing.T) {})
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
	t.Run("OBHNB-B02: A node that carries the trigger and does appear passes", func(t *testing.T) {})
	// screaming-snake: MetadataEntry → METADATA_ENTRY
	root := tmpRoot(t, "const t = [METADATA_ENTRY_TABLE]\n")
	content := "<!-- @anchors\n  carries: pii\n-->\n"
	if v, msg := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Errorf("dever cumprido deveria passar: %v (%s)", v, msg)
	}
}

func TestObligation_semGatilhoPassa(t *testing.T) {
	t.Run("OBHNB-B03: A node without the trigger contracts no obligation", func(t *testing.T) {})
	root := tmpRoot(t, "vazio\n")
	content := "<!-- @anchors\n  code: ABCDX\n-->\n" // não declara carries: pii
	if v, _ := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Errorf("nó sem o atributo-gatilho não contrai a obrigação")
	}
}

func TestObligation_waiverComMotivoEximeMasSemMotivoNao(t *testing.T) {
	t.Run("OBHNB-B04: A waiver exempts only when it carries a written reason", func(t *testing.T) {})
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
	t.Run("OBHNB-B05: A project with no declared obligation is skipped", func(t *testing.T) {})
	root := tmpRoot(t, "")
	if v, _ := checkObligationHonored("carries: pii", obligNode(), root, nil, &config.Config{}); v != Skip {
		t.Errorf("projeto sem obrigações declaradas deveria dar Skip")
	}
}

func TestApplyIdentifierForm(t *testing.T) {
	t.Run("OBHNB-I01: The token is derived through the declared form", func(t *testing.T) {})
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
func TestObligation_aMensagemOfereceAsTresSaidas(t *testing.T) {
	t.Run("OBHNB-B09: The failing verdict offers the three ways out", func(t *testing.T) {})
	root := tmpRoot(t, "// nada aqui\n")
	content := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n# x\n"
	v, d := checkObligationHonored(content, obligNode(), root, nil, obligCfg())
	if v != Fail {
		t.Fatalf("obrigação descumprida sem declaração deveria falhar, foi %s (%s)", v, d)
	}
	for _, saida := range []string{"obligation_pending", "obligation_waived"} {
		if !strings.Contains(d, saida) {
			t.Errorf("a mensagem não oferece a saída %q: %s", saida, d)
		}
	}
	if !strings.Contains(d, "CUMPRA") && !strings.Contains(d, "FULFILL") {
		t.Errorf("a mensagem não oferece a saída CUMPRA/FULFILL: %s", d)
	}
}

// O TERCEIRO ESTADO: a dívida assumida.
//
// O gate só oferecia "cumpra" ou "dispense", e nenhuma servia ao caso mais comum: a
// obrigação é REAL e será cumprida noutra fase. Dispensar seria mentira (o dever não
// deixou de existir), e deixar vermelho confunde dívida assumida com esquecimento, que é
// justamente a distinção que o pilar existe para preservar.
func TestObligation_dividaAssumidaComQuandoEhPendente(t *testing.T) {
	t.Run("OBHNB-B06: An acknowledged debt with a written when yields Pending", func(t *testing.T) {})
	root := tmpRoot(t, "// nada aqui\n")
	comQuando := "<!-- @anchors\n  code: MTENX\n  carries: pii\n" +
		"  obligation_pending: pii-purgavel — o handler nasce na fase 2 do plano\n-->\n# x\n"
	v, d := checkObligationHonored(comQuando, obligNode(), root, nil, obligCfg())
	if v != Pending {
		t.Fatalf("dívida assumida deveria ser Pendente, foi %s (%s)", v, d)
	}
	if (!strings.Contains(d, "DÍVIDA ASSUMIDA") && !strings.Contains(d, "ASSUMED DEBT")) || !strings.Contains(d, "fase 2") {
		t.Errorf("o registro não mostra o compromisso: %s", d)
	}
}

// O marcador NU não assume dívida nenhuma — só esconde melhor.
func TestObligation_dividaSemQuandoContinuaFalhando(t *testing.T) {
	t.Run("OBHNB-B07: A bare debt marker keeps failing", func(t *testing.T) {})
	root := tmpRoot(t, "// nada aqui\n")
	nu := "<!-- @anchors\n  code: MTENX\n  carries: pii\n" +
		"  obligation_pending: pii-purgavel\n-->\n# x\n"
	if v, d := checkObligationHonored(nu, obligNode(), root, nil, obligCfg()); v != Fail {
		t.Fatalf("dívida sem QUANDO deveria continuar falhando, foi %s (%s)", v, d)
	}
}

// Dispensa e dívida continuam distintas: dispensar diz "o dever não se aplica";
// assumir diz "o dever vale e será pago". Só a primeira passa como resolvida.
func TestObligation_dispensaEDividaSaoDistintas(t *testing.T) {
	t.Run("OBHNB-B08: Waiver and debt stay distinct", func(t *testing.T) {})
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

// GLOB QUE NÃO CASA ARQUIVO NENHUM não produz violação. Acusar onde não havia o que ler
// seria carimbar o que nunca foi medido — e o relatório passaria a dizer "descumprido"
// sobre um destino que não existe.
func TestObligation_globSemArquivoNaoInventaViolacao(t *testing.T) {
	t.Run("OBHNB-I02: A glob that matches no file produces no violation", func(t *testing.T) {})
	root := t.TempDir() // nenhum arquivo: o glob `purge.ts` não casa nada
	cfg := &config.Config{Obligations: []config.Obligation{{
		Name:         "pii-purgavel",
		When:         "carries: pii",
		MustAppearIn: []string{"purge.ts"},
		IdentifiedBy: "screaming-snake",
	}}}
	content := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n"
	if v, d := checkObligationHonored(content, obligNode(), root, nil, cfg); v != Pass {
		t.Fatalf("sem arquivo para checar não há violação a declarar, foi %s (%s)", v, d)
	}
}

// PRECEDÊNCIA: o `identified_as` do NÓ ganha da forma automática da obrigação. É a única
// fonte que conhece a irregularidade real do projeto — `MetadataEntry` referenciado no
// plural, uma irmã no singular, sem regra derivável. Inverter esta ordem acusa 28 modelos
// que estão corretos (medido).
func TestObligation_identifiedAsDoNoGanhaDaFormaAutomatica(t *testing.T) {
	t.Run("OBHNB-I03: The node's own identified_as wins over the automatic form", func(t *testing.T) {})
	// O destino nomeia só o PLURAL. A forma automática derivaria METADATA_ENTRY_TABLE_NAME
	// (singular), que ali não existe — se o gate a usasse, acusaria um nó correto.
	root := tmpRoot(t, "const tables = [METADATA_ENTRIES_TABLE_NAME]\n")
	cfg := &config.Config{Obligations: []config.Obligation{{
		Name:         "pii-purgavel",
		When:         "carries: pii",
		MustAppearIn: []string{"purge.ts"},
		IdentifiedBy: "{{SCREAMING}}_TABLE_NAME",
	}}}
	comDeclaracao := "<!-- @anchors\n  code: MTENX\n  carries: pii\n" +
		"  identified_as: METADATA_ENTRIES_TABLE_NAME\n-->\n"
	if v, d := checkObligationHonored(comDeclaracao, obligNode(), root, nil, cfg); v != Pass {
		t.Fatalf("o `identified_as` do nó devia ser o token procurado, foi %s (%s)", v, d)
	}
	// A contraparte que prova que a precedência é real: sem a declaração, a forma
	// automática procura o singular e o nó correto é acusado.
	semDeclaracao := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n"
	if v, _ := checkObligationHonored(semDeclaracao, obligNode(), root, nil, cfg); v != Fail {
		t.Fatalf("sem a declaração a forma automática procura o singular e devia acusar, foi %s", v)
	}
}

// O gate NÃO DECIDE quais obrigações existem. Sem declaração na Estrutura, um nó que
// qualquer revisor chamaria de obviamente purgável não é cobrado — inventar deveres
// cobraria o que ninguém se comprometeu a fazer.
func TestObligation_naoInventaDever(t *testing.T) {
	t.Run("OBHNB-X01: The gate does not decide which obligations exist", func(t *testing.T) {})
	root := tmpRoot(t, "const tables = []\n")
	content := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n"
	if v, d := checkObligationHonored(content, obligNode(), root, nil, &config.Config{}); v != Skip {
		t.Fatalf("sem obrigação declarada nada é cobrado, foi %s (%s)", v, d)
	}
}

// A régua é a PRESENÇA, e só ela. Um script que nomeia o token num ramo morto e não apaga
// nada passa — separar esquecido de lembrado já é o defeito que este gate existe para
// pegar; julgar a implementação é outra régua.
func TestObligation_naoJulgaOQueODestinoFazComOToken(t *testing.T) {
	t.Run("OBHNB-X02: The gate does not understand what the destination does with the token", func(t *testing.T) {})
	root := tmpRoot(t, "if (false) { const morto = METADATA_ENTRY_TABLE }\n")
	content := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n"
	if v, d := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Pass {
		t.Fatalf("a presença basta — julgar o uso é outra régua, foi %s (%s)", v, d)
	}
}

// A declaração vale no HEADER, e só nele. Sem esse corte, uma citação na prosa — um
// exemplo, um trecho de outro documento — eximiria uma obrigação que ninguém quis eximir.
func TestObligation_declaracaoNoCorpoNaoVale(t *testing.T) {
	t.Run("OBHNB-X03: A declaration written in the body is not read", func(t *testing.T) {})
	root := tmpRoot(t, "// nada aqui\n")
	enchimento := strings.Repeat("Prosa da spec que empurra o corpo para longe do header.\n", 40)
	content := "<!-- @anchors\n  code: MTENX\n  carries: pii\n-->\n# x\n" + enchimento +
		"obligation_waived: pii-purgavel — este nó não guarda dado do titular\n"
	if v, d := checkObligationHonored(content, obligNode(), root, nil, obligCfg()); v != Fail {
		t.Fatalf("declaração fora do header não exime, foi %s (%s)", v, d)
	}
}

func TestObligationHonored_Errors(t *testing.T) {
	t.Run("OBHNB-E01: An unreadable destination file is not proof, and the others are still searched", func(t *testing.T) {
		cfg := &config.Config{Obligations: []config.Obligation{{
			Name:         "pii-purgavel",
			When:         "carries: pii",
			MustAppearIn: []string{"purge/*.sql"},
			IdentifiedBy: "screaming-snake",
		}}}
		content := "<!-- @anchors\n  carries: pii\n-->\n"
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "purge"), 0o755); err != nil {
			t.Fatal(err)
		}
		locked := filepath.Join(root, "purge", "a.sql")
		if err := os.WriteFile(locked, []byte("DELETE FROM METADATA_ENTRY;\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(locked, 0o000); err != nil {
			t.Skip("cannot chmod 0000 on this platform")
		}
		defer os.Chmod(locked, 0o644)
		if _, err := os.ReadFile(locked); err == nil {
			t.Skip("running with privileges that read a 0000 file")
		}
		if v, d := checkObligationHonored(content, obligNode(), root, nil, cfg); v != Fail || !strings.Contains(d, "purge/*.sql") {
			t.Fatalf("a token only inside an unreadable file must not count as present; got %v: %s", v, d)
		}
		if err := os.WriteFile(filepath.Join(root, "purge", "b.sql"), []byte("DELETE FROM METADATA_ENTRY;\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if v, d := checkObligationHonored(content, obligNode(), root, nil, cfg); v != Pass {
			t.Fatalf("the readable file that carries the token must still be found; got %v: %s", v, d)
		}
	})
}
