<!-- @anchors
  code: DCGDP
  layer: gate
-->

# DependencyChain — every import flagged with the code it uses, and every symbol with who uses it

> **Code**: `DCGDP`

## Overview

Every import is a dependency, and the map knew only the ones somebody wrote in a spec's
Dependencies table: a project imported its tokens, a toast hook and its components without
declaring them, and every gate stayed green. This unit closes the chain in both directions,
with flags the agent writes beside the code and Anchors confronts with the code
(DESIGN-dependencies-and-navigation.md):

- `dep-declared` — every import of a governed file carries the dependency flag naming the
  code of the file it uses, or a waiver with its reason;
- `dep-honored` — every dependency flag names the code of the file its import resolves to;
- `used-by-declared` — every symbol another file imports carries the used-by flag with exactly
  the codes of who imports it.

Anchors does not parse the language: it reads each import statement with the project's
patterns — the dialect's `imports`, the path in the group `path`, and `import_names` for what
it brings — and resolves the path with `import_resolve` (relative paths, the project's aliases,
the extensions). The core knows no syntax of an import (DESIGN-chain-patterns.md).
An import that resolves to no file of the map — a package of the ecosystem — is no link of
the chain. A test takes no part: it is tied to its unit by its `ref:`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | code of the project | a test, a support file, a vendored file | this unit, which skips them |
| the import patterns | the dialect's `imports`, each with a `path` group | a dialect with none: no family gives them by itself | this unit: pending, naming the field |
| the import's path | relative to the importer, or under an alias the project declares | a package outside the project | this unit: resolves to nothing, no link |

## Effects

| Effect | Description |
| --- | --- |
| `DCGDP-B01` | The import statements of a file are the matches of the dialect's `imports` that read its files — a match may span lines, its first line the statement's and the line of its `path` the one its flag goes on; what it brings are the captures of `import_names` on its text — and each path is resolved to a file of the map: relative to the importer or through an alias, as written or with each extension of `import_resolve`; a path naming a directory of code files is a package; anything else is outside the project. (`ImportsOf`, `resolveImport`) |
| `DCGDP-B02` | `dep-declared` fails naming each import of a governed file with no dependency flag on its statement, by line, path and the code it would carry; an import outside the project, or one with a waiver, is not charged, and a vendored file is skipped. Each failure of the chain points at `anchors guide header`, where the flags are. (`checkDepDeclared`) |
| `DCGDP-B03` | `dep-honored` fails naming each dependency flag whose code is not the code of the file its import resolves to — with the right one —, and each flag on a line with no import. (`checkDepHonored`) |
| `DCGDP-B04` | `used-by-declared` fails naming each symbol another code file imports whose used-by flag is missing or names other codes, and each used-by flag on a symbol nobody imports. (`checkUsedByDeclared`, `UsedByOf`) |
| `DCGDP-B05` | The fixers write the dependency flag on each unflagged import of a single governed file, correct a flag naming another code, write above each imported symbol the used-by flag with exactly who imports it, and remove a used-by flag nobody answers; an import of a package of several files is left for the author. (`fixDepFlags`, `fixUsedBy`) |
| `DCGDP-B06` | A re-export declares the names it lists — with the `ts` example, `export { X } from`, `export type { A, B }`, or a list spanning lines —, and the fixer writes each name's used-by flag above the list, naming its symbol in parentheses; an inline import (`import('…').Name`) brings the member it reads and carries its dependency flag on its line. (`declarationLine`) |
| `DCGDP-B07` | Every artifact that imports takes part in the chain — code, tests, test support, flows —: a test's imports carry the dependency flag like the code's; a flow composing another imports it when a pattern of `imports` narrowed to the flows reads it, its path relative to the flow when the pattern says `relative` (the `ts` example: `runFlow:`, inline or by `file:`); the fixer writes the flag in each artifact's own comment, and — when the used-by gate confronts tests (`on:` holds `test`) — a symbol a test imports lists the test. |
| `DCGDP-B08` | A kinded dependency flag (`@dep[<kind>]: <name>`) is no import: `dep-honored` asks only that its kind is one the project declares in `dependency_kinds:`, naming the line, the kind and the name when it is not. |
| `DCGDP-B09` | An import of types only — the statement the dialect's `type_import_pattern` reads (the `ts` example: `import type …`, `export type … from`) — is flagged `@dep[type]: CODE` by the fixer and a plain flag on it is turned into one; a type flag on a statement whose pattern says it `runs` is turned back, and on any other — an inline `import('…')`, which reads like a dynamic import — it is the author's and stays; `dep-honored` confronts a type flag as an import. |
| `DCGDP-B10` | Where the dialect resolves imports by name (`import_resolve.by_name`, the `ocaml` example), every module the `imports` patterns read on a line outside its comments is an import, read once per file and resolved to the code file named after it (`name_case`, `extensions`): the one beside the importer, or the only one in the project; a name of several files, none beside the importer, resolves to none. The fixer writes the file's flags one per line below its header, and removes a flag of a module it no longer uses. |
| `DCGDP-B11` | The fixer finds the declaration of a symbol whose name ends in a letter outside ASCII (`café`), and does not take a longer name that begins with it (`caféForte`) for it. (`declarationLine`) |
| `DCGDP-B12` | Without `dialect.imports` the chain gates are Pending, naming the field — no family reads a project's imports by itself —; and an import commented out, on its own line or after code, is no import. (`hasImports`, `blankComments`) |
| `DCGDP-B13` | A region of `imports` reads a list of imports as one each: with the `go` example, every path of an import block is an import on its own line, flagged there, and a slice of strings outside the block is none. |
| `DCGDP-B14` | An import path whose directories `import_resolve.separator` divides (Python's `.`) resolves as if divided by `/`: `app.core.money` through an alias, `.helpers` beside the importer, `..core` one directory up, a package by its extensions (`/__init__.py`). (`withSlashes`) |
| `DCGDP-B15` | The used-by flag the fixer writes names a listed symbol inside the file's comment: in a block-comment language the closer (`*)`, `-->`) follows the symbol, never precedes it. (`usedByFlag`) |
| `DCGDP-B16` | dep-honored names a flag whose kind contradicts its import: `@dep:` on a statement the dialect's `type_import_pattern` reads as types only — to be `@dep[type]:` —, and `@dep[type]:` on one whose import pattern says it plainly runs — to be `@dep:`. What reads the map trusts the kind (`depends_on_layer` counts value edges); the fixer rewrites it, and the gate no longer lets it pass silent. |

## Errors

none — an unreadable file is skipped by the readers it belongs to, and a pattern that does not compile reads no import; neither is a failure of this unit.
