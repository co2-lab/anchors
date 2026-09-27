<!-- @anchors
  code: RCRWR
  updated_at: 2026-09-26
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
| `RCRWR-B05` | A bare mention of the old code is replaced, and the character right after it is kept. |
| `RCRWR-B06` | The dry-run listing classifies each occurrence as a scenario code, a header or a bare reference, counting each occurrence once. (`Find`) |

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

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `CodeLengthPattern` | config — the lengths a code may have come from the project |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
