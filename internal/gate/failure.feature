# language: en
# @anchors
#   code: FLFTB
#   ref: FLRAI
#   updated_at: 2026-10-07
#   layer: feature

@FLRAI
Feature: Failure — the failure a spec declares must be handled, recorded, and every handling declared

  @FLRAI-B01 @unit-level
  Scenario: Every failure gate skips an artifact that is not a spec
    Given a code, a test and a feature node
    When failure-handled, failure-logged and failure-declared confront each of them
    Then every verdict is Skip

  @FLRAI-B02 @unit-level
  Scenario: A failure rule is read in the heading, table row and bullet forms, not in prose
    Given a spec declaring "CRED-E01" in a heading, "CRED-E02" in a table row and "CRED-E03" in a bold bullet
    And governed code with no handling path
    When failure-handled confronts the spec
    Then it returns Fail naming "CRED-E01, CRED-E02, CRED-E03"
    And a spec citing "CRED-E01" only in a prose sentence is skipped

  @FLRAI-B03 @unit-level
  Scenario: A spec that declares no failure is skipped by failure-handled and failure-logged
    Given a spec whose only rule is the behaviour "CRED-B01"
    When failure-handled and failure-logged confront it
    Then both return Skip

  @FLRAI-B04 @unit-level
  Scenario: Without the dialect patterns the failure gates are Pending
    Given a project with no dialect declared
    When the three failure gates confront a spec declaring "CRED-E01" over code that checks an error
    Then every verdict is Pending
    And failure-logged is Pending too when the dialect declares handling patterns but no recording patterns

  @FLRAI-B05 @unit-level
  Scenario: Without governed code that can be read the failure gates are Pending
    Given a spec declaring "CRED-E01"
    And either no map, a specifies edge to a file that is not on disk, or only an edge of another type
    When the three failure gates confront the spec
    Then every verdict is Pending

  @FLRAI-B06 @unit-level
  Scenario: A handling written only in a comment line does not count
    Given a spec declaring "CRED-E01"
    And governed code whose only error check sits in a comment line
    When failure-handled confronts the spec
    Then it returns Fail

  @FLRAI-B07 @unit-level
  Scenario: A declared failure with no handling in the governed code fails, naming it
    Given a spec declaring "CRED-E01"
    And governed code that only returns a number
    When failure-handled confronts the spec
    Then it returns Fail naming "CRED-E01"

  @FLRAI-B08 @unit-level
  Scenario: Any handling path in the governed code passes failure-handled
    Given a spec declaring "CRED-E01" and "CRED-E02"
    And governed code that checks an error once
    When failure-handled confronts the spec
    Then it returns Pass

  @FLRAI-B09 @unit-level
  Scenario: A failure marked resilient with a reason is not charged
    Given a spec whose only failure "CRED-E01" is marked "@resilient: the partner returns null during migration"
    When failure-handled confronts it over code with no handling
    Then it returns Pass
    And failure-logged over code that handles silently also returns Pass

  @FLRAI-B10 @unit-level
  Scenario: A bare resilient marker exempts nothing
    Given a spec whose failure "CRED-E01" carries "@resilient" with no reason
    When failure-handled confronts it over code with no handling
    Then it returns Fail

  @FLRAI-B11 @unit-level
  Scenario: A handling that records nothing fails failure-logged, naming the failure
    Given a spec declaring "CRED-E01"
    And governed code that checks an error and returns without logging
    When failure-logged confronts the spec
    Then it returns Fail naming "CRED-E01"

  @FLRAI-B12 @unit-level
  Scenario: A handling that records the occurrence passes failure-logged
    Given a spec declaring "CRED-E01"
    And governed code that checks an error and logs it with "log.Error"
    When failure-logged confronts the spec
    Then it returns Pass

  @FLRAI-B13 @unit-level
  Scenario: Handling in the code with no failure declared in the spec fails
    Given a spec with no failure rule and no failure section
    And governed code that checks an error
    When failure-declared confronts the spec
    Then it returns Fail

  @FLRAI-B14 @unit-level
  Scenario: Code with no handling is skipped by failure-declared
    Given a spec with no failure rule
    And governed code that only returns a number
    When failure-declared confronts the spec
    Then it returns Skip

  @FLRAI-B15 @unit-level
  Scenario: Handling in the code and a failure declared in the spec passes
    Given a spec declaring "CRED-E01"
    And governed code that checks an error
    When failure-declared confronts the spec
    Then it returns Pass

  @FLRAI-B16 @unit-level
  Scenario: An Errors section closed with none and a reason passes failure-declared
    Given governed code whose only handling match is a lazy map initialisation
    When failure-declared confronts specs whose Errors section opens with "none — the only nil check is a lazy map init" or "nenhuma: o nil check inicializa o mapa"
    Then both return Pass
    And a bare "none", a "none" outside the section, and a section opening with prose each return Fail

  @FLRAI-B17 @unit-level
  Scenario: The conclusions read the resilient and observing reasons of each failure whole
    Given a spec with "CRED-E01" unmarked, "CRED-E02" marked "@resilient: the partner restarts at 3am daily; the retry covers it" and "CRED-E03" marked "@observing: ruled out partner retry and network latency; only on migrated accounts"
    When the conclusions are read
    Then "CRED-E01" carries no conclusion
    And "CRED-E02" carries a resilient reason containing "retry covers it"
    And "CRED-E03" carries an observing reason containing "migrated accounts"

  @FLRAI-B18 @unit-level
  Scenario: A conclusion reason ends at its table cell
    Given a table row "| `CRED-E01` | cond | @resilient: the real reason | another column |"
    When the conclusions are read
    Then the resilient reason contains "the real reason" and not "another column"

  @FLRAI-I01 @unit-level
  Scenario: failure-handled and failure-logged never both fail the same spec
    Given a spec declaring "CRED-E01"
    When both gates confront it over code with no handling, and over code that handles without recording
    Then the first gives Fail and Skip, and the second gives Pass and Fail

  @FLRAI-X01 @unit-level
  Scenario: One handling path answers for every declared failure
    Given a spec declaring "CRED-E01" and "CRED-E02"
    And governed code with a single error check
    When failure-handled confronts the spec
    Then it returns Pass for both failures together

  @FLRAI-B19 @unit-level
  Scenario: A failure rule is read at the code length the project declares
    Given a project that declares code length 7
    When a spec with the table row "CREDITS-E01" marked resilient and the bullet "CREDITS-E02" is read
    Then both are declared failures, "CREDITS-E01" is resilient, and its resilient reason is read

  @FLRAI-B20 @unit-level
  Scenario: A unit with a fallible source and no declared failure fails, naming the source
    Given a unit whose code calls a fallible query on line 5, and mentions it in a comment on line 3
    When failure-declared runs on its spec with no failure, with a failure that does not name the query, with one that names it, with a waiver, and with a section closed as none
    Then the first two fail naming line 5 and not line 3, and the other three pass

  @FLRAI-B21 @unit-level
  Scenario: A dependency on a file of a fallible layer is a fallible source
    Given a spec depending on a hook of a layer marked fallible and on a utility of another layer
    When failure-declared runs on the spec with no failure, then with a failure whose uses name the DEP
    Then the first fails naming the hook's dependency and not the utility's, and the second passes
