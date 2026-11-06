<!-- @anchors
  code: BCSBR
  updated_at: 2026-11-06
  layer: comando
-->
# Branch — confront merges and commits with the project's branch flow

> **Code**: `BCSBR`

## Overview

`branch_flow:` (BFSBR) says which branch merges into which, and which branches take no commit of their
own. This command confronts it where the merge happens.

On the remote, a merge into a guarded branch goes through a pull request, and a pull request names
both sides: `anchors branch check`, in a job on the pull request event, refuses a pair the flow does not
allow. A repository with no server-side protection still takes a direct push — a private repository
on GitHub's free plan applies no ruleset —, so the same command on the push event refuses a push that
is not the merge of an allowed pull request, asking GitHub for the pull request whose merge commit is
the pushed sha.

On the machine, merges are not confronted: they reach the remote as pull requests. What the hooks
guard is that a protected branch takes no commit of its own (`anchors branch guard`).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a pull request's sides | `--head`, `--base`, or `$GITHUB_HEAD_REF`, `$GITHUB_BASE_REF` | a side missing | this unit: refuses, naming what is missing |
| a push | `--push --ref --sha`, or `$GITHUB_EVENT_NAME=push` with `$GITHUB_REF_NAME`, `$GITHUB_SHA` | — | the CI |
| the pull requests of a commit | what `gh api repos/{owner}/{repo}/commits/<sha>/pulls` answers | no `gh`, no answer | this unit: refuses the push |
| the current branch | `git symbolic-ref --short HEAD`, empty on a detached HEAD | — | git |

## Effects

| Effect | Description |
| --- | --- |
| `BCSBR-B01` | `branch check` with a pull request's sides passes when the flow allows its head into its base, and fails naming the pair and the heads the base takes. (`checkMerge`) |
| `BCSBR-B02` | `branch check` on a push passes for a branch the flow does not guard, and for a guarded one only when a merged pull request into that branch, whose merge commit is the pushed sha, comes from a head the flow allows; otherwise it fails, naming the refused pull requests or that none made the push. A push it cannot confront — no `gh`, no answer — fails. (`checkPush`) |
| `BCSBR-B03` | `branch guard` fails on a branch the flow protects, naming the heads it takes; on another branch, or a detached HEAD, it passes. (`guardCommit`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BCSBR-E01` | A pull request's head or base, or a push's branch or sha, is given neither by flag nor by the CI's environment. | The check fails naming the flags and the variables. | A check with a side missing would confront nothing and pass. |
| `BCSBR-E02` | A push cannot be confronted: `gh` is not in the PATH, GitHub does not answer, or its answer is not read. | The check fails naming the sha. | A push no one could confront would pass as allowed. |
| `BCSBR-E03` | The configuration does not load. | `branch check` and `branch guard` fail with the load error. | A flow that was not read guards nothing. |
