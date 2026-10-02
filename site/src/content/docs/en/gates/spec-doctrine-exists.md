---
title: "Gate: spec-doctrine-exists"
description: "Ensures that a doutrina referenciada pelThe spec com @realizes realmente existe no projeto."
---

> **Gate Identifier:** `spec-doctrine-exists` / `doutrina-existe`  
> **Code:** `DCTRN` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a doutrina referenciada pelThe spec com @realizes realmente existe no projeto.

---

## 2. Why It Matters

Impede apontar para uma regra de produto inexistente ou renomeada.

---

## 3. How It Works

Verifies whether o identificador após `@realizes` resolve para uma regra nos arquivos de `product/*.doctrine.md`.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as doutrinas citadas existem no repositório. | None. Pipeline proceeds. |
| **`FAIL`** | Menção a regra de doutrina que não existe. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Specs sem anotação @realizes. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: spec-doctrine-exists
    on: [spec]
    check: spec-doctrine-exists
    blocking: true
    measures: "a doutrina referenciada pela spec existe no projeto"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Crie a regra nThe file de doutrina ou corrija The code na anotação @realizes.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
