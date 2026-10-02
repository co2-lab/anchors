---
title: "Gate: rule-types"
description: "Ensures that os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto."
---

> **Gate Identifier:** `rule-types` / `tipos-de-regra`  
> **Code:** `RLTYR` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec`, `feature` | **Layers:** `All Layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto.

---

## 2. Why It Matters

Mantém o vocabulário de letras de regras padronizado em todo o repositório.

---

## 3. How It Works

Compara as letras de regras usadas nas specs contra o vocabulário `rule_types` do anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as letras de códigos pertencem ao vocabulário oficial. | None. Pipeline proceeds. |
| **`FAIL`** | Letra de código não cadastrada no anchors.yaml. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: rule-types
    blocking: true
    measures: "as letras dos códigos estão no vocabulário declarado"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a consistência da taxonomia.
- **How to Fix:** Use as letras padrão (B, V, E, S, DS, VR) ou declare a nova letra em `rule_types`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
