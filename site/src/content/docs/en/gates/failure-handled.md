---
title: "Gate: failure-handled"
description: "Ensures that toda falha declarada nThe spec é tratada nThe code fonte."
---

> **Gate Identifier:** `failure-handled` / `falha-tratada`  
> **Code:** `FLRAI` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that toda falha declarada nThe spec é tratada nThe code fonte.

---

## 2. Why It Matters

Não adianta documentar o erro se The code não tem `if err != nil` ou `try/catch` corresponding.

---

## 3. How It Works

Verifies whether os blocos de erro catalogados nThe spec possuem correspondência nThe code governado.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as falhas declaradas possuem tratamento nThe code. | None. Pipeline proceeds. |
| **`FAIL`** | Falha declarada nThe spec não tratada nThe code. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código pendente. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: failure-handled
    blocking: true
    measures: "toda falha declarada na spec é tratada no código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em código de produção.
- **How to Fix:** Implemente o tratamento do erro nThe code fonte.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
