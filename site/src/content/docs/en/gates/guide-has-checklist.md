---
title: "Gate: guide-has-checklist"
description: "Ensures that guias de governança destilam suas regras em uma seção de checklist com itens CK1, CK2..."
---

> **Gate Identifier:** `guide-has-checklist` / `guia-tem-checklist`  
> **Code:** `INCHN` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `guide` | **Layers:** `doutrina`, `guide`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that guias de governança destilam suas regras em uma seção de checklist com itens CK1, CK2...

---

## 2. Why It Matters

Guias que são apenas prosa não fornecem critérios objetivos para que agentes de IA façam checagens.

---

## 3. How It Works

Verifica a presença do título '## Pontos de conformidade' (ou equivalente em inglês/espanhol) e itens com formato `CK1`, `CK2`.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | O guia contém a seção de pontos de conformidade e itens CK. | None. Pipeline proceeds. |
| **`FAIL`** | Guia sem seção de conformidade ou sem itens CK. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Arquivos que não são do kind `guide`. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: guide-has-checklist
    on: [guide]
    check: guide-has-checklist
    blocking: true
    measures: "um guia de governança destila regras em pontos CK"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para documentos que regem o projeto.
- **How to Fix:** Adicione a seção `## Compliance points` com itens `- CK1: ...` ao final do guia.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
