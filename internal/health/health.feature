# language: en
# @anchors
#   ref: DCTRO
#   updated_at: 2026-10-02
#   layer: feature

@DCTRO
Feature: Doctor — the global health check that hunts the systemic loose ends of a project

  @DCTRO-B01 @unit-level
  Scenario: The report is sorted by check then subject, and counts the map
    Given a toy project with several findings
    When the doctor diagnoses it
    Then each finding is in order of check then subject, and the report counts the map's nodes, edges and declared layers

  @DCTRO-B02 @unit-level
  Scenario: Warnings keeps only the warnings
    Given a report with two warnings and one information
    When its warnings are asked
    Then the two warnings come back, in order

  @DCTRO-B03 @unit-level
  Scenario: A node whose file is gone is a ghost
    Given a spec node Sumido.spec.md with no file on disk
    When the doctor diagnoses the project
    Then there is a no-fantasma finding on Sumido.spec.md

  @DCTRO-B04 @unit-level
  Scenario: An edge to an unknown node is dead
    Given an edge from Login.spec.md to Nowhere.tsx, which is not a node
    When the doctor diagnoses the project
    Then there is an aresta-morta finding on "Login.spec.md → Nowhere.tsx", and none on edges between known nodes

  @DCTRO-B05 @unit-level
  Scenario: A spec that points at nothing has no realization
    Given a spec that is the source of no edge, and one that specifies its code
    When the doctor diagnoses the project
    Then only the first gets a spec-sem-realizacao finding

  @DCTRO-B06 @unit-level
  Scenario: A spec without a code has no identity
    Given a spec node without an identity code
    When the doctor diagnoses the project
    Then there is an identidade-ausente finding on it

  @DCTRO-B07 @unit-level
  Scenario: A code owned by units of different domains is a duplicate identity
    Given code CROSS owned by features/auth/screens/S.tsx and features/dash/screens/H.tsx
    When the doctor diagnoses the project
    Then there is a warning identidade-duplicada on CROSS naming exactly the two owning units

  @DCTRO-B08 @unit-level
  Scenario: Units of the same domain may share a code
    Given code INTR owned by a lib file and a screen of the same feature
    When the doctor diagnoses the project
    Then there is no identidade-duplicada finding on INTR

  @DCTRO-B09 @unit-level
  Scenario: Only specs, features, tests and code own a code
    Given a document in another folder that references a code owned by a single unit
    When the doctor diagnoses the project
    Then the code is not a duplicate identity

  @DCTRO-B10 @unit-level
  Scenario: The files of one unit count as one owner
    Given the spec, feature, test and code of Login at the project root, all with code LOGIN
    When the duplicate identities are checked
    Then there is no finding

  @DCTRO-B11 @unit-level
  Scenario: A file that declares a shared code is not an owner
    Given a code file and a test in another domain with the same code, the test declaring the code as shared
    When the duplicate identities are checked
    Then there is no finding, and without the declaration there is one

  @DCTRO-B12 @unit-level
  Scenario: A declared layer with no node of its kind is empty
    Given a declared feature layer and no feature node
    When the doctor diagnoses the project
    Then there is a camada-vazia finding on feature

  @DCTRO-B13 @unit-level
  Scenario: A guide that governs nothing is reported
    Given a guide node that is the source of no governs edge
    When the doctor diagnoses the project
    Then there is a guide-sem-governo finding on it

  @DCTRO-B14 @unit-level
  Scenario: A spec, feature or test kind that no gate confronts is reported
    Given spec, feature, test, code and document nodes, and one gate on feature
    When the gate coverage is checked
    Then there are kind-sem-gate warnings on spec and test only

  @DCTRO-B15 @unit-level
  Scenario: An unknown perspective in skip_on names the gate and the value
    Given a gate with skip_on chnage, one with skip_on all and one without
    When the skip_on values are checked
    Then there is one finding, on the gate with the typo, quoting chnage

  @DCTRO-B16 @unit-level
  Scenario: A gate whose tool is not on the PATH is reported with its install hint
    Given a gate needing a tool that does not exist with the hint brew install foo, and a gate needing sh
    When the tools are checked
    Then there is one warning ferramenta-ausente on the first gate carrying brew install foo

  @DCTRO-B17 @unit-level
  Scenario: A missing git binary is told to install, not to initialize
    Given no git binary
    When the git presence is checked
    Then the finding is on git and does not mention git init

  @DCTRO-B18 @unit-level
  Scenario: A project under no repository is told to initialize one
    Given git installed and a root under no .git
    When the git presence is checked
    Then there is one warning git-ausente saying git init and not mentioning the PATH

  @DCTRO-B19 @unit-level
  Scenario: A .git in the root or an ancestor, as a folder or a file, is a repository
    Given a root with a .git folder, a subfolder of it, and a root with a .git pointer file
    When the git presence is checked for each
    Then there is no finding

  @DCTRO-B20 @unit-level
  Scenario: In GitHub mode a missing repository names the work queue
    Given a GitHub mode configuration and a root under no .git
    When the git presence is checked
    Then the finding says the work queue has nowhere to come from

  @DCTRO-B21 @unit-level
  Scenario: A needs to a plan that does not exist is broken
    Given a plan whose needs names plans/0099-nao-existe.md, which does not exist
    When the plan needs are checked
    Then there is one needs-quebrado finding naming 0099-nao-existe

  @DCTRO-B22 @unit-level
  Scenario: A cycle of needs gives one finding naming only the cycle, the same on every run
    Given plans c ⇄ d, a ⇄ b, and "plans/0.md" (first in path order) needing a, declared in a shuffled order
    When the plan needs are checked fifty times
    Then every run gives the same single needs-ciclo finding, on "plans/b.md", whose path is "plans/a.md → plans/b.md → plans/a.md"

  @DCTRO-B23 @unit-level
  Scenario: A chain in order, or no plan at all, gives nothing
    Given three plans each needing the previous one
    When the plan needs are checked
    Then there is no finding, and a map without plans gives none either

  @DCTRO-B24 @unit-level
  Scenario: Tests without results and code without coverage are warnings
    Given a test and a code node with no signal
    When the signals are checked
    Then there are warnings for the missing execution and coverage, and none once one node has each signal

  @DCTRO-B25 @unit-level
  Scenario: Missing or partial mutation is informational
    Given two code files, none with mutation, then one of two with mutation
    When the signals are checked
    Then there is an informational missing mutation, then an informational partial mutation saying 1 of 2, and nothing when both have it

  @DCTRO-I01 @unit-level
  Scenario: A healthy project has no warnings
    Given a project with its files on disk, a realized spec, a governing guide, gated kinds, a repository and every signal
    When the doctor diagnoses it
    Then the report has no warnings

  @DCTRO-X01 @unit-level
  Scenario: Code without a spec is not reported
    Given a code file with no spec
    When the doctor diagnoses the project
    Then no finding names that file

  @DCTRO-B26 @unit-level
  Scenario: An ingested report that reached no test is told apart from no report
    Given tests with no result and a spec holding scenarios a suite proved
    When the signals are checked
    Then the warning says the ingested report reached no test file and names the `file` attribute
    And with no spec proven the warning is the one to run the ingest
