# language: en
# @anchors
#   code: IMFTM
#   ref: MPCTI
#   updated_at: 2026-10-03
#   layer: feature

@MPCTI
Feature: Impact — what a change to one file reaches, in both directions of the map

  @MPCTI-B01 @unit-level
  Scenario: A change to a spec propagates down to the code and the test it specifies
    Given a project whose login spec specifies a login code file and a login test
    When the impact of the login spec is asked
    Then the propagate direction lists the login code file and the login test

  @MPCTI-B02 @unit-level
  Scenario: A change to the code is validated up against its spec and its guide
    Given a project whose login code file is specified by the login spec and governed by the code guide
    When the impact of the login code file is asked
    Then the validate direction lists the login spec and the code guide

  @MPCTI-B03 @unit-level
  Scenario: An empty direction is said explicitly instead of printed as an empty list
    Given a login spec that nobody governs and a login code file nothing depends on
    When the impact of each is asked
    Then the spec reads that it is not governed by anyone and the code reads that no child depends on it

  @MPCTI-B04 @unit-level
  Scenario: A root-relative, native-separator or absolute argument resolves to the same node id
    Given a project with the file packages/backend/services/auth.spec.md
    When the argument is given root-relative with slashes, with the native separator, and absolute
    Then each one resolves to the id packages/backend/services/auth.spec.md

  @MPCTI-B05 @unit-level
  Scenario: A relative argument that does not exist under the root is resolved from the working directory
    Given a project with the file packages/backend/x.spec.md and the command running from its packages folder
    When the argument backend/x.spec.md is resolved
    Then it resolves to the id packages/backend/x.spec.md

  @MPCTI-I01 @unit-level
  Scenario: The resolved node id always uses forward slashes
    Given a native-separator path that names nothing on disk
    When the argument is resolved
    Then the result carries no backslash

  @MPCTI-X01 @unit-level
  Scenario: The impact query leaves the map unchanged
    Given a project with a built map
    When the impact of the login spec is asked
    Then the map is byte for byte what it was

  @MPCTI-E01 @unit-level
  Scenario: A file that is not a node of the map is refused
    Given a project whose map has no node for src/ghost.ts
    When the impact of src/ghost.ts is asked
    Then it fails saying the file is not in the map

  @MPCTI-E02 @unit-level
  Scenario: Without a map the impact query fails and asks for the map build
    Given a directory with no map
    When the impact of a file is asked there
    Then it fails with a hint to run the map build
