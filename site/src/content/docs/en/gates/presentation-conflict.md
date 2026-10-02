---
title: "Gate: presentation-conflict"
description: "Ensures that um prop e uma condição não levam a duas aparências conflitantes."
---

> **Gate Identifier:** `presentation-conflict` / `apresentacao-sem-conflito`  
> **Code:** `PRSNT` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`, `Componentes`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that um prop e uma condição não levam a duas aparências conflitantes.

---

## 2. Why It Matters

Impede regras ambíguas que dizem ao mesmo tempo 'botão fica vermelho' e 'botão fica desabilitado cinza'.

---

## 3. How It Works

Analisa as matrizes de decisão visual procurando sobreposição de condições contraditórias.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhum conflito de aparência encontrado. | None. Pipeline proceeds. |
| **`FAIL`** | Mesma condição conduz a duas aparências visuais incompatíveis. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: presentation-conflict
    on: [spec]
    check: presentation-conflict
    blocking: true
    measures: "um prop e uma condição não levam a duas aparências"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para design systems e componentes visuais.
- **How to Fix:** Defina a ordem de precedência clara entre os estados visuais nThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
