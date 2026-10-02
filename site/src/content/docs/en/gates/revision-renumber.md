---
title: "Gate: revision-renumber"
description: "Renumera revisões em branch para evitar conflito quando a branch base já utilizou o mesmo número."
---

> **Gate Identifier:** `revision-renumber` / `renumerar-revisoes`  
> **Code:** `RVRNR` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `spec`, `plan` | **Layers:** `governed layers`, `Planos`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Renumera revisões em branch para evitar conflito quando a branch base já utilizou o mesmo número.

---

## 2. Why It Matters

Evita colisões de números de revisão em merges concorrentes.

---

## 3. How It Works

Compara os números de revisão da branch com os commits recentes da branch base.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Números de revisão não colidem com a branch principal. | None. Pipeline proceeds. |
| **`FAIL`** | Colisão de numeração de revisão detectada. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: revision-renumber
    blocking: true
    measures: "revisões em branch movem para número livre se houver colisão"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em CI para validação de PRs.
- **How to Fix:** Renomeie o número da revisão na branch para o próximo número livre.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
