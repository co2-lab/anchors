<!-- @anchors
  code: CHPTC
  layer: desenho
-->

# Chain patterns — imports and declarations are read by the project's patterns; the core reads no syntax

> APPROVED on 2026-10-16 (the user: "faça as recomendações visando a melhor solução sem se
> preocupar com compatibilidade"). Proposed by MIF after 0.1.348: "it's the same principle as
> the name group, propagates and test_handles".
>
> IMPLEMENTED on 2026-10-16. Two departures from the proposal, both from "no compatibility":
> - `import_pattern` (one regular expression) became `imports`, a list of patterns: Go's
>   blocks and TS's several forms need a list.
> - A pattern may carry `runs`. The fixer turned a type flag on a static import back to a plain
>   one by knowing what a static import looks like; now the pattern that reads such statements
>   says they run.
>
> Measured on MIF's HEAD (bc0fba5aa), in a clone with the `ts` example declared as the doctor
> offers it:
> - `dep-declared`, `dep-honored`, `used-by-declared` and the three handle gates gave the same
>   table as 0.1.348;
> - `check --all --fix` changed no file.
>
> The first measurement found five `export default function X` read as a list of two (the
> default and the name); two patterns reading one statement are now no list.

## The problem

`used-by-declared`, `dep-declared` and their fixers read a file's chain with syntax the core
carries. Measured in the code on 2026-10-16:

| What | Where | Whose syntax |
| --- | --- | --- |
| The names an import brings: `{ A, B as C }`, a default import, `import type`, a destructuring of `await import()`, `import('./x').Name` | `scan.importSymbols`: `importBracesRE`, `importDefaultRE`, `importDestructRE`, `importInlineRE` | TS/JS |
| What the line below a `@used-by` declares: `function`, `const`, `class`…, an export list, the default export | `scan.extractUsedBy`: `declaredNameRE`, `exportListRE`, `defaultExportRE` | TS/JS, with `func` and `def` added |
| The same, for the fixer | `gate.declarationLine`: its own regex, `exportListLineRE`, `defaultDeclRE` (copies) | TS/JS |
| A statement over several lines | `gate.ImportsOf`: `listOpenRE` (`import {` … `} from`), and Go's `import ( … )` keyed on `d.Family == "go"` | TS, Go |
| The import's path when the project declares none | `gate.familyImportPattern`, `familyTypeImportPattern`: a table of families inside the gate, applied by itself | ts, go, ocaml… |
| A flow's import | `flowImportRE` (`runFlow:`, `file:`) and `isFlowFile` (`.yaml`) | Maestro |

0.1.348 made the identifiers Unicode-aware, which helped one bug, but the grammar is still in
the core. A Go, Python or Rust project cannot be measured by these gates unless the core learns
its grammar: exported `func Name(`, `def name(`, `pub fn`, `from x import y`, `use crate::x`.

The dialect already holds part of this:
- `import_pattern`: the path, in a `path` group;
- `type_import_pattern`;
- `import_resolve`;
- `definition`: the names a unit defines, which `test-exercises-unit` reads;
- `exported_func`.

## The design

### 1. One engine: patterns over the file, lists by `within` / `each`

Every chain pattern runs over the whole content, not line by line. Each capture carries its
offset, which becomes its line. A pattern may then span lines, so `import {` … `} from` and
Go's `import ( … )` need no code of their own.

A list of names cannot come from one group: a regex captures once per match. The same limit
left a nested conditional out of the `ts` test handle example. So a pattern entry is either:
- a string, which captures once per match; or
- a pair `{within, each}`. `within` captures a region in its `within` group, and `each` runs
  repeatedly inside that region.

      import_names:
        - within: 'import\s+(?:type\s+)?\{(?P<within>[^}]*)\}\s*from'
          each: '(?:^|,)\s*(?:type\s+)?(?P<name>[\p{L}_$][\p{L}\p{N}_$]*)'
        - 'import\s+(?P<default>[\p{L}_$][\p{L}\p{N}_$]*)\s*(?:,|\s+from)'

The engine is generic. It also serves `test_handles`: a conditional's branches become
`{within: <the expression>, each: <a literal>}` (D05).

### 2. The fields

| Field | Groups | Read by |
| --- | --- | --- |
| `imports` (was `import_pattern`) | `path` | dep-declared, the resolution |
| `type_import_pattern` (exists) | — | the fixer's `@dep[type]:` |
| `import_names` (new) | `name`, or `default` (the module's default) | used-by-declared: what a file imports of another |
| `declarations` (new) | `name`, or `default` | used-by-declared and its fixer: what the line below a flag declares |

The core keeps only the rules that are Anchors' own:
- a declaration that yields several names needs the name on the flag (`@used-by: X (Name)`);
- the name inside the flag's parentheses is whatever stands there, up to `)`, so the flag needs
  no identifier syntax.

### 3. Flows are patterns too

A flow's import (`runFlow:`, `file:`) is one more `imports` entry, narrowed by `files` to the
flows and marked `relative`: its paths are relative to the flow without `./`. Which files are
flows is the pattern's glob, not a `.yaml` test in the core.

### 4. The families' syntax is an example

`FamilyChain(family)` holds today's behavior for `ts` and `ocaml`, and the `go` and `python`
examples written from their grammars. Rust is left out: no project is measured with it yet.
`anchors init` writes the example, and the governance tip offers it, as with `test_handles`
(D02).

### 5. Measured before it ships

MIF's HEAD runs in a clone with the `ts` example declared. `used-by-declared` and
`dep-declared` must give the same findings as 0.1.348, and any difference is explained before
the release.

## Phases

| Phase | What | Proof |
| --- | --- | --- |
| `CHPTC-W01` | **The engine.** Patterns over the content with lines, `{within, each}`, validation (compiles, has a group of its list). `test_handles` accepts the pair; the `ts` example reads nested conditionals. | ImportScreen's nested conditional reads all its branches; a pattern spanning lines gives the line of each capture. |
| `CHPTC-W02` | **The fields.** `import_names` and `declarations`; the eight core regexes and `familyImportPattern` removed; flows by the surface. | The `ts` example gives this repository's and MIF's findings unchanged; a Go and a Python fixture are measured with their examples. |
| `CHPTC-W03` | **The way in.** `init` writes the family's example, the doctor offers it, and the message tells MIF what to declare. | A project with neither field reads `pending` with the field named, and the tip's block loads back. |

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `CHPTC-D01` | New fields beside `import_pattern`, or one `chain:` block holding all four? | **Beside**, under `dialect:` (2026-10-16). With compatibility no concern, `import_pattern` became `imports`, a list. |
| `CHPTC-D02` | Is the family's example applied by itself when the project declares nothing? Today `familyImportPattern` is, so dep-declared measures every ts and go project out of the box. | **Example only** (2026-10-16), as `test_handles` (the user, 2026-10-15: "se não for configurado no dialect o gate não funciona"). It breaks every project that relies on the family until it declares the patterns; `init` and the doctor give the block. |
| `CHPTC-D03` | Merge `declarations` with `definition` (test-exercises-unit)? | **No.** `definition` is every name a unit defines; `declarations` is what the line below a flag exports, the default and export lists included. The questions differ, and a pattern made for both answers neither well. |
| `CHPTC-D04` | Multi-line statements: patterns over the whole content, or line by line with a list continuation in the core? | **Over the whole content**: the continuation is syntax (`{`, `(`), and the pattern can say it. |
| `CHPTC-D05` | Does `test_handles` take `{within, each}` too? | **Yes**, and the `ts` example uses it for conditionals: one engine, and MIF's nested conditional needs no project pattern. |

## Limits, measured on 2026-10-16

- A flow is read by a pattern narrowed with `files` to the flows and marked `relative`. The ts
  example narrows it to `**/*.{yaml,yml}`; a project narrows it to its own e2e surface.
- Python's dotted paths (`app.core`) resolve to no file: `import_resolve` maps prefixes, not
  dots to directories. The example reads the statements; resolving them is future work.
- Still syntax in the core, out of this design's scope:
  - comment stripping (the comment markers, OCaml's `(* *)` in `by_name`);
  - `navigation_call`'s family fallback;
  - `code-cataloged`'s default TS export regex when `exported_func` is empty.

## What NOT to do

- **Keep the family table in the gate as a fallback.** That is the syntax in the core under
  another name.
- **Learn each language's grammar in Go code.** A language Anchors has not seen must be
  measurable by declaring patterns, with no release.
