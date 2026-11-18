# language: en
# @anchors
#   code: SCFTB
#   ref: LGSCL
#   updated_at: 2026-10-03
#   layer: feature

@LGSCL
Feature: LogScan — finding the occurrences of declared failures in the project's logs

  @LGSCL-B01 @unit-level
  Scenario: Without a declared log path nothing is scanned
    Given a project with a log file and no declared log path
    When the logs are scanned
    Then there is no result and no error

  @LGSCL-B02 @unit-level
  Scenario: The format does not matter
    Given the failure "CRED-E01" logged once in JSON, once in plain text and once in syslog, in three files
    When the logs are scanned
    Then "CRED-E01" has 3 occurrences, from 3 files and 3 lines

  @LGSCL-B03 @unit-level
  Scenario: Look-alikes of a code are not captured
    Given a log with a build id, a cache key, a SKU and a trace id shaped like codes, and one "#[CRED-E01]"
    When the logs are scanned with "CRED-E01" declared
    Then only "CRED-E01" is an occurrence and nothing is reported unknown

  @LGSCL-B04 @unit-level
  Scenario: An undeclared failure is reported apart
    Given a log with "ERROR #[ORPHA-E99]" and "INFO #[ORPHA-E98]", and a bare log with "ERROR ORPHA-E77" and "INFO ORPHA-E76"
    When each is scanned with only "CRED-E01" declared, the bare one without delimiters
    Then the delimited ones are both reported, and of the bare ones only "ORPHA-E77"

  @LGSCL-B05 @unit-level
  Scenario: The delimiter removes the ambiguity
    Given a log citing "CRED-E01" in prose, in an array index and in JSON tags, plus one "#[CRED-E01]"
    When it is scanned with the default delimiters, and with an empty pair
    Then the first counts 1 occurrence and the second more than 1

  @LGSCL-B06 @unit-level
  Scenario: The alias only serves lines that carry no code
    Given the alias "INSUFFICIENT_BALANCE" for "CRED-E01" and a log with that legacy line and a "#[CRED-E02]" line
    When the logs are scanned
    Then "CRED-E01" and "CRED-E02" have one occurrence each

  @LGSCL-B07 @unit-level
  Scenario: Each occurrence records its time window
    Given "CRED-E01" logged on 2026-09-01, 2026-09-05 and 2026-09-19
    When the logs are scanned with a timestamp pattern
    Then its first moment is on 2026-09-01 and its last on 2026-09-19

  @LGSCL-B08 @unit-level
  Scenario: Occurrences are sorted by rule
    Given a log with "#[CRED-E02]" before "#[CRED-E01]", both declared
    When the logs are scanned
    Then the occurrences are "CRED-E01" then "CRED-E02"

  @LGSCL-B09 @unit-level
  Scenario: A spec's failure rules are listed in every form
    Given a spec with "### ABCDE-E01" as a heading, "| `ABCDE-E02` |" as a table row, "- **ABCDE-E03**" as a bullet, and "ABCDE-B01"
    When its failure rules are listed
    Then the listed failure rules are "ABCDE-E01", "ABCDE-E02" and "ABCDE-E03", and the behaviour rule is left out

  @LGSCL-E01 @unit-level
  Scenario: An invalid alias or timestamp pattern fails the scan
    Given the alias pattern "(" and, separately, the timestamp pattern "("
    When the logs are scanned
    Then each scan fails with no result
