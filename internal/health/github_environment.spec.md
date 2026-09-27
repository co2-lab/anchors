<!-- @anchors
  code: GHEGT
  updated_at: 2026-09-26
  layer: infra
-->
# GitHubEnvironment — the doctor warns, before the work starts, about the pieces the GitHub flow silently needs

> **Code**: `GHEGT`

## Overview

The GitHub mode of the flow assumes pipelines and a protected branch are in place. When any piece is missing,
nothing fails loudly: the flow simply DOES NOT HAPPEN. A missing identification pipeline gives no error, it
gives silence, and the artefacts stay without a card forever. A branch without protection accepts a direct
push that skips the card, the review and the pipeline that runs when a pull request opens.

This unit is the doctor's part that looks for those pieces, the same way the doctor looks for a missing git:
it warns ahead, before someone discovers the gap in the middle of a piece of work. It only runs in GitHub mode,
because charging pipelines to a project that declared local mode would be guaranteed noise, and recurring
noise trains the team to ignore the doctor.

It reads the disk for the flow's pipelines (missing, out of date with the template, present without
serialization) and asks the platform for the branch protection and for the approval reachability of
`APRCP`. The board is deliberately NOT checked: the state of work is a label, and the project board is an
optional mirror.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | a configuration in GitHub mode | local mode, or no configuration, which checks nothing | this unit: it returns before any check |
| the project root | a directory whose workflows folder may or may not hold the flow's pipelines | — | the caller, which passes the project root |
| the pipelines | the flow's pipeline files the initializer seeds, with or without the template marker | pipelines the team adopted by removing the marker, which are theirs and are not compared with the template | the initializer's catalogue of flow pipelines |
| the platform CLI | the `gh` binary on the PATH | its absence, which another doctor finding reports | this unit: without `gh` the platform is not asked |

## Effects

| Effect | Description |
| --- | --- |
| `GHEGT-B01` | Outside GitHub mode, or with no configuration, nothing is checked. |
| `GHEGT-B02` | Each flow pipeline missing on disk gives a warning `pipeline-ausente` on its file, saying what stops happening and pointing at the doctor's fix. |
| `GHEGT-B03` | A pipeline present without serialization gives its own warning `pipeline-sem-serializacao`, saying the same card can reach two agents; it is never reported as missing. |
| `GHEGT-B04` | A pipeline that still carries the template marker and differs from the current template gives a warning `pipeline-desatualizado`; one whose marker the team removed is not reported. |
| `GHEGT-B05` | When the branch protection cannot be read (the platform answers that there is none), the finding is a warning `main-sem-protecao` on the repository saying a direct push skips the card, the review and the pipeline. |
| `GHEGT-B06` | A branch protection that exists but does not require pull request reviews gives the same check with the "partially protected" message. |
| `GHEGT-B07` | Without the platform CLI, the branch protection is not asked and gives no finding: another finding already names the missing tool. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GHEGT-I01` | What the doctor charges for pipelines and what its fix seeds come from the same list: after seeding, no pipeline finding remains. | seeds the pipelines into an empty project and checks the pipelines again |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GHEGT-X01` | The project board is not checked, and no finding mentions it. | The state of work is a label and the board is an optional mirror; charging a piece the flow does not need gives a finding nobody needs to resolve. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GHEGT-E01` | REF[GHEGT-B05]: a failed read of the branch protection is the answer for an unprotected branch, so it becomes the finding, not an error | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config`, `GitHubMode` | config — the mode and the repository |
| DEP2 | `internal/initx/workflows.go` | `MissingWorkflow`, `OutdatedWorkflows`, `SemConcurrency` | the flow's pipeline catalogue and its comparisons |
| DEP3 | `internal/i18n/i18n.go` | `T` | apoio — localized finding texts |
| DEP4 | `internal/health/approval.go` | `checkApprovalReachable` | infra — the approval reachability joins the GitHub checks |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
