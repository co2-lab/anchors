# language: en
# @anchors
#   code: HGFHD
#   ref: HDGDH
#   updated_at: 2026-11-02
#   layer: feature

@HDGDH
Feature: HeaderGuide — render the project's header guide in the language's comment dialect, passing the gate that init itself declares

  @HDGDH-B01 @unit-level
  Scenario: The comment dialect follows the language family
    Given the families python, ruby, go, ts, java and none
    When the header guide is rendered for each
    Then python and ruby show "# @anchors"
    And every other family, and none, shows "// @anchors"

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
    Given no family and no modules
    When the header guide is rendered
    Then it mentions "code:", "updated_at:", "header-valid" and "anchors guide header"

  @HDGDH-B05 @unit-level
  Scenario: The compliance-points section is always present with five points
    Given no family and no modules
    When the header guide is rendered
    Then it has a section titled with the compliance-points title of the current language
    And it lists the points CK1, CK2, CK3, CK4 and CK5

  @HDGDH-B06 @unit-level
  Scenario: The title falls back to project
    Given the family go and no family
    When the header guide is rendered for each
    Then the first is titled "Header guide — go project"
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

  @HDGDH-B07 @unit-level
  Scenario: A seeded guide carries its header under a free code
    Given a header guide seeded twice into the same set of taken codes
    When each is given its header
    Then each opens with its @anchors header, and the second takes another code

  @HDGDH-B08 @unit-level
  Scenario: A seeded guide governs the layers it rules
    Given a guide layer, a spec layer, a code layer and a headless doc
    When the header and spec guides are given their rules, twice, and a file of no guide layer too
    Then the header guide governs the spec and code tags, the spec guide the spec, each once, and the other file nothing
