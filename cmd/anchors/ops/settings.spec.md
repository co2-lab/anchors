<!-- @anchors
  code: STCMS
  updated_at: 2026-10-08
  layer: comando
-->
# SettingsCommand — one agent's local decisions, declared with a date and kept out of the project's configuration

> **Code**: `STCMS`

## Overview

`anchors.yaml` is the project's structure: versioned, reviewed, the same for everyone. Some
decisions hold for ONE agent on one machine instead, and `settings` is where they are
declared: they live in the agent's local settings file, which does not go to git.

Two decisions are recorded here. The ROLE of whoever operates the agent says which work is
theirs and derives the agent's capabilities; only some roles act on escalated cards, which
await a decision from whoever knows the product. The older, narrower decision, whether the
agent acts on escalated cards at all, is still accepted on its own. Each declaration is
stamped with the agent's identity and a date given by whoever declares it: the command never
reads the clock, and a decision with no date is refused.

The value can be given as an argument or answered on the terminal. On the terminal, a reply
the command does not understand is pointed out and asked again, a limited number of times;
it is never assumed to mean "no", because whoever typed something was answering. `show`
prints what is recorded, with the role's capabilities, or teaches how to declare a role when
none is.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the date | any non-empty text given with `--date` | an absent or empty date | this unit: refuses before recording anything |
| the role | a role the settings package recognizes | any other word | the settings package parses it; this unit refuses the unrecognized one |
| the escalated-cards answer | a yes or no the settings package recognizes | any other word | the settings package parses it; this unit refuses, or asks again on the terminal |

## Effects

| Effect | Description |
| --- | --- |
| `STCMS-B01` | `role` and `user-issues` without `--date` are refused with an error naming `--date`, and nothing is recorded. |
| `STCMS-B02` | An unrecognized role or answer given as an argument is refused naming the value, and nothing is recorded. |
| `STCMS-B03` | Declaring a role records the role, the agent's identity and the date, removes the legacy escalated-cards answer, and says where it was recorded. |
| `STCMS-B04` | Declaring the escalated-cards answer records it with the agent's identity and the date. |
| `STCMS-B05` | Without an argument, `role` asks for the role on the terminal and records the answer. |
| `STCMS-B06` | Without an argument, `user-issues` asks on the terminal; an unrecognized reply is pointed out and asked again, and after three unrecognized replies the command fails instead of assuming an answer. |
| `STCMS-B07` | `show` with a role prints the role and every capability it allows. |
| `STCMS-B08` | `show` without a role teaches how to declare one and lists no capabilities. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `STCMS-I01` | After a role is declared there is one source for the escalated-cards question: the legacy answer never survives a role declaration. | records a legacy answer, declares a role, and reads the settings back |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `STCMS-X01` | The decisions are written to the agent's local settings file, never to `anchors.yaml`. | They hold for one agent on one machine; the project's configuration is shared and reviewed. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `STCMS-E01` | The terminal input closes before an answer is read. | The command fails with the read error, not a silent "no". | A closed input is not an answer; recording "no" would make the agent decide on its own. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
