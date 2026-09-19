<!-- @anchors
  code: RTEXR
  updated_at: 2026-09-19
  layer: gate
-->
# RouteExists — declared route in specification must exist in application route registry

> **Code**: `RTEXR`

## Overview

Confronts the route declared by a specification against actual application code: **the route that the
specification declares must exist where the application registers its routes.**

The neighbouring `route-declared` gate confronts a specification against itself — does the specification
declare a route? This gate confronts the specification against the CODE: does that route actually exist
where the application registers its routes?

The critical distinction emerged in a real end-to-end delivery defect. A specification for a new screen
declared `> **Rota**: MetadataEdit` (or `route: MetadataEdit`), and another specification promised navigation
to it. The blocking `route-declared` gate gave a green checkmark to both specifications because both declared
routes. Yet the route existed nowhere in the application code. Both specifications described a path leading to
an unreachable screen, and the entire verification pipeline remained green.

This represents an anchor that lies in the most elusive way possible: nothing is missing and every document
references each other, but nobody asked whether the destination actually exists in the code.

Measured across 96 screen specifications in a real project prior to enabling this gate: 2 findings, both
true defects (the aforementioned unreachable end-to-end route, and another specification whose screen the
application registered under a different name). Zero false positives.

When a specification declares no route, this gate skips confrontation, as enforcing route declaration is the
exclusive responsibility of `route-declared`. Furthermore, when project configuration does not specify where
routes are registered (`route_registry`), or when zero routes can be extracted, the gate returns **Pending**
rather than falsely approving uninspected routes.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `spec` | nodes of kind `code`, `feature`, or `test` | this unit: skips non-spec nodes |
| the declared route | specification content declaring a route via `route:` or `rota:` | specifications containing no route declaration | this unit: skips specifications with no declared route |
| the route registry configuration | project configuration defining one or more `route_registry` glob patterns | unconfigured or empty `route_registry` globs | this unit: returns Pending when registry globs are missing |
| the route registry files | existing source files matching configured globs | missing or unreadable route files, or invalid glob patterns | `doublestar.Glob` and filesystem: returns read error causing Pending |
| the registered route definitions | routes matching standard navigation patterns or custom `derived.route_pattern` regex | files containing zero identifiable route definitions | this unit: returns Pending when zero routes are found |

## Effects

| Effect | Description |
| --- | --- |
| `RTEXR-B01` | When the confronted node is not of kind spec, the gate skips confrontation. |
| `RTEXR-B02` | When the specification does not declare a route, the gate skips confrontation. |
| `RTEXR-B03` | When project configuration defines no route registry glob patterns, confrontation returns Pending. |
| `RTEXR-B04` | When route registry glob pattern is invalid, the gate returns Pending reporting the glob error. |
| `RTEXR-B05` | When route registry files are inspected but zero registered routes are discovered, the gate returns Pending. |
| `RTEXR-B06` | When a declared screen name matches a component navigation prop in registered route files, the gate passes. |
| `RTEXR-B07` | When a declared screen name matches a navigation stack parameter type entry in registered route files, the gate passes. |
| `RTEXR-B08` | When a declared backend route matches an added HTTP resource registration, the gate passes. |
| `RTEXR-B09` | When a declared route includes an HTTP method verb prefix, the gate strips the verb and passes if the route path exists. |
| `RTEXR-B10` | When a declared route matches a registered route path regardless of leading slash differences, the gate passes. |
| `RTEXR-B11` | When project configuration specifies a custom route pattern regex, routes matching that pattern pass. |
| `RTEXR-B12` | When a custom route pattern regex is invalid or contains no capture group, the gate returns Pending. |
| `RTEXR-B13` | When a declared route is absent from all registered routes, the gate fails citing the missing route, total registered count, and searched globs. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RTEXR-I01` | Route presence is validated only for specifications; non-spec nodes skip confrontation. | confronts non-spec nodes and verifies the verdict is Skip |
| `RTEXR-I02` | The gate never approves route existence without inspecting registry files, returning Pending when registry configuration or routes are missing. | confronts specifications when route registry is empty or missing and verifies verdict is Pending |
| `RTEXR-I03` | Route matching is slash-normalized, treating route paths with and without a leading slash as equivalent representations. | confronts matching routes with alternating leading slashes and verifies they pass |
| `RTEXR-I04` | Missing routes always result in a blocking Fail verdict to prevent shipping unreachable destinations. | confronts an unregistered route and verifies the verdict is Fail |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RTEXR-X01` | Does not mandate that every specification declare a route. | Demanding route declaration is the exclusive responsibility of `route-declared`; duplicating it here would produce redundant gate failures. |
| `RTEXR-X02` | Does not validate route parameter schemas, HTTP payload structures, or response codes. | This gate validates whether the route exists in the application registry; parameter and payload validation belong to interface and contract gates. |
| `RTEXR-X03` | Does not evaluate authentication or access permissions attached to the route. | Verifying user roles and permissions belongs to authorization gates rather than registry existence checks. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — route registry globs and custom route pattern configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindSpec`, `Node` | core — graph model and node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
