<!-- @anchors
  code: TLSTT
  updated_at: 2026-10-08
  layer: comando
-->
# TelemetrySetup — every command starts telemetry the same way: the opt-outs first, then the notice, then the emitter

> **Code**: `TLSTT`

## Overview

Telemetry is on by default, and what makes that honest is that nobody is counted without first being
told, and that turning it off is always in reach. This unit is the start-up step every command runs before
its own work. It finds the project root, reads the project's declared opt-out from that root's
configuration, honours the environment's opt-out over it, shows the notice the first time, and only then
builds the emitter that later commands use to send decision events.

The project root is the one the person asked for with `--root`, taken as given, or, with no `--root`, the
nearest directory above the working directory that holds the project configuration. Reading the opt-out
from that root, not from the working directory, is what makes `telemetry: off` hold when a command runs
from a subdirectory or points at another project.

The emitter authenticates only with a key given in the environment; the repository never holds one.
The emitter is built with unit codes allowed: the emitter can drop them (`NoCodes`), but no project
setting declares that yet, so this step never asks for it.
At exit, the command waits for the events still in flight, and does nothing when telemetry never started.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the `--root` flag | an explicit directory, or empty / the default | — | this unit: explicit is respected, default walks up |
| the project's opt-out | the `telemetry:` field of the root's configuration | — | the telemetry unit decides which values mean off |
| the environment opt-out | `ANCHORS_TELEMETRY` | — | the telemetry unit: when set, it overrides the project |
| the ingestion key | `ANCHORS_TELEMETRY_KEY` | a key committed anywhere in the project | this unit: reads it from the environment only |

## Effects

| Effect | Description |
| --- | --- |
| `TLSTT-B01` | A project's `telemetry: off` is read from the project root's configuration, so it holds from the root, from a subdirectory, and with `--root` pointing at the project from elsewhere. |
| `TLSTT-B02` | With telemetry turned off by the environment or by the project, no notice is shown, no emitter is built and nothing is written in the project. |
| `TLSTT-B03` | `NoticeTelemetry`: with telemetry on, the notice is shown once, saying how to turn it off, the fact that it was shown is marked in the project, the emitter is built, and the next run shows no notice. |
| `TLSTT-B04` | The emitter carries an authentication header only when `ANCHORS_TELEMETRY_KEY` is set, and then only that header. |
| `TLSTT-B05` | `ProjectRoot`: the project root is the explicit `--root` as given; without it, the nearest directory above the working directory that holds the project configuration, and empty when no directory above holds one. |
| `TLSTT-B06` | `FlushTelemetry`: waiting for events in flight at exit returns at once when no emitter was built. |
| `TLSTT-B07` | When the project configuration does not load, a top-level `telemetry:` line is still read from its text, so a declared `telemetry: off` still builds nothing. |
| `TLSTT-B08` | Outside a project — no root found, or a root without the project configuration — no notice is shown, no emitter is built and no mark is written anywhere. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TLSTT-I01` | No emitter exists unless the notice was shown in this project at some point: turning telemetry off stops both together. | runs the start-up with each opt-out and checks there is neither notice, nor emitter, nor mark; runs it with none and checks all three exist |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TLSTT-X01` | Does not send any event while starting; it only builds the emitter. | Events are the commands' business; the start-up must never reach the network. |

## Errors

none — a project configuration that does not load is not an error here: its `telemetry:` line is read as text (see TLSTT-B07), and the command itself reports the configuration error when it loads it.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
