---
title: Traceability
---

> This document defines the Anchors **Traceability pillar**. It
> presupposes the general mechanism from [`CONCEPT.md`](/en/docs/conceito/) —
> anchor, graph, synchrony, issues — and specializes it for one purpose:
> guaranteeing that every piece of the project has a **continuous
> identity** and is **connected**, such that the project is a single
> organism, not a pile of files that happen to live in the same folder.
>
> It is framework theory, already validated by a tool that implements
> Anchors and by a real-world proof of concept (incomplete and immature
> on purpose). Both appear only as instances.

---

## 1. The forest, not the trees

In a forest, what makes it a single living organism isn't the trees — it's
the **web of fungi** that runs under the soil and connects the roots of
all of them. That network (the mycorrhizal network) carries nutrients and
signals from one tree to another. Without it, you don't have a forest; you
have isolated trees sharing the same ground.

A software project is the same. The **anchors** (spec, feature, test,
code) are the trees. **Traceability** is the web underneath: what makes
every requirement keep the same identity across all its forms, and what
guarantees no piece becomes a disconnected island. It's the glue that
turns co-located files into a single fabric.

> While the **graph** (`CONCEPT.md` §3) is the visible structure — which
> nodes exist and how they connect —, **Traceability** is what makes each
> connection reliable, unique and without gaps. The graph is the map;
> Traceability is what guarantees the map matches reality and that no tree
> was left out of it.

### The two halves of Traceability

Traceability has **two sides**, and a project can have one without the
other:

- **Identity** — "this test fulfills *that* requirement". Links a
  requirement to its incarnations (spec → feature → test → code) via a
  stable key (§3). Answers: *the same requirement, across its forms.*
- **Dependency** — "this file depends on those; if it changes, these are
  impacted". This is the **dependency map** between files (§4). Answers:
  *what links to what.*

The two are the same glue seen from different angles: identity connects
the *forms of a requirement*; dependency connects the *files to each
other*. And both materialize in the **same artifact** — the map (§4) —
where a node can be the target of an identity edge (spec→test) and of a
dependency edge (file→file) at the same time.

> Without the dependency half, there's nothing for **Propagation** to
> operate on — it traverses the map to compute the impact path
> (`PROPAGATION.md`). Traceability **maintains** the map (the structure);
> Propagation **consumes** it (the dynamic). One artifact, two pillars.

---

## 2. Why Traceability is a pillar

Traceability is **transversal**: every other Anchors mechanism
presupposes it.

- The **graph** only exists because there's a way to link `spec ↔
  feature ↔ test ↔ code`. Without traceability, there are no edges to
  draw.
- The **Quality** pillar (`QUALITY.md` §4) depends entirely on knowing
  "did this scenario become a test?". Without the glue, the
  feature→test gate is blind.
- **Anti-drift** (`CONCEPT.md` §2) is only verifiable because you can
  trace "this identifier cited in the anchor → this real piece".
- **Issues** name the edge that broke — again, traceability.

When the glue fails, everything shuts down silently: the graph gets
holes, the gates go blind, and pieces pile up unpaired. Because it's
presupposed by all of them, if it stays implicit it rots without anyone
noticing. It's the same criterion that elevated Quality to a pillar:
**name it so it isn't handled carelessly.**

And it passes the "is it a real pillar" test: it has **its own gates** and
a **characteristic failure** that belongs to no other pillar (§6). There
is a symmetry with Quality: **Quality measures whether the work is good;
Traceability measures whether the work is connected.** Both are
"measuring" — one measures the level, the other measures whether the
pieces form a single fabric.

---

## 3. Continuous identity: the stable code

The heart of Traceability is **preserving and propagating** the **stable
identity** that each requirement receives in the spec (`SPEC.md` §4), so
it survives across all its incarnations. A requirement is born in a spec,
described as behavior in a feature, implemented as a test, and exists as
code. That's four forms — but it's **the same thing**. What ties them
together as the same thing is a **stable code** carried by all of them.

```
spec.md   declares  →  CODE-S01     (the requirement is born with identity)
feature   marks     →  @CODE-S01     (the same identity, another form)
test      titles    →  CODE-S01: …   (the same identity, executable)
code      fulfills  →  (traceable via the chain above)
```

The chain is **universal** — it holds for behavioral and declarative
specs. What varies is the *type* of test at the end: behavioral (Gherkin
scenario) for behavior specs; **conformance** (does the declared resource
match the shape?) for declarative specs (`SPEC.md` §6). "Test" in Anchors
covers both. A requirement the user doesn't want to test is marked
**`@no-test`** — the tag exempts it from the test requirement without
breaking the chain (identity remains traceable in the other forms).

Properties that make identity work (distilled from the POC):

- **The code is the key, not the file name.** Linking by file name is
  fragile (renamed, broken; two similar files cause confusion). The
  stable code survives renames and reorganizations — identity doesn't
  live in the location.
- **The relationship is N:M.** A requirement can have several incarnations
  of the same type (a scenario covered by several tests) and several
  forms (the same code appears in the spec, in the feature, and in tests
  at different levels). The code reconciles them all.
- **The code has a single-source grammar.** There is *one* place that
  defines how a code is formed and read, imported by all the validators
  of all the languages. It's the coupling point that makes pieces from
  different worlds (a markdown spec, a jest test, a go test) speak the
  **same identity language**. Without a single source, every validator
  invents its own notion of code and the glue breaks apart.

### When identity diverges: silent drift

The identity failure isn't abstract. There's a real case in the POC: the
login screen was renamed, and its stable code changed from `LGNN` to
`LOGI`. The spec, the feature and most of the tests were updated — but
**one test was left behind**, still citing `LGNN-S02`. The feature
declares `@LOGI-S02`; the gate looks for a test citing `LOGI-S02`; the
stale `LGNN-S02` **doesn't match**. Result: the feature appears
uncovered, even though the test exists — the glue broke silently because
the identity string diverged between two incarnations of the same
requirement.

That's exactly why the grammar needs to be a **single source** and identity
needs to be **the key, not the name**: the divergence of an identity string
is invisible to the naked eye, but visible to the gate — which confronts
by exact intersection of codes. Identity drift is the identity-side
analog of what anti-drift is for the anchor's content.

> **Boundary with the graph:** the graph (`CONCEPT.md`) defines *what an
> edge is* — nodes, direction, type, synchrony. Traceability defines *how
> edges come to exist and stay honest* — the stable key that lets you say
> "this node is the incarnation of that one", the uniqueness of that key,
> and the absence of orphan nodes. The graph uses Traceability to
> materialize itself; Traceability guarantees the graph has no gaps or
> duplicates.

---

## 4. The map: the material artifact

The dependency half of Traceability (§1) needs a place to live. That
place is **the map** — the material artifact, versioned in the repository,
that lists the project's files and the edges between them. Without this
file, "the project knows its dependencies" is an unsupported claim: the
information exists scattered across imports, but there's no single,
queryable, confrontable place that materializes it. The map is that
place.

The map **is the materialization of the graph** (`CONCEPT.md` §3) — not a
second artifact. CONCEPT defines the graph in the abstract (nodes, typed
edges, versioned); the map is the concrete file where that graph lives
(the `anchors.graph.yaml` from the example in `CONCEPT.md` §3), and it's
Traceability that maintains it.

The map is the **projection of the Structure onto the real files**.
Project Structure (`STRUCTURE.md` §2.1) is already the *diagram* of
dependency — the rules of who depends on whom, per layer, that exist
without any file. The map instantiates those rules on the files that
actually exist, with a version and a stamp per edge — and that's what
enables tracing the **minimum path** of impact when a file changes (the
Structure says *what type* depends on *what type*; only the map says
*which file* exactly). The Structure gives the skeleton; the map gives the
path.

> This is the gap the POC hasn't yet filled: it does identity very well
> (the scenario code crossing spec→feature→test), but **it doesn't have a
> file with the map of files and dependencies**, nor a version per file,
> nor a stamp of when each relation was validated. The map is what
> materializes dependency Traceability.

### Anatomy of the map

The map carries three pieces of information that are missing today — each
answers a question that Propagation needs to ask:

**Nodes** — one per file that participates in the graph:

| node field | answers |
|---|---|
| `id` · `kind` | which file, of what type (the literal `kind` values —
`spec`/`feature`/`test`/`code`/`doc`/`guide` — in `CONCEPT.md` §3) |
| **`rev`** (content version) | *"is this the same version as when I validated it?"* — the **file's version** |
| **`updated_at`** | *"when did this file last change?"* — the **change stamp** |

**Edges** — one per relationship, identity or dependency:

| edge field | answers |
|---|---|
| `from` · `to` · `type` | who depends on/governs whom (`CONCEPT.md` §3) |
| `origin` | how the edge entered the map (below) |
| **validation stamp** (`validated_from_rev`, `validated_to_rev`, `last_validated`, `verdict`) | *"against which versions of each end was this relation confronted, when, and with what verdict?"* |

This is the complete set of stamp fields — Traceability is its single
source. The `verdict` (`ok`/`issue`/`pending`) is written by the gate that
confronts the edge (`QUALITY.md`) and read by Propagation; the
`validated_*` fields and `last_validated` hold what and when.

The **node's version** and the **edge's stamp** are the raw material of
synchrony: Propagation (`PROPAGATION.md` §3) compares the node's current
`rev` against the `rev` the stamp holds — if it advanced, the relation is
stale. Traceability **describes and maintains** these fields in the map;
Propagation **reads** them to compute the impact path. It's the same
structure/dynamic boundary as the two halves.

### How edges enter the map

Edges are born from three origins, which **coexist** in the same material
map — with **inference as the engine**:

- **`inferred` (the engine)** — the framework reads imports/symbols from
  the code and **proposes** most edges automatically. It's what makes the
  map viable in a large project: nobody draws the graph by hand.
- **`convention`** — co-location generates the obvious edges
  (`login.tsx` ↔ `login.spec.md` ↔ `login.feature`), without needing
  inference or declaration.
- **`declared`** — the spec (or the map itself) explicitly declares edges
  that inference doesn't catch — a conceptual dependency that doesn't
  appear as an import, an architecture spec that governs N interfaces.

Inference and convention *propose*; the material map is where the edges
come to actually exist, versioned and reviewable in a PR. The map is the
source of truth; any in-memory index is a cache reconstructible from it
(`CONCEPT.md` §3).

**The reuse edge between layers (`depends-on`) is `declared`.** When a
unit consumes another layer — a screen that uses a hook/store, a usecase
that uses a repository — that dependency isn't co-location (the pieces
don't live together) and isn't always inferable (the relation is
conceptual). It's **declared in the consuming spec**, via `SPEC_TYPES.md`
§5's coded format (the Dependencies Table: `DEPn · File · Method ·
Layer`). Each row becomes a `depends-on` edge from the spec to the
referenced **file**, with the **method** as metadata (for fine-grained
impact). It's through these edges that Propagation flows down through the
data layers — without them, a change in a repository spec would never
reach the screens that depend on it (an **invisible dependency**, below).
Declaring them is what gives track to reuse propagation.

### The map also can't lie

Like every anchor, the map has its own anti-drift. An edge pointing to a
file that no longer exists is a structural lie; a real dependency the map
didn't record is an **invisible dependency** — the map-side analog of the
identity orphan (§6). The invisible dependency is dangerous because it
breaks Propagation silently: the wave doesn't traverse an edge that isn't
in the map, so a file that *should* have gone stale is never rechecked.
Keeping the map honest (re-proposing inferred edges when code changes,
removing dead edges) is an obligation of the pillar.

### The map is a self-referential anchor

A natural question: if every anchor has its dependencies recorded in the
map, what are the dependencies *of the map itself*? The answer dissolves
the regress: **the map's dependencies are all the documents it
contains.** The map doesn't need a list of dependencies about itself —
it *is* the list. It's self-referential by nature: it depends on
everything it maps, and that's already written into it. There's no "who
maps the map"; it contains itself.

### The maintenance law: whoever touches a file updates the map

The map isn't *built* by an owning agent — it's **maintained by every
agent execution**, as a transversal obligation. The rule is simple and
holds for agents of any pillar:

> **Every agent that creates, moves or removes a file updates the map in
> the same act** — recording the nodes and edges that were born, moving
> the ones that changed location, removing the ones that vanished.

That's why creating the map didn't need a special mechanism (the
supposed loose end "the framework traverses the graph but doesn't create
it"): creation is **distributed**. When the planning agent generates a
spec, it records the `plan→spec` edge; when the propagation agent creates
the code and the feature for a spec, it records those nodes and their
edges. The edge is born from the **act of creation**, not from an
external source — it's the agent that made the change declaring what it
did.

This is what keeps the map **always true**: it's updated in the same
instant the files change, never falling behind waiting for a
reconstruction. It's the continuous maintenance that gives Propagation a
map it can trust. (In the POC, this maintenance is a deterministic
side-effect of each trigger — see the concrete cascade in
`PROPAGATION.md` §6.)

---

## 5. Traceability is what routes

Stable identity doesn't just serve to *link* — it **routes**. In the POC,
each scenario declares not only its code but also its **level** and its
**priority**, and these attributes, attached to identity, tell the
framework what to do with that piece:

- the **level** (`@nivel-unit`, `@nivel-integration`, `@nivel-e2e`) routes
  the scenario to the correct test runner — it's Traceability saying
  "this identity must have an incarnation *here*";
- the **priority** governs what's required and what blocks release.

In other words: the same glue that connects also **carries the contract**
of where each piece needs to exist. It's Traceability that lets the
Quality pillar ask "where should there be a test for this identity, and
does it exist?".

---

## 6. The characteristic failure: the orphan

Each pillar has a failure that's its own. Traceability's is the
**orphan**: a piece that became disconnected from the network — a tree
outside the web of fungi. Since the pillar has two halves (§1), the orphan
has two families.

**Identity orphans** (the requirement and its forms don't connect):

- **Requirement without incarnation** — a spec or scenario nobody
  implements; the identity exists but has no realization.
- **Incarnation without requirement** — a test or code that doesn't map
  to any identity; it exists but nobody knows *what* it guarantees.
- **Missing identity** — a piece without a stable code. This is the most
  dangerous orphan on the identity side because it's **invisible to the
  gates**: without a code, no validator can cover it or link it. The POC
  documents this as the root cause of orphans — a scenario without a code
  is never charged for, so it drops off the radar without raising an
  alarm.

**Dependency orphans** (the file and the map don't match):

- **Invisible dependency** — a real dependency between files that the map
  did **not** record. It's the map-side analog of missing identity:
  without the edge, Propagation doesn't traverse that path, and a file
  that should have gone stale is never rechecked. Breaks the wave
  silently.
- **Dead edge** — an edge in the map pointing to a file that no longer
  exists. The map lies about a dependency that's no longer there.

The orphan isn't a quality failure (the test may even pass) nor an
architecture failure (the code may even respect the layers). It's a
**connection** failure — the piece isn't part of the organism, or the map
can't see it. That's why Traceability deserves its own pillar: its
failure has its own nature, and without a pillar to hunt it down, it grows
silently until the graph is more hole than fabric.

---

## 7. The Traceability gates

Like every pillar, Traceability is enforced through **gates** (which are
confrontation anchors, `QUALITY.md` §2) — only what they measure is
*connection*, not level of quality:

| gate | confrontation question | failure = |
|---|---|---|
| **identity present** | does every traceable piece have a stable code? | invisible orphan |
| **unique identity** | does each code identify only one thing (no collision)? | ambiguity |
| **requirement fulfilled** | does every declared identity have the incarnations its contract requires? (behavioral or conformance test — except what's `@no-test`) | orphan requirement |
| **anchored incarnation** | does every test/code map to a known identity? | orphan incarnation |
| **single source honored** | do all validators use the same code grammar? | broken glue |
| **faithful map** | does every edge in the map point to an existing file, and is every real dependency in the map? | invisible dependency / dead edge |

These gates follow the same model as the rest of the framework: they run
on the **impact path** (`PROPAGATION.md` §3 — Propagation traverses the
dependency map that Traceability maintains and rechecks only what the
change touched), emit **material issues** when they fail (`CONCEPT.md`
§5), and have the same cycle of
**`informative → blocking` maturation** (`QUALITY.md` §7) — a
traceability gate can be born report-only ("we have 66 known gaps") and
be promoted to blocking once coverage closes.

The traceability issue is typically of kind `violation` (a piece violates
the connection contract) and names exactly which identity became an
orphan and where — such that resolving means reconnecting (giving it a
code, implementing the missing incarnation, or anchoring the loose
incarnation).

---

## 8. Relationship with the other pillars

Traceability is the base on which the other pillars operate:

- **Enables the graph** (`CONCEPT.md`): without stable identity, there's
  no way to affirm that two nodes are incarnations of the same
  requirement. The glue is what allows the edge to be drawn.
- **Enables Quality** (`QUALITY.md`): the feature→test gate only knows
  "did this scenario become a test?" because the scenario's identity
  appears in the test. Without Traceability, Quality's gates measure
  against nothing.
- **Enables anti-drift** (`CONCEPT.md` §2): "every path/identifier cited
  in an anchor must exist" is a traceability check — following the trail
  from the citation to the real piece.
- **Respects Project Structure** (`STRUCTURE.md`): the dependency map
  links layers the Structure declares and respects the blueprint's
  boundaries. The Structure says which layers exist; Traceability links
  the pieces within and between them.

That's why, on a maturation roadmap, Traceability tends to come early: it
is the soil the other pillars take root in.

---

## 9. Pillar summary

- **Traceability = the glue.** It's the web that connects the pieces and
  makes the project a single organism, not isolated trees. The graph is
  the visible structure; Traceability is what makes it reliable, unique
  and without gaps.
- **Two halves.** *Identity* (the same requirement across its forms) and
  *dependency* (what links to what). The same glue, two angles,
  materialized in the same map.
- **Continuous identity via stable code.** Each requirement carries a
  code that crosses its incarnations (spec → feature → test → code). The
  key is the code, not the file name; the relationship is N:M; the code's
  grammar is a single source. When the string diverges (the POC's
  `LGNN`/`LOGI` drift), the glue breaks silently.
- **The map is the material artifact.** Versioned, lists files and
  edges. Each node carries its **version** (`rev`) and the **change
  stamp** (`updated_at`); each edge, the **validation stamp**.
  Traceability maintains the map; Propagation consumes it. Edges enter via
  inference (the engine) + convention + declaration.
- **Identity routes.** The code, with level and priority, carries the
  contract of where each piece needs to exist — it's what lets Quality
  know where to require a test.
- **Characteristic failure = orphan, in two families.** Identity
  (requirement/incarnation without a pair; missing identity) and
  dependency (invisible dependency; dead edge). The invisible one —
  missing code, or missing edge — is the most dangerous, because it
  escapes the gates and breaks the wave silently.
- **Connection gates.** Identity present, unique, fulfilled, anchored,
  single source honored, and faithful map. Same model: impact path,
  material issues, informative→blocking maturation.
- **Base of the other pillars.** Enables the graph, Quality and
  anti-drift. Tends to mature early — it's the soil for the rest.

What the pillar delivers: a project where nothing is lost and nothing is
an island. Every requirement can be followed from wish to realization and
back; every piece knows what it belongs to; and the organism stays one as
it grows. It's what keeps a large project from becoming a pile of files
nobody knows how to connect anymore.
