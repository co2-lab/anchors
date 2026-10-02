---
title: "Gate: revision-orphans"
description: "Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente."
---

> **Gate Identifier:** `revision-orphans` / `orfas-de-revisao`  
> **Code:** `RVORP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente.

---

## 2. Why It Matters

Altera o significado de uma regra sem avisar quebra os testes e clientes que dependiam da semântica anterior.

---

## 3. How It Works

Detecta mudanças no texto de regras existentes sem nova numeração de versão ou nota de revisão.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Revisões de regras são explícitas e numeradas. | None. Pipeline proceeds. |
| **`FAIL`** | Regra alterada silenciosamente sem declaração de revisão. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: revision-orphans
    on: [spec]
    check: revision-orphans
    blocking: true
    measures: "regras cujo significado mudou são declaradas"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking parThe specifications maduras.
- **How to Fix:** Incremente o número da regra ou adicione uma nota de revisão nThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
