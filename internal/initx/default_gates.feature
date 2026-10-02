# language: en
# @anchors
#   ref: DFGTD
#   updated_at: 2026-10-02
#   layer: feature

@DFGTD
Feature: DefaultGates — the gates a project is born with, by artifact and by project age, and the canonical gate catalog

  @DFGTD-B01 @unit-level
  Scenario: A project with no artifact chosen is seeded with no gate
    Given no artifact chosen
    When the default gates are seeded for an existing and for a new project
    Then both lists are empty

  @DFGTD-B02 @unit-level
  Scenario: A project with spec, feature and test is born with the gates of each
    Given an existing project that chose spec, feature and test
    When the default gates are seeded
    Then the spec-complete, feature-not-empty, tests-green, line-coverage and scenario-coverage gates are present
    And none is blocking except the gates blocking by nature

  @DFGTD-B03 @unit-level
  Scenario: The gates that cross spec and feature are seeded only when both are chosen
    Given the choices test alone, spec with test, and feature with test
    When the default gates are seeded for each
    Then none has scenario-coverage nor spec-feature-match
    And the choice spec with feature has both
    And feature-spec-match is seeded only with spec and feature, and test-feature-match only with test and feature

  @DFGTD-B04 @unit-level
  Scenario: Choosing only guides seeds only the guide checklist gate
    Given a project that chose only guides
    When the default gates are seeded
    Then the only gate is guide-checklist

  @DFGTD-B05 @unit-level
  Scenario: The parent gate confronts exactly the chosen artifacts among spec, plan and code
    Given a project that chose spec, code and feature
    When the default gates are seeded
    Then parent-valid runs on spec and code
    And a project that chose only feature and test has no parent-valid

  @DFGTD-B06 @unit-level
  Scenario: An existing project is born with its gates informative, except the two blocking by nature
    Given an existing project that chose every artifact
    When the default gates are seeded
    Then the blocking gates are exactly no-secret-leaked and doc-required

  @DFGTD-B07 @unit-level
  Scenario: A new project is born with its gates blocking
    Given a new project that chose spec, feature and test
    When the default gates are seeded
    Then every gate that does not read an ingested report is blocking

  @DFGTD-B08 @unit-level
  Scenario: A gate that depends on an ingested signal stays informative even in a new project
    Given a new project that chose spec, feature, test and code
    When the default gates are seeded
    Then tests-green, line-coverage, coverage-delta, mutation-score, scenario-coverage, sbom-generated, dependency-vulnerable, no-duplication, license-compatible, circular, deadcode and spellcheck are seeded and informative

  @DFGTD-B09 @unit-level
  Scenario: Every judgment gate that asks about code or a test carries the @TBD instruction
    Given a project that chose every artifact
    When the judgment gates are inspected
    Then every judgment gate question mentions @TBD, except the one about a permanent test waiver, which does not
    And at least one judgment gate was inspected

  @DFGTD-B10 @unit-level
  Scenario: The canonical declaration of a gate is found by name in the catalog of every artifact
    Given the canonical gate catalog
    When guide-checklist, plan-seeds-valid, tests-green and no-such-gate are resolved
    Then the first three are found under their own names
    And no-such-gate is not found

  @DFGTD-B11 @unit-level
  Scenario: Loading a configuration completes a canonical gate declared by name alone
    Given a configuration file that declares only the name guide-checklist
    When the configuration is loaded
    Then the gate carries the check, the measure and the targets of the canonical declaration

  @DFGTD-I01 @unit-level
  Scenario: The list of seeded gates does not change with the age of the project
    Given the choice spec, feature and test
    When the default gates are seeded for a new and for an existing project
    Then both lists have the same gates in the same order

  @DFGTD-I02 @unit-level
  Scenario: Every gate of the full catalog has a unique name that is also its id
    Given a project that chose every artifact
    When the default gates are seeded
    Then no name appears twice
    And every id equals its name

  @DFGTD-I03 @unit-level
  Scenario: Every canonical name the migration renames a legacy gate to is a default gate
    Given the legacy-to-canonical gate renames of every migration step
    When each target is followed through the later steps and looked up in the full catalog
    Then every name a target lands on is found
    And format 2's `trinca-completa` lands on `unit-complete` through format 6

  @DFGTD-X01 @unit-level
  Scenario: No default gate carries a legacy name
    Given the legacy gate names of every migration step
    When the full catalog is compared with them
    Then no default gate has a legacy name

  @DFGTD-B12 @unit-level
  Scenario: The gate is of the blocking class
    Given the default gates of a new project and of an existing one
    When revision-orphans is seeded
    Then it is blocking in the new project and informative in the existing one
    And it is not among the gates that depend on an ingested signal

  @DFGTD-B13 @unit-level
  Scenario: Choosing every artifact init offers seeds every gate of the catalog
    Given every artifact name init offers, chosen
    When the default gates are seeded
    Then every gate of the canonical catalog is seeded, no-secret-leaked and layer-boundary among them

  @DFGTD-B14 @unit-level
  Scenario: The gates that run on specs, features and tests are seeded without plans
    Given a project that chose spec, feature and test, and no plans
    When the default gates are seeded
    Then docs-fresh, the doctrine, flag and failure gates, feature-spec-match, test-feature-match, doc-required and plan-change-justified are seeded
    And plan-change-justified runs on specs alone
    And plan-seeds-valid, plan-source-declared, plan-doctrine-exists, plan-revised and phase-ordered are not seeded

  @DFGTD-B15 @unit-level
  Scenario: The gate names registered for the vocabulary check are the full catalog
    Given the names registered with the configuration package
    When they are compared with the canonical catalog
    Then they are the same names in the same order

  @DFGTD-B16 @unit-level
  Scenario: Choosing specs seeds header-valid on every governed kind
    Given a project that chooses specs
    When the default gates are seeded
    Then header-valid is among them, informative, on spec, feature, code, test, guide, doc, plan, product and flag
    And a configuration naming header-valid alone inherits that on

  @DFGTD-B17 @unit-level
  Scenario: no-duplication is the native duplication check on code files
    Given the canonical catalog
    When no-duplication is looked up
    Then it runs the duplication check on code, needs npx, and declares no command or project scope

  @DFGTD-I04 @unit-level
  Scenario: Choosing plans adds only gates that run on plans
    Given the unit, and the unit with code and guides
    When each is seeded with and without plans
    Then every gate that plans added runs on plans

  @DFGTD-X02 @unit-level
  Scenario: No default gate writes Portuguese into the project
    Given every default gate of every artifact
    When its question, measure and install hint are read
    Then none carries Portuguese letters or words

  @DFGTD-B18 @unit-level
  Scenario: The mock dialect judgment presupposes the pattern it asks about
    Given the default gates of a project with tests
    When the mock gates are inspected
    Then mock-detect-covers-dialect and mock-stamped presuppose derived.mock_detect, and mock-typed presupposes derived.mock_contract

  @DFGTD-B19 @unit-level
  Scenario: rule-fulfilled is judged and marked to review
    Given the default gates of a project with specs and code
    When rule-fulfilled is inspected
    Then it measures by judgment and declares a review with its own question

  @DFGTD-B20 @unit-level
  Scenario: The catalog carries the checkers that measure the unit's kinds
    Given every artifact chosen
    When the default gates are listed
    Then evidence-fresh, feature-test-match, rule-implemented, scenario-identity, scenario-letter-declared, placeholder-filled and updated-at-atual are among them
    And scenario-type-aligned presupposes rule_types, and route-declared is scoped by the screen tag

  @DFGTD-B21 @unit-level
  Scenario: Init seeds what relates to the project and has its premise
    Given a project with spec, feature and test layers and no rule_types, that chose code with no code layer yet
    When the default gates are filtered for it
    Then feature-test-match and the security gates over code are seeded
    And scenario-type-aligned, presupposing rule_types, and route-declared, scoped by the screen tag, are not
