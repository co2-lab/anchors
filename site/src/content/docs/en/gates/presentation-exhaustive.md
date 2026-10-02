---
title: "Gate: presentation-exhaustive"
description: "Ensures that todo valor de prop ou estado lido pela apresentação tem aparência explicitamente decidida."
---

> **Gate Identifier:** `presentation-exhaustive` / `apresentacao-exaustiva`  
> **Code:** `PRSNT` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`, `Componentes`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo valor de prop ou estado lido pela apresentação tem aparência explicitamente decidida.

---

## 2. Why It Matters

Evita telas quebradas quando um prop inesperado (ex: estado de carregamento ou lista vazia) é recebido.

---

## 3. How It Works

Cruza os estados e props da interface com as variantes visuais catalogadas nThe spec.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os estados e combinações possuem aparência documentada. | None. Pipeline proceeds. |
| **`FAIL`** | Prop ou estado da tela sem comportamento visual catalogado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Camadas que não são de apresentação. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: presentation-exhaustive
    on: [spec]
    check: presentation-exhaustive
    blocking: true
    measures: "todo valor de prop ou estado tem aparência decidida"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em frontend e mobile.
- **How to Fix:** Descreva o comportamento visual para todos os estados nThe spec (carregando, vazio, sucesso, erro).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
