<!-- @anchors
  code: DTAUI
  updated_at: 2026-09-26
  layer: comando
-->
# Audit — the dossier of everything pending on one file, for fixing it in one pass

> **Code**: `DTAUI`

## Overview

`anchors audit <file>` gathers in one place everything pending on a target: the quality gates that apply to its kind, which is what `check` would run, and the systemic findings of `doctor` that cite it. It exists for parallel sweeping. One agent takes one file, reads the whole dossier, fixes every pending item at once — header, spec, identity, coverage — and moves on, instead of opening the same file again for each gate.

By default the scope is the file alone. With `--impact` the scope grows to the impact path of the target: the rest of its unit and what it propagates to or is validated by, so a single worker can fix the whole triad.

The dossier separates what is actionable from what is only shown. A failing gate and a warning finding are work for whoever holds the file; a gate waiting for judgment or data, and an informational finding, are shown so they are not forgotten but are not counted as items to fix.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the target | one path of a file the map knows | a file the map does not know | this unit: it refuses a file outside the map |
| the project | a root with a loadable configuration and a built map | a root without configuration or without a map | this unit: it fails naming what it could not load |
| the gate results and doctor findings | whatever the gates and the doctor return for the scope | — | the gates and the doctor; this unit only filters and prints |

## Effects

| Effect | Description |
| --- | --- |
| `DTAUI-B01` | A target the map does not know is refused with "file … is not in the map", and nothing is audited. |
| `DTAUI-B02` | Without `--impact` the scope is the target alone, and the header reads "audit: <target> — the file". |
| `DTAUI-B03` | With `--impact` the scope adds the nodes the target propagates to and the nodes that validate it, and the header reads "the unit (N nodes on the impact path)". |
| `DTAUI-B04` | A gate that passed or was skipped is dropped. A failing gate is marked ✗, a pending one ~ and one awaiting judgment ⏳, each with only the first line of its detail. |
| `DTAUI-B05` | A doctor finding is shown only when its subject is in the scope, marked ⚠ for a warning and ℹ for information. |
| `DTAUI-B06` | Only a failing gate and a warning finding count toward "N actionable pending item(s)"; pending, judgment and informational lines are shown but not counted. |
| `DTAUI-B07` | When nothing is left to show, the dossier says "✓ nothing pending". |
| `DTAUI-B08` | The target's block prints first, marked ●; every other node in scope prints after it, marked ○ and "(impact)". |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DTAUI-I01` | Nothing about a node outside the audited scope reaches the dossier: neither a gate line nor a doctor finding. | audits a file with other nodes and out-of-scope findings present, and checks that none of them is printed |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DTAUI-E01` | The root has no loadable configuration. | The command fails with "load config: …". | Without the configuration there are no gates to run; an empty dossier would read as a clean file. |
| `DTAUI-E02` | The root has no loadable map. | The command fails with "load map: … (run `anchors map build`)". | The scope and the doctor both read the map; the message names the command that creates it. |
| `DTAUI-E03` | REF[DTAUI-B01]: a target outside the map is the failure B01 answers: the command refuses it by name | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/gate.go` | `RunWithConfig` | gate — the gates over the scope |
| DEP2 | `internal/health/health.go` | `Diagnose` | the doctor's systemic findings |
| DEP3 | `internal/mapx/impact.go` | `AnalyzeImpact` | the impact path for `--impact` |
| DEP4 | `cmd/anchors/common/path.go` | `RelTo`, `NodeExists` | the target resolved against the root and the map |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
