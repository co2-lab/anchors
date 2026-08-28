package testsig

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// FileCoverage é a cobertura de linha de um arquivo, do lcov.
type FileCoverage struct {
	File         string
	CoveredLines int
	TotalLines   int
	// Lines: por linha instrumentada, coberta (true) ou não (false). Permite a
	// cobertura de DIFF — cruzar as linhas mudadas com as descobertas.
	Lines map[int]bool
}

// UncoveredIn devolve, dentre `changed`, as linhas deste arquivo que são
// INSTRUMENTADAS mas NÃO cobertas. Linhas mudadas que o lcov não instrumentou (ex.:
// comentários, linhas em branco) não contam — não são executáveis.
func (f FileCoverage) UncoveredIn(changed map[int]bool) []int {
	var out []int
	for line := range changed {
		if covered, instrumented := f.Lines[line]; instrumented && !covered {
			out = append(out, line)
		}
	}
	return out
}

// InstrumentedIn conta, dentre `changed`, quantas são instrumentadas (a base do %
// de diff coverage — comentários/blanks não entram no denominador).
func (f FileCoverage) InstrumentedIn(changed map[int]bool) int {
	n := 0
	for line := range changed {
		if _, ok := f.Lines[line]; ok {
			n++
		}
	}
	return n
}

// Percent devolve a cobertura de linha em 0..100.
func (f FileCoverage) Percent() float64 {
	if f.TotalLines == 0 {
		return 0
	}
	return float64(f.CoveredLines) / float64(f.TotalLines) * 100
}

// CoverageReport é a cobertura de todos os arquivos.
type CoverageReport struct {
	Files []FileCoverage
}

// ParseLCOV lê um arquivo lcov (o formato .info, universal — jest, go tool cover
// -func não, mas gcov/istanbul/pytest-cov sim). Cada registro:
//
//	SF:<arquivo>
//	DA:<linha>,<hits>      (uma por linha instrumentada)
//	LF:<total de linhas>   (opcional; se ausente, contamos as DA)
//	LH:<linhas cobertas>   (opcional)
//	end_of_record
//
// Contamos a partir das DA (robusto) e usamos LF/LH se presentes.
func ParseLCOV(path string) (*CoverageReport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rep := &CoverageReport{}
	var cur *FileCoverage
	var daCovered, daTotal int
	var lf, lh int

	flush := func() {
		if cur == nil {
			return
		}
		fc := *cur
		if lf > 0 { // prefere LF/LH quando o reporter os emitiu
			fc.TotalLines = lf
			fc.CoveredLines = lh
		} else {
			fc.TotalLines = daTotal
			fc.CoveredLines = daCovered
		}
		rep.Files = append(rep.Files, fc)
		cur, daCovered, daTotal, lf, lh = nil, 0, 0, 0, 0
	}

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "SF:"):
			flush()
			cur = &FileCoverage{File: strings.TrimPrefix(line, "SF:"), Lines: map[int]bool{}}
		case strings.HasPrefix(line, "DA:"):
			parts := strings.SplitN(strings.TrimPrefix(line, "DA:"), ",", 2)
			if len(parts) == 2 {
				daTotal++
				ln, _ := strconv.Atoi(parts[0])
				hits, err := strconv.Atoi(parts[1])
				covered := err == nil && hits > 0
				if covered {
					daCovered++
				}
				if cur != nil && ln > 0 {
					cur.Lines[ln] = covered
				}
			}
		case strings.HasPrefix(line, "LF:"):
			lf, _ = strconv.Atoi(strings.TrimPrefix(line, "LF:"))
		case strings.HasPrefix(line, "LH:"):
			lh, _ = strconv.Atoi(strings.TrimPrefix(line, "LH:"))
		case line == "end_of_record":
			flush()
		}
	}
	flush()
	return rep, sc.Err()
}
