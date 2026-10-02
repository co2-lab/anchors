---
title: "The Graph & Dependency Map"
description: "How Anchors builds a directed acyclic graph (DAG) of your entire project to calculate impact and propagation."
---

Modern software systems are not flat lists of files. They are interconnected webs of dependencies: a database query depends on a domain model, an API handler depends on a service, and a user interface depends on an API contract.

Anchors models your entire repository as a **Directed Acyclic Graph (DAG)** of anchors and files.

---

## 1. What is the Anchors Graph?

When you run `anchors analyze` or evaluate gates, Anchors scans your project and builds an in-memory graph where:
- **Nodes** are project artifacts: specs, feature files, test files, source code files, doctrine guides, plans, and feature flags.
- **Edges** are formal relationships between artifacts:
  - `governs`: Spec -> Code (the spec dictates how the code behaves).
  - `covered-by`: Spec -> Feature (the feature file tests the spec).
  - `tested-by`: Feature -> Test (the test code executes the feature scenario).
  - `depends_on`: Code -> Code (import/call dependency between modules).
  - `realizes`: Spec -> Doctrine (the spec realizes a higher-level product doctrine).
  - `flagged-by`: Unit -> Flag (the unit is governed by a feature flag).

```
  [ Product Doctrine ]
           │
           │ (realizes)
           ▼
    [ Order Spec ] ──────(governs)───────► [ Order Code ]
           │                                      │
           │ (covered-by)                         │ (depends_on)
           ▼                                      ▼
   [ Order Feature ] ──(tested-by)──► [ Payment Service ]
           │
           ▼
     [ Order Test ]
```

---

## 2. Why the Graph Matters

Traditional linters only look at one file at a time in isolation. They cannot detect systemic architectural issues.

Because Anchors maintains a project-wide graph, it can execute **relational gates**:
1. **Detecting Circular Dependencies**: The [`circular`](/docs/gates/circular/) gate traces topological cycles across layers.
2. **Boundary Enforcement**: The [`layer-boundary`](/docs/gates/layer-boundary/) gate prevents domain entities from importing infrastructure or UI layers.
3. **Change Propagation**: When an anchor changes, Anchors traverses downstream edges to identify all impacted units that must be re-verified.
4. **Ghost Detection**: The [`revision-orphans`](/docs/gates/revision-orphans/) gate discovers specs that were removed or renamed while leaving orphaned test code behind.

---

## 3. Querying the Graph with the CLI

You can inspect the dependency graph directly using the Anchors CLI:

```bash
# Visualize the topological order of all units
anchors graph --format=tree

# Check impact wave if a specific unit changes
anchors wave src/services/billing/invoice.spec.md

# Inspect inbound and outbound edges for an anchor
anchors inspect BILL-B01
```

---

## 4. Related Concepts

- [The Unit](/docs/concepts/unit/): The local cluster of nodes forming each unit.
- [Project Layers](/docs/layers/): The architectural boundaries defining allowed edges in the graph.
- [Propagation & Impact](/docs/concepts/propagation-and-impact/): How changes travel across graph edges.
