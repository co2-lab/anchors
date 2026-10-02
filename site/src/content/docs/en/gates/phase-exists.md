---
title: "Gate: phase-exists"
description: "Ensures that as fases citadas nas tarefas e dependencies existem no plano."
---

> **Gate Identifier:** `phase-exists` / `fase-existe`  
> **Code:** `PHORP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that as fases citadas nas tarefas e dependencies existem no plano.

---

## 2. Why It Matters

Impede apontar dependencies para fases inexistentes ou renomeadas.

---

## 3. How It Works

Cruza referências `phase:` com a lista de identificadores de fases declaradas no plano.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as fases citadas existem. | None. Pipeline proceeds. |
| **`FAIL`** | Referência a fase inexistente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: phase-exists
    on: [plan]
    check: phase-exists
    blocking: true
    measures: "as fases citadas existem no plano"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para governança de planos.
- **How to Fix:** Crie a fase no plano ou corrija o nome da referência.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
