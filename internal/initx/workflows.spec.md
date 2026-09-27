<!-- @anchors
  code: FLWRF
  updated_at: 2026-09-26
  layer: infra
-->
# FlowWorkflows — declare the pipelines of the work flow, find what is missing or broken, and seed them without taking over what the team owns

> **Code**: `FLWRF`

## Overview

The work flow runs in pipelines Anchors carries in the binary: the card creation, the gates on every pull
request, the claim that hands work to agents, the move to review, the queue resolution, the board, the
release of abandoned cards, the guard against state changes from outside the flow, and the return of a
card whose blockers were all delivered. This unit holds the single list of those pipelines — the one the
doctor checks and the one `--fix` seeds — and the Go side of what is done with it. The pipeline templates
themselves (their triggers, jobs and scripts) are not this unit's rules.

Most pipelines must run one at a time: serialization is what stops two runs from assigning the same card
or creating it twice, and cancelling a run in progress could kill it after it created half the cards. A
pipeline that is present but not serialized is the worst case, because it looks configured and brings the
race back in silence, so it is a finding of its own, distinct from "missing".

Ownership is decided by a marker. A seeded file carries it, and while it does, the file is Anchors' own:
seeding brings it up to date, and the doctor calls it outdated when it differs from what seeding would
write now. A file whose marker was removed belongs to the team — another stale rhythm, one more
permission — and is never rewritten nor called outdated. The board page is seeded with its pipeline
(seeding one without the other leaves the flow halfway), outside the pipelines folder, and the seeding
says what it did with it; an identical page is never reported as updated, because a notice that fires
when nothing changed trains people to ignore it.

The project's integration branch — where the work arrives — is written into the branch lines the templates
mark for it. The unit also holds the flow's vocabulary: the work states, the board columns, and the
labels, including the per-card ones built from a prefix and the card.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | an existing folder | a folder where the pipelines folder cannot be created or written | this unit: the seeding returns the error |
| the configuration | the loaded project configuration, with or without a work-flow block | — | the configuration, which defaults the integration branch to `main` |
| the carried templates | one template per listed pipeline, plus the board page | a binary that lacks one | this unit: the seeding returns the error |
| the files on disk | any content | — | this unit: presence, the two serialization words and the marker are all it reads |

## Effects

| Effect | Description |
| --- | --- |
| `FLWRF-B01` | Every pipeline in the flow's list has a carried template and a role, and every one requires serialization except the gates and the board pipelines. |
| `FLWRF-B02` | A pipeline is reported missing (`MissingWorkflow`) only when its file is absent from the pipelines folder; any content counts as present. |
| `FLWRF-B03` | A present pipeline (`SemConcurrency`) that requires serialization and lacks either the concurrency block or the refusal to cancel the run in progress is flagged; absent files and pipelines that need no serialization are never flagged. |
| `FLWRF-B04` | Seeding (`SemeiaWorkflows`) writes every missing pipeline, rewrites every one that still carries the marker with the current template, and returns the names it wrote, sorted. |
| `FLWRF-B05` | Seeding leaves a pipeline without the marker (`ÉTemplateIntacto` is false) untouched and does not list it as written. |
| `FLWRF-B06` | Seeding also writes the board page, outside the pipelines folder, carrying the marker. |
| `FLWRF-B07` | The board page seeding reports "created" when there was no page, "updated" when Anchors' own page differed, and "unchanged" when the page was identical or belongs to the team. |
| `FLWRF-B08` | A pipeline is outdated (`OutdatedWorkflows`) when it still carries the marker and its content differs from the template as seeding would write it now; team-owned and absent pipelines are never outdated. |
| `FLWRF-B09` | When the project declares an integration branch other than `main`, every template line marked as the integration branch becomes a branch list with that branch, keeping its indentation. |
| `FLWRF-B10` | Anchors writes the board columns (`ColumnsAnchorsWrites`) only from TO DO up to READY TO TEST. |
| `FLWRF-B11` | A per-card label (`LabelSob`, `LabelDesbloqueia`, `LabelDePR`, `LabelBlockedBy`) is its prefix followed by the card. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLWRF-I01` | What seeding writes is never outdated for the same configuration, whatever the integration branch. | seeds with the default and with a `develop` branch, then asks for the outdated pipelines with the same configuration |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLWRF-X01` | Does not take over a file the team owns: a pipeline or a board page without the marker is never rewritten, whatever it holds. | Rewriting the team's customization with the default would erase deliberate work without warning. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLWRF-E01` | The pipelines folder cannot be created. | The seeding returns an error naming the folder and writes nothing. | Reporting success with no pipeline on disk would leave the flow silently off. |
| `FLWRF-E02` | A pipeline cannot be written. | The seeding stops with an error naming the file and returns what it wrote before. | The caller must know which pipelines exist and which do not. |
| `FLWRF-E03` | The binary does not carry a listed template. | The seeding stops with an error naming the template. | A doctor that demands a pipeline Anchors cannot create would loop forever. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `Workflow.IntegrationBranchOrDefault` | config — the project's integration branch |
| DEP2 | `internal/scan/upstream.go` | `UpstreamDir`, `UpstreamMarker` | scan — where the pipelines live and the ownership marker |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
