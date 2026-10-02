---
title: "Gate: scenario-identity"
description: "Ensures that cada scenario é distinguível, com código próprio e passos que não são cópias de outro."
---

> **Gate Identifier:** `scenario-identity` / `identidade-de-cenario`  
> **Code:** `SCIDS` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that cada scenario é distinguível, com código próprio e passos que não são cópias de outro.

---

## 2. Why It Matters

Evita duplicação de scenarios com títulos diferentes que fingem cobrir mais coisas do que realmente cobrem.

---

## 3. How It Works

Compara os códigos e os passos Given/When/Then dos scenarios da feature, acusando repetições literais sob títulos diferentes.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os scenarios têm códigos únicos e passos diferenciados. | None. Pipeline proceeds. |
| **`FAIL`** | Dois scenarios compartilham o mesmThe code ou possuem passos idênticos. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: scenario-identity
    on: [feature]
    check: scenario-identity
    blocking: true
    measures: "cada cenário é distinguível por código e passos"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Dê um código único a cada scenario e diferencie seus passos ou parâmetros de teste.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
