# language: en
# @anchors
#   ref: INWZN
#   updated_at: 2026-09-27
#   layer: feature

@INWZN
Feature: InitWizard — the interactive init walks a person from an unconfigured directory to a reviewed anchors.yaml, and writes nothing on answers nobody gave

  @INWZN-B01 @unit-level
  Scenario: An existing config is kept when the overwrite is not confirmed
    Given an anchors.yaml holding "version: 1" and a comment, and no terminal
    When init runs
    Then it returns without error
    And anchors.yaml is byte for byte the original and nothing announces a write
    And answering yes to the overwrite, with every other question answered, rewrites anchors.yaml

  @INWZN-B02 @unit-level
  Scenario: A ready repository makes the git step silent
    Given a repository with a commit
    When the git step runs
    Then it prints nothing and lets the init go on

  @INWZN-B03 @unit-level
  Scenario: Without git the step warns and the init goes on
    Given a PATH with no git
    When the git step runs
    Then it prints the not-installed warning, asks nothing and lets the init go on

  @INWZN-B04 @unit-level
  Scenario: Declining git names what stays off
    Given a directory outside git and the answer "n"
    When the git step runs
    Then no .git exists
    And it says it proceeds without git and names coverage --diff, install-hooks and git init later

  @INWZN-B05 @unit-level
  Scenario: Accepting git leaves a repository with HEAD
    Given a directory outside git with a .DS_Store, and the answer "y"
    When the git step runs
    Then the repository is ready with the first commit
    And the seeded .gitignore covers .DS_Store
    And and a repository without a commit only gets the commit

  @INWZN-B06 @unit-level
  Scenario: An existing gitignore is kept
    Given a .gitignore holding "# mine" and "*.log"
    When git is initialized
    Then the .gitignore is unchanged

  @INWZN-B07 @unit-level
  Scenario: A failed git initialization is reported and names the fix
    Given git with no author identity and the answer "y"
    When the git step runs
    Then it reports it could not initialize git and that git does not know who you are
    And the fix names git config --global user.email
    And the init goes on

  @INWZN-B08 @unit-level
  Scenario: The findings report only what exists on disk
    Given a repository with specs, features, tests, guides and code, and an empty proposal
    When the findings are printed
    Then the first lists specs, features, tests, guides and code
    And the empty proposal lists none of them

  @INWZN-B09 @unit-level
  Scenario: The DISCOVER step is silent when the phase already happened
    Given a proposal with code, one with specs, and a directory with PROJECT.md
    When the DISCOVER step runs
    Then it prints nothing and lets the init go on

  @INWZN-B10 @unit-level
  Scenario: An AI operator gets the DISCOVER work order
    Given an empty project operated by an AI
    When the DISCOVER step runs
    Then it prints the work order naming anchors guide project, PROJECT.md, INSIGHTS.md and "never in a background worker"
    And the init goes on

  @INWZN-B11 @unit-level
  Scenario: A person without a known AI gets the prompt to paste
    Given an empty project, a person, and no AI tool detected
    When the person is instructed
    Then the step-by-step holds the whole interview prompt
    And nothing offers to open a tool
    And and the prompt is wrapped at the width with every word kept

  @INWZN-B12 @unit-level
  Scenario: A detected AI is opened in the root with the prompt, or declined for the step-by-step
    Given a person with Gemini CLI detected
    When the offer is answered "y", then "n"
    Then "y" runs the tool in the project root with the interview prompt as its one argument
    And "n" prints the step-by-step

  @INWZN-B13 @unit-level
  Scenario: A modular preset finds the module directories
    Given src/modules/users, src/modules/billing and the file src/modules/README.md
    When the modules of the node-ts preset are detected
    Then they are src/modules/billing and src/modules/users
    And and a non-modular preset has none

  @INWZN-B14 @unit-level
  Scenario: Non-interactive routes to the JSON mode
    Given an empty directory
    When init runs with --non-interactive
    Then it prints the questions as JSON

  @INWZN-B15 @unit-level
  Scenario: A stack preset picked from the menu is applied and announced
    Given a repository with code and every question answered
    When the express-ts preset is picked from the menu
    Then the init announces the preset by title
    And its code layers are in the written anchors.yaml

  @INWZN-B16 @unit-level
  Scenario: The header guide is seeded in guides/ when the project has no guide directory
    Given a repository with code and no guide directory
    When the header guide is accepted and every other question is answered
    Then guides/HEADER_GUIDE.md exists
    And no HEADER_GUIDE.md is written at the root

  @INWZN-B17 @unit-level
  Scenario: The default gates are offered only when the chosen artifacts have any, and accepted ones are written
    Given a repository with a spec, whose artifacts yield default gates
    When the gates are accepted and every other question is answered
    Then the written anchors.yaml holds gates
    And an empty project, with no artifact chosen, is never offered gates

  @INWZN-B18 @unit-level
  Scenario: A project with no code, spec, feature or test is announced as new
    Given an empty repository, and a repository with code
    When init runs on an input that already ended
    Then the empty one is announced as new
    And the one with code is not

  @INWZN-B19 @unit-level
  Scenario: The code-layer question is asked only when there are code layers, and a new project is told to declare them later
    Given an empty repository, and a repository with code
    When init runs on an input that already ended
    Then the empty one is not asked about code directories and is told to declare them once they exist
    And the one with code is asked which code directories are layers

  @INWZN-B20 @unit-level
  Scenario: Each guide found is asked which tag it governs
    Given a repository with code and guides/STYLE_GUIDE.md
    When init runs on an input that already ended
    Then it asks which tag STYLE_GUIDE.md governs

  @INWZN-I01 @unit-level
  Scenario: A prompt that cannot run makes the init write nothing
    Given no terminal, or an input that ended
    When init runs outside git, in a repository with code, and in an empty repository
    Then it fails naming --non-interactive
    And no anchors.yaml, no guides/HEADER_GUIDE.md and no .git were created
    And and every prompt helper flags its failure and returns no choice

  @INWZN-X01 @unit-level
  Scenario: Git is never initialized nor committed without a yes
    Given a repository without a commit and no terminal
    When the git step runs
    Then it tells the init to stop
    And no commit exists

  @INWZN-E01 @unit-level
  Scenario: The no-terminal error offers the non-interactive mode
    Given the context "Nothing was written."
    When the no-terminal error is built
    Then it holds the context and names --non-interactive twice with an --artifacts example

  @INWZN-E02 @unit-level
  Scenario: End of input in line mode is refused, not taken as the defaults
    Given TERM=dumb and an input that is already at its end
    When init runs in a repository with code
    Then it fails naming --non-interactive
    And no anchors.yaml and no header guide were written
