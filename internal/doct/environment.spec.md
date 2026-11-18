<!-- @anchors
  code: DCENV
  updated_at: 2026-10-03
  layer: doct
-->
# Environment — the project's environment variables page, compiled from its specs

> **Code**: `DCENV`

## Overview

Whoever deploys needs one list: every environment variable the project reads, its type, whether it is
required, its default, its possible values, whether it is deprecated and what to use instead. Each
unit declares the variables it reads in its spec's `Environment Variables` section — the contract
`env-declared` confronts with the code —, and this unit compiles them into the project's page:
`{{range envVars}}` in a template under `doct/`. A variable two units read is listed once, with both.
`anchors docs init` seeds `environment.md.tmpl` when a spec declares a variable.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the specs | the specs the compiler loaded | — | the compiler |
| a variable row | a row of an Environment Variables table, in any language of the catalog | a row with no name, or a `TODO` placeholder | this unit: the row is no variable (`DCENV-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `DCENV-B01` | `envVars` lists every variable the specs declare, read under the section's title and columns in any language of the catalog, ordered by name; a `TODO` placeholder is no variable. |
| `DCENV-B02` | A variable declared by several specs is listed once, with every unit that declares it, and the first description given. |
| `DCENV-B03` | `ScaffoldEnvironment`: The environment page template, in the project's language, renders one table row per variable with the units that read it; init writes it when a spec declares a variable. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCENV-E01` | A row has no variable name, or carries the template's `TODO` placeholder. | The row is no variable and is left out of the page. | A placeholder listed as a variable tells whoever deploys to set something that does not exist. <!-- @resilient: the row is the template's own text, not a failure anyone needs to read about --> |
