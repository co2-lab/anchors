package gate

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestExternal_PassExitZero(t *testing.T) {
	t.Run("EXCMX-B01: An external command exiting with status zero returns Pass", func(t *testing.T) {})
	v, detail := RunExternalArgs("exit 0", []string{"a.ts"}, t.TempDir())
	if v != Pass {
		t.Fatalf("expected Pass, got %v", v)
	}
	if detail != "" {
		t.Fatalf("expected empty detail, got %q", detail)
	}
}

func TestExternal_FailExitNonZeroWithOutput(t *testing.T) {
	t.Run("EXCMX-B02: An external command exiting with non-zero status returns Fail with output", func(t *testing.T) {})
	v, detail := RunExternalArgs("echo 'failure message' && exit 1", []string{"a.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "failure message") {
		t.Fatalf("expected detail to contain failure message, got %q", detail)
	}
}

// TestReprovaSemSaidaDizPorQue: um gate que reprova sem imprimir nada deixava o laudo
// vazio, e o operador lia "violação de código" onde havia problema de ambiente (sem
// `sh` no PATH, linha de comando estourada). O laudo tem que dizer o que houve.
func TestReprovaSemSaidaDizPorQue(t *testing.T) {
	t.Run("EXCMX-B03: An external command failing with empty output reports execution error", func(t *testing.T) {})
	v, detalhe := RunExternalArgs("exit 1", []string{"a.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("exit 1 reprova; veio %v", v)
	}
	if strings.TrimSpace(detalhe) == "" {
		t.Error("reprovação sem saída precisa dizer o motivo, não deixar o laudo em branco")
	}
}

func TestExternal_RunExternalDelegatesWithNodeID(t *testing.T) {
	t.Run("EXCMX-B04: Single node execution delegates to RunExternalArgs", func(t *testing.T) {})
	node := mapx.Node{ID: "src/foo.ts"}
	v, detail := runExternal("echo $1 && exit 1", node, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "src/foo.ts") {
		t.Fatalf("expected node ID in output, got %q", detail)
	}
}

func TestExternal_PlaceholderFileRewritten(t *testing.T) {
	t.Run("EXCMX-B05: Placeholder file is rewritten to positional parameter", func(t *testing.T) {})
	v, detail := RunExternalArgs("echo {{file}} && exit 1", []string{"file1.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "file1.ts") {
		t.Fatalf("expected output to contain file1.ts, got %q", detail)
	}
}

func TestExternal_PlaceholderFilesRewritten(t *testing.T) {
	t.Run("EXCMX-B06: Placeholder files is rewritten to all positional parameters", func(t *testing.T) {})
	v, detail := RunExternalArgs("echo {{files}} && exit 1", []string{"a.ts", "b.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "a.ts b.ts") {
		t.Fatalf("expected output to contain both files, got %q", detail)
	}
}

// TestSemAlvosRodaUmaVez fixa o `scope: project`: a ferramenta olha o projeto inteiro
// e não recebe alvo. Devolver lote nenhum faria o gate não rodar — e um gate que não
// roda passa por omissão, que é o oposto do que ele existe para fazer.
func TestSemAlvosRodaUmaVez(t *testing.T) {
	t.Run("EXCMX-B07: Execution without targets runs once for project scope", func(t *testing.T) {})
	lotes := sliceTargets(nil, 24000)
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
	t.Run("EXCMX-B08: Targets within budget run in a single batch", func(t *testing.T) {})
	alvos := []string{"a.ts", "b.ts", "c.ts"}
	lotes := sliceTargets(alvos, 24000)
	if len(lotes) != 1 {
		t.Fatalf("3 alvos curtos cabem numa execução; veio %d lotes", len(lotes))
	}
	if len(lotes[0]) != 3 {
		t.Errorf("a execução única leva os 3 alvos; veio %v", lotes[0])
	}
}

// TestLoteNaoEstouraALinhaDeComando guarda o caso que fazia todo `scope: batch`
// reprovar no Windows: os alvos iam num argv só, o CreateProcess cortava em 32767
// chars, o exec falhava ANTES de rodar e — como a saída voltava vazia — o gate
// reprovava mudo, indistinguível de violação real. Aqui os alvos somam ~85 KB, o
// tamanho real do app de referência.
func TestLoteNaoEstouraALinhaDeComando(t *testing.T) {
	t.Run("EXCMX-B09: Targets exceeding budget are partitioned across batches", func(t *testing.T) {})
	var alvos []string
	for i := 0; i < 1500; i++ {
		alvos = append(alvos, fmt.Sprintf("apps/mobile/src/components/atoms/Componente%04d.tsx", i))
	}

	const teto = 24000
	lotes := sliceTargets(alvos, teto)
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

func TestExternal_BatchFailureCombinesOutput(t *testing.T) {
	t.Run("EXCMX-B10: Failure in any batch causes entire execution to fail", func(t *testing.T) {})
	t.Setenv("ANCHORS_ARGV_MAX", "20")
	v, detail := RunExternalArgs("if [ \"$1\" = \"bad.ts\" ]; then echo 'bad file failed'; exit 1; fi", []string{"good.ts", "bad.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "bad file failed") {
		t.Fatalf("expected detail to contain batch failure, got %q", detail)
	}
}

func TestExternal_SingleTargetOutputTruncatedAt500(t *testing.T) {
	t.Run("EXCMX-B11: Single target failure detail is truncated at five hundred characters", func(t *testing.T) {})
	longOutput := strings.Repeat("e", 600)
	v, detail := RunExternalArgs(fmt.Sprintf("echo '%s' && exit 1", longOutput), []string{"a.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "… (saída truncada)") {
		t.Fatalf("expected truncation marker, got %q", detail)
	}
	prefix := strings.Split(detail, "\n")[0]
	if len(prefix) != 500 {
		t.Fatalf("expected truncated prefix length 500, got %d", len(prefix))
	}
}

func TestExternal_BatchOutputTruncatedAt4000(t *testing.T) {
	t.Run("EXCMX-B12: Batch failure detail is truncated at four thousand characters", func(t *testing.T) {})
	longOutput := strings.Repeat("b", 4500)
	v, detail := RunExternalArgs(fmt.Sprintf("echo '%s' && exit 1", longOutput), []string{"a.ts", "b.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "… (saída truncada)") {
		t.Fatalf("expected truncation marker, got %q", detail)
	}
	prefix := strings.Split(detail, "\n")[0]
	if len(prefix) != 4000 {
		t.Fatalf("expected truncated prefix length 4000, got %d", len(prefix))
	}
}

func TestExternal_ArgvLimitEnvOverride(t *testing.T) {
	t.Run("EXCMX-B13: Environment variable overrides argv limit", func(t *testing.T) {})
	t.Setenv("ANCHORS_ARGV_MAX", "12345")
	if limit := argvLimit(); limit != 12345 {
		t.Fatalf("expected 12345, got %d", limit)
	}
}

func TestExternal_TargetPathsPassedStrictlyInArgv(t *testing.T) {
	t.Run("EXCMX-I01: Target paths are passed strictly in argv preventing injection", func(t *testing.T) {})
	// Filename with shell metacharacters
	malicious := "foo;touch INJECTED;echo.ts"
	v, detail := RunExternalArgs("echo \"target is: $1\"", []string{malicious}, t.TempDir())
	if v != Pass {
		t.Fatalf("expected Pass, got %v (%s)", v, detail)
	}
}

// TestAlvoMaiorQueOTetoVaiSozinho — não há como partir um caminho ao meio. Ele vai
// só no lote e o SO recusa com a mensagem dele; o que não pode é sumir em silêncio.
func TestAlvoMaiorQueOTetoVaiSozinho(t *testing.T) {
	t.Run("EXCMX-I02: Oversized single target is isolated in its own batch", func(t *testing.T) {})
	gigante := strings.Repeat("x", 200)
	lotes := sliceTargets([]string{"a.ts", gigante, "b.ts"}, 100)

	var vistos int
	for _, l := range lotes {
		vistos += len(l)
	}
	if vistos != 3 {
		t.Errorf("nenhum alvo pode sumir na fatia; cobertos %d de 3", vistos)
	}
}

func TestExternal_PlatformDefaultArgvLimit(t *testing.T) {
	t.Run("EXCMX-I03: Platform default argv limits are enforced", func(t *testing.T) {})
	t.Setenv("ANCHORS_ARGV_MAX", "")
	limit := argvLimit()
	if runtime.GOOS == "windows" {
		if limit != 6000 {
			t.Fatalf("expected 6000 on windows, got %d", limit)
		}
	} else {
		if limit != 100000 {
			t.Fatalf("expected 100000 on unix, got %d", limit)
		}
	}
}

func TestExternal_DoesNotParseDiagnostics(t *testing.T) {
	t.Run("EXCMX-X01: The gate does not parse or interpret linter diagnostics", func(t *testing.T) {})
	v, detail := RunExternalArgs("echo 'src/a.ts:1:1: error: something wrong' && exit 1", []string{"src/a.ts"}, t.TempDir())
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(detail, "src/a.ts:1:1: error: something wrong") {
		t.Fatalf("expected raw output preserved, got %q", detail)
	}
}

func TestExternal_DoesNotAggregateCrossFileState(t *testing.T) {
	t.Run("EXCMX-X02: The gate does not aggregate cross-file state across partitioned batches", func(t *testing.T) {})
	t.Setenv("ANCHORS_ARGV_MAX", "20")
	// Each batch runs in a separate subshell, so environment changes in one do not persist to another
	v, detail := RunExternalArgs("if [ -z \"$SEEN\" ]; then export SEEN=1; exit 0; else exit 1; fi", []string{"a.ts", "b.ts"}, t.TempDir())
	if v != Pass {
		t.Fatalf("expected Pass because each batch runs in its own shell, got %v (%s)", v, detail)
	}
}
