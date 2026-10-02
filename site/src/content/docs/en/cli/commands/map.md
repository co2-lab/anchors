---
title: "anchors map & impact"
description: "Building the dependency graph, analyzing impact waves, and recoding."
---

The Anchors dependency map (`anchors.graph.yaml`) models the entire repository as a Directed Acyclic Graph (DAG).

---

## 1. Managing the Map

```bash
# Rebuild the dependency graph from files on disk
anchors map build

# Display graph summary and top-level node metrics
anchors map show

# Check for circular dependency loops
anchors map circular
```

---

## 2. Impact Wave Analysis (`anchors impact`)

Before modifying a spec or core domain file, calculate the downstream blast radius:

```bash
anchors impact src/domain/user.spec.md
```

Outputs:
- **Direct Dependents**: Services and handlers directly calling this unit.
- **Downstream Impact**: Integration tests and contracts that will become **stale**.
- **Recommended Test Scope**: The minimal set of test suites required to validate the change.

---

## 3. Atomic Identity Renaming (`anchors recode`)

Need to rebrand or reorganize a domain module from `USER` to `ACCOUNT`?

```bash
anchors recode USER-B01 ACC-B01
```

This atomically updates the code in the spec (`*.spec.md`), the feature (`*.feature`), test assertions (`*_test.go`), code annotations, and the dependency graph in one step.
