---
title: "Gate: flag-scenario-governs"
description: "Ensures that a flag governa as regras corretas nThe specification."
---

> **Gate Identifier:** `flag-scenario-governs` / `flag-governa-regras`  
> **Code:** `FLSCF` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `flag` | **Layers:** `flags`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a flag governa as regras corretas nThe specification.

---

## 2. Why It Matters

Evita flags órfãs que existem no repositório mas não condicionam nenhuma regra real.

---

## 3. How It Works

Checks whether a flag é referenciada por ao menos umThe spec através do grafo.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A flag governa ao menos uma regra no projeto. | None. Pipeline proceeds. |
| **`FAIL`** | Flag criada mas nunca usada em nenhumThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Flags novas em rascunho. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: flag-scenario-governs
    on: [flag]
    check: flag-scenario-governs
    blocking: true
    measures: "a flag governa as regras corretas na spec"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a árvore de flags limpa.
- **How to Fix:** Adicione `@gated-by` na regra corresponding ou apague a flag desnecessária.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
