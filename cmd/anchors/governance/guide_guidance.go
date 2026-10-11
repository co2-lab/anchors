// @anchors
//   code: GDGDG
//   ref: GVGDG

package governance

// guidanceGuide is the ruler of the guidance layer — the files at the project's root that tell
// whoever arrives, person or agent, how the project is worked on (DESIGN-guidance-layer.md).
// It is what `guidance-honored` reads each guidance file against.
const guidanceGuide = `# Guidance guide (README, CONTRIBUTING, AGENTS, CLAUDE: how the project says it is worked on)

Whoever arrives at a project — a person or a coding agent — learns how to work in it from a
handful of files at its root, and does what they say. Under Anchors those files are the
GUIDANCE layer, a layer every project has: Anchors knows it by its files, wherever the code
lives. Its files carry no @anchors header.

## What each file is for

- CONTRIBUTING.md — the working agreement, in one place. The other three send the reader here.
- AGENTS.md — read first by most coding agents: it sends the agent to CONTRIBUTING.md before
  any change, and names the few commands an agent needs at once.
- CLAUDE.md — read first by Claude Code: it sends the agent to CONTRIBUTING.md and imports
  AGENTS.md (a line '@AGENTS.md'), so both agents read the same text.
- README.md — what the project is; for whoever will change it, a pointer to CONTRIBUTING.md.

'anchors init' seeds CONTRIBUTING.md, AGENTS.md and CLAUDE.md from the configuration when they
are missing, and never touches one that exists.

## What the agreement must say

The guidance says what the project's configuration (anchors.yaml) decides — never the
opposite. A guidance file under Anchors:

1. THE ORDER OF THE WORK — the spec is the anchor: a change starts in the spec, the feature
   turns each rule into scenarios, the tests prove the scenarios, the code realises them.
2. THE LAYERS — the ones anchors.yaml declares, kept apart. It names no layer the project
   does not declare, and proposes no folder layout.
3. WHERE WORK COMES FROM — 'anchors next' gives the next task; 'anchors work <stage> --for
   <target>' frames one.
4. WHAT BARS A COMMIT — the gates anchors.yaml declares blocking. A gate that fails is
   fixed, or waived in the commit message with its reason ('[skip-<rule>@<CODE>: why]').
   NEVER '--no-verify': the hooks are where the agreement is enforced.
5. WATCHING LONG WORK — every long command (a suite, a mutation run, a build, a script, an
   agent's task) launched through 'anchors monitor run -- <command>' in the background, with
   'anchors monitor' running while it works; 'anchors install-hooks --agent' for a coding
   agent.

## What contradicts the agreement (the gate 'guidance-honored' fails it)

- an instruction to skip a gate, the hooks or the tests ('commit with --no-verify', 'tests
  can come later');
- an order of work other than the anchor's (code first, the spec written after);
- a layer named that anchors.yaml does not declare, or a folder layout proposed as the
  project's structure;
- a gate called blocking that the configuration does not block, or a blocking one presented
  as optional;
- long work left unwatched as the way to do it.

The judgment asks for contradictions, never for style or length: a short guide that says
nothing wrong passes. Each failure names the passage and proposes the text that replaces it.

## The gates

- 'guidance-linked' (informs) — CONTRIBUTING.md exists, and each other guidance file present
  points to it.
- 'guidance-honored' (informs, AI judgment) — each guidance file read against this guide and
  against anchors.yaml; asked again when the file or anchors.yaml changes.
`
