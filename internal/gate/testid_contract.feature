# language: en
# @anchors
#   ref: TICTS
#   updated_at: 2026-09-26
#   layer: feature

@TICTS
Feature: TestIDContract — a test handle is one contract with four ends: the code exposes it, the spec declares it, a consumer queries it

  @TICTS-B01 @unit-level
  Scenario: The gate skips what is not a spec, and is Pending without a map
    Given a code node, and in turn the spec x.spec.md with no map
    When testid-consistent confronts each
    Then the first returns Skip and the second Pending

  @TICTS-B02 @unit-level
  Scenario: Without a declared handle attribute the gate skips
    Given a spec whose code exposes ":abcd-screen", in a project that declares no test handle attribute
    When testid-consistent confronts the spec
    Then it returns Skip

  @TICTS-B03 @unit-level
  Scenario: A spec with no readable linked code skips
    Given a spec with no specifies edge, and in turn a spec whose specifies edge points to a missing file
    When testid-consistent confronts each
    Then both return Skip, saying the spec has no linked code

  @TICTS-B04 @unit-level
  Scenario: A unit that exposes no handle and declares none skips
    Given a spec with no inventory whose code is "<View />"
    When testid-consistent confronts the spec
    Then it returns Skip

  @TICTS-B05 @unit-level
  Scenario: A handle exposed, declared and queried passes
    Given code exposing ":abcd-screen", a spec declaring it, and the linked test calling getByTestId(':abcd-screen')
    When testid-consistent confronts the spec
    Then it returns Pass

  @TICTS-B06 @unit-level
  Scenario: A handle the code exposes and the spec does not declare fails
    Given code exposing ":abcd-screen" and ":abcd-hidden", and a spec declaring only ":abcd-screen"
    When testid-consistent confronts the spec
    Then it returns Fail naming "abcd-hidden"

  @TICTS-B07 @unit-level
  Scenario: A handle the spec declares and the code does not expose fails
    Given code exposing ":abcd-screen", and a spec declaring ":abcd-screen" and ":abcd-ghost"
    When testid-consistent confronts the spec
    Then it returns Fail naming "abcd-ghost"

  @TICTS-B08 @unit-level
  Scenario: A handle no consumer queries fails when a consumer surface exists
    Given code exposing and a spec declaring ":abcd-screen", and a linked test that only renders the screen
    When testid-consistent confronts the spec
    Then it returns Fail naming "abcd-screen"

  @TICTS-B09 @unit-level
  Scenario: With no consumer surface the queried end is not charged
    Given code exposing ":abcd-screen" and ":abcd-hidden", a spec declaring ":abcd-screen", and no test, flow or neighbour
    When testid-consistent confronts the spec
    Then the line of ":abcd-hidden" shows "queried —", and a fully declared inventory passes

  @TICTS-B10 @unit-level
  Scenario: The handle attribute is the one the project declares
    Given a project declaring "data-testid", code with data-testid "abcd-root" and "abcd-item", and a spec declaring "abcd-root"
    When testid-consistent confronts the spec
    Then it returns Fail naming "abcd-item"

  @TICTS-B11 @unit-level
  Scenario: A literal handle counts, also through a derived prop or an object key
    Given code with backTestID=":abcd-back" and an object key confirmTestID: ':abcd-ok', and a spec declaring both
    When testid-consistent confronts the spec
    Then it returns Pass

  @TICTS-B12 @unit-level
  Scenario: A template handle, or a prefix prop, exposes the prefix as a wildcard
    Given code with testID={`:abcd-item-${id}`}, and in turn code with testIDPrefix=":otp-input"
    When the exposed handles are read
    Then they are "abcd-item-*" and "otp-input-*", covering a spec that declares "abcd-item-*" or "otp-input-0" and "otp-input-5"

  @TICTS-B13 @unit-level
  Scenario: Only the branches of a conditional handle are handles, never its condition
    Given code with testID={k === 'push' ? ':abcd-toggle' : undefined}, and a spec declaring "abcd-toggle"
    When testid-consistent confronts the spec
    Then it returns Pass, and "push" is not an exposed handle

  @TICTS-B14 @unit-level
  Scenario: The inventory is read only inside the test surface section
    Given a spec whose "## Test IDs (Maestro)" section declares "abcd-screen", and in turn one that mentions it only under "## Notes"
    When the inventory is read
    Then the first declares it, the second declares nothing, and an id after the next heading is not declared

  @TICTS-B15 @unit-level
  Scenario: In a table only the first cell is the id; on any other line every quoted id counts
    Given an inventory row "| `abcd-screen` | Root of the `TouchableOpacity` | `ABCDX-VR` |", and a list line "- `abcd-screen`, `abcd-close`"
    When the inventory is read
    Then the row declares only "abcd-screen", and the list line declares "abcd-screen" and "abcd-close"

  @TICTS-B16 @unit-level
  Scenario: The attribute's own name is not a declared id
    Given an inventory row "| `testID` (prop) | Root |" in a project whose attribute is testID
    When the inventory is read
    Then no id is declared

  @TICTS-B17 @unit-level
  Scenario: A wildcard at either end covers the concrete ids it opens
    Given the wildcard "abcd-item-*" and the concrete id "abcd-item-3"
    When coverage is asked in both directions
    Then the wildcard covers the concrete id and not the reverse, and a spec declaring the wildcard passes over code exposing ":abcd-item-1" and ":abcd-item-2"

  @TICTS-B18 @unit-level
  Scenario: A handle is queried when a consumer mentions it; a wildcard, when it mentions the prefix
    Given a linked test calling getByTestId('abcd-screen') without the mark, and in turn one building `:abcd-item-${id}` for the declared "abcd-item-*"
    When testid-consistent confronts the spec
    Then both return Pass

  @TICTS-B19 @unit-level
  Scenario: The e2e flows the project declares are consumers
    Given an e2e surface whose file template is "e2e/{{name}}.yaml" with a flow tapping ":abcd-screen", a surface with no path, and a surface given only by an override whose flow does not query the id
    When testid-consistent confronts the spec
    Then the first passes, the second claims nothing about the surface, and the third reports "queried ✗"

  @TICTS-B20 @unit-level
  Scenario: The test files beside the spec, and in its sibling folders, are consumers
    Given a spec in trends/components with no edge to any test, and trends/screens/TrendsScreen.test.tsx querying ":abcd-screen"
    When testid-consistent confronts the spec
    Then it returns Pass, also when e2e flows exist that do not query the id

  @TICTS-I01 @unit-level
  Scenario: One handle is one line of the report, whatever the spelling at each end
    Given code exposing ":abcd-screen", a spec declaring "abcd-screen", and a linked test that does not query it
    When testid-consistent confronts the spec
    Then the report names the handle once, spelled ":abcd-screen"

  @TICTS-X01 @unit-level
  Scenario: A handle only a consumer mentions is not charged to the spec
    Given a linked test querying ":abcd-screen" and ":other-screen-ghost", with code and spec holding only ":abcd-screen"
    When testid-consistent confronts the spec
    Then it returns Pass

  @TICTS-E01 @unit-level
  Scenario: Linked code that cannot be read is not an end
    Given a spec whose only specifies edge points to a missing file
    When testid-consistent confronts the spec
    Then it returns Skip as a spec without linked code
