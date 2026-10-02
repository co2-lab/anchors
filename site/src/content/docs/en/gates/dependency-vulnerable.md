---
title: "Gate: dependency-vulnerable"
description: "Audita dependencies e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs)."
---

> **Gate Identifier:** `dependency-vulnerable` / `dependencia-vulneravel`  
> **Code:** `EXCMX` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Audita dependencies e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs).

---

## 2. Why It Matters

Usar pacotes vulneráveis abre portas para invasões e violações de segurança.

---

## 3. How It Works

Invoca ferramentas como `osv-scanner`, `govulncheck`, `pnpm audit` ou `pip-audit`.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhuma CVE conhecida encontrada nas dependencies. | None. Pipeline proceeds. |
| **`FAIL`** | Vulnerabilidade conhecida detectada. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Ferramenta de auditoria não instalada. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: dependency-vulnerable
    on: [code]
    scope: project
    run: "osv-scanner scan source -r ."
    needs_tool: osv-scanner
    install_hint: "brew install osv-scanner"
    blocking: true
    when: [pre-push, ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no pre-push e CI.
- **How to Fix:** Atualize a versão da biblioteca vulnerável no seu package.json, go.mod ou requirements.txt.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
