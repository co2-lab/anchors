---
title: Planning
---

> This document defines the Anchors **Planning pillar**. It presupposes
> the mechanism from [`CONCEPT.md`](/en/docs/conceito/) — anchor (in the
> climbing metaphor), graph, propagation, issues. Planning is the
> **origin**: the anchor that points the route, the input of the workflow,
> where the first change comes from.
>
> It is framework theory, already validated by a tool that implements it
> (the *Planner* agent produces the *Plan* that triggers the whole chain)
> and by a real-world proof of concept.

---

## 1. Reading the route before climbing

Every climber reads the wall before touching it: where the holds are,
which way to climb, in what order to drive the anchors. **Planning** is
that reading. It's the anchor that **points where to go** — the first of
the three climbing functions (`CONCEPT.md` §2).

It is the pillar of the **origin**. Everything else in the framework is
reactive: Propagation reacts to a change, Quality reacts to a
confrontation, Traceability links what already exists. But none of that
answers the question that precedes all of them: **where does the first
change come from?** The workflow begins at "a change happens" — as if it
fell from the sky. Planning is the answer: it **seeds the starting
specs**. It's the origin of the wave.

Crucial point: the plan **doesn't plan code** — it plans **specs**. Its
output is the set of specs that need to be born or changed. Code is never
a product of the plan; it's a product of the *propagation of the specs*
(§3). Planning is the origin, but operates **one level above the code**:
it seeds the base-anchors, and it's those that flow down to the code.

Without Planning, you have a system that **reacts** perfectly but doesn't
**advance with direction**. It propagates, validates, measures — but
responds to loose changes instead of to a plan. It's the difference
between climbing a route and grabbing rocks at random.

---

## 2. The plan is an anchor

The plan fulfills the anchor functions like any other:

- **Points** — says what to build or evolve, decomposed and ordered. It's
  the route.
- **Holds the rope** — it's confrontable: *"does what was done correspond
  to what was planned? was the plan fulfilled?"*. A plan that nobody
  confronts is just a wish list; a plan-anchor is the criterion against
  which progress is measured.
- **Marks** — records the intent and its order, leaving the trail of
  *why* we're climbing this way.

As an anchor, the plan lives in the graph: it **governs** the specs that
derive from it (a `governs` edge from the plan to each spec it mandates
creating), it propagates when it changes, and it goes stale when reality
diverges from it. Note that the plan governs **specs**, not code — the
plan's boundary is the spec (§3).

---

## 3. The plan is a seed, not a complete blueprint

A plan isn't a loose task list, nor an exhaustive enumeration of
everything the project needs. It **starts** the necessary specs — the
minimum to get going — and lets propagation do the rest.

What a plan captures:

- **Seeding** — the specs that need to be born or changed to get going.
  This is what links Planning to Spec: the plan says *which* starting
  specs, the spec says *what* each one is.
- **Order / phases (optional)** — the sequence in which the parts are
  built. A plan *can* organize itself by layer and phase, but it's **not a
  requirement**: the layers and their order come from **Project
  Structure** (`STRUCTURE.md`), not from the plan. The plan uses the
  layers the Structure defines; organizing by them is the planner's
  choice.
- **Progress** — the state of each part (to do / in progress / done).
  This is what makes the plan confrontable and what answers "where did we
  stop?".
- **Direction / the why** — the intent that justifies the route. This is
  what carries the compass between sessions (§6).

### Propagation on demand

The plan **doesn't need** to enumerate the entire cascade of layers. An
interface spec, once implemented, can demand a change in a usecase spec —
and that propagation **happens on the spot, by the agent**, following
Project Structure, without needing to be in the plan. A spec propagates to
other specs, to code, to features, to tests, to execution — all through
the same **Propagation** machine (`PROPAGATION.md`). It's propagation **on
demand**: the plan seeds; the entire cascade emerges from propagation
starting at the seed.

That's why Planning and Spec are adjacent but distinct: the plan seeds the
*starting* specs; the specs (via Propagation + Structure) generate the
*remaining* ones on demand. The plan doesn't need to know the entire tree
— just enough for propagation to take over.

The proof of concept materializes the seeding: the *Planner* produces a *Plan*;
the *Plan Executor* **generates the specs** in the corresponding folders
(never source code) and updates progress with each generated spec. From
there on, it's the specs that propagate.

### Two modes: seed and reinforcement (open point)

The plan has two modes of operation:

- **Seed mode (default)** — starts the minimum specs and lets propagation
  and the agents resolve the rest on demand.
- **Reinforcement mode (exception)** — sometimes the agent would
  misinterpret a spec while propagating. This is usually a sign of a
  **spec or guide failure** (something wasn't clear enough — and not
  everything can be foreseen). The workaround is for the plan to
  **pin/dictate** details of the propagation chain's behavior, to force
  the right result and avoid the wrong interpretation.

> **Open for future refinement.** Reinforcement mode is not yet refined.
> The open questions remain unanswered for now: *when to reinforce vs. let
> it propagate? is a reinforcement in the plan a smell of a weak spec/guide
> (and should it become a fix to the spec/guide, not a pin in the plan)?
> how does reinforcement interact with propagation on demand?* These are
> points to be polished once real usage reveals the pattern.

---

## 4. When the plan pays off (and where else it comes from)

Planning isn't mandatory for every change — it's an **efficiency
recommendation**, not a rule. Two things feed it besides the user's
initial intent: high-degree changes and converted issues.

### High-degree change: the plan is the more efficient path, not the only one

Changing a **high-degree** node — a guide, an architecture spec, something
that governs many files — can be done in two ways, and **both work**:

- **Direct** — the guide is edited, and Propagation does the rest: a
  global wave, all governed specs go stale, the impact is discovered *on
  the fly*, anchor by anchor (`PROPAGATION.md` §3). Reaches the right
  result, but is **expensive** — it reprocesses the tree reactively.
- **Via plan** — the plan already **brings the spec revisions** the
  change implies and **maps the affected tree ahead of time**. The
  propagation that follows is directed, not blind. Same result, **more
  efficient**.

That's why Anchors **recommends** (not requires) that high-degree changes
be born from a plan: not because the direct approach is wrong, but because
the plan avoids the cost of the blind global wave. The node's degree
(`PROPAGATION.md` §3) doesn't just decide the *size of the wave* — it
signals *when the plan pays off*. High degree → the plan pays for itself.

### The issue as origin: the user's third route

The plan can also be born from a **converted issue**. Faced with an issue,
the user has three routes (`CONCEPT.md` §5): resolve it themselves,
delegate to an agent, or **convert into a plan** when the resolution is
structured work, not a touch-up. The conversion feeds back into the flow:
the issue closes (`done/`, "converted into plan 00XX") and the plan is
born with the purpose ("opened to resolve issue 00XX") — a bidirectional
cross-reference, the debt transferred with a trail, not erased.

This closes the framework's cycle: Planning is the **origin** (seeds
specs) and also the **destination for reentry** — a large issue goes back
to the beginning by becoming a plan, by the operator's choice. The flow
isn't linear with a loop-back; it's genuinely cyclical.

---

## 5. The characteristic failure: work without direction

Each pillar has its own failure. Planning's is **work without
direction** — changes that happen, propagate, pass the gates, and yet the
project doesn't go anywhere coherent because nobody decided *what* to
build and *in what order*.

It's a failure distinct from all the others:

- it's not disconnection (**Traceability**) — the pieces can all be
  linked;
- it's not a lying anchor (anti-drift) — every anchor can be true;
- it's not low quality (**Quality**) — every gate can be passing;
- it's not an incomplete wave (**Propagation**) — everything can have
  propagated correctly.

It's the **absence of structured intent**. A project can have all the
other pillars vigorous and still be a perpetual-motion machine with no
destination: lots of movement, no advancement. Planning is the pillar that
gives **vector** to what would otherwise be mere agitation.

---

## 6. Planning carries the compass between sessions

The original pain that motivated Anchors: *future sessions need to
continue with the same rigor, following the same direction as previous
sessions.* Planning is where that **direction** lives.

A future session — another agent, another day, another human — reading
the plan knows three things that no other artifact delivers together:

- **where the project is going** (the intent, the route);
- **in what order** (the phases);
- **where we stopped** (the progress).

The "what I can't break" and the "where to resume" are born from the
plan. That's why Planning **absorbs the continuity of direction between
sessions**: it doesn't need to be a separate pillar — it's a *consequence*
of the plan being a well-made, living, confronted anchor. An updated plan
is the handoff: the next session picks up the rope right where the
previous one left it.

---

## 7. Relationship with the other pillars

- **Uses Project Structure** (`STRUCTURE.md`): the layers the plan
  organizes itself by — and the order among them — come from the
  Structure, not from the plan. The Structure is the template; the plan
  seeds specs *within* it.
- **Precedes the Spec** (`SPEC.md`): the plan says *which* starting specs
  are born; the spec says *what* each one is. Planning is the vector (the
  direction of movement); Spec is the target (the true destination). The
  plan starts the route; the spec gives each point a confrontable
  destination.
- **Gives rise to Propagation** (`PROPAGATION.md`): the first change of
  any cycle is born from a plan — which seeds specs. Planning is the
  upstream trigger of the wave; Propagation then spreads from the seeded
  specs, on demand.
- **Defines the scope of Quality** (`QUALITY.md`): the priority and
  target of each part of the plan inform what needs to be ready to be
  promoted — the plan says what's critical; Quality defends the
  threshold.
- **Is confronted like any anchor** (`CONCEPT.md` §2): "does what was
  done correspond to what was planned?" is the plan's confrontation.
  Divergence generates an issue, through the common mechanism.
- **Gives the ruler to the plan-validation gate** (`QUALITY.md` §5.1): is
  the newly-created plan coherent/feasible? — an AI review, at the flow's
  entry edge. Quality executes; Planning gives the criterion. (The
  *loop-closing* issue→plan check isn't a one-off gate — it's checked by
  the **ecosystem health validator**, `QUALITY.md` §5.2, which sweeps
  loops that never closed.)

---

## 8. Pillar summary

- **Planning = the origin.** Seeds the starting specs; is the input of
  the flow, where the first wave comes from. Nothing else in the
  framework answers "where does the change come from?".
- **The plan generates specs, not code.** Its boundary is the spec — it
  operates one level above the code. Code is a product of the
  propagation of the specs, never of the plan.
- **The plan is an anchor.** Points the route, holds the rope (is
  confrontable: "does what was done correspond to what was planned?"),
  and marks the intent. Lives in the graph, governs the specs that derive
  from it, propagates.
- **It's a seed, not a blueprint.** Starts the minimum specs; the cascade
  emerges via **propagation on demand** (a spec propagates to other
  specs/code/features/tests, following the Structure). **Seed** mode
  (default) vs. **reinforcement** (pin details in the plan when
  propagation would get it wrong — an open point for refinement).
- **Captures** seeding, order/phases (optional, coming from the
  Structure), progress, and direction.
- **Recommended (not required) for high-degree changes.** Editing a guide
  directly also propagates and works — but the plan, which brings the
  revisions and maps the tree, is *more efficient* than the blind global
  wave. Cost, not correctness.
- **Origin and destination for reentry.** Born from intent, from
  high-degree changes, and from **converted issues** (the user's 3rd
  route, `CONCEPT.md` §5): the issue closes pointing to the plan, the
  plan is born with the purpose. The flow is cyclical, not linear.
- **Characteristic failure = work without direction.** Lots of movement,
  no advancement. Gives vector to what would otherwise be mere agitation.
- **Carries the compass between sessions.** Where we're going + in what
  order + where we stopped. Absorbs the continuity of direction — the
  updated plan is the handoff.
- **Precedes the Spec.** Vector (Planning) and target (Spec) are the two
  origins: one of movement, the other of truth.

What the pillar delivers: a project that **advances with direction**, not
just reacts. Every work cycle begins with a structured intent, and any
future session knows the route, the order, and the point of resumption.
It's what keeps a project from being a healthy organism that just walks in
circles.
