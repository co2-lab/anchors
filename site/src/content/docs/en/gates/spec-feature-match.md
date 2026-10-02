---
title: "Gate: spec-feature-match"
description: "Ensures that toda regra catalogada nThe spec possui ao menos um scenario corresponding na feature."
---

> **Gate Identifier:** `spec-feature-match` / `spec-feature-match`  
> **Code:** `SFMSP` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that toda regra catalogada nThe spec possui ao menos um scenario corresponding na feature.

---

## 2. Why It Matters

Evita que regras sejam esquecidas no papel sem que ninguém especifique como elas devem ser testadas.

---

## 3. How It Works

Extrai todos os códigos de regra definidos nThe spec e Checks whether cada um é citado como tag (@CODE-B01) na feature ligada.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as regras dThe spec possuem tags de scenario corresponding na feature. | None. Pipeline proceeds. |
| **`FAIL`** | Existe alguma regra definida nThe spec sem scenario na feature. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | A feature ligada ainda não foi criada. | Review before the next release cycle. |
| **`SKIP`** | Regras marcadas com @no-scenario: motivo ou camadas que dispensam feature. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: spec-feature-match
    on: [spec]
    check: spec-feature-match
    blocking: true
    measures: "cada regra da spec tem cenário na feature"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos onde BDD/Gherkin é parte da governança.
- **How to Fix:** Adicione um scenario na feature com a tag da regra pendente, ou declare @no-scenario nThe spec com justificativa.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
