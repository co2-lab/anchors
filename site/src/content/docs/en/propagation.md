---
title: Propagation
---

> This document defines the Anchors **Propagation pillar**. It presupposes
> the general mechanism from [`CONCEPT.md`](/en/docs/conceito/) — anchor,
> graph, issues — and the [`TRACEABILITY.md`](/en/docs/rastreabilidade/)
> pillar, on which it operates. Propagation is the engine: it's what makes
> a change at one point run through the organism, and it's what makes
> development advance.
>
> It is framework theory, already validated by a tool that implements
> Anchors and by a real-world proof of concept that exercises the concept.

---

## 1. The engine, not the map

The other pillars describe **properties** of the organism: how the pieces
connect (Traceability), whether they're good (Quality), whether they tell
the truth (anti-drift, the law from `CONCEPT.md` §2). Propagation is
different — it's the **metabolism**. It doesn't describe a state; it
produces movement. It's what makes the project leave a coherent state and
reach the next one after each change.

If **Traceability** is the web that connects the forest's roots,
**Propagation** is the signal that runs through that web — the nutrient
that flows from one root until it reaches the whole network. One is the
wiring; the other is the current running through it. One is structure; the
other is life.

The sentence that defines the pillar: **it's the propagation of changes
that makes development advance.** Without it, the graph is a static map of
dependencies. With it, the project is a system that **completes itself**
with every change — the change at one point pushes everything that
depends on it until the organism is coherent again.

---

## 2. Synchrony is diagnosis; propagation is movement

It's worth separating two things that look like one:

- **Synchrony** answers *"what's out of date right now?"*. It's a
  **state**, a snapshot. Diagnostic.
- **Propagation** answers *"when this changes, what else becomes
  outdated, and how does the wave travel through the graph?"*. It's a
  **movement**, a flow. Generative.

Synchrony is what propagation **produces and reads** along the way.
Propagation traverses the **dependency map** that Traceability maintains
in order to *compute* synchrony: at each edge it touches, it determines
whether that relation went stale. The synchrony stamp is structural (it
lives on the edge, `CONCEPT.md` §4); propagation **recomputes** it as it
passes. Stale detection isn't a separate mechanism — it's the **effect**
of the wave traversing the map.

> **The pillar's lean scope.** Propagation **points**, it doesn't resolve.
> Its job ends at "these relations became stale and need attention".
> What happens next — running the gates, confronting, generating an
> issue — isn't propagation; it's the continuation of the workflow (§6).
> Propagation is *one step* of a larger movement, not the whole movement.

---

## 3. How the wave works

### The trigger: a revision advances

Every node in the graph has a **revision** (`rev`) of its content. When a
file changes, its `rev` advances. That's the stimulus that starts the
wave.

### The stamp lives on the edge

"Being up to date" isn't a property of an isolated file — it's a property
of **a relationship**. A spec can be in sync with the code it describes
and, at the same time, behind the architecture that governs it. Two
verdicts, same spec, two edges. That's why the validation stamp hangs on
the **edge**, not the node.

The stamp's fields (`validated_from_rev`, `validated_to_rev`,
`last_validated`, `verdict`) are defined and maintained by Traceability,
which owns the map — see the full enumeration in `TRACEABILITY.md` §4.
Propagation **reads** them to compute staleness; that's the usage this
pillar describes.

### The staleness rule

An edge is **stale** when either end advanced since the last
confrontation, or when it was never validated:

```
stale(edge) :=
     node[edge.from].rev > edge.validated_from_rev    -- the anchor changed
  OR node[edge.to].rev   > edge.validated_to_rev       -- the target changed
  OR validation absent                                 -- never validated
```

```
code             spec (anchor)
 rev 7   ◄─specifies─  validated_to_rev 7   → IN SYNC  ✔
 rev 9   ◄─specifies─  validated_to_rev 7   → STALE    ⚠  (target 2 revs ahead)
```

### Impact analysis: the impact path, never the whole project

When a node advances a rev, Propagation performs **impact analysis** — it
traverses the dependency map from Traceability starting from what changed
and computes **exactly** what went stale. The result is the **impact
path**: the exact set of confrontations that need to run, and only those.
Revalidating the entire project on every change is expensive; expensive
things don't get done; things that don't get done turn into drift. Rigor
only scales if the wave touches **the minimum necessary**.

It's the same principle as an incremental build system (Make, Bazel):
don't recompile everything, recompile only what's out of date relative to
its dependencies. The impact analysis:

- **Reaches the immediate neighborhood first** — all the edges where the
  changed node is `to` (whoever governs it needs to recheck) *and* where
  it is `from` (whatever it governs may have gone stale). Both directions,
  only the immediate neighbors.
- **Propagates breadth-first, pruning what didn't go stale** — a search
  through the graph starting from the changed node, that **stops** at
  every edge that *didn't* go stale. The wave only advances ahead of the
  changes; it never sweeps the whole repository.

### The wave goes up and down the graph

Propagation follows the direction of the dependencies, and that naturally
makes it **go up and down**:

```
code (rev↑)
   │ specifies
   ▼
 spec  ── goes stale ──►  gets updated (rev↑)
   │ governs                      │
   ▼                              ▼ now the spec advanced
architecture / guide  ── may go stale ──►  rechecks
```

You touch the code → the spec that describes it goes stale. When the spec
is updated to reabsorb the change, the *spec itself* advances a rev — and
then whatever **governs the spec** (the architecture, a guide) can go
stale in turn. The wave goes up. But only the nodes actually touched enter
the queue: a branch of the graph the change didn't reach stays in sync and
is pruned.

### Wave size emerges from node degree

The same propagation engine produces waves of radically different
sizes — and the size **isn't** a second mechanism, it's a consequence of
**how many edges the changed node has**:

- Changing a screen spec (a node of **low out-degree**, that governs only
  its code, feature and test) generates a **local** wave: it reaches the
  handful of co-located derivatives and settles.
- Changing a **guide** or an architecture spec (a node of **high
  out-degree** — *fan-out* —, that *governs* dozens or hundreds of files
  via `governs` edges) generates a **global** wave: when the ruler
  changes, *all* the nodes it governs go stale at once, and the wave is a
  fan-out across the whole domain.

> The degree that matters here is **out-degree** (how many `governs` edges
> leave the node). It's different from the "node of highest convergence
> degree" that Spec is (`SPEC.md` §2) — there, many incarnations *hang*
> from a spec (in-degree). Fan-out and convergence are distinct
> dimensions of degree.

The distinction "data change" (local) vs. "criterion/ruler change"
(global) that the POC materializes is, in Anchors, **the same wave seen
at two nodes of different degree**. There's no special propagation for
guides — there's a node that governs many, and propagation does what it
always does: reach whoever went stale. That's why changing a guide is
**expensive**: the cost comes from the fan-out — a ruler change rechecks
the entire domain.
(And that's why **maturation** — deciding *when* the gate starts blocking,
`QUALITY.md` §7 — weighs so heavily for a gate that confronts a guide: it
blocks proportionally to the fan-out.)

---

## 4. The characteristic failure: the wave that dies halfway

Each pillar has a failure that's its own. Propagation's is the
**incomplete wave** — the ripple that stops before reaching everything it
should.

The result is a project in a **half-coherent** state: part reached the
new state, part stayed on the old one. And it's the most insidious
failure of all because it's **invisible point by point** — each piece,
looked at in isolation, seems fine. Only someone who sees the whole wave
notices the change didn't finish spreading.

This failure is distinct from all the others:

- it isn't a **Traceability** failure — the edges exist, the connections
  are there;
- it isn't an anti-drift failure — each anchor, in isolation, can be
  telling the truth;
- it isn't a **Quality** failure — each local gate can be passing.

It's a failure of **propagation completeness**: the change didn't finish
traversing the graph. Development stopped in the middle of a transition
without anyone noticing. Without a pillar to watch the wave as a whole,
the project accumulates unfinished transitions until it becomes a mosaic
of inconsistent states that "locally look fine".

That's why Propagation needs its own pillar: its failure is only
detectable by whoever looks at the whole movement, not the pieces.

---

## 5. Quiescence: how you know a change is finished

Propagation brings an objective criterion Anchors wouldn't have without
it: **when a change is actually complete.**

It's not when you saved the file. It's when the wave has **settled** —
when it has traversed the graph and there's no longer any stale edge that
the change left behind. That resting state is **quiescence**.

```
change → wave travels → edges go stale → get rechecked/updated
        → generate new waves → ... → no stale edges remaining
        → QUIESCENCE (the change finished propagating)
```

> **Important boundary.** Quiescence *of propagation* means only "the
> staleness wave has settled" — nothing else is out of date by
> *dependency*. This is **not** the same as "the work is done". After
> propagation settles, the flow continues: quality gates run, they can
> find divergences, generate issues, and start new changes — which
> propagate again. The *complete* quiescence *of the project* (wave
> settled **and** gates passed **and** zero open issues) is a property of
> the **workflow** (§6), not of Propagation in isolation. Propagation
> guarantees only its part: nothing stayed stale by dependency.

---

## 6. Propagation within the workflow

Propagation is **one step** of a larger movement. The workflow is the
design that links the steps — and Propagation is the first of them, not
the only one:

```
   change
      │
      ▼
 ┌──────────────┐
 │ PROPAGATION  │  ← impact analysis over Traceability's map
 │ (this pillar)│     marks what went stale — POINTS, doesn't resolve
 └──────┬───────┘
        │
        ▼
 ┌──────────────┐
 │ GATES        │  ← run, measure, confront
 │ (Quality,    │     NOT propagation — it's the continuation of the flow
 │  anti-drift) │
 └──────┬───────┘
        │
   divergence? ──► ISSUE ──► new change ──┐
        │                                 │
        └─────────────────────────────────┘
                    (and the new change CAN propagate again)
```

When a gate is called, that's **no longer propagation of a change** — it's
execution, a step of a different nature. Propagation delivers the trigger
("these edges went stale, confront them") and exits the scene; the gates
take over. If a gate finds divergence, it generates an issue; the issue
can start a new change; and that change re-enters at the top of the wave.
The cycle is the flow; Propagation is the stretch that spreads.

### Intermediate barriers: the gate between links, not just at the end

The diagram above simplifies by putting "the gates" as a single block
after Propagation. In practice, when the wave traverses a **chain** of
anchors (spec → code → feature → test), the gate can sit **between each
link**, not just at the end. A flow can require that spec→feature
propagation only **advance** to feature→test when the feature has been
confronted and passed — a **barrier** per link.

```
spec ──►│gate│──► feature ──►│gate│──► test ──►│gate│──► ...
        (advances only if passed)  (same)              (same)
```

This doesn't change what Propagation *is* — it still only points at the
stale state. What changes is the *design of the flow*: where the barriers
sit is a choice of the flow (which isn't a pillar). The POC materializes
this by putting a gate (approval/merge) between each link of the chain —
the wave from one link only triggers the next when the previous one is
accepted. Anchors recognizes that propagation can be **gated per link**,
not only at the end.

> This is the position of the barrier *within the chain* (between-links
> vs. at-the-end) — an axis orthogonal to the *local vs. promotion*
> distinction from `QUALITY.md` §8, which is about the *moment* of
> enforcement (pre-registration vs. pre-integration). One doesn't
> partition the other: you can have between-link barriers both in the
> local gate and in the promotion gate.

### Inverse propagation: the wave that generates an issue, not a mutation

The wave doesn't only run in the "anchor → target" direction (spec flows
down to code). It also runs **inverse**: when the *target* changes outside
the flow (someone edits the code directly, without touching the spec),
the `specifies` edge goes stale from the bottom end, and the wave rises up
to the spec.

But inverse propagation has its own behavior: it **does not
self-resolve**. Unlike the direct wave (where the spec *rules* and the
code follows it), here the framework doesn't know whether the divergent
code is correct (the spec became outdated) or wrong (the code violated
the spec). So the inverse wave **stops and generates a `stale` issue** —
the bottom end advanced a rev, and revalidation isn't automatic
(`CONCEPT.md` §5). It isn't a `violation` issue: inverse propagation
*stops before confronting*, so the framework doesn't yet know whether a
violation occurred — it only records that the ends went out of sync. The
decision between "fix the target" and "update the anchor" is the
operator's (`CONCEPT.md` §2, the bidirectional flow). It's a wave that
produces *recorded desync*, not *automatic mutation*. Anchors never
rewrites an anchor to legitimize a target change without human decision.

### How the wave materializes: a watcher that knows the map and the pillars' interrelation

Abstractly, the wave "traverses the graph". Concretely, this materializes
into a **watcher** — and the ideal isn't a swarm of blind watchers (one
per file type), but **one intelligent watcher** that knows three things:

- **the virtual graph** (`STRUCTURE.md` §2.1) — through file/folder
  patterns, it knows *which layer* the changed file belongs to and what
  structure it should have, even when the material map is still empty
  (the bootstrap);
- **the dependency map** (`TRACEABILITY.md` §4) — upon seeing a file
  change, it knows *what actually depends on it*, without needing
  separate patterns per type;
- **the pillars' interrelation** — it knows *which agent* to call for that
  change (spec changed → propagation agent, which creates code +
  feature; feature changed → test and doc agents).

The watcher is impact analysis embodied: it reads the map to know *what*
went stale, and the pillars' interrelation to know *whom* to dispatch. The
wave then becomes a cascade — each triggered agent creates/updates
artifacts, which the watcher notices and routes to the next agent:

```
file change
   │  watcher consults the map + the pillars' interrelation
   ▼
spec created/changed   → propagation agent → creates code + feature
   │
feature created        → test agent (creates/updates tests)
   │                    → doc agent (generates docs for devs + stakeholders — terminal link)
   ▼
(at every link) the agent UPDATES THE MAP — records created, moved, removed nodes/edges
```

Two points close the seams:

- **Every agent creates AND maps.** Propagation doesn't only *detect*
  stale state — the triggered agent **creates** the missing artifacts
  (the code, the feature, the test) and, in the same act, updates the map
  (the maintenance law, `TRACEABILITY.md` §4). This is how the graph
  *grows*, not just gets traversed.
- **The doc is the terminal link.** The documentation agent is triggered
  at the end of the cascade, generates the consumption doc for
  stakeholders (devs and users), and the wave settles there — the doc is
  a leaf of the graph, it doesn't trigger anything upstream
  (`CONCEPT.md` §2).

This is *one* materialization (the one the POC adopts); Anchors defines
the wave in the abstract, and the watcher-that-knows-the-map is the
concrete form of executing it. Another tool could trigger differently —
the concept of propagation doesn't depend on the watcher.

The **flow itself isn't a pillar** — it's a design of how things link
together, which can change or have several perspectives. The components
that operate within it (Project Structure, Planning, Spec, Traceability,
Propagation, Quality) are what constitute the pillars: without them, the
flow would just be a flow. Propagation is one of those components — the
one that gives movement to the design.

---

## 7. Relationship with the other pillars

- **Operates on Traceability** (`TRACEABILITY.md`): the wave only
  traverses what's connected. Propagation *uses* the traceable glue to
  know where to spread — without continuous identity, the wave has no
  rails. If Traceability is the wiring, Propagation is the current.
- **Follows Project Structure** (`STRUCTURE.md`): the *order* in which
  the wave rises and falls through the dependencies is the order the
  Structure declares. Traceability says which pieces link together; the
  Structure says in what order layers depend on each other — and the
  wave respects that blueprint.
- **Feeds Quality** (`QUALITY.md`): the impact path that Propagation
  computes is exactly the set of gates the quality pipeline needs to run.
  Propagation says *what* to recheck; Quality *rechecks*.
- **Feeds anti-drift** (`CONCEPT.md` §2): the wave points at which
  anchors may have started lying (the target changed, the anchor didn't
  keep up) — giving the anti-drift confrontation its impact path.
- **Generates issues via the common mechanism** (`CONCEPT.md` §5): when a
  confrontation triggered by the wave fails, the result is a material
  issue — the same record model as the rest of the framework.

Propagation is the pillar that **connects the others through time**: it's
what turns "here's a set of desirable properties" into "here's a system
that stays coherent as it changes".

---

## 8. Pillar summary

- **Propagation = the engine.** It doesn't describe a state; it produces
  movement. It's the propagation of changes that makes development
  advance.
- **Synchrony is diagnosis; propagation is movement.** Propagation
  traverses Traceability's dependency map to *compute* synchrony —
  detecting stale state is the effect of the wave passing, not a separate
  mechanism.
- **Impact analysis = impact path.** A node advances a rev → only the
  edges actually touched go stale → breadth-first search pruning what
  didn't change. Never the whole project. Incremental like a build
  system.
- **The wave rises and falls** through the graph, following the
  dependencies, but only through the branches the change reached.
- **Wave size emerges from node degree.** Changing a screen spec (low
  degree) = local wave; changing a guide/architecture (high degree,
  governs many) = global wave. Not a special mechanism — it's the same
  wave reaching whoever went stale.
- **Characteristic failure = incomplete wave.** Half-coherent state,
  invisible point by point. Only someone who sees the whole movement
  detects it.
- **Barriers and inversion.** The flow can put a gate *between each link*
  of the chain (advances only if it passes), not only at the end. And the
  inverse wave (target changes outside the flow) **stops and generates an
  issue** instead of self-mutating — the decision is human.
- **Quiescence** = the wave has settled (nothing stale by dependency).
  It's the criterion for "the change finished propagating" — but **not**
  for "the work is finished": that's the flow's, which continues through
  the gates.
- **Lean scope: points, doesn't resolve.** Delivers the impact path and
  exits the scene; the gates take over. It's *one step* of the flow, not
  the flow.
- **The flow isn't a pillar** — it's the design; the pillars are the
  components that operate within it.

What the pillar delivers: a project that **evolves like an organism**.
Every change doesn't stay stuck where it was born — it travels through
the network, makes visible everything that needs to keep up, and the
system completes itself. It's what keeps a project from becoming an
accumulation of local changes that never finish integrating into the
whole.
