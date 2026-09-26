<!-- @anchors
  code: CRVCD
  updated_at: 2026-09-26
  layer: gate
-->
# CodeReferenceValid — cross-referenced requirement codes must resolve to existing units

> **Code**: `CRVCD`

## Overview

Confronts specification content against the project universe of identities: **every requirement code cited by a specification must resolve to an existing unit in the map.**

A lying anchor is the exact failure this framework exists to prevent. It takes a form that other gates miss: a specification cites an external requirement identifier — in a note, a narrative explanation, or a cross-reference — and the referenced unit does not exist anywhere in the repository. The citation simulates traceability while pointing to empty space.

A measured incident demonstrated this defect: a schema specification asserted creation of indices on 2026-08-11 and referenced four requirement codes belonging to specifications that were never created. The file described by that schema did not contain a single line implementing the entity. Yet the specification passed all gates: it possessed a code, a header, and valid section headings, and dependency checks only evaluate method symbols rather than requirement prose. Any future reader — human or autonomous agent — reads the text as an authentic historical record of completed work.

This gate enforces referential truth by extracting every requirement token shaped like an identity code followed by a requirement suffix and verifying that its owning unit exists in the project graph.

This gate operates in distinct territory from neighbouring gates:
- Unlike `dependency-honored`, which inspects dependency tables to confirm that promised Go symbols and methods appear in code, this gate validates requirement codes against the project identity universe. A unit can legitimately consume existing Go packages while citing non-existent requirement codes.
- Unlike `code-cataloged`, which ensures that exported code symbols are catalogued within their own unit specification, this gate governs external references pointing to other specifications.
- Unlike `triad-complete`, which enforces local triad structure within a single unit, this gate maintains referential integrity across the collective graph of specifications.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind spec | nodes of kind feature, code, test, plan, or config | this unit: skips non-spec nodes |
| the map graph | a built graph containing project identity nodes | nil or empty graph without identities | this unit: returns pending when graph or identities are absent |
| requirement citations | tokens matching requirement code patterns | arbitrary text, symbol names, or file paths | this unit: parses tokens matching configured code length patterns |
| identity universe | identity codes declared in headers or graph nodes | uncatalogued or orphaned codes | `mapx.Graph` and header metadata extraction |

## Effects

| Effect | Description |
| --- | --- |
| `CRVCD-B01` | When the confronted node kind is not a specification, the gate skips confrontation. |
| `CRVCD-B02` | When the map graph is nil, confrontation returns a pending verdict because external identities cannot be resolved. |
| `CRVCD-B03` | When the map graph contains no declared identity codes, confrontation returns pending rather than approving blindly. |
| `CRVCD-B04` | When all cited external requirement codes resolve to known units in the map graph, the gate passes. |
| `CRVCD-B05` | When a specification cites an external requirement code whose unit is not in the map graph, the gate fails. |
| `CRVCD-B06` | Citations matching the specification's own identity code are recognized as self-references and are not charged as orphans. |
| `CRVCD-B07` | Multiple orphaned requirement codes are reported in alphabetical order formatted as identifiers. |
| `CRVCD-B08` | Identity ownership is determined from graph node metadata and specification header comments on disk. |
| `CRVCD-B09` | Text tokens that do not match requirement code structure are ignored and never charged as external references. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CRVCD-I01` | A specification is sovereign over its own identity code, so internal cross-references to its own requirements never fail. | confronts a specification citing its own requirement codes and verifies pass |
| `CRVCD-I02` | An absent map graph or an empty identity universe always yields Pending, never Pass, refusing to approve what cannot be measured. | confronts citations with nil graph and empty graph and verifies pending |
| `CRVCD-I03` | Any unresolvable external requirement citation produces a blocking Fail verdict. | confronts unresolvable requirement codes and verifies fail |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CRVCD-X01` | Does not evaluate whether the referenced requirement behavior is implemented correctly. | Semantic correctness belongs to unit tests and linters; this gate enforces referential existence. |
| `CRVCD-X02` | Does not inspect non-specification artifacts like code or tests for dangling requirement citations. | Code and tests are verified against their owning specifications through feature-test matching and marker parity. |
| `CRVCD-X03` | Does not mandate that a specification must cite external requirements. | Isolated self-contained units legitimately reference no other specifications. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern`, `Config` | core — identity code length regex pattern and project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized gate verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindSpec`, `Node` | core — map graph structure, spec kind definition, and node representations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
