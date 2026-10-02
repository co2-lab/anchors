---
title: "Gate: doc-self-contained"
description: "Ensures that The spec é autossuficiente e compreensível sem depender de contextos orais ou implícitos."
---

> **Gate Identifier:** `doc-self-contained` / `doc-autossuficiente`  
> **Code:** `DSCDC` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that The spec é autossuficiente e compreensível sem depender de contextos orais ou implícitos.

---

## 2. Why It Matters

Uma IA ou desenvolvedor novo deve conseguir entender a regra apenas lendo The spec.

---

## 3. How It Works

Verifies whether todos os termos e dependencies citados estão definidos ou linkados.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The spec é autossuficiente. | None. Pipeline proceeds. |
| **`FAIL`** | Spec com referências opacas não explicadas. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: doc-self-contained
    on: [spec]
    check: doc-self-contained
    blocking: true
    measures: "a spec se sustenta sem depender de contexto implícito"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking parThe specifications de produto.
- **How to Fix:** Adicione as definições dos termos utilizados na seção de Visão Geral dThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
