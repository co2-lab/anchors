package main

import (
	"strings"
	"testing"
)

// O DEFEITO QUE ISTO PEGA: um binário mais velho não falha — ele grava o formato que
// conhece e DESFAZ o que a versão nova escreveu. Medido depois de renomear um campo do
// carimbo: o build local anterior revertia 26 linhas a cada `check`, e o mapa oscilava
// entre dois formatos com conflito a cada PR.
func TestAvisaQuandoQuemGravouEhOutro(t *testing.T) {
	a := avisoDeBinarioVelho("0.1.9", "0.1.8")
	if a == "" {
		t.Fatal("versões diferentes têm de avisar")
	}
	// O aviso precisa NOMEAR as duas, senão quem lê não sabe qual atualizar.
	if !strings.Contains(a, "0.1.9") || !strings.Contains(a, "0.1.8") {
		t.Errorf("o aviso deve nomear as duas versões; veio: %s", a)
	}
	// E dizer o que fazer: um aviso que só constata o problema é ruído.
	if !strings.Contains(a, "map build") {
		t.Errorf("o aviso deve dizer como resolver; veio: %s", a)
	}
}

// COMPARA POR IGUALDADE, não por ordem — e isso é o que faz o aviso funcionar no caso
// real. Os dois builds que produziram o defeito se identificavam como "dev": ordenar
// versões não teria pego nenhum dos dois, porque "dev" não é ordenável contra "0.1.9".
func TestDoisDevDiferentesSaoIndistinguiveis(t *testing.T) {
	// Mesmo nome, ainda que sejam builds distintos: aqui não há o que fazer, e é a
	// limitação conhecida. O que NÃO pode é o inverso — calar quando os nomes diferem.
	if avisoDeBinarioVelho("dev", "dev") != "" {
		t.Error("nomes iguais não têm como ser distinguidos — avisar aqui seria ruído em todo build local")
	}
	if avisoDeBinarioVelho("dev", "0.1.9") == "" {
		t.Error("mapa gravado por build local e binário publicado é justamente o caso a avisar")
	}
	if avisoDeBinarioVelho("0.1.9", "dev") == "" {
		t.Error("o inverso também: quem roda build local sobre mapa publicado precisa saber")
	}
}

// MAPA SEM O CAMPO não pode avisar: todo mapa gerado antes desta versão não tem
// `gerado_por`, e acusá-los faria o primeiro `check` de todo projeto existente gritar
// sobre algo que ninguém pode consertar.
func TestMapaAntigoNaoAcusa(t *testing.T) {
	if avisoDeBinarioVelho("", "0.1.9") != "" {
		t.Error("mapa sem `gerado_por` é o de toda versão anterior — avisar seria acusar quem não errou")
	}
}
