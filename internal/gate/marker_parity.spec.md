<!-- @anchors
  code: MRPRM
  updated_at: 2026-09-26
  layer: gate
-->
# MarkerParity — the same rule has to appear at BOTH ends that fulfil it

> **Code**: `MRPRM`

## Overview

Confronts a rule that lives in TWO places at once: **the promise at one end, and the
fulfilment at the other — are they still the same rule?**

The defect class is the mapping that comes undone on one side only. A rule split across
two ends cannot be checked by any per-file gate: each side, looked at alone, is
impeccable. What breaks is the RELATION, and it disappears without leaving an error.

The case that motivated it, measured (MIF, `EXSC-Q01`): the data-deletion page LISTS
what will be erased in each scope, and the backend ERASES. The two lists were born
together and nothing binds them. If a scope starts erasing more (or less), the page goes
on showing the old version — and the data subject consents on the basis of it. That is
informed consent over the exercise of a right, so the mismatch is not cosmetic.

**This gate is GENERIC on purpose, and that is what separates it from every neighbour.**
It is not canonical: a project may declare SEVERAL instances of it, each with its own
`id`, its own `marker_prefix` and its own `marker_scopes`. What it confronts is a
decision of the project — WHICH rules live at two ends — and that decision does not fit
a universal catalogue. Every other gate in the framework asks a question the framework
itself wrote; this one asks the question the project wrote.

The suffix after the prefix is the NAME of the rule, and it is the name that pairs the
ends. **Two markings on the same side do NOT satisfy the gate:** `marker_scopes` exists
precisely because counting without looking WHERE would let through the very case the
gate exists to catch.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the gate declaration | a `marker-parity` entry of the project's Structure, with its own prefix and id | — (the gate does not choose which rules live at two ends) | the project, in `anchors.yaml`: the decision is its own, not the framework's |
| `marker_prefix` | any non-empty string | an empty prefix, which would match every marking in the tree | this unit: an absent prefix returns Pending, never Pass |
| `marker_scopes` | zero or more glob patterns over the tree | — (no scope is the weaker counting mode, not an error) | this unit: with no scope it falls back to `marker_count` |
| `marker_count` | how many occurrences EACH rule needs | zero with no scopes declared either, which measures nothing | this unit: neither count nor scopes returns Pending |
| the scanned tree | the project root, walked for text files by extension | binaries, and the directories of the ignore list | this unit: an unreadable directory does not bring the whole walk down |

## Effects

| Effect | Description |
| --- | --- |
| `MRPRM-B01` | A rule marked at both declared scopes passes: the mapping still has its two ends. |
| `MRPRM-B02` | A rule marked at one scope and missing from the other fails, and the verdict NAMES the scope that was left empty. |
| `MRPRM-B03` | Two markings on the SAME side do not satisfy the gate — they add up to two and would hide exactly the mismatch it exists to catch. |
| `MRPRM-B04` | A rule that survives at one end after being removed from the other fails, and the verdict names the leftover rule. |
| `MRPRM-B05` | TOTAL absence of the prefix is not approval: it is almost always a typo in the declaration, and returns Pending asking to check `marker_prefix`. |
| `MRPRM-B06` | A declaration with no `marker_prefix` returns Pending: with no prefix there is nothing to confront. |
| `MRPRM-B07` | With no scopes declared, the ruler is the COUNT — the weaker mode, which the gate documents as such. |
| `MRPRM-B08` | The marking crosses LANGUAGE: the mapping between a `.ts` page and a `.go` handler is the same mapping. |
| `MRPRM-B09` | Directories of the ignore list, `node_modules` among them, never count towards parity. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MRPRM-I01` | The failing verdict NAMES what is missing — the rule and the scope left empty. A gate that fails without saying what transfers the diagnostic work to whoever reads it, and here the whole point is that neither side looks wrong on its own. | fails with one end missing and verifies the scope's name appears in the verdict |
| `MRPRM-I02` | Nothing measured is never approved. Missing prefix, missing count and scopes, and total absence all return Pending, not Pass: approving without having looked would make the gate look vigilant while watching nothing. | runs the three cases and verifies none returns Pass |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MRPRM-X01` | Does not read the CONTENT of what each end does. | It confronts that the marking is present at both ends, not that the two implementations agree. Judging whether the page's list and the handler's list say the same thing is semantics, and semantics is not deterministic; presence is. The gate separates "only one end" from "both ends", which is the failure that has no other watcher. |
| `MRPRM-X02` | Does not decide WHICH rules live at two ends. | That is a decision of the project, declared in its Structure. A gate that invented parities would charge mappings nobody committed to — and the whole reason this gate is generic is that the catalogue of such rules belongs to the project, not to the framework. |
| `MRPRM-X03` | Does not read files whose extension is not in the text list. | The cost of erring low here is not seeing a marking in an exotic place; the cost of erring high is reading megabytes of image on every walk. The marking lives in code and documentation, and that is what the list covers. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Gate` | core — prefix, count and scopes all come from the gate the project declared |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — the verdict has to name the missing end in the reader's language |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
