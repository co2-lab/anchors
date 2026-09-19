# language: en
# @anchors
#   ref: OPQSP
#   updated_at: 2026-09-19
#   layer: feature

@OPQSP
Feature: OpenQuestions — a spec with an open question is not ready to implement

  @OPQSP-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has open decisions to demand

  @OPQSP-B02 @unit-level
  Scenario: Whoever OPENED the section is confronted by its content
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because the content is what the ruler reads

  @OPQSP-B03 @unit-level
  Scenario: An open item bars the spec
    Given a spec whose open-decisions section carries one catalogued question
    When the gate confronts it
    Then it returns Fail, because while there is a question the spec does not pass as ready

  @OPQSP-B04 @unit-level
  Scenario: A section closed honestly releases the spec
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because saying "there is no question" differs from not having looked

  @OPQSP-B05 @unit-level
  Scenario: An item marked as resolved does not block
    Given a spec whose question is marked resolved, citing the rule born from it
    When the gate confronts it
    Then it returns Pass, and the question stays in the trail instead of being swept away

  @OPQSP-B06 @unit-level
  Scenario: A question with no code is charged
    Given a spec whose open-decisions section carries a question written without a code
    When the gate confronts it
    Then it returns Fail, because without identity the question is neither a traceable
      item nor survives a rewrite of the spec

  @OPQSP-B07 @unit-level
  Scenario: The count of pending decisions reads the project's own lexicon
    Given a spec whose open-decisions section carries two questions
    And the project names that section with its own wording
    When the number of pending decisions is counted
    Then it answers two, because counting zero over a section named otherwise would
      assert "no pending decision" about a spec full of them

  @OPQSP-I01 @unit-level
  Scenario: Prose is not an item
    Given a spec whose open-decisions section carries explanatory text and no catalogued item
    When the gate confronts it
    Then it returns Pass, because otherwise the author would learn to explain nothing

  @OPQSP-I02 @unit-level
  Scenario: The section boundary is respected
    Given a spec carrying catalogued items in a section that FOLLOWS the open decisions
    And the open-decisions section itself is empty
    When the gate confronts it
    Then it returns Pass, because otherwise the whole spec would read as a section of decisions

  @OPQSP-I03 @unit-level
  Scenario: Filling in what the question BECOMES does not close the question
    Given a spec whose question names the rule it is expected to become
    And the question is not marked resolved
    When the gate confronts it
    Then it returns Fail, because the intended destination and the answer given are two
      different things

  @OPQSP-X01 @unit-level
  Scenario: The gate does not judge whether the question is good
    Given a spec whose only open question is trivial
    When the gate confronts it
    Then it returns Fail all the same, because the ruler is deterministic — an open item
      exists, or it does not; judging the merit of a doubt belongs to another gate

  @OPQSP-X02 @unit-level
  Scenario: A spec with no section is a pending item, and the verdict teaches the way out
    Given a spec with rules catalogued and no open-decisions section at all
    When the gate confronts it
    Then it returns Pending, because the absence does not tell "everything was decided"
      apart from "the section was deleted"
    And the verdict says how to close it: declare that there is no question, or write
      what is not yet decided
