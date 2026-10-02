---
title: "Gate: marker-parity"
description: "Ensures that a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend)."
---

> **Gate Identifier:** `marker-parity` / `paridade-de-marcador`  
> **Code:** `MRPRM` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `code`, `spec` | **Layers:** `All Layers`  
> **Execution Model:** `Composto / ComGate` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend).

---

## 2. Why It Matters

Evita assimetria regulatória: a tela diz que apaga o dado do usuário, mas o backend esquece de apagar.

---

## 3. How It Works

Verifies whether regras marcadas com prefixos simétricos possuem a contraparte corresponding na outra ponta.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A mesma regra está presente em ambas as pontas. | None. Pipeline proceeds. |
| **`FAIL`** | Regra simétrica presente em um lado mas ausente no outro. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: marker-parity
    blocking: true
    marker_prefix: "LGPD-"
    measures: "a mesma regra aparece nas DUAS pontas que a cumprem"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos com dados regulados (LGPD, saúde, financeiro).
- **How to Fix:** Implemente e marque a regra na ponta que faltava.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
