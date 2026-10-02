---
title: "Gate: rule-uses-implemented"
description: "Verifies whether os campos que a regra diz usar aparecem e são consumidos nThe code governado."
---

> **Gate Identifier:** `rule-uses-implemented` / `uso-de-regra-implementado`  
> **Code:** `RLIMR` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Verifies whether os campos que a regra diz usar aparecem e são consumidos nThe code governado.

---

## 2. Why It Matters

Evita regras de papel: The spec diz que usa o campo, mas The code nunca o lê.

---

## 3. How It Works

Verifica nThe code fonte corresponding se os identificadores dos campos declarados na regra são referenciados.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Os campos declarados nThe spec aparecem nThe code governado. | None. Pipeline proceeds. |
| **`FAIL`** | Campos declarados nThe spec não existem nThe code. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código ainda não implementado (@TBD). | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: rule-uses-implemented
    on: [spec]
    check: rule-uses-implemented
    blocking: false
    measures: "os campos que a regra usa aparecem no código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo.
- **How to Fix:** Implemente a leitura do campo nThe code fonte ou remova o campo não utilizado dThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
