---
title: "Gate: code-language"
description: "Ensures that The code não mistura idiomas de identificadores e comentários."
---

> **Gate Identifier:** `code-language` / `idioma-de-codigo`  
> **Code:** `CDLNG` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that The code não mistura idiomas de identificadores e comentários.

---

## 2. Why It Matters

Metade das variáveis em inglês e metade em português cria confusão e dificulta refatorações.

---

## 3. How It Works

Inspeciona identificadores e comentários verificando coerência idiomática com o padrão configurado.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Identificadores e comentários coerentes com o idioma padrão. | None. Pipeline proceeds. |
| **`FAIL`** | Mistura desordenada de idiomas no mesmThe file. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: code-language
    on: [code]
    blocking: true
    measures: "o código não mistura idiomas"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a higiene dThe code.
- **How to Fix:** Padronize os nomes de variáveis e comentários no idioma adotado pelo projeto.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
