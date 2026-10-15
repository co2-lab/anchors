<!-- @anchors
  code: THPTS
  layer: desenho
-->

# Test handle patterns — the project says how a handle is written; the core reads no syntax

> APPROVED on 2026-10-15 (the user: "tira a sintaxe e passa para configuração do anchors, sem
> fallback pra projetos antigos funcionarem… use o que tem na sintaxe só pra dar exemplo de
> como configurar o gate na lang do projeto, mas se não for configurado no dialect o gate não
> funciona").

## The problem

Reported from MIF on 2026-10-15. `testid-consistent` and `testid-queried-exists` did not see
the handles of SegmentedControl, which a helper composes: `getOptionTestID(base, option)` gives
`${base}-${option}`. A spec that declared them read "missing in the code".

The cause sat in the core. `internal/gate` knew JSX:
- the attribute and its derived props (`backTestID`);
- a template literal;
- a prop ending in `Prefix`;
- the branches of a conditional;
- a suffix composed onto a prop (`${testID}-row`);
- a `:`-marked literal anywhere.

It also knew Maestro's `id:` in a flow, which extensions are source, and that a `.test.` file
is a test. A project's own way of writing a handle had nowhere to go. This is the same defect
as the fallible call's name and propagation (0.1.344, 0.1.346).

## The design

### 1. `dialect.test_handles`

    dialect:
      test_handles:
        exposed:     # where the code exposes a handle — what a unit owns
          - '(?i)\b[a-z]*testID[a-z]*\s*[=:]\s*\{?["''](?P<testid>:?[a-zA-Z][a-zA-Z0-9._-]*)["'']'
          - 'getOptionTestID\(\s*["''](?P<prefix>:?[a-z][a-z0-9-]*)["'']'
        mentioned:   # what else makes a handle exist: a fallback, a table, a composed suffix
          - '`\$\{[a-zA-Z_][a-zA-Z0-9_.]*\}(?P<suffix>-[a-zA-Z][a-zA-Z0-9-]*)'
        queried:     # how an end-to-end flow queries one
          - '(?m)^\s*id:\s*["''](?P<testid>:?[a-zA-Z][a-zA-Z0-9._-]*)["'']\s*$'

Each pattern says what it captures in a named group:
- `testid` is a handle, whole.
- `prefix` is the head of a family of handles whose rest is runtime data; it reads as
  `head-*`.
- `suffix` (`mentioned` only) is a segment composed onto a head from outside the file; it
  reads as `*-seg-*`.

A value one pattern captures as a prefix is not also a whole handle by another. A pattern
that does not compile, or captures with no group, refuses the config on load.

### 2. No pattern, no measurement

The three gates presuppose `dialect.test_handles.exposed`, and `testid-queried-exists` also
presupposes `queried`. Without them the gates are Pending and name the field.

There is no fallback. `derived.test_handle` alone no longer makes them measure: the attribute
does not say how a handle is written. A project upgrading adapts, by declaring its patterns.

### 3. The family's syntax is an example

`config.FamilyTestHandles(family, attr)` holds what the core used to read, as patterns, for the
attribute the project marks with (`ts` only). It is never applied by itself:
- `anchors init` writes it into a new project whose inference found the attribute;
- the governance tips offer it, as the block to copy, to a project that declares a handle gate
  and no patterns.

### 4. What the core still reads

What is Anchors' own, not a language's, stays in the core:
- the spec's inventory (its section and table);
- the `:` mark, which is not part of a handle's identity;
- the matching of wildcards;
- the floor below which a queried head identifies nothing.

The files read for handles are the map's code nodes: the project's layers say which files are
source and which are tests, not a list of extensions. The flows read are every file of the e2e
surface; what a flow is, is what the `queried` patterns find.

## Rules

- `DLCTI-B27`, `DLCTI-B28`: the config and the family's example.
- `TICTS-B02`, `TICTS-B10`…`B13`, `TICTS-B25`: testid-consistent.
- `TQETS-B01`, `B09`…`B15`, `X03`: testid-queried-exists.
- `IDCND-B13`: identity-consistent.
- `GVOPG-B10`: the governance tip.
- `BLCNB-B10`: init.

## What NOT to do

- **Fall back to the family's syntax when the project declares none.** A project would be
  measured by patterns it never chose, and a gate would answer over code it never read.
- **Keep a list of source extensions or test file names in the gate.** The map already knows.

## Limits

The `ts` example reads a conditional handle's two branches with two patterns, the `then`
literal and the `else` literal. A deeper expression is the project's to write.
