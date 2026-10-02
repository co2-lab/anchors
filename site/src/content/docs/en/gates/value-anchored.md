---
title: "Gate: value-anchored"
description: "Ensures that constantes ou chaves replicadas em vários arquivos carregam exatamente o mesmo valor."
---

> **Gate Identifier:** `value-anchored` / `valor-ancorado`  
> **Code:** `VLANV` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `code`, `spec` | **Layers:** `All Layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that constantes ou chaves replicadas em vários arquivos carregam exatamente o mesmo valor.

---

## 2. Why It Matters

Impede bugs onde uma chave de evento ou header HTTP é copiado com letras minúsculas num arquivo e maiúsculas noutro.

---

## 3. How It Works

Busca chaves ancoradas marcadas e valida a igualdade exata de valor em todas as cópias do projeto.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as cópias do valor possuem conteúdo idêntico. | None. Pipeline proceeds. |
| **`FAIL`** | Foi encontrada divergência de valor entre arquivos que deveriam compartilhar a mesma chave. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: value-anchored
    blocking: true
    measures: "chaves replicadas carregam o mesmo valor em todo o projeto"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para integrações e contratos distribuídos.
- **How to Fix:** Sincronize os valores ou unifique a definição em um módulo compartilhado.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
