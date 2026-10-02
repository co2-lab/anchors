---
title: "Gate: scenario-asserts"
description: "Ensures that o passo de desfecho (Then/Então) afirma um resultado observável concreto."
---

> **Gate Identifier:** `scenario-asserts` / `cenario-afirma`  
> **Code:** `SCASS` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o passo de desfecho (Then/Então) afirma um resultado observável concreto.

---

## 2. Why It Matters

Um passo Then que diz 'o sistema funciona' é uma tautologia que não prova nada. Este gate exige asserções concretas.

---

## 3. How It Works

Analisa as frases dos passos Then/Então procurando verbos e substantivos de validação observável.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | O passo final expressa uma asserção verificável (ex: 'deve retornar status 200', 'o botão fica desabilitado'). | None. Pipeline proceeds. |
| **`FAIL`** | O passo de resultado é vago ou tautológico. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: scenario-asserts
    on: [feature]
    check: scenario-asserts
    blocking: true
    measures: "o passo Então afirma algo observável"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter o rigor dos scenarios em Gherkin.
- **How to Fix:** Reescreva o passo Então com um resultado claro e mensurável.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
