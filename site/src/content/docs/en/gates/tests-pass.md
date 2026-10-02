---
title: "Gate: tests-pass"
description: "Verifies whether a suíte de testes passou com zero falhas no relatório ingerido."
---

> **Gate Identifier:** `tests-pass` / `tests-green`  
> **Code:** `PRJTS` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `All Layers`, `test`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether a suíte de testes passou com zero falhas no relatório ingerido.

---

## 2. Why It Matters

Não adianta ter testes se eles estão falhando no pipeline.

---

## 3. How It Works

Lê The file JUnit XML gerado pelo runner de testes após a ingestão (`anchors ingest`).

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | ZerThe tests falharam (failed: 0). | None. Pipeline proceeds. |
| **`FAIL`** | Um ou mais testes falharam. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Nenhum resultado de teste foi ingerido ainda. | Review before the next release cycle. |
| **`SKIP`** | Arquivos de suporte a testes (helpers). | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: tests-pass
    on: [test]
    check: tests-pass
    blocking: true
    measures: "a suíte de testes passa com zero falhas"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking obrigatório em pre-commit, pre-push e CI.
- **How to Fix:** Corrija The code que causou a quebra dThe test e rode a suíte novamente.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
