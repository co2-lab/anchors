# language: en
# @anchors
#   ref: PSVPL
#   updated_at: 2026-09-19
#   layer: feature

@PSVPL
Feature: PlanSeedsValid — specifications seeded in a plan must target valid governed layers

  @PSVPL-B01 @unit-level
  Scenario: Non-plan artifacts skip confrontation
    Given an artifact node whose kind is not plan
    When the gate confronts it
    Then it returns Skip, because seed validation is evaluated on plans

  @PSVPL-B02 @unit-level
  Scenario: Missing project configuration returns Pending
    Given an implementation plan artifact
    And a nil project configuration
    When the gate confronts it
    Then it returns Pending, because layer rules cannot be checked without configuration

  @PSVPL-B03 @unit-level
  Scenario: Plans without seeded specifications skip confrontation
    Given a plan containing no backticked specification paths
    When the gate confronts it
    Then it returns Skip, because there are no seeds to validate

  @PSVPL-B04 @unit-level
  Scenario: Template specification references are ignored as templates
    Given a plan citing template files prefixed with _TEMPLATE
    When the gate confronts it
    Then it ignores the templates and skips confrontation if no other seeds exist

  @PSVPL-B05 @unit-level
  Scenario: Bare specification file names without directory paths are ignored
    Given a plan mentioning a specification by bare file name without directory slashes
    When the gate confronts it
    Then it treats the mention as casual prose rather than an actionable seed

  @PSVPL-B06 @unit-level
  Scenario: Informal path abbreviations without real top-level directories are ignored
    Given a plan referencing a path fragment whose top directory does not exist on disk
    When the gate confronts it
    Then it treats the path as an informal prose abbreviation and ignores it

  @PSVPL-B07 @unit-level
  Scenario: Seeded specifications targeting valid governed layers pass
    Given a plan seeding specification paths in declared governed layers
    When the gate confronts it
    Then it returns Pass, confirming all promised paths target valid layers

  @PSVPL-B08 @unit-level
  Scenario: Multiple target source file extensions resolve the governed layer
    Given a plan seeding specifications for source files across diverse language extensions
    When the gate confronts it
    Then it tests candidate extensions to correctly classify the target layer and passes

  @PSVPL-B09 @unit-level
  Scenario: Seeded specifications targeting declarative layers fail
    Given a plan seeding a specification in a recognized declarative layer
    When the gate confronts it
    Then it returns Fail, explaining that declarative layers do not take specifications

  @PSVPL-B10 @unit-level
  Scenario: Seeded specifications matching no declared layer in a real directory fail
    Given a plan seeding a specification in a real repository directory with no matching layer
    When the gate confronts it
    Then it returns Fail, reporting the undeclared or misspelled layer path

  @PSVPL-B11 @unit-level
  Scenario: Multiple seed defects across declarative and undeclared layers are aggregated
    Given a plan seeding specifications in both declarative layers and undeclared layers
    When the gate confronts it
    Then it returns Fail, reporting both categories sorted and combined

  @PSVPL-I01 @unit-level
  Scenario: Plan seed validation applies exclusively to plan artifacts
    Given code, feature, test, or specification artifacts confronted by the gate
    When the gate confronts them
    Then it returns Skip, keeping verification scoped to implementation plans

  @PSVPL-I02 @unit-level
  Scenario: Declarative layers reject specification seeding
    Given a plan promising specifications in declarative models or schema layers
    When the gate confronts it
    Then it returns Fail, preventing erroneous specification creation in declarative domains

  @PSVPL-I03 @unit-level
  Scenario: Seed validation evaluates structural layer validity without requiring file existence
    Given a plan promising future specifications that do not yet exist on disk
    When the gate confronts it against valid governed layers
    Then it passes, verifying path structure independently of file existence

  @PSVPL-I04 @unit-level
  Scenario: Casual prose citations and templates are not treated as seeded paths
    Given a plan containing prose discussions and template references
    When the gate confronts it
    Then it ignores non-seeded citations and avoids false alarms

  @PSVPL-X01 @unit-level
  Scenario: Existing file presence is not required for seeded specifications
    Given a plan describing future specifications
    When the gate runs
    Then it does not verify whether the files are already present on disk

  @PSVPL-X02 @unit-level
  Scenario: Plan progress synchronization is not evaluated by this gate
    Given a plan artifact with ongoing tasks
    When the gate runs
    Then it validates seeded paths without checking progress tracking synchronization

  @PSVPL-X03 @unit-level
  Scenario: Specification content and scenarios within seeded files are not verified
    Given a plan with valid seeded paths
    When the gate confronts it
    Then it checks destination layers without inspecting specification internals
