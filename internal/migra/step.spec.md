<!-- @anchors
  code: MSCMG
  updated_at: 2026-09-28
  layer: apoio
-->
# MigrationStepChain — one step per format version, chained in order, and a hole in the chain is an error

> **Code**: `MSCMG`

## Overview

The migrator does not know how to convert "from the old format to the current one" — it knows how
to convert from N to N+1, and applies those steps in sequence. A project stuck on format 1 with a
binary on format 4 crosses 1→2, 2→3 and 3→4, in that order.

A single function that looks at the state and repairs it looks simpler and ages badly: it must
recognise every possible intermediate state, and each new format adds one more combination to tell
apart. An isolated step only knows ONE transition; once written it never changes, and the next
format is a new file that does not touch it. It is the same reason database migrations are written
this way: whoever writes a step does not know the starting state of whoever upgrades.

This unit is the registry of those steps and the chain that links them. Each step declares the
format it produces, one line saying what changed (shown to whoever reviews the migration), the keys
renamed per file and the values renamed per file and key. Each format registers its own step from
its own file, in whatever order the files load; the registry keeps them sorted. Asking for the
steps between two formats returns exactly the ones in between, in order — and a missing step is an
error, never silence: a project on format 2 cannot go to 4 pretending 3 never existed, or the file
would get the new number with the old content.

The registry also answers whether a key is one some step renames. The configuration uses it to
explain an unknown key: an old-format key the migration renames is fixed by migrating, while a real
typo is not, and telling the two apart keeps the error message trustworthy.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a registered step | a step with the format it produces and its renames | two steps producing the same format | each format file registers exactly one step for its own number (see the migration steps spec) |
| the origin and target formats | two format numbers, the origin at most the target | an origin above the target | the caller only asks when the file is below the target |
| the key asked about | any key name | — | this unit answers for any text |

## Effects

| Effect | Description |
| --- | --- |
| `MSCMG-B01` | `Register` keeps the steps sorted by the format they produce, whatever the order they were registered in. |
| `MSCMG-B02` | `StepsFrom`, asked for the steps from format A to format B, returns the steps producing A+1 through B, in ascending order. |
| `MSCMG-B03` | A file already at the target format needs no step: the answer is an empty list, not an error. |
| `MSCMG-B05` | A step may declare code letters renamed by kind of unit, which it does not apply to the YAML files itself (`RenameLetters`, `LetterRename`). |
| `MSCMG-B04` | `RenamedKey` reports a key as renamed when any registered step renames it in any file; any other key is not. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MSCMG-I01` | `AllSteps`, listing the registered steps, hands out a copy: changing the list a caller received never changes the registry. | alters the listed steps and lists again, verifying the registry is intact |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MSCMG-X01` | Steps outside the requested interval are never returned: neither those producing the origin or earlier formats, nor those producing formats beyond the target. | Applying a step twice, or one the target does not reach, would convert a file past the format the binary understands. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MSCMG-E01` | A step inside the requested interval is missing — the next registered step produces a later format. | An error naming the format whose step is missing. | Skipping it would give the file the new number with the old content; the message names the step so whoever reads knows what to write. |
| `MSCMG-E02` | The target is beyond the last registered step. | An error naming the first format that has no step. | The same hole at the end of the chain: the file cannot claim a format no step produces. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the registry holds the steps the format files register at load time and depends on nothing.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
