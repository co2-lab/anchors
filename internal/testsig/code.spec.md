<!-- @anchors
  code: RCGRL
  updated_at: 2026-09-26
  layer: infra
-->
# RuleCodeGrammar — the grammar that recognizes a scenario code in a test's name, in the project's vocabulary

> **Code**: `RCGRL`

## Overview

A test proves a rule when its name carries the rule's code, and the run report is the only place Anchors
reads that name. This unit holds the grammar that recognizes such a code in free text: the unit's identity
(a run of uppercase letters and digits of the project's code length), a dash, and then either a rule
(a declared rule letter followed by two digits, optionally with a lowercase slug), a design-system item, or
the visual-regression marker.

The vocabulary belongs to the PROJECT, not to the engine. The rule letters and the code length start at the
canonical defaults and are replaced by the project's configuration when it loads. The defaults must be the
same as the configuration's own defaults: a report read before the configuration is applied uses them, and
a letter missing here makes a passing test look like no test at all. This is one of the three copies of the
rule letters that have to move together.

The unit keeps the package free of any dependency on the scanner or the configuration: the values arrive
through two setters, and an empty value leaves the current vocabulary in place.

## Effects

| Effect | Description |
| --- | --- |
| `RCGRL-B01` | A rule code — the unit's identity, a dash, a declared rule letter and two digits — is recognized in a test's name, wherever it appears in the text. |
| `RCGRL-B02` | A rule code followed by a lowercase slug (`-some-slug`) is recognized with the slug as part of the code. |
| `RCGRL-B03` | A design-system item (`CODE-DS-<name>`) and the visual-regression marker (`CODE-VR`) are recognized as codes. |
| `RCGRL-B04` | `SetRuleLetters`: Declaring the rule letters replaces the vocabulary: a declared letter is recognized, and a letter outside the declaration is not. |
| `RCGRL-B05` | `SetCodeLenPattern`: Declaring the code length replaces the accepted identity length: an identity of a newly declared length is recognized. |
| `RCGRL-B06` | An empty declaration of letters or of length keeps the vocabulary that was in place: a configuration that declares nothing changes nothing. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared rule letters | a string of uppercase rule letters, or empty | lowercase letters or punctuation meant as letters | the configuration, which produces the project's letters |
| the declared code length | a length pattern as the configuration renders it, or empty | a malformed pattern, which would make the grammar unbuildable | the configuration's length pattern, the only producer |
| the text read | any test name, in any language | — | this unit reads any text and only extracts what matches |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RCGRL-I01` | The default rule letters are the configuration's default rule letters, letter for letter. A report read before the configuration loads recognizes the same rules the rest of the engine does. | compares the default letters against the value the configuration declares, and reads a code of every canonical letter |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RCGRL-X01` | An identity longer than the accepted length, or glued to a longer word, is not a code. | A code is a whole token; matching the tail of a longer word would prove rules no test names. |

## Errors

none — the grammar only reads text and never fails: a text without a code yields no code, and a malformed length pattern is outside the domain (the configuration is its only producer).

## Dependencies

none — the unit depends on nothing of the project; the configuration's values arrive through the setters.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
