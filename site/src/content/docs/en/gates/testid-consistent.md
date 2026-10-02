---
title: "Gate: testid-consistent"
description: "Garante o contrato de testID: The code expõe, The spec declara e The test consome exatamente o mesmo identificador."
---

> **Gate Identifier:** `testid-consistent` / `testid-coerente`  
> **Code:** `TICTS` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec`, `code`, `test` | **Layers:** `UI`, `Telas`, `Componentes`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Garante o contrato de testID: The code expõe, The spec declara e The test consome exatamente o mesmo identificador.

---

## 2. Why It Matters

Elimina testes quebrados por erro de digitação no nome do testID.

---

## 3. How It Works

Valida as quatro pontas do testID no grafo de dependencies.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Contrato de testID íntegro entre código, spec e teste. | None. Pipeline proceeds. |
| **`FAIL`** | TestID consumido nThe test não é exposto pelThe code ou não consta nThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: testid-consistent
    blocking: true
    measures: "contrato de testID coerente entre spec, código e testes"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para testes E2E e de integração visual.
- **How to Fix:** Corrija o nome do testID nThe test para bater com o exposto no componente.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
