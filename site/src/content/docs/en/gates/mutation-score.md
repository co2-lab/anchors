---
title: "Gate: mutation-score"
description: "Measures test suite defect-catching quality by generating code mutants and checking if tests kill them."
---

> **Gate Identifier:** `mutation-score` / `escore-de-mutacao`  
> **Code:** `PRJTS` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Measures test suite defect-catching quality by generating code mutants and checking if tests kill them.

---

## 2. Why It Matters

Line coverage can easily be gamed by executing code without asserting behavior. Mutation testing mathematically proves test effectiveness.

---

## 3. How It Works

Injects syntactic mutations into source code (e.g. changing > to <, or replacing return values) and runs tests. If all tests pass, the mutant survived.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Mutant kill ratio meets or exceeds the configured threshold (e.g. >= 80%). | None. Pipeline proceeds. |
| **`FAIL`** | Mutation score falls below the required threshold. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Tests fail on pristine code prior to mutation. | Review before the next release cycle. |
| **`SKIP`** | Layers tagged as declarative or non-executable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: mutation-score
    on: [code]
    check: mutation-score
    blocking: false
    measures: "se a linha mudasse, algum teste quebraria"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Configure at 70% during initial development, raising to 80-85% as units mature.
- **How to Fix:** Add missing assertions that specifically inspect return values, edge conditions, and error cases.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
