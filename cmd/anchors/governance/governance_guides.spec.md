<!-- @anchors
  code: GVGDG
  updated_at: 2026-10-07
  layer: comando
-->
# GovernanceGuides — the guides an agent reads to operate Anchors, and the contracts other code relies on

> **Code**: `GVGDG`

## Overview

Anchors does not call an AI; the AI, in its own client, calls Anchors as a tool. `anchors guide` is how it learns to: with no subcommand it prints the operating playbook (the development flow, the commands, what to report), and each subcommand prints the ruler for one artifact or one phase — plan, spec, code, feature, test, guide, header, product doctrine, feature flag, flow, the discover phase of a new project, reviewing a PR and working a card. The texts are embedded in the binary on purpose: a guide that travels with the version is read as it is today, while an instruction copied into a card freezes on the day the card was born.

This spec governs the set as a whole: `register.go` (the four commands this package adds to the root), `guide.go` (the playbook and the command tree), and the files that each hold one guide's text. Most of those texts are prose with no rule of their own. A few are not: they carry facts that other code or the pipeline depends on, and a text that drifts from them sends every reader the wrong way with no error anywhere. Those facts are stated here as rules:

- every `anchors <command>` a guide tells the reader to run must exist;
- the review guide teaches the exact verdict line the pull-request pipeline parses, and who counts;
- the work guide uses the real label names the workflows create;
- the review guide's conformance points are continuous and each one is anchored in the prose above it;
- the work and review guides state the board facts an agent must not break: the claim comes first, the CI verdict is waited for, and nobody moves or closes a card by hand.

The review and work guides are also the two that tell an agent what to do with what it does not know, so they append the autonomy section (`AutonomyGuide`), read from the role declared in the project the reader points at.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the subcommand | none, or one of the thirteen guide names | a name that is not a guide | the command tree refuses it as an unknown command |
| the root of `guide review` and `guide work` | any directory; by default the project found from the current directory | — | the autonomy section: a directory with no declaration reads as no role |

## Effects

| Effect | Description |
| --- | --- |
| `GVGDG-B01` | `anchors guide` with no subcommand prints the operating playbook. |
| `GVGDG-B02` | `anchors guide` has exactly fifteen subcommands — changelog, code, feature, flag, flow, guide, header, plan, product, project, report-bug, review, spec, test, work — each with a short description, each printing its own guide. |
| `GVGDG-B03` | `guide review` and `guide work` print their guide followed by the autonomy section of the root given by `--root`; a root that is not a project still prints, with the section for no role declared. |
| `GVGDG-B04` | Registering this package adds exactly `guide`, `audit`, `governs` and `compliance` to the root. |
| `GVGDG-B05` | The review guide teaches that only the reviewer the claim assigned counts, only with a line posted after the assignment, that the last line wins, and that a line inside a code block is an example, not a verdict. |
| `GVGDG-B06` | The review guide says the reviewer does not move the card: the checks move it to `ready-to-review` and the merge to `ready-to-test`, and a wrong state costs more than a late one. |
| `GVGDG-B07` | The review guide carries a conformance points section, under a heading the `guide-checklist` gate recognises, with points `REV-CK1` to `REV-CK18` and no gap, each anchored in the prose above the list, and says the list is not a substitute for the checks. |
| `GVGDG-B08` | The work guide teaches the claim (`anchors next`, with `ANCHORS_SESSION` declared) before the board order, and why: a card that never enters the review column leaves the next agent finding it empty. |
| `GVGDG-B09` | The work guide says pushing is not a stopping point: the agent waits for the CI with `--watch`, and the only exits are green or an escalation. |
| `GVGDG-B10` | The work guide says the agent does not close the card — the merge does — and what closing early breaks. |
| `GVGDG-B11` | The work guide sends a finding that is not the card's own through `anchors escalate` tied to the card, and has the open decision queue read before a new `--for-user` escalation is taught. |
| `GVGDG-B12` | The project guide covers the discover phase — `PROJECT.md`, `INSIGHTS.md`, the five technical stages, one question at a time — and the playbook points to it before the plan phase. |
| `GVGDG-B13` | The test guide names the proving instrument for each shape of input space and teaches the stamp refresh for a doubled function. |
| `GVGDG-B14` | `anchors guide --help` lists every subcommand once, with what it teaches, and lists nothing that is not a subcommand. |
| `GVGDG-B16` | The changelog guide says `anchors changelog` builds a technical changelog, not the product's, and recommends a product changelog an agent synthesizes from it: breaking changes, visible features and bugs fixed go in, fixes without `Bug:` stay out, and chores only when they matter to the product. |
| `GVGDG-B17` | The spec guide asks for the spec in four passes — every section, with its questions (every input, every effect and its limit, every failure, every state shared with another unit); the variations of each rule; the generalization of cases into an invariant; a review that each rule states the intent and not the mechanism —, and says a defect is a rule not written, to be written with its invariant before the fix. |
| `GVGDG-B22` | The spec guide asks every unit that reads environment variables to declare them — type, required, default, values, deprecated — as the contract `env-declared` confronts with the code and the variables page is compiled from. |
| `GVGDG-B23` | The spec guide asks a unit that loads data for four states at least — loading, empty, load error and loaded —, the load error its own and never written as "no data" or "not found", its failure a rule naming the source it answers (`failure-declared`), each fallible call reading its error beside it (`failure-handled`), and points at the screen preset that writes them. |
| `GVGDG-B24` | The header guide names the flags beside the code — the dependency flag on each import line and its waiver, the used-by flag above each imported symbol (naming its symbol above an export list), the navigation flag on each call and its waiver —, the gates that confront each, and that the fixer writes what the code proves; the navigation guide, a subcommand of its own, shows a screen spec's In and Out tables, the flag on every navigation call, a back navigation included, the four navigation gates, the entry routes, and the page and command that show the map. Their examples are written so this repository does not read them as its own flags. |
| `GVGDG-B21` | The spec guide ties every validation of a visual unit to a State Flow transition (or `@no-state: <reason>`) and every error to the message it shows (or `@no-message: <reason>`), and the test and feature guides ask a visual-regression capture of every message as of every state. |
| `GVGDG-B20` | The test guide strongly recommends a contract test for every API unit, validating the API against the OpenAPI compiled from the specs with the language's own tool, named by `{CODE}-CT`, loading the document instead of copying it, and asserting each error response's status, code and message. |
| `GVGDG-B19` | The test guide strongly recommends a visual-regression test for every state a visual unit's spec registers, tied both ways — a VR scenario per state, a VR test naming `{CODE}-VR-<state>` and a baseline image per scenario, every scenario of a state and every test of a scenario, a state with no visual value exempted by `@no-vr: <reason>` on its line —, the spec guide points every state of a visual unit to it, and the feature guide asks for one VR scenario per state. |
| `GVGDG-B18` | The report-bug guide tells a bug in Anchors from a problem of the project, asks for a made-up minimal case in Anchors' terms because the repository is public, a `--dry-run` before sending, and while the fix does not come: no hand edits to Anchors' files, a waiver that names the issue, and the user told. The playbook points to it and to `anchors report-bug`. |
| `GVGDG-B15` | The work guide, in every mode, tells a fix from a bug — a bug is a defect that shipped —, asks for each fix as its own `fix` commit with a `Bug:` footer only on a bug, and for the failing test first; the code and review guides point to it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GVGDG-I01` | Every `anchors <command>` that the playbook or any guide cites is a command of the real command tree. | builds the whole tree from every package's registration and confronts each command every guide prints |
| `GVGDG-I02` | The work guide names the labels by the names the workflows create — `anchors:under-` and `anchors:needs-user` — and never the old `anchors:sob-`. | confronts the printed work guide with the label constants the workflows are seeded from |
| `GVGDG-I03` | Each verdict line the review guide teaches, with a name in place of `<you>`, is read by the pull-request pipeline as that verdict by that name. | extracts the verdict expression from the seeded `anchors-pr-checks.yml` and applies it to the lines the guide prints |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GVGDG-X01` | Only `guide review` and `guide work` depend on the project: they alone take `--root`, and every other guide prints its text alone. | The other rulers are the framework's and read the same everywhere; only the autonomy section depends on a local declaration. |

## Errors

none — every guide prints an embedded text and has no failure to handle; a `--root` that is not a project is not a failure either, and `GVGDG-B03` states what it prints.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/governance/guide_autonomy.go` | `printAutonomy` | the section appended to the review and work guides |
| DEP2 | `internal/config/root.go` | `AbsRoot` | the root of the review and work guides |

The label names and the verdict line are not called by this code: they are facts the texts must match. Their sources — the label constants of `internal/initx/workflows.go` and the verdict expression of the seeded `anchors-pr-checks.yml` — are confronted by the proofs of `GVGDG-I02` and `GVGDG-I03`.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
