# language: en
# @anchors
#   ref: DCTRN
#   updated_at: 2026-09-26
#   layer: feature

@DCTRN
Feature: Doctrine — the vertical axis: product doctrine exists, is realized, and is never copied

  @DCTRN-B01 @unit-level
  Scenario: A TBD marker defers only with a reason and outside backticks
    Given a plan citing the missing doctrine "product/a.doctrine.md"
    When plan-doctrine-exists confronts the plan with "@TBD: written next cycle" on the line, with a bare "@TBD", and with the marker quoted in backticks
    Then the first is Pending
    And the bare marker and the quoted marker both give Fail

  @DCTRN-B02 @unit-level
  Scenario: Without a map the gates that read edges are Pending
    Given no map loaded
    When doctrine-realized, spec-doctrine-exists and doctrine-not-duplicated confront their artifacts
    Then each returns Pending with the "no map loaded" message

  @DCTRN-B03 @unit-level
  Scenario: plan-doctrine-exists skips an artifact that is not a plan
    Given a spec citing "product/a.doctrine.md"
    When plan-doctrine-exists confronts it
    Then it returns Skip

  @DCTRN-B04 @unit-level
  Scenario: A plan whose cited doctrine exists passes
    Given a plan citing "product/a.doctrine.md", which exists on disk
    When plan-doctrine-exists confronts the plan
    Then it returns Pass

  @DCTRN-B05 @unit-level
  Scenario: A plan citing missing doctrines fails, naming each once, in order
    Given a plan citing "product/z.doctrine.md" twice, "product/b.doctrine.md" once, and the existing "product/a.doctrine.md"
    When plan-doctrine-exists confronts the plan
    Then it returns Fail naming "product/b.doctrine.md, product/z.doctrine.md" with the count 2

  @DCTRN-B06 @unit-level
  Scenario: A missing doctrine cited on a TBD line is Pending, naming it
    Given a plan line citing the missing "product/a.doctrine.md" with "@TBD: written next cycle"
    When plan-doctrine-exists confronts the plan
    Then it returns Pending naming "product/a.doctrine.md"

  @DCTRN-B07 @unit-level
  Scenario: A citation with no directory, or of a template, seeds nothing
    Given a plan citing only "x.doctrine.md" and "product/_TEMPLATE.doctrine.md"
    When plan-doctrine-exists confronts the plan
    Then it returns Skip

  @DCTRN-B08 @unit-level
  Scenario: doctrine-realized skips an artifact that is not a doctrine
    Given a spec node
    When doctrine-realized confronts it
    Then it returns Skip

  @DCTRN-B09 @unit-level
  Scenario: A doctrine with no rules is skipped
    Given a doctrine holding only prose
    When doctrine-realized confronts it
    Then it returns Skip

  @DCTRN-B10 @unit-level
  Scenario: A doctrine whose every rule is realized passes
    Given a doctrine with "LIMIT-R01" and "LIMIT-R02"
    And realizes edges naming both rules
    When doctrine-realized confronts it
    Then it returns Pass

  @DCTRN-B11 @unit-level
  Scenario: The unrealized rules of a doctrine fail, and only they are named
    Given a doctrine with "LIMIT-R01" and "LIMIT-R02"
    And a realizes edge naming only "LIMIT-R01"
    When doctrine-realized confronts it
    Then it returns Fail naming "LIMIT-R02" and not "LIMIT-R01"

  @DCTRN-B12 @unit-level
  Scenario: Unrealized rules that are all deferred with TBD are Pending
    Given a doctrine whose only rule "LIMIT-R03" carries "@TBD: not realized yet"
    When doctrine-realized confronts it with no realizer
    Then it returns Pending naming "LIMIT-R03"

  @DCTRN-B13 @unit-level
  Scenario: An open question is not a rule to realize
    Given a doctrine with the open question "LIMIT-Q01" and no realizer for it
    When doctrine-realized confronts it
    Then the verdict does not mention "LIMIT-Q01"

  @DCTRN-B14 @unit-level
  Scenario: spec-doctrine-exists skips an artifact that is not a spec
    Given a code node declaring a realizes tag for "LIMIT-R03"
    When spec-doctrine-exists confronts it
    Then it returns Skip

  @DCTRN-B15 @unit-level
  Scenario: A spec that declares no realization outside TBD lines is skipped
    Given a spec with no realizes tag, and a spec whose only tag sits on a line with "@TBD: doctrine being written"
    When spec-doctrine-exists confronts each of them
    Then both return Skip

  @DCTRN-B16 @unit-level
  Scenario: A declared rule the map resolved, or present in a resolved doctrine, passes
    Given a spec with a realizes edge naming "LIMIT-R03" to a doctrine that also catalogues "LIMIT-R06"
    When spec-doctrine-exists confronts a spec declaring "LIMIT-R03" and "LIMIT-R06"
    Then it returns Pass

  @DCTRN-B17 @unit-level
  Scenario: A declared rule no resolved doctrine holds fails, naming it
    Given a spec whose reached doctrine does not catalogue "LIMIT-R09"
    When spec-doctrine-exists confronts a spec declaring "LIMIT-R03" and "LIMIT-R09"
    Then it returns Fail naming "LIMIT-R09" and not "LIMIT-R03"

  @DCTRN-B18 @unit-level
  Scenario: doctrine-not-duplicated skips an artifact that is not a spec
    Given a code, feature, test and product node
    When doctrine-not-duplicated confronts each of them
    Then each returns Skip

  @DCTRN-B19 @unit-level
  Scenario: A spec with no readable realized doctrine is skipped
    Given a spec with no realizes edge, and a spec whose realized doctrine file was removed
    When doctrine-not-duplicated confronts each of them
    Then both return Skip

  @DCTRN-B20 @unit-level
  Scenario: A corpus under four rules is Pending, and four rules are measured
    Given a doctrine with two rules and a spec copying one of them
    When doctrine-not-duplicated confronts the spec
    Then it returns Pending
    And with three doctrine rules the same copy returns Fail

  @DCTRN-B21 @unit-level
  Scenario: A spec rule that copies the doctrine rule it realizes fails, naming both and the score
    Given a spec rule "CRED-V01" whose text is the text of "LIMIT-R03", which it realizes
    When doctrine-not-duplicated confronts the spec
    Then it returns Fail naming "CRED-V01" and "LIMIT-R03" with "(100%)"

  @DCTRN-B22 @unit-level
  Scenario: A near copy with a word changed still fails
    Given a spec rule that repeats "LIMIT-R03" with one word replaced and one dropped
    When doctrine-not-duplicated confronts the spec
    Then it returns Fail

  @DCTRN-B23 @unit-level
  Scenario: A spec rule with text specific to its unit passes
    Given a spec rule that says it disables the send button when the amount field is empty, realizing "LIMIT-R03"
    When doctrine-not-duplicated confronts the spec
    Then it returns Pass

  @DCTRN-B24 @unit-level
  Scenario: A copy on a line deferred with TBD is not charged
    Given a spec rule copying "LIMIT-R03" on a line with "@TBD: redacao em revisao"
    When doctrine-not-duplicated confronts the spec
    Then it returns Pass

  @DCTRN-B25 @unit-level
  Scenario: spec-realizes-doctrine skips a non-spec and a project with no configuration
    Given a code node of a demanding layer, and a spec in a project with no configuration
    When spec-realizes-doctrine confronts each of them
    Then both return Skip

  @DCTRN-B26 @unit-level
  Scenario: spec-realizes-doctrine skips a spec whose layer does not demand doctrine
    Given a spec of the layer "gate", which does not require doctrine
    When spec-realizes-doctrine confronts a rule with no declaration
    Then it returns Skip

  @DCTRN-B27 @unit-level
  Scenario: A rule with no realizes tag fails where the layer demands doctrine, naming it
    Given a spec of the layer "screen", which requires doctrine
    When spec-realizes-doctrine confronts the rule "CRED-V01" with no declaration
    Then it returns Fail naming "CRED-V01"

  @DCTRN-B28 @unit-level
  Scenario: A rule that declares what it realizes passes
    Given a spec of a demanding layer whose rule "CRED-V01" carries a realizes tag for "LIMIT-R03" on its line, or on the line below
    When spec-realizes-doctrine confronts it
    Then it returns Pass

  @DCTRN-B29 @unit-level
  Scenario: A rule deferred with TBD is debt, not failure
    Given a spec of a demanding layer whose rule carries "@TBD: doctrine being written"
    When spec-realizes-doctrine confronts it
    Then it returns Pending

  @DCTRN-B30 @unit-level
  Scenario: A realizes tag after a blank line declares nothing for the rule above
    Given a spec of a demanding layer with the rule "CRED-V01", a blank line, and then a realizes tag for "LIMIT-R03"
    When spec-realizes-doctrine confronts it
    Then it returns Fail naming "CRED-V01"

  @DCTRN-B31 @unit-level
  Scenario: The demanding layer is the one of the specified target, by edge or by path
    Given a demanding layer "screen" for "screens/**/*.tsx" and a plain layer "lib"
    When the layer of "docs/home.spec.md" is resolved through a specifies edge to a screen, and of "screens/home.spec.md" by path with no map
    Then both demand doctrine
    And a realizes edge, a specifies edge to a lib file, or a lib path do not

  @DCTRN-I01 @unit-level
  Scenario: Only a realizes edge realizes a rule
    Given a doctrine with "LIMIT-R01" and "LIMIT-R02"
    And a realizes edge naming "LIMIT-R01" and a specifies edge carrying "LIMIT-R02"
    When doctrine-realized confronts it
    Then it returns Fail naming "LIMIT-R02"

  @DCTRN-X01 @unit-level
  Scenario: A pair that only shares a rare word is not a copy
    Given a doctrine of eight rules where "LIMIT-R03" is the only one saying "excedido"
    And a spec rule "exibe a mensagem excedido no rodape" realizing "LIMIT-R03"
    When doctrine-not-duplicated confronts the spec
    Then it returns Pass

  @DCTRN-E01 @unit-level
  Scenario: A resolved doctrine that cannot be read confirms none of its rules
    Given a spec with a realizes edge naming "LIMIT-R03" to a doctrine file that was removed
    When spec-doctrine-exists confronts the spec declaring "LIMIT-R03" and "LIMIT-R06"
    Then it returns Fail naming "LIMIT-R06" and not "LIMIT-R03"
