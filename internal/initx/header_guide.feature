# language: en
# @anchors
#   ref: HDGDH
#   updated_at: 2026-09-26
#   layer: feature

@HDGDH
Feature: HeaderGuide — render the project's header guide in the stack's comment dialect, passing the gate that init itself declares

  @HDGDH-B01 @unit-level
  Scenario: The comment dialect follows the preset
    Given every preset of the catalog
    When the header guide is rendered for each
    Then django, fastapi, python-lib, rails and phoenix show "# @anchors"
    And every other preset shows "// @anchors"

  @HDGDH-B02 @unit-level
  Scenario: The grouping example names the first module, or auth
    Given the modules "billing" and "auth"
    When the header guide is rendered
    Then the example reads "@feature: billing"
    And with no modules the example reads "@feature: auth"

  @HDGDH-B03 @unit-level
  Scenario: The module list appears only when there are modules
    Given the modules "auth" and "billing"
    When the header guide is rendered
    Then it reads "In this project: auth, billing."
    And with no modules the guide has no "In this project:" line

  @HDGDH-B04 @unit-level
  Scenario: The essentials are always present
    Given no preset and no modules
    When the header guide is rendered
    Then it mentions "code:", "updated_at:", "header-valid" and "anchors guide header"

  @HDGDH-B05 @unit-level
  Scenario: The compliance-points section is always present with five points
    Given no preset and no modules
    When the header guide is rendered
    Then it has a section titled with the compliance-points title of the current language
    And it lists the points CK1, CK2, CK3, CK4 and CK5

  @HDGDH-B06 @unit-level
  Scenario: The title falls back to project
    Given a preset with the title "Next.js (App Router)" and a preset with no title
    When the header guide is rendered for each
    Then the first is titled "Header guide — Next.js (App Router)"
    And the second is titled "Header guide — project"

  @HDGDH-I01 @unit-level
  Scenario: The seeded guide passes the checklist heading in every language
    Given each language the translation catalog supports
    When the header guide is rendered in it
    Then the gate's compliance-section heading matches the guide
    And the guide has at least one CK point

  @HDGDH-X01 @unit-level
  Scenario: Rendering writes nothing to disk
    Given an empty working folder
    When the header guide is rendered
    Then the folder is still empty and the guide comes back as text
