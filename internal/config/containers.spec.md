<!-- @anchors
  code: CNTNR
  updated_at: 2026-09-26
  layer: config
-->
# Containers — what runs separately, and which layers run inside each

> **Code**: `CNTNR`

## Overview

The project's structure declares layers: what each piece is. It does not say what runs
SEPARATELY — the app, the API, the infrastructure, the database — and that is the question
of the container level of an architecture diagram, and the one that decides how many
component diagrams there are. This unit reads that declaration and answers the questions the
architecture documentation and the doctor ask of it: which containers exist, which of them
are ours inside, in which container a layer runs, and which layers run in none.

The distinction cannot be derived from the folder layout: a shared package runs inside the
API, and an infrastructure package runs nowhere. Guessing from paths would be right by
chance in one project and wrong in the next, with nothing saying so. So the unit never
guesses: a layer is in a container only because the project said so.

An external container (the database, a payment gateway, an identity provider) talks to us
and is not ours inside. It appears at the container level, and it gets no component
diagram, because drawing one would claim a knowledge we do not have. Its layers are still
declared, and still claimed by it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project configuration | a loaded configuration, with or without a containers block, or none at all | — | this unit: a missing configuration or block answers with no containers |
| a layer name | any layer name, in any case | — | this unit: the declared names are compared ignoring case and surrounding spaces |
| the existing layers | the layer names the caller found in the project, in the caller's order | — | the caller: this unit does not discover layers, it classifies the ones it is given |

## Effects

| Effect | Description |
| --- | --- |
| `CNTNR-B01` | The declared containers come back as written, with their talks; a missing configuration has none. |
| `CNTNR-B02` | The internal containers are the declared ones without the external ones, in declared order (`InternalContainers`). |
| `CNTNR-B03` | A layer is found in the container that declares it, ignoring case and the spaces around the declared name (`ContainerOfLayer`). |
| `CNTNR-B04` | A layer no container claims has no container, and that is an answer, not an error. |
| `CNTNR-B05` | The layers of an external container are still claimed by it. |
| `CNTNR-B06` | The orphan layers are the given ones that no container claims, in the order they were given (`OrphanLayers`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CNTNR-I01` | A layer is an orphan exactly when it has no container: the two answers never disagree about the same layer. | classifies a set of layers both ways and checks every layer is in exactly one of the two answers |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CNTNR-X01` | With no container declared, no layer is placed by guessing: every layer is an orphan. | The container a layer runs in is not derivable from its folder; an inferred answer would be right by chance and wrong without warning. |

## Errors

none — an undeclared layer or a missing block is information the unit returns, not a failure it handles.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | config — holds the declared containers |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
