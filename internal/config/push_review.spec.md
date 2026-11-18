<!-- @anchors
  code: PSRSA
  updated_at: 2026-11-11
  layer: config
-->
# PushReview — the evidence each push carries: a review of its commits and the judgments they reach

> **Code**: `PSRSA`

## Overview

The commit hook runs the gates a script computes; the gates that ask for judgment, and the review
that reads what a change means, ran only when someone remembered them. `push_review:` makes the
push carry that evidence. The agent pushing reviews the commits it sends, runs the judgment gates
while it does, and records the verdict with `anchors push-review`. The pre-push hook confronts it:

```yaml
push_review:
  required: true          # each pushed commit is covered by a review whose verdict is ok
  file: REVIEWS.md        # where the reviews are written (the default)
  retention: 30d          # writing a review drops those older; empty keeps them all
  require_judgment: true  # each judgment gate the pushed commits reach has a verdict in the map
  require_gates: true     # the blocking gates pass over what the pushed commits reach, no waiver
```

The demands are independent: a project may ask for the judgments or the gates without a review,
or the reverse. A commit-message waiver (`[skip-<gate>@<CODE>: why]`) frees a commit in the
middle of the work; with `require_gates` it is settled before the push, and a waiver meant to
reach the remote is one the artifact carries (`@no-*`), which the gate reads. Without the section, the push confronts neither.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `push_review` | `required`, `require_judgment` and `require_gates` as booleans, a `file` inside the project, a `retention` in days (`30d`) | a `retention` that is not days, a `file` that is absolute or leaves the project | this unit: refuses the configuration (`PSRSA-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `PSRSA-B01` | The section declares whether a push needs a review, whether it needs the judgments, whether it needs the blocking gates to pass, where the reviews are written — `REVIEWS.md` when it names no file — and how many days a review stays; with no retention every review stays, and with neither demand the push confronts nothing. (`PushReview`, `FileOrDefault`, `RetentionDays`, `Enforced`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PSRSA-E01` | A `retention` that is not a number of days, or a `file` that is absolute or outside the project. | The configuration does not load, naming the field. | A retention nobody can read would keep or drop reviews by accident, and a file outside the project is not pushed with the commits it proves. |
