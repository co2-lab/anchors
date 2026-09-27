<!-- @anchors
  code: FLBLF
  updated_at: 2026-09-26
  layer: apoio
-->
# FlowBuild — assembling the work-flow graph from the project's flow and action files

> **Code**: `FLBLF`

## Overview

A project writes its work flows as markdown in `flows/*.flow.md`, and the reusable pieces those flows
fit, the actions, in `flows/actions/*.action.md`. An action declares what it does and the results it
offers, and deliberately does not know who comes next; a flow fits actions together and says where each
result goes. This unit reads both kinds of file and assembles one graph of states and transitions (asked
by `FLMDF`).

A state opens a heading section, and everything below it until the next heading belongs to it: its
exits, the piece it fits, the reaction a result suggests, whether it is terminal. The association is
positional, so a flow is written the way it reads. An exit's condition is prose on purpose: whoever works
judges whether it holds; what the graph guarantees is the SET of exits. A terminal state is declared,
never inferred from having no exit, since a state with no exit may be the end of the work or an
oversight. A suggestion is recorded as a suggestion, not a transition: presenting it as an ordinary exit
would make the flow lie about who decides.

The keywords a person writes in a flow (the piece a step fits, the reaction a result suggests) come from
the translation catalog, so a project writes flows in its own language. A project without flows has
simply not declared any, and files are read in a stable order.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a folder with or without `flows/` | — | this unit: no flows means no graph |
| the flow and action files | markdown with coded headings, bullet exits and keyword lines | exits written before any state | this unit: an exit with no state above it is ignored |

## Effects

| Effect | Description |
| --- | --- |
| `FLBLF-B01` | The graph is built from every `flows/actions/*.action.md` and then every `flows/*.flow.md`, each set sorted by file name; a project without flows gives no graph and no error. |
| `FLBLF-B02` | A state is a heading carrying a code, recorded with its title and the file it comes from. |
| `FLBLF-B03` | An exit is a bullet carrying a code in backticks; it belongs to the last state above it, and the rest of the line is its condition. An exit above every state is ignored. |
| `FLBLF-B04` | A line carrying two codes routes a result: a transition from the current step, on the first code, to the second code. |
| `FLBLF-B05` | A state is terminal only when its section declares `@terminal`. |
| `FLBLF-B06` | The piece a step fits is read from its section's "fits" line, written with the keyword of any supported language; a step without one fits nothing. |
| `FLBLF-B07` | The reaction a result suggests is read from its section's "suggests" line and recorded on the state, never as a transition; a result without one suggests nothing. |
| `FLBLF-B08` | The condition of a routed result is the prose of its line before, between and after the two codes, without the codes, their backticks, the arrow (`→` or `->`), the bullet and trailing punctuation. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLBLF-E01` | `flows/` or `flows/actions/` exists but cannot be read as a folder, or a flow or action file in them cannot be read. | The build returns the read error and no graph. Only a folder that does not exist is read as "no flows" (B01). | Skipping it would drop its states from the graph in silence, and the flows routed to them would look dangling. |

A flow's findings (unreachable, unhandled, dangling) are answered by `FLMDF`, not failures of the build.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/i18n/i18n.go` | `AllTranslations` | apoio — the flow keywords in every language (`INCTA`) |
| DEP2 | `internal/mapx/model.go` | `FlowGraph` | mapa — the graph the flows are written into |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
