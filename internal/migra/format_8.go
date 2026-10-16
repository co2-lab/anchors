// @anchors
//   code: F8AFR
//   ref: MGSTM

package migra

// FORMAT 8 — the core reads no syntax (DESIGN-core-without-syntax.md).
//
// How a file writes a comment, the lexicon of its language and the shapes the gates used to
// know by themselves are now the project's, under `dialect:` and `derived:`. The two files
// carry nothing renamed by this step; crossing it writes the families' examples into the
// project's anchors.yaml, once, where it declared nothing (the command's
// `migrateToDeclaredSyntax`, read from `DeclaredSyntax`), so the project measures what it
// measured before. A project not migrated is refused with the message to run
// `anchors migrate`: read by half, its headers and flags would be read in no file.
func init() {
	Register(Step{
		To:  8,
		Why: "the core reads no syntax: comments, the language's lexicon and the gates' shapes are declared in the dialect",
	})
}
