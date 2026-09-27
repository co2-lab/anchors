<!-- @anchors
  code: CGPCH
  updated_at: 2026-09-27
  layer: comando
-->
# CheckGatePipeline — confronts the map's nodes against the declared gates, records the verdicts and reports the profile

> **Code**: `CGPCH`

## Overview

The check is the quality pipeline. It loads the configuration and the map, decides which nodes to confront
and which gates to charge, runs the gates, prints the verdict profile, records what it found, and exits 1 when
a blocking gate failed. Everything else in the project that says "passing" — the commit hooks through the
verify facade, the CI pipelines — goes through it.

The scope is either the full sweep or the impact path of the changed files: each file, what propagates from
it, what validates it, and the tests whose contract stamps point at it. A changed file the project does not
govern (it matches no layer, it is the Anchors record, it is a plan's progress companion) is answered with a
distinct not-governed code, so a hook can let a configuration-only commit through; a governed file that is not
yet in the map is an error, because outside the map no gate confronts it.

The gates charged are those of the requested phase and category, minus the slow ones when asked, minus those
that declared they skip the current perspective, minus those waived for this run with a written reason. The
deterministic mode drops the judgment gates, which a commit cannot wait for.

Recording is the loop from reporting to remembering: the confronted edges are stamped in the map, and —
depending on the workflow mode — blocking failures open violation issues, passes resolve them, open
decisions and assumed debts open in their own folders, and a full sweep closes the violations it no longer
reproduces. Judgments left pending become tasks in the local queue, and bar the incremental check of whoever
touched the file; in github mode the queue is the board, and the check instead prints a brief for the
reviewer.

The profile table is read by eye over dozens of gates, so its format is part of the behaviour: columns that
align, a drift column only when there is drift, the clean gates omitted on request but counted, the drift
listed in full on request and grouped by reason, a legend of the symbols used, and a last line that says what
is still open even when the check passes. Whether the time measurement per gate is printed is governed by the
timing-metrics flag's own spec; what the measurement shows is stated here.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | the project's configuration with at least one gate | no configuration, or no gate | this unit: it refuses (`CGPCH-E01`, `CGPCH-E02`) |
| the map | the project's map at its default path or the one given | no map | this unit: it refuses and points at the map build (`CGPCH-E03`) |
| the scope | the full sweep, or one or more changed paths, relative or absolute | neither; a path on neither disk nor map | this unit: it refuses (`CGPCH-E05`, `CGPCH-E07`) |
| the waivers | rule id and reason, from the flag, the environment or the commit message file | a waiver without a reason | this unit: it refuses (`CGPCH-E04`) |
| the phase and category | any declared phase or category, or none | a phase or category no gate declares | this unit: it says the slice is empty and runs nothing (`CGPCH-B07`) |

## Effects

### Scope — which nodes are confronted

| Effect | Description |
| --- | --- |
| `CGPCH-B01` | With the full-sweep flag, every node of the map is confronted, and the header names the scope, the number of nodes and of gates. |
| `CGPCH-B02` | With changed files, the check confronts the union of their impact paths — each file, what propagates from it and what validates it — and no node off those paths. |
| `CGPCH-B03` | The tests whose contract stamps point at a changed file enter that file's impact path; tests with no stamp on it do not. |
| `CGPCH-B70` | The other pieces of a changed file's unit enter its impact path by name, whether or not the map links them: spec, feature, code and test, the test named with the Go suffix or the TypeScript one, and the code with the changed file's own extension. |
| `CGPCH-B04` | Changed paths are compared in the map's form: relative to the root, cleaned, with forward slashes, whether the caller passed them absolute or relative. |
| `CGPCH-B05` | A changed file is not governed when it matches no layer, or when it lies under a directory the scan ignores, such as the issues folder that holds the Anchors record. |
| `CGPCH-B06` | In a batch of changed files, the not-governed ones are dropped and the rest confronted; when every file is not governed, the answer is not-governed, naming the first file and counting the others. |
| `CGPCH-B67` | A plan's progress companion is not governed, while the plan beside it stays governed. |

### Gate selection — which gates are charged

| Effect | Description |
| --- | --- |
| `CGPCH-B07` | Only the gates of the requested phase and category are charged; when none is left, the check says so and runs nothing. |
| `CGPCH-B08` | A gate that declares no perspective to skip runs both on the full sweep and on changed files. |
| `CGPCH-B09` | A gate that declares it skips the change perspective is not charged on changed files and is charged on the full sweep. |
| `CGPCH-B10` | A gate that declares it skips the full-sweep perspective is charged on changed files and not on the full sweep. |
| `CGPCH-B11` | A gate that declares both perspectives to skip is charged in neither. |
| `CGPCH-B12` | A skip marker with a reason in the commit message file waives that gate for this run. |
| `CGPCH-B13` | A waiver given through the environment, as the git hook passes it, drops that gate and prints it with its reason. |
| `CGPCH-B14` | With the deterministic flag the judgment gates are dropped; when only judgment gates are declared, the check says there is no deterministic gate and runs nothing. |

### Recording — stamps and issues

| Effect | Description |
| --- | --- |
| `CGPCH-B15` | Whether the check writes issues depends on the mode: always in local mode; in github mode only in CI or when asked to record issues; in manual mode only when asked. |
| `CGPCH-B16` | The edges the check confronted are stamped in the map at the current revisions of their ends. |
| `CGPCH-B17` | With the no-record flag the check neither stamps the map nor writes issues: it only reports. |
| `CGPCH-B18` | A blocking failure opens a violation issue only when issues are on; the map is stamped either way. |
| `CGPCH-B19` | A pass resolves the open issue of the same gate and target; a pending open decision opens a user-owned decision in todo and is resolved by its own kind when the gate passes; an assumed debt opens in the future folder; a plain pending opens nothing. |
| `CGPCH-B20` | A full check closes the open violations it did not reproduce; an incremental check leaves them alone. |
| `CGPCH-B72` | When the stamps or the issues cannot be written, the check warns on the error output and still prints the verdict; a record that succeeds raises no warning. |
| `CGPCH-B75` | When issues are off, the record says why no issue was written: the manual mode, or a local run that leaves the board issues to CI. |
| `CGPCH-B76` | In github mode the recorded issues go to the board, never to the local issue folders. |
| `CGPCH-B77` | The record summary counts the edges stamped, the issues opened, open decisions included, and the issues resolved, closed decisions and reconciled violations included; the reconciled violations and the assumed debts get a line of their own only when there are any. |

### Judgment — the queue and the brief

| Effect | Description |
| --- | --- |
| `CGPCH-B21` | In the local modes, each target a judgment gate left pending becomes one judge task in the queue, named after the gate and the target. |
| `CGPCH-B69` | A judge task is of the judgment kind and suggests the review stage, a verb the work command composes; its reason names the command that records the verdict with the task's own gate, and the guide and question when the gate declares them. A queued task with the legacy judge verb still counts as a judgment. |
| `CGPCH-B22` | Judgments waiting in the queue bar an incremental check with exit 1, and never the full sweep. |
| `CGPCH-B73` | After recording, the check says how many judge tasks this run queued, and says nothing when it queued none. |
| `CGPCH-B23` | On a full sweep, a judge task that came from the check, for a gate that ran, whose target this run did not enqueue, is dropped from the queue; tasks of another origin stay. |
| `CGPCH-B68` | A judge task is read back as the gate whose name and the task's target rebuild its exact name, so a gate whose name prefixes another's never takes the other's tasks. |
| `CGPCH-B24` | The drop of stale judge tasks happens only on the full sweep; an incremental check keeps the judge tasks of the targets it did not look at. |
| `CGPCH-B25` | In github mode the check prints the judgment brief even without recording, and queues no judgment locally: the queue there is the board. |
| `CGPCH-B26` | The judgment brief names every target awaiting judgment. |
| `CGPCH-B27` | The judgment brief carries each gate's question and the guide where its ruler is. |
| `CGPCH-B28` | The judgment brief groups the targets by gate, so each question is asked once. |
| `CGPCH-B29` | The judgment brief says the pipeline does not judge and points at the review checklist item. |
| `CGPCH-B30` | The judgment brief lists at most ten targets per gate and counts the rest. |
| `CGPCH-B71` | The judgment brief is printed only when some target awaits judgment. |

### Warnings that do not bar

| Effect | Description |
| --- | --- |
| `CGPCH-B31` | Governed files the map does not hold are warned about as a stale map, measured by content, never by file times. |
| `CGPCH-B32` | On an incremental check, the map nodes whose content moved since the map build are warned about on the error output, the first three named; the full sweep does not warn. |
| `CGPCH-B33` | When the map was written by another version than the running binary, a warning names both versions and says to rebuild the map. |
| `CGPCH-B34` | A map that does not record which version wrote it raises no version warning. |
| `CGPCH-B35` | A declared gate that no node reached is named in a warning instead of vanishing from the table. |
| `CGPCH-B36` | The local backlog is printed on the full sweep of a mode other than github, and never on an incremental check. |
| `CGPCH-B37` | The governance tips are printed on the full sweep only, with the pointer to the doctor. |
| `CGPCH-B74` | When there is no governance tip to give, the full sweep prints neither a tip nor the pointer to the doctor. |

### The mirror and the exit

| Effect | Description |
| --- | --- |
| `CGPCH-B38` | The report is mirrored to a file under the project's work folder, and the check says where. |
| `CGPCH-B39` | A blocking failure ends the check with exit 1, after the mirror has received the end of the report. |

### The profile table

| Effect | Description |
| --- | --- |
| `CGPCH-B40` | The gate-name column is as wide as the longest name, with a floor of 18. |
| `CGPCH-B41` | Each counter column has its own width, that of its largest number. |
| `CGPCH-B42` | The pass, fail, skip and judgment columns are at least 1 wide; the drift column is 0 wide when no gate has drift. |
| `CGPCH-B43` | When no gate has drift, the drift column and its separator are absent from every line. |
| `CGPCH-B78` | A gate's indeterminate counter is its skipped and pending results less its drift, which the drift column counts apart. |
| `CGPCH-B44` | An empty drift cell is as many terminal columns wide as a filled one: the symbol plus the digits. |
| `CGPCH-B45` | A gate is clean when it has no failure, no drift, nothing skipped or pending and nothing awaiting judgment. |
| `CGPCH-B46` | By default every gate appears in the table, clean or not. |
| `CGPCH-B47` | With only-issues the clean gates are omitted from the table and counted in a footer line. |
| `CGPCH-B48` | With show-drift every drift item is listed with its address, with no cap. |
| `CGPCH-B49` | The drift listing is not suppressed on a large scan, unlike the skip reasons. |
| `CGPCH-B50` | Without show-drift the drift is only counted in the table, never listed. |
| `CGPCH-B51` | On a scan of at most forty results, the reason of each skip that has one is listed under the table. |
| `CGPCH-B52` | Above forty results the skip reasons are not listed; the counter is enough. |
| `CGPCH-B53` | The legend explains only the symbols the table used. |
| `CGPCH-B54` | In the drift listing, a reason shared by several targets is written once, followed by the count and one target per line. |
| `CGPCH-B55` | In the drift listing, targets with distinct reasons are listed one by one, each with its reason. |
| `CGPCH-B56` | The drift listing's heading counts the items and the gates they are in. |
| `CGPCH-B57` | The findings heading counts the failures, split into blocking and informative, and the divergences, and says how to list the divergences. |
| `CGPCH-B82` | Under the findings heading each failure is listed with its blocking or informative mark, its gate and its target, and its detail indented below when it has one; with neither failure nor divergence there is no findings heading. |
| `CGPCH-B58` | The last line says what is still open even when the check passes: informative findings, pending items, confrontations that did not happen, or nothing; a failed check says how many blocking gates failed and how many informative findings came with them. |
| `CGPCH-B83` | An informative gate that passed and failed nothing is named as ready to become blocking; when no gate is, nothing is said. |
| `CGPCH-B59` | A detail with occurrences joined by a semicolon and a space is printed one occurrence per line, every line indented. |
| `CGPCH-B60` | One separator is enough to break a detail. |
| `CGPCH-B61` | A detail with a single occurrence stays on one line. |
| `CGPCH-B62` | A semicolon with no space after it does not break a detail. |
| `CGPCH-B63` | A line longer than 110 characters that lists at least four comma-separated items, most of them without inner spaces, is broken one item per line. |
| `CGPCH-B64` | A short sentence with commas stays on one line. |
| `CGPCH-B65` | Long prose with commas is not broken, since its items carry spaces. |
| `CGPCH-B66` | A long list of paths joined by commas is broken. |

### The time measurement

| Effect | Description |
| --- | --- |
| `CGPCH-B79` | The time table counts each gate's targets over every verdict, as one target or as N targets with the numbers right-aligned to the widest, pads the gate names to the longest, and shows the worst target's time only for a gate with more than one target. |
| `CGPCH-B80` | The slowest targets are listed slowest first, and the list is left out when no target recorded any time. |
| `CGPCH-B84` | The check stamps the copy of the map it confronted, and carries only the stamps it changed to the map as it is on disk when it records, under the map's lock: what another process wrote meanwhile — an ingestion's signals, another stamp — is kept. |
| `CGPCH-B81` | Times are rounded to what a decision needs: to ten milliseconds from one second up, to a tenth of a millisecond from one millisecond up, and to the microsecond below. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CGPCH-I01` | The filtering axes are independent: declaring a perspective to skip does not change what phase, category or cost do. | filters a slow gate that skips the change perspective with and without skip-slow |
| `CGPCH-I02` | The version warning compares by equality, never by order: two different names always warn and two equal names never do. | computes the warning for equal and different version names |
| `CGPCH-I03` | The skip column sits in the same terminal column on every line of the table, with or without drift on that line. | prints a table with a drift line and a line without drift and measures the skip column in runes |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CGPCH-X01` | The stale-map warnings never bar the check: they go to the output and the verdict is the gates'. | A stale map does not invalidate the work, only the picture the gates confront; barring there would stop someone from saving work in progress, which is exactly when the map goes stale. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CGPCH-E01` | The project has no configuration. | Error naming the configuration load. | Every gate, layer and mode comes from the configuration. |
| `CGPCH-E02` | The configuration declares no gate. | Error saying no gate is declared. | A check with no gate would pass having confronted nothing. |
| `CGPCH-E03` | The map does not exist or cannot be read. | Error pointing at `anchors map build`. | The gates confront the map's nodes; without it there is nothing to confront. |
| `CGPCH-E04` | A waiver, from the flag, the environment or the commit message, has no reason. | Error "invalid --skip-rule" explaining the reason is mandatory. | A waiver without a written reason cannot be told from someone dodging a gate that found a defect. |
| `CGPCH-E05` | Neither changed files nor the full sweep is asked for. | Error asking for one of the two. | There is no default scope: guessing one would confront either too little or everything. |
| `CGPCH-E06` | A changed file is governed by a layer but is not in the map. | Error saying the file is GOVERNED and to run the map build first. | Outside the map no gate confronts it: a new governed file would pass with no spec, feature or test. |
| `CGPCH-E07` | A changed path exists neither on disk nor in the map. | Error asking to check the path. | A typo must not become "not governed" and let the hook continue in silence. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/gate.go` | `RunWithWaiver` | gate — running the gates under the waivers |
| DEP9 | `internal/gate/profile.go` | `Aggregate` | gate — the verdict profile |
| DEP2 | `internal/gate/rule.go` | `ParseWaiver`, `WaiverFromMessage` | gate — the waivers and their reasons |
| DEP3 | `internal/mapx/stamp.go` | `StampEdges` | mapa — stamping the confronted edges |
| DEP4 | `internal/issue/issue.go` | `Open`, `OpenAt`, `Resolve`, `ReconcileViolations` | infra — the issue folders |
| DEP5 | `internal/queue/queue.go` | `Enqueue`, `List`, `Drop` | infra — the judgment queue |
| DEP6 | `internal/scan/scan.go` | `Classify`, `Walk` | scan — what the project governs |
| DEP10 | `internal/scan/progress.go` | `IsProgressFile` | scan — the plan's progress companion |
| DEP7 | `cmd/anchors/mapcmd/map_staleness.go` | `StaleMapNodes` | comando — the nodes edited after the map build |
| DEP8 | `internal/config/config.go` | `Load` | config — gates, layers and workflow mode |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
