<!-- @anchors
  code: TSSTT
  updated_at: 2026-10-08
  layer: comando
-->
# TaskStatus — discover what the machine knows about the task at hand, so the agent's report does not have to

> **Code**: `TSSTT`

## Overview

An agent works and answers, and what it answers was entirely up to it: there was no format for a
progress report, and what varies first is what decides whether someone continues. Measured: an agent
diagnosed a CI failure, fixed it, pushed, and ended its turn with "waiting for the new run". The report
was right; what was missing — the card still in progress under its name, and the verdict of the check
it triggered unread — was not missing by carelessness: nothing asked for it.

`anchors task-status` discovers what the machine knows and prints it; the agent adds only what the
machine cannot know. This unit is the discovery: the card (given by number, or the one the board says
this agent owns) and its state, the decisions stopped waiting for a person, the working tree, the pull
request of the branch and the verdict of its checks, and the moves the state lock undid on the card.
Every lookup names the configured repository and the root's branch, never what the working directory
would imply. Every source fails silently on purpose: a partial report is useful, and a command that
aborts because one lookup did not answer leaves the agent with no format at all.

It also records, as telemetry, the state the turn ended in — only numbers and product vocabulary.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the card | a number given with `--card`, or none | — | this unit: without a number, the board is asked for the agent's card |
| the project root | a directory, possibly a git repository on a branch | a detached head | this unit: no branch, no pull request lookup |
| the workflow mode | github or local | — | this unit: cards are only looked up in github mode |

## Effects

| Effect | Description |
| --- | --- |
| `TSSTT-B01` | The working tree is read from the root: its branch, and whether it has any uncommitted or untracked change. |
| `TSSTT-B02` | In github mode, a card given by number is read with its labels, and its state is its `anchors:` label — or "closed" when the card is closed, whatever label it kept. |
| `TSSTT-B03` | In github mode without a number, the card is the one the board says this agent owns. |
| `TSSTT-B04` | In github mode, the open cards waiting on a person's decision are listed with their numbers and titles. |
| `TSSTT-B05` | The pull request is the one of the root's branch, looked up from the root; each check counts as running when it has no conclusion or has not completed, as passed when it succeeded, was neutral or was skipped, and as failed otherwise; a legacy commit status that is pending or expected counts as running, and its state is the verdict otherwise. |
| `TSSTT-B06` | The moves the state lock undid on the card are the comments that start with the lock's marker AND were written by the automation's account; each is reduced to its first line, without bold or code marks and without the marker. |
| `TSSTT-B07` | The command prints the report of the state it discovered. |
| `TSSTT-B08` | The state the turn ended in is recorded as telemetry: the card's state without its prefix, whether there is a pull request, its state, and the counts of checks — never the card's title nor the labels' full names. The attribute names are stable English identifiers (`has_card`, `has_pr`, `clean_tree`, `unpushed`, `card_state`, `pr_state`, `checks_total`, `checks_failed`, `checks_running`), whatever the user's language. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `TSSTT-I01` | Each check falls in exactly one class, so the classes add up to the total. | six checks of every kind are classified, and running, failed and passed add up to six |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `TSSTT-X01` | Every platform lookup names the configured repository; none lets the tool infer it from the working directory. | From a fork, or with the root elsewhere, the inferred repository is another one, and the report would describe its card and its pull request as this one's. |
| `TSSTT-X02` | Outside github mode no card is looked up. | In local mode a card is not an issue. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TSSTT-E01` | A lookup fails or answers with something unreadable. | That part of the state is absent (no card, no pull request, no waiting list, no undone move), and the report is still printed. | A partial report beats none; nothing is invented to fill the gap. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
