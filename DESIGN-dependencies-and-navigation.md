<!-- @anchors
  code: DNDDP
  layer: desenho
-->

# Dependencies and navigation — a map of what each file uses, and where each screen leads

> IN PROGRESS — W01 (v0.1.290), W02 (v0.1.291) and W03 delivered; approved on 2026-10-07 with the decisions below. Two chains, one mechanism: **the agent marks in the code, with
> comment flags, what the file uses and where the screen navigates; Anchors confronts the
> marks with the real code and with the spec, and builds the map.**

## The problem

**Dependencies.** Every `import` is a dependency, and the map does not know them. It holds
only what someone declared by hand — the spec's Dependencies table and the `dep: <path>` of
layers without a spec —, and `dependency-honored` confronts in one direction only (the spec
promises; does the code use it?). An import nobody declared shows up nowhere: jokenpo
imported `tokens.ts`, `useToast` and components without declaring them, and everything stayed
green. Without the chain there is no dependency tree, and the VR evidence rules (a capture
goes stale when its hook changes) see only what was declared.

**Navigation.** A screen spec already declares its route and its navigation tables
(`route-declared`), and `route-exists` checks the route is registered. But nothing ties the
**Out** table to the navigation call in the code, nor one screen's **Out** to the other's
**In**. A screen can navigate where its spec does not say, or the spec can promise a
destination the code never reaches, with every gate green. And the app has no navigation map.

## The design

One principle for both: **the comment flag is the declaration; the code is the proof.**
Anchors does not interpret the language. It reads comments, which are text, in whatever
dialect the header already uses. It confronts them with the code through the patterns the
PROJECT declares (`dialect`), like `import_pattern`, `env_read` and `http_status`. Since
format 7 every governed file has its own `code:`, and the flags address files by it, never by
path: renaming or moving a file does not break the chain.

### The flags

| Flag | Where | Says | Example |
| --- | --- | --- | --- |
| `@dep:` | on each `import` line, beside it | this import ties this file to that one — the symbols are the import's own | `import { PALETTE } from './tokens' // @dep: TOKNS` |
| `@used-by:` | comment right above each exported symbol | who uses this symbol | `// @used-by: ARNAA, WLLTW` |
| `@navigates:` | comment on the navigation call's line, or the line before — every call: `navigate`, `push`, `replace`, `goBack`, `reset`, `popToTop`… | this call leads to that screen (or those, for a back navigation that returns to more than one), through that rule | `// @navigates: GLDTG [GLETG-A02]` |
| `@no-dep: <reason>` | on the import line | this import stays out of the chain | `// @no-dep: types only, erased at compile time` |
| `@no-nav: <reason>` | on the call's line | this call is no screen edge | `// @no-nav: closes a modal of this same screen` |

The header `dep: <path>`, which recognized layers use today, is migrated to the flag on each
import line: the relation is explicit where the tie is made, and removing or changing the import
shows at once which flag goes with it.

### The edges

- `depends-on` (exists): also born from `@dep:` by code, and carrying the **symbols** of its import line. Origin
  `declared` when it comes from the flag, `inferred` when it comes from a real import with no
  flag. The map then shows the gap instead of hiding it.
- `navigates-to` (new): screen → screen, with the rule that triggers it and its origin —
  `spec` when it comes from the Out table, `code` when it comes from `@navigates:`.

### The gates

All are born **informative**; a project turns each one blocking once its chain is complete.

**Dependencies**
- `dep-declared`: every real import that resolves to a governed file carries `@dep:` naming
  that file's code, or `@no-dep:`. External packages and the standard library stay out: they have no code.
- `dep-honored`: every `@dep:` names the code of the file its import line resolves to. It
  catches the flag left wrong when the import changed.
- `used-by-declared`: the `@used-by:` of each exported symbol lists exactly the codes of the
  non-test files that import it — the way back of `@dep:`. A test is tied by its `ref:` and is
  not listed. A new importer with no `@used-by` is flagged on the symbol's file.

**Navigation**
- `nav-annotated`: every real navigation call (pattern `dialect.navigation_call`), back and
  reset included, carries `@navigates:` or `@no-nav:`. The call's destination (the route name) resolves through the
  specs' `route:`, and the flag must name the same screen.
- `nav-matches-spec`: each row of the spec's Out table has a `@navigates:` in the unit's code
  (the screen and what it `specifies`), and each `@navigates:` has its row in the Out table.
  The Rule column, when the table has one, matches the flag's `[CODE-X01]`.
- `nav-symmetric`: if A's Out leads to B, B's In lists A, and vice versa.
- `nav-reachable`: every screen is reachable from the declared initial routes
  (`navigation.entry` in anchors.yaml). An orphan screen is flagged.

### The `--fix`

`check --fix` writes what is mechanical and unambiguous:
- the `@dep:` of each resolved import line;
- the `@used-by:` of each export;
- the `@navigates:` of a call whose destination resolves to a single screen.

What is ambiguous it lists and does not write: a call with a dynamic destination
(`navigate(name)`), a back navigation (its destination is wherever the screen was opened
from, which the In table says and the call does not), and a re-exported symbol.

### Resolution, declared by the project

```yaml
dialect:
  import_pattern: "^import .* from ['\"](?P<path>[^'\"]+)['\"]"   # exists (proof-crosses-boundary)
  import_resolve:
    extensions: [.ts, .tsx, /index.ts, /index.tsx]
    aliases: { "@backend/": "packages/backend/", "@/": "apps/mobile/src/" }
  navigation_call: "\\.(?:navigate|push|replace)\\(\\s*['\"](?P<route>[A-Za-z0-9_/-]+)['\"]"
navigation:
  entry: [Splash, Login]          # initial routes, for nav-reachable
```

With no declaration the gate answers "nothing to measure", like the other dialect gates.
`family:` brings defaults (TS/React Native, Next, Go) for projects that do not want to write
the patterns.

### The map

- `anchors map deps <CODE|file> [--up|--down] [--depth N]`: the dependency tree — down is
  what it uses, up is who uses it.
- `anchors map nav [<screen>]`: the navigation graph, of one screen or of the whole app.
- doct: `{{ navigation }}` generates the navigation map page (Mermaid, screens grouped by
  feature, each edge with the rule that triggers it), and `{{ dependencies }}` the per-unit page.

## The phases

| Phase | Delivers | Done when |
| --- | --- | --- |
| `DNDDP-W01` | **Common base.** The scan reads `@dep:` on import lines, `@used-by:`, `@navigates:`, `@no-dep:` and `@no-nav:`; the `navigates-to` edge; `depends-on` with symbols and origin; `import_resolve`, `navigation_call` and `navigation.entry` in the config. | Spec, feature and test for each piece; this repository's map shows the declared edges. |
| `DNDDP-W02` | **Dependency chain.** The gates `dep-declared`, `dep-honored` and `used-by-declared`, informative; the `@dep:` and `@used-by:` fixer, and the migration of header `dep:` to import-line flags; `anchors map deps`. Depends on `DNDDP-W01`. | On clones of jokenpo and MIF, `--fix` closes the chain and the check after it flags nothing caused by the new gates. This repository (Go) with its chain complete. |
| `DNDDP-W03` | **Navigation.** Edges from the spec's In/Out tables and from `@navigates:`; the gates `nav-annotated`, `nav-matches-spec`, `nav-symmetric` and `nav-reachable`, informative; the `@navigates:` fixer; `anchors map nav`. Depends on `DNDDP-W01`; can run alongside `DNDDP-W02`. | On clones of MIF and jokenpo the navigation map matches the app, and each divergence flagged is real (a sample checked by hand). |
| `DNDDP-W04` | **Visualization and guides.** The `{{ navigation }}` and `{{ dependencies }}` doct pages; `anchors guide header` with the flags; `anchors guide navigation`; `family:` defaults. Depends on `DNDDP-W02` and `DNDDP-W03`. | MIF's navigation page generated and readable; the guides cited by `check --fix`. |
| `DNDDP-W05` | **Consumers.** The VR evidence rules follow the full chain (today they follow the declared `depends-on`); `--changed` reaches whoever imports what changed; a change to a screen's Out table stales the e2e flows that pass through it. Depends on `DNDDP-W02` and `DNDDP-W03`. | Measured on a clone: changing a hook stales the captures of the screens that use it, and only those. |

Each phase: full triad, `check --all` green, CI green on the three systems, tag, and any
migration tested on a clone of each peer project before release. The peers are told at every
version.

## What NOT to do

- **Do not parse the language.** Anchors reads text and declared patterns. A parser per
  language would be the engine choosing the project's structure, and it would break on the
  first dialect it does not know.
- **Do not infer navigation from prose.** Only the table, the flag and the call count.
- **Do not address by path.** The code is the address; paths change.
- **Do not be born blocking.** A new gate that turns every project red the next day teaches
  people to switch it off.

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `DNDDP-D01` | Where does a dependency flag live? | On each `import` line (`// @dep: TOKNS`): the relation is explicit where the tie is made, and removing or changing the import shows which flag goes with it. |
| `DNDDP-D02` | Does `@used-by:` list test files? | No: a test is tied by its `ref:`. |
| `DNDDP-D03` | Do back and reset navigations count? | Yes, every navigation counts — `goBack`, `reset`, `popToTop` and the rest; a back navigation names the screens it returns to. |

## Open Decisions

none — the three questions were decided on 2026-10-07 (above).

## What the implementation taught

- **A default import is the module's default.** `import palette from './palette'` binds a
  local name the importer chose; the symbol it uses is `default`, declared by `export default`
  or `module.exports =`, and that is where its `@used-by:` stands.
- **`--fix` carries the evidence.** The chain's fixer flagged 147 files of a peer; a map rebuild
  would have dropped the proof of every one, for comments. `check --fix` now carries each
  repaired file's evidence to its new content, the line-level signals only when no line moved
  (a flag appended to an import line moves none; an inserted `@used-by:` line does).
- **Aliases by longest prefix.** A monorepo maps `@/` to one workspace and `@/backend/` to
  another; one alias map serves both, the longest matching prefix winning.
- **A package of several files is the author's.** An import that resolves to a directory — a
  Go package — may be flagged with the code of any of its files, and the fixer does not choose
  one. A Go repository's chain is therefore not closed by `--fix` alone; flagging package
  imports by their main file is a follow-up.
- **Measured.** On a jokenpo clone: `dep-declared` 101 → 0 and `used-by-declared` 125 → 0
  after `--fix`, with `evidence-fresh` and `tests-pass` unchanged. On a MIF clone: `dep-declared`
  649 → 0 and `used-by-declared` 817 → 27, 953 files keeping their evidence.
- **A statement opens above only when a list of names does.** The reader first climbed from any
  line with a path to the import above it, and a `require(…)` or an `export … from` borrowed
  the names and the flag of the import before it (MIF: phantom symbols `case:`, `member:`, and
  flags read on the wrong import). Only a line that closes a list (`} from '…'`) opens above.
- **Navigation, measured.** On a jokenpo clone (21 screens), `--fix` flagged every call whose
  route names one screen (`nav-annotated` 19 → 11); what remains is the back and reset
  navigations the author names by hand. `nav-matches-spec` and `nav-symmetric` named real
  divergences: an Out table that leads to Login while the code goes to Register; an Out row
  that says "previous screen" instead of naming it; a code navigation to a profile the Out
  table does not declare. Reachability waits on the reset navigations being flagged — the
  entry screen navigates only by `reset`.

