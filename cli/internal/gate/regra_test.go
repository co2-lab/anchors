package gate

import "testing"

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
		"trinca-completa@packages/novo/A.spec.md=spec nova do plano 0007," +
			"trinca-completa@packages/novo/B.spec.md=spec nova do plano 0007")
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}

	id := RegraID("trinca-completa")

	// Os alvos NOMEADOS estão dispensados.
	for _, alvo := range []string{"packages/novo/A.spec.md", "packages/novo/B.spec.md"} {
		if motivo, ok := d.DispensouAlvo(id, alvo, ""); !ok || motivo == "" {
			t.Errorf("%s deveria estar dispensado com motivo", alvo)
		}
	}

	// E O RESTO CONTINUA SENDO CONFRONTADO. É o ponto inteiro do recurso: a trinca que
	// quebrou por descuido em `packages/antigo` não pode passar de carona.
	if _, ok := d.DispensouAlvo(id, "packages/antigo/Quebrada.spec.md", "QBRDA"); ok {
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
	if _, ok := d.DispensouAlvo(id, "qualquer/coisa.spec.md", ""); !ok {
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

// O CÓDIGO como alvo é a forma preferida: ele é a identidade do artefato e sobrevive a
// mover ou renomear o arquivo. Uma dispensa presa ao caminho deixa de valer em silêncio
// quando alguém reorganiza pastas, e o commit seguinte reprova sem explicação.
func TestDispensaAceitaCodigoComoAlvo(t *testing.T) {
	d, erros := ParseDispensa("trinca-completa@WRKSP=spec nova do plano 0007")
	if len(erros) > 0 {
		t.Fatalf("não deveria haver erro: %v", erros)
	}
	id := RegraID("trinca-completa")

	// Casa pelo CÓDIGO, seja qual for o caminho — inclusive um que mudou de lugar.
	if _, ok := d.DispensouAlvo(id, "packages/shared/Workspace.spec.md", "WRKSP"); !ok {
		t.Error("deveria dispensar pelo código")
	}
	if _, ok := d.DispensouAlvo(id, "outro/lugar/Workspace.spec.md", "WRKSP"); !ok {
		t.Error("o código sobrevive a mover o arquivo — é o motivo de preferi-lo")
	}
	// E não alcança outro artefato.
	if _, ok := d.DispensouAlvo(id, "packages/shared/Outra.spec.md", "OUTRA"); ok {
		t.Error("o código dispensa UM artefato, não os vizinhos")
	}
	// Um alvo sem código declarado (um package.json, um config) não é alcançado por
	// engano quando a dispensa é por código.
	if _, ok := d.DispensouAlvo(id, "package.json", ""); ok {
		t.Error("alvo sem código não pode casar uma dispensa por código")
	}
}
