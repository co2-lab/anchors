# language: en
# @anchors
#   ref: CNSTC
#   updated_at: 2026-09-29
#   layer: feature

@CNSTC
Feature: StampCommand — writes the missing contract stamps on test doubles, and refreshes them after a change

  @CNSTC-B01 @unit-level
  Scenario: With no test named, every test of the map is considered
    Given a map with a test holding an unstamped double and a test with no double
    When the stamp command runs with no test named
    Then the test with the double is listed with its stamp
    And the test with no double is not listed

  @CNSTC-B02 @unit-level
  Scenario: The missing stamp is written above the double and totalled
    Given a test whose double of the balance hook has no stamp
    When the stamp command runs
    Then the test file starts with the contract stamp of the balance hook
    And the output closes with "wrote 1 stamp(s) in 1 file(s); 0 double(s) left unstamped"

  @CNSTC-B03 @unit-level
  Scenario: The dry run says what it would write and writes nothing
    Given a test whose double has no stamp
    When the stamp command runs in dry run
    Then it prints "would write 1 stamp(s) in 1 file(s)"
    And the test file is unchanged

  @CNSTC-B04 @unit-level
  Scenario: The refresh lists the doubles stamped against the old version with the block change
    Given a committed module and a test double stamped against it
    And the module's body changed after the commit
    When the stamp command refreshes the module
    Then it names the test and line, the member, and the removed and added lines of the block
    And it tells the user to adjust the double to the new contract

  @CNSTC-B05 @unit-level
  Scenario: The refresh of a file no double is stamped against says so
    Given a module no double is stamped against
    When the stamp command refreshes the module
    Then it prints that no double is stamped against a previous version of it

  @CNSTC-B06 @unit-level
  Scenario: A refreshed block with no HEAD version says there is nothing to compare
    Given a stamped double of a module outside any git repository
    And the module's body changed
    When the stamp command refreshes the module
    Then it prints "(no HEAD version of the block to compare)"
    And the stamp is updated

  @CNSTC-B07 @unit-level
  Scenario: The refresh in dry-run lists the doubles and writes nothing
    Given a stamped double of a module whose body changed
    When the stamp command refreshes the module in dry run
    Then it lists the double and prints "would update the stamps above (--dry-run: nothing written)."
    And the test file is unchanged

  @CNSTC-B08 @unit-level
  Scenario: The block diff lists what left and what came, counting repeats
    Given an old block "a, b, c" and a new block "a, B, c, d"
    When the block diff is computed
    Then it lists "- b", "+ B" and "+ d" in that order
    And an old block with a line twice against a new block with it once lists that line once as left

  @CNSTC-I01 @unit-level
  Scenario: Stamping twice writes nothing the second time
    Given a test the stamp command has just stamped
    When the stamp command runs again over that test
    Then it reports "wrote 0 stamp(s) in 0 file(s)"

  @CNSTC-X01 @unit-level
  Scenario: A divergent stamp is left as it is
    Given a test whose double carries a stamp with a wrong hash
    When the stamp command runs over that test
    Then it writes no stamp and the test file keeps the wrong hash

  @CNSTC-E01 @unit-level
  Scenario: The stamp command without configuration fails
    Given a directory with no anchors.yaml
    When the stamp command runs
    Then it fails

  @CNSTC-E02 @unit-level
  Scenario: The stamp command without a map points at the map build
    Given a project with configuration and no map
    When the stamp command runs
    Then it fails with an error naming "anchors map build"

  @CNSTC-E03 @unit-level
  Scenario: A stamp whose anchor is gone is reported and not refreshed
    Given a stamped double of the useBalance member
    And the module renamed useBalance to useSaldo
    When the stamp command refreshes the module
    Then the stamp is listed as NOT refreshed
    And the test file is unchanged

  @CNSTC-E04 @unit-level
  Scenario: A test of the map that cannot be read is reported and skipped
    Given a map naming a test that is gone from disk
    When the stamp command runs over every test
    Then it prints that the gone test could not be read
    And it still writes the stamps of the other tests

  @CNSTC-B09 @unit-level
  Scenario: The refresh handed a test names the modules to refresh
    Given a test whose double is stamped against a module that changed
    When the refresh is handed the test instead of the module
    Then it names the module and the command, and the stamp is left as it was
