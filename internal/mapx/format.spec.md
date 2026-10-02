<!-- @anchors
  code: MPFRM
  updated_at: 2026-10-01
  layer: mapa
-->
# MapFormat — the map's format number decides whether this binary may read it

> **Code**: `MPFRM`

## Overview

The map file carries a format number, and that number is the contract between the file and the binary
that reads it. It says in which shape the file is written, not which release of Anchors wrote it. The
binary writes exactly one format and reads a closed range of formats; anything outside the range is
refused before a single field is interpreted.

The refusal exists because partial reading loses data without a trace: a binary that meets an unknown
shape keeps the fields it recognises, drops the rest, and the next save writes only what survived. The
judgment stamps are the concrete case — renamed keys would vanish and `check` would ask again what a
person had already answered.

The two directions of refusal are told apart, because they send the reader to opposite remedies. A file
from the future asks for a newer binary; a file from the past asks for a migration. A single "could not
load the map" for both would send the person looking for corruption where there is only a version.

The format goes up only when a change makes the file unreadable to the previous version (a renamed key, a
changed value shape, a removed required field) — never for a new optional field, which an old binary simply
ignores.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the format found in the file | any integer, including zero for a file with no version | nothing: every integer gets an answer | this unit: zero is read as format 1, every other value is compared to the readable range |
| the path of the file | any text, used only to name the file in the message | — | the caller passes the path it loaded |

## Effects

| Effect | Description |
| --- | --- |
| `MPFRM-B01` | Checking the format (`ConfereFormato`), a map in the format this binary writes, or in the oldest format it still reads, is accepted with no error. |
| `MPFRM-B02` | A map in a format newer than the one this binary writes is refused with a message that says the map was written by a newer Anchors, that continuing would lose what it wrote in silence (judgment stamps included), and how to upgrade. |
| `MPFRM-B03` | A map older than the oldest readable format is refused with a message that names `anchors migrate` as the fix. |
| `MPFRM-B04` | A map with no version is format 1: it is refused as needing migration, and the message says format 1, never format 0. |
| `MPFRM-B05` | The two refusals do not mix: the refusal of a newer map never names the migration command. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPFRM-I01` | The binary writes format 6 and reads only format 6 — format 6 renamed the gate `triad-complete` to `unit-complete`, and a project not migrated would declare a gate this binary does not know, its pieces charged by nothing; every other format is refused. | walks every format from below the range to above it and checks accept or refuse for each |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPFRM-X01` | Format 1 is not read, even though earlier binaries wrote it: it is migrated. | Its keys were renamed; reading both names forever would mean the file never repairs itself. |

## Errors

Each failure the unit reports is already stated as a behaviour; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPFRM-E01` | REF[MPFRM-B02]: a map newer than the binary is the failure B02 answers with the upgrade message | — | — <!-- @resilient: the newer map is refused with a typed error carrying the upgrade message, which the loader returns and the command prints --> |
| `MPFRM-E02` | REF[MPFRM-B03]: a map older than the readable range is the failure B03 answers with the migration message | — | — <!-- @resilient: the older map is refused with a typed error carrying the migration message, which the loader returns and the command prints --> |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the check uses only the standard library.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
