<!-- @anchors
  code: TLCNT
  updated_at: 2026-09-26
  layer: apoio
-->
# TelemetryConfig — the opt-out is easy to find, easy to get right, and the environment overrides the file

> **Code**: `TLCNT`

## Overview

Decides whether Anchors emits its decision events at all. Telemetry is on by default — the
maintainer's choice — and what makes that choice honest is that turning it off is trivial and
announced. There are two ways, and the first needs no file edit: the `ANCHORS_TELEMETRY`
environment variable, valid for everything including CI, and `telemetry: off` in the project's
configuration, versioned with the project.

Two decisions shape the answer. First, the environment wins over the file: someone running a
command on a CI machine must be able to turn telemetry off without committing, because committing
to turn it off would make one person's decision a change in the team's repository. Second, the
opt-out accepts the words people actually try — `off`, `0`, `false`, `no`, in any case and with
surrounding spaces. Requiring the exact spelling would be a trap: the person would believe they had
turned it off.

The unit also carries the emission settings the emitter reads (whether it is enabled, the endpoint
and the authentication headers, which come from the environment and never from the versioned file).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the environment variable `ANCHORS_TELEMETRY` | any text, including absent or blank | — (every text has an answer) | this unit: blank or absent defers to the file, any other text decides |
| the file's `telemetry:` value | any text, passed in by the caller, including empty | — (every text has an answer) | the caller reads the configuration and hands the value over; this unit only interprets it |

## Effects

| Effect | Description |
| --- | --- |
| `TLCNT-B01` | `Disabled` answers yes for the words `off`, `0`, `false` and `no` turn telemetry off, whatever their case and with surrounding spaces ignored. |
| `TLCNT-B02` | With nothing declared in the environment nor in the file, telemetry is on. |
| `TLCNT-B03` | A non-blank environment value decides alone: `on` in the environment keeps telemetry on over `off` in the file, and `off` in the environment turns it off over `on` in the file. |
| `TLCNT-B04` | With the environment blank or absent, the file's value decides with the same words. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TLCNT-I01` | The same value means the same thing in both places: a word that disables through the environment disables through the file, and one that keeps telemetry on does so in both. | evaluates each off-word and a non-off word through the environment and through the file and compares the answers |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TLCNT-X01` | Only the four off-words turn telemetry off; any other text — `on`, `yes`, `disabled`, a typo — leaves it on. | The switch is a closed list of the obvious forms, not a guess at intent; a value outside the list keeps the maintainer's default. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |

none — the decision is total: every text, blank or not, from either source, has an answer, and
nothing is read or parsed that could fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the decision reads only the process environment and the value the caller passes in.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
