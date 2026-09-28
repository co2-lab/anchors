package migra

// FORMAT 5 — the letters of the artifacts that are not specs, in English.
//
// A plan's phase was `F` (*Fase*) and a flow's step `P` (*Passo*): Portuguese, in a
// vocabulary that is English. An action's result was `R`, the letter of a spec's
// permission rule. They are now `W` (Wave), `T` (Task) and `O` (Outcome) — the decision,
// and why these three, is recorded beside the constants (`config.PhaseLetter`,
// `config.StepLetter`, `config.OutcomeLetter`).
//
// LOSS IF NOT CONVERTED: SILENT. A plan whose phases stay `-F01` has no phase the gates
// recognize — `phase-ordered` reads "plan without declared phases" and passes, and every
// `needs:` pointing at a phase stops resolving. A flow whose results stay `-R01` loses
// them: `flow next` no longer tells a result from a step. So the minimum readable format
// goes up with this one: a project not migrated is refused, with the message that says
// to run `anchors migrate`, instead of being read by half.
//
// No key changes: the two YAML files only record the format; the codes are rewritten in
// the project's files by the command (`RewriteCodeLetters`).

func init() {
	Register(Step{
		To:  5,
		Why: "the letters of plans, flows and actions in English: phase F→W, step P→T, result R→O (of actions and of flows fitted as pieces)",
		RenameLetters: []LetterRename{
			{Kind: "plan", From: "F", To: "W"},
			{Kind: "flow", From: "P", To: "T"},
			// A flow fits as a piece inside another flow, and then has RESULTS of its own
			// (`JUDGE-R01` in the worker flow). Left as `R`, `flow next` would read it as a
			// step once results are `O` — silently.
			{Kind: "flow", From: "R", To: "O"},
			{Kind: "action", From: "R", To: "O"},
		},
	})
}
