---
title: "Layers & Regimes"
description: "Understand how Anchors organizes code into Governed and Recognized layers, and applies testing regimes."
---

Not all code in a project serves the same purpose or requires the same level of ceremony. Writing full Gherkin specifications and BDD scenarios for a simple database adapter that runs `SELECT * FROM users` creates unnecessary overhead and encourages developers to invent fake rules.

Anchors solves this by categorizing files into **Layers** and assigning them distinct **Regimes**.

---

## 1. Governed Layers vs. Recognized Layers

```
┌────────────────────────────────────────────────────────┐
│                   PROJECT LAYERS                       │
└───────────────────────────┬────────────────────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        ▼                                       ▼
┌───────────────────────────────┐ ┌───────────────────────────────┐
│        GOVERNED LAYERS        │ │       RECOGNIZED LAYERS       │
│     (Behavioral Regime)       │ │     (Declarative Regime)      │
├───────────────────────────────┤ ├───────────────────────────────┤
│ • Contains business rules     │ │ • No business rules originated│
│ • Requires The Unit           │ │ • Pure adapters / glue code   │
│   (Spec + Feature + Test)     │ │ • Declares layer identity     │
│ • e.g.: usecase, service      │ │ • e.g.: dao, infra, types     │
└───────────────────────────────┘ └───────────────────────────────┘
```

### A. Governed Layers (Behavioral Regime)
Governed layers contain the **decision-making logic** of your application: validation rules, financial formulas, state machines, and business workflows.
- **Regime**: `regime: rule` (behavioral).
- **Mandate**: Every unit in these layers **must** have a complete [Unit](/docs/concepts/unit/) (Spec, Feature, Test, Code).
- **Typical Layers**: `usecase`, `service`, `domain`, `command`.

### B. Recognized Layers (Declarative Regime)
Recognized layers contain **technical plumbing**: database connections, HTTP clients, cloud infrastructure scripts, and type definitions.
- **Regime**: `regime: declarative`.
- **Mandate**: These files do **not** require business specifications or Gherkin features. They only require structural identity (e.g. declaring `layer: dao` in their `@anchors` header).
- **Typical Layers**: `dao`, `infra`, `dto`, `types`, `doc`.

---

## 2. Testing Regimes

Different layers require different types of verification. Anchors recognizes four canonical **Testing Regimes**:

| Regime | What it validates | Typical Location | Recommended Layers |
| --- | --- | --- | --- |
| **Unit** | Pure algorithmic logic without I/O or network dependencies. | `src/**/unit/*_test.go` | `domain`, `service`, `usecase` |
| **Contract** | Request/response schemas, API contracts, and serialized payloads. | `tests/contracts/**/*.spec.ts` | `api`, `dto`, `client` |
| **Integration** | Communication between units, databases, and message brokers. | `tests/integration/**/*.test.ts`| `service`, `dao`, `infra` |
| **Visual / E2E** | Rendered user interfaces, UI state transitions, and pixel regressions. | `tests/visual/**/*.spec.ts` | `ui`, `component`, `page` |

---

## 3. Boundary Gates

To keep your architecture clean and maintainable, Anchors provides boundary gates:
- [`layer-boundary`](/docs/gates/layer-boundary/): Enforces architectural direction (e.g. domain layers cannot import infrastructure).
- [`proof-crosses-boundary`](/docs/gates/proof-crosses-boundary/): Ensures integration tests verify that cross-layer contracts are respected.
- [`sibling-guard`](/docs/gates/sibling-guard/): Prevents tight coupling between parallel modules within the same architectural layer.

---

## 4. Further Reading

- Read the complete [Project Layers Guide](/docs/layers/) for configuration examples and boundary matrices.
- Learn about [The Unit](/docs/concepts/unit/).
