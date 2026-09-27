package testsig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLCOV(t *testing.T) {
	t.Run("LCINL-B01: Each record becomes one file's coverage in report order", func(t *testing.T) {})
	t.Run("LCINL-B02: A line with hits is covered and a line without is not", func(t *testing.T) {})
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
	if len(rep.Files) != 2 || rep.Files[0].File != "src/a.ts" {
		t.Fatalf("expected 2 files, src/a.ts first, got %+v", rep.Files)
	}
	// a.ts: 2 of 3 covered (via DA)
	if rep.Files[0].CoveredLines != 2 || rep.Files[0].TotalLines != 3 {
		t.Errorf("a.ts expected 2/3, got %d/%d", rep.Files[0].CoveredLines, rep.Files[0].TotalLines)
	}
	if l := rep.Files[0].Lines; !l[1] || l[2] || !l[3] || len(l) != 3 {
		t.Errorf("a.ts lines = %v, want 1 and 3 covered, 2 uncovered", l)
	}
	// b.ts: 8 of 10 (via LF/LH)
	if rep.Files[1].CoveredLines != 8 || rep.Files[1].TotalLines != 10 {
		t.Errorf("b.ts expected 8/10, got %d/%d", rep.Files[1].CoveredLines, rep.Files[1].TotalLines)
	}
	if p := rep.Files[1].Percent(); p != 80 {
		t.Errorf("b.ts expected 80%%, got %.0f", p)
	}
}

func TestUncoveredIn(t *testing.T) {
	t.Run("LCINL-B04: The uncovered lines of a change are instrumented and not covered", func(t *testing.T) {})
	t.Run("LCINL-B05: The instrumented count ignores changed lines the report did not instrument", func(t *testing.T) {})
	t.Run("LCINL-I01: Uncovered changed lines never exceed the instrumented changed lines", func(t *testing.T) {})
	fc := FileCoverage{Lines: map[int]bool{10: true, 11: false, 12: true}}
	changed := map[int]bool{10: true, 11: true, 13: true} // 13 is not instrumented
	un := fc.UncoveredIn(changed)
	if len(un) != 1 || un[0] != 11 {
		t.Fatalf("expected [11] uncovered, got %v", un)
	}
	if fc.InstrumentedIn(changed) != 2 { // 10 and 11 (13 is not instrumented)
		t.Errorf("expected 2 instrumented in the diff, got %d", fc.InstrumentedIn(changed))
	}
	if len(un) > fc.InstrumentedIn(changed) {
		t.Error("uncovered lines cannot exceed the instrumented ones")
	}
}

func TestLCOVTotalsAndRecords(t *testing.T) {
	t.Run("LCINL-B03: Stated totals win over the line entries", func(t *testing.T) {
		rep, err := ParseLCOV(write(t, "c.info", "SF:a.ts\nDA:1,1\nDA:2,0\nLF:10\nLH:8\nend_of_record\n"))
		if err != nil {
			t.Fatal(err)
		}
		if f := rep.Files[0]; f.CoveredLines != 8 || f.TotalLines != 10 {
			t.Errorf("want 8/10 from the stated totals, got %d/%d", f.CoveredLines, f.TotalLines)
		}
	})
	t.Run("LCINL-B06: The percentage is covered over total, and zero for a file without lines", func(t *testing.T) {
		if p := (FileCoverage{CoveredLines: 8, TotalLines: 10}).Percent(); p != 80 {
			t.Errorf("want 80, got %v", p)
		}
		if p := (FileCoverage{}).Percent(); p != 0 {
			t.Errorf("a file of no line is 0, got %v", p)
		}
	})
	t.Run("LCINL-B07: A record without its end line is closed by the next record or the end of the report", func(t *testing.T) {
		rep, err := ParseLCOV(write(t, "c.info", "SF:a.ts\nDA:1,1\nSF:b.ts\nDA:1,0\nDA:2,1\n"))
		if err != nil {
			t.Fatal(err)
		}
		if len(rep.Files) != 2 {
			t.Fatalf("want 2 files, got %+v", rep.Files)
		}
		if a, b := rep.Files[0], rep.Files[1]; a.File != "a.ts" || a.TotalLines != 1 || a.CoveredLines != 1 ||
			b.File != "b.ts" || b.TotalLines != 2 || b.CoveredLines != 1 {
			t.Errorf("each record keeps its own lines, got %+v", rep.Files)
		}
	})
	t.Run("LCINL-X01: Branch and function entries do not change the line counts", func(t *testing.T) {
		rep, err := ParseLCOV(write(t, "c.info", "SF:a.ts\nFN:1,f\nFNDA:3,f\nDA:1,1\nDA:2,0\nBRDA:1,0,0,1\nBRF:2\nBRH:1\nend_of_record\n"))
		if err != nil {
			t.Fatal(err)
		}
		if f := rep.Files[0]; f.CoveredLines != 1 || f.TotalLines != 2 || len(f.Lines) != 2 {
			t.Errorf("want 1/2 lines, got %+v", f)
		}
	})
	t.Run("LCINL-E01: A missing report returns the error", func(t *testing.T) {
		if _, err := ParseLCOV(filepath.Join(t.TempDir(), "none.info")); !os.IsNotExist(err) {
			t.Errorf("want the not-exist error, got %v", err)
		}
	})
}

// `DA:<line>,<hits>,<checksum>` is legal lcov (geninfo --checksum). Splitting in two left
// "5,abc" as the hit count, which failed to parse and marked a covered line uncovered.
func TestLCOVChecksumField(t *testing.T) {
	t.Run("LCINL-B08: A line entry with a checksum field keeps its hit count", func(t *testing.T) {
		rep, err := ParseLCOV(write(t, "c.info", "SF:a.ts\nDA:1,5,PF4Rz2r7RTliO9u6bZ7h6g\nDA:2,0,XyZ\nend_of_record\n"))
		if err != nil {
			t.Fatal(err)
		}
		if f := rep.Files[0]; !f.Lines[1] || f.Lines[2] || f.CoveredLines != 1 || f.TotalLines != 2 {
			t.Errorf("want line 1 covered, line 2 not, 1/2; got %+v", f)
		}
	})
}

// Entries before the first `SF:` belong to no file. They were counted and the counts were
// not reset, so they leaked into the first record's totals.
func TestLCOVEntriesBeforeFirstRecord(t *testing.T) {
	t.Run("LCINL-B09: Entries before the first source-file line belong to no file", func(t *testing.T) {
		rep, err := ParseLCOV(write(t, "c.info", "DA:1,1\nDA:2,1\nLF:9\nLH:9\nSF:a.ts\nDA:1,0\nend_of_record\n"))
		if err != nil {
			t.Fatal(err)
		}
		if len(rep.Files) != 1 {
			t.Fatalf("want one file, got %+v", rep.Files)
		}
		if f := rep.Files[0]; f.TotalLines != 1 || f.CoveredLines != 0 {
			t.Errorf("a.ts is 0/1, got %d/%d", f.CoveredLines, f.TotalLines)
		}
	})
}
