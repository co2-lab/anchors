---
title: "Gate: code-cataloged"
description: "Ensures that todo símbolo público (função, tipo, export) exportado pelThe code está catalogado nThe spec."
---

> **Gate Identifier:** `code-cataloged` / `codigo-catalogado`  
> **Code:** `CDCTC` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo símbolo público (função, tipo, export) exportado pelThe code está catalogado nThe spec.

---

## 2. Why It Matters

Impede o surgimento de código fantasma ou funções públicas criadas sem specification que ninguém sabe o que fazem.

---

## 3. How It Works

Analisa The code procurando funções/tipos exportados e confronta com The spec que governa The file.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os símbolos exportados têm regra corresponding nThe spec ou dispensa com @no-rule. | None. Pipeline proceeds. |
| **`FAIL`** | Função pública exportada sem nenhuma menção nThe spec dona. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Spec ainda não criada. | Review before the next release cycle. |
| **`SKIP`** | recognized layers (infra, dao, types) ou arquivos com regime declarativo. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: code-cataloged
    blocking: true
    measures: "todo símbolo exportado tem regra na spec ou dispensa"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em bibliotecas e serviços de negócio.
- **How to Fix:** Documente a função nThe spec corresponding ou marque com @no-rule nThe code explicando o motivo.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
