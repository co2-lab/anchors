# language: en
# @anchors
#   ref: TLSTT
#   updated_at: 2026-09-26
#   layer: feature

@TLSTT
Feature: TelemetrySetup — every command starts telemetry the same way: the opt-outs first, then the notice, then the emitter

  @TLSTT-B01 @unit-level
  Scenario: The project opt-out is read from the project root
    Given a project whose configuration declares "telemetry: off"
    When a command starts from the project root, from its subdirectory "sub", and from elsewhere with --root pointing at it
    Then no run shows the notice or builds the emitter

  @TLSTT-B02 @unit-level
  Scenario: The environment opt-out builds nothing
    Given ANCHORS_TELEMETRY is "off" and a project with no marks
    When a command starts
    Then no notice is printed, no emitter is built and the project has no ".anchors" directory

  @TLSTT-B03 @unit-level
  Scenario: The notice is shown once, then telemetry runs quietly
    Given telemetry is not turned off and the project has never shown the notice
    When a command starts twice
    Then the first start prints a notice containing "ANCHORS_TELEMETRY=off", builds the emitter and writes ".anchors/telemetry-noticed"
    And the second start prints nothing

  @TLSTT-B04 @unit-level
  Scenario: The authentication header comes only from the environment
    Given ANCHORS_TELEMETRY_KEY is empty, and then "secret"
    When the emitter headers are built
    Then there are no headers, and then only "x-honeycomb-team" with "secret"

  @TLSTT-B05 @unit-level
  Scenario: The project root is the explicit root or the nearest configured directory above
    Given an explicit --root, and then no --root with the working directory at "pkg/a" under a configured project
    When the project root is resolved
    Then it is the explicit directory, and then the configured project directory

  @TLSTT-B06 @unit-level
  Scenario: Waiting at exit with no emitter returns at once
    Given telemetry turned off, so no emitter was built
    When the command waits for events in flight at exit
    Then it returns without error

  @TLSTT-B07 @unit-level
  Scenario: A declared opt-out holds when the configuration does not load
    Given a project whose configuration declares "telemetry: off" next to a key Anchors does not know
    When a command starts with --root pointing at it
    Then no notice is printed and no emitter is built

  @TLSTT-B08 @unit-level
  Scenario: Outside a project there is no notice, no emitter and no mark
    Given a working directory with no project configuration in it or above it
    When the project root is resolved, and a command starts with no --root and with --root pointing at that directory
    Then the root is empty, no notice is printed, no emitter is built, and the directory has no ".anchors"

  @TLSTT-I01 @unit-level
  Scenario: No emitter without the notice
    Given each opt-out, and then none
    When a command starts
    Then with an opt-out there is neither notice, nor emitter, nor mark, and without one all three exist

  @TLSTT-X01 @unit-level
  Scenario: Starting telemetry sends nothing
    Given the telemetry endpoint is a local server that counts the requests it receives
    When a command starts with telemetry on and then waits at exit
    Then the emitter is built and the server has received no request
