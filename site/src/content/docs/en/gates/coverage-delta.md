---
title: "Gate: coverage-delta"
description: "Ensures that a alteração atual não reduziu a cobertura de linhas em relação ao baseline."
---

> **Gate Identifier:** `coverage-delta` / `delta-de-cobertura`  
> **Code:** `INCHN` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`, `code`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a alteração atual não reduziu a cobertura de linhas em relação ao baseline.

---

## 2. Why It Matters

Impede a regressão de cobertura: uma nova PR não pode diminuir a porcentagem que o projeto já havia conquistado.

---

## 3. How It Works

Compara a cobertura atual com o baseline registrado no snapshot anterior do projeto.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Delta de cobertura >= 0. | None. Pipeline proceeds. |
| **`FAIL`** | A cobertura caiu em relação ao commit base. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Sem baseline anterior para comparar. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: coverage-delta
    on: [code]
    check: coverage-delta
    blocking: true
    measures: "a cobertura não caiu com esta alteração"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no CI para evitar degradação de testes.
- **How to Fix:** Adicione testes para o novThe code para compensar ou manter a cobertura geral.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
