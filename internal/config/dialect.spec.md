<!-- @anchors
  code: DLCTI
  updated_at: 2026-09-26
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
| `DLCTI-B01` | A project that declares no dialect gets only the naming defaults: no language lexicon at all. |
| `DLCTI-B02` | The family fills every lexicon field the project left empty, and a field the project declared wins over the family (`DialectFor`). |
| `DLCTI-B03` | The family name is matched ignoring case. |
| `DLCTI-B04` | An unknown family contributes nothing; the known families are listed in alphabetical order for the messages that name them (`KnownDialectFamilies`). |
| `DLCTI-B05` | The two naming conventions (the set-promise and the set-slice verbs) apply whatever the family, and a declared one replaces its default. |
| `DLCTI-B06` | The Gherkin language defaults to English and is looked up ignoring case, answering the table's own spelling of the code, so `PT` is `pt` and `zh-cn` is `zh-CN`; a language outside the table keeps its code as declared and gets the English keywords (`GherkinFor`). |
| `DLCTI-B07` | The recognised ways to open a scenario are those of every language in the table — scenario, outline and synonyms — each once, longest first, so an outline that contains the word for scenario is tried before it (`GherkinScenarioAlternatives`). |
| `DLCTI-B08` | The recognised result keywords are those of every language in the table, in alphabetical order (`GherkinThenAlternatives`). |
| `DLCTI-B09` | An empty or invalid pattern compiles to nothing, which the gates read as "not declared" (`Compile`). |
| `DLCTI-B10` | A lexicon field is waived only when the opt-out list names it by its configuration name, ignoring case and spaces; its name in the code does not waive it (`WaivedField`). |
| `DLCTI-B11` | The default set-promise convention recognises the verb after a provider prefix (a cloud-provider list call) and never inside another word (allocate, callback, enlistment). |
| `DLCTI-B12` | The default set-slice convention is a query verb opening the name, followed by a slice word (first, recent, top, page…). |

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
| `DLCTI-E01` | REF[DLCTI-B09]: a pattern that does not compile is the failure B09 absorbs: it compiles to nothing, and the gate answers as for an undeclared pattern | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | config — holds the declared dialect block |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
