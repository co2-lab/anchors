package governance

// planGuide é a régua da fase PLANEJAR (Planejamento — a origem do movimento). Um
// plano no Anchors NÃO cria código; ele decide QUAIS specs precisam nascer ou mudar,
// e por quê. É a âncora que semeia as outras âncoras. Este guia embutido casa com a
// versão do binário — a IA o lê antes de escrever qualquer plano.
const planGuide = `# Plan guide (the ruler of the PLAN phase)

A plan is the ORIGIN of movement in Anchors. It writes neither code nor features —
it decides WHICH specs need to be born or change, and why. The plan seeds the specs;
the specs seed the rest. If you are about to write code from a
plan, the plan failed: it should have pointed at a spec.

## GOLDEN RULE: approve the scope BEFORE writing the file

Writing the plan file is an IRREVERSIBLE ACT of scope: as soon as it is saved in the
plan layer, the watcher detects it and QUEUES the first task ("specify") — the
machine starts executing. That is why the user's approval comes BEFORE the file exists,
not after:

1. In the CONVERSATION, draft the scope (objective, reason, the LIST of specs, phases, out of
   scope, definition of done) and present it to the user IN TEXT in the reply.
2. ITERATE until they approve the scope explicitly.
3. ONLY THEN write the plan file in plans/.

Do not write the file to "show how it turned out" and ask afterwards — that already starts the
belt. If the watcher is on, saving = starting. Treat the plan file as the
commit of the scope: it only exists after the "yes".

## What a plan IS and what it is NOT

- It IS: a scope decision. "These N specs will be born/change, in this order, for these
  reasons, and this is how we will know it finished."
- It is NOT: a detailed design (that is the spec), nor a list of coding tasks.
- A plan can be small (one spec) or large (a slice of product). The plan's size
  is the size of the scope change — not of the implementation effort.

## The structure of a plan

Write the plan as a document in the project's plan layer (see anchors.yaml
— typically plans/). A plan has:

### 1. Objective
One or two sentences: what changes in the product and for whom. No implementation jargon.

### 2. Reason
Why now. What is wrong, missing, or requested. It is what survives when
someone reopens the plan in six months asking "why did we do this?".

### 3. Seeded specs
The heart of the plan. A LIST of the specs that will be born or change. For each one:
  • the spec file (new or existing) — e.g.: features/larder/AddItem.spec.md
  • born or changed? if changed, what changes in it
  • one line of what it must cover (the detail goes in the spec, not here)
Do not list code, screens, endpoints. If you feel the urge to list those, it is a sign
that the corresponding spec has not been thought through yet — add the spec to the list.

### 4. Phases (order and dependencies)
Group the specs into small, independently deliverable PHASES, in the order in which
they must be born. If one spec depends on another first, it comes in a later phase. The
plan is a sequence, not a bucket:

  ## Fase 1 — <name>
  - features/larder/AddItem.spec.md — item registration
  - features/larder/ItemCard.spec.md — item card
  ## Fase 2 — <name> (depende da Fase 1)
  - features/larder/ExpiryAlert.spec.md — expiry alert

**PROGRESS DOES NOT LIVE HERE.** It lives in the companion file ending in
'-progress.md', beside the plan, and that is where '[x]' is marked. 'anchors new' creates it
alongside.

The reason is the difference the plan must preserve: the plan is a DECISION, and changing it
must mean the decision changed. While the checkboxes lived here, marking a
phase done was CHANGING the plan — and the gate that demands a justification for a change could
not distinguish "I finished phase 1" from "I changed the project's direction". Worse: the
judgment of the spec the phase delivered fell with it, because it stores the plan's rev.
Finishing the phase invalidated the verification of what the phase produced.

The progress file stays out of the map on purpose: it exists to change, and
no gate should demand a justification from a file whose job is to record that the
work moved.

### 5. Out of scope
What this plan deliberately does NOT do. It closes the door to leaking scope. Each
item here is a temptation named and refused.

### 6. Definition of done
How to know the plan finished. Normally: "all the seeded specs exist,
implemented, with a feature and a test, and 'anchors check' passes the blocking gates".
If there is an accepted debt (an honest opt-out), name it here — do not leave it implicit.

## Plan rules

- ALWAYS by spec. Each line of work in the plan ends in a spec, never in a
  code file. That is the spec-first discipline: the plan seeds specs, the spec is
  the source of truth.
- THE WHAT, not the HOW. The plan says what needs to exist and why. The HOW (behaviour,
  contracts, invariants) belongs to the spec; the implementation to the code. Do not describe
  an algorithm or a data structure in the plan — if you felt the urge, it is spec content.
- RECONCILE before seeding (anti-rework). Before listing a new spec, look at
  what ALREADY exists — other plans and specs. Do not plan what another plan already delivers, do not
  create overlapping scope nor an impossible ordering. If you change a decision because it
  conflicts with something existing, note the WHY in the plan — the reason must survive.
- RESPECT the structure. The seeded specs fall into the layers anchors.yaml
  declares. If a spec has no layer to live in, either the structure is incomplete
  (report it to the user) or the spec is in the wrong place.
- IDENTITY early. If the project uses scenario codes (identity — TRACEABILITY),
  the spec will carry them; the plan does not need to invent them, but it must know that each
  requirement in the spec will receive one.
- AMBIGUITY becomes a TODO, never a guess. If a spec's scope is uncertain, write
  an explicit TODO in the plan instead of guessing. An honest TODO is better than an
  invented decision the worker will implement wrongly.
- BUT the TODO cannot govern the EXISTENCE of an item. This is the distinction that decides
  whether the plan is executable:

    ✓ delegating WHERE it is decided — "the indexes go in the SCHEMA spec, decided upon seeing
      the current schema". The item exists; the CONTENT of one of its decisions lives in another
      artifact, which will be born. Whoever executes knows where the answer will be.

    ✗ deferring WHETHER the item exists — "it is born IF mobile calls via API Gateway; otherwise
      mark it out of scope". Here whoever executes must RESOLVE the condition before
      knowing what to do.

  The second has two outcomes, and both are bad: the agent decides alone (and the product
  decision becomes an implementation choice, without going through whoever should have), or it stops and
  asks (and the automated flow dies there). Measured in a real project: the agent
  decided — it got it right, and nobody had asked it to decide.

  Before seeding an item, ask: "does whoever executes this need to DECIDE something
  in order to start?". If so, decide now, or seed first the item that ANSWERS the
  question — the decision becomes work, not an invisible prerequisite.
- AN OPEN QUESTION stays OUTSIDE the checklist. The plan CAN (and should) record what is still
  unknown — in prose, in its own section, with who decides. What it cannot do is
  hide the question INSIDE a seeded item, where it looks like work and is a decision.
- "A OR B" IS A DEFERRED GUESS — decide it in the plan, not in the spec. A plan that says "use
  repository OR service", "write here OR there", "sync OR async" is NOT an honest
  TODO: it is an ARCHITECTURE fork that the spec will resolve in silence, with
  nobody weighing the consequence (e.g.: "repository" may mean duplicating a rule
  outside its place). Close every "or" NOW — with the user — or mark it as a NAMED TODO
  that blocks descending to the spec until it is resolved. Never leave the choice to
  whoever specifies.
- HONESTY about debt. If the plan decides to skip a gate (e.g.: no test for now),
  that is an explicit decision in the Definition of done — never a silence.

## After writing the plan (already approved)

The file is only born after the user's "yes" (see the GOLDEN RULE above). On saving it,
the watcher queues the first task ("specify"). From there on the work flows through the
queue: you (or a worker in background) run 'anchors next' and follow the playbook's
cycle. The plan is the scope contract; the specs are its execution. If the scope
changes midway, that is a conscious EDIT of the plan (and a new conversation), not a silent
drift of the specs.

## Anti-patterns (refuse them)

- Writing the plan file BEFORE approval → it starts the belt without a "yes"; present
  the scope in text in the conversation first. The file is the commit of the scope.
- A plan that lists code files → the spec is missing; go up one level.
- A plan without "out of scope" → it will leak; name what it does not do.
- A plan without a "definition of done" → there is no way to know when to stop; add one.
- A plan that changes the ruler (a guide) in passing → changing a ruler is a plan of its own,
  not a side effect. Separate them.
`
