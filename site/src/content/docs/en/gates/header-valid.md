---
title: "Gate: header-valid"
description: "Every file Anchors governs carries the @anchors header at its top, with its identity."
---

> **Gate Identifier:** `header-valid`  
> **Code:** `INCHN` | **Category:** [The Unit & Structure](/docs/gates/)  
> **Evaluates:** `spec`, `feature`, `code`, `test`, `guide`, `doc`, `plan`, `product`, `flag` | **Layers:** `All Layers`  
> **Execution Model:** `Deterministic, internal` | **Fixable:** `anchors check --fix`

---

## 1. What This Gate Measures

Every file Anchors governs carries the `@anchors` header **at its top**, with its identity.

---

## 2. Why It Matters

The header is how a file enters the map: it says which unit the file belongs to. A file without it is invisible to the gates that read a file's unit — the reach of a test, the references of a code file, the layer it declares.

---

## 3. How It Works

The header is the `@anchors` block at the top of the file, read exactly as the map reads it: before it, only blank lines, comments and a shebang. An `@anchors` further down — an example in a guide, a string in the code — is text, not the header. A comment opens the header only when `@anchors` is its first word.

The identity it demands depends on the file's role:

| File | Identity |
| --- | --- |
| a spec, a plan, a product doctrine, a flag | `code:` — it owns its identity |
| code, a test, a feature | `ref:` — the code of the spec it belongs to (`ref: A, B` when shared) |
| a guide, a document, a test support file | `layer:` — it belongs to no unit |

The header is read in every comment dialect: `//`, `#`, `--`, `<!-- -->` and a block comment's ` * `.

---

## 4. Verdict Conditions

| Verdict | Condition | Action Required |
| --- | --- | --- |
| **`PASS`** | The header is at the top, with the identity the file's role demands. | None. |
| **`FAIL`** | No header at the top, a header below the top with no `@fixed-header: <why>`, or a header without identity. | Run `anchors check --fix`, or write the header. |
| **`SKIP`** | Binary files and external runner scripts. | None. |

---

## 5. Configuration Example (`anchors.yaml`)

```yaml
gates:
  - name: header-valid
    blocking: true
```

A bare entry inherits the canonical `on:` — every governed kind.

---

## 6. Remediation

`anchors check --all --fix` writes the header a file lacks, from the map: the `ref:` of the unit the file belongs to — the spec that specifies it, or through its feature the spec of a test — or the `layer:` of a guide, a document or a test support file. A header with no identity gets the missing line; nothing already written is changed. A shebang stays first, and Go's `//go:build` stays before the package clause.

A file the map ties to no unit is left as it is: its identity is a decision to make.

```go
// @anchors
//   ref: INVCE

package billing
```

---

## 7. Related Concepts & Gates

- [The Unit](/docs/concepts/unit/): the spec and every sibling it regulates.
- [Project Layers](/docs/layers/): architectural boundaries and testing regimes.
- [Gates & Verdicts](/docs/concepts/gates-and-verdicts/): how checks operate.
- [Gates Catalog](/docs/gates/): every verification gate.
