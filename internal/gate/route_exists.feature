# language: en
# @anchors
#   ref: RTEXR
#   updated_at: 2026-09-19
#   layer: feature

@RTEXR
Feature: RouteExists — declared route in specification must exist in application route registry

  @RTEXR-B01 @unit-level
  Scenario: Non-spec artifacts skip confrontation
    Given an artifact node whose kind is not spec
    When the gate confronts it
    Then it returns Skip, because route existence is evaluated on specifications

  @RTEXR-B02 @unit-level
  Scenario: Specifications without declared routes skip confrontation
    Given a specification containing no declared route
    When the gate confronts it
    Then it returns Skip, leaving route declaration requirements to route-declared

  @RTEXR-B03 @unit-level
  Scenario: Specifications with declared routes return Pending when route registry is unconfigured
    Given a specification declaring a route
    And a project configuration with no route registry globs
    When the gate confronts it
    Then it returns Pending, because route existence cannot be verified without registry locations

  @RTEXR-B04 @unit-level
  Scenario: Invalid route registry glob patterns return Pending
    Given a specification declaring a route
    And a project configuration with a malformed route registry glob pattern
    When the gate confronts it
    Then it returns Pending, reporting the glob matching error

  @RTEXR-B05 @unit-level
  Scenario: Route registry files containing zero registered routes return Pending
    Given a specification declaring a route
    And route registry files that define no recognized routes
    When the gate confronts it
    Then it returns Pending, because no routes could be discovered

  @RTEXR-B06 @unit-level
  Scenario: Declared screen route matching a component navigation prop passes
    Given a specification declaring a screen route
    And a route registry file defining that screen in a navigation component prop
    When the gate confronts it
    Then it returns Pass, confirming the screen exists in the application registry

  @RTEXR-B07 @unit-level
  Scenario: Declared screen route matching a navigation stack parameter type entry passes
    Given a specification declaring a screen route
    And a route registry file defining that screen in a navigation stack parameter type
    When the gate confronts it
    Then it returns Pass, confirming the route exists in the navigation stack types

  @RTEXR-B08 @unit-level
  Scenario: Declared backend route matching an HTTP resource registration passes
    Given a specification declaring an HTTP backend route
    And a route registry file registering that route with addResource
    When the gate confronts it
    Then it returns Pass, confirming the HTTP resource is registered

  @RTEXR-B09 @unit-level
  Scenario: HTTP method verb prefixes are stripped when evaluating declared routes
    Given a specification declaring a route with an HTTP method prefix
    And a route registry file containing the path without the HTTP method
    When the gate confronts it
    Then it returns Pass, stripping the method verb to match the resource path

  @RTEXR-B10 @unit-level
  Scenario: Route paths match regardless of leading slash differences between specification and code
    Given a specification declaring a route path
    And a route registry registering the route with differing leading slash notation
    When the gate confronts it
    Then it returns Pass, normalizing slashes during matching

  @RTEXR-B11 @unit-level
  Scenario: Custom route pattern regex matches custom route registration patterns
    Given a project configuration defining a custom route pattern regex
    And route registry files matching that custom pattern
    When the gate confronts a specification declaring that route
    Then it returns Pass, matching routes extracted via the custom regex

  @RTEXR-B12 @unit-level
  Scenario: Custom route pattern regex that is malformed or lacks capture groups returns Pending
    Given a project configuration with an invalid custom route pattern regex
    When the gate confronts a specification declaring a route
    Then it returns Pending, reporting the invalid pattern configuration

  @RTEXR-B13 @unit-level
  Scenario: Declared routes missing from all registered route definitions fail
    Given a specification declaring a route absent from all registry files
    When the gate confronts it
    Then it returns Fail, citing the missing route, registered count, and searched globs

  @RTEXR-I01 @unit-level
  Scenario: Route presence is validated only for specifications
    Given code, feature, or test artifacts confronted by the gate
    When the gate confronts them
    Then it returns Skip, keeping route verification scoped to specifications

  @RTEXR-I02 @unit-level
  Scenario: The gate never approves route existence without inspecting registry files
    Given a specification declaring a route in a project with missing or empty registries
    When the gate confronts it
    Then it returns Pending, preventing false approvals of uninspected routes

  @RTEXR-I03 @unit-level
  Scenario: Route matching is slash-normalized between specification and code
    Given specifications declaring routes with or without leading slashes
    When the gate confronts them against registered routes
    Then it accepts either format symmetrically

  @RTEXR-I04 @unit-level
  Scenario: Missing routes produce a blocking Fail verdict
    Given a specification promising a route that does not exist in code
    When the gate confronts it
    Then it returns Fail, preventing unreachable destinations from reaching production

  @RTEXR-X01 @unit-level
  Scenario: The gate does not require every specification to declare a route
    Given a specification that defines no route
    When the gate confronts it
    Then it returns Skip, leaving route declaration obligations to route-declared

  @RTEXR-X02 @unit-level
  Scenario: Route parameter schemas and payload contracts are not evaluated by this gate
    Given a specification declaring a route that exists in the registry
    When the gate confronts it
    Then it evaluates existence only, without validating payload or parameter schemas

  @RTEXR-X03 @unit-level
  Scenario: Route access permissions and authentication middlewares are outside evaluation scope
    Given a specification declaring an authenticated route
    When the gate confronts it
    Then it validates registry existence without checking route guard middlewares
