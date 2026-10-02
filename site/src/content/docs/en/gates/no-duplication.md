---
title: "Gate: no-duplication"
description: "Detecta blocos de código idênticos ou quase idênticos copiados em múltiplos arquivos."
---

> **Gate Identifier:** `no-duplication` / `sem-duplicacao`  
> **Code:** `DUPLC` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Detecta blocos de código idênticos ou quase idênticos copiados em múltiplos arquivos.

---

## 2. Why It Matters

Código copiado e colado multiplica bugs: ao corrigir em um lugar, os outros continuam vulneráveis.

---

## 3. How It Works

Executa ferramentas como `jscpd`, `pmd cpd` ou `dupl` com limiar de tokens duplicados.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhum bloco duplicado acima do limiar encontrado. | None. Pipeline proceeds. |
| **`FAIL`** | Bloco duplicado copiado em dois ou mais arquivos. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Ferramenta não instalada. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: no-duplication
    on: [code]
    scope: project
    run: "npx --yes jscpd . --reporters console --silent"
    needs_tool: jscpd
    blocking: false
    when: [ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; promova a blocking após refatorar trechos repetidos.
- **How to Fix:** Extraia o bloco duplicado para uma função auxiliar ou módulo compartilhado.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
