<!-- @anchors
  code: ARCGN
  updated_at: 2026-10-08
  layer: comando
-->
# AgentRoleCLI — who this agent is, and the role it declared, as the commands show and ask it

> **Code**: `ARCGN`

## Overview

Several agents can work on one project, and the commands must tell them apart and know what each one is
allowed to decide. This unit gives the running agent an identity (its machine and its session), shows
the list of known roles and what the declared role means, asks for a role in the terminal when there is
somebody there to answer, and says whether the agent declared that it decides the product's direction.

The role catalogue itself (which roles exist, what each does, who decides the product, each role's
review lens) belongs to the settings unit; this unit only presents it and asks for it.

The question is never asked in the dark: when the input is a pipe or the null device there is nobody to
answer, and the question is not asked. When somebody answers, an answer that is not a role is repeated
back and asked again, never guessed.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the session name | `ANCHORS_SESSION`, trimmed; else `USER`; else `USERNAME` (Windows) | all empty | this unit: falls back to `default` |
| the typed answer | a role name or an abbreviation the settings unit recognises | any other text | this unit: echoes it back and asks again, at most three times |
| the input stream | a terminal with somebody on the other side | a pipe, the null device, a closed stream | this unit: not interactive; a closed stream while asking is an error |
| the project settings | a readable settings file with a declared role | a missing or unreadable file | this unit: the agent does not decide the product |

## Effects

| Effect | Description |
| --- | --- |
| `ARCGN-B01` | `AgentID`: the agent's identity is the machine name and the session joined by a slash; the session is `ANCHORS_SESSION` trimmed, else the `USER` name, else the `USERNAME` name (where Windows keeps it), else `default`. |
| `ARCGN-B02` | `RoleList`: the role list shows every known role with what it does. |
| `ARCGN-B03` | `PrintRole`: showing a declared role gives its title and what it does, says whether it acts on the escalated cards or must escalate instead, and shows the role's review lens when it has one. |
| `ARCGN-B04` | `AskRole`: asking for the role accepts the first answer the settings unit recognises, abbreviations included; an unrecognised answer is repeated back in quotes before asking again. |
| `ARCGN-B05` | After three unrecognised answers the question gives up with an error that points at the command that declares a role. |
| `ARCGN-B06` | `InteractiveTerminal`: the input is interactive only when it is a terminal: a pipe or the null device is not. |
| `ARCGN-B07` | `DecidesProduct`: the agent decides the product only when its declared role handles the escalated cards; with no settings it does not. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ARCGN-X01` | Does not define the roles, their capabilities or their lenses; it reads them from the settings unit. | One catalogue of roles; a copy here would drift from the one the other commands obey. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ARCGN-E01` | The input closes before any answer while the role is being asked. | An error saying the response could not be read; no role is assumed. | Assuming a role would grant or deny capabilities nobody chose. |
| `ARCGN-E02` | The settings file exists but cannot be read. | The agent does not decide the product. | The capability is closed by default: a broken file must never unlock it. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
