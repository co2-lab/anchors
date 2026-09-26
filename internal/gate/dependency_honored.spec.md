<!-- @anchors
  code: DEPHN
  updated_at: 2026-09-26
  layer: gate
-->
# DependencyHonored — methods promised in the dependency table are consumed in code

> **Code**: `DEPHN`

## Overview

Confronts a spec's **Dependency Table** against actual use in the unit's code: **every method
symbol a spec promises to consume must genuinely appear in the non-comment content of the code
it governs.**

This is the **relational gate of the spec→code edge**. It catches the divergence that sibling gates
(such as `feature-test-match`) cannot see: a spec declares `DEP2 → metadataVersioning · resolveVersion, applyEdit`,
but the unit's code calls only `resolveVersion` — the promised `applyEdit` is never invoked.
That exact divergence was the bug that slipped past 65 green tests in the originating project's
end-to-end suite.

The ruler is **static, without execution**:
- **Only symbols are confronted**: identifiers enclosed in `backticks` in the Method field. A prose
  description (such as `"CRUD + queries"` or `"requests"`) is not confrontable and is deliberately ignored;
  prose explains the relationship rather than establishing a verifiable contractual promise.
- Each declared symbol must appear as a token within the non-comment content of the code the spec
  `specifies`. Line comments are stripped before checking so that comments cannot mask unused dependencies.
  A declared symbol absent from code fails.
- When a promised symbol is missing from code, the gate checks for a near rename (e.g. `ping` → `pingHeartbeat`
  or `computeSeatsAmount` → `computeSeatsAmountCents`). If an identifier in code extends the symbol as a prefix
  or suffix and the symbol has at least four characters, the verdict actively suggests the rename instead
  of merely reporting an error.
- Undetermined and unapplicable states are explicit:
  - An artifact that is not a spec skips with reason (`not a spec — only spec has Dependency Table`).
  - When no relational graph is loaded, the verdict is undetermined (`Pending`).
  - When the Dependency Table contains no confrontable symbols, the check skips (`Dependency Table promises no confrontable SYMBOL`).
  - When the spec governs no code (`specifies` edge absent), the verdict is undetermined (`Pending`).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node in the map graph | — (the gate does not select the target) | the gate engine, routing by declared `on: [spec]` |
| the relational graph | a populated map graph with spec edges | `nil` graph (verdict is undetermined as Pending) | the caller / map loader |
| the dependency promises | backticked identifier symbols in `depends-on` edges | prose descriptions without backticks (ignored without confrontation) | this unit, extracting symbols via `backtickedSymbols` |
| the governed code | non-comment content of code files linked via `specifies` edges | specs with no `specifies` edge (returns Pending) | this unit, stripping comments and checking word-boundary tokens |

## Effects

| Effect | Description |
| --- | --- |
| `DEPHN-B01` | An artifact that is not a spec leaves without a verdict: only specs have a Dependency Table. |
| `DEPHN-B02` | Without a relational map the verdict is UNDETERMINED (Pending), because dependency and specification edges cannot be traversed. |
| `DEPHN-B03` | A spec declaring no confrontable symbols in its Dependency Table leaves without a verdict (Skip): prose descriptions and empty tables promise no verifiable identifiers. |
| `DEPHN-B04` | A spec that specifies no code files leaves the verdict UNDETERMINED (Pending): there is no governed code to confront yet. |
| `DEPHN-B05` | When every promised symbol appears in the non-comment content of the governed code, the gate passes. |
| `DEPHN-B06` | When a promised symbol is absent from the governed code, the gate fails and names the unused symbol and dependency target file. |
| `DEPHN-B07` | When an absent symbol resembles an identifier in code (sharing a prefix or suffix extension), the failure verdict suggests the candidate rename. |
| `DEPHN-B08` | Line comments in governed code are stripped before confrontation, so symbols appearing exclusively within comments do not fulfill the promise. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DEPHN-I01` | Prose descriptions in dependency methods are never treated as contracts; only backticked identifiers constitute promises to verify. | confronts dependency methods with and without backticks and verifies prose is ignored |
| `DEPHN-I02` | Symbol presence in code is matched strictly on token word boundaries, never as a substring of a larger identifier name. | confronts a code file containing a longer identifier containing the symbol as a substring and verifies failure |
| `DEPHN-I03` | Near-symbol rename suggestions are strictly conservative, requiring prefix or suffix containment and a minimum symbol length of four characters. | evaluates candidate identifiers of various lengths and forms, verifying only valid extensions trigger suggestions |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DEPHN-X01` | Static textual confrontation without runtime execution. | The gate inspects non-comment source tokens rather than executing target code or inspecting call graphs; dynamic verification belongs to test suites. |
| `DEPHN-X02` | Does not interpret dependency semantics, parameter signatures, or method types. | The gate enforces relational honesty between declared symbols and source references; semantic and type checking belongs to language compilers. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/query.go` | `Neighbors` | core — querying outgoing edges from the spec node |
| DEP2 | `internal/mapx/model.go` | `KindSpec`, `EdgeDependsOn`, `EdgeSpecifies` | core — identifying spec nodes and dependency/governance edge types |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — localized messages for gate verdicts |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
