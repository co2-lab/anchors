package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func featNode() mapx.Node { return mapx.Node{Kind: mapx.KindFeature, ID: "x.feature"} }

// Dois cenários com o MESMO código são indistinguíveis: nada liga um deles a um
// teste específico, e os gates relacionais comparam N títulos contra o mesmo teste.
func TestCenarioIdentidadeAcusaCodigoRepetido(t *testing.T) {
	t.Run("SCIDS-B01: Two scenarios sharing one code are reported", func(t *testing.T) {})
	v, msg := checkScenarioIdentity(`
  @USBPX-B01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`, featNode(), "", nil, nil)

	if v == Pass {
		t.Fatalf("código repetido passou (%s)", msg)
	}
	if v != Pending {
		t.Fatalf("veredito %v, queria Pending", v)
	}
}

// Apontar sem NOMEAR o código manda quem lê varrer a feature à procura do repetido.
func TestCenarioIdentidadeNomeiaOCodigoRepetido(t *testing.T) {
	t.Run("SCIDS-B02: The report names the repeated code", func(t *testing.T) {})
	_, msg := checkScenarioIdentity(`
  @SAUTX-B01 @nivel-unit
  Cenário: um caso só

  @USBPX-B01 @nivel-unit
  Cenário: busca pontos

  @USBPX-B01 @nivel-unit
  Cenário: sem usuário logado
`, featNode(), "", nil, nil)

	if !strings.Contains(msg, "USBPX-B01") {
		t.Errorf("mensagem sem o código repetido: %s", msg)
	}
	if strings.Contains(msg, "SAUTX-B01") {
		t.Errorf("mensagem acusa um código que aparece uma vez só: %s", msg)
	}
}

// QUANTOS cenários dividem o código separa a duplicata acidental do código emprestado
// numa regra inteira — e as duas pedem trabalhos diferentes.
func TestCenarioIdentidadeDizQuantosDividemOCodigo(t *testing.T) {
	t.Run("SCIDS-B03: The report says how many scenarios share the code", func(t *testing.T) {})
	_, msg := checkScenarioIdentity(`
  @USBPX-B01 @nivel-unit
  Cenário: primeiro

  @USBPX-B01 @nivel-unit
  Cenário: segundo

  @USBPX-B01 @nivel-unit
  Cenário: terceiro
`, featNode(), "", nil, nil)

	if !strings.Contains(msg, "3") {
		t.Errorf("a mensagem não diz quantos cenários dividem o código: %s", msg)
	}
	for _, titulo := range []string{"primeiro", "segundo", "terceiro"} {
		if !strings.Contains(msg, titulo) {
			t.Errorf("a mensagem não ajuda a reconhecer qual cenário é qual (falta %q): %s", titulo, msg)
		}
	}
}

// A saída é ensinada com o código REAL do projeto: um exemplo genérico obriga a
// traduzir a instrução antes de aplicá-la.
func TestCenarioIdentidadeEnsinaComOCodigoDoProjeto(t *testing.T) {
	t.Run("SCIDS-B04: The report teaches the way out with the project's own code", func(t *testing.T) {})
	_, msg := checkScenarioIdentity(`
  @USBPX-B01 @nivel-unit
  Cenário: busca pontos

  @USBPX-B01 @nivel-unit
  Cenário: sem usuário logado
`, featNode(), "", nil, nil)

	if !strings.Contains(msg, "USBPX-B01#01") {
		t.Errorf("a mensagem não ensina o sufixo sobre o código real do projeto: %s", msg)
	}
	if strings.Contains(msg, "XXXXX") {
		t.Errorf("a mensagem caiu no exemplo genérico: %s", msg)
	}
}

// Numerados, os dois cenários passam a ter identidade própria — que é o ponto.
func TestCenarioIdentidadeAceitaSufixo(t *testing.T) {
	t.Run("SCIDS-B05: The suffix gives each scenario its own identity", func(t *testing.T) {})
	v, msg := checkScenarioIdentity(`
  @USBPX-B01#01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01#02 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`, featNode(), "", nil, nil)

	if v != Pass {
		t.Errorf("veredito %v (%s), queria Pass — os sufixos distinguem os cenários", v, msg)
	}
}

// Uma regra com UM cenário é o caso comum: não pode acusar nada.
func TestCenarioIdentidadeNaoAcusaCodigoUnico(t *testing.T) {
	t.Run("SCIDS-B06: Distinct codes pass", func(t *testing.T) {})
	v, _ := checkScenarioIdentity(`
  @SAUTX-B01 @nivel-unit
  Cenário: Hidratar carrega a sessão
    Então o usuário fica disponível

  @SAUTX-B02 @nivel-unit
  Cenário: Login popula o usuário
    Então o usuário fica disponível
`, featNode(), "", nil, nil)
	if v != Pass {
		t.Errorf("veredito %v, queria Pass", v)
	}
}

// PENDING e não FAIL: numerar cenários é migração, e o gate nasce sobre uma base que
// não conhecia a notação. Reprovar travaria o projeto inteiro de uma vez.
func TestCenarioIdentidadeEhPendenteNaoReprovacao(t *testing.T) {
	t.Run("SCIDS-B07: The verdict is Pending and never a failure", func(t *testing.T) {})
	v, _ := checkScenarioIdentity(`
  @USBPX-B01 @nivel-unit
  Cenário: primeiro

  @USBPX-B01 @nivel-unit
  Cenário: segundo
`, featNode(), "", nil, nil)
	if v == Fail {
		t.Fatal("o gate reprovou — numerar cenários é migração, e isso travaria a base inteira")
	}
	if v != Pending {
		t.Fatalf("veredito %v, queria Pending", v)
	}
}

// O gate só fala de feature. Um nó de código ou spec não é assunto dele.
func TestCenarioIdentidadeSoOlhaFeature(t *testing.T) {
	t.Run("SCIDS-B08: An artifact that is not a feature leaves without a verdict", func(t *testing.T) {})
	repetido := "\n  @USBPX-B01\n  Cenário: a\n\n  @USBPX-B01\n  Cenário: b\n"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindSpec, mapx.KindTest, mapx.KindGuide} {
		if v, _ := checkScenarioIdentity(repetido, mapx.Node{Kind: k}, "", nil, nil); v != Skip {
			t.Errorf("veredito %v para kind %s, queria Skip", v, k)
		}
	}
}

// Feature sem cenário COM CÓDIGO não tem o que confrontar: a ausência do código é
// cobrança de outro gate.
func TestCenarioIdentidadeSemCenarioComCodigoNaoConfronta(t *testing.T) {
	t.Run("SCIDS-B09: A feature with no coded scenario leaves without a verdict", func(t *testing.T) {})
	v, _ := checkScenarioIdentity(`
  Cenário: um cenário sem código nenhum
    Então algo acontece

  Cenário: outro cenário sem código nenhum
    Então outra coisa acontece
`, featNode(), "", nil, nil)
	if v != Skip {
		t.Errorf("veredito %v, queria Skip — não há código a confrontar", v)
	}
}

// Vários repetidos vêm JUNTOS e em ordem estável: um relatório que muda de ordem entre
// execuções faz o achado parecer novo a cada rodada.
func TestCenarioIdentidadeAgrupaVariosEmOrdemEstavel(t *testing.T) {
	t.Run("SCIDS-B10: Several repeated codes are reported together in a stable order", func(t *testing.T) {})
	feature := `
  @ZZZZX-B01
  Cenário: z um

  @ZZZZX-B01
  Cenário: z dois

  @AAAAX-B01
  Cenário: a um

  @AAAAX-B01
  Cenário: a dois

  @MMMMX-B01
  Cenário: m um

  @MMMMX-B01
  Cenário: m dois
`
	_, msg := checkScenarioIdentity(feature, featNode(), "", nil, nil)
	for _, cod := range []string{"AAAAX-B01", "MMMMX-B01", "ZZZZX-B01"} {
		if !strings.Contains(msg, cod) {
			t.Fatalf("a mensagem não traz todos os repetidos (falta %s): %s", cod, msg)
		}
	}
	iA := strings.Index(msg, "AAAAX-B01")
	iM := strings.Index(msg, "MMMMX-B01")
	iZ := strings.Index(msg, "ZZZZX-B01")
	if !(iA < iM && iM < iZ) {
		t.Errorf("os repetidos não saem em ordem estável: %s", msg)
	}
}

// O endereço é o CÓDIGO; o título só ajuda a reconhecer. Um título longo inteiro na
// mensagem enterra o que importa.
func TestCenarioIdentidadeEncurtaTitulosLongos(t *testing.T) {
	t.Run("SCIDS-B11: Long titles are shortened in the report", func(t *testing.T) {})
	longo := strings.Repeat("a", 80)
	_, msg := checkScenarioIdentity("\n  @USBPX-B01\n  Cenário: "+longo+
		"\n\n  @USBPX-B01\n  Cenário: "+longo+"\n", featNode(), "", nil, nil)

	if strings.Contains(msg, longo) {
		t.Errorf("o título longo entrou inteiro na mensagem: %s", msg)
	}
	if !strings.Contains(msg, strings.Repeat("a", 40)) {
		t.Errorf("a mensagem não guardou o começo do título para reconhecer o cenário: %s", msg)
	}
	if !strings.Contains(msg, "…") {
		t.Errorf("o truncamento não é sinalizado: %s", msg)
	}

	// O título que CABE inteiro sai inteiro: marcar reticências sem ter cortado nada
	// afirma um corte que não houve, e manda quem lê procurar um resto inexistente.
	cabe := strings.Repeat("b", 40)
	_, msgCabe := checkScenarioIdentity("\n  @USBPX-B01\n  Cenário: "+cabe+
		"\n\n  @USBPX-B01\n  Cenário: "+cabe+"\n", featNode(), "", nil, nil)
	if !strings.Contains(msgCabe, cabe) {
		t.Errorf("o título que cabe foi cortado: %s", msgCabe)
	}
	if strings.Contains(msgCabe, "…") {
		t.Errorf("reticências sem corte — afirma um resto que não existe: %s", msgCabe)
	}
}

// Agrupar pelo PREFIXO acusaria exatamente quem já fez a migração que o gate pede.
func TestCenarioIdentidadeAgrupaPeloCodigoCompleto(t *testing.T) {
	t.Run("SCIDS-I01: Grouping is by the complete code, suffix included", func(t *testing.T) {})
	v, msg := checkScenarioIdentity(`
  @USBPX-B01#01
  Cenário: primeiro

  @USBPX-B01#02
  Cenário: segundo

  @USBPX-B01#03
  Cenário: terceiro
`, featNode(), "", nil, nil)
	if v != Pass {
		t.Fatalf("veredito %v (%s) — quem já migrou não pode ser acusado", v, msg)
	}
}

// A mensagem tem de ser DETERMINÍSTICA: a varredura de um mapa em Go sai em ordem
// aleatória, e um achado que muda de texto a cada rodada parece novo toda vez.
func TestCenarioIdentidadeMensagemDeterministica(t *testing.T) {
	t.Run("SCIDS-I02: The message is deterministic", func(t *testing.T) {})
	feature := `
  @ZZZZX-B01
  Cenário: z um

  @ZZZZX-B01
  Cenário: z dois

  @AAAAX-B01
  Cenário: a um

  @AAAAX-B01
  Cenário: a dois

  @MMMMX-B01
  Cenário: m um

  @MMMMX-B01
  Cenário: m dois
`
	_, primeira := checkScenarioIdentity(feature, featNode(), "", nil, nil)
	for i := 0; i < 20; i++ {
		if _, outra := checkScenarioIdentity(feature, featNode(), "", nil, nil); outra != primeira {
			t.Fatalf("a mensagem mudou entre execuções:\n%s\n%s", primeira, outra)
		}
	}
}

// A régua é IDENTIDADE, não julgamento: o gate não decide se os dois cenários descrevem
// comportamentos diferentes — só que o código não distingue os dois.
func TestCenarioIdentidadeNaoJulgaOsComportamentos(t *testing.T) {
	t.Run("SCIDS-X01: The gate does not judge whether the two scenarios describe different behaviours", func(t *testing.T) {})
	// dois cenários descrevendo a MESMA coisa: acusado do mesmo jeito
	mesmo, msgMesmo := checkScenarioIdentity("\n  @USBPX-B01\n  Cenário: o mesmo caso\n\n  @USBPX-B01\n  Cenário: o mesmo caso\n",
		featNode(), "", nil, nil)
	// dois descrevendo coisas claramente diferentes: acusado do mesmo jeito
	difer, _ := checkScenarioIdentity("\n  @USBPX-B01\n  Cenário: salva no banco\n\n  @USBPX-B01\n  Cenário: envia um email\n",
		featNode(), "", nil, nil)

	if mesmo != Pending || difer != Pending {
		t.Fatalf("o gate distinguiu os casos: mesmo=%s difer=%s (%s)", mesmo, difer, msgMesmo)
	}
}

// Um código repetido em DUAS features é outro defeito, de outro dono. Este gate não tem
// grafo — julgar isso aqui exigiria o que ele não recebe.
func TestCenarioIdentidadeNaoOlhaEntreFeatures(t *testing.T) {
	t.Run("SCIDS-X02: The gate does not look across features", func(t *testing.T) {})
	v, msg := checkScenarioIdentity(`
  @USBPX-B01
  Cenário: o único desta feature
`, featNode(), "", nil, nil)
	if v != Skip && v != Pass {
		t.Fatalf("veredito %v (%s) — dentro desta feature não há repetição", v, msg)
	}
	// e o mesmo código noutra feature, confrontada em separado, também passa
	if outra, _ := checkScenarioIdentity("\n  @USBPX-B01\n  Cenário: o único da outra\n",
		mapx.Node{Kind: mapx.KindFeature, ID: "y.feature"}, "", nil, nil); outra == Pending {
		t.Error("o gate enxergou entre features — ele não recebe grafo para isso")
	}
}

// Cenário SEM código é invisível para o parser e para os gates relacionais. Cobrar a
// ausência é régua do gate que pareia cenário e teste.
func TestCenarioIdentidadeNaoCobraAusenciaDeCodigo(t *testing.T) {
	t.Run("SCIDS-X03: The gate does not charge the absence of a code on a scenario", func(t *testing.T) {})
	v, msg := checkScenarioIdentity(`
  @USBPX-B01
  Cenário: este tem código

  Cenário: este não tem código nenhum

  @USBPX-B02
  Cenário: este tem outro código
`, featNode(), "", nil, nil)
	if v != Pass {
		t.Errorf("veredito %v (%s) — todo código presente aparece uma vez só", v, msg)
	}
}

// O gate aponta, não renumera: o sufixo carrega significado e reescrever a feature
// invalidaria todo teste já amarrado ao código antigo.
func TestCenarioIdentidadeNaoRenumera(t *testing.T) {
	t.Run("SCIDS-X04: The gate does not renumber the scenarios", func(t *testing.T) {})
	feature := "\n  @USBPX-B01\n  Cenário: primeiro\n\n  @USBPX-B01\n  Cenário: segundo\n"
	antes := feature

	checkScenarioIdentity(feature, featNode(), "", nil, nil)

	if feature != antes {
		t.Error("o gate alterou a feature — ele aponta, não renumera")
	}
	// e a mensagem não afirma ter consertado: ela ENSINA o sufixo
	_, msg := checkScenarioIdentity(feature, featNode(), "", nil, nil)
	if !strings.Contains(msg, "#01") {
		t.Errorf("a mensagem deveria ensinar a numeração em vez de aplicá-la: %s", msg)
	}
}
