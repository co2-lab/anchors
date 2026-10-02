---
title: "Gate: testid-queried-exists"
description: "Ensures that todo testID buscado por um roteiro de fluxo E2E existe de verdade nThe code."
---

> **Gate Identifier:** `testid-queried-exists` / `testid-buscado-existe`  
> **Code:** `TQETS` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `test` | **Layers:** `UI`, `E2E`, `test`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that todo testID buscado por um roteiro de fluxo E2E existe de verdade nThe code.

---

## 2. Why It Matters

Evita testes E2E falhando por buscar elementos que foram renomeados ou excluídos da tela.

---

## 3. How It Works

Examina as queries `getByTestId` nos testes E2E e verifica a presença da string nThe code dos componentes.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os testIDs buscados existem nThe code fonte. | None. Pipeline proceeds. |
| **`FAIL`** | Teste E2E buscando testID que não existe em nenhum componente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: testid-queried-exists
    on: [test]
    check: testid-queried-exists
    blocking: true
    measures: "todo testID buscado por fluxo E2E existe no código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no CI.
- **How to Fix:** Adicione o `testID` no componente de destino ou corrija a query nThe test.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
