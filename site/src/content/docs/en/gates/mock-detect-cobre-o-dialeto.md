---
title: "Gate: mock-detect-cobre-o-dialeto"
description: "Evaluates whether a regex declarada para detectar dublês alcança todas as formas que o projeto usa."
---

> **Gate Identifier:** `mock-detect-cobre-o-dialeto` / `mock-detect-cobre-dialeto`  
> **Code:** `MCSTM` | **Category:** [Test Doubles & Mocks](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `governed layers`, `test`  
> **Execution Model:** `Julgamento por IA (ask)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Evaluates whether a regex declarada para detectar dublês alcança todas as formas que o projeto usa.

---

## 2. Why It Matters

Se a regex de mock-detect falhar, dublês passam sem validação de contrato.

---

## 3. How It Works

Pergunta ao modelo se existem padrões de mock no repositório que escaparam da regex configurada.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Veredito `pass` emitido por julgamento. | None. Pipeline proceeds. |
| **`FAIL`** | Veredito `fail` acusando padrões de mock não cobertos. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Julgamento não registrado. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: mock-detect-cobre-o-dialeto
    on: [test]
    ask: "o padrão de detecção de mocks alcança todas as formas usadas no repositório?"
    blocking: false
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo; avalie periodicamente com `anchors judge`.
- **How to Fix:** Ajuste o regex `mock_detect` no anchors.yaml e emita novo veredito.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
