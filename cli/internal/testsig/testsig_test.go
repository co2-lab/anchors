package testsig

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseJUnit(t *testing.T) {
	xml := `<?xml version="1.0"?>
<testsuites>
  <testsuite name="Spacer" file="src/Spacer.test.tsx">
    <testcase name="SPCRX-V01: eixo"/>
    <testcase name="SPCRX-V02: horizontal"><failure message="x"/></testcase>
    <testcase name="SPCRX-X01: sem conteudo"><skipped/></testcase>
  </testsuite>
</testsuites>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Cases) != 3 {
		t.Fatalf("esperava 3 casos, veio %d", len(rep.Cases))
	}
	proven := rep.PassedCodes()
	if !proven["SPCRX-V01"] {
		t.Error("SPCRX-V01 passou, deveria estar provado")
	}
	if proven["SPCRX-V02"] {
		t.Error("SPCRX-V02 falhou, NÃO deveria estar provado")
	}
	if proven["SPCRX-X01"] {
		t.Error("SPCRX-X01 foi pulado, NÃO deveria estar provado")
	}
}

func TestParseJUnitSingleSuite(t *testing.T) {
	// alguns reporters emitem <testsuite> na raiz
	xml := `<testsuite name="X" file="a.test.ts"><testcase name="AB01X-S01: ok"/></testsuite>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil || len(rep.Cases) != 1 {
		t.Fatalf("esperava 1 caso, veio %d err=%v", len(rep.Cases), err)
	}
}

func TestParseLCOV(t *testing.T) {
	lcov := `SF:src/a.ts
DA:1,1
DA:2,0
DA:3,1
end_of_record
SF:src/b.ts
LF:10
LH:8
end_of_record`
	rep, err := ParseLCOV(write(t, "c.info", lcov))
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Files) != 2 {
		t.Fatalf("esperava 2 arquivos, veio %d", len(rep.Files))
	}
	// a.ts: 2 de 3 cobertas (via DA)
	if rep.Files[0].CoveredLines != 2 || rep.Files[0].TotalLines != 3 {
		t.Errorf("a.ts esperava 2/3, veio %d/%d", rep.Files[0].CoveredLines, rep.Files[0].TotalLines)
	}
	// b.ts: 8 de 10 (via LF/LH)
	if rep.Files[1].CoveredLines != 8 || rep.Files[1].TotalLines != 10 {
		t.Errorf("b.ts esperava 8/10, veio %d/%d", rep.Files[1].CoveredLines, rep.Files[1].TotalLines)
	}
	if p := rep.Files[1].Percent(); p != 80 {
		t.Errorf("b.ts esperava 80%%, veio %.0f", p)
	}
}

func TestCodesInCase(t *testing.T) {
	got := CodesInCase("SPCRX-V01: eixo vertical estende SPCRX-X01")
	if len(got) != 2 {
		t.Fatalf("esperava 2 códigos, veio %v", got)
	}
}
