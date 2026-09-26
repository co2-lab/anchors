<!-- @anchors
  code: RFRSR
  updated_at: 2026-09-26
  layer: gate
-->
# RefResolves — the reference points at the spec that REALLY describes the unit

> **Code**: `RFRSR`

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

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RFRSR-I01` | The ruler is the SIBLING spec on disk, never the map. A reference is confronted against the file the convention co-locates, so the gate answers the same with a graph and without one. | confronts the same pair with no graph and verifies the verdict |
| `RFRSR-I02` | The gate never writes and never repairs. It reads the artifact and the sibling and returns a verdict; a gate that fixed what it points at would pass on the second run. | confronts a divergent pair and verifies the artifact on disk is untouched |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RFRSR-X01` | Does not charge the ABSENCE of the reference field. | That is the header gate's ruler, and it already states it. Two gates on the same defect produce two messages for one fix — the reader turns both off. |
| `RFRSR-X02` | Does not charge the absence of the sibling spec. | The missing piece of a triad is the triad gate's charge. Here the absence is simply a case where there is nothing to compare. |
| `RFRSR-X03` | Does not consult the map to resolve the reference. | The convention that co-locates spec and code is the ruler, and it is legible from the filesystem alone. Going through the graph would make the verdict depend on a build that may be stale, and stale is exactly the defect this gate exists to catch. |
| `RFRSR-X04` | Does not judge whether the sibling spec DESCRIBES the unit well. | The ruler here is identity: which spec owns this file. Whether the spec's content matches the code is judgement, and judgement belongs to another class of gate. |

## Errors / Failures

| Rule | Condition | Effect |
| --- | --- | --- |
| `RFRSR-E01` | Underlying I/O or parsing failure | Returns Skip or Pending with error description | @no-scenario: error paths are handled by returning early verdict without panic @resilient: returns early without panic |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindCode` | core — the kind decides whether the artifact references or owns |
| DEP2 | `internal/config/config.go` | `Config` | core — the project's Structure travels with the confrontation, and the accepted identity length is read from it per call |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
