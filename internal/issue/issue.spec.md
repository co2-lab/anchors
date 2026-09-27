<!-- @anchors
  code: ISLFS
  updated_at: 2026-09-26
  layer: infra
-->
# IssueLifecycle — a divergence recorded so it survives the session, with its state as a folder

> **Code**: `ISLFS`

## Overview

An issue is the record of a divergence that must outlive the session that found it: a gate violation, a
stale edge (one end moved ahead), a conflict between blocking anchors, or an open decision a spec admits
it has not taken. Locally an issue is a markdown file in the project's `issues/` folder, and its state is
the folder it lives in: `future/` (an assumed debt, due later), `todo/`, `doing/`, `done/`. Moving the
file is changing the state; there is no status field that could disagree with the folder.

An issue is found by a stable key made of its kind, its gate and its edge, never by its date: the same
problem detected on two different days is the same issue. That is what makes opening idempotent across
every state (a problem already being handled or already solved is not opened again) and what lets the
confrontation close the issue on its own when it passes again. A different finding on the same edge is
not the same detection: it reopens the issue and appends the new report, so a second verdict never
disappears in silence.

Each issue has an owner, independent of its kind: the agent, who resolves it by changing the repository,
or the user, who must decide or answer something the agent cannot. An issue can change hands without
ceasing to be what it is, and the reason travels with it. A full check closes the agent's violations it
no longer reproduces, which is how issues of a renamed gate, a fixed target or a deleted file stop
lingering. When the project works on GitHub, the same lifecycle is routed to cards (`GHIGT`).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the issue | a kind, a target, optionally an anchor and a gate, a detail, a date, optionally a due moment and an owner | — | the gates and commands that open issues |
| the date | year-month-day, stamped by the caller | — | the caller; this unit never reads the clock |
| the owner | agent or user, or empty | — | this unit: empty means the agent |
| the alive keys | the violation keys a FULL check reproduced | keys from a partial check | the caller (`anchors check --all`) |

## Effects

| Effect | Description |
| --- | --- |
| `ISLFS-B01` | The key of an issue is its kind, then its gate (when there is one), then its edge (the target, or the anchor, `--vs--` and the target), each path with its slashes turned into hyphens; it does not depend on the date. |
| `ISLFS-B02` | The file name of an issue is its date, two hyphens and its key, with `.md`, and holds no slash. |
| `ISLFS-B03` | The body names the kind and the target in its title and lists the kind, the anchor (when there is one), the target, the gate (when there is one), the owner and the detection date, then the detail under a heading of the kind (left out when the detail brings its own headings). |
| `ISLFS-B04` | Opening an issue writes it to `todo/`, where the issues of that state are listed by file name (`List`). |
| `ISLFS-B05` | Opening an issue whose key already exists in any state creates nothing and answers the state it is in, whatever the dates. (`Exists`) |
| `ISLFS-B06` | An assumed debt is opened in `future/`, and its body says when it will be paid and that it is an assumed debt. (`OpenAt`) |
| `ISLFS-B07` | The body of an open decision closes with how to close it (promote the answer to a rule) and what not to do (delete the item without a rule). |
| `ISLFS-B08` | Resolving moves a live issue (future, todo or doing) to `done/`; resolving an issue already done, or not found, does nothing. (`Resolve`) |
| `ISLFS-B09` | Reopening with a different finding moves the issue to `todo/` and appends the new report under an "Additional finding" heading, keeping the old one; reopening with a finding the issue already holds does nothing. |
| `ISLFS-B10` | An issue with no owner is the agent's, including an issue written before the owner field existed; listing by owner gives only the issues of that owner. (`OwnerOf`, `FileOwner`, `ListByOwner`) |
| `ISLFS-B11` | Reassigning changes the owner line of an open issue, keeps the rest of the file (kind included), and appends who it became and why; an issue without an owner line, even with the old Portuguese labels, gets one. |
| `ISLFS-B12` | Reconciling closes every open (todo or doing) violation whose key a full check did not reproduce, and answers the files it closed. (`ReconcileViolations`) |
| `ISLFS-B13` | By default issues are files; configuring GitHub routes the lifecycle to the declared repository's cards. (`UseGitHub`, `UseFiles`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `ISLFS-I01` | An issue lives in exactly one state folder: resolving moves it, it never copies it, and a resolved issue is never resurrected by opening it again. | opens, resolves and reopens-by-open an issue, checking the folders after each step |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ISLFS-X01` | Reconciling never closes an issue that is not a violation, nor a violation owned by the user. | A decision, a debt, a stale or conflict issue is not born of a gate failure, and an issue waiting for the user is theirs to close. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ISLFS-E01` | The state folder cannot be created when an issue is opened. | Opening returns the error and reports nothing created. | A gate that accuses and does not record would say "recorded" about a finding nobody will find tomorrow. |
| `ISLFS-E02` | The issue to reassign does not exist in the given state. | Reassigning returns the error. | Handing over an issue that is not there would tell the user they have work they do not. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/issue/github.go` | `GitHub` | infra — the card backend of the same lifecycle (`GHIGT`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
