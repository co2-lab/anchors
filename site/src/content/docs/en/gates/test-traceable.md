---
title: "Gate: test-traceable"
description: "Ensures that todThe test ligado a uma feature declara no título The code do scenario que prova."
---

> **Gate Identifier:** `test-traceable` / `teste-rastreavel`  
> **Code:** `TSTRT` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `governed layers`, `test`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todThe test ligado a uma feature declara no título The code do scenario que prova.

---

## 2. Why It Matters

Sem The code no título dThe test, o runner não consegue ligar o resultado da execução ao requisito dThe spec.

---

## 3. How It Works

Examina as funções e blocos it/test dThe file de teste procurandThe codes de scenario válidos.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os testes declaram The code do scenario que comprovam. | None. Pipeline proceeds. |
| **`FAIL`** | Teste sem menção a código de scenario ou citandThe code que não existe. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Arquivos de teste em camadas declarativas ou testes auxiliares. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: test-traceable
    blocking: true
    measures: "o teste cita o código do cenário que prova"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para fechar a fiação de rastreabilidade.
- **How to Fix:** Inclua The code do scenario no nome da função de teste (ex: `TestAuth_AUTH_B01_Bloqueio`).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
