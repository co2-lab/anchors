---
title: "Gate: line-coverage"
description: "Verifies whether a cobertura de linhas dThe file atinge o piso mínimo exigido (ex: >= 70%)."
---

> **Gate Identifier:** `line-coverage` / `cobertura-de-linha`  
> **Code:** `INCHN` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`, `code`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Verifies whether a cobertura de linhas dThe file atinge o piso mínimo exigido (ex: >= 70%).

---

## 2. Why It Matters

Ensures that as linhas de código foram pelo menos exercitadas pela suíte durante os testes.

---

## 3. How It Works

Lê o relatório LCOV ingerido e calcula a porcentagem de linhas cobertas daquele arquivo.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Cobertura de linha maior ou igual ao limiar configurado. | None. Pipeline proceeds. |
| **`FAIL`** | Cobertura de linha abaixo do mínimo exigido. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Nenhum relatório de cobertura foi ingerido. | Review before the next release cycle. |
| **`SKIP`** | Arquivos sem linhas executáveis. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: line-coverage
    on: [code]
    check: line-coverage
    blocking: false
    measures: "a cobertura de linhas atinge o piso mínimo"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; aumente o rigor conforme a base estabiliza.
- **How to Fix:** Escreva testes unitários cobrindo as linhas e ramos que não foram executados.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
