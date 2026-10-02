---
title: "CLI Commands Reference"
description: "Comprehensive searchable reference index of all commands available in the Anchors CLI."
---

The **Anchors CLI** provides an integrated toolchain for managing the lifecycle, quality gates, dependency graph, and AI agent coordination.

Below is the complete reference of all commands organized by functional domain.

---

## 1. Quality & Verification Pipeline

Commands used to evaluate code, verify quality gates, and assess project health:

| Command | Usage | Description | Detailed Guide |
| :--- | :--- | :--- | :--- |
| **`check`** | `anchors check [--all]` | Runs declared quality gates against staged or modified files. | [Guide: check](/docs/cli/commands/check/) |
| **`verify`** | `anchors verify` | Runs EVERYTHING required for the phase: internal gates + external tools. | [Guide: check](/docs/cli/commands/check/) |
| **`doctor`** | `anchors doctor [--verbose]` | Complete X-ray of project health, orphaned artifacts, and systemic drift. | [Guide: doctor](/docs/cli/commands/doctor/) |
| **`audit`** | `anchors audit <file>` | Generates a focused dossier of pending items for a single file or module. | [Guide: doctor](/docs/cli/commands/doctor/) |
| **`test`** | `anchors test` | Executes declared test suites and binds raw test results to graph nodes. | [Guide: check](/docs/cli/commands/check/) |
| **`mutation`** | `anchors mutation` | Runs mutation testing engines and ingests mutant survival scores. | [Gate: mutation-score](/docs/gates/mutation-score/) |
| **`coverage`** | `anchors coverage` | Reports test coverage by scenario, line, and diff delta. | [Gate: line-coverage](/docs/gates/line-coverage/) |
| **`judge`** | `anchors judge <gate>` | Records AI evaluation verdicts for synthetic judgment gates. | [AI Judgment](/docs/concepts/ai-judgment/) |
| **`review`** | `anchors review <target> --gate <g> --by <who>` | Records a second look at a target of a reviewed gate, with its findings; `--pending` lists what is to review. | [Guide: review](/docs/cli/commands/review/) |

---

## 2. Graph & Traceability

Commands that inspect, modify, and calculate relationships in `anchors.graph.yaml`:

| Command | Usage | Description | Detailed Guide |
| :--- | :--- | :--- | :--- |
| **`map`** | `anchors map <build\|show>` | Operates the dependency graph; inspects nodes and directional edges. | [Guide: map](/docs/cli/commands/map/) |
| **`impact`** | `anchors impact <target>` | Computes the downstream impact wave if the given file or spec changes. | [Guide: map](/docs/cli/commands/map/) |
| **`stale`** | `anchors stale` | Lists stale edges: dependencies whose upstream changed without reconfrontation. | [Propagation](/docs/concepts/propagation-and-impact/) |
| **`code`** | `anchors code <layer>` | Generates the next sequential unique identity code (e.g. `AUTH-B03`). | [Traceability](/docs/concepts/traceability-and-codes/) |
| **`recode`** | `anchors recode <old> <new>` | Atomically renames an identity code across specs, features, tests, and code. | [Guide: map](/docs/cli/commands/map/) |
| **`stamp`** | `anchors stamp <file>` | Writes missing `@contract` verification stamps on test doubles. | [Gate: mock-stamped](/docs/gates/mock-stamped/) |
| **`failures`** | `anchors failures` | Audits observed test failures that are not yet cataloged in specifications. | [Gate: failure-declared](/docs/gates/failure-declared/) |
| **`compliance`**| `anchors compliance` | Evaluates regulatory status: how many nodes subject vs. how many comply. | [Gate: obligation-honored](/docs/gates/obligation-honored/) |
| **`governs`** | `anchors governs` | Displays which files each doctrine guide governs across the codebase. | [Product Doctrine](/docs/concepts/doctrine/) |

---

## 3. AI Agent Flow & Task Queue

Commands that power continuous AI pair programming via background queues:

| Command | Usage | Description | Detailed Guide |
| :--- | :--- | :--- | :--- |
| **`watch`** | `anchors watch` | Background filesystem watcher that detects edits and enqueues tasks. | [Guide: flow](/docs/cli/commands/flow/) |
| **`queue`** | `anchors queue` | Lists active tasks in the queue waiting to be processed by a worker. | [Guide: flow](/docs/cli/commands/flow/) |
| **`next`** | `anchors next` | Claims and pulls the next scheduled task from the queue. | [Guide: flow](/docs/cli/commands/flow/) |
| **`work`** | `anchors work <target>` | Emits an optimized instruction prompt for a specific task stage. | [Guide: flow](/docs/cli/commands/flow/) |
| **`deliver`** | `anchors deliver` | Records stage delivery and triggers automated peer confrontation. | [Guide: flow](/docs/cli/commands/flow/) |
| **`done`** | `anchors done <task-id>` | Closes claimed tasks and moves them into historical archive. | [Guide: flow](/docs/cli/commands/flow/) |
| **`discard`** | `anchors discard <task-id>`| Takes a card off the board without permanently deleting its metadata. | [Guide: flow](/docs/cli/commands/flow/) |
| **`drop`** | `anchors drop <task-id>` | Immediately drops a task from the queue without archiving. | [Guide: flow](/docs/cli/commands/flow/) |
| **`reclaim`** | `anchors reclaim` | Returns tasks claimed by dead or disconnected workers back to the queue. | [Guide: flow](/docs/cli/commands/flow/) |
| **`suggest`** | `anchors suggest` | Lists, applies, or rejects proposed automated code fixes. | [Guide: flow](/docs/cli/commands/flow/) |
| **`escalate`** | `anchors escalate` | Opens an issue when a task requires human intervention or spec change. | [Guide: flow](/docs/cli/commands/flow/) |
| **`unblock`** | `anchors unblock` | Opens a work card to resolve an item parked in `needs-user` status. | [Guide: flow](/docs/cli/commands/flow/) |
| **`decided`** | `anchors decided <card>` | Releases an escalated card after a human decision has been made. | [Guide: flow](/docs/cli/commands/flow/) |

---

## 4. Setup, Operations & Governance

Commands for initializing, configuring, and freezing the repository:

| Command | Usage | Description | Detailed Guide |
| :--- | :--- | :--- | :--- |
| **`init`** | `anchors init` | Interactive setup wizard to configure `anchors.yaml` and discover layers. | [Guide: init](/docs/cli/commands/init/) |
| **`new`** | `anchors new <kind>` | Generates skeleton files (spec, feature, test, plan) per layer standards. | [Guide: init](/docs/cli/commands/init/) |
| **`freeze`** | `anchors freeze` | Freezes the repository: stops all work until explicitly thawed. | [Guide: freeze](/docs/cli/commands/freeze/) |
| **`thaw`** | `anchors thaw` | Unfreezes the project and resumes normal development. | [Guide: freeze](/docs/cli/commands/freeze/) |
| **`install-hooks`** | `anchors install-hooks` | Installs the git `pre-commit` hook to enforce quality gates on stage. | [Guide: init](/docs/cli/commands/init/) |
| **`changelog`** | `anchors changelog` | Generates a clean technical changelog from commits release by release. | [Project Pillars](/docs/structure/) |
| **`commit-msg`** | `anchors commit-msg` | Validates git commit messages against format expected by changelog. | [Gate: header-valid](/docs/gates/header-valid/) |
| **`migrate`** | `anchors migrate` | Upgrades legacy Anchors configuration files to the current syntax. | [The anchors.yaml](/docs/anchors-yaml/) |
| **`settings`** | `anchors settings` | Configures local agent-specific options without polluting repo git config. | [Guide: init](/docs/cli/commands/init/) |
| **`guide`** | `anchors guide` | Prints self-contained machine guides intended for AI agents to digest. | [Core Concepts](/docs/concept/) |
| **`board`** | `anchors board [--live]`| Publishes or serves the live web Kanban board of project state. | [Guide: flow](/docs/cli/commands/flow/) |
