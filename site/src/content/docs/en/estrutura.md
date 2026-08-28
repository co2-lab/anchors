---
title: Project Structure
---

> This document defines the Anchors **Project Structure pillar**. It
> presupposes the mechanism from [`CONCEPT.md`](/en/docs/conceito/) and is the
> **template** on which the other pillars operate: it defines which layers
> the project has, how they're organized, and in what order they relate.
>
> It is framework theory, already validated by a tool that implements it
> (it decomposes the project into layers — infra, architecture, interface,
> usecase, repository, service, devops — and the agent chain respects that
> order) and by a real-world proof of concept.

---

## 1. The blueprint of the house

For the framework to work, a prior agreement is needed: **which layers the
project has and how they relate to each other.** Without this agreement,
every plan and every agent invents its own layers, and the project turns
into a shapeless heap.

**Project Structure** is that agreement. If Traceability is the wiring and
Propagation is the current, Project Structure is the **blueprint of the
house** — where the rooms are and how they connect. You can wire and run
current, but it's the blueprint that says which rooms exist and how they
connect. Without a blueprint, everyone builds their own add-on.

It **suggests or dictates** the structure, and that structure must be
**respected and documented**. Like every anchor, it guides (the plan and
the agents consult it to know where things live) and is confrontable (is a
spec in the right layer? did propagation follow the right order?).

---

## 2. What Project Structure defines

- **Which layers exist** — the project's layer vocabulary (e.g.: infra,
  architecture, interface, usecase, repository, service, devops). It's not
  a universal list; each project declares its own.
- **The order / dependency between layers** — what comes before what,
  what can depend on what. This is what says "usecase comes before
  interface" or "interface depends on usecase".
- **Where each type of anchor lives** — in which layer a spec, a feature, a
  test lives; the organization convention (co-location, per-layer
  folders, etc.), expressed as **file/folder patterns**
  (`screens/**/*.spec.md` → screen spec; `repositories/**/*.ts` →
  repository).
- **The boundary rules** — what a layer can or cannot do/import. It's the
  template against which it's validated whether something respected its
  layer's limits.

> **A layer that cuts across applications: the INTERFACE.** Among the
> layers a project declares, there is almost always the **interface** — the
> entry point where a trigger (internal or external) starts the process.
> Its dialect changes with the app type (CLI: command; GUI: screen; API:
> route/handler; Job: trigger), but it's always the same conceptual layer,
> with the same spec ruler and the same position in propagation (it's the
> origin of the wave; it declares what it consumes and flows down to
> usecase/repository/service). A project can have more than one interface
> dialect (a mobile app with its own backend might have `screen` in the
> app and `handler`/`trigger` in the backend) — these are distinct tags,
> same ruler. See `SPEC_TYPES.md` §2.5.

---

## 2.1 The Structure is the virtual graph (bootstrapping the map)

The file/folder patterns per layer don't just serve organization — they
give the Structure a role that solves the **empty-map paradox**. If the
dependency map (`TRACEABILITY.md` §4) is built incrementally by each agent,
at the start it is empty: on the first file, there is no edge at all. How
does the watcher know what to do?

The answer: the Structure is the **virtual graph** — the *expected* graph,
before any material edge exists. Through the patterns, upon seeing a file,
the framework already knows, without consulting the map:

- **which layer it belongs to** (the pattern classifies:
  `screens/Login.spec.md` → screen spec);
- **what structure it should have around it** (a screen spec *should* have
  a `.tsx`, a `.feature`, a co-located `.test.tsx`);
- **which agent should handle it** (it's a spec → the propagation agent).

In other words, there are two complementary graphs:

| | source | answers |
|---|---|---|
| **virtual graph** (Structure) | file/folder patterns per layer | "which layer does this file belong to and what structure it *should* have" |
| **material graph** (the map, Traceability) | edges recorded by each agent | "what actually depends on what, right now" |

The virtual graph gives the **expected skeleton**; the material graph
gives the **realized relationships**. A new project has a complete virtual
graph (the rules already exist) and an empty material graph (nothing has
been created yet) — and it's the virtual graph that lets the watcher
classify the first file and trigger the first agent, which then starts
filling the material graph. No bootstrap paradox: the Structure knows the
shape before there's any content.

**The synthesis:** the **Structure is already the dependency diagram** —
the rules of who depends on whom, per layer, that exist without any file.
The **map** is the projection of those rules onto the real files, and it
serves to trace the **minimum path** of impact when a file changes
(`TRACEABILITY.md` §4). The Structure gives the skeleton; the map gives the
path.

## 2.2 The surface of the trio: where the pieces *must* live

The virtual graph says that an anchor *should* have its trio around it —
spec, feature, test. But **where** these pieces live is a project layout
decision, and there are two legitimate conventions:

- **Co-located** — the pieces are siblings in the same directory
  (`Login.tsx`, `Login.spec.md`, `Login.feature`, `Login.test.tsx`). This
  is the simple case: the stem `{{dir}}/{{name}}` resolves everything.
- **Centralized / by region** — the piece lives in its own tree, away from
  the anchor (e.g.: tests for all the Lambdas in
  `__tests__/unit/lambdas/`, not next to each `handler.ts`). Common in
  backends.

In both cases, **the material link continues to be by CODE**
(`TRACEABILITY.md`: "the code is the key, not the file name") — who
fulfills which requirement is the shared scenario-code, not the location.
What changes is **discovery**: when it's not co-located, the framework
needs to know *where to expect* each piece to be able to **confront its
absence**. This is where the **trio location pattern** comes in.

> **The path isn't identity — but it must respect a pattern.** A file's
> path never identifies the unit (that's the code). But it **must obey a
> declared structural pattern**: the spec for a given layer lives *here*,
> the feature *there*, the test *over there*. This pattern is the
> **validation surface** for structure gates — the ruler against which the
> question "are the pieces of this unit's trio in the places the pattern
> mandates?" is asked. A test that exists and links by code, but lives
> outside the pattern, is a **structure finding** (not a traceability gap,
> but a layout violation) — Quality decides whether to flag it.

That's why the declaration of derivatives (`Derived` in the config)
admits, in addition to the co-located stem (the default), **location
patterns per anchor-layer**: for an anchor of that layer, the test/feature
*must* match a given region template. The template can capture parts of
the anchor's path (e.g.: the **module** — the parent directory when the
file is a `handler.ts`) to compose the expected region. That way the
virtual graph knows the expected shape of the trio even when it isn't
co-located, and the gate has a surface to confront — without ever treating
the path as identity.

> This is the same classifier as the "declared layer" gate (§6): asking
> "which layer does this piece belong to?" is applying the virtual graph's
> patterns.

## 2.3 Verification regimes: each scenario confronts the surface of its regime

The "test piece" of the trio isn't a single thing — a requirement can be
verified across different **test regimes**, and each regime lives on its
own **surface**:

| Regime (canonical) | What it verifies | Typical surface |
|---|---|---|
| `unit` | a pure rule/function, no UI or I/O | unit test file |
| `integration` | behavior of a unit with its close dependencies | integration test file |
| `e2e` | user flow crossing screens/services | end-to-end script (e.g.: `.yaml`) |
| `vr` | visual appearance/regression | screenshot baseline |

A **feature mixes regimes**: its scenarios each declare, individually, in
which regime they'll be verified (a tag per scenario; a scenario can have
more than one regime). The matching gate does **not** run every scenario
against the test file — it routes **each scenario to the surface of the
regime that scenario declares**. An `e2e` scenario is confronted against
the e2e script, not against the unit test; a `vr` scenario against the
visual baseline. Confronting a visual scenario against the unit test is a
structural false-negative — the piece exists, it just lives on a different
surface.

> **The regime vocabulary belongs to the PROJECT; the canonical regime
> belongs to the framework.** A project can call its regimes whatever it
> wants (`@nivel-unit`, `@smoke`, `@wip`…). The Structure declares a
> **mapping** `project-tag → canonical-regime` (unit/integration/e2e/vr) so
> the engine can route without knowing the local naming — the same
> principle as recognized layers (`regime: declarativo`) and location
> patterns (§2.2): the mechanism is universal, the vocabulary is local. A
> new project would already name its scenarios with the canonical regimes
> and skip the mapping. A scenario tag with no mapped regime isn't
> confronted by any surface (it stays outside scrutiny — an honest
> opt-out).

Each regime resolves to a **surface** (a key in `Derived`: where that test
piece lives, co-located or via location pattern §2.2). This way the
virtual graph knows not only *that* a requirement must be tested, but *in
which regime* and *where* that verification lives — and the gate confronts
each scenario against the right surface.

---

## 3. Meta-level: the Structure defines the layers; specs live inside them

There is a subtle but essential distinction between Project Structure and
the anchors that describe each part:

- **Project Structure** is the **meta-level**: it declares that an
  architecture layer, a usecase layer, an interface layer *exist* — and
  the order among them.
- A **spec** (including the architecture spec) is **content within** one
  of the layers the Structure defined.

In other words: the architecture spec describes *that project's
architecture*; Project Structure is the template that says *that an
architecture layer exists, and where it fits in the order*. The template
precedes the content. Confusing the two is confusing the blueprint of the
house with the decoration of a room.

---

## 4. The other pillars operate on top of the Structure

Project Structure is transversal — the other pillars presuppose it:

- **Planning** (`PLANNING.md`) organizes according to it. The plan
  *can* organize itself by layer, but the layers themselves come from the
  Structure. The requirement isn't the plan's — it's the Structure's: it's
  the one that says which layers exist and their order; the plan just
  uses them.
- **Propagation** (`PROPAGATION.md`) flows in the order the Structure
  defines. When an interface spec demands a change in a usecase spec, it's
  the Structure that says that dependency exists and in which direction —
  the wave follows this blueprint.
- **Traceability** (`TRACEABILITY.md`) links layers that the Structure
  declares. The dependency map respects the boundaries of the blueprint.
- **Quality** (`QUALITY.md`) has gates that confront against the
  Structure: "does this piece respect its layer's limits?".

The Structure is the template; the other pillars draw on top of it.

---

## 5. The characteristic failure: structural anarchy

Each pillar has its own failure. The Structure's is **structural
anarchy** — without a defined and respected structure, every plan and
every agent invents its own layers and boundaries. The project loses its
common shape: pieces live in inconsistent places, layers mix, the
dependency order is violated without anyone knowing, because there was no
agreement on what the order was.

It's a distinct failure:

- it's not disconnection (**Traceability**) — the pieces can be linked,
  just linked haphazardly;
- it's not an incomplete wave (**Propagation**) — the wave can propagate,
  but through paths that shouldn't exist;
- it's not work without direction (**Planning**) — there can be a clear
  plan, but executed on a structure that changes at every step.

It's the **absence of a template**. A project without a respected
structure has no shape that a future session can recognize and continue —
each session would redesign the blueprint. Project Structure is what
guarantees that everyone builds the same house.

---

## 6. The ruler the Structure gives to conformance gates

The Structure has no gate mechanics of its own — **gate mechanics belong
to Quality** (`QUALITY.md` §2). What the Structure provides is the
**ruler**: the blueprint-conformance criterion that Quality's architecture
gate (`QUALITY.md` §3) confronts. Just as Traceability provides the ruler
of *connection*, the Structure provides the ruler of *structural
conformance* — the one who runs the confrontation is Quality.

What the Structure provides to confront:

| dimension measured | confrontation question | failure = |
|---|---|---|
| **declared layer** | does every piece belong to a layer the Structure defines? | piece outside the blueprint |
| **respected boundary** | does the piece respect what its layer can/cannot do? | layer violation |
| **honored order** | do dependencies follow the order the Structure defines? | inverted dependency |
| **honored organization** | are the anchors where convention mandates (co-location, folders)? | piece in the wrong place |

Like every Quality gate, the architecture gate runs on Propagation's
impact path (`PROPAGATION.md` §3), emits material issues when it fails,
and matures from `informative` to `blocking` (`QUALITY.md` §7) — a project
can start with a loose structure (report-only) and tighten it as it
matures. The Structure owns the *criterion*; Quality owns the *execution*.

---

## 7. Relationship with the other pillars

Project Structure is the base template:

- **Precedes Planning**: the plan can only organize by layer because the
  Structure has already defined the layers.
- **Guides Propagation**: the wave follows the dependency order the
  Structure declares.
- **Frames Traceability**: the dependency map respects the boundaries of
  the blueprint.
- **Gives Quality its ruler**: the boundary/layer gates confront against
  the Structure.

On a maturation roadmap, the Structure tends to come very early — together
with or before Planning —, because it's the template all the others
presuppose.

---

## 8. Pillar summary

- **Project Structure = the blueprint of the house.** Defines which
  layers exist, their order/dependency, where each anchor lives, and the
  boundary rules. Suggests or dictates, and must be respected and
  documented.
- **It's the meta-level.** Declares *that* layers exist; specs (including
  the architecture spec) are content *within* the layers. The template
  precedes the content.
- **The other pillars operate on it.** Planning organizes by it,
  Propagation flows in its order, Traceability links layers it declares,
  Quality confronts against it.
- **Characteristic failure = structural anarchy.** Without a respected
  template, each session redesigns the blueprint and the project loses
  its common shape.
- **Gives the ruler to conformance gates.** Declared layer, respected
  boundary, honored order, honored organization. Gate mechanics belong to
  Quality; the Structure supplies the criterion the architecture gate
  confronts.

What the pillar delivers: a project with a **recognizable shape** — a
blueprint every session respects, such that what one builds the next
understands and continues. It's the template that prevents the project
from turning into a heap of add-ons.
