<!-- @anchors
  code: USSTS
  updated_at: 2026-09-26
  layer: apoio
-->
# UserSettings — the agent's local settings: the declared role, kept out of git

> **Code**: `USSTS`

## Overview

`anchors.yaml` is the project's structure: versioned, reviewed, the same for everyone. The settings are
the opposite: they hold for ONE agent on one machine and never reach git, in `.anchors/settings.yaml`.
They exist because, in a project with several developers, not everyone may decide for the product: an
agent that takes an escalated card and asks the developer running it gets an answer, and that answer
may not be the product owner's. Escalation exists to take the question to whoever decides, and an overly
helpful agent short-circuits it.

So the default is closed: with nothing declared, the agent does not act on escalated cards. The person
declares a role (`AGRLG`), and the capabilities derive from it. Settings written before roles existed
carry a yes/no field about escalated cards; it keeps working, granting only the product decision, until
a role is declared, and a declared role wins over it. The field distinguishes three states on purpose:
never asked (the only case where the agent asks), declared no, and declared yes.

A missing settings file is not an error: a freshly cloned project has none, and that is exactly when the
agent must ask. The saved file explains itself, because whoever finds it later was not in the
conversation.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the settings file | YAML written by this unit, or absent | a hand-edited file that is not YAML | this unit: loading fails naming the file |
| the typed answer | yes or no, in English or Portuguese, short or long, any case | anything else | this unit: no answer is returned, and the caller asks again |

## Effects

| Effect | Description |
| --- | --- |
| `USSTS-B01` | The settings live in `.anchors/settings.yaml` under the project root. (`Path`) |
| `USSTS-B02` | A missing settings file loads as empty settings without error: nothing decided, and escalated cards not handled. |
| `USSTS-B03` | Without a role, the legacy field has three states: never asked (undecided, not handling escalated cards), declared no (decided, not handling), declared yes (decided, handling). |
| `USSTS-B04` | Without a role, a legacy yes grants only the product decision, no other capability. |
| `USSTS-B05` | A declared role decides the capabilities, winning over the legacy field. |
| `USSTS-B06` | Saving creates the state folder when missing, and loading gives back what was saved. (`Save`) |
| `USSTS-B07` | The saved file opens with a header saying it is local and kept out of git, that only product-owner and architect act on escalated cards, and how to declare the role. |
| `USSTS-B08` | A typed answer of s, sim, y or yes is yes, and n, nao, não or no is no, ignoring case and surrounding spaces; anything else is no answer. (`ParseAnswer`) |
| `USSTS-B09` | The description of the settings names the role's title and the agent when a role is declared, and points to declaring the role when only the legacy field is set. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `USSTS-E01` | The settings file exists but is not valid YAML. | Loading fails with an error naming the file. | Reading it as empty would silently drop a decision someone took, and the agent would ask again or act without its role. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/settings/roles.go` | `Role`, `Capability` | apoio — the roles and their capabilities (`AGRLG`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
