# language: en
# @anchors
#   code: PTHFA
#   ref: PTHSA
#   updated_at: 2026-11-06
#   layer: feature

@PTHSA
Feature: Paths — the aliases a path of anchors.yaml may open with, translated in one place

  @PTHSA-B01 @unit-level
  Scenario: A path that opens with an alias reads the path the alias stands for, in every field that holds a path
    Given aliases for the app's sources, the guides and an external repository, cited by a layer, a guide, a judgment gate's guide, a derived template, a gate exception and the changelog template
    When the configuration is loaded
    Then every one of those fields holds the path the alias stands for, the judgment guide its URL, and a folder named with @ that no alias names is kept

  @PTHSA-E01 @unit-level
  Scenario: An alias badly named, empty or opening with another, a path citing an undeclared alias, or a URL in a glob or a file written, refuses the configuration
    Given an alias with no @, an empty one, one that opens with another, a layer citing an alias nobody declared, a layer's glob and the changelog's path that are URLs
    When the configuration is loaded
    Then it does not load, naming the alias or the field
