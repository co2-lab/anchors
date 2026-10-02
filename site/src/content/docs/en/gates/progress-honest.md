---
title: "Gate: progress-honest"
description: "Ensures that The file de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco."
---

> **Gate Identifier:** `progress-honest` / `progresso-honesto`  
> **Code:** `PRHNP` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan`, `doc` | **Layers:** `Planos`, `doc`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that The file de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco.

---

## 2. Why It Matters

Evita checklists com tarefas marcadas como concluídas [x] quando The file ou teste ainda não existe no repositório.

---

## 3. How It Works

Cruza itens marcados como concluídos com a existência real dos arquivos e testes no disco.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os itens marcados como prontos existem de fato no disco. | None. Pipeline proceeds. |
| **`FAIL`** | Item marcado como concluído no checklist aponta para arquivo inexistente. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: progress-honest
    blocking: true
    measures: "o progresso diz a verdade sobre o disco"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para manter a gestão transparente.
- **How to Fix:** Crie The file que falta ou desmarque o item [ ] até que o trabalho seja concluído.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
