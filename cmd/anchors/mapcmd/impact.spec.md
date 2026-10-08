<!-- @anchors
  code: MPCTI
  updated_at: 2026-10-08
  layer: comando
-->
# Impact — what a change to one file reaches, in both directions of the map

> **Code**: `MPCTI`

## Overview

Before or after changing a file, whoever works on it needs two answers the map already
holds. Downwards: which artifacts depend on this one and have to be redone (the change
propagates to them). Upwards: which artifacts govern this one and have to be confronted (a
divergence there is an issue, not a propagation). This command queries the map for a single
file and prints both directions, saying explicitly when one of them is empty.

The file can be named the way a shell completes it (relative to where the command runs, with
the machine's separator, or absolute) or the way Anchors' own prompts print it (relative to
the project root). Both must reach the same node, and the node id is always written with
forward slashes, because the map is versioned and travels between operating systems.

The command is a query: it opens no issue and changes nothing.

## Effects

| Effect | Description |
| --- | --- |
| `MPCTI-B01` | For a spec, the propagate direction lists the artifacts that depend on it: the code and the test it specifies. |
| `MPCTI-B02` | For a code file, the validate direction lists what governs it: its spec and its guide. |
| `MPCTI-B03` | An empty direction is stated explicitly (nothing depends on it, or nobody governs it) instead of being printed as an empty list. |
| `MPCTI-B04` | A root-relative argument, the same path with the machine's native separator, and the absolute path all resolve to the same node id. |
| `MPCTI-B05` | A relative argument that does not exist under the project root is resolved from the directory the command was run in. |

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

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
