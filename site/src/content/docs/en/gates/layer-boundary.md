---
title: "Gate: layer-boundary"
description: "Enforces architectural import limits and dependency directions between project layers."
---

> **Gate Identifier:** `layer-boundary` / `fronteira-de-camada`  
> **Code:** `LYBNL` | **Category:** [Architectural Boundaries](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Enforces architectural import limits and dependency directions between project layers.

---

## 2. Why It Matters

Prevents architectural decay, such as domain entities importing SQL drivers or UI components calling repositories directly.

---

## 3. How It Works

Inspects import statements in AST and cross-references them with declared layer boundaries in anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | All imports respect the dependency direction configured in anchors.yaml. | None. Pipeline proceeds. |
| **`FAIL`** | A file imports a module from a forbidden or higher-level layer. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Layers with no defined boundary constraints. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: layer-boundary
    blocking: true
    measures: "as fronteiras entre camadas, verificadas por import real"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Always blocking in CI. Protects architectural integrity across AI pair programming sessions.
- **How to Fix:** Refactor the import by injecting an interface or moving shared types to a common layer.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
