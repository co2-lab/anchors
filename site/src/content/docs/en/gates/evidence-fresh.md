---
title: "Gate: evidence-fresh"
description: "Ensures that o placar dThe test continua fresco e foi medido contra The code atual."
---

> **Gate Identifier:** `evidence-fresh` / `evidencia-fresca`  
> **Code:** `EVFRV` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `test`, `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o placar dThe test continua fresco e foi medido contra The code atual.

---

## 2. Why It Matters

Prevents o pipeline aprove commits usando relatórios de testes antigos medidos antes de alterações recentes.

---

## 3. How It Works

Compara o hash ou data de modificação dThe test ingerido com a versão atual dos arquivos no disco.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | O relatório de execução foi gerado na mesma revisão dThe code atual. | None. Pipeline proceeds. |
| **`FAIL`** | The code ou The spec mudaram após a última execução dThe test. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Relatório ausente. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: evidence-fresh
    blocking: true
    measures: "o resultado do teste vale contra o código de hoje"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para impedir testes stale.
- **How to Fix:** Execute a suíte de testes e ingira o novo relatório com `anchors ingest`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
