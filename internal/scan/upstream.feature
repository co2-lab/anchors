# language: en
# @anchors
#   ref: UPOWP
#   updated_at: 2026-10-02
#   layer: feature

@UPOWP
Feature: UpstreamOwnership — which files Anchors still owns in a project, and where a file's `@anchors` header begins and ends

  @UPOWP-B01 @unit-level
  Scenario: Only a marked workflow is owned upstream
    Given ".github/workflows/anchors-claim.yml" with the marker, ".github/workflows/ci.yml" without it, and "docs/anchors-board.yml" with it
    When ownership is decided for each
    Then only the first is owned upstream

  @UPOWP-B02 @unit-level
  Scenario: A Windows path under the workflow directory is recognised
    Given the path ".github\workflows\a.yml" with the marker
    When ownership is decided
    Then the file is owned upstream

  @UPOWP-B03 @unit-level
  Scenario: The shared-code flag and prose do not open a header
    Given a file whose only comment is "@anchors-shared-code", and a file that mentions "@anchors" in prose
    When the header of each is read
    Then neither has a header

  @UPOWP-B04 @unit-level
  Scenario: An HTML or block-comment header ends with its comment
    Given an HTML header closed by "-->", a block-comment header closed by "*/", and a one-line HTML header, each followed by a body line "code: Y"
    When the header of each is read
    Then each header ends at its closing line and none contains the body line

  @UPOWP-B05 @unit-level
  Scenario: A line-comment header ends at the first line that is not a comment
    Given a "#" header followed by a blank line and another comment
    When the header is read
    Then it holds only the two header lines

  @UPOWP-B06 @unit-level
  Scenario: A header comment that never closes runs to the end of the file
    Given an HTML header without "-->"
    When the header is read
    Then it holds every line to the end of the file

  @UPOWP-I01 @unit-level
  Scenario: Nothing after the header's comment belongs to it
    Given a block-comment header followed by a body line "code: Y"
    When the header is read
    Then the body line is not part of it

  @UPOWP-X01 @unit-level
  Scenario: Only the marker and the directory decide ownership
    Given a workflow whose name starts with "anchors-" but whose content lacks the marker
    When ownership is decided
    Then the file is the project's

  @UPOWP-B07 @unit-level
  Scenario: Only a header at the top is the header, unless it says why it stands lower
    Given headers after comments and a shebang, one after a directive, one after a directive declaring @fixed-header with a reason, and one with a bare @fixed-header
    When the header is read
    Then the first ones and the declared one are headers, and the others are not and are reported off the top

  @UPOWP-B08 @unit-level
  Scenario: Only a comment whose first word is @anchors opens the header
    Given a package comment that names the header in its prose, and no header
    When the header is read
    Then there is none
    And a comment opening with @anchors at the top opens it
