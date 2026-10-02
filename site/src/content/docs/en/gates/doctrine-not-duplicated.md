---
title: "Gate: doctrine-not-duplicated"
description: "Ensures that regras de produto não foram duplicadas no corpo de specs locais."
---

> **Gate Identifier:** `doctrine-not-duplicated` / `doutrina-nao-duplicada`  
> **Code:** `DCTRN` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that regras de produto não foram duplicadas no corpo de specs locais.

---

## 2. Why It Matters

A regra de produto deve ter uma única fonte de verdade. Se for copiada em várias specs, divergirá com o tempo.

---

## 3. How It Works

Detecta frases ou blocos de texto idênticos aos dos arquivos de doutrina dentro de specs comuns.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The spec referencia a doutrina via `@realizes` sem copiar o texto integral. | None. Pipeline proceeds. |
| **`FAIL`** | Texto da doutrina duplicado no corpo dThe spec local. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: doctrine-not-duplicated
    on: [spec]
    check: doctrine-not-duplicated
    blocking: true
    measures: "regras de doutrina não são copiadas no corpo das specs locais"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a arquitetura limpa.
- **How to Fix:** Remova o texto duplicado dThe spec local e aponte apenas `@realizes CODIGO`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
