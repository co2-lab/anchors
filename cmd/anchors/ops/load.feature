# language: en
# @anchors
#   code: LDCMB
#   ref: LDCMS
#   updated_at: 2026-11-19
#   layer: feature

@LDCMS
Feature: LoadCommand — what the machine lends its runs, and who holds it

  @LDCMS-B01 @unit-level
  Scenario: anchors load shows the capacity, what is in use, the runs holding the machine and the runs waiting, and which waited past their limit
    Given 5 CPU and 8 GB, app's test unit holding 4 CPU and 6 GB for 90 seconds, and web's mutation unit waiting 10 minutes past a limit of 5
    When anchors load runs, and on an idle machine
    Then it shows each of them with the waiter past its limit, and dashes when idle
