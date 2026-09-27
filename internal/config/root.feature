# language: en
# @anchors
#   ref: PRRPR
#   updated_at: 2026-09-26
#   layer: feature

@PRRPR
Feature: ProjectRootResolution — the project root a command works on

  @PRRPR-B01 @unit-level
  Scenario: The project root is the nearest directory above the start that holds the config
    Given a project "outer" with the subdirectory "pkg/deep", and a nested project "outer/inner" with the subdirectory "src"
    When the project root is resolved from "outer/pkg/deep", from "outer" and from "outer/inner/src"
    Then the roots are "outer", "outer" and "outer/inner"

  @PRRPR-B02 @unit-level
  Scenario: With no project above the start, the start comes back unchanged
    Given the directory "loose/dir", with no configuration file in it or above it
    When the project root is resolved from its absolute path, and from "." inside it
    Then the absolute path comes back as is, and "." comes back as "."

  @PRRPR-B03 @unit-level
  Scenario: With no root given, the root is found by walking up from the working directory
    Given the working directory "outer/pkg/deep" inside the project "outer"
    When the root is resolved from the option "." and from an empty option
    Then both give the absolute path of "outer"

  @PRRPR-I01 @unit-level
  Scenario: The resolved root is always absolute, even outside any project
    Given the working directory "loose/dir", with no project above it
    When the root is resolved from ".", from an empty option and from "sub"
    Then every answer is an absolute path

  @PRRPR-X01 @unit-level
  Scenario: An explicit root is made absolute and never walked above
    Given the working directory "outer", which holds a project, and its subdirectory "pkg/deep", which does not
    When the root is resolved from the explicit option "pkg/deep"
    Then it is the absolute path of "outer/pkg/deep", not "outer"
