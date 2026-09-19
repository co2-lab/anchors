<!-- @anchors
  code: RPHRG
  updated_at: 2026-09-19
  layer: gate
-->
# RegionPairHonored — every opened source region must close with its own identity code

> **Code**: `RPHRG`

## Overview

Confronts source code regions against their pairing markers: **every opened code region must close,
and must close specifying its own exact identity code.**

A region marker pair gives identity a concrete line interval within the source file. Without region
boundaries, one only knows that a given file realizes a requirement, but not WHERE within the file that
implementation resides (TRACEABILITY §3). The region pair is the single fragile component of this mechanism,
and its fragility is particularly perilous because it raises no compiler or syntax errors: closing in the
wrong location produces an interval that is syntactically VALID and completely WRONG.

This gate detects the three critical structural pairing defects that cannot be reliably spotted by human
review when reading large source files top-to-bottom:
1. **Unclosed region (`sem-fecho`)**: A region opened and never closed, causing the requirement interval
   to bleed silently to the end of the file.
2. **Orphan close (`fecho-orfao`)**: An end region marker appearing without a preceding opening marker,
   typically left behind from careless cut-and-paste refactorings.
3. **Mismatched close (`fecho-trocado`)**: An end marker carrying a DIFFERENT identity code than the currently
   open region, indicating inverted or twisted nesting.

The third defect is the exact reason why region close markers must explicitly carry the identity code rather
than using an anonymous end marker. With anonymous end markers, the count of open and close markers would
balance perfectly, the gate would remain green, and the freshness stamp would measure the NEIGHBOR interval:
requirement A would appear to change when neighbor B was edited. A silent, crossed error is the worst kind
of defect: it vanishes during initial measurement and resurfaces weeks later as an erroneous conclusion.

Unlike gates that report non-blocking warnings for domain ambiguities (such as an uncatalogued requirement
letter), this gate issues a definitive blocking **Fail**. A malformed region pair has no two legitimate
interpretations: it is an objective syntax defect in code structure with a single unambiguous, local fix.

Finally, the ABSENCE of region markers is NOT a defect. Granular region delimitation is entirely optional;
when a file does not define regions, the whole file revision applies cleanly.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind `code` or `test` | nodes of kind `spec`, `feature`, or documentation | this unit: skips non-code and non-test nodes |
| the source content | source code containing zero or more region markers | non-text binary assets | the gate engine, routing text artifacts |
| the region markers | comments matching `#region [code]` and `#endregion [code]` syntax | arbitrary comments without region keywords | `scan.Regioes`: parses line comments and extracts markers |
| the region errors | pairing defects identified during stack-based region scan | valid balanced and properly nested regions | `scan.Regioes`: reports unclosed, orphan, and mismatched errors |

## Effects

| Effect | Description |
| --- | --- |
| `RPHRG-B01` | When the confronted node is neither code nor test, the gate skips confrontation. |
| `RPHRG-B02` | When a code or test file contains no region markers, the gate skips without asserting failure. |
| `RPHRG-B03` | When all opened regions close with their matching identity code, the gate passes. |
| `RPHRG-B04` | An opened region that is never closed fails, reporting the line and missing close marker. |
| `RPHRG-B05` | An end region marker without an opening marker fails as an orphan close. |
| `RPHRG-B06` | An end region marker closing with a different code than the open region fails, citing both codes and line number. |
| `RPHRG-B07` | Multiple pairing errors within a file are ordered sequentially by line number in the verdict. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RPHRG-I01` | Region absence is never charged as a failure. Delimitation is optional and files without regions evaluate at whole-file scope. | confronts a code file without region markers and verifies it returns Skip |
| `RPHRG-I02` | Pairing defects always result in a blocking Fail verdict, never Pending, because malformed nesting has only one correct repair. | confronts mismatched or unclosed regions and verifies the verdict is Fail |
| `RPHRG-I03` | End markers must explicitly match opening codes to prevent silent interval cross-over between neighboring requirements. | confronts inverted nesting tags and verifies it fails naming both codes |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RPHRG-X01` | Does not mandate region markers in code or test files. | Region markers provide fine-grained interval traceability, but requiring them everywhere would penalize small single-purpose files. |
| `RPHRG-X02` | Does not inspect region markers inside specs or markdown files. | Specifications and guides frequently contain documentation examples showing region syntax; parsing them would trigger false alarms. |
| `RPHRG-X03` | Does not enforce semantic validity of the code enclosed within a region. | This gate validates the integrity of the interval boundary markers; code logic and semantics belong to compiler, linters, and unit tests. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized verdict and defect messages |
| DEP3 | `internal/mapx/model.go` | `KindCode`, `KindTest`, `Node` | core — node representations and artifact kinds |
| DEP4 | `internal/scan/region.go` | `Regioes` | core — extracts region intervals and detects pairing errors |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
