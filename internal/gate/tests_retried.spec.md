<!-- @anchors
  code: RTSRT
  layer: gate
-->
# TestsRetried — a test that passed only on a retry under load stays in sight

> **Code**: `RTSRT`

## Overview

A suite with `retry_on_load:` runs once more a test that ran out of time while the machine was
loaded, and the run passes when the retry does (DESIGN-retry-on-load.md). The retry lets the
work go on. This informative gate keeps it from hiding a test that keeps running out of time:
over a test file, it lists the tests that passed only on a retry in the last 7 days, how many
times, and under what load the last time. A test that needed 3 retries reads as unstable.

## Effects

| Effect | Description |
| --- | --- |
| `RTSRT-B01` | A test file with no test passed on retry in the last 7 days passes. (`checkTestsRetried`) |
| `RTSRT-B02` | A test file whose tests passed on retry in the last 7 days fails, informing: each test, how many times, the load of the last one; a test retried 3 times or more reads as unstable. A retry older than 7 days, or of another file, does not count. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the record of retries | `.anchors/retried.json`, as `anchors test` writes it | a missing or unreadable record, read as empty | `testsig.LoadRetried` |
