<!-- @anchors
  code: TLEMT
  updated_at: 2026-09-26
  layer: apoio
-->
# TelemetryEmitter — decision events leave as OTLP logs, never block the work, and never carry who uses the product

> **Code**: `TLEMT`

## Overview

Sends the decision events to a collector over the OpenTelemetry logs protocol (HTTP with a JSON
body). The protocol is written by hand rather than through the official SDK, which brings dozens of
modules into a tool that installs with a single command and runs on agent machines, CI and laptops;
the gain is that the destination is swappable — the same emission goes to Honeycomb, Grafana,
Elastic or a local collector by changing the address, not the code.

Two decisions rule the unit. Telemetry never delays the work: sending happens in the background,
failures are swallowed, and waiting for in-flight sends when the process ends has its own short
deadline, so a slow or dead collector costs at most that deadline and the event is lost — losing an
event is better than making someone wait for it. And the payload identifies the PRODUCT, never who
uses it: the resource carries only the product name and version, and each attribute value is
either a number, a boolean or a string of Anchors' vocabulary; anything else is reduced to the name
of its type, so a map or an error carrying a file path cannot leak through a careless future call.

When telemetry is disabled there is no emitter at all, and the absent emitter accepts every call as
a no-op, so callers need no guard around each call.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | enabled or not, an endpoint (blank means the default) and authentication headers | — | the caller builds it from the environment and the opt-out |
| the event | a named event with attributes of any Go value | attribute values that are content (paths, errors, maps) | this unit reduces any value that is not a string, an integer or a boolean to its type description |
| the collector | any OTLP/HTTP logs endpoint, reachable or not | — | this unit tolerates any answer, or none |

## Effects

| Effect | Description |
| --- | --- |
| `TLEMT-B01` | `NewEmitter` with a disabled configuration yields no emitter, and emitting or flushing through the absent emitter does nothing and does not crash. |
| `TLEMT-B02` | A blank endpoint means the default, Honeycomb's logs endpoint. |
| `TLEMT-B03` | Each event is sent as one OTLP log record inside one resource: the record's body is the event name, its time is the event's instant in Unix nanoseconds, and its attributes are the event's attributes. |
| `TLEMT-B04` | Sending happens in the background: emitting returns at once, and a collector that is down neither blocks nor fails the caller. |
| `TLEMT-B05` | Flushing waits for the sends still in flight, so an event emitted just before the process ends reaches the collector. |
| `TLEMT-B06` | Flushing waits at most two seconds, on a deadline of its own that does not depend on the HTTP client's timeout. |
| `TLEMT-B07` | Each request is a POST with a JSON content type and carries the authentication headers of the configuration. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TLEMT-I01` | No attribute value leaves as arbitrary content: strings, integers and booleans keep their typed form, and any other value is sent as its type description in parentheses. | converts a map holding a path and verifies only the type description comes out, while an integer and a vocabulary string pass typed |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TLEMT-X01` | The resource carries exactly two attributes, the product name `anchors` and its version — never a host, user or repository. | The first attribute someone added for convenience would identify who uses the product, and nothing else would accuse it. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TLEMT-E01` | REF[TLEMT-B04]: an unreachable collector is the failure B04 answers: the send fails in the background and the caller is never told | — | — |
| `TLEMT-E02` | The collector answers with an error status. | The answer is discarded; nothing is retried and nothing fails. | Telemetry that fails cannot become a problem for whoever is working; the event is lost in silence. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/telemetry/config.go` | `Config` | apoio — whether to emit, where, and with which headers |
| DEP2 | `internal/telemetry/event.go` | `Event` | apoio — the decision event being sent |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
