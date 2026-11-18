<!-- @anchors
  code: INQSN
  updated_at: 2026-10-08
  layer: infra
-->
# InitQuestions — describe the human decisions of `init` so an agent can answer them without the terminal UI, and judge every answer

> **Code**: `INQSN`

## Overview

`init` is interactive, and the terminal guard aborts it outside a terminal. That left an agent with no way
to start a project, precisely in the flow where the user asked it to. The unit turns the decisions of the
terminal UI into a two-call contract: the first call returns the questions, each with what the agent needs
to DECIDE (the accepted values, what Anchors inferred from the disk, and what the answer changes in the
project); the second call brings the answers and gets back a verdict for every one of them.

The questions come in the order of the terminal UI, because each answer narrows the next. Their texts, the
reasons and the refusal details are in the project's language: the agent relays them to the user, who is the
one who really decides. The defaults are
what the inference read from the real project, so an agent with no ground to disagree should accept them.

The verdict covers EVERY question, not only the invalid ones: that is what lets the agent check that Anchors
understood what it meant — a silently ignored answer would be indistinguishable from an accepted one. An
answer not given is told apart from an answer given empty, because "no artifacts" and "use what was
detected" are opposite decisions. And one refused answer refuses the whole set: writing only the valid ones
would produce a configuration nobody decided in full.

The work-queue mode brings a dependency between answers: in `github` mode the repository and the labels
are required, and outside it a repository is refused, because a declared repository makes whoever reads the
file conclude that the integration is active.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the inference proposal | the result of the disk inference, possibly with no configuration built, or nothing | — | this unit: a missing proposal or configuration gives empty defaults, never a crash |
| the answers | for each question, a value or nothing (not answered) | values outside the question's options; `github` mode without repository or labels; a repository outside `github` mode | this unit: each is refused in its verdict |

## Effects

| Effect | Description |
| --- | --- |
| `INQSN-B01` | The questions are, in this order: header, contributing, artifacts, gates, colocation, layers, workflow, repo, labels, governs. |
| `INQSN-B02` | Every question carries its identifier, its text, its answer type and what the answer changes; every single-choice question carries its options. |
| `INQSN-B03` | The defaults come from the inference: the artifacts detected on disk (in name order), whether the project is colocated, and the code layers of the inferred configuration as both the options and the default of the layers question; with no proposal or no configuration, those defaults are empty. |
| `INQSN-B04` | The work queue is chosen among `local`, `manual` and `github`, `local` by default. |
| `INQSN-B05` | The verdict (`ValidateAnswers`) has one entry per question, in the questions' order; an answer not given takes the question's default and is marked as default. |
| `INQSN-B06` | An answer outside the question's options is refused, and the verdict lists the accepted values; a multiple-choice question with no declared options accepts any value. |
| `INQSN-B07` | In `github` mode, a missing or empty repository and a missing or empty label list are each refused. |
| `INQSN-B08` | Outside `github` mode, a repository answer is refused. |
| `INQSN-B09` | A single refused answer makes the whole set refused (`TudoAceito`). |
| `INQSN-B10` | An answer given empty is a decision, not the default: it is kept as given and not marked as default. |
| `INQSN-B11` | The artifacts question offers exactly the artifact options (ARCHR-B01), code included, so an answer of code is accepted. |
| `INQSN-B12` | The question texts, their reasons and the refusal details are written in the project's language. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INQSN-I01` | Every answer in the set shows up in the verdict, so an answer that was not understood can never pass as accepted in silence. | validates a set with one explicit answer and checks that every question has its entry and only the unanswered ones are marked as default |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INQSN-X01` | Does not correct an invalid answer to a valid one. | A corrected answer is a decision the user never made; it is reported as refused instead. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INQSN-E01` | REF[INQSN-B06]: an answer outside the options is the input failure B06 refuses, with the accepted values | — | — <!-- @resilient: the refusal is the answer's status, with the accepted values, returned for every answer; nothing is written until every answer is accepted --> |
| `INQSN-E02` | REF[INQSN-B07]: `github` mode without a repository or labels is the dependency failure B07 refuses | — | — <!-- @resilient: the refusal is the answer's status, with the reason, returned for every answer; nothing is written until every answer is accepted --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
