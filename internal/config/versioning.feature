# language: en
# @anchors
#   code: VRSNF
#   ref: VRSNN
#   updated_at: 2026-11-18
#   layer: feature

@VRSNN
Feature: Versioning — the products a project versions, where each writes its version, and how each branch moves it

  @VRSNN-B01 @unit-level
  Scenario: The versioning block loads its products, and a configuration without it versions nothing
    Given an anchors.yaml with a product app of scheme semver, a source, a tag, a default strategy and a branch rule
    When the configuration is loaded, and another without versioning
    Then app holds each of them, and the other has no product

  @VRSNN-B02 @unit-level
  Scenario: Next is a strategy's name or a map, and manual when none is written
    Given the products next none, next with strategy increment and part patch, and one with no next
    When their strategies are resolved
    Then they are none, increment of the patch, and manual

  @VRSNN-B03 @unit-level
  Scenario: A branch takes the first rule it matches, in the order written, else the default
    Given a default semantic, and the branches hotfix/* increment and re:hotfix/.* none, in that order
    When hotfix/login and develop resolve their strategies
    Then hotfix/login takes increment naming hotfix/*, and develop takes semantic naming no pattern

  @VRSNN-B04 @unit-level
  Scenario: The tag is v{version} by default, and a tag is read back only with the product's text around it
    Given a product with no tag, and one with the tag app-v{version}
    When 1.2.0 is recorded, and the tags app-v1.2.0, api-v1.2.0 and app-v are read back
    Then the tags are v1.2.0 and app-v1.2.0, and only app-v1.2.0 reads back as 1.2.0

  @VRSNN-B05 @unit-level
  Scenario: A source is read at its path or in its group, and the sources must agree
    Given a JSON file with version 1.2.0 under expo, a YAML list whose first item holds 1.2.0, and a text file with VERSION=1.3.0
    When the sources are read and compared
    Then each version is read, and the comparison names the files that disagree with their versions

  @VRSNN-B06 @unit-level
  Scenario: A source is rewritten in place, keeping every byte around the version
    Given a JSON file with its keys, spacing and quotes, and a text file with VERSION=1.2.0 among other lines
    When each is rewritten with 1.3.0
    Then only the version changed

  @VRSNN-E01 @unit-level
  Scenario: A product whose scheme does not build refuses the configuration
    Given a product of scheme date, and an empty product
    When the configuration is loaded
    Then it does not load, naming versioning.app

  @VRSNN-E02 @unit-level
  Scenario: A source that names no place refuses the configuration
    Given sources with no file, with both path and regex, with neither, and with a regex of two groups
    When the configuration is loaded
    Then it does not load, naming the source

  @VRSNN-E03 @unit-level
  Scenario: A tag without its version, or with it twice, refuses the configuration
    Given the tags release and v{version}-{version}
    When the configuration is loaded
    Then it does not load, naming the tag

  @VRSNN-E04 @unit-level
  Scenario: A rule that cannot decide refuses the configuration
    Given the strategies latest, increment of build, semantic on number, a types level huge, manual with no source, part on semantic, a key bump in next, a key label in a branch rule, and the branch pattern re:(
    When the configuration is loaded
    Then it does not load, naming the field

  @VRSNN-E05 @unit-level
  Scenario: A version not found is an error naming the file
    Given a source file missing, a path no document holds, a regex that matches nothing, and a YAML file that is not one
    When the sources are read
    Then each fails, naming the file
