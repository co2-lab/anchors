# language: en
# @anchors
#   ref: APISP
#   updated_at: 2026-10-03
#   layer: feature

@APISP
Feature: APISpec — the coherence of an API spec, and its error codes in the code

  @APISP-B01 @unit-level
  Scenario: What is no API unit leaves every gate without a verdict
    Given a spec node and a code file whose spec has no Endpoint
    When the three gates confront each
    Then each leaves without a verdict

  @APISP-B02 @unit-level
  Scenario: Every contract cited resolves to a spec with a Domain
    Given a body citing a contract with no Domain, a 201 citing a code no spec has, a 400 with an empty contract and a 204 with "—"
    When api-contracts-resolve confronts the unit's code
    Then it fails naming the 400, the missing code and the contract with no Domain, and not the 204

  @APISP-B03 @unit-level
  Scenario: Every error response is complete and under a declared status
    Given Responses declaring 201 and 4xx, an error under 404, one under 503, and one with no message
    When api-errors-declared confronts the unit's code
    Then it fails naming the 503 one and the one with no message, and not the 404

  @APISP-B04 @unit-level
  Scenario: Every declared error code is emitted by the unit's code
    Given errors WALLET_NOT_FOUND and WALLET_FORBIDDEN, the second only in a comment
    When error-codes-honored confronts the unit's code
    Then it fails naming WALLET_FORBIDDEN

  @APISP-B05 @unit-level
  Scenario: Sections and columns are read in any language
    Given an API spec written in Portuguese
    When api-errors-declared confronts the unit's code
    Then it reads Respostas and Respostas de Erro and passes

  @APISP-E01 @unit-level
  Scenario: A code file with no spec beside it leaves without a verdict
    Given a code file with no spec beside it
    When the three gates confront it
    Then each leaves without a verdict, saying so

  @APISP-E02 @unit-level
  Scenario: A code file the spec specifies that cannot be read emits nothing
    Given a spec that specifies a file the disk does not have
    When error-codes-honored confronts the unit's code
    Then the codes only that file would emit are named
