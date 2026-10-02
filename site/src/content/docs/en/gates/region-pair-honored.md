---
title: "Gate: region-pair-honored"
description: "Ensures that todo bloco #region nThe code fecha com um #endregion carregando o mesmThe code."
---

> **Gate Identifier:** `region-pair-honored` / `par-de-region-honrado`  
> **Code:** `RPHRG` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo bloco #region nThe code fecha com um #endregion carregando o mesmThe code.

---

## 2. Why It Matters

Impede regiões de código abertas sem fechamento que quebram o parser de rastreabilidade de código.

---

## 3. How It Works

Examina marcadores `#region CODE` e valida se cada abertura tem seu `#endregion CODE` com o mesmThe code.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as regiões fecham com o mesmThe code corresponding. | None. Pipeline proceeds. |
| **`FAIL`** | Região aberta sem fechar, ou código do fechamento diverge da abertura. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Arquivos sem marcadores #region. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: region-pair-honored
    on: [code]
    check: region-pair-honored
    blocking: true
    measures: "todo #region fecha com #endregion de mesmo código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Feche a região nThe code com `// #endregion CODE`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
