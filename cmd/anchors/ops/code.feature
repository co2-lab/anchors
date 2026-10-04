# language: en
# @anchors
#   code: CDFTC
#   ref: CDCMC
#   updated_at: 2026-10-03
#   layer: feature

@CDCMC
Feature: CodeCommand — a new unit gets an identity code that no other unit in the map already owns, and the codes in use are listed from the map

  @CDCMC-B01 @unit-level
  Scenario: A free canonical code is the suggestion
    Given a map with no code taken
    When code runs for Spacer
    Then it suggests the canonical code of Spacer

  @CDCMC-B02 @unit-level
  Scenario: A taken canonical is adjusted to a free code naming its owner
    Given a map where ui/Spacer.spec.md holds the canonical code of Spacer
    When code runs for Spacer
    Then it suggests another code
    And it says the canonical already belongs to ui

  @CDCMC-B03 @unit-level
  Scenario: A path in a layer with a code prefix gets the module prefix
    Given a layer auth with code prefix AU and the file src/auth/Session.ts
    When code runs for src/auth/Session.ts
    Then it suggests the code with prefix AU for Session and mentions the module prefix 'AU'

  @CDCMC-B04 @unit-level
  Scenario: A generic basename takes its identity from the parent directory
    Given the path src/billing/index.ts
    When code runs for it
    Then it says the identity came from the parent dir and does not suggest the code of "index"

  @CDCMC-B05 @unit-level
  Scenario: Check answers whether a code is free, ignoring case, and names the owners
    Given a map where src/auth/Login.spec.md holds LOGNS
    When code runs with --check logns, then with --check FREEX
    Then the first says "LOGNS is already used by: src/auth" and fails with the collision error
    And the second says FREEX is free

  @CDCMC-B06 @unit-level
  Scenario: The list prints one sorted code per line with its folder, the summary kept off stdout
    Given a map where LOGI is on a spec and a tsx of the same folder, plus SGSB and BTTN
    When code list runs
    Then stdout holds "LOGI<TAB>apps/mobile/src/features/auth" once, BTTN before LOGI
    And stdout holds no summary line

  @CDCMC-B07 @unit-level
  Scenario: The list filters by path prefix and names a filter that matched nothing
    Given the same map
    When code list runs with --in packages/backend, then with --in apps/web
    Then the first lists only SGSB
    And the second says no code in use under "apps/web"

  @CDCMC-B08 @unit-level
  Scenario: The JSON list carries each code's folder, file, kind, title and work order fields
    Given a plan PLNBB titled "Plano 0002 — Build" with needs, parent and revises
    When code list runs with --json
    Then the entry has code PLNBB, onde plans, the file, kind plan, titulo Build, needs PLNAA, parent PRDCT, revises PLNZZ
    And an empty map gives []

  @CDCMC-B09 @unit-level
  Scenario: The title drops the text before the dash
    Given markdown files titled "# Login", "# Plano 0001 — Fundação" and one with no heading
    When their titles are read
    Then they are "Login", "Fundação" and empty
    And and a .go file has no title

  @CDCMC-B10 @unit-level
  Scenario: The length check accuses only declared codes and proposes the canonical code
    Given a map with declared WLTX (4 letters), declared BDGTS and cited FXTR
    When code list runs with --check
    Then it lists "WLTX → " the canonical code of Wallet and fails
    And FXTR is not accused and "1 conforming" is reported
    And and filtered to a conforming folder it passes counting the cited code

  @CDCMC-B11 @unit-level
  Scenario: The unit name drops the artifact suffixes
    Given the paths Login.spec.md, Login.feature, Login.test.ts and NewLogin.tsx
    When their unit names are taken
    Then they are Login, Login, Login and NewLogin

  @CDCMC-B12 @unit-level
  Scenario: An empty map says no node has an identity
    Given a map with no nodes
    When code list runs
    Then it says the map has no node with identity

  @CDCMC-X01 @unit-level
  Scenario: The codes come from the map's identity field
    Given a map whose nodes point at files that do not exist on disk
    When code list runs
    Then it lists every code of the map

  @CDCMC-E01 @unit-level
  Scenario: Without a name or a map the command fails and says what to do
    Given a project with a map, and one without
    When code runs with no name, then for X where there is no map
    Then the first asks for the unit name
    And the second asks to run anchors map build

  @CDCMC-E02 @unit-level
  Scenario: The list refuses a project without config or without map
    Given a directory with no anchors.yaml, then one with a config and no map
    When code list runs
    Then it fails with "load anchors.yaml", then with "read the map"

  @CDCMC-B13 @unit-level
  Scenario: A second check in the same process judges only its own map
    Given a first map that declares WLTX and a second map that only cites WLTX
    When code list --check runs on the first and then on the second
    Then the second passes, counting WLTX as only cited
    And the JSON of the second names app/fixture_test.go as the file

  @CDCMC-B14 @unit-level
  Scenario: The length check points each divergence to anchors recode, and there is no --fix
    Given a map that declares WLTX for app/Wallet.spec.md
    When code list --check runs
    Then the output has "anchors recode WLTX" followed by the canonical code, and never "--fix"
    And the command declares no --fix flag

  @CDCMC-B15 @unit-level
  Scenario: A name shaped like a code also gets that code's status
    Given a map where one code is taken
    When a code is generated for that taken code's shape, for a free code's shape, and for a unit name
    Then the first two also say the code is taken or free and point at --check, and the unit name gets no note
