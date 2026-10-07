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

- `failure-declared`: **each fallible source of a unit is named by a declared failure** — an
  `-E` rule whose row, or whose row of the rules' uses, cites the source: the name called
  (`useBudget`), or the dependency's `DEPn` or file stem. "At least one failure" is not enough,
  and measuring proved it: the screen that showed "Budget not found" declared exactly that
  failure, and the load that failed was no failure of its. Unanswered sources fail, named by
  file and line or by `DEPn`, unless the spec waives with `@no-failure: <reason>` or closes its
  Errors section with `none — <reason>`.
- `failure-handled`: besides the set-level check it does today, **each fallible call in the
  code has its handling within its window**. "Reads `data` from `useQuery` and never looks at
  `isError`" is exactly this, and so is `?? []` over query data with nothing beside it. A call
  waived with `@no-handle: <reason>` on its line is not charged.
- `failure-logged`: unchanged — the handling it finds now includes the per-call handling.
- For a visual unit, the existing `error-message-declared` then asks what the user sees:
  the declared failure names its message, so "the load failed" cannot be shown as "not found".

All of it stays **informative** by default, as the family is today.

### Defaults per family — known, offered, never imposed

Each family knows what can fail in its stack and what handles it: React Query (`useQuery`
and siblings, `useMutation`), SWR, `fetch` and axios for `ts`; an HTTP request and a database
call for `go`; `requests` for `python`. They are **not applied by themselves**: a project whose
failure gates already block would wake up, the day it updates, with every fetch charged —
measured on MIF, 65 screens, blocking. Instead:

- `anchors init` writes them into a new project's `dialect:` — a project with nothing to
  unlearn starts measured;
- the governance tips (`check --all`, `anchors doctor`) offer them to a project that declared
  none, as the block to copy; the project adopts them, measures with the gates informative,
  then lets them bar — or declines with `fallible_patterns: []`, which silences the tip.

A call whose result is returned (`return useQuery(…)`, `=> useQuery(…)`) hands its failure to
the caller: the handling is the caller's, and the unit's spec still names the failure.

### The spec and the guide

- The screen preset suggests the four states of a screen that loads: **loading**, **empty**,
  **load error** (distinct from empty and from not-found) and **loaded**, each with its own
  code, and the load-error failure in the Errors section.
- `anchors guide` for screens gets the rule in one paragraph: a screen that loads data has a
  load-error state, and "no data" never stands for "the load failed".

## The phases

| Phase | Delivers | Done when |
| --- | --- | --- |
| `DFLSR-W01` | **Sources.** `fallible: true` on layers; `dialect.fallible_patterns` with `call`, `handled` and `window`; `failure-declared` charges each fallible source no failure names, waived by `@no-failure:`. | Spec, feature and test; on a MIF clone with its data hooks declared, the screens of the incident are named, and a sample of the rest checked by hand. |
| `DFLSR-W02` | **Per-call handling.** `failure-handled` checks each fallible call's window; `@no-handle:` waives a call. Depends on `DFLSR-W01`. | On a MIF clone, the screens that read `data` without the error are named, with the call's line; the ones that handle it pass. |
| `DFLSR-W03` | **Family defaults, offered.** The patterns per `family:`, seeded by `anchors init` and offered by the governance tips; a returned call is the caller's. Depends on `DFLSR-W01`. | A new ts project starts with them; an existing one is offered them, and nothing changes until it adopts them. |
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

## What the implementation taught

- **Per source, not per unit.** W01 first asked for "at least one failure", and a MIF clone
  answered: the screen of the incident declared one — "not found" — and passed. A failure
  answers a source only when it names it.
- **Offered, not imposed.** W03 first applied the family's patterns by themselves; MIF, whose
  failure gates already blocked, would have had 65 screens barred the day it updated. A new
  project gets them from `init`; an existing one adopts them.
- **Hook names, not file names.** A project declaring its data hooks lists the names called
  (`useGoal`), which are not always the file's (`useGoals.ts`); the fallible layer and the
  dependency table cover the files.
