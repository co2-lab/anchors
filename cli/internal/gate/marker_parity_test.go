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
	g := gateDePara()
	g.MarkerPrefix = ""
	v, d := checkMarkerParity(g, arvore(t, map[string]string{"a.ts": "x"}), nil, nil)
	if v != Pending || !strings.Contains(d, "marker_prefix") {
		t.Fatalf("esperava Pending pedindo o prefixo; veio %v: %s", v, d)
	}
}

// Sem escopos, resta a contagem — o modo mais fraco, documentado como tal.
func TestMarkerParitySemEscoposUsaContagem(t *testing.T) {
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
	root := arvore(t, map[string]string{
		"web/scopes.ts":   "// @data-purge-rule-conta: promete\n",
		"server/purge.go": "// @data-purge-rule-conta: apaga\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("esperava Pass entre .ts e .go; veio %v: %s", v, d)
	}
}

func TestMarkerParityIgnoraNodeModules(t *testing.T) {
	root := arvore(t, map[string]string{
		"web/scopes.ts":              "// @data-purge-rule-conta: promete\n",
		"server/purge.ts":            "// @data-purge-rule-conta: apaga\n",
		"node_modules/lixo/index.ts": "// @data-purge-rule-conta: ruido\n",
	})
	if v, d := checkMarkerParity(gateDePara(), root, nil, nil); v != Pass {
		t.Fatalf("node_modules nao pode contar; veio %v: %s", v, d)
	}
}
