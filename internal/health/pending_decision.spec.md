<!-- @anchors
  code: PNDCP
  updated_at: 2026-09-26
  layer: infra
-->
# PendingDecisions — the doctor lists the specs that still hold open decisions, the heaviest first

> **Code**: `PNDCP`

## Overview

The doctor is the "give me the overview" command, and an open decision did not show in it. The hierarchy was
upside down: what is still to be MEASURED (a missing mutation signal) was listed as a point of attention,
while a rule still waiting for a decision lived alongside "0 points". Yet the open decision is what most needs
to survive the session: it depends on a person and can take weeks.

This unit reads every spec the map knows and counts its open decisions with the same rule the check uses, so
the doctor and the check never disagree about what is open. Each spec with at least one open decision becomes
one warning carrying its count, and the list starts with the spec that holds the most, because the list
exists to be TAKEN to whoever decides, and that person needs to know where to start.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a graph whose spec nodes point at files under the root | a nil graph, which gives nothing; nodes of other kinds, which are not read | this unit |
| the spec file | a readable spec, with or without an open decisions section | a spec the map knows but the disk lost, which another doctor finding reports | this unit: it skips the spec |
| the configuration | the project's section titles and rule letters, or nil for the defaults | — | the open decisions counter of the check, which accepts nil |

## Effects

| Effect | Description |
| --- | --- |
| `PNDCP-B01` | A spec with open decisions gives one warning `decisao-pendente` on that spec, carrying how many decisions are open. |
| `PNDCP-B02` | A spec that closed its open decisions section with "none" is not pending and gives nothing. |
| `PNDCP-B03` | The findings are ordered by count, the spec with the most open decisions first. |
| `PNDCP-B04` | A nil map, or a map with no open decision anywhere, gives no finding; nodes that are not specs are not read. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PNDCP-I01` | The count the doctor reports for a spec is the count the check's open decisions rule reads from the same text. | the doctor's count for a spec is compared with the check's counter over the same file |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PNDCP-X01` | The unit does not decide what an open decision is: it delegates the count to the check's rule. | Two readings of the same section would let the doctor say "none pending" about a spec the check blocks, or the reverse. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PNDCP-E01` | A spec node's file cannot be read. | The spec is skipped silently, and the other specs are still reported. | The map knowing a file the disk lost is already the doctor's `no-fantasma` finding; repeating it here would report one problem twice. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/open_questions.go` | `OpenDecisions` | gate — the rule that counts open decisions |
| DEP2 | `internal/mapx/model.go` | `Graph`, `KindSpec` | mapa — the specs of the project |
| DEP3 | `internal/config/config.go` | `Config` | config — section titles and rule letters for the count |
| DEP4 | `internal/i18n/i18n.go` | `T` | apoio — the localized finding text |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
