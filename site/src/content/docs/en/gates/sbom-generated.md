---
title: "Gate: sbom-generated"
description: "Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias."
---

> **Gate Identifier:** `sbom-generated` / `sbom-gerado`  
> **Code:** `EXCMX` | **Category:** [Security & Hygiene](/docs/gates/)  
> **Evaluates:** `code` | **Layers:** `All Layers`, `code`  
> **Execution Model:** `Externo (run)` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias.

---

## 2. Why It Matters

Essencial para conformidade com padrões de segurança em distribuição de software corporativo.

---

## 3. How It Works

Executa ferramentas como `syft` gerando The file `sbom.json`.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Arquivo de SBOM gerado com sucesso. | None. Pipeline proceeds. |
| **`FAIL`** | Falha na geração do inventário. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Syft não instalado. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: sbom-generated
    on: [code]
    scope: project
    run: "syft scan dir:. -o cyclonedx-json=sbom.json -q"
    needs_tool: syft
    install_hint: "brew install syft"
    blocking: true
    when: [ci]
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking no CI de release.
- **How to Fix:** Instale o `syft` e verifique as permissões de escrita dThe file sbom.json.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
