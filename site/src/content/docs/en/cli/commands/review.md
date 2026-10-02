---
title: "anchors review"
description: "Record a second look at what an agent decided, with who looked and what they found."
---

A gate that declares `review:` marks its targets **to review**, apart from how it measures. A target stays to review until a review is recorded at its current revision; a change to it makes the review due again. Reviews inform and never block.

```bash
# What is to review, grouped by gate, with the question
anchors review --pending [--gate <gate>]

# A review that found nothing
anchors review <target> --gate <gate> --by human:ana

# A review with findings: the whole report, each finding with what, where and why
anchors review <target> --gate <gate> --by agent:<vendor>/<model> --findings "<report>"
```

| Flag | What it does |
| --- | --- |
| `--pending` | lists the targets to review |
| `--gate` | the gate whose review this is — a gate that declares `review:` |
| `--by` | who reviewed, as you name them; the record keeps it and does not grade it |
| `--findings` | the findings report; it opens the target's issue for the gate |
| `--record-issues` | in manual mode, also write the issue of the findings |

The product of a review is findings, not a stamp: each one becomes a failing test and a fix, or is dismissed with its reason in the issue. There is no `waived` — a review that did not happen is still due. See [judgment and review](/docs/concepts/ai-judgment/).
