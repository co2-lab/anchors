---
title: "Gate: identity-consistent"
description: "Ensures that a identidade dThe spec bate com o testID exposto e a imagem de baseline visual."
---

> **Gate Identifier:** `identity-consistent` / `identidade-consistente`  
> **Code:** `IDCND` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `UI`, `Telas`, `Componentes`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a identidade dThe spec bate com o testID exposto e a imagem de baseline visual.

---

## 2. Why It Matters

Impede divergência de nomenclatura entre o design system, The code e os testes visuais.

---

## 3. How It Works

Cruza o identificador da unidade com o atributo testID exposto nThe code e o nome dThe file PNG de baseline.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Identidade idêntica nThe spec, nThe code e no baseline visual. | None. Pipeline proceeds. |
| **`FAIL`** | Divergência entre o testID exposto e a identidade dThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: identity-consistent
    blocking: true
    measures: "identidade da spec bate com testID e baseline visual"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em UI.
- **How to Fix:** Alinhe o `testID` no componente com a identidade declarada nThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
