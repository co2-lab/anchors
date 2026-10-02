---
title: "Gate: no-test-prova-real"
description: "Evaluates whether a justificativa apontada para uma dispensa @no-test é real e legítima."
---

> **Gate Identifier:** `no-test-prova-real` / `no-test-prova-real`  
> **Code:** `RLUEX` | **Category:** [AI Judgment & Synthetic](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Julgamento por IA (ask)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Evaluates whether a justificativa apontada para uma dispensa @no-test é real e legítima.

---

## 2. Why It Matters

Impede desenvolvedores de usarem `@no-test: testado em outro lugar` quando na verdade ninguém testou.

---

## 3. How It Works

A IA inspeciona o alvo apontado pela dispensa e valida se a prova realmente ocorre lá.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Veredito `pass` confirmando que a prova é real. | None. Pipeline proceeds. |
| **`FAIL`** | Veredito `fail` acusando que a justificativa é falsa ou insuficiente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Julgamento pendente. | Review before the next release cycle. |
| **`SKIP`** | Specs sem anotação @no-test. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: no-test-prova-real
    on: [spec]
    ask: "a prova apontada pelo @no-test exercita mesmo o comportamento?"
    blocking: false
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo; use para auditar dispensas de teste.
- **How to Fix:** Escreva The test corresponding ou aponte a referência exata da suíte onde a prova ocorre.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
