<!-- @anchors
  code: CMMSC
  updated_at: 2026-09-27
  layer: comando
-->
# CommitMsg — the commit subject is confronted with the format the changelog will read, before the commit exists

> **Code**: `CMMSC`

## Overview

The changelog is born from the commits, so the format has to hold before the commit is
made: a commit already made cannot be fixed, and a history where half the subjects miss the
format produces a changelog with holes nobody can fill later. The case that motivated the
command was a squash merge whose message was the pull request title, written in the card
format; the commit that introduced a whole plan would never have reached the changelog.

`commit-msg <file>` reads the message file git hands to the commit-msg hook and confronts
its subject with Conventional Commits, the format changelog tools already read. The ruler
was confronted with commitlint, the mature tool for this format, and agrees with it in every
case but one, where it is deliberately stricter (an empty scope), so a project that later
moves to commitlint never discovers a history the new tool refuses. It deliberately leaves
out commitlint's subject-case check, which refuses legitimate acronyms at the start of the
subject.

The checks go from the most structural to the most cosmetic, and only the first defect is
reported, each with its own diagnosis: whoever got the type wrong does not need to hear
about the final period in the same round. The rejection repeats the subject, teaches the
format with examples, and lists the accepted types. Messages git writes itself pass, and so
does a message with no subject, which git refuses on its own with a better message.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the message file | a readable file, possibly with git's `#` comment lines and blank lines | a path that cannot be read | this unit: fails naming the read |
| the subject | any single line of text | — | this unit: every line is either accepted or refused with a diagnosis |

## Effects

| Effect | Description |
| --- | --- |
| `CMMSC-B01` | The subject is the first line of the file that is neither blank nor a `#` comment, trimmed. |
| `CMMSC-B02` | A subject git generates (starting with `Merge `, `Revert `, `fixup! `, `squash! ` or `Reapply `) passes; a human subject merely starting alike does not escape. |
| `CMMSC-B03` | A subject of the form `type: text`, `type(scope): text`, with an optional `!` before the colon, passes when the type is known. |
| `CMMSC-B04` | The type must be lowercase and one of feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert; anything else is refused naming the type. |
| `CMMSC-B05` | An empty scope, `type():`, is refused. |
| `CMMSC-B06` | A subject longer than the limit of 100 characters (see `CMMSC-B14`) is refused, and the diagnosis says the detail belongs in the body; a subject of exactly the limit passes. |
| `CMMSC-B07` | A subject whose text ends in a period is refused. |
| `CMMSC-B08` | Each defect has its own diagnosis, distinct from every other. |
| `CMMSC-B09` | A capital letter at the start of the subject text is allowed. |
| `CMMSC-B10` | A rejection repeats the subject, teaches the `type(scope): what changed` format with examples, and lists the accepted types. |
| `CMMSC-B11` | A message with no subject (empty, or only comments) passes. |
| `CMMSC-B12` | The checks run in order (format, lowercase type, known type, non-empty scope, non-empty text, a space after the colon, length, final period) and only the first defect is reported. |
| `CMMSC-B13` | A known type followed by nothing after the colon is refused. |
| `CMMSC-B14` | The subject limit counts characters, not bytes: a subject of 100 accented letters passes, and the diagnosis of a longer one gives its length in characters. |
| `CMMSC-B15` | A subject with no space after the colon (`feat:x`) is refused with its own diagnosis, as commitlint refuses it. |
| `CMMSC-B16` | A `Bug:` footer in the message's last paragraph marks a fix of a defect that shipped: it is refused when empty, when spelled other than `Bug:`, or on a commit that is not a `fix`; a message without it, or with a body sentence starting with "Bug:", passes. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CMMSC-I01` | REF[CMMSC-B03]: every subject commitlint accepts in the confronted cases, except the empty scope, passes here — the ruler never diverges to the looser side | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CMMSC-X01` | The command does not refuse an empty message and does not rewrite the message: it only accepts or refuses. | Git refuses an empty message with a better diagnosis; rewriting a message would put words in the author's mouth. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CMMSC-E01` | The message file cannot be read. | The command fails with "read the message" and the cause. | A hook that silently passed an unread message would let any subject through. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/ops/install_hooks.go` | the commit-msg hook | comando — the caller that runs this command on every commit |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
