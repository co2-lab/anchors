---
title: "Gate: count-honored"
description: "Ensures that asserções numéricas escritas nThe spec batem com os números reais nThe code."
---

> **Gate Identifier:** `count-honored` / `contagem-honrada`  
> **Code:** `CNHNC` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `governed layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `false`

---

## 1. What This Gate Measures

Ensures that asserções numéricas escritas nThe spec batem com os números reais nThe code.

---

## 2. Why It Matters

Se The spec diz 'o lote processa no máximo 50 itens', The code não pode ter uma constante `MAX_BATCH = 100`.

---

## 3. How It Works

Extrai números e limites citados nas regras dThe spec e confere com as constantes nThe code corresponding.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Os limites numéricos batem entre spec e código. | None. Pipeline proceeds. |
| **`FAIL`** | Divergência numérica entre The spec e a constante nThe code. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Código pendente. | Review before the next release cycle. |
| **`SKIP`** | Specs sem asserções numéricas. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: count-honored
    blocking: false
    measures: "asserções numéricas da spec batem com o código"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Informativo no início; promova a blocking ao estabilizar constantes de negócio.
- **How to Fix:** Sincronize o número nThe spec ou ajuste o valor da constante nThe code.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
