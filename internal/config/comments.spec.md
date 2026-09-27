<!-- @anchors
  code: CMMRC
  updated_at: 2026-09-26
  layer: config
-->
# CommentMarkers — which text opens a line comment in each kind of file

> **Code**: `CMMRC`

## Overview

Anchors never parses code. What it reads inside a code file (an annotation, a header, a
stamp) lives in a comment, so the only thing it needs to know about a language is where a
line comment starts. This unit is that knowledge: a table from file extension to the
prefixes that open a line comment, and the choice of the one prefix to use when Anchors has
to WRITE a comment into a file.

Reading and writing ask different questions. A reader wants every prefix the language
accepts, and nothing for an extension it does not know, so that it does not claim to have
found an annotation it has no reliable anchor for. A writer always needs an answer, and
when the extension is unknown it gets the hash comment: most script and configuration
languages use it, and a hash in a C-like file is visibly wrong to whoever reads it, while a
double slash in a Python file is a syntax error that breaks the file. Between a visible
mistake and a breaking one, the unit chooses the visible one.

Markup files have only block comments. The writer is given the opening of the block, and
closing it is the caller's job.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the extension asked of the reader | an extension exactly as a path library returns it: a leading dot, lower case | a whole path, an extension without its dot, upper case | the caller: it passes the extension of the path; this unit looks it up as given |
| the path asked of the writer | any path, with or without an extension, in any case | nothing: every path gets an answer | this unit: an unknown or missing extension falls back to the hash comment |

## Effects

| Effect | Description |
| --- | --- |
| `CMMRC-B01` | A known extension answers with all of its line-comment prefixes, in declared order (PHP accepts both the double slash and the hash) (`MarkersFor`). |
| `CMMRC-B02` | An unknown extension answers with no prefix at all, so the reader has no anchor to claim. |
| `CMMRC-B03` | The comment used to write into a path follows the path's last extension, ignoring case: a test file named after its subject is the language of its last extension (`LineCommentFor`). |
| `CMMRC-B04` | A markup file (Markdown, HTML, XML) gets the opening of a block comment; closing it is the caller's job. |
| `CMMRC-B05` | A path with no extension, or with one the table does not know, gets the hash comment. |
| `CMMRC-B06` | An extension with several prefixes gets the first one declared. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CMMRC-I01` | For every extension in the table, the comment the writer uses is the first prefix the reader lists: reading and writing never disagree about a language. | walks the whole table and compares the writer's answer with the reader's first prefix |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CMMRC-X01` | The reader looks the extension up exactly as given: it neither lowers its case nor extracts it from a path. | Normalising is the writer's job, which receives a path; the reader receives an extension, and guessing there would hide a caller that passes the wrong thing. |

## Errors

none — an unknown extension is answered, not refused: the reader answers with no prefix and the writer with the hash comment, and neither is a failure the unit handles.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the unit uses only the standard library.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
