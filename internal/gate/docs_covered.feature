# language: en
# @anchors
#   ref: DCCVD
#   updated_at: 2026-09-26
#   layer: feature

@DCCVD
Feature: DocsCovered — every spec must reach some page of the compiled documentation

  @DCCVD-B01 @unit-level
  Scenario: An artifact that is not a spec is skipped
    Given a code node in a project whose templates leave one spec out
    When the docs-covered gate confronts the code node
    Then it returns Skip

  @DCCVD-B02 @unit-level
  Scenario: A project with no templates directory is skipped
    Given a project whose documentation templates directory was removed
    And a spec that no template would reach
    When the docs-covered gate confronts the spec
    Then it returns Skip

  @DCCVD-B03 @unit-level
  Scenario: A spec no template reaches fails, naming the spec
    Given a template that selects only the specs of layer gate
    And the spec "pkg/B.spec.md" of another layer
    When the docs-covered gate confronts "pkg/B.spec.md"
    Then it returns Fail
    And the verdict names "pkg/B.spec.md"

  @DCCVD-B04 @unit-level
  Scenario: A spec a template reaches passes
    Given a template that selects the specs of layer gate
    And the spec "pkg/A.spec.md" of layer gate
    When the docs-covered gate confronts "pkg/A.spec.md"
    Then it returns Pass

  @DCCVD-B05 @unit-level
  Scenario: A template that does not compile skips with no message
    Given a template with an unterminated range action
    When the docs-covered gate confronts a spec no template would reach
    Then it returns Skip with an empty message

  @DCCVD-B06 @unit-level
  Scenario: The coverage is computed once per project root and map
    Given the spec "pkg/B.spec.md" already judged an orphan with a map
    And the template widened afterwards to select every spec
    When the docs-covered gate confronts "pkg/B.spec.md" again with the same map
    Then it still returns Fail
    And with a new map it returns Pass

  @DCCVD-I01 @unit-level
  Scenario: Another spec's orphan status never fails the confronted spec
    Given a project where "pkg/B.spec.md" is an orphan
    When the docs-covered gate confronts the reached spec "pkg/A.spec.md"
    Then it returns Pass

  @DCCVD-X01 @unit-level
  Scenario: A reached spec passes without any page having been built
    Given a project whose documentation was never built
    And a template that selects the spec "pkg/A.spec.md"
    When the docs-covered gate confronts "pkg/A.spec.md"
    Then it returns Pass
