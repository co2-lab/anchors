# language: en
# @anchors
#   code: CUCFC
#   ref: CRUCT
#   layer: feature

@CRUCT
Feature: CrossUnitCitation — a rule of another unit a spec cites lives in the product

  @CRUCT-B01 @unit-level
  Scenario: A cited rule of another unit must realize the product; references are no citation
    Given a home spec citing a wallet rule in a rule, realizing a product rule, aliasing its own, retiring a rule, naming the wallet in its navigation and in a revision
    When cross-unit-citation runs, and again once the wallet rule realizes a product rule
    Then only the cited wallet rule is named, with its line, and then the spec passes
