# language: en
# @anchors
#   code: VRSNJ
#   ref: VRSNI
#   updated_at: 2026-11-18
#   layer: feature

@VRSNI
Feature: VersionCommand — a product's version, its latest release and the next version its branch asks for

  @VRSNI-B01 @unit-level
  Scenario: The command shows the sources, the release, the branch's rule and the next version
    Given a product app at 1.2.0 in app/version.json, released as app-v1.1.0, two feat since, and hotfix/* on increment patch
    When anchors version runs on develop and on hotfix/login
    Then it shows 1.2.0 in app/version.json, 1.1.0 at app-v1.1.0, semantic by default, the next 1.2.0 saying 2 feat, and 2 commits; and on hotfix/login increment by hotfix/*

  @VRSNI-B02 @unit-level
  Scenario: The latest release is the highest version of the product, not the last name
    Given the tags app-v1.9.0, app-v1.10.0 and api-v9.0.0
    When anchors version runs
    Then the release is 1.10.0 at app-v1.10.0

  @VRSNI-B03 @unit-level
  Scenario: JSON prints the state with the commits classified
    Given a product released as app-v1.1.0 and the commits feat!: a new store and fix: a crash with a Bug footer
    When anchors version --json runs
    Then the JSON holds the product, the versions, the next version, and each commit with its type, whether it breaks, and its footers

  @VRSNI-B04 @unit-level
  Scenario: Bump writes the strategy's version into every source, and only the version
    Given sources app/version.json and app/build.txt at 1.1.0, released as app-v1.1.0, and a feat since
    When anchors version bump runs
    Then both hold 1.2.0, every other byte kept, and the output names them

  @VRSNI-B05 @unit-level
  Scenario: Bump raises a part, or writes the version given
    Given the sources at 1.2.0 and the latest release 1.1.0
    When anchors version bump major runs, and then bump 3.0.1
    Then the sources hold 2.0.0, and then 3.0.1

  @VRSNI-B06 @unit-level
  Scenario: Bump to the version the sources hold writes nothing
    Given strategy manual and the sources at 1.2.0
    When anchors version bump runs
    Then it says app is already at 1.2.0, and no file changed

  @VRSNI-B07 @unit-level
  Scenario: A single product needs no name, and the current branch decides
    Given a single product app and the current branch hotfix/login
    When anchors version runs with no flag
    Then it shows app and the branch hotfix/login with its rule

  @VRSNI-E01 @unit-level
  Scenario: No block, an unknown product, or several products without a name
    Given a project with no versioning, one with app and api
    When anchors version runs, with --product web, and with no product
    Then each fails, naming the products there are

  @VRSNI-E02 @unit-level
  Scenario: A detached HEAD needs a branch named
    Given a repository at a detached HEAD
    When anchors version runs with no --branch
    Then it fails, asking for --branch

  @VRSNI-E03 @unit-level
  Scenario: Sources that disagree or a tag the scheme cannot read stop the command
    Given sources at 1.2.0 and 1.3.0, and in another project a tag app-vnext
    When anchors version runs
    Then it fails naming both files with their versions, and the tag

  @VRSNI-E04 @unit-level
  Scenario: Bump refuses a version that does not move forward, and writes nothing
    Given the sources at 1.2.0, the latest release app-v1.3.0, feat/* on none, and another product with no source
    When bump 1.1.0, bump 1.2.x and bump 1.2.5 run, bump runs on feat/x, and bump runs on the other product
    Then each fails naming the versions or what is missing, and the sources still hold 1.2.0

  @VRSNI-B08 @unit-level
  Scenario: A branch's range shows with its rule, and bump refuses a version outside it
    Given hotfix/* on increment patch with range from_base, and app released as app-v1.2.0
    When anchors version runs on hotfix/login, then bump minor and bump
    Then it shows the range from_base by hotfix/*, bump minor fails saying 1.3.0 is outside 1.2.x, and bump writes 1.2.1
