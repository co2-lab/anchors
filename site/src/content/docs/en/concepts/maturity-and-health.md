---
title: "Maturity & Project Health"
description: "How Anchors measures codebase maturity, tracks regression scores, and drives continuous architectural health."
---

How do you know if your codebase is getting healthier or accumulating invisible tech debt?

Most teams rely on gut feeling or crude metrics like code coverage. Anchors provides a comprehensive, multi-dimensional **Maturity & Health Framework**.

---

## 1. The Four Levels of Project Maturity

Anchors classifies projects into four progressive maturity tiers:

| Level | Name | Characteristics | Active Gates |
| --- | --- | --- | --- |
| **L1** | **Seeded** | Project initialized with `anchors.yaml` and baseline specs. | [`has-code`](/docs/gates/has-code/), [`header-valid`](/docs/gates/header-valid/) |
| **L2** | **Governed** | Core units governed by [The Unit](/docs/concepts/unit/); specs and tests co-located. | [`unit-complete`](/docs/gates/unit-complete/), [`spec-feature-match`](/docs/gates/spec-feature-match/) |
| **L3** | **Robust** | Mutation testing enabled; architectural boundaries enforced. | [`mutation-score`](/docs/gates/mutation-score/), [`layer-boundary`](/docs/gates/layer-boundary/) |
| **L4** | **Hardened** | Full propagation waves; security gates blocking; release freezing active. | [`contract-impact`](/docs/gates/contract-impact/), [`no-secret-leaked`](/docs/gates/no-secret-leaked/), [`freeze`](/docs/freeze/) |

---

## 2. Mutation Score vs. Line Coverage

High line coverage often gives a false sense of security. A test can execute 100% of statements without actually asserting anything.

Anchors incorporates mutation testing via the [`mutation-score`](/docs/gates/mutation-score/) gate:
- It mutates production code (e.g. changing `>` to `<`, or deleting function calls).
- If your tests still pass, the mutant survived, indicating weak test assertions.
- A high mutation score (>80%) proves that your tests genuinely catch defects.

---

## 3. Monitoring Health with the CLI

Check your project health score at any time:

```bash
# Display project health scorecard and gate summary
anchors status

# Audit all gates and list warnings
anchors audit
```
