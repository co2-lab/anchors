<!-- @anchors
  code: CNHNC
  updated_at: 2026-09-26
  layer: gate
-->
# CountHonored — a numerical assertion written in a spec must match reality in code

> **Code**: `CNHNC`

## Overview

Confronts numerical claims made in specs against the codebase: **does a number that a spec asserts about
the code match the actual count in code?**

The "lying anchor" defect has a numerical variant that ages completely on its own: a spec asserts "the 50
models of the product", an engineer introduces the 51st model, and the sentence instantly becomes a lie
without anyone ever touching or modifying the spec. In real project measurements, a single spec file held 7
numeric assertions, and adding ONE model rendered 10 sentences obsolete at once.

Unlike other anchor defects caused by developer carelessness or skipped steps, numerical drift depends
strictly on the passage of time. Every count written in prose represents debt accumulating compound interest.

The gate does NOT guess what to count. In the measured project, "51 models", "51 authorization clauses", and
"46 indexed models" represented three completely different questions evaluated over the exact same set of
files; no static heuristic can guess all three correctly. Instead, the spec explicitly DECLARES how to count
using the declaration contract:
- A glob pattern alone counts matching FILES (for example, counting model files).
- A glob pattern followed by a regex pattern counts OCCURRENCES matching the regex across those files (for
  example, counting authorization clauses).

Any undeclared number in prose (such as a data retention policy stating "90 days") is deliberately ignored:
the gate verifies what was explicitly contracted, not every arbitrary number appearing in text.

Crucially, the gate also confronts the PROSE situated alongside the count declaration marker, not merely the
marker itself. In real inspection, a marker declared 51 (matching code), yet a prose sentence three lines below
asserted "50 authorization clauses". Passing the marker while ignoring the prose would certify a spec that lies
to any human reader opening the document.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a node of kind `spec` | nodes of any kind other than spec | this unit: skips non-spec nodes |
| the count declarations | valid count markers declaring expected count and glob | specs without count declarations | this unit: skips specs with no count declarations |
| the glob expressions | valid doublestar glob expressions matching repository files | invalid globs or expressions matching zero files | this unit: returns errors for invalid syntax or empty matches |
| the regex patterns | valid regular expressions or empty string | invalid regex syntax | this unit: returns errors for invalid regex expressions |
| the prose statements | numbers preceding declared labels with non-restrictive complements | numbers followed by qualifying phrases describing subsets | this unit: excludes qualified subsets from confrontation |

## Effects

| Effect | Description |
| --- | --- |
| `CNHNC-B01` | When the confronted node is not a spec, the gate skips confrontation. |
| `CNHNC-B02` | When the spec contains no count declarations, the gate skips without asserting green. |
| `CNHNC-B03` | When a declared glob expression matches matching files and count equals expected, the gate passes. |
| `CNHNC-B04` | When the file count on disk differs from the expected count declared in the marker, the gate fails. |
| `CNHNC-B05` | When a regex pattern is provided, the gate counts regex occurrences across files rather than file counts. |
| `CNHNC-B06` | When regex occurrence count differs from the expected count declared in the marker, the gate fails. |
| `CNHNC-B07` | When a declared glob matches zero files on disk, the gate fails and warns to check the path. |
| `CNHNC-B08` | When a glob pattern contains invalid glob syntax, the gate fails reporting the glob error. |
| `CNHNC-B09` | When a count pattern contains invalid regex syntax, the gate fails reporting the regex error. |
| `CNHNC-B10` | When the marker count matches disk but adjacent prose asserts a conflicting count, the gate fails. |
| `CNHNC-B11` | Non-restrictive complements attached to prose labels are confronted as claims about the total count. |
| `CNHNC-B12` | Qualifying words following the label indicate subset claims and are excluded from total count confrontation. |
| `CNHNC-B13` | Arbitrary numbers in prose without an explicit count declaration are ignored and skip confrontation. |
| `CNHNC-B14` | Only FILES are counted: a directory the glob also matches is not a file. Counting it made `models/*` say 2 for one model and one subfolder, while the pattern mode skipped the same directory — the two modes disagreed on the same glob. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CNHNC-I01` | The gate only verifies explicit contracts. Numerical assertions without declaration markers never trigger a check. | confronts a spec with undeclared numbers in prose and verifies it returns Skip |
| `CNHNC-I02` | Both the declaration marker and the adjacent prose must agree with code reality; a correct marker cannot mask lying prose. | confronts a matching marker accompanied by divergent prose and verifies it returns Fail |
| `CNHNC-I03` | Zero file matches are treated as path or glob errors rather than a genuine count of zero files. | confronts a glob that matches zero files and verifies it returns Fail with a path check message |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CNHNC-X01` | Does not guess or infer what to count from arbitrary text in prose. | Multiple domain counts evaluate over the same files; guessing heuristics would create rampant false positives. |
| `CNHNC-X02` | Does not charge specs that contain no count declaration markers. | Not every spec makes numerical claims, and forcing every spec to declare counts would impose empty bureaucracy. |
| `CNHNC-X03` | Does not confront prose statements that qualify a subset rather than the total. | Differentiating subsets from totals requires deep semantic understanding; conservatively silencing alerts on qualified phrases prevents noisy false positives. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CNHNC-E01` | REF[CNHNC-B08]: a declared glob that does not parse is the configuration failure B08 answers: the gate fails carrying the glob error | — | — |
| `CNHNC-E02` | REF[CNHNC-B09]: a count pattern that does not compile is the configuration failure B09 answers: the gate fails carrying the regex error | — | — |
| `CNHNC-E03` | A file the glob matches cannot be read while occurrences of a pattern are counted. | `Fail` naming the file. | An unread file would undercount, and the gate would accuse the declared number of a divergence it never measured. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — project configuration |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — localized error and skip messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `KindSpec`, `Node` | core — graph model and spec node representation |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
