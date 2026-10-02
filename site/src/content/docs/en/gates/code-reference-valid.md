---
title: "Gate: code-reference-valid"
description: "Ensures that referências cruzadas a códigos de outras regras apontam para regras que realmente existem."
---

> **Gate Identifier:** `code-reference-valid` / `referencia-de-codigo-valida`  
> **Code:** `CRVCD` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec`, `feature`, `code`, `test` | **Layers:** `All Layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that referências cruzadas a códigos de outras regras apontam para regras que realmente existem.

---

## 2. Why It Matters

Quando umThe spec cita '@realizes REQ-01' ou '@ref: AUTH-B01', essa regra precisa existir. Caso contrário, é um link quebrado.

---

## 3. How It Works

Varre o repositório indexando todos os códigos válidos e Verifies whether todas as menções cruzadas resolvem para um código existente.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as referências a códigos apontam para regras catalogadas no projeto. | None. Pipeline proceeds. |
| **`FAIL`** | Foi encontrada referência a um código que não existe em nenhumThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: code-reference-valid
    blocking: true
    measures: "as referências a códigos apontam para regras que existem"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking obrigatório para manter o grafo íntegro.
- **How to Fix:** Corrija The code digitado ou crie a regra corresponding nThe spec dona.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
