# language: en
# @anchors
#   ref: TLEMT
#   updated_at: 2026-09-26
#   layer: feature

@TLEMT
Feature: TelemetryEmitter — decision events leave as OTLP logs, never block the work, and never carry who uses the product

  @TLEMT-B01 @unit-level
  Scenario: A disabled configuration yields no emitter that still accepts calls
    Given a configuration with telemetry disabled
    When the emitter is created and an event is emitted and flushed through it
    Then there is no emitter and neither call crashes

  @TLEMT-B02 @unit-level
  Scenario: A blank endpoint means Honeycomb's logs endpoint
    Given an enabled configuration with no endpoint
    When the emitter is created
    Then it sends to "https://api.honeycomb.io/v1/logs"

  @TLEMT-B03 @unit-level
  Scenario: An event is sent as one log record whose body is the event name
    Given an enabled emitter pointing at a local collector
    When the event "claim.served" with candidates 14 is emitted at 2026-09-14 18:30 UTC and flushed
    Then the collector receives one resource with one log record
    And the record's body is "claim.served", its time is the instant in Unix nanoseconds and it carries the attribute candidates 14

  @TLEMT-B04 @unit-level
  Scenario: Emitting to a dead collector neither blocks nor fails
    Given an enabled emitter pointing at an address nobody listens on
    When an event is emitted
    Then the call returns within half a second without an error

  @TLEMT-B05 @unit-level
  Scenario: Flushing waits for the send in flight
    Given a collector that takes 50 milliseconds to answer
    When an event is emitted and the emitter is flushed
    Then the collector has received the event when the flush returns

  @TLEMT-B06 @unit-level
  Scenario: Flushing has its own deadline independent of the HTTP client
    Given a send in flight that never finishes and never goes through the HTTP client
    When the emitter is flushed
    Then the flush returns after about two seconds

  @TLEMT-B07 @unit-level
  Scenario: Each request is a JSON POST carrying the configured headers
    Given an enabled emitter whose configuration has the header "x-honeycomb-team" set to "k"
    When an event is emitted and flushed
    Then the collector sees a POST with content type "application/json" and that header

  @TLEMT-B08 @unit-level
  Scenario: With NoCodes a unit or rule code never leaves in an attribute
    Given an event with the attributes unit "RLSGR", rule "RLSGR-B01", gate "triad-complete", card_state "in-progress" and n 3
    When it is built by an emitter with NoCodes, and by one without
    Then the first carries only gate, card_state and n, and the second carries all five

  @TLEMT-I01 @unit-level
  Scenario: A value that is not vocabulary becomes its type description
    Given an attribute value that is a map holding a file path
    When it is converted for the protocol
    Then the value sent is "(map[string]string)" and the path does not appear
    And the integer 14 and the string "pr-aberto" pass typed

  @TLEMT-X01 @unit-level
  Scenario: The resource carries only the product name and version
    Given an enabled emitter for version "0.1.99"
    When an event body is built
    Then the resource has exactly the attributes "service.name" = "anchors" and "service.version" = "0.1.99"

  @TLEMT-E02 @unit-level
  Scenario: An error status from the collector is discarded
    Given a collector that answers 500
    When an event is emitted and flushed
    Then the collector received the event once and nothing fails
