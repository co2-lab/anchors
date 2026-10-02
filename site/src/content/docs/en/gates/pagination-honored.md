---
title: "Gate: pagination-honored"
description: "Ensures that funções ou telas que prometem conjuntos paginados não retornam apenas a primeira página em silêncio."
---

> **Gate Identifier:** `pagination-honored` / `paginacao-honrada`  
> **Code:** `PGNHN` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`, `API`, `UI`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that funções ou telas que prometem conjuntos paginados não retornam apenas a primeira página em silêncio.

---

## 2. Why It Matters

Evita o bug comum de listar dados fingindo que paginou, mas descartando os resultados além da primeira página.

---

## 3. How It Works

Verifies whether contratos que prometem paginação possuem parâmetros de cursor/offset e controle de próxima página.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Mecanismo de paginação completo declarado e implementado. | None. Pipeline proceeds. |
| **`FAIL`** | Promete paginação mas não implementa controle de página. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Funções que não operam sobre coleções paginadas. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: pagination-honored
    blocking: false
    measures: "o que promete paginação, pagina de verdade"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Recomendado para APIs e telas de listagem.
- **How to Fix:** Implemente o tratamento de limite, offset/cursor e indicador de próxima página.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
