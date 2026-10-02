---
title: "Gate: domain-declared"
description: "Verifies whether The spec declara o que a unidade aceita e quem bloqueia entradas inválidas."
---

> **Gate Identifier:** `domain-declared` / `dominio-declarado`  
> **Code:** `DMDCD` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether The spec declara o que a unidade aceita e quem bloqueia entradas inválidas.

---

## 2. Why It Matters

Ensures that valores fora do domínio (ex: idade negativa, string nula) tenham tratamento explícito.

---

## 3. How It Works

Lê a seção de Domínio dThe spec e verifica a presença da coluna 'Quem bloqueia' ou tratamento de erro.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A tabela de domínio define entradas válidas, entradas inválidas e o componente responsável pelo bloqueio. | None. Pipeline proceeds. |
| **`FAIL`** | Tabela de domínio ausente ou sem declaração de quem bloqueia entradas inválidas. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Specs puramente visuais. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: domain-declared
    on: [spec]
    check: domain-declared
    blocking: true
    measures: "a spec declara o que aceita e quem bloqueia o inválido"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em camadas de usecase e domínio.
- **How to Fix:** Preencha a tabela de Domínio nThe spec especificando entradas aceitas e quem rejeita o inválido.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
