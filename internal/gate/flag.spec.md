<!-- @anchors
  code: FLSCF
  updated_at: 2026-09-27
  layer: gate
-->
# FlagScenarios — the scenarios a feature flag declares are written, complete, cited and tested

> **Code**: `FLSCF`

## Overview

A feature flag multiplies the paths of the code without multiplying the spec. A check of a flag's
value creates two behaviours, and the spec usually describes one of them, or both mixed in a sentence
that does not say which holds when. The cost shows up as silences no other gate sees: in review nobody
knows which branch is current and which is dying, the test covers the ON path and leaves the OFF path
unproven, and the "temporary" flag grows old until removing it is archaeology.

A flag file declares its scenarios in a table: a code with the `G` letter, the condition on the value,
and what holds then. Five gates confront that declaration, each answering one question:

- flag-scenario-grammar: is each condition written in the fixed grammar?
- flag-scenarios-complete: does the flag say what happens when it is absent?
- flag-scenario-exists: does every scenario a spec cites exist?
- flag-scenario-governs: is every scenario cited by some rule?
- flag-covered: does every scenario have a test, and did that test pass?

What is deliberately not read is the flag's real value: Anchors has no access to the flag service,
must not have, and the value changes per user and per minute. The gates confront the declared scenarios.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | a flag file for four gates, a spec for flag-scenario-exists | any other kind of node | this unit: each gate skips a node of another kind |
| the flag's scenarios | a scenario table with codes, conditions and outcomes | a flag file with no scenario | this unit: the flag gates skip it |
| the spec's citations | `@gated-by` followed by a code with the `G` letter | a citation of a code with another letter | this unit: it is not read as a citation |
| the map | the graph with the incoming citation edges, the test nodes and the ingested proven codes | no map | this unit: the gates that need it answer Pending |

## Effects

| Effect | Description |
| --- | --- |
| `FLSCF-B01` | The four gates that confront a flag skip a node that is not a flag, and skip a flag that declares no scenario; the citation gate skips a node that is not a spec. |

### flag-scenario-grammar

| Effect | Description |
| --- | --- |
| `FLSCF-B02` | A flag whose every condition is in the grammar passes. |
| `FLSCF-B03` | A condition the grammar refuses fails the gate, and the message names the scenario, its line and the parse error. |

### flag-scenarios-complete

| Effect | Description |
| --- | --- |
| `FLSCF-B04` | A flag that declares a scenario for the absent value passes. |
| `FLSCF-B05` | A flag with no absent scenario fails, and the message names the waiver that would close it. |
| `FLSCF-B06` | A written waiver for the absent case, with a reason, passes; a bare marker, or one quoted in backticks, does not waive. |

### flag-scenario-exists

| Effect | Description |
| --- | --- |
| `FLSCF-B07` | A spec with no citation of a flag scenario skips; only a citation of a `G` code counts as one. |
| `FLSCF-B08` | A spec whose every cited scenario is declared by some flag of the project passes. |
| `FLSCF-B09` | Citations of scenarios no flag declares fail, and the message names each unknown code once, sorted. |

### flag-scenario-governs

| Effect | Description |
| --- | --- |
| `FLSCF-B10` | Without a map the gate answers Pending, because it cannot know who cites. |
| `FLSCF-B11` | Each scenario needs an incoming citation that names it; the scenarios nobody cites fail the gate, named, and the cited ones are not accused. |
| `FLSCF-B12` | A scenario whose outcome carries a written waiver of governance, with a reason, is not charged; a bare marker does not waive. |

### flag-covered

| Effect | Description |
| --- | --- |
| `FLSCF-B13` | Without a map the gate answers Pending, because it cannot know what was proven. |
| `FLSCF-B14` | Coverage is judged per scenario: a scenario among the flag's ingested proven codes is green, and a flag whose every scenario is green passes. |
| `FLSCF-B15` | A scenario that no test names fails as having no test; a code that appears only in a comment of a test does not count as written. |
| `FLSCF-B16` | A scenario a test names but that is not proven fails with a different message: written but not ingested when no execution was ingested, written and not passing when it was. |
| `FLSCF-B17` | A `@gated-by` citation is read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, a citation of a 7-character scenario is confronted. |
| `FLSCF-B18` | When the project says how its tests are written, a flag scenario is written in a file the source lists tests in only when a test TITLE cites its code; without a source, or in a file the source lists no test in, a code anywhere in the file outside comments counts. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLSCF-I01` | A gate that could not measure never answers Pass: without a map, governance and coverage are Pending. | runs both gates with no map and asserts Pending, not Pass |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLSCF-X01` | No waiver is accepted without a written reason. | The reason is what answers the question somebody asks six months later; a bare marker turns the gate off without saying why. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLSCF-E01` | REF[FLSCF-B03]: a condition the parser refuses is the failure the grammar gate reports, naming the scenario and the parse error | — | — |
| `FLSCF-E02` | The `flags/` folder, or a flag file in it, cannot be read when a spec cites a flag scenario. | Pending, naming the read error and the `flags/` folder — not the "no map" message. | Without the flags the gate cannot say which scenarios exist: Pass would approve a citation never checked, Fail would accuse one that may exist, and blaming the map would send the reader to build a map this gate never reads. |
| `FLSCF-E03` | The project's tests source fails, or answers outside its contract | Fail, naming the source's error | Which flag scenarios have a test cannot be told, and the gate says why instead of answering about tests nobody read |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/flagx/parse.go` | `ParseContent`, `Load`, `ByCode` | scan — the flag file's scenarios |
| DEP2 | `internal/mapx/model.go` | `Graph`, `Node`, `EdgeGatedBy` | mapa — the citation edges and the ingested signal |
| DEP3 | `internal/gate/internal_checks.go` | `codesNamedByTests` | gate — which scenario codes a test body names |
| DEP4 | `internal/i18n/i18n.go` | `T` | apoio — the localized messages |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
