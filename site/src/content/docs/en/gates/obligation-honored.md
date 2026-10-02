---
title: "Gate: obligation-honored"
description: "Ensures that deveres regulatórios declarados fora da unidade são cumpridos."
---

> **Gate Identifier:** `obligation-honored` / `obrigacao-honrada`  
> **Code:** `OBHNB` | **Category:** [Failure & Error Handling](/docs/gates/)  
> **Evaluates:** `spec`, `code` | **Layers:** `All Layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that deveres regulatórios declarados fora da unidade são cumpridos.

---

## 2. Why It Matters

Obrigações regulatórias não podem ser esquecidas durante o desenvolvimento ágil.

---

## 3. How It Works

Verifica o cumprimento de deveres mapeados a partir de pacotes de compliance.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todos os deveres regulatórios declarados estão honrados. | None. Pipeline proceeds. |
| **`FAIL`** | Dever regulatório não cumprido ou sem evidência. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Projetos sem catálogo de obrigações regulatórias. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: obligation-honored
    blocking: true
    measures: "os deveres regulatórios declarados são cumpridos"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para sistemas regulados (fintechs, saúde, privacidade).
- **How to Fix:** Implemente a exigência regulatória e forneça The test corresponding.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
