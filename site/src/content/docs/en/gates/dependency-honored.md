---
title: "Gate: dependency-honored"
description: "Verifies whether os métodos prometidos na tabela de dependencies dThe spec são realmente consumidos nThe code."
---

> **Gate Identifier:** `dependency-honored` / `dependencia-honrada`  
> **Code:** `DEPHN` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether os métodos prometidos na tabela de dependencies dThe spec são realmente consumidos nThe code.

---

## 2. Why It Matters

Evita dependencies zumbis que constam na arquitetura mas não têm utilidade.

---

## 3. How It Works

Extrai a tabela de dependencies dThe spec e procura as chamadas corresponding nThe code fonte.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os métodos prometidos nThe spec são chamados nThe code. | None. Pipeline proceeds. |
| **`FAIL`** | Método prometido na tabela de dependencies nunca é invocado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código pendente. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: dependency-honored
    blocking: true
    measures: "métodos prometidos na dependência são consumidos no código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a arquitetura limpa.
- **How to Fix:** Invoque o método nThe code ou remova a dependency desnecessária dThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
