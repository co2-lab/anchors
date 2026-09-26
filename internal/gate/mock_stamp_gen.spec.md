<!-- @anchors
  code: MKSTP
  updated_at: 2026-09-26
  layer: gate
-->
# MockStampGenerator — writes the missing `@contract` stamps, and never rewrites one

> **Code**: `MKSTP`

## Overview

A project adopting `mock-stamped` has every double of a governed module unstamped at once.
Measured in the project that asked for this: 1278 `jest.mock`/`vi.mock` calls, none stamped.
Stamping by hand is not viable, and a hand-written stamp is also the one most likely to point at
the wrong snippet.

The generator writes, above each double that has none, the stamp the gate recomputes. It is built
from the gate's own functions, so a stamp it writes is exactly the stamp the gate checks. Measured
on that project, on a copy: 1214 stamps in 278 files, and afterwards `mock-stamped` 278 ✓, 0 ✗.

**It only ADDS.** An existing stamp is never rewritten, even when it diverges. A divergence is the
gate saying the double may have drifted; a tool that refreshed it would let whoever edits the test
regenerate the stamp to match their own mock, and the stamp would certify itself.

Exposed as `anchors stamp [tests...]`, with `--dry-run`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the test content | any text | — | this unit |
| the double pattern | `derived.mock_detect`, with one capture group | an undeclared pattern | this unit: refuses to run, since no double can be found |
| the export pattern | `derived.export_detect` or the dialect's | none declared | this unit: every stamp then covers the whole module |
| the map | a built graph | — | the command: it loads the map before calling |

## Effects

| Effect | Description |
| --- | --- |
| `MKSTP-B01` | `GenerateStamps` writes a stamp that passes the gate, and a later change to the stamped snippet fails it. |
| `MKSTP-B02` | A specifier matching several files resolves to the one sharing the longest directory prefix with the test; a tie is skipped and reported, never guessed. |
| `MKSTP-B03` | Each factory key that names an export of the module gets its own stamp, anchored on that export and covering its block — a change to another member does not make it diverge. |
| `MKSTP-B04` | With no factory key naming an export (an automock, for instance), one stamp covers the whole module — which is what an automock replaces. |
| `MKSTP-B05` | A double of a module outside the map (a third-party library) is not stamped. |
| `MKSTP-B06` | A whole-module stamp starts after the module's `@anchors` header, so an edit that only rewrites the header (`updated_at:`, which `check --fix` bumps on every edit) does not make it diverge. |
| `MKSTP-B07` | A key added later to an already-stamped double gets its own stamp when no existing stamp of that module covers its export (same anchor, or a range containing it); the existing stamps are left untouched, and a second run adds nothing. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MKSTP-I01` | An existing stamp is never rewritten, even when it diverges. | runs the generator over a test with a divergent stamp and verifies the content is untouched |
| `MKSTP-I02` | A line that occurs more than once in the module never anchors a stamp — the gate refuses an ambiguous anchor. | stamps a module whose first line is repeated and verifies the gate passes |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MKSTP-X01` | Does not refresh a divergent stamp. | Refreshing would let the stamp certify itself; a divergence calls for a person looking at the double. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MKSTP-E01` | The project's `derived.mock_detect` does not compile as a regular expression, or has no capture group. | `GenerateStamps` returns the error and the content unchanged; no stamp is written. | The pattern is the only way to find a double: a broken one would find none, and an empty result would read as "nothing to stamp" in a test full of doubles. The error sends the author to fix the configuration. |
| `MKSTP-E02` | The project does not declare `derived.mock_detect`. | `GenerateStamps` returns an error naming `derived.mock_detect`, and the content unchanged. | Without the pattern no double can be found; returning no stamps and no error would claim the test has no double to stamp. |
| `MKSTP-E03` | The file the map resolves a double to is no longer on disk. | That double is listed as skipped, with a reason naming the file that could not be read; the other doubles of the test are still stamped. | A stamp is a hash of the real module's snippet, and a file that is not there has no snippet to hash. Reporting it as skipped tells the author which double is still unstamped, and one stale map node must not block the stamps the rest of the test can get. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/mock_stamped.go` | `declaredStamps`, `snippetHash`, `moduleHasStamp` | gate — the stamp is written with the functions that check it |
| DEP2 | `internal/gate/mock_typed.go` | `isGovernedModule` | gate — the same reading of "a module of this project" |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
