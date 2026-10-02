---
title: "Gate: failure-declared"
description: "Ensures that possíveis falhas da unidade estão declaradas nThe specification."
---

> **Gate Identifier:** `failure-declared` / `falha-declarada`  
> **Code:** `FLRAI` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that possíveis falhas da unidade estão declaradas nThe specification.

---

## 2. Why It Matters

Evita a ilusão do 'caminho feliz': todo software falha, e as falhas esperadas precisam ser documentadas.

---

## 3. How It Works

Lê a seção de Erros/Falhas (`## Errors`) nThe spec e verifica a catalogação de códigos com prefixo E.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The spec possui catálogo de falhas declaradas. | None. Pipeline proceeds. |
| **`FAIL`** | Spec sem declaração de erros ou falhas possíveis. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Camadas declarativas ou funções puras que não falham. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: failure-declared
    on: [spec]
    check: failure-declared
    blocking: true
    measures: "possíveis falhas da unidade são declaradas na spec"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para usecases e serviços.
- **How to Fix:** Adicione a seção `## Errors` com ao menos uma regra de erro catalogada (ex: `CODE-E01`).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
