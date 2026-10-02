---
title: "Gate: contract-status-declared"
description: "Ensures that o contrato lista os códigos de status que The code realmente retorna, e apenas esses."
---

> **Gate Identifier:** `contract-status-declared` / `status-de-contrato-declarado`  
> **Code:** `CSDCN` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`, `API`, `comando`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o contrato lista os códigos de status que The code realmente retorna, e apenas esses.

---

## 2. Why It Matters

Evita APIs mentirosas que dizem retornar apenas 200 e 400, mas nThe code lançam 403, 404 e 500 sem documentar.

---

## 3. How It Works

Extrai retornos HTTP/status dThe code e confronta com a tabela de status declarada nThe spec.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Os status retornados nThe code batem exatamente com os declarados nThe spec. | None. Pipeline proceeds. |
| **`FAIL`** | Código retorna status não documentado nThe spec, ou spec documenta status nunca retornado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código ainda não implementado. | Review before the next release cycle. |
| **`SKIP`** | Camadas sem contratos de saída. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: contract-status-declared
    blocking: true
    measures: "o contrato de saída lista os status reais retornados"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos de API REST / gRPC.
- **How to Fix:** Declare o status na tabela de contrato dThe spec ou trate o erro nThe code.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
