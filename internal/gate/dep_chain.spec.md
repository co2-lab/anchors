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

Anchors does not parse the language: it reads the import's path with a pattern — the
dialect's `import_pattern` with the path in the group `path`, or the family's (`ts`, `go`) —
and resolves it with `import_resolve` (relative paths, the project's aliases, the extensions).
An import that resolves to no file of the map — a package of the ecosystem — is no link of
the chain. A test takes no part: it is tied to its unit by its `ref:`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted node | code of the project | a test, a support file, a vendored file | this unit, which skips them |
| the import pattern | the dialect's `import_pattern` with a `path` group, or a family that knows one | a dialect with neither | this unit: pending, "nothing to read" |
| the import's path | relative to the importer, or under an alias the project declares | a package outside the project | this unit: resolves to nothing, no link |

## Effects

| Effect | Description |
| --- | --- |
| `DCGDP-B01` | The import statements of a file are read by the dialect's pattern — a statement whose names span lines is one, its path on its last line, an alias read as the name it aliases — and each path resolved to a file of the map: relative to the importer or through an alias, as written or with each extension; a path naming a directory of code files is a package; anything else is outside the project. With no pattern for the dialect, nothing is read. (`ImportsOf`, `RealImport`) |
| `DCGDP-B02` | `dep-declared` fails naming each import of a governed file with no dependency flag on its statement, by line, path and the code it would carry; an import outside the project, or one with a waiver, is not charged, and a vendored file is skipped. Each failure of the chain points at `anchors guide header`, where the flags are. (`checkDepDeclared`) |
| `DCGDP-B03` | `dep-honored` fails naming each dependency flag whose code is not the code of the file its import resolves to — with the right one —, and each flag on a line with no import. (`checkDepHonored`) |
| `DCGDP-B04` | `used-by-declared` fails naming each symbol another code file imports whose used-by flag is missing or names other codes, and each used-by flag on a symbol nobody imports. (`checkUsedByDeclared`, `UsedByOf`) |
| `DCGDP-B05` | The fixers write the dependency flag on each unflagged import of a single governed file, correct a flag naming another code, write above each imported symbol the used-by flag with exactly who imports it, and remove a used-by flag nobody answers; an import of a package of several files is left for the author. (`fixDepFlags`, `fixUsedBy`) |
| `DCGDP-B06` | A re-export declares the names it lists — `export { X } from`, `export type { A, B }`, or a list spanning lines —, and the fixer writes each name's used-by flag above the list, naming its symbol in parentheses; an inline import (`import('…').Name`) brings the member it reads and carries its dependency flag on its line. (`declarationLine`) |
| `DCGDP-B07` | Every artifact that imports takes part in the chain — code, tests, test support, flows —: a test's imports carry the dependency flag like the code's, a flow composing another (`runFlow:`, inline or by `file:`, relative to the flow) imports it, the fixer writes the flag in each artifact's own comment, and — when the used-by gate confronts tests (`on:` holds `test`) — a symbol a test imports lists the test in its used-by flag; with the gate on code alone, the used-by flags list the code that imports the symbol, as before. (`chainUnit`) |
| `DCGDP-B08` | A kinded dependency flag (`@dep[<kind>]: <name>`) is no import: `dep-honored` asks only that its kind is one the project declares in `dependency_kinds:`, naming the line, the kind and the name when it is not. |
| `DCGDP-B09` | An import of types only — the dialect's `type_import_pattern`, or its family's (TypeScript: `import type …`, `export type … from`) — is flagged `@dep[type]: CODE` by the fixer, a plain flag on it is turned into one, and a type flag on a static import that runs is turned back — a type flag on an inline `import('…')`, which reads like a dynamic import, is the author's and stays; `dep-honored` confronts a type flag as it does an import's. (`fixDepFlags`, `familyTypeImportPattern`) |

## Errors

none — an unreadable file is skipped by the readers it belongs to, and a pattern that does not compile reads no import; neither is a failure of this unit.
