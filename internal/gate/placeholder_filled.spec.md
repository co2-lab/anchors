<!-- @anchors
  code: PLCFL
  updated_at: 2026-09-20
  layer: gate
-->
# PlaceholderFilled — the skeleton the generator emits must be FILLED IN

> **Code**: `PLCFL`

## Overview

Confronts an artifact against the question no other gate asks: **was this actually
written, or is it still the frame the generator handed over?**

The work prompt promises, textually, that an unfilled placeholder fails. It did not.
Measured: a freshly generated spec — with the layer field, the date, the title and the
first rule all still carrying the generator's marker — crossed EVERY blocking gate with
"can promote", including the header gate, which read a placeholder as a layer and
approved it.

The reason is structural, and it is why this gate cannot be folded into its neighbours:
**header gates validate FORM** — does the field exist, is the shape right — and never ask
whether the value MEANS anything. A marker is a well-formed value. And the completeness
gate counts sections, which the skeleton has all of, empty.

The cost of that silence is the worst kind: an artifact nobody wrote passes as written.
The relational gates that follow find the spec, find the code, confront two things that
reference each other, and the whole pipeline certifies work that does not exist.

Measured before switching it on, against the real repository: **zero findings across 590
specs**. Nothing alive carries a generator marker — which confirms both that whoever
writes, fills in, and that the gate charges only what was left behind.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the content | any text, including empty | — (absent text is a case, not an error) | this unit: nothing written is nothing to charge |
| the marker vocabulary | what the project declares as the pending word | a word the project never declared | the Structure, so the marker is the project's and not the framework's |

## Effects

| Effect | Description |
| --- | --- |
| `PLCFL-B01` | A raw skeleton is FAILED: the artifact was generated and never written. |
| `PLCFL-B02` | A header FIELD whose value is the marker fails — it is the gravest shape, because a placeholder is not a layer and the header gate approves it as well-formed. |
| `PLCFL-B03` | A table CELL holding only the marker fails: the rule exists as a code and says nothing. |
| `PLCFL-B04` | A TITLE or body line opening with the marker fails. |
| `PLCFL-B05` | The verdict NAMES what was left behind, so the reader does not hunt the file for it. |
| `PLCFL-B06` | An artifact with every marker replaced passes — the gate charges the generator's leftovers and nothing else. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PLCFL-I01` | A section the author wrote ON PURPOSE to list pending work is legitimate and is never accused. Measured: 77 specs of a real project carry one. Confusing "the author listed what is missing" with "the author wrote nothing" would punish the honesty the framework asks for everywhere else. | confronts a spec carrying a deliberate pending-work section and verifies it is not accused |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PLCFL-X01` | Does not judge whether what replaced the marker is GOOD. | The ruler is deterministic: the marker is there, or it is not. Whether the sentence that replaced it says something worth saying is judgment, and judgment belongs to another class of gate. |
| `PLCFL-X02` | Does not invent the marker vocabulary. | The pending word comes from the project's Structure. Hardcoding one here would charge a project that writes its pending items differently, and miss the one that matches. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — the marker vocabulary is declared in the Structure |
| DEP2 | `internal/mapx/model.go` | `KindSpec` | core — the kind is what routes the jurisdiction |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
