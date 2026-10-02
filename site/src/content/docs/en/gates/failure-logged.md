---
title: "Gate: failure-logged"
description: "Ensures that falhas tratadas emitem log adequado com contexto."
---

> **Gate Identifier:** `failure-logged` / `falha-registrada`  
> **Code:** `FLRAI` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that falhas tratadas emitem log adequado com contexto.

---

## 2. Why It Matters

Erros silenciados sem log dificultam a observabilidade e o diagnóstico em produção.

---

## 3. How It Works

Verifies whether os tratamentos de erro executam chamadas de logging estruturado com contexto do erro.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Falhas emitem log corresponding. | None. Pipeline proceeds. |
| **`FAIL`** | Tratamento de erro silencia falha sem log. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: failure-logged
    on: [code]
    check: failure-logged
    blocking: false
    measures: "falhas tratadas geram log com contexto"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo para observabilidade.
- **How to Fix:** Adicione chamada ao logger dentro do bloco de tratamento de erro.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
