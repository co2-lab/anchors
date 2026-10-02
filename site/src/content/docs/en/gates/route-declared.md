---
title: "Gate: route-declared"
description: "Ensures that uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação."
---

> **Gate Identifier:** `route-declared` / `rota-declarada`  
> **Code:** `RTDCL` | **Category:** [Rules Usage & Contracts](/docs/gates/)  
> **Evaluates:** `spec` | **Layers:** `UI`, `Telas`, `API`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação.

---

## 2. Why It Matters

Impede o surgimento de telas órfãs que existem no repositório mas ninguém sabe como acessar.

---

## 3. How It Works

Verifica a presença da seção de Rotas ou Cabeçalho de Navegação nThe spec.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | A rota de entrada e vizinhos estão formalmente declarados. | None. Pipeline proceeds. |
| **`FAIL`** | Spec de tela ou rota sem declaração de caminho de acesso. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Camadas que não representam pontos de entrada de usuário ou API. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: route-declared
    on: [spec]
    check: route-declared
    blocking: true
    measures: "uma tela declara como se chega nela e seus vizinhos"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking em projetos de frontend (React, Flutter, Vue) e APIs.
- **How to Fix:** Declare a rota na seção de Identidade ou Rotas dThe spec.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
