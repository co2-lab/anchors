---
title: The CLI
---

A single CLI for the Anchors framework, written in Go. It's the tool an AI
operates to exercise the cycle — the AI doesn't need to know Anchors by
heart, it asks the binary (`anchors guide`), learns the flow, and operates
with the commands. Anchors **does not embed AI**: it is the tool that the AI
uses, in any client (Claude Code, GPT, Gemini…).

## Installation

```sh
git clone https://github.com/co2-lab/anchors.git
cd anchors/cli
go build -o anchors ./cmd/anchors
./anchors --help
```

## Use in a project

```sh
anchors init            # configures anchors.yaml (via Q&A; suggests stack presets)
anchors map build       # builds the dependency map from the files
anchors doctor          # ecosystem health: orphans, collisions, coverage holes
anchors check --all     # runs the quality gates; opens issues; stamps the map
anchors report all      # generates the reports in docs/anchors/
```

The CLI only reads **text** — it never parses code. The annotations it
understands (identity, `@noPropagation`…) live in comments; the
per-language markers are configurable. This is what makes it stack-agnostic.

## The flow, driven by the AI

The AI reads `anchors guide` and operates this cycle — the watcher
**enqueues** the work, the AI **pulls** from the queue (the conversation
never gets stuck):

```
  plan ──▶ specify ──▶ map ──▶ implement ──▶ test ──▶ confront
    │          │          │          │           │          │
plan guide  .spec.md   map build  code+feature  tests   check / doctor
                                                              │
                                        issue ◀── divergence the AI doesn't resolve
```

Every saved file makes the watcher enqueue the next task (spec→implement,
feature→test…). The AI doesn't need to remember what comes next; the queue
tells it.

## What the CLI does today

| area | commands | what it delivers |
|---|---|---|
| **Structure** | `init` | configures `anchors.yaml`; structure presets for ~17 stacks |
| **Map** | `map build`, `map show`, `governs` | the dependency graph; who governs whom |
| **Propagation** | `impact`, `stale` | the wave of a change; what became out of date |
| **Queue** | `watch`, `queue`, `next`, `done`, `drop`, `reclaim` | the background watcher enqueues; the AI pulls |
| **Quality** | `check`, `judge`, `doctor` | deterministic gates **and AI-judgment gates**; systemic health |
| **Identity** | `code` | generates/validates a unique scenario code (avoids collisions) |
| **Confidence** | `ingest`, `coverage` | ingests JUnit/lcov from the runner; coverage by **scenario**, by the **diff**, and **delta** |
| **Reports** | `report` | 6 perspectives under `docs/`: tests, quality, structure, config, issues, inconsistencies |
| **The AI bridge** | `guide` (+ `guide plan/spec/code/feature/test/guide`) | the playbook and embedded rulers the AI reads to operate |

### What gives confidence in the deliverable

The best indicator of a successful cycle is **not having bugs at the end**.
Beyond requiring tests, Anchors measures their *quality* from the artifact
the runner already produces:

- **by scenario** — does each requirement in the spec (`SPCR-V01`…) have a
  test that **passed**? (semantic, not line coverage)
- **from the diff** — are the lines you **changed** covered? (catches the
  bug in the new line)
- **delta** — has coverage **dropped** since the last measurement? (catches
  regression)

And gates that a script cannot compute ("does this screen respect the
architecture?") become **AI-judgment gates**: the AI reads the guide's
*conformance points*, confronts the target item by item, and the verdict
enters the same mechanics (stamp + issue) — aging if the target changes.

## Project status

Under construction, and honest about it. The **doctrine** of the 6 pillars
is written and reviewed; the **CLI** exercises the whole cycle and has been
validated against a real-world proof of concept — a mobile app with a
serverless backend, with a non-trivial graph.
