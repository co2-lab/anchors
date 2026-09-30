<!-- @anchors
  code: CLRTC
  updated_at: 2026-09-30
  layer: comando
-->
# CliRoot — every command passes through one root that speaks the project's language and honours the freeze

> **Code**: `CLRTC`

## Overview

The root assembles the CLI from its domains (governance, map, quality, flow, operations) and is
the one point every command passes through before it runs. Two cross-cutting duties live
there, because spreading them over the commands guarantees the next command is born without
them.

The first is the language. The project's top-level `lang:` is applied before the command
prints anything, so even the lines a command prints before loading its configuration come out
in the project's language. The read is deliberately loose: a configuration that does not
parse, or declares a language Anchors does not have, must not stop the command — the loader
reports those with the line and the key.

The second is the freeze. When the project declares itself frozen, every command that
produces project state is refused, naming the command and the reason. Reading stays allowed,
because whoever investigates the problem needs the status, the doctor and the guides; `thaw`
and `freeze` stay allowed, or the freeze could not be undone by Anchors itself. The allowance
covers every subcommand of an allowed command. A missing or broken configuration is not a
freeze: answering "frozen" there would send whoever investigates the wrong way. The root also
leaves printing to the entry point, which prints the error once and decides the exit code.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | the command's `--root` when it has one, otherwise the working directory | — | this unit |
| the configuration | a loadable `anchors.yaml`, possibly frozen | a missing or broken one | this unit: treated as not frozen |
| the language | a top-level `lang:` Anchors supports | a nested `lang:`, an unsupported value, an unparsable file | this unit: ignored, the default language stays |

## Effects

| Effect | Description |
| --- | --- |
| `CLRTC-B01` | In a frozen project, a command outside the allowance is refused with a message naming the command and the freeze reason. |
| `CLRTC-B02` | While frozen, thaw, freeze, status, doctor, guide, help, completion, coverage and impact still run. |
| `CLRTC-B03` | A subcommand of an allowed command runs while frozen, because every command up the parent chain is checked. |
| `CLRTC-B04` | A missing configuration, or one that does not load, is not frozen: the command proceeds. |
| `CLRTC-B05` | The project read for the freeze and the language is the one `--root` points at, when the command has that flag. |
| `CLRTC-B06` | A top-level `lang:` of the project is applied before the command runs, even when the rest of the file does not parse; a nested or unsupported `lang:` is ignored. |
| `CLRTC-B07` | The root prints neither the error nor the usage of a failing command. |
| `CLRTC-B08` | The project's `lang:` is read from a line ending in `\r\n` as from one ending in `\n`. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CLRTC-I01` | Every `anchors <command>` that the work guide and the pipeline workflows teach is a registered command. | collects the registered names and confronts every command the guide and the workflows cite |
| `CLRTC-I02` | Every command that takes several files — a usage that names files and repeats them, such as `<file>...` — reads them as every other does, through `TakesFiles`: separate arguments or a comma-separated list. | walks every registered command and confronts each whose usage takes several files with the mark `TakesFiles` leaves |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CLRTC-X01` | The freeze refusal is written in the project's language, not in the default one. | The refusal is the one message everyone in the project must read during a freeze. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CLRTC-E01` | REF[CLRTC-B01]: a command run in a frozen project is the failure the root raises | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `AbsRoot`, `Load`, `Frozen`, `FreezeReasonText` | config — the freeze state |
| DEP2 | `internal/i18n` | `Set`, `T` | apoio — the language |
| DEP3 | `cmd/anchors/common` | `NoticeTelemetry` | comando — the telemetry notice every command passes through |
| DEP4 | `cmd/anchors/ops/register.go` | `Register` | comando — OPRGP, and the other domains' registrations |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
