<!-- @anchors
  code: EFTAV
  layer: desenho
-->

# Evidence follows the asserts — a test goes stale for what it asserts, not for the wiring that gets it there

> APPROVED on 2026-10-09 as the rule for every project (the user, on MIF's request).
> "What affects a test is its asserts, not the wiring" (MIF's user). It is the user's own
> decision for navigation (DNDDP-D04: "a flow asserts; it does not pass through"), applied to the
> whole evidence closure.

## The problem

A test's evidence closure — the files whose change stales its proof — followed every edge out
of the test, and out of everything it reached:

- **Navigation.** A screen's `navigates-to` edges took the closure to every screen it leads to,
  and on from there: MIF's ShoppingList capture had InventoryAnalytics and Recurrences in its
  closure of 93 files.
- **The wiring.** A flow depends on its utils — login, navigate-to-inventory —, and the closure
  went on through everything those utils reach: 37 MIF flows had `RootNavigator` in theirs.
A path screen that breaks the path makes the flow fail the next time it runs; it needs no
preemptive stale. MIF counted about 4 hours of Maestro for a header alignment.

## The design

- **The closure never follows navigation.** A navigation the test asserts is in it already, by
  its Out row's revision (EVFRA-B11).
- **A util is wiring: it counts as a file.** A test file the test depends on — a util, a
  sub-flow — enters the closure with the test files it composes, and the walk does not descend
  from it into code.
- **What has no side effect is said by its flag** (the user, 2026-10-09: "precisamos de um
  controle de sideeffect"; the config names flags, never the map's structure). An edge whose
  flag the project lists in `evidence.no_side_effect.flags` — `@navigates` and `@dep[type]` by
  default — is not followed. An import of types only is flagged `@dep[type]: CODE`: `check
  --fix` writes it where the dialect reads a type import (`type_import_pattern`; TypeScript:
  `import type …`), and the author writes it where no pattern can tell (an inline
  `import('…').T` in a type position). `evidence.no_side_effect.sections` names the spec
  sections whose change proves nothing (navigation and change history by default).
- **Existing proofs follow the rule at once.** A proof's stored closure is read through the rule:
  a file it holds that the rule leaves out stales nothing; a file the rule adds is stamped at the
  test's next run.

## What was measured

On a clone of MIF at HEAD (916 tests with a stored closure):

- **25 closures shrink**, 734 files leave them (18,465 → 17,731), none enters.
- A change to `GoalsScreen.tsx` stales **3** flows by either rule: the "122 flows" were counted
  with `anchors impact`, the doctrine's propagation, not with what goes stale.
- A change to `RootNavigator.tsx` stales **150** by either rule: each captured screen imports it
  — `import type { RootStackParamList }`, 87 screens —, and what a screen depends on is in its
  capture's closure.
- A change to the login util stales 176 by either rule: the util is in the closure as a file,
  as MIF asked.

## Lessons

- **The asserted unit with all it reaches is the wrong widening.** Adding the code of the
  test's unit, with what it depends on, grew the closures from 18,465 files to 109,647; it was
  dropped before delivery.
- **The cost MIF counted is mostly type imports.** A screen's `import type` of the navigator's
  route list carries no behaviour, and is what tied 150 flows to the navigator. With
  `@dep[type]`, written by `check --fix` on a clone of MIF (194 imports in 195 files), the 150
  fell to 115; the rest came through a header `dep:` line repeating a type import, and two
  inline `import('…').T` the fixer cannot read as types. With those two edits by hand: **0**.
