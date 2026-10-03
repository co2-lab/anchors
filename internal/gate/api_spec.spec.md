<!-- @anchors
  code: APISP
  updated_at: 2026-10-03
  layer: gate
-->
# APISpec — the coherence of an API spec, and its error codes in the code

> **Code**: `APISP`

## Overview

An API spec ties three things the client relies on: which contract each body is, which status each
refusal answers with, and which error code and message come with it. Three gates ask whether the spec
holds them together, before anything is compiled or run:

| Gate | Question |
| --- | --- |
| `api-contracts-resolve` | Is every contract the body and the responses cite a spec with a Domain, and does every response say its contract? |
| `api-errors-declared` | Does every error response carry its error code and message, under a status the Responses declare? |
| `error-codes-honored` | Is every error code the spec declares one the unit's code emits? |

The OpenAPI build already fails on a contract that is no spec; asked here, per unit, the failure names
the unit while it is being written, and the pre-commit stops the commit that broke it. They run on an
API unit's main code file — the layers tagged `interface` — and read the spec beside it, when it has an
`Endpoint` section; sections and their columns are read in any language of the catalog.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: leaves without a verdict |
| the spec | `<Unit>.spec.md` beside it, with `## Endpoint` | none, or none with Endpoint | this unit: leaves without a verdict, saying why |
| a Contract cell | a spec's code in backticks, or `—` for no body | empty, `TODO`, plain text | `api-contracts-resolve`: a failure |
| the unit's code | the main file and the files its spec `specifies` | — | the map |

## Effects

| Effect | Description |
| --- | --- |
| `APISP-B01` | A node that is not code and a code file whose spec is missing or has no `Endpoint` leave every gate without a verdict. |
| `APISP-B02` | `api-contracts-resolve` fails naming each body or response (by status) whose contract is empty or `TODO`, each cited code no spec of the map carries, and each contract whose spec has no Domain table; `—` is a response with no body. |
| `APISP-B03` | `api-errors-declared` fails naming each error response whose status the Responses do not declare — exactly, or by its range (`4xx`) —, and each one without an error code or a message; a spec with no error response leaves without a verdict. |
| `APISP-B04` | `error-codes-honored` fails naming each declared error code that appears in none of the unit's code — the main file and the files its spec specifies, comments removed —, and the files it read. |
| `APISP-B05` | Sections are found under their title in any language of the catalog, and their columns by their header in those languages. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `APISP-E01` | The spec beside the code file cannot be read, or there is none. | No verdict, saying there is no spec beside the unit's main file. | A part of the unit has no spec; failing it would charge every file of an API. |
| `APISP-E02` | A contract's spec or one of the unit's code files cannot be read. | The contract counts as having no Domain; the code file is left out of what is read. | Unread is not proven: a contract that cannot be read gives no fields, and a code file that cannot be read emits nothing. |
