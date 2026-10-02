---
title: "Gate: open-questions-resolved"
description: "Ensures that specifications com perguntas ou decisões em aberto não sejam liberadas para implementação."
---

> **Gate Identifier:** `open-questions-resolved` / `decisoes-em-aberto-resolvidas`  
> **Code:** `OPQSP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `spec`, `plan` | **Layers:** `governed layers`, `Planos`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that specifications com perguntas ou decisões em aberto não sejam liberadas para implementação.

---

## 2. Why It Matters

UmThe spec com perguntas em aberto ainda não foi decidida. Codificar antes de decidir gera retrabalho garantido.

---

## 3. How It Works

Verifies whether a tabela de Decisões em Aberto (`## Open Decisions`) nThe spec contém itens sem resolução.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as perguntas em aberto foram respondidas ou a tabela está vazia. | None. Pipeline proceeds. |
| **`FAIL`** | Existem perguntas em aberto não resolvidas nThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: open-questions-resolved
    on: [spec]
    check: open-questions-resolved
    blocking: true
    measures: "as decisões em aberto foram decididas"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking antes de iniciar a implementação dThe code.
- **How to Fix:** Responda à pergunta nThe spec e mova a decisão para a seção de Decisões Tomadas.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
