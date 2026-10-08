<!-- @anchors
  code: APRCP
  updated_at: 2026-10-08
  layer: infra
-->
# ApprovalReachable — the doctor says when the required approval can never be given, and how to get out

> **Code**: `APRCP`

## Overview

The GitHub flow assumes that ANOTHER agent reviews, and a required approval is what stops a merge without
review. But the platform refuses to let the author approve their own pull request, and agents on the same
machine share one account. The result is a flow that stops at the last step: the card reaches review, the
reviewer confronts and approves… and cannot. No pull request of the project advances by the normal path.

This is neither a defect of Anchors nor of the project: it is a known restriction with two ways out that
depend on who operates the repository. This unit gives the doctor the means to DETECT the situation before
someone discovers it in the middle of a merge, and to apply the second way out on request.

It answers three questions. Can the current account merge over the requirement? Only an administrator of the
repository whose protection does not also bind administrators can, and each refusal says which of the
conditions failed. Is the requirement reachable? When it is not, one warning names the requirement and both
ways out. And, when the team chooses it, the unit turns the requirement off on the platform, while the
project's configuration remains the source of truth.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the repository | the `owner/name` declared in the workflow configuration | a repository the account cannot read, which becomes a refusal or a warning, never a crash | this unit: every unreadable answer of the platform CLI is read as "could not ask" |
| the branch | the integration branch the configuration resolves, `main` by default | — | the configuration, through its default |
| the required approvals | a count of zero or more from the workflow configuration | a configuration without a workflow, which requires nothing | this unit: a nil configuration or workflow is silent |
| the platform CLI | the `gh` binary on the PATH, logged in | its absence, which another doctor finding already reports | this unit: without `gh` the reachability check is silent |

## Effects

| Effect | Description |
| --- | --- |
| `APRCP-B01` | `CanBypassProtection`: The current account can bypass the protection only when it is an administrator of the repository AND the protection does not enforce its rules on administrators. |
| `APRCP-B02` | Each refusal of the bypass names its reason: the platform CLI is missing, the current user cannot be read, the account is not an administrator (or its permission cannot be read), the protection cannot be read, or the protection binds administrators. |
| `APRCP-B03` | When zero approvals are required, or there is no configuration, the reachability check reports nothing: there is nothing to be unreachable. |
| `APRCP-B04` | Without the platform CLI the reachability check reports nothing, so the doctor does not duplicate the finding that already names the missing tool. |
| `APRCP-B05` | When approvals are required and the current account cannot bypass the protection, the check gives one warning `aprovacao-inalcancavel` on the repository, carrying the required count. An administrator with an escape gets no finding: the escape is the expected way around. |
| `APRCP-B06` | The warning names both ways out: a service account for the agents, or zero required approvals in the configuration applied by the doctor's fix. |
| `APRCP-B07` | `DisableApprovalRequirement`: Turning the requirement off replaces the branch protection on the platform with zero required approving reviews and no enforcement on administrators, sending the new protection as the request body. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `APRCP-I01` | The reachability warning and the bypass answer agree: the warning appears exactly when the bypass is refused. | the same scripted platform answers are given to both, for an administrator with an escape and for a writer |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `APRCP-X01` | The reachability check only READS the platform: it never changes the branch protection. | Changing the repository's protection is the team's decision, applied only when someone runs the fix; a diagnosis that rewrote it would decide for the team. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `APRCP-E01` | The platform refuses the protection update. | An error naming the branch and carrying the platform's own output. | The fix must say which branch failed and why, or the operator retries blind. |
| `APRCP-E02` | REF[APRCP-B02]: an unreadable answer of the platform while deciding the bypass is a refusal with its reason, not an error | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
