<!-- @anchors
  code: GTVCG
  updated_at: 2026-09-28
  layer: config
-->
# GateVocabulary — the list of default gate names, injected into the configuration layer

> **Code**: `GTVCG`

## Overview

Gate names are identifiers: they go into every project's configuration file, into
tutorials and answers, and they are fixed in English. A table that once accepted Portuguese
gate names at load was removed; the conversion became a one-time migration step, and a file
in the old format is refused with a message that says to migrate.

What remains here is a hook. The list of default gates lives in the package that seeds new
projects, and the configuration package cannot import it without an import cycle, because
that package builds its gates from the configuration's types. So the seeding package
registers a source of the list when it initialises, and the configuration layer's tests ask
this unit for it — without it they would have nothing to confront a gate name against. When
no source is registered the answer is simply absent, never a failure. The other hooks the
configuration layer receives by injection (the canonical gate resolver, the slots hook, the
renamed-key hint) belong to the configuration spec, `CNFGO`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the registered source | a function that gives the default gate names, or nothing | — | the seeding package: it registers the source at initialisation; with none, this unit answers with nothing |

## Effects

| Effect | Description |
| --- | --- |
| `GTVCG-B01` | With no source registered, the default gate names are absent, and asking is not a failure (`DefaultGateNamesForTest`). |
| `GTVCG-B02` | With a source registered, the default gate names are the ones it gives. |
| `GTVCG-B03` | The letters of the artifacts that are not specs are fixed and English: a plan's phase is `W`, a flow's step `T`, an action's result `O`, and none is a canonical spec letter with another meaning; the phase letter is among the canonical rule letters, since specs cite phases in `needs:` (`PhaseLetter`, `StepLetter`, `OutcomeLetter`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GTVCG-I01` | The answer always comes from the source registered last: registering again replaces the previous source. | registers two sources in turn and asks |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GTVCG-X01` | The names are not cached: each question asks the registered source again. | The list belongs to the package that registers it; a copy kept here could diverge from it. |

## Errors

none — an unregistered source is answered with nothing, not handled as a failure.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx/default_gates.go` | `RegisterGateNames` (caller) | infra — registers the source of the default gate names at initialisation |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
