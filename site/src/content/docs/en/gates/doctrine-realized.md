---
title: "Gate: doctrine-realized"
description: "Ensures that as regras da doutrina de produto foram concretizadas em código e specs do projeto."
---

> **Gate Identifier:** `doctrine-realized` / `doutrina-realizada`  
> **Code:** `DCTRN` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `product` | **Layers:** `produto`, `doutrina`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that as regras da doutrina de produto foram concretizadas em código e specs do projeto.

---

## 2. Why It Matters

Evita doutrinas de produto no papel que nenhuma funcionalidade do sistema cumpre na prática.

---

## 3. How It Works

Verifies whether cada regra dThe file de doutrina possui ao menos umThe spec local com a anotação `@realizes`.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Toda regra da doutrina é realizada por ao menos uma unidade. | None. Pipeline proceeds. |
| **`FAIL`** | Regra de doutrina órfã sem nenhuma unidade que a realize. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: doctrine-realized
    on: [product]
    check: doctrine-realized
    blocking: false
    measures: "a doutrina de produto é realizada pelas unidades"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo até que todo o produto esteja coberto por specs.
- **How to Fix:** Adicione `@realizes CODIGO-DOUT` nas specs das unidades que implementam a regra.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
