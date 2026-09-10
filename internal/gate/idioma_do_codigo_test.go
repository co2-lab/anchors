package gate

import "testing"

// A DETECÇÃO acerta os dois lados — e o segundo é o que decide se o gate sobrevive.
//
// Um gate que acusa código correto é pior que gate nenhum: a saída barata vira desligá-lo.
// A primeira tentativa desta detecção comparava contra o dicionário do sistema e tinha 90%
// de falso positivo (`AbsRoot`, `FileOwner`, `Classify` — palavras COMPOSTAS em inglês,
// que nenhum dicionário comum tem).
func TestIdioma_acusaPortuguesEDeixaOInglesEmPaz(t *testing.T) {
	portugues := []string{
		"carimbo", "DonoDoArquivo", "checkFaseExiste", "alvoDaSpec",
		"pecasPorDesenvolver", "escreveCongelamento", "trincaCompleta",
		"validaVeredito", "imprimeOrdemDeServico", "linhaDeRegraRE",
	}
	for _, id := range portugues {
		if _, ehPT := IdentificadorEhPT(id); !ehPT {
			t.Errorf("%q é português e passou", id)
		}
	}

	// O INGLÊS que a detecção ingênua acusaria. Cada um destes existe no projeto hoje.
	ingles := []string{
		"AbsRoot", "FileOwner", "Classify", "stamp", "Rule", "Waiver",
		"checkPhaseExists", "writeFreeze", "commitAndPush", "progressPath",
		"IsProgressFile", "RequiredApprovals", "SectionTitle", "targetsOf",
		"declaredStamps", "firstLine", "hasCommit", "listValue",
		// os que a morfologia sozinha acusaria
		"metadata", "schemaVersion", "deltaRange", "areaCode", "upgradePath",
	}
	for _, id := range ingles {
		if w, ehPT := IdentificadorEhPT(id); ehPT {
			t.Errorf("%q é inglês legítimo e foi acusado (pela palavra %q)", id, w)
		}
	}
}

// A palavra ACUSADORA é devolvida, e não só o veredito.
//
// Sem ela, quem lê a reprovação de `checkPlanoAlteradoJustificado` tem quatro candidatas
// para adivinhar qual está errada.
func TestIdioma_devolveAPalavraQueAcusou(t *testing.T) {
	w, ehPT := IdentificadorEhPT("checkPlanoAlteradoJustificado")
	if !ehPT {
		t.Fatal("não acusou")
	}
	if w != "plano" && w != "alterado" {
		t.Errorf("a palavra acusadora foi %q — esperava a primeira em português", w)
	}
}

// SÓ DECLARAÇÕES em coluna zero.
//
// Um `func` dentro de string ou indentado num comentário não é declaração, e acusá-lo
// faria o gate reprovar exemplo de documentação — que é justamente onde o português é
// legítimo.
func TestIdioma_ignoraOQueNaoEhDeclaracao(t *testing.T) {
	fonte := "package x\n" +
		"// func carimbo() — o exemplo em prosa NÃO conta\n" +
		"const exemplo = \"func alvoDaSpec() string\"\n" +
		"\tfunc indentadoNaoEhDeclaracao() {}\n" +
		"func stamp() {}\n"
	achados := IdentificadoresPT(fonte)
	if len(achados) != 0 {
		t.Errorf("acusou o que não é declaração de topo: %v", achados)
	}
}

func TestIdioma_achaAsDeclaracoesDeVerdade(t *testing.T) {
	fonte := "package x\n" +
		"func carimbo() {}\n" +
		"type DonoDoArquivo struct{}\n" +
		"var alvoPadrao = 1\n" +
		"func stamp() {}\n" +
		"type FileOwner struct{}\n"
	achados := IdentificadoresPT(fonte)
	if len(achados) != 3 {
		t.Fatalf("achou %d, esperava 3 (carimbo, DonoDoArquivo, alvoPadrao): %v", len(achados), achados)
	}
	for _, esperado := range []string{"carimbo", "DonoDoArquivo", "alvoPadrao"} {
		if _, ok := achados[esperado]; !ok {
			t.Errorf("não achou %q", esperado)
		}
	}
}

// Palavra curta não é português: siglas e contadores (`id`, `n`, `ok`) atravessariam a
// morfologia por acaso.
func TestIdioma_palavraCurtaNaoConta(t *testing.T) {
	for _, w := range []string{"id", "n", "ok", "a", "eh"} {
		if PalavraEhPT(w) {
			t.Errorf("%q é curta demais para ser acusada", w)
		}
	}
}
