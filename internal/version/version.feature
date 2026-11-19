# language: en
# @anchors
#   code: VRSNB
#   ref: VRSNE
#   updated_at: 2026-11-18
#   layer: feature

@VRSNE
Feature: Version — a product's version read, ordered and moved, and the next one decided

  @VRSNE-B01 @unit-level
  Scenario: Semver versions are ordered as SemVer 2.0 says
    Given the versions 1.0.0-alpha, 1.0.0-alpha.1, 1.0.0-beta.2, 1.0.0-beta.11, 1.0.0-rc.1, 1.0.0 and 1.10.0
    When each is compared with the next
    Then each is below the next, and 1.0.0+build equals 1.0.0

  @VRSNE-B02 @unit-level
  Scenario: A calver version is read in the widths its format gives
    Given the format YYYY.0M.MICRO
    When 2026.11.2, 2026.09.0, 2026.9.0 and 2026.11 are read
    Then the first two are read and ordered by their fields, and the others are refused

  @VRSNE-B03 @unit-level
  Scenario: A number version is a non-negative integer
    Given the scheme number
    When 9, 412, -1 and 4a are read
    Then 9 is below 412, and the others are refused

  @VRSNE-B04 @unit-level
  Scenario: A pattern version is ordered by its groups
    Given the pattern r(\d+)-(\d+)
    When r14-3, r14-10 and r2-99 are read
    Then r2-99 is below r14-3, which is below r14-10

  @VRSNE-B05 @unit-level
  Scenario: Semver raises a part and resets the ones after it, and a pre-release is the release it precedes
    Given the versions 1.4.2 and 1.1.0-rc.1
    When each is raised in its major, its minor and its patch
    Then 1.4.2 gives 2.0.0, 1.5.0 and 1.4.3, and 1.1.0-rc.1 gives 2.0.0, 1.1.0 and 1.1.0

  @VRSNE-B06 @unit-level
  Scenario: Calver takes the date from today and its counter goes up only on the same date
    Given the format YYYY.0M.MICRO and the version 2026.11.2
    When it is raised on 2026-11-18 and on 2026-12-01
    Then it gives 2026.11.3 and 2026.12.0

  @VRSNE-B07 @unit-level
  Scenario: A pattern raises the group it names and resets the later ones
    Given the pattern r(?P<major>\d+)-(\d+) and the version r14-3
    When it is raised in group major, and in group 2
    Then it gives r15-0 and r14-4

  @VRSNE-B08 @unit-level
  Scenario: Manual is the version the sources hold
    Given strategy manual, the sources at 3.1.0 and the latest release 3.0.0
    When the next version is decided
    Then it is 3.1.0

  @VRSNE-B09 @unit-level
  Scenario: Increment raises the latest release, and stays there with no commit since it
    Given strategy increment of the patch and the latest release 1.4.2
    When the next version is decided with one commit since the release, and with none
    Then it is 1.4.3, and then 1.4.2

  @VRSNE-B10 @unit-level
  Scenario: Semantic raises the release by the largest change since it
    Given the latest release 1.4.2
    When the commits since it are two feat and a fix, a fix alone, a breaking change, or only chores
    Then the next version is 1.5.0 saying 2 feat and 1 fix, 1.4.3, 2.0.0, and 1.4.2

  @VRSNE-B11 @unit-level
  Scenario: Semantic with initial bumps the minor for a breaking change, and types give a type its level
    Given the latest release 0.3.1, initial true, and types perf patch and feat none
    When the commits since it are a breaking change, a perf, or a feat
    Then the next version is 0.4.0, 0.3.2, and 0.3.1

  @VRSNE-B12 @unit-level
  Scenario: Before any release, the sources hold the version, or the scheme gives its first one
    Given no release
    When semantic decides with the sources at 0.9.0, with no source, with no source and initial true, and increment decides a number with no source
    Then the next version is 0.9.0, 1.0.0, 0.1.0 and 1

  @VRSNE-B13 @unit-level
  Scenario: A Version footer forces the next version, the newest one winning
    Given strategy semantic, the latest release 1.4.2, and commits with the footers Version 2.0.0 and then Version 3.0.0
    When the next version is decided
    Then it is 3.0.0, forced, naming the newest commit

  @VRSNE-B14 @unit-level
  Scenario: None gives no version, whatever the footers say
    Given strategy none and a commit with the footer Version 2.0.0
    When the next version is decided
    Then there is no version, and the reason says the branch never versions

  @VRSNE-E01 @unit-level
  Scenario: A scheme that reads nothing is refused
    Given the schemes date, calver with no format, calver YYYY.MICRO.0M, calver YYYY.WW, a pattern ( and a pattern with no group
    When each is built
    Then each fails, naming what is wrong

  @VRSNE-E02 @unit-level
  Scenario: A version the scheme cannot read is an error, never smaller nor equal
    Given strategy manual with the sources at 1.2 under semver, and a latest release v1.0.0
    When the next version is decided
    Then each fails, naming the version and the scheme

  @VRSNE-E03 @unit-level
  Scenario: A next version below the latest release is refused
    Given the latest release 2.0.0, and manual with the sources at 1.9.0, or a footer Version 1.0.0
    When the next version is decided
    Then it fails, naming both versions

  @VRSNE-E04 @unit-level
  Scenario: A strategy that would have to guess is refused
    Given increment of a part semver lacks, semantic under number, manual with no source, increment of a pattern with no release and no source, a calver with no counter raised on its own date, and a strategy named latest
    When the next version is decided
    Then each fails, naming the strategy and what it lacks

  @VRSNE-B15 @unit-level
  Scenario: intents raises the latest release by the largest bump of the intents waiting; none or no intent releases nothing
    Given the release 1.2.3 and the intents patch, minor and none; then patch and major; then none alone; a number at 41; and no release with the sources at 0.1.0
    When the intents strategy computes
    Then 1.3.0 by 2 intents, 2.0.0, nothing to release, 42, and 0.1.0 as no release yet

  @VRSNE-B16 @unit-level
  Scenario: promote and script take the version decided outside, saying by what, and fail with the reason when none was
    Given promote from review deciding 1.1.0, a script that printed nothing, and a script deciding 1.9.0 below the release 2.0.0
    When the next version is computed
    Then it is 1.1.0 by what review carries, and the others fail naming the script and why
