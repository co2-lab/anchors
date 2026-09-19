<!-- @anchors
  code: MCSTM
  updated_at: 2026-09-19
  layer: gate
-->
# MockStamped — the double carries the mark of the snippet it replaces, and the gate RECOMPUTES it

> **Code**: `MCSTM`

## Overview

Confronts a test against the question no other ruler asks about a test double: **the double
claims to replace a snippet of the real module — is it still the snippet that exists today?**

It is the AGNOSTIC half of the pair that attacks mock drift. `mock-typed` solves the problem
where the language helps — structural types let the compiler check the module's surface — and
only there: most ecosystems have no `Partial<typeof X>`, and not even in TypeScript does it
reach the SHAPE of the returned VALUE. Measured: 202 tests doubled a React Query hook with
two fields where the real one returns about 15, and the type does not tell "deliberate
partial" apart from "out of date".

The stamp does not interpret the code: it READS it. A text hash catches any change —
signature, body, type, constant — and works in Python, Ruby, Go or plain JS, with no
per-language extractor.

**The gate RECOMPUTES instead of validating format**, and that is what makes the mechanism
resistant to whoever writes it. A stamp nobody confronts is theatre: whoever edits the test
would regenerate it to match their own mock, and it would start certifying itself.

**What separates it from its neighbours:** `triad-complete` asks whether the test exists,
`feature-test-match` confronts scenario against case, `tests-green` reads the run — all three
answer "yes" about a double that froze a contract which no longer exists. `mock-typed` demands
the TIE to the real module; this gate demands the recomputable MARK of the snippet.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | only `mapx.KindTest` is judged | spec, code, feature — they leave without a verdict | this unit: a kind that is not a test returns Skip |
| the double dialect | the regex the PROJECT declares in `derived.mock_detect`, with at least one capture group | any built-in dialect default — there is none | the project's Structure; an undeclared dialect switches the gate off, and one that does not compile fails loudly |
| the doubled module | the specifier written in the test, in any spelling the dialect captures | third-party library modules, which the graph does not govern | this unit, resolving the specifier against the map |
| the stamp annotation | `@contract: <path> \| <anchor line> \| <line count> \| <hash>`, count greater than zero | an annotation with a non-numeric or non-positive count, which is discarded | this unit, via the stamp pattern |
| the anchor | a line of the real module occurring EXACTLY once | a line occurring zero times, or more than once — both are findings, not errors | this unit: the ambiguity is reported, never resolved by guessing |
| the real module | a file readable from the project root | a module that no longer exists, which is a finding | this unit, which names the unreadable file instead of crashing |

## Effects

| Effect | Description |
| --- | --- |
| `MCSTM-B01` | An artifact that is not a test leaves without a verdict: only a test declares doubles. |
| `MCSTM-B02` | A stamp whose recomputed hash matches the module today passes. |
| `MCSTM-B03` | A snippet that changed since the stamp was written fails, and the verdict shows the value the contract has today. |
| `MCSTM-B04` | The stamp is immune to DISPLACEMENT: the anchor is searched by CONTENT, so editing lines above does not invalidate the stamp of a snippet that did not change. |
| `MCSTM-B05` | An anchor that vanished — renamed, removed or rewritten — fails with its own message: it is a finding, not a tool error, because the double is certainly out of date. |
| `MCSTM-B06` | An anchor occurring more than once fails as ambiguous: a stamp pointing at "one of the two" proves nothing, and the gate reports rather than choosing. |
| `MCSTM-B07` | The declared line count delimits the window: a change beyond it is not reached, and one inside it is. The reach stays in plain sight of whoever reads, and the gate needs no per-language parser to find where the block ends. |
| `MCSTM-B08` | Without the dialect declared by the project the gate goes quiet: adopting the stamp is the project's decision. |
| `MCSTM-B09` | A double of a module the project does not govern is not charged. |
| `MCSTM-B10` | A stamp whose module no longer exists on disk is a finding with its own message, not a crash. |
| `MCSTM-B11` | The ABSENCE of a stamp on a governed double is accused, not skipped — and this is the most important half of the gate. |
| `MCSTM-B12` | A stamp that is present and correct satisfies both charges at once: the absence and the correspondence. |
| `MCSTM-B13` | A dialect regex that does not compile fails LOUDLY: it is a configuration error, and silencing it would make the gate sweep zero doubles and report green. |
| `MCSTM-B14` | A dialect regex with no capture group fails too: without it the gate cannot know WHICH module was doubled. |
| `MCSTM-B15` | A different ecosystem's dialect is charged exactly the same way once declared — the stamp is agnostic in fact, not in intention. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MCSTM-I01` | The gate RECOMPUTES the hash against the real module; it never validates the stamp's format alone. A stamp nobody confronts would certify itself, because whoever edits the test regenerates it to match their own mock. | changes the real module without touching the stamp and verifies the verdict turns |
| `MCSTM-I02` | The double and the stamp are tied by the module's PATH, matched by suffix without extension. The specifier and the stamped path describe the same file through alias and disk path, and a per-ecosystem alias resolver would be needed otherwise. | declares the double through an alias and the stamp through the disk path, and verifies the charge is satisfied |
| `MCSTM-I03` | A configuration fault fails; a project decision skips. Compiling error and missing capture group are Fail, an undeclared dialect is Skip — the difference is whether somebody CHOSE the silence. | confronts the three configurations and verifies Fail, Fail and Skip |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MCSTM-X01` | Does not interpret the code of the stamped module. | It hashes text. That is what makes it work in Python, Ruby, Go or plain JS with no per-language extractor — and what lets it catch a change of body, constant or type, which a signature extractor would miss. |
| `MCSTM-X02` | Does not carry a built-in dialect for detecting doubles. | The dialect comes from `derived.mock_detect` in the project's Structure. Embedding the jest/vitest pattern as a default would make the gate run in a Python project, match zero doubles and report GREEN over what it never checked — the worst possible failure in a measuring device. The gate is language-agnostic BY DESIGN: in Go there is no call to detect at all, because the double is a satisfied interface, and the right verdict there is that the gate does not apply. |
| `MCSTM-X03` | Does not skip the absence of a stamp to accommodate legacy code. | Turning the charge on in an old project produces hundreds of findings at once, and designing for that case would turn a MIGRATION problem into a permanent property of the framework: every future project would inherit the slack. A project born with Anchors has no debt — the first double is written after the gate exists. Legacy is handled with the vocabulary that already exists: a non-blocking gate during adoption, and per-unit opt-out with a written reason. |
| `MCSTM-X04` | Does not guarantee cryptographic strength: the hash is truncated. | The stamp lives in a comment line and is read by a human. Thirty-two bits are enough to detect the accidental change this gate pursues — there is no adversary forging a collision against their own test. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindTest` | core — jurisdiction starts from the node's kind, and the `Graph` says which modules the project governs |
| DEP2 | `internal/config/config.go` | `Config` | core — the double dialect is declared by the project, never assumed by the gate |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — the findings name the drift in the reader's language |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
