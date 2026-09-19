# language: en
# @anchors
#   ref: MRPRM
#   updated_at: 2026-09-19
#   layer: feature

@MRPRM
Feature: MarkerParity — the same rule has to appear at both ends that fulfil it

  @MRPRM-B01 @unit-level
  Scenario: A rule marked at both declared scopes passes
    Given a declaration whose scopes are the page tree and the server tree
    And the same rule name marked once in each of them
    When the gate confronts it
    Then it returns Pass, because the mapping still has its two ends

  @MRPRM-B02 @unit-level
  Scenario: A rule missing from one end fails, and the verdict names the empty scope
    Given the rule marked in the page tree and absent from the server tree
    When the gate confronts it
    Then it returns Fail
    And the verdict names the server scope, because neither side looks wrong on its own

  @MRPRM-B03 @unit-level
  Scenario: Two markings on the same side do not satisfy the gate
    Given the same rule marked twice inside the page tree and never in the server tree
    When the gate confronts it
    Then it returns Fail, because the two add up to the expected count and would hide
      exactly the mismatch the gate exists to catch

  @MRPRM-B04 @unit-level
  Scenario: A rule left over at one end fails and is named
    Given both ends marked for one rule
    And a second rule marked only in the page tree, left behind when the server dropped it
    When the gate confronts it
    Then it returns Fail naming that leftover rule

  @MRPRM-B05 @unit-level
  Scenario: Total absence of the prefix is not approval
    Given a declared prefix that appears nowhere in the tree
    When the gate confronts it
    Then it returns Pending asking to check marker_prefix, because absence is almost
      always a typo in the declaration and Pass would make the gate look vigilant while
      watching nothing

  @MRPRM-B06 @unit-level
  Scenario: A declaration with no prefix returns Pending
    Given a gate declaration whose marker_prefix is empty
    When the gate confronts it
    Then it returns Pending asking for the prefix, because with no prefix there is
      nothing to confront

  @MRPRM-B07 @unit-level
  Scenario: With no scopes declared the ruler is the count
    Given a declaration with no scopes and a required count of two
    And the rule marked in two files anywhere in the tree
    When the gate confronts it
    Then it returns Pass
    And a tree carrying only one of those markings returns Fail

  @MRPRM-B08 @unit-level
  Scenario: The marking crosses language
    Given the rule marked in a TypeScript file of the page tree
    And marked again in a Go file of the server tree
    When the gate confronts it
    Then it returns Pass, because the mapping between the two ends is the same mapping
      whatever language writes each end

  @MRPRM-B09 @unit-level
  Scenario: Ignored directories never count towards parity
    Given both declared ends marked
    And a third copy of the marking inside node_modules
    When the gate confronts it
    Then it returns Pass, and the vendored copy is not counted as an end

  @MRPRM-I01 @unit-level
  Scenario: The failing verdict names the rule and the empty scope
    Given a rule whose server end was never marked
    When the gate confronts it
    Then the verdict carries both the rule name and the scope left empty, so the reader
      does not have to diff the two trees to find the orphan

  @MRPRM-I02 @unit-level
  Scenario: What was not measured is never approved
    Given in turn a declaration with no prefix, one with neither count nor scopes, and
      a prefix that appears nowhere
    When the gate confronts each of them
    Then none of them returns Pass, because approving without having looked would stamp
      what was never measured

  @MRPRM-X01 @unit-level
  Scenario: The gate does not read what each end actually does
    Given both ends marked with the same rule name
    And the page promising a list of items the handler does not erase
    When the gate confronts it
    Then it returns Pass, because presence is deterministic and agreement of meaning is
      not — this gate separates "one end" from "both ends"

  @MRPRM-X02 @unit-level
  Scenario: The gate does not decide which rules live at two ends
    Given a project whose Structure declares no marker-parity gate for a rule a reviewer
      would consider obviously two-ended
    When the gate is asked to run
    Then nothing is charged, because the catalogue of two-ended rules belongs to the
      project and a gate that invented parities would charge what nobody committed to

  @MRPRM-X03 @unit-level
  Scenario: Files outside the text extension list are not read
    Given both declared ends marked
    And a binary asset carrying the same byte sequence
    When the gate confronts it
    Then the asset is not scanned, because erring low costs a marking in an exotic place
      and erring high costs megabytes read on every walk
