package gate

import (
	"strings"
	"testing"
)

// O motivo é a única coisa que separa dispensa DELIBERADA de gate ignorado. Aceitá-lo
// ausente esvaziaria a garantia — e o relatório passaria a mostrar "waived" sem
// dizer por quê.
func TestDispensaExigeMotivo(t *testing.T) {
	t.Run("RLUEX-B05: A waiver with no reason is refused", func(t *testing.T) {})

	_, erros := ParseWaiver("triad-complete")

	if len(erros) != 1 {
		t.Fatalf("dispensa sem motivo tem de ser recusada, veio %d erro(s)", len(erros))
	}
	if _, erros := ParseWaiver("triad-complete="); len(erros) != 1 {
		t.Error("motivo vazio é o mesmo que ausente")
	}
}

// Sem a regra não há o que dispensar — e aceitar a entrada em silêncio guardaria um
// motivo sob a chave vazia, que nenhuma pergunta alcança.
func TestDispensaExigeARegra(t *testing.T) {
	t.Run("RLUEX-B06: A waiver with no rule name is refused", func(t *testing.T) {})

	_, erros := ParseWaiver("=um motivo qualquer")
	if len(erros) != 1 {
		t.Fatalf("entrada sem regra tem de ser recusada, veio %d erro(s): %v", len(erros), erros)
	}
	if !strings.Contains(erros[0], "regra") {
		t.Errorf("o erro deveria dizer o que falta: %s", erros[0])
	}
}

// Dispensar o GATE cobre todas as regras dele; dispensar a REGRA preserva o resto. As
// duas granularidades existem porque quem não conhece as regras precisa da saída grossa.
func TestDispensaAceitaAsDuasGranularidades(t *testing.T) {
	t.Run("RLUEX-B09: The waiver accepts both granularities", func(t *testing.T) {})

	d, erros := ParseWaiver("spec-complete/sem-placeholder=a spec nasce em rascunho")
	if len(erros) != 0 {
		t.Fatalf("erros inesperados: %v", erros)
	}

	if _, ok := d.Waived("spec-complete/sem-placeholder"); !ok {
		t.Error("a regra dispensada não foi reconhecida")
	}
	// A OUTRA regra do mesmo gate continua valendo — é o ponto de dispensar por regra.
	if _, ok := d.Waived("spec-complete/tem-regra-catalogada"); ok {
		t.Error("dispensar uma regra não pode desligar as demais do mesmo gate")
	}

	// E dispensar o gate inteiro cobre as regras dele.
	dg, _ := ParseWaiver("spec-complete=projeto em bootstrap")
	if _, ok := dg.Waived("spec-complete/sem-placeholder"); !ok {
		t.Error("dispensar o gate tem de cobrir suas regras")
	}
}

// Um gate que não foi dispensado precisa continuar rodando — o erro que mais custaria
// aqui é uma dispensa vazando para o que ninguém pediu.
func TestDispensaNaoAlcancaOQueNaoFoiPedido(t *testing.T) {
	t.Run("RLUEX-I01: A waiver never reaches what nobody waived", func(t *testing.T) {})

	d, _ := ParseWaiver("triad-complete=a feature ainda é um card")

	if _, ok := d.Waived("guide-checklist"); ok {
		t.Error("a dispensa alcançou um gate que ninguém dispensou")
	}
	if _, ok := (Waiver{}).Waived("qualquer-coisa"); ok {
		t.Error("dispensa vazia não dispensa nada")
	}
}

// O ID é `<gate>/<regra>`, e o gate faz parte de propósito: nomes de regra curtos se
// repetiriam entre gates, e um ID que colide não identifica nada.
func TestRegraIDMontaOIdentificador(t *testing.T) {
	t.Run("RLUEX-B01: Building an identifier joins the gate and the rule", func(t *testing.T) {})

	id := NewRuleID("spec-complete", "sem-placeholder")
	if id != "spec-complete/sem-placeholder" {
		t.Errorf("ID montado errado: %q", id)
	}
	// Duas regras de mesmo nome em gates diferentes NÃO colidem — é o que o prefixo
	// de gate existe para garantir.
	if NewRuleID("spec-complete", "tem-codigo") == NewRuleID("header-valid", "tem-codigo") {
		t.Error("regras homônimas de gates distintos não podem produzir o mesmo ID")
	}
}

// Gate com uma verificação só: o nome já a identifica, e não há barra a inventar.
func TestRegraIDSemRegraNaoGanhaBarra(t *testing.T) {
	t.Run("RLUEX-B02: A gate with a single verification gains no separator", func(t *testing.T) {})

	if s := NewRuleID("layer-boundary", ""); s != "layer-boundary" {
		t.Errorf("gate sem regra não deve ganhar barra: %q", s)
	}
}

// A decomposição é o que permite falar do gate sem reler o ID inteiro — e é como o
// `Waived` acha a dispensa grossa a partir da regra fina.
func TestRegraIDSeparaGateDeRegra(t *testing.T) {
	t.Run("RLUEX-B03: The identifier decomposes into gate and rule", func(t *testing.T) {})

	id := NewRuleID("spec-complete", "sem-placeholder")
	if id.Gate() != "spec-complete" || id.Rule() != "sem-placeholder" {
		t.Errorf("decomposição errada: gate=%q regra=%q", id.Gate(), id.Rule())
	}
}

// Sem a parte de regra o veredito é do GATE INTEIRO, e a parte de regra tem de vir
// vazia — não o nome do gate repetido, que faria o relatório inventar uma regra.
func TestRegraVaziaQuandoOIDEhSoOGate(t *testing.T) {
	t.Run("RLUEX-B04: The rule half is empty when the identifier carries only a gate", func(t *testing.T) {})

	id := RuleID("layer-boundary")
	if id.Rule() != "" {
		t.Errorf("sem barra não há regra a nomear, veio %q", id.Rule())
	}
	if id.Gate() != "layer-boundary" {
		t.Errorf("o gate continua sendo o ID inteiro, veio %q", id.Gate())
	}
}

// O CASO QUE MOTIVOU A DISPENSA POR ALVO.
//
// Num projeto com a trinca completa e tudo passando, um plano novo semeia specs sem
// código. Dispensar `trinca-completa` para commitá-las apagava o gate para o REPOSITÓRIO
// INTEIRO — e uma trinca que quebrou por descuido noutro lugar passava junto, sem que
// nada acusasse. É o mascaramento que a dispensa por regra existe para evitar, um nível
// acima.
func TestDispensaPorAlvoNaoApagaOResto(t *testing.T) {
	t.Run("RLUEX-B11: A waiver by target spares the named codes and confronts the rest", func(t *testing.T) {})

	d, erros := ParseWaiver(
		"triad-complete@NOVOA=spec nova do plano 0007," +
			"triad-complete@NOVOB=spec nova do plano 0007")
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}

	id := RuleID("triad-complete")

	// Os alvos NOMEADOS estão dispensados.
	for _, cod := range []string{"NOVOA", "NOVOB"} {
		if motivo, ok := d.WaivedTarget(id, cod); !ok || motivo == "" {
			t.Errorf("%s deveria estar dispensado com motivo", cod)
		}
	}

	// E O RESTO CONTINUA SENDO CONFRONTADO. É o ponto inteiro do recurso: a trinca que
	// quebrou por descuido noutra unidade não pode passar de carona.
	if _, ok := d.WaivedTarget(id, "QBRDA"); ok {
		t.Error("um alvo não nomeado NÃO pode ser dispensado — é o mascaramento que este " +
			"recurso existe para impedir")
	}
}

// A pergunta SEM alvo responde "não dispensado", para que o filtro de gates não remova o
// gate da lista: ele precisa RODAR para confrontar os outros alvos.
func TestDispensaComAlvosNaoValeParaOGateInteiro(t *testing.T) {
	t.Run("RLUEX-B10: A waiver that declares targets does not hold for the whole gate", func(t *testing.T) {})

	d, _ := ParseWaiver(
		"triad-complete@NOVOA=spec nova do plano 0007," +
			"triad-complete@NOVOB=spec nova do plano 0007")
	if _, ok := d.Waived(RuleID("triad-complete")); ok {
		t.Error("dispensa COM alvos não pode valer para o gate inteiro — sairia da lista " +
			"e não confrontaria ninguém")
	}
}

// A dispensa SEM alvo continua valendo para tudo: há casos legítimos, como um gate
// recém-declarado que o projeto ainda não cumpre em lugar nenhum.
func TestDispensaSemAlvoValeParaTudo(t *testing.T) {
	t.Run("RLUEX-B12: A waiver with no declared target holds for every code", func(t *testing.T) {})

	d, _ := ParseWaiver("triad-complete=gate novo, nenhuma unidade o cumpre ainda")
	id := RuleID("triad-complete")
	if _, ok := d.Waived(id); !ok {
		t.Error("sem alvo declarado, a dispensa vale para o gate inteiro")
	}
	if _, ok := d.WaivedTarget(id, "QUALQ"); !ok {
		t.Error("sem alvo declarado, qualquer caminho está dispensado")
	}
}

// `@` sem caminho é engano de digitação, e aceitá-lo em silêncio produziria uma dispensa
// que não dispensa nada — o commit reprovaria sem explicação aparente.
func TestDispensaAlvoVazioEhRecusada(t *testing.T) {
	t.Run("RLUEX-B08: A target marker with nothing after it is refused", func(t *testing.T) {})

	_, erros := ParseWaiver("triad-complete@=motivo qualquer")
	if len(erros) == 0 {
		t.Error("`regra@=motivo` deveria ser recusado: falta o caminho")
	}
}

// O CAMINHO é RECUSADO como alvo, e a recusa é o ponto: ele não é identidade. Muda
// quando alguém reorganiza pastas, e a dispensa deixaria de valer em silêncio — o commit
// seguinte reprovaria sem que nada explicasse o que mudou.
//
// Aceitá-lo e nunca casar seria pior: uma dispensa que não dispensa, sem erro visível.
func TestDispensaRecusaCaminhoComoAlvo(t *testing.T) {
	t.Run("RLUEX-B07: A path is refused as the waiver target", func(t *testing.T) {})

	for _, bruto := range []string{
		"triad-complete@packages/shared/Workspace.spec.md=motivo",
		"triad-complete@packages/*=motivo",
		"triad-complete@arquivo.spec.md=motivo",
	} {
		_, erros := ParseWaiver(bruto)
		if len(erros) == 0 {
			t.Errorf("%q deveria ser recusado: o alvo é o CÓDIGO, não o caminho", bruto)
			continue
		}
		if !strings.Contains(erros[0], "CÓDIGO") {
			t.Errorf("o erro deveria dizer o que usar no lugar: %s", erros[0])
		}
	}
}

// O mesmo caminho que um revisor consideraria óbvio: ele é recusado, e nada fica
// dispensado por ele. A recusa vale mais do que a conveniência.
func TestCaminhoObvioNaoDispensaNada(t *testing.T) {
	t.Run("RLUEX-X01: The unit does not accept a path as the waiver target", func(t *testing.T) {})

	d, erros := ParseWaiver("triad-complete@internal/gate/rule.spec.md=parece óbvio")
	if len(erros) == 0 {
		t.Fatal("o caminho tem de ser recusado, por mais natural que pareça")
	}
	// E nada foi guardado: a entrada recusada não pode dispensar de lado.
	if _, ok := d.Waived(RuleID("triad-complete")); ok {
		t.Error("uma entrada recusada não pode produzir dispensa nenhuma")
	}
}

// Um alvo SEM CÓDIGO não é alcançado por uma dispensa restrita. Artefato sem identidade
// é um problema anterior — quem cobra isso é o `codigo-catalogado`, e dar uma saída
// lateral aqui esconderia a causa.
func TestDispensaPorAlvoNaoAlcancaQuemNaoTemCodigo(t *testing.T) {
	t.Run("RLUEX-B13: An artifact with no code is not reached by a waiver restricted to targets", func(t *testing.T) {})

	d, _ := ParseWaiver("triad-complete@WRKSP=spec nova")
	if _, ok := d.WaivedTarget(RuleID("triad-complete"), ""); ok {
		t.Error("sem código não há alvo a dispensar")
	}
}

// A MENSAGEM DE COMMIT é a forma preferida de dispensar: ela fica no histórico, ao lado
// do porquê da mudança. A variável de ambiente some junto com o shell — quem ler o commit
// meses depois vê um gate que não rodou, sem saber por quê nem quem decidiu.
func TestDispensaDaMensagemDeCommit(t *testing.T) {
	t.Run("RLUEX-B15: The commit message declares waivers that survive in the history", func(t *testing.T) {})

	msg := `feat(plano): libera o plano 0007

As specs nascem antes do código, como sempre na primeira rodada.

[skip-triad-complete@NOVOA: spec nova do plano 0007]
[skip-triad-complete@NOVOB: spec nova do plano 0007]`

	d, erros := WaiverFromMessage(msg)
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}
	id := RuleID("triad-complete")
	for _, cod := range []string{"NOVOA", "NOVOB"} {
		if motivo, ok := d.WaivedTarget(id, cod); !ok || motivo != "spec nova do plano 0007" {
			t.Errorf("%s deveria estar dispensado com o motivo escrito, veio %q/%v", cod, motivo, ok)
		}
	}
	// E o resto continua confrontado — é o mesmo ponto da dispensa por alvo.
	if _, ok := d.WaivedTarget(id, "QBRDA"); ok {
		t.Error("um código não nomeado não pode ser dispensado")
	}
}

// Sem o motivo o marcador é recusado: é a mesma garantia da forma por variável, e
// aceitá-lo vazio faria o relatório dizer "waived" sem dizer por quê.
func TestMarcadorSemMotivoEhRecusado(t *testing.T) {
	t.Run("RLUEX-B16: A commit marker whose reason is blank is refused", func(t *testing.T) {})

	if _, erros := WaiverFromMessage("fix: algo\n\n[skip-triad-complete@WRKSP: ]"); len(erros) == 0 {
		t.Error("marcador sem motivo deveria ser recusado")
	}
}

// A MESMA garantia nas DUAS formas. Uma valendo e a outra não daria uma porta lateral:
// quem quisesse fugir do gate usaria a forma permissiva.
func TestAsDuasFormasExigemMotivoEscrito(t *testing.T) {
	t.Run("RLUEX-I02: Every accepted waiver carries a written reason", func(t *testing.T) {})

	dTexto, errosTexto := ParseWaiver("triad-complete=")
	if len(errosTexto) == 0 {
		t.Error("a forma textual com motivo vazio tem de ser recusada")
	}
	if _, ok := dTexto.Waived(RuleID("triad-complete")); ok {
		t.Error("a entrada recusada não pode ter sido guardada mesmo assim")
	}

	dMsg, errosMsg := WaiverFromMessage("fix: algo\n\n[skip-triad-complete:    ]")
	if len(errosMsg) == 0 {
		t.Error("o marcador com motivo em branco tem de ser recusado")
	}
	if _, ok := dMsg.Waived(RuleID("triad-complete")); ok {
		t.Error("o marcador recusado não pode ter sido guardado mesmo assim")
	}
}

// Toda recusa CHEGA a quem chamou. Uma entrada malformada aceita em silêncio produziria
// uma dispensa que não dispensa — e o commit reprovaria sem explicação aparente.
func TestCadaEntradaMalformadaProduzUmErro(t *testing.T) {
	t.Run("RLUEX-I03: A refusal always reaches the caller as an error", func(t *testing.T) {})

	_, erros := ParseWaiver("sem-motivo,=sem regra,triad-complete@=sem código," +
		"triad-complete@pasta/arquivo.md=caminho")
	if len(erros) != 4 {
		t.Fatalf("cada entrada malformada tem de produzir o seu erro, veio %d: %v", len(erros), erros)
	}
}

// O marcador SEM código dispensa a regra inteira — a saída grossa continua existindo,
// para o gate recém-declarado que o projeto ainda não cumpre em lugar nenhum.
func TestMarcadorSemCodigoValeParaTudo(t *testing.T) {
	d, _ := WaiverFromMessage("chore: liga o gate\n\n[skip-header-valid: nenhum arquivo tem header ainda]")
	if _, ok := d.Waived(RuleID("header-valid")); !ok {
		t.Error("sem código, o marcador vale para a regra inteira")
	}
}

// Cada alvo carrega o SEU motivo. `PorRegra` guarda um motivo por regra, e duas
// dispensas da mesma regra faziam a segunda sobrescrever a primeira — o relatório
// mostrava o mesmo motivo para os dois alvos, e deixava de dizer a verdade sobre um.
func TestMotivoEhPorAlvo(t *testing.T) {
	t.Run("RLUEX-B14: Each target carries its own reason", func(t *testing.T) {})

	d, _ := WaiverFromMessage(`chore: libera duas

[skip-triad-complete@FRMTT: spec nova, é o card #6]
[skip-triad-complete@TSHRT: spec nova, é o card #8]`)

	id := RuleID("triad-complete")
	if m, _ := d.WaivedTarget(id, "FRMTT"); m != "spec nova, é o card #6" {
		t.Errorf("FRMTT deveria trazer o motivo dele, veio %q", m)
	}
	if m, _ := d.WaivedTarget(id, "TSHRT"); m != "spec nova, é o card #8" {
		t.Errorf("TSHRT deveria trazer o motivo dele, veio %q", m)
	}
}

// A dispensa da mensagem de commit e a da variável de ambiente CONVIVEM: um projeto pode
// ter um hook de CI que usa a variável e um autor que escreve o marcador, e recusar a
// combinação obrigaria a escolher sem motivo.
func TestMergeJuntaAsDuasOrigens(t *testing.T) {
	t.Run("RLUEX-B17: Two waivers merge instead of forcing a choice", func(t *testing.T) {})

	daMensagem, _ := WaiverFromMessage("feat: x\n\n[skip-triad-complete@NOVOA: spec nova]")
	daVariavel, _ := ParseWaiver("header-valid=projeto em bootstrap")

	junta := daMensagem.Merge(daVariavel)

	if _, ok := junta.WaivedTarget(RuleID("triad-complete"), "NOVOA"); !ok {
		t.Error("a dispensa da mensagem tem de sobreviver ao merge, com o alvo dela")
	}
	if _, ok := junta.Waived(RuleID("header-valid")); !ok {
		t.Error("a dispensa da variável tem de sobreviver ao merge")
	}
	// E o alvo não nomeado segue confrontado: juntar não afrouxa nenhuma das duas.
	if _, ok := junta.WaivedTarget(RuleID("triad-complete"), "QBRDA"); ok {
		t.Error("juntar duas dispensas não pode afrouxar a restrição por alvo de nenhuma")
	}
}

// A dispensa só sabe dizer SE dispensou — ela não emite veredito. Misturar as duas coisas
// poria a saída de emergência dentro da régua.
func TestDispensaNaoEmiteVeredito(t *testing.T) {
	t.Run("RLUEX-X02: The unit does not decide whether a rule passes", func(t *testing.T) {})

	d, _ := ParseWaiver("header-valid=projeto em bootstrap")
	motivo, ok := d.Waived(RuleID("triad-complete"))
	if ok {
		t.Error("uma regra que ninguém dispensou não pode vir dispensada")
	}
	if motivo != "" {
		t.Errorf("a resposta é o motivo da dispensa, nunca um veredito: %q", motivo)
	}
}

// A dispensa é TEXTO PURO. Construída duas vezes do mesmo texto, em diretórios de
// trabalho diferentes, ela responde igual — porque não olha o disco. Se olhasse, a mesma
// dispensa significaria coisas diferentes em dois checkouts.
func TestDispensaNaoDependeDoDisco(t *testing.T) {
	t.Run("RLUEX-X03: The unit reads neither files nor the map", func(t *testing.T) {})

	const bruto = "triad-complete@WRKSP=spec nova do plano 0007"

	aqui, _ := ParseWaiver(bruto)
	t.Chdir(t.TempDir())
	ali, _ := ParseWaiver(bruto)

	id := RuleID("triad-complete")
	mAqui, okAqui := aqui.WaivedTarget(id, "WRKSP")
	mAli, okAli := ali.WaivedTarget(id, "WRKSP")
	if okAqui != okAli || mAqui != mAli {
		t.Errorf("a mesma dispensa respondeu diferente em dois diretórios: %v/%q vs %v/%q",
			okAqui, mAqui, okAli, mAli)
	}
	if !okAqui {
		t.Error("o alvo nomeado deveria estar dispensado nas duas leituras")
	}
}
