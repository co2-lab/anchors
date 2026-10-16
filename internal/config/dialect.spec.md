<!-- @anchors
  code: DLCTI
  updated_at: 2026-10-16
  layer: config
-->
# Dialect — the lexicon of the project's language, between an agnostic gate and concrete code

> **Code**: `DLCTI`

## Overview

Some gates confront a truth that holds in any language — "a function that promises the
whole set must not return the first page" — but to see it they must recognise things that
are syntax: an exported function, a loop, a pagination cursor, a failure being handled, a
failure being logged. Syntax belongs to the project, not to the tool. This unit is the
bridge: it gives each gate the effective lexicon of the project's language.

A wrong lexicon does not raise an error: it produces silence, and the gate passes green
because it saw nothing. So the lexicon is declared, never assumed. The project names a
family (a known language whose lexicon comes built in), declares patterns of its own, or
both; a declared pattern always wins over the family. Two naming conventions — the verbs
that promise a whole set and the ones that declare a slice — are not language, they are
naming, and they apply as the floor to every project whatever its family. A family the
unit does not know contributes nothing, and the gates that need the lexicon then say what
is missing instead of pretending they checked.

The unit also carries the Gherkin keyword table. The keywords written into a new feature
follow the project's declared language, English by default. The keywords recognised when
READING a feature are those of every language in the table, because a reader that knew only
the declared language would go blind on a feature written in another one — and blind here
means reporting green over scenarios it cannot see.

Finally it answers two small questions the gates ask: whether a declared pattern is usable
(an empty or invalid one is the same as not declared), and whether the project declared that
a lexicon field does not apply to it (an explicit opt-out, read by the field's name as it is
written in the configuration file).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared dialect | a dialect block with a family, explicit patterns, an opt-out list, or none at all | — | this unit: a missing block yields the naming defaults alone |
| the family name | any name, in any case | a name that is not a known family | this unit: an unknown family contributes nothing, and the known families are listed for the gates' messages |
| the Gherkin language | a language code, in any case, or nothing | — | this unit: nothing means English; a code outside the table keeps its code and uses English keywords |
| a pattern to compile | any text | text that is not a valid regular expression | this unit: an empty or invalid pattern compiles to nothing; the configuration load refuses invalid declared patterns before this |
| the field asked about for an opt-out | the field's name as written in the configuration file | the field's name in the code | the caller: it asks with the configuration name |

## Effects

| Effect | Description |
| --- | --- |
| `DLCTI-B01` | A project's dialect is what it declares: no family's lexicon and no naming default fills a field it left empty, a family named alone included (`DialectFor`). |
| `DLCTI-B02` | Each family has an example of the lexicon the gates read — with its statement end, signature end, limits and guards —, handed out as a copy and never applied by itself: `anchors init` and `anchors migrate` write it, the doctor offers it (`FamilyLexicon`). |
| `DLCTI-B03` | The family name is matched ignoring case. (`FamilyLexicon`) |
| `DLCTI-B04` | An unknown family has no example; the known families are listed in alphabetical order for the messages that name them (`KnownDialectFamilies`). |
| `DLCTI-B05` | Every family's example carries the two naming conventions (the set-promise and the set-slice verbs, `NamingExample`); a project declares its own, or none. |
| `DLCTI-B06` | The Gherkin language defaults to English and is looked up ignoring case, answering the table's own spelling of the code, so `PT` is `pt` and `zh-cn` is `zh-CN`; a language outside the table keeps its code as declared and gets the English keywords (`GherkinFor`). |
| `DLCTI-B07` | The recognised ways to open a scenario are those of every language in the table — scenario, outline and synonyms — each once, longest first, so an outline that contains the word for scenario is tried before it (`GherkinScenarioAlternatives`). |
| `DLCTI-B08` | The recognised result keywords are those of every language in the table, in alphabetical order (`GherkinThenAlternatives`). |
| `DLCTI-B09` | An empty or invalid pattern compiles to nothing, which the gates read as "not declared" (`Compile`). |
| `DLCTI-B10` | A lexicon field is waived only when the opt-out list names it by its configuration name, ignoring case and spaces; its name in the code does not waive it (`WaivedField`). |
| `DLCTI-B11` | The default set-promise convention recognises the verb after a provider prefix (a cloud-provider list call) and never inside another word (allocate, callback, enlistment). |
| `DLCTI-B12` | The default set-slice convention is a query verb opening the name, followed by a slice word (first, recent, top, page…). |
| `DLCTI-B13` | The Go family recognises both shapes of error handling: `if err != nil` and `if err := f(); err != nil`. The second is the commoner, and matching only the first left the handling of most Go code invisible to the failure gates. |
| `DLCTI-B20` | Each language family's example says how its code reads an environment variable (`env_read`), capturing the name; what is read is the project's own pattern, and with none `EnvReadPattern` is empty. |
| `DLCTI-B21` | `fallible_patterns` (`FalliblePatterns`, each a `FalliblePattern` with `call`, `handled` and `window` — `DefaultFallibleWindow` lines when it declares none) are the project's; the family's fill them only when the project declares none. |
| `DLCTI-B22` | A family knows what can fail in its stack and what handles it (`FamilyFalliblePatterns`: React Query, SWR, fetch and axios for `ts`, a type argument taken; HTTP and SQL for `go`; `requests` for `python`; each with the `name` group that says the function a failure cites), and does not impose them: a project that declared no `fallible_patterns` is measured with none until it adopts them. |
| `DLCTI-B23` | The dialect declares how an import's path becomes a file (`import_resolve`: `extensions` tried in order, `aliases` of path prefixes, and the `separator` of a path whose directories are not divided by `/` — `ImportResolve`) and what a navigation call looks like (`navigation_call`, its destination in the group `route` — `NavigationCall`); no family gives a navigation call by itself, and the `ts` family has an example of one (`FamilyNavigationCall`). |
| `DLCTI-B24` | The OCaml family reads an exported function as a top-level `let` or `let rec`, an environment read as `Sys.getenv` or `Sys.getenv_opt` with a literal name, and a loop as `for` or `while`; its example of the chain reads a module opened (`open`, `open!`) or qualifying a name, and is never the project's import pattern by itself. |
| `DLCTI-B25` | The import resolution may declare that imports are module names (`by_name`) and how a name becomes a file name (`name_case`: `lower_first`, `lower`, or as written); both are read from `import_resolve`. (`ImportResolve`) |
| `DLCTI-B26` | A family's fallible patterns say, in their `propagates` group, when a call hands its failure to the caller: for `ts` a call `return`ed or an arrow's body, with or without `await`; for `go`'s HTTP and `python`'s `requests`, a call `return`ed. A call kept in a variable propagates nothing. |
| `DLCTI-B27` | `dialect.test_handles` declares the patterns the handle gates read — `exposed`, `mentioned`, `queried`; a pattern that does not compile, or captures with none of its list's groups (`testid`, `prefix`, and `suffix` for `mentioned`), refuses the config, naming it. (`TestHandles`) |
| `DLCTI-B28` | The `ts` family has an example of the test handle patterns for the attribute a project marks its elements with — a conditional handle read as a region, every literal of its branches, nested or not, and none of its condition —; another family, or no attribute, has none; and the example is never applied by itself — a project that declared none is measured with none. (`FamilyTestHandles`) |
| `DLCTI-B29` | A dialect pattern is a regular expression run over the whole content — a match may span lines, and each capture keeps its offset and its line —, or a region (`within`, its `within` group) with what repeats in it (`each`), every capture of one region sharing its start and read as one of a list; `files` narrows it to the files its glob matches, and it is written as text or as a map. A pattern that does not compile, a region missing its half or its `within` group, a capture with none of its list's groups, or a glob that is no glob refuses the config, naming it. (`Pattern`, `Patterns`) |
| `DLCTI-B30` | `dialect.imports` (a `path` group), `import_names` and `declarations` (a `name` or a `default` group) are the patterns the chain gates read, checked on load like any dialect pattern; a family gives none of them, nor its import resolution, by itself. (`Imports`, `ImportNames`, `Declarations`) |
| `DLCTI-B31` | The `ts`, `go`, `python` and `ocaml` families have an example of the chain patterns — `ts`: imports with lists across lines, side effects, `require`/`import()` and a flow's `runFlow:`, the names they bring, the declarations at the top level and export lists, type imports, the extensions; `go`: single imports and import blocks, top-level declarations; `python`: `from … import` and `import`; `ocaml`: modules opened or qualifying a name, resolved by name —; another family has none. (`FamilyChain`) |
| `DLCTI-B32` | The `ts` family has an example of the mock shapes — `mock_detect` for `jest.mock`/`vi.mock`, the contract, a factory passed and a factory typed — and of `route_pattern` — `name="X"`, a stack's param list, `addResource('x')` —; `python` one of `mock_detect` (`patch`); another family none; never applied by themselves. (`FamilyMocks`, `FamilyRoutePattern`) |
| `DLCTI-B33` | `import_resolve.scoped` declares aliases for the importers a glob matches — each workspace's own `paths` —: an importer's aliases are the global ones with those of every entry whose glob matches it over them, the first entry winning a prefix two declare; an entry with no glob, a glob that is no glob, or no aliases refuses the config, naming it. (`AliasesFor`, `ScopedAliases`) |
| `DLCTI-B19` | The Go family also sees an error held in a field and a sentinel error: `if result.Error != nil` and `return ErrNotFound` handle a failure, and `return nil, result.Error` and `return ErrNotFound` propagate it, which records it as `fmt.Errorf` does; `return e.Error()`, which turns the error into text, propagates nothing. |
| `DLCTI-B14` | The Go and TS families' examples say how a test is written, its title in a `title` group — Go by `t.Run(`, TS by `it`/`test`/`describe`, also as `.only`, `.skip` or `.each(table)` —, and a project's own `tests` is read as declared. |
| `DLCTI-B15` | The TS family recognises a `catch` as handling a failure with or without its binding: `catch (e) {` and `catch {`. |
| `DLCTI-B16` | The Go and TS families' examples say what an assertion is (`tests.assertion`). A project's `tests` is read as declared: one declaring only its assertion has no way to open a test, one with its own script no assertion. |
| `DLCTI-B17` | The Go, TS and Python families say how code defines a name (`definition`) — Go's functions, methods and types, TS's functions, classes and constants bound to an arrow function, Python's `def` and `class` — and a project's own is read as declared. A capture group named `owner` is the type a member belongs to, not the name: Go's captures a method's receiver type. |
| `DLCTI-B18` | The recognised keywords that open an outline's examples are those of every language in the table, each once, in alphabetical order (`GherkinExamplesAlternatives`). |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DLCTI-I01` | Every pattern a built-in family or a naming default brings compiles. A built-in pattern that did not compile would be read as undeclared, and the gate would approve in silence. | resolves every known family and compiles each of its non-empty patterns |
| `DLCTI-I02` | The keywords written for any language of the table are among those every reader recognises. | resolves the keywords of each language and looks them up in the recognised scenario and result keywords |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DLCTI-X01` | No family brings a collection query: that pattern is always the project's to declare. | The query that returns many records is vendor-specific (a DynamoDB project and a SQL project look nothing alike); a built-in guess would name one vendor's API and be blind to the rest. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DLCTI-E01` | REF[DLCTI-B09]: a pattern that does not compile is the failure B09 absorbs: it compiles to nothing, and the gate answers as for an undeclared pattern | — | — <!-- @resilient: loading the configuration already refuses a pattern that does not compile and names it, so the nil here only reaches a hand-built dialect --> |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
