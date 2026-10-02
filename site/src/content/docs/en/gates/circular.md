---
title: "Gate: circular"
description: "Detects circular dependency loops between units, packages, or architectural layers."
---

> **Gate Identifier:** `circular` / `dependencia-circular`  
> **Code:** `EXCMX` | **Category:** [Architectural Boundaries](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Detects circular dependency loops between units, packages, or architectural layers.

---

## 2. Why It Matters

Cycles prevent clean modular testing, complicate compilation, and make isolated refactoring impossible.

---

## 3. How It Works

Executes Tarjan's strongly connected components algorithm over the project dependency DAG.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | No cycles exist in the dependency graph. | None. Pipeline proceeds. |
| **`FAIL`** | A circular dependency chain is detected (e.g. A -> B -> C -> A). | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: circular
    on: [code]
    scope: project
    run: "npx madge --circular --extensions ts,tsx src/"
    needs_tool: madge
    blocking: true
    when: [ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Strictly blocking. Never allow cycles into the mainline branch.
- **How to Fix:** Break the loop by extracting shared contracts or using event-driven communication.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
