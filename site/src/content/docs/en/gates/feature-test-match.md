---
title: "Gate: feature-test-match"
description: "Ensures that cada scenario da feature está implementado nThe test por código e descrição."
---

> **Gate Identifier:** `feature-test-match` / `feature-test-match`  
> **Code:** `FTMFT` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`, `Test`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that cada scenario da feature está implementado nThe test por código e descrição.

---

## 2. Why It Matters

Prevents o desenvolvedor ou a IA escreva um teste que cita The code do scenario, mas testa outra coisa completamente diferente.

---

## 3. How It Works

Lê os scenarios da feature e procura nos testes ligados tanto a menção aThe code quanto trechos significativos da descrição do scenario.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os scenarios da feature estão presentes e descritos nos testes corresponding. | None. Pipeline proceeds. |
| **`FAIL`** | scenarios da feature não encontrados nThe test ou descrição divergente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | The file de teste ainda não foi criado ou ingerido. | Review before the next release cycle. |
| **`SKIP`** | scenarios com tag @no-test. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: feature-test-match
    on: [feature]
    check: feature-test-match
    blocking: true
    measures: "cada cenário da feature está no teste ligado por código e descrição"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking assim que a suíte de testes estiver rodando.
- **How to Fix:** Atualize a função de teste ou o describe/it para incluir The code e o título do scenario da feature.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
