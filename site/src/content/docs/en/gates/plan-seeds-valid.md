---
title: "Gate: plan-seeds-valid"
description: "Ensures that as specs semeadas pelo plano miram camadas governadas válidas."
---

> **Gate Identifier:** `plan-seeds-valid` / `sementes-de-plano-validas`  
> **Code:** `PSVPL` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that as specs semeadas pelo plano miram camadas governadas válidas.

---

## 2. Why It Matters

Impede um plano de semear código em pastas inexistentes ou fora da planta da casa.

---

## 3. How It Works

Confere os caminhos das specs semeadas no plano contra as regras de camadas do anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as specs semeadas miram camadas governadas válidas. | None. Pipeline proceeds. |
| **`FAIL`** | Spec semeada mirando camada inexistente ou proibida. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: plan-seeds-valid
    on: [plan]
    check: plan-seeds-valid
    blocking: true
    measures: "as specs semeadas pelo plano existem nas camadas certas"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking ao semear novas fases.
- **How to Fix:** Ajuste o caminho dThe spec no plano para casar com uma camada declarada no `anchors.yaml`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
