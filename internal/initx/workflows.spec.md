<!-- @anchors
  code: FLWRF
  updated_at: 2026-10-08
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
themselves (their triggers, jobs and scripts) are not this unit's rules, except the two parts whose
scripts its tests run: the release of idle cards and the review verdict.

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

The project's integration branch — where the work arrives, `main` when none is declared — is written into
the branch lines the templates mark for it, whatever branch the template itself carries. The unit also holds the flow's vocabulary: the work states, the board columns, and the
labels, including the per-card ones built from a prefix and the card.

The release of idle cards frees a card nobody moves: an owned card with no progress within its window is
released to any agent, once, and a card sent back by a rejected review gets a shorter window, because its
author only has to resume it. The review verdict is a line the assigned reviewer posts on the pull request;
the pipeline turns it into a commit status on the pull request's card, frees the reviewer, sends a rejected
card back to its author, and at the merge says on the card when the work landed without that outcome.

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
| `FLWRF-B09` | Every template line marked as the integration branch becomes a branch list with the project's integration branch — `main` when none is declared — keeping its indentation; this holds for `main` too, so no template's own branch survives the seeding. |
| `FLWRF-B10` | Anchors writes the board columns (`ColumnsAnchorsWrites`) only from TO DO up to READY TO TEST. |
| `FLWRF-B11` | A per-card label (`LabelSob`, `LabelDesbloqueia`, `LabelDePR`, `LabelBlockedBy`) is its prefix followed by the card. |
| `FLWRF-B12` | The stale pipeline releases an owned in-progress card with no progress within its window by commenting the release with the previous owner, which it reads per card, never from the card list; it never releases a card already released, and it skips and names a card whose comments cannot be read. |
| `FLWRF-B13` | A to-do card whose owner came back from a rejected review is released after the rework window, shorter than the waiting window of an ordinary to-do card. |
| `FLWRF-B14` | The review job publishes the review status of the pull request's card: success or failure on the last approved or rejected verdict line of the reviewer who owned the card when it moved to review, posted after that move, at the start of a line and outside a code block, by someone with write access; pending otherwise, awaiting a reviewer when the card is back in ready-to-review; and no status for a pull request that declares no card or whose card was never moved to review. |
| `FLWRF-B15` | The review job runs on pull request comments that carry a verdict line and may publish statuses, and the mover runs after it, even when it failed or was skipped, never on a comment, receiving its outcome. |
| `FLWRF-B16` | A merge moves the card to ready-to-test, and when the review status was not success it also says on the card that the pull request merged without the review outcome, naming the reviewer and the status, or saying no reviewer had been assigned. |
| `FLWRF-B17` | When a green pull request moves its card to ready-to-review, the review status is published as pending, awaiting a reviewer. |
| `FLWRF-B18` | The claim teaches the reviewer to post the approved or rejected verdict line, and never to move the card on approval. |
| `FLWRF-B19` | An approved or rejected verdict releases the reviewer's ownership of the card, once; a review with no verdict releases nothing. |
| `FLWRF-B20` | When the pull request body references one card and closes another, the review status and the move at the merge are for the closed card. |
| `FLWRF-B21` | A rejected verdict sends the card back to to-do, owned again by the owner before the reviewer; an approved verdict leaves the card for the merge. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLWRF-I01` | What seeding writes is never outdated for the same configuration, whatever the integration branch. | seeds with the default and with a `develop` branch, then asks for the outdated pipelines with the same configuration |
| `FLWRF-I02` | The verdict line the review job parses is the one the review guide teaches. | reads the review guide for the approved and rejected lines, and the review job's script for the pattern that parses them |

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

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
