---
title: "Gate: rule-uses-resolve"
description: "Ensures that o que a regra diz usar realmente existe declarado nThe spec."
---

> **Gate Identifier:** `rule-uses-resolve` / `uso-de-regra-resolve`  
> **Code:** `RLUSG` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that o que a regra diz usar realmente existe declarado nThe spec.

---

## 2. Why It Matters

Evita que uma regra afirme usar um campo 'dataNascimento' se The spec só define 'data_criacao'.

---

## 3. How It Works

Cruza a lista de campos que a regra declara usar com a tabela de domínio da própriThe spec.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os campos citados pela regra existem na tabela de campos dThe spec. | None. Pipeline proceeds. |
| **`FAIL`** | A regra afirma usar um campo que não está declarado nThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: rule-uses-resolve
    on: [spec]
    check: rule-uses-resolve
    blocking: false
    measures: "o que a regra diz usar existe na spec"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo até que o vocabulário de dados dThe spec esteja maduro.
- **How to Fix:** Declare o campo na tabela de Domínio dThe spec ou corrija o nome do campo na regra.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
