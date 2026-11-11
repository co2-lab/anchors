<!-- @anchors
  code: BYPRS
  updated_at: 2026-11-10
  layer: gate
-->
# BypassProven — a path the code skips while a test runs is run by a suite that reaches it

> **Code**: `BYPRS`

## Overview

An effect the code skips while a test runs is proven by no test that skips it. Measured in MIF
(MIF-1S): a hook returned early under `EXPO_PUBLIC_E2E` so the visual captures stayed clean; the
unit tests doubled the write, the end-to-end skipped it, and nothing ran it — the write reached
the store review broken.

The project declares its bypasses and the suite layers that run the skipped path
(`test_bypass:`, TBYSA). A code file with a bypass asks a test of its own — derived beside it,
importing it, or testing its spec's feature — that passed in one of those layers. A bypass whose
path is proven elsewhere, or that guards nothing a test must run, waives it with
`@no-bypass-proof: <reason>` on its line or the line above.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the bypasses | `test_bypass:` | — | the configuration (TBYSA) |
| the file's tests and their proofs | the map's tests of the file and their ingested signals, by suite layer | — | the map and the ingestion |

## Effects

| Effect | Description |
| --- | --- |
| `BYPRS-B01` | A code file with a bypass not waived, no test of which passed in a `proven_by` layer, fails, naming the bypasses' lines and the layers; one with such a test passes. (`checkBypassProven`) |
| `BYPRS-B02` | A project that declares no bypass, a file that is no code, or one whose every bypass is waived with a reason on its line or the line above, is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BYPRS-E01` | No map is built. | Pending, saying there is no map. | The file's tests and their signals are the map's; without them, Pass would claim a proof nobody checked. |
