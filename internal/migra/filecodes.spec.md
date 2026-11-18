<!-- @anchors
  code: MGFCD
  updated_at: 2026-10-04
  layer: migra
-->
# FileCodes — every file a code of its own, of five characters

> **Code**: `MGFCD`

## Overview

A code always named a file: the spec carried its `code:`, and the code, the feature and the test only
`ref:`d it, so a file could be named only through its unit. The dependency chain needs every file to be
addressable, so format 7 gives every governed file a `code:` of its own, generated from its name and its
type, beside the `ref:` it keeps; and it widens the four-character codes a project kept from before the
default became five. This unit holds the pure part of that migration: what a four-character code
becomes, what a file's code is generated from, which file can carry the line, and where it goes.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| an old code | a four-character code | — | the migration, which widens only those |
| the taken codes | every code of the project, new ones included as they are made | — | the migration |
| a file | a scanned file with a path and its layer or kind | — | the scan |
| the content | any text | a binary | `CanCarryCode`, which refuses it |

## Effects

| Effect | Description |
| --- | --- |
| `MGFCD-B01` | `WidenedCode`: A four-character code becomes a five-character one that keeps it as the prefix and adds a letter of the unit's name; a taken code is never returned. |
| `MGFCD-B02` | `FileCodeName`: A file's code is generated from its name without the artifact suffixes (`.spec.md`, `.feature`, `.test`, `_test`) and extension, and its layer — or its kind when it has no layer —, so the files of one unit differ. |
| `MGFCD-B03` | `FileCode`: The file's code has five characters and is not one of the taken codes. |
| `MGFCD-B04` | `CanCarryCode`: A file can carry the line when it is text and its extension has a comment syntax — a Gherkin `.feature` included, which comments with `#` —; a binary, a JSON or an unknown type cannot. |
| `MGFCD-B05` | `WithHeaderCode`: The `code:` line goes right below the `@anchors` opener of an existing header, in the file's comment syntax; a file with no header gets one at the top, after a shebang when there is one. |
| `MGFCD-B06` | `AsRefWithOwnCode`: A header whose `code:` is its unit's has that line turned into `code: <own>` followed by `ref: <unit>`, in the header's own syntax; a file with no header, or with no such line, is returned as it was. |
| `MGFCD-B07` | `Renames`: The codes recorded in `anchors.renames.yaml` (`RenamesFile`) are read old → current, a code renamed twice resolving to the last one; a missing file is no rename. |
| `MGFCD-B08` | `WithHeaderDate`: The `updated_at` of the `@anchors` header is set to the day given; one outside the header, or a header without it, is left as it was. |

## Errors

none — every early return is normal flow: a file with no header, or a header with no such line, is returned as it was; a renames file that cannot be read is no rename (MGFCD-B07).
