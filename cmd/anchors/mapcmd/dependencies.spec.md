<!-- @anchors
  code: DPNDS
  updated_at: 2026-11-15
  layer: comando
-->
# Dependencies — everything a file depends on, the other way of impact

> **Code**: `DPNDS`

## Overview

`impact` answers who depends on a file. Its other end is this command: everything the file
depends on — what it imports or declares `@dep` of, the writers of what it reads, what it
captures or realizes, and the specs, features and guides that rule it. It walks the same ties,
each from the side that depends to the side depended on.

By default only the direct dependencies show. `--propagate indirect` goes one level further and
`--propagate cascade` to the end of the chain, each item past the first level with the path that
ties it to the file; an item that depends on many others says how many it pulls, and `--stop-at
<glob>` lists it and goes no further from it.

The command is a query: it opens no issue and changes nothing.

## Effects

| Effect | Description |
| --- | --- |
| `DPNDS-B01` | Without `--propagate` the direct dependencies show, each with how the file depends on it; `indirect` and `cascade` go further, each item with its path back to the file, and a `--stop-at` glob lists what it matches and goes no further from it. One or more files are asked about at once. (`newDependenciesCmd`) |
| `DPNDS-B02` | `--format files` prints the dependencies one per line, and `--format json` the origins and the dependencies with level, relation and path. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file arguments | one or more paths relative to the root, relative to the working directory, or absolute | a file that is not a node of the map | `impact`'s reading of the query (`readTiesQuery`): fails naming the resolved path |
| the flags | `--propagate indirect|cascade`, `--format text|files|json`, `--stop-at` globs | any other value | `readTiesQuery`: fails naming the value (MPCTI-E03) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DPNDS-X01` | The query leaves the map unchanged and opens no issue. | It is read before deciding what to do. |

## Errors

none — what the command cannot answer (a file outside the map, a missing map, an unknown flag value) is refused by the query reading it shares with `impact` (MPCTI-E01, MPCTI-E02, MPCTI-E03)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
