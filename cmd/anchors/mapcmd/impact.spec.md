<!-- @anchors
  code: MPCTI
  updated_at: 2026-11-18
  layer: comando
-->
# Impact — who depends on a file, the tests to run again, and what it is confronted with

> **Code**: `MPCTI`

## Overview

Before or after changing a file, whoever works on it needs the answers the map already holds.
Who depends on it — always the side that depends, whatever the edge type: whoever imports it or
declares `@dep` of it, reads what it writes, captures it or realizes it; for a spec, the code,
feature and tests it rules. Which tests to run again — those reached and those whose evidence
closure holds a reached file, the same closure `evidence-fresh` judges, grouped by the layer of
the suite that runs them. And what it is confronted with: the spec, feature and guides that
govern it (a divergence there is an issue, not a propagation).

By default only the direct dependents show. `--propagate indirect` climbs one level more and
`--propagate cascade` to the end of the chain; each item past the first level carries the path
that ties it to the file, and an item many others depend on says how many it pulls, so the one
who reads can judge whether that branch belongs — `--stop-at <glob>` lists it and climbs no
further. A tie the project marks with no side effect on the tests, an import of types, is listed
and not climbed past. The list of what the files depend on, the other way, is `dependencies`.

The file can be named the way a shell completes it (relative to where the command runs, with
the machine's separator, or absolute) or the way Anchors' own prompts print it (relative to
the project root). Both must reach the same node, and the node id is always written with
forward slashes, because the map is versioned and travels between operating systems.

The command is a query: it opens no issue and changes nothing.

## Effects

| Effect | Description |
| --- | --- |
| `MPCTI-B01` | For a spec, the direct dependents are the artifacts it rules: the code and the test it specifies. |
| `MPCTI-B02` | For a code file, the validate direction lists what governs it: its spec and its guide. |
| `MPCTI-B03` | An empty direction is stated explicitly (nothing depends on it, or nobody governs it) instead of being printed as an empty list. |
| `MPCTI-B04` | A root-relative argument, the same path with the machine's native separator, and the absolute path all resolve to the same node id. |
| `MPCTI-B05` | A relative argument that does not exist under the project root is resolved from the directory the command was run in. |
| `MPCTI-B06` | What the change reaches that a freeze covers whole is listed apart, under its own heading with each freeze, and left out of the parts to redo and to confront. (`printLevels`) |
| `MPCTI-B07` | Without `--propagate` only the direct dependents show, each with how it depends on the file, and a line saying how to climb further; `indirect` adds one level and `cascade` every level, each item with its path back to the file. An item with ten direct dependents or more says how many it pulls; a `--stop-at` glob lists the items it matches and climbs no further from them; an item tied by types alone is marked and not climbed past. A level asked for and empty is said so. (`impactReport`, `tieItems`, `printLevels`) |
| `MPCTI-B08` | The tests to run again — the tests reached, and those whose evidence closure holds the file or a dependent reached by more than types — are grouped by the layer of the suites that run them: the layers their last result names, or the suites whose paths cover them, or "unassigned". `--format files` prints them one per line, `--layer` keeps one layer's, and `--format json` gives the origins, the dependents with level, relation, path and pull, the tests with their layers, and what the files validate against. (`testLayers`) |
| `MPCTI-B09` | More than one file is asked about at once: the dependents, tests and rulers of all of them, none of the files among the results. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file argument | a path relative to the root, relative to the working directory, or absolute, with either separator | a file that is not a node of the map | this unit: fails naming the resolved path |
| the map | the project's map, or one given explicitly | a missing or unreadable map | this unit: fails asking for the map build |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPCTI-I01` | The resolved node id always uses forward slashes, even when the argument names nothing on disk. | resolve a native-separator path that does not exist and verify the result carries no backslash |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPCTI-X01` | The query leaves the map unchanged and opens no issue. | Impact is read before deciding what to do; acting on it belongs to the check and the work commands. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPCTI-E01` | The resolved file is not a node of the map. | The command fails saying the file is not in the map. | There is no edge to follow from a file the map does not know. |
| `MPCTI-E02` | The map cannot be loaded. | The command fails with the load error and a hint to build the map. | Impact is a query over the map, and the fix is to build it. |
| `MPCTI-E03` | `--propagate` is not `indirect` or `cascade`, `--format` is not `text`, `files` or `json`, or a `--stop-at` glob does not compile. | The command fails naming the value and what is accepted. | A depth or a format guessed would answer another question than the one asked. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
