package config

import (
	"os"
	"path/filepath"
)

// RaizDoProjeto sobe a partir de `inicio` até achar o diretório que tem o `anchors.yaml`.
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
func RaizDoProjeto(inicio string) string {
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

// AbsRaiz resolve a raiz do projeto a partir do valor de `--root`: absolutiza e, quando o
// usuário não passou nada (o default "."), sobe até o diretório que tem o `anchors.yaml`.
//
// Um `--root` EXPLÍCITO é respeitado como dado: quem aponta para um diretório específico
// está dizendo onde quer trabalhar, e subir por cima disso seria ignorar a instrução.
func AbsRaiz(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if root == "." || root == "" {
		return RaizDoProjeto(abs), nil
	}
	return abs, nil
}
