---
title: "Gate: scenario-letter-declared"
description: "Verifies whether a letra usada nThe code do scenario existe no vocabulário declarado do projeto."
---

> **Gate Identifier:** `scenario-letter-declared` / `letra-cenario-declarada`  
> **Code:** `SCLTR` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether a letra usada nThe code do scenario existe no vocabulário declarado do projeto.

---

## 2. Why It Matters

Prevents alguém invente letras aleatórias em códigos (ex: AUTH-X01) sem padronização.

---

## 3. How It Works

Valida a letra dThe code contra a lista de letras permitidas configurada em rule_types no anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A letra pertence ao vocabulário oficial (ex: B, V, E, S, DS, VR). | None. Pipeline proceeds. |
| **`FAIL`** | Código usando letra não declarada no vocabulário. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: scenario-letter-declared
    on: [feature]
    check: scenario-letter-declared
    blocking: true
    measures: "as letras dos códigos estão no vocabulário declarado"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Use uma das letras canônicas do projeto ou declare a nova letra em rule_types no anchors.yaml.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
