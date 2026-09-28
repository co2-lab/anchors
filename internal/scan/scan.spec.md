<!-- @anchors
  code: RPSCR
  updated_at: 2026-09-28
  layer: scan
-->
# RepoScan — the repository read as text: which files exist, of which layer, and what each declares

> **Code**: `RPSCR`

## Overview

The scan is the map's entry door. It walks the repository reading TEXT — never parsing code — and
extracts the facts the map is built from: which files exist, which declared layer each belongs to,
which rule codes each carries, and what each declares in its header and body (identity, unit layer,
dependencies, the plans it waits for or revises, its parent, the doctrine rules it realizes, the flag
scenarios that gate it, the specs a plan seeds). Every gate stands on these facts.

The failure mode the unit is built against is SILENCE. A file the scan misses is not reported as
missing: it never exists for any gate. A wrong vocabulary, a native path separator, an unanchored
ignore pattern or a verbose layer pattern all produced a smaller map without a single error line, and
each rule below closes one of those holes. Where the scan must choose — two layers matching one file —
the default is a heuristic, the project's declaration wins over it, and the choice is reported when the
heuristic decided.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a readable directory | a path that does not exist or cannot be walked | this unit: the walk error is returned to the caller |
| the configuration | the loaded project configuration with its layers, their patterns, priorities, exclusions and kinds, and the rule letters | an absent configuration | the caller loads it before walking |
| a path to classify | a root-relative path, with `/` or `\` separators | an absolute path | the callers pass root-relative paths; the separator is normalized here |
| the file content | any text | — | this unit: content without codes or declarations yields empty facts |
| a spec's dependency table | a markdown table under the dependencies heading, in any catalogue language, with a code column and a file column | a table without a code column, which is documental | this unit: such a table yields no dependency |

## Effects

### Walk — which files enter the map

| Effect | Description |
| --- | --- |
| `RPSCR-B01` | `Walk`: Only files that match a declared layer are returned, each with the layer and the kind of that layer; a file outside every layer is left out. |
| `RPSCR-B02` | Directories and files the ignore set excludes (the project's `.gitignore`, the built-in list, editor ephemera) are not returned. |
| `RPSCR-B03` | A directory that is another checkout — it holds its own `.git`, a file in a worktree or a directory in a clone — is not descended into. |
| `RPSCR-B04` | A plan's progress companion is not returned, even when a layer matches it. |
| `RPSCR-B05` | A workflow owned upstream enters with no rule codes of its own: the codes in its comments are the Anchors project's examples. |
| `RPSCR-B06` | `ShortHash`: A file's revision is a short hash of its content with CRLF line endings normalized to LF, so the same content has the same revision on every machine, and different content a different one. |

### Classification — which layer a file belongs to

| Effect | Description |
| --- | --- |
| `RPSCR-B07` | `Classify`, `ClassifyPath`: When several layers match, the highest declared priority wins; at equal priority the longer pattern wins; at equal length the lexically smaller layer name wins. |
| `RPSCR-B08` | A layer's exclusion globs remove a path from that layer; the path may still belong to another layer. |
| `RPSCR-B09` | A path written with Windows separators is classified the same as its slash form. |
| `RPSCR-B10` | `Ambiguities`: A file matched by two or more layers whose winner was not decided by a higher declared priority is reported as an ambiguity naming the winner and the losers; a winner by declared priority is not reported. |
| `RPSCR-B11` | `LayerOfUnit`: The layer of a file's UNIT is the `layer:` its header declares; without one it is the file's classified layer. |

### Codes — what a file owns

| Effect | Description |
| --- | --- |
| `RPSCR-B12` | The rule codes a file owns are those outside `//` line comments and `/* */` block comments, each listed once in order of first appearance. |
| `RPSCR-B13` | `ScenarioCodeRE`, `SetRuleLetters`: Rule codes are recognised with the rule letters the project declares, not a fixed set. |
| `RPSCR-B14` | The identity declared on the header's `code:` line is recorded apart from the codes cited in the body, and the `@anchors-shared-code` and `@noPropagation` annotations are recorded as flags; the same reading is offered to other packages (`HeaderCodeOf`). |

### Header declarations — relations a file declares

| Effect | Description |
| --- | --- |
| `RPSCR-B15` | A file's `parent:` is read only inside its `@anchors` header; a body line starting with `parent:` is not a declaration. |
| `RPSCR-B16` | A plan's `needs:` lists plan paths, each resolved against the root; a spec's `needs:` keeps only phase codes (`CODE-Fnn`) and drops anything else; other kinds need nothing. |
| `RPSCR-B17` | Only a plan declares `revises:`, whose paths are resolved like `needs:`; any other kind revises nothing. |
| `RPSCR-B18` | A file that is not a spec declares its dependencies on the header's `dep:` line, a comma-separated list of paths resolved against the root. |
| `RPSCR-B19` | A YAML test script's dependencies are the scripts it composes, each recorded with the method `runFlow`. |

### Dependency table — what a spec consumes

| Effect | Description |
| --- | --- |
| `RPSCR-B20` | A spec's dependencies are the rows of the first table under the dependencies heading, recognised in every catalogue language; the columns are found by the meaning of their header in any order, and a spec without the heading has none. |
| `RPSCR-B21` | A row without a file, or whose code is not `DEP` followed by digits, is not a dependency. |
| `RPSCR-B22` | The method cell keeps its backticks, the sign that it names a symbol; the other cells lose them. |
| `RPSCR-B23` | A declared file resolves to the first existing of: the path as written, the path under the spec's `src/` directory, and the path next to that `src/`; when none exists, the path is kept as written. |

### Rule tags — `@realizes` and `@gated-by`

| Effect | Description |
| --- | --- |
| `RPSCR-B24` | An `@realizes` tag is paired with the rule code of the same line or of the nearest rule line above it, in all three rule forms: heading, table row and bold bullet. |
| `RPSCR-B25` | A blank line closes a rule's scope: a tag after it has no owning rule. |
| `RPSCR-B26` | Only a spec declares rule tags; in any other kind they are not read. |
| `RPSCR-B27` | A repeated (rule, target) pair is recorded once; one rule realizing several targets, and several rules realizing one target, are all kept. |
| `RPSCR-B28` | An `@gated-by` tag is read only when it names a flag scenario code (letter `G`); a tag naming any other letter is not a declaration. |

### Seeds — what a plan promises

| Effect | Description |
| --- | --- |
| `RPSCR-B29` | A plan seeds the backticked `.spec.md` and `.doctrine.md` paths it cites that contain a directory, carry no glob character and are not templates, each once; no other kind seeds. |

### Header keys — where a declaration is read

| Effect | Description |
| --- | --- |
| `RPSCR-B30` | The header keys `code:`, `layer:`, `needs:`, `revises:` and `dep:` are read only inside the file's `@anchors` header, like `parent:`; a body line starting with one of them is not a declaration. |
| `RPSCR-B31` | The rule tags `@realizes` and `@gated-by`, and the rule a tag belongs to, are read with the code length the project declares (`code_lengths`): a fixed length left the edges of every other length undrawn. |
| `RPSCR-B32` | A file that matches its own layer's `support` list is marked as support; a file of the layer outside the list, or matching only another layer's list, is not. |
| `RPSCR-B33` | Given paths, the scan reads only them, each exactly as the walk reads it in the tree, and leaves out one that is ignored, in an ignored directory, the progress file, in no layer, or missing; the governed paths of the tree can be listed without reading any file. (`ScanPaths`, `GovernedPaths`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RPSCR-I01` | The classification of a path is the same on every run with the same configuration. | classifies the same path fifty times under a tie of patterns and verifies a single answer |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RPSCR-X01` | The dependency table is read only from specs: a guide's table under a dependencies heading is not a dependency. | A document can list related files under that heading; only a spec's table is a reuse contract. |
| `RPSCR-X02` | A spec's dependencies come from its table only; a `dep:` line in a spec is not read. | The spec has the table as its contract; two sources for the same relation would disagree. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RPSCR-E01` | The root cannot be walked (it does not exist or cannot be read). | The walk error is returned to the caller. | A map built from an unreadable root would be empty without saying so. |
| `RPSCR-E02` | A file that matches a layer cannot be read (other than having vanished since the listing). | The walk fails with an error naming the file. | Dropping it would shrink the map without a line of error; a file that vanished has nothing to map and is skipped. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | config — declared layers, rule letters and code length |
| DEP2 | `internal/i18n/i18n.go` | `AllTranslations` | apoio — the dependencies heading in every catalogue language |
| DEP3 | `internal/scan/ignore.go` | `LoadIgnoreFor` | scan — what the walk never sees |
| DEP4 | `internal/scan/progress.go` | `IsProgressFile` | scan — the progress companion kept out of the map |
| DEP5 | `internal/scan/upstream.go` | `IsUpstreamOwned`, `AnchorsHeader` | scan — upstream ownership and the header block |
| DEP6 | `internal/scan/region.go` | `ComposeRefs` | scan — the composition of test scripts |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
