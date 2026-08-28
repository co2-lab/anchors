package testsig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseUnifiedDiff(t *testing.T) {
	diff := `diff --git a/src/A.tsx b/src/A.tsx
--- a/src/A.tsx
+++ b/src/A.tsx
@@ -10,0 +11,2 @@
+const x = 1
+const y = 2
@@ -20,1 +22,1 @@
-old line
+new line`
	changed := parseUnifiedDiff(diff)
	lines := changed["src/A.tsx"]
	if lines == nil {
		t.Fatal("deveria ter linhas mudadas em src/A.tsx")
	}
	for _, want := range []int{11, 12, 22} {
		if !lines[want] {
			t.Errorf("linha %d deveria estar marcada como mudada; veio %v", want, lines)
		}
	}
	if lines[21] {
		t.Error("linha 21 não foi tocada")
	}
}

func TestUncoveredIn(t *testing.T) {
	fc := FileCoverage{Lines: map[int]bool{10: true, 11: false, 12: true}}
	changed := map[int]bool{10: true, 11: true, 13: true} // 13 não é instrumentada
	un := fc.UncoveredIn(changed)
	if len(un) != 1 || un[0] != 11 {
		t.Fatalf("esperava [11] descoberta, veio %v", un)
	}
	if fc.InstrumentedIn(changed) != 2 { // 10 e 11 (13 não é instrumentada)
		t.Errorf("esperava 2 instrumentadas no diff, veio %d", fc.InstrumentedIn(changed))
	}
}

func TestParseDiffFileDeletedFile(t *testing.T) {
	// arquivo deletado (+++ /dev/null) não deve gerar entrada
	diff := "--- a/gone.ts\n+++ /dev/null\n@@ -1,2 +0,0 @@\n-a\n-b\n"
	p := filepath.Join(t.TempDir(), "d.diff")
	os.WriteFile(p, []byte(diff), 0o644)
	changed, _ := ParseDiffFile(p)
	if len(changed) != 0 {
		t.Errorf("arquivo deletado não deveria ter linhas mudadas, veio %v", changed)
	}
}
