---
title: "Gate: doc-required"
description: "Ensures that a unidade alimenta o documento agregado obrigatório do projeto."
---

> **Gate Identifier:** `doc-required` / `documento-obrigatorio`  
> **Code:** `DCRQD` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Composto / Batch` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that a unidade alimenta o documento agregado obrigatório do projeto.

---

## 2. Why It Matters

Ensures that relatórios de auditoria agregados recebam dados de todas as unidades.

---

## 3. How It Works

Checks whether a unidade alimenta os artefatos agregados configurados.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Documento agregado alimentado. | None. Pipeline proceeds. |
| **`FAIL`** | Unidade não alimenta o documento obrigatório. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: doc-required
    blocking: false
    measures: "a unidade alimenta o documento agregado obrigatório"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo em auditorias de conformidade.
- **How to Fix:** Inclua a seção de alimentação do documento agregado nThe spec da unidade.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
