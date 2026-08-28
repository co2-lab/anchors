---
title: Spec
---

> This document defines the Anchors **Spec pillar**. It presupposes the
> mechanism from [`CONCEPT.md`](/en/docs/conceito/) and connects to all the
> other pillars — because the spec is where they tie together. The Spec is
> the **base-anchor**: the origin of truth and the point of support from
> which the rest hangs.
>
> It is framework theory, already validated by a tool that implements it
> (it's a *spec-first* tool) and by a real-world proof of concept. (Don't
> confuse this document — the conceptual pillar — with a `SPEC.md` a
> concrete tool might have at its own root, which would be that tool's
> own specification.)

---

## 1. The safepoint from which everything hangs

In a climb, some anchors are passage points and others are **safepoints** —
the fixed point, well driven in, on which the entire rope depends. If it
fails, the fall is long. The **Spec** is Anchors' safepoint. It's the
anchor that **holds the rope** with the most force, and everything else is
attached to it.

While **Planning** (`PLANNING.md`) is the origin of *movement* — where the
change comes from —, the **Spec** is the origin of *truth* — the criterion
of what the thing should be. The plan is the vector; the spec is the
target. The plan says *what to do and in what order*; the spec says *what
the thing should be*. Together they are the two origins: the plan starts
the wave, the spec gives it a true destination.

---

## 2. The spec-first discipline is the pillar

The pillar isn't the spec *file*. It's the **spec-first discipline**: the
law that **everything is born from a spec**, that the spec is the source
of truth, and that nothing exists in the system without a spec that
governs it and ties it to the other pillars.

This distinction matters. A loose `.spec.md` file is just a format. What
*operates* — what makes the Spec a pillar — is the discipline:

- the spec comes **before** the code (the code fulfills the spec, not the
  other way around);
- the spec is the **source of truth** (when there's divergence, it's the
  spec that arbitrates what's correct, until a decision is made to update
  it);
- **nothing enters the system without going through a spec** that
  connects it to the network.

This law is what gives the other pillars a place to anchor to. That's why
Spec is the most **central** pillar, even though it isn't the most
dynamic one: it's the node of highest degree in the graph, the point of
convergence.

---

## 3. The spec ties the pillars together

What you said when proposing this pillar: *the specs tie propagation
together with traceability and serve as an anchor for the quality gates.*
The spec is where the pillars **meet**.

- **Ties Traceability together** (`TRACEABILITY.md`): it's in the spec
  that requirements are born with their **stable identity** (the
  scenario/state code). That identity is what traceability propagates
  through the whole chain (spec → feature → test → code). The spec is the
  wellspring of continuous identity.
- **Ties Propagation together** (`PROPAGATION.md`): the spec is the node
  through which the wave passes in both directions. A change in the code
  makes the spec stale; a change in the spec propagates to everything that
  derives from it. It's the pivot of propagation — the wave goes up and
  down from it.
- **Is the anchor for Quality gates** (`QUALITY.md`): the gates confront
  the work *against the spec*. It's the spec that declares the behavior
  the tests must cover and that the code must fulfill. Without the spec as
  ruler, the gates measure against nothing.

That's why the spec is the **tying point**: remove it, and traceability
has no wellspring, propagation has no pivot, and quality has no ruler. The
other pillars don't cease to exist — they cease to have *anywhere to
attach*.

---

## 4. What a spec is (and isn't)

The spec captures **the what**, not **the how** — behavior, contracts and
invariants, never implementation details. This is anti-drift applied to
the spec (`CONCEPT.md` §2): implementation details drift on the first
refactor; contracts survive. A spec that describes implementation becomes
a second source of truth that competes with the code and rots.

Consequences of the rule (distilled from the POC):

- **Describes behavior, not code.** No code blocks in the spec — they
  would create a target that drifts. Prose, tables, contracts,
  invariants.
- **Every requirement gains identity in the spec.** States, rules and
  requirements receive their stable code here — it's the wellspring of
  traceability.
- **Preserves what's stable when evolving.** When updating, it preserves
  identifiers, history, and the parts still correct — it doesn't rewrite
  wholesale.
- **No placeholders.** A section is either filled in or removed. A
  half-driven safepoint is worse than no safepoint at all: it gives false
  confidence.

---

## 5. Three levels: guide, template, spec

The Spec isn't a loose artifact — it's a **chain of anchors governing
anchors**, with three levels:

```
spec GUIDE      → the universal rules: how every spec behaves           [1 per project]
  └─ governs
TEMPLATE        → the skeleton of ONE spec type (per layer/file type)   [1 per type]
  └─ governs
SPEC            → the concrete instance describing a target             [N in the project]
```

The **guide** governs the templates; the templates govern the specs — the
same `governs` from the graph (`CONCEPT.md` §3), applied within the pillar
itself. This links Spec to **Project Structure** (`STRUCTURE.md`): the
*template types* correspond to the *layers the Structure declares*.
Typically one layer corresponds to one spec type and one template — but
the relationship isn't always 1:1 (a layer can yield two types, one type
can cover two files; see [`SPEC_TYPES.md`](/en/docs/tipos-de-spec/)). The
Structure says "layer X exists"; the Spec provides the template(s) for X.

The reference proof of concept proves the three levels: it has one spec guide
and two templates (screen and component); each target's `.spec.md` follows the
template of its type.

### The template = common skeleton + variation block

Surveying spec patterns by language/layer (frontend, backend, infra, CLI,
lib) reveals that templates **are not independent structures** — they are
**a common skeleton + a variation block**:

**Common skeleton** — every spec, of any type, carries:

| common block | role |
|---|---|
| **Identity** | stable code + file + type/layer (the wellspring of traceability) |
| **Purpose / responsibility** | what it does and what it does **not** do (the layer's boundary) |
| **Input contract** | what it receives (params / props / request / event / args / inputs) |
| **Output / effects contract** | what it returns or changes in the world |
| **Rules / behavior** | pre/post-conditions, invariants |
| **Errors / failures** | how and why it fails (and who signals) |
| **Dependencies** | which layers it depends on — and which it's **forbidden** to call (the boundaries come from the Structure, `STRUCTURE.md`; the spec just instantiates/narrows them for this target) |
| **Traceability** | the codes that flow down to feature/test |

**Variation block** — each type adds *one* block, dictated by its "source
of variation". What changes isn't the structure; it's the **vocabulary**
that fills input/output/contract:

| spec type | source of variation → distinguishing block |
|---|---|
| screen | state + navigation |
| component | props + variants + events |
| hook | effects (queries/mutations) |
| repository | queries / indexes / data access |
| handler / interface | route (or event) + status |
| usecase | business rules |
| entity | invariants + state machine |
| CLI command | args / flags / exit codes |
| library | exported public API |
| infra / devops | resources / secrets / permissions |

That's why the guide defines **one base template** and the types are
**thin overlays** on top of it — not N templates from scratch. The
concrete list of types and the content of each variation block lives in a
separate catalog ([`SPEC_TYPES.md`](/en/docs/tipos-de-spec/)), because it
grows per project/language and would bloat the pillar.

---

## 6. Two regimes: behavioral spec and declarative spec

Not every spec describes *behavior*. There are two **regimes**, orthogonal
to the type:

| | **behavioral** | **declarative** |
|---|---|---|
| describes | action: "given input X, Y happens" | desired state: "these resources will exist, connected like this" |
| examples | screen, hook, usecase, CLI command, lib | Dockerfile, K8s manifest, Terraform module, data schema |
| natural test type | behavioral (given/when/then scenarios) | **conformance** (does the declared resource exist and comply?) |
| failure mode | runtime, edge case, exit code | reconciliation, invalid config, drift |

The regime doesn't fork the propagation path — the chain is the same
(spec → feature → test → code, `TRACEABILITY.md`). What it changes is the
**type of test** that makes sense at the end. To that end, Anchors
**expands the notion of test**: a test isn't only the behavioral Gherkin
scenario — it includes the **conformance test** (structural validation:
does the declared resource match the specified shape?). A declarative spec
is testable — by conformance.

### Testable by default; `@no-test` as opt-out

**Every spec is testable by default**, including declarative ones. But
testing infra conformance isn't always desired — it depends on the user.
So the mechanism is an opt-out, via **tag**, the same way level and
priority are tagged. It's the Spec's instance of the **honest opt-out**
(`CONCEPT.md` §5.1): explicit, recorded and with a dated justification
(the tag carries the *why* it isn't tested).

- a spec (or a requirement within it) marked **`@no-test`** is exempted
  from the test requirement — the traceability gate "requirement
  fulfilled" (`TRACEABILITY.md`) does **not** treat it as an orphan;
- without the tag, the spec is testable and requires its test
  incarnation (behavioral or conformance, depending on the regime).

That way the Traceability chain stays **universal** — there's no special
branch for declarative specs; what exists is a different type of test and
an explicit opt-out.

The axis is orthogonal to the type: the same artifact can have parts of
both regimes (a CI pipeline is declarative in its stage definitions,
behavioral in its gates). The spec declares its regime — and the regime
determines *what type of test* the requirement demands; Propagation
follows the edges normally, with no parallel routing.

### The regime is per requirement, not per spec

The regime isn't a label on the whole spec — it's a property **of each
requirement**, declared as a **tag**, the same way as level and priority
(`TRACEABILITY.md`). The behavioral requirement carries its behavioral
regime tag; the declarative one, its own. So a **mixed** spec isn't a
third regime — it's a spec whose *requirements* have different regimes.

The canonical example is the **entity**: its *fields* are declarative (does
the shape exist and comply?) and its *invariants / state machine* are
behavioral (given X, does the transition happen?). The entity's spec
doesn't choose one regime — each of its requirements carries its own, and
each pulls the type of test that fits it (conformance for the fields,
behavioral for the invariants). It reuses the existing tag mechanism; it
doesn't invent anything.

---

## 7. The characteristic failure: what exists without a spec that governs it

The pillar's own failure is the **thing without a spec** — code, behavior
or artifact that exists in the system without a spec that governs it. It's
climbing without driving in the safepoint: the thing is there, it may
work, but **nothing sustains it** — there's no source of truth to measure
it against, no identity linking it to the network, no ruler to confront
it.

This failure is distinct:

- it isn't a **Traceability** orphan (a piece can have edges and still not
  have a spec that governs it — it's linked, but to nothing that *rules*
  it);
- it isn't low **Quality** (the thing can pass every generic gate and
  still not have its own spec saying what *it* should be);
- it isn't a lack of **Planning** (it could have been planned and still
  skip the spec step).

It's **absence of governance**. A thing without a spec is a thing the
project doesn't control — it evolves without criteria, and when it
diverges, there's nothing to judge it against. The Spec pillar guarantees
that **everything that exists is governed by something**.

---

## 8. The Spec's gates

The spec-first discipline is enforced through gates (confrontation
anchors, `QUALITY.md` §2) that measure *governance and completeness*, not
any other dimension:

| gate | confrontation question | failure = |
|---|---|---|
| **spec present** | does everything requiring governance have its spec? | ungoverned thing |
| **spec complete** | does the spec have the required sections, no placeholders? | half-driven safepoint |
| **spec-first honored** | does the spec describe the what (not the how), no code? | second source of truth |
| **declared coverage** | did everything the spec declares (states, rules) become a traceable identity? | requirement the spec promises but doesn't anchor |

These gates follow the common model: they run on Propagation's impact path
(`PROPAGATION.md` §3), emit material issues when they fail, and mature
from `informative` to `blocking`
(`QUALITY.md` §7) — a project can start with many things without a spec
(report-only gate, "we have N gaps") and promote the gate to blocking once
governance closes.

---

## 9. Relationship with the other pillars

The Spec is the center of gravity:

- **Lives within Project Structure** (`STRUCTURE.md`): the Structure
  declares *that* the layer exists; the spec (including the architecture
  spec) is the content *within* it. The template precedes the content.
- **Depends on Planning** (`PLANNING.md`): the plan decides *which*
  starting specs are born; the spec fills in *what* each one is (the layer
  order comes from the Structure, not from the plan). Planning is the
  vector, Spec is the target.
- **Wellspring of Traceability**: the requirements' stable identity is
  born here.
- **Pivot of Propagation**: the wave rises and falls starting from the
  spec.
- **Ruler of Quality**: the gates confront the work against the spec.
- **Governs Documentation**: the consumption doc (`CONCEPT.md` §2) is
  derived, among other sources, from the spec — the spec provides the
  true content the doc presents outward.

On a maturation roadmap, Spec comes right after Planning in the route
(after Structure and Planning): it's the safepoint driven in before the
rope's weight is trusted to it.

---

## 10. Pillar summary

- **Spec = the base-anchor, the safepoint.** The anchor that holds the
  rope with the most force; everything else is attached to it. Origin of
  *truth* (the target), as Planning is the origin of *movement* (the
  vector).
- **The pillar is the spec-first discipline**, not the file: everything
  is born from a spec, the spec is the source of truth, nothing exists
  without a spec that governs it.
- **Ties the pillars together.** Wellspring of Traceability (identity),
  pivot of Propagation (the wave rises and falls through it), ruler of
  Quality (the gates confront against it). It's the point of convergence.
- **Captures the what, not the how.** Behavior and contracts; no code; no
  placeholders; preserves what's stable.
- **Three levels: guide → template → spec.** The guide governs the
  templates; the templates govern the specs. Template types correspond to
  the layers the Structure declares. The template is **common skeleton +
  one variation block** per type — not N structures from scratch.
  (Type catalog in [`SPEC_TYPES.md`](/en/docs/tipos-de-spec/).)
- **Two regimes: behavioral and declarative.** The regime is **per
  requirement** (tag, like level/priority), not for the whole spec — a
  mixed spec (e.g.: entity) has requirements with different regimes. The
  regime changes how it propagates: behavioral → features/tests;
  declarative → shape conformance. Axis orthogonal to the type.
- **Characteristic failure = a thing without a spec that governs it.**
  Absence of governance — the thing exists but nothing sustains or judges
  it.
- **Governance and completeness gates.** Spec present, complete,
  spec-first honored, declared coverage. Same model: impact path, issues,
  maturation.

What the pillar delivers: a project where **everything that exists is
governed by a source of truth**, and where that source is the point where
traceability, propagation and quality meet. It's the safepoint that makes
the entire climb safe — the knot that, well driven in, holds all the
weight of the rest.
