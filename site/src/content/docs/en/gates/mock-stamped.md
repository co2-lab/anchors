---
title: "Gate: mock-stamped"
description: "Ensures that todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa."
---

> **Gate Identifier:** `mock-stamped` / `mock-carimbado`  
> **Code:** `MCSTM` | **Category:** [Test Doubles & Mocks](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `governed layers`, `test`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa.

---

## 2. Why It Matters

Quando The code real de um serviço muda, o mock precisa ser atualizado; se o hash divergir, o mock mentiu.

---

## 3. How It Works

Verifica a presença do carimbo `@contract <hash>` no mock e compara o hash com a assinatura do método real substituído.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | O carimbo bate com o hash da assinatura do método original. | None. Pipeline proceeds. |
| **`FAIL`** | Mock sem carimbo ou hash desatualizado (o método real mudou). | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Mocks declarados em camadas de suporte com dispensa. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: mock-stamped
    on: [test]
    check: mock-stamped
    blocking: true
    measures: "o dublê carrega a marca do snippet que substitui"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para impedir mocks que mentem.
- **How to Fix:** Atualize o mock para refletir o novo contrato e regenere o carimbo com `anchors check --fix`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
