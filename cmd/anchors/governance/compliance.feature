# language: en
# @anchors
#   ref: CMPLN
#   updated_at: 2026-09-26
#   layer: feature

@CMPLN
Feature: Compliance — the state of each regulatory duty, grouped by the norm that imposes it

  @CMPLN-B01 @unit-level
  Scenario: Each duty is reported under the norm that originates it
    Given a pack with authority "Local Privacy Act", a pack with no authority named local-privacy, and inline obligations
    When anchors compliance runs
    Then the pack duty is under "── Local Privacy Act", or under "── local-privacy" when the pack has no authority
    And the inline duties are under "── declared in the project", norms and duties sorted by name

  @CMPLN-B02 @unit-level
  Scenario: Each duty is marked by how many of its subjects comply
    Given a duty with 3 subjects, 1 complying and 1 debt, a duty no node triggers, and a duty whose 2 subjects comply one by waiver
    When anchors compliance runs
    Then the first reads "✗ erasure" with "Art. 9" and "3 subject,   1 complying, 1 assumed debt(s)"
    And the second is marked "·" and the third "✓" with "2 subject,   2 complying, 1 waiver(s)"

  @CMPLN-B03 @unit-level
  Scenario: A duty that no subject complies with and no debt explains warns of a disconnected target
    Given an inline duty whose only subject does not comply and another whose only subject declared the debt
    When anchors compliance runs
    Then only the first gets "⚠ NONE complies — check whether `handlers/log.go` is still the right path in the obligation's `must_appear_in:`"

  @CMPLN-B04 @unit-level
  Scenario: Verbose lists the missing nodes and the plain report only hints at them
    Given a duty that models/profile.go does not comply with
    When anchors compliance runs with --verbose
    Then "✗ models/profile.go" is listed under the duty and the "(use --verbose" hint is not printed
    And without --verbose the node is not listed and the hint is printed

  @CMPLN-B05 @unit-level
  Scenario: A project without duties says so
    Given a project with no pack and no obligation
    When anchors compliance runs
    Then it starts with "No duty declared." and lists the packs available and not adopted

  @CMPLN-B06 @unit-level
  Scenario: The embedded packs the project did not adopt are listed
    Given a project that adopted one embedded pack as ./packs/<name>.yaml
    When the available packs are printed
    Then every other embedded pack is offered and the adopted one is not
    And with every pack adopted nothing is printed
    And with every pack adopted as packs/<name>.yaml, ./packs/<name>.yml or packs/<name>.yml nothing is printed either

  @CMPLN-B07 @unit-level
  Scenario: The total counts node-duty pairs, and each hint names where its target is declared
    Given a pack duty with three subjects, an inline duty with one and another inline duty with two, one node subject to both inline duties, and nobody complying
    When anchors compliance runs
    Then the report ends with "total: 3 duty(ies), 6 subject node-duty pair(s), 0 fulfilled"
    And the pack duty's hint says "in `pack_values:`" and the inline duty's says "in the obligation's `must_appear_in:`"

  @CMPLN-I01 @unit-level
  Scenario: Every duty in force has its line even when no node is subject
    Given an inline duty ghost whose condition no file carries
    When anchors compliance runs
    Then the report shows "· ghost" with "0 subject,   0 complying"

  @CMPLN-E01 @unit-level
  Scenario: A pack missing a required value fails naming it
    Given an adopted pack that requires erasure_handler and no value for it
    When anchors compliance runs
    Then it fails with an error naming erasure_handler

  @CMPLN-E02 @unit-level
  Scenario: A project without configuration fails loading it
    Given a directory with no anchors.yaml
    When anchors compliance runs there
    Then it fails with "load config"

  @CMPLN-E03 @unit-level
  Scenario: A project without a map fails pointing at the map build
    Given a project with anchors.yaml and no map
    When anchors compliance runs there
    Then it fails with "load map" and the hint "run `anchors map build`"
