<!-- @anchors
  code: QLCMQ
  updated_at: 2026-09-29
  layer: comando
-->
# QualityCommands — the quality domain puts its twelve commands under the root command

> **Code**: `QLCMQ`

## Overview

The CLI is assembled from domains, and each domain hands its commands to the root command in one place.
The quality domain owns the commands that judge and measure the project: the check and its commit-time
wrapper verify, the doctor, the status, the stale-edge listing, coverage, the reports, the test and mutation
suites, the contract stamps and the date touch.

This unit is that one place. It holds no behaviour of its own: each command's rules live in the spec of the
file that builds it. What this unit guarantees is the set: exactly those twelve commands are registered,
each one is reachable by its name from the root, and no two of them answer to the same name. A command
built but never registered would exist in the code and be invisible to the user; a name registered twice
would make one of the two unreachable.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the root command | the CLI's root command, into which the domains register | a root that already holds a command with one of these names | the CLI's root assembly: each domain owns disjoint names |

## Effects

| Effect | Description |
| --- | --- |
| `QLCMQ-B01` | Registers exactly twelve commands under the root: check, verify, doctor, status, stale, coverage, report, test, mutation, stamp, touch and keep-evidence. |
| `QLCMQ-B02` | Each registered command is reachable from the root by its name. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `QLCMQ-I01` | No two registered commands answer to the same name. | registers into an empty root and checks every command name is distinct |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `QLCMQ-X01` | Registering runs no command and prints nothing; it only attaches the commands. | The root is assembled on every invocation, including `--help`; work done at registration would run for every command of every domain. |

## Errors

none — registration attaches commands to the root and handles no failure: there is no input to reject and nothing that can be missing at that point.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/quality/check.go` | `newCheckCmd` | comando — the check command |
| DEP2 | `cmd/anchors/quality/suite.go` | `newTestCmd`, `newMutationCmd` | comando — the suite commands |
| DEP3 | `cmd/anchors/quality/report.go` | `newReportCmd` | comando — the report command |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
