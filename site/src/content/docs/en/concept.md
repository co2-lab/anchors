---
title: Concept
---

> **Anchors** is a continuity framework for AI-assisted development.
> This document defines the **concept**: what an anchor is, how anchors
> relate to each other, how the system knows what's up to date, and how
> desynchronization is handled. It is pure theory — independent of any
> tool.
>
> There is a tool that implements Anchors, and a real-world proof of
> concept where the concept is validated in practice. Both appear here
> only as instances — Anchors doesn't depend on either.

---

## 1. The problem

All AI-assisted work suffers from **structural amnesia**. An AI session is
stateless between invocations: the session ends, the context vanishes, and
the next session — another agent, another day, another human — starts over
without the rules, without knowing what has already been decided and why,
without knowing where the work left off. The project grows in code, but
**direction** doesn't persist.

Agent tools already know how to manage amnesia *within* a session (history
compression, summarization, intentional context reset). What's missing is
continuity *between* sessions and *across the life* of the project —
rigorously maintaining the same patterns, the same structure, and the same
direction that previous sessions established.

What persists between sessions isn't memory. It's **artifact anchored in
the repository**. Anchors is about designing these artifacts: what they
are, what role they play, how they don't rot, and how a future session is
forced to respect them.

---

## 1.1 Maturity: the umbrella concept

Delivering isn't just building the right thing — it's delivering with a
level of quality that makes it effective. A project can *work* and still be
ineffective: going back and forth in failures, never reaching a reliable
plateau. The difference between a project that works and a project that is
**ready** is **maturity**.

In Anchors, maturity is the concept at the top. A project is mature when it
has **all of its pillars implemented and vigorous**. An immature project is
one that is missing pillars, or whose pillars are loose. Maturity isn't a
number measured on an isolated artifact — it's the **presence and vigor of
the pillars** across the project as a whole.

This isn't merely descriptive: maturity is **measured** by the **ecosystem
health validator** (`QUALITY.md` §5.2) — a global-view meta-gate
(materialized as a CLI command, something like `anchors doctor`) that
sweeps the state of the pillars and the systemic loose ends, and answers
"is the framework applied to this project sound and mature?". It's what
turns "presence and vigor of the pillars" from a notion into a verdict.

The pillars are the structural mechanisms that Anchors defines. This
document specifies the mechanism common to all of them — the **anchor**
(§2), the **graph** (§3), **incremental synchrony** (§4), **desync
tracking** (§5), and the **living/historical split** (§6). Each thematic
pillar is a specialized application of that mechanism.

Pillars named so far, in route order (from origin to finish):

- **Project Structure** — [`STRUCTURE.md`](/en/docs/estrutura/). The blueprint
  of the house: defines which layers exist, their order/dependency, and
  where each anchor lives. It's the template on which the other pillars
  operate — the meta-level that declares the layers before any spec fills
  them.
- **Planning** — [`PLANNING.md`](/en/docs/planejamento/). The origin of
  *movement*: seeds the starting specs (never code), is the input of the
  flow and carries the compass between sessions (where we're going, in
  what order, where we stopped). Without it, the project reacts but
  doesn't advance with direction.
- **Spec** — [`SPEC.md`](/en/docs/spec/). The origin of *truth*: the
  base-anchor, the safepoint from which everything hangs. The spec-first
  discipline ties the other pillars together — it's the wellspring of
  Traceability, the pivot of Propagation, and the ruler of Quality.
- **Traceability** — [`TRACEABILITY.md`](/en/docs/rastreabilidade/). The glue,
  in two halves: gives each requirement a continuous identity across its
  forms (spec → feature → test → code) and maintains the dependency map
  between files, ensuring no piece becomes an island. It's the soil the
  other pillars take root in.
- **Propagation** — [`PROPAGATION.md`](/en/docs/propagacao/). The engine:
  makes a change at one point run through the organism via Traceability,
  marking what became stale, until everything is coherent again. It's the
  propagation of changes that makes development advance.
- **Quality** — [`QUALITY.md`](/en/docs/qualidade/). Without *measured*
  quality, "well done" is a feeling that doesn't survive between sessions.
  Defines gates that measure whether the work reached a threshold, and how
  those gates compose maturity.

Other pillars will be named as the framework matures. Nothing here is
carved in stone — the pillars are ideas finding their place, refined as the
concept evolves.

The relationship between the documents: CONCEPT defines what an anchor is
(in the climbing metaphor) and how it points, holds the rope and marks the
route; each pillar applies this to a type of anchor or essential dynamic —
and defines how they compose the project's maturity. **Documentation** is
an anchor, but not a pillar: it's for consumption, not structural (§2).

---

## 2. The anchor

The framework's name comes from **climbing**. To climb a wall, you drive in
anchors — fixed points that do three things: they tell you **where to go**
(the next point on the unsupported surface), **hold the rope so you don't
fall** (the safepoint), and **mark where you've been** (the route's
record). In Anchors it's the same.

> An **anchor** is any document that we write and that accompanies us in
> the project's climb — pointing to the next step, holding the rope so we
> don't fall, or marking where we've been.

This is the root definition, and it's **deliberately generous**: every
document that we generate and that is used by the application or the
framework is an anchor. Specs, features, plans, guides, docs — all are
anchors, because all accompany us in the climb. The question isn't "does
this document confront something?"; it's "is this document a fixed point
we depend on to climb safely?".

### The three functions of an anchor

Every anchor fulfills one or more of the three climbing functions:

- **Points (says where to go).** The anchor guides: the AI reads it and
  generates/changes the work following it. A plan points the route; a spec
  points to what to build.
- **Holds (the rope, the safepoint).** The anchor confronts: review,
  validation, lint and check run *against* it. If the work diverges, it
  holds — it prevents the fall, forces the correction. This is the
  **blocking** function.
- **Marks (the trail).** The anchor records where we've been — the route
  climbed, the decisions, the face of the application for those who arrive
  later. This function is **informative**: it marks without blocking.

The same anchor can fulfill several functions, and it's the **strength**
with which it holds the rope that decides its weight in the confrontation —
which already shows up in the graph's edge types (§3): `blocking` edges
hold the rope; `references` edges only mark the trail. The climbing
metaphor is the model; the graph is its materialization.

### Guide + confront: the anchor that holds the rope

The most rigorous form of an anchor — the one that **holds the rope** —
combines pointing and holding in a closed cycle: it's **executable as a
criterion**, a compass on the way out and a ruler on the way back.

```
        ANCHOR (spec / feature / plan / guide)
         ▲                              │
         │ holds the rope               │ points
         │ (review, validate,           │ (the AI generates/changes
         │  lint, check)                │  following the anchor)
         │                              ▼
        TARGET (code / other artifact) ── diverges? ──┐
                                                      │
                            ┌─────────────────────────┴───────────┐
                      fix the target                    update the anchor
                    (anchor right,                     (anchor superseded,
                     execution wrong)                   intentional change)
```

Anchors never resolves this divergence on its own. It **detects it and
presents it**. The decision between "fix the target" and "update the
anchor" is always the operator's. The framework is a desync detector and
tracker, not a judge.

### Structural anchors and consumption anchors

Anchors are distinguished by the **direction** in which they operate:

- **Structural** — govern *inward*, sustain the system. The spec governs
  the code, the plan governs the spec, the architecture governs the
  layers. These are the ones that have things hanging off them. This is
  where the pillars live.
- **Consumption (terminal)** — are generated *outward*, from the system,
  for those arriving from outside. **Documentation** is the exemplar: it
  gives a **face** to the application for those who will get to know it
  without exercising it. It's a **leaf** of the graph — nothing depends on
  it, it depends on everything; the propagation wave always *ends* at it,
  never triggers anything upstream. Its only confrontation obligation is
  **freshness**: "is it up to date with what it describes?". If what it
  describes has changed and it hasn't, it lies to the external observer.

Documentation is a first-class anchor — a citizen of the graph, it
propagates, it goes stale, it generates an issue when it goes out of
date — but **it isn't a pillar**: its importance is large for the end user
and small for the framework's mechanics, and a pillar is a measure of
mechanics. It is sustained; it doesn't sustain.

### Anti-drift is the law

An anchor that holds the rope or points the way needs a mechanism that
keeps it from **lying** about what it describes. Without this, anchors rot:
the spec describes behavior the code no longer has, the guide cites a file
that was renamed, the doc promises an API that disappeared. An anchor that
lies is worse than no anchor at all — it's a loose safepoint, one that
gives false confidence and drops you when you trust it most.

Anti-drift isn't an optional best practice; it's what keeps an anchor
trustworthy — what keeps it from turning into a dead note. It manifests in
three rules that apply to every anchor:

- **Ground it in reality.** Every path, symbol or identifier cited in an
  anchor must exist at the moment it is written. If it doesn't exist, the
  anchor is wrong.
- **Describe the WHAT, not the HOW.** Anchors capture behavior, contracts
  and invariants — never implementation details. Implementation details
  drift on the first refactor; contracts survive.
- **Preserve what is stable.** When updating an anchor, preserve stable
  identifiers, history, and the parts that are still correct. Rewriting
  wholesale destroys traceability and breaks cross-references.

---

## 3. The graph of anchors

Anchors don't live in isolation. An architecture spec governs many screen
specs; a screen is governed by its spec, by the architecture, and by a
layer guide. The relationship is **many-to-many and directed** — it's a
**graph**, not a list of pairs.

The graph is **material**: a versioned artifact in the repository,
reviewable in a Pull Request, that survives without any database. Any
in-memory index or `.db` is just a cache reconstructible from it. The graph
is the source of truth; the cache, a convenience.

### Nodes

Every file that participates in the graph is a node. Many nodes are both
anchor *and* target at the same time (a spec is the target of the guide
that governs it *and* the anchor of the code that it governs).

| field | meaning |
|---|---|
| `id` | file path (`login.spec.md`) |
| `kind` | `spec` \| `feature` \| `test` \| `code` \| `doc` \| `guide` |
| `rev` | current revision of the node's content |
| `updated_at` | when the content last changed |

### Edges

Each "governs / depends on" relationship is **its own directed edge**.
This is where the many-to-many lives: a node appears in as many edges as
needed, as origin (`from`) and as destination (`to`).

```yaml
# anchors.graph.yaml (material, versioned)
- from: architecture.spec.md
  to:   login.spec.md
  type: governs        # architecture governs the screen spec
  origin: declared

- from: login.spec.md
  to:   login.tsx
  type: specifies      # the spec describes the code
  origin: convention   # came from name co-location

- from: login.spec.md
  to:   login.feature
  type: covered-by     # the feature covers the spec

- from: login.feature
  to:   login.test.tsx
  type: tested-by       # the test exercises the feature

- from: docs/auth.md
  to:   login.tsx
  type: references     # the doc only points — doesn't confront
```

### Typed edge: the type carries the semantics

Each edge has a **semantic type**. The type carries two things at once:
the **strength** (whether divergence blocks) and the **confrontation
question** (what the verifier asks when comparing target and anchor).

| type | confrontation question | strength |
|---|---|---|
| `governs` | "does the child respect the parent's limits?" | blocking |
| `specifies` | "does the code do what the spec describes?" | blocking |
| `covered-by` | "does the feature cover all the spec's scenarios?" | blocking |
| `tested-by` | "does the test exercise the feature?" | blocking |
| `references` | (none — traceability only) | informative |

Strength is **derived from the type**, not a separate field. A blocking
edge that fails a confrontation becomes an issue (§5); an informative edge
never generates an issue — at most a warning.

The type says **what to confront**, never **who decides**. There is no
built-in hierarchy: two `governs` edges that disagree about the same
target aren't resolved by precedence — they become a conflict issue (§5).
Anchors doesn't arbitrate which anchor wins.

### Edge origin

Each edge records **how it entered the graph**, because the three sources
coexist in the same material graph:

| origin | how the edge arises |
|---|---|
| `convention` | name co-location (`login.tsx` → `login.spec.md`) |
| `declared` | the anchor explicitly declares what it governs |
| `inferred` | derived from imports/symbols by a tool |

Regardless of origin, the edge is materialized and versioned in the graph.
Inference and convention *propose* edges; the material graph is where they
come to actually exist.

---

## 4. Synchrony

The graph says *who* depends on whom. What's missing is knowing *whether*
something fell out of sync. Each node carries a **revision** (`rev`) of its
content, which advances when the file changes; and each **edge** carries a
stamp of "validated against which rev of each end". An edge is **stale**
when either end has advanced since the last confrontation — the anchor
changed, or the target changed, or the edge was never validated.

Detecting staleness without revalidating the entire project, and making a
change run through the graph until everything is coherent again, is the
system's **dynamic** — and it is its own pillar: **Propagation**, specified
in [`PROPAGATION.md`](/en/docs/propagacao/).

Here it's enough to retain the structure: **the synchrony stamp lives on
the edge** (not the node), because "being up to date" is a property of a
*relationship*, not of a file — the same spec can be in sync with the code
it describes and behind the architecture that governs it. How synchrony is
calculated, how the **impact analysis** traverses the dependency map and
produces the **impact path**, and what it means for a change to "finish
propagating" (quiescence) are the content of the Propagation pillar.

---

## 5. Issues — the record of desynchronization

When a blocking edge confrontation fails, or when the graph reveals a
desync that Anchors can't resolve on its own, the result is an **issue**:
the material record of a divergence that needs human decision.

Anchors **does not arbitrate**. Every non-mechanically-resolvable
divergence becomes an issue and stops there. The operator decides.

### Three kinds

An issue is born from one of three origins, and each kind has its own
template because each requires a different body:

| kind | trigger | what the body describes |
|---|---|---|
| `stale` | one end of the edge advanced a rev; revalidation isn't automatic | who is behind whom (`anchor_rev` vs `target_rev`) |
| `conflict` | two or more blocking anchors disagree about the same target | "anchor A says X, anchor B says Y — decide which one changes" |
| `violation` | the confrontation ran and the target **violates** the anchor | which invariant of the anchor was broken |

### Issues are files; the state is the folder

An issue is a markdown file. Its state is the **folder** where it sits.
Moving to a folder = changing state — a `git mv` visible in the PR. The
state lives on the filesystem, not in a database: it's impossible to
desync from the repository.

```
issues/
  _templates/
    stale.md
    conflict.md
    violation.md
  todo/       ← detected, nobody has picked it up
  doing/      ← someone is resolving it
  done/       ← handled (a dated fact)
```

The file name encodes date, kind and the edge, to be unique and legible:

```
issues/todo/2026-08-05-a1--stale--login.spec.md--vs--login.tsx.md
```

Example of a `stale` issue:

```markdown
---
id: 2026-08-05-a1
kind: stale
edge:
  from: login.spec.md
  to:   login.tsx
  type: specifies
detected_at: 2026-08-05
detected_by: <confrontation that opened it — lint | check | agent>
anchor_rev: 4
target_rev: 9          # the target advanced → the spec fell behind
report: <pointer to the full report, the "why">
---

## What is out of sync
login.tsx (rev 9) was changed after the spec's last validation (rev 4).
The spec describes login by email; the code now accepts a phone number.

## How to resolve (one of the two)
- [ ] **Fix the target** — revert the phone number in login.tsx (the spec rules)
- [ ] **Update the anchor** — describe the phone number in login.spec.md
      (⚠ can propagate stale upward to whatever governs the spec)

## Resolution
<!-- filled in when moved to done/: which action was taken, by whom, when -->
```

### Resolution is always one of two actions

An issue is never resolved by writing "resolved". It is resolved by
executing one of the two actions in the bidirectional flow (§2):

- **fix the target** → the target resyncs with the anchor, or
- **update the anchor** → the anchor comes to reflect the new reality (and
  this can propagate stale upward in the graph).

In both cases, a new confrontation runs, the edge's stamp reaches the
current revs, and the issue goes to `done/`.

### The issue is a user intervention point: three routes

The issue is the point where the **human intervenes** — Anchors detects
and presents, but the operator decides what to do. Faced with an issue, the
user chooses among **three routes**:

1. **Resolve it themselves** — perform one of the two actions by hand (the
   manual touch-up).
2. **Delegate to an agent** — have an agent execute the fix (the automated
   touch-up).
3. **Convert into a plan** — when the resolution is structured work, not a
   touch-up, the issue **feeds back into the development flow** by
   becoming a plan (`PLANNING.md`).

Conversion into a plan is **never automatic** — not by size, not by
heuristic. It's a choice by the operator, one of the three. When they
convert:

- the **issue closes** (goes to `done/`) with the message *"converted into
  plan 00XX"*;
- the **plan is born with the purpose** *"opened to resolve issue 00XX"* —
  a bidirectional cross-reference.

This does **not** violate "done can't lie": the issue isn't resolved in the
sense of synchronized — it's **forwarded**. The debt didn't vanish; it
changed owner (now the plan carries it and decomposes it), and the trail
shows where it went. It's the difference between *resolving* and
*transferring with a record*.

### Issues are immutable. `done` is done forever.

An issue is a **dated event**: "on this date, this desync was detected; it
was handled like this". `todo → doing → done` is a one-way path. An issue
**never reopens** — reopening would erase history and make `done` lie about
the past.

If the same desync **comes back** (`done` was false, or the target
regressed afterward), that's not the same issue again — it's a **new
event**. The next confrontation (lint, check, or agent) that detects the
conflict opens a **new issue** in `todo/`, with a new date and new id,
optionally pointing to "recurrence of <previous id>".

```
issues/done/2026-08-05-a1--stale--login.spec.md--vs--login.tsx.md   (handled on 08/05)
                    │
              (regressed / done was false)
                    ▼
   next confrontation detects it again
                    ▼
issues/todo/2026-08-09-c7--stale--login.spec.md--vs--login.tsx.md   (new fact on 08/09)
```

There is no "we're living with this" state. A desync that is consciously
not going to be resolved right now simply **stays in `todo/`, being
annoying** — and that's the correct behavior. Persistent annoyance is the
feature, not the bug: the only way to silence it is to resolve it, never
to archive it.

This gives, for free, a quality signal: since the file name encodes the
edge, scanning `done/` reveals **chronic** desyncs ("this edge has already
generated three issues") — a sign that the anchor may be poorly designed.

---

## 5.1 Honest opt-out: waiving without hiding

There are moments when a validation requirement should **not** apply to a
piece: a declarative spec the team doesn't want to test, a block that needs
to be skipped this time, a gate that isn't yet enforced with force. Anchors
allows waiving — but only in one way: the **honest opt-out**.

An honest opt-out is an **explicit, recorded, and dated** waiver of a
validation requirement. It has three invariants:

1. **Explicit** — someone wrote the mark (a tag, a flag, a policy change).
   It never emerges from silence or forgetfulness.
2. **Recorded** — it waives *the requirement* (the test, the block, the
   rigor), **never the record**. The debt remains visible. You can choose
   not to block, never to not record.
3. **Dated and justified** — carries *why* it was waived, when, and by
   whom, in a versioned artifact. It's auditable: you can list "everything
   we waived and the reason", and reassess whether it still makes sense.

This is what separates a **legitimate exit** from validation from the
**coverage hole** (`QUALITY.md` §5.1). The hole is the exact opposite:
*implicit* (nobody decided), *unrecorded* (nothing shows up), *invisible*
(you can't find it). The honest opt-out is the front door — decided, in
plain sight; the hole is the broken window.

Don't confuse this with "we're living with this" (§5): that is a *detected
desync* that keeps being annoying until resolved. The honest opt-out
comes earlier — it's waiving the *requirement* to validate that piece,
such that the desync never even gets flagged. One waives the
confrontation; the other is the confrontation pending.

The concept is transversal — each pillar instantiates it over what it
waives:

| instance | waives | where | scope |
|---|---|---|---|
| **`@no-test`** (`SPEC.md` §6) | the test requirement | tag on the spec | permanent, per anchor |
| **`--no-block`** (`QUALITY.md` §7) | the block, this time | invocation flag | one-off, single run |
| **informative maturation** (`QUALITY.md` §7) | the rigor (measures, doesn't block) | gate policy | permanent, per gate |

All three are the same concept — honest waiver — applied to different
things (the test, the block, the rigor). In all of them, the mark is
explicit, the record is born, and the why is dated.

---

## 6. Two planes: the living and the historical

Anchors maintains two separate planes, with opposite temporal natures.
They never contradict each other because each one holds a different
truth.

| plane | nature | holds the truth… | contains |
|---|---|---|---|
| **anchors + graph** | **living** — reflects the present, rewrites | …*current* about synchrony | specs, features, docs, guides, edges, revs, stamps |
| **record** | **append-only** — dated, immutable facts | …*historical* about events | resolved issues, reports, patches |

The *current* truth about whether something is in sync lives in the
**graph** (calculated from revs and stamps). The *historical* truth about
what happened lives in the **record** (dated, immutable issues and
reports).

That's why `done` doesn't need to be eternally true about the graph — it
needs to be true about *what happened on that day*. The honesty of an
issue isn't "this desync is resolved forever" (impossible to guarantee);
it's "this desync was handled on this date" (an immutable fact). Each
plane is honest on its own time axis.

Mixing the two planes — letting dated decisions rot inside a living
artifact, or rewriting history to reflect the present — is the mistake
that makes any memory system rot. Anchors keeps them rigorously separate.

---

## 7. Model summary

- **Anchor** = a document that accompanies us in the climb: it **points**
  (says where to go), **holds the rope** (confronts/blocks) or **marks**
  (leaves a trail). A generous definition — every anchor used is an
  anchor, including documentation (a consumption anchor, not a pillar).
  "Guide + confront" is the **blocking** form (the one that holds the
  rope), not the existence test. Anti-drift is law: ground it in reality,
  describe the what not the how, preserve what's stable.

- **Graph** = nodes + directed edges, many-to-many, **material and
  versioned**. The edge is **typed** (the type carries strength +
  confrontation question); strength is derived from type. No built-in
  hierarchy — the type says what to confront, never who decides.

- **Synchrony** = a validation stamp per **edge** (not per node), because
  "being up to date" is a property of the relationship. How the wave
  traverses the dependency map, performs **impact analysis** and produces
  the **impact path** (the minimum to recheck) is the **Propagation**
  pillar (`PROPAGATION.md`).

- **Issue** = the material record of desynchronization, in **three
  kinds** (stale, conflict, violation), as **files in state-folders**
  (todo/doing/done). Opened only by a blocking edge confrontation.
  **Immutable and one-way** — never reopens; recurrence is a new issue.
  Resolution = one of the two actions of the bidirectional flow.

- **Honest opt-out** = waiving a validation requirement in an **explicit,
  recorded, and dated** manner (`@no-test`, `--no-block`, maturation). It
  waives the requirement, never the record. It's the legitimate door — the
  opposite — implicit and silent waiver — is the coverage hole
  (`QUALITY.md` §5.1).

- **Two planes** = the **living** one (anchors + graph, holds the current
  truth) and the **historical** one (append-only record, holds the dated
  facts). Separated by principle; never contradict each other.

What all of this delivers is the answer to the original pain: a project
that grows without losing its rigor. Future sessions don't rebuild
context — they read anchors that guide them, confront what they do against
those anchors, and the cost of that rigor stays low because validation is
always incremental. Direction persists because it's anchored, and anchors
cannot lie.

---

## 8. Instances (non-normative)

The concept above is tool-independent. Two real instances exercise it:

- A tool that implements Anchors with AI agents: watchers detect changes,
  agents confront artifacts against guides/specs, and every change goes
  through a PR. That tool's "guide" is an anchor; its reaction chain is
  the guide+confront cycle; its dated reports and patches are the
  historical plane. What Anchors adds to it is the **material dependency
  graph** and **incremental synchrony** — today it rechecks by name
  convention, not by a graph with a freshness stamp.

- A real-world proof-of-concept workspace. Already practices anchors by
  hand (`guides/` with its own anti-drift policy, typed memories, plans
  with progress). It's where the concept is validated before being
  generalized.

These instances illustrate the concept; they don't define it. Anchors
remains true without them.
