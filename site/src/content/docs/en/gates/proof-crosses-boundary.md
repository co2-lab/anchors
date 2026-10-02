---
title: "Gate: proof-crosses-boundary"
description: "Ensures that quando uma regra afirma relação entre módulos, a prova de teste alcança a outra ponta."
---

> **Gate Identifier:** `proof-crosses-boundary` / `prova-cruza-fronteira`  
> **Code:** `PCBPR` | **Category:** [Architectural Boundaries](/docs/gates/)  
> **Evaluates:** `spec`, `test` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that quando uma regra afirma relação entre módulos, a prova de teste alcança a outra ponta.

---

## 2. Why It Matters

Evita testes que mockam a outra ponta sem verificar se o contrato do vizinho ainda aceita aquela chamada.

---

## 3. How It Works

Verifies whether testes com anotação de fronteira exercitam a interação ou usam dublês carimbados.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A interação com a outra ponta é comprovada. | None. Pipeline proceeds. |
| **`FAIL`** | A regra afirma relação com outro módulo mas The test não alcança a fronteira. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Testes puramente unitários isolados. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: proof-crosses-boundary
    blocking: false
    measures: "teste que cruza fronteira a declara e comprova"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo em suítes de integração.
- **How to Fix:** Adicione um teste de integração ou utilize um mock com contrato verificado.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
