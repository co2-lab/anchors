<!-- @anchors
  code: GHRNG
  updated_at: 2026-09-26
  layer: apoio
-->
# GhRunner — running `gh` for the board so that a failure names its cause and never waits for input

> **Code**: `GHRNG`

## Overview

Every conversation with the board goes through `gh`, an external process Anchors does not control. On
some paths `gh` WAITS instead of failing, for instance on an authentication prompt; measured, a new
developer ran a command without `gh auth login`, it hung, and had to be killed. A command that never
returns is worse than one that fails: the one that fails says what to do.

So `gh` is run with no input: a `gh` that decides to ask for a credential finds the input closed and
gives up with its own message instead of waiting. And the failure `gh` reports without credentials,
exit code 4, says nothing by itself (a new developer got `exit status 4` and had no way to know login
was missing), so that code, and only that code, gets the hint that names the cause and the fix.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the arguments | the arguments of one `gh` call | — | the board client (`BRCRB`) |

## Effects

| Effect | Description |
| --- | --- |
| `GHRNG-B01` | `gh` runs with no input: a `gh` that reads its input gets its end at once instead of waiting. |
| `GHRNG-B02` | A `gh` that exits with code 4 fails with the hint that it is NOT AUTHENTICATED, that `gh auth login` fixes it, and that the login is interactive. |
| `GHRNG-B03` | Any other failure, another exit code or a command that could not start, carries no authentication hint. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GHRNG-E01` | REF[GHRNG-B02]: `gh` without credentials is the failure B02 explains with its fix | — | — |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
