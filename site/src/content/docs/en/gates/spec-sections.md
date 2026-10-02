---
title: "Gate: spec-sections"
description: "Verifies whether The spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto."
---

> **Gate Identifier:** `spec-sections` / `spec-completa`  
> **Code:** `SFMSP` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `governed layers`, `Spec`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether The spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto.

---

## 2. Why It Matters

UmThe spec com texto solto em prosa não permite que a IA ou o CLI indexem as regras. Além disso, evita misturar idiomas nas seções.

---

## 3. How It Works

Busca regras com código em cabeçalhos (###), linhas de tabela (|) ou bullets em negrito. Verifica também se os títulos batem com o lang configurado.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Pelo menos uma regra está catalogada de forma estruturada e os títulos seguem o idioma configurado. | None. Pipeline proceeds. |
| **`FAIL`** | The spec só contém prosa corrida sem regras catalogadas, ou possui seções em idioma diferente do lang declarado. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Camadas declarativas ou specs de coordenação. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: spec-sections
    on: [spec]
    check: spec-sections
    blocking: true
    measures: "a spec tem ao menos uma regra catalogada estruturada"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking desde o início. Ensures that os templates de spec sejam preenchidos de verdade.
- **How to Fix:** Estruture suas regras usando títulos com código (### CODE-B01) ou tabelas, e use o idioma padrão do projeto nas seções.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
