<!-- @anchors
  code: DFLSR
  layer: desenho
-->

# Fallible sources — the failure nobody declared

> PLAN, approved on 2026-10-07. It extends the failure family (`failure-declared`,
> `failure-handled`, `failure-logged`) with what it lacked: a source that says a unit CAN
> fail, independent of anyone having written it down.

## The problem

Six detail screens in a real project ignored the error of their data fetch. With no
network, one said "Budget not found", one spun forever, one opened empty (`fetched ?? []`),
and one opened the form with default values — saving it erased the real data. Spec,
feature, test and code agreed with one another, because no spec declared a loading error.

The failure family did not see it, and by construction could not:

- `failure-handled` starts from a failure the spec DECLARES (`-E`) and asks for its handling;
- `failure-declared` starts from handling that EXISTS in the code and asks for its failure;
- `failure-logged` starts from the handling, too.

When the failure is neither declared nor handled, all three are silent. And they read
handling over the whole file: one `catch` anywhere satisfies a unit with five fetches.

## The design

### What can fail — two sources, both declared by the project

1. **The unit's dependencies.** A layer can be marked `fallible: true` in anchors.yaml — the
   layers that reach a network, a database or a device: hooks that read, repositories,
   services. A spec whose Dependencies table names a file of such a layer consumes a
   fallible source. (When the dependency chain of `DNDDP` lands, the `@dep:` flag on each
   import line becomes this source too, with no table to keep.)
2. **The unit's code.** `dialect.fallible_patterns` lists the calls that can fail, each with
   what counts as handling it:

   ```yaml
   dialect:
     fallible_patterns:
       - call: "\\buse(?:Query|SuspenseQuery|Infinite\\w*)\\("      # React Query
         handled: "\\bisError\\b|\\berror\\b|\\bonError\\b|\\bthrowOnError\\b"
         window: 25          # lines after the call where the handling must appear
       - call: "\\bfetch\\("
         handled: "\\.catch\\(|\\bcatch\\b|\\bok\\b"
   ```

   Anchors does not parse the language. The project says what fails and what handles it,
   as it already does with `handle_patterns`.

### What each gate asks — the same family, one question more each

- `failure-declared`: **a unit with a fallible source declares how it fails** — at least one
  `-E` rule —, or waives it with `@no-failure: <reason>` in the spec. The verdict names the
  source: the dependency (`DEP2 useBudget`) or the call and its line.
- `failure-handled`: besides the set-level check it does today, **each fallible call in the
  code has its handling within its window**. "Reads `data` from `useQuery` and never looks at
  `isError`" is exactly this, and so is `?? []` over query data with nothing beside it. A call
  waived with `@no-handle: <reason>` on its line is not charged.
- `failure-logged`: unchanged — the handling it finds now includes the per-call handling.
- For a visual unit, the existing `error-message-declared` then asks what the user sees:
  the declared failure names its message, so "the load failed" cannot be shown as "not found".

All of it stays **informative** by default, as the family is today.

### Defaults per family

A project declaring `dialect.family` gets `fallible_patterns` for its stack without writing
them: React Query (`useQuery` and siblings), SWR (`useSWR`), `fetch`, axios, and Go (a call
whose error result is discarded). A project with its own data hooks adds them — or marks the
hooks' layer `fallible: true`, which covers every screen that depends on them through the
table. `anchors doctor` says when a family is declared and its fallible defaults are in use,
so a new project gets the warning without knowing the check exists.

### The spec and the guide

- The screen preset suggests the four states of a screen that loads: **loading**, **empty**,
  **load error** (distinct from empty and from not-found) and **loaded**, each with its own
  code, and the load-error failure in the Errors section.
- `anchors guide` for screens gets the rule in one paragraph: a screen that loads data has a
  load-error state, and "no data" never stands for "the load failed".

## The phases

| Phase | Delivers | Done when |
| --- | --- | --- |
| `DFLSR-W01` | **Sources.** `fallible: true` on layers; `dialect.fallible_patterns` with `call`, `handled` and `window`; `failure-declared` charges a unit with a fallible source and no `-E`, naming the source, waived by `@no-failure:`. | Spec, feature and test; on clones of MIF and jokenpo, the six known screens are named, and a sample of the rest checked by hand. |
| `DFLSR-W02` | **Per-call handling.** `failure-handled` checks each fallible call's window; `@no-handle:` waives a call. Depends on `DFLSR-W01`. | On a MIF clone, the screens that read `data` without the error are named, with the call's line; the ones that handle it pass. |
| `DFLSR-W03` | **Family defaults and doctor.** The patterns per `family:`; `anchors doctor` says what is in use. Depends on `DFLSR-W01`. | A project declaring the React Query family is warned with nothing else written. |
| `DFLSR-W04` | **Preset and guide.** The screen preset's four states and the load-error failure; the guide's rule. Can run alongside the others. | `anchors new spec --preset screen` writes the four states; the guide is cited by the verdicts. |

Each phase: full triad, `check --all` green, CI green on the three systems, tag, measured on
peer clones before release, and the peers told at every version.

## What NOT to do

- **Do not hardcode a stack in a gate.** The defaults live in the family, as data, and a
  project can replace them; the gate only applies patterns.
- **Do not demand a handling per rule.** The code does not cite the failure's code inside an
  `if`; the gates tie the failure to the unit and the handling to the call, never a rule to a
  line.
- **Do not be born blocking.**

## Open Decisions

none — the approach was approved on 2026-10-07; the details above are the plan's.
