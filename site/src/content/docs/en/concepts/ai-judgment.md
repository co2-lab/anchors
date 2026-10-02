---
title: "AI Judgment & Synthetic Gates"
description: "How Anchors uses calibrated LLM-as-judge evaluations to verify semantic clarity, documentation honesty, and edge cases."
---

Some aspects of software engineering cannot be verified by regex or abstract syntax tree (AST) analysis:
- *"Does this user guide genuinely explain how to recover from payment errors, or is it just fluff?"*
- *"Is the plan description honest about the current completion status?"*
- *"Are the requirements in this spec ambiguous or self-contradictory?"*

Anchors introduces **Synthetic / AI-Judged Gates** to evaluate semantic quality.

---

## 1. How AI Judgment Works in Anchors

Synthetic gates run automated evaluations using calibrated Large Language Models:

```
  [ Target Artifact ] ──┐
                        ├─► [ Deterministic Rubric ] ──► [ LLM Evaluator ] ──► [ Verdict ]
  [ System Context  ] ──┘                                                        (OK / FAIL)
```

To ensure evaluation reproducibility:
1. **Strict Rubrics**: The prompt provides a deterministic checklist with binary criteria.
2. **Zero-Temperature Evaluation**: Models are invoked with `temperature: 0` for consistent verdicts.
3. **Structured Outputs**: The evaluator returns JSON containing the verdict, citation of evidence, and remediation suggestions.

---

## 2. Key Synthetic Gates in Anchors

- [`doc-self-contained`](/docs/gates/doc-self-contained/): Evaluates whether a documentation page is clear and understandable without requiring tribal knowledge.
- [`progress-honest`](/docs/gates/progress-honest/): Cross-references committed code and test passing rates against checklist completion claims.
- [`open-questions-resolved`](/docs/gates/open-questions-resolved/): Ensures no unanswered "TODO" questions remain in milestone specs.
- [`spellcheck`](/docs/gates/spellcheck/): Evaluates terminology consistency and grammar in customer-facing specs.

---

## 3. Judgment and Review

A **judgment** answers a gate's question with a verdict that stamps the map. When the agent that wrote the code judges its own work, a `pass` guarantees little — but a `fail` still finds real bugs. A **review** is the other thing: a second look at what an agent decided, whose product is **findings**, not a stamp.

Any gate may declare `review:`, apart from how it measures, so a gate can be judged *and* reviewed. A target is **to review** until a review is recorded at its current revision; a change to it makes the review due again. Reviews inform and never block.

```yaml
gates:
  - name: rule-fulfilled
    review:
      ask: "Does each marked snippet do what its rule says?"
```

```bash
anchors review --pending                                   # what is to review, and the question
anchors review src/pay.go --gate rule-fulfilled --by human:ana
anchors review src/pay.go --gate rule-fulfilled --by agent:<vendor>/<model> \
  --findings "### 1. B02 refunds twice (src/pay.go:41) — ..."
```

`--by` records who reviewed, as they name themselves. Findings open the target's issue, the way a failing judgment does: each one becomes a failing test and a fix, or is dismissed with its reason. There is no `waived` — a review that did not happen is still due. `anchors check` shows `🔍 N target(s) to review` on a line of its own.

