<!-- @anchors
  code: RFRSR
  updated_at: 2026-09-26
  layer: gate
-->
# RefResolves — the reference points at the spec that REALLY describes the unit

> **Code**: `RFRSR`

> **RFRSR-R0001:** when no sibling spec exists on disk, the gate now consults the project
> graph to prevent references to invented or phantom identities. Previously, absence of a sibling
> spec caused the gate to stay silent (Skip), which left 42% of references uninspected and allowed
> phantom identities through as undetermined. An infra file with no sibling spec remains
> valid (Skip) only if the referenced identity actually exists in the project; references
> to codes declared nowhere in the project fail.
>
> **Revises:** `I01`, `X03`, `B14`, `B15`, `B16`, `B17`
> **Checked:** `B01`, `B02`, `B03`, `B04`, `B05`, `B06`, `B07`, `B08`, `B09`, `B10`, `B11`, `B12`, `B13`, `I02`, `X01`, `X02`, `X04`, `E01`

## Overview

Confronts an artifact against the spec co-located with it: **the reference field is
filled, but does it name the right owner?**

The header gate verifies that the field EXISTS. Nobody verified that it points at the
right place — and a wrong reference is worse than a missing one: it looks like
traceability, the gate goes green, and the whole unit is attributed to the wrong spec.
Every relational gate that depends on that edge starts confronting the wrong pair, in
silence.

The characteristic failure mode is the REFACTORING nobody propagated. Measured on a real
project: 49 model files still referencing the identity from back when all the models lived
in a single file. After the split, each one got its own spec, and not one reference was
updated. The 49 kept pointing at the whole schema, and nothing raised a hand.

**The ruler**: if a SIBLING spec exists — the one the Structure co-locates with this file —
the reference must be that spec's own identity. With no sibling spec the gate goes quiet:
charging the absence of the piece is the triad gate's job, and two gates accusing the same
defect become noise.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the node kind | code, feature or test — the kinds that REFERENCE | a spec, which OWNS an identity instead of citing one | this unit: a kind outside the set leaves without a verdict |
| the artifact's content | any text, empty included | — (text with no reference is a case, not an error) | this unit: the absence belongs to the header gate, and is skipped here |
| the project root | a readable directory path | — (an unreadable sibling is read as absent) | this unit: a read error means "no sibling spec", never a failure |
| the code length | the lengths the project declares in its Structure | — | `internal/config/config.go`, which serves the pattern at call time |

## Effects

| Effect | Description |
| --- | --- |
| `RFRSR-B01` | A reference that does not match the sibling spec's identity FAILS — the refactoring that nobody propagated. |
| `RFRSR-B02` | The failing verdict names BOTH sides: what is written and what the sibling spec declares, so whoever reads it does not have to open two files. |
| `RFRSR-B03` | A reference equal to the sibling spec's identity passes. |
| `RFRSR-B04` | A spec is not confronted: it OWNS an identity, it does not reference one. |
| `RFRSR-B05` | An artifact with no reference declared leaves without a verdict — the absence is the header gate's charge, not this one's. |
| `RFRSR-B06` | With no sibling spec on disk the gate goes quiet: the missing piece is the triad gate's charge, and two gates on one defect become noise. |
| `RFRSR-B07` | A sibling spec that exists and declares no identity of its own counts as no sibling: there is nothing to compare against. |
| `RFRSR-B08` | The sibling is found by NAME convention — same stem, same directory, spec suffix. |
| `RFRSR-B09` | The test suffix of each supported language is stripped before the stem is computed, so a test file finds the same sibling its code does. |
| `RFRSR-B10` | An intermediate extension is dropped from the stem too, so a file that carries a kind in its name still lands on the unit's spec. |
| `RFRSR-B11` | The reference is read from the header whatever the comment syntax of the language — the three families of line marker are accepted. |
| `RFRSR-B13` | A leading dot is not a stem separator: a hidden file keeps its whole name, so it is never attributed to a spec named only by the suffix. |
| `RFRSR-B12` | The accepted identity length comes from the project's Structure, read at confrontation time and not frozen at process start. |
| `RFRSR-B14` | A reference to a code that exists nowhere in the project fails when no sibling spec exists, reporting the unknown code. |
| `RFRSR-B15` | Without a sibling spec on disk, an existing declared code in the graph skips. |
| `RFRSR-B16` | An inferred identity (`CodeDeclarado: false`) does not satisfy the reference. |
| `RFRSR-B17` | Without a graph, absence is not asserted and the gate skips. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RFRSR-I01` | The sibling spec on disk is the primary ruler. When a sibling spec exists, the verdict depends on the filesystem alone. When no sibling spec exists, the graph is consulted solely to verify that the cited reference exists as a declared identity. | confronts the same pair with no graph and verifies the verdict |
| `RFRSR-I02` | The gate never writes and never repairs. It reads the artifact and the sibling and returns a verdict; a gate that fixed what it points at would pass on the second run. | confronts a divergent pair and verifies the artifact on disk is untouched |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RFRSR-X01` | Does not charge the ABSENCE of the reference field. | That is the header gate's ruler, and it already states it. Two gates on the same defect produce two messages for one fix — the reader turns both off. |
| `RFRSR-X02` | Does not charge the absence of the sibling spec. | The missing piece of a triad is the triad gate's charge. Here the absence is simply a case where there is nothing to compare. |
| `RFRSR-X03` | Does not consult the map when a sibling spec is present on disk. | The co-located file takes precedence so stale graph builds cannot override the filesystem truth. When no sibling spec exists, the map is queried solely to prevent references to phantom identities. |
| `RFRSR-X04` | Does not judge whether the sibling spec DESCRIBES the unit well. | The ruler here is identity: which spec owns this file. Whether the spec's content matches the code is judgement, and judgement belongs to another class of gate. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RFRSR-E01` | REF[RFRSR-B06]: a sibling spec that cannot be read is answered as the missing sibling of B06: the gate goes quiet | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindCode` | core — the kind decides whether the artifact references or owns |
| DEP2 | `internal/config/config.go` | `Config` | core — the project's Structure travels with the confrontation, and the accepted identity length is read from it per call |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
