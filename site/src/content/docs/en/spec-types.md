---
title: Spec Types Catalog
---

> A reference catalog accompanying the [`SPEC.md`](/en/docs/spec/) pillar.
> While SPEC.md defines the **model** (guide → template → spec; common
> skeleton + variation block; behavioral/declarative regime), this catalog
> surveys the **concrete types** of spec by project family and what each
> one's variation block contains.
>
> It's a **living sketch**, not a closed list. The taxonomy of types comes
> from each project's Project Structure (`STRUCTURE.md`) — this catalog
> gathers the main patterns across languages/stacks to *standardize*, not
> to catalog everything. Distilled from a survey of frontend, backend and
> infra/CLI/lib, anchored in real-world examples: a mobile app with a
> serverless backend, and a multi-stack desktop tool (Go on the backend,
> React on the interface).

---

## 1. The common skeleton (every type carries this)

Regardless of type, every spec has this skeleton. What changes between
types is only the **vocabulary** that fills input/output/contract, and the
**variation block** (§3).

| common block | role |
|---|---|
| **Identity** | stable code + file + type/layer + date. The wellspring of traceability. |
| **Purpose / responsibility** | one sentence: what it does and what it does **not** do (the layer's boundary). |
| **Input contract** | what it receives. |
| **Output / effects contract** | what it returns or changes in the world. |
| **Rules / behavior** | pre/post-conditions, invariants. |
| **Errors / failures** | how and why it fails; who signals. |
| **Dependencies** | which layers it depends on — and which it's **forbidden** to call (layer boundary). |
| **Traceability** | the codes that flow down to feature/test. |

---

## 2. The regimes (orthogonal to the type)

Each type tends toward a regime, but the axis is orthogonal — an artifact
can have parts of both (see [`SPEC.md`](/en/docs/spec/) §6).

- **Behavioral** — describes action ("given X, Y happens"); propagates to
  features → tests (Gherkin).
- **Declarative** — describes desired state ("these resources will
  exist"); propagates to shape conformance (diff / drift / structural
  validation).

---

## 2.5 The INTERFACE: the entry point (one concept, many dialects)

Before the table by family, a concept that **cuts across all of them** and
is the wellspring of almost every process: the **interface**.

**Interface is the point where a trigger — internal or external — enters
the process for interaction.** It's where behavior begins. Every
application has one, and the *role* is always the same; only the **name
and the dialect** change according to the type of application:

| Application | The interface is the… | Who/what triggers it |
|---|---|---|
| CLI | **command** | the user or another program |
| Web / Desktop / Mobile | **screen** (screen/page) | the user |
| HTTP API | **route** (handler) | an external request |
| Job / Worker / event-driven serverless | **trigger** (cron/event) | a scheduler or event (internal/external) |

`screen`, `handler`, `command`, `trigger` **are not independent types —
they are DIALECTS of the same conceptual type: the interface.** Confusing
the dialect with the concept is the mistake that makes a project treat
"screen" and "HTTP handler" as unrelated things, when they are the same
architectural thing: the entry door of the process.

### Why unifying matters

- **Same content ruler.** Every interface, in any dialect, has the same
  spec structure: **input contract** (what comes in — props+navigation in
  a screen, path/query/body/headers in a route, args+flags in a command,
  event payload in a trigger), **output contract** (what returns or
  changes), **states** (loading/error/data, or the status codes),
  **errors** (and their mapping), and **auth/access**. The concrete
  vocabulary changes; the *governed aspects* don't. That's why dialect
  guides (`SCREEN_GUIDE`, `HANDLER_GUIDE`, …) are **siblings**: same
  content, different dialects — each project has the guide for the
  dialect it uses, but all of them describe an interface.
- **Same position in propagation.** The interface is the **origin** of
  the data wave: it declares what it consumes and propagation flows down
  from it to the data layers (interface → usecase/business-logic →
  repository/service). `STRUCTURE.md` §4 talks about "an interface spec
  demands a change in a usecase spec" — it's this concept, not a
  backend-specific type. The Dependencies Table (§5) of the interface is
  what ties this chain together, whether the interface is a screen or a
  route.
- **Same regime.** Behavioral: "given input X, the process does Y and
  responds Z" — propagates naturally to feature → test, in any dialect.

### How a project materializes it

**Project Structure** (`STRUCTURE.md`) declares which interface dialect(s)
the project has, with the name the team uses. A project that is mobile
**and** serverless backend, for example, has two: `screen` (the app's
interface) and `handler`/`auth-trigger`/`job-trigger` (the backend's
interfaces — HTTP route, auth event, cron/event). They are **distinct tags in the graph**
(the dialect has its own name), but **governed by the same ruler** (the
interface's). The guide for each can exist separately (a `HANDLER_GUIDE`
in addition to `SCREEN_GUIDE`), and its content mirrors its sibling
because both describe an interface — only the input/output contract
dialect changes.

> Rule of thumb: when finding an "entry point" in a new project — a
> screen, a route, a command, a job trigger — treat it as an
> **interface**. Give it the name of the project's dialect, but write the
> spec by the interface ruler (input, output, states, errors, auth), and
> link the Dependencies Table so propagation can flow down from it.

---

## 3. Types by family

For each type: the **source of variation** (what distinguishes it), the
**variation block** (the extra sections on top of the common skeleton),
and the predominant **regime**.

### Frontend (React / React Native — anchored in a real mobile app)

| type | source of variation | variation block | regime |
|---|---|---|---|
| **screen / page** _(interface)_ | state + navigation | states (loading/error/data), data contract, data states, navigation, messages | behavioral |
| **component** | props + variants | props API, variants (per axis), visual states, events/callbacks, slots | behavioral |
| **hook** | effects | signature (in/out), effects (queries/mutations, queryKeys, invalidations), error strategy | behavioral |
| **store (global state)** | state + actions | state shape, actions, selectors, hydration/persistence, invariants | behavioral |
| **service** | remote operation | operation/endpoint, input/output, server-side effects, errors (throw) | behavioral |
| **repository** | data access | model accessed, queries/mutations (CRUD), data shape, errors (throw) | behavioral |
| **lib / pure helper** | pure transformation | signature, rule/calculation, purity (no effects) | behavioral |

> Native mobile (iOS SwiftUI/MVVM, Android Compose/MVVM): the "screen
> state" (spread across screen+hook+store in React) is concentrated in a
> **ViewModel** with a typed **UiState**. Adds a `viewmodel` type (block:
> UiState, intents/actions, effects); the View spec becomes thinner
> (bindings + visual states).

### Backend (Clean Architecture / serverless — anchored in real-world examples)

| type | source of variation | variation block | regime |
|---|---|---|---|
| **entity / domain** | invariants | fields and types, invariants, state machine, intrinsic validations | mixed (fields = declarative; invariants/transitions = behavioral) |
| **usecase / interactor** | business rules | input/output (DTO), pre-conditions, step-by-step rules, domain errors, consumed ports | behavioral |
| **repository (port+adapter)** | data access | operations (CRUD), keys/queries/indexes, return format, pagination/consistency, data errors | behavioral |
| **service** | external operation / orchestration | exposed operation, input/output, external dependencies, idempotency, retry/timeout | behavioral |
| **handler / trigger** _(interface)_ | route **or** event + status | route/method (or event type), request (path/query/body/headers), response, status codes, auth, error→status mapping | behavioral |
| **data schema** _(data interface)_ | model + access | models and fields, keys and INDEXES, relations, AUTHORIZATION (who reads/writes what), migration/version | behavioral |
| **infrastructure / resource** | resources + wiring | provisioned resources, env/secrets, permissions (IAM), stack grouping, cron/schedule | declarative |

> Real-world nuances: **"service" is an overloaded term** — classic Clean
> Architecture = domain service; in other projects = remote access (calls
> a serverless function). The project chooses the definition and the spec
> follows it. **Serverless handler** describes an *event*, not a route
> (API GW / AppSync / EventBridge / auth events have distinct shapes).
> **Repository** can be split in two (remote `repositories/` + direct
> database-access `models/`, for example). **Go** materializes the port
> (code interface) and the adapter in separate files — the spec can cover
> both.

> **The `screen` (frontend) and the `handler`/`trigger` (backend) are the
> SAME concept: the interface (§2.5)** — the entry point of the process.
> They appear in different families of this table only because they have
> different dialects (props+navigation vs. request/response/status), but
> the spec's content ruler is the same. A `HANDLER_GUIDE` mirrors the
> `SCREEN_GUIDE`; both describe an interface. In the CLI, the dialect is
> the **command** (the "CLI command" row below) — also an interface.

> **The `data schema` is GOVERNED, not declarative — and it's easy to
> confuse it with `infrastructure/resource`.** The distinction: infra
> provisions (a table exists, with such capacity); the schema DECIDES
> (which fields, which indexes, and **who can read or write each
> model**). A missing index makes a query impossible months later; an
> `allow.owner()` is a business rule about visibility. None of this is
> deducible from the code that consumes the data — that's why the
> decision needs a spec.
>
> Don't confuse it with the **DAO** (a `repository` split in two): the DAO
> *translates* table↔object and decides nothing — it's a RECOGNIZED
> layer, without a spec. The schema is what decides. Dialects:
> `amplify/data/resource.ts` (Amplify), `schema.prisma` (Prisma),
> migrations (Rails/Django), versioned `CREATE TABLE` (plain SQL), table +
> index IaC (DynamoDB/CDK).

### Infra / CLI / Lib

| type | source of variation | variation block | regime |
|---|---|---|---|
| **Terraform / IaC module** | inputs → resources | variables (name/type/default/required), outputs, resources created, providers, preconditions, cost/blast-radius | declarative |
| **K8s / Helm manifest** | workload + network | workload (type, replicas), resources (cpu/mem), network (service/ingress), config/secrets, health probes, storage, policies | declarative |
| **Dockerfile / image** | execution environment | base image, included artifacts, entrypoint/cmd, ports, env vars, volumes, user/permissions | declarative |
| **CI/CD pipeline** | stages + gates | triggers, stages/jobs and order, gates (lint/tests/approval), inputs/secrets, artifacts, deploy targets, rollback | mixed (stages declarative, gates behavioral) |
| **deploy config** | target + services | environment, services/replicas, dependencies, env/secrets, ports/volumes, strategy (rolling/blue-green), healthcheck | declarative |
| **CLI command** _(interface)_ | invocation signature | positional args, flags/options (default, env equivalent, mutual exclusion), behavior, stdin, stdout (format), exit codes, preconditions | behavioral |
| **library / package** | public API | exported surface (functions/types), per-symbol contract (signature, errors, pre/post), invariants (thread-safety, complexity), SemVer/deprecations | behavioral |

> A **declarative** artifact describes shape/state (inputs → guaranteed
> resources, no temporal sequence) and fits poorly with `.feature`
> scenarios with temporal steps — it propagates to structural conformance
> gates, not Gherkin tests. A **behavioral** artifact fits naturally with
> given/when/then.

---

## 4. How to use this catalog

- The project's **Project Structure** (`STRUCTURE.md`) declares *which*
  of these layers/types exist in that project.
- For each declared type, the project has **one template** = common
  skeleton (§1) + the variation block from the corresponding row (§3).
- The project's **spec guide** defines the universal rules and points to
  the templates.
- Each concrete target's `.spec.md` follows the template of its type, and
  declares its **regime** (§2), which determines the *type of test* the
  requirement demands (behavioral or conformance). Testable by default;
  `@no-test` waives it. Propagation follows the edges normally — the
  regime doesn't route the wave.

This catalog is a starting point — each project adjusts, removes and adds
types according to its reality. What Anchors standardizes is the **model**
(skeleton + variation + regime), not the list.

---

## 5. Declaring dependencies: how a spec points to the layers it consumes

The skeleton (§1) has the **Dependencies** block ("which layers it depends
on — and which it's forbidden to call"). This is the **concrete format**
of that block, and it's what materializes the reuse edge that the pillars
presuppose: `STRUCTURE.md` §4 ("when an interface spec demands a change in
a usecase spec, it's the Structure that says that dependency exists and in
which direction") and `TRACEABILITY.md` §"How edges enter the map" (the
**`declared`** edge — the conceptual dependency that inference doesn't
catch).

### Why declare, and not infer

A layer (usecase, repository, service, hook, store) exists because it's
**reused**: consumed by N units. Reusable behavior needs to be **regulated
by its own spec** — otherwise every consumer redefines it and diverges.
It's the same reason a component reused by N screens gets its own spec
(Atomic), just on the **data** axis instead of the **visual** one. The
dependency is **declared by whoever has it** (the consuming spec), not
inferred from the import: the declaration is honest (the author states "I
depend on this"), stable (doesn't break when imports are refactored) and
language-agnostic (Anchors only reads text).

### The mechanism: two tables linked by code

The consuming spec (a screen, a usecase) already has a **data contract**
table (input/output). To avoid bloating that table with the origin of
each field, the dependency lives in its **own, coded table**, and the
contract just **points to the code**:

**Dependencies Table** — one row per data origin, with a code local to the
spec:

| Code | File                      | Method         | Layer      |
| ---- | ------------------------- | -------------- | ---------- |
| DEP1 | `stores/auth.store.ts`    | `useAuthStore` | store      |
| DEP2 | `hooks/useAuth.ts`        | `signIn`       | hook       |

- **Code** — `DEPn`, **local to the spec** (numbered from DEP1; unique
  within the spec). The graph's real edge comes from the **File**, not the
  code — the `DEPn` is just the pointer the contract uses.
- **File** and **Method** are **separate** columns: the same method
  (`signIn`, `list`, `get`) can live in different files depending on the
  situation; merging them would lose that reference. The edge links to
  the **file** (the graph node); the **method** is **edge metadata** — it
  enables fine-grained impact ("changed `signIn` → only screens that use
  `signIn`").
- **Layer** — the type (§3) of the dependency, for layer-boundary gates
  ("a screen can depend on hook/store/service; a repository can't depend
  on a screen").

**Data contract** — the origin column now references the `DEPn` instead of
free text:

| Field       | Origin | Required | … |
| ----------- | ------ | -------- | - |
| `isLoading` | DEP1   | ✅       | … |
| `user`      | DEP2   | ✅       | … |

### The edge born from this

Each row of the Dependencies Table becomes a **`depends-on`** (declared)
edge from the consuming spec to the referenced **file**, carrying the
**method** as metadata. It's through these edges that Propagation **flows
down** through the data layers: the repository's spec changes → the
screens that consume it go stale (the reuse wave). It's the track that was
missing for the data-propagation described in `PROPAGATION.md` §7 to flow
beyond a single unit's trio.

> Variable depth per project. The dependency chain has whatever length
> that project's Structure declares: one project goes `screen →
> repository/service` directly; another goes `screen → usecase →
> repository → service`. The table is the same at each link — a usecase
> declares its repository/service dependencies the same way a screen
> declares its own. The recursion is propagation following the blueprint
> (`STRUCTURE.md` §4).
