<!-- @anchors
  code: TRDCT
  updated_at: 2026-09-26
  layer: gate
-->
# TriggerDeclared — cited compliance triggers and obligations must exist in the declared vocabulary

> **Code**: `TRDCT`

## Overview

Confronts compliance triggers and obligations cited within specification text against the active
vocabulary declared in project configuration and adopted compliance packs: **every cited obligation
trigger or obligation name must exist in the declared vocabulary.**

Compliance obligations are triggered by declarations in artifact headers (such as `carries: personal-data`).
When a specification instructs an author on what trigger to declare, it teaches the vocabulary. If it
teaches an identifier that no pack declares, the author follows the instruction, writes the header, and
**no obligation ever triggers**.

This produces no structural failure in ordinary gates. The header is syntactically well-formed, header
checks pass, the unit appears covered by regulatory frameworks (such as LGPD or GDPR), yet actual
enforcement is zero. It is the most insidious form of a lying anchor: the instructional text is wrong,
causing every developer and automated agent relying on it to reproduce the defect.

Measured in doctrine: in a real project, 47 model specifications instructed authors to declare `carries: pii`
and cited an obligation named `pii-purgavel`. Neither identifier existed in any pack — the canonical
vocabulary declared across packs was `carries: personal-data`, with the obligations `lgpd-eliminacao` and
`lgpd-portabilidade`. All 47 instances originated from the same boilerplate snippet, copied from spec to spec,
without any check raising an alarm.

What separates this gate from neighbouring gates:
- Header gates (`header-conforms`, `header-valid`): only verify that header syntax is well-formed and fields
  are structured; they do not validate whether the trigger values activate real obligations.
- Obligation gates (`obligation-honored`): verify that active obligations attached to a unit are fulfilled;
  but if a misnamed trigger never fired the obligation in the first place, `obligation-honored` has nothing
  to enforce.
- `trigger-declared`: confronts cited trigger symbols and obligation names mentioned in spec prose against
  the declared universe of packs and configuration.

Finally, citing triggers is entirely optional. Specifications that cite no compliance triggers skip cleanly.
When a project has declared no compliance obligations and loaded no packs, the gate returns **Pending**
rather than approving unverified vocabulary or penalizing projects without compliance.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `spec` | nodes of kind `code`, `feature`, or `test` | this unit: skips confrontation on non-spec nodes |
| the specification content | specification text containing backticked trigger or obligation citations | specifications containing no trigger or obligation citations | this unit: skips confrontation when citations are absent |
| the trigger predicates | recognized trigger keys (`carries`, `processing`, `renders`, `shared-with`, `retains`, `transfers`) | arbitrary metadata keys such as `layer` or `code` | this unit: filters citations against closed key set |
| the compliance vocabulary | declared triggers and obligations from configuration and loaded packs | empty vocabulary when no obligations or packs are configured | this unit: returns Pending when vocabulary is empty |
| the pack loader | valid packs loaded via `pack.LoadAll` | unreadable or missing pack files | `pack.LoadAll`: returns loaded packs and error |

## Effects

| Effect | Description |
| --- | --- |
| `TRDCT-B01` | When the confronted node is not of kind spec, the gate skips confrontation. |
| `TRDCT-B02` | When the specification text contains no cited obligation triggers, the gate skips confrontation. |
| `TRDCT-B03` | When no compliance vocabulary is declared in configuration or packs, confronting cited triggers returns Pending. |
| `TRDCT-B04` | Non-trigger key-value citations (such as layer or code markers) are ignored. |
| `TRDCT-B05` | Natural language prose mentioning compliance terms without backtick quotes is ignored. |
| `TRDCT-B06` | When every cited trigger value exists in the declared vocabulary, the gate passes. |
| `TRDCT-B07` | When cited obligation names exist in the declared vocabulary, the gate passes. |
| `TRDCT-B08` | A cited trigger value absent from the declared vocabulary fails, providing a nearest-match suggestion. |
| `TRDCT-B09` | When no close match exists for an undeclared trigger, the failure suggests available declared triggers. |
| `TRDCT-B10` | A cited obligation name absent from the declared vocabulary fails. |
| `TRDCT-B11` | Duplicate citations of the same trigger or obligation within a file are deduplicated in defect reporting. |
| `TRDCT-B12` | Multiple vocabulary errors are sorted deterministically and aggregated in the failure verdict. |
| `TRDCT-B13` | When no declared trigger is close enough to suggest, the verdict lists the declared ones — capped at four. A list of every trigger a large project declares would bury the advice it exists to give. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TRDCT-I01` | Trigger and obligation citations are validated only in specifications; other artifacts skip confrontation. | confronts non-spec nodes and verifies the verdict is Skip |
| `TRDCT-I02` | Missing compliance packs and obligations produce Pending rather than Pass, preventing silent approval of unverified claims. | confronts cited triggers when configuration has no vocabulary and verifies verdict is Pending |
| `TRDCT-I03` | Trigger keys are restricted to a closed set of recognized compliance predicates to prevent false positives on general key-value metadata. | confronts specifications citing non-trigger metadata keys and verifies they are skipped |
| `TRDCT-I04` | Undeclared triggers produce a definitive Fail verdict because teaching invalid triggers breaks downstream compliance enforcement. | confronts an unknown trigger and verifies the verdict is Fail |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TRDCT-X01` | Does not mandate that specifications cite compliance triggers or obligations. | Compliance triggers are domain-specific; non-regulated units have no obligation to cite them. |
| `TRDCT-X02` | Does not validate compliance triggers in source code, features, or test files. | Instructional trigger definitions belong to specifications; code and tests realize requirements rather than teaching vocabulary. |
| `TRDCT-X03` | Does not enforce implementation of the cited obligations within code. | Verifying that declared obligations are satisfied is the dedicated responsibility of `obligation-honored`. |
| `TRDCT-X04` | Does not inspect unquoted natural language mentions of compliance concepts. | Only explicit backticked symbols represent actionable syntax instructions. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `TRDCT-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration and obligation definitions |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindSpec`, `Node` | core — graph model and artifact representations |
| DEP4 | `internal/pack/pack.go` | `LoadAll` | core — loads adopted compliance packs and their obligation triggers |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
