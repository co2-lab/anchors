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
