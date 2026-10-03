<!-- @anchors
  code: RPBUG
  updated_at: 2026-10-03
  layer: comando
-->
# ReportBug — report a bug in Anchors itself, where it is fixed for every project

> **Code**: `RPBUG`

## Overview

An agent using Anchors in a project meets bugs that are Anchors' own: a gate that misreads a file, a
command that writes nothing when it should, a seeded file that is wrong. Recorded only in the project,
the bug is worked around there and met again in the next project. `anchors report-bug` sends it to
github.com/co2-lab/anchors, where the fix reaches every project.

The issue has the sections of the repository's bug form — what happened, what should happen, a
minimal case, the version and the platform — so a report from an agent reads like one from a person,
and the release and platform are filled in by the binary, not remembered by the agent. The same
defect met by several projects is one issue: an open issue with the same title receives a "seen
again" comment instead of a duplicate.

That repository is public, and the agent that writes the report is inside a private project. So the
command refuses a text that carries the project's path, the user's home folder or the project's
repository, and `--dry-run` prints the issue without sending it.

`anchors escalate --bug --upstream` uses the same report, with the reason alone and a warning
instead of a refusal: there the finding is already recorded in the project.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| what happened | one or more words: the command or gate, and what it did | nothing | this unit: refused without it |
| `--expected` | what should have happened | empty | this unit: refused without it (`RPBUG-E01`) |
| `--repro` | a file path, or `-` for standard input | an unreadable file | this unit: refused (`RPBUG-E03`) |
| the project | any directory, with or without an `anchors.yaml`, in any mode | — | this unit: the configuration only adds the repository to the marks |

## Effects

| Effect | Description |
| --- | --- |
| `RPBUG-B01` | The title is `[bug] ` and the first line of what happened; the body has the sections `What happened`, `What should happen` (only when given), `Minimal case` (only when given, in a shell block), `Version` and `Platform`, and says it was reported by `anchors report-bug`. |
| `RPBUG-B02` | When an OPEN issue at co2-lab/anchors has the same title, ignoring case, it receives a comment saying it was seen again with the release and the platform, no issue is created, and the output says "already reported to Anchors" with its address. |
| `RPBUG-B03` | Otherwise an issue is created at co2-lab/anchors with the title, the body and the `bug` label, and the output says "reported to Anchors" with its address. |
| `RPBUG-B04` | `--dry-run` prints the title and the body and sends nothing. |
| `RPBUG-B05` | `--repro` puts the file's content, or standard input's with `-`, in the minimal case. |
| `RPBUG-B06` | It works in a directory with no `anchors.yaml`. |
| `RPBUG-B07` | `escalate --upstream` whose reason names the project does not report it: a warning names what was found and points to `anchors report-bug`. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RPBUG-X01` | A report whose title or body carries the project's absolute path, the user's home folder, the project's `workflow.repo` or the repository of its `origin` remote (other than co2-lab/anchors) is refused, naming what was found, and nothing is sent. | The repository is public, and the text is written from inside a private project. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RPBUG-E01` | `--expected` is missing. | Refused, saying it is required; nothing is sent. | What should have happened is half of a bug report; without it the issue asks the reader to guess. |
| `RPBUG-E02` | The platform refuses to create the issue. | The output carries the prefilled new-issue link, and the command fails with "could not report to Anchors". | A person can still file it; a success would say it was reported. |
| `RPBUG-E03` | The `--repro` file cannot be read. | Refused with the read error; nothing is sent. | A report without the case it promised is worse than none sent yet. |
