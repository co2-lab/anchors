<!-- @anchors
  code: PCSDP
  updated_at: 2026-09-26
  layer: infra
-->
# PackSeeding — copy the compliance packs carried in the binary into the project, never over an adapted one

> **Code**: `PCSDP`

## Overview

The compliance packs (privacy, payment, health, accessibility) travel inside the binary and are COPIED
into the project's `packs/` folder by `init`; they are not read from inside the binary. The difference
matters: in the project the pack is material, versioned and auditable — whoever must justify compliance
opens the file and reads the duty next to the article it comes from. It also allows what the law demands
in practice: adapting. A norm changes, a project has its own legal reading, a market has an extra demand —
the pack in the project is edited; a pack living only in the binary would have to wait for a release.

Every pack is copied, not only the ones of the jurisdiction the project declared: a project that opens a
new market tomorrow finds the pack ready and declares one line. What filters is adoption, declared in the
configuration, not presence on disk. And a pack that already exists is never overwritten: the project may
have adapted it, and overwriting someone's adaptation is the silent loss this framework exists to prevent.

The unit also lists the packs it carries, grouped by domain, so `init` can ask about them and the
compliance command can say what exists but was not adopted.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | an existing folder the process may write to | a folder where `packs/` cannot be created | this unit: the write failure is returned |
| the packs carried in the binary | the pack files grouped in one folder per domain | — | the build, which embeds the pack folder |

## Effects

| Effect | Description |
| --- | --- |
| `PCSDP-B01` | Seeding (`SeedPacks`) copies every pack the binary carries to `packs/<domain>/<file>` in the project, byte for byte, whatever jurisdiction the project declared. |
| `PCSDP-B02` | A pack file that already exists in the project is left untouched and reported as preserved, not as created. |
| `PCSDP-B03` | The lists of created and of preserved packs come back sorted. |
| `PCSDP-B04` | The list of available packs (`AvailablePacks`) groups them by domain, each named `<domain>/<name>` without the file extension, and each group is sorted. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PCSDP-I01` | Seeding twice is the same as seeding once: the second run creates nothing and reports every pack as preserved. | seeds a folder, seeds it again, and compares the two reports |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PCSDP-X01` | Does not decide which packs the project adopts; presence on disk is not adoption. | Adoption is a line in the configuration; copying every pack keeps the next market one line away. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PCSDP-E01` | The pack folders or files cannot be written in the project. | The seeding returns the write error. | A silent partial seeding would let `init` report packs that are not on disk. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the unit reads only the packs the binary carries and writes only the project folder.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
