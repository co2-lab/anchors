// @anchors
//   code: RTCNR
//   ref: PRRPR

package config

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectRoot sobe a partir de `inicio` até achar o diretório que tem o `anchors.yaml`.
// Devolve `inicio` quando não acha — o chamador então falha com a mensagem normal de
// config ausente, que é o comportamento certo fora de um projeto Anchors.
//
// Existe porque o default de `--root` era "." e todo comando falhava ao ser chamado de um
// subdiretório:
//
//	$ cd packages/backend && anchors check --all
//	Error: carregar config: open .../packages/backend/anchors.yaml: no such file
//
// O relato veio de uma execução real: o diretório de trabalho persistia entre comandos, um
// `cd` anterior fazia o comando seguinte falhar, e a mensagem não dizia que a causa era o
// CWD — mostrava um caminho que ninguém pediu. Ferramentas de repositório (git, npm, cargo)
// sobem até a raiz justamente porque quem trabalha não fica parado nela.
func ProjectRoot(inicio string) string {
	abs, err := filepath.Abs(inicio)
	if err != nil {
		return inicio
	}
	dir := abs
	for {
		if _, err := os.Stat(filepath.Join(dir, DefaultFile)); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir { // chegou em / sem achar
			return inicio
		}
		dir = pai
	}
}

// AbsRoot resolve a raiz do projeto a partir do valor de `--root`: absolutiza e, quando o
// usuário não passou nada (o default "."), sobe até o diretório que tem o `anchors.yaml`.
//
// Um `--root` EXPLÍCITO é respeitado como dado: quem aponta para um diretório específico
// está dizendo onde quer trabalhar, e subir por cima disso seria ignorar a instrução.
//
// The root is written as the disk names it (TrueCase): on a file system that ignores case,
// the working directory comes as the shell typed it, and a tool handed paths built on it
// does not find its own files under them.
func AbsRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if root == "." || root == "" {
		return TrueCase(ProjectRoot(abs)), nil
	}
	return TrueCase(abs), nil
}

// TrueCase is an absolute path as the disk names it: on a file system that ignores case,
// `/Users/x/mif` and `/Users/x/MIF` are one folder, and the disk calls it `MIF`. Only the
// case changes — a symbolic link is kept as written —, and where the system cannot say, the
// path comes back as given (diskPath, per system).
//
// Measured in MIF on macOS: the shell's working directory was `.../mif`, the folder is
// `.../MIF`, and `jest --findRelatedTests` handed `.../mif/...` matched none of the files it
// knows under `.../MIF/`; it ran no test and exited 0, in five suites.
func TrueCase(path string) string {
	if !filepath.IsAbs(path) {
		return path
	}
	disk, ok := diskPath(path)
	if !ok || disk == path {
		return path
	}
	// The kernel's path also resolves links (`/var` is `/private/var` on macOS): the case is
	// taken from its last parts, as many as stand under no link — those the resolved path
	// names as the given one does, but for case.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	sep := string(filepath.Separator)
	given, real, typed := strings.Split(path, sep), strings.Split(disk, sep), strings.Split(resolved, sep)
	for i, j, k := len(given)-1, len(real)-1, len(typed)-1; i >= 0 && j >= 0 && k >= 0; i, j, k = i-1, j-1, k-1 {
		if given[i] != typed[k] || !strings.EqualFold(typed[k], real[j]) {
			break
		}
		given[i] = real[j]
	}
	return strings.Join(given, sep)
}
