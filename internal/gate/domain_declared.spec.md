<!-- @anchors
  code: DMDCD
  updated_at: 2026-09-26
  layer: gate
-->
# DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid

> **Code**: `DMDCD`

## Overview

Confronts a spec against the question the rest of the framework does not ask: **what does this
unit accept as input, and who guarantees that the invalid never arrives?**

The gap it closes was measured: 71% of the specs of a real project had a section of
rules or effects, and only 14% said what the unit accepts. The whole framework is
built on cataloguing EFFECTS — what the unit does —, and every edge defect
found in three rounds of adversarial review lived in what nobody had declared.

The distinction that gives the gate its reason: `## Constraints` says what the unit does NOT do, and pushes
the duty OUTWARD; `## Domain` says what it ACCEPTS, and NAMES WHO is left with it.
Writing more constraints closes nothing — it creates orphans, because every "not mine" needs
someone on the other side.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the artifact's content | any text, empty included | — (absent text is a case, not an error) | this unit: content without the section is a FAILURE, not an exception |
| the domain section | the title in any language of the catalogue, and the spellings the project uses | a title the catalogue does not name | this unit, by the vocabulary of accepted titles |

> `Who guarantees` cannot be left empty nor say only "not mine": if nobody guarantees, the
> duty is orphaned — and that is exactly where the invalid input gets through. It is the same demand that
> this gate makes of the specs it confronts, applied to itself.

## Effects

| Effect | Description |
| --- | --- |
| `DMDCD-B01` | An artifact that is not a spec leaves the confrontation without a verdict: the gate has no jurisdiction over code, test or guide. |
| `DMDCD-B02` | A spec WITHOUT the domain section FAILS. The absence is not silence: it is the assertion not made. |
| `DMDCD-B03` | A spec with the section OPENED and EMPTY fails too — opening the title without declaring anything is the same hole wearing the appearance of compliance. |
| `DMDCD-B04` | Every declared input must name WHO guarantees it. An input without an owner fails, and the verdict names which ones were left orphaned. |
| `DMDCD-B05` | The waiver is DECLARED and with a written reason. Whoever has no external input records that in the spec, and the gate goes quiet — but the trace that someone looked remains. |
| `DMDCD-B06` | A line filled only with a pending marker is not a declaration: the untouched mould asserts nothing. |
| `DMDCD-B07` | An input whose owner is named passes — it is the other side of the same ruler, and what makes it satisfiable. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DMDCD-I01` | The waiver requires a REASON. A bare waiver mark does not silence the gate — silence without a why is what it exists to prevent. | confronts a spec with the waiver without a reason and verifies that the verdict still demands |
| `DMDCD-I02` | The failing verdict NAMES what is wrong — which input was left without an owner, or that the section is missing. A gate that fails without saying what transfers the diagnostic work to whoever reads it. | confronts a spec with an orphaned input and verifies that its name appears in the verdict |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DMDCD-X01` | Does not judge whether the declared input is RIGHT — only whether it exists and has an owner. | Whether the set of accepted values corresponds to the real domain is judgement, and judgement belongs to another class of gate. Here the ruler is the PRESENCE of the declaration, which is deterministic. |
| `DMDCD-X02` | Does not confront the code to check whether the validation in fact exists. | This layer reads TEXT. Crossing the declaration with the implementation is the work of the relational gate, which has the map; doing it here would duplicate the ruler in two places that would diverge. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `DMDCD-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec` | core — the gate needs the node's KIND to know whether it has jurisdiction |
| DEP2 | `internal/config/config.go` | `Config` | core — the lexicon of accepted titles comes from the project's Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
