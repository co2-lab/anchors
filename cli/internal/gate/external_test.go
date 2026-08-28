package gate

import (
	"fmt"
	"strings"
	"testing"
)

// TestLoteNaoEstouraALinhaDeComando guarda o caso que fazia todo `scope: batch`
// reprovar no Windows: os alvos iam num argv só, o CreateProcess cortava em 32767
// chars, o exec falhava ANTES de rodar e — como a saída voltava vazia — o gate
// reprovava mudo, indistinguível de violação real. Aqui os alvos somam ~85 KB, o
// tamanho real do app de referência.
func TestLoteNaoEstouraALinhaDeComando(t *testing.T) {
	var alvos []string
	for i := 0; i < 1500; i++ {
		alvos = append(alvos, fmt.Sprintf("apps/mobile/src/components/atoms/Componente%04d.tsx", i))
	}

	const teto = 24000
	lotes := fatiarAlvos(alvos, teto)
	if len(lotes) < 2 {
		t.Fatalf("85 KB de alvos precisam ser fatiados; veio %d lote(s)", len(lotes))
	}

	var vistos int
	for i, lote := range lotes {
		tam := 0
		for _, a := range lote {
			tam += len(a) + 1
		}
		if tam > teto {
			t.Errorf("lote %d tem %d bytes, acima do teto de %d", i, tam, teto)
		}
		if len(lote) == 0 {
			t.Errorf("lote %d está vazio — lote vazio vira execução sem alvo, que é outro gate", i)
		}
		vistos += len(lote)
	}
	// Fatiar não pode PERDER alvo: um arquivo que some do lote é um arquivo que
	// ninguém confronta, e o gate passa certificando o que não olhou.
	if vistos != len(alvos) {
		t.Errorf("os lotes cobrem %d alvos; deviam cobrir os %d", vistos, len(alvos))
	}
}

// TestSemAlvosRodaUmaVez fixa o `scope: project`: a ferramenta olha o projeto inteiro
// e não recebe alvo. Devolver lote nenhum faria o gate não rodar — e um gate que não
// roda passa por omissão, que é o oposto do que ele existe para fazer.
func TestSemAlvosRodaUmaVez(t *testing.T) {
	lotes := fatiarAlvos(nil, 24000)
	if len(lotes) != 1 {
		t.Fatalf("sem alvos deve haver exatamente 1 execução; veio %d", len(lotes))
	}
	if len(lotes[0]) != 0 {
		t.Errorf("a execução de projeto não leva alvo; veio %v", lotes[0])
	}
}

// TestLoteCabendoNaoEhFatiado protege o caminho comum (e o macOS, onde o ARG_MAX é
// de centenas de KB): onde já cabia numa execução, continua sendo uma só.
func TestLoteCabendoNaoEhFatiado(t *testing.T) {
	alvos := []string{"a.ts", "b.ts", "c.ts"}
	lotes := fatiarAlvos(alvos, 24000)
	if len(lotes) != 1 {
		t.Fatalf("3 alvos curtos cabem numa execução; veio %d lotes", len(lotes))
	}
	if len(lotes[0]) != 3 {
		t.Errorf("a execução única leva os 3 alvos; veio %v", lotes[0])
	}
}

// TestAlvoMaiorQueOTetoVaiSozinho — não há como partir um caminho ao meio. Ele vai
// só no lote e o SO recusa com a mensagem dele; o que não pode é sumir em silêncio.
func TestAlvoMaiorQueOTetoVaiSozinho(t *testing.T) {
	gigante := strings.Repeat("x", 200)
	lotes := fatiarAlvos([]string{"a.ts", gigante, "b.ts"}, 100)

	var vistos int
	for _, l := range lotes {
		vistos += len(l)
	}
	if vistos != 3 {
		t.Errorf("nenhum alvo pode sumir na fatia; cobertos %d de 3", vistos)
	}
}

// TestReprovaSemSaidaDizPorQue: um gate que reprova sem imprimir nada deixava o laudo
// vazio, e o operador lia "violação de código" onde havia problema de ambiente (sem
// `sh` no PATH, linha de comando estourada). O laudo tem que dizer o que houve.
func TestReprovaSemSaidaDizPorQue(t *testing.T) {
	v, detalhe := RunExternalArgs("exit 1", []string{"a.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("exit 1 reprova; veio %v", v)
	}
	if strings.TrimSpace(detalhe) == "" {
		t.Error("reprovação sem saída precisa dizer o motivo, não deixar o laudo em branco")
	}
}
