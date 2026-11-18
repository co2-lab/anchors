<!-- @anchors
  code: PWSPP
  updated_at: 2026-11-06
  layer: gate
-->
# PipelineWired — the stages the workflows call are the stages the configuration declares

> **Code**: `PWSPP`

## Overview

A stage the pipeline asks about is refused at run time when the configuration lacks it, which is late:
the release is already running. A declared stage no workflow calls is worse, a guard that guards nothing
while the configuration says it does. And a workflow that calls no stage runs unconfronted. This gate
reads the workflows `pipeline.workflows` names (GitHub's by default) for the calls `anchors pipeline
can-run|record|done|run <stage>`, and confronts them with `pipeline.stages` (PCSPP) before anything runs.

It is a question of the project, never of one file: declared with `scope: project`, it gives one
verdict.

```yaml
gates:
  - name: pipeline-wired
    check: pipeline-wired
    scope: project
    on: [spec, code]
    blocking: true
```

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the stages | `pipeline.stages` | none declared | this unit: pending (`PWSPP-B03`) |
| the workflows | the files `pipeline.workflows` matches | none | this unit: pending (`PWSPP-B03`) |
| a call | `anchors pipeline can-run`, `record` (or its alias `done`) or `run`, and a literal stage name | a name built at run time (`${{ … }}`) | not read: it names no stage the text says |

## Effects

| Effect | Description |
| --- | --- |
| `PWSPP-B01` | The gate passes when every call names a declared stage, every declared stage is called and recorded — by `record`, `done` or `run` —, and every workflow calls a stage or is under `pipeline.unguarded`, by its path or its file name; it says how many workflows and stages it read. |
| `PWSPP-B02` | The gate fails naming each gap: a workflow calling a stage `pipeline.stages` does not declare; a declared stage no workflow calls; a stage asked with `can-run` and never recorded, which no later stage could see passed; a workflow calling no stage and not unguarded. |
| `PWSPP-B03` | With no stage declared, or no workflow matching `pipeline.workflows`, the gate is pending, naming the globs. |

## Errors

none — a workflow that cannot be read is not one of the files the gate reads; the gaps it finds are its verdict, not failures of its own.
