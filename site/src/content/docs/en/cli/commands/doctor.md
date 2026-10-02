---
title: "anchors doctor & audit"
description: "Ecosystem diagnostics, maturity checks, and single-file audits."
---

While `check` evaluates local gate rules, `anchors doctor` provides a global, holistic **X-ray of your entire ecosystem**.

---

## 1. Running `anchors doctor`

```bash
# Run standard health diagnosis
anchors doctor

# Run with verbose diagnostic breakdown
anchors doctor --verbose
```

### What `doctor` Detects:
1. **Orphaned Artifacts**: Specs without implementation code, or test files that verify deleted rules.
2. **Identity Code Collisions**: Duplicated codes like `AUTH-B01` declared in two separate files.
3. **Architectural Drift**: Discrepancies between declared layer boundaries in `anchors.yaml` and real-world imports.
4. **Coverage Voids**: Governed layers that have fallen below required maturity levels.

---

## 2. File Audit with `anchors audit`

When a specific file or module fails multiple checks, use `audit` to generate an actionable punch list:

```bash
anchors audit src/services/billing/invoice.spec.md
```

This outputs a unified dossier of every missing scenario, failed assertion, stale dependency, and layer violation associated with that unit.
