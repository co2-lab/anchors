<!-- @anchors
  code: PRRPR
  updated_at: 2026-09-26
  layer: config
-->
# ProjectRootResolution — the project root a command works on

> **Code**: `PRRPR`

## Overview

Every command works on a project root: the directory that holds the project's configuration
file. The default of the root option used to be the working directory, and every command
failed when run from a subdirectory — the report came from a real session, where a `cd` left
by one command made the next fail with a message that showed a path nobody asked for and did
not say the cause was the working directory. Repository tools walk up to the root precisely
because people do not stay in it.

This unit resolves the root. With no root given, it walks up from the working directory to
the nearest directory holding the configuration file — the nearest, so a project nested in
another is its own. With none found, it answers with the directory it started from, and the
caller then fails with the ordinary "no configuration" message, which is the right behaviour
outside a project. A root given explicitly is respected as given: whoever points at a
directory is saying where they want to work, and walking above it would ignore the
instruction.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the starting directory | any directory path, relative or absolute | — | this unit: with no project above it, the start is returned unchanged |
| the root option | empty or `.` for "find it", or an explicit directory | — | this unit: only empty and `.` walk up; anything else is made absolute and kept |

## Effects

| Effect | Description |
| --- | --- |
| `PRRPR-B01` | The project root is the nearest directory, from the start upwards, that holds the configuration file; a project nested in another is its own. |
| `PRRPR-B02` | With no project above the start, the start comes back unchanged, relative if it was relative. |
| `PRRPR-B03` | With no root given (empty or `.`), the root is found by walking up from the working directory (`AbsRoot`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRRPR-I01` | The resolved root is always absolute, even outside any project. | resolves the default and an explicit relative root from a directory with no project above it |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRRPR-X01` | An explicit root is made absolute and never walked above, even when a directory above it holds a project. | An explicit root is an instruction about where to work; replacing it with an ancestor would silently act on another directory. |

## Errors

none — not finding a project is answered with the starting directory, and the failure to load its configuration belongs to the caller.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `DefaultFile` | config — the name of the configuration file looked for |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
