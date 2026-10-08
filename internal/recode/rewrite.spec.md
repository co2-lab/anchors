<!-- @anchors
  code: RCRWR
  updated_at: 2026-10-08
  layer: infra
-->
# RecodeRewrite — renaming an identity code inside a text, on every surface where it appears

> **Code**: `RCRWR`

## Overview

A unit's code is stable only if it can be renamed safely; otherwise "stable" just means "impossible to
change". This unit is the pure text engine of the rename: given a text, the old code and the new one, it
rewrites every place the old code appears, and it lists those places for a dry run.

A code appears on three surfaces: in the @anchors header (`code:` or `ref:`, where `ref:` may be a list
of several codes), in scenario codes derived from it (the code followed by a hyphen and a suffix, such as
`-B01`, `-DS-receipt-none` or `-VR`), and as a bare mention in the prose or comments of other units. The
three are rewritten in that order, scenario codes first, so that one step does not corrupt the next.
Every match is anchored on word boundaries: a longer code that contains the old one, or a word that ends
with it, is never touched.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the old and new codes | well-formed codes of the project's declared lengths | malformed codes | the planner (`RCPLR`), which refuses them before any rewrite |
| the text | any file content | — | the planner reads the files |

## Effects

| Effect | Description |
| --- | --- |
| `RCRWR-B01` | A code is well formed when it is made of upper-case letters and digits and its length is one the project declares. (`ValidCode`) |
| `RCRWR-B02` | The old code in a `code:` or `ref:` header line is replaced, in the markdown, line-comment and hash-comment header styles alike. |
| `RCRWR-B03` | In a `ref:` list only the old code changes; the other codes of the list stay. |
| `RCRWR-B04` | Every scenario code built on the old code keeps its suffix and takes the new code, and each one counts as one replacement. |
| `RCRWR-B05` | A bare mention of the old code is replaced, the character right after it is kept, and each mention counts as one replacement. |
| `RCRWR-B06` | The dry-run listing classifies each occurrence as a scenario code, a header or a bare reference, counting each occurrence once. (`Find`) |
| `RCRWR-B07` | Each occurrence in the dry-run listing carries the text it matched and the line it starts on, counted from one, also for headers and mentions that come after scenario codes. (`Find`) |
| `RCRWR-B08` | `RewriteRuleCodes` rewrites only the rule and scenario codes of OLD — `OLD-B08`, `OLD-S06#02`, `OLD-VR-S01`, `OLD-CT` — and never the bare code, which outside the governed files is as likely an ordinary word or part of an identifier (`DATA_URL`). |
| `RCRWR-B09` | `RewriteCited` rewrites a code only where it is cited as one — its rule and scenario codes and any name it heads joined by a hyphen (`CODE-DS-method-email`, `CODE-perm-suite`, the stem of a file named by it), the header fields that hold codes (`code:`, `file_code:`, `ref:`/`refs:`, `dep:`, `needs:`), the code in backticks and a Gherkin tag `@CODE` — and leaves the bare word; `CitedSet` (`NewCitedSet`) does it for many codes in one pass over the text (`Rewrite`, and `RewriteRuleCodes` for the rule codes alone) and lists the lines where each still stands bare, an identifier (`CODE_X`, `X_CODE`) not counted. |
| `RCRWR-B10` | An underscore is no boundary for a bare mention: `CODE_TABLE` and `X_CODE` are identifiers, not mentions of the code, and `Rewrite` leaves them. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RCRWR-I01` | A text without the old code comes back unchanged with zero replacements, so rewriting a second time after a rewrite changes nothing. | rewrites a text that does not hold the old code, and rewrites a rewritten text again, checking both are unchanged with zero replacements |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RCRWR-X01` | Never rewrites a longer code that contains the old one, nor a word that merely ends or starts with it. | A rename in mass that mutilated a neighbour code would break a unit nobody asked to touch. |

## Errors

none — the unit is pure text rewriting; a text without the code is the invariant I01, not a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
