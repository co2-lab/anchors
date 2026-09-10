package testsig

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ChangedLines mapeia arquivo (relativo à raiz) → conjunto de linhas ADICIONADAS ou
// ALTERADAS naquele arquivo, segundo um diff. É a metade "o que mudei" da cobertura
// de diff; a outra metade (o que está coberto) vem do lcov ingerido.
type ChangedLines map[string]map[int]bool

// GitDiff roda `git diff --unified=0 <ref>` na raiz e devolve as linhas mudadas por
// arquivo. `ref` é a base da comparação (ex.: "main", "HEAD~1"); vazio compara o
// working tree contra HEAD. Só linhas do lado NOVO (adições) — são as que precisam
// de teste. Requer git; erro se não for um repositório.
func GitDiff(root, ref string) (ChangedLines, error) {
	args := []string{"-C", root, "diff", "--unified=0", "--no-color"}
	if ref != "" {
		args = append(args, ref)
	}
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	return parseUnifiedDiff(string(out)), nil
}

// ParseDiffFile lê um arquivo de unified diff (o fallback não-git: qualquer fonte que
// produza `git diff`/`diff -u` serve). Mesma extração do GitDiff.
func ParseDiffFile(path string) (ChangedLines, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseUnifiedDiff(string(data)), nil
}

// parseUnifiedDiff extrai as linhas do lado NOVO de um unified diff. Acompanha o
// arquivo corrente (linha "+++ b/<path>") e o cursor de linha nova (do hunk
// "@@ -a,b +c,d @@"), incrementando a cada linha de adição ("+").
func parseUnifiedDiff(diff string) ChangedLines {
	changed := ChangedLines{}
	var curFile string
	newLine := 0

	sc := bufio.NewScanner(strings.NewReader(diff))
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "+++ "):
			curFile = stripDiffPath(strings.TrimPrefix(line, "+++ "))
			if curFile != "" && changed[curFile] == nil {
				changed[curFile] = map[int]bool{}
			}
		case strings.HasPrefix(line, "@@"):
			newLine = hunkNewStart(line)
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			if curFile != "" {
				changed[curFile][newLine] = true
			}
			newLine++
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			// remoção: não avança o cursor do lado novo
		case strings.HasPrefix(line, " "):
			newLine++ // contexto (só existe com --unified>0; inofensivo)
		}
	}
	// remove arquivos sem linhas mudadas (ex.: só remoções)
	for f, lines := range changed {
		if len(lines) == 0 {
			delete(changed, f)
		}
	}
	return changed
}

// stripDiffPath remove o prefixo "a/" ou "b/" e o sufixo de timestamp de um caminho
// de diff. "/dev/null" (arquivo deletado) vira vazio.
func stripDiffPath(p string) string {
	p = strings.TrimSpace(p)
	if i := strings.IndexAny(p, "\t"); i >= 0 {
		p = p[:i]
	}
	if p == "/dev/null" {
		return ""
	}
	p = strings.TrimPrefix(p, "a/")
	p = strings.TrimPrefix(p, "b/")
	return filepath.ToSlash(p)
}

// hunkNewStart extrai o início do lado novo de um cabeçalho de hunk "@@ -a,b +c,d @@".
func hunkNewStart(hunk string) int {
	plus := strings.Index(hunk, "+")
	if plus < 0 {
		return 0
	}
	rest := hunk[plus+1:]
	end := strings.IndexAny(rest, ", ")
	if end < 0 {
		end = len(rest)
	}
	n, _ := strconv.Atoi(rest[:end])
	return n
}
