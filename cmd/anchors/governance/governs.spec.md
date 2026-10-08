<!-- @anchors
  code: GVRNS
  updated_at: 2026-10-08
  layer: comando
-->
# Governs — who each guide governs, and how many, read from the map

> **Code**: `GVRNS`

## Overview

`anchors governs` answers "whom does each guide govern, and how many" from the map. With no argument it prints the board of every guide and the number of files it governs directly, which sizes a judgment audit per guide and exposes redundancy: guides that govern the same set are candidates to narrow by tag. With a guide as argument it lists the files that guide governs directly, grouped by kind — the list a batch of judgments is cut from.

"Directly" means the governs edges leaving the guide, not the transitive wave; the full reach of a change belongs to `anchors impact`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the guide | none, or one path, relative to the root or absolute | more than one argument | the command's argument rule refuses it |
| the map | the project's map, or the one `--map` names | a missing or unreadable map | this unit: it fails naming the command that builds the map |

## Effects

| Effect | Description |
| --- | --- |
| `GVRNS-B01` | With no argument, the board lists each guide with its count of governed files, the most governing first, then the total of (guide, governed) pairs. |
| `GVRNS-B02` | With no argument and no governance in the map, it prints "no guide governs anything". |
| `GVRNS-B03` | With a guide, it prints "<guide> governs N file(s):" and the files grouped by node kind, the kinds sorted, each with its count. |
| `GVRNS-B04` | A file that governs nothing is answered with "<file> governs nobody", and the command succeeds. |
| `GVRNS-B05` | The guide argument is made relative to the root before the map is asked, and `--map` reads the map from the path it names instead of the project's default. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GVRNS-I01` | Only governs edges count as governance; any other edge leaving a node adds nothing to its count or to the total. | a map where a spec also specifies code still totals only the governs edges |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GVRNS-E01` | The map cannot be loaded. | The command fails with "load map: … (run `anchors map build`)". | Governance exists only in the map; the message names the command that builds it. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
