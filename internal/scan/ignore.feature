# language: en
# @anchors
#   ref: SCIGS
#   updated_at: 2026-09-26
#   layer: feature

@SCIGS
Feature: ScanIgnore — what the scan never sees: the built-in list, the project's `.gitignore`, and editor noise

  @SCIGS-B01 @unit-level
  Scenario: The built-in directories are skipped when nothing is declared
    Given a project whose .gitignore only lists "*.log" and no layer points inside "build"
    When the scan asks whether to skip the directory "build"
    Then the directory is skipped

  @SCIGS-B02 @unit-level
  Scenario: A layer pointing inside a built-in directory re-enables it, a catch-all does not
    Given a layer with the pattern "build/**/*.ts" and another project with only the pattern "**/*.ts"
    When the scan asks whether to skip "build" in the first and "node_modules" in the second
    Then "build" is scanned and "node_modules" is still skipped

  @SCIGS-B03 @unit-level
  Scenario: A gitignore negation re-enables a built-in directory
    Given a .gitignore containing "!build/"
    When the scan asks whether to skip the directory "build"
    Then the directory is scanned

  @SCIGS-B04 @unit-level
  Scenario: The records Anchors writes are never scanned
    Given a .gitignore containing "!issues/"
    When the scan asks whether to skip "issues" at the root and "changes" nested under "a/"
    Then both directories are skipped

  @SCIGS-B05 @unit-level
  Scenario: Editor and system ephemera never become files to scan
    Given the paths "amplify/data/.!21662!resource.spec.md", "src/a.ts.swp", "src/a.ts~", "src/.#a.ts", "src/#a.ts#", "src/x.tmp" and ".DS_Store"
    And the paths "src/a.ts" and "src/tmp/util.ts"
    When the scan asks whether to skip each file
    Then every ephemeral file is skipped and both material files are kept

  @SCIGS-B06 @unit-level
  Scenario: A slash anchors a gitignore pattern at the root, no slash matches at any depth
    Given a .gitignore with "/data/", "node_modules", "*.log" and "docs/saida"
    When the scan asks about "data", "amplify/data", "a/node_modules", "a/b/c.log", "docs/saida" and "apps/docs/saida"
    Then "data", "a/node_modules", "a/b/c.log" and "docs/saida" are ignored, while "amplify/data" and "apps/docs/saida" are not

  @SCIGS-B07 @unit-level
  Scenario: A trailing slash ignores the directory and what is below it, but not a file of that name
    Given a .gitignore with "data/"
    When the scan asks about the directory "data", the file "data/dump.json" and a plain file named "data"
    Then the directory and the file below it are ignored, and the plain file is not

  @SCIGS-B08 @unit-level
  Scenario: The last matching gitignore rule decides
    Given a .gitignore with "*.log" followed by "!keep.log"
    When the scan asks about "a.log" and "keep.log"
    Then "a.log" is ignored and "keep.log" is scanned

  @SCIGS-B09 @unit-level
  Scenario: Without a loaded ignore set the fixed exclusions still hold
    Given no ignore set at all
    When the scan asks about the directories "node_modules", "issues", ".git" and ".anchors" and the files "a.swp" and "a.ts"
    Then the four directories and "a.swp" are skipped, and "a.ts" is kept

  @SCIGS-I01 @unit-level
  Scenario: No declaration re-enables the machinery directories
    Given a .gitignore containing "!.git/" and "!.anchors/"
    When the scan asks whether to skip ".git" and ".anchors"
    Then both are skipped

  @SCIGS-X01 @unit-level
  Scenario: A nested gitignore does not change what the scan sees
    Given a root without a .gitignore and a "sub/.gitignore" that lists "x.ts"
    When the scan asks whether to skip "sub/x.ts"
    Then the file is kept
