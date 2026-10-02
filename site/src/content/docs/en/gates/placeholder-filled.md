---
title: "Gate: placeholder-filled"
description: "Verifies whether os placeholders deixados por geradores ou templates foram preenchidos."
---

> **Gate Identifier:** `placeholder-filled` / `placeholder-preenchido`  
> **Code:** `PLCFL` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec`, `plan`, `feature`, `doc` | **Layers:** `All Layers`  
> **Execution Model:** `Relacional com Grafo` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Verifies whether os placeholders deixados por geradores ou templates foram preenchidos.

---

## 2. Why It Matters

Impede commitar arquivos com textos como '[Descreva aqui]' ou 'TODO: preencher'.

---

## 3. How It Works

Busca marcadores universais de template não preenchidos no corpo do documento.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | Nenhum placeholder de template esquecido. | None. Pipeline proceeds. |
| **`FAIL`** | Encontrado texto de placeholder não substituído. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Not applicable. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: placeholder-filled
    on: [spec, feature]
    check: placeholder-filled
    blocking: true
    measures: "o esqueleto emitido pelo gerador foi preenchido"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** blocking para evitar entrega de documentação pela metade.
- **How to Fix:** Substitua os textos de placeholder pelo conteúdo real da regra.

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
