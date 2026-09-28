<!-- @anchors
  code: FLMDF
  updated_at: 2026-09-28
  layer: apoio
-->
# FlowModel — the questions a work-flow graph answers: what comes next, and what is broken

> **Code**: `FLMDF`

## Overview

A catalogue of "what to do in case X" makes whoever works remember to consult it, choose right among
many options, and not skip a step, three chances to err on every round. A flow inverts the burden:
instead of "choose among eight", it says "from here, the exits are these". This unit answers the
questions asked of the flow graph built from the project's flow and action files (`FLBLF`): the exits of
a state, the states of a flow in the order the author wrote them, where a flow starts, and the name of
the action a step fits.

It also answers the three questions that find a broken flow. A STEP no transition reaches (other than
where a flow starts) is written, documented and unreachable. A RESULT an action declares that no flow
routes leaves whoever gets it with nowhere to go. An exit that points at a state that does not exist
looks like a path and leads nowhere. Steps are reached; results are declared by their action, so a
result is never asked "who reaches you?".

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the flow graph | the graph built from the project's flows, or none | — | the flow builder (`FLBLF`); a project without flows gives no graph |
| the codes | step codes and result codes (letter R) | — | the flow files |

## Effects

| Effect | Description |
| --- | --- |
| `FLMDF-B01` | The exits of a state are exactly the transitions declared from it; a state with no exit has none. (`Next`) |
| `FLMDF-B02` | A state is found by its code, and a code with no state is not found. (`StateByCode`) |
| `FLMDF-B03` | The states of a flow come in the order its file declares them, never alphabetical, and the flows are listed once each in order of appearance. (`StatesOf`, `Flows`) |
| `FLMDF-B04` | The entry of a flow is its first declared step that is not a result; a file that declares only results has no entry. (`Entry`) |
| `FLMDF-B05` | A code is a result when the letter after its hyphen is the outcome letter, O (`config.OutcomeLetter`; it was R until format 5). (`IsResult`) |
| `FLMDF-B06` | The title of an action is the name of the action file that declares its results; an action nobody wrote has none. (`ActionTitle`) |
| `FLMDF-B07` | A step no transition reaches is unreachable, except the first state of each flow; results are never unreachable. (`Unreachable`) |
| `FLMDF-B08` | A result that no transition routes is unhandled. (`Unhandled`) |
| `FLMDF-B09` | A transition to a state that does not exist is dangling. (`Dangling`) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLMDF-X01` | A project without flows answers nothing to every question, and none of them fails. | A project that has not declared a flow yet is not broken; its commands must keep working. |

## Errors

none — every question over a missing graph or an unknown code answers "nothing" (X01), and the broken flows are findings the unit returns, not failures it handles.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `FlowGraph`, `FlowState`, `FlowTransition` | mapa — the flow graph's shape |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
