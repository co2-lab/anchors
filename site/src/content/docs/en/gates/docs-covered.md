---
title: "Gate: docs-covered"
description: "Ensures that todThe spec alcança alguma página da documentação compilada."
---

> **Gate Identifier:** `docs-covered` / `docs-cobertos`  
> **Code:** `DCCVD` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that todThe spec alcança alguma página da documentação compilada.

---

## 2. Why It Matters

Impede specifications esquecidas que nunca aparecem no portal de documentação do projeto.

---

## 3. How It Works

Verifies whether The spec está indexada no sumário ou em alguma página gerada da documentação.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The spec aparece na documentação compilada. | None. Pipeline proceeds. |
| **`FAIL`** | Spec não referenciada na documentação. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: docs-covered
    blocking: false
    measures: "toda spec alcança alguma página da documentação compilada"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo para publicação de portais de documentação.
- **How to Fix:** Inclua The spec no índice de documentação ou execute `anchors docs build`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
