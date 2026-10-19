<!-- @anchors
  code: SCRQD
  layer: desenho
-->

# Section required — a spec section the code calls for, and its absence bars

> APPROVED on 2026-10-18: D01–D04 as proposed (the user). Asked by MIF, from measurements on
> its 211 component specs. Built in v0.1.355.

## The problem

A spec section is charged today only by the gates that read it once it exists. When the
code calls for a section and the spec has none, nothing is said:

- `testid-consistent` returns `~` (nothing to confront) when the spec declares no handle.
  Its spec end is never crossed.
- A data contract that is absent is not checked against the code that depends on a store.
- A screen that navigates, or shows a message, with no section saying so passes every gate.

Measured at MIF on 2026-10-18, across 211 component specs:

| Case | Specs |
| --- | --- |
| Expose a test handle and have no `## Test IDs` | 67 |
| Depend on the hook/store layers and have no data contract section | 34 of 35 |
| Navigate and have no `## Navegação` | 8 |
| Show a message and have no `## Mensagens ao Usuário` | 8 |

The user's rule: an absence that the code calls for bars; nothing stays silent.

## The design

A new gate, `section-required`, on specs. It measures one thing: does each section that the
unit's code calls for exist in the spec?

The project declares the requirement per layer, in its own vocabulary. The trigger is
something Anchors already measures, never a syntax of its own:

```yaml
layers:
  component:
    requires_sections:
      - section: "Test IDs"
        when: {exposes: test_handles}          # the code matches dialect.test_handles.exposed
      - section: "Contrato de Dados"
        when: {depends_on_layer: [hook, store]} # the code depends on a node of these layers
      - section: "Navegação"
        when: {flag: "@navigates"}             # the code carries the flag
      - section: "Mensagens ao Usuário"
        when: {code_matches: 'Alert\.alert\(|showToast\('}  # a pattern of the project's
```

Triggers, one per entry, from what the core measures:

| `when` | Holds when |
| --- | --- |
| `exposes: test_handles` | some file the spec specifies exposes a handle (`dialect.test_handles.exposed`) |
| `depends_on_layer: [..]` | the unit's code depends on (`depends-on`) a node of a named layer — on its values by default; `kinds: [value, type]` counts an import of its types only (`@dep[type]`) too (v0.1.357, asked by MIF: 22 of 36 charges were components importing only a store's types) |
| `flag: "@x"` | some file the spec specifies carries the flag |
| `code_matches: <regex>` | some file the spec specifies matches the project's pattern |

### Verdicts

- **Fail:** a triggered section is missing. The message names the section and the reason,
  e.g. "`Login.spec.md` declares no `## Test IDs`, and `Login.tsx` exposes `:login-submit`".
  The fix is in the spec.
- **Pass:** every triggered section is there.
- **Skip:**
  - the spec's layer declares no `requires_sections`;
  - nothing triggers;
  - the spec has no code. That absence is `unit-complete`'s to charge.
- **Pending:** a trigger cannot be measured — `exposes: test_handles` with no
  `dialect.test_handles.exposed`, or a map missing. It is never a silent Skip.

A section is found by its heading, at any level, with its title compared ignoring case and
surrounding spaces. The project's own vocabulary is used as written. No translation table
is involved, because the titles are the project's.

### Configuration checks

The config refuses, naming the entry:
- an entry with no `section` or no `when`, or with two triggers;
- an unknown trigger;
- a `depends_on_layer` naming a layer that does not exist;
- a `code_matches` that does not compile;
- a `flag` that does not start with `@`.

### What it does not do

- It does not check the section's content. `testid-consistent`, the data gates and
  `nav-matches-spec` do that once the section exists. This gate only makes them have
  something to confront.
- It knows no framework. The navigation trigger is the project's flag, and the message
  trigger the project's pattern.

## Decisions

| Code | Question | Proposal |
| --- | --- | --- |
| `SCRQD-D01` | A new gate, or a check inside `spec-sections`? | **A new gate.** Its failure has its own reason and fix, and a project can make it informative while adopting it without touching `spec-sections`. |
| `SCRQD-D02` | Declared per layer, or once for every spec? | **Per layer.** A component's spec and a handler's spec call for different sections. |
| `SCRQD-D03` | Blocking by default? | **Blocking**, as the project's rule asks. A project adopting it over an existing tree sets `blocking: false` while it writes the sections. |
| `SCRQD-D04` | May a spec state that a triggered section does not apply? | **No escape in the spec.** If the code triggers it, the section is written. A section that says "none — the handle belongs to the caller" is written, and passes. |

## Phases

| Phase | What | Proof |
| --- | --- | --- |
| `SCRQD-W01` | `layers.<l>.requires_sections`, validated; the gate with its four triggers; the spec guide. | Each trigger fails without its section and passes with it. MIF declares its four requirements and the gate names the 67 + 34 + 8 + 8. |
