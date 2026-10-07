# language: en
# @anchors
#   code: NCFNV
#   ref: NCGNV
#   layer: feature

@NCGNV
Feature: NavigationChain — every navigation flagged with the screen it leads to, and the screens' tables confronted with it

  @NCGNV-B01 @unit-level
  Scenario: nav-annotated names each navigation call with no flag, and each flag naming another screen than its route
    Given an app whose home navigates unflagged, whose detail flags its call right and waives a back navigation, and whose edit flags its call with another screen
    When nav-annotated runs on each
    Then the home's call and the edit's wrong flag are named, and the detail passes

  @NCGNV-B02 @unit-level
  Scenario: nav-matches-spec names each Out row no flag answers, and each flag the Out table does not declare
    Given the same app, whose edit's Out leads to the detail while its code is flagged to the home
    When nav-matches-spec runs on the detail and the edit
    Then the detail passes, and the edit names both the unanswered row and the undeclared flag

  @NCGNV-B03 @unit-level
  Scenario: nav-symmetric names each Out with no matching In, and each In with no matching Out
    Given the same app, whose detail's In forgets the edit that leads to it
    When nav-symmetric runs on the edit and the home
    Then the edit names the detail, and the home passes

  @NCGNV-B04 @unit-level
  Scenario: nav-reachable fails a screen no entry route reaches, and is pending with no entry declared
    Given the same app with Home as its entry, and a lonely screen nothing leads to
    When nav-reachable runs on the edit and the lonely screen, then with no entry declared
    Then the edit passes, the lonely screen fails, and with no entry it is pending

  @NCGNV-B05 @unit-level
  Scenario: The fixer flags a navigation call whose route names one screen, and leaves a back navigation to the author
    Given the same app
    When check --fix runs the navigation fixer
    Then the home's call is flagged with the detail's code, and the detail's waived back navigation stays as it is
