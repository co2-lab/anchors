<!-- @anchors
  code: PRVLS
  updated_at: 2026-11-11
  layer: infra
-->
# PushReview — the reviews a push carries, as the repository keeps them

> **Code**: `PRVLS`

## Overview

`anchors push-review` writes a review of the commits a push sends; the pre-push hook reads it
back. This unit is the file between them, newest first. Each entry says when, by whom, with
which verdict, the commits it covers, the judgment gates' answers, and what was validated and
how, or what failed, with the observations. A line the machine reads opens each entry, and the
rest is the text as written:

```markdown
## 2026-11-11 12:40 UTC · ok · 56997e0f..dd5f234b
<!-- push-review at="2026-11-11T12:40:00Z" verdict="ok" by="agent:claude" -->

**Commits**

- `dd5f234b…` `p:1a2b3c4d5e6f` fix(scan): a block comment inside a line of code…

**Validated**

The evidence tests, with a probe against the project's dialect.
```

A commit is covered by its hash, or by its patch id when a rebase gave it another hash. A commit
that changes only the reviews' file and the map owes no review: those are a review's evidence,
and reviewing them would never end.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the reviews' file | entries opened by the machine line, any other text | a section with no machine line | this unit: leaves it out of the reviews (`PRVLS-B02`) |
| a retention | days, or 0 | — | `PSRSA` |

## Effects

| Effect | Description |
| --- | --- |
| `PRVLS-B01` | A review is written as an entry: a heading with its date, verdict and range, the machine line with its date, verdict and author, the commits with hash, patch id and subject, the judgment gates' answers, what was validated — or what failed — and the notes. (`Review.Entry`) |
| `PRVLS-B02` | The file is read back into its entries, in its order, each with its date, verdict, author and commits and its text as written; a section with no machine line is no review. Written again, the file holds its header and the same entries. (`Parse`, `Render`) |
| `PRVLS-B03` | Writing a review drops the entries older than the retention; with no retention, or with a date nobody can read, an entry stays. (`Keep`) |
| `PRVLS-B04` | A commit is covered by the newest entry that names it by hash, a short hash included, or by patch id; a commit no entry names is not covered. (`Covering`) |
| `PRVLS-B05` | A commit that changes only the reviews' file and the map owes no review; one that changes anything else, or nothing, owes one. (`Exempt`) |
| `PRVLS-B06` | The commits of a range are read oldest first with their patch id and subject, merges left out, and a commit's files as the paths it changes; a commit reworded keeps its patch id. (`Commits`, `Files`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRVLS-E01` | git cannot list the range: a revision it does not know, or no repository. | Reading the commits fails naming the range. | A range read as empty would let every commit of the push through as reviewed. |
| `PRVLS-E02` | git gives a commit no files or no patch id: an empty commit, or git failing on it. | Its files are none and its patch id empty: it is covered by its hash alone, and with no files it still owes a review. | A commit nobody could read is not evidence; treating it as exempt would push it unreviewed. |

