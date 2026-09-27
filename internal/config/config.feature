# language: en
# @anchors
#   ref: CNFGO
#   updated_at: 2026-09-27
#   layer: feature

@CNFGO
Feature: Config — loads the project's anchors.yaml, refuses what it cannot honour, and answers every setting with its default

  @CNFGO-B01 @unit-level
  Scenario: An unknown key in the file is a load error naming the key
    Given an anchors.yaml with the top-level key "governz", and another whose governs item uses "guide" and "tags"
    When each file is loaded
    Then each load fails, and the error names the unknown key

  @CNFGO-B02 @unit-level
  Scenario: An unknown key gets the misspelling and old-binary hypotheses
    Given an anchors.yaml in format 1 with the unknown key "chaveQueNaoExiste"
    When the file is loaded
    Then the error names the key, says the binary may be OLD, and gives the go install command

  @CNFGO-B03 @unit-level
  Scenario: A renamed key in an older-format file advises the migration
    Given a migration registry that knows "trinca_opcional" as renamed
    And an anchors.yaml in format 1, or with no version field, that declares "trinca_opcional"
    When the file is loaded
    Then the error advises "anchors migrate" and does not advise go install

  @CNFGO-B04 @unit-level
  Scenario: The migration advice needs both an older format and a renamed key
    Given an older-format file with a misspelled key, a current-format file with a renamed key, and a file loaded with no migration registry
    When each file is loaded
    Then none of the errors advises "anchors migrate", and the misspelled key gets the misspelling hypothesis

  @CNFGO-B05 @unit-level
  Scenario: An error that is not an unknown key gets no version hint
    Given an anchors.yaml with a YAML syntax error
    When the file is loaded
    Then the load fails, and the error carries neither the old-binary hint nor the migration advice

  @CNFGO-B06 @unit-level
  Scenario: Two gates sharing an ID fail the load
    Given an anchors.yaml with gate "a" and a gate "b" that declares id "a"
    When the file is loaded
    Then the load fails naming both gates, while gates with distinct IDs load

  @CNFGO-B07 @unit-level
  Scenario: An unknown scope, full scope, cost, phase or perspective value fails the load
    Given gates declaring scope "repo", cost "lento", phase "precommit", scope_full "projct" or "node", and skip_on "chnage"
    When each file is loaded
    Then each load fails naming the gate and the wrong value, and the valid values of the five lists load

  @CNFGO-B08 @unit-level
  Scenario: Local and manual modes refuse the GitHub fields
    Given an anchors.yaml with no workflow block, one in manual mode, and one in local or manual mode that declares a repo
    When each file is loaded
    Then the first two load as non-GitHub workflows, and the ones with a repo fail

  @CNFGO-B09 @unit-level
  Scenario: GitHub mode requires an owner/name repository and a label
    Given GitHub-mode workflows without a repo, with repo "exemplo", without labels, and one complete
    When each file is loaded
    Then the three incomplete ones fail naming what is missing, and the complete one loads in GitHub mode

  @CNFGO-B10 @unit-level
  Scenario: An unknown workflow mode fails with no fallback
    Given a workflow with mode "guithub"
    When the file is loaded
    Then the load fails and says there is no fallback between modes

  @CNFGO-B44 @unit-level
  Scenario: A test level's code filter accepts by allow and refuses by exclude
    Given a level with no filter, a level that allows only VR codes, and a level that allows its unit prefix but excludes VR codes
    When each is asked about a rule code and a VR code
    Then the first accepts both, the second only the VR code, the third only the rule code
    And a filter pattern that does not compile fails the load naming the gate, the level, the list and the index

  @CNFGO-B11 @unit-level
  Scenario: A declared pattern that does not compile fails the load naming the field
    Given dialect and derived patterns that do not compile, and a valid guard pattern with the parameter placeholder
    When each file is loaded
    Then each broken one fails naming its field, and the valid one loads

  @CNFGO-B12 @unit-level
  Scenario: An unsupported language fails the load
    Given an anchors.yaml declaring lang "xx", and another declaring lang "es"
    When each file is loaded
    Then the first fails listing the supported languages, and after the second the current language is "es"

  @CNFGO-B13 @unit-level
  Scenario: Code lengths outside two to eight fail, and valid ones reach the engine
    Given anchors.yaml files declaring code lengths [1], [9] and [4, 5], and a registered generator hook
    When each file is loaded
    Then the first two fail naming the length, and after the third the engine and the hook both hold [4, 5]

  @CNFGO-B42 @unit-level
  Scenario: A file with no code lengths restores the default, whatever an earlier load set
    Given an anchors.yaml declaring code lengths [4] already loaded, and a registered generator hook
    When an anchors.yaml that declares no code lengths is loaded
    Then the engine and the hook both hold [5]

  @CNFGO-B43 @unit-level
  Scenario: Every load refusal and the header Save writes are in the project's language
    Given a pt-BR file loaded first, and files with a bad cost, an incomplete GitHub workflow, an unknown mode, GitHub fields in local mode, a length of 9 and a repeated gate ID
    When each is loaded declaring lang "en", and again declaring lang "pt-BR", and configurations with no lang, "en" and "pt-BR" are saved
    Then every refusal under "en" has no Portuguese and every one under "pt-BR" is Portuguese, and each saved file starts with the header of its own language

  @CNFGO-B14 @unit-level
  Scenario: A canonical gate inherits every field the project omitted
    Given a catalog that declares "gate-canonico" with check, measures, ask, scope, cost, category, targets, phases and blocking
    And an anchors.yaml that declares that gate by name only
    When the file is loaded
    Then the loaded gate carries every one of those fields from the catalog

  @CNFGO-B15 @unit-level
  Scenario: A field the project declared wins over the canonical one
    Given a canonical gate with measures "canônico" and blocking true
    And a project gate declaring measures "do projeto", or blocking false
    When the file is loaded
    Then the gate keeps the project's measures and does not block, while the omitted ask is inherited

  @CNFGO-B16 @unit-level
  Scenario: Declaring run or check inherits neither of the pair
    Given a canonical gate with a check and a project gate with only a run, and the reverse
    When each file is loaded
    Then each gate keeps only the dispatch the project declared

  @CNFGO-B17 @unit-level
  Scenario: A gate the catalog does not know loads untouched
    Given a catalog that knows no gate, and a project gate "meu-gate" with a run
    When the file is loaded
    Then the gate has no inherited check nor ask

  @CNFGO-B18 @unit-level
  Scenario: A gate with no declared severity does not block
    Given a gate declaring no blocking value
    When its severity is read
    Then it does not block

  @CNFGO-B19 @unit-level
  Scenario: The scope defaults to one run per target, and the full scan uses scope_full only when it is batch or project
    Given gates with no scope, an unknown scope, scope batch, and scope_full project or node
    When the scope is read for an incremental and for a full scan
    Then an undeclared or unknown scope is node, and the full scan answers scope_full only for batch or project

  @CNFGO-B20 @unit-level
  Scenario: A gate with no phases runs in every phase
    Given a gate with no phases and a gate declaring only "manual"
    When each is asked about the pre-commit, manual and unnamed phases
    Then the first runs in all of them, and the second only in manual and the unnamed phase

  @CNFGO-B21 @unit-level
  Scenario: A gate participates in every perspective unless skip_on excludes it
    Given a gate with no skip_on and a gate with skip_on [change]
    When each is asked about the change and all perspectives
    Then only the second skips, and only the change perspective

  @CNFGO-B22 @unit-level
  Scenario: The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format
    Given configurations with no format, format "  GREMLINS  " on mutation-score, and format "gremlins" on another gate
    When the mutation format is read
    Then it is gremlins only for the mutation-score gate, and mutation-testing-elements otherwise

  @CNFGO-B23 @unit-level
  Scenario: Section language is checked unless the gate turns it off
    Given a gate with no section language setting, one set true and one set false
    When the section language setting is read
    Then only the gate set false does not check

  @CNFGO-B24 @unit-level
  Scenario: The integration branch defaults to main, and a non-main integration branch protects main too
    Given no workflow, a workflow with integration branch "develop", and one declaring three protected branches
    When the branches are read
    Then they are main and [main]; develop and [develop, main]; and the three declared ones

  @CNFGO-B25 @unit-level
  Scenario: One approval is required unless the project declares another number, zero included
    Given no workflow, a workflow with no approvals declared, one declaring 0 and one declaring 2
    When the required approvals are read
    Then they are 1, 1, 0 and 2

  @CNFGO-B26 @unit-level
  Scenario: A stale pipeline or a manual ingest blocks only when the project asks
    Given no workflow, an empty workflow, and one declaring both blocks true
    When the two switches are read
    Then only the last one blocks

  @CNFGO-B27 @unit-level
  Scenario: Only an explicit enabled false freezes the project
    Given configurations with no enabled value, enabled true, enabled false, and no configuration at all
    When the freeze is read
    Then only enabled false is frozen

  @CNFGO-B28 @unit-level
  Scenario: The freeze reason is shown trimmed, and a missing one asks for freeze_reason
    Given a frozen configuration with a reason surrounded by spaces, one with a blank reason, and one with none
    When the freeze reason text is read
    Then the first shows the trimmed reason, and the others show a message naming freeze_reason

  @CNFGO-B29 @unit-level
  Scenario: A section title comes from the layer, then the project, then the framework
    Given a project renaming "rules" to "Regras" and layer "screen" renaming it to "Comportamentos", with a blank rename in layer "blank"
    When the title of "rules" is asked for layers "screen", "blank" and "other", and for key "states"
    Then the answers are Comportamentos, Regras, Regras and the framework default

  @CNFGO-B30 @unit-level
  Scenario: Placeholder markers default to the templates' marker word
    Given a configuration declaring " FIXME " and a blank marker, one declaring only blanks, and one declaring none
    When the placeholder markers are read
    Then the first is [FIXME], and the other two are the templates' default marker list

  @CNFGO-B31 @unit-level
  Scenario: Rule letters come from the declared rule types, or the canonical set
    Given rule types with letters "b", "B", "XY", "s" and "", and a configuration with only an invalid letter
    When the rule letters are read
    Then the first answers "BS" and the second the canonical SRVAXBNMDEIQFG

  @CNFGO-B32 @unit-level
  Scenario: A scenario tag maps to every letter that declares it
    Given rule types S and V both declaring the tag "@estado-dado"
    When the tag " @Estado-Dado " and the tag "@smoke" are looked up
    Then the first maps to S and V, and the second is unknown

  @CNFGO-B33 @unit-level
  Scenario: The code length pattern matches exactly the declared lengths, contiguous or not
    Given the engine configured with lengths [5], [4, 5], [4, 6], [5, 7] and [3, 5, 8]
    When the pattern is placed after one character class and matched against codes of 2 to 9 characters
    Then exactly the declared lengths match each time: 4 and 6 for [4, 6], never 5 or 7

  @CNFGO-B41 @unit-level
  Scenario: A rule type catalogues the sections it declares, ignoring case and surrounding spaces
    Given a rule type that declares the section "Business Rules" as requiring a code
    When it is asked about " business rules ", and about "Notes"
    Then the first section requires a code and the second does not

  @CNFGO-B34 @unit-level
  Scenario: With no filter every suite is selected
    Given four suites across two workspaces, and suites with no scope
    When suites are selected with no filter, or with a filter on another axis
    Then all four come back, and the suites with no scope are not excluded

  @CNFGO-B35 @unit-level
  Scenario: The filter axes intersect
    Given suites for unit and integration in backend, unit and e2e in mobile, and mutation suites with scopes full and isolated
    When suites are selected by layer, by workspace, by scope and by their combinations
    Then each selection holds exactly the suites matching every filtered axis

  @CNFGO-B36 @unit-level
  Scenario: Selected suites keep the order of the file
    Given suites declared as b-unit, m-unit, m-e2e among others
    When the layers "e2e" and "unit" are requested in that order
    Then the selection is b-unit, m-unit, m-e2e

  @CNFGO-B37 @unit-level
  Scenario: Suite names match ignoring case and surrounding spaces
    Given suites with layer "unit" and workspace "backend"
    When layer " Unit " and workspace "BackEnd" are requested
    Then one suite is selected and nothing is reported missing

  @CNFGO-B38 @unit-level
  Scenario: A name missing from the file is reported with its axis, and an empty combination is not a missing name
    Given suites with no layer "smoke", no workspace "web", no scope "parcial", and no integration suite in mobile
    When those names, and the integration-in-mobile combination, are requested
    Then smoke, web and parcial come back labelled by axis, and the combination selects nothing with nothing missing

  @CNFGO-B39 @unit-level
  Scenario: The declared vocabulary lists each name once, in file order
    Given suites repeating the layer "unit" across workspaces, and suites with no workspace
    When the declared layers and workspaces are listed
    Then the layers are unit, integration, e2e, the workspaces backend, mobile, and no empty workspace appears

  @CNFGO-B40 @unit-level
  Scenario: Patterns replace the code template and keep the other derived files
    Given derived files declaring code, feature and two patterns, and derived files declaring only code
    When the patterns of each are resolved
    Then the first has the two patterns as code, keeps feature, and has no derived file named patterns, and the second keeps its code

  @CNFGO-I01 @unit-level
  Scenario: What Save writes, Load reads back
    Given a loaded configuration with one layer and one blocking gate
    When it is saved and the written file is loaded again
    Then the reloaded configuration has the same layer and a blocking gate

  @CNFGO-I02 @unit-level
  Scenario: A configuration with every key known keeps loading
    Given an anchors.yaml with layers, governs, gates and a touch block
    When the file is loaded
    Then it loads with one layer, one governs, one gate and the touch settings as written

  @CNFGO-I03 @unit-level
  Scenario: Every reader answers its default on a nil configuration
    Given no configuration and no workflow
    When every reader is asked
    Then each answers its default: canonical letters, the default marker list, the framework title, main, one approval, not frozen, no GitHub mode

  @CNFGO-X01 @unit-level
  Scenario: GitHub mode never infers the repository from the git remote
    Given a GitHub-mode workflow with labels and no repo, loaded inside a directory
    When the file is loaded
    Then the load fails asking for repo in owner/name form instead of reading a remote

  @CNFGO-X02 @unit-level
  Scenario: Suite names are the project's vocabulary, never a fixed list
    Given suites with workspace "cobranca" and layers "contrato" and "carga"
    When layer "carga" in workspace "cobranca" is requested
    Then the carga suite is selected and nothing is reported missing

  @CNFGO-E01 @unit-level
  Scenario: A file that cannot be read fails the load with the read error
    Given a path where no anchors.yaml exists
    When the file is loaded
    Then the load fails with a file-not-found error and no configuration

  @CNFGO-E02 @unit-level
  Scenario: Save to a path that cannot be written returns the write error
    Given a path inside a directory that does not exist
    When a configuration is saved there
    Then saving returns a file-not-found error
