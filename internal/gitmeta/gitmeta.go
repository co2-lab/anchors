// @anchors
//   code: GTING
//   ref: GTMTG

// Package gitmeta lê metadados do git de um arquivo — hoje, a data do último commit
// que o tocou. É a fonte de verdade para o carimbo de alteração (updated_at): mais
// confiável que uma data mantida à mão, que sempre desatualiza.
package gitmeta

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Today devolve a data do sistema (AAAA-MM-DD). Usada como desempate quando o arquivo
// tem alterações não-commitadas: a "última alteração" é hoje, não o último commit. Um
// gate de data legitimamente precisa saber que dia é — é da natureza dele.
func Today() string { return time.Now().Format("2006-01-02") }

// LastCommitDate devolve a data (AAAA-MM-DD) do ÚLTIMO commit que modificou `rel`
// (caminho relativo à raiz). Só a data — nunca a hora: um arquivo pode ser editado
// num horário e commitado em outro do MESMO dia, e isso é irrelevante; só a mudança
// de DIA importa. Devolve ok=false se não é repo git, o arquivo nunca foi commitado,
// ou o git não está disponível.
func LastCommitDate(root, rel string) (date string, ok bool) {
	cmd := exec.Command("git", "-C", root, "log", "-1", "--format=%ad", "--date=short", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	d := strings.TrimSpace(string(out))
	if d == "" {
		return "", false // arquivo sem commit (novo/untracked)
	}
	return d, true
}

// AllCommitDates devolve, num só `git log`, a data (AAAA-MM-DD) do último commit que
// tocou CADA arquivo do repo — o carimbo de alteração de todos, sem 1 chamada git por
// arquivo. Percorre o log com --name-only e associa a data de cada commit ao primeiro
// aparecimento de cada arquivo (o mais recente, pois o log vem do mais novo ao antigo).
// Mapa vazio se não for repo git.
// AllCommitDatesUnder is AllCommitDates with each path relative to root, where root may be a
// folder inside the repository.
func AllCommitDatesUnder(root string) map[string]string {
	all := AllCommitDates(root)
	prefix := ""
	if p, err := exec.Command("git", "-C", root, "rev-parse", "--show-prefix").Output(); err == nil {
		prefix = strings.TrimSpace(string(p))
	}
	if prefix == "" {
		return all
	}
	out := map[string]string{}
	for f, d := range all {
		if strings.HasPrefix(f, prefix) {
			out[strings.TrimPrefix(f, prefix)] = d
		}
	}
	return out
}

// AllCommitDates reads the dates once per HEAD: the answer only changes with a new commit,
// and the whole history's log cost seconds on every check of a project with a long history
// (4.7 s of a pre-commit, reported). The answer is kept in the repository's common git
// folder, so every worktree at the same HEAD reads it; a handful of HEADs are kept.
func AllCommitDates(root string) map[string]string {
	head, okHead := Commit(root, "HEAD")
	common, okCommon := CommonDir(root)
	// A shallow clone's history deepens under the same HEAD (`git fetch --unshallow`): its
	// answer is not kept.
	shallow, _ := exec.Command("git", "-C", root, "rev-parse", "--is-shallow-repository").Output()
	if !okHead || !okCommon || strings.TrimSpace(string(shallow)) == "true" {
		return allCommitDates(root)
	}
	top := ""
	if out, err := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel").Output(); err == nil {
		top = strings.TrimSpace(string(out))
	}
	dir := filepath.Join(common, "anchors", "commit-dates")
	file := filepath.Join(dir, head+".json")
	if b, err := os.ReadFile(file); err == nil {
		var dates map[string]string
		if json.Unmarshal(b, &dates) == nil {
			return dates
		}
	}
	dates := allCommitDates(top)
	if top == "" {
		return dates
	}
	if b, err := json.Marshal(dates); err == nil && os.MkdirAll(dir, 0o755) == nil {
		tmp := file + ".tmp"
		if os.WriteFile(tmp, b, 0o644) == nil {
			_ = os.Rename(tmp, file)
		}
		pruneCommitDates(dir, 8)
	}
	return dates
}

// pruneCommitDates keeps the newest `keep` answers of the folder.
func pruneCommitDates(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) <= keep {
		return
	}
	type aged struct {
		name string
		at   time.Time
	}
	var all []aged
	for _, e := range entries {
		if info, err := e.Info(); err == nil {
			all = append(all, aged{e.Name(), info.ModTime()})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].at.After(all[j].at) })
	for _, a := range all[keep:] {
		_ = os.Remove(filepath.Join(dir, a.name))
	}
}

func allCommitDates(root string) map[string]string {
	// `-c`: a merge lists the files it changed against every parent — a conflict resolved, a
	// header the hook dated in it — as `git log -- <file>` counts it. Without it a merge lists
	// nothing, and a file dated in a merge read the date of the commit before (reported from a project).
	// A file that came from one side unchanged is not listed: its date is that side's commit.
	cmd := exec.Command("git", "-C", root, "log", "-c", "--format=D:%ad", "--date=short", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return map[string]string{}
	}
	dates := map[string]string{}
	var cur string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "D:") {
			cur = strings.TrimPrefix(line, "D:")
			continue
		}
		f := strings.TrimSpace(line)
		if f == "" || cur == "" {
			continue
		}
		if _, seen := dates[f]; !seen { // 1º aparecimento = commit mais recente
			dates[f] = cur
		}
	}
	return dates
}

// HasUncommittedChanges diz se `rel` tem alterações não-commitadas no working tree.
// Útil para não reclamar de um arquivo em edição ativa (a "última alteração" é hoje,
// não o último commit).
//
// ATENÇÃO ao `false`: ele significa "sem alteração pendente" E TAMBÉM "não deu para
// perguntar" (sem git, sem repo). Quem toma DECISÃO a partir disso deve usar
// `UncommittedChanges`, que separa os dois — confundi-los faz o Anchors afirmar sobre
// um estado que ele não conseguiu ler.
func HasUncommittedChanges(root, rel string) bool {
	mudou, _ := UncommittedChanges(root, rel)
	return mudou
}

// UncommittedChanges é o `HasUncommittedChanges` que não esconde a ignorância: o
// segundo retorno diz se a pergunta pôde ser FEITA. Sem ele, "o arquivo está limpo" e
// "não há repositório para consultar" chegam ao chamador como o mesmo `false`, e a
// diferença entre os dois é justamente o que decide se há um veredito a dar.
func UncommittedChanges(root, rel string) (mudou, sabido bool) {
	cmd := exec.Command("git", "-C", root, "status", "--porcelain", "--", rel)
	out, err := cmd.Output()
	if err != nil {
		return false, false
	}
	return strings.TrimSpace(string(out)) != "", true
}

// StagedChanges is UncommittedChanges over the INDEX: whether the commit being made
// changes the file — staged against HEAD, or added. An unstaged edit does not count: it is
// not in the commit. known=false outside git.
func StagedChanges(root, rel string) (changed, known bool) {
	out, err := exec.Command("git", "-C", root, "diff", "--cached", "--name-only", "--", rel).Output()
	if err != nil {
		return false, false
	}
	return strings.TrimSpace(string(out)) != "", true
}

// StagedChangesOnce is StagedChanges for a whole run: one `git diff --cached` lists every
// staged file, and each question is answered from the list. Asked a file at a time, the commit
// hook ran one git process per target — 853 of them for a project's commit of 8 files.
func StagedChangesOnce() func(root, rel string) (changed, known bool) {
	var mu sync.Mutex
	lists := map[string]map[string]bool{}
	return func(root, rel string) (bool, bool) {
		mu.Lock()
		defer mu.Unlock()
		staged, ok := lists[root]
		if !ok {
			out, err := exec.Command("git", "-C", root, "diff", "--cached", "--name-only", "--relative", "-z").Output()
			if err != nil {
				lists[root] = nil
				return false, false
			}
			staged = map[string]bool{}
			for _, f := range strings.Split(string(out), "\x00") {
				if f != "" {
					staged[f] = true
				}
			}
			lists[root] = staged
		}
		if staged == nil {
			return false, false
		}
		return staged[filepath.ToSlash(rel)], true
	}
}

// Head devolve o hash curto do HEAD e o assunto do commit. Serve para carimbar um
// relatório com a versão do código que ele descreve: sem isso, uma leitura salva
// não diz sobre QUE código ela fala, e envelhece sem avisar.
//
// ok=false quando não é repo git, o git não está disponível, ou o repo ainda não
// tem commit nenhum.
func Head(root string) (short, subject string, ok bool) {
	out, err := exec.Command("git", "-C", root, "log", "-1", "--format=%h%x00%s").Output()
	if err != nil {
		return "", "", false
	}
	partes := strings.SplitN(strings.TrimRight(string(out), "\n"), "\x00", 2)
	if len(partes) != 2 || partes[0] == "" {
		return "", "", false
	}
	return partes[0], partes[1], true
}

// DirtyCount conta os arquivos com alteração não-commitada na árvore inteira.
// Um relatório rodado sobre árvore suja descreve um estado que não está em commit
// nenhum — quem relê precisa saber que não vai reencontrá-lo pelo hash.
//
// Devolve -1 quando não deu para contar (sem git, sem repo). NÃO é 0: zero afirma
// "árvore limpa", e carimbar um relatório com essa afirmação sem tê-la verificado é
// exatamente o silêncio que tranquiliza — o pior tipo.
func DirtyCount(root string) int {
	out, err := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	if err != nil {
		return -1
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}

// AtHead returns a file's content as the last commit has it. ok=false when there is no
// repository, no commit, or the file is not in HEAD (a new file): there is no earlier
// version to compare with.
func AtHead(root, rel string) (content string, ok bool) {
	out, err := exec.Command("git", "-C", root, "show", "HEAD:./"+rel).Output()
	if err != nil {
		return "", false
	}
	return string(out), true
}

// CommitAt is the commit a run at that moment measured: the last commit of HEAD made at or
// before it. ok=false outside git, or when no commit is that old.
func CommitAt(root string, at time.Time) (sha string, ok bool) {
	out, err := exec.Command("git", "-C", root, "rev-list", "-1", "--before="+at.Format(time.RFC3339), "HEAD").Output()
	sha = strings.TrimSpace(string(out))
	return sha, err == nil && sha != ""
}

// Commit is the commit a reference names (a hash, a tag, a branch), in full; ok=false when it
// names none.
func Commit(root, ref string) (sha string, ok bool) {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--verify", "--quiet", ref+"^{commit}").Output()
	sha = strings.TrimSpace(string(out))
	return sha, err == nil && sha != ""
}

// CommitTime is when a commit was made; ok=false when git cannot say.
func CommitTime(root, sha string) (time.Time, bool) {
	out, err := exec.Command("git", "-C", root, "show", "-s", "--format=%cI", sha).Output()
	if err != nil {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
	return t, err == nil
}

// AtCommit returns a file's content as a commit has it; ok=false when the file is not in it.
func AtCommit(root, sha, rel string) (content []byte, ok bool) {
	out, err := exec.Command("git", "-C", root, "show", sha+":./"+rel).Output()
	if err != nil {
		return nil, false
	}
	return out, true
}

// ProjectFiles are the project's files as git sees them — tracked, and new ones it does not
// ignore —, by their path relative to the root; known=false outside git.
func ProjectFiles(root string) (files []string, known bool) {
	out, err := exec.Command("git", "-C", root, "ls-files", "-z", "--cached", "--others", "--exclude-standard").Output()
	if err != nil {
		return nil, false
	}
	seen := map[string]bool{}
	for _, f := range strings.Split(string(out), "\x00") {
		if f != "" && !seen[f] {
			seen[f] = true
			files = append(files, f)
		}
	}
	return files, true
}

// DirtyFiles are, in one `git status`, the files with uncommitted changes — modified,
// staged, added or untracked —, by their path relative to the root; known=false outside
// git. It answers for a whole run what UncommittedChanges answers for one file, without a
// git process per file (thousands in a large repository).
func DirtyFiles(root string) (dirty map[string]bool, known bool) {
	out, err := exec.Command("git", "-C", root, "status", "--porcelain", "-z", "--untracked-files=all").Output()
	if err != nil {
		return nil, false
	}
	prefix := ""
	if p, err := exec.Command("git", "-C", root, "rev-parse", "--show-prefix").Output(); err == nil {
		prefix = strings.TrimSpace(string(p))
	}
	dirty = map[string]bool{}
	entries := strings.Split(string(out), "\x00")
	for i := 0; i < len(entries); i++ {
		e := entries[i]
		if len(e) < 4 {
			continue
		}
		status, path := e[:2], e[3:]
		if status[0] == 'R' || status[0] == 'C' {
			i++ // a rename or a copy carries its origin in the next entry
		}
		if prefix != "" {
			if !strings.HasPrefix(path, prefix) {
				continue
			}
			path = strings.TrimPrefix(path, prefix)
		}
		dirty[path] = true
	}
	return dirty, true
}

// Worktrees are the folders of every worktree of the repository root belongs to, each at the
// same place inside it as root is in its own: agents working in parallel run in worktrees of
// their own, and what they do there is the project's too (reported from a project). Without git, or
// outside a repository, it is root alone.
func Worktrees(root string) []string {
	top, err := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return []string{root}
	}
	rel, err := filepath.Rel(strings.TrimSpace(string(top)), root)
	if err != nil {
		rel = "."
	}
	if r, err := filepath.EvalSymlinks(root); err == nil {
		if t, err := filepath.EvalSymlinks(strings.TrimSpace(string(top))); err == nil {
			if x, err := filepath.Rel(t, r); err == nil {
				rel = x
			}
		}
	}
	out, err := exec.Command("git", "-C", root, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return []string{root}
	}
	var dirs []string
	for _, line := range strings.Split(string(out), "\n") {
		if p, ok := strings.CutPrefix(line, "worktree "); ok {
			dirs = append(dirs, filepath.Join(filepath.FromSlash(p), rel))
		}
	}
	if len(dirs) == 0 {
		return []string{root}
	}
	return dirs
}

// CommonDir is the folder git shares between every worktree of the repository root belongs
// to — the main `.git` —, absolute; false outside a repository, or with no git.
func CommonDir(root string) (string, bool) {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return "", false
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", false
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(root, dir)
	}
	return dir, true
}

// HooksDir is the folder git runs the repository's hooks from — `core.hooksPath` when set,
// the common `.git/hooks` otherwise —, absolute; false outside a repository, or with no git.
func HooksDir(root string) (string, bool) {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--git-path", "hooks").Output()
	if err != nil {
		return "", false
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", false
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(root, dir)
	}
	return dir, true
}

// Renames are the files git reads as moved since HEAD — staged or in the tree —, from their
// old path to their new one, both relative to root; nil outside git, before the first commit,
// or with nothing moved. git decides a move by content, as `git mv` and `git status` say it.
func Renames(root string) map[string]string {
	out, err := exec.Command("git", "-C", root, "diff", "HEAD", "--relative", "-M", "--name-status", "-z").Output()
	if err != nil {
		return nil // @resilient: no git, or no HEAD yet: nothing was moved since a commit
	}
	var moves map[string]string
	entries := strings.Split(string(out), "\x00")
	for i := 0; i+2 < len(entries); i++ {
		if !strings.HasPrefix(entries[i], "R") {
			continue
		}
		if moves == nil {
			moves = map[string]string{}
		}
		moves[entries[i+1]] = entries[i+2]
		i += 2
	}
	return moves
}
