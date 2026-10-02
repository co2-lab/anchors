---
title: "What is an Anchor"
description: "Understand what an Anchor is, its anatomy, identity codes, and how it anchors software development."
---

If you have ever been rock climbing, you know that climbing without fixed support points on the rock face is a recipe for disaster. An **anchor** in mountaineering is a secure, solid attachment point that:
1. **Points where to go** (the next target to reach).
2. **Holds the safety rope** (prevents a fatal fall if you slip).
3. **Marks where you've been** (leaves a verifiable route for anyone following).

In **Anchors**, we brought this exact philosophy to software engineering and AI-assisted development.

---

## 1. What is an Anchor in Software?

In software, an **Anchor** is an indivisible, verifiable, and contract-bound point of business truth. It ties together requirements, specifications, tests, and source code so that no piece can drift or lie about what it does.

Instead of keeping requirements in issue trackers or wiki pages disconnected from the codebase, an anchor lives **inside the repository alongside the code**.

```
             ┌────────────────────────────────────────────────────────┐
             │                         ANCHOR                         │
             │       Identity Code: AUTH-B01 (e.g. Session Token)     │
             └───────────────────────────┬────────────────────────────┘
                                         │
                   ┌─────────────────────┼─────────────────────┐
                   ▼                     ▼                     ▼
          ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
          │      SPEC       │   │     FEATURE     │   │      TEST       │
          │  What to build  │   │  Human Scenario │   │ Automated Proof │
          └─────────────────┘   └─────────────────┘   └─────────────────┘
                                         │
                                         ▼
                                ┌─────────────────┐
                                │      CODE       │
                                │ Implementation  │
                                └─────────────────┘
```

An anchor is not just a file; it is the **formal relationship** that connects intent to implementation.

---

## 2. The Anchor Lifecycle

Every anchor travels through a defined lifecycle, tracked deterministically by the Anchors CLI:

1. **Seed (Planning)**: The requirement is conceived during planning. It is tagged with an identity prefix (e.g., `SEED: AUTH-B01`) in `PLANNING.md` or a feature plan.
2. **Draft (Spec First)**: A specification (`*.spec.md`) is written *before* any implementation code exists. The rules and constraints are declared.
3. **Realization (The Unit)**: The feature scenario (`*.feature`), executable test (`*_test.go`, `*.test.ts`), and implementation code are created, all stamped with the anchor's identity code.
4. **Governed (Gates Active)**: Automated quality [gates](/docs/gates/) inspect the anchor on every commit and pull request. If the code drifts from the spec, the build fails.
5. **Frozen (Release Baseline)**: When a milestone is reached, anchors are frozen to record the baseline. Changes now require an intentional revision wave.

---

## 3. Why Anchors are Essential for AI Pair Programming

When you build software with Large Language Models (LLMs), AI models suffer from two primary failure modes:
- **Structural Amnesia**: The model forgets decisions made in previous prompts and introduces random directory structures, duplicated helper functions, or contradictory business logic.
- **Plausible Hallucination**: The model produces syntactically valid code that sounds convincing but quietly ignores subtle domain edge cases.

Anchors act as **guardrails that hold the rope**:
- **On the way up**: The spec gives the AI unambiguous requirements and constraints before it writes a single line of code.
- **On the way back**: Quality [gates](/docs/gates/) evaluate the generated code against the spec. If the AI hallucinates, the gate catches the divergence immediately.

---

## 4. Next Steps

- Explore [The Unit](/docs/concepts/unit/): The indivisible unit of Spec, Feature, Test, and Code.
- Learn about [The Graph & Dependency Map](/docs/concepts/graph-and-map/): How anchors connect across the system.
- Understand [Traceability & Identity Codes](/docs/concepts/traceability-and-codes/): How codes like `@rule: AUTH-B01` tie everything together.
