<!-- @anchors
  code: CHNGL
  updated_at: 2026-10-08
  layer: infra
-->
# Changelog — the releases of a project, read from its commits

> **Code**: `CHNGL`

## Overview

The commits are the source, in Conventional Commits (the `commit-msg` hook keeps them that way).
The TYPE says what a change is to whoever reads the changelog, and two footers say what the type
cannot: `BREAKING CHANGE:` (or a `!` after the type) marks a change that breaks what was there, and
`Bug:` marks a `fix` of a defect that shipped, saying where it was seen. A `fix` without it corrects
work that never reached anyone.

A release is what the history holds between two tags. It lists its breaking changes, its features,
the bugs it fixed and its other fixes; the internal types (`refactor`, `test`, `chore`, …) are
history, not changelog. A release is rendered by a text/template — the built-in one or the
project's — whose headings are translated by the function it is given, and written into a file
that keeps what it already holds.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the commits | any subject and body; one outside Conventional Commits is left out | — | this unit |
| the repository | a git work tree | a directory git does not know | git: the command fails and the error says so |
| the template | text/template text over a release, with a `t` function | — | this unit: a template that does not parse is an error |
| the existing file | any text, possibly written by hand, possibly with markers | — | this unit |

## Effects

### Classification

| Effect | Description |
| --- | --- |
| `CHNGL-B01` | A commit that declares a breaking change — a `!` after the type, or a `BREAKING CHANGE:` footer — goes to the breaking changes whatever its type, and the footer's text is added to its subject (`Classify`). |
| `CHNGL-B02` | A `feat` goes to the features, with its scope and short hash. |
| `CHNGL-B03` | A `fix` with a `Bug:` footer goes to the bugs fixed, carrying where the bug was seen; a `fix` without it goes to the fixes. |
| `CHNGL-B04` | A commit of another type, or outside Conventional Commits, is left out; a release with nothing listed is empty. |
| `CHNGL-B05` | Footers are the paragraphs at the end of the body made only of `Key: value` lines, read from the last one back — a `Bug:` block followed by a `Co-Authored-By:` block is still the footer —; a `Bug:` line in a paragraph of prose, or before one, is not a footer (`FooterLines`). |

### Releases from the history

| Effect | Description |
| --- | --- |
| `CHNGL-B06` | The tags reachable from a ref are listed in version order, and each release holds the commits after the previous tag of that whole list, oldest first; the releases come back newest first. |
| `CHNGL-B07` | With `unreleased`, what came after the last tag becomes a release with no version, on top. |

### Rendering and writing

| Effect | Description |
| --- | --- |
| `CHNGL-B08` | The built-in template writes the version and date, or the translated "unreleased" heading, and one translated section per non-empty group; a bug entry ends with where it was seen. |
| `CHNGL-B09` | A release written into an incremental file starts with a marker naming its version (`Marker`), and `Written` lists the versions a file holds by those markers. |
| `CHNGL-B10` | Writing puts the releases the file does not hold on its top, below a leading `# ` title, and keeps everything else as it was, hand edits included; a release the file holds is not written again (`Prepend`). |
| `CHNGL-B11` | The unreleased block a previous write left is replaced by the new one, up to the next marker. |
| `CHNGL-B12` | An empty file gets the given title before the releases. |
| `CHNGL-B13` | The markers of a changelog with CRLF line endings are found as in one with LF. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CHNGL-E01` | The template does not parse, or fails while rendering. | `Render` returns the error, prefixed with `changelog template`. | A changelog written from half a template would look complete. |
| `CHNGL-E02` | A git command fails (not a repository, an unknown ref). | The error names the command and carries git's own message. | The history is the source; without it there is nothing honest to write. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
