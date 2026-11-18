<!-- @anchors
  code: LCBCL
  updated_at: 2026-10-08
  layer: comando
-->
# LocalBacklog — what is still open locally after a full check, said in two lines

> **Code**: `LCBCL`

## Overview

In the local workflow mode the work lives in two folders: the issues under the issues directory (the
`todo` and `doing` states) and the task queue. A full check used to report only what that one run opened
or resolved; the work already open from earlier runs stayed silent, visible only to whoever opened the
folder, and a clean-looking check could sit on top of a backlog nobody was working on.

This unit reads that backlog and says it. It counts the issues waiting in `todo` (and, among them, how
many are owned by the user), the issues in `doing`, the tasks still pending, the tasks claimed, and among
the claimed ones how many are past the work window. It then prints a short block: a heading, one line
for the issues and one line for the tasks, each line only when its side has something open. When nothing
is open it prints nothing at all.

The unit only reads and prints. Deciding WHEN the backlog is shown (the full sweep of the local mode, and
never on an incremental check) belongs to the check command that calls it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory, with or without an issues directory and a task queue | a path that is not the project | the caller: the check command passes its resolved root |
| the issue folders | `todo` and `doing` folders holding one file per issue, or absent | a state folder that cannot be listed | this unit: an unlistable folder counts as zero (`LCBCL-E01`) |
| the task queue | task files in pending or claimed state, or no queue | tasks already done | the queue: done tasks are not in the pending or claimed states and are not counted |

## Effects

| Effect | Description |
| --- | --- |
| `LCBCL-B01` | Counts the issues in `todo` and in `doing`, and among the `todo` ones those owned by the user. |
| `LCBCL-B02` | Counts the pending tasks and the claimed tasks, and among the claimed ones those past the work window. |
| `LCBCL-B03` | When nothing is open, the backlog is empty and nothing is printed. |
| `LCBCL-B04` | When something is open, it prints a heading, the issues line only when there are issues in `todo` or `doing`, and the tasks line only when there are pending or claimed tasks. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `LCBCL-I01` | The backlog is empty exactly when `todo`, `doing`, pending and claimed are all zero; the user-owned and past-window counts are parts of those and never make it non-empty on their own. | a backlog holding only user-owned and past-window counts is empty and prints nothing |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `LCBCL-X01` | Reads the issues and the queue and never changes them: no issue moves and no task is claimed or released. | It is a report printed at the end of a check; changing the backlog while describing it would make the check the owner of work it only observes. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `LCBCL-E01` | An issue state folder (or the queue) exists but cannot be listed. | That part counts as zero and the check goes on. | The backlog is informative and never blocks: a check that failed because its closing summary could not read a folder would hide the verdict it already reached. <!-- @resilient: the backlog is an informative summary after the verdict, and an unlistable folder can only lower a count, never change what the check decided --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
