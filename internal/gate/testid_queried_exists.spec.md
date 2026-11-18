<!-- @anchors
  code: TQETS
  updated_at: 2026-10-15
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
| the project configuration | a `Config` declaring `dialect.test_handles` (`exposed`, `queried`, optionally `mentioned`) and an E2E surface | a configuration missing any of them | this unit: without the patterns the gate is pending, without the surface it skips; it never assumes a syntax |
| the E2E surface flows | YAML or YML flow definition files located in the declared directory | non-YAML files (such as `.js` runner helpers) or unconfigured paths | this unit, filtering files during surface walk |
| the queried handles | static literal IDs or pattern prefixes extracted from flow definitions | dynamic runtime-interpolated IDs (`${...}`) | this unit, discarding expressions requiring runtime evaluation |
| the exposed handles | what the `exposed` and `mentioned` patterns capture in the map's code files: handles, heads (`prefix`), composed segments (`suffix`) | the files the map holds as tests, and every file outside the map | this unit, reading the map's code nodes |

## Effects

| Effect | Description |
| --- | --- |
| `TQETS-B01` | Without the patterns of `dialect.test_handles.exposed` and `dialect.test_handles.queried` declared the gate is Pending, naming the field missing: the core reads no syntax of a handle, in the code or in a flow. |
| `TQETS-B02` | When the project does not declare an E2E surface, the gate skips honestly rather than reporting compliance. |
| `TQETS-B03` | When no exposed handles are found in project code, the gate skips confrontation. |
| `TQETS-B04` | A handle queried by a flow that matches an exposed handle in code passes. |
| `TQETS-B05` | A handle queried by a flow that does not exist in any code file fails. |
| `TQETS-B06` | The failing verdict names the missing handle and lists all flow files that query it. |
| `TQETS-B07` | A non-existent handle inside a negative assertion (`assertNotVisible`) fails, preventing false green by vacuity. |
| `TQETS-B08` | A code template handle (`:item-*`) satisfies a concrete instance queried by a flow (`:item-3`). |
| `TQETS-B09` | A query the `queried` patterns capture in a `prefix` group — a flow querying by regular expression — is confronted as a head against the exposed handles. |
| `TQETS-B10` | A query the `queried` patterns do not capture — one the runner interpolates at run time (`${...}`) — is not confronted, and fails nothing. |
| `TQETS-B11` | A handle the `mentioned` patterns find — a marked literal in a lookup table or a constant object — counts as existing in the code. |
| `TQETS-B12` | A head the `mentioned` patterns find behind a fallback (`??`) counts as existing, and covers its instances. |
| `TQETS-B13` | A segment the `mentioned` patterns capture in a `suffix` group — composed in a child onto a prop from its parent (`*-row-*`) — covers the queries of composite ids holding it. |
| `TQETS-B14` | A `suffix` segment covers the composite ids that end with it (`-toggle`), not only the ones holding it in the middle. |
| `TQETS-B15` | Only the map's code files are read for handles: a file the map holds as a test exposes none, whatever its name — a test queries a handle, it does not create it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TQETS-I01` | Lack of configuration never produces a false approval. Without the patterns, or an E2E surface, the gate returns Pending or Skip, never Pass. | confronts unconfigured fixtures and verifies neither returns Pass |
| `TQETS-I02` | Findings are grouped by handle ID across flows so that a single missing handle queried multiple times does not multiply reported defects. | queries the same missing ID across flows and verifies the verdict groups them under that ID |
| `TQETS-I03` | Negative assertions are confronted with the same rigor as positive queries to eliminate green-by-vacuity passes. | executes an `assertNotVisible` with an unexposed handle and verifies it fails |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TQETS-X01` | Does not charge code for exposing handles that no flow queries. | Not every marked element needs an automated flow, and spec-level coverage is already measured by `testid-consistent`. |
| `TQETS-X02` | Does not execute flows or evaluate runtime JavaScript expressions. | Dynamic runtime string evaluation requires a running emulator or engine; static analysis inspects only the static prefix or skips interpolation. |
| `TQETS-X03` | Does not assume default patterns for a handle or for a flow's query (such as JSX's `testID=` or Maestro's `id:`). | Assuming defaults without project declaration would stamp approval over unverified codebases in unfamiliar stacks. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TQETS-E01` | The project declares an E2E surface, but the directory its pattern points to does not exist on disk. | `Skip` with the no-E2E-surface message; nothing is confronted. | A directory that is not there holds no flow, so there is nothing to confront: `Skip` says "not measured", where `Pass` would claim flows were checked and found clean. |
| `TQETS-E02` | A flow of the E2E surface, or a code file of the map, cannot be read. | `Pending` naming every unreadable file — never `Pass`, and no handle is accused while one is unread. | An unreadable flow may be the one that queries a missing handle, and an unreadable source the one that exposes it. Skipping them answered `Pass` over a flow the gate never read, or accused a handle that the unread source exposes. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
