<!-- @anchors
  code: CFPCD
  updated_at: 2026-09-26
  layer: infra
-->
# CodeFromPath — the most meaningful unique code for a unit, given its file path

> **Code**: `CFPCD`

## Overview

When a code is generated for a file rather than a name, the file name is not always the unit's name.
In `functions/manage-metadata/handler.ts` the unit is "manage-metadata": "handler" is a role inside the
folder, shared by every sibling unit. For such generic file names (handler, index, main, types, route,
routes, schema, config, resource, mod) the parent folder is the identity, and the code is generated from
it. Otherwise the code comes from the file name.

The unit's base name drops the artifact suffixes, so the spec, the feature and the test of a unit all
compress the same name as its code file. The result is always unique against the codes already taken,
using the generator's deterministic collision resolution (`CDGNC`): given the same taken codes the answer
is the same, although which of two colliding units gets the clean code depends on which was created
first.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the unit path | a file path relative to the project root | — | the caller (`anchors new`, `anchors init`) |
| the codes taken | the codes already in use in the project's map | — | the caller |

## Effects

| Effect | Description |
| --- | --- |
| `CFPCD-B01` | The generic file names are handler, index, resource, mod, main, types, route, routes, schema and config, recognised in any letter case. (`IsGenericBasename`) |
| `CFPCD-B02` | A unit whose file name is generic takes its code from its parent folder's name; one with no parent folder keeps its own name. (`GenerateFromPath`) |
| `CFPCD-B03` | The unit's base name drops `.spec.md`, `.feature`, or everything from `.test.` on, and otherwise the extension, so every artifact of a unit compresses the same name. |
| `CFPCD-B04` | A file name that is not generic gives the generator's code for that name when it is free. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CFPCD-I01` | Two units with the same base name get different codes, and the same path with the same taken codes gets the same code. | generates for `models/metadata`, takes it, generates for `repositories/metadata` twice, and compares |

## Errors

none — every path yields a code; collisions are resolved by the generator (`CDGNC-B08`), not failed.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/code/code.go` | `GenerateUnique` | infra — the compression and the collision resolution (`CDGNC`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
