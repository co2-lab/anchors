# language: en
# @anchors
#   code: OPFTP
#   ref: OPNAP
#   updated_at: 2026-11-05
#   layer: feature

@OPNAP
Feature: OpenAPI — the project's API document, compiled from the specs of its API units

  @OPNAP-B01 @unit-level
  Scenario: Each Endpoint row of a spec is an operation, with its parameters, body, responses, errors, security and limits
    Given an API spec with an Endpoint, parameters, a body, responses, an error response, security and limits
    When the OpenAPI template is compiled
    Then the document has the operation with each of them

  @OPNAP-B02 @unit-level
  Scenario: A cited contract is a shared schema built from its Domain, in any language
    Given the body's contract written in Portuguese, with Tipo and Obrigatório columns
    When the OpenAPI template is compiled
    Then the schema has each field with its type, and the required ones

  @OPNAP-B03 @unit-level
  Scenario: The document is written in OpenAPI's order, two spaces deep
    Given an API spec with a body-less 204 response
    When the OpenAPI template is compiled
    Then the keys come in OpenAPI's order, indented by two spaces, and the 204 has no content

  @OPNAP-B04 @unit-level
  Scenario: The compiled document carries the generated marker as a YAML comment and is valid YAML
    When the OpenAPI template is compiled
    Then the file opens with "# anchors:generated" and parses as YAML

  @OPNAP-B05 @unit-level
  Scenario: A compiled YAML is recognised as generated and its freshness checked
    Given a compiled OpenAPI
    When a contract it reads changes
    Then the document is stale, and the next build overwrites it

  @OPNAP-B06 @unit-level
  Scenario: A Type cell citing a contract's code is a reference to its schema, and the contract a shared schema
    Given a contract whose fields are typed `QRCDM`, array of `QRCDM` and `QRCDM`[]
    When the OpenAPI template is compiled
    Then each field is a $ref to the contract's schema, in items for an array, and the contract is a schema

  @OPNAP-B07 @unit-level
  Scenario: A contract that cites itself, or one citing it back, is a reference and compiles
    Given a contract with a field typed with its own code
    When the OpenAPI template is compiled
    Then the field is a $ref to its own schema

  @OPNAP-B08 @unit-level
  Scenario: The rows of a Security table are one requirement, every scheme together
    Given an API spec whose Security table names an API key and a bearer session with two scopes
    When the OpenAPI template is compiled
    Then the operation has one security requirement holding both schemes, with their scopes

  @OPNAP-E01 @unit-level
  Scenario: A contract cited and not found fails the build, naming it
    Given an API spec whose body cites a contract no spec has
    When the OpenAPI template is compiled
    Then the build fails naming the contract

  @OPNAP-E02 @unit-level
  Scenario: An Endpoint row with no path fails the build
    Given an API spec whose Endpoint row has a method and no path
    When the OpenAPI template is compiled
    Then the build fails naming the spec

  @OPNAP-E03 @unit-level
  Scenario: A Type cell citing a code no spec has fails the build, naming it
    Given a contract whose field cites a code no spec has
    When the OpenAPI template is compiled
    Then the build fails naming the spec and the code
