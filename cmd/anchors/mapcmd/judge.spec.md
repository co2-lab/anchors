<!-- @anchors
  code: JDGUE
  updated_at: 2026-09-26
  layer: comando
-->
# Judge — records an AI's verdict on a judgment gate with the same bookkeeping as a deterministic gate

> **Code**: `JDGUE`

## Overview

A judgment gate asks a question no program answers: does this screen follow the atomic design guide, does
this excerpt realize the rule it cites. The CLI does not judge. The AI that operates Anchors reads the
gate's guide, confronts the target, and reports the verdict here; this command then does what a
deterministic gate's run does — it stamps the map and opens or resolves the issue — so a subjective
verdict enters the same anti-drift loop. The stamp carries the target's revision, so the verdict goes
stale when the target changes afterwards.

There are three verdicts. A pass says the target was verified; a fail says it was not, and its reason is
the full report — the issue body is that reason verbatim, so nobody has to reread the target to learn
what to fix. The third, waived, exists because the other two lie when the thing the question is about
does not exist yet: a spec that declares its code as still to be written has no excerpt to judge, and a
pass would leave a stamp that looks like a verification that never happened. The older spelling of the
waiver is still accepted, because it sits in committed maps and scripts, and it becomes the canonical one
everywhere — including the stamp.

A target that is not in the map yet falls back to the piece of the same unit that is — typically the
spec, when the code is not written — so a review prescribed at the start of a unit has somewhere to be
recorded. The cycle's own `review` is accepted without being declared as a gate, because it runs on what
was just delivered and not over every node. In manual workflow mode the verdict is stamped but no issue
file is written unless asked; the report is printed instead. A fail may carry a patch, which becomes an
applicable suggestion, and recording a verdict closes the judge task the queue held for it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the verdict | pass, fail, waived, or the legacy spelling of waived, in any letter case | anything else | this unit: refuses it |
| the reason | the full report; required for fail and waived | an empty reason on fail or waived | this unit: refuses it |
| the gate | a declared judgment gate by its exact name, or `review` | a gate that is not a judgment gate, or a legacy name | this unit: refuses it |
| the target | a node of the map, or a file whose unit has a piece in the map | a file whose unit has no piece in the map | this unit: refuses it |
| the project | a root with a map and a configuration file | a root with neither | this unit: refuses, naming what is missing |
| the patch | a readable file with the fixing diff | a path that cannot be read | this unit: refuses it |

## Effects

| Effect | Description |
| --- | --- |
| `JDGUE-B01` | The verdict is read case-insensitively, and the legacy spelling of the waiver is a waiver from validation to the stamp and the printed word. |
| `JDGUE-B02` | A pass needs no reason; a fail and a waiver must name one (`ValidateVerdict` refuses them otherwise). |
| `JDGUE-B03` | The gate must be a declared judgment gate matched by its exact name; the cycle's `review` is accepted without being declared. |
| `JDGUE-B04` | A gate that declares a guide is stamped on the edge from the guide to the target, with `ok` for pass, `issue` for fail and `waived` for waiver, and the judgment counts as answered for that gate. |
| `JDGUE-B05` | A gate with no guide edge to the target is stamped on every edge of the target, with the gate's name. |
| `JDGUE-B06` | A target not in the map is recorded on the piece of the same unit that is in the map, and the output says which. |
| `JDGUE-B12` | The piece of the same unit is found by name from a code or a test target: the test suffixes (`.test`, `.spec`, `_test`) are dropped, and the candidates are the spec, the feature, and the Go and TypeScript code and test names, plus the target's own extension. |
| `JDGUE-B07` | A fail opens an issue whose body is the report; the same report again changes nothing, and a different report reopens the issue and appends it. |
| `JDGUE-B08` | A pass resolves the open issue of that gate and target; a waiver resolves it too, announcing a waiver and never a pass. |
| `JDGUE-B09` | In manual workflow mode the verdict is stamped and the report printed, but no issue is written unless recording issues is asked; in every other mode the issue is written. |
| `JDGUE-B10` | A fail with a patch opens a fix suggestion named after the gate and the target, carrying the patch. |
| `JDGUE-B11` | Recording a verdict closes the queue's judge task for that gate and target; the pending switch lists the queued judge tasks — those of the judgment kind, and those with the legacy judge verb — or says none is waiting. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `JDGUE-I01` | The stamp always tells the three verdicts apart: a waiver is never stamped `ok`, whatever spelling it came in. | records a waiver in both spellings and reads `waived` on the guide edge |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `JDGUE-X01` | Does not announce a waiver as a pass, in any output. | A pass in the output undoes the point of the third verdict: the reader would believe the target was verified. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `JDGUE-E01` | No target is given and the pending switch is off. | Refused, showing how to call it. | There is nothing to record the verdict on. |
| `JDGUE-E02` | No gate is given. | Refused: the gate is mandatory. | A verdict without its gate answers no question. |
| `JDGUE-E03` | The verdict is none of the three. | Refused, naming the accepted verdicts. | Any other word would be stamped as a pass. |
| `JDGUE-E04` | REF[JDGUE-B02]: a fail or a waiver without a reason is the refusal B02 states | — | — |
| `JDGUE-E05` | The gate is not a declared judgment gate. | Refused, saying it is not a judgment gate. | Stamping an undeclared gate would create answers no check ever asks for. |
| `JDGUE-E06` | There is no map. | Refused, pointing at `anchors map build`. | The stamp lives in the map. |
| `JDGUE-E07` | The target is not in the map and no piece of its unit is. | Refused, naming the target. | A stamp on nothing would be lost at once. |
| `JDGUE-E08` | There is a map but no configuration. | Refused: the configuration could not be loaded. | The gate cannot be validated without its declaration. |
| `JDGUE-E09` | The patch file cannot be read. | Refused before any issue or suggestion is written; the verdict was already stamped on the map. | A suggestion without its diff is not applicable. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/stamp.go` | `StampEdgeByGate`, `StampNodeByGate` | mapa — stamping the verdict |
| DEP2 | `internal/issue/issue.go` | `Open`, `Reopen`, `Resolve` | infra — the issue of the finding |
| DEP3 | `internal/suggestion/suggestion.go` | `Open` | infra — the fix suggestion |
| DEP4 | `internal/queue/queue.go` | `List`, `MarkDone` | infra — the judge tasks |
| DEP5 | `internal/config/config.go` | `Load` | config — the judgment gates and the workflow mode |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
