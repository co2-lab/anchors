<!-- @anchors
  code: TFSTM
  layer: infra
-->
# TimeoutFailure — telling a test that ran out of time from one that failed its assertion

> **Code**: `TFSTM`

## Overview

A test may fail because it waited too long, not because it is wrong: under load a wait runs out
before the work it waits for ends. `anchors test` retries such a test once (DESIGN-retry-on-load.md),
and never a test whose assertion failed. This unit tells the two apart, from the failure the
report carries.

The failure's message is not enough. Measured at MIF, a `waitFor` of the testing library that runs
out rethrows its last assertion (`Expected: true / Received: false`) and never says "timeout"; only
the line of the test where it failed, `await waitFor(...)`, tells the wait. So the pattern is
matched against the failure and against the lines of the case's own test file that the failure's
stack points at.

## Effects

| Effect | Description |
| --- | --- |
| `TFSTM-B01` | A failed case whose failure matches the pattern — a test or a hook past its timeout — ran out of time. (`TimedOut`, `DefaultTimeoutPattern`) |
| `TFSTM-B02` | A failed case whose failure does not say it, but whose stack points at a line of its own test file that matches the pattern — a wait — ran out of time. |
| `TFSTM-B03` | A failed case whose failure and pointed lines match nothing failed its assertion; a line of another file the stack points at does not count; a case that passed did not run out of time. |
| `TFSTM-B05` | A test file the stack points at that cannot be read counts for nothing: the case is told by its failure alone — a missing file never turns a failure into a retry. |
| `TFSTM-B04` | The tests that passed only on a retry are kept, each with when and under what load; added entries join the earlier ones, the latest 500 kept; a missing or unreadable record is empty. (`Retried`, `AddRetried`, `LoadRetried`, `RetriedFile`) |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the case | a case read from a JUnit report, its file relative to the root | — | `ParseJUnit` |
| the pattern | a regular expression, the default or the suite's `timeout_pattern` | — | the suite's configuration |

## Errors

none — a test file that cannot be read is normal flow: its line counts for nothing (TFSTM-B05), and the case is told by its failure alone.
