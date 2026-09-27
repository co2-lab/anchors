# language: en
# @anchors
#   ref: CVCMC
#   updated_at: 2026-09-26
#   layer: feature

@CVCMC
Feature: CoverageCommand — answers the confidence questions from the ingested signals: by scenario, by line, of the diff and the delta

  @CVCMC-B01 @unit-level
  Scenario: The scenarios of one spec are listed as proven or not
    Given a spec declaring BBBBB-B01 and BBBBB-B02, with only BBBBB-B01 proven by a green test
    When the coverage command runs for that spec
    Then it lists "✓ BBBBB-B01 — proven by a green test" and "✗ BBBBB-B02 — no passing test"
    And it prints "1/2 scenario(s) proven" and "green tests missing for: [BBBBB-B02]"

  @CVCMC-B02 @unit-level
  Scenario: A code another unit owns, cited in prose, is not a declared scenario
    Given the InfraList spec declaring INLSN-B01 and INLSN-B02 and citing QSCOP-B02 in prose
    When its declared codes are read for the unit INLSN
    Then they are INLSN-B01 and INLSN-B02 only

  @CVCMC-B03 @unit-level
  Scenario: A spec with no unit code keeps every code it declares
    Given a spec file holding ABCDE-B01 and FGHIJ-B02
    When its declared codes are read with no unit code
    Then both codes are kept

  @CVCMC-B04 @unit-level
  Scenario: A spec changed since ingestion flags its signal as stale
    Given a spec whose revision is newer than the revision of its ingested signal
    When the coverage command runs for that spec
    Then the proven count is followed by "⚠ STALE SIGNAL"

  @CVCMC-B05 @unit-level
  Scenario: A spec with no scenario code has nothing to cover
    Given a spec in the map whose file declares no scenario code
    When the coverage command runs for that spec
    Then it prints that the spec declares no scenario codes

  @CVCMC-B06 @unit-level
  Scenario: The panorama answers by scenario, by line and by mutation
    Given a map with a spec missing a proof, a code file at 30% of lines and 25% mutation score with a stale signal, and a well covered file
    When the coverage command runs with no argument
    Then the spec is listed as "1/2 unproven", the low file under line coverage and mutation, both marked stale
    And the summary counts 1 spec, 1 file below 70% of lines, 3 surviving mutants and 1 stale signal
    And the well covered file and the fully proven spec are not listed

  @CVCMC-B07 @unit-level
  Scenario: Every measured file above the threshold is said with the count, survivors included
    Given two code files measured above the threshold, one with a surviving mutant
    When the coverage command runs with no argument
    Then it prints "✓ none of the 2 measured file(s) below the threshold"
    And it prints that the 2 measured files still have 1 surviving mutant

  @CVCMC-B08 @unit-level
  Scenario: Each panorama section shows at most fifteen entries and counts the rest
    Given seventeen code files below the threshold by line and by mutation
    When the coverage command runs with no argument
    Then each section shows fifteen files and "… and 2 more file(s) below the threshold"
    And the summary counts 17 files below 70% of lines

  @CVCMC-B09 @unit-level
  Scenario: The changed lines are crossed with the coverage report
    Given a diff adding two code lines and a comment to pkg/a.go and a line to README.md
    And a coverage report under an absolute path covering the first line and not the second
    When the coverage command runs for the diff file with a 50% threshold
    Then it prints "pkg/a.go — 2 instrumented changed line(s), 50% covered — NO test: [3]"
    And it prints "✓ what you changed is covered"
    And README.md is not listed

  @CVCMC-B10 @unit-level
  Scenario: A diff with no instrumented line has nothing to cover
    Given a diff whose changed files have no coverage in the report
    When the coverage command runs for the diff file
    Then it prints "(no instrumented code line in the diff — nothing to cover)"

  @CVCMC-B11 @unit-level
  Scenario: Changed lines covered below the threshold fail the diff coverage
    Given a diff whose instrumented changed lines are 50% covered
    When the coverage command runs for the diff file with a 70% threshold
    Then it prints "✗ below the threshold (70%)" and exits with status 1

  @CVCMC-B12 @unit-level
  Scenario: The delta with no drop says so and counts the improvements
    Given a file that improved, a file unchanged and a file with no previous coverage
    When the coverage command runs the delta
    Then it prints "✓ no file lost coverage" and "(1 file(s) improved)"

  @CVCMC-B13 @unit-level
  Scenario: A file that lost line coverage fails the delta
    Given low.go at 30% of lines after 50% at the previous ingestion
    When the coverage command runs the delta
    Then it lists "⚠ low.go — 50% → 30% (-20)" and "✗ 1 file(s) lost coverage (worst: -20 points)"
    And it exits with status 1

  @CVCMC-I01 @unit-level
  Scenario: Nothing measured is never reported as nothing below the threshold
    Given a map whose one code file has no line or mutation signal
    When the coverage command runs with no argument
    Then it prints "⚠ no line coverage ingested" and "⚠ no mutation signal ingested (1 code file(s))"
    And the summary reads "line coverage NOT measured; mutation NOT measured"

  @CVCMC-X01 @unit-level
  Scenario: The diff coverage needs no map
    Given a project with a diff file and a coverage report and no map
    When the coverage command runs for the diff file
    Then it answers without error

  @CVCMC-E01 @unit-level
  Scenario: The coverage command without a map points at the map build
    Given a project with configuration and no map
    When the coverage command runs with no argument
    Then it fails with an error naming "anchors map build"

  @CVCMC-E02 @unit-level
  Scenario: A spec outside the map is refused
    Given a map without ghost.spec.md
    When the coverage command runs for ghost.spec.md
    Then it fails saying the spec is not in the map

  @CVCMC-E03 @unit-level
  Scenario: The diff coverage without a coverage report is refused
    Given a diff file and no coverage report given
    When the coverage command runs for the diff file
    Then it fails saying "--lcov <file> is mandatory"

  @CVCMC-E04 @unit-level
  Scenario: A coverage report that cannot be read fails the diff coverage
    Given a diff file and a coverage report path that does not exist
    When the coverage command runs for the diff file
    Then it fails naming the lcov parse

  @CVCMC-E05 @unit-level
  Scenario: A git diff outside a repository explains why and offers the diff file
    Given a project that is not under git
    When the coverage command runs the diff against main
    Then it fails with an explanation that points at "--diff-file"

  @CVCMC-E06 @unit-level
  Scenario: A diff file that cannot be read fails the diff coverage
    Given a diff file path that does not exist
    When the coverage command runs for that diff file
    Then it fails
