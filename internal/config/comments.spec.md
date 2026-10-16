<!-- @anchors
  code: CMMRC
  updated_at: 2026-10-16
  layer: config
-->
# Comments — how the project's files write a comment, declared in the dialect

> **Code**: `CMMRC`

## Overview

Anchors never parses code. What it reads inside a code file lives in a comment:
- an annotation;
- the `@anchors` header;
- a flag of the chains.

The evidence of a proof leaves comments out. So Anchors must know where a comment is. That
knowledge is the language's, not Anchors': the project declares it in `dialect.comments`
(DESIGN-core-without-syntax.md).

Each entry is for the files a glob matches, and gives:
- the markers of a line comment;
- the delimiters of a block comment;
- the delimiters of a string, inside which a marker is text;
- whether a string may span lines;
- the lines whose end takes no comment.

One reader turns a file's comments into spaces, so every caller reads the same:
- the gates;
- the map's header and flags;
- the evidence;
- the fixers that write a flag.

The writer takes its marker from the same entry. A file no entry covers has no comment: no
header is read from it and none is written, and the gates that need one say so.

Anchors' own formats are the exception: its specs, plans and docs in Markdown, its features
in Gherkin. Their comments are Anchors' syntax. Each family has an example of the entries,
for its stack and the files beside it; `anchors init` and `anchors migrate` write it, and the
doctor offers it. It is never applied by itself.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared entries | a list under `dialect.comments`, each with a glob, at least one marker, valid delimiters | an entry with no marker, an empty marker, a glob or a pattern that does not compile | the config's load, which refuses it naming the entry |
| the file asked about | a path relative to the root, with `/` or `\` | nothing: a path no entry covers gets no syntax | this unit |
| the content read | any text | — | this unit |

## Effects

| Effect | Description |
| --- | --- |
| `CMMRC-B01` | A file's syntax is Anchors' own for its formats — Markdown (`<!-- -->`), Gherkin (`#`) —, else the first entry whose glob matches it, a glob written `**/x` reading files at the root too; a file nothing covers has none. (`CommentsFor`) |
| `CMMRC-B02` | Blanking a content turns every comment — a line marker to the end of its line, a block from its opener to its closer, across lines — into spaces, keeping each line and offset; a marker inside a string is text, a string ends at its delimiter or, unless the delimiter spans lines, at the end of its line, and a backslash escapes in a one-character string. No syntax blanks nothing. (`Blank`) |
| `CMMRC-B03` | A line is only a comment when it is not blank and blanks entirely, a block's inner lines included; the code of a content is its lines blanked, trimmed on the right, the lines that were only a comment dropped; a line cut is the line without its comment. (`CommentOnly`, `Code`, `CutLine`) |
| `CMMRC-B04` | The writer's marker is the first line marker, else the first block's opener, closed with its closer after a space; a line the entry's `no_line_after` matches takes no comment at its end. (`LineMarker`, `CloseFor`, `TakesLineComment`) |
| `CMMRC-B05` | An entry with no marker, an empty marker, a glob that is no glob, or a `no_line_after` that does not compile refuses the config, naming it; a block is written `[open, close]`. (`Delims`) |
| `CMMRC-B06` | Each family has an example of the entries for its stack and the configuration and scripts beside it — `ts` with `//`, `{/* */}` and `/* */`, a template literal across lines and no comment after a JSX tag; `python` with `#` and triple-quoted strings; `ocaml` with `(* *)` and dune's `;` —, and every family's together for migrating a project; none is applied by itself. (`FamilyComments`, `AllComments`) |
| `CMMRC-B07` | The entries of the examples that cover a set of files are those of every family whose glob matches one of them, in the examples' order, Anchors' own formats needing none; the extensions no example covers — a file's name when it has none — are given sorted, for `init` and `migrate` to write and name. (`CommentsCovering`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CMMRC-I01` | Blanking never changes a content's length nor the place of a line break: an offset or a line number read on the blanked content is the original's. | blanks contents of every example and compares lengths and the positions of `\n` |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CMMRC-X01` | No language's comment is built in. | A built-in table guessed for files the project never declared, and a project in a language Anchors had not seen could not say how it comments. |

## Errors

none — a file no entry covers is answered with no syntax, and a malformed entry is refused by the config's load (`CMMRC-B05`), not handled here.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the unit uses only the standard library and the glob matcher.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
