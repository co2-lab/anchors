---
title: "Gate: flag-scenario-grammar"
description: "Ensures that os scenarios declarados em arquivos de feature flags seguem a gramática correta."
---

> **Gate Identifier:** `flag-scenario-grammar` / `gramatica-de-flag`  
> **Code:** `FLSCF` | **Category:** [Doctrine & Feature Flags](/docs/gates/)  
> **Evaluates:** `flag` | **Layers:** `flags`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that os scenarios declarados em arquivos de feature flags seguem a gramática correta.

---

## 2. Why It Matters

Padroniza os nomes de estados de flags (ex: ON, OFF, ROLLOUT) para que ferramentas possam consumi-los.

---

## 3. How It Works

Analisa o cabeçalho e os subtítulos dThe file `flags/*.flag.md` contra o padrão de gramática de flags.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A flag segue o formato padrão de declaração. | None. Pipeline proceeds. |
| **`FAIL`** | Arquivo de flag mal formatado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: flag-scenario-grammar
    on: [flag]
    check: flag-scenario-grammar
    blocking: true
    measures: "cenários de flag seguem a gramática correta"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato para arquivos de flag.
- **How to Fix:** Siga o formato padrão de cabeçalho `### ON — Descrição` nThe file de flag.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
