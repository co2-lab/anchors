---
title: "Gate: docs-fresh"
description: "Ensures that a documentação compilada reflete as specifications atuais sem defasagem."
---

> **Gate Identifier:** `docs-fresh` / `documentos-frescos`  
> **Code:** `DCFRD` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `doc` | **Layers:** `doc`, `docs`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a documentação compilada reflete as specifications atuais sem defasagem.

---

## 2. Why It Matters

Se The spec mudou e a documentação compilada não foi regerada, o site de documentação mente para os usuários.

---

## 3. How It Works

Verifies whether os arquivos gerados em docs/ são mais novos do que as specs das quais derivam.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A documentação compilada está em dia com as specs. | None. Pipeline proceeds. |
| **`FAIL`** | A documentação compilada está desatualizada (stale). | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Projetos sem compilação de docs. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: docs-fresh
    blocking: true
    measures: "a documentação compilada reflete a spec atual"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em CI.
- **How to Fix:** Recompile a documentação com `anchors docs build`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
