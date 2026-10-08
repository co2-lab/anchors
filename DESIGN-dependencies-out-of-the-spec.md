<!-- @anchors
  code: DOOSD
  layer: desenho
-->

# Dependencies out of the spec — declared where the import is, in every artifact that imports

> IN PROGRESS — W01 (v0.1.303), W02 and W03 delivered; approved on 2026-10-08. The spec comes before the code, so it cannot hold what the code
> produces: the files a unit imports and the methods it calls are born with the implementation.
> They are declared where the import happens — in the code, in the tests, in any artifact that
> imports —, by the `@dep:` flag on the import line, and the spec's Dependencies table goes.

## The problem

The spec carries a Dependencies table: `| DEP1 | packages/backend/models/userProfile.ts |
getProfile | dao |`. A spec is written first; when it is born, no file it will import exists,
and the table is empty or invented. When the code is written, the table has to be filled from
it, and every change to an import sends the author back to the spec. The spec's content then
depends on what derives from it — a cycle: the spec is rewritten to follow the code it
specifies.

Since DNDDP the same relation is also declared in the code, by the `@dep:` flag on each import
line, and confronted with the import itself. The two copies drift: on MIF, a spec held two
Dependencies tables, `DEP1` naming one file in the first and another in the second
(DPDCD-W02).

And the chain stops at the code: a test imports what it tests, its fixtures and its helpers,
and that relation is declared nowhere (DNDDP-D02 left the tests out).

## The design

- **The relation lives where the import is.** Every artifact that imports — code, tests,
  test support, e2e flows that compose another flow — carries `@dep:` on each import of a
  governed file, or `@no-dep: <reason>`. `dep-declared` and `dep-honored` confront them in all of
  them; the map's `depends-on` edges come from them alone.
- **`@used-by:` lists who imports a symbol, tests included**, so a symbol a test imports is
  confronted like one a screen imports (DNDDP-D02 revised).
- **The spec's Dependencies table goes.** Nothing in a spec names a file or a method of the
  code; the template no longer writes the section; the spec guide says where dependencies are
  declared instead.
- **A rule's uses name what it reads in the spec's own vocabulary.** `Rule uses` cites the
  fields, states and codes the spec declares; a `DEPn` is no longer a name of anything.
  `contract-impact` follows a changed field to the rules of the dependent units through the
  map's `depends-on` edges from the flags — code to code, each end taken to its unit.
- **A failure names its source by the call or the file** (`failure-declared`), as it already
  may; `DEPn` stops being one of the names.
- **A spec relies on another unit only through the product.** When two units share a concept,
  the rule lives in the product doctrine, which comes first, and each spec realizes it
  (`@realizes`). No spec cites another spec: none has to watch whether the other changed to
  stay valid, only follow the product above both. A spec citing another unit's spec is a
  finding, answered by moving the shared rule up to the product.
- **What the code reaches that is no import is a `@dep` too, with its kind.** A table of the
  database, an external API, a queue, a bucket: flagged where the code calls it, as
  `@dep[<kind>]: <name>` — `// @dep[db]: transactions`, `// @dep[api]: stripe.charges`. The
  kinds are the project's, declared in `anchors.yaml` (`dependency_kinds:`); an import is the
  implicit kind. The map's `depends-on` edge carries the kind, so the same chain answers new
  questions later — which units touch a table, which call an API — without a second mechanism.
- **`dependency-honored` retires**: its question — does the code use what the spec says it
  depends on — has no spec side left; `dep-declared` and `dep-honored` ask it of the code.

## The phases

| Phase | What | Proof |
| --- | --- | --- |
| `DOOSD-W01` | **The chain reaches every importer.** `dep-declared`, `dep-honored` and the fixer on tests, test support and flows; `@used-by:` counting them. | On clones of jokenpo and MIF, `--fix` flags the tests' imports, and evidence, stamps and coverage are unchanged. |
| `DOOSD-W02` | **The spec stops depending on the code.** The template, the spec guide and the presets without the Dependencies section; `Rule uses` and `failure-declared` without `DEPn`; `contract-impact` through the flags' edges; `dependency-honored` retired; the map no longer reads the spec's table. | The impact of a changed field still names the dependent units' rules, measured on the clones against today's answer. |
| `DOOSD-W03` | **Migration.** `anchors migrate` removes the Dependencies tables and rewrites each `DEPn` a rule cites to what it reads; whatever it cannot rewrite it names for the author. | On clones of jokenpo (80 specs with the table) and MIF (278), the migration then `check --all`: no new failure beyond those it names. |
| `DOOSD-W04` | **Kinds of dependency.** `@dep[<kind>]: <name>`, `dependency_kinds:` in `anchors.yaml`, the kind on the `depends-on` edge, `anchors map deps --kind`; the spec-to-spec citation flagged and answered through the product. | On a clone, the database tables and external APIs a project's code calls are flagged by kind and listed by `map deps --kind`. |

## What NOT to do

- Keep the table "for documentation". A copy of what the code says is what drifted.
- Make the spec name a file, a method or a file's code. All three are born with the code.
- Leave the tests out of the chain. They import too, and a moved fixture breaks them as a moved
  hook breaks a screen.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `DOOSD-D01` | Where is a dependency declared? | Where the import is, in every artifact that imports — code, tests, support, flows —; never in the spec, which precedes the code (the user, 2026-10-08). |
| `DOOSD-D02` | May a spec rely on another unit? | Only through the product: the shared rule lives in the product doctrine, which comes first, and each spec realizes it; specs never cite each other, so none has to watch the other to stay valid (the user, 2026-10-08). |
| `DOOSD-D03` | What the code reaches that is no import — a table, an external API, a queue? | A `@dep` too, with a kind the project declares, flagged where it is called: the map gains the kind and later uses without a second mechanism (the user, 2026-10-08). |

## Open Decisions

none

## What the implementation taught

- **W01, measured.** On a jokenpo clone, `--fix` flagged every import of the tests:
  `dep-declared` 88 → 0, `used-by-declared` 45 → 0; on a MIF clone 1746 → 0 and 773 → 0, once
  the forms tests import with were read: an inline import bound whole, a `require(…)`'s member
  and a destructuring's `A: B`. `evidence-fresh`,
  `mock-stamped` and `tests-pass` unchanged on both.
- **A symbol only a test imported gains its first `@used-by:` line**, which moves the lines of
  the file: its line coverage waits for the next run (11 files on jokenpo, 83 on MIF). One
  run per project, the day it adopts the chain on its tests.
- **`typeof import('./x')` brings the module, not its default.** MIF's tests type their mocks
  with it, and reading it as the default asked a `@used-by:` of an `export default` that does
  not exist — 303 symbols no fixer could flag. An import bound whole names no symbol, like
  `import * as`.
- **A migration that rewrites a spec must carry its evidence.** Taking the table out changed
  80 specs of jokenpo, and with them the revision their scenarios were proven at: 73 failed
  `scenario-coverage` until the migration carried what held at each spec's revision, as
  `check --fix` does. Taking a table out proves nothing new.
- **What the migration leaves is the author's, and it says so.** On jokenpo: 80 tables and
  480 `DEPn` out of the rules' uses; 34 specs left with rules whose only use was a `DEPn` — every
  one named by the migration, and no other new failure. The rows that named no file (Cognito,
  a push service) and the data origins citing a `DEPn` are named too: what a datum comes from,
  and an external dependency, are the author's to say (W04 gives the latter a kind).
- **This repository** carried 221 tables; two rows named no file of it (the flag library, git).

