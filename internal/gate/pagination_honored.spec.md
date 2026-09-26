<!-- @anchors
  code: PGNHN
  updated_at: 2026-09-26
  layer: gate
-->
# PaginationHonored — what promises a SET does not return the first page in silence

> **Code**: `PGNHN`

## Overview

Confronts a function against the promise its name makes: **whoever says "list all" cannot
deliver the first hundred without warning.**

It is the COST/SCALE class of defect — the one that shows up in no test, because the test
runs with three records and production runs with three thousand. Nothing in the code is wrong:
the query is valid, the type is the expected one, the suite passes. The defect is the difference between what
the name promises and what the function delivers when the data grows.

The distinction that gives the ruler: **a CALLER's limit** is not a defect — whoever passed the limit
knows there is more. **A HIDDEN limit** is: a default value the caller does not see makes the
function promise the set and return a slice. The hundred-and-first row is never
processed, and no one is notified.

When the module has sisters that paginate, the proof is by ASYMMETRY: the author knew the
pattern, and the one that does not paginate is forgetfulness, not decision. Without sisters paginating the verdict is
weaker, and it only accuses when the name promises a set unambiguously.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the stack's dialect | declared in the project's Structure | — (absent is a case, not an error) | this unit: without a dialect the verdict is indeterminate, never approval |
| the function's name | any exported identifier the dialect recognizes | a name that does not promise a set | this unit, by the vocabulary of promise |
| the waiver | the mark on the function, with a written reason | a bare mark, without a reason | this unit: a waiver without a why does not waive |

## Effects

| Effect | Description |
| --- | --- |
| `PGNHN-B01` | A function with a limit received from the CALLER passes: the page is deliberate, and whoever asked for it knows there is more. |
| `PGNHN-B02` | A function with a limit HIDDEN in a default value is accused: the name promises the set and the return is a slice. |
| `PGNHN-B03` | The NAME bounds the promise: whoever does not promise a set is not charged, because there is no promise to break. |
| `PGNHN-B04` | When there are sisters that paginate in the same module, the ASYMMETRY is the proof: the author knew the pattern. |
| `PGNHN-B05` | The waiver is DECLARED and with a written reason, and leaves the report. |
| `PGNHN-B06` | The verdict OFFERS the way out: rename to what the function does, expose the limit, return the cursor, or waive with a reason. |
| `PGNHN-B07` | Without a declared dialect the verdict is INDETERMINATE, and says so — approving without being able to read the code would be stamping what was not checked. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PGNHN-I01` | The ruler is AGNOSTIC: the confronted truth — the name promises a set, the return is partial — belongs to no language. What comes from the project is only how function, loop and cursor are recognized. | confronts the same defect written in different dialects and verifies the same verdict |
| `PGNHN-I02` | Where the unit does not recognize the construct, it stays silent. Silence is better than an invented accusation — a false positive here trains the team to ignore the gate. | confronts code whose pattern the dialect does not reach and verifies that nothing is accused |
| `PGNHN-I03` | A cursor without a loop does not count as pagination: returning the cursor and not walking it leaves the consumer with the same slice, only with the appearance of completeness. | declares the cursor without the loop and verifies that the accusation remains |
| `PGNHN-I04` | A provider prefix in the name does not hide the promise: what the name says counts, wherever it comes from. | confronts a function with a provider prefix and verifies that the promise is still charged |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PGNHN-X01` | Does not invent a cursor where the provider offers none. | Accusing the absence of a mechanism the dependency does not have would transfer to the author a defect that is not theirs — and the advice would be impossible to follow. |
| `PGNHN-X02` | Does not measure PERFORMANCE nor page size. | The ruler is the broken promise, not the cost of the query. Judging whether a hundred is a lot or a little depends on the domain, and that is the project's decision, not the gate's. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `PGNHN-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/dialect.go` | `DialectFor` | core — recognizing function, loop and cursor belongs to the project's dialect |
| DEP2 | `internal/mapx/model.go` | `KindCode` | core — the gate only has jurisdiction over code |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
