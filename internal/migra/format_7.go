// @anchors
//   code: F7AFR
//   ref: MGSTM

package migra

// FORMAT 7 — every file a code of its own, of five characters (see filecodes.go).
//
// The two files carry nothing renamed by this step; the step exists so that crossing it
// rewrites the project — the four-character codes widened to five, a `code:` on every
// governed file that can carry one, the measurements carried to the new revisions (the
// command's `migrateToFileCodes`). A project not migrated is refused with the message to
// run `anchors migrate`: read by half, its files would carry no code of their own.
func init() {
	Register(Step{
		To:  7,
		Why: "every governed file gets a code of its own, and four-character codes are widened to five",
	})
}
