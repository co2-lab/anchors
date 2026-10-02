---
title: "Gate: presentation-copy-single-source"
description: "Ensures that o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido."
---

> **Gate Identifier:** `presentation-copy-single-source` / `texto-de-apresentacao-fonte-unica`  
> **Code:** `PRSNT` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido.

---

## 2. Why It Matters

Facilita internacionalização (i18n) e Prevents o mesmo texto mude em uma tela e continue antigo noutra.

---

## 3. How It Works

Verifies whether os textos de interface apontam para identificadores de strings ou chaves de i18n.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Textos referenciados por códigos de mensagens. | None. Pipeline proceeds. |
| **`FAIL`** | Textos literais repetidos na regra sem chave de catálogo. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: presentation-copy-single-source
    on: [spec]
    check: presentation-copy-single-source
    blocking: true
    measures: "o texto que a tela mostra vem de um código de mensagem"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Recomendado para apps multilíngues.
- **How to Fix:** Extraia o texto para o catálogo de mensagens/i18n e use a chave corresponding nThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
