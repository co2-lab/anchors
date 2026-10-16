<!-- @anchors
  code: MGSYN
  updated_at: 2026-10-16
  layer: migra
-->
# Declared syntax — what the core knew of a language, written into the project at format 8

> **Code**: `MGSYN`

## Overview

Before format 8 the core read some syntax by itself:
- comments, with a table of its own;
- the lexicon of a declared family, filling every field the project left empty;
- the shapes of a route, a mock, a statement's end, a limit and a test's block.

Format 8 reads only what the project declares (DESIGN-core-without-syntax.md). Crossing it
writes the examples the old binary applied into the project's anchors.yaml, once, and only
where the project declared nothing. The project then measures what it measured before, and
its yaml says what is read. What the old binary applied and no example can say is named, to
be declared by hand.

This unit holds two pure parts:
- what to add, from the configuration and the files of its layers;
- writing it into the text of the file, keeping the rest.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | a loaded anchors.yaml | none: nil adds nothing | the config's load |
| the files | the paths of the project's layers, relative to the root | — | the scan |
| the text | the anchors.yaml the configuration was loaded from | a section written in flow style that cannot take a line | `WithSyntax`, which refuses it naming the section |

## Effects

| Effect | Description |
| --- | --- |
| `MGSYN-B01` | With no `dialect.comments`, the entries of the examples whose glob covers one of the files are added, in the examples' order, Anchors' own formats left out; an extension no example covers is named, to be declared by hand. A project that declares comments gets none. (`DeclaredSyntax`) |
| `MGSYN-B02` | Every lexicon field the project left empty is filled from its family's example (`config.FamilyLexicon`); with no family, only the two naming conventions, which the old binary applied to any project, and `env-declared` with no `env_read` is named. (`DeclaredSyntax`) |
| `MGSYN-B03` | A project with no `dialect.tests` gets the family's whole; inside its own, only what the old binary filled: the family's way to open a test and its assertion when it declared neither a pattern nor a script, and the shape of a block. A pattern with no `title` group is named. (`DeclaredSyntax`) |
| `MGSYN-B04` | A project declaring `route-exists` with no `route_pattern` gets the family's, else the ts example the old binary fell back to; one declaring `mock-typed` gets the mock shapes it left empty, the ts example's when its family has none; a test layer of YAML files not yet `headless` becomes so. (`DeclaredSyntax`) |
| `MGSYN-B05` | The additions are written at the end of their block, indented as its children are; a list of short texts on one line; a one-line flow mapping takes them inline. The rest of the text — comments, order, quoting — is kept, and the result reads back with nothing left to add. The fields added are named a line each. (`WithSyntax`, `Fields`) |
| `MGSYN-B06` | A section the file lacks is written whole: a top-level one at the end, after a blank line. One written in flow style that cannot take the fields is refused, naming it. (`WithSyntax`) |
| `MGSYN-B07` | The example runners (`config.ExampleRunners`) whose command line some suite's `run` is — tests or mutation — are added under `monitor.runners`, after the project's own, one it declares by name left out: the old monitor knew them by itself. (`DeclaredSyntax`, `WithSyntax`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MGSYN-E01` | REF[MGSYN-B06]: a section in flow style cannot take the fields | an error naming the section and the fields | rewriting the user's style would change lines nobody asked to change |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the unit reads the configuration package's examples and the standard library.
