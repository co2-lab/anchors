---
title: "Gate: updated-at-atual"
description: "Ensures that o campo updated_at no cabeçalho reflete a data real da alteração ou commit."
---

> **Gate Identifier:** `updated-at-atual` / `updated-at-atual`  
> **Code:** `INCHN` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec`, `feature`, `code`, `test`, `doc` | **Layers:** `All Layers`  
> **Execution Model:** `Com Root / Git` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that o campo updated_at no cabeçalho reflete a data real da alteração ou commit.

---

## 2. Why It Matters

Prevents o cabeçalho minta sobre quando The file foi modificado pela última vez.

---

## 3. How It Works

Compara a data `updated_at: YYYY-MM-DD` com a data do sistema (se houver alteração não commitada) ou a data do último commit Git.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A data declarada bate com a data da alteração ou commit. | None. Pipeline proceeds. |
| **`FAIL`** | Data desatualizada em relação à data da alteração. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Arquivo sem commits prévios. | Review before the next release cycle. |
| **`SKIP`** | Ambientes sem Git ou arquivos sem updated_at declarado. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: updated-at-atual
    blocking: true
    measures: "o updated_at do cabeçalho bate com a data do commit"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking; use `anchors check --fix` para atualizar a data automaticamente.
- **How to Fix:** Atualize a linha `updated_at: AAAA-MM-DD` no cabeçalho para a data de hoje, ou rode `anchors check --fix`.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
