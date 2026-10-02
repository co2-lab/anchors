---
title: "Gate: parent-valid"
description: "Ensures that o campo parent aponta para uma fase ou plano pai existente."
---

> **Gate Identifier:** `parent-valid` / `parent-valido`  
> **Code:** `PHORP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o campo parent aponta para uma fase ou plano pai existente.

---

## 2. Why It Matters

Evita subtarefas órfãs cujo pai foi excluído ou digitado errado.

---

## 3. How It Works

Verifies whether o identificador em `parent:` resolve para uma fase válida.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Fase pai existe. | None. Pipeline proceeds. |
| **`FAIL`** | Fase pai inexistente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Fases raiz (sem parent). | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: parent-valid
    on: [plan]
    check: parent-valid
    blocking: true
    measures: "o parent: aponta para uma fase que existe"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Corrija o identificador do parent no plano.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
