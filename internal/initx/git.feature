# language: en
# @anchors
#   code: GTFTB
#   ref: GTSTG
#   updated_at: 2026-10-03
#   layer: feature

@GTSTG
Feature: GitState — classify the project's versioning before `init` scans it, and say what to do about it

  @GTSTG-B01 @unit-level
  Scenario: With git not installed the state is not installed
    Given an isolated empty project folder
    And the caller reports that git is not installed
    When the folder is classified
    Then the state is "not installed"

  @GTSTG-B02 @unit-level
  Scenario: With git installed and no repository anywhere the state is not initialised
    Given an isolated empty project folder with no repository in it or above it
    And the caller reports that git is installed
    When the folder is classified
    Then the state is "not initialised"

  @GTSTG-B03 @unit-level
  Scenario: A repository with a branch reference is ready
    Given a project folder whose repository holds the branch reference "main"
    When the folder is classified with git installed
    Then the state is "ready"

  @GTSTG-B04 @unit-level
  Scenario: A repository with packed references is ready
    Given a project folder whose repository has no loose branch reference
    And a non-empty packed reference list naming "refs/heads/main"
    When the folder is classified with git installed
    Then the state is "ready"

  @GTSTG-B05 @unit-level
  Scenario: A repository with no reference at all has no commit yet
    Given a project folder whose repository has an empty branch folder
    And either no packed reference list or an empty one
    When the folder is classified with git installed
    Then the state is "no commit yet"

  @GTSTG-B06 @unit-level
  Scenario: A subfolder of an existing repository is ready
    Given a repository with one commit
    And the project folder "pacotes/app" inside it, with no repository of its own
    When the subfolder is classified with git installed
    Then the state is "ready", so no nested repository is offered

  @GTSTG-B07 @unit-level
  Scenario: A repository marker that is a file is ready
    Given a project folder whose repository marker is a file pointing at a worktree
    When the folder is classified with git installed
    Then the state is "ready"

  @GTSTG-B08 @unit-level
  Scenario: Only the not-initialised and no-commit states have an action to offer
    Given each of the four versioning states
    When the offer is asked for each
    Then "not initialised" and "no commit yet" have an action to offer
    And "not installed" and "ready" have none

  @GTSTG-B09 @unit-level
  Scenario: Each unready state has its own warning and ready has none
    Given each of the four versioning states, in English and in Portuguese
    When the warning is asked for each
    Then it is written in the project's language
    And "not installed" says git is not installed
    And "not initialised" says the project is not under git
    And "no commit yet" says there is no commit
    And "ready" has an empty warning

  @GTSTG-B10 @unit-level
  Scenario: The seeded ignore list covers what Anchors generates
    Given the ignore list seeded with the first commit
    When its entries are read
    Then it ignores ".DS_Store" and ".anchors/"
    And it has no exception inside ".anchors/"

  @GTSTG-I01 @unit-level
  Scenario: Not installed and not initialised stay two states with different offers
    Given one isolated empty project folder
    When it is classified once without git and once with git
    Then the states are "not installed" and "not initialised"
    And only "not initialised" has an action to offer

  @GTSTG-X01 @unit-level
  Scenario: The seeded ignore list does not guess the stack
    Given the ignore list seeded with the first commit
    When its entries are read
    Then it names none of "node_modules", "target/" or "vendor/"

  @GTSTG-X02 @unit-level
  Scenario: Whether git is installed comes from the caller
    Given a machine where git is installed
    And the caller reports that git is not installed
    When an isolated empty folder is classified
    Then the state is "not installed", because the unit never looks the program up itself
