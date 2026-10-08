<!-- @anchors
  code: NWTMN
  updated_at: 2026-10-08
  layer: comando
-->
# NewTemplates — the catalog of artifact skeletons: which kinds exist, their headers, their sections and the presets that pick them

> **Code**: `NWTMN`

## Overview

`anchors new` emits skeletons; this unit is what it emits from. It holds the catalog of
artifact kinds, each with its file extension, whether its header owns an identity (`code`) or
references a spec's (`ref`), its header in the comment dialect of its file, and its sections.

Specs, plans, flows and actions are Markdown and their header is an HTML comment. A feature's
header opens with the Gherkin language line, and a test's header uses the line comment of the
file being created, deduced from its extension, so a Python test is never born with a
JavaScript comment. Product doctrine declares no layer, because a doctrine has no target, and
its sections have their own keys so it never inherits a spec's body.

Every section that catalogues rules declares the rule letter it realizes, and those letters
must be the engine's canonical ones; a section realizing another letter would make `new` emit
a spec whose rules the engine cannot see. The test body is written in the idiom of the
project's language family, and it always carries the first scenario code in the test's name,
which is what binds a test result to the spec requirement; for a family it does not know, it
writes an instruction instead of guessing a syntax. The spec presets name common unit kinds
(backend logic, screen, component, store and so on) and fix both the set of sections and the
reading order.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the kind | spec, feature, test, plan, product, flow, action | any other kind | the new command refuses it |
| the language family | python, go, java, kotlin, csharp, rust, ruby, php, ts or empty | any other family | this unit: an instruction instead of a body |
| the output path | any path; its extension picks the test comment dialect | — | the config's comment table |

## Effects

| Effect | Description |
| --- | --- |
| `NWTMN-B01` | The catalog holds seven kinds — action, feature, flow, plan, product, spec, test — and `new` gives birth to each at the given path with the given identity. |
| `NWTMN-B02` | A Markdown header is an HTML comment holding the identity, and `updated_at` and `layer` placeholders to fill. |
| `NWTMN-B03` | A feature header opens with `# language: <lang>` (English when none is given), then the `ref`, and declares the feature layer. |
| `NWTMN-B04` | A test header uses the line comment of the output file's extension, holds the `ref`, and declares the test layer. |
| `NWTMN-B05` | Every section that realizes a rule letter realizes only letters among the engine's canonical rule letters. |
| `NWTMN-B06` | The test body is written in the idiom of each known family (a Python test function, a Go test function, a JUnit test, an xUnit fact, a Rust test, an RSpec block, a PHPUnit method, or a describe/it block), with no other language's syntax. |
| `NWTMN-B07` | Every test body, of every family including an unknown one, carries the first scenario code of the unit. |
| `NWTMN-B08` | For a family it does not know, the test body is an instruction to write a case with the scenario code in its name, not a guessed syntax. |
| `NWTMN-B09` | The twelve spec presets (api, backend-logic, component, handler, hook, mobile-logic, repository, schema, screen, service, store, validation) are each an ordered set of catalog sections, opening with the title and closing with the open decisions. |
| `NWTMN-B10` | Product doctrine has sections of its own, emits `-R` rules, and its header declares no layer. |
| `NWTMN-B11` | The Python and Rust test functions are named by the snake_case of the unit name: each run of separators (`-`, space, `.`, `_`) becomes one `_`, an uppercase letter starts a new word only after a lowercase letter or a digit or when it opens a word after an acronym, and there is no leading or trailing `_` (`My-Name` → `my_name`, `HTTPServer` → `http_server`). |
| `NWTMN-B12` | The spec catalog ties rules to what they use: `validations` realizes `V`, `presentation-validations` realizes `P`, and `rule-uses` realizes no letter; each emits a table whose first column is the rule's code and whose second is what it uses; `rule-uses` is in every preset, `validations` in the screen and validation presets, `presentation-validations` in the screen and component ones. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `NWTMN-I01` | REF[NWTMN-B05]: the letters the catalog emits and the letters the engine reads are the same set | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `NWTMN-X01` | REF[NWTMN-B08]: the catalog never guesses the syntax of a language it does not know | — |

## Errors

none — the catalog is data and pure renderers; an unknown kind is refused by the new command, and an unknown family is a documented instruction (`NWTMN-B08`), not a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
