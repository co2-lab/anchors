package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// arvore monta um projeto temporário com os arquivos dados.
func arvore(t *testing.T, arquivos map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for nome, conteudo := range arquivos {
		caminho := filepath.Join(root, nome)
		if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func gateDePara() config.Gate {
	return config.Gate{
		Name:         "paridade-exclusao",
		ID:           "data-purge-parity",
		Check:        "marker-parity",
		MarkerPrefix: "data-purge-rule",
		MarkerCount:  2,
		MarkerScopes: []string{"web/**", "server/**"},
	}
}

func TestMarkerParityPassaComAsDuasPontas(t *testing.T) {
	t.Run("MRPRM-B01: A rule marked at both declared scopes passes", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: o que a página promete apagar\n",
		"server/purge.ts": "// @data-purge-rule-conta: o que o handler apaga\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Pass {
		t.Fatalf("esperava Pass, veio %v: %s", v, d)
	}
}

// O CASO QUE O GATE EXISTE PARA PEGAR: alguém mexeu num lado só.
func TestMarkerParityPegaLadoQueFaltou(t *testing.T) {
	t.Run("MRPRM-B02: A rule missing from one end fails, and the verdict names the empty scope", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts": "// @data-purge-rule-conta: a página promete apagar\n",
		// o servidor NÃO tem a marcação — a regra existe só de um lado
		"server/purge.ts": "// nada aqui\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v: %s", v, d)
	}
	if !strings.Contains(d, "server/**") {
		t.Fatalf("o laudo tem de NOMEAR o lado que faltou; veio: %s", d)
	}
}

// A razão de `marker_scopes` existir: contar sem olhar ONDE deixaria isto passar.
func TestMarkerParityNaoAceitaDuasNoMesmoLado(t *testing.T) {
	t.Run("MRPRM-B03: Two markings on the same side do not satisfy the gate", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts": "// @data-purge-rule-conta: promete\n",
		"web/outro.ts":  "// @data-purge-rule-conta: promete de novo\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Fail {
		t.Fatalf("duas do MESMO lado somam 2 e nao podem passar; veio %v: %s", v, d)
	}
}

// Regra que sobra: apagada de um lado, esquecida no outro.
func TestMarkerParityPegaRegraOrfa(t *testing.T) {
	t.Run("MRPRM-B04: A rule left over at one end fails and is named", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: ok\n// @data-purge-rule-antiga: sobrou\n",
		"server/purge.ts": "// @data-purge-rule-conta: ok\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Fail || !strings.Contains(d, "antiga") {
		t.Fatalf("esperava Fail nomeando `antiga`; veio %v: %s", v, d)
	}
}

// Ausência TOTAL não é aprovação — quase sempre é prefixo errado na declaração.
func TestMarkerParityAusenciaTotalNaoAprova(t *testing.T) {
	t.Run("MRPRM-B05: Total absence of the prefix is not approval", func(t *testing.T) {})
	root := arvore(t, map[string]string{"web/scopes.ts": "// nada\n"})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Pending {
		t.Fatalf("esperava Pending, veio %v: %s", v, d)
	}
	if !strings.Contains(d, "marker_prefix") {
		t.Fatalf("o laudo tem de sugerir conferir o prefixo; veio: %s", d)
	}
}

func TestMarkerParitySemPrefixoNaoAprova(t *testing.T) {
	t.Run("MRPRM-B06: A declaration with no prefix returns Pending", func(t *testing.T) {})
	g := gateDePara()
	g.MarkerPrefix = ""
	v, d := checkMarkerParity(g, arvore(t, map[string]string{"a.ts": "x"}), nil, nil)
	if v != Pending || !strings.Contains(d, "marker_prefix") {
		t.Fatalf("esperava Pending pedindo o prefixo; veio %v: %s", v, d)
	}
}

// Sem escopos, resta a contagem — o modo mais fraco, documentado como tal.
func TestMarkerParitySemEscoposUsaContagem(t *testing.T) {
	t.Run("MRPRM-B07: With no scopes declared the ruler is the count", func(t *testing.T) {})
	g := gateDePara()
	g.MarkerScopes = nil
	root := arvore(t, map[string]string{
		"a.ts": "// @data-purge-rule-conta: um\n",
		"b.ts": "// @data-purge-rule-conta: dois\n",
	})
	if v, d := checkMarkerParity(g, root, nil, nil); v != Pass {
		t.Fatalf("2 ocorrencias com count=2 devia passar; veio %v: %s", v, d)
	}
	root2 := arvore(t, map[string]string{"a.ts": "// @data-purge-rule-conta: so uma\n"})
	if v, _ := checkMarkerParity(g, root2, nil, nil); v != Fail {
		t.Fatalf("1 ocorrencia com count=2 devia reprovar; veio %v", v)
	}
}

// A marcação vale em qualquer linguagem da árvore — o de-para cruza fronteira.
func TestMarkerParityCruzaLinguagens(t *testing.T) {
	t.Run("MRPRM-B08: The marking crosses language", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: promete\n",
		"server/purge.go": "// @data-purge-rule-conta: apaga\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("esperava Pass entre .ts e .go; veio %v: %s", v, d)
	}
}

func TestMarkerParityIgnoraNodeModules(t *testing.T) {
	t.Run("MRPRM-B09: Ignored directories never count towards parity", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":              "// @data-purge-rule-conta: promete\n",
		"server/purge.ts":            "// @data-purge-rule-conta: apaga\n",
		"node_modules/lixo/index.ts": "// @data-purge-rule-conta: ruido\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("node_modules nao pode contar; veio %v: %s", v, d)
	}
}

// O laudo tem de NOMEAR a regra e o escopo vazio: nenhum dos dois lados, olhado
// sozinho, parece errado — quem le precisa do nome para achar o orfao.
func TestMarkerParityLaudoNomeiaRegraEEscopo(t *testing.T) {
	t.Run("MRPRM-I01: The failing verdict names the rule and the empty scope", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta-completa: a pagina promete\n",
		"server/purge.ts": "// nada aqui\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v: %s", v, d)
	}
	if !strings.Contains(d, "conta-completa") {
		t.Fatalf("o laudo tem de nomear a REGRA; veio: %s", d)
	}
	if !strings.Contains(d, "server/**") {
		t.Fatalf("o laudo tem de nomear o ESCOPO vazio; veio: %s", d)
	}
}

// O que nao foi medido nunca e aprovado: sem prefixo, sem contagem nem escopos, e
// ausencia total — os tres devolvem Pending, jamais Pass.
func TestMarkerParityNadaMedidoNuncaAprova(t *testing.T) {
	t.Run("MRPRM-I02: What was not measured is never approved", func(t *testing.T) {})
	semPrefixo := gateDePara()
	semPrefixo.MarkerPrefix = ""

	semRegua := gateDePara()
	semRegua.MarkerCount = 0
	semRegua.MarkerScopes = nil

	casos := map[string]struct {
		g     config.Gate
		files map[string]string
	}{
		"sem prefixo":    {semPrefixo, map[string]string{"web/a.ts": "// @data-purge-rule-conta: x\n"}},
		"sem regua":      {semRegua, map[string]string{"web/a.ts": "// @data-purge-rule-conta: x\n"}},
		"ausencia total": {gateDePara(), map[string]string{"web/a.ts": "// nada\n"}},
	}
	for nome, caso := range casos {
		v, d := checkMarkerParity(caso.g, arvore(t, caso.files), nil, nil)
		if v == Pass {
			t.Fatalf("%s: aprovar sem ter olhado carimba o que nao foi medido; veio %v: %s", nome, v, d)
		}
		if v != Pending {
			t.Fatalf("%s: esperava Pending, veio %v: %s", nome, v, d)
		}
	}
}

// O gate confronta PRESENCA nas duas pontas, nao concordancia de sentido: a pagina
// promete uma lista que o handler nao apaga, e ainda assim passa.
func TestMarkerParityNaoLeOQueCadaPontaFaz(t *testing.T) {
	t.Run("MRPRM-X01: The gate does not read what each end actually does", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: apaga perfil, historico e recibos\n",
		"server/purge.ts": "// @data-purge-rule-conta: apaga so o perfil\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("presenca e deterministica, sentido nao; esperava Pass, veio %v: %s", v, d)
	}
}

// Quais regras vivem em duas pontas e decisao do PROJETO: um prefixo que o projeto
// nao declarou nao e cobrado, mesmo com a marcacao obvia so de um lado.
func TestMarkerParityNaoInventaParidade(t *testing.T) {
	t.Run("MRPRM-X02: The gate does not decide which rules live at two ends", func(t *testing.T) {})
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: promete\n// @outra-regra-qualquer: so de um lado\n",
		"server/purge.ts": "// @data-purge-rule-conta: apaga\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("cobrar prefixo nao declarado seria inventar paridade; veio %v: %s", v, d)
	}
}

// Extensao fora da lista de texto nao e lida — errar para menos custa uma marcacao
// em lugar exotico; errar para mais custa megabytes por varredura.
func TestMarkerParityNaoLeForaDaListaDeTexto(t *testing.T) {
	t.Run("MRPRM-X03: Files outside the text extension list are not read", func(t *testing.T) {})
	marcado := "// @data-purge-rule-conta: ponta\n"
	root := arvore(t, map[string]string{
		"web/scopes.ts":   marcado,
		"server/purge.ts": marcado,
		// mesma sequencia de bytes, extensao fora da lista: nao pode ser lida
		"web/asset.png": "@data-purge-rule-fantasma: ruido binario\n",
	})
	v, d := checkMarkerParity(gateDePara(), root, nil, nil)
	if v != Pass {
		t.Fatalf("o .png nao podia ser varrido; veio %v: %s", v, d)
	}
	if strings.Contains(d, "fantasma") {
		t.Fatalf("a regra do binario nao podia aparecer no laudo; veio: %s", d)
	}
}

func TestMarkerParity_Errors(t *testing.T) {
	t.Run("MRPRM-E01: A marking inside an unreadable directory counts as absent from its end", func(t *testing.T) {
		root := arvore(t, map[string]string{
			"web/scopes.ts":          "// @data-purge-rule-conta: o que a página promete\n",
			"server/locked/purge.ts": "// @data-purge-rule-conta: o que o handler apaga\n",
		})
		locked := filepath.Join(root, "server", "locked")
		if err := os.Chmod(locked, 0o000); err != nil {
			t.Skip("cannot chmod 0000 on this platform")
		}
		defer os.Chmod(locked, 0o755)
		if _, err := os.ReadDir(locked); err == nil {
			t.Skip("running with privileges that read a 0000 directory")
		}
		v, d := checkMarkerParity(gateDePara(), root, nil, nil)
		if v != Fail || !strings.Contains(d, "conta") || !strings.Contains(d, "server/**") {
			t.Fatalf("the unreadable end must be charged as missing, naming rule and scope; got %v: %s", v, d)
		}
	})
}
