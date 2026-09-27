<!-- @anchors
  code: RCDLR
  updated_at: 2026-09-26
  layer: infra
-->
# RecodeDialect — the project's own surfaces of a code: testID prefixes and file names

> **Code**: `RCDLR`

## Overview

Beyond the header, the scenario codes and the bare mentions, a project may carry a unit's code in two
more places of its own making, and it declares them in its `recode:` block: a testID prefix derived from
the code (the code projected on the observable edge of a screen, such as `tcdt-amount`), and file names
built from the code (end-to-end flows, visual snapshots). This unit is the text and path engine for those
two surfaces, used by the recode planner (`RCPLR`).

The testID prefix is only rewritten where it is a testID: inside a string literal, right after the
opening quote, and followed by a hyphen. A word that merely starts with the same letters is left alone.
When the expected prefix is absent, the planner can still tell that a file holds testIDs, so it warns
about a divergent prefix instead of guessing it.

A file name is renamed only when the code sits on the name's boundaries (the start, a dot, a hyphen, the
end); the folder never changes, and the path comes back with forward slashes, the form the project's map
uses.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the testID convention | `lower` (the code in lower case) | any other value | this unit: an unknown convention derives no prefix, so nothing is rewritten |
| the file patterns | glob patterns holding `{{code}}` where the code goes | — | the project's `recode:` block |

## Effects

| Effect | Description |
| --- | --- |
| `RCDLR-B01` | The `lower` convention derives the testID prefix as the code in lower case; any other convention derives none. (`TestIDPrefix`) |
| `RCDLR-B02` | A testID prefix is rewritten only right after a quote or backtick and when followed by a hyphen and a letter or digit; the number of rewrites is counted. (`RewriteTestIDs`) |
| `RCDLR-B03` | With no old prefix, or the same prefix on both sides, the text is left unchanged with zero rewrites. |
| `RCDLR-B04` | The occurrences of a testID prefix are counted with the same rule, and an empty prefix counts zero. (`CountTestIDPrefix`) |
| `RCDLR-B05` | The `testID=` attributes of a text are counted whatever their prefix. (`CountAnyTestID`) |
| `RCDLR-B06` | A path matches the project's file patterns for a code when one pattern, with `{{code}}` replaced by the code, matches it; the code is matched with its exact case. (`FileMatchesCode`) |
| `RCDLR-B07` | A file whose name holds the old code on its boundaries (start, dot, hyphen, end) is renamed to the new code in the name only, keeping its folder; a name without it keeps its path. (`RenameFilePath`) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RCDLR-X01` | Does not rename a file whose name holds a longer code that contains the old one. | A mass rename that turned `TCDTXX` into `TCTXXX` would move another unit's files. |

## Errors

none — the unit transforms text and paths; a pattern that does not match or a prefix that is absent is a count of zero, not a failure.

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
