<!-- @anchors
  code: TSUNT
  updated_at: 2026-10-08
  layer: mapa
-->
# TestedUnits — which code a test tests, found by the project's own derivation

> **Code**: `TSUNT`

## Overview

The map links a unit's test only through the unit: a spec, its feature, the feature's test. A unit
with no spec — a util, a model, what a `regime: declarativo` layer holds — has no edge to its test,
and a gate asking "what does this test test?" had no answer: in the reference app, 130 tests with no
feature had 12 edges among them. The project's derivation answers it without the engine assuming
anything about the language: the templates that say where a unit's test lives (`derived.files`, and
the overrides of the code's layer), applied to each code file.

With `anchor: code` the code file is the unit, and its dir, name and extension fill the templates as
in the map build. With `anchor: spec` the templates hang off the spec's name, so the code file's path
is matched against the `code` templates first, to recover the variables.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | the scan's map, or nil | — | this unit: nil answers nothing |
| the configuration | a project configuration with a derivation, or without one | — | this unit: without `derived` it answers nothing |

## Effects

| Effect | Description |
| --- | --- |
| `TSUNT-B01` | With code as the anchor, a code file tests through the `test` templates filled with its own dir, name and extension. |
| `TSUNT-B02` | With the spec as the anchor, a code file's variables come from matching its path against the `code` templates, and the `test` templates are filled with them. |
| `TSUNT-B03` | An override whose `when` is the code file's layer replaces the default templates of the kinds it declares, for that file. |
| `TSUNT-B04` | A `test` template that resolves to a glob matches every test file of the map it covers; a `code` template is read literally around its variables, so a directory named `[slug]` matches itself and glob characters match only themselves. |
| `TSUNT-B05` | Only test files the map has are answered; each test's units are listed once, in path order, and a project without a derivation answers nothing. |

## Errors

none — the lookup reads only the map and the configuration it is given; a template that does not
match is a file it does not describe, not a failure

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
