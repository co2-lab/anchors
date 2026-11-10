<!-- @anchors
  code: PRFGP
  layer: desenho
-->

# Proof gaps — the defects that reached a release with every gate green, and what each would have needed

> APPROVED on 2026-11-10 (the user: "todos do meu ponto de vista sao uteis, vamos fazer e deixar
> pro usuario decidir qual usar"). Source: MIF (session mif-e3), the case MIF-1S and the analysis of
> 123 commits with a `Bug:` footer, `~/hd/co2lab/MIF/pendencias/anchors-bugs-nao-dedurados.md`
> (read-only; not ours to edit).
>
> The direction: build every proposal that is Anchors' to build, each one **opt-in** — a gate the
> project declares, a field it sets, a tag it writes —, so each project decides which to use. None
> joins the default gates of `anchors init` until a peer has used it and asked.
>
> This file is the plan of record across sessions: the status column of [The items](#the-items) is
> kept current as each ships.

## The case that opened it (MIF-1S)

The Alerts screen lists two kinds of item normalised into one type (`source: 'audit' |
'notification'`). A tap marked any item as an audit (`markAuditSeen(it.id)`); for a notification the
id does not exist as an AuditItem and AppSync refused it. Every gate passed: scenarios proven, line
and branch coverage, mutation 100%. Four holes:

| Hole | What happened |
| --- | --- |
| H1 — variant not proven | the rule speaks of two cases; `scenario-coverage` accepts ONE scenario with the tag; a sibling error rule (E02) stays "proven" by the only test it had |
| H2 — mutation blind to shared behaviour | invalidation and error propagation have no branch per variant; every mutant dies by variant A's test; 100% with no proof of B |
| H3 — the double checks shape, not domain | `mock-stamped` proves the double's signature matches the real one, but `markAuditSeen(notificationId)` is accepted; the E2E was blind because the hook returns early under `EXPO_PUBLIC_E2E` |
| H4 — spec and test contradict | rule B02 says "markSeen records in notificationLog"; the test with the same tag asserts `expect(markAuditSeen)…`; the gate sees tag + green test |

## The analysis (123 commits)

| Reason | Commits | Hole |
| --- | --- | --- |
| UI text/visual no gate judges | 25 | — |
| Rule absent from the spec | 15 | — |
| Boundary with a real system, the double accepted the wrong thing | 14 | H3 |
| Concurrency / order / idempotency | 12 | — |
| Wrong rule (spec wrong, proven) | 10 | H4 |
| Variant not proven | 10 | H1 |
| Real data the fixture did not represent (time zone, file format, text) | 10 | — |
| Internal contract: writer and reader of the same data diverge | 7 | H3 |
| Infra wiring (IAM, env, table/secret not granted) | 5 | — |
| Security / authorisation (no negative test) | 5 | — |
| External contract changed | 2 | H3 |
| Unhandled failure | 2 | — |
| Blind E2E | 1 | H3 |
| Other | 3 | — |

MIF's reading: the boundary is the largest technical hole (~29, 24%); gates born late found many
bugs afterwards and would have barred them at the origin (Stryker campaign 11, VR per state 15);
specs written after the code prove the code, not the desired behaviour; concurrency does not show
in deterministic tests.

## The items

Each item: what MIF asked, whose it is, what Anchors already has, what is built, how it is opted
into. Status: `planned`, `in progress`, `shipped vX`, `MIF's` (nothing for Anchors to build).

| ID | From | Item | Already in Anchors | To build (Anchors) | Opt-in | Status |
| --- | --- | --- | --- | --- | --- | --- |
| G01 | P3 | **A new gate bars the legacy too, after a deadline** | `gates[].legacy: {before, severity}` softens targets last changed before a date — forever | `legacy.until: <date>`: from that date the legacy severity no longer applies and the gate's own does; before it, a softened verdict says until when and what it will do after (GTENG-B37, CNFGO-B77/E13) | the field | shipped v0.1.389 |
| G02 | A, P5 | **Variants as units of proof** | `spec-feature-match` (a scenario per requirement), `data_states` (`DS-` states of a unit), `examples-match` (each Examples row is a run case) | a rule declares its variants (`@variants: audit, notification` on its line, or a discriminant the guide names); the gate asks one scenario per variant, or an Outline whose Examples cover each; counts Examples rows | the tag / a gate `variants-proven` | planned |
| G03 | A (2nd half) | **Siblings proven only on the old variant** | revisions and stamped edges | when a rule gains a variant, the other rules of the unit whose proofs do not cite the new variant are flagged divergence ("proven only on <old>") until a proof names it | part of `variants-proven` | planned |
| G04 | B, P11 | **Rule × test coherence, judged** | judgment gates (`measures: judgment`, `ask`, `guide`), `anchors judge` | a built-in question `rule-test-coherent`: for each tagged test, what it asserts (expect targets, function/table names) against the rule's text; an explicit contradiction is pending; a default `ask` shipped, the project may replace it | declaring the gate | planned |
| G05 | P11 | **Negative cases required** | `scenario-asserts` | a rule tagged `@matching` or `@timing` (names configurable) asks for at least one scenario whose outcome is a refusal/no-match | the tags | planned |
| G06 | P4 | **Order, idempotency, concurrency** | — | rule tags `@order`, `@idempotent`, `@concurrent` (names configurable); each asks a proof of its kind: a test run in both orders / twice / interleaved, recognised by a declared test-name or helper pattern; guide text for the five patterns MIF lists (double tap on a writing action, read-then-write with conditional write, matching in both orders, total tie-break in sort, timer stopped on every exit) | the tags | planned |
| G07 | P6 | **Internal contract writer → reader** | `@dep` edges, `evidence-fresh`, `proof-crosses-boundary` | a `@writes: <data>` / `@reads: <data>` pair on the units that persist and consume the same data; an edge writer→reader in the map, so a new revision of the writer's shape stales the readers' evidence; a gate flags a fixture literal for data a declared builder produces | the marks | planned |
| G08 | P1 (Anchors part) | **Boundary proof** | `proof-crosses-boundary` (a mirrored rule's proof reaches the other unit), `mock-stamped`, `mock-typed` | boundaries declared in config (`boundaries: [{name: appsync, calls: <pattern>}]`); a unit that calls a declared boundary needs ≥1 proof from a suite the project marks as integration/live (`suites[].reaches: [appsync]`) | declaring boundaries | planned |
| G09 | P2 (Anchors part) | **Recorded doubles expire with the API version** | `mock-stamped` recomputes a stamp from the real snippet | a recorded fixture carries `@recorded: <api> <version> <date>`; the project declares each API's current version; a different version or an age past a declared limit makes the proofs that use it stale | the mark + config | planned |
| G10 | C, P14 | **Effect erased in test mode** | `fallible_patterns`-style dialect patterns | a dialect pattern list `test_bypass` (e.g. an early return under `E2E`/`TEST`/`CI` env); each match asks an alternative proof of the skipped path (an integration suite proof of the unit) or `@no-bypass-proof: <reason>` on the line | the dialect patterns | planned |
| G11 | P7 | **Backend enum vs presentation** | `value-anchored`, `categoria-canonica` (MIF's own script gate) | `value-anchored` also reads comparisons, `case` labels and object keys; a gate "raw enum on screen" by a declared enum source and the presentation layers | config | planned |
| G12 | P8 | **Data states for screens** | `data_states` (DS-), `vr-states-covered`, `state-transitions` | a configurable required set of states per layer kind (loading, error, not found, empty, offline, no access…); error ≠ empty; error offers retry — the set is the project's, no default list imposed | config | planned |
| G13 | P9 | **Date / time zone patterns** | MIF's `padrao-inseguro` (a pattern gate) | nothing new in the engine: ship the patterns as an example pack (`new Date('YYYY-MM-DD')`, `toISOString().split('T')[0]`, `-31T23:59:59Z`) | the pack | planned |
| G14 | P10 | **Coherence between specs** | `cross-unit-citation`, doctrine consumers | an input "from X" asks a matching output in X; an output names its destination; a UI action with a persisted effect points to the backend rule that produces it | a gate `specs-coherent` | planned |
| G15 | P12 | **Transitive wiring** | `dep-chain` | `anchors deps <unit> --closure` prints the `@dep` closure (machine-readable), for the project's script to cross env/tables/grants | the command | planned |
| G16 | P13 | **Presentation lints** | — | the Anchors part only: a contrast (AA) check over VR baselines and a pt-BR text gate are out of scope until asked; the lints are MIF's eslint | — | MIF's |
| G17 | P15 | **Minors** | — | authorisation guard declared → a negative 403 scenario; CSS class / `var(--x)` exist; a late response after unmount; a cross-cutting invariant with its implementers listed | each its own small gate or tag | planned |
| G18 | D | **Domain-aware doubles** | `mock-stamped` | out of scope for now (MIF agreed): a double fed ids records the model they come from | — | deferred |
| — | P1, P2 (MIF part), P9 (TZ in CI), P13 | fakes generated from the schema, recorded responses, `TZ=America/Sao_Paulo`, eslint rules | — | — | — | MIF's |

## Order

By MIF's priority (P1, P3, P4, P5, P6 cover ~64 of 123) and by cost, deterministic first:

1. **G01** (legacy deadline) — small, and it turns every existing gate into one that ends the legacy.
2. **G02 + G03** (variants) — the hole of the opening case.
3. **G04 + G05** (rule × test, negatives) — the only partial remedy for wrong rules.
4. **G06** (order/idempotent/concurrent).
5. **G07** (writer → reader).
6. **G08 + G09** (boundary proof, recorded doubles).
7. **G10** (test bypass).
8. **G11–G15, G17** in the order peers ask.

Each item ships alone, with the release workflow, and MIF is told what to declare to use it.

## Rules that bind every item

- Opt-in: a project that declares nothing sees no new finding.
- Generic: tag names, patterns, suites, boundaries are configured; nothing names AppSync, Stripe or
  Expo in the engine.
- Each item's spec rows and scenarios before its code; each new test fails on the old code.
- No item edits MIF's repository; MIF adopts it.

## Resuming this work

Read [The items](#the-items) for the next `planned` row in [Order](#order); update its status when it
ships, with the version. A new proposal from a peer joins the table as a new row, `planned` only after
the user approves it.
