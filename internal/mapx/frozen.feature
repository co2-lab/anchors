# language: en
# @anchors
#   code: FSFFR
#   ref: FRSCF
#   updated_at: 2026-10-20
#   layer: feature

@FRSCF
Feature: Frozen scope — the map carries each freeze to everything it covers

  @FRSCF-B01 @unit-level
  Scenario: A frozen unit covers its files, and nothing that only cites it
    Given a unit whose spec is frozen by QLDTQ-F01, its code beside it, a test that refs it with a freeze of its own, and a file citing its code in a string
    When the map is built
    Then the spec and the code are frozen by QLDTQ-F01, the test by its own, and the citing file is not frozen

  @FRSCF-B02 @unit-level
  Scenario: A frozen doctrine freezes the rules that realize it
    Given a doctrine with one frozen rule, a frozen doctrine document, and specs realizing them, one rule with its own mark
    When the map is built
    Then each realizing rule carries the doctrine's freeze, the one with its own mark keeps it, and a rule realizing nothing frozen is not frozen

  @FRSCF-B03 @unit-level
  Scenario: A freeze reads with its reason
    Given a doctrine recording QLDTQ-F01 with its reason
    When the code, an unrecorded code and a reason are described
    Then the code reads with its reason, the unrecorded code as itself, and the reason as itself

  @FRSCF-B04 @unit-level
  Scenario: The frozen parts of the project, and of one freeze
    Given a frozen unit and a spec with one frozen rule
    When the frozen parts are listed, all and for one freeze
    Then each frozen file is listed once without its rules, the frozen rule of the other spec is listed, and one freeze lists only its own

  @FRSCF-I01 @unit-level
  Scenario: No mark, no freeze
    Given the files of a project with no freeze mark
    When the map is built
    Then no node is frozen and no rule is frozen

  @FRSCF-B05 @unit-level
  Scenario: A frozen folder freezes its files and the units it owns
    Given a folder frozen in anchors.yaml holding a unit's spec and code, the unit's test in another folder, and a file of the folder with a mark of its own
    When the map is built
    Then the spec, the code and the test carry the folder's freeze, and the marked file keeps its own
