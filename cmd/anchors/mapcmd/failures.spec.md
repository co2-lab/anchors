<!-- @anchors
  code: FLRSA
  updated_at: 2026-10-08
  layer: comando
-->
# Failures — the observed failures that the spec has not explained yet

> **Code**: `FLRSA`

## Overview

A spec declares a failure as possible without knowing how it will happen. Once the
application runs, the log has the context that answers it, and the log ingestion binds each
occurrence of a declared failure code to the spec that declares it. This command is the
third layer of that circuit: it lists the failures that happened and that the spec has not
concluded about yet, so that whoever understands the domain has a target (which failures),
material (how many times, when, in which spec) and a place to record the conclusion (the
spec itself, beside the rule).

The command does not analyse the log and does not guess a root cause. A failure leaves the
list when the spec carries one of two conclusions beside the rule: resilient (the failure
is handled, with the reason) or under observation (what was already ruled out). Both stay
visible on demand, because "still unknown, and this is what was ruled out" is knowledge the
next person starts from.

An occurrence was measured against one version of the spec. When the spec changed since the
ingestion, the rule may be another one, and the occurrence is flagged so nobody concludes
about what was not observed.

## Effects

| Effect | Description |
| --- | --- |
| `FLRSA-B01` | The review lists, as open, only the ingested failures whose rule carries no conclusion in the spec with the spec that declares each one. |
| `FLRSA-B02` | The open failures are listed from the most frequent to the least, whatever order the map holds them in. |
| `FLRSA-B03` | With the option to show everything, the failures that carry a conclusion are listed too, each with its resilient reason or what is under observation. Without it they are left out. |
| `FLRSA-B04` | An occurrence whose recorded spec revision differs from the spec's current revision is flagged as measured against an older spec, asking for a new ingestion before concluding. |
| `FLRSA-B05` | When nothing is open (and, with the option to show everything, nothing is concluded either), the review prints the none message instead of an empty list. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | a map whose spec nodes may carry ingested failure occurrences | a missing or unreadable map | this unit: fails with the load error |
| the spec files | the specs named by the nodes that carry failures | a spec deleted or unreadable since the ingestion | this unit: skips that spec |
| the conclusions | a resilient or under-observation mark beside the rule in the spec | a rule with no mark | this unit: lists that rule's failure as open |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLRSA-I01` | What the log ingestion binds to a spec is exactly what the review lists for it: a failure that happened and carries no conclusion is never missing from the list. | closed cycle: ingest a log with the real ingestion, then run the review and find the open failure with its spec |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLRSA-X01` | The review writes nothing: neither the map nor the spec is changed by listing the failures. | The conclusion is the reader's work, recorded by hand in the spec; the command only asks for it, like `judge` asks for a verdict it does not compute. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLRSA-E01` | The map cannot be loaded. | The command fails with the load error. | Without the map there are no ingested occurrences to review. |
| `FLRSA-E02` | A spec that carries failures can no longer be read. | That spec is skipped; the review answers for the others. | Its conclusions cannot be read, and one missing file must not hide every other open failure. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
