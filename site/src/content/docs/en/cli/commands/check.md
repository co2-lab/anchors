---
title: "anchors check & verify"
description: "How to execute quality gates in CLI and CI pipelines."
---

The `anchors check` command is the primary verification engine of the Anchors framework. It evaluates the project's quality gates against modified, staged, or entire repository files.

---

## 1. Syntax & Common Options

```bash
# Check modified files in working directory
anchors check

# Check staged files (used by git pre-commit hook)
anchors check --staged

# Check all governed files across the entire project
anchors check --all

# Check only a specific gate
anchors check --gate=unit-complete

# Run in strict mode (fails on warnings)
anchors check --strict
```

### Exit Codes:
- `0`: All gates passed cleanly (`OK`).
- `1`: One or more blocking gates returned `FAIL`.
- `2`: Configuration or file parsing error.

---

## 2. Difference between `check` and `verify`

- **`anchors check`**: Focuses on **internal Anchors gates** (unit completeness, identity code integrity, boundary violations, doc parity).
- **`anchors verify`**: The **full-phase orchestrator**. It runs `anchors check`, executes project test suites (`anchors test`), runs linter tools, and performs mutation checks before approving a release milestone.

---

## 3. CI/CD Usage Example

```yaml
# In GitHub Actions or GitLab CI
- name: Evaluate Anchors Pipeline
  run: anchors check --all
```
