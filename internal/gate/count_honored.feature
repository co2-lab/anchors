# language: en
# @anchors
#   ref: CNHNC
#   updated_at: 2026-09-26
#   layer: feature

@CNHNC
Feature: CountHonored — a numerical assertion written in a spec must match reality in code

  @CNHNC-B01 @unit-level
  Scenario: Confronting an artifact that is not a spec skips
    Given a node whose kind is not spec
    When the gate confronts it
    Then it returns Skip, because numerical assertions belong to specs

  @CNHNC-B02 @unit-level
  Scenario: A spec containing no count declarations skips
    Given a spec with no count declaration markers
    When the gate confronts it
    Then it returns Skip, because the gate only confronts explicit count contracts

  @CNHNC-B03 @unit-level
  Scenario: A declared file count matching the number of files on disk passes
    Given a spec declaring an expected file count for a glob pattern
    And the actual matching files on disk equal the expected count
    When the gate confronts it
    Then it returns Pass, confirming the file count matches reality

  @CNHNC-B04 @unit-level
  Scenario: A declared file count that differs from disk fails
    Given a spec declaring an expected file count for a glob pattern
    And the actual matching files on disk differ from the expected count
    When the gate confronts it
    Then it returns Fail, citing the expected and actual file counts

  @CNHNC-B05 @unit-level
  Scenario: A declared regex pattern counts occurrences across files instead of file count
    Given a spec declaring a regex pattern alongside a glob pattern
    And the regex occurrence count across matching files equals the expected count
    When the gate confronts it
    Then it returns Pass, confirming regex occurrences match reality

  @CNHNC-B06 @unit-level
  Scenario: A declared regex occurrence count that differs from reality fails
    Given a spec declaring a regex pattern alongside a glob pattern
    And the regex occurrence count across matching files differs from the expected count
    When the gate confronts it
    Then it returns Fail, reporting the mismatch in regex occurrences

  @CNHNC-B07 @unit-level
  Scenario: A glob pattern matching zero files fails with a path warning
    Given a spec with a glob pattern matching zero files on disk
    When the gate confronts it
    Then it returns Fail, advising to check the path before modifying the number

  @CNHNC-B08 @unit-level
  Scenario: An invalid glob expression fails reporting the glob syntax error
    Given a spec with a syntactically invalid glob expression
    When the gate confronts it
    Then it returns Fail, reporting the invalid glob syntax error

  @CNHNC-B09 @unit-level
  Scenario: An invalid regex pattern fails reporting the regex syntax error
    Given a spec with a count declaration containing an invalid regex pattern
    When the gate confronts it
    Then it returns Fail, reporting the invalid regex syntax error

  @CNHNC-B10 @unit-level
  Scenario: Divergent count stated in adjacent prose fails even if marker matches
    Given a spec where the count marker matches disk but adjacent prose asserts a different count
    When the gate confronts it
    Then it returns Fail, reporting that the prose read by humans contradicts reality

  @CNHNC-B11 @unit-level
  Scenario: Non-restrictive complements attached to prose labels are confronted as total count
    Given a spec containing prose with non-restrictive complements like of authorization or in total
    When the gate confronts it
    Then it treats the statement as a claim about the total count

  @CNHNC-B12 @unit-level
  Scenario: Qualifying words attached to prose labels indicate subsets and are not accused
    Given a spec with prose containing qualifying words like indexed or financial data after the label
    When the gate confronts it
    Then it skips confronting the subset claim, avoiding false positives on partial counts

  @CNHNC-B13 @unit-level
  Scenario: Arbitrary numbers in prose without declaration markers are ignored
    Given a spec containing numbers in prose without any count declaration markers
    When the gate confronts it
    Then it returns Skip, ignoring undeclared numbers like retention days or arbitrary limits

  @CNHNC-I01 @unit-level
  Scenario: Numerical claims without declaration markers never trigger confrontation
    Given a spec containing arbitrary numerical assertions without anchors-count markers
    When the gate confronts it
    Then it returns Skip, guaranteeing that undeclared numbers are never falsely accused

  @CNHNC-I02 @unit-level
  Scenario: Matching markers cannot conceal lying prose
    Given a spec whose declaration marker is accurate but whose prose states an outdated count
    When the gate confronts it
    Then it returns Fail, ensuring human-facing text is held to the same standard as markers

  @CNHNC-I03 @unit-level
  Scenario: Zero matched files indicates path error rather than legitimate zero count
    Given a glob pattern that matches no files in the workspace
    When the gate confronts it
    Then it returns Fail, ensuring missing directories or wrong paths are not mistaken for zero items

  @CNHNC-X01 @unit-level
  Scenario: The gate does not guess what to count from arbitrary text
    Given a spec with multiple numerical claims and no declaration markers
    When the gate confronts it
    Then it returns Skip, leaving count definitions to explicit contracts

  @CNHNC-X02 @unit-level
  Scenario: Specs without count markers are not penalized
    Given a spec that makes no numerical assertions and has no count markers
    When the gate confronts it
    Then it returns Skip, avoiding unnecessary marker mandates on qualitative specs

  @CNHNC-X03 @unit-level
  Scenario: Prose phrases qualifying subsets are not accused as total count divergences
    Given prose referring to subsets of models or items
    When the gate confronts it
    Then it avoids flagging the subset statements as mismatches against the total

  @CNHNC-B14 @unit-level
  Scenario: Only files are counted, never directories
    Given a glob that matches one file and one subdirectory
    When the gate counts the files
    Then the count is 1

  @CNHNC-E03 @unit-level
  Scenario: An unreadable file fails the pattern count naming it
    Given a glob matching a file that cannot be read
    And a count of occurrences of a pattern
    When the gate counts
    Then it returns Fail naming the unreadable file
