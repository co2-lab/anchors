<!-- @anchors
  code: INTNA
  updated_at: 2026-11-18
  layer: infra
-->
# Intents — the change intents, read from their files and rendered as a release

> **Code**: `INTNA`

## Overview

A change intent says what a change is to the people a changelog is for, in the product's words. The
commits say what changed in the code. Whoever makes the change writes the intent, one file per
change, under `changelog.intents` (`.anchors/changes/` by default):

```markdown
---
product: app
bump: minor          # a part of the product's scheme, or none
type: feat           # one of the declared types
---
Receipts can now be shared from the list.
```

The intents of a release are grouped by type, in the order the types are declared. They are
rendered through a template: the project's own, or the built-in one, which writes the heading the
commits' template writes and one section per type.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| an intent file | a `---` front matter with `product`, `bump`, `type`, then the text | no front matter, one not closed, another key | this unit (`INTNA-E01`) |
| the directory | `*.md` files; subdirectories (`released/`) are not read | a directory that does not exist | this unit: no intents |
| the template | text/template text over an `IntentRelease`, with a `t` function | one that does not parse | this unit: an error (`INTNA-B03`) |

## Effects

| Effect | Description |
| --- | --- |
| `INTNA-B01` | An intent reads its product, bump and type from the front matter, and its text after it, trimmed; CRLF reads as LF. (`ParseIntent`) |
| `INTNA-B02` | A directory's intents are read sorted by file name; a file that does not parse is an error of its own and the others are still read; a directory that does not exist holds none. (`ReadIntents`) |
| `INTNA-B03` | Intents are grouped by type, the declared types first in their order, another type after them, and empty types left out; the built-in template writes the release's heading — version and date, or the candidate and "unreleased" — and one section per type, its heading translated (`changelog.type.<type>`). A template that does not parse is an error. (`GroupIntents`, `RenderIntents`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INTNA-E01` | An intent with no front matter, one not closed by `---`, a key it does not know, or YAML that does not parse. | An error naming the file. | An intent half read would ship half a change. |
