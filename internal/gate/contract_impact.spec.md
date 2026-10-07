<!-- @anchors
  code: CTRIM
  updated_at: 2026-10-07
  layer: gate
-->
# ContractImpact — a changed field names the rules that use it, and their tests

> **Code**: `CTRIM`

## Overview

A spec's file revision says the spec changed, not WHAT changed: editing one field made every
relation of the spec stale, and nothing told which rules read that field. The rule-use sections say
what each rule reads, so a field whose row differs from the last commit names its rules — those of
this spec that use it, and those of the specs that depend on this unit's code and use a field of
that name —, and the tests that cite those rules. The change is read from git, so nothing new is
kept in the map. The same answer feeds the test selection: a test of an affected rule runs even
when its own file did not move.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the spec | a spec in a git repository | a spec with no committed version (new, or no repository) | this unit: no impact |
| the fields | rows of any table outside the rule-use sections, named by their first cell — a field, or a code a rule uses (a state) | `DEPn` as a first cell | this unit: not a name a rule uses |

## Effects

| Effect | Description |
| --- | --- |
| `CTRIM-B01` | A field whose row differs from the spec at HEAD, or that HEAD had and the spec no longer has, is changed; a field only the new version has is not. (`ContractImpacts`, `Impact`) |
| `CTRIM-B02` | A changed field names the rules that use it — by name or first segment — in this spec and in the specs that depend on the code it governs, and the test files whose titles cite those rules; a changed field no rule uses names nothing. |
| `CTRIM-B03` | `contract-impact` is a divergence with each changed field, its rules and its tests; with no impact it passes, and a node that is not a spec is skipped. |
| `CTRIM-B04` | The test files every impacted rule reaches, across the specs with uncommitted changes, are listed for the test selection, which adds those its suite runs. (`ImpactedTests`) |
| `CTRIM-B05` | The rules a revision added since HEAD names in `Revises:` or `Checked:` are answered: an impact whose rules are all answered is not reported, and one with a rule nobody answered still is. The impact lives only while the change is uncommitted, and the change's own revision is where whoever changed the field says they looked. (`acknowledgedRules`) |
| `CTRIM-B06` | A rule is answered where it lives: by a revision added since the last commit to its own spec — the changed one or a spec that reads its data, every revision of a spec not yet in git counting as added —, naming it by its short or full code; or by one added to the changed spec naming it by its full code. A short code answers only its own spec's rule, never another unit's rule of the same letter and number. (`acknowledgedFullCodes`) |

## Errors

none — no repository, no committed version, or no change is no impact, not a failure

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gitmeta/gitmeta.go` | `AtHead`, `HasUncommittedChanges` | infra — the spec at the last commit |
| DEP2 | `internal/gate/rule_uses.go` | `ruleUsesOf`, `cellsOf` | gate — what each rule reads |
| DEP3 | `internal/gate/project_tests.go` | `projectTests` | gate — the tests and their titles |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
