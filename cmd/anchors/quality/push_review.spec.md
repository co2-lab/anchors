<!-- @anchors
  code: PRCMS
  updated_at: 2026-11-11
  layer: comando
-->
# PushReview — the review of what a push sends, recorded by whoever pushes and confronted by the pre-push hook

> **Code**: `PRCMS`

## Overview

The commit hook runs the gates a script computes. The gates that ask for judgment, and the review
that reads what a change means, ran only when someone remembered them. With `push_review:`
(PSRSA) the push carries that evidence: the agent pushing reviews what it sends, with
`anchors guide review` as its ruler, runs the judgment gates while it does, and records the
verdict:

    anchors push-review --ok "<what was validated, and how>" --by <who> [--note "<observation>"]...
    anchors push-review --fail "<what failed>" --by <who> [--note "<observation>"]...

The review goes to the file `push_review.file` names (REVIEWS.md by default), newest first, as
PRVLS writes it, and is committed alone. The pre-push hook runs `anchors push-review verify`
with the refs git hands it. A commit no ok review covers stops the push, and so does one that
reaches a judgment gate no verdict answers in the map the push sends, and — with `require_gates` — one a blocking gate fails over, whatever a commit message waived: that waiver frees the commit, not the push.

The range is what the branch has beyond its upstream, or beyond every remote branch when it has
none, unless `--range <base>..<tip>` names it. A commit that changes only the reviews' file and the
map owes no review, so the review's own commit, and the map's judgments, push without one.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the verdict | exactly one of `--ok "<text>"` and `--fail "<text>"` | neither, or both | this unit: refuses (`PRCMS-E01`) |
| `--by` | who reviewed, as the project names them | empty | this unit: refuses (`PRCMS-E01`) |
| the range | `<base>..<tip>`, or what the branch has beyond its upstream | a range with no commit that owes a review | this unit: refuses (`PRCMS-E01`) |
| the hook's input | the lines git hands the pre-push hook: local ref, local hash, remote ref, remote hash | lines of another shape | this unit: ignores them |

## Effects

| Effect | Description |
| --- | --- |
| `PRCMS-B01` | A review is written at the top of the reviews' file with its verdict, author, range, commits, the judgment gates' answers, its text and notes, and committed alone as `chore(review): <verdict> <range>`; whatever else is staged or changed stays out. `--no-commit` writes without committing. |
| `PRCMS-B02` | With `push_review.required`, a push whose commits no review covers — in the reviews' file as the pushed tip has it — is refused, naming each commit and how to record the review. |
| `PRCMS-B03` | A push whose commits' newest review failed is refused, showing that review. |
| `PRCMS-B04` | A push whose commits an ok review covers goes through; a commit that changes only the reviews' file and the map owes none. |
| `PRCMS-B05` | The refs git hands the pre-push hook say what each ref sends: the commits after the remote's tip, or, for a new branch, those no remote branch has; a deletion sends nothing. With no input, the range is what the branch has beyond its upstream. |
| `PRCMS-B06` | Without `push_review`, or with neither demand, the push confronts nothing. |
| `PRCMS-B07` | Writing a review drops those older than `push_review.retention`. |
| `PRCMS-B08` | With `push_review.require_judgment`, a push whose changed files reach a judgment gate with no verdict in the map the pushed tip holds is refused, naming each gate and target, and so is an ok review recorded over them; a verdict in a map the push does not send answers nothing. |
| `PRCMS-B09` | With `push_review.require_gates`, the blocking gates a commit or a push charges — judgment left to `require_judgment` — run over what the pushed files reach, every file read as the pushed tip holds it, with no waiver a commit message declared; each failure refuses the push naming the gate, the target and the pushed commit whose message waived it, and an edit not committed is not what the push sends. (`gatesOverPush`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRCMS-E01` | No verdict, both verdicts, no `--by`, or a range with no commit that owes a review. | The review is not written, saying what is missing. | A review without a verdict or an author proves nothing, and one of no commit covers nothing. |
