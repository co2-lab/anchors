<!-- @anchors
  code: TICTS
  updated_at: 2026-09-26
  layer: gate
-->
# TestIDContract — a test handle is one contract with four ends: the code exposes it, the spec declares it, a consumer queries it

> **Code**: `TICTS`

## Overview

A test handle (the attribute a project uses to mark an element for tests and automated flows) is one
contract with four ends: the code exposes it, the spec declares it in its test surface inventory, the
feature may describe it, and a test or a flow queries it. The testid-consistent gate confronts the ends
of that contract at once, starting from the spec, and reports per handle, with every end on the same
line. One gate and not several partial ones, because a handle renamed in the code used to produce two
findings in two gates ("exposed without declaring" for the new name, "declared without exposing" for the
old one) when the fact is one: somebody renamed and did not propagate.

Three files hold the contract, and they must agree on what a handle is:

- The gate itself decides when there is something to confront and what fails.
- The parsing half recognises the handle where it is written: in the code (a literal value, a template
  whose suffix is runtime data, a prop that carries a handle to a child, the branches of a conditional)
  and in the spec (the inventory table or list inside the test surface section).
- The consumer half finds who queries a handle: the test linked through the spec's feature, the tests
  beside the spec and in its sibling folders, and the end-to-end flows the project declares.

The attribute is the project's declaration, never a guess: without it the gate skips, because inferring
one would report green over what was never checked. The fourth edge, a consumer querying a handle no
code exposes, is deliberately left out: the flows are the whole project's tree, and charging one spec
for another screen's handle produced 829 findings in a single spec when it was tried. That question
belongs to a project-scope gate.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | a spec | any other kind | this unit: the gate skips it |
| the handle attribute | the attribute the project declares for test handles | no declaration | this unit: the gate skips |
| the map | the spec's links to its code and to its feature, and the feature's link to its test | no map | this unit: the gate is Pending |
| the linked code | code files the spec specifies, readable under the root | none, or none readable | this unit: the gate skips |
| the e2e surface | a declared surface whose file template, or an override's, gives a path | a surface with no path | this unit: it contributes no consumer |

## Effects

### When the gate confronts

| Effect | Description |
| --- | --- |
| `TICTS-B01` | The gate skips a node that is not a spec, and is Pending without a map. |
| `TICTS-B02` | Without a declared handle attribute the gate skips. |
| `TICTS-B03` | A spec with no readable linked code skips: there is no end that exposes a handle. |
| `TICTS-B04` | A unit whose code exposes no handle and whose spec declares none skips: it has no test surface to contract. |

### The ends of the contract

| Effect | Description |
| --- | --- |
| `TICTS-B05` | A handle exposed by the code, declared by the spec and queried by a consumer passes. |
| `TICTS-B06` | A handle the code exposes and the spec does not declare fails, named in the report, also when the spec has no inventory at all. |
| `TICTS-B07` | A handle the spec declares and the code does not expose fails, named in the report. |
| `TICTS-B08` | When some consumer surface exists, a handle no consumer queries fails, named in the report. |
| `TICTS-B09` | When there is no consumer surface at all, the queried end is not charged, and the report shows a dash in its place. |

### How a handle is written

| Effect | Description |
| --- | --- |
| `TICTS-B10` | The handle is read under the attribute the project declares, whatever its ecosystem. |
| `TICTS-B11` | A literal handle value counts, also when it reaches the code through a derived prop or an object key whose name contains the attribute. |
| `TICTS-B12` | A template handle, whose suffix is runtime data, and a prop that gives a child the head of its handles, both expose the static prefix as a wildcard. |
| `TICTS-B13` | In a conditional handle only the literals of the branches are handles; the literals of the condition are not. |
| `TICTS-B14` | The spec's inventory is read only inside its test surface section, whose title may carry a qualifier, and the section ends at the next heading. |
| `TICTS-B15` | In a table row of the inventory only the first cell is the id; on any other line every quoted id counts. |
| `TICTS-B16` | The attribute's own name, quoted in the inventory, is not a declared id. |
| `TICTS-B17` | A wildcard at either end covers the concrete ids it opens; a concrete id does not cover a wildcard. |

### Who consumes a handle

| Effect | Description |
| --- | --- |
| `TICTS-B18` | A handle is queried when a consumer mentions it, with or without the mark; a wildcard handle, when a consumer mentions its prefix. A mention counts only at an id boundary (`TICTS-B22`). |
| `TICTS-B19` | The end-to-end flows are consumers when the project declares the surface with a path: the surface's file template, or else the first override that gives one, whose static prefix is read as a whole; a surface with no path contributes nothing. |
| `TICTS-B20` | The test files beside the spec and in the sibling folders of its folder are consumers, even when flows exist, without any edge to them. |
| `TICTS-B21` | The static root read for a surface is the directory before the first placeholder or glob wildcard (`*`, `?`, `[`) of its path: `e2e/**/*.yaml` and `e2e/login-*.yaml` read `e2e`, `apps/x-{{module}}/flows` reads `apps`, and a path that opens with a wildcard reads the project root. |
| `TICTS-B22` | A mention is a query only when no id character (letter, digit, `.`, `_`, `-`) touches it before, and, for a concrete handle, after: a consumer naming only `abcd-screen-header` or `my-abcd-screen` does not query `abcd-screen`. |
| `TICTS-B23` | The test the spec's feature is tested by is a consumer wherever it lives: the link runs spec → feature → test, two hops, and it needs no neighbour folder to reach the test. |
| `TICTS-B24` | Each line of the report shows whether the spec's feature describes the handle — a ✓ when the feature mentions it, a ✗ when it does not. It is information, never the reason for the failure: not every handle belongs in a written scenario. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TICTS-I01` | One handle is one line of the report, whatever its spelling at each end: the mark is not part of the identity, and the line spells it as the code does. | a handle exposed with the mark and declared without it, and queried by nobody, appears once, with the mark |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TICTS-X01` | A handle only a consumer mentions, which the spec does not declare and the code does not expose, is not charged to the spec. | The flows span the whole project; a spec cannot know which of its handles a foreign flow meant, and the question belongs to a project-scope gate. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TICTS-E01` | A linked code file cannot be read. | It is not an end of the contract; when no linked code can be read, the gate skips as a spec without linked code. | Reading nothing and then charging every declared handle as "not exposed" would accuse the spec of a divergence the gate never measured. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Derived` | config — the handle attribute and the e2e surface |
| DEP2 | `internal/mapx/model.go` | `Graph`, `EdgeSpecifies`, `EdgeCoveredBy`, `EdgeTestedBy` | mapa — the spec's code, feature and test |
| DEP3 | `internal/i18n/i18n.go` | `T` | apoio — the localized report |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
