---
title: "Gate: deadcode"
description: "Identifica funções, tipos, exports e arquivos órfãos não consumidos."
---

> **Gate Identifier:** `deadcode` / `codigo-morto`  
> **Code:** `EXCMX` | **Category:** [Architectural Boundaries](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Identifica funções, tipos, exports e arquivos órfãos não consumidos.

---

## 2. Why It Matters

Código morto gera débito técnico, confunde IAs e aumenta o tempo de compilação sem agregar valor.

---

## 3. How It Works

Executa ferramentas como `deadcode` (Go), `knip` (TS/JS) ou `vulture` (Python).

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhum export ou código órfão encontrado. | None. Pipeline proceeds. |
| **`FAIL`** | Código morto detectado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Ferramenta não instalada. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: deadcode
    on: [code]
    scope: project
    run: "deadcode ./..."
    needs_tool: deadcode
    blocking: false
    when: [ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; limpe os órfãos periodicamente.
- **How to Fix:** Apague as funções e exports não utilizados ou consuma-os onde for necessário.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
