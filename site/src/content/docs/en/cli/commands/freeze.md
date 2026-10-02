---
title: "anchors freeze & thaw"
description: "How to freeze the repository for releases and unfreeze for development."
---

When a project reaches a release candidate or production milestone, you must prevent accidental code modifications while allowing inspection and audit commands to run.

---

## 1. Freezing the Project (`anchors freeze`)

```bash
# Freeze the project with an audit rationale
anchors freeze --reason="v2.0.0 Release Candidate verification"
```

When frozen:
- All mutating commands (`new`, `work`, `deliver`, `map build`) are strictly **blocked**.
- Read-only diagnostics (`status`, `doctor`, `guide`, `coverage`, `version`) **continue to function**.

---

## 2. Unfreezing (`anchors thaw`)

```bash
anchors thaw
```

Resumes normal development activities and unlocks the task queue.
