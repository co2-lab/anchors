<!-- @anchors
  code: IDCND
  updated_at: 2026-09-26
  layer: gate
-->
# IdentityConsistent — a unit's spec identity must match its exposed testID and visual baseline

> **Code**: `IDCND`

## Overview

Confronts a unit's declared identity against the external surfaces where it reappears: **does the spec
code agree with the testID prefix exposed in code and the visual regression baseline filename?**

The spec code is the canonical identity of the unit. Other identity gates verify mere PRESENCE
(whether a spec declares a code, or whether a scenario carries a code), never CONCORDANCE across surfaces.
Under presence-only checking, a single unit could declare a spec code of `BGET`, expose `testID=":bdge-screen"`,
and store its visual regression baseline as `BDEDX-VR-*.png` — three conflicting identities for the exact same
unit — while every pipeline gate remains green because each gate inspects its own surface in isolation.

The REAL COST of this discordance is not aesthetic. The project code dictionary consumed by spellcheck is
AUTOMATICALLY GENERATED from the map. When an acronym exists only inside a testID, it is missing from the
generated dictionary. Spellcheck immediately flags the acronym as a typo, and the natural developer
workaround — adding the rogue acronym to the manual spelling dictionary — CRISTALLIZES the divergence.
The symptom becomes permanent vocabulary while the underlying defect disappears from sight. This exact
failure mode occurred in the reference application when `bdgc` and `bdge` entered the manual dictionary,
motivating the creation of this gate.

The gate exercises crucial DISCERNMENT regarding what does NOT constitute divergence: a component may
legitimately use the identity code of ANOTHER unit as its testID prefix when the testID designates WHERE
the element appears. For instance, `SpendingMonthCard` (`SMCD`) exposes `:home-spending-card` because it
lives inside `HomeScreen` (`HOME`), which is how end-to-end flows locate the element starting from the screen.
Demanding strict local identity unification there would break end-to-end navigation flows.

The rule distinguishing the two cases is precise:
1. A prefix that matches the code of ANY unit declared in the map represents a KNOWN IDENTITY (deliberate reuse).
2. A prefix shaped like an identity code (4-5 letters) that belongs to NO unit in the map is an ORPHAN IDENTITY.
   Because it does not exist in the map, it cannot exist in the generated dictionary, making it the exact case
   that forces rogue acronyms into manual spelling dictionaries. Only orphan identities are rejected.

For visual regression baselines (`<Unit>.<CODE>-VR-<variant>.png`), cross-unit delegation is not permitted:
a baseline is the physical proof of THIS specific unit, not a pointer to where it appears.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a node of kind `spec` | nodes of any kind other than spec | this unit: skips non-spec nodes |
| the graph | a built map graph, or nil | — | this unit: without a graph returns Pending, never approves in the dark |
| the spec code | a non-empty identity code | a spec with empty or whitespace code | this unit: skips empty code so absence is caught by the dedicated code presence gate |
| the specified code units | files reached via `specifies` edges in the map graph | units lacking a graph edge or unreadable files | the map graph and file system |
| the testID handles | testID attributes in code formatted with 4-5 letter code prefixes | short prefixes of 1-3 letters or non-code handles | this unit, filtering prefixes by code length pattern |
| the visual regression baselines | baseline image files matching the naming pattern `<Unit>.*-VR*.png` | non-image files or baselines not following the visual regression naming convention | this unit, matching filenames via glob |

## Effects

| Effect | Description |
| --- | --- |
| `IDCND-B01` | When the confronted node is not a spec, the gate skips confrontation. |
| `IDCND-B02` | When the graph is nil, the gate returns Pending without approving. |
| `IDCND-B03` | When the spec has no declared code, the gate skips without duplicating findings from the code presence gate. |
| `IDCND-B04` | A testID prefix matching the spec's own identity code passes. |
| `IDCND-B05` | An orphan testID prefix with code shape that belongs to no unit in the map fails. |
| `IDCND-B06` | A testID prefix matching the code of another declared map unit is accepted as legitimate container reuse. |
| `IDCND-B07` | A visual regression baseline whose code differs from the spec code fails. |
| `IDCND-B08` | Short testID prefixes of three letters or fewer do not have code shape and pass without accusation. |
| `IDCND-B09` | The failure verdict cites the conflicting acronyms and their origin files. |
| `IDCND-B10` | When no orphan testID prefixes or baseline discrepancies exist, the gate passes. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `IDCND-I01` | Without a map graph the gate never approves. It returns Pending because concordance cannot be evaluated without the global code inventory. | confronts a spec with nil graph and verifies it returns Pending |
| `IDCND-I02` | Cross-unit reuse is allowed only for testIDs, never for visual regression baselines. The baseline must prove this specific unit. | confronts a baseline carrying another unit's code and verifies it returns Fail |
| `IDCND-I03` | Absence of a spec code is never double-charged. It skips here so that the dedicated code presence gate reports the single defect. | confronts a spec with empty code and verifies it returns Skip |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `IDCND-X01` | Does not charge testID prefixes of three letters or fewer. | Common shorthand prefixes like `tab-` or `btn-` are not code-shaped acronyms; scrutinizing them would create noise across every UI component. |
| `IDCND-X02` | Does not charge specs for code absence. | Code presence is governed by `spec-has-code`; charging it here would create two findings for a single defect. |
| `IDCND-X03` | Does not forbid components from referencing parent screen codes in testIDs. | Child components regularly identify their screen context for end-to-end flow navigation; forbidding this would break flow selectors. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `IDCND-E01` | A file the spec governs (`specifies`) is gone from disk when its testIDs are read. | That file is left out; the other governed files are still confronted, and an orphan acronym in any of them still fails (`IDCND-B05`). | A file that no longer exists carries no testID, so there is nothing in it to be inconsistent; and one stale edge must not hide the orphan identity in the files that are there. <!-- @resilient: a stale map edge is expected between builds, and the next map build removes it --> |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized messages for skips and failure verdicts |
| DEP3 | `internal/mapx/model.go` | `EdgeSpecifies`, `Graph`, `KindSpec`, `Node` | core — graph model, node representations, and specifies edge type |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
