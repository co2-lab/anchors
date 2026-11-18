---
title: "Project Layers"
description: "What layers are in Anchors, how governed and recognized layers work, and which gates apply to each."
---

If you have never used **Anchors**, imagine you are building a house. You don't start by laying bricks randomly in arbitrary corners of the lot: first you have a **blueprint**. The blueprint defines where the foundation sits, where the plumbing runs, where electrical wiring is placed, and where structural walls stand.

In Anchors, **Layers** are the architectural **blueprint** of your software. They tell developers and artificial intelligence **where each piece of code must live**, **who can talk to whom**, and **what level of rigor is required for each file**.

Without a clear blueprint, an AI pair programmer suffers from structural amnesia: it invents new directories on every prompt, imports databases directly inside UI components, duplicates business rules, and silently degrades your architecture.

---

## 1. What is a Layer?

A **Layer** is an architectural grouping of files sharing the same role in the system. Rather than treating all files in the repository as equals, Anchors classifies every file using directory and extension patterns (`pattern`) declared in [`anchors.yaml`](/docs/anchors-yaml/).

Every layer possesses well-defined properties:
- **`pattern`**: The file mask (glob) belonging to it (e.g. `src/services/**/*.ts`).
- **`kind`**: The artifact type (`code`, `spec`, `feature`, `test`, `doc`, `guide`, `plan`, `product`, `flag`).
- **`regime`**: Whether the layer is behavioral (`rule`) or structural (`declarative`).
- **`tags`**: Keywords used by [gates](/docs/gates/) to know how to evaluate the layer.
- **`depends_on`** / **`governs`**: Import rules and architectural boundary limits.

---

## 2. The Two Major Layer Families

This is the most critical distinction in Anchors: **not every file requires a full business specification**. Anchors divides layers into two distinct families:

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
Governed layers house the **intelligence, business decisions, and domain rules** of your system: validation rules, financial logic, decision flows, and core workflows.

- **Regime**: `regime: rule` (behavioral).
- **Mandate**: Every code unit in these layers requires a complete [Unit](/docs/concepts/unit/):
  1. A **Spec** (`.spec.md`) cataloging its business rules with [identity codes](/docs/concepts/traceability-and-codes/).
  2. A **Feature** (`.feature`) with Gherkin scenarios testing the rules.
  3. A **Test** (`_test.go`, `.test.ts`) executing and proving the scenarios.
  4. The **Code** implementing the functions.
- **Common Examples**:
  - `usecase`: Application use cases and business orchestrations.
  - `service`: Domain services.
  - `domain`: Core domain models and entities.
  - `command`: CLI executable commands and API endpoint handlers.

### B. Recognized Layers (Declarative Regime)
Recognized layers contain technical plumbing that **does not originate business rules**. They exist to connect parts of the system: database drivers, HTTP clients, low-level type declarations, infrastructure code, or documentation files.

- **Regime**: `regime: declarative`.
- **Mandate**: They **DO NOT** require functional specifications or Gherkin features. Forcing a database adapter that executes `SELECT * FROM users` to have a business spec would encourage fake, meaningless rules.
- **Minimum Identity**: In their `@anchors` header, they simply declare their layer membership (`layer: infra` or `layer: dao`).
- **Common Examples**:
  - `infra`: Terraform, CDK, or cloud driver configurations.
  - `dao` / `repository`: Raw database queries and data access objects.
  - `domain-types`: Pure DTOs, interfaces, and struct definitions.
  - `doc`: Manuals, guides, and architectural documentation.

> [!TIP]
> **The Golden Rule:** If the file decides *how the business behaves*, the layer is **Governed**. If the file merely executes a technical instruction without domain decisions, the layer is **Recognized**.

---

## 3. Special Governance Layers

Beyond application code, Anchors governs its own orientation artifacts through dedicated layers:

| Special Layer | Typical Pattern | Kind | Purpose in Project |
| --- | --- | --- | --- |
| **Doctrine** | `{CONCEPT,QUALITY,...}.md` | `guide` | Governance handbooks and non-negotiable project principles. |
| **Product** | `product/*.doctrine.md` | `product` | [Product Doctrine](/docs/concepts/doctrine/): Cross-cutting business rules spanning multiple screens and units. |
| **Flags** | `flags/*.flag.md` | `flag` | [Feature Flags](/docs/concepts/feature-flags/): Conditional activation scenarios without duplicating specs. |
| **Design** | `DESIGN-*.md` | `doc` | Architectural decision records (ADRs, RFCs, technical trade-offs). |
| **Plans** | `PLANNING.md`, `plans/*.plan.md` | `plan` | Milestone sequencing and seeding of future specs. |

---

## 4. Testing Regimes and Surfaces

A governed layer is verified according to the **testing regime** appropriate to its architectural nature:

| Canonical Regime | What it validates | Where it lives | Recommended Layers |
| --- | --- | --- | --- |
| **Unit** | Pure deterministic logic without I/O or network dependencies. | `src/**/unit/*_test.go` | `domain`, `service`, `usecase` |
| **Contract** | API schemas, payload serializations, and cross-service boundary models. | `tests/contracts/**/*.spec.ts` | `api`, `dto`, `client` |
| **Integration** | Communication with actual databases, caches, and queues. | `tests/integration/**/*.test.ts` | `service`, `dao`, `infra` |
| **Visual / E2E** | Screen layouts, theme transitions, accessibility, and visual regressions. | `tests/visual/**/*.spec.ts` | `ui`, `component`, `page` |

---

## 5. Architectural Boundary Enforcement

Layers exist not just to organize folders, but to enforce **direction of dependency**:

```
 [ UI / Presentation ] ────► [ Use Cases ] ────► [ Domain ]
          │                         │               ▲
          ▼                         ▼               │
    [ API Client ]           [ Services ] ──────────┘
          │                         │
          ▼                         ▼
   [ Infrastructure ] ────► [ Database DAO ]
```

The [`layer-boundary`](/docs/gates/layer-boundary/) gate evaluates your project graph against declared `depends_on` constraints:
- Domain entities cannot import UI or Infrastructure layers.
- DAOs cannot import presentation components.
- Services can only communicate with peers through approved interfaces.

---

## 6. Layer-to-Gate Applicability Matrix

| Gate Category | Governed Layers (`usecase`, `service`, `domain`) | Technical Recognized (`dao`, `infra`, `types`) | Visual / UI Layers (`component`, `page`) | Documentation (`doc`, `guide`) |
| --- | --- | --- | --- | --- |
| **The Unit & Structure**<br>([`unit-complete`](/docs/gates/unit-complete/), [`has-code`](/docs/gates/has-code/)) | **Mandatory** (`blocking: true`) | *Exempt* | **Mandatory** if business rules present | *Exempt* |
| **Traceability**<br>([`spec-feature-match`](/docs/gates/spec-feature-match/), [`code-cataloged`](/docs/gates/code-cataloged/)) | **Mandatory** (`blocking: true`) | *Exempt* | Optional | *Exempt* |
| **Boundaries & Architecture**<br>([`layer-boundary`](/docs/gates/layer-boundary/), [`circular`](/docs/gates/circular/)) | **Mandatory** (`blocking: true`) | **Mandatory** (`blocking: true`) | **Mandatory** (`blocking: true`) | *Exempt* |
| **Testing & Mutation**<br>([`mutation-score`](/docs/gates/mutation-score/), [`tests-pass`](/docs/gates/tests-pass/)) | **Mandatory** (threshold >= 80%) | Integration only | Visual baselines | *Exempt* |
| **Hygiene & Security**<br>([`no-secret-leaked`](/docs/gates/no-secret-leaked/), [`deadcode`](/docs/gates/deadcode/)) | **Mandatory** (`blocking: true`) | **Mandatory** (`blocking: true`) | **Mandatory** (`blocking: true`) | Spellcheck & format |

---

## 7. Next Steps

- Explore [The Unit](/docs/concepts/unit/) to see how specs, features, and code connect.
- Search all gates in the [Gates Catalog](/docs/gates/).
- Learn how to configure layers in [`anchors.yaml`](/docs/anchors-yaml/).
