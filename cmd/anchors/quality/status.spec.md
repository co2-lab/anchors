<!-- @anchors
  code: PRSTP
  updated_at: 2026-10-08
  layer: comando
-->
# ProjectStatus — where the project stands in the cycle, and the one next step

> **Code**: `PRSTP`

## Overview

Coming back to the work is always started by someone who does not remember where it stopped — often an
agent opening a conversation days later. Without an answer it starts by guessing: did the interview
happen? is the map built? is there work in progress?

The status command answers "where am I?" by walking the cycle in order and stopping at the first step
that is missing, naming it as the next step. Listing everything pending at once would make the reader
choose where to start, and the cycle's order is what they should not have to rebuild. The ladder is: a
git repository; the discovery interview (the project document) and the configuration; the map; then the
work itself, read from where the workflow mode says the queue lives — the task and issue folders in local
mode, the repository's cards in github mode. Along the way it names the informative gates that are already
clean everywhere and could be promoted to blocking.

In github mode it first checks the pipelines are in place, states the pull-request flow, and before telling
an agent to claim new work it lists the open cards that agent already owns. A project that holds only the
guides the init seeded has no work yet, and the next step there is the first plan, not "nothing pending".

It is read by a person and by an agent alike, and it changes nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | any directory, in any phase of the cycle | — | this unit: every phase, including "not started", has an answer |
| the configuration | absent, or one that loads | one that exists and does not load | this unit: it fails (`PRSTP-E01`) |
| the agent identity | the agent name in the environment, or none | — | this unit: without it no card is attributed to anyone (`PRSTP-B12`) |

## Effects

| Effect | Description |
| --- | --- |
| `PRSTP-B01` | A directory with no git repository stops at git init. |
| `PRSTP-B02` | Without the git binary status warns and goes on. |
| `PRSTP-B03` | A project with neither PROJECT.md nor configuration is not started: the next step is the DISCOVER phase. |
| `PRSTP-B04` | A project with PROJECT.md and no configuration is sent to init, with the discovery shown as done. |
| `PRSTP-B05` | A configured project with no map is sent to the map build. |
| `PRSTP-B06` | The configuration is shown with its numbers of layers and gates and the map with its numbers of nodes and edges; informative gates clean everywhere are named as candidates for blocking. |
| `PRSTP-B07` | The local queue shows the pending tasks and the issues in todo and doing, and names one next step in this order: finish the work in doing, pick from todo, take the next task, or nothing pending. |
| `PRSTP-B08` | An assembled local project with no work is sent to the first plan. |
| `PRSTP-B09` | The github queue names the repository and label and, when workflow pipelines are missing, stops at the doctor fix. |
| `PRSTP-B10` | With the pipelines in place, the github queue states that work enters by pull request to the integration branch, and lists the protected branches when there is more than one. |
| `PRSTP-B11` | The agent's own open cards, those whose last owner comment names this agent, are listed with number, title and state, and the agent is told to finish them instead of claiming new work. |
| `PRSTP-B12` | Without an agent identity no card is attributed and the next step is to claim work. |
| `PRSTP-B13` | A github project whose map holds only guides is sent to the first plan before any card. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PRSTP-I01` | Status never names a step beyond the first one missing. | a project with PROJECT.md and no configuration names init and never the map build |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PRSTP-X01` | Status leaves the project as it found it: no file is created, changed or removed. | It is the question asked before deciding anything; answering it must not be a step of the cycle. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PRSTP-E01` | The configuration file exists and does not load. | Error naming the configuration file. | A broken configuration is not a phase of the cycle; reporting a next step on top of it would send the reader past the real problem. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
