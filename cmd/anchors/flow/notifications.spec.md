<!-- @anchors
  code: NTFCT
  updated_at: 2026-10-08
  layer: comando
-->
# Notifications — a message to every agent, read from one file and printed on top of `next`

> **Code**: `NTFCT`

## Overview

Whoever runs a project had no way to tell every agent something at once: a comment on a card reaches
only the agent holding that card, and only if it rereads the card. Measured on 2026-09-24: five agents
on another machine kept opening bugs as decisions after the bug exit of `escalate` existed, and the
only channel to say "update the binary" was a comment on each escalation after the fact.

The message lives in one file at the repository root, `notifications.md`. `anchors next`, the command
every agent runs between two pieces of work, prints its content on top, before the card, every time;
emptying the file stops it. The file may explain itself in HTML comments and still count as empty,
because an explanation printed on every `next` would teach agents to skip the block.

In github mode the file is read from the integration branch on the platform, not from the agent's
checkout, which can be days behind: a message reaches every agent as soon as it is merged. In local
mode the file at the project root is the message. A read that fails is one line and never costs the
agent its work.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file's content | any text, with or without HTML comments | — | this unit: comments are removed before deciding whether there is a message |
| the integration branch | a branch name, possibly with characters a URL escapes; or none | — | this unit: escapes it in the address; without a branch, the platform's default is read |
| the repository | the configured `owner/name` | — | the caller (`next`): passes the configured repository |

## Effects

| Effect | Description |
| --- | --- |
| `NTFCT-B01` | In local mode the message is the content of `notifications.md` at the project root, and a missing file is an empty message. |
| `NTFCT-B02` | HTML comments are removed and the rest trimmed; a file with only comments, or empty, prints nothing. |
| `NTFCT-B03` | In github mode the file is read raw from the repository on the platform at the integration branch, the branch escaped in the address; without a branch no ref is sent. |
| `NTFCT-B04` | A file missing on the platform (not found) is silence, not an error. |
| `NTFCT-B05` | A message is printed as a block that names the file and where it was read from, with every line of the message indented. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `NTFCT-I01` | Text inside HTML comments is never printed. | a file with an explanatory comment and a message prints the message and not the comment |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `NTFCT-X01` | In github mode the agent's checkout is not read. | A worktree can be days behind; the message must reach every agent as soon as it is merged. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `NTFCT-E01` | The platform read fails for any reason other than a missing file. | The failure is returned with the platform's answer, and printed as one line saying the file could not be read, with no message block. | A broken read must not look like "no message", and it must not cost the agent its work either. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
