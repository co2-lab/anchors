<!-- @anchors
  code: VTRST
  updated_at: 2026-10-03
  layer: gate
-->
# StateTransitions — every change of a visual unit is proven through what the screen shows

> **Code**: `VTRST`

## Overview

A validation changes the screen: the form refuses, a field turns red, a button disables, the flow
takes another path. Each of those is a state, and the states are already proven by their visual
capture. So a validation does not need a capture of its own; it needs to say which state it leads to.
An error, too, shows on the screen — as the message the spec catalogs; the state is the same, what
changes is the text, and the messages are captured.

Two gates tie each change of a visual unit to what the screen shows:

| Gate | Question |
| --- | --- |
| `validation-transitions` | Is every validation the trigger of a State Flow transition, from the state it is checked in to the state it leads to — or does it say `@no-state: <reason>`? |
| `error-message-declared` | Does every error name the message it shows — or say `@no-message: <reason>`? |

They run, like the visual-regression gates, on a visual unit's main code file and read the spec beside
it; sections and columns are read in any language of the catalog.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: leaves without a verdict |
| the spec | `<Unit>.spec.md` beside it, with a header `code:` | none, or none with a code | this unit: leaves without a verdict, saying why |
| a validation | a row of `Validations` or `Presentation validations` whose first cell carries a unit code | — | this unit |
| a transition | a row of `State Flow` whose trigger names validation codes | — | this unit |
| an error | a row of `Errors / Failures` whose first cell carries a unit code | — | this unit |

## Effects

| Effect | Description |
| --- | --- |
| `VTRST-B01` | A node that is not code, a code file with no spec beside it and a spec with no code leave both gates without a verdict. |
| `VTRST-B02` | `validation-transitions` fails naming each validation — of `Validations` and of `Presentation validations` — that no State Flow row names as its trigger; a spec with no validation leaves without a verdict. |
| `VTRST-B03` | A transition counts only when its From and To are states the spec registers; one from or to an unknown state is named with both ends. |
| `VTRST-B04` | A validation whose row says `@no-state: <reason>` is not asked; one that says `@no-state` with no reason is still asked, and named as an exemption with no reason. |
| `VTRST-B06` | A section is found at any heading level — `### Validações` nested under `## Rules` — and ends at the next heading of its level or above. |
| `VTRST-B05` | `error-message-declared` fails naming each error that cites no message code of the unit, and each that cites codes the Messages section does not catalog; `@no-message: <reason>` on its row exempts it; a spec with no error leaves without a verdict. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VTRST-E01` | The spec beside the code file cannot be read, or there is none. | No verdict, saying there is no spec beside the unit's main file. | A part of the unit has no spec; failing it would charge every file of a screen. |
