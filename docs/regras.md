<!-- anchors:generated from doct/regras.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Regras

Todas as regras do sistema, de todas as unidades. Cada uma leva ao texto na página da sua
camada.

Para ver uma camada inteira de uma vez — com a visão geral de cada unidade e os cenários —
abra a página dela em `camadas/`.

## gate

### [CDCTC — CodeCataloged — what the code EXPORTS must be in the spec, or waived in the code](camadas/gate.md#cdctc--codecataloged--what-the-code-exports-must-be-in-the-spec-or-waived-in-the-code)

### [CDLNG — CodeLanguage — the code does not go back to mixing languages](camadas/gate.md#cdlng--codelanguage--the-code-does-not-go-back-to-mixing-languages)

### [CRVCD — CodeReferenceValid — cross-referenced requirement codes must resolve to existing units](camadas/gate.md#crvcd--codereferencevalid--cross-referenced-requirement-codes-must-resolve-to-existing-units)

### [CSDCN — ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those](camadas/gate.md#csdcn--contractstatusdeclared--the-output-contract-lists-the-status-codes-the-code-really-returns-and-only-those)

### [CNHNC — CountHonored — a numerical assertion written in a spec must match reality in code](camadas/gate.md#cnhnc--counthonored--a-numerical-assertion-written-in-a-spec-must-match-reality-in-code)

### [DCRQD — DocRequired — the aggregated document the unit must feed](camadas/gate.md#dcrqd--docrequired--the-aggregated-document-the-unit-must-feed)

### [DSCDC — DocSelfContained — the spec has to stand on its own](camadas/gate.md#dscdc--docselfcontained--the-spec-has-to-stand-on-its-own)

### [DCFRD — DocsFresh — the compiled document has to reflect the spec](camadas/gate.md#dcfrd--docsfresh--the-compiled-document-has-to-reflect-the-spec)

### [DMDCD — DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid](camadas/gate.md#dmdcd--domaindeclared--the-spec-declares-what-the-unit-accepts-and-who-blocks-the-invalid)

### [EVFRV — EvidenceFresh — the score of this test holds against TODAY's code](camadas/gate.md#evfrv--evidencefresh--the-score-of-this-test-holds-against-todays-code)

### [EXCMX — ExternalCommand — executes external tools via shell passing targets as positional arguments](camadas/gate.md#excmx--externalcommand--executes-external-tools-via-shell-passing-targets-as-positional-arguments)

### [FTMFT — FeatureTestMatch — scenarios in feature must be implemented in test by code and description](camadas/gate.md#ftmft--featuretestmatch--scenarios-in-feature-must-be-implemented-in-test-by-code-and-description)

### [IDCND — IdentityConsistent — a unit's spec identity must match its exposed testID and visual baseline](camadas/gate.md#idcnd--identityconsistent--a-units-spec-identity-must-match-its-exposed-testid-and-visual-baseline)

### [LYBNL — LayerBoundary — a layer does not reach what is not its own](camadas/gate.md#lybnl--layerboundary--a-layer-does-not-reach-what-is-not-its-own)

### [MRPRM — MarkerParity — the same rule has to appear at BOTH ends that fulfil it](camadas/gate.md#mrprm--markerparity--the-same-rule-has-to-appear-at-both-ends-that-fulfil-it)

### [MCSTM — MockStamped — the double carries the mark of the snippet it replaces, and the gate RECOMPUTES it](camadas/gate.md#mcstm--mockstamped--the-double-carries-the-mark-of-the-snippet-it-replaces-and-the-gate-recomputes-it)

### [MCTYM — MockTyped — every test double must DERIVE from the module it replaces](camadas/gate.md#mctym--mocktyped--every-test-double-must-derive-from-the-module-it-replaces)

### [OBHNB — ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit](camadas/gate.md#obhnb--obligationhonored--the-cross-cutting-duty-that-lives-outside-the-unit)

### [OPQSP — OpenQuestions — a spec with an open question is not ready to implement](camadas/gate.md#opqsp--openquestions--a-spec-with-an-open-question-is-not-ready-to-implement)

### [PGNHN — PaginationHonored — what promises a SET does not return the first page in silence](camadas/gate.md#pgnhn--paginationhonored--what-promises-a-set-does-not-return-the-first-page-in-silence)

### [PHORP — PhaseOrdered — plan phases and phase dependencies must be ordered and consistent](camadas/gate.md#phorp--phaseordered--plan-phases-and-phase-dependencies-must-be-ordered-and-consistent)

### [PCJPL — PlanChangeJustified — a modified plan or spec must declare why it changed](camadas/gate.md#pcjpl--planchangejustified--a-modified-plan-or-spec-must-declare-why-it-changed)

### [PLRVP — PlanRevised — mutual revision visibility between superseded and revising plans](camadas/gate.md#plrvp--planrevised--mutual-revision-visibility-between-superseded-and-revising-plans)

### [PSVPL — PlanSeedsValid — specifications seeded in a plan must target valid governed layers](camadas/gate.md#psvpl--planseedsvalid--specifications-seeded-in-a-plan-must-target-valid-governed-layers)

### [PSDPL — PlanSourceDeclared — a plan that NAMES a source has to declare who builds it](camadas/gate.md#psdpl--plansourcedeclared--a-plan-that-names-a-source-has-to-declare-who-builds-it)

### [PRHNP — ProgressHonest — the progress file tells the truth about the disk](camadas/gate.md#prhnp--progresshonest--the-progress-file-tells-the-truth-about-the-disk)

### [PRGTP — PromotableGates — identifies clean informative gates ready for promotion to blocking](camadas/gate.md#prgtp--promotablegates--identifies-clean-informative-gates-ready-for-promotion-to-blocking)

### [PCBPR — ProofCrossesBoundary — when a rule claims a relation, the proof must reach the other side](camadas/gate.md#pcbpr--proofcrossesboundary--when-a-rule-claims-a-relation-the-proof-must-reach-the-other-side)

### [RFRSR — RefResolves — the reference points at the spec that REALLY describes the unit](camadas/gate.md#rfrsr--refresolves--the-reference-points-at-the-spec-that-really-describes-the-unit)

### [RPHRG — RegionPairHonored — every opened source region must close with its own identity code](camadas/gate.md#rphrg--regionpairhonored--every-opened-source-region-must-close-with-its-own-identity-code)

### [RTEXR — RouteExists — declared route in specification must exist in application route registry](camadas/gate.md#rtexr--routeexists--declared-route-in-specification-must-exist-in-application-route-registry)

### [RLIMR — RuleImplemented — a spec catalogues rules, and the code shows it realized them](camadas/gate.md#rlimr--ruleimplemented--a-spec-catalogues-rules-and-the-code-shows-it-realized-them)

### [RLTYR — RuleTypes — the rule VOCABULARY is extensible, but it must be DECLARED](camadas/gate.md#rltyr--ruletypes--the-rule-vocabulary-is-extensible-but-it-must-be-declared)

### [SCASS — ScenarioAsserts — scenario outcome steps must assert concrete verifiable outcomes](camadas/gate.md#scass--scenarioasserts--scenario-outcome-steps-must-assert-concrete-verifiable-outcomes)

### [SCIDS — ScenarioIdentity — two scenarios of the same feature cannot share one code](camadas/gate.md#scids--scenarioidentity--two-scenarios-of-the-same-feature-cannot-share-one-code)

### [STASC — ScenarioTypeAligned — scenario classification tags must match the code nature letter](camadas/gate.md#stasc--scenariotypealigned--scenario-classification-tags-must-match-the-code-nature-letter)

### [SFMSP — SpecFeatureMatch — every requirement the spec DEFINES has at least one scenario](camadas/gate.md#sfmsp--specfeaturematch--every-requirement-the-spec-defines-has-at-least-one-scenario)

### [TSTRT — TestTraceable — a test linked to a feature must declare what scenario it proves](camadas/gate.md#tstrt--testtraceable--a-test-linked-to-a-feature-must-declare-what-scenario-it-proves)

### [TQETS — TestidQueriedExists — every handle queried by an E2E flow must exist in code](camadas/gate.md#tqets--testidqueriedexists--every-handle-queried-by-an-e2e-flow-must-exist-in-code)

### [TRCMT — TriadComplete — the pieces that realize a spec EXIST](camadas/gate.md#trcmt--triadcomplete--the-pieces-that-realize-a-spec-exist)

### [TRDCT — TriggerDeclared — cited compliance triggers and obligations must exist in the declared vocabulary](camadas/gate.md#trdct--triggerdeclared--cited-compliance-triggers-and-obligations-must-exist-in-the-declared-vocabulary)

### [VLANV — ValueAnchored — every value of a closed set points at the rule that justifies it, and the anchor carries the value](camadas/gate.md#vlanv--valueanchored--every-value-of-a-closed-set-points-at-the-rule-that-justifies-it-and-the-anchor-carries-the-value)

### [VRBSV — VRBaseline — ensures visual regression scenarios have captured reference baseline images](camadas/gate.md#vrbsv--vrbaseline--ensures-visual-regression-scenarios-have-captured-reference-baseline-images)

