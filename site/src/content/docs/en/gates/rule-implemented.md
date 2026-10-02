---
title: "Gate: rule-implemented"
description: "Ensures that umThe spec cataloga regras e The code mostra que as realizou."
---

> **Gate Identifier:** `rule-implemented` / `regra-implementada`  
> **Code:** `RLIMR` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that umThe spec cataloga regras e The code mostra que as realizou.

---

## 2. Why It Matters

Impede regras de fachada que constam nThe spec mas não existem nThe code fonte.

---

## 3. How It Works

Verifies whether há marcações de código ou símbolos corresponding nThe file de implementação.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | As regras catalogadas estão materializadas nThe code. | None. Pipeline proceeds. |
| **`FAIL`** | Regra catalogada sem implementação nThe code corresponding. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código pendente (@TBD). | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: rule-implemented
    blocking: true
    measures: "a spec cataloga regras e o código mostra que as realizou"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em código pronto.
- **How to Fix:** Implemente a regra nThe code ou declare @TBD: code se a entrega for em outra fase.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
