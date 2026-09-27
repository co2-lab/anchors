<!-- @anchors
  code: DCCVD
  updated_at: 2026-09-26
  layer: gate
-->
# DocsCovered — every spec must reach some page of the compiled documentation

> **Code**: `DCCVD`

## Overview

The `docs-covered` gate asks whether a spec reaches ANY page of the documentation the project
compiles from its templates. Its sibling, `docs-fresh`, asks whether each page that exists still
reflects the spec it came from; it cannot see a spec that no page asks for.

That is the silent defect this gate exists for. The templates select specs by filter (by layer,
for instance), so a spec that no filter selects compiles to nowhere, and nothing complains: every
page that exists is correct, and the unit has a complete triad and passes every relational gate.
The first spec of a new layer drops out of the documentation without a word, exactly when the
project grows and nobody is watching.

The gate is informative: the fix is a person's decision, to widen a template's filter or to add
the layer's page. Rebuilding the documentation does not fix it, because the build only produces
the pages the templates ask for.

The answer is a property of the whole set of specs, so it is computed once per project root and
map and reused for every spec of the scan. The verdict of each spec, however, is about that spec
only.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a node of kind `spec` | any other kind | this unit: skips anything that is not a spec |
| the templates | a project with a documentation templates directory | a project that declares no templates | this unit: skips when the directory is absent |
| the map | the graph of the current scan, whose spec nodes the templates select from | a missing or stale map | the caller: the scan hands over the map it built |

## Effects

| Effect | Description |
| --- | --- |
| `DCCVD-B01` | An artifact that is not a spec is skipped: only a spec is compiled into documentation. |
| `DCCVD-B02` | A project with no templates directory is skipped: demanding documentation from a project that declared none would invent a duty. |
| `DCCVD-B03` | A spec that no template reaches fails, and the verdict names the spec. |
| `DCCVD-B04` | A spec that some template reaches passes. |
| `DCCVD-B05` | When the templates do not compile, the gate skips with no message: the compile failure is the sibling gate's finding, reported there with the compiler's own words, and repeating it here would look like a second defect. |
| `DCCVD-B06` | The set of unreached specs is computed once per project root and map, and reused for every spec of the same scan; a different map recomputes it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCCVD-I01` | The verdict of a spec depends only on whether THAT spec is reached; another spec being an orphan never fails it. | confronts a reached spec in a project that also holds an orphan, and verifies it passes |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCCVD-X01` | Does not judge whether a page is built or up to date; it only asks whether some template selects the spec. | Freshness is the sibling gate's question, and a spec that is selected is covered even before the first build. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCCVD-E01` | REF[DCCVD-B05]: the templates failing to compile is the one failure the gate handles, and B05 states its answer: skip, and leave the report to the sibling gate | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/doct/doct.go` | `New`, `Uncovered`, `Dir` | the documentation compiler that knows what each template selects |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node`, `KindSpec` | core — the map the templates select from |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
