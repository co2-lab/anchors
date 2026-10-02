---
title: "Gate: trigger-declared"
description: "Ensures that triggers de conformidade e eventos citados existem no vocabulário declarado do projeto."
---

> **Gate Identifier:** `trigger-declared` / `trigger-declarado`  
> **Code:** `TRDCT` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that triggers de conformidade e eventos citados existem no vocabulário declarado do projeto.

---

## 2. Why It Matters

Impede o uso de nomes arbitrários de eventos e gatilhos que ninguém consegue correlacionar.

---

## 3. How It Works

Confere eventos e gatilhos citados contra a lista de triggers válidos no anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os triggers pertencem ao vocabulário configurado. | None. Pipeline proceeds. |
| **`FAIL`** | Trigger ou obrigação não cadastrada encontrada nThe spec. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: trigger-declared
    on: [spec]
    check: trigger-declared
    blocking: true
    measures: "triggers e obrigações citados existem no vocabulário"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos orientados a eventos ou compliance regulatório.
- **How to Fix:** Declare o trigger no vocabulário do anchors.yaml ou corrija o nome nThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
