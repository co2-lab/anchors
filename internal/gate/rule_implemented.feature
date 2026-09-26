# language: en
# @anchors
#   ref: RLIMR
#   updated_at: 2026-09-26
#   layer: feature

@RLIMR
Feature: RuleImplemented — the spec catalogues rules, and the code shows it realised them

  @RLIMR-B01 @unit-level
  Scenario: A spec whose rules the code ignores is accused, and the verdict names them
    Given a spec cataloguing three rules, of which the code marks only the first
    When the gate confronts it
    Then it returns Fail
    And the verdict names the rule left without realisation, so the reader does not
      have to diff spec against code by hand

  @RLIMR-B02 @unit-level
  Scenario: A rule waived with a written reason closes the account
    Given a spec cataloguing a restriction the code cannot mark, because it is satisfied
      by the ABSENCE of code
    And the rule's row carries a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, because the declaration is the answer the gate asked for

  @RLIMR-B03 @unit-level
  Scenario: A unit that predates the practice is a pending item, not a failure
    Given a spec whose rules carry no mark anywhere in the code
    And the project does not declare that it requires marking
    When the gate confronts it
    Then it returns Pending
    And the verdict NAMES the debt, instead of pretending approval

  @RLIMR-B04 @unit-level
  Scenario: Declaring the requirement turns the pending item into a failure
    Given the same unmarked spec
    And the project declares in its Structure that marking is required
    When the gate confronts it
    Then it returns Fail, because declaring the requirement is the act of saying
      "here the migration is over"

  @RLIMR-B05 @unit-level
  Scenario: A spec with no linked code is not this gate's subject
    Given a spec that no code realises yet
    When the gate confronts it
    Then it returns Skip, because without the piece on the other side there is no
      confrontation to make — and who accuses the absence is the triad gate

  @RLIMR-B06 @unit-level
  Scenario: A waiver with no named rule covers every rule of the spec
    Given a spec whose waiver marker names no rule in particular
    When the gate confronts it
    Then it returns Pass, because the waiver was declared for the unit as a whole

  @RLIMR-I01 @unit-level
  Scenario: Requiring the marking never punishes whoever already marks
    Given a spec whose every rule is marked in the code
    And the project declares that marking is required
    When the gate confronts it
    Then it returns Pass, because whoever did the work before the requirement cannot
      fail for having done it

  @RLIMR-I02 @unit-level
  Scenario: The identity survives a rename
    Given a code marked with the identity the unit carried before being renamed
    When the gate confronts it
    Then it returns Pass, because losing the mark on a rename would turn identity
      stability into new debt

  @RLIMR-X01 @unit-level
  Scenario: The gate does not judge whether the implementation is correct
    Given a spec whose rules are all marked in the code
    And the marked code does something other than what the rule describes
    When the gate confronts it
    Then it returns Pass, because the ruler here is deterministic — the mark exists, or
      the waiver exists with a reason; whether the code honours the rule is judgment

  @RLIMR-X02 @unit-level
  Scenario: The gate does not demand a mark on EVERY rule
    Given a spec whose restrictions are satisfied by the absence of code
    And those rows declare their waiver with a reason
    When the gate confronts it
    Then it returns Pass, because absence has nowhere to receive a comment — demanding
      it would produce thousands of findings and teach the team to ignore the list

  @RLIMR-E01 @unit-level
  Scenario: A code file that cannot be read is pending, naming the file
    Given a spec whose code file is on disk without read permission
    When the gate confronts it
    Then it returns Pending naming the file that could not be read
