---
title: "Gate: mock-typed"
description: "Ensures that todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui."
---

> **Gate Identifier:** `mock-typed` / `mock-tipado`  
> **Code:** `MCTYM` | **Category:** [Test Doubles & Mocks](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `governed layers`, `test`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui.

---

## 2. Why It Matters

Impede a criação de mocks ad-hoc com métodos inventados que a classe real não possui.

---

## 3. How It Works

Verifies whether a estrutura do mock implementa a interface ou estende a classe original.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | O mock é formalmente tipado conforme o módulo real. | None. Pipeline proceeds. |
| **`FAIL`** | Dublê sem tipagem estrita ou divergindo da interface original. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Linguagens sem tipagem estática. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: mock-typed
    on: [test]
    check: mock-typed
    blocking: true
    measures: "todo dublê de teste deriva do módulo que substitui"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em Go, TypeScript, Java e Rust.
- **How to Fix:** Faça o mock implementar a interface oficial do serviço substituído.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
