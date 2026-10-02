---
title: "Gate: spec-realizes-doctrine"
description: "Ensures that The spec declara e comprova como ela realiza a doutrina."
---

> **Gate Identifier:** `spec-realizes-doctrine` / `spec-realiza-doutrina`  
> **Code:** `DCTRN` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that The spec declara e comprova como ela realiza a doutrina.

---

## 2. Why It Matters

Não basta dizer `@realizes`: The spec precisa descrever a regra local corresponding.

---

## 3. How It Works

Verifies whether a linha que carrega `@realizes` possui contexto funcional associado.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A realização da doutrina é contextualmente válida. | None. Pipeline proceeds. |
| **`FAIL`** | Anotação `@realizes` solta sem regra corresponding. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: spec-realizes-doctrine
    on: [spec]
    check: spec-realizes-doctrine
    blocking: true
    measures: "a spec comprova como realiza a doutrina"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos com doutrina ativa.
- **How to Fix:** Associe a anotação `@realizes` ao cabeçalho ou linha da regra que a implementa.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
