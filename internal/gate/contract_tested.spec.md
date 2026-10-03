<!-- @anchors
  code: CTTST
  updated_at: 2026-10-03
  layer: gate
-->
# ContractTested — an API unit is proven against the project's OpenAPI document

> **Code**: `CTTST`

## Overview

The project's OpenAPI is compiled from the specs of its API units (`anchors docs build`), so it says
what the API promises. A contract test is what says the implementation keeps it: each language has
the tool that runs requests against an OpenAPI document and validates the answers — Schemathesis or
Dredd for any stack, kin-openapi in Go, jest-openapi in JavaScript, openapi-core in Python,
swagger-request-validator in Java. Anchors does not choose the tool; it asks for the three things that
make the proof traceable: the contract scenario `{CODE}-CT` in the unit's feature, with the project's
contract regime; a test of the unit that names `{CODE}-CT`; and that test loading the OpenAPI
document, so it validates against the compiled contract and not against a copy written in the test.

It runs on the API unit's main code file — the layers the project tags `interface` — and reads the
spec beside it; a spec with no `Endpoint` section is no API unit.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: leaves without a verdict |
| the spec | `<Unit>.spec.md` beside the code, with `## Endpoint` and an `@anchors` header `code:` | none, none with Endpoint, none with a code | this unit: leaves without a verdict, saying why |
| the feature | `<Unit>.feature` beside the spec | missing | this unit: it has no contract scenario |
| the tests | the map's test nodes of the unit | none | this unit: no contract test is found |

## Effects

| Effect | Description |
| --- | --- |
| `CTTST-B01` | A node that is not code, a code file with no spec beside it, and a spec with no `Endpoint` section or no code leave without a verdict. |
| `CTTST-B02` | The feature must carry a scenario line with both `@{CODE}-CT` and the contract regime tag; otherwise the failure names both. |
| `CTTST-B03` | A test of the unit — its path names the unit's code, it sits beside the unit under its name, or a folder of its path is named after the unit — must name `{CODE}-CT`; otherwise the failure says no test names it. |
| `CTTST-B04` | A test naming `{CODE}-CT` must mention an OpenAPI document; otherwise the failure names the tests that do not load it. |
| `CTTST-B05` | The contract regime tag is the one `derived.regimes` maps to a regime naming a contract; without one, `contract-level`. |
| `CTTST-B06` | With the scenario and a contract test that loads the document, the gate passes. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CTTST-E01` | The spec beside the code file cannot be read, or there is none. | No verdict, saying there is no spec beside the unit's main file. | A part of the unit has no spec; failing it would charge every file of an API. |
| `CTTST-E02` | The feature or a test file cannot be read. | The feature has no contract scenario; the test is read by its path alone. | Unread is not proven, and a test named by its code still says what it validates. |
