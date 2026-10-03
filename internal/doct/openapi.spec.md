<!-- @anchors
  code: OPNAP
  updated_at: 2026-10-03
  layer: doct
-->
# OpenAPI — the project's API document, compiled from the specs of its API units

> **Code**: `OPNAP`

## Overview

An OpenAPI written by hand next to the specs is a second copy of the same contract, and the two drift in
silence. Here the spec is the source and the OpenAPI is compiled, like every page of `docs/`: the
template `doct/openapi.yaml.tmpl` holds `{{ openapi "Title" "1.0.0" "https://server" }}`, and
`anchors docs build` writes `docs/openapi.yaml` with the generated marker, so `docs-fresh` says when it
is out of date and the contract tests that validate the API against it run against what the specs say.

Every spec with an `Endpoint` section is an API unit, and each row of it an operation. Its `Parameters`,
`Request Body`, `Responses`, `Error Responses`, `Security` and `Limits` become the operation's fields,
read under their title in any language of the catalog. The body and each response CITE their contract
by code — the model's spec, in the domain — and the contract's `Domain` table, with its `Type` and
`Required` columns, becomes a shared schema, written once however many operations use it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the specs | the specs the map knows, loaded by the compiler | — | the compiler |
| the title, version and servers | the template's arguments | — | the template |
| a Contract cell | a spec's code in backticks, or `—` for no body | any other text | this unit: fails the build (`OPNAP-E01`) |
| a Type cell | `string`, `integer`, `number`, `boolean`, `object`, a format in parentheses, `array of X` or `X[]` | any other word | this unit: kept as the field's description, with no type |

## Effects

| Effect | Description |
| --- | --- |
| `OPNAP-B01` | Each Endpoint row of a spec is an operation under its path and method, with the operation name, the spec's title as summary and its overview as description, `deprecated` when the row says so, and the spec's parameters, body, responses (a range `4xx` as `4XX`), the error codes and messages of `Error Responses` under each one's status, its security schemes and its limits. |
| `OPNAP-B02` | A contract cited by the body or a response is a shared schema named after the contract's spec file, built from its Domain — one property per input, the type from `Type`, the required ones from `Required`, what it accepts as the description —, whatever language the section and its columns are written in. |
| `OPNAP-B03` | The document's keys come in OpenAPI's order and are indented by two spaces; a response with no contract has no content. |
| `OPNAP-B04` | The compiled YAML opens with the generated marker as a `#` comment, and is valid YAML. |
| `OPNAP-B05` | A compiled YAML is recognised as generated — overwritten, never skipped — and goes stale when a spec it reads changes. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `OPNAP-E01` | A Contract cell cites no spec's code, or a code no spec of the project has. | The build fails naming the spec and the contract. | An OpenAPI with a hole would compile green. |
| `OPNAP-E02` | An Endpoint row has no method or no path. | The build fails naming the spec. | An operation with no address cannot be written. |
