# language: en
# @anchors
#   code: MNFMP
#   ref: MNCMP
#   layer: feature

@MNCMP
Feature: MapNav — the app's navigation, screen by screen

  @MNCMP-B01 @unit-level
  Scenario: The navigation is printed screen by screen, each with where it leads, sorted
    Given a home leading to a detail, which leads to an edit and back home
    When the navigation is printed
    Then each screen has its line with its destinations, sorted

  @MNCMP-B02 @unit-level
  Scenario: One screen shows where it comes from and where it leads; a name that is no screen is refused
    Given the same navigation
    When the detail is named, then a name nobody owns
    Then the detail shows the home as origin and the edit and home as destinations, and the other name is refused
