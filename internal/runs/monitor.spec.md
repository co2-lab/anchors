<!-- @anchors
  code: MNTRS
  layer: infra
-->
# RunMonitor — reading the project's runs tick by tick: what started, stalled, died, ended, and what runs

> **Code**: `MNTRS`

## Overview

Each tick, the monitor reads a snapshot — the time, the system's processes, the records of the
runs, the load, the test reports — against what it remembers, and returns the events that call
for a reaction, and the records to write (DESIGN-process-monitor.md). A run is a record a
launcher wrote, or a process under the project that a runner recognizes. A run moves while its
CPU time advances, its output grows or new processes appear under it; it is stalled when none
of that happens for its runner's quiet time; it ended when its process is gone or its record
says so — and how it ended comes from its exit, its output's summary, or nothing. The memory
survives the watch: a re-armed monitor resumes it. Runners — how a kind of process reads, and
how its output says it passed, failed or advanced — have defaults and come from the project's
configuration too.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the snapshot | the processes listed in one reading, the records, the reports | — | the caller (`anchors monitor`) |
| the monitor's state | the state a previous tick wrote, or none | a state another tool wrote | this unit: an unreadable state is a fresh start |
| the `monitor:` block | durations in Go's form, patterns as Go regular expressions | — | this unit: one that does not read is an error naming it |

## Effects

| Effect | Description |
| --- | --- |
| `MNTRS-B01` | A process under the project that a runner recognizes and no record accounts for is a run found running: a started line and a record by the process table. A launcher and the workers it starts are one run — the topmost; a process of another folder, the monitor and what is above it, Anchors' own wrapping and watching processes, and the launcher of a process a run already accounts for — the shell its command was typed in — are no run. (`Tick`, `projectRunners`, `launchesClaimed`) |
| `MNTRS-B02` | A run whose CPU time did not advance, whose output did not grow and under which no new process appeared for its runner's quiet time — the default one when it has no runner — gives one stalled line; when it moves again, a line says so. A process under it that ended is no sign of life. |
| `MNTRS-B03` | A run that ended is said once: with its exit code when its record has one; by its output's summary when it does not; as died, with its output's last lines, when the output ends without a summary or when its launcher was to record its exit and did not; and as ended with its exit unknown when only the process table knew of it; a command of the agent the monitor never saw running — launched and gone while none ran — is said to have ended before the monitor saw it, with when it was launched. (`endOf`) |
| `MNTRS-B04` | A record of the agent takes the topmost process under the project that runs the program its command's first line names, past the launchers in front of it, or that the same runner recognizes; a command of the agent no process runs, past its first tick, ended once its output — when known — stopped growing. (`findProcess`, `firstWord`) |
| `MNTRS-B05` | A test report modified since the last reading is a line with its tests, its failures and the first failure; a report already there before the monitor's memory began is no event. |
| `MNTRS-B06` | Every heartbeat, a line says what runs — by name, the stalled marked —, or `running=[NOTHING]`; how many runs ended since the last heartbeat; and the load against the number of CPUs (`load=8.4/10`), marked HIGH above it — processes wait for a CPU. |
| `MNTRS-B07` | The memory is kept between watches: a re-armed monitor says once what happened while nobody watched — a death included — and repeats nothing; a run that ended before the memory began is old news. (`LoadState`, `SaveState`) |
| `MNTRS-B08` | A run whose runner reads progress gives, at most once per progress interval, a line counting the progress marks in the end of its output. |
| `MNTRS-B09` | The settings are the defaults — a tick of 30s, a heartbeat of 5m, progress every 5m, a quiet time of 5m, the built-in runners, 50 ended runs within 7 days —, then the project's `monitor:` block: its timing, its runners read first and replacing a built-in of the same name, its reports beside the suites' JUnit reports, and its bounds. A duration or a pattern that does not read, or a runner with no match, is an error naming it. (`Configure`, `DefaultTiming`) |
| `MNTRS-B10` | An output's end says failed when the runner's failure pattern is in it — even beside a pass —, passed when only the pass pattern is, and nothing when neither is; its summary line is the line carrying the verdict, and its last lines are the last non-blank ones. Reading an output gives its size and its end. (`Verdict`, `SummaryLine`, `LastLines`, `ReadOutput`) |
| `MNTRS-B11` | Reading the reports the globs find gives each one's path from the root and its time; a report modified since its known time is parsed for its tests, failures and first failure. (`ReadReports`) |
| `MNTRS-B12` | The built-in runners recognize jest, vitest, go test, pytest, maestro, playwright, stryker, gremlins and Anchors' long commands, each with its kind and quiet time; a command line no runner matches has none. (`DefaultRunners`, `RunnerFor`) |
| `MNTRS-B13` | A command is recognized by its first line only — what a heredoc carries after it is no command —, and a line of the monitor shows it by that line, cut at a hundred characters, with how many lines it left out (`[+N lines]`). (`FirstLine`, `Brief`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MNTRS-I01` | Every run that ended is said, and said once: across ticks and re-armed monitors, no ended run is silent and none is repeated. | ends runs in each of the four ways across ticks and a reloaded state, and counts the end lines of each |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MNTRS-E01` | A `monitor:` value that does not read — a duration that is not positive, a pattern that does not compile. | The settings are not read, and the error names the key. | A monitor running on a value it guessed would stall or flood without anyone knowing why. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
