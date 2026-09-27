<!-- @anchors
  code: ININT
  updated_at: 2026-09-26
  layer: comando
-->
# InitNonInteractive — the init that an agent answers with flags: it asks in JSON, and writes only a complete, valid set of answers

> **Code**: `ININT`

## Overview

An agent cannot drive the interactive init: it has no terminal to answer prompts on. With
`--non-interactive`, `init` becomes a two-step exchange in JSON. Called with no answer, it
infers what it can from the disk and returns the questions — with their options and inferred
defaults, whether the project still needs the DISCOVER phase, and how to answer — and writes
nothing, because writing the defaults of a new project would produce a configuration that
governs nothing. Called again with the answers as flags, it validates them all and writes the
configuration.

An answer counts only when its flag was actually given, so a deliberate "no" or an empty
value is honored rather than replaced by a default; `--defaults` is the explicit way to
accept every inferred answer. One invalid answer refuses the whole set, because writing the
valid ones would produce a file nobody fully decided. A stack preset fills the code layers
from the stack's structure, and the code layers are pruned only when the agent chose them.
The success document names the written file, echoes the accepted answers, and says what comes
next: the DISCOVER phase when the project has nothing to learn from yet, otherwise building
the map.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the answers | the flags of the init questions, each counted only when given | answers the question validation refuses (a github workflow with no repository, an unknown option) | the init package validates them; this unit writes nothing when any is refused |
| the governs rules | `GUIDE=tag1,tag2` | a rule with no `=` | this unit: refuses with the expected form |
| the project root | any directory, empty or with code | — | the init package infers what it can |

## Effects

| Effect | Description |
| --- | --- |
| `ININT-B01` | With no answer and no `--defaults`, the command prints the questions as JSON with `escrito: false`, whether the DISCOVER phase is needed, and how to answer, and writes nothing. |
| `ININT-B02` | The given answers reach the configuration: the chosen artifact layers, co-location, the default gates, the manual workflow, and the header guide in `guides/` when no guide directory was detected. |
| `ININT-B03` | A flag given as false or empty is a deliberate answer: `--gates=false` seeds no gate and `--header=false` writes no header guide. |
| `ININT-B04` | The github workflow writes the mode with the repository and the labels. |
| `ININT-B05` | One refused answer refuses the whole set: the command prints JSON with `escrito: false`, every answer's status and the error, fails, and writes nothing. |
| `ININT-B06` | `--defaults` with no other answer writes the configuration from the inferred defaults. |
| `ININT-B07` | A stack preset fills the code layers of that stack, whether or not the project already has code. |
| `ININT-B08` | `--layers` keeps only the chosen code layers and prunes the others; without it, no layer is pruned. |
| `ININT-B09` | The success JSON has `escrito: true`, the written file, the answers' status and the next step: the DISCOVER phase when the project needs it, otherwise `anchors map build`. |
| `ININT-B10` | Each `--governs GUIDE=tag1,tag2` rule is written to the configuration as one governs rule per non-blank tag, guides in sorted order. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `ININT-I01` | Either every answer is accepted and the file is written, or nothing is written. | refuses one answer of a set and checks no file exists; accepts a set and loads the file |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ININT-X01` | The non-interactive mode never prompts: every exchange is JSON on standard output. | It exists for callers with no terminal; a prompt would hang or fail them. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ININT-E01` | A `--governs` rule has no `=`. | The command fails naming the expected `GUIDE=tag1,tag2` form. | A rule without its tags cannot be read either way, and guessing would govern the wrong guide. |
| `ININT-E02` | REF[ININT-B05]: a refused answer is the failure this mode handles, by refusing the set and writing nothing | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx` | `Infer`, `Questions`, `ValidateAnswers`, `TudoAceito`, `PrecisaDescobrir`, `ApplyPreset`, `ApplyArtifactChoice`, `ApplyColocation`, `PruneCodeLayers`, `DefaultGates`, `RenderHeaderGuide` | apoio — the init questions and how each answer shapes the config |
| DEP2 | `internal/config/config.go` | `Save`, the workflow modes | config — the written file |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
