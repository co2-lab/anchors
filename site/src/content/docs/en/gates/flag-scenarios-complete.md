---
title: "Gate: flag-scenarios-complete"
description: "Ensures that todos os valores possíveis da feature flag possuem scenarios documentados."
---

> **Gate Identifier:** `flag-scenarios-complete` / `cenarios-de-flag-completos`  
> **Code:** `FLSCF` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `flag` | **Layers:** `flags`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todos os valores possíveis da feature flag possuem scenarios documentados.

---

## 2. Why It Matters

Uma flag booleana não pode documentar apenas o caso ON e esquecer o comportamento do OFF.

---

## 3. How It Works

Checks whether todos os estados possíveis daquele tipo de flag (ex: ON e OFF) estão presentes.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os ramos da flag estão catalogados. | None. Pipeline proceeds. |
| **`FAIL`** | Falta documentar um dos estados da flag. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: flag-scenarios-complete
    on: [flag]
    check: flag-scenarios-complete
    blocking: true
    measures: "todos os valores da flag possuem cenários documentados"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para governança de feature flags.
- **How to Fix:** Adicione o scenario para os estados que faltam (ex: `### OFF — Comportamento desativado`).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
