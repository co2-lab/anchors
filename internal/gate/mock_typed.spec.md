<!-- @anchors
  code: MCTYM
  updated_at: 2026-09-19
  layer: gate
-->
# MockTyped — every test double must DERIVE from the module it replaces

> **Code**: `MCTYM`

## Overview

Confronts a test against the hole that is the most treacherous of the whole gate family,
because it is not the absence of proof — it is FALSE proof. **A test that doubles its
neighbour keeps passing after the neighbour changes signature, return or name:** the double
became a frozen copy of a contract that no longer exists. The green certifies the old
version, and nobody goes looking for a defect where there is a green test.

The neighbouring gates do not reach it, and each of them answers "yes" about a test that
lies:

| gate | what it asks | its answer |
| --- | --- | --- |
| `triad-complete` | does the test EXIST? | it does |
| `feature-test-match` | does the scenario match a test case? | it does |
| `tests-green` | did the run pass? | it did |

The defence is not the gate reimplementing type checking: it is demanding that the double be
TIED to the original by a mechanism the language already knows how to check. In TypeScript,
annotating the factory with a partial of the real module's type makes the compiler accuse
both a non-existent method and a divergent return; in Python it is autospec; in Go it is the
interface. **The framework does not know which — the project declares it in
`derived.mock_contract`, and this gate charges the presence of the tie.**

**What separates it from `mock-stamped`:** this one solves the problem where the LANGUAGE
helps, and only there. `mock-stamped` is the agnostic half, which hashes text and works where
there is no structural type to lean on.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | only `mapx.KindTest` is judged | spec, code, feature — the double lives in the test, and charging elsewhere would accuse the wrong file | this unit: a kind that is not a test returns Skip |
| the tie shape | whatever the project declares in `derived.mock_contract`, literal or with the module placeholder | any inferred default — there is none | the project's Structure; an undeclared shape switches the gate off rather than guessing an ecosystem |
| the double dialect | the project's declared detection pattern, or the factory form of the JS/TS runners as a fallback | a double written in a form no pattern describes | the project's Structure, via `derived.mock_detect` |
| the doubled module | a specifier that resolves to a node of the map — alias, relative path or bare name | a third-party library, which the graph does not govern | this unit, resolving against the graph instead of a prefix list in the config |
| the map | a built graph | — (with no graph nothing is governed, and nothing is charged) | this unit: what does not resolve is not charged |

## Effects

| Effect | Description |
| --- | --- |
| `MCTYM-B01` | A double with NO tie fails, and the verdict names the loose module — this is the false proof the gate exists to catch. |
| `MCTYM-B02` | A double whose factory carries the declared tie passes: the annotation is what makes the compiler check name, signature and return against the real module. |
| `MCTYM-B03` | The charge is per MODULE, not per file: one loose double among several is enough to fail, and the verdict counts only the loose ones. |
| `MCTYM-B04` | A double with no factory is not charged: the automock derives from the real module by construction, so it cannot drift, and charging it would be noise one learns to ignore. |
| `MCTYM-B05` | Without the tie shape declared by the project the gate goes quiet, naming what has to be declared. |
| `MCTYM-B06` | An artifact that is not a test leaves without a verdict. |
| `MCTYM-B07` | A test that doubles nobody leaves without a verdict — Skip and not Pass, because nothing was checked and a Pass would inflate the count of greens with nothing. |
| `MCTYM-B08` | Another runner of the same ecosystem is recognised the same way: the double is the double regardless of which library spells it. |
| `MCTYM-B09` | A third-party library double is not charged. |
| `MCTYM-B10` | The third-party exemption is not an escape hatch: one OWN loose double among third-party ones still fails, and only the own one is named. |
| `MCTYM-B11` | A relative import that resolves to a node of the map is an own module like any other — the criterion is the graph, not the shape of the specifier. |
| `MCTYM-B12` | Another ecosystem's dialect is charged the same way once declared, with its own detection pattern and its own tie shape. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MCTYM-I01` | What is governed is decided by the GRAPH, never by a prefix list in the config. That asks for no new configuration, assumes no alias convention — which varies per ecosystem — and follows the project on its own: code that is born enters the map and starts being charged. | doubles a module through an alias, a relative path and a bare name, and verifies all three are charged because all three resolve |
| `MCTYM-I02` | An undeclared tie shape SKIPS instead of guessing one. Inferring the TypeScript form would assume the ecosystem and report GREEN over what was never checked in any other — the worst possible failure in a measuring device. | runs with no declared shape over a file carrying a loose double and verifies Skip |
| `MCTYM-I03` | The verdict COUNTS the loose doubles, not just names them. Naming one without saying how many leaves the reader unable to tell whether the others are on the list too. | confronts a file with one tied double and one loose, and verifies the count is one |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MCTYM-X01` | Does not check whether the annotated type actually MATCHES the real module. | Whoever checks that is the language's compiler, which already does it and does it better. Here the ruler is that the tie WAS WRITTEN — deterministic, and enough, because writing it is what hands the check to the compiler. |
| `MCTYM-X02` | Does not cover drift of BEHAVIOUR — only drift of SHAPE: name, signature, return type. | If the real module starts throwing in one case, or changes semantics while keeping the signature, the double goes on lying and no compiler sees it. For that remainder there is judgement and integration testing at the edge; promising more than the shape would be selling a green this gate does not give. |
| `MCTYM-X03` | Carries no built-in tie shape and no built-in ecosystem. | The shape comes from `derived.mock_contract` and the dialect from `derived.mock_detect`, both in the project's Structure. The gate is language-agnostic BY DESIGN: in Go there is no call to detect at all, because the double is a satisfied interface, and the right verdict there is that the gate does not apply — not a green over an unchecked file. |
| `MCTYM-X04` | Does not charge third-party library doubles. | The drift it pursues is "the neighbour changed and the double did not know", and the neighbour that changes every week is the own module: an external dependency has its version pinned in the lockfile, and its double usually swaps a component for a stub instead of reproducing a contract. There is a harder, measured reason: in the reference app the gate accused 305 files at once, 245 of them third-party doubles. A gate that accuses everything is not read — it is switched off, and it takes the legitimate findings with it. |
| `MCTYM-X05` | When the tie is a type on the factory (the form carries `{{module}}`), does not charge a double with no factory — an automock (no second argument) or an options object such as Vitest's `{ spy: true }`. When the tie is an option on the call (`autospec=True`), the bare call is still charged: Python's `patch` without it is a `MagicMock`, not the module. | Both derive from the real module BY CONSTRUCTION: the automock mirrors its exports and the spy mock IS the module with spies, so there is no hand-written contract to go stale and nothing a `Partial<typeof X>` could bind. `derived.mock_detect` says what a DOUBLE is — a project declares it for `mock-stamped` too — and it must not decide what needs a type: measured in the project that declared it, 38 tests failed at once on 51 doubles, 42 automocks and 9 spy mocks, not one hand-written factory. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindTest` | core — jurisdiction starts from the node's kind, and the `Graph` is what decides which modules the project governs |
| DEP2 | `internal/config/config.go` | `Config` | core — the tie shape and the double dialect are declared by the project, never assumed by the gate |
| DEP3 | `internal/i18n/i18n.go` | `T` | core — the verdict names the loose modules in the reader's language |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
