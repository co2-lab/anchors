<!-- @anchors
  code: DSCDC
  updated_at: 2026-09-19
  layer: gate
-->
# DocSelfContained — the spec has to stand on its own

> **Code**: `DSCDC`

## Overview

Confronts a spec against the reader who does not have the repository open: **the reference
points somewhere else — does it bring what it points at?**

The spec is read by two audiences, and one of them has no checkout. The compiled `docs/*.md`
inherits the spec's text word for word, and whoever reads the documentation to learn what the
system does has no use for the name of a plan file — navigating to it is exactly what the
documentation mechanism exists to eliminate.

**What the gate does NOT charge.** The reference WITH the text alongside is right and stays:
the reader has the argument in hand. What is left over is the SCAFFOLD — the sentence that
announces the quotation and then does not quote — and it is the scaffold that sends the
person away. Measured in a real project: 48 mentions across 37 specs, almost all of them
followed by the quoted passage.

**No vocabulary, anywhere.** The gate does not look for "plan", "see" or "per": Anchors
governs projects in any language, and a gate that matches words passes in silence over the
project written in the other one — which is worse than not existing, because the spec then
LOOKS protected. What it matches is STRUCTURE, and only what Anchors itself defines: the
PATH of a map node written in the body, and the `{CODE}-R000N` form that is the doctrine's
revision identity.

Unlike `docs-fresh`, this one is informative and offers no command that fixes it: rewriting
a sentence is the work of whoever wrote it, and blocking a commit over a question of form
would stop the flow. The gate marks, and the correction rides along with the card that
already touches the spec.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | only a spec is confronted | plan, feature, test, code | this unit: the plan references sibling plans by function, and the feature does not become documentation prose |
| the map | a built graph, or none | — | this unit: with no map there is no list of what counts as a reference, and the confrontation is skipped |
| what counts as a reference | the id and the base name of any map node that is not this one, and the `{CODE}-R000N` revision form | a path invented by hand, and this node's own path | the map: it is what knows which files exist, in whatever language the project names its folders |
| what counts as accompanying content | a quotation in any written tradition, a table, a list item, or — for a revision only — substantive prose on the same line | prose measured by vocabulary instead of length | this unit: a keyword list would pass in silence over the project written in the other language |

## Effects

| Effect | Description |
| --- | --- |
| `DSCDC-B01` | A reference that comes WITH the passage it announces passes: the reader has the argument in hand and does not leave the page. |
| `DSCDC-B02` | A reference that only POINTS is accused, and the finding names the line and shows it, so the correction does not need a hunt. |
| `DSCDC-B03` | A quotation counts in any written tradition — straight quotes, typographic ones, guillemets, German low quotes, CJK corner brackets. |
| `DSCDC-B04` | A path inside a code fence is not a reference: there it is an EXAMPLE, a command to run, not a sentence sending the reader away. |
| `DSCDC-B05` | The spec citing its OWN path is identifying itself, not sending anybody anywhere, and is not accused. |
| `DSCDC-B06` | A rule code is not a revision: the `-R000N` form is the one the doctrine reserves, and accusing the other would charge the spec for naming its own subject. |
| `DSCDC-B07` | Only the spec is charged; every other kind leaves without a verdict. |
| `DSCDC-B08` | A revision cited with an EXPLANATION on the same line passes: the code is the label and the sentence is the content. |
| `DSCDC-B09` | With no map the confrontation is skipped, because the list of what counts as a reference comes from the map and nowhere else. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DSCDC-I01` | The ruler matches STRUCTURE, never vocabulary. A spec written in English, Spanish, German or Japanese is charged by exactly the same code, because no language changes a path. | writes the same empty reference in several languages and verifies the same accusation |
| `DSCDC-I02` | The on-line-explanation escape belongs to the REVISION and not to the PATH. A revision code is a label whose sentence is the content; a path is the place the person would have to go, and no amount of surrounding prose says what is there. | puts long prose around a bare path and verifies the accusation stands |
| `DSCDC-I03` | The verdict names the LINE and shows what it says. A gate that accuses without pointing hands the diagnostic work to whoever reads it. | confronts an empty reference and verifies the line number and the excerpt in the message |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DSCDC-X01` | Does not judge whether the accompanying content is FAITHFUL to what the reference announces. | Whether the quoted passage really says what the sentence claims is judgement, and judgement belongs to another class of gate. Here the ruler is whether the reader is left with something to read, which is deterministic. |
| `DSCDC-X02` | Does not offer a command that fixes it, and does not block. | Rewriting a sentence is the work of whoever wrote it — there is nothing to generate. Blocking a commit over a question of form would stop the flow, and the gate that stops the flow is the gate that gets turned off. |
| `DSCDC-X03` | Errs on the side of letting things through when measuring whether a line explains something. | The threshold is a coarse and admittedly imperfect ruler. This is an informative gate, and mass false positives are what make someone switch it off — so the error is deliberately pushed to the permissive side. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec`, `KindCode` | core — the kind gives jurisdiction, and the code nodes are left out of what counts as a reference |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — the findings are written in the reader's language, since the gate itself matches no vocabulary |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
