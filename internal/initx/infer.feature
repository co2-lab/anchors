# language: en
# @anchors
#   ref: INPRN
#   updated_at: 2026-09-26
#   layer: feature

@INPRN
Feature: InferProposal — walks the project and proposes its structure deterministically, for init to confirm

  @INPRN-B01 @unit-level
  Scenario: Dependency, build and tool directories are not walked
    Given a project whose only code, twelve files per folder, and specs sit under node_modules, .git, dist, build, vendor, .next, coverage, .expo and .anchors
    When the project is inferred
    Then no code directory, no code extension and no spec is found

  @INPRN-B02 @unit-level
  Scenario: The presence of specs, features and tests is detected
    Given a project with a colocated screen: its code, spec, feature and test
    When the project is inferred
    Then specs, features and tests are all detected

  @INPRN-B03 @unit-level
  Scenario: A markdown file in a plans directory is a plan, even inside a guides directory
    Given a project with docs/guides/plans/0001.md and docs/guides/STYLE.md
    When the project is inferred
    Then the plans folder is docs/guides/plans
    And the only guide file is docs/guides/STYLE.md

  @INPRN-B04 @unit-level
  Scenario: The first guides directory is detected with every guide file in it
    Given a project with guides/FRONTEND_GUIDE.md and guides/BACKEND_GUIDE.md
    When the project is inferred
    Then the guides folder is guides, with two guide files

  @INPRN-B05 @unit-level
  Scenario: A code directory is a top directory of up to two segments holding at least ten code files, ordered by volume
    Given nine code files under small/x, ten under mid/y and fifteen under big/z
    When the project is inferred
    Then the code directories are big/z then mid/y

  @INPRN-B06 @unit-level
  Scenario: The code extensions are the five most frequent, most frequent first
    Given twelve .ts, eleven .go, ten .py, nine .rb, eight .java and seven .rs files
    When the project is inferred
    Then the code extensions are .ts, .go, .py, .rb, .java

  @INPRN-B07 @unit-level
  Scenario: Colocation is detected when at least three stems pair code with a spec, test or feature
    Given a project where only two stems pair code with a test or a spec, and a feature sits in another folder
    When the project is inferred
    Then no colocation is detected
    And a project with three paired stems is detected as colocated

  @INPRN-B08 @unit-level
  Scenario: The test handle is the known attribute used most, and only from five uses
    Given code files using data-testid more than testID, code using testID four times, and code tying data-testid with testID
    When each project is inferred
    Then the handles are data-testid, none, and testID

  @INPRN-B09 @unit-level
  Scenario: Inference carries a proposed configuration with code layers and no artifact layer
    Given a project with code directories, specs, features, tests and guides
    When the project is inferred
    Then the proposed configuration has code layers, no artifact layer and no governs rule
    And the detected artifacts are reported

  @INPRN-X01 @unit-level
  Scenario: Only the first three hundred code files are read in search of the test handle
    Given three hundred unmarked code files walked first, then twenty files marked with testID
    When the project is inferred
    Then no test handle is proposed

  @INPRN-E01 @unit-level
  Scenario: A root that cannot be walked fails the inference with no proposal
    Given a project root that does not exist
    When the project is inferred
    Then the inference fails with no proposal

  @INPRN-B10 @unit-level
  Scenario: A test named in any known dialect pairs with the code of the same stem
    Given three Go files each beside its `_test.go`, and the same with `_test.py` and `.spec.ts`
    When each project is inferred
    Then each is detected as having tests and as colocated

  @INPRN-I02 @unit-level
  Scenario: The same tree always gives the same code extensions and code directories, ties included
    Given four code directories of eleven files each, one extension per directory
    When the project is inferred twenty times
    Then every run lists the extensions and the directories in name order
    And the layer pattern lists the extensions in name order
