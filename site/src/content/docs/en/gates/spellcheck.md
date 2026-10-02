---
title: "Gate: spellcheck"
description: "Elimina erros de digitação e ortografia em identificadores, comentários e textos."
---

> **Gate Identifier:** `spellcheck` / `verificador-ortografico`  
> **Code:** `EXCMX` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code`, `doc`, `spec`, `feature` | **Layers:** `All Layers`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Elimina erros de digitação e ortografia em identificadores, comentários e textos.

---

## 2. Why It Matters

Erros de digitação geram nomes de variáveis bizarros, quebram buscas textuais e passam impressão de descuido.

---

## 3. How It Works

Executa o `typos` (binário nativo ultra-rápido) ou `cspell` sobre o repositório.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Sem erros ortográficos detectados. | None. Pipeline proceeds. |
| **`FAIL`** | Palavra com erro de digitação encontrada. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Ferramenta não instalada. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: spellcheck
    on: [code, doc, spec]
    scope: project
    run: "typos"
    needs_tool: typos
    install_hint: "brew install typos"
    blocking: false
    when: [pre-commit, ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no pre-commit; rode com `typos -w` para correção automática.
- **How to Fix:** Corrija a palavra no texto ou cadastre-a no dicionário do projeto (`.typos.toml`).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
