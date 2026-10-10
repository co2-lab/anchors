# language: en
# @anchors
#   code: SCFTC
#   ref: RPSCR
#   updated_at: 2026-10-09
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
    Given a marked workflow ".github/workflows/anchors-board.yml" citing "FNDTN-W04" in a comment and a body line starting with "parent:"
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
    Given the text "KVALX-Z01" and the default rule letters
    When the codes are extracted before and after the project declares the letter "Z"
    Then nothing is found before and "KVALX-Z01" is found after

  @RPSCR-B14 @unit-level
  Scenario: The declared identity and the annotations are recorded
    Given a spec whose header declares "code: OWNRX", whose body cites "DTAXX-B11" first, and which carries "@anchors-shared-code" and "@noPropagation"
    When the repository is walked
    Then the file's declared identity is "OWNRX" and both flags are set

  @RPSCR-B15 @unit-level
  Scenario: The parent is read only inside the header
    Given a workflow whose body has a "parent:" line, and headers in HTML, line-comment and one-line forms, each followed by a body "parent: NOPE-W09"
    When the parent of each is read
    Then the workflow has none, each header's own parent is read, and "NOPE-W09" is never taken

  @RPSCR-B16 @unit-level
  Scenario: Needs are plan paths for a plan and phase codes for a spec
    Given a plan needing "`plans/a.md`, plans/b.md" and a spec needing "FNDTN-W02, plans/a.md, FNDTN-B01"
    When the needs are read
    Then the plan needs "plans/a.md" and "plans/b.md", the spec needs only "FNDTN-W02", and a code file needs nothing

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

  @RPSCR-B30 @unit-level
  Scenario: Header keys are read only inside the header
    Given a spec whose header declares nothing but whose body has the lines "code: BOGUS", "layer: bogus", "needs: FNDTN-W02, plans/a.md" and "revises: plans/a.md", and a code file whose body has the comment "// dep: a.ts"
    When its identity, unit layer, needs, revisions and dependencies are read
    Then none of them is taken: no identity, no declared layer, no needs, no revision, no dependency, and the unit layer is the file's layer "spec"

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

  @RPSCR-E02 @unit-level
  Scenario: A layer file that cannot be read fails the walk
    Given the file "src/a.ts" in the layer "src/*.ts", with no read permission
    When the repository is walked
    Then the walk returns an error naming "src/a.ts"

  @RPSCR-B31 @unit-level
  Scenario: Rule tags follow the declared code length
    Given a project that declares 7-character codes
    And a rule carrying @realizes and @gated-by tags of that length
    When the spec is scanned
    Then both edges are read, from the rule to each tagged code

  @RPSCR-B32 @unit-level
  Scenario: A file in its layer's support list is marked as support
    Given a test layer whose support list covers its utils folder, and another layer whose support list covers everything
    When the project is scanned
    Then the utils file is support, the flow beside it is not, and the other layer's list marks nothing here

  @RPSCR-B33 @unit-level
  Scenario: Only the given files are read, as the walk reads them
    Given a tree with a spec, a file in no layer, a spec under node_modules, and a missing path
    When those paths are scanned, and the governed paths are listed
    Then only the spec comes back, equal to what the walk gives, and the listing names the walk's paths

  @RPSCR-B34 @unit-level
  Scenario: The staged walk reads the index, not the tree
    Given a repository with a committed spec edited but not staged, a new spec staged, and an untracked spec
    When the staged walk runs
    Then the edited spec has its committed revision, the staged one is read, and the untracked one is not

  @RPSCR-B35 @unit-level
  Scenario: The index reader reads what the commit records
    Given a committed file edited and not staged, a new staged file, an untracked one and an untouched one
    When each is read through the index reader
    Then the edited one reads as committed, the staged one as staged, the untracked one is absent, and the untouched one reads from the tree

  @RPSCR-B36 @unit-level
  Scenario: The governed files where the tree and the index part
    Given a governed file edited and not staged, one staged, an untracked governed file and an untracked file of no layer
    When the tree's changes are listed
    Then the unstaged and the untracked governed files are listed, and nothing else

  @RPSCR-B37 @unit-level
  Scenario: The index reader confronts the tree at each read
    Given a committed file deleted from the tree and another edited, both after the index reader was made
    When each is read through the index reader
    Then both read as the index has them

  @RPSCR-B38 @unit-level
  Scenario: A rule is defined in any of the three forms
    Given the lines "### ABCDE-B01 — x", "| `ABCDE-B02` | y |", "- **ABCDE-B03** — z" and "as ABCDE-B04 says"
    When each is matched as a rule definition
    Then the first three define ABCDE-B01, ABCDE-B02 and ABCDE-B03
    And the prose line defines nothing

  @RPSCR-B39 @unit-level
  Scenario: A spec's Parts Used names its components
    Given a spec whose "## Componentes Utilizados", the title the project declares, lists `BottomSheet` and `MoveIcon`
    When it is scanned
    Then its Composes are BottomSheet and MoveIcon, and a code file's are none

  @RPSCR-B40 @unit-level
  Scenario: The header's ref names the units the file realizes
    Given a header with "ref: ARENA, WALLT"
    When it is scanned
    Then its HeaderRefs are ARENA and WALLT

  @RPSCR-B41 @unit-level
  Scenario: The @dep and @no-dep flags of import lines are read, with the symbols each import brings
    Given a file whose imports carry dependency flags — one with named and aliased symbols, one default — and one a waiver
    When its flags are read
    Then each flag has its code, its symbols and its line, and the waiver its reason

  @RPSCR-B42 @unit-level
  Scenario: Each @used-by flag is read with the symbol declared below it
    Given a file with two exported symbols, each under a used-by flag
    When its flags are read
    Then each flag has the codes that use it and the symbol below it

  @RPSCR-B43 @unit-level
  Scenario: Each @navigates and @no-nav flag is read with its screens, its rule and its call's line
    Given a navigation call flagged on its own line with a rule, a back navigation flagged on the line above with two screens, and a call waived
    When its flags are read
    Then each has its screens, its rule and its call's line, and the waiver its reason

  @RPSCR-B44 @unit-level
  Scenario: A spec's Out rows are read by rule, each with a revision of the row alone
    Given a Portuguese Out table with rows ARNAA-A03 and ARNAA-A04, a row with no rule, and a rule table after it
    When its Out rows are read, and then with A03's destination changed
    Then two rows are read, and only A03's revision changes

  @RPSCR-B45 @unit-level
  Scenario: A spec's Dependencies table is no dependency of the map
    Given a spec with a Dependencies table naming a file
    When its dependencies are read for the map, and then for the migration
    Then the map gets none, and the migration reads the row

  @RPSCR-B46 @unit-level
  Scenario: A kinded dependency is read with its kind and name
    Given a line flagged `@dep[db]: transactions` and an import flagged with a code
    When the dependency flags are read
    Then the first has kind db and name transactions, and the second its code

  @RPSCR-B47 @unit-level
  Scenario: A dormant navigation flag is read with its reason
    Given a navigation call flagged with its screen, its rule and a dormant reason
    When the navigation flags are read
    Then the flag has its screen, its rule and the reason

  @RPSCR-B48 @unit-level
  Scenario: A spec's evidence leaves out its header, navigation, history and spacing
    Given a spec with a header date, rules in rows and under a heading, an Out table and a change history
    When its evidence is read, and then with another date, Out row, history line, spacing, rule, and title
    Then the date, the Out row, the history and the spacing move nothing, the rule moves its own revision alone, and the title moves the rest

  @RPSCR-B49 @unit-level
  Scenario: A file's evidence leaves out its header and the chain's flags, and its line revision keeps the lines
    Given a module, the same with flags at the end of its lines, the same with a flag line of its own, and the same under headers of two dates
    When their evidence is read, and then with a change to the code
    Then they have one evidence, the line-end flags and the header's date keep the lines and the flag line moves them, and the code's change moves the evidence

  @RPSCR-B50 @unit-level
  Scenario: A type dependency is an import of types
    Given an import line flagged `@dep[type]: TYPSU`
    When the dependency flags are read
    Then it has the code TYPSU, the kind type and the import's symbols, and names no resource

  @RPSCR-B51 @unit-level
  Scenario: The project's sections with no side effect are left out of the evidence
    Given a project declaring its Implementation Notes section with no side effect
    When a spec's notes change, and then its history
    Then the notes move nothing, and the history — no longer in the list — moves the evidence

  @RPSCR-B52 @unit-level
  Scenario: A line that is only a comment is no evidence
    Given a module gaining a rule citation on a line of its own, a block comment, a JSX comment, and a flow gaining a comment line
    When their evidence is read, and then with a comment at the end of a line of code changed, code after a block's end changed, and a file of no known language
    Then the comments move no evidence but move the lines, and the others move the evidence
