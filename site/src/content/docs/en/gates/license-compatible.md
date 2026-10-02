---
title: "Gate: license-compatible"
description: "Impede a inclusão de dependencies com licenças incompatíveis ou copyleft forte (como AGPL)."
---

> **Gate Identifier:** `license-compatible` / `licenca-compativel`  
> **Code:** `EXCMX` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Impede a inclusão de dependencies com licenças incompatíveis ou copyleft forte (como AGPL).

---

## 2. Why It Matters

Evita riscos jurídicos de contaminação de código proprietário por licenças não comerciais.

---

## 3. How It Works

Executa `go-licenses`, `license-checker` ou `cargo-deny` comparando com a política de licenças.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Todas as dependencies possuem licenças permitidas. | None. Pipeline proceeds. |
| **`FAIL`** | dependency com licença proibida ou desconhecida detectada. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Ferramenta de licenças não instalada. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: license-compatible
    on: [code]
    scope: project
    run: "go-licenses check ./... --disallowed_types=forbidden,restricted"
    needs_tool: go-licenses
    blocking: true
    when: [pre-push, ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no CI.
- **How to Fix:** Substitua a biblioteca com licença proibida por uma alternativa de licença permissiva (MIT, Apache-2.0, BSD).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
