---
title: "Gate: unit-complete"
description: "Verifies that every governed spec has its complete unit: code, feature, and test."
---

> **Gate Identifier:** `unit-complete` (formerly `triad-complete`; `anchors migrate` renames it)  
> **Code:** `UNTCP` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`, `usecase`, `service`, `domain`, `comando`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies that every governed spec has its complete unit: code, feature, and test.

---

## 2. Why It Matters

Relational gates fail open if there are no tests or code linked. A standalone spec would pass all downstream checks silently. This gate closes that hole by requiring material artifacts.

---

## 3. How It Works

Queries the spec node in the project graph and verifies that outbound 'governs' (to code), 'covered-by' (to feature), and 'tested-by' (to test) edges exist.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The spec is bound to a matching code file, feature file, and test file. | None. Pipeline proceeds. |
| **`FAIL`** | Any piece of the unit is missing, or the relationship is uncataloged. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | The dependency graph is not yet initialized or an artifact is tagged @TBD. | Review before the next release cycle. |
| **`SKIP`** | The node belongs to a Recognized layer with declarative regime. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: unit-complete
    on: [spec]
    check: unit-complete
    blocking: true
    measures: "a spec tem código, feature e teste que a realizam"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Enable as blocking as soon as your first units are seeded. It is the most fundamental gate in Anchors.
- **How to Fix:** Create the missing files (e.g. the feature or test file) in the same directory, or declare an explicit documented waiver (@no-test: reason).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
