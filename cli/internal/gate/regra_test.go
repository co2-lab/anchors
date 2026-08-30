package gate

import (
	"strings"
	"testing"
)

// O motivo é a única coisa que separa dispensa DELIBERADA de gate ignorado. Aceitá-lo
// ausente esvaziaria a garantia — e o relatório passaria a mostrar "dispensado" sem
// dizer por quê.
func TestDispensaExigeMotivo(t *testing.T) {
	_, erros := ParseDispensa("trinca-completa")

	if len(erros) != 1 {
		t.Fatalf("dispensa sem motivo tem de ser recusada, veio %d erro(s)", len(erros))
	}
	if _, erros := ParseDispensa("trinca-completa="); len(erros) != 1 {
		t.Error("motivo vazio é o mesmo que ausente")
	}
}

// Dispensar o GATE cobre todas as regras dele; dispensar a REGRA preserva o resto. As
// duas granularidades existem porque quem não conhece as regras precisa da saída grossa.
func TestDispensaAceitaAsDuasGranularidades(t *testing.T) {
	d, erros := ParseDispensa("spec-completa/sem-placeholder=a spec nasce em rascunho")
	if len(erros) != 0 {
		t.Fatalf("erros inesperados: %v", erros)
	}

	if _, ok := d.Dispensou("spec-completa/sem-placeholder"); !ok {
		t.Error("a regra dispensada não foi reconhecida")
	}
	// A OUTRA regra do mesmo gate continua valendo — é o ponto de dispensar por regra.
	if _, ok := d.Dispensou("spec-completa/tem-regra-catalogada"); ok {
		t.Error("dispensar uma regra não pode desligar as demais do mesmo gate")
	}

	// E dispensar o gate inteiro cobre as regras dele.
	dg, _ := ParseDispensa("spec-completa=projeto em bootstrap")
	if _, ok := dg.Dispensou("spec-completa/sem-placeholder"); !ok {
		t.Error("dispensar o gate tem de cobrir suas regras")
	}
}

// Um gate que não foi dispensado precisa continuar rodando — o erro que mais custaria
// aqui é uma dispensa vazando para o que ninguém pediu.
func TestDispensaNaoAlcancaOQueNaoFoiPedido(t *testing.T) {
	d, _ := ParseDispensa("trinca-completa=a feature ainda é um card")

	if _, ok := d.Dispensou("guide-checklist"); ok {
		t.Error("a dispensa alcançou um gate que ninguém dispensou")
	}
	if _, ok := (Dispensa{}).Dispensou("qualquer-coisa"); ok {
		t.Error("dispensa vazia não dispensa nada")
	}
}

func TestRegraIDSeparaGateDeRegra(t *testing.T) {
	id := NovaRegraID("spec-completa", "sem-placeholder")
	if id != "spec-completa/sem-placeholder" {
		t.Errorf("ID montado errado: %q", id)
	}
	if id.Gate() != "spec-completa" || id.Regra() != "sem-placeholder" {
		t.Errorf("decomposição errada: gate=%q regra=%q", id.Gate(), id.Regra())
	}
	// Gate com uma verificação só: o nome já a identifica, e não há barra a inventar.
	if s := NovaRegraID("layer-boundary", ""); s != "layer-boundary" {
		t.Errorf("gate sem regra não deve ganhar barra: %q", s)
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
	d, erros := ParseDispensa(
		"trinca-completa@NOVOA=spec nova do plano 0007," +
			"trinca-completa@NOVOB=spec nova do plano 0007")
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}

	id := RegraID("trinca-completa")

	// Os alvos NOMEADOS estão dispensados.
	for _, cod := range []string{"NOVOA", "NOVOB"} {
		if motivo, ok := d.DispensouAlvo(id, cod); !ok || motivo == "" {
			t.Errorf("%s deveria estar dispensado com motivo", cod)
		}
	}

	// E O RESTO CONTINUA SENDO CONFRONTADO. É o ponto inteiro do recurso: a trinca que
	// quebrou por descuido noutra unidade não pode passar de carona.
	if _, ok := d.DispensouAlvo(id, "QBRDA"); ok {
		t.Error("um alvo não nomeado NÃO pode ser dispensado — é o mascaramento que este " +
			"recurso existe para impedir")
	}

	// E a pergunta sem alvo responde "não dispensado", para que o filtro de gates não
	// remova o gate da lista: ele precisa RODAR para confrontar os outros alvos.
	if _, ok := d.Dispensou(id); ok {
		t.Error("dispensa COM alvos não pode valer para o gate inteiro — sairia da lista " +
			"e não confrontaria ninguém")
	}
}

// A dispensa SEM alvo continua valendo para tudo: há casos legítimos, como um gate
// recém-declarado que o projeto ainda não cumpre em lugar nenhum.
func TestDispensaSemAlvoValeParaTudo(t *testing.T) {
	d, _ := ParseDispensa("trinca-completa=gate novo, nenhuma unidade o cumpre ainda")
	id := RegraID("trinca-completa")
	if _, ok := d.Dispensou(id); !ok {
		t.Error("sem alvo declarado, a dispensa vale para o gate inteiro")
	}
	if _, ok := d.DispensouAlvo(id, "QUALQ"); !ok {
		t.Error("sem alvo declarado, qualquer caminho está dispensado")
	}
}

// `@` sem caminho é engano de digitação, e aceitá-lo em silêncio produziria uma dispensa
// que não dispensa nada — o commit reprovaria sem explicação aparente.
func TestDispensaAlvoVazioEhRecusada(t *testing.T) {
	_, erros := ParseDispensa("trinca-completa@=motivo qualquer")
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
	for _, bruto := range []string{
		"trinca-completa@packages/shared/Workspace.spec.md=motivo",
		"trinca-completa@packages/*=motivo",
		"trinca-completa@arquivo.spec.md=motivo",
	} {
		_, erros := ParseDispensa(bruto)
		if len(erros) == 0 {
			t.Errorf("%q deveria ser recusado: o alvo é o CÓDIGO, não o caminho", bruto)
			continue
		}
		if !strings.Contains(erros[0], "CÓDIGO") {
			t.Errorf("o erro deveria dizer o que usar no lugar: %s", erros[0])
		}
	}
}

// Um alvo SEM CÓDIGO não é alcançado por uma dispensa restrita. Artefato sem identidade
// é um problema anterior — quem cobra isso é o `codigo-catalogado`, e dar uma saída
// lateral aqui esconderia a causa.
func TestDispensaPorAlvoNaoAlcancaQuemNaoTemCodigo(t *testing.T) {
	d, _ := ParseDispensa("trinca-completa@WRKSP=spec nova")
	if _, ok := d.DispensouAlvo(RegraID("trinca-completa"), ""); ok {
		t.Error("sem código não há alvo a dispensar")
	}
}

// A MENSAGEM DE COMMIT é a forma preferida de dispensar: ela fica no histórico, ao lado
// do porquê da mudança. A variável de ambiente some junto com o shell — quem ler o commit
// meses depois vê um gate que não rodou, sem saber por quê nem quem decidiu.
func TestDispensaDaMensagemDeCommit(t *testing.T) {
	msg := `feat(plano): libera o plano 0007

As specs nascem antes do código, como sempre na primeira rodada.

[skip-trinca-completa@NOVOA: spec nova do plano 0007]
[skip-trinca-completa@NOVOB: spec nova do plano 0007]`

	d, erros := DispensaDaMensagem(msg)
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}
	id := RegraID("trinca-completa")
	for _, cod := range []string{"NOVOA", "NOVOB"} {
		if motivo, ok := d.DispensouAlvo(id, cod); !ok || motivo != "spec nova do plano 0007" {
			t.Errorf("%s deveria estar dispensado com o motivo escrito, veio %q/%v", cod, motivo, ok)
		}
	}
	// E o resto continua confrontado — é o mesmo ponto da dispensa por alvo.
	if _, ok := d.DispensouAlvo(id, "QBRDA"); ok {
		t.Error("um código não nomeado não pode ser dispensado")
	}
}

// Sem o motivo o marcador é recusado: é a mesma garantia da forma por variável, e
// aceitá-lo vazio faria o relatório dizer "dispensado" sem dizer por quê.
func TestMarcadorSemMotivoEhRecusado(t *testing.T) {
	if _, erros := DispensaDaMensagem("fix: algo\n\n[skip-trinca-completa@WRKSP: ]"); len(erros) == 0 {
		t.Error("marcador sem motivo deveria ser recusado")
	}
}

// O marcador SEM código dispensa a regra inteira — a saída grossa continua existindo,
// para o gate recém-declarado que o projeto ainda não cumpre em lugar nenhum.
func TestMarcadorSemCodigoValeParaTudo(t *testing.T) {
	d, _ := DispensaDaMensagem("chore: liga o gate\n\n[skip-header-conforme: nenhum arquivo tem header ainda]")
	if _, ok := d.Dispensou(RegraID("header-conforme")); !ok {
		t.Error("sem código, o marcador vale para a regra inteira")
	}
}

// Cada alvo carrega o SEU motivo. `PorRegra` guarda um motivo por regra, e duas
// dispensas da mesma regra faziam a segunda sobrescrever a primeira — o relatório
// mostrava o mesmo motivo para os dois alvos, e deixava de dizer a verdade sobre um.
func TestMotivoEhPorAlvo(t *testing.T) {
	d, _ := DispensaDaMensagem(`chore: libera duas

[skip-trinca-completa@FRMTT: spec nova, é o card #6]
[skip-trinca-completa@TSHRT: spec nova, é o card #8]`)

	id := RegraID("trinca-completa")
	if m, _ := d.DispensouAlvo(id, "FRMTT"); m != "spec nova, é o card #6" {
		t.Errorf("FRMTT deveria trazer o motivo dele, veio %q", m)
	}
	if m, _ := d.DispensouAlvo(id, "TSHRT"); m != "spec nova, é o card #8" {
		t.Errorf("TSHRT deveria trazer o motivo dele, veio %q", m)
	}
}
