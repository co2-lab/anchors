<!-- anchors:generated from doct/regras.md.tmpl — inputs:0a800c6b19e3eafa — DO NOT EDIT: run `anchors docs build` -->


# Regras

Todas as regras do sistema, de todas as unidades. Cada uma leva ao texto na página da sua
camada.

Para ver uma camada inteira de uma vez — com a visão geral de cada unidade e os cenários —
abra a página dela em `camadas/`.

## apoio

### [BRCRB — BoardCards — reading the repository's board, where the work queue lives in github mode](camadas/apoio.md#brcrb--boardcards--reading-the-repositorys-board-where-the-work-queue-lives-in-github-mode)

### [CLWTC — ClaimWait — asking the claim pipeline for a card and waiting, bounded, for the answer](camadas/apoio.md#clwtc--claimwait--asking-the-claim-pipeline-for-a-card-and-waiting-bounded-for-the-answer)

### [GHRNG — GhRunner — running `gh` for the board so that a failure names its cause and never waits for input](camadas/apoio.md#ghrng--ghrunner--running-gh-for-the-board-so-that-a-failure-names-its-cause-and-never-waits-for-input)

### [DCLND — DocLinks — the anchors, links, sizes and layer arrows the documentation templates are given](camadas/apoio.md#dclnd--doclinks--the-anchors-links-sizes-and-layer-arrows-the-documentation-templates-are-given)

### [C4CNC — C4Containers — the declared containers, each with the units that run in it, for the architecture page](camadas/apoio.md#c4cnc--c4containers--the-declared-containers-each-with-the-units-that-run-in-it-for-the-architecture-page)

### [DTCDC — DocTemplateCompiler — compiles documentation pages from templates that reference the specs' content](camadas/apoio.md#dtcdc--doctemplatecompiler--compiles-documentation-pages-from-templates-that-reference-the-specs-content)

### [DCKND — DocKinds — what each kind of project documentation must answer, told to the agent that writes it](camadas/apoio.md#dcknd--dockinds--what-each-kind-of-project-documentation-must-answer-told-to-the-agent-that-writes-it)

### [DCOXX — DocLayout — the one decision of whether a documentation page shows everything or summarizes](camadas/apoio.md#dcoxx--doclayout--the-one-decision-of-whether-a-documentation-page-shows-everything-or-summarizes)

### [DCSCD — DocScaffolds — the starting templates `anchors docs init` proposes, one page per question and per layer](camadas/apoio.md#dcscd--docscaffolds--the-starting-templates-anchors-docs-init-proposes-one-page-per-question-and-per-layer)

### [GSRGH — GherkinScenarioReader — the scenarios of a unit's feature, with their steps, for the documentation](camadas/apoio.md#gsrgh--gherkinscenarioreader--the-scenarios-of-a-units-feature-with-their-steps-for-the-documentation)

### [FLGRF — FlagGrammar — the fixed grammar of a feature-flag scenario's condition](camadas/apoio.md#flgrf--flaggrammar--the-fixed-grammar-of-a-feature-flag-scenarios-condition)

### [FLPRF — FlagParse — reading a project's flag files and the scenarios their values open](camadas/apoio.md#flprf--flagparse--reading-a-projects-flag-files-and-the-scenarios-their-values-open)

### [FLBLF — FlowBuild — assembling the work-flow graph from the project's flow and action files](camadas/apoio.md#flblf--flowbuild--assembling-the-work-flow-graph-from-the-projects-flow-and-action-files)

### [FLMDF — FlowModel — the questions a work-flow graph answers: what comes next, and what is broken](camadas/apoio.md#flmdf--flowmodel--the-questions-a-work-flow-graph-answers-what-comes-next-and-what-is-broken)

### [INCTA — I18nCatalog — every message a person reads, in the project's language, with a fallback that never goes blank](camadas/apoio.md#incta--i18ncatalog--every-message-a-person-reads-in-the-projects-language-with-a-fallback-that-never-goes-blank)

### [LGSCL — LogScan — finding the occurrences of declared failures in the project's logs](camadas/apoio.md#lgscl--logscan--finding-the-occurrences-of-declared-failures-in-the-projects-logs)

### [MGLTR — CodeLetters — rewriting the letter of the codes of one kind of unit](camadas/apoio.md#mgltr--codeletters--rewriting-the-letter-of-the-codes-of-one-kind-of-unit)

### [MGFLM — MigrateFile — takes one project file from its declared format to the current one, renaming only what is a key](camadas/apoio.md#mgflm--migratefile--takes-one-project-file-from-its-declared-format-to-the-current-one-renaming-only-what-is-a-key)

### [MGSTM — MigrationSteps — the registered steps, one per format, take any project from format 1 to the current format](camadas/apoio.md#mgstm--migrationsteps--the-registered-steps-one-per-format-take-any-project-from-format-1-to-the-current-format)

### [MSCMG — MigrationStepChain — one step per format version, chained in order, and a hole in the chain is an error](camadas/apoio.md#mscmg--migrationstepchain--one-step-per-format-version-chained-in-order-and-a-hole-in-the-chain-is-an-error)

### [AGRLG — AgentRoles — who is who in a project, and the capabilities each role carries](camadas/apoio.md#agrlg--agentroles--who-is-who-in-a-project-and-the-capabilities-each-role-carries)

### [USSTS — UserSettings — the agent's local settings: the declared role, kept out of git](camadas/apoio.md#ussts--usersettings--the-agents-local-settings-the-declared-role-kept-out-of-git)

### [TLCNT — TelemetryConfig — the opt-out is easy to find, easy to get right, and the environment overrides the file](camadas/apoio.md#tlcnt--telemetryconfig--the-opt-out-is-easy-to-find-easy-to-get-right-and-the-environment-overrides-the-file)

### [TLEVT — TelemetryEvent — a decision event has a name from a closed vocabulary, a caller-stamped instant and attributes](camadas/apoio.md#tlevt--telemetryevent--a-decision-event-has-a-name-from-a-closed-vocabulary-a-caller-stamped-instant-and-attributes)

### [TLNTT — TelemetryNotice — the telemetry notice reaches whoever did not ask for it, once, and says how to turn it off](camadas/apoio.md#tlntt--telemetrynotice--the-telemetry-notice-reaches-whoever-did-not-ask-for-it-once-and-says-how-to-turn-it-off)

### [TLEMT — TelemetryEmitter — decision events leave as OTLP logs, never block the work, and never carry who uses the product](camadas/apoio.md#tlemt--telemetryemitter--decision-events-leave-as-otlp-logs-never-block-the-work-and-never-carry-who-uses-the-product)

## comando

### [AGCRG — AgentCards — the cards this agent owns, and the card a pull request declares, read from the tracker](camadas/comando.md#agcrg--agentcards--the-cards-this-agent-owns-and-the-card-a-pull-request-declares-read-from-the-tracker)

### [CMCLC — CommonCLI — the contract every command shares: how a path becomes a node, how "not governed" is signalled, what the binary says it is](camadas/comando.md#cmclc--commoncli--the-contract-every-command-shares-how-a-path-becomes-a-node-how-not-governed-is-signalled-what-the-binary-says-it-is)

### [FLALF — FlagAliases — a renamed command flag keeps answering to its old name](camadas/comando.md#flalf--flagaliases--a-renamed-command-flag-keeps-answering-to-its-old-name)

### [ARCGN — AgentRoleCLI — who this agent is, and the role it declared, as the commands show and ask it](camadas/comando.md#arcgn--agentrolecli--who-this-agent-is-and-the-role-it-declared-as-the-commands-show-and-ask-it)

### [UNCDN — UnitCodes — the identity code of a unit, read from a header, from the map, or from the codes a file names](camadas/comando.md#uncdn--unitcodes--the-identity-code-of-a-unit-read-from-a-header-from-the-map-or-from-the-codes-a-file-names)

### [TLSTT — TelemetrySetup — every command starts telemetry the same way: the opt-outs first, then the notice, then the emitter](camadas/comando.md#tlstt--telemetrysetup--every-command-starts-telemetry-the-same-way-the-opt-outs-first-then-the-notice-then-the-emitter)

### [BCLBB — BackfillLabels — write into open cards the blocking and provenance links the board already implies](camadas/comando.md#bclbb--backfilllabels--write-into-open-cards-the-blocking-and-provenance-links-the-board-already-implies)

### [DCDDE — Decided — release the card an escalation stopped for a person, once the decision became a rule](camadas/comando.md#dcdde--decided--release-the-card-an-escalation-stopped-for-a-person-once-the-decision-became-a-rule)

### [DLVRE — Deliver — record what a stage delivered, where the reviewer reads it, and send the author to the review](camadas/comando.md#dlvre--deliver--record-what-a-stage-delivered-where-the-reviewer-reads-it-and-send-the-author-to-the-review)

### [DLCND — DeliveryConfront — confront what a delivery declares against the disk, at the moment it is declared](camadas/comando.md#dlcnd--deliveryconfront--confront-what-a-delivery-declares-against-the-disk-at-the-moment-it-is-declared)

### [DSCRD — Discard — take off the board a card that no longer makes sense, without deleting it](camadas/comando.md#dscrd--discard--take-off-the-board-a-card-that-no-longer-makes-sense-without-deleting-it)

### [SCLTE — Escalate — open the right card for a change the plan, the spec or the tool needs, and stop the work only when it must](camadas/comando.md#sclte--escalate--open-the-right-card-for-a-change-the-plan-the-spec-or-the-tool-needs-and-stop-the-work-only-when-it-must)

### [ESDPS — EscalateDuplicate — find the open cards that already deal with the target of an escalation](camadas/comando.md#esdps--escalateduplicate--find-the-open-cards-that-already-deal-with-the-target-of-an-escalation)

### [NTFCT — Notifications — a message to every agent, read from one file and printed on top of `next`](camadas/comando.md#ntfct--notifications--a-message-to-every-agent-read-from-one-file-and-printed-on-top-of-next)

### [PRBDP — PRBody — write the lines that link a pull request to its cards, in the platform's syntax](camadas/comando.md#prbdp--prbody--write-the-lines-that-link-a-pull-request-to-its-cards-in-the-platforms-syntax)

### [PLPRP — PlanProgress — create a plan's progress file, the state that lives beside the decision and outside the map](camadas/comando.md#plprp--planprogress--create-a-plans-progress-file-the-state-that-lives-beside-the-decision-and-outside-the-map)

### [PRMRP — ProgressMerge — the git merge driver of a plan's progress file: unite both sides, done beats pending](camadas/comando.md#prmrp--progressmerge--the-git-merge-driver-of-a-plans-progress-file-unite-both-sides-done-beats-pending)

### [WRQUW — WorkQueue — list, pull, close and discard the work, from the local queue or from the board](camadas/comando.md#wrquw--workqueue--list-pull-close-and-discard-the-work-from-the-local-queue-or-from-the-board)

### [FLRGF — FlowRegister — attach the flow domain's commands to the root command, once each](camadas/comando.md#flrgf--flowregister--attach-the-flow-domains-commands-to-the-root-command-once-each)

### [TSSTT — TaskStatus — discover what the machine knows about the task at hand, so the agent's report does not have to](camadas/comando.md#tsstt--taskstatus--discover-what-the-machine-knows-about-the-task-at-hand-so-the-agents-report-does-not-have-to)

### [TSRTS — TaskStatusReport — the format of the round's report: where the task is, the verdict, what is missing, and what comes next](camadas/comando.md#tsrts--taskstatusreport--the-format-of-the-rounds-report-where-the-task-is-the-verdict-what-is-missing-and-what-comes-next)

### [NBLCK — Unblock — open the work card a decision demanded, linked to the card stuck waiting for a person](camadas/comando.md#nblck--unblock--open-the-work-card-a-decision-demanded-linked-to-the-card-stuck-waiting-for-a-person)

### [WTCHA — Watch — the background watcher that turns "a file changed" into "there is work in the queue"](camadas/comando.md#wtcha--watch--the-background-watcher-that-turns-a-file-changed-into-there-is-work-in-the-queue)

### [WTDMW — WatchDaemon — the watcher started in the background survives the terminal that started it](camadas/comando.md#wtdmw--watchdaemon--the-watcher-started-in-the-background-survives-the-terminal-that-started-it)

### [WRPRW — WorkPrompt — compose the work prompt of one stage over one target, from what the project declares](camadas/comando.md#wrprw--workprompt--compose-the-work-prompt-of-one-stage-over-one-target-from-what-the-project-declares)

### [DTAUI — Audit — the dossier of everything pending on one file, for fixing it in one pass](camadas/comando.md#dtaui--audit--the-dossier-of-everything-pending-on-one-file-for-fixing-it-in-one-pass)

### [CMPLN — Compliance — the state of each regulatory duty, grouped by the norm that imposes it](camadas/comando.md#cmpln--compliance--the-state-of-each-regulatory-duty-grouped-by-the-norm-that-imposes-it)

### [GVGDG — GovernanceGuides — the guides an agent reads to operate Anchors, and the contracts other code relies on](camadas/comando.md#gvgdg--governanceguides--the-guides-an-agent-reads-to-operate-anchors-and-the-contracts-other-code-relies-on)

### [GVRNS — Governs — who each guide governs, and how many, read from the map](camadas/comando.md#gvrns--governs--who-each-guide-governs-and-how-many-read-from-the-map)

### [ATGDT — AutonomyGuide — what an agent does with what it does not know, by the role declared locally](camadas/comando.md#atgdt--autonomyguide--what-an-agent-does-with-what-it-does-not-know-by-the-role-declared-locally)

### [SPGDS — SpecGuide — the project's own spec guide, instantiated with its dialect and a complete example](camadas/comando.md#spgds--specguide--the-projects-own-spec-guide-instantiated-with-its-dialect-and-a-complete-example)

### [CLMNC — CliMain — the entry point that stamps the build identity, prints a failure once and turns it into the exit code the hooks read](camadas/comando.md#clmnc--climain--the-entry-point-that-stamps-the-build-identity-prints-a-failure-once-and-turns-it-into-the-exit-code-the-hooks-read)

### [FLRSA — Failures — the observed failures that the spec has not explained yet](camadas/comando.md#flrsa--failures--the-observed-failures-that-the-spec-has-not-explained-yet)

### [FLWOX — Flow — build, draw and navigate the work flows kept in the map](camadas/comando.md#flwox--flow--build-draw-and-navigate-the-work-flows-kept-in-the-map)

### [MPCTI — Impact — what a change to one file reaches, in both directions of the map](camadas/comando.md#mpcti--impact--what-a-change-to-one-file-reaches-in-both-directions-of-the-map)

### [NGSTI — Ingest — binds the test and log signals the project produced to the nodes of the map](camadas/comando.md#ngsti--ingest--binds-the-test-and-log-signals-the-project-produced-to-the-nodes-of-the-map)

### [JDGUE — Judge — records an AI's verdict on a judgment gate with the same bookkeeping as a deterministic gate](camadas/comando.md#jdgue--judge--records-an-ais-verdict-on-a-judgment-gate-with-the-same-bookkeeping-as-a-deterministic-gate)

### [MPCMM — MapCommand — builds the dependency map from the project and answers questions about it](camadas/comando.md#mpcmm--mapcommand--builds-the-dependency-map-from-the-project-and-answers-questions-about-it)

### [MPMRM — MapMerge — the git merge driver that unites two versions of the map instead of merging text](camadas/comando.md#mpmrm--mapmerge--the-git-merge-driver-that-unites-two-versions-of-the-map-instead-of-merging-text)

### [MPSTM — MapStaleness — names the files of the map whose content changed after the map was built](camadas/comando.md#mpstm--mapstaleness--names-the-files-of-the-map-whose-content-changed-after-the-map-was-built)

### [RCDEO — Recode — renames an identity code and carries the change to every textual surface of the project](camadas/comando.md#rcdeo--recode--renames-an-identity-code-and-carries-the-change-to-every-textual-surface-of-the-project)

### [MPRGM — MapRegister — hangs the map domain's commands on the root command](camadas/comando.md#mprgm--mapregister--hangs-the-map-domains-commands-on-the-root-command)

### [RNMBR — Renumber — moves the revisions a branch added when the base already took their number](camadas/comando.md#rnmbr--renumber--moves-the-revisions-a-branch-added-when-the-base-already-took-their-number)

### [BRSRB — BoardServe — the board page served locally with live state, read from the host only when something changed](camadas/comando.md#brsrb--boardserve--the-board-page-served-locally-with-live-state-read-from-the-host-only-when-something-changed)

### [CLGCM — anchors changelog — the technical changelog, printed or written](camadas/comando.md#clgcm--anchors-changelog--the-technical-changelog-printed-or-written)

### [CDCMC — CodeCommand — a new unit gets an identity code that no other unit in the map already owns, and the codes in use are listed from the map](camadas/comando.md#cdcmc--codecommand--a-new-unit-gets-an-identity-code-that-no-other-unit-in-the-map-already-owns-and-the-codes-in-use-are-listed-from-the-map)

### [CMMSC — CommitMsg — the commit subject is confronted with the format the changelog will read, before the commit exists](camadas/comando.md#cmmsc--commitmsg--the-commit-subject-is-confronted-with-the-format-the-changelog-will-read-before-the-commit-exists)

### [DCCMD — DocsCommand — the documentation is compiled from the specs through templates, against a map rebuilt from the tree](camadas/comando.md#dccmd--docscommand--the-documentation-is-compiled-from-the-specs-through-templates-against-a-map-rebuilt-from-the-tree)

### [FRZEX — Freeze — the project is stopped in three layers with a written reason, and thawed by undoing exactly those layers](camadas/comando.md#frzex--freeze--the-project-is-stopped-in-three-layers-with-a-written-reason-and-thawed-by-undoing-exactly-those-layers)

### [GNPTG — GeneratedPaths — the product names the files it derives, so a conflict in them is rebuilt, not merged](camadas/comando.md#gnptg--generatedpaths--the-product-names-the-files-it-derives-so-a-conflict-in-them-is-rebuilt-not-merged)

### [ININT — InitNonInteractive — the init that an agent answers with flags: it asks in JSON, and writes only a complete, valid set of answers](camadas/comando.md#inint--initnoninteractive--the-init-that-an-agent-answers-with-flags-it-asks-in-json-and-writes-only-a-complete-valid-set-of-answers)

### [INWZN — InitWizard — the interactive init walks a person from an unconfigured directory to a reviewed anchors.yaml, and writes nothing on answers nobody gave](camadas/comando.md#inwzn--initwizard--the-interactive-init-walks-a-person-from-an-unconfigured-directory-to-a-reviewed-anchorsyaml-and-writes-nothing-on-answers-nobody-gave)

### [INHKN — InstallHooks — the git hooks that confront every commit and push with the gates and the freeze, installed without taking a hook the user wrote](camadas/comando.md#inhkn--installhooks--the-git-hooks-that-confront-every-commit-and-push-with-the-gates-and-the-freeze-installed-without-taking-a-hook-the-user-wrote)

### [MGCMM — MigrateCommand — the command the format error promises, bringing the map and the config up to this binary's format](camadas/comando.md#mgcmm--migratecommand--the-command-the-format-error-promises-bringing-the-map-and-the-config-up-to-this-binarys-format)

### [NWARN — NewArtifact — a new artifact is born beside its unit, with a resolved identity and the sections of the project's ruler](camadas/comando.md#nwarn--newartifact--a-new-artifact-is-born-beside-its-unit-with-a-resolved-identity-and-the-sections-of-the-projects-ruler)

### [NWTMN — NewTemplates — the catalog of artifact skeletons: which kinds exist, their headers, their sections and the presets that pick them](camadas/comando.md#nwtmn--newtemplates--the-catalog-of-artifact-skeletons-which-kinds-exist-their-headers-their-sections-and-the-presets-that-pick-them)

### [OPRGP — OpsRegister — the operation commands reach the CLI through one registration point](camadas/comando.md#oprgp--opsregister--the-operation-commands-reach-the-cli-through-one-registration-point)

### [STCMS — SettingsCommand — one agent's local decisions, declared with a date and kept out of the project's configuration](camadas/comando.md#stcms--settingscommand--one-agents-local-decisions-declared-with-a-date-and-kept-out-of-the-projects-configuration)

### [SGCMS — SuggestCommand — the proposed fixes are listed, shown, applied or rejected, and every decision keeps its record](camadas/comando.md#sgcms--suggestcommand--the-proposed-fixes-are-listed-shown-applied-or-rejected-and-every-decision-keeps-its-record)

### [SYCMS — SynthesizeCommand — two pull requests in content conflict become one card that asks for the best of each, and every end points at it](camadas/comando.md#sycms--synthesizecommand--two-pull-requests-in-content-conflict-become-one-card-that-asks-for-the-best-of-each-and-every-end-points-at-it)

### [BDGRN — BudgetRun — run a suite's files fastest first until a time budget is spent](camadas/comando.md#bdgrn--budgetrun--run-a-suites-files-fastest-first-until-a-time-budget-is-spent)

### [CGPCH — CheckGatePipeline — confronts the map's nodes against the declared gates, records the verdicts and reports the profile](camadas/comando.md#cgpch--checkgatepipeline--confronts-the-maps-nodes-against-the-declared-gates-records-the-verdicts-and-reports-the-profile)

### [CVCMC — CoverageCommand — answers the confidence questions from the ingested signals: by scenario, by line, of the diff and the delta](camadas/comando.md#cvcmc--coveragecommand--answers-the-confidence-questions-from-the-ingested-signals-by-scenario-by-line-of-the-diff-and-the-delta)

### [HLDCH — DoctorCommand — the global health x-ray, and the repair of the github-mode environment](camadas/comando.md#hldch--doctorcommand--the-global-health-x-ray-and-the-repair-of-the-github-mode-environment)

### [LCBCL — LocalBacklog — what is still open locally after a full check, said in two lines](camadas/comando.md#lcbcl--localbacklog--what-is-still-open-locally-after-a-full-check-said-in-two-lines)

### [QLCMQ — QualityCommands — the quality domain puts its eleven commands under the root command](camadas/comando.md#qlcmq--qualitycommands--the-quality-domain-puts-its-eleven-commands-under-the-root-command)

### [RPRTS — Reports — markdown perspectives on what Anchors already measures, written into docs](camadas/comando.md#rprts--reports--markdown-perspectives-on-what-anchors-already-measures-written-into-docs)

### [SLCTN — RunSelection — a run takes only what is stale and below the minimum, unless told otherwise](camadas/comando.md#slctn--runselection--a-run-takes-only-what-is-stale-and-below-the-minimum-unless-told-otherwise)

### [STEDS — StaleEdges — lists the confrontation debt: expired test evidence and stale edges](camadas/comando.md#steds--staleedges--lists-the-confrontation-debt-expired-test-evidence-and-stale-edges)

### [CNSTC — StampCommand — writes the missing contract stamps on test doubles, and refreshes them after a change](camadas/comando.md#cnstc--stampcommand--writes-the-missing-contract-stamps-on-test-doubles-and-refreshes-them-after-a-change)

### [PRSTP — ProjectStatus — where the project stands in the cycle, and the one next step](camadas/comando.md#prstp--projectstatus--where-the-project-stands-in-the-cycle-and-the-one-next-step)

### [STPRS — SuiteProxy — runs the test and mutation suites the project declared, and binds their reports to the map](camadas/comando.md#stprs--suiteproxy--runs-the-test-and-mutation-suites-the-project-declared-and-binds-their-reports-to-the-map)

### [HDTHD — HeaderDateTouch — bumps the header date of the files that changed, and only of those](camadas/comando.md#hdthd--headerdatetouch--bumps-the-header-date-of-the-files-that-changed-and-only-of-those)

### [VPFVR — VerifyPhaseFacade — one invocation per phase, delegated to the check pipeline](camadas/comando.md#vpfvr--verifyphasefacade--one-invocation-per-phase-delegated-to-the-check-pipeline)

### [CLRTC — CliRoot — every command passes through one root that speaks the project's language and honours the freeze](camadas/comando.md#clrtc--cliroot--every-command-passes-through-one-root-that-speaks-the-projects-language-and-honours-the-freeze)

## config

### [CMMRC — CommentMarkers — which text opens a line comment in each kind of file](camadas/config.md#cmmrc--commentmarkers--which-text-opens-a-line-comment-in-each-kind-of-file)

### [CNFGO — Config — loads the project's anchors.yaml, refuses what it cannot honour, and answers every setting with its default](camadas/config.md#cnfgo--config--loads-the-projects-anchorsyaml-refuses-what-it-cannot-honour-and-answers-every-setting-with-its-default)

### [CNTNR — Containers — what runs separately, and which layers run inside each](camadas/config.md#cntnr--containers--what-runs-separately-and-which-layers-run-inside-each)

### [DLCTI — Dialect — the lexicon of the project's language, between an agnostic gate and concrete code](camadas/config.md#dlcti--dialect--the-lexicon-of-the-projects-language-between-an-agnostic-gate-and-concrete-code)

### [DCRQA — DocsRequired — which aggregate documentation a change of a unit obliges to touch](camadas/config.md#dcrqa--docsrequired--which-aggregate-documentation-a-change-of-a-unit-obliges-to-touch)

### [MNVRM — MinVersion — whether the running binary meets the minimum version the project declares](camadas/config.md#mnvrm--minversion--whether-the-running-binary-meets-the-minimum-version-the-project-declares)

### [DRPTD — DerivedPatterns — a derived file declared as one pattern or as a list of them](camadas/config.md#drptd--derivedpatterns--a-derived-file-declared-as-one-pattern-or-as-a-list-of-them)

### [PRRPR — ProjectRootResolution — the project root a command works on](camadas/config.md#prrpr--projectrootresolution--the-project-root-a-command-works-on)

### [GTVCG — GateVocabulary — the list of default gate names, injected into the configuration layer](camadas/config.md#gtvcg--gatevocabulary--the-list-of-default-gate-names-injected-into-the-configuration-layer)

## gate

### [CDCTC — CodeCataloged — what the code EXPORTS must be in the spec, or waived in the code](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code)

### [CDLNG — CodeLanguage — the code does not go back to mixing languages](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages)

### [CRVCD — CodeReferenceValid — cross-referenced requirement codes must resolve to existing units](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units)

### [CSDCN — ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those)

### [CNHNC — CountHonored — a numerical assertion written in a spec must match reality in code](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code)

### [DEPHN — DependencyHonored — methods promised in the dependency table are consumed in code](camadas/gate.md#dephn--dependencyhonored--methods-promised-in-the-dependency-table-are-consumed-in-code)

### [DCRQD — DocRequired — the aggregated document the unit must feed](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed)

### [DSCDC — DocSelfContained — the spec has to stand on its own](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own)

### [DCCVD — DocsCovered — every spec must reach some page of the compiled documentation](camadas/gate.md#dccvd--docscovered--every-spec-must-reach-some-page-of-the-compiled-documentation)

### [DCFRD — DocsFresh — the compiled document has to reflect the spec](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec)

### [DCTRN — Doctrine — the vertical axis: product doctrine exists, is realized, and is never copied](camadas/gate.md#dctrn--doctrine--the-vertical-axis-product-doctrine-exists-is-realized-and-is-never-copied)

### [DMDCD — DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid)

### [DUPLC — Duplication — no code file holds a block copied from somewhere else](camadas/gate.md#duplc--duplication--no-code-file-holds-a-block-copied-from-somewhere-else)

### [EVFRV — EvidenceFresh — the score of this test holds against TODAY's code](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code)

### [EXCMX — ExternalCommand — executes external tools via shell passing targets as positional arguments](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments)

### [FLRAI — Failure — the failure a spec declares must be handled, recorded, and every handling declared](camadas/gate.md#flrai--failure--the-failure-a-spec-declares-must-be-handled-recorded-and-every-handling-declared)

### [FTMFT — FeatureTestMatch — scenarios in feature must be implemented in test by code and description](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description)

### [FXIXX — Fix — the self-healer that applies the mechanical, safe repairs of `check --fix`](camadas/gate.md#fxixx--fix--the-self-healer-that-applies-the-mechanical-safe-repairs-of-check---fix)

### [FLSCF — FlagScenarios — the scenarios a feature flag declares are written, complete, cited and tested](camadas/gate.md#flscf--flagscenarios--the-scenarios-a-feature-flag-declares-are-written-complete-cited-and-tested)

### [GTENG — GateEngine — which gates reach which node, and what the run concludes](camadas/gate.md#gteng--gateengine--which-gates-reach-which-node-and-what-the-run-concludes)

### [IDCND — IdentityConsistent — a unit's spec identity must match its exposed testID and visual baseline](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline)

### [INCHN — InternalChecks — the registry that routes a declared check name to a function](camadas/gate.md#inchn--internalchecks--the-registry-that-routes-a-declared-check-name-to-a-function)

### [LYBNL — LayerBoundary — a layer does not reach what is not its own](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own)

### [MRPRM — MarkerParity — the same rule has to appear at BOTH ends that fulfil it](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it)

### [MKSTP — MockStampGenerator — writes the missing `@contract` stamps, and never rewrites one](camadas/gate.md#mkstp--mockstampgenerator--writes-the-missing-contract-stamps-and-never-rewrites-one)

### [MCSTM — MockStamped — the double carries the mark of the snippet it replaces, and the gate RECOMPUTES it](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it)

### [MCTYM — MockTyped — every test double must DERIVE from the module it replaces](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces)

### [OBHNB — ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit)

### [BLGTN — Obligations — the duties in force, resolved from packs and config, and their status across the project](camadas/gate.md#blgtn--obligations--the-duties-in-force-resolved-from-packs-and-config-and-their-status-across-the-project)

### [OPQSP — OpenQuestions — a spec with an open question is not ready to implement](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement)

### [PGNHN — PaginationHonored — what promises a SET does not return the first page in silence](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence)

### [PHORP — PhaseOrdered — plan phases and phase dependencies must be ordered and consistent](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent)

### [PLCFL — PlaceholderFilled — the skeleton the generator emits must be FILLED IN](camadas/gate.md#plcfl--placeholderfilled--the-skeleton-the-generator-emits-must-be-filled-in)

### [PCJPL — PlanChangeJustified — a modified plan or spec must declare why it changed](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed)

### [PLRVP — PlanRevised — mutual revision visibility between superseded and revising plans](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans)

### [PSVPL — PlanSeedsValid — specifications seeded in a plan must target valid governed layers](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers)

### [PSDPL — PlanSourceDeclared — a plan that NAMES a source has to declare who builds it](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it)

### [PRFLO — Profile — the verdicts of a run, gathered per gate and per node](camadas/gate.md#prflo--profile--the-verdicts-of-a-run-gathered-per-gate-and-per-node)

### [PRHNP — ProgressHonest — the progress file tells the truth about the disk](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk)

### [PRJTS — ProjectTests — the gates read the project's tests through the source the project declares](camadas/gate.md#prjts--projecttests--the-gates-read-the-projects-tests-through-the-source-the-project-declares)

### [PRGTP — PromotableGates — identifies clean informative gates ready for promotion to blocking](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking)

### [PCBPR — ProofCrossesBoundary — when a rule claims a relation, the proof must reach the other side](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side)

### [RFRSR — RefResolves — the reference points at the spec that REALLY describes the unit](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit)

### [RPHRG — RegionPairHonored — every opened source region must close with its own identity code](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code)

### [RVMTR — ReverseMatch — every scenario still has its rule, and every proven code still has its scenario](camadas/gate.md#rvmtr--reversematch--every-scenario-still-has-its-rule-and-every-proven-code-still-has-its-scenario)

### [RVORP — RevisionOrphans — the rules a revision changed the meaning of, without saying so](camadas/gate.md#rvorp--revisionorphans--the-rules-a-revision-changed-the-meaning-of-without-saying-so)

### [RVRNR — RevisionRenumber — the revisions a branch added move to a free number when the base took theirs](camadas/gate.md#rvrnr--revisionrenumber--the-revisions-a-branch-added-move-to-a-free-number-when-the-base-took-theirs)

### [RTDCL — RouteDeclared — a screen declares how one arrives, and names its neighbours](camadas/gate.md#rtdcl--routedeclared--a-screen-declares-how-one-arrives-and-names-its-neighbours)

### [RTEXR — RouteExists — declared route in specification must exist in application route registry](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry)

### [RLUEX — Rule — the identity of a verification INSIDE a gate, and the waiver that names it](camadas/gate.md#rluex--rule--the-identity-of-a-verification-inside-a-gate-and-the-waiver-that-names-it)

### [RLIMR — RuleImplemented — a spec catalogues rules, and the code shows it realized them](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them)

### [RLTYR — RuleTypes — the rule VOCABULARY is extensible, but it must be DECLARED](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared)

### [RLUSG — RuleUses — each rule says what it uses, and what it uses exists](camadas/gate.md#rlusg--ruleuses--each-rule-says-what-it-uses-and-what-it-uses-exists)

### [SCASS — ScenarioAsserts — scenario outcome steps must assert concrete verifiable outcomes](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes)

### [SCIDS — ScenarioIdentity — two scenarios of the same feature cannot share one code](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code)

### [SCLTR — ScenarioLetterDeclared — the letter of a scenario code exists in the vocabulary](camadas/gate.md#scltr--scenarioletterdeclared--the-letter-of-a-scenario-code-exists-in-the-vocabulary)

### [STASC — ScenarioTypeAligned — scenario classification tags must match the code nature letter](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter)

### [SBGRD — SiblingGuard — sibling functions treat the same parameter consistently](camadas/gate.md#sbgrd--siblingguard--sibling-functions-treat-the-same-parameter-consistently)

### [SFMSP — SpecFeatureMatch — every requirement the spec DEFINES has at least one scenario](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario)

### [TLVCD — TestLevelCodes — each scenario references only codes its test level accepts](camadas/gate.md#tlvcd--testlevelcodes--each-scenario-references-only-codes-its-test-level-accepts)

### [TSTRT — TestTraceable — a test linked to a feature must declare what scenario it proves](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves)

### [TICTS — TestIDContract — a test handle is one contract with four ends: the code exposes it, the spec declares it, a consumer queries it](camadas/gate.md#ticts--testidcontract--a-test-handle-is-one-contract-with-four-ends-the-code-exposes-it-the-spec-declares-it-a-consumer-queries-it)

### [TQETS — TestidQueriedExists — every handle queried by an E2E flow must exist in code](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code)

### [TRCMT — TriadComplete — the pieces that realize a spec EXIST](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist)

### [TRDCT — TriggerDeclared — cited compliance triggers and obligations must exist in the declared vocabulary](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary)

### [VLANV — ValueAnchored — a replicated key is declared where it is used, and every copy carries the same value](camadas/gate.md#vlanv--valueanchored--a-replicated-key-is-declared-where-it-is-used-and-every-copy-carries-the-same-value)

### [VRBSV — VRBaseline — ensures visual regression scenarios have captured reference baseline images](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images)

## infra

### [CHRCC — ChangeRecord — the delivery record an agent leaves when it finishes a stage](camadas/infra.md#chrcc--changerecord--the-delivery-record-an-agent-leaves-when-it-finishes-a-stage)

### [CHNGL — Changelog — the releases of a project, read from its commits](camadas/infra.md#chngl--changelog--the-releases-of-a-project-read-from-its-commits)

### [CHLGC — CheckLog — the check's output mirrored to a file, so it can be reread without re-running](camadas/infra.md#chlgc--checklog--the-checks-output-mirrored-to-a-file-so-it-can-be-reread-without-re-running)

### [CDGNC — CodeGenerator — the short, stable identity code suggested for a unit's name](camadas/infra.md#cdgnc--codegenerator--the-short-stable-identity-code-suggested-for-a-units-name)

### [CFPCD — CodeFromPath — the most meaningful unique code for a unit, given its file path](camadas/infra.md#cfpcd--codefrompath--the-most-meaningful-unique-code-for-a-unit-given-its-file-path)

### [DMSTD — DaemonState — the background watcher's state files: PID, pause flag, log and meta](camadas/infra.md#dmstd--daemonstate--the-background-watchers-state-files-pid-pause-flag-log-and-meta)

### [DMRND — DaemonRuntime — how each platform probes and terminates the background watcher](camadas/infra.md#dmrnd--daemonruntime--how-each-platform-probes-and-terminates-the-background-watcher)

### [GTAVG — GitAvailability — why a git operation cannot happen, named with its fix](camadas/infra.md#gtavg--gitavailability--why-a-git-operation-cannot-happen-named-with-its-fix)

### [GTMTG — GitMeta — what git knows about the files: last commit dates, pending changes, HEAD, dirty count](camadas/infra.md#gtmtg--gitmeta--what-git-knows-about-the-files-last-commit-dates-pending-changes-head-dirty-count)

### [APRCP — ApprovalReachable — the doctor says when the required approval can never be given, and how to get out](camadas/infra.md#aprcp--approvalreachable--the-doctor-says-when-the-required-approval-can-never-be-given-and-how-to-get-out)

### [GHEGT — GitHubEnvironment — the doctor warns, before the work starts, about the pieces the GitHub flow silently needs](camadas/infra.md#ghegt--githubenvironment--the-doctor-warns-before-the-work-starts-about-the-pieces-the-github-flow-silently-needs)

### [GVOPG — GovernanceOpportunities — the doctor suggests the canonical gates and settings a project has not adopted yet](camadas/infra.md#gvopg--governanceopportunities--the-doctor-suggests-the-canonical-gates-and-settings-a-project-has-not-adopted-yet)

### [DCTRO — Doctor — the global health check that hunts the systemic loose ends of a project](camadas/infra.md#dctro--doctor--the-global-health-check-that-hunts-the-systemic-loose-ends-of-a-project)

### [PNDCP — PendingDecisions — the doctor lists the specs that still hold open decisions, the heaviest first](camadas/infra.md#pndcp--pendingdecisions--the-doctor-lists-the-specs-that-still-hold-open-decisions-the-heaviest-first)

### [SPSCS — SpecSections — the doctor tells when most specs of a layer lack a section a gate needs to see](camadas/infra.md#spscs--specsections--the-doctor-tells-when-most-specs-of-a-layer-lack-a-section-a-gate-needs-to-see)

### [APPRP — ApplyPreset — writes a stack preset's layers into the configuration and deduces one identity prefix per module](camadas/infra.md#apprp--applypreset--writes-a-stack-presets-layers-into-the-configuration-and-deduces-one-identity-prefix-per-module)

### [ARCHR — ArtifactChoice — turns the artifacts the user chose at init into artifact layers and colocation](camadas/infra.md#archr--artifactchoice--turns-the-artifacts-the-user-chose-at-init-into-artifact-layers-and-colocation)

### [BLCNB — BuildConfig — builds the configuration that inference proposes as the default for the init questions](camadas/infra.md#blcnb--buildconfig--builds-the-configuration-that-inference-proposes-as-the-default-for-the-init-questions)

### [INDCN — InitDecisions — the pure decisions of init over the proposed configuration: code layers, tags and governs rules](camadas/infra.md#indcn--initdecisions--the-pure-decisions-of-init-over-the-proposed-configuration-code-layers-tags-and-governs-rules)

### [DFGTD — DefaultGates — the gates a project is born with, by artifact and by project age, and the canonical gate catalog](camadas/infra.md#dfgtd--defaultgates--the-gates-a-project-is-born-with-by-artifact-and-by-project-age-and-the-canonical-gate-catalog)

### [BREXB — BoardExposure — hand the local board the same page and the same collect contract the pipeline publishes](camadas/infra.md#brexb--boardexposure--hand-the-local-board-the-same-page-and-the-same-collect-contract-the-pipeline-publishes)

### [GTSTG — GitState — classify the project's versioning before `init` scans it, and say what to do about it](camadas/infra.md#gtstg--gitstate--classify-the-projects-versioning-before-init-scans-it-and-say-what-to-do-about-it)

### [HDGDH — HeaderGuide — render the project's header guide in the stack's comment dialect, passing the gate that init itself declares](camadas/infra.md#hdgdh--headerguide--render-the-projects-header-guide-in-the-stacks-comment-dialect-passing-the-gate-that-init-itself-declares)

### [INPRN — InferProposal — walks the project and proposes its structure deterministically, for init to confirm](camadas/infra.md#inprn--inferproposal--walks-the-project-and-proposes-its-structure-deterministically-for-init-to-confirm)

### [INCTN — InitCatalogs — the stack preset catalog and the @TBD instruction that init seeds into judgment gates](camadas/infra.md#inctn--initcatalogs--the-stack-preset-catalog-and-the-tbd-instruction-that-init-seeds-into-judgment-gates)

### [OPDTP — OperatorDetection — tell whether a person or an AI is running `init`, and whether the discovery phase is still to be done](camadas/infra.md#opdtp--operatordetection--tell-whether-a-person-or-an-ai-is-running-init-and-whether-the-discovery-phase-is-still-to-be-done)

### [PCSDP — PackSeeding — copy the compliance packs carried in the binary into the project, never over an adapted one](camadas/infra.md#pcsdp--packseeding--copy-the-compliance-packs-carried-in-the-binary-into-the-project-never-over-an-adapted-one)

### [INQSN — InitQuestions — describe the human decisions of `init` so an agent can answer them without the terminal UI, and judge every answer](camadas/infra.md#inqsn--initquestions--describe-the-human-decisions-of-init-so-an-agent-can-answer-them-without-the-terminal-ui-and-judge-every-answer)

### [FLWRF — FlowWorkflows — declare the pipelines of the work flow, find what is missing or broken, and seed them without taking over what the team owns](camadas/infra.md#flwrf--flowworkflows--declare-the-pipelines-of-the-work-flow-find-what-is-missing-or-broken-and-seed-them-without-taking-over-what-the-team-owns)

### [GHIGT — GitHubIssues — the issue lifecycle on the repository's cards, when the project works on GitHub](camadas/infra.md#ghigt--githubissues--the-issue-lifecycle-on-the-repositorys-cards-when-the-project-works-on-github)

### [ISLFS — IssueLifecycle — a divergence recorded so it survives the session, with its state as a folder](camadas/infra.md#islfs--issuelifecycle--a-divergence-recorded-so-it-survives-the-session-with-its-state-as-a-folder)

### [OBPCB — ObligationPack — distributable sets of obligations from a norm, resolved against the project](camadas/infra.md#obpcb--obligationpack--distributable-sets-of-obligations-from-a-norm-resolved-against-the-project)

### [TSQUT — TaskQueue — the file-backed queue between "something changed" and "someone works on it"](camadas/infra.md#tsqut--taskqueue--the-file-backed-queue-between-something-changed-and-someone-works-on-it)

### [RCDLR — RecodeDialect — the project's own surfaces of a code: testID prefixes and file names](camadas/infra.md#rcdlr--recodedialect--the-projects-own-surfaces-of-a-code-testid-prefixes-and-file-names)

### [RCPLR — RecodePlan — planning and applying the rename of a code across the whole project](camadas/infra.md#rcplr--recodeplan--planning-and-applying-the-rename-of-a-code-across-the-whole-project)

### [RCRWR — RecodeRewrite — renaming an identity code inside a text, on every surface where it appears](camadas/infra.md#rcrwr--recoderewrite--renaming-an-identity-code-inside-a-text-on-every-surface-where-it-appears)

### [TXSMT — TextSimilarity — how close two texts that should be equal are, weighted by what each word discriminates](camadas/infra.md#txsmt--textsimilarity--how-close-two-texts-that-should-be-equal-are-weighted-by-what-each-word-discriminates)

### [SGSTS — SuggestionStore — a proposed fix, as a patch plus its reason, waiting for someone to decide](camadas/infra.md#sgsts--suggestionstore--a-proposed-fix-as-a-patch-plus-its-reason-waiting-for-someone-to-decide)

### [TSTLS — TestList — the project's tests, read the way the project says they are written](camadas/infra.md#tstls--testlist--the-projects-tests-read-the-way-the-project-says-they-are-written)

### [RCGRL — RuleCodeGrammar — the grammar that recognizes a scenario code in a test's name, in the project's vocabulary](camadas/infra.md#rcgrl--rulecodegrammar--the-grammar-that-recognizes-a-scenario-code-in-a-tests-name-in-the-projects-vocabulary)

### [DCLDF — DiffChangedLines — which lines of which files a change added, read from a unified diff](camadas/infra.md#dcldf--diffchangedlines--which-lines-of-which-files-a-change-added-read-from-a-unified-diff)

### [JUIJN — JUnitIngest — the run's outcome per test case, and the scenario codes each case proves, read from a JUnit report](camadas/infra.md#juijn--junitingest--the-runs-outcome-per-test-case-and-the-scenario-codes-each-case-proves-read-from-a-junit-report)

### [LCINL — LcovIngest — line coverage per file, and the uncovered lines of a change, read from an lcov report](camadas/infra.md#lcinl--lcovingest--line-coverage-per-file-and-the-uncovered-lines-of-a-change-read-from-an-lcov-report)

### [MTINM — MutationIngest — the mutation score per file, read from a Mutation Testing Elements report](camadas/infra.md#mtinm--mutationingest--the-mutation-score-per-file-read-from-a-mutation-testing-elements-report)

### [GRING — GremlinsIngest — the mutation score per file, read from a gremlins report](camadas/infra.md#gring--gremlinsingest--the-mutation-score-per-file-read-from-a-gremlins-report)

## mapa

### [GRBLG — GraphBuild — projecting the declared structure onto the scanned files: one node per file, and the relations between them](camadas/mapa.md#grblg--graphbuild--projecting-the-declared-structure-onto-the-scanned-files-one-node-per-file-and-the-relations-between-them)

### [EVFRA — EvidenceFreshness — a test's evidence expires when anything it exercises changes, not only its own file](camadas/mapa.md#evfra--evidencefreshness--a-tests-evidence-expires-when-anything-it-exercises-changes-not-only-its-own-file)

### [MPFRM — MapFormat — the map's format number decides whether this binary may read it](camadas/mapa.md#mpfrm--mapformat--the-maps-format-number-decides-whether-this-binary-may-read-it)

### [IMANM — ImpactAnalysis — what changing one file propagates to, and what it must be confronted against](camadas/mapa.md#imanm--impactanalysis--what-changing-one-file-propagates-to-and-what-it-must-be-confronted-against)

### [SGINA — SignalIngestion — hanging the runner's results on the map's nodes: executions, proven rules, coverage and mutation](camadas/mapa.md#sgina--signalingestion--hanging-the-runners-results-on-the-maps-nodes-executions-proven-rules-coverage-and-mutation)

### [MPLCK — MapLock — the map changed by one writer at a time, each applying only what it changes](camadas/mapa.md#mplck--maplock--the-map-changed-by-one-writer-at-a-time-each-applying-only-what-it-changes)

### [GRMDG — GraphModel — the shape of the map file, and when a validated relation goes stale](camadas/mapa.md#grmdg--graphmodel--the-shape-of-the-map-file-and-when-a-validated-relation-goes-stale)

### [GRQRG — GraphQueries — read-only questions over the loaded map: who governs what, neighbours, orphans, counts and a parents-first order](camadas/mapa.md#grqrg--graphqueries--read-only-questions-over-the-loaded-map-who-governs-what-neighbours-orphans-counts-and-a-parents-first-order)

### [EDSTD — EdgeStamping — recording on each relation that it was confronted, with what result, and since when](camadas/mapa.md#edstd--edgestamping--recording-on-each-relation-that-it-was-confronted-with-what-result-and-since-when)

### [GRPRG — GraphPersistence — saving and loading the map file without churn and without partial reads](camadas/mapa.md#grprg--graphpersistence--saving-and-loading-the-map-file-without-churn-and-without-partial-reads)

### [TSUNT — TestedUnits — which code a test tests, found by the project's own derivation](camadas/mapa.md#tsunt--testedunits--which-code-a-test-tests-found-by-the-projects-own-derivation)

## scan

### [SCIGS — ScanIgnore — what the scan never sees: the built-in list, the project's `.gitignore`, and editor noise](camadas/scan.md#scigs--scanignore--what-the-scan-never-sees-the-built-in-list-the-projects-gitignore-and-editor-noise)

### [PRFLP — ProgressFile — a plan's progress companion: how it is named, and how two sides of it merge](camadas/scan.md#prflp--progressfile--a-plans-progress-companion-how-it-is-named-and-how-two-sides-of-it-merge)

### [SRRGS — SourceRegion — the declared interval of a rule code in source, and the composition of test scripts](camadas/scan.md#srrgs--sourceregion--the-declared-interval-of-a-rule-code-in-source-and-the-composition-of-test-scripts)

### [RPSCR — RepoScan — the repository read as text: which files exist, of which layer, and what each declares](camadas/scan.md#rpscr--reposcan--the-repository-read-as-text-which-files-exist-of-which-layer-and-what-each-declares)

### [UPOWP — UpstreamOwnership — which files Anchors still owns in a project, and where a file's `@anchors` header begins and ends](camadas/scan.md#upowp--upstreamownership--which-files-anchors-still-owns-in-a-project-and-where-a-files-anchors-header-begins-and-ends)

