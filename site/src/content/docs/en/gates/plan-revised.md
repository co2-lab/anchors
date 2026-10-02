---
title: "Gate: plan-revised"
description: "Garante visibilidade mútua de revisão entre planos substituídos e planos revisores."
---

> **Gate Identifier:** `plan-revised` / `plano-revisado`  
> **Code:** `PLRVP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Garante visibilidade mútua de revisão entre planos substituídos e planos revisores.

---

## 2. Why It Matters

Impede revisões fantasmas: se a revisão 2 substitui a 1, a 1 deve apontar para a 2 e vice-versa.

---

## 3. How It Works

Confere os links bidirecionais de revisão entre os planos.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Revisões mútuas explicitamente vinculadas. | None. Pipeline proceeds. |
| **`FAIL`** | Plano revisor sem link para o plano original ou vice-versa. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Planos originais sem revisões. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: plan-revised
    on: [plan]
    check: plan-revised
    blocking: true
    measures: "a revisão está numerada e vinculada mutuamente"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para governança de planos.
- **How to Fix:** Adicione as anotações mútuas de revisão (`revises:` e `superseded_by:`) nos planos.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
