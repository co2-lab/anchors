<!-- @anchors
  code: EXPRS
  updated_at: 2026-11-10
  layer: gate
-->
# ExternalProven — a file that calls a real system is proven by a test that reached it

> **Code**: `EXPRS`

## Overview

A unit test proves a file against doubles. A double of a database or an API accepts whatever it
is given, so a file that sends the real system a query it refuses stays green in every unit test.
Measured in MIF: fourteen defects crossed such a boundary — an index that did not exist, a
condition outside the key schema, an owner field the API refused —, each green against a double.
The unit tests were right about the unit; no gate asked whether any proof of the file ever met
the real system.

The project declares its real systems (`external_systems:`, EXSSA): how a call to each is
recognised, and the suite layers that reach it. A code file that calls one asks a test of its
own — derived beside it by the project's templates, importing it, or testing the feature that
covers its spec — that passed in one of those layers. A file whose calls stay behind a flag, or
that only a CI may run against the system, waives it with `@no-external-proof: <reason>`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the systems | `external_systems:` | — | the configuration (EXSSA) |
| the file's tests | derived beside it (`derived`), reaching it by an edge, or testing the feature of its spec | — | the map |
| a test's proof | its ingested signal, by suite layer | — | the ingestion |

## Effects

| Effect | Description |
| --- | --- |
| `EXPRS-B01` | A code file whose content matches a declared system's `calls`, and no test of which passed in a suite of that system's `proven_by` layers, fails, naming the system, the layers and how many tests of the file were seen; one with such a test passes. (`checkExternalProven`, `testsOfUnit`, `provenIn`) |
| `EXPRS-B02` | A project that declares no system, a file that is no code, a file that calls no declared system, or one that waives the proof with a reason, is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `EXPRS-E01` | No map is built. | Pending, saying there is no map. | The file's tests and their signals are the map's; without them, Pass would claim a proof nobody checked. |
