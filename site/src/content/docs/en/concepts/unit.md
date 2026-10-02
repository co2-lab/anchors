---
title: "The Unit"
description: "Understand how Spec, Feature, Test, Code, and Docs form the indivisible governance unit in Anchors."
---

If you ask ten developers how they guarantee that a feature has been delivered with high quality, most will answer: *"We achieved 100% test coverage"*.

In Anchors, we know that line coverage alone is a dangerous trap. You can have 100% line coverage executing tests against code that implements the completely wrong business logic, or empty assertions that never fail on bugs.

To solve this fundamentally, every [Governed Layer](/docs/layers/) in Anchors operates on the concept of **The Unit**.

---

## 1. What is The Unit?

**The Unit** is the fundamental, indivisible building block of development in Anchors. Instead of scattering rules in Jira, acceptance criteria in emails, tests in distant directories, and code by itself, a unit is a **cohesive bundle of artifacts that together fulfill a capability**:

```
                      ┌─────────────────────────┐
                      │          SPEC           │
                      │   (The Written Rule)    │
                      └────────────┬────────────┘
                                   │
                                   ▼
                      ┌─────────────────────────┐
                      │         FEATURE         │
                      │  (The Gherkin Scenario) │
                      └────────────┬────────────┘
                                   │
               ┌───────────────────┼───────────────────┐
               ▼                   ▼                   ▼
  ┌─────────────────────────┐ ┌───────────┐ ┌─────────────────────────┐
  │          TEST           │ │   CODE    │ │      DOCS / VISUAL      │
  │   (Executable Proof)    │ │  (Logic)  │ │ (Baselines / Handbooks) │
  │                         │ │           │ │                         │
  └─────────────────────────┘ └───────────┘ └─────────────────────────┘
```

| Piece | What it does | Typical Format | Question it answers |
| --- | --- | --- | --- |
| **1. Spec** | Defines business rules and constraints in a structured layout. | `*.spec.md` | *"What is the expected behavior and boundaries?"* |
| **2. Feature** | Translates rules into readable, executable scenarios. | `*.feature` | *"What steps does the user/system execute to exercise the rule?"* |
| **3. Test** | Automates scenarios into executable test code. | `*_test.go`, `*.test.ts` | *"Does the running code genuinely satisfy the scenario?"* |
| **4. Code** | Implements the functions and types realizing the logic. | `*.go`, `*.ts`, `*.py` | *"How does the system accomplish what was requested?"* |
| **5. Docs / Assets** | External documentation or visual regression baselines (when applicable). | `*.doc.md`, `*.png` | *"How is this consumed externally or how should the UI appear?"* |

> [!NOTE]
> **Historical note on the name "The Unit" (A Trinca):** In early versions of Anchors, this concept was informally centered on uniting all delivery artifacts (code, feature, test) revolving around the spec. The automated gate check is named [`unit-complete`](/docs/gates/unit-complete/) for backwards compatibility, but the standard and natural concept is **The Unit**, uniting all indispensable artifacts of delivery.

---

## 2. Why Must All Pieces Coexist?

Consider what happens when any of the pieces is missing:

- **Code + Test (without Spec and Feature)**: You have working code today, but no one knows *why* it behaves this way. When another developer or an AI modifies it months later, they cannot distinguish deliberate edge cases from accidental bugs.
- **Spec + Code (without Tests)**: It is merely a paper promise. Nothing proves that the code adheres to what the spec specified.
- **Spec alone (without Code or Tests)**: It is a ghost plan. Worse, without the [`unit-complete`](/docs/gates/unit-complete/) gate, an orphaned spec could pass CI silently.

The Unit turns a business capability into an **auditable, non-negotiable fact**.

---

## 3. Co-location vs. Virtual Unit

Anchors strongly favors **co-location**: placing all artifacts of the unit in the same folder:

```
src/services/billing/
├── invoice_generator.spec.md    # The business specification
├── invoice_generator.feature    # The BDD test scenarios
├── invoice_generator_test.go    # The automated test code
└── invoice_generator.go         # The production code
```

When files are co-located, any developer or AI modifying the logic immediately sees the spec, feature, and tests side by side.

When legacy frameworks force tests into a top-level `tests/` directory and specs into `docs/`, Anchors uses **virtual unit mapping** configured in `anchors.yaml`, so the graph and gates still treat them as a unified unit.

---

## 4. Key Gates Governing The Unit

- [`unit-complete`](/docs/gates/unit-complete/): Verifies that every governed spec has its matching feature, test, and code.
- [`spec-feature-match`](/docs/gates/spec-feature-match/): Ensures all rules in the spec are covered by scenarios in the feature file.
- [`feature-test-match`](/docs/gates/feature-test-match/): Verifies that every scenario in the feature file is implemented in test code.
- [`has-code`](/docs/gates/has-code/): Verifies that all rules and scenarios carry standard identity codes (e.g. `BILL-B01`).
- [`code-cataloged`](/docs/gates/code-cataloged/): Ensures every production function is registered to an identity code in the spec.
