<!-- @anchors
  code: TLEVT
  updated_at: 2026-09-26
  layer: apoio
-->
# TelemetryEvent — a decision event has a name from a closed vocabulary, a caller-stamped instant and attributes

> **Code**: `TLEVT`

## Overview

What an agent decides is invisible once the command ends. Measured in a session with six agents:
three took the same issue and nobody knew until three identical pull requests appeared; one ended
its turn with the card still in progress and the work stood still for hours. None of these is an
error — every command returned success. They are decisions that, seen in sequence, form a pattern,
and this unit is the shape in which each decision is recorded.

An event has a name, an instant and attributes. The name comes from a CLOSED vocabulary of five
decisions — a claim served a card, a claim found no work, a check finished, work was escalated to a
person, a turn ended — because an open set would become free text, and free text is where a leak
enters without anyone noticing. The attributes are meant to be numbers and Anchors' own vocabulary
(how many candidates, which gate, which state), never content.

The unit does not invent time: whoever creates the event passes the clock, so the instant is the
caller's and a test can pass a fixed clock instead of sleeping.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the event name | one of the five names of the vocabulary | any other text | the callers, which only use the declared names; the type is a named string, so a raw literal is still accepted by the language |
| the attributes | a key-value set, or none at all | free text, file content, paths | the callers; the emitter additionally turns any value that is not a number, a boolean or a string into its type description |
| the clock | a function returning the instant | — | the caller |

## Effects

| Effect | Description |
| --- | --- |
| `TLEVT-B01` | An event created by `New` carries the name and the attributes it was given, unchanged. |
| `TLEVT-B02` | An event created without attributes carries an empty attribute set, never an absent one. |
| `TLEVT-B03` | The event's instant is the one returned by the clock the caller passed. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TLEVT-I01` | The vocabulary is exactly five names, each with its fixed wire value: `claim.served`, `claim.empty`, `check.finished`, `escalate.raised` and `turn.ended`. | compares every declared name against its wire value |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TLEVT-X01` | Does not read the system clock: the instant comes only from the clock passed in. | The package does not invent time; a fixed clock makes the event reproducible. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |

none — creating an event only assembles values the caller passed; nothing is read, parsed or
sent here.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the event is a plain value; sending it is the emitter's duty.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
