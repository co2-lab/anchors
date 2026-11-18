<!-- @anchors
  code: TLNTT
  updated_at: 2026-09-26
  layer: apoio
-->
# TelemetryNotice — the telemetry notice reaches whoever did not ask for it, once, and says how to turn it off

> **Code**: `TLNTT`

## Overview

Telemetry is on by default, and what makes that choice honest is that the notice reaches people who
never asked for it — including someone who installed the binary into a project another person
configured, and so never ran `init`. The caller shows it before any command runs, so it appears
before the first event is sent.

The text says three things, in this order: what is collected (decision events — how many
candidates a claim saw, which gate failed, in which state a turn ended), what is NOT collected (file
content, specs, diffs, error messages, repository, user or branch names), and how to turn it off
(the environment variable, or the configuration key). Someone who reads "we collect data" and does
not find the way out on the same screen assumes the worst.

It is shown once. A marker file under the project's `.anchors/` directory, which is not
versioned, records that the notice was shown there, and later calls stay silent. Repeating the
notice every time would be worse than not showing it: whoever reads the same thing every time stops
reading. Failing to record the marker never stops a command — at worst the notice shows again.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the destination of the text | any writer; the caller passes the error stream | — | the caller |
| the project root | a directory, possibly one where `.anchors/` cannot be created | — (an unwritable root is tolerated, see the Errors) | the caller resolves the root; this unit tolerates any |

## Effects

| Effect | Description |
| --- | --- |
| `TLNTT-B01` | `Notice`: The first call writes the notice; a later call for the same project root writes nothing, because the first one left its marker. |
| `TLNTT-B02` | The text states what is sent, what is not sent — naming file content — both ways to turn telemetry off: the `ANCHORS_TELEMETRY=off` environment variable and `telemetry: off` in the configuration — and that it appears once per project on this machine, which is where its marker lives. |
| `TLNTT-B03` | The marker, recorded by `MarkNoticed`, is written at `.anchors/telemetry-noticed` under the project root, a directory that is kept out of version control, so one person seeing the notice does not silence it for the team. |
| `TLNTT-B04` | The text comes from the translation catalog, so it is written in the project's language. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TLNTT-I01` | Showing the notice and asking whether it was shown agree: after a call that wrote the notice, `AlreadyNoticed` answers yes for that root. | calls the notice on a fresh root, then asks whether it was already shown |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TLNTT-X01` | Does not decide whether telemetry is on: it only announces it; the caller shows it only when telemetry is enabled. | Announcing is a separate duty from the opt-out, which belongs to the configuration unit. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TLNTT-E01` | The marker cannot be written (the `.anchors` directory cannot be created under the root). | The notice is still written, nothing fails, and the next call shows it again. | Showing the notice twice is better than blocking a command over an announcement. <!-- @resilient: showing the notice again on the next call is the whole cost, and blocking a command over an announcement would be worse --> |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the notice writes to the writer it receives and to a marker file under the given root.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
