---
title: "Gate: sibling-guard"
description: "Prevents módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos."
---

> **Gate Identifier:** `sibling-guard` / `guarda-de-irmaos`  
> **Code:** `SBGRD` | **Category:** [Architectural Boundaries](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Prevents módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos.

---

## 2. Why It Matters

Evita acoplamento cruzado desordenado entre serviços que estão no mesmo nível hierárquico.

---

## 3. How It Works

Verifica os caminhos de chamada e tipos de parâmetros entre módulos do mesmo diretório ou nível.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Tratamento de parâmetros consistente e sem acoplamento proibido. | None. Pipeline proceeds. |
| **`FAIL`** | Inconsistência de parâmetros entre funções irmãs ou acesso indevido. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: sibling-guard
    blocking: true
    measures: "módulos irmãos tratam o mesmo parâmetro de forma consistente"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos médios e grandes.
- **How to Fix:** Padronize os parâmetros das funções irmãs ou utilize um DTO comum.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
