---
title: "Gate: presentation-observable"
description: "Ensures that o elemento alterado pela apresentação possui identificador que testes conseguem apontar."
---

> **Gate Identifier:** `presentation-observable` / `apresentacao-observavel`  
> **Code:** `PRSNT` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`, `Componentes`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o elemento alterado pela apresentação possui identificador que testes conseguem apontar.

---

## 2. Why It Matters

Impede criar regras visuais que nenhum teste automatizado consegue inspecionar.

---

## 3. How It Works

Verifies whether o elemento visual alterado possui um testID ou seletor semântico declarado.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os elementos alterados possuem seletores observáveis. | None. Pipeline proceeds. |
| **`FAIL`** | Elemento sem seletor de teste identificável. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: presentation-observable
    on: [spec]
    check: presentation-observable
    blocking: true
    measures: "o que a tela muda possui elemento observável por teste"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para frontend e mobile.
- **How to Fix:** Declare o `testID` ou seletor acessível do elemento nThe spec e nThe code.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
