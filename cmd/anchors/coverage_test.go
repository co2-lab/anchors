package main

import (
	"os"
	"path/filepath"
	"testing"
)

// A spec CITA regras vizinhas em prosa, e isso e' estilo do projeto -- nao defeito.
// O defeito e' o `codesInFile` tratar a citacao como DECLARACAO.
//
// MEDIDO no blue-eyes: 37 dos 55 nos com `proven_codes` carregavam codigo de outra
// unidade. O `InfraList.spec.md` cita `QSCOP-B02` UMA vez, em prosa ("o efetivo
// (`QSCOP-B02`), nao o pedido"), nenhum teste do InfraList o menciona -- e mesmo
// assim o mapa afirmava que o InfraList o provou.
//
// O estrago e' pior que uma lacuna: a lacuna aparece no relatorio, e isto some.
// Quem lesse o mapa veria a regra do vizinho como provada por um teste que nunca a
// tocou, e o `stale` nao cobraria a unidade que de fato a deve.
//
// A causa e' de endereco: `CodesInCase` foi escrita para o NOME de um caso de teste
// (`"SPCRX-V01: ..."`), onde todo codigo presente E' o codigo do caso. Aplicada ao
// arquivo INTEIRO da spec, ela colhe tambem o que a prosa menciona.
func TestCodigoCitadoEmProsaNaoContaComoDeclarado(t *testing.T) {
	dir := t.TempDir()
	spec := filepath.Join(dir, "InfraList.spec.md")
	conteudo := `# InfraList

## INLSN-B01 — a lista ordena por gravidade

A tela usa o efetivo (` + "`QSCOP-B02`" + `), nao o pedido.

## INLSN-B02 — o veredito vem da integracao
`
	if err := os.WriteFile(spec, []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}

	codigos, err := codesInFileOfUnit(spec, "INLSN")
	if err != nil {
		t.Fatal(err)
	}

	temB01, temB02, temAlheio := false, false, ""
	for _, c := range codigos {
		switch c {
		case "INLSN-B01":
			temB01 = true
		case "INLSN-B02":
			temB02 = true
		default:
			temAlheio = c
		}
	}

	if !temB01 || !temB02 {
		t.Errorf("as regras da PROPRIA unidade sumiram: %v — o filtro cortou demais", codigos)
	}
	if temAlheio != "" {
		t.Errorf("`%s` e' citado em PROSA e entrou como declarado — o mapa vai afirmar "+
			"que esta unidade provou a regra do vizinho, e nenhum teste dela a toca", temAlheio)
	}
}

// Sem o codigo da unidade nao ha o que filtrar, e cortar tudo seria pior que nao
// filtrar: o no perderia os proprios cenarios. O comportamento aqui e' passar reto.
func TestSemCodigoDaUnidadeNaoFiltra(t *testing.T) {
	dir := t.TempDir()
	spec := filepath.Join(dir, "x.spec.md")
	if err := os.WriteFile(spec, []byte("## ABCDE-B01\n## FGHIJ-B02\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	codigos, err := codesInFileOfUnit(spec, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(codigos) != 2 {
		t.Errorf("sem codigo da unidade o filtro deveria passar reto, veio %v", codigos)
	}
}
