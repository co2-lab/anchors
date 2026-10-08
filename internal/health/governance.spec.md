<!-- @anchors
  code: GVOPG
  updated_at: 2026-10-08
  layer: infra
-->
# GovernanceOpportunities — the doctor suggests the canonical gates and settings a project has not adopted yet

> **Code**: `GVOPG`

## Overview

Canonical gates and good settings are born after a project is created, or are simply never adopted. Nothing
fails because of that: the project is only less governed than it could be, and nobody is told. This unit
looks at what the project HAS (the kinds of artefact in its map) and what it DECLARES (its gates and test
suites), and suggests what is missing for what it has.

A project with tests and code is offered the gate that watches evidence freshness. A project with code is
offered the gates that catch leaked secrets, vulnerable dependencies and duplication, and is told when the
secret gate is declared but does not block, since a secret in history has no simple way back. A project with
tests but no suite that produces JUnit results is told its test configuration is suboptimal.

Every finding is informational: these are opportunities, not problems. The same list, cut to its first two
items, is the short hint other commands can show without getting in the way.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a graph of nodes with kinds | a nil graph, which gives nothing | this unit: it returns before looking |
| the configuration | the declared gates (by name, with their blocking flag) and test suites | a nil configuration, which gives nothing | this unit: it returns before looking |

## Effects

| Effect | Description |
| --- | --- |
| `GVOPG-B01` | A map with both test and code nodes and no gate named `evidence-fresh` gives an informational `sugestao-gate` on `evidence-fresh`. |
| `GVOPG-B02` | A map with code and no gate named `no-secret-leaked` gives an informational `sugestao-gate` on `no-secret-leaked`. |
| `GVOPG-B03` | When `no-secret-leaked` is declared but does not block, the finding is an informational `gate-subotimo` on it instead of a suggestion. |
| `GVOPG-B04` | A map with code gives an informational `sugestao-gate` for `dependency-vulnerable` and for `no-duplication`, each when that gate is not declared. |
| `GVOPG-B05` | A map with test nodes where no test suite declares JUnit output gives an informational `config-subotima` on `tests.junit`. |
| `GVOPG-B06` | A nil map or a nil configuration gives no finding. |
| `GVOPG-B07` | `QuickGovernanceHints`: The quick hints are the first two opportunities in order, or all of them when there are fewer, leaving out a catalog suggestion whose gate presupposes a field the project has not declared — the doctor still lists it. |
| `GVOPG-B08` | Every catalog gate that relates to a declared layer and that the project does not declare, beyond the ones above, is an informational `sugestao-gate` on it, saying what it measures, how to declare it, and the field it presupposes when it has one. |
| `GVOPG-B09` | A project with code, a family that knows fallible patterns and no `dialect.fallible_patterns` declared is offered them, as the block to copy under `dialect:`; one that declined (`fallible_patterns: []`) or declared its own is not. (`fallibleYAML`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GVOPG-I01` | Every opportunity is informational: this unit never raises a warning. | a project where every opportunity fires is checked, and each finding's severity is read |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GVOPG-X01` | A gate the project already declares (blocking, for the secret gate) is never suggested again, and a project with JUnit output is not told to add it. | A suggestion the project already followed is noise that trains the team to ignore the doctor. |

## Errors

none — the unit only reads a graph and a configuration already in memory; nil inputs are the empty case of `GVOPG-B06`, not a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
