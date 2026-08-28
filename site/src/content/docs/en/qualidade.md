---
title: Quality
---

> This document defines the Anchors **Quality pillar**. It presupposes the
> general mechanism established in [`CONCEPT.md`](/en/docs/conceito/) —
> anchor, graph, incremental synchrony, issues, and the living/historical
> split — and specializes it for one purpose: making sure the work isn't
> just well done, but that "well done" is **measured**.
>
> It is framework theory, already validated by a tool that implements
> Anchors and by a real-world proof of concept (incomplete and immature on
> purpose — it proves the mechanism, not completeness). Both appear only as
> instances.

---

## 1. Why quality is a pillar

Spec-first frameworks already exist. They guarantee you build the right
thing (the spec guides) and, with anti-drift, that the anchor doesn't lie
about the code. What they don't have — and what Anchors adds as a
first-class pillar — is the guarantee that the work reached a **measured
level of quality** before advancing.

No project goes to production without some level of quality or maturity.
But "quality" treated as the developer's subjective daily care is
fragile: it's a feeling, it doesn't survive between sessions, and it
vanishes when the author changes. A project like that can work and still
be ineffective — it keeps going back and forth in failures.

The Quality pillar turns that feeling into a **measured, versioned
property of the artifact**, which any future session can reassess.
That's why it's named as its own pillar: to reinforce that quality isn't
diluted into the other mechanisms nor handled carelessly. Without this
pillar well structured, the project isn't mature (see
[`CONCEPT.md` §1.1](/en/docs/conceito/)).

---

## 2. A gate is a quality anchor

A **quality gate** is an anchor (`CONCEPT.md` §2) whose target is
confronted against a **measured threshold**, not against a description.

- **Guides (way out):** the gate declares the quality target — "features
  in this file must have passing tests", "coverage ≥ 80%", "zero layer
  violations", "every screen has a `.feature`".
- **Confronts (way back):** the gate **performs the measurement** — runs
  the tests, computes coverage, runs the linter, validates structure —
  and emits an issue if the result is below the threshold.

The difference between a gate and a regular anchor is the type of
confrontation question. A spec asks *"does the code do what I
described?"* (correspondence). A gate asks *"did the target reach the
measured level?"* (threshold). Mechanically, it's the same
confrontation: the result becomes the `verdict` in the edge's stamp, and
failure becomes a `violation`-kind issue (`CONCEPT.md` §5). The Quality
pillar **doesn't invent new mechanics** — it reuses the anchor, the graph
and the issue.

### The dual output: issue and block

When a gate fails, it produces **two** things, not one:

- **An issue** — the material record of the failure, for correction
  (`CONCEPT.md` §5). Traceable, dated, append-only.
- **A block** — the barrier that prevents advancement: nothing goes up
  while the failure exists.

And there is a **deliberate asymmetry** between the two: the **issue is
non-negotiable**, the **block is negotiable**. The desync gets recorded
either way — the truth doesn't erase itself. But the block has an
opt-out (`--no-block`, §7): you can *choose to let it pass*, consciously,
without erasing the issue. The work advances, the debt stays visible.
It's the same spirit as the Spec's `@no-test` and the "done can't lie" of
issues (`CONCEPT.md` §5): **you can choose not to block, never to not
record.**

What the pillar **adds** to the common mechanism is:

1. **The link to external measurers** — how to plug test runners, linters
   and language validators in as gates (§4).
2. **The quality pipeline** — the orchestration of the set of gates (§5).
3. **Completeness** — closing all the loose ends: guaranteeing nothing
   escapes validation (§5.1).
4. **Aggregation** — quality as the set of verdicts, not a single gate
   (§6).
5. **Gate maturation** — the `informative → blocking` cycle, and the
   one-off `--no-block` opt-out — which link the pillar to the umbrella
   concept of maturity (§7).

---

## 3. "Quality gate" is a set, not a single gate

"Quality gate" is an umbrella name. In practice it's a **set of tools and
prompts that measure and evaluate** the work. **Quality emerges from the
set** — each gate measures one specific thing, and no single gate *is*
quality.

The pillar organizes itself into **gate categories**, each confronting a
different type of anchor. (Don't confuse this with the project's
*architecture layers*, which are the domain of Project Structure,
`STRUCTURE.md` — here "category" is the type of thing the gate measures.)
The reference proof of concept already exercises four categories, which
serve as a reference — not a closed list:

| gate category | what the gate measures | measurer type |
|---|---|---|
| **Architecture** | does the target respect layer/structure rules? (the Structure's ruler, `STRUCTURE.md`) | deterministic (static analyzer) |
| **Artifact** | do the anchors exist and are they complete? do they cover what the spec declares? | deterministic (spec/feature validator) |
| **Execution** | do the tests derived from the features pass? | deterministic (native test runner) |
| **Traceability** | did every feature scenario become a test in the right runner? (Traceability's ruler) | deterministic (scenario↔test coverage validator) |

### Two types of measurer: deterministic and AI judgment

What changes between gates isn't just *what* they measure — it's **how**
they measure. There are two types of measurer, and the pillar
accommodates both **with the same model** (dual output issue+block, §2;
edge strength; maturation):

- **Deterministic** — runs a command and reads the exit code. Objective,
  cheap, binary. Measures what can be *computed*: lint passed, spellcheck
  clean, coverage ≥ X%, build green, spec has the required sections, every
  scenario became a test.
- **AI judgment** — an agent confronts the target and issues the same type
  of verdict. Measures what *no script computes*: is the spec **complete**
  and coherent? is the code **readable**? does the solution respect the
  spirit of the architecture, not just the letter? This is the **AI
  review gate** — the exemplar of the qualitative side.

The difference is only *who measures*. Edge strength (blocking/
informative), dual output (issue + block) and the issue's kind are
identical in both. An AI review gate that fails generates the same issue
and the same block (with `--no-block`) as a failing lint — only the
nature of the measurer changes, not the gate's mechanics.

> **What about conformance?** The declarative regime's conformance test
> (`SPEC.md` §6 — "does the declared resource match the shape?") **isn't a
> third type of measurer**. It distributes across the two above depending
> on whether the *criterion is fixed*: comparing the Dockerfile's `EXPOSE
> 8080` against the spec, or the GSI declared in Terraform against the
> query the spec requires, is **deterministic** (shape ↔ shape,
> computable). It only becomes **judgment** when the criterion is vague
> ("do 3 replicas satisfy 'high availability'?" — unless the spec fixes
> "HA = ≥3"). Conformance is *what* the gate confronts (one shape against
> another), not *how* it measures.

### The orthogonal axis: local gate vs. gate with external state

Regardless of measurer type, there's a second axis — **where the gate
gets the truth it confronts against**:

- **Local gate** — needs only the **repository**. Lint, spellcheck,
  validate-spec, comparing Dockerfile↔spec. Cheap, reproducible, runs
  anywhere.
- **Gate with external state** — needs to look at the **world outside the
  repo**. "Does the GSI declared in Terraform actually exist in the
  cloud?" requires querying AWS. This is the *declared ↔ realized*
  confrontation — the basis of infra **drift** detection.

This axis matters because gates with external state are more expensive,
less reproducible, and can fail for reasons outside the code (the cloud
changed, the credential expired). A project typically runs the local
ones always (cheap) and the external-state ones on promotion or on demand
(§8). The axis is orthogonal to measurer type: an external-state gate can
still be deterministic or judgment — it's just that the truth it
confronts lives outside the repository.

### Catalog of common gates (reference, not a closed list)

| gate | measures | measurer | source |
|---|---|---|---|
| **lint** | code style/patterns beyond the compiler | deterministic | local |
| **spellcheck** | spelling in code/docs | deterministic | local |
| **format** | canonical formatting | deterministic | local |
| **typecheck / build** | compiles, types match | deterministic | local |
| **validate-spec** | spec present, complete, spec-first (`SPEC.md` §8) | deterministic | local |
| **validate-feature** | spec→feature coverage; valid feature | deterministic | local |
| **test / coverage** | tests pass; coverage ≥ threshold (`§4`) | deterministic | local |
| **traceability** | every requirement became an incarnation (`TRACEABILITY.md` §7) | deterministic | local |
| **architecture** | respects layer boundaries (`STRUCTURE.md`'s ruler) | deterministic | local |
| **conformance (shape)** | does the declared thing (Dockerfile/`.tf`) match the spec? | deterministic | local |
| **conformance (drift)** | does the declared resource actually exist in the environment? | deterministic | **external state** |
| **AI review** | coherence, readability, spec completeness, adherence to spirit | **AI judgment** | local |
| **AI best-practices** | follows the domain's patterns/anti-patterns | **AI judgment** | local |

A project declares the gates it wants (§4, *wiring a gate is
declarative*) — this table is a starting point, anchored in what the POC
exercises, not a fixed list. What matters is that **every type of thing
that needs validation has a gate covering it** (§5.1) — and where a
script can't reach, the AI review gate closes the gap.

---

## 4. The central bridge: feature → test → gate

The link that closes the spec-first loop with quality is the chain
**spec → feature → test**. It's the core of the pillar.

```
spec  ──generates──►  feature (scenarios)  ──implemented as──►  test (native runner)
 (the what)            (expected                   (jest / go test / pytest / …)
                        behavior, executable)              │
                                              is a quality GATE
```

The **feature isn't just a document** — it's the executable specification
of behavior, which **must be implemented as tests** in the language's
native test frameworks. Those tests **are a quality gate**. The
framework must offer a way to **wire those tests as a gate** and to
assemble a **quality pipeline** that includes them.

> **Test covers conformance, not just behavior.** The chain above is that
> of the **behavioral** regime (`SPEC.md` §6) — Gherkin scenario →
> executable test. A **declarative** spec (Terraform, Dockerfile, K8s) is
> tested by **conformance**: does the declared resource match the
> specified shape? It's a different type of test, but it enters the same
> pipeline as a gate. Every spec is testable by default; a spec the user
> doesn't want to test is marked **`@no-test`** and the gate waives it
> (see `TRACEABILITY.md`).

### The feature↔test relationship is N:M, mediated by a scenario code

The link can't be by file-name convention (fragile) nor 1:1
(insufficient). The pattern validated in practice is more robust: **every scenario
carries a stable code**, and that code is the key that ties the four
artifacts together.

```
spec.md   declares  →  CODE-S01 (a state/rule/requirement)
feature   marks     →  @CODE-S01  @nivel-integration  @P2
                                   │  (the level routes the scenario to the right runner)
                    ┌──────────────┴───────────────┐
                    ▼                              ▼
          @nivel-unit / @nivel-integration    @nivel-e2e
          unit/component test:                flow test:
          it('CODE-S01: …')  in the native runner   E2E file named CODE-S01
```

Properties of this model (distilled from the POC):

- **A scenario can declare multiple levels** (`@nivel-unit @nivel-e2e`) →
  generates a test in more than one runner. Hence the relationship is
  N:M, not 1:1.
- **The level routes the scenario to the measurer** — each test level
  corresponds to a runner: a pure rule → unit test; an isolated
  component → component test; a flow across screens → E2E test.
- **The linking key is the code**, not the file name. "Did the scenario
  become a test?" is answered by searching for the code in the declared
  level's test artifact. A scenario **without a code is invisible to the
  gate** — it's never charged for, and this is the root of orphans: the
  code is what makes traceability verifiable.
- **The code is a single source.** A single module defines the scenario
  code's grammar, imported by all validators (of all languages). It's the
  coupling point that keeps gates from different languages speaking the
  same language.

### Wiring a gate into the framework is declarative

Since the test is implemented in the native runner (which varies by
language), Anchors defines the **concept** of a test gate, but the
**concrete wiring** is declared per project. The framework offers the
mechanism; the project declares the gate:

- **which command** runs the measurer (the runner, the linter, the
  validator);
- **in what scope** it runs (which package/directory/language);
- **how the result is read** (exit code, coverage report, violation
  count);
- **what the threshold is** and whether the gate is **blocking or
  informative**.

The POC today codifies this list **imperatively** in two places (the CI
workflow and an aggregator script) — there's no single declarative
manifest. **That's exactly the gap the pillar fills:** Anchors must offer
the declarative gate manifest, so "these commands are the gates" is a
versioned anchor and not scattered code.

---

## 5. The quality pipeline

The **quality pipeline** is the orchestration of the set of gates over a
target. It transforms the spec into measured quality:

```
                          target changed
                               │
        impact path: gates to run (via Propagation, PROPAGATION §3)
                               │
      ┌───────────┬────────────┼────────────┬───────────────┐
      ▼           ▼            ▼             ▼               ▼
 architecture  artifact    execution    coverage      (project gates)
 (layers)      (spec/feat) (tests)      (scenario↔test)
      │           │            │             │               │
      └───────────┴────────────┼─────────────┴───────────────┘
                               ▼
                    aggregation = quality profile
                               │
              did any blocking gate fail? → violation issue(s)
                               │
                    (else) target passes the pipeline
```

Two points inherited from the common mechanism make the pipeline cheap
and honest:

- **Incremental** (`PROPAGATION.md` §3): the pipeline doesn't run every
  gate over the entire project on every change. Propagation's **impact
  analysis**, over Traceability's dependency map, gives the **impact
  path** — only the gates whose edges went stale enter. This is what
  makes rigor sustainable in a growing project: what's too expensive
  never gets done, and what doesn't get done turns into drift.
- **Material issues** (`CONCEPT.md` §5): a failing blocking gate isn't a
  log that vanishes — it's a `violation` issue, a file in
  `issues/todo/`, immutable, that only leaves once the desync is
  resolved (fix the target or consciously lower the threshold, which is
  updating the gate's anchor).

---

## 5.1 Closing all the loose ends: nothing without validation

Quality emerges from the *set* of gates (§3) — but this is only true if
the set **covers every loose end**. If there's an anchor, an edge, or a
requirement that **no gate confronts**, there's a hole through which
unvalidated work slips through. Closing every loose end is the goal: every
thing that exists has someone validating it.

### The characteristic failure: the coverage hole

Quality's own failure **isn't "a gate failed"** — that's the system
working (the gate did its job). The failure is the opposite: **something
exists that no gate covers** — the *open end*. And it's more dangerous
than a failure, because it's **silent**: nobody fails it, but nobody
validated it either. The work looks sound only because it wasn't looked
at.

It's a failure distinct from the others:

- it isn't Traceability's **orphan** (the piece can be perfectly
  connected and still not have a gate that measures its quality);
- it isn't Propagation's **incomplete wave** (the wave may have reached
  the piece, but there was no gate at the end to confront it).

It's **absence of coverage**. The *legitimate* door out of coverage is the
**honest opt-out** (`CONCEPT.md` §5.1) — the explicit, recorded and
justified waiver (`@no-test`, `--no-block`, informative maturation). The
coverage hole is the *illegitimate* exit — the exact opposite: implicit,
unrecorded, invisible; escaping without anyone deciding. This comparison
is what gives the meta-gate its criterion: an end without a gate is only
acceptable if there's an honest opt-out covering it; otherwise, it's a
hole.

### Artifact gates and flow gates

Closing every loose end requires distinguishing *what* is being validated.
There are two gate families:

- **Artifact gates** — validate an **anchor** (§3's catalog: lint,
  validate-spec, test, architecture…). "Is this artifact good?"
- **Flow gates** — validate an **execution or propagation step**: the
  entry, the exit, and the transversal steps of the cascade. "Was this
  *step* validated?"

The core of the cascade (`spec → code → feature → test`) is covered by
artifact gates. But the **edges** (entry, exit) and the **transversal
steps** need flow gates — otherwise they're open ends. Each execution step
has its own:

| flow gate | validates the step | measurer | ruler |
|---|---|---|---|
| **plan validation** | is the newly created plan coherent, complete, feasible? | AI review | Planning |
| **doc freshness** | does the doc reflect the source now? (synchrony triggers; the AI confronts) | AI review | — (consumption anchor) |
| **faithful-map (incremental)** | on every agent execution: are the edges the act should have created in the map? (a script that traverses the touched segment) | deterministic | Traceability |

Observed pattern: the flow gates at the **edges** (plan, doc) tend
toward **AI review** — validating intent and prose is judgment, not
computation; the **transversal** ones (map) tend toward deterministic. All
of them are Quality gates with the owning pillar's ruler, just like the
artifact ones — same mechanics (dual output, maturation), nothing new. The
audit that finds them is the meta-gate itself running over the flow
(example in `simulation/larder/cobertura.md`). Not every loose end,
however, closes via a gate at the moment of the step — some are
**systemic** and only the health validator (§5.2) catches them.

> **The opt-out doesn't need a gate.** A gate that audits old opt-outs
> (`@no-test`/`--no-block`) was considered. But the opt-out is a *user*
> decision — they validate whether it still makes sense and react when it
> isn't what was expected. A gate that enforces it would be the framework
> arbitrating a decision that belongs to the operator. The opt-out is
> already honest by being explicit, recorded and dated; the trail is
> enough for the user to review *whenever they want* — the framework
> doesn't enforce it.

### The completeness meta-gate

That's why Quality needs a gate that confronts the **completeness of the
set itself**: *"does every anchor/edge that should have a gate have
one?"*. It measures *gate* coverage, not code coverage — it hunts for open
ends before they become the path through which drift enters. When it
finds an open end, it produces the dual output like any gate (§2): an
issue (recording the end) and a block (with opt-out). Closing the end
means either wiring a gate that covers it, or explicitly marking it as
waived (`@no-test` / no requirement) — never leaving it open silently.

---

## 5.2 The ecosystem health validator

The gates — artifact and flow — cover **local** loose ends: an artifact, a
step. But there are loose ends that are **systemic** — they don't live in
a step, they live in the *integration between the pillars and the state of
the ecosystem as a whole*:

- an **issue→plan loop** that never closed (the plan concluded but the
  edge that generated the issue is still stale);
- the **overall integrity of the map** (not "this edge I just created",
  but "is the entire graph consistent, with no islands or spurious
  cycles");
- a **loose pillar** (the Structure declares layers no spec uses; there
  are specs with no gate at all; Traceability has regions with no
  identity);
- **gates that should exist and weren't declared** (the completeness
  meta-gate, §5.1, at the ecosystem level).

These loose ends don't close via a gate at the moment of a step — no
isolated step sees them. They need a **global view**. That's why Anchors
has the **ecosystem health validator**: a meta-gate elevated from the
*step* level to the *whole framework* level. It confronts *"are the
pillars sound and integrated? is the project mature?"* — sweeping the
state of all the pillars and the systemic loose ends at once.

Unlike the other gates, it:

- runs **periodically / on demand** (global view), not on every change
  (the incremental gates already cover the local level);
- measures **integration and maturity**, not a target — it's the
  **executable embodiment of maturity** (`CONCEPT.md` §1.1). "Maturity =
  presence and vigor of the pillars" was descriptive; the health
  validator is what *measures* that and turns it into a verdict;
- **presents and records, but doesn't block.** It's a special case in
  gate mechanics (§2): since it runs on demand and outside the path of a
  merge, it doesn't *lock* anything — it **reports** the health state and
  **opens issues** for the systemic loose ends it finds (reusing the
  material record). Blocking remains the job of the local gates, at the
  moment of the step. It's global diagnosis, faithful to "detect and
  present, don't arbitrate" (`CONCEPT.md` §2): the operator decides
  whether to act on the verdict.

The health validator is what **closes the last loose ends** — the ones no
local gate reaches. It's the guardian ensuring the framework applied to a
project remains a sound organism, not a set of pillars that work in
isolation but don't integrate.

### The concrete form: `doctor` / `status` in the CLI

It materializes as a **CLI command** — something like `anchors doctor`,
`status` or `validate` — that runs the ecosystem sweep and reports:
problems found, pillar states, open systemic loose ends, alerts, and the
current maturity level. It's the "X-ray" of the framework applied to the
project, invoked on demand by the operator. (There is precedent in tools
of this kind: a `/doctor` that checks the index and offers a reindex; the
health validator generalizes that to *all* pillars.) Like the watcher
(`PROPAGATION.md` §6), the command is *one* materialization — Anchors
defines the health validator in the abstract; the executor exposes it as
a CLI.

---

## 6. Aggregation: quality is the set of verdicts

The quality of a target **isn't a single number**. It's the **set of
verdicts** from its gates — a profile per gate category:

```
login.tsx  →  architecture ✔   artifact ✔   execution ✔   coverage ⚠(72%)
```

Anchors models quality as this profile, not as a 0–100 score. The
promotion decision ("can this target advance?") is **all blocking gates
passed** — informative gates enter the profile but don't lock anything.
An aggregate score, if desired, is derivable from the profile, but the
profile is the source of truth: it says *exactly what's* below, not just
*how much*.

This preserves the pillar's rule: each gate measures one specific thing;
quality is given by the set. No single gate carries the quality verdict
alone.

---

## 7. Gate maturation: `informative → blocking`

This is the link between the Quality pillar and the umbrella concept of
maturity.

A gate isn't born blocking. A real project almost never meets, on day
one, the threshold it wants to reach — imposing the gate as blocking
immediately would stop the project. That's why every gate has a
**maturation state**:

| state | behavior | when |
|---|---|---|
| **informative** (report-only) | measures and reports, but **does not block** | the project doesn't yet meet the threshold; it's on the way |
| **blocking** | measures and **prevents advancement** if below the threshold | the project reached the threshold; now it's defended |

**Promotion** from informative to blocking is an **explicit, versioned
decision** — the gate's declaration is changed in the manifest, reviewable
in a PR. It happens when the project's reality reaches the gate's
threshold ("the phase that zeroes the gaps has closed"). From then on, the
gate stops being a goal and becomes a defended frontier: nothing regresses
below it.

This cycle **is** maturity in action, at the gate level:

- An **informative** gate is a promise ("we want to get here").
- A **blocking** gate is a guarantee ("we don't go below here").
- A **mature** project is one whose critical gates have been promoted to
  blocking — the pillars are implemented and defended.
- An **immature** project still has its gates in report-only — it
  measures, it knows where it stands, but doesn't yet defend the
  threshold.

The POC demonstrates this running: its gates are born in report-only mode
and are promoted to blocking one by one, as each front matures. The POC
being "incomplete and immature" isn't a flaw to hide — it's the **living
demonstration that the pillar has maturation stages**. A finished project
would have all its gates already blocking; a POC has them in transition.

### The cold start of adoption: batch maturation

There's an asymmetry between being born with Anchors and adopting it
later, and it matters most for the **AI judgment gate** (§5.2), whose
measurer is expensive:

- A **project born with Anchors** matures for free: every artifact is
  confronted and **stamped** (`PROPAGATION.md` §3) as it's created.
  Judgment coverage grows alongside the code; there's never accumulated
  debt.
- A **project that adopts Anchors** inherits a map where *everything* is
  born without a stamp. A complete judgment audit can cost millions of AI
  tokens — infeasible all at once. This **isn't a flaw of the pillar**;
  it's the same cold start the informative→blocking maturation already
  presupposes, now on the *coverage* axis, not the *threshold* one.

The adoption strategy is to **mature coverage in batches**, and the
pillar's own mechanism sustains it:

1. **Structural pruning before measuring.** Governance redundancy
   (several guides governing the same set) and guides without governance
   inflate the work without generating value. Refining the Structure (the
   tags) cuts the total at the root, without spending on measurement.
2. **Batch by ruler, not by target.** The cost of judgment is dominated
   by *reading the ruler*, not the target. Confronting many targets of
   the same guide in one pass amortizes that fixed cost. Parallelizing,
   if applicable, is **by ruler** — never by target (which would reread
   the ruler each time).
3. **The stamp is the cursor.** The first pass doesn't need to finish all
   at once: the per-edge stamp records what's already been confronted, so
   resuming skips whatever has a valid verdict at the current rev. The
   audit becomes incremental and resumable — and, once done, only the
   *drift* (what changes afterward) returns to the queue.

The §7 honesty applies here: partial coverage must be **reported as
partial** (the ecosystem health validator, §5.1, shows what hasn't been
confronted yet), never disguised as a complete audit.

### `--no-block`: the one-off override, not the policy

Maturation is the gate's **permanent policy** — informative or blocking,
holding for the whole project until changed in a PR. **`--no-block` is
different**: it's a **one-off** override, for a single run, on a gate
that *is* blocking. "I know this would block; let it pass this time."

The difference matters:

| | maturation (informative) | `--no-block` (one-off) |
|---|---|---|
| scope | permanent, whole project | a single run |
| where it lives | the gate's declaration, versioned | the invocation, recorded |
| the issue | is born (it's a record, not a block) | is born the same |

In both cases the **issue is born** — the dual output (§2) guarantees the
truth is always recorded. What `--no-block` waives is only that time's
block, and in an **explicit and dated** way with the *why* (not silence).
Both `--no-block` and informative maturation are instances of the
**honest opt-out** (`CONCEPT.md` §5.1) — the explicit, recorded and
justified waiver, sibling of the Spec's `@no-test`: letting it pass is a
visible and auditable decision, never a hole. Permanently lowering a gate
is maturation (it changes the policy); skipping the block once is
`--no-block` (it doesn't touch the policy, and the gate stays blocking
next time).

---

## 8. Where quality is enforced

The pillar distinguishes **where** each gate runs, because this decides
what's cheap and what's definitive:

- **Locally**, before registering work (e.g.: a commit hook): cheap,
  incremental gates over what changed — fast feedback, small scope.
- **On promotion** (e.g.: before merging into the main line): the
  complete pipeline of blocking gates — the real frontier, "nothing goes
  up if it doesn't pass".

The POC proves the "nothing goes up without passing the gates" policy
with a cheap local hook (architecture over changed files) and a complete
integration pipeline that locks promotion. The distinction matters: the
local gate gives agility; the promotion gate gives the guarantee. Anchors
defines both as part of the pillar, with different granularities.

---

## 9. Pillar summary

- **Gate = quality anchor.** Guides a measured threshold, confronts by
  performing the measurement, emits a `violation` issue if below.
  Reuses anchor + graph + issue from CONCEPT.
- **Quality = a set.** "Quality gate" is an umbrella term for many gates,
  each measuring one thing; quality emerges from the set, not from a
  single gate.
- **Two types of measurer.** Deterministic (runs a command, reads the
  exit code: lint, spellcheck, validators, tests) and **AI judgment** (an
  agent measures what a script doesn't compute: coherence, readability,
  spec completeness — the AI review gate). Same mechanics, different
  measurer.
- **Dual output: issue and block.** Every failing gate generates both —
  but the **issue is non-negotiable** and the **block is negotiable**
  (`--no-block`). You can choose not to block, never to not record.
- **The feature → test bridge.** The feature is an executable
  specification; it becomes a test in the native runner; the test is a
  gate. The relationship is **N:M mediated by scenario code**, the
  verifiable key that ties spec ↔ feature ↔ test together.
- **Declarative wiring.** The framework offers the mechanism to wire
  external measurers (runners, linters, validators, agents) as gates; the
  project declares command, scope, result reading, threshold and
  strength. Replaces the POC's imperative list.
- **Closing every loose end.** The pillar's own failure is the
  **coverage hole** — something no gate covers (the open end, silent),
  distinct from "a gate failed". A meta-gate confronts completeness;
  where a script can't reach, AI review closes it.
- **Incremental pipeline.** The set of gates runs over the impact path
  (Propagation); failures become material issues. Cheap enough to be
  law.
- **Aggregation by profile.** Quality is the profile of verdicts by gate
  category, not a score. Promotion = all blocking gates passed.
- **`informative → blocking` maturation.** Every gate has a maturation
  state (permanent policy); promotion is a versioned decision.
  `--no-block` is the distinct *one-off* override — lets it pass once,
  without touching the policy or erasing the issue.

What the pillar delivers: quality that doesn't depend on memory or the
care of a single session. It's anchored in versioned gates, measured on
every change incrementally, and the project's maturity level is legible
in the state of those gates. "Well done" stops being a feeling and
becomes a verifiable property.
