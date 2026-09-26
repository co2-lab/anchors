<!-- @anchors
  code: RLIMR
  updated_at: 2026-09-26
  layer: gate
-->
# RuleImplemented — a spec catalogues rules, and the code shows it realized them

> **Code**: `RLIMR`

## Overview

Confronts the spec against the code in the direction that was missing: **did the spec end up
talking to itself?**

It is the inverse of the gate that validates references. That one checks that the codes CITED by
the code exist in the spec; this one checks that the rules DECLARED in the spec got an
implementation. Without it, a spec can declare five new rules and the code gain
not a single line — with all gates green, because the spec exists, the code exists, and the
two reference each other through the header.

Measured: an interface spec gained five rules and 98 lines, and the corresponding file
had ZERO occurrence of the subject. Half of the delivery was dead code declared as
done, and none of the 26 gates asked. The defect only showed up when someone read spec and
code in the same pass.

**The ruler is the declaration, not the guesswork.** Demanding every rule be marked would be false by
construction — measured against 592 units, it would yield 3.121 findings, and not even the well-made units
would pass: among those that mark the code, none marks 100%. The reason is a good one: a constraint ("the
unit does NOT do Y") is satisfied by the ABSENCE of code, and absence has nowhere to receive a
mark. But "at least one" does not serve either — it separates those who implemented from those who did not
implement and says nothing about the other fifteen rules. So whoever writes the spec
DECLARES, rule by rule, whether it has code.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the spec | with catalogued rules, or with none at all | — (a spec without a rule is a case, not an error) | this unit: with no rule there is nothing to charge |
| the linked code | whatever the realization edge points to, or none | — | this unit: with no linked code the subject does not exist |
| the per-rule waiver | the mark on the rule's line, with a written reason | a bare mark, without a reason | this unit: a waiver without a why does not waive |
| the marking requirement | declared in the project's Structure | — (omitted counts as not required) | the Structure: while the project does not declare it, the debt is PENDING and not a failure |

## Effects

| Effect | Description |
| --- | --- |
| `RLIMR-B01` | A spec whose rules do not appear in the code is ACCUSED, and the verdict names which ones were left without realization. |
| `RLIMR-B02` | A rule waived with a written reason settles the account: the declaration counts as the answer, and what it waives stops being charged. |
| `RLIMR-B03` | A unit that predates the practice becomes PENDING, not a failure — while the project does not declare that it requires the marking. |
| `RLIMR-B04` | Once the requirement is declared in the Structure, the pending item becomes a failure: it is the act of saying "the migration ended here". |
| `RLIMR-B05` | A spec with no linked code is not this gate's subject: without the piece on the other side there is no confrontation to make. |
| `RLIMR-B06` | The waiver may name the rule it covers, or count for all of them when it names none. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RLIMR-I01` | Requiring the marking never punishes whoever already marks. Whoever did the work before the requirement cannot fail for having done it. | declares the requirement over a unit that marks everything and verifies that it passes |
| `RLIMR-I02` | Identity survives the rename: code marked with the previous name keeps counting. Losing the mark in a rename would turn identity stability into new debt. | marks the code with the old name and verifies that the rule stays recognized |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RLIMR-X01` | Does not judge whether the implementation is RIGHT — only whether it exists and declares itself. | Whether the code fulfills what the rule says is judgment, and judgment belongs to another class of gate. The ruler here is deterministic: the mark exists, or the waiver exists with a reason. |
| `RLIMR-X02` | Does not require a mark on EVERY rule. | A constraint is satisfied by the absence of code, and absence has nowhere to receive a comment. Charging the 3.121 occurrences that the naive ruler would produce would train the team to ignore the whole list. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `RLIMR-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph` | core — the linked code is reached by the edge, not by name convention |
| DEP2 | `internal/config/config.go` | `Config` | core — the marking requirement is declared in the Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
