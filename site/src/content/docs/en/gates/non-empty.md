---
title: "Gate: non-empty"
description: "Ensures that The file não é um esqueleto vazio e que a feature possui scenarios de verdade."
---

> **Gate Identifier:** `non-empty` / `feature-nao-vazia`  
> **Code:** `FTMFT` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `feature`, `spec`, `doc` | **Layers:** `governed layers`, `Feature`, `Spec`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that The file não é um esqueleto vazio e que a feature possui scenarios de verdade.

---

## 2. Why It Matters

Um arquivo .feature com 8 linhas de comentários ou cabeçalho Gherkin não é vazio em bytes, mas não declara scenario nenhum. Este gate pega essas cascas vazias.

---

## 3. How It Works

Verifies whether The file tem conteúdo não-espaço e, para features, valida a presença de palavras-chave de scenario (Scenario: ou scenario:).

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The file possui conteúdo real e, no caso de feature, contém scenarios declarados. | None. Pipeline proceeds. |
| **`FAIL`** | Arquivo vazio, só com espaços ou feature sem nenhum scenario declarado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: feature-not-empty
    on: [feature]
    check: non-empty
    blocking: true
    measures: "a feature tem cenários de verdade"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking imediato.
- **How to Fix:** Escreva ao menos um scenario concreto com Dado/Quando/Então na feature.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
