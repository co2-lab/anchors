---
title: "Gate: flag-scenario-exists"
description: "Ensures that uma regra que cita @gated-by aponta para uma flag e scenario que realmente existem."
---

> **Gate Identifier:** `flag-scenario-exists` / `cenario-de-flag-existe`  
> **Code:** `FLSCF` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that uma regra que cita @gated-by aponta para uma flag e scenario que realmente existem.

---

## 2. Why It Matters

Evita links quebrados para feature flags que foram removidas ou digitadas errado.

---

## 3. How It Works

Verifies whether o identificador em `@gated-by FLG.ESTADO` resolve para um arquivo de flag existente.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A flag e o estado citado existem. | None. Pipeline proceeds. |
| **`FAIL`** | Referência a flag ou estado inexistente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Specs sem anotação @gated-by. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: flag-scenario-exists
    on: [spec]
    check: flag-scenario-exists
    blocking: true
    measures: "a anotação @gated-by aponta para flag e cenário existentes"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Crie The file da flag ou corrija a referência em `@gated-by`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
