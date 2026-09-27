<!-- @anchors
  code: NWARN
  updated_at: 2026-09-26
  layer: comando
-->
# NewArtifact — a new artifact is born beside its unit, with a resolved identity and the sections of the project's ruler

> **Code**: `NWARN`

## Overview

`anchors new <kind> <name> --out <path>` emits the skeleton of an artifact with its header and
identity already resolved, so nobody writes a header or picks a code by hand. The artifact is
born where `--out` says — beside the unit it describes, never at the repository root — and an
existing file is never overwritten.

A spec owns an identity: it gets a code no unit in the map uses (with the module prefix of its
layer when the target path is in one). A feature or a test references the spec instead: it
reads the code from the sibling spec's header, and without one it still gets a code, with a
warning that it is born orphaned — the likeliest cause is a wrong `--out`. `--code` pins the
identity by hand. A spec for a unit in a declarative layer is refused, because such layers
originate no rule and are declared precisely to leave scrutiny.

The sections come from the kind's catalog: its defaults, adjusted by `--with` and `--without`,
or a spec preset that fixes the set and the reading order. The text follows the project: a
section title comes from the target layer's titles, then the project's titles, then the rule
type of its letter, then the project's language; the bodies are translated too; the feature
keywords and the test syntax follow the project's declared dialect; the unit regime tag comes
from the project's mapping, or a visible placeholder asks for it. A plan is born with its
progress companion.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the kind | a kind of the template catalog, in any case | any other word | this unit: refuses it |
| the name | any name | none given | this unit: refuses it |
| the output path | a path, relative to the root or absolute | none given, or a file that already exists | this unit: refuses it |
| the sections | keys of the kind's catalog | unknown keys | this unit: refuses them |
| the preset | a spec preset name, for a spec | an unknown preset, or any preset for another kind | this unit: refuses it |

## Effects

| Effect | Description |
| --- | --- |
| `NWARN-B01` | An unknown kind is refused. |
| `NWARN-B02` | Without a name, or without `--out`, the command is refused and asks for it. |
| `NWARN-B03` | A spec is born with a code no node of the map holds; a target path in a layer with a module prefix gets that prefix, a generic basename takes the parent directory's identity, and a bare name gets its canonical code. |
| `NWARN-B04` | A feature or test takes its code from the header of the sibling spec with the same stem, and says where it read it; with no sibling spec it gets a generated code and a warning that it is born orphaned. |
| `NWARN-B05` | `--code` pins the identity written in the header. |
| `NWARN-B06` | `--with` adds and `--without` removes sections from the kind's defaults; an unknown section key is refused. |
| `NWARN-B07` | `--preset` applies only to specs and must name a known preset; it fixes the sections and their order, and a section added with `--with` comes after the preset's. |
| `NWARN-B08` | A spec whose target lies in a declarative layer is refused, and nothing is written. |
| `NWARN-B09` | A section title comes from the target layer's titles, then the project's titles, then the rule type of the section's letter, then the project's language (English by default); section bodies follow the language too. |
| `NWARN-B10` | The feature's language line and keywords follow the project's declared Gherkin language, and the test body follows the project's declared language family. |
| `NWARN-B11` | The unit regime tag is the first project tag mapped to the unit regime; with none, a visible placeholder asks for the mapping. |
| `NWARN-B12` | The artifact is written at `--out`, the output names its identity and the check to run next, and a second `new` onto the same path is refused. |
| `NWARN-B13` | A plan is born with its progress companion beside it. |
| `NWARN-B14` | `--list-sections` prints the kind's sections, marking defaults and optional ones, and the presets only for specs. |
| `NWARN-B15` | The target layer of a spec or feature is the layer of the unit it describes: an existing code file wins, otherwise the most specific extension that a layer claims. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `NWARN-I01` | A refused `new` leaves nothing behind in the tree. | runs every refusal against an empty root and lists it afterwards |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `NWARN-X01` | A feature or test never owns an identity: its header references the spec's code. | Minting a new code for a feature produced orphans by construction, a `ref` pointing at no spec. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `NWARN-E01` | REF[NWARN-B12]: an output path that already exists is refused, never overwritten | — | — |
| `NWARN-E02` | REF[NWARN-B08]: a spec for a declarative layer is refused before any file exists | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/ops/new_templates.go` | `templates`, `specPresets` | comando — NWTMN |
| DEP2 | `cmd/anchors/ops/code.go` | `takenCodes`, `unitName` | comando — CDCMC |
| DEP3 | `internal/code` | `GenerateUnique`, `GenerateUniqueWithPrefix`, `GenerateFromPath` | apoio — the code algorithm |
| DEP4 | `internal/config/config.go` | `SectionTitle`, `RuleTypes`, `DialectFor`, `Derived.Regimes` | config — the project's lexicon and dialect |
| DEP5 | `cmd/anchors/flow` | `WriteInitialProgress` | comando — the plan's progress companion |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
