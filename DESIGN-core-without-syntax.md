<!-- @anchors
  code: CWSCR
  layer: desenho
-->

# The core without syntax — what is left of a language in the core moves to the project's dialect

> APPROVED on 2026-10-16. The user's decisions, which replace the proposals in the table below:
> - D01: comments live in `dialect`, in anchors.yaml.
> - D02: a file whose comments nobody declared shows as an ALERT in each gate that cannot read
>   it.
> - D03: the floor's verbs become an example, and the real list is the one in anchors.yaml.
> - D04: `anchors migrate` writes what it can; what it cannot goes in the documentation, in the
>   release notes.
> - D05: all at once, one release.
> - D06: the monitor's runners follow the same rule — an example, the real list in anchors.yaml.
> It follows DESIGN-chain-patterns.md and DESIGN-test-handle-patterns.md, under the same rule
> (the user, 2026-10-15/16):
> - the core reads no syntax;
> - the family's syntax is an example, written by `init` and offered by the doctor;
> - with nothing declared, the gate does not measure;
> - no compatibility to keep.

## What is left

Inventoried in the code on 2026-10-16, by a read-only sweep of `internal/` and `cmd/`. A gate
that already reads only the dialect (test handles, chain, navigation, fallible calls) is not
listed.

### 1. Comments — always applied, with nothing to change them

The core knows the comments of about 25 languages, in about ten places that do not agree:

| Where | What |
| --- | --- |
| `config.CommentMarkers`, `CommentNamed`, `LineCommentFor`, `CommentCloseFor` | A built-in table: extension → markers. "anchors.yaml has no key that extends or overrides it." Unknown extension → `#`. |
| `gate.stripLineComments`, `cutInlineComment`, `commentLine` | Strip `//`, `*`, `/*`, `--`, `#` — minus `#include/#define/#if…` —, whatever the file's language: `#` in TS and `//` in Python are cut too. Quote-aware for `"`, `'`, backtick. |
| `scan.stripComments` | `/* */` and `//`, not quote-aware, not by extension. Used by `extractCodes`: the codes a file owns. |
| `scan/evidence_rev.go` | `ChainFlagOnlyRE`, `chainFlagEndRE`, `chainFlagBlockRE`, `commentLines` (with JSX's `{/* */}`), `headerSpan`. |
| `scan/upstream.go` | The `@anchors` header's openers: `// # /* * -- <!-- (* ;`, a `#!` shebang. |
| `gate.value_anchored`, `code_cataloged`, `mock_stamp_gen.commentPrefixFor` (its own table: `#` or else `//`), `quality/touch.go`, `ops/freeze_scope.go`, `scan.extractNavigates` | Each its own list of openers and closers. |
| `gate.statementStart` | A statement starts after a line ending in `;`, `}` or `)`. |
| `nav_chain` fixer | A line ending in `>` is JSX; the flag is written as `/* @navigates: X */`. |

The header and the chain's flags live inside comments, so this layer is under everything:
the map, the evidence, every gate that reads code, every fixer that writes a flag.

### 2. The family lexicon `DialectFor` still fills by itself

When `dialect.family` names a known family, every field the project left empty is taken from
the family:

| Field | Families | Gates that read it |
| --- | --- | --- |
| `exported_func`, `param_name` | 10 | sibling-guard, pagination-honored, code-cataloged, `stamp` |
| `loop`, `cursor` | 10, 9 | pagination-honored |
| `http_status`, `http_status_dynamic` | ts, go | contract-status-declared |
| `env_read` | 10 — and with **no family**, all of them together | env-declared |
| `handle_patterns`, `log_patterns` | ts, go | failure-handled, failure-logged, failure-declared |
| `tests` (pattern, assertion) | ts, go | feature-test-match, test-has-assertion, test-traceable, examples-match, contract-impact, scenario-coverage, flag-covered |
| `definition` | ts, go, python | test-exercises-unit, test-ref-matches-unit, `ingest` |

There is also a floor applied whatever the family: `set_promise` and `set_slice`. These are
English and Portuguese verbs in camel or snake case, read by pagination-honored.

### 3. Syntax inline in gates (not in the dialect at all)

| Gate | What it knows |
| --- | --- |
| mock-typed, mock-stamped, `stamp` | `jest.mock`/`vi.mock` with an arrow factory (when `derived.mock_detect` is empty), a TS return annotation, JS object literals, `./`/`@scope/` specifiers, a declaration's end at column 0 |
| sibling-guard | The signature ends at the first `{`; default guards (`.filter(`, `??`, `throw/raise/panic`) when `guard_patterns` is empty |
| pagination-honored | Limit names (`limit`, `pageSize`, `take`…), `=`/`:` defaults, type names of TS/Go/Rust/Java/C |
| route-exists | React Navigation's `name="X"`, a TS param list, CDK's `addResource` (when `derived.route_pattern` is empty) |
| proof-crosses-boundary, test-exercises-unit | Besides `dialect.imports`, always any line with `import`, `require`, `from `, `use `, `using `, `#include` |
| contract-status-declared | Go's `http.StatusOK`… names |
| test-has-assertion, `testlist` | Blocks by `{ ( [ do` / `} ) ] end`; literals by `'` `"` backtick `"""` |
| feature-test-match | Jest's `%s`, `$var`, `${…}` placeholders in a title |
| failure-declared, obligation-honored | A JS identifier `[A-Za-z_$][\w$]*` |
| map (`scan`) | Maestro's `runFlow:`/`file:` as `depends-on` edges; a `.yaml` test is a runner script |

### 4. File names of tests and code

- **The map pairs derived files by name.** `mapx.stemOfDerived` strips `.test.ts`, `.spec.tsx`
  and `_test.go`, among others.
- **Several commands and gates find a test by `.test.`:** `health`, `testid_surfaces`,
  `obligation_honored`, `code.go`, `deliver`, `watch`.
- **Some gates guess a spec's code file by extension:** `.ts/.tsx/.go/.py/.js`.

### 5. Out of the gates

The monitor's `DefaultRunners` recognise jest, vitest, go test, pytest, maestro, playwright,
stryker and gremlins by command line and output. The project's `monitor.runners` are read first.

## The design

One rule for all of it: what is a language's or a tool's is declared in the dialect, with a
pattern of the engine (`Pattern`: an expression, or `within`/`each`, narrowed by `files`). The
family's values become `Family…` examples that `init` writes and the doctor offers. A gate
whose field is missing is Pending and names it (`presupposes`).

### 1. `dialect.comments`

    dialect:
      comments:
        - files: '**/*.{ts,tsx,js,jsx,go}'
          line: ['//']
          block: [['/*', '*/']]
          strings: ['"', "'", '`']      # where a marker is not a comment
        - files: '**/*.py'
          line: ['#']
          strings: ['"', "'", '"""']

- One reader replaces the ten (`config.Comments`). It finds a file's comments with quotes taken
  into account, and keeps offsets so callers can blank them.
- One writer, `LineCommentFor`/`CommentCloseFor`, takes its markers from the same entry.
- A file no entry covers has no comments: no header, no flag, no code is cut out of it. A gate
  that needs to read one is Pending and says so.
- JSX's `{/* */}` is the ts entry's first block, and `no_line_after` says the lines whose end
  takes no comment (after a JSX tag): the nav fixer writes its flag in a block there.

### 2. No family lexicon by itself

`DialectFor` fills nothing.
- Every field of section 2 is the project's.
- `FamilyLexicon(family)` is the example, written by `init` and offered by the doctor, as
  FamilyChain is.
- The `set_promise`/`set_slice` floor becomes the example too (D03).
- `env_read` with no family stops being every family's together.

### 3. The inline syntax of section 3 becomes dialect fields, or goes

| Was | Is (v0.1.351) |
| --- | --- |
| mock-typed's `jest.mock` fallback, factory and typed-factory shapes | `derived.mock_detect`, `mock_factory`, `mock_typed`; the ts example (`FamilyMocks`) has them |
| `stamp`'s object-literal keys | the real module's exports the double's call names, its strings and comments blanked by `dialect.comments` |
| sibling-guard's default guards and `{` | `guard_patterns` and `signature_end`; the family's lexicon example has both |
| pagination-honored's limit names and types | `dialect.limits` (patterns with `value`, `param`, `type` groups) |
| failure's statement start (`;`, `}`, `)`) | `dialect.statement_end` |
| route-exists's fallback | `derived.route_pattern` only (`FamilyRoutePattern`) |
| the `import`/`require` lines added to `imports` | gone: `imports` only |
| test-has-assertion's blocks and empty body | `tests.block_open`, `block_close`, `empty_body` |
| `testlist`'s string literals | the pattern's `title` groups, read as captured |
| Jest's title placeholders | `tests.placeholder` |
| a JS identifier filtering imported names | gone: the names are what the project's `import_names` captured |
| a `.yaml` test as a runner script | the layer's `headless: true` |
| Go's `http.Status` names | kept: they are HTTP's reason phrases, captured by the go example's `http_status` |

### 4. Names of tests by the layers

- A file's unit name is its basename without what a layer adds to it: the part after
  `{{name}}` in a `derived.files` template (`config.UnitName`).
- A unit's pieces are the templates rendered for it (`config.PiecePaths`).
- A spec's code is looked for in the extensions the code layers declare (`CodeExtensions`).
- A neighbour test is a file the layers classify as a test.
- No list of extensions or test suffixes remains in the gates or the commands.
- `init`'s inference still recognises the conventions it proposes, as examples.

### 5. The monitor's runners

D06 reversed the proposal. The monitor knows only Anchors' own long commands. Every other tool
is the project's `monitor.runners`, and `config.ExampleRunners` are the examples:
- `anchors migrate` writes the ones its suites run with;
- the agent hook names a run by the project's runners.

### 6. Migration and release notes

- `anchors migrate` crosses **format 8** once (`migra.DeclaredSyntax`). Where the project
  declared nothing, it writes:
  - the comments of the examples covering its files;
  - its family's lexicon (only the naming conventions with no family);
  - the tests' block shape;
  - the route and mock shapes for the gates that read them;
  - `headless` on test layers of YAML;
  - the example runners of its suites.
- It names what no example can say: an extension no example covers, a tests pattern with no
  `title` group, `env_read` with no family.
- A map not migrated is refused, asking for `anchors migrate`.
- The release workflow opens a release with `release-notes/<tag>.md`, where the manual steps
  are.

## Phases

| Phase | What | Proof |
| --- | --- | --- |
| `CWSCR-W01` | **Comments.** `dialect.comments`, the one reader and writer; the ten places moved to it. | This repository's and MIF's map, evidence revisions and check table unchanged with the examples declared; a file of an undeclared extension has no header read. |
| `CWSCR-W02` | **The lexicon.** `DialectFor` fills nothing; `FamilyLexicon`; init and the doctor. | Each gate of section 2 Pending without its field, the same table with the example. |
| `CWSCR-W03` | **Inline syntax.** Section 3's fields. | Each gate's tests on the example; MIF's table unchanged. |
| `CWSCR-W04` | **Test names by the layers.** Section 4. | The map of this repository and MIF's pairs the same files. |
| `CWSCR-W05` | **Migration.** `anchors migrate` writes the family's examples into a project's dialect, once, so its table does not change on upgrade. | MIF's clone migrated: same table as before; the doctor offers nothing. |

## Decisions

| Code | Question | Proposal |
| --- | --- | --- |
| `CWSCR-D01` | Comments: a dialect field, or a separate top-level `comments:` (languages mix in a project)? | **Under `dialect:`**, a list of entries narrowed by `files`, so a TS and a Python app in one repo each get their own. |
| `CWSCR-D02` | A file whose comments nobody declared: read nothing, or keep the built-in table as an example only? | **Read nothing**, and the doctor names the extensions found without an entry, with the example's block. |
| `CWSCR-D03` | The `set_promise`/`set_slice` floor (verbs in two human languages) — an example too? | **Yes**: it is a project's naming convention, not Anchors'. |
| `CWSCR-D04` | `anchors migrate` writes the examples into existing projects (W05), or each project adapts by hand from the doctor? | **Migrate writes them once**, tested in clones of the peers before release. The upgrade then changes no table, and the yaml says what is read. |
| `CWSCR-D05` | Order of delivery: one release per phase, or all at once? | **All at once** (the user, D05), announced before install with its yaml (the MIF rule of 2026-10-16). |
| `CWSCR-D06` | The monitor's `DefaultRunners`: move them to examples too? | **Yes** (the user, D06): the same rule as the rest — examples, the real list in anchors.yaml. |
