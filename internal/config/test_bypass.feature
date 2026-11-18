# language: en
# @anchors
#   code: TBYFA
#   ref: TBYSA
#   updated_at: 2026-11-10
#   layer: feature

@TBYSA
Feature: TestBypass — how the code skips an effect while a test runs, and the suites that run it

  @TBYSA-B01 @unit-level
  Scenario: The project declares its bypasses and the layers that run the skipped path
    Given a configuration declaring a bypass pattern and a layer
    When the configuration is loaded
    Then the pattern is compiled and the layer kept

  @TBYSA-E01 @unit-level
  Scenario: A bypass declaration that cannot hold refuses the configuration
    Given declarations with no pattern, an invalid pattern, and no layer
    When each configuration is loaded
    Then each is refused naming the field
