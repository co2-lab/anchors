<!-- @anchors
  code: BLGTN
  updated_at: 2026-09-26
  layer: gate
-->
# Obligations — the duties in force, resolved from packs and config, and their status across the project

> **Code**: `BLGTN`

## Overview

A cross-cutting obligation is a duty that lives outside the unit: a node whose header carries a
trigger attribute must appear in the files the duty names (the account-deletion script, the audit
log). The obligation-honored gate judges one node against one duty. This unit holds the two pieces
around that gate, which must agree with it.

The first resolves which duties are in force. A project declares duties inline in its configuration
and may adopt packs, sets of duties that come from a norm (a privacy law, for instance). Pack duties
come first and inline ones last, so that a local declaration, the most specific one, is the last word.
A pack duty carries its source into its reason, so the gate's message cites the norm instead of only
asserting the duty. A pack that fails to load is a configuration error and is said on the error
output; it drops only its own duties, never the project's inline ones, because silencing it would turn
a whole set of duties into nothing with a green report. The resolution is read once per project root,
because a full check asks for it for every node.

The second is the compliance report. Its unit is the duty, not the file: "42 nodes are subject, 41
comply" answers an auditor, "file X violates" repeated 42 times does not. For each duty it counts the
subjects, how many fulfil it, how many declared it as acknowledged debt, how many waived it with a
reason, and names the ones that do not comply. It judges each node through the same path as the gate,
because two implementations of one rule diverge, and the wrong one would be the report, where people
trust without checking.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | inline duties and the adopted packs, or no configuration at all | — | the config loader validates the fields; this unit reads what arrives |
| the packs | pack files the loader can read and resolve | a pack that does not load | this unit: the error is reported and the inline duties remain |
| the map | nodes whose files are under the project root | a node whose file cannot be read | this unit: that node is not a subject |

## Effects

### The duties in force

| Effect | Description |
| --- | --- |
| `BLGTN-B01` | The duties in force, as seen from outside the package (`ObligationsInForce`), are the resolved list the gate uses, pack duties included. |
| `BLGTN-B02` | With no configuration there is no duty; with no pack adopted, the duties in force are the inline list as declared. |
| `BLGTN-B03` | The pack duties come first and the inline duties last, and each pack duty keeps its trigger, its target files and how it is identified. |
| `BLGTN-B04` | A pack duty's reason gets its source appended in parentheses: the authority and the article when both exist, either one alone otherwise; with no reason, the source alone is the reason. |
| `BLGTN-B05` | A pack that fails to load is reported on the error output, in the project's language, and the inline duties remain in force. |
| `BLGTN-B06` | The pack duties are read once per project root and pack set (the adopted packs, their values and the jurisdictions); later calls with the same root and set are served from that first read. |
| `BLGTN-B12` | A call with the same root and another pack set reads that set, instead of being served the list of the first. |

### The compliance report

| Effect | Description |
| --- | --- |
| `BLGTN-B07` | The report (`EvaluateObligations`) gives one status per duty, in the order the duties are handed, each with its name and target files. |
| `BLGTN-B08` | A node is a subject of a duty only when its header carries the duty's trigger attribute; a duty with no trigger has no subject. |
| `BLGTN-B09` | A subject that waives the duty with a reason counts as fulfilled and as waived. |
| `BLGTN-B10` | A subject whose judgement is Pending, because it declared the duty as acknowledged debt, counts as debt. |
| `BLGTN-B11` | Every other subject that does not pass is named as missing, and the missing list is sorted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BLGTN-I01` | Every subject is counted exactly once: fulfilled, debt, or missing add up to the subjects. | a project with two fulfilling nodes, one waived, one in debt and two missing gives 6 subjects = 3 + 1 + 2 |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BLGTN-X01` | The report evaluates only the duties it is handed: it does not resolve the packs of the configuration it receives. | Who resolves the duties is the resolution above; resolving twice would let the report and the gate see different lists. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BLGTN-E01` | REF[BLGTN-B05]: a pack that fails to load is the configuration failure B05 reports on the error output while keeping the inline duties | — | — |
| `BLGTN-E02` | A node of the map whose file cannot be read. | It is not counted as a subject of any duty. | A node that cannot be read cannot be shown to carry the trigger; counting it as missing would accuse a file that is not there. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/obligation_honored.go` | `checkObligationHonored`, `headerHasAttr`, `waiverFor` | gate — the judgement of one node against one duty |
| DEP2 | `internal/pack/pack.go` | `LoadAll` | infra — the packs adopted by the project |
| DEP3 | `internal/config/config.go` | `Config`, `Obligation` | config — the inline duties and the adopted packs |
| DEP4 | `internal/mapx/model.go` | `Graph` | mapa — the nodes the report walks |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
