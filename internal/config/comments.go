// Package config carrega os defaults embutidos e a config do projeto.
package config

import "strings"

// CommentMarkers mapeia extensão de arquivo → prefixos de comentário de linha.
// O CLI nunca parseia código; ele só precisa saber onde uma anotação (que vive
// num comentário) começa, para extraí-la do texto. Ver DECISIONS.md D2/D4.
//
// Defaults embutidos cobrem as linguagens comuns; o projeto pode estender ou
// sobrescrever via anchors.yaml (a carregar em config.Load — TODO).
var CommentMarkers = map[string][]string{
	// C-family / chaves
	".go":    {"//"},
	".ts":    {"//"},
	".tsx":   {"//"},
	".js":    {"//"},
	".jsx":   {"//"},
	".java":  {"//"},
	".c":     {"//"},
	".h":     {"//"},
	".cpp":   {"//"},
	".cs":    {"//"},
	".rs":    {"//"},
	".swift": {"//"},
	".kt":    {"//"},
	".scala": {"//"},
	".php":   {"//", "#"},
	// hash-family
	".py":   {"#"},
	".rb":   {"#"},
	".sh":   {"#"},
	".yaml": {"#"},
	".yml":  {"#"},
	".toml": {"#"},
	".r":    {"#"},
	// SQL / Lua / Haskell
	".sql": {"--"},
	".lua": {"--"},
	".hs":  {"--"},
	// markup — o comentário é de bloco; tratado à parte quando necessário
	".md":   {"<!--"},
	".html": {"<!--"},
	".xml":  {"<!--"},
}

// MarkersFor devolve os prefixos de comentário de linha para uma extensão.
// Extensão desconhecida → nil (o scanner ainda pode buscar a anotação como
// texto solto, mas sem âncora de comentário confiável).
func MarkersFor(ext string) []string {
	return CommentMarkers[ext]
}

// LineCommentFor devolve o marcador de comentário de LINHA para um caminho, pela
// extensão. Default `#`: é o marcador da maioria das linguagens de script e de
// configuração, e — o que decide — um `#` num arquivo C-like é visivelmente errado para
// quem lê, enquanto um `//` num arquivo Python é sintaxe INVÁLIDA que quebra o arquivo.
// Entre errar visível e errar quebrando, erra-se visível.
func LineCommentFor(path string) string {
	ext := path
	if i := strings.LastIndex(path, "."); i >= 0 {
		ext = strings.ToLower(path[i:])
	}
	// `.test.ts` / `_test.py`: a extensão real é sempre a última.
	if markers, ok := CommentMarkers[ext]; ok && len(markers) > 0 {
		if markers[0] == "<!--" {
			return "<!--" // markup: quem chama trata o fechamento
		}
		return markers[0]
	}
	return "#"
}
