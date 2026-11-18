<!-- @anchors
  code: THSAS
  updated_at: 2026-10-08
  layer: gate
-->
# TestHasAssertion — every test asserts something in its body

> **Code**: `THSAS`

## Overview

A test that asserts nothing passes whatever the code does. It carries the scenario's code, its
title matches the scenario, the scenario counts as proved — and the body calls the code and checks
nothing. Every gate that reads titles is green over it.

What an assertion looks like is the test library's (`expect(`, `t.Errorf`, a project helper), so the
project declares it in `dialect.tests.assertion`, with a default for the families that have one.
Where a test's body ends is the language's too: a tests script that knows it says it (`end`), and
otherwise the body is read by its layout, which every language keeps whatever its syntax.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a test file that is not a support file | any other kind, a support file | this unit: Skip |
| the tests source | the dialect's `tests`, pattern or script | none declared | this unit: Skip, nothing to measure |
| the assertion | the dialect's `tests.assertion` | none declared | this unit: Skip, nothing to measure |
| the gate entry | `labels: true` or nothing | — | the configuration load |

## Effects

| Effect | Description |
| --- | --- |
| `THSAS-B01` | A test whose body holds no assertion fails, each named by its line and title; a file whose tests all assert passes. |
| `THSAS-B02` | A test's body is the block its line opens: a line ending in a bracket or `do` closes on the first later line back at its indentation that starts with a closer; any other line that opens a block ends it before the first later line back at its indentation; a line whose next line is not deeper is a block by itself. (`blockEnd`) |
| `THSAS-B03` | A line that starts inside a multi-line literal — opened by a backtick or a triple quote — is text, and never ends a block; a backtick quoted alone is a character, not a delimiter. |
| `THSAS-B04` | With `labels: true` on the gate, a test opened and closed on its own line with an empty block is a label for the block around it, and that block is what must assert; without it, the empty test is a test with no assertion. |
| `THSAS-B05` | When the tests source says where a test ends (`end`), its body runs from its line to that one. |
| `THSAS-B06` | A node that is not a test or is a support file, a project with no tests source or no assertion, and a file where the source lists no test are skipped. |
| `THSAS-B07` | A test that carries `@no-assert: <why>` in its body or on the line above it is left out; a waiver with no reason waives nothing. |
| `THSAS-B08` | A test listed at a line the content does not have is read as an empty body, never past the content's end. (`checkTestHasAssertion`, `testBody`) |
| `THSAS-B09` | A line back at the opener's indentation that starts with a closer and ends opening a block — the table of an `it.each([` closing into the test's body, a `} else {` — goes on with the block; only a closer that opens nothing ends it. (`blockEnd`) |

## Errors

| Error | When | What the user sees |
| --- | --- | --- |
| `THSAS-E01` | The tests cannot be listed (a script that fails, output outside the contract), or the assertion does not compile | the gate fails with the reason, instead of answering as if there were no test |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
