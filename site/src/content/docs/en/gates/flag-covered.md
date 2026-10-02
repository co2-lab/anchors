---
title: "Gate: flag-covered"
description: "Ensures that scenarios governados por flags possuem testes em todos os seus ramos."
---

> **Gate Identifier:** `flag-covered` / `flag-coberta`  
> **Code:** `FLSCF` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `flag` | **Layers:** `flags`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that scenarios governados por flags possuem testes em todos os seus ramos.

---

## 2. Why It Matters

Testar apenas o caminho ON e nunca o OFF leva a quebras graves quando a flag for desligada.

---

## 3. How It Works

Verifies whether há testes automatizados cobrindo as regras associadas a cada estado da flag.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os ramos da flag possuem testes verdes corresponding. | None. Pipeline proceeds. |
| **`FAIL`** | Um dos ramos da flag não possui teste automatizado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Relatório de testes pendente. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: flag-covered
    on: [flag]
    check: flag-covered
    blocking: false
    measures: "todos os ramos da flag possuem testes"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; promova a blocking ao estabilizar a funcionalidade.
- **How to Fix:** Escreva testes cobrindo tanto o comportamento da flag ativada quanto desativada.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
