---
title: "Gate: single-test-per-unit"
description: "Ensures that uma unidade tem um únicThe file de teste por camada de teste (ou declara divisão com @split-test)."
---

> **Gate Identifier:** `single-test-per-unit` / `unico-teste-por-unidade`  
> **Code:** `SNGTU` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that uma unidade tem um únicThe file de teste por camada de teste (ou declara divisão com @split-test).

---

## 2. Why It Matters

Evita testes duplicados ou espalhados em vários arquivos que testam a mesma unidade sem critério.

---

## 3. How It Works

Conta quantos arquivos de teste apontam para a mesma unidade e valida a presença de @split-test se houver mais de um.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Existe exatamente um arquivo de teste por camada para a unidade, ou a divisão foi justificada. | None. Pipeline proceeds. |
| **`FAIL`** | Múltiplos arquivos de teste encontrados para a mesma unidade sem declaração @split-test. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: single-test-per-unit
    on: [code]
    check: single-test-per-unit
    blocking: true
    measures: "uma unidade tem um arquivo de teste por camada de teste"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a organização dos testes.
- **How to Fix:** Una os arquivos de teste ou declare @split-test no cabeçalho justificando a divisão.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
