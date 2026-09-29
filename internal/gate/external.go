package gate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/shell"
)

// argvLimit é quanto os alvos podem ocupar, em bytes, numa ÚNICA chamada do shell.
//
// O Windows monta a linha de comando como uma string só e o CreateProcess a corta em
// 32767 chars: um `scope: batch` sobre um projeto inteiro (1484 alvos ≈ 85 KB) falhava
// ANTES de rodar, com "o nome do arquivo ou a extensão é muito grande". Como o exec
// nem chegava a executar, a saída voltava vazia e o gate reprovava MUDO — indistinguível
// de violação real, e verde no macOS, onde o ARG_MAX é de centenas de KB.
//
// O teto do Windows é MUITO menor que os 32767 porque o orçamento não é gasto só aqui:
// o gate costuma empilhar wrappers (sh → yarn → node → binário), e cada camada reexpande
// os caminhos — o sh do MSYS2 converte cada relativo em absoluto ao chamar um .exe
// nativo. Medido no app de referência com o gate eslint: 8379 bytes de alvos ainda passam, 11389 já
// estouram DENTRO do script, longe da vista do Anchors. 6000 deixa a margem.
//
// O teto do Unix fica alto de propósito: onde o lote já cabia, continua sendo UMA
// execução, e o comportamento não muda.
//
// ANCHORS_ARGV_MAX ajusta o teto para quem empilha mais (ou menos) wrapper que isso —
// o número certo é propriedade dos gates do projeto, não do Anchors.
func argvLimit() int {
	if v := os.Getenv("ANCHORS_ARGV_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if runtime.GOOS == "windows" {
		return 6000
	}
	return 100000
}

// RunExternalArgs é o motor: roda o comando com N alvos como argumentos posicionais
// ($1, $2, … e "$@"). Um alvo é o caso por-nó; vários, o `scope: batch`; nenhum, o
// `scope: project` (a ferramenta olha o projeto inteiro e não recebe alvo).
//
// A garantia de segurança é a mesma qualquer que seja a quantidade: os caminhos vão
// em argv, nunca interpolados no script. `{{files}}` é reescrito para "$@" pela mesma
// razão que `{{file}}` vira "$1" — conveniência, sem abrir caminho para injeção.
//
// Quando os alvos não cabem numa linha de comando, o gate roda EM LOTES e os vereditos
// se somam: reprova se qualquer lote reprovar, e o laudo junta a saída de todos. Isso
// pressupõe um gate que julga arquivo a arquivo — o que todo gate de `scope: batch` é,
// por construção (recebe um recorte arbitrário do projeto, nunca o conjunto fechado).
// Um gate que precise ver todos os alvos DE UMA VEZ (achar duplicata entre arquivos,
// por exemplo) deve declarar `scope: project` e varrer sozinho, não depender do lote.
func RunExternalArgs(command string, targets []string, root string) (Verdict, string) {
	script := strings.ReplaceAll(command, "{{file}}", `"$1"`)
	script = strings.ReplaceAll(script, "{{files}}", `"$@"`)

	// Under `--index` a tool reads the commit's content, not the tree's: a target whose
	// tree differs from the index is handed over as a copy of what the index holds.
	targets, clean, done := indexedTargets(root, targets)
	defer done()

	var falhas []string
	for _, lote := range sliceTargets(targets, argvLimit()-len(script)) {
		out, err := execShell(script, lote, root)
		out = clean(out)
		if err == nil {
			continue
		}
		if detalhe := strings.TrimSpace(out); detalhe != "" {
			falhas = append(falhas, detalhe)
			continue
		}
		// Reprovou sem dizer nada. Acontece quando o próprio exec falha (não há `sh`
		// no PATH, a linha de comando estourou) — sem esta linha o laudo fica vazio e
		// o operador lê "violação de código" onde há problema de ambiente.
		falhas = append(falhas, "gate não produziu saída: "+err.Error())
	}
	if len(falhas) == 0 {
		return Pass, ""
	}

	detail := strings.Join(falhas, "\n")
	// Ferramenta de projeto (tsc/eslint) reporta MUITOS erros de uma vez. O corte em
	// 500 do caso por-nó deixaria o laudo inútil justamente onde ele mais importa —
	// some a lista de arquivo:linha, que é o conteúdo que resolve o problema.
	limite := 500
	if len(targets) != 1 {
		limite = 4000
	}
	if len(detail) > limite {
		detail = detail[:limite] + "\n… (saída truncada)"
	}
	return Fail, detail
}

// execShell roda UM lote. `sh -c '<script>' sh <alvo...>` → os alvos viram $1..$N
// (e "$@") dentro do script, sem passar pela tokenização do shell (não é injeção).
func execShell(script string, targets []string, root string) (string, error) {
	cmd, err := shell.Command(script, append([]string{"sh"}, targets...)...) // os alvos vão como argv, não interpolados
	if err != nil {
		return err.Error(), err
	}
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// sliceTargets divide os alvos em lotes que cabem numa linha de comando.
//
// Sem alvos devolve UM lote vazio, não lote nenhum: `scope: project` é justamente a
// execução sem alvo, e devolver nada faria o gate não rodar (passando por omissão).
// Um alvo sozinho maior que o teto vai só no seu lote — não há como parti-lo, e é
// melhor deixar o SO recusar com a mensagem dele do que silenciar o alvo.
func sliceTargets(targets []string, teto int) [][]string {
	if len(targets) == 0 {
		return [][]string{nil}
	}
	if teto < 1 {
		teto = 1
	}

	var lotes [][]string
	var atual []string
	tam := 0
	for _, t := range targets {
		custo := len(t) + 1 // o separador que o SO conta entre um argumento e o próximo
		if len(atual) > 0 && tam+custo > teto {
			lotes = append(lotes, atual)
			atual, tam = nil, 0
		}
		atual = append(atual, t)
		tam += custo
	}
	return append(lotes, atual)
}

// runGateCommand invoca um comando externo (jest, eslint, tsc…). O CLI NÃO reimplementa
// a ferramenta — roda e lê o exit code (D5: reimplementamos parsing, não ferramentas
// de teste/lint). O comando do gate vem do anchors.yaml (config do projeto) e pode
// ser composto (ex.: "cd apps/mobile && jest \"$1\""), por isso roda via `sh -c`.
//
// SEGURANÇA: o caminho do alvo NUNCA é interpolado na string do shell — ele é
// passado como ARGUMENTO POSICIONAL (`$1`), que o shell não reinterpreta. Assim um
// nome de arquivo com metacaracteres (foo;rm.tsx) é dado puro, não injeção. O
// comando deve referenciar o alvo como "$1" (com aspas). `{{file}}` é aceito por
// conveniência e reescrito para "$1" antes de rodar.
//
// exit 0 = pass; qualquer outro = fail (com stderr/stdout no detalhe).
//
// It runs where the gate's `workdir` says: the tree, or — under `--index`, with
// `workdir: index` — a copy of what the commit records.
func runGateCommand(g config.Gate, targets []string, root string) (Verdict, string) {
	if g.Workdir != config.WorkdirIndex {
		return RunExternalArgs(g.Run, targets, root)
	}
	dir, clean, err := indexWorkdir(root)
	if err != nil {
		return Skip, i18n.T("gate.workdir_index_failed", err.Error())
	}
	v, out := RunExternalArgs(g.Run, targets, dir)
	return v, clean(out)
}

// indexCopy is the copy of the git index this run's `workdir: index` gates share: made
// once, on the first gate that asks, and removed when the gates are done.
var indexCopy struct {
	top, dir string // the repository's top, and where its copy is
	root     string // the project root the copy was made for, and what it answers
	answer   string
	err      error
}

// indexWorkdir is where a `workdir: index` command runs for the project at root, and what
// turns the copy's paths in its output back into the project's. With no index source set
// (no `--index`), or a tree that holds nothing the index does not, it is root itself.
func indexWorkdir(root string) (dir string, clean func(string) string, err error) {
	same := func(s string) string { return s }
	fileSourceMu.RLock()
	src := fileSource
	fileSourceMu.RUnlock()
	if src == nil {
		return root, same, nil
	}
	if indexCopy.root != root {
		releaseIndexWorkdir()
		indexCopy.root = root
		indexCopy.answer, indexCopy.err = copyIndex(root)
	}
	if indexCopy.err != nil {
		return "", same, indexCopy.err
	}
	if indexCopy.dir == "" {
		return root, same, nil
	}
	prefixes := []string{indexCopy.dir}
	if real, err := filepath.EvalSymlinks(indexCopy.dir); err == nil && real != indexCopy.dir {
		prefixes = append(prefixes, real) // a tool may print the resolved path (/private/var on macOS)
	}
	top := indexCopy.top
	return indexCopy.answer, func(s string) string {
		for _, d := range prefixes {
			s = strings.ReplaceAll(s, d, top)
			s = strings.ReplaceAll(s, filepath.ToSlash(d), filepath.ToSlash(top))
		}
		return s
	}, nil
}

// releaseIndexWorkdir removes this run's copy of the index.
func releaseIndexWorkdir() {
	if indexCopy.dir != "" {
		_ = os.RemoveAll(indexCopy.dir)
	}
	indexCopy.top, indexCopy.dir, indexCopy.root, indexCopy.answer, indexCopy.err = "", "", "", "", nil
}

// copyIndex writes what the git index holds to a temporary folder and answers where the
// project root is inside it. When the tree holds nothing the index does not — no unstaged
// change, no untracked file — nothing is copied and the answer is root itself.
//
// The folders and files git ignores are the project's tooling (dependencies, build caches,
// local settings): they are linked into the copy, not copied, so the command finds them.
func copyIndex(root string) (string, error) {
	git := func(args ...string) (string, error) {
		out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
		return string(out), err
	}
	top, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}
	top = filepath.Clean(strings.TrimSpace(top))
	prefix, _ := git("rev-parse", "--show-prefix")
	prefix = strings.TrimSpace(prefix)
	if err := exec.Command("git", "-C", top, "diff", "--quiet").Run(); err == nil {
		if untracked, err := git("ls-files", "--others", "--exclude-standard", "--full-name", "--", ":/"); err == nil && strings.TrimSpace(untracked) == "" {
			return root, nil // the tree is the commit
		}
	}
	dir, err := os.MkdirTemp("", "anchors-index-")
	if err != nil {
		return "", err
	}
	indexCopy.top, indexCopy.dir = top, dir
	if out, err := exec.Command("git", "-C", top, "checkout-index", "--all", "--force", "--prefix="+filepath.ToSlash(dir)+"/").CombinedOutput(); err != nil {
		return "", fmt.Errorf("git checkout-index: %v %s", err, strings.TrimSpace(string(out)))
	}
	ignored, err := exec.Command("git", "-C", top, "ls-files", "--others", "--ignored", "--exclude-standard", "--directory", "-z").Output()
	if err != nil {
		return "", fmt.Errorf("git ls-files --ignored: %w", err)
	}
	for _, rel := range strings.Split(string(ignored), "\x00") {
		rel = strings.TrimSuffix(rel, "/")
		if rel == "" {
			continue
		}
		src, dst := filepath.Join(top, filepath.FromSlash(rel)), filepath.Join(dir, filepath.FromSlash(rel))
		if err := mirrorIgnored(top, dir, src, dst, mirrorDepth); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, filepath.FromSlash(prefix)), nil
}

// mirrorDepth is how deep an ignored folder is looked into for links back into the project:
// a workspace package is linked as `deps/<name>` or `deps/@scope/<name>`.
const mirrorDepth = 2

// mirrorIgnored puts the ignored src at dst in the copy. It is a link to src, unless a link
// inside it (to mirrorDepth levels) points back into the project: that one would reach the
// tree's file, not the commit's. Then the folder is made for real, the link is made again
// to the same place in the copy, and everything else is linked.
func mirrorIgnored(top, dir, src, dst string, depth int) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	if depth == 0 || !linksInto(top, src, depth) {
		return linkTo(src, dst)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s, d := filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())
		if target, ok := innerTarget(top, s); ok {
			rel, _ := filepath.Rel(top, target)
			if err := linkTo(filepath.Join(dir, rel), d); err != nil {
				return err
			}
			continue
		}
		if e.IsDir() {
			if err := mirrorIgnored(top, dir, s, d, depth-1); err != nil {
				return err
			}
			continue
		}
		if err := linkTo(s, d); err != nil {
			return err
		}
	}
	return nil
}

// linksInto says whether a link inside src, to depth levels, points into the project.
func linksInto(top, src string, depth int) bool {
	entries, err := os.ReadDir(src)
	if err != nil {
		return false
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		if _, ok := innerTarget(top, s); ok {
			return true
		}
		if e.IsDir() && depth > 1 && linksInto(top, s, depth-1) {
			return true
		}
	}
	return false
}

// innerTarget is where the link at p points, when p is a link and it points into the
// project at top.
func innerTarget(top, p string) (string, bool) {
	t, err := os.Readlink(p)
	if err != nil {
		return "", false
	}
	if !filepath.IsAbs(t) {
		t = filepath.Join(filepath.Dir(p), t)
	}
	t = filepath.Clean(t)
	rel, err := filepath.Rel(top, t)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return t, true
}

// linkTo makes dst point at src: a symbolic link, or on Windows — where one needs a
// privilege most accounts lack — a junction for a folder.
func linkTo(src, dst string) error {
	err := os.Symlink(src, dst)
	if err == nil || runtime.GOOS != "windows" {
		return err
	}
	if fi, statErr := os.Stat(src); statErr == nil && fi.IsDir() {
		if out, jerr := exec.Command("cmd", "/c", "mklink", "/J", dst, src).CombinedOutput(); jerr != nil {
			return fmt.Errorf("mklink /J %s: %v %s", dst, jerr, strings.TrimSpace(string(out)))
		}
		return nil
	}
	return err
}
