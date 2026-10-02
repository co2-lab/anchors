---
title: "Gate: vr-baseline"
description: "Ensures that scenarios de regressão visual (-VR) possuem imagens de baseline capturadas."
---

> **Gate Identifier:** `vr-baseline` / `baseline-visual`  
> **Code:** `VRBSV` | **Category:** [Presentation & UI](/docs/gates/)  
> **Evaluates:** `feature`, `test` | **Layers:** `UI`, `Telas`, `VR`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that scenarios de regressão visual (-VR) possuem imagens de baseline capturadas.

---

## 2. Why It Matters

Um scenario visual sem imagem de baseline não pode ser verificado contra regressão de aparência.

---

## 3. How It Works

Verifies whether existe um arquivo PNG corresponding na pasta de baselines para cada scenario com código -VR.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as variantes visuais possuem imagem de baseline arquivada. | None. Pipeline proceeds. |
| **`FAIL`** | scenario -VR sem imagem de baseline capturada. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Baseline ainda não capturado. | Review before the next release cycle. |
| **`SKIP`** | scenarios não visuais. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: vr-baseline
    blocking: true
    measures: "cenários visuais possuem imagem de baseline capturada"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para equipes de frontend e design system.
- **How to Fix:** Capture o screenshot de baseline executando a suíte de teste visual.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
