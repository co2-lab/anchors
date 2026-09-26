<!-- @anchors
  code: EXCMX
  updated_at: 2026-09-26
  layer: gate
-->
# ExternalCommand — executes external tools via shell passing targets as positional arguments

> **Code**: `EXCMX`

## Overview

Executes an external command (such as jest, eslint, or tsc) defined in `anchors.yaml` against targets without reimplementing the tool, reading its exit code. To prevent command injection, target paths are never interpolated into the shell command string; instead, they are passed as positional arguments (`$1`, `$2`, `"$@"`) via `sh -c`. When target arguments exceed command line length limits (6000 bytes on Windows, 100000 bytes on Unix, or configurable via `ANCHORS_ARGV_MAX`), targets are partitioned into batches and executed sequentially, combining any failures.

## Signature

| Parameter | Type | Description |
| --- | --- | --- |
| `command` | `string` | Shell command template from configuration, which may contain `{{file}}` or `{{files}}` placeholders |
| `targets` | `[]string` | Target file paths to pass to the external command |
| `root` | `string` | Working directory where the command executes |

**Returns**: `(Verdict, string)` where exit code 0 yields `Pass` with an empty detail string, and non-zero yields `Fail` with execution output or error details.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `command` | any shell command string | empty command string | the configuration loader, which parses valid gate definitions |
| `targets` | slice of target path strings | nil slice or paths with control characters | this unit: nil slice runs as a single project-level invocation |
| `root` | valid directory path | inaccessible directory | the caller, which resolves the repository root directory |

## Effects

| Effect | Description |
| --- | --- |
| `EXCMX-B01` | An external command exiting with status zero returns `Pass` with empty detail. |
| `EXCMX-B02` | An external command exiting with non-zero status returns `Fail` containing the trimmed output. |
| `EXCMX-B03` | If an external command fails with empty output, the failure detail reports that the gate produced no output along with the execution error. |
| `EXCMX-B04` | When executing for a single node via `runExternal`, it delegates to `RunExternalArgs` passing the node ID as the single target. |
| `EXCMX-B05` | The placeholder `{{file}}` in the command template is rewritten to `"$1"` before shell invocation. |
| `EXCMX-B06` | The placeholder `{{files}}` in the command template is rewritten to `"$@"` before shell invocation. |
| `EXCMX-B07` | When targets are empty, execution runs exactly once without positional target arguments representing project scope. |
| `EXCMX-B08` | Targets fitting within the command line budget run in a single shell execution. |
| `EXCMX-B09` | Targets exceeding the budget are partitioned across multiple batches with none dropped or duplicated. |
| `EXCMX-B10` | When targets are partitioned, failure in any batch causes the entire gate to fail combining failure outputs. |
| `EXCMX-B11` | For a single target, failure output exceeding 500 characters is truncated with a truncation marker. |
| `EXCMX-B12` | For batch or project executions, failure output exceeding 4000 characters is truncated with a truncation marker. |
| `EXCMX-B13` | Environment variable `ANCHORS_ARGV_MAX` overrides the default argv limit when set to a positive integer. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `EXCMX-I01` | Target paths are passed strictly as positional argv arguments and never interpolated directly into the shell script string. | inspects shell command arguments ensuring target filenames with shell metacharacters do not execute injected commands |
| `EXCMX-I02` | A single target whose path length exceeds the limit is assigned its own batch and never silenced or dropped. | slices an oversized target with a small budget and verifies it is preserved intact in its own slice |
| `EXCMX-I03` | On Windows runtime without environment override the default target limit is 6000 bytes, while Unix defaults to 100000 bytes. | calls the argv limit calculation under platform defaults and asserts the respective bound |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `EXCMX-X01` | Does not parse or interpret linter or test tool diagnostics beyond reading stdout and stderr. | Following design principle D5, the framework reimplements spec and relation parsing but delegates tool execution directly to external linters and test runners. |
| `EXCMX-X02` | Does not aggregate cross-file state across partitioned batches. | Batch scope assumes file-by-file independent validation; cross-file analysis must declare project scope and scan all targets itself. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `EXCMX-E01` | REF[EXCMX-B02]: a command that exits non-zero is answered by B02: Fail with its trimmed output | — | — |
| `EXCMX-E02` | REF[EXCMX-B03]: a command that fails with no output is answered by B03: the detail says it produced none, with the execution error | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Node` | core — provides the node abstraction and its identifier for per-node execution |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
