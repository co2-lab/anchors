<!-- @anchors
  code: ENVDC
  updated_at: 2026-10-03
  layer: gate
-->
# EnvDeclared — the environment variables a unit reads are the ones its spec declares

> **Code**: `ENVDC`

## Overview

An environment variable is a contract with whoever deploys. A variable the code reads and no spec names
is found missing in production; a variable the spec names and the code no longer reads is configured
forever for nothing. Like `contract-status-declared` with the status codes, this gate confronts both
sides: the `Environment Variables` section of the spec against the reads in the code it specifies. A
variable declared deprecated may stop being read — it is on its way out.

The code is read with the dialect's `env_read`: the project's own pattern, its language family's
(`os.Getenv`, `process.env`, `os.environ`, `System.getenv`, `ENV[...]`, `getenv`,
`Environment.GetEnvironmentVariable`), or every family's together when the project declares none. A
read by a computed name is out of reach of the text.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a spec | any other kind | this unit: leaves without a verdict |
| the declared variables | the `Variable` cells of the Environment Variables table, in any language | a `TODO` placeholder | this unit: not a variable |
| the code | the files the spec `specifies`, comments removed | — | the map |
| the read pattern | a regular expression whose first non-empty group is the name | one that does not compile | this unit: reads nothing |

## Effects

| Effect | Description |
| --- | --- |
| `ENVDC-B01` | A node that is not a spec, and a spec with no variable declared whose code reads none, leave without a verdict. |
| `ENVDC-B02` | Each variable the code reads and the spec does not declare is named. |
| `ENVDC-B03` | Each variable the spec declares and the code does not read is named — unless it is declared deprecated (`yes`, or `yes: use X`). |
| `ENVDC-B04` | A read inside a comment is no read. |
| `ENVDC-B05` | With no language family declared, the reads of every family are recognised. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ENVDC-E01` | A specified code file cannot be read. | It is left out of what is read. | A file that cannot be read reads no variable; the declaration then stands alone and is confronted as such. |
