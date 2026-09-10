package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// monta um projeto com N arquivos e roda o gate sobre a spec dada.
func rodaContagem(t *testing.T, spec string, arquivos map[string]string) (Verdict, string) {
	t.Helper()
	dir := t.TempDir()
	for nome, conteudo := range arquivos {
		p := filepath.Join(dir, nome)
		os.MkdirAll(filepath.Dir(p), 0o755)
		os.WriteFile(p, []byte(conteudo), 0o644)
	}
	return checkCountHonored(spec, mapx.Node{Kind: mapx.KindSpec}, dir, nil, nil)
}

// O caso real: a spec dizia "os 50 modelos do produto", alguém adicionou o 51º, e a frase
// virou mentira sem que ninguém a editasse. Num único arquivo havia 7 afirmações
// numéricas — adicionar UM modelo tornava 10 frases obsoletas de uma vez.
func TestCountPegaContagemEnvelhecida(t *testing.T) {
	modelos := map[string]string{
		"models/A.ts": "x", "models/B.ts": "y", "models/C.ts": "z",
	}
	desatualizada := "# Spec\n<!-- @anchors-count: 2 modelos = models/*.ts -->\nO produto tem **2 modelos**.\n"
	v, d := rodaContagem(t, desatualizada, modelos)
	if v != Fail {
		t.Fatalf("contagem desatualizada deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "afirma **2 modelos**") || !strings.Contains(d, "tem **3**") {
		t.Errorf("a mensagem não mostra os dois números: %s", d)
	}

	// Corrigir SÓ o marcador não basta — a prosa continuaria mentindo, e é ela que o
	// leitor lê. Este teste mudou quando o gate passou a confrontar as duas: antes ele
	// trocava só o marcador e esperava verde, o que era o próprio defeito.
	certa := strings.NewReplacer("count: 2", "count: 3", "**2 modelos**", "**3 modelos**").
		Replace(desatualizada)
	if v, d := rodaContagem(t, certa, modelos); v != Pass {
		t.Fatalf("contagem correta deveria passar, foi %s (%s)", v, d)
	}
}

// Contar OCORRÊNCIAS, não arquivos: "51 cláusulas de autorização" é pergunta diferente de
// "51 modelos" sobre os mesmos arquivos. Nenhuma heurística acerta as duas — a spec diz.
func TestCountComPadraoContaOcorrencias(t *testing.T) {
	arquivos := map[string]string{
		"models/A.ts": "allow.owner()\nallow.group('admin')\n",
		"models/B.ts": "allow.owner()\n",
	}
	spec := "# Spec\n<!-- @anchors-count: 3 cláusulas = models/*.ts /allow\\.(owner|group)/ -->\n"
	if v, d := rodaContagem(t, spec, arquivos); v != Pass {
		t.Fatalf("3 ocorrências em 2 arquivos deveria passar, foi %s (%s)", v, d)
	}
	errada := strings.Replace(spec, "count: 3", "count: 2", 1)
	if v, _ := rodaContagem(t, errada, arquivos); v != Fail {
		t.Fatal("contar arquivos em vez de ocorrências deveria reprovar")
	}
}

// Glob que não casa nada quase nunca é "o código tem zero" — é o caminho errado. Dizer
// isso evita a conclusão errada: a correção é o glob, não o número.
func TestCountGlobVazioExplica(t *testing.T) {
	spec := "# Spec\n<!-- @anchors-count: 5 modelos = caminho/que/nao/existe/*.ts -->\n"
	v, d := rodaContagem(t, spec, map[string]string{"models/A.ts": "x"})
	if v != Fail {
		t.Fatalf("glob vazio deveria reprovar, foi %s", v)
	}
	if !strings.Contains(d, "confira o caminho") {
		t.Errorf("não explicou que o problema é o glob: %s", d)
	}
}

// O gate cobra o que foi DECLARADO, não todo número do texto. Um "90 dias" de retenção
// não é contagem de código, e acusá-lo seria ruído que faz desligar o gate.
func TestCountIgnoraNumeroNaoDeclarado(t *testing.T) {
	spec := "# Spec\n\nRetemos por 90 dias e o limite é 500 itens. O produto tem 3 modelos.\n"
	if v, _ := rodaContagem(t, spec, map[string]string{"models/A.ts": "x"}); v != Skip {
		t.Fatal("número em prosa sem declaração não é confrontável — deveria ser Skip")
	}
}

func TestCountSoSpec(t *testing.T) {
	spec := "<!-- @anchors-count: 9 = models/*.ts -->"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
		if v, _ := checkCountHonored(spec, mapx.Node{Kind: k}, t.TempDir(), nil, nil); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

// A PROSA ao lado do marcador também afirma o número — e é ela que o leitor lê.
// Achado real: o marcador dizia 51 (certo) e uma frase três linhas abaixo dizia "50
// cláusulas de autorização". O gate ficava ✓ e a spec mentia para quem a abrisse.
func TestCountConfrontaAProsaTambem(t *testing.T) {
	arquivos := map[string]string{
		"models/A.ts": "allow.owner()\n", "models/B.ts": "allow.owner()\n",
	}
	spec := `# Spec
<!-- @anchors-count: 2 cláusulas = models/*.ts /allow\.owner/ -->

O schema declara **2 cláusulas de autorização**.
`
	if v, d := rodaContagem(t, spec, arquivos); v != Pass {
		t.Fatalf("marcador e prosa de acordo deveriam passar, foi %s (%s)", v, d)
	}

	mentindo := strings.Replace(spec, "**2 cláusulas de autorização**", "**3 cláusulas de autorização**", 1)
	v, d := rodaContagem(t, mentindo, arquivos)
	if v != Fail {
		t.Fatalf("prosa divergente do disco deveria reprovar, foi %s (%s)", v, d)
	}
	if !strings.Contains(d, "PROSA diz **3**") {
		t.Errorf("não apontou a prosa: %s", d)
	}
}

// Frase sobre SUBCONJUNTO não é afirmação sobre o total. Acusá-la é o ruído que faz
// desligar o gate — e apareceu de imediato no texto real ("nos 46 modelos INDEXADOS",
// "em 16 modelos DE DADO FINANCEIRO", num arquivo cujo total é 50).
func TestCountNaoAcusaSubconjunto(t *testing.T) {
	arquivos := map[string]string{"models/A.ts": "x", "models/B.ts": "y", "models/C.ts": "z"}
	spec := `# Spec
<!-- @anchors-count: 3 modelos = models/*.ts -->

O produto tem **3 modelos**. Convenção sem exceção nos 2 modelos indexados.
Presente em 1 modelos de dado financeiro do usuário.
`
	if v, d := rodaContagem(t, spec, arquivos); v != Pass {
		t.Fatalf("subconjunto qualificado não é afirmação sobre o total, foi %s (%s)", v, d)
	}
}
