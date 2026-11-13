<!-- @anchors
  code: STRQS
  updated_at: 2026-11-13
  layer: gate
-->
# StatesRequired — every spec of a layer registers the states the project requires of it

> **Code**: `STRQS`

## Overview

A spec registers the states its author thought of, and every state gate confronts only those.
Measured in MIF: nine released defects were a screen state nobody wrote — an error shown as an
empty list, a not-found that spun forever, an offline screen with no way back.

The project names the states each layer owes, by the layer a spec's header declares:

```yaml
required_states:
  screen: [loading, error, empty, offline]
```

Each spec of that layer registers a state by each name — a line carrying the state's code
(`<CODE>-S03`) and the name, its heading or its row —, or says why not, on a line of its own:
`@no-state: offline — the screen reads only local data`. The names are the project's: Anchors
proposes none.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the required states | `required_states:`, layer → names | — | the configuration |
| a registered state | a line of the spec with a state code of the unit (`CODE-S01`, the letter of the project's State rule type) naming the state, without case | the name on a line with no state code | this unit |
| a waiver | `@no-state: <name> — <reason>` | a name with no reason | this unit: it waives nothing |

## Effects

| Effect | Description |
| --- | --- |
| `STRQS-B01` | A spec of a layer with required states that registers none by one of the names, and does not waive it with a reason, fails, naming the layer and the missing states; registered or waived, it passes. (`checkStatesRequired`) |
| `STRQS-B02` | A project that requires no state, a spec of a layer with none required, or an artifact that is no spec, is skipped. |
| `STRQS-B03` | A state is registered where the spec defines it, by its name: a heading that opens with the state's code (its whole text), a table row whose first cell is the code (that cell and the next), a list item that opens with the code (up to the first spaced dash or period). A line that only cites a state's code registers no name. (`stateNames`) |

## Errors

none — the gate reads only the spec's text; a spec with no code defines no state, and what it does not register is a verdict, not a failure
