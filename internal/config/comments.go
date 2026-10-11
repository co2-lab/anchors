// @anchors
//   code: CMCNC
//   ref: CMMRC

// Package config carrega os defaults embutidos e a config do projeto.
package config

import "strings"

// CommentMarkers mapeia extensão de arquivo → prefixos de comentário de linha.
// O CLI nunca parseia código; ele só precisa saber onde uma anotação (que vive
// num comentário) começa, para extraí-la do texto. Ver DECISIONS.md D2/D4.
//
// The table is built in and fixed: anchors.yaml has no key that extends or overrides it,
// and Load does not touch it. An extension missing here gets no prefix from MarkersFor
// and the hash comment from LineCommentFor.
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
	// OCaml has only the block comment `(* … *)`; Menhir's grammars take it and `/* … */`.
	// Like markup, a flag or a header field written in it is closed (CommentCloseFor).
	".ml":  {"(*"},
	".mli": {"(*"},
	".mll": {"(*"},
	".mly": {"(*", "/*"},
	// markup — o comentário é de bloco; tratado à parte quando necessário
	".md":   {"<!--"},
	".html": {"<!--"},
	".xml":  {"<!--"},
}

// CommentNamed are the comment markers of the files known by NAME, with no extension: a
// dune build file comments with `;`.
var CommentNamed = map[string][]string{
	"dune":         {";"},
	"dune-project": {";"},
	"dune-file":    {";"},
}

// CommentCloseFor is what closes a comment opened with the marker, with its leading space —
// ` -->` after `<!--`, ` *)` after `(*` —, or empty for a line comment, which the end of
// the line closes. A flag written at a line's end in a block-comment language is closed
// with it, or the comment would swallow the code that follows.
func CommentCloseFor(marker string) string {
	switch marker {
	case "<!--":
		return " -->"
	case "(*":
		return " *)"
	case "/*":
		return " */"
	}
	return ""
}

// MarkersFor devolve os prefixos de comentário de linha para uma extensão.
// Extensão desconhecida → nil (o scanner ainda pode buscar a anotação como
// texto solto, mas sem âncora de comentário confiável).
func MarkersFor(ext string) []string {
	return CommentMarkers[ext]
}

// MarkersForPath are the comment markers of a file: by its extension, or by its name for
// the files known by name (CommentNamed).
func MarkersForPath(path string) []string {
	base := path
	if i := strings.LastIndexAny(path, "/\\"); i >= 0 {
		base = path[i+1:]
	}
	if m, ok := CommentNamed[base]; ok {
		return m
	}
	ext := ""
	if i := strings.LastIndex(base, "."); i >= 0 {
		ext = strings.ToLower(base[i:])
	}
	return CommentMarkers[ext]
}

// LineCommentFor devolve o marcador de comentário de LINHA para um caminho, pela
// extensão. Default `#`: é o marcador da maioria das linguagens de script e de
// configuração, e — o que decide — um `#` num arquivo C-like é visivelmente errado para
// quem lê, enquanto um `//` num arquivo Python é sintaxe INVÁLIDA que quebra o arquivo.
// Entre errar visível e errar quebrando, erra-se visível.
func LineCommentFor(path string) string {
	// `.test.ts` / `_test.py`: a extensão real é sempre a última.
	if markers := MarkersForPath(path); len(markers) > 0 { // every entry of the table has a prefix
		return markers[0] // a block opener (`<!--`, `(*`): the caller closes it (CommentCloseFor)
	}
	return "#"
}
