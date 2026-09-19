# language: en
# @anchors
#   ref: CSDCN
#   updated_at: 2026-09-19
#   layer: feature

@CSDCN
Feature: ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those

  @CSDCN-B01 @unit-level
  Scenario: A status emitted and not declared is accused by number
    Given an output contract declaring 200, 404 and the 5xx range
    And a handler that also emits 401, 403 and 409 on its refusal branches
    When the gate confronts it
    Then it returns Fail naming 401, 403 and 409, because the client programmed from the
      table does not handle a refusal it was never told about

  @CSDCN-B02 @unit-level
  Scenario: A status declared and emitted by no path is accused as a phantom
    Given an output contract declaring 402 for exceeded quota
    And a handler whose quota branch answers 429 and never 402
    When the gate confronts it
    Then it returns Fail naming 402 as dead code in the client, which disappears with
      nobody noticing

  @CSDCN-B03 @unit-level
  Scenario: A faithful table passes
    Given an output contract whose concrete status codes are exactly the ones the handler emits
    When the gate confronts it
    Then it returns Pass, because a gate that only accuses is a noise generator

  @CSDCN-B04 @unit-level
  Scenario: The 500 of the top-level try/catch is not charged
    Given an output contract declaring 200 and the 5xx range
    And a handler whose only other status is the 500 of its catch block
    When the gate confronts it
    Then it returns Pass, because that 500 is infrastructure every handler carries, not a
      decision of this one

  @CSDCN-B05 @unit-level
  Scenario: The 5xx range covers, the 4xx range does not
    Given an output contract declaring 200, the 4xx range and the 5xx range
    And a handler that emits 503 and 403
    When the gate confronts it
    Then it returns Fail naming 403 and staying silent about 503, because a generic 4xx
      would hide exactly the access refusals this gate hunts

  @CSDCN-B06 @unit-level
  Scenario: A status that lives only in a comment is not emitted
    Given an output contract declaring 200 and 404
    And a handler whose 404 appears only in a comment describing the old behaviour, the
      live branch answering 403
    When the gate confronts it
    Then it returns Fail naming 404 as a phantom, because a comment is not behaviour

  @CSDCN-B07 @unit-level
  Scenario: Without the contract section there is nothing to confront
    Given a spec that catalogues effects and opens no output contract section
    When the gate confronts it
    Then it returns Skip, because charging the section's existence belongs to spec-complete

  @CSDCN-B08 @unit-level
  Scenario: Code that returns no status is skipped
    Given an output contract that declares void, the contract of a cron handler
    And code that persists items and returns no status
    When the gate confronts it
    Then it returns Skip, because there are no numbers on either side to compare

  @CSDCN-B09 @unit-level
  Scenario: A literal status passed to a local helper counts as emitted
    Given an output contract declaring 200 and 400
    And a handler that builds its 400 through a locally defined fail helper called with the literal
    When the gate confronts it
    Then it returns Pass, because a 400 is a 400 wherever the envelope is built

  @CSDCN-B10 @unit-level
  Scenario: With a dynamic status the phantom side goes quiet and the literals still count
    Given an output contract declaring 200 and 404
    And a handler with a helper that takes the status by parameter and one literal 403
    When the gate confronts it
    Then it returns Fail naming 403 and never naming 404, because a declared value may be
      emitted through a call textual reading cannot reach

  @CSDCN-B11 @unit-level
  Scenario: Without a declared dialect the verdict is Pending
    Given a project whose Structure declares no http_status lexicon
    And a handler that emits 200
    When the gate confronts it
    Then it returns Pending, because the meter does not fake conformity nor guess the stack

  @CSDCN-B12 @unit-level
  Scenario: An explicit opt-out of the http_status field is honoured
    Given a project whose Structure waives the http_status field of the dialect
    When the gate confronts a spec with an output contract
    Then it returns Skip, because the waiver is declared and localised, not a silent absence

  @CSDCN-I01 @unit-level
  Scenario: The lexicon comes from the project's dialect, not from the gate
    Given a Go handler that writes its refusal through the net/http writer
    And a Structure declaring the go dialect family
    When the gate confronts it
    Then it returns Fail naming 403, because embedding one stack's syntax would make the
      gate silent on every other one

  @CSDCN-I02 @unit-level
  Scenario: A dialect declared by hand teaches the gate its own lexicon
    Given a Structure declaring the http_status pattern directly, with no family
    And a Ruby handler that renders 422 through that pattern
    When the gate confronts it
    Then it returns Fail naming 422, because the agnosticism cannot stop at the built-in families

  @CSDCN-I03 @unit-level
  Scenario: A named constant is worth the number it means
    Given a handler whose refusal is written as the named forbidden constant and never as digits
    When the gate confronts it
    Then it returns Fail naming 403, because reading only digits would approve every
      handler written with constants

  @CSDCN-X01 @unit-level
  Scenario: The gate does not demand the generic ranges
    Given an output contract that declares no 5xx range at all
    And a handler whose only failure path is the 500 of its catch block
    When the gate confronts it
    Then it returns Pass, because treating a range as a status would mean guessing which
      numbers it covers

  @CSDCN-X02 @unit-level
  Scenario: The gate does not judge when each status is right
    Given an output contract declaring 403 for a branch a reviewer would call a 404
    And a handler that emits exactly that 403
    When the gate confronts it
    Then it returns Pass, because the ruler is the correspondence between two sets of
      numbers, not judgement about the design

  @CSDCN-X03 @unit-level
  Scenario: The gate does not charge the phantom side under a dynamic status
    Given a handler whose envelope helper receives the status by parameter
    And an output contract declaring a status no literal in the code shows
    When the gate confronts it
    Then it never names that status, because mass false positives are what makes a team
      turn the gate off
