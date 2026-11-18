# language: en
# @anchors
#   code: SGFTA
#   ref: SGSTS
#   updated_at: 2026-10-03
#   layer: feature

@SGSTS
Feature: SuggestionStore — a proposed fix, as a patch plus its reason, waiting for someone to decide

  @SGSTS-B01 @unit-level
  Scenario: A new suggestion is born in pending with its context and reason
    Given a project with no suggestions
    When the suggestion "s1" for target "anchors.yaml" is opened
    Then it is created under "suggestions/pending" holding the gate, the origin, the target and the reason

  @SGSTS-B02 @unit-level
  Scenario: A suggestion with a patch carries it in a diff block
    Given the suggestion "s1" with a patch to "anchors.yaml"
    When it is opened
    Then its file holds the patch inside a diff block

  @SGSTS-B03 @unit-level
  Scenario: A suggestion without a patch says the fix needs a human decision
    Given the suggestion "s2" with no patch
    When it is opened
    Then its file has no diff block and says "needs a human decision"

  @SGSTS-B04 @unit-level
  Scenario: Opening the same pending suggestion twice creates it once
    Given the suggestion "s1" already pending
    When it is opened again
    Then nothing is created

  @SGSTS-B05 @unit-level
  Scenario: A rejected suggestion is never reopened
    Given the suggestion "s1" opened and then rejected with a reason
    When it is opened again
    Then nothing is created and pending is empty

  @SGSTS-B06 @unit-level
  Scenario: Deciding moves the suggestion and records the reason and the decider
    Given the suggestion "s1" pending
    When a person approves it with the reason "the pattern really was wrong"
    Then approved lists "s1" and its file records the reason and "by:** pessoa"

  @SGSTS-B07 @unit-level
  Scenario: An automatic decision is marked as the AI's
    Given the suggestion "s1" pending
    When it is approved under automatic judgment
    Then its file records "IA (auto_judgment)" as the decider

  @SGSTS-B08 @unit-level
  Scenario: Listing a state gives its sorted IDs
    Given pending suggestions "b" and "a", a text file, and no rejected folder
    When pending and rejected are listed
    Then pending gives "a" then "b", and rejected gives nothing without error

  @SGSTS-B09 @unit-level
  Scenario: The patch comes out clean for git apply
    Given the pending suggestion "s1" with a patch starting "--- a/anchors.yaml"
    When its patch is extracted
    Then it starts with "--- a/anchors.yaml" and has no diff fence

  @SGSTS-I01 @unit-level
  Scenario: A decided suggestion is no longer pending
    Given the suggestion "s1" pending
    When it is approved
    Then approved lists it and pending does not

  @SGSTS-E01 @unit-level
  Scenario: A suggestion without an ID is refused
    Given a suggestion with an empty ID
    When it is opened
    Then it is refused with "suggestion without ID" and no suggestions folder exists

  @SGSTS-E02 @unit-level
  Scenario: A decision to an unknown state is refused
    Given the suggestion "s1" pending
    When it is decided to the state "pending"
    Then it is refused with "invalid decision state" and "s1" stays pending

  @SGSTS-E03 @unit-level
  Scenario: A decision without a reason is refused
    Given the suggestion "s1" pending
    When it is approved with a blank reason
    Then it is refused and "s1" stays pending

  @SGSTS-E04 @unit-level
  Scenario: Deciding a suggestion that is not pending is refused
    Given no pending suggestion "ghost"
    When it is approved with a reason
    Then it is refused with "pending suggestion not found"

  @SGSTS-E05 @unit-level
  Scenario: Asking the patch of a suggestion without one is refused
    Given the pending suggestion "s2" with no patch
    When its patch is extracted
    Then it is refused with "has no patch"
