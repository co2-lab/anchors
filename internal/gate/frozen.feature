# language: en
# @anchors
#   code: FGFFR
#   ref: FRGTF
#   updated_at: 2026-10-20
#   layer: feature

@FRGTF
Feature: Frozen gates — a frozen part is quiet, and told when touched

  @FRGTF-B01 @unit-level
  Scenario: A gate on a frozen part answers frozen, with the reason
    Given a failing gate on two files, one frozen by a recorded freeze, and a batch gate over the same two
    When the gates run
    Then the frozen file is frozen with the code and its reason, the other fails, the batch gate is given the other alone, and with both frozen it does not run

  @FRGTF-B02 @unit-level
  Scenario: A frozen rule leaves the questions asked of a rule
    Given a spec with a frozen rule and a live one, its feature with both scenarios, and a doctrine with a freeze entry and a frozen rule
    When the requirements, the feature's scenarios and the doctrine's realizers are asked
    Then only the live rule is required, and the doctrine asks no realizer of its freeze entry or its frozen rule

  @FRGTF-B03 @unit-level
  Scenario: The freeze marks and records are confronted
    Given a mark of an unrecorded code, a record no mark carries, and a mark with a reason
    When freeze-consistent runs, and then with the marks and records matching, and with no freeze
    Then it fails naming the unknown mark and the orphan record, then passes, then skips

  @FRGTF-B04 @unit-level
  Scenario: A change to a frozen part is told and goes on
    Given a frozen file and a spec with a frozen rule, both changed, an unchanged frozen file, and a file with no freeze
    When frozen-touched runs on each
    Then the changed ones are a declared divergence naming the freeze, the unchanged passes, and the unfrozen is skipped

  @FRGTF-B05 @unit-level
  Scenario: Frozen verdicts are counted apart
    Given a run with a frozen verdict and a pass
    When the profile is aggregated
    Then the gate counts one frozen and one pass, and the run is not barred

  @FRGTF-I01 @unit-level
  Scenario: Thawing gives back the verdicts
    Given a failing gate on a frozen node
    When it runs frozen and then thawed
    Then it is frozen, then it fails as before

  @FRGTF-E01 @unit-level
  Scenario: With no map the freezes are not confronted
    Given no map
    When freeze-consistent runs
    Then it is pending, naming the missing map

  @FRGTF-B06 @unit-level
  Scenario: A frozen folder that holds no file is named
    Given a project whose anchors.yaml freezes a folder holding files and one holding none
    When freeze-consistent confronts it
    Then it fails naming the empty glob and its freeze, and passes once that glob is gone
