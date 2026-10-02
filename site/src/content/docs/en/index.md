---
title: Anchors — Official Documentation
description: "A spec-first continuity framework for AI-assisted development. Learn core concepts, layers, and quality gates."
template: splash
hero:
  tagline: A spec-first framework — the spec comes before the code, and stays the truth after it. Keeps a project coherent over time through anchors that guide the work on the way out and confront it on the way back, and cannot lie about the code.
  actions:
    - text: Read the Concept
      link: /docs/concept/
      icon: right-arrow
      variant: primary
    - text: Explore Gates Catalog
      link: /docs/gates/
      icon: document
    - text: View on GitHub
      link: https://github.com/co2-lab/anchors
      icon: external
---

## 🌟 What is Anchors?

If you are new to **Anchors**, think of it as the **immune system and the compass of your codebase**.

When you write software with Artificial Intelligence (ChatGPT, Claude, Copilot, etc.), the AI frequently suffers from **structural amnesia**: with every new session, it forgets architectural decisions made yesterday, silently breaks business rules, and leaves behind orphaned code that no one knows how to verify.

Anchors solves this through three non-negotiable principles:
1. **Spec-First**: The specification is written *before* any implementation code, and remains the single source of truth afterwards.
2. **[The Unit](/docs/concepts/unit/)**: No business rule exists without [Spec](/docs/spec/), [Feature](/docs/concepts/unit/), [Test](/docs/quality/), and Code bound together by [identity codes](/docs/concepts/traceability-and-codes/).
3. **[Automated Gates](/docs/concepts/gates-and-verdicts/)**: Automated quality portals that hold the safety rope before git commit or during CI, preventing regressions and uncontrolled drift.

---

## 🧭 Where to Start?

### 1. If you want to UNDERSTAND how it works (Core Concepts)
Read the conceptual guides written plainly for newcomers:
- [**What is an Anchor?**](/docs/concepts/anchor/) — The climbing metaphor: how anchors point the route, hold the safety rope, and mark the trail.
- [**The Unit**](/docs/concepts/unit/) — Why code without a spec or test is immediate technical debt.
- [**Project Layers**](/docs/layers/) — The blueprint: the difference between Governed Layers (business rules) and Recognized Layers (infrastructure/types).
- [**Traceability & Identity Codes**](/docs/concepts/traceability-and-codes/) — How concise identity codes (e.g. `AUTH-B01`) connect requirements end-to-end.
- [**Gates & Verdicts**](/docs/concepts/gates-and-verdicts/) — The 4 verdict outcomes (`OK`, `FAIL`, `WARN`, `SKIP`) and why every gate starts informative.
- [**Maturity & Project Health**](/docs/concepts/maturity-and-health/) — The `anchors doctor` inspection and how to turn a fragile repo into a fortress.

### 2. If you are APPLYING Anchors to your code today (Hands-on Workflow)
Follow the practical roadmap:
- [**The CLI**](/docs/cli/) — Installing the binary and running your first checks.
- [**The Workflow**](/docs/workflow/) — Day-to-day routine from seeding a spec to safe merging.
- [**The anchors.yaml**](/docs/anchors-yaml/) — The central configuration file of your project.
- [**Gates Catalog**](/docs/gates/) — Search all 95 automated checkers ready to activate.

### 3. The Complete 7 Pillars Doctrine
For those seeking deep architectural formulation:
1. [**Project Structure**](/docs/structure/) — Blueprint layout and boundary enforcement.
2. [**Planning**](/docs/planning/) — Seeding future specs and milestone phases.
3. [**Spec**](/docs/spec/) — Executable specification discipline.
4. [**Spec Types**](/docs/spec-types/) — Interface, use case, domain service, and architectural specs.
5. [**Traceability**](/docs/traceability/) — The continuous wiring between requirements and tests.
6. [**Propagation**](/docs/propagation/) — The change wave that keeps the system coherent.
7. [**Quality**](/docs/quality/) — Measured quality theory, mutation scores, and blocking thresholds.
