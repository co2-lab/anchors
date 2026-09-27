# language: en
# @anchors
#   ref: RPSCR
#   updated_at: 2026-09-26
#   layer: feature

@RPSCR
Feature: RepoScan — the repository read as text: which files exist, of which layer, and what each declares

  @RPSCR-B01 @unit-level
  Scenario: Only files in a declared layer enter the scan
    Given a layer "code" with the pattern "src/*.ts", and the files "src/a.ts" and "docs/readme.md"
    When the repository is walked
    Then only "src/a.ts" is returned, with the layer "code" and the kind "code"

  @RPSCR-B02 @unit-level
  Scenario: What the ignore set excludes is not scanned
    Given a .gitignore listing "probe*.ts", and the files "src/a.ts", "src/probe1.ts", "src/a.ts.swp" and "node_modules/x.ts" under a layer "**/*.ts"
    When the repository is walked
    Then only "src/a.ts" is returned

  @RPSCR-B03 @unit-level
  Scenario: A nested checkout is not scanned
    Given the project file "src/a.ts", a worktree "tools/worktrees/agent-1" with a ".git" file, and a clone "vendor-clone" with a ".git" directory
    When the repository is walked
    Then "src/a.ts" is returned and nothing from the worktree or the clone is

  @RPSCR-B04 @unit-level
  Scenario: A progress companion stays out of the map
    Given the plan "plans/0017-mutacao.md" and its companion "plans/0017-mutacao-progress.md", both matched by the layer "plans/*.md"
    When the repository is walked
    Then the plan is returned and the companion is not

  @RPSCR-B05 @unit-level
  Scenario: An upstream workflow carries no codes
    Given a marked workflow ".github/workflows/anchors-board.yml" citing "FNDTN-F04" in a comment and a body line starting with "parent:"
    When the repository is walked
    Then the workflow is returned as upstream, with no codes and no parent

  @RPSCR-B06 @unit-level
  Scenario: The revision ignores line endings but not content
    Given the same two lines ending in LF and in CRLF, and a third text with a different line
    When the revision of each is computed
    Then the LF and CRLF revisions are equal and the different text has another revision

  @RPSCR-B07 @unit-level
  Scenario: Priority, then pattern length, then layer name decide the layer
    Given the layers "broad" with "pkg/**/*.ts" and "specific" with "pkg/models/**/*.ts", and later "broad" with priority 10
    And two layers "a" and "b" with the same pattern "src/*.ts"
    When "pkg/models/a.ts" and "src/x.ts" are classified
    Then "specific" wins by length, "broad" wins once it declares the priority, and "a" wins the tie of names

  @RPSCR-B08 @unit-level
  Scenario: An exclusion removes a path from its layer
    Given the layer "c" with the pattern "src/**" excluding "src/gen/**"
    When "src/gen/y.ts" and "src/z.ts" are classified
    Then "src/gen/y.ts" has no layer and "src/z.ts" is in "c"

  @RPSCR-B09 @unit-level
  Scenario: A Windows path is classified like its slash form
    Given the layer "code" with the pattern "apps/mobile/src/**/*.ts"
    When "apps/mobile/src/business-logic/a.ts" and its backslash form are classified
    Then both are in the layer "code" with the kind "code"

  @RPSCR-B10 @unit-level
  Scenario: A heuristic decision is reported as an ambiguity, a declared priority is not
    Given the layers "a" and "b" with "src/*.ts" and "c" with "src/**", none with a priority
    And another configuration where "b" declares priority 1 over "a"
    When the ambiguities of "src/x.ts" and "src/z.go" are listed
    Then "src/x.ts" is reported with the winner "a" and the losers "b" and "c", "src/z.go" is not, and nothing is reported under the declared priority

  @RPSCR-B11 @unit-level
  Scenario: The unit's layer comes from the header before the path
    Given a spec whose header declares "layer: lambdas", the code file beside it, and a spec without a header
    When the unit layer of each is asked
    Then they are "lambdas", "lambdas" and "spec"

  @RPSCR-B12 @unit-level
  Scenario: Codes cited in comments are not owned, and each code is listed once
    Given code citing "FOOOX-VR" and "BARRX-S01" in line comments, "AAAAX-S01" in a block comment, and "BAZZX-B01" in a string
    And a text repeating "ABCDX-B01" around "ABCDX-B02"
    When the owned codes are extracted
    Then only "BAZZX-B01" is owned in the first, and the second gives "ABCDX-B01" then "ABCDX-B02"

  @RPSCR-B13 @unit-level
  Scenario: The project's rule letters are recognised
    Given the text "KVALX-P01" and the default rule letters
    When the codes are extracted before and after the project declares the letter "P"
    Then nothing is found before and "KVALX-P01" is found after

  @RPSCR-B14 @unit-level
  Scenario: The declared identity and the annotations are recorded
    Given a spec whose header declares "code: OWNRX", whose body cites "DTAXX-B11" first, and which carries "@anchors-shared-code" and "@noPropagation"
    When the repository is walked
    Then the file's declared identity is "OWNRX" and both flags are set

  @RPSCR-B15 @unit-level
  Scenario: The parent is read only inside the header
    Given a workflow whose body has a "parent:" line, and headers in HTML, line-comment and one-line forms, each followed by a body "parent: NOPE-F09"
    When the parent of each is read
    Then the workflow has none, each header's own parent is read, and "NOPE-F09" is never taken

  @RPSCR-B16 @unit-level
  Scenario: Needs are plan paths for a plan and phase codes for a spec
    Given a plan needing "`plans/a.md`, plans/b.md" and a spec needing "FNDTN-F02, plans/a.md, FNDTN-B01"
    When the needs are read
    Then the plan needs "plans/a.md" and "plans/b.md", the spec needs only "FNDTN-F02", and a code file needs nothing

  @RPSCR-B17 @unit-level
  Scenario: Only a plan revises
    Given a header declaring "revises: plans/a.md"
    When it is read as a plan and as a spec
    Then the plan revises "plans/a.md" and the spec revises nothing

  @RPSCR-B18 @unit-level
  Scenario: A non-spec file declares dependencies in its header
    Given a presentation file under "apps/mobile/src/" whose header declares "dep: theme/tokens.ts, utils/cn.ts"
    When its header dependencies are read
    Then they are "apps/mobile/src/theme/tokens.ts" and "apps/mobile/src/utils/cn.ts"

  @RPSCR-B19 @unit-level
  Scenario: A YAML test script depends on the scripts it composes
    Given the test script "t/f/x.yaml" running "../u/l.yaml"
    When its dependencies are read
    Then it depends on "t/u/l.yaml" with the method "runFlow"

  @RPSCR-B20 @unit-level
  Scenario: A spec's dependency table is read in any catalogue language
    Given a spec with a Portuguese dependencies heading and two rows, a spec with an English "Dependencies" heading, and a spec without the heading
    When their dependencies are read
    Then the first has "DEP1" and "DEP2", the second has its row, and the third has none

  @RPSCR-B21 @unit-level
  Scenario: A row with a malformed code is not a dependency
    Given a dependency table with a row coded "notdep" and a row coded "DEP1"
    When the dependencies are read
    Then only "DEP1" is a dependency

  @RPSCR-B22 @unit-level
  Scenario: The method keeps its backticks
    Given a dependency row whose method cell is "`useAuthStore`" and whose layer cell is "store"
    When the dependencies are read
    Then the method is "`useAuthStore`" with backticks and the layer is "store"

  @RPSCR-B23 @unit-level
  Scenario: A declared file resolves through the src fallbacks
    Given the file "apps/mobile/src/hooks/x.ts" and a spec under "apps/mobile/src/"
    When "src/hooks/x.ts", "hooks/x.ts" and "nope/x.ts" are resolved
    Then the first two give "apps/mobile/src/hooks/x.ts" and the third stays "nope/x.ts"

  @RPSCR-B24 @unit-level
  Scenario: A realizes tag pairs with its rule in the three rule forms
    Given a heading rule, a table row rule and a bold bullet rule, each with an "@realizes" tag on its line or the next
    When the realizations are read
    Then each tag is paired with its own rule

  @RPSCR-B25 @unit-level
  Scenario: A tag after a blank line has no owning rule
    Given a rule heading, a blank line, and a paragraph with "@realizes ORFA-R01"
    When the realizations are read
    Then the tag has no owning rule

  @RPSCR-B26 @unit-level
  Scenario: Only a spec declares rule tags
    Given a heading rule with an "@realizes" tag
    When it is read as code, test, feature, plan and product
    Then no realization is read

  @RPSCR-B27 @unit-level
  Scenario: A repeated pair is recorded once
    Given the pair "CRED-V01" to "LIMIT-R03" written twice, "CRED-V02" to "LIMIT-R03", and "CRED-V03" to "LIMIT-R03" and "LIMIT-R09"
    When the realizations are read
    Then there are 4 realizations

  @RPSCR-B28 @unit-level
  Scenario: A gated-by tag names only a flag scenario
    Given a rule tagged "@gated-by FLAGX-G01" and "@gated-by FLAGX-B02"
    When the gates are read in a spec and in a code file
    Then the spec gives only "FLAGX-G01" owned by the rule, and the code file gives nothing

  @RPSCR-B29 @unit-level
  Scenario: A plan seeds only concrete spec and doctrine paths
    Given a plan citing "`*.spec.md`", "`apps/x/Tela.spec.md`", "`_TEMPLATE_SCREEN.spec.md`", a bare "`MutualTls.spec.md`", "`a/b.doctrine.md`" and "`apps/x/Tela.spec.md`" again
    When its seeds are read, and the same text is read as a spec
    Then the plan seeds "apps/x/Tela.spec.md" and "a/b.doctrine.md" once each, and the spec seeds nothing

  @RPSCR-I01 @unit-level
  Scenario: The classification is stable across runs
    Given two layers that both match "pkg/models/a.ts"
    When the path is classified fifty times
    Then every run gives the same layer

  @RPSCR-X01 @unit-level
  Scenario: A guide's dependencies table is not a dependency
    Given a guide with a dependencies table that has a code and a file column, and a documental table without a code column
    When their dependencies are read
    Then neither yields a dependency, while the same table in a spec yields one

  @RPSCR-X02 @unit-level
  Scenario: A spec does not read a dep header line
    Given a header declaring "dep: theme/tokens.ts"
    When it is read as code and as a spec
    Then the code file has one dependency and the spec has none

  @RPSCR-E01 @unit-level
  Scenario: A root that cannot be walked returns the error
    Given a root path that does not exist
    When the repository is walked
    Then the walk returns an error
