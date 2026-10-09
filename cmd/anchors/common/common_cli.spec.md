<!-- @anchors
  code: CMCLC
  updated_at: 2026-10-09
  layer: comando
-->
# CommonCLI — the contract every command shares: how a path becomes a node, how "not governed" is signalled, what the binary says it is

> **Code**: `CMCLC`

## Overview

Every command receives paths from people and from scripts, answers about nodes of the project map, and
must agree with the others on two things a script depends on: the signal for "this path is not governed"
and the version the binary reports. This set of small pieces holds that shared contract, so that each
command does not reinvent it.

A path arrives in two legitimate shapes: relative to the project root (what the tool's own prompts print)
or relative to the working directory or absolute (what a shell hands over). Both must land on the same node
identifier, with forward slashes, because the map names nodes that way on every system. A node exists only
when the map holds that exact identifier. A task slug is the path without its last extension.

A path that no layer governs is not an error of the tool; it is an answer. The commands signal it with a
dedicated error that names the path, and the process exits with a dedicated code so a script can tell
"nothing to confront" apart from a failure.

The binary reports a version, a commit and a build date stamped at build time. A build that was not stamped
says so plainly (`dev`, `none`, `unknown`) instead of pretending to be a release.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the path argument | a root-relative, working-directory-relative or absolute path | — | this unit: every shape is normalised to a root-relative, forward-slash path |
| the project root | an absolute directory | a relative root | the caller: the commands resolve the root to an absolute directory first; this unit keeps the argument as given when it cannot relate the two |
| the map | a loaded map, or none | — | this unit: no map means no node exists |
| the build stamp | values injected at build time | — | the release build; an unstamped build keeps the defaults |

## Effects

| Effect | Description |
| --- | --- |
| `CMCLC-B01` | The "not governed" error names the path, quoted, says it is not governed, and keeps the path for the caller to read back. |
| `CMCLC-B02` | The exit code that means "not governed" is 3. |
| `CMCLC-B03` | `RelTo`: a path that exists relative to the root is kept as that relative path, cleaned and with forward slashes. |
| `CMCLC-B04` | An absolute path, or a path that does not exist under the root, is resolved from the working directory and expressed relative to the root, with forward slashes. |
| `CMCLC-B05` | `NodeExists`: a node exists only when the map holds a node with exactly that identifier; with no map, no node exists. |
| `CMCLC-B06` | `RelSlug`: a task slug is the path without its last extension only. |
| `CMCLC-B07` | A binary built without a stamp reports the version `dev`, the commit `none` and the date `unknown`. |
| `CMCLC-B08` | A command that takes files reads each argument as a file or as several separated by commas — the list a `--changed` flag takes —, drops blanks, keeps a file named twice once, and keeps the order; `TakesFiles` gives a command that reading and marks it. (`FileArgs`, `TakesFiles`) |
| `CMCLC-B09` | An exit code a command ends with is an error carrying the code, which says it: `exit status <code>`. (`ExitCode`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CMCLC-I01` | The same file named root-relative or absolute resolves to the same node identifier. | resolves one file through both shapes and compares the two results |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CMCLC-E01` | The path cannot be expressed relative to the root (the root is relative and the path absolute). | The path is returned as given, with forward slashes. | A path the caller can still show is better than an empty one that would match nothing and hide the cause. <!-- @resilient: the path as given is still a correct name for the file, and the caller that fails to find it as a node reports that with the path in hand --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
