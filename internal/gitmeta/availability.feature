# language: en
# @anchors
#   ref: GTAVG
#   updated_at: 2026-09-26
#   layer: feature

@GTAVG
Feature: GitAvailability — why a git operation cannot happen, named with its fix

  @GTAVG-B01 @unit-level
  Scenario: Without git on the PATH the answer is no binary
    Given a PATH that holds no git executable
    And a project folder that holds a ".git" folder
    When the root is classified
    Then the answer is "no git binary"

  @GTAVG-B02 @unit-level
  Scenario: A repository above the root is seen
    Given git on the PATH and a ".git" folder two levels above the folder "pacotes/app"
    When "pacotes/app" is classified
    Then the answer is "available"

  @GTAVG-B03 @unit-level
  Scenario: A folder outside any repository is told apart from one inside
    Given git on the PATH and a temporary folder outside any repository
    When it is classified before and after a ".git" folder is created in it
    Then the answers are "no repository" and then "available"

  @GTAVG-B04 @unit-level
  Scenario: The missing binary is explained with the install fix
    Given the answer "no git binary" and the action "list the staged files"
    When it is explained
    Then the sentence names the action, says git is not installed, and does not mention "git init"

  @GTAVG-B05 @unit-level
  Scenario: The missing repository is explained with the init fix
    Given the answer "no repository" and the action "list the staged files"
    When it is explained
    Then the sentence names the action, mentions "git init", and does not say git is not installed

  @GTAVG-B06 @unit-level
  Scenario: An available git is not explained
    Given the answer "available"
    When it is explained
    Then the explanation is empty
