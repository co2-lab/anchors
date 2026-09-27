<!-- @anchors
  code: WRPRW
  updated_at: 2026-09-26
  layer: comando
-->
# WorkPrompt — compose the work prompt of one stage over one target, from what the project declares

> **Code**: `WRPRW`

## Overview

`anchors guide` teaches the doctrine — what a spec is, how a test is written — permanent and without a
target. Whoever executes ONE stage over ONE file needs something else: what to read now and in what
order, which layer the target is in and what it demands, where the pieces of the triad are born, what
is not the stage's scope, and how to verify and record the work. Without it every orchestrator rewrote
that prompt by hand and it came out different each time.

`anchors work <artifact> --for <target>` composes that prompt, and invents nothing: it reads the layers,
the guides that govern them, the derived paths, the gates, and the layer's own extra steps from the
project's configuration. The stages are spec, code, feature and test (which produce a piece), and the
reviews — of a unit, of a whole plan, and of a plan draft — which produce findings instead.

The target is the unit, the code file; pointing at a derived piece is the predictable mistake (it is
the file that already exists), so the command redirects it and says so. When the stage's piece is one
the layer waives, or the layer is declarative and the stage is a triad piece, the prompt is only a
STOP: a production script under a heading that forbids the piece made workers create it anyway.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the artifact | spec, code, feature, test, review, review-plan, review-plan-draft, in any case | anything else | this unit: refuses naming the valid ones |
| the target | a path relative to the root or absolute; the unit or one of its derived pieces | no target | this unit: refuses with an example |
| the project | a root with a loadable configuration; a map is optional | a root without configuration | this unit: refuses; without a map it falls back to the naming convention |

## Effects

### The command

| Effect | Description |
| --- | --- |
| `WRPRW-B01` | The command refuses an artifact outside the stages it composes, a missing `--for`, and a project whose configuration does not load. |
| `WRPRW-B02` | A target in a spec, feature or test layer is redirected to its unit — the file the map says that spec specifies (also through the spec a feature or test is linked to), else the code file with the same stem beside it — with a note on standard error; a code file is not redirected, and a piece whose unit does not exist is not redirected. |

### When the stage must not happen

| Effect | Description |
| --- | --- |
| `WRPRW-B03` | When the target's layer waives the stage's piece, the prompt is only a STOP saying the layer waives it and not to create it. |
| `WRPRW-B04` | When the target's layer is declarative and the stage is spec, feature or test, the prompt is only a STOP saying the layer is recognized and has no piece of its own. |

### The prompt's sections

| Effect | Description |
| --- | --- |
| `WRPRW-B05` | The heading names the stage and the target — "Work: <artifact> of", "Review of", or "WHOLE review —" — and the role follows: producer, reviewer of the unit, or reviewer of the whole. |
| `WRPRW-B06` | The target section names the file, its layer with its regime and its tags, or says the layer is unclassified and to confirm it. |
| `WRPRW-B07` | "Read first" lists the artifact's guide, then every project guide that governs the artifact or one of the layer's tags, each once and sorted, then the target and its neighbours; without such guides it says no project guide governs the layer. |
| `WRPRW-B08` | The pieces are listed at the paths the project derives for the target's layer — the stage's own piece marked, the existing ones (under the project root) marked, overridden paths noted, waived pieces marked not to be created — with the unit's name cut at its compound suffix; without derived paths it says so, and a declarative layer's code stage lists only the file itself. |
| `WRPRW-B09` | The feature and test stages list the project's regime tags, sorted, with the surface each is proved on. |
| `WRPRW-B10` | Each stage states what is NOT its scope. |
| `WRPRW-B11` | Each stage has its procedure; a declarative layer gets the declarative procedure, and the layer's own extra steps for the stage are appended. |
| `WRPRW-B12` | "What the gates will demand" lists, once per check, the gates that apply to the stage (the test stage also owes the feature-test match), marks informative ones, and omits gates it has no description for. |
| `WRPRW-B13` | A review confronts the pieces that exist and the pending delivery records of the unit, sorted; without one it says the record is in the card's comments in github mode, and that there is no record otherwise. |
| `WRPRW-B14` | The open findings recorded for the unit (to do and in progress) are listed before writing: an issue counts when one of the paths in its name is a piece of the unit — the unit's stem followed by the piece's suffix — so a unit whose name only contains the stem (`src/pricing-v2.ts`, `lib/src/pricing.ts`) is another unit. |
| `WRPRW-B15` | The test stage and the unit review explain the execution signals: ingesting the suite's report and the unit's mutation report. |
| `WRPRW-B16` | Producing stages end with how to record the delivery with `anchors deliver`; reviews end with how to close the review with `anchors judge`, including findings in another unit. |
| `WRPRW-B17` | The verification runs the checks over the piece this stage produces — not over a code file that may not exist yet — first without recording, then once more recording. |
| `WRPRW-B18` | The legitimate waivers are listed for the stage: `@no-mark` and `@no-scenario` for a spec, `@no-paginate` and `@allow-boundary` for code. |
| `WRPRW-B19` | When the ruler does not decide, the spec stage is told to write the question in the spec's open-decisions section, or its "none" value; the other stages are told to read that section. |
| `WRPRW-B20` | The code stage requires marking each rule in the code when the project's rule-marking policy is required, and only suggests it otherwise. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `WRPRW-I01` | A prompt never both records a delivery and closes a review. | every stage is composed, and the producing stages carry only the record section while the reviews carry only the close section |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `WRPRW-X01` | Composing the prompt writes nothing to the project. | The prompt describes work for someone else; a prompt that created files would do part of the stage behind the worker's back. |
| `WRPRW-X02` | Wherever the prompt cites the open-decisions section it uses one title and one "none" value, the catalogue's; and it cites configuration keys and checks by their current names (`optional_triad_edges`, `rule-implemented`, `tests-pass`). | A prompt that names the same section two ways, or a key no configuration accepts, sends the worker to write what no gate reads. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `WRPRW-E01` | REF[WRPRW-B01]: an unknown artifact, a missing target or an unloadable configuration, which B01 refuses before composing anything | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Layer`, `Derived`, `Gate` | config — the declared structure |
| DEP2 | `internal/queue` | `ValidWorkArtifact`, `ArtefatosDeTrabalho` | the stages that are composed |
| DEP3 | `internal/change` | `Pending` | the pending delivery records |
| DEP4 | `internal/issue` | `List` | the recorded findings |
| DEP5 | `internal/mapx` | `Load`, `StemOfAnchor` | mapa — the specifies edges and the unit's stem |
| DEP6 | `internal/i18n` | `TIn` | the open-decisions section's title and "none" value |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
