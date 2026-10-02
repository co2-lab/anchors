---
title: "Gate: no-secret-leaked"
description: "Ensures that nenhum segredo, token, senha ou chave privada entre no histórico do repositório."
---

> **Gate Identifier:** `no-secret-leaked` / `secret-nao-vazado`  
> **Code:** `EXCMX` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code`, `test`, `doc`, `spec`, `feature`, `guide`, `plan` | **Layers:** `All Layers`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that nenhum segredo, token, senha ou chave privada entre no histórico do repositório.

---

## 2. Why It Matters

Segredo commitado fica eternamente no Git. Revogar credenciais custa caro e expõe a empresa.

---

## 3. How It Works

Executa o `gitleaks` sobre o diff ou o projeto inteiro procurando entropia de senhas e padrões de chaves conhecidas.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Zero segredos detectados. | None. Pipeline proceeds. |
| **`FAIL`** | Segredo ou chave privada encontrada em arquivos do commit. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Gitleaks não instalado. | Review before the next release cycle. |
| **`SKIP`** | Linhas cadastradas na allowlist em .gitleaks.toml. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: no-secret-leaked
    on: [code, test, doc, spec, feature, guide, plan]
    scope: batch
    run: "gitleaks git --no-banner --redact -v"
    needs_tool: gitleaks
    install_hint: "brew install gitleaks"
    blocking: true
    when: [pre-commit, pre-push, ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking desde o minuto zero em pre-commit e CI.
- **How to Fix:** Remova a credencial dThe code, use variáveis de ambiente e configure o `.gitleaks.toml` se for um falso positivo comprovado.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
