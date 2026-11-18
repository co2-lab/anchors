# language: en
# @anchors
#   code: CHFTE
#   ref: CHLGC
#   updated_at: 2026-10-03
#   layer: feature

@CHLGC
Feature: CheckLog — the check's output mirrored to a file, so it can be reread without re-running

  @CHLGC-B01 @unit-level
  Scenario: Each scope is mirrored to its own file
    Given a project root with no state folder
    When a full check and then a changed-files check are mirrored
    Then the outputs land in ".anchors/check-all.txt" and ".anchors/check-changed.txt"

  @CHLGC-B02 @unit-level
  Scenario: The header comes first and the output is copied after it
    Given a mirror opened with the header "# header"
    When the command prints "report line" and the mirror is closed
    Then the mirror file holds the header and "report line"

  @CHLGC-B03 @unit-level
  Scenario: A long output is copied whole without hanging
    Given a mirror opened for the full check
    When the command prints 2000 lines of 200 characters and the mirror is closed
    Then closing returns within 10 seconds and the file holds all 2000 lines

  @CHLGC-B04 @unit-level
  Scenario: A mirror that cannot be opened does not stop the check
    Given a project root where ".anchors" is a file, not a folder
    When the mirror is opened
    Then no mirror is returned, closing it is safe, its path is empty, and standard output is still usable

  @CHLGC-B05 @unit-level
  Scenario: The header records the command, the moment and the HEAD
    Given the command "anchors check --all" run at 2026-08-18 14:32:00 UTC on HEAD "abc1234" with subject "fix: something"
    When the header is built
    Then it holds the command, "2026-08-18 14:32:00" and "abc1234 fix: something"
    And a header built with no HEAD has no HEAD line

  @CHLGC-B06 @unit-level
  Scenario: The tree line reports the real count
    Given dirty-file counts of 0, 1 and 5
    When a header is built for each
    Then the tree lines read "clean", "1 modified file" and "5 modified file(s)"

  @CHLGC-I01 @unit-level
  Scenario: The changed-files mirror leaves the full snapshot intact
    Given a full check mirrored with the output "full snapshot"
    When a changed-files check is mirrored with the output "incremental"
    Then the full file still holds "full snapshot" and the changed file does not

  @CHLGC-X01 @unit-level
  Scenario: A tree that could not be counted is never called clean
    Given a dirty-file count of -1
    When the header is built
    Then the tree line says "unknown" and never "tree: clean"
