---
title: "Gate: rule-uses-declared"
description: "Verifies whether cada regra declara explicitamente quais campos, validações e dependencies utiliza."
---

> **Gate Identifier:** `rule-uses-declared` / `regra-declara-uso`  
> **Code:** `RLUSG` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Verifies whether cada regra declara explicitamente quais campos, validações e dependencies utiliza.

---

## 2. Why It Matters

Permite calcular o impacto de mudanças de campos: se um campo mudar, o Anchors sabe quais regras são afetadas.

---

## 3. How It Works

Examina as seções de validação e uso de regras nThe spec procurando a declaração de campos consumidos.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | As regras declaram as variáveis e dependencies que utilizam. | None. Pipeline proceeds. |
| **`FAIL`** | Regra complexa sem declaração de dependencies de campos. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Specs puramente declarativas. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: rule-uses-declared
    on: [spec]
    check: rule-uses-declared
    blocking: false
    measures: "toda regra diz o que usa"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; promova a blocking ao refinar as specs.
- **How to Fix:** Adicione a seção de campos usados na regra dentro dThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
