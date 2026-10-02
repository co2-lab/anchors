---
title: "Gate: test-level-codes"
description: "Ensures that cada nível de teste referencia apenas códigos permitidos para seu escopo."
---

> **Gate Identifier:** `test-level-codes` / `codigos-de-nivel-de-teste`  
> **Code:** `TLVCD` | **Category:** [Proof & Execution](/docs/gates/)  
> **Evaluates:** `feature` | **Layers:** `governed layers`, `Feature`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that cada nível de teste referencia apenas códigos permitidos para seu escopo.

---

## 2. Why It Matters

Impede testes unitários de referenciarem regras visuais (-VR) e vice-versa.

---

## 3. How It Works

Filtra os scenarios de acordo com as regras de allow e exclude configuradas para cada nível no anchors.yaml.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhum código proibido pelo nível foi referenciado. | None. Pipeline proceeds. |
| **`FAIL`** | scenario de um nível (ex: @unit) referencia código proibido (ex: LOGIN-VR). | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Níveis sem restrição declarada. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: test-level-codes
    on: [feature]
    check: test-level-codes
    blocking: true
    levels:
      nivel-unit: { exclude: ['-VR$'] }
      nivel-vr:   { allow: ['-VR$'] }
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos com múltiplos regimes de teste.
- **How to Fix:** Mova o scenario para a feature do nível corresponding ou ajuste o identity code da regra.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
