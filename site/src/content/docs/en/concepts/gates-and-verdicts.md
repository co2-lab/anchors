---
title: "Gates & Verdicts"
description: "Understand how Anchors gates evaluate code, the four verdict states, and how to configure blocking policies."
---

In mountaineering, safety equipment is useless if it doesn't hold when subjected to stress. In Anchors, **Gates** are the automated inspection points that verify whether your project's code, specs, tests, and architecture adhere to declared standards.

---

## 1. What is a Gate?

A **Gate** is an automated checker that evaluates an artifact or relationship in the project graph against an invariant rule.

Gates are configured in [`anchors.yaml`](/docs/anchors-yaml/):

```yaml
gates:
  - name: unit-complete
    check: unit-complete
    on: [spec]
    blocking: true
    measures: "spec has matching feature, test, and code"

  - name: mutation-score
    check: mutation-score
    on: [code]
    threshold: 80
    blocking: true
    measures: "test suite kills at least 80% of generated mutants"
```

Each gate declares:
- **`name`**: The human-readable identifier.
- **`check`**: The underlying verification algorithm.
- **`on`**: The artifact types to evaluate (`spec`, `code`, `feature`, `test`, `plan`, `doctrine`, etc.).
- **`blocking`**: Whether a failure stops the CI/CD pipeline.
- **`threshold`**: Optional quantitative criteria (e.g. coverage percentage or mutation score).

---

## 2. The Four Verdict States

When a gate executes, it returns one of four definitive **Verdicts**:

| Verdict | Meaning | CI Impact | What to do |
| --- | --- | --- | --- |
| **`OK`** | All criteria are fully satisfied. | Passes cleanly. | Proceed to next step. |
| **`FAIL`** | An invariant was violated. | **Blocks pipeline** (if `blocking: true`). | Fix the code or update the spec to resolve the discrepancy. |
| **`WARN`** | Advisory discrepancy or impending deprecation. | Logs warning; does not block CI. | Schedule remediation before the gate becomes blocking. |
| **`SKIP`** | Gate was skipped due to an explicit layer waiver or filter. | Neutral. | Recorded in audit log for transparency. |

---

## 3. Deterministic vs. Relational vs. Synthetic Gates

Anchors categorizes its 95 gates into three execution models:

1. **Deterministic Internal Gates**: Inspect a single file's syntax and contents without requiring graph context. (e.g. [`has-code`](/docs/gates/has-code/), [`header-valid`](/docs/gates/header-valid/)).
2. **Relational Graph Gates**: Inspect relationships and edges across multiple files in the dependency graph. (e.g. [`unit-complete`](/docs/gates/unit-complete/), [`layer-boundary`](/docs/gates/layer-boundary/), [`circular`](/docs/gates/circular/)).
3. **Synthetic / AI-Judged Gates**: Use Large Language Models with calibrated prompts to evaluate semantic consistency and clarity. (e.g. [`doc-self-contained`](/docs/gates/doc-self-contained/), [`progress-honest`](/docs/gates/progress-honest/)).

---

## 4. Explore All Gates

Browse all 95 available verification gates in the [Gates Catalog](/docs/gates/).
