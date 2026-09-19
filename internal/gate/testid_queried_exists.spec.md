<!-- @anchors
  code: TQETS
  updated_at: 2026-09-19
  layer: gate
-->
# TestidQueriedExists — every handle queried by an E2E flow must exist in code

> **Code**: `TQETS`

## Overview

Confronts the end-to-end execution surface against the application code: **does an automated flow
query a test handle that no code file exposes?**

This gate closes the fourth edge of the testID contract. The sibling gate `testid-consistent` starts
from a single spec and confronts inventory coherence; however, the E2E surface encompasses the entire
tree of flows across the whole project. Charging an individual screen's spec for an external flow handle
would accuse innocent specs of foreign handles (measured when attempted: 829 findings in a single spec,
all belonging to other screens). Therefore, the question belongs to the PROJECT scope, confronting both
sides collectively.

The measured defect in the reference application (2026-08-25) revealed **13 invented IDs in flows**,
frequently caused by incorrect screen prefixes (`review-*` where the component actually emits `revi-*`).
The assumption that "an invalid ID will simply fail at runtime in the test runner" is FALSE for two
independent reasons:
1. **Flow execution reachability**: The runner must reach the step. In a test suite where step 3 fails,
   phantom IDs in subsequent steps remain completely hidden — surfacing one per run across seven iterations.
2. **Vacuum passing in negative assertions**: In `assertNotVisible`, an invented or non-existent handle
   causes the test to PASS immediately. In `REVI-R03`, asserting `:review-edit-controls` (which never existed)
   proved the exact opposite of reality: green by VACUITY. This is the most dangerous defect, and the test
   runner can never catch it.

The gate enforces a STRICT SINGLE DIRECTION (flow → code). The reverse (code exposing a handle not queried
by any flow) is not an issue: not every marked UI element requires an automated scenario, and `testid-consistent`
governs spec-level declaration.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any project node | — (evaluated at project scope) | the gate engine, routing by the declared `on:` |
| the project configuration | a `Config` declaring `derived.test_handle` and an E2E surface | a configuration missing either property | this unit: without configuration the gate skips, never assumes defaults |
| the E2E surface flows | YAML or YML flow definition files located in the declared directory | non-YAML files (such as `.js` runner helpers) or unconfigured paths | this unit, filtering files during surface walk |
| the queried handles | static literal IDs or pattern prefixes extracted from flow definitions | dynamic runtime-interpolated IDs (`${...}`) | this unit, discarding expressions requiring runtime evaluation |
| the exposed handles | handles matching `test_handle`, marked literals (`:id`), or prop composition suffixes | test files (`.test.`, `.spec.`), dependencies, and build directories | this unit, walking source code and ignoring non-production files |

## Effects

| Effect | Description |
| --- | --- |
| `TQETS-B01` | When the project does not declare a test handle attribute, the gate skips without asserting green. |
| `TQETS-B02` | When the project does not declare an E2E surface, the gate skips honestly rather than reporting compliance. |
| `TQETS-B03` | When no exposed handles are found in project code, the gate skips confrontation. |
| `TQETS-B04` | A handle queried by a flow that matches an exposed handle in code passes. |
| `TQETS-B05` | A handle queried by a flow that does not exist in any code file fails. |
| `TQETS-B06` | The failing verdict names the missing handle and lists all flow files that query it. |
| `TQETS-B07` | A non-existent handle inside a negative assertion (`assertNotVisible`) fails, preventing false green by vacuity. |
| `TQETS-B08` | A code template handle (`:item-*`) satisfies a concrete instance queried by a flow (`:item-3`). |
| `TQETS-B09` | A regex pattern queried by a flow matches against the exposed prefix head in code. |
| `TQETS-B10` | A dynamic runtime expression interpolated by the runner (`${...}`) is skipped without failing. |
| `TQETS-B11` | A marked handle literal defined in a lookup table or constant object counts as exposed. |
| `TQETS-B12` | A template handle defined behind a nullish coalescing fallback (`??`) counts as exposed. |
| `TQETS-B13` | A suffix composed in a child component from a parent prop (`*-row-*`) satisfies queries against composite IDs. |
| `TQETS-B14` | A terminal suffix composed from a prop (`-toggle`) satisfies composite IDs ending with that segment. |
| `TQETS-B15` | Handles referenced only within test files (`.test.tsx`, `.spec.ts`) do not count as exposed in application code. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TQETS-I01` | Lack of configuration never produces a false approval. Without an explicit test handle or E2E surface, the gate returns Skip. | confronts unconfigured fixtures and verifies neither returns Pass |
| `TQETS-I02` | Findings are grouped by handle ID across flows so that a single missing handle queried multiple times does not multiply reported defects. | queries the same missing ID across flows and verifies the verdict groups them under that ID |
| `TQETS-I03` | Negative assertions are confronted with the same rigor as positive queries to eliminate green-by-vacuity passes. | executes an `assertNotVisible` with an unexposed handle and verifies it fails |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TQETS-X01` | Does not charge code for exposing handles that no flow queries. | Not every marked element needs an automated flow, and spec-level coverage is already measured by `testid-consistent`. |
| `TQETS-X02` | Does not execute flows or evaluate runtime JavaScript expressions. | Dynamic runtime string evaluation requires a running emulator or engine; static analysis inspects only the static prefix or skips interpolation. |
| `TQETS-X03` | Does not assume a default test handle attribute (such as `testID`). | Assuming defaults without project declaration would stamp approval over unverified codebases in unfamiliar stacks. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration containing derived surfaces and test handle settings |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized messages for skips and failure verdicts |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node` | core — graph model and node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
