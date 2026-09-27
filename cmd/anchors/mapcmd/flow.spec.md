<!-- @anchors
  code: FLWOX
  updated_at: 2026-09-27
  layer: comando
-->
# Flow — build, draw and navigate the work flows kept in the map

> **Code**: `FLWOX`

## Overview

A work flow is work driven by shape instead of memory. Actions declare what they do and the
results they offer; a flow fits actions together and says where each result goes. This unit
is the command side of that model, with three subcommands.

Building reads the declared flows and actions and writes the flow graph into the existing
map, next to the nodes, counting steps, results and links. Its finding is the result that no
flow routes: whoever receives it would improvise. That is reported as a warning, because a
project may legitimately not cover every branch; what it must not do is not know.

Navigating answers, for one step, where the work can go from there: the valid exits, the
condition of each, and the suggestion the result carries. It also makes the two defects of a
flow visible at the point of use: a step that fits an action nobody wrote, and a step with no
exit that does not declare itself terminal.

Drawing derives the picture of a flow at the moment it is asked, from the graph, in text or
as a mermaid diagram. It is never stored, so it cannot age against the flow it describes.

## Effects

| Effect | Description |
| --- | --- |
| `FLWOX-B01` | Building a project with no flow declared says so and is not an error. |
| `FLWOX-B02` | Building writes the flow graph into the existing map, keeps the map's nodes, and reports how many steps, results and links it has. |
| `FLWOX-B03` | Building names every result no flow routes, as a warning: the build still succeeds. |
| `FLWOX-B04` | Navigating a step (its code is matched regardless of case) prints the step, the action it fits, and each valid exit with its destination's title, its condition, and the suggestion carried by the result. An exit with no result is labelled with a dash. |
| `FLWOX-B05` | Navigating a step that fits an action no action file declares flags it. |
| `FLWOX-B06` | Navigating a step with no exit says the work ends when the step is terminal, and says whoever arrives is stuck when it is not. |
| `FLWOX-B07` | The text drawing lists the steps in the file's order, never draws a result as a step, marks the terminal steps, and lists each step's exits sorted by result; an exit that no result triggers is labelled `sempre`, followed by its condition. |
| `FLWOX-B08` | The mermaid drawing turns step codes into valid identifiers, replaces double quotes in labels, puts the fitted action on a second line with a native line break, draws terminals in the stadium shape, labels each exit with the short name of its result (or its condition), and highlights the entry and the terminals. |
| `FLWOX-B09` | The mermaid direction accepts top-down, left-right, right-left and bottom-top in any case; any other value falls back to top-down. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared flows | flow and action files the flow builder reads | a project with no flow file | this unit: says so and builds nothing |
| the map | an existing map, to receive or provide the flow graph | a missing map, or a map with no flow | this unit: refuses and names the command to run first |
| the step code | any case of a step code present in the flow graph | a code that is not a step of the graph | this unit: fails naming the code |
| the flow name | a text contained in a flow's name, or nothing for every flow | a text no flow name contains | this unit: fails listing the available flows |
| the direction | TD, LR, RL or BT, any case | any other text | this unit: falls back to TD |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLWOX-I01` | What the build writes into the map is what drawing and navigating read: the drawing is derived from the stored graph, never from a stored picture. | closed cycle: build a flow, then show and navigate it and find the declared steps and exits |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLWOX-X01` | Only building writes the map; drawing and navigating leave it unchanged. | They are queries over the graph; a query that rewrote the map would change what the gates confront. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLWOX-E01` | Building a declared flow when the project has no map. | The build is refused, asking for the map to be built first, and no map is created. | A map created only to hold the flow would have no nodes, and the relational gates would confront a void. |
| `FLWOX-E02` | Navigating a step that is not in the flow graph. | The command fails naming the step. | There is no exit to answer for a step the graph does not have. |
| `FLWOX-E03` | Showing a flow when no flow name contains the given text. | The command fails listing the available flows. | The reader needs the names to ask again. |
| `FLWOX-E04` | Showing or navigating when the map holds no flow. | The command fails asking for the flow build. | The flow is read from the map, and nothing was built into it yet. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/flowx/build.go` | `Build` | apoio — reads the declared flows and actions into a graph |
| DEP2 | `internal/flowx/model.go` | `Next`, `StatesOf`, `Unhandled`, `StateByCode`, `Entry` | apoio — the queries over the flow graph |
| DEP3 | `internal/mapx/lock.go` | `Update` | mapa — the map that stores the flow graph, changed under its lock |
| DEP4 | `internal/config/root.go` | `AbsRoot` | config — the project root |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
