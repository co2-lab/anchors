---
title: "Gate: plan-source-declared"
description: "Ensures that um plano que cita uma fonte externa declara explicitamente quem a constrói."
---

> **Gate Identifier:** `plan-source-declared` / `fonte-de-plano-declarada`  
> **Code:** `PSDPL` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that um plano que cita uma fonte externa declara explicitamente quem a constrói.

---

## 2. Why It Matters

Evita planos baseados em premissas ou artefatos que ninguém se responsabilizou por gerar.

---

## 3. How It Works

Verifies whether referências a fontes possuem a declaração de builder ou responsável.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as fontes possuem construtor declarado. | None. Pipeline proceeds. |
| **`FAIL`** | Fonte citada sem declaração de quem a constrói. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Planos autossuficientes. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: plan-source-declared
    on: [plan]
    check: plan-source-declared
    blocking: true
    measures: "um plano que cita uma fonte declara quem a constrói"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para planos complexos.
- **How to Fix:** Declare o responsável ou passo de geração da fonte no plano.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
