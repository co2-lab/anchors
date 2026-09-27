<!-- @anchors
  code: CDGNC
  updated_at: 2026-09-26
  layer: infra
-->
# CodeGenerator — the short, stable identity code suggested for a unit's name

> **Code**: `CDGNC`

## Overview

Every specifiable unit carries a short, unique and stable code that prefixes its rules and crosses
spec, feature and test. This unit suggests that code from the unit's name: it compresses the name into
the project's code length, and resolves a collision against the codes already taken. The suggestion is
not aesthetics: what matters is uniqueness in the project's namespace, and a person may still pick a
more readable code by hand.

The compression follows the shape of the name. Generic suffixes such as "Screen" carry no identity and
are dropped first. A name of many words keeps the initials; a name of a few words takes letters from
each word; a single word takes its consonants, then its vowels. Short results are padded with X. A
project that groups codes by module may ask for a module prefix, which the code then starts with.

The generated length follows the project: it is the smallest length the project declares, so during a
migration from one length to another the new codes stay readable by every tool of the project. The same
padding rule completes an existing code of the wrong length, which preserves the choice of whoever named
it instead of generating a different code.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the unit name | a name made of ASCII letters, digits and the separators `-`, `_` and space | — | this unit: any name yields a code, an empty one yields only padding |
| the codes taken | the codes already in use in the project's map | — | the caller (`anchors new`, `anchors code`) |
| the declared code lengths | the project's `code_lengths`, each between 2 and 8 | lengths out of range | the configuration loader, which refuses them before they arrive |

## Effects

| Effect | Description |
| --- | --- |
| `CDGNC-B01` | A generic suffix (Screen, Component, Modal, Sheet, Layout, Overlay, View, Page) is removed before compressing; a name that is only the suffix is kept whole. (`StripGeneric`) |
| `CDGNC-B02` | The name is split into words at `-`, `_` and spaces, where a lowercase letter or a digit meets an uppercase letter, and between an acronym and the capitalised word after it. |
| `CDGNC-B03` | A name with at least as many words as the code length takes the initials of its first words. (`Generate`) |
| `CDGNC-B04` | A name of two words up to one fewer than the code length gives each word the same share of letters (its initial, then its consonants), and any slot left is filled from the joined words, consonants first. |
| `CDGNC-B05` | A single word gives its consonants in order, then its vowels and digits, so a word that starts with a vowel does not start its code with it. |
| `CDGNC-B06` | The code is upper case, padded with X up to the code length; padding an existing code applies the same rule and cuts it to the code length. (`Pad`) |
| `CDGNC-B07` | With a module prefix, the code starts with the prefix in upper case (cut to the code length) and the rest comes from the name's initial and consonants. (`GenerateWithPrefix`) |
| `CDGNC-B08` | A code not taken is returned as generated; a taken one is varied deterministically, trying every other letter in the last position, then in the one before, never changing the positions of the module prefix. (`GenerateUnique`, `GenerateUniqueWithPrefix`) |
| `CDGNC-B09` | The module prefix of a module name is its initial and its first consonant, padded with X to two characters. (`ModulePrefix`) |
| `CDGNC-B10` | The generated length is the smallest declared code length of at least 2; loading the project's configuration applies it, and an empty list leaves the length unchanged. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CDGNC-I01` | The same name and the same taken codes always give the same code, and a resolved collision is never a taken code. | generates twice for the same inputs and compares, and checks the resolved code is not in the taken set |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CDGNC-X01` | Does not guarantee a free code when every variation is taken: it then returns the generated code, taken as it is. | The search is bounded to single-position variations; a namespace that full is a project decision (a longer code length), not something the generator can invent. |

## Errors

none — every name yields a code (an empty name yields only padding), and a saturated namespace is the boundary X01, not a failure.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `SetSlotsHook` | config — the project's declared code lengths reach the generator when the configuration loads |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
