---
title: "Gate: scenario-coverage"
description: "Ensures that cada scenario declarado nThe spec tem um teste corresponding que rodou e passou."
---

> **Gate Identifier:** `scenario-coverage` / `cobertura-de-cenario`  
> **Code:** `SFMSP` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that cada scenario declarado nThe spec tem um teste corresponding que rodou e passou.

---

## 2. Why It Matters

Fecha a rastreabilidade semântica: não basta The test existir nThe file, ele precisa ter rodado e sido aprovado.

---

## 3. How It Works

Cruza os códigos de scenario dThe spec com a lista de códigos comprovados (ProvenCodes) no relatório ingerido.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os scenarios declarados possuem teste verde no relatório. | None. Pipeline proceeds. |
| **`FAIL`** | scenario declarado nThe spec não possui teste ou The test falhou. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Relatório de testes ainda não ingerido. | Review before the next release cycle. |
| **`SKIP`** | Camadas com dispensa de teste declarada. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: scenario-coverage
    blocking: true
    measures: "cada cenário da spec tem teste verde executado"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no CI.
- **How to Fix:** Execute a suíte com `anchors test` e ingira o resultado com `anchors ingest`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
