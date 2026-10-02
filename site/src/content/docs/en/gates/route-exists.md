---
title: "Gate: route-exists"
description: "Verifies whether a rota declarada nThe specification existe no registro de rotas da aplicação."
---

> **Gate Identifier:** `route-exists` / `rota-existe`  
> **Code:** `RTEXR` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`, `API`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether a rota declarada nThe specification existe no registro de rotas da aplicação.

---

## 2. Why It Matters

Prevents umThe spec prometa a rota `/checkout/pix` e a aplicação registre `/pagamento/pix`.

---

## 3. How It Works

Cruza a string da rota nThe spec com The file de rotas central da aplicação (declarado em router_file).

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A rota existe nThe file de rotas da aplicação. | None. Pipeline proceeds. |
| **`FAIL`** | A rota não foi encontrada no registro de rotas dThe code. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Registro de rotas pendente. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: route-exists
    on: [spec]
    check: route-exists
    blocking: true
    measures: "a rota declarada na spec existe no registro da aplicação"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em frontend e backend.
- **How to Fix:** Registre a rota no roteador da aplicação ou corrija The spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
