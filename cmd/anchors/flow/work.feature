# language: en
# @anchors
#   code: WRFTW
#   ref: WRPRW
#   updated_at: 2026-10-03
#   layer: feature

@WRPRW
Feature: WorkPrompt — compose the work prompt of one stage over one target, from what the project declares

  @WRPRW-B01 @unit-level
  Scenario: The command refuses what it cannot compose
    Given a project with a configuration, and a root without one
    When `anchors work deploy --for src/pricing.ts`, `anchors work code` without `--for`, and `anchors work code --for x.ts` on the root without configuration run
    Then they fail with `unknown artifact "deploy"`, naming `--for`, and with "load config"

  @WRPRW-B02 @unit-level
  Scenario: A derived piece is redirected to its unit
    Given the unit "src/metadataVersioning.ts" with its spec, feature and test beside it
    When the unit of each piece is resolved
    Then each resolves to "src/metadataVersioning.ts"
    And a map edge makes "a/different-name.spec.md" resolve to "b/otherName.ts", also from a feature linked to that spec
    And "src/x.ts" is not redirected, and "src/ghost.spec.md" with no code resolves to nothing
    And `anchors work code --for src/pricing.spec.md` says on standard error that the piece is not the unit

  @WRPRW-B03 @unit-level
  Scenario: A waived piece stops the prompt
    Given the layer "model" waiving the feature
    When the feature prompt of "models/user.ts" is composed
    Then it has "## STOP" and "the layer **model** WAIVES the piece `feature`"
    And it has no "## Procedure"

  @WRPRW-B04 @unit-level
  Scenario: A declarative layer has no unit piece
    Given the declarative layer "dao"
    When the test prompt of "daos/user.ts" is composed
    Then it has "## STOP" and "declared as RECOGNIZED"

  @WRPRW-B05 @unit-level
  Scenario: The heading and the role follow the stage
    Given the unit "src/pricing.ts" and the plan "plans/0001-f.md"
    When the spec, review and review-plan prompts are composed
    Then they start with "# Work: spec of `src/pricing.ts`", "# Review of `src/pricing.ts`" and "# WHOLE review — `plans/0001-f.md`"
    And the reviews say "You are the REVIEWER of this unit." and "You are the REVIEWER OF THE WHOLE."

  @WRPRW-B06 @unit-level
  Scenario: The target section names the layer, or says it is unclassified
    Given the unit "src/pricing.ts" in the layer "logic", and "scripts/tool.py" in no layer
    When their prompts are composed
    Then the first has "- Layer: **logic** (regime: comportamental)" and "- Tags: backend"
    And the second has "- Layer: **unclassified**"

  @WRPRW-B07 @unit-level
  Scenario: The guides to read are the artifact's and the layer's, each once and sorted
    Given "guides/SPEC_GUIDE.md" governing spec, and "guides/BACKEND.md" governing both backend and spec
    When the spec prompt of "src/pricing.ts" is composed
    Then it lists "2. `guides/BACKEND.md`" and "3. `guides/SPEC_GUIDE.md`", BACKEND once
    And a target with no governing guide says "No project guide governs this layer"

  @WRPRW-B08 @unit-level
  Scenario: The pieces are listed where the project derives them
    Given derived paths beside the code, a test override for the layer "handler", and "src/billing/charge.spec.md" on disk
    When the prompts are composed
    Then the spec stage marks "→ `src/pricing.spec.md` — spec"
    And the handler's test is "tests/unit/push.test.ts" noted as a layer override
    And "src/billing/charge.spec.md" is marked "(already exists)" and its feature is not
    And "packages/infra/GoLiveChecklist.spec.md" derives "GoLiveChecklist.feature", not "GoLiveChecklist.spec.feature"
    And a project without derived paths is told so, and a declarative code stage lists "**Only the file itself.**"

  @WRPRW-B09 @unit-level
  Scenario: Feature and test stages list the regime tags
    Given the regimes "unit-level" and "integration-level", the second proved on "api"
    When the test prompt is composed
    Then it lists "`@integration-level` (regime integration) → proved on the surface `api`" before "`@unit-level` (regime unit)"
    And the spec prompt lists no regime tags

  @WRPRW-B10 @unit-level
  Scenario: Each stage says what is not its scope
    Given the unit "src/pricing.ts"
    When the spec, test and review prompts are composed
    Then they say "**Do not write code, feature or test**", "**Do not mock the logic under test**" and "**Do not correct the code**"

  @WRPRW-B11 @unit-level
  Scenario: The procedure follows the stage and the layer
    Given the layer "logic" with the extra code step "run the migrations before the handler", and the declarative layer "dao"
    When the code prompts are composed
    Then the logic prompt appends "- run the migrations before the handler" and keeps "Read the whole spec"
    And the dao prompt says "Translate, transport or declare — **do not decide**" and not "Read the whole spec"

  @WRPRW-B12 @unit-level
  Scenario: The gates' demands are listed once, for the stage
    Given two gates with the check spec-sections on spec, an informative feature-test-match on feature, and a code gate with an undescribed check
    When the spec, test and code prompts are composed
    Then the spec prompt says "**Every rule must be CATALOGUED**" once
    And the test prompt carries "declares which scenario it proves" marked "*(informative)*"
    And the code prompt has no "What the gates will demand"

  @WRPRW-B13 @unit-level
  Scenario: A review confronts the unit's delivery records
    Given two pending records for "src/pricing.ts" and one for "src/tax.ts"
    When the review prompt of "src/pricing.ts" is composed
    Then it has "## Delivery records of this unit (2)" in order, and not the other unit's
    And without a record, local mode says "**No delivery record**" and github mode says the record "is in the COMMENTS of the issue"

  @WRPRW-B14 @unit-level
  Scenario: The findings already recorded for the unit are listed
    Given open issues about "src/pricing.spec.md" vs "src/pricing.ts", about "src/pricing.test.ts", about "src/tax.ts", about "src/pricing-v2.ts" and about "lib/src/pricing.ts"
    When the review prompt of "src/pricing.ts" is composed
    Then it lists the two pricing issues under "## Findings ALREADY RECORDED about this unit"
    And it does not list the tax, the pricing-v2 nor the lib/src/pricing issue

  @WRPRW-B15 @unit-level
  Scenario: Tests and unit reviews explain the execution signals
    Given the unit "src/pricing.ts"
    When the test, review and spec prompts are composed
    Then the test and review prompts have "## Execution signals" and the spec prompt does not

  @WRPRW-B16 @unit-level
  Scenario: Producing stages record the delivery, reviews close with a verdict
    Given the unit "src/pricing.ts"
    When the spec and review prompts are composed
    Then the spec prompt has "anchors deliver --stage spec --unit src/pricing.ts"
    And the review prompt has "anchors judge src/pricing.ts --gate review --verdict pass"

  @WRPRW-B17 @unit-level
  Scenario: The verification runs over the stage's own piece
    Given the unit "src/pricing.ts"
    When the spec prompt is composed
    Then it has "anchors check --changed src/pricing.spec.md --no-record --deterministic"
    And then "anchors check --changed src/pricing.spec.md --deterministic"

  @WRPRW-B18 @unit-level
  Scenario: The waivers are listed for the stage
    Given the unit "src/pricing.ts"
    When the spec and code prompts are composed
    Then the spec prompt names "`@no-mark: <reason>`" and "`@no-scenario: <reason>`"
    And the code prompt names "`@no-paginate: <reason>`" and "`@allow-boundary: <reason>`"

  @WRPRW-B19 @unit-level
  Scenario: When the ruler does not decide, the open-decisions section is named
    Given the unit "src/pricing.ts"
    When the spec and code prompts are composed
    Then the spec prompt says to write the question in the open-decisions section, or its "none" value
    And the code prompt says to stop and read that section in the spec

  @WRPRW-B20 @unit-level
  Scenario: Rule marking follows the project's policy
    Given a project whose rule-marking policy is required, and one without a policy
    When the marking step is written
    Then the first says "this project REQUIRES the marking"
    And the second says "if the project uses that pattern"

  @WRPRW-I01 @unit-level
  Scenario: No prompt both records a delivery and closes a review
    Given the unit "src/pricing.ts"
    When every stage's prompt is composed
    Then spec, code, feature and test have only the record section, and the three reviews only the close section

  @WRPRW-X01 @unit-level
  Scenario: Composing writes nothing
    Given an empty project root
    When the review prompt is composed
    Then the root is still empty

  @WRPRW-X02 @unit-level
  Scenario: The prompt cites one open-decisions title and value and the current names
    Given a project whose gates include open-questions-resolved on specs and rule-implemented on code
    When the spec, code and test prompts of "src/pricing.ts" and the feature and spec prompts of the waiving "models/user.ts" are composed
    Then the spec prompt names "## Open Decisions" and "none" in both the procedure and the gates' demands, and never "Decisões em aberto" nor "nenhuma"
    And the prompts cite "optional_unit_edges", "rule-implemented" and "tests-pass", and never "trinca_opcional", "regra-implementada" nor "testes-passam"
