---
title: "Gate: scenario-type-aligned"
description: "Ensures that as tags de classificação do scenario batem com a letra dThe code (ex: @unit para regra B)."
---

> **Gate Identifier:** `scenario-type-aligned` / `tipo-cenario-alinhado`  
> **Code:** `STASC` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that as tags de classificação do scenario batem com a letra dThe code (ex: @unit para regra B).

---

## 2. Why It Matters

Impede misturar tipos de teste (ex: marcar um scenario visual como unitário ou vice-versa).

---

## 3. How It Works

Compara as tags de regime e escopo do scenario com a letra central dThe code (B, V, E, VR, DS).

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Tags e letras de código estão alinhadas conforme as convenções do projeto. | None. Pipeline proceeds. |
| **`FAIL`** | Inconsistência entre a tag do scenario e a letra dThe code. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | scenarios sem classificação especial. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: scenario-type-aligned
    on: [feature]
    check: scenario-type-aligned
    blocking: true
    measures: "a tag do cenário bate com a letra do código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Ative como informativo e promova a blocking ao padronizar as tags.
- **How to Fix:** Ajuste a tag do scenario (ex: use @unit para regras B e @vr para regras VR).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
