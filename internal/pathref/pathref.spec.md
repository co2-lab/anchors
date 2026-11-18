<!-- @anchors
  code: PTRFS
  updated_at: 2026-11-06
  layer: apoio
-->
# PathRef — the file a path of anchors.yaml names, in the project or at a URL

> **Code**: `PTRFS`

## Overview

A guide, a changelog template or a pack may live in another repository, shared by several projects.
The fields of anchors.yaml that read a file accept a URL (PTHSA), and this unit is the one place that
reads it: no command opens a URL on its own.

A URL may carry `{{branch}}`, replaced by the project's current branch, so a project on `develop` reads
the `develop` of the external repository. `{{branch|develop}}` names the branch read when the mirror
does not exist:

```yaml
aliases:
  "@std": https://raw.githubusercontent.com/co2-lab/standards/{{branch|develop}}
governs:
  - from: "@std/SCREEN_GUIDE.md"
    governs: tela
```

The current branch is a pull request's head in a CI (`GITHUB_HEAD_REF`), else the branch checked out,
else the ref a CI ran on (`GITHUB_REF_NAME`). A URL is read once per process and its copy kept in
`.anchors/cache/urls/`; with the server unreachable, the copy answers. A raw GitHub file of a private
repository is asked with `GH_TOKEN` or `GITHUB_TOKEN`, sent to GitHub's hosts only.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a path | relative to the project's root | — | the configuration (PTHSA) |
| a URL | `http://` or `https://`, with `{{branch}}` or `{{branch\|default}}` anywhere | a URL no server answers | this unit: fails the read (`PTRFS-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `PTRFS-B01` | A path reads the project's file, and a URL the external file: a 200 is its content, a 404 a file that does not exist, any other answer or no answer an error. (`Read`, `IsURL`) |
| `PTRFS-B02` | `{{branch}}` reads the current branch's mirror first and, when the server answers it does not exist or no branch is checked out, the default's; the reader is named the URL that answered. (`Resolve`, `Candidates`) |
| `PTRFS-B03` | With the server unreachable, the last copy kept of the URL answers — the branch's, else the default's. (`Read`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PTRFS-E01` | The server answers that neither the branch's nor the default's file exists; the URL asks for the branch with none checked out and no default; the server is unreachable and no copy is kept. | The read fails naming each URL tried, and so does whatever command reads it. | A guide that reads nothing would govern nothing, and the project would never hear of it. |
