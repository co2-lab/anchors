<!-- @anchors
  code: GRMDG
  updated_at: 2026-10-07
  layer: mapa
-->
# GraphModel — the shape of the map file, and when a validated relation goes stale

> **Code**: `GRMDG`

## Overview

The map is the material form of the dependency graph: which file depends on which, with a revision per
file and a validation stamp per relation. This unit fixes its shape — the kinds of node, the types of
relation, and the fields a node, a relation, a stamp and a judgment carry — and the one rule that reads
it: when a validated relation stops being trustworthy.

The key names are part of the file format. The map is committed, read back by later binaries and merged
between branches, so renaming a key is a format change (see the format contract), not a refactor. Every key
is English, like every configuration key of the product: language settings translate what is READ,
never what is WRITTEN. Optional fields that hold nothing are left out, so the committed file does not
grow with empty lines.

A relation goes stale when it was never validated, or when either end has moved to another revision since
it was stamped. Only revisions decide it: re-confronting and finding the same result is not news, so the
stamp's date and verdict play no part in staleness. The same logic applies to a mutation measurement kept
per scope: it is stale when it was measured against another revision of the file — and a measurement with
no recorded revision is not called stale, because accusing it without evidence would ask for a
re-measurement of what may be correct.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a relation to judge | any relation, stamped or not | — | this unit: an unstamped relation is stale |
| the ends of the relation | node ids present in the graph | an id with no node | this unit: a missing end has no revision, so a stamp recorded against a real revision reads stale |
| a mutation scope and the file's revision | any measurement, with or without a recorded revision | — | this unit: no recorded revision is never stale |

## Effects

| Effect | Description |
| --- | --- |
| `GRMDG-B01` | A relation that was never stamped is stale. |
| `GRMDG-B02` | A stamped relation is fresh while both ends keep the revisions recorded in the stamp, and stale as soon as either end moves. |
| `GRMDG-B03` | A mutation measurement of one scope is stale only when it recorded a revision and that revision differs from the file's current one. |
| `GRMDG-B04` | The navigation between screens is an edge type of its own, `navigates-to` (`EdgeNavigatesTo`), apart from the dependency between files: a screen that leads to another does not use it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GRMDG-I01` | The map file's keys are fixed and English: the graph writes version, generated_by, nodes, edges and flow; a relation writes from, to, judgments, type, origin, method, dep and stamp; a stamp writes validated_from_rev, validated_to_rev, changed_at, verdict and gate; a judgment writes gate, verdict, validated_from_rev, validated_to_rev and changed_at; a node writes id, kind, rev, updated_at, code, layer, code_declared, tags, regime, no_propagation, shared_code, needs, parent, upstream, revises, signal and failures. | serialises fully populated values and compares the exact set of keys |
| `GRMDG-I02` | Optional fields that hold nothing are left out: a bare node writes only id, kind and rev, and a bare relation only from, to, type and origin. | serialises bare values and compares the exact set of keys |
| `GRMDG-I03` | The node kinds are spec, feature, test, code, doc, guide, plan, product and flag; the relation types are governs, specifies, covered-by, tested-by, references, depends-on, seeds, needs, realizes and gated-by. | compares the written values of every kind and type |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRMDG-X01` | Staleness reads revisions only: a stamp's date and verdict never make a relation stale or fresh. | Confronting again and finding the same result is not a new fact; the change of an end is. |

## Errors

none — the model is data plus two comparisons; a missing end or a missing revision is normal input answered by B02 and B03.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the model is the base of the package and depends on nothing in it.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
