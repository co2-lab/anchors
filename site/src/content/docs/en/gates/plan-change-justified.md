---
title: "Gate: plan-change-justified"
description: "Ensures that um plano ou spec modificado declare no diff por que a mudança aconteceu."
---

> **Gate Identifier:** `plan-change-justified` / `plano-alterado-justificado`  
> **Code:** `PCJPL` | **Category:** [Planning & Progress](/docs/gates/)  
> **Evaluates:** `plan`, `spec` | **Layers:** `Planos`, `Spec`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that um plano ou spec modificado declare no diff por que a mudança aconteceu.

---

## 2. Why It Matters

A deriva silenciosa é perigosa: alterar um plano sem justificativa apaga a inconsistência sem deixar rastro do motivo.

---

## 3. How It Works

Analisa o git diff e Checks whether há nota explicativa justificando a alteração dThe specification.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A alteração nThe file acompanha justificativa de mudança. | None. Pipeline proceeds. |
| **`FAIL`** | Alteração estrutural de plano/spec sem justificativa no commit/diff. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Sem git para inspecionar diff. | Review before the next release cycle. |
| **`SKIP`** | Criação de arquivos novos. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: plan-change-justified
    on: [plan, spec]
    check: plan-change-justified
    blocking: true
    measures: "um plano/spec que mudou declara por quê"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para equipes que precisam de auditoria rigorosa de mudanças.
- **How to Fix:** Adicione uma nota ou entrada de revisão explicando a razão da alteração nThe spec ou plano.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
