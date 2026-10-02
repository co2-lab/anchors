---
title: "Gate: phase-ordered"
description: "Ensures that a ordem e dependencies entre fases do plano são consistentes e sem ciclos."
---

> **Gate Identifier:** `phase-ordered` / `fase-ordenada`  
> **Code:** `PHORP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a ordem e dependencies entre fases do plano são consistentes e sem ciclos.

---

## 2. Why It Matters

Impede dependencies circulares entre fases de projeto que travam a execução.

---

## 3. How It Works

Monta o grafo de fases e valida se é um grafo acíclico direcionado (DAG) respeitando a sequência.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Fases ordenadas e sem dependencies circulares. | None. Pipeline proceeds. |
| **`FAIL`** | Ciclo de fases ou ordem temporal incoerente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: phase-ordered
    on: [plan]
    check: phase-ordered
    blocking: true
    measures: "a ordem declarada entre fases é coerente"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em planejamento.
- **How to Fix:** Ajuste as dependencies `depends_on` das fases para eliminar o ciclo.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
