---
title: "Gate: has-code"
description: "Ensures that specification and feature files declare formal rule and scenario identity codes."
---

> **Gate Identifier:** `has-code` / `spec-tem-codigo`  
> **Code:** `SCIDS` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec`, `feature` | **Layers:** `governed layers`, `All Layers`  
> **Execution Model:** `Interno Determinístico` | **Default Blocking:** `true`

---

## 1. What This Gate Measures

Ensures that specification and feature files declare formal rule and scenario identity codes.

---

## 2. Why It Matters

Without identity codes (e.g. AUTH-B01), requirements are invisible to relational traceability checkers.

---

## 3. How It Works

Scans file content using regex matching against the project's identity code grammar.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`OK`** | The file contains at least one valid identity code (e.g. PREFIX-B01). | None. Pipeline proceeds. |
| **`FAIL`** | The file contains no rule or scenario identity codes. | **Pipeline blocked.** Correct the discrepancy. |
| **`WARN`** | Not applicable. | Review before the next release cycle. |
| **`SKIP`** | Binary files or declarative recognized layers. | Exemption recorded in audit report. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: has-code
    on: [spec, feature]
    check: has-code
    blocking: true
    measures: "o arquivo carrega um código de cenário"
```

---

## 6. Recommendations & Remediation

- **Best Practice:** Blocking from day one. Without identity codes, automated traceability cannot operate.
- **How to Fix:** Add formal identity codes to the rules in the file (e.g. `### AUTH-B01 — Rule Title`).

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): The core bundle of Spec, Feature, Test, and Code.
- [Project Layers](/docs/layers/): Architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): Understanding how checks operate.
- [Gates Catalog](/docs/gates/): Master index of all 95 verification gates.
