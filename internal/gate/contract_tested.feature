# language: en
# @anchors
#   ref: CTTST
#   updated_at: 2026-10-03
#   layer: feature

@CTTST
Feature: ContractTested — an API unit is proven against the project's OpenAPI document

  @CTTST-B01 @unit-level
  Scenario: What is no API unit leaves without a verdict
    Given a spec node, a code file with no spec beside it, and a spec with no Endpoint section
    When contract-tested confronts each
    Then each leaves without a verdict

  @CTTST-B02 @unit-level
  Scenario: The feature must carry the contract scenario
    Given an API unit whose feature has no @QRGEN-CT scenario
    When contract-tested confronts the unit's code
    Then it fails naming QRGEN-CT and the contract regime tag

  @CTTST-B03 @unit-level
  Scenario: A test of the unit must name the contract code
    Given an API unit with the scenario and no test naming QRGEN-CT
    When contract-tested confronts the unit's code
    Then it fails saying no test names QRGEN-CT

  @CTTST-B04 @unit-level
  Scenario: The contract test must load the OpenAPI document
    Given a test naming QRGEN-CT that never mentions an OpenAPI document
    When contract-tested confronts the unit's code
    Then it fails naming the test that does not load it

  @CTTST-B05 @unit-level
  Scenario: The contract regime tag comes from the project
    Given derived.regimes maps "nivel-contrato" to "contract"
    When the contract regime tag is read
    Then it is nivel-contrato, and contract-level without the mapping

  @CTTST-B06 @unit-level
  Scenario: Scenario and a contract test that loads the document pass
    Given the scenario and a test naming QRGEN-CT that loads docs/openapi.yaml
    When contract-tested confronts the unit's code
    Then it passes

  @CTTST-E01 @unit-level
  Scenario: A code file with no spec beside it leaves without a verdict
    Given a code file with no spec beside it
    When contract-tested confronts it
    Then it leaves without a verdict, saying so

  @CTTST-E02 @unit-level
  Scenario: A test that cannot be read is read by its path
    Given a flow QRGEN-CT-openapi.yaml the map lists and the disk does not have
    When contract-tested confronts the unit's code
    Then the flow counts as the contract test that loads the document
