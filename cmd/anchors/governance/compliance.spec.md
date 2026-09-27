<!-- @anchors
  code: CMPLN
  updated_at: 2026-09-26
  layer: comando
-->
# Compliance — the state of each regulatory duty, grouped by the norm that imposes it

> **Code**: `CMPLN`

## Overview

`anchors compliance` answers the question an auditor asks: which duties apply to this project, and where does it stand on each one. Before it, the answer was scattered — `check` said whether ONE file broke ONE obligation, and the state of a whole regime had to be walked by hand.

The report is per DUTY, not per file. The duties come from the adopted packs and from the obligations declared inline in the configuration. Each duty line says how many nodes are subject to it and how many comply, cites the pack's article when there is one, and adds the declared debts and waivers. The duties are grouped under the norm that originates them, which is what makes the report presentable to whoever audits and not only to whoever programs.

A duty with subjects and nobody complying is rarely total violation; it is usually a target that moved. The report says so next to the duty instead of letting the reader conclude the worst. Finally, it lists the embedded packs the project did not adopt: not as advice to adopt them, but so the project can tell "does not apply to me" from "I forgot".

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project | a root with a loadable configuration and a built map | a root without configuration or without a map | this unit: it fails naming what it could not load |
| the adopted packs | pack references the pack loader resolves, with every value they require declared | a pack whose required value is missing | the pack loader rejects it, and this unit fails with that error |
| the duty statuses | what the obligation evaluation returns per duty | — | the obligation evaluation; this unit only groups and prints |

## Effects

| Effect | Description |
| --- | --- |
| `CMPLN-B01` | Each duty is printed under its norm: a pack duty under the pack's authority, or the pack's name when it declares none; an inline obligation under "declared in the project". Norms and the duties inside each norm are sorted by name. |
| `CMPLN-B02` | Each duty line is marked · when no node is subject, ✗ when fewer subjects comply than are subject, and ✓ otherwise. It shows the subject and complying counts, the pack article when there is one, and appends the assumed debts and the waivers when there are any. |
| `CMPLN-B03` | A duty with subjects, none complying and no declared debt gets a "⚠ NONE complies" line naming its targets, warning that a disconnected target looks like total violation. A duty whose only gap is declared debt does not warn. |
| `CMPLN-B04` | With `--verbose`, each duty lists the nodes that do not comply. Without it, when some duty has a gap, the report ends with the hint "(use --verbose …)". |
| `CMPLN-B05` | With no duty in force the report says "No duty declared." and still lists the packs available. |
| `CMPLN-B06` | The embedded packs the project did not adopt are listed, sorted, after the report. A pack adopted by its short name or as `./packs/<name>.yaml` is not offered again, and when every pack is adopted nothing is printed. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CMPLN-I01` | Every duty in force has its line in the report, even one that no node triggers. | a report with a duty whose condition no file carries still shows that duty, marked · |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CMPLN-E01` | An adopted pack is missing a value it requires. | The command fails with the pack loader's error, naming the missing value. | A duty whose target is a placeholder would be evaluated against no file and read as total violation. |
| `CMPLN-E02` | The root has no loadable configuration. | The command fails with "load config: …". | Without the configuration there are no packs or obligations; "No duty declared" would be a false answer. |
| `CMPLN-E03` | The root has no loadable map. | The command fails with "load map: … (run `anchors map build`)". | The subjects are the map's nodes; the message names the command that creates it. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/pack/pack.go` | `LoadAll` | the adopted packs, with their values resolved |
| DEP2 | `internal/gate/obligations_report.go` | `ObligationsInForce`, `EvaluateObligations` | gate — the duties in force and their status per node |
| DEP3 | `internal/initx/packs.go` | `AvailablePacks` | the embedded packs, to list the ones not adopted |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
