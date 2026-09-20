<!-- @anchors
  code: GTENG
  updated_at: 2026-09-19
  layer: gate
-->
# GateEngine — which gates reach which node, and what the run concludes

> **Code**: `GTENG`

## Overview

This is the ENGINE. It answers two questions and nothing else: **which gates apply to
which node**, and **what a run over a set of nodes concludes**. Every ruler belongs to
somebody else — the checker measures, the registry routes the name, the rule unit names
the verification and reads the waiver. What the engine owns is routing and aggregation,
and both have measured failure modes of their own.

**Routing.** A gate reaches a node when the node's kind is in what the gate declares, AND
the node carries none of the excluded labels, AND — when the gate names labels — the node
carries at least one of them, AND — when the gate demands a mark in the content — the
target's text contains it. The ORDER inside that conjunction is the decision, not the
implementation: **exclusion comes BEFORE the positive filter**. Layers carry transversal
labels alongside their own, so a node that matches both a named label and an excluded one
must stay OUT. If the positive filter won, declaring the exception would have no effect
exactly where it matters.

The content filter exists for a measured reason of its own: without it, a judgment gate
declared over specs queues an AI question for EVERY spec in the project, and the pending
counter stops measuring the work and starts measuring the size of the repository. And an
UNREADABLE target does not apply: better to stop charging than to charge blind against a
target whose content nobody knows.

**Aggregation, and the four verdicts that are not two.** A run concludes with a
promotion decision, and the engine's job is to keep four distinct states from collapsing
into pass/fail:

- a gate whose required BINARY is missing steps aside — never fails. The gate did not
  measure, and "I did not measure" is neither "clean" nor "dirty". Failing would say the
  project violated something when what is missing is the tool.
- a JUDGMENT gate does not compute at all. It asks whether somebody already answered, by
  reading the stamp a judgement run left in the map. Without that lookup the gate asked
  forever: the verdict was written, the next run asked again, the pending counter never
  came down, and work already done was invisible.
- a pending item that says "there is a decision STILL TO TAKE" bars promotion; a pending
  item that says "I had nothing to confront" does not. Measured in a real repository:
  treating them alike failed 411 nodes at once, and 410 of them were gates with no signal
  ingested.
- a debt that is ASSUMED — a known duty with a declared moment to be paid — becomes
  recorded work; the other pending items become no issue at all.

**What separates this unit from its neighbours.** The registry knows which function
answers a name; this one never looks inside a checker. The rule unit knows how a
verification is named and whether it was waived; this one only ASKS it, per target, and
turns the answer into a verdict that carries the written reason. And the waiver by target
is applied HERE rather than by filtering the gate out of the list, because a waiver
restricted to codes must leave the gate RUNNING to confront everybody else.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared gates | what the Structure declares, or nothing | — (no gate is a case, not an error) | the Structure: what is not declared is not run |
| the nodes | any slice of the map, from one node to the whole project | — | the caller, which decides the slice |
| the project root | a path, used to read target content and to invoke external tooling | — | the caller |
| the map | a built graph, or none | — (an absent map is a case: relational checkers answer undetermined) | this unit: it passes what it has, and never fabricates a graph |
| the Structure | a loaded configuration, or none | — | this unit: the relational checkers treat absence as "no mapping declared" |
| the waiver | what arrived by flag, environment or commit message, or nothing | — | the rule unit, which parsed and validated it |
| the sweep kind | whether the slice is the WHOLE project or an incremental cut | — | the caller, and only it knows which it asked for |

## Effects

| Effect | Description |
| --- | --- |
| `GTENG-B01` | A gate reaches a node only when the node's kind is among the kinds the gate declares. |
| `GTENG-B02` | A gate that names labels reaches only the nodes carrying at least one of them. |
| `GTENG-B03` | A gate that excludes labels never reaches a node carrying one, and ONE excluded label is enough. |
| `GTENG-B04` | Exclusion wins over the positive label filter, because layers carry transversal labels and the exception must hold exactly where it matters. |
| `GTENG-B05` | A gate demanding a mark in the content reaches only the targets whose text carries it — otherwise the pending counter measures the size of the project, not the work. |
| `GTENG-B06` | An UNREADABLE target does not apply: better to stop charging than to charge blind. |
| `GTENG-B07` | A gate with no applicable target does not run at all, so a commit touching one document does not fire a whole-project check. |
| `GTENG-B08` | A gate whose required binary is absent steps aside and never fails — the gate did not measure, and the missing piece is the tool. |
| `GTENG-B09` | A waiver naming THIS target spares only this node, with the written reason in the verdict, and the gate keeps confronting the others. |
| `GTENG-B10` | An aggregate-scope gate runs ONCE and reports a single verdict against the scope itself, not against one of the files. |
| `GTENG-B11` | A judgment gate emits a request for judgment when nothing has answered it yet. |
| `GTENG-B12` | A judgment gate reads the stamp left by an earlier judgement and turns it into the verdict, so work already done stops being invisible. |
| `GTENG-B13` | A judgement recorded as waived becomes Skip and never Pass, because the gate did not measure and Pass would assert an approval nobody gave. |
| `GTENG-B14` | A gate declaring neither a command nor a check answers undetermined, naming the omission. |
| `GTENG-B15` | A pending item that says "there is a decision still to take" bars promotion, and a pending item that says "I had nothing to confront" does not. |
| `GTENG-B16` | Only the obligations gate produces ASSUMED DEBT, and only that pending item carries a deadline into the record. |
| `GTENG-B17` | The engine reconfigures the code grammar from the project's vocabulary before running anything. |
| `GTENG-B18` | `Run` is the entry point that confronts the gates with the map alone, delegating with no Structure — the relational checkers then read absence as "no mapping declared". |
| `GTENG-B19` | `RunWithConfig` is the same entry point carrying the Structure, which the relational checkers need to read the regimes and the surfaces of the triad. |
| `GTENG-B20` | `RunFull` is the one that also knows whether the sweep is the WHOLE project, which is the only thing that lets a gate able to sweep on its own run ONCE instead of receiving thousands of targets in batches. |
| `GTENG-B21` | `RunWithWaiver` is the one that honours a waiver BY TARGET, which cannot be served by filtering the gate out of the list. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GTENG-I01` | A waiver by target is applied while the gate RUNS, never by removing the gate from the list. Removing it would erase the ruler for the whole repository, and a defect elsewhere would pass along. | waives one target, runs over two nodes, verifies the other is still confronted |
| `GTENG-I02` | Skip, Pending and Fail are three different answers and never collapse into two. Each says something the others do not: the gate does not apply, the gate could not measure, the gate measured and the target failed. | runs a gate with a missing tool and one that measured and failed, and verifies the verdicts differ |
| `GTENG-I03` | The verdict of a waived target is Skip WITH the reason written, never silence. It leaves the failure tally without leaving the report — that is the difference between waiving and hiding. | waives a target and verifies the verdict carries the declared reason |
| `GTENG-I04` | The reported target of an aggregate gate is the SCOPE, never one of the files. Blaming one of many files for a verdict about the set would be a statement the engine cannot support. | runs an aggregate gate over several nodes and verifies the reported target |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GTENG-X01` | Does not decide WHETHER a target is correct. | The ruler belongs to the checker or to the external tool. An engine that judged would put the measurement inside the router, and the same defect would have two owners. |
| `GTENG-X02` | Does not compute the verdict of a judgment gate. | The question is judgment, and the CLI cannot answer it. All the engine knows is whether SOMEBODY already answered — and that answer lives in a stamp in the map, carrying the revisions of both ends so it ages when the target changes. |
| `GTENG-X03` | Does not invent a map, a configuration or a waiver when it receives none. | Absence is a case, not an error. Fabricating one would make the run answer about a project state that does not exist, and the verdict would be about nothing. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Gate` | core — what a gate declares is where routing reads its filters and its scope |
| DEP2 | `internal/mapx/model.go` | `Node` | core — the kind, the labels and the identity that routing matches against |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — every verdict the engine writes itself is a translated message, never a fixed phrase |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
