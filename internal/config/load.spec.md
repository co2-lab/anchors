<!-- @anchors
  code: LDCFS
  updated_at: 2026-11-19
  layer: config
-->
# Workload — what each kind of run costs, so heavy runs wait for one another

> **Code**: `LDCFS`

## Overview

Heavy runs of one project, or of several worktrees of it, fight over one machine: a project
reported two worktrees' suites running at once at a load of 30 to 40 on 10 cores, and timing tests
breaking under it. `load:` says what each kind of run costs, in real resources, so the queue
(LDQUS) admits a run only when its cost fits in what the machine has free. The costs are the
project's and the same on every machine; the capacity is each machine's own.

```yaml
load:
  levels:                            # optional: light, medium and heavy have defaults
    heavy: { cpu: 6, memory: 8GB }
  wait_limit: 5m                     # default for every task
  ci: false                          # CI skips the queue unless this is true
  tasks:
    mutation:
      default: medium
      wait_limit: 20m
      scopes:                        # the first pattern a run's files match wins
        "packages/backend/**": heavy
        "apps/mobile/**": medium
    test: { default: medium }
    gate: { default: light }
```

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `levels` | a name → `cpu` (positive) and `memory` (a number and B, KB, MB, GB or TB) | no CPU, an unreadable size | this unit (`LDCFS-E01`) |
| `tasks` | `test`, `mutation`, `gate`, each with `default`, `wait_limit`, `scopes` (file pattern → level, in order) | another task, an undeclared level, a pattern that is no glob, another key | this unit and the strict load (`LDCFS-E01`) |
| `wait_limit` | a duration (`30s`, `5m`, `1h`) | anything else | this unit (`LDCFS-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `LDCFS-B01` | `load:` declares levels over the defaults (light 1 CPU / 512 MB, medium 2 CPU / 2 GB, heavy 4 CPU / 6 GB); a run of a task takes the level of the first scope one of its files matches, else the task's default, else medium; it may wait the task's limit, else the project's, else 5 minutes; and CI admits through the queue only with `ci: true`. (`LevelCost`, `LevelFor`, `WaitLimitFor`, `ParseBytes`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `LDCFS-E01` | A level with no CPU or an unreadable memory, a task the queue does not admit, a level no one declares, a scope pattern that is no glob, a duration that is not one, or a key the block does not know. | The configuration does not load, naming the field. | A cost that reads nothing would admit a run the machine cannot hold, or none. |
