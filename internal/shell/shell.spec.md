<!-- @anchors
  code: PSXSH
  updated_at: 2026-10-08
  layer: infra
-->
# Shell — the POSIX shell that runs a project's commands

> **Code**: `PSXSH`

## Overview

The commands a project declares — a gate's `run:`, a suite's `run:`, a tests `script:` — are POSIX
shell (`"$@"`, `&&`, `VAR=x cmd`), so they need `sh` on every system. On Linux and macOS it is always
there. On Windows it comes with Git for Windows, whose default install puts only `git.exe` on PATH:
the commands failed with a bare "executable file not found", read as the gate failing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the system | any | — | this unit: the lookup |
| PATH | with or without `sh` | — | this unit: falls back to git's shell on Windows |

## Effects

| Effect | Description |
| --- | --- |
| `PSXSH-B01` | The shell is the `sh` found on PATH, on every system. (`Path`) |
| `PSXSH-B02` | On Windows with no `sh` on PATH, the shell is the one beside the installed git: `bin\sh.exe` or `usr\bin\sh.exe` under the root three levels above `git --exec-path`. |
| `PSXSH-B03` | A command is `sh -c <script>` with the arguments after it as `$0 $1…`. (`Command`) |
| `PSXSH-B04` | A command runs without the variables git exports to a hook that tie it to the repository being committed (`GIT_DIR`, `GIT_WORK_TREE`, `GIT_INDEX_FILE` and their kin): a test suite the hook runs that uses git in a directory of its own acts on that directory, never on the repository; every other variable passes. (`WithoutRepoEnv`) |

## Errors

| Error | When | What the user sees |
| --- | --- | --- |
| `PSXSH-E01` | No shell on PATH, and on Windows none beside git either | an environment error saying no POSIX shell was found and how to get one — never a gate failing (`ErrNoShell`) |

## Dependencies

none

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
