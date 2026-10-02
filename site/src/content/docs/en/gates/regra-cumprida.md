---
title: "Gate: regra-cumprida"
description: "Pergunta ao modelo se o trecho marcado nThe code realmente realiza o que a regra descreve."
---

> **Gate Identifier:** `regra-cumprida` / `regra-cumprida`  
> **Code:** `RLUEX` | **Category:** [AI Judgment & Synthetic](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`  
> **Execution Model:** `Julgamento por IA (ask)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Pergunta ao modelo se o trecho marcado nThe code realmente realiza o que a regra descreve.

---

## 2. Why It Matters

Regexes e linters não sabem se uma função faz o que promete semânticamente. A IA avalia a correspondência real de intenção.

---

## 3. How It Works

Avalia o prompt semântico comparando o texto da regra nThe spec com The code entre os marcadores nThe file fonte.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Veredito `pass` gravado com `anchors judge`. | None. Pipeline proceeds. |
| **`FAIL`** | Veredito `fail` acusando que The code diverge da regra. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Nenhum veredito emitido ainda. | Review before the next release cycle. |
| **`SKIP`** | Regra marcada como pendente de implementação (@TBD: code) recebe veredito `dispensado`. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: regra-cumprida
    on: [spec]
    ask: "o trecho marcado REALIZA o que a regra descreve?"
    blocking: false
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo; avalie e registre com `anchors judge` em revisões de PR.
- **How to Fix:** Ajuste The code para cumprir a regra dThe spec e grave o novo laudo com `anchors judge`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
