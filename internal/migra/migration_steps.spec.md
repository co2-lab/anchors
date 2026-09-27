<!-- @anchors
  code: MGSTM
  updated_at: 2026-09-26
  layer: apoio
-->
# MigrationSteps — the registered steps, one per format, take any project from format 1 to the current format

> **Code**: `MGSTM`

## Overview

The format files are the CONTENT of the migration chain: each one registers the single step that
produces its format, with one line saying what changed and the renames that step performs, scoped by
file name. Together they must take a project from format 1 to the binary's current format without a
hole, and each keeps its own reasons, so that writing format N+1 is a new file that touches none of
the previous ones.

The steps, and why each exists:

- Format 2 turns the identifiers born in Portuguese into English. In the map, the keys for who
  generated it, the declared code and — the costliest to lose — the judgment stamps, whose loss would
  make the check redo every judgment. In the configuration, the key for the optional triad edges.
  It also renames the Portuguese gate names wherever they are stored as values: in the `gate` of the
  map's judgments and in the `name`, `id` and `check` of the configuration's gate declarations. The
  read side used to normalise those names through an alias table while the write side did not, and a
  project ended up with the same gate stamped under two names.
- Format 3 renames the eight gate names format 2 left behind, in the configuration's `name` and
  `check` only. The map's `gate` and the declaration `id` are left out on purpose: none of those
  eight has a default gate, so no project declares them by id nor stamps judgments with them, and
  renaming where the value cannot exist would only point at gates that are not declared.
- Format 4 renames four configuration keys whose names lied about what they hold: a switch that read
  like a mode, a list that read like a boolean, a key without the subject it constrains, and a policy
  that read like the data. Only the configuration: none of them appears in the map.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the files a step renames in | the map (`anchors.graph.yaml`) and the configuration (`anchors.yaml`), by file name | any other file | the file migrator applies each rename only to the file named for it |
| the values a step renames | whole gate names stored under the keys the step names | a gate name inside prose or a comment | the file migrator matches whole values under the named key only |

## Effects

| Effect | Description |
| --- | --- |
| `MGSTM-B01` | Format 2 renames the map keys `gerado_por`, `code_declarado` and `julgamentos` to `generated_by`, `code_declared` and `judgments`, and the configuration key `trinca_opcional` to `triad_optional`. |
| `MGSTM-B02` | Format 2 renames the Portuguese gate names to their English names as the value of `gate` in the map and of `name`, `id` and `check` in the configuration. |
| `MGSTM-B03` | Format 3 renames the eight remaining Portuguese gate names — among them `regra-implementada` to `rule-implemented` and `header-conforme` to `header-valid` — as the value of `name` and `check` in the configuration. |
| `MGSTM-B04` | Format 4 renames the configuration keys `auto_judgment`, `triad_optional`, `requires_code` and `rule_marking` to `enable_auto_judgment`, `optional_triad_edges`, `sections_require_code` and `rule_marking_policy`. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MGSTM-I01` | The registered steps form an unbroken chain from format 1 to the current format: each format has exactly one step, and every step says what it changed. | asks the chain for the steps from 1 to the current format and checks one step per format, each with its reason |
| `MGSTM-I02` | A key renamed by two formats ends under its latest name: a configuration on format 1 with `trinca_opcional` ends with `optional_triad_edges`, keeping its value. | migrates such a configuration from format 1 to the current format through the real steps |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MGSTM-X01` | Format 3 does not rename the map's `gate` values nor the configuration's `id` values. | None of those eight names has a default gate, so no project declares them by id nor stamps judgments with them; a rename there would point at undeclared gates. |
| `MGSTM-X02` | Format 4 renames nothing in the map. | The four keys are policy declarations, which live only in the configuration; the map holds measurements. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |

none — the format files only declare data registered at load time; a hole in the chain or a file
that cannot be migrated is refused by the step chain and the file migrator, which own those
failures.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/migra/step.go` | `Register`, `Step` | apoio — the registry each format file adds its step to |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
