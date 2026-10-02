---
title: "Gate: contract-impact"
description: "Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes."
---

> **Gate Identifier:** `contract-impact` / `impacto-de-contrato`  
> **Code:** `CTRIM` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes.

---

## 2. Why It Matters

Ensures that mudanças de contrato disparem exatamente os testes que dependem do campo alterado.

---

## 3. How It Works

Lê o git diff, identifica campos de contrato alterados e cruza com a tabela de usos de regras.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as regras afetadas pela alteração de campo tiveram seus testes executados. | None. Pipeline proceeds. |
| **`FAIL`** | Campos de contrato mudaram e regras dependentes não foram revalidadas. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Sem alteração de contrato no commit. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: contract-impact
    on: [spec]
    check: contract-impact
    blocking: false
    measures: "campo alterado nomeia as regras que o usam e seus testes"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo durante desenvolvimento de APIs e schemas.
- **How to Fix:** Execute os testes das regras impactadas com `anchors test`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
