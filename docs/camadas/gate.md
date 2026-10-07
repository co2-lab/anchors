<!-- anchors:generated from doct/camadas/gate.md.tmpl — inputs:eec4a64acfe012b4 — DO NOT EDIT: run `anchors docs build` -->


# Camada: gate



> Esta camada tem 81 unidades e 1347 regras — acima do corte de 20 unidades / 2000 linhas, então esta página traz o RESUMO de cada unidade. O texto completo está na spec.


## APISP — APISpec — the coherence of an API spec, and its error codes in the code

An API spec ties three things the client relies on: which contract each body is, which status each
refusal answers with, and which error code and message come with it. Three gates ask whether the spec
holds them together, before anything is compiled or run:

| Gate | Question |
| --- | --- |
| `api-contracts-resolve` | Is every contract the body and the responses cite a spec with a Domain, and does every response say its contract? |
| `api-errors-declared` | Does every error response carry its error code and message, under a status the Responses declare? |
| `error-codes-honored` | Is every error code the spec declares one the unit's code emits? |

The OpenAPI build already fails on a contract that is no spec; asked here, per unit, the failure names
the unit while it is being written, and the pre-commit stops the commit that broke it. They run on an
API unit's main code file — the layers tagged `interface` — and read the spec beside it, when it has an
`Endpoint` section; sections and their columns are read in any language of the catalog.


- **APISP-B01** — A node that is not code and a code file whose spec is missing or has no `Endpoint` leave every gate without a verdict.

- **APISP-B02** — `api-contracts-resolve` fails naming each body or response (by status) whose contract is empty or `TODO`, each cited code no spec of the map carries, and each contract whose spec has no Domain table; `—` is a response with no body.

- **APISP-B03** — `api-errors-declared` fails naming each error response whose status the Responses do not declare — exactly, or by its range (`4xx`) —, and each one without an error code or a message; a spec with no error response leaves without a verdict.

- **APISP-B04** — `error-codes-honored` fails naming each declared error code that appears in none of the unit's code — the main file and the files its spec specifies, comments removed —, and the files it read.

- **APISP-B06** — `SetProjectSectionTitles`: The titles the project's `section_titles` give a catalog section — its own and each layer's — are found as that section too, besides the catalog's translations.

- **APISP-B05** — Sections are found under their title in any language of the catalog, and their columns by their header in those languages.

- **APISP-E01** — The spec beside the code file cannot be read, or there is none.

- **APISP-E02** — A contract's spec or one of the unit's code files cannot be read.


## BRCOV — BranchCoverage — the tests take the branches the code has

Line coverage says a line ran. A line with a condition runs whichever way the condition goes,
and the other way can stay untested for good: in the reference app a "red" state was unreachable
(the slider's maximum was below the base it had to cross), a percentage was always zero, an
empty-state skeleton never rendered — every line covered, every branch not.

The lcov format carries branches (`BRDA`) whatever the language that wrote it, and the coverage
ingestion keeps them on the node: how many, and the ones no suite took. A branch never taken on a
line where the mutation run also found mutants no test ran is reported apart as likely dead —
nothing the tests do reaches it.


- **BRCOV-B01** — A file whose share of branches taken is below the floor fails, naming the share, the floor, the branches missed out of the total and their lines in order; at or above the floor it passes. With no `min_percent` the floor is 100: every branch.

- **BRCOV-B02** — A branch on a line carrying `@no-branch: <why>`, or on the line after it, is left out of the missed; a waiver with no reason waives nothing.

- **BRCOV-B03** — A missed branch on a line where the mutation, measured at the node's revision, found a mutant no test ran fails as likely dead, naming the lines — whatever the floor. A mutation measured at another revision says nothing about these lines.

- **BRCOV-B04** — A node that is not code is skipped; a node with no coverage, or coverage of another revision, is pending; a coverage with no branch is skipped.

- **BRCOV-B05** — A file a coverage report listed with no instrumentable line is skipped, and one a whole run of its suite left out of the report is a divergence, as `line-coverage` reads them. (`coverageAbsence`)


## CDCTC — CodeCataloged — what the code EXPORTS must be in the spec, or waived in the code

Confronts a spec against the code it governs, starting from the OTHER side: **this symbol
is public — does it have a rule?**

It is the inverse of `rule-implemented`. That one starts from the spec and asks "does this
rule have code?"; this one starts from the code. Without both, the divergence escapes on
one of the sides — measured in a real project: **a spec catalogued 2 rules for 7 exported
functions, and no gate asked about the remaining 5.**

**Noise is not an argument for not building** — and this line has already been wrong here.
The earlier version said that charging a spec of every symbol would produce hundreds of
legitimate-but-useless findings, and a gate that accuses everything is switched off. The
fear was right; the conclusion was not. Compare the two outcomes: a granular gate
switched off by the project does not protect, **and the project KNOWS**; a gate too coarse
and left on does not protect either, **and it reports GREEN**. The outcome is the same;
what changes is the honesty. And there is a decisive asymmetry: a noisy gate is
CALIBRATABLE by whoever uses it, while a gate that is too coarse cannot be sharpened by
the project — the decision was taken inside and there is no way to recover it. When in
doubt between granular-with-noise and coarse-with-silence, the default is GRANULAR.

The same principle governs the language: the export pattern used to be TypeScript syntax
built in, so in a Go, Python or Ruby project it matched zero symbols and the gate reported
VERDE — stamping approval, with `blocking: true`, over what it had never read. **Without
knowing how to read, the gate goes quiet; it never approves.**


- **CDCTC-B01** — An artifact that is not a spec leaves without a verdict: the ruler starts from the spec that governs the code.

- **CDCTC-B02** — An exported symbol the spec never names FAILS, and the verdict names the orphan and its LINE — the name alone would make the reader hunt for the symbol.

- **CDCTC-B03** — What the spec already catalogues is never accused.

- **CDCTC-B04** — `@no-rule: <reason>` on the symbol waives it: not every export deserves a rule, and the waiver is what makes the noise manageable without lying.

- **CDCTC-B05** — A BARE marker, with no written reason, does not waive — it would be a silent way to quiet the gate, and the trace that a decision was taken would vanish.

- **CDCTC-B06** — A spec cataloguing every exported symbol passes.

- **CDCTC-B07** — With no code linked the gate leaves without a verdict: the absence belongs to `unit-complete`, and accusing it in both places would duplicate the debt.

- **CDCTC-B08** — Without a declared export pattern the gate SKIPS and says it skipped, naming how to enable it. It never approves what it cannot read.

- **CDCTC-B09** — With the pattern declared the gate confronts for real, in any language — the project's own `export_detect` is the ruler.

- **CDCTC-B10** — The declared dialect family also supplies the pattern: a Go project needs only name its family.

- **CDCTC-B11** — Every file the spec governs (`specifies`) is confronted, not only the first; the verdict names the files that hold orphans, and when more than one does, each orphan carries its file beside its line.

- **CDCTC-I01** — The waiver holds in the COMMENT BLOCK above the symbol, not only on its own line. The declaration is documentation: whoever writes it puts the explanation alongside, and the explanation rarely fits on one line. Looking only at `i-1` made the gate ignore the declaration and keep accusing — measured in a file where it sat on the second line of a two-line comment.

- **CDCTC-I02** — The waiver does NOT leak between symbols. A symbol with no comment above it does not inherit another symbol's declaration — inheriting would let one marker exempt the whole file, which is the opposite of what it is.

- **CDCTC-I03** — Green over what was never read is the worst possible failure in a measuring instrument. Not knowing how to read is a reason to go quiet and SAY SO, never to approve — and never silently, because the bias of the built-in ecosystem would hide inside that silence.

- **CDCTC-X01** — Does not judge whether the catalogued rule DESCRIBES the symbol well.

- **CDCTC-X02** — Does not decide which symbols deserve a rule.

- **CDCTC-X03** — Does not know any language: whoever declares what is public is the project.

- **CDCTC-X04** — Does not charge the absence of code — that belongs to `unit-complete`.

- **CDCTC-E01** — No map has been built when the gate confronts a spec.

- **CDCTC-E02** — A code file the spec `specifies` is no longer on disk.

- **CDCTC-E03** — `derived.export_detect` is declared but has no capture group.


## CDLNG — CodeLanguage — the code does not go back to mixing languages

Defends a decision that, without a gate, undoes itself: **what is WRITTEN in the code is
English; what is READ is translated.**

Anchors was born with the code in Portuguese and migrated. An external contributor should not
need Portuguese to read a function's name. But that work is lost with a
single PR from someone who does not know the rule — and the rule is nowhere the compiler
reads. This gate is where it comes to be.

**What it looks at, and what it ignores on purpose:**

- IDENTIFIERS — yes. They are what the contributor needs to read in order to work.
- COMMENTS — no. They carry the measurements and the whys, in the team's language.
- USER-FACING TEXT — no. That goes through the translation catalogue, and charging it here would
  duplicate the ruler in two places that would diverge.

The detection **is not by dictionary**, and the reason was measured: comparing each word against the
system dictionary accused compound identifiers in English, which no common dictionary
has — 491 accusations for 47 real cases, 90% false positive. A gate like that is
turned off on the first day.


- **CDLNG-B01** — An identifier in Portuguese is ACCUSED; an identifier in English passes without noise.

- **CDLNG-B02** — The verdict RETURNS the word that accused — without it, whoever reads looks for the needle in the whole file.

- **CDLNG-B03** — Only a DECLARATION is the subject: what does not declare an identifier is not read.

- **CDLNG-B04** — The declarations are truly found, in every form the language offers — not only in the most common one. `PortugueseIdentifiers` sweeps the content and returns what accused.

- **CDLNG-B05** — `WordIsPortuguese` decides ONE word, and `IdentifierIsPortuguese` decides a whole identifier by breaking it into the words that compose it — it is the separation that lets the length floor hold per word, and not for the whole identifier.

- **CDLNG-B06** — The project's own production code declares no identifier this ruler accuses: every Go file outside tests, vendored code and test data is swept, and each accusation names the identifier, the word that caused it and the file.

- **CDLNG-B07** — No gate of this package decides by confronting prose in the team's language: the text a gate matches is stable vocabulary — a marker, a code or the document's structure — so translating the project never silences a gate.

- **CDLNG-I01** — A short word does not count. Below the length floor there is no language to infer, and accusing there would be noise — which is how a gate loses the trust of whoever reads it.

- **CDLNG-X01** — Does not read a COMMENT.

- **CDLNG-X02** — Does not read USER-FACING TEXT.

- **CDLNG-X03** — Does not use a dictionary to decide the language.


## CRVCD — CodeReferenceValid — cross-referenced requirement codes must resolve to existing units

Confronts specification content against the project universe of identities: **every requirement code cited by a specification must resolve to an existing unit in the map.**

A lying anchor is the exact failure this framework exists to prevent. It takes a form that other gates miss: a specification cites an external requirement identifier — in a note, a narrative explanation, or a cross-reference — and the referenced unit does not exist anywhere in the repository. The citation simulates traceability while pointing to empty space.

A measured incident demonstrated this defect: a schema specification asserted creation of indices on 2026-08-11 and referenced four requirement codes belonging to specifications that were never created. The file described by that schema did not contain a single line implementing the entity. Yet the specification passed all gates: it possessed a code, a header, and valid section headings, and dependency checks only evaluate method symbols rather than requirement prose. Any future reader — human or autonomous agent — reads the text as an authentic historical record of completed work.

This gate enforces referential truth by extracting every requirement token shaped like an identity code followed by a requirement suffix and verifying that its owning unit exists in the project graph.

This gate operates in distinct territory from neighbouring gates:
- Unlike `dependency-honored`, which inspects dependency tables to confirm that promised Go symbols and methods appear in code, this gate validates requirement codes against the project identity universe. A unit can legitimately consume existing Go packages while citing non-existent requirement codes.
- Unlike `code-cataloged`, which ensures that exported code symbols are catalogued within their own unit specification, this gate governs external references pointing to other specifications.
- Unlike `unit-complete`, which enforces local unit structure within a single unit, this gate maintains referential integrity across the collective graph of specifications.


- **CRVCD-B01** — When the confronted node kind is not a specification, the gate skips confrontation.

- **CRVCD-B02** — When the map graph is nil, confrontation returns a pending verdict because external identities cannot be resolved.

- **CRVCD-B03** — When the map graph contains no declared identity codes, confrontation returns pending rather than approving blindly.

- **CRVCD-B04** — When all cited external requirement codes resolve to known units in the map graph, the gate passes.

- **CRVCD-B05** — When a specification cites an external requirement code whose unit is not in the map graph, the gate fails.

- **CRVCD-B06** — Citations matching the specification's own identity code are recognized as self-references and are not charged as orphans.

- **CRVCD-B07** — Multiple orphaned requirement codes are reported in alphabetical order formatted as identifiers.

- **CRVCD-B08** — Identity ownership is determined from graph node metadata and specification header comments on disk.

- **CRVCD-B09** — Text tokens that do not match requirement code structure are ignored and never charged as external references.

- **CRVCD-I01** — A specification is sovereign over its own identity code, so internal cross-references to its own requirements never fail.

- **CRVCD-I02** — An absent map graph or an empty identity universe always yields Pending, never Pass, refusing to approve what cannot be measured.

- **CRVCD-I03** — Any unresolvable external requirement citation produces a blocking Fail verdict.

- **CRVCD-X01** — Does not evaluate whether the referenced requirement behavior is implemented correctly.

- **CRVCD-X02** — Does not inspect non-specification artifacts like code or tests for dangling requirement citations.

- **CRVCD-X03** — Does not mandate that a specification must cite external requirements.

- **CRVCD-E01** — A spec the map carries without its code is no longer on disk, so its header cannot be read when the declared codes are collected.


## CTRIM — ContractImpact — a changed field names the rules that use it, and their tests

A spec's file revision says the spec changed, not WHAT changed: editing one field made every
relation of the spec stale, and nothing told which rules read that field. The rule-use sections say
what each rule reads, so a field whose row differs from the last commit names its rules — those of
this spec that use it, and those of the specs that depend on this unit's code and use a field of
that name —, and the tests that cite those rules. The change is read from git, so nothing new is
kept in the map. The same answer feeds the test selection: a test of an affected rule runs even
when its own file did not move.


- **CTRIM-B01** — A field whose row differs from the spec at HEAD, or that HEAD had and the spec no longer has, is changed; a field only the new version has is not. (`ContractImpacts`, `Impact`)

- **CTRIM-B02** — A changed field names the rules that use it — by name or first segment — in this spec and in the specs that depend on the code it governs, and the test files whose titles cite those rules; a changed field no rule uses names nothing.

- **CTRIM-B03** — `contract-impact` is a divergence with each changed field, its rules and its tests; with no impact it passes, and a node that is not a spec is skipped.

- **CTRIM-B04** — The test files every impacted rule reaches, across the specs with uncommitted changes, are listed for the test selection, which adds those its suite runs. (`ImpactedTests`)

- **CTRIM-B05** — The rules a revision added since HEAD names in `Revises:` or `Checked:` are answered: an impact whose rules are all answered is not reported, and one with a rule nobody answered still is. The impact lives only while the change is uncommitted, and the change's own revision is where whoever changed the field says they looked. (`acknowledgedRules`)


## CSDCN — ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those

Confronts the `Output Contract` table of an INTERFACE spec — a handler, a route — against the
status codes the governed code actually emits: **the table was written, but does the code
still keep it?**

Why the gate exists, measured: in an audit of 51 spec-versus-code divergences in the
reference app (2026-08), this was the MOST REPEATED pattern — eight handlers declared a
contract the code did not honour. And the omitted status was almost always the SECURITY
one: the 403 of ownership, the 409 of conflict. Whoever writes the table thinks about the
happy path and about the "business" errors, not about the refusals of access.

**The two sides of the error are different, and both matter.** A status EMITTED and not
declared leaves the client — programmed from the table — unable to handle the refusal: the
user sees a generic error where there was a specific reason. That was `accept-org-invite`,
which declared 3 status codes and emitted 8; the two missing 403s and the 409 were the
defence against invite hijacking. A status DECLARED and never emitted is worse in another
way: it is dead code in the client, and it disappears with nobody noticing.
`reanalyse-metadata` declared 402 for exceeded quota and no path of the handler emits 402
(the quota answers 429) — a client treating 402 as "needs to pay" would never fire that
branch.

What separates it from its neighbours: `dependency-honored` confronts the symbols the spec
promises to consume; this one confronts the status codes the spec promises to return. And
unlike `spec-complete`, it never charges the EXISTENCE of the section — with no table there
is nothing to confront.


- **CSDCN-B01** — A status EMITTED by the code and absent from the table fails, and the verdict names it — the client programmed from the table does not handle the refusal.

- **CSDCN-B02** — A status DECLARED in the table and emitted by no path fails too, as dead code in the client that disappears unnoticed.

- **CSDCN-B03** — A faithful table passes: the gate that only accuses is a noise generator, and nobody keeps one.

- **CSDCN-B04** — The 500 of the top-level try/catch is not charged for absence: it is infrastructure every handler carries, not a decision of this one.

- **CSDCN-B05** — A declared `5xx` covers the 5xx codes the code emits; a generic `4xx` covers nothing, because it would hide exactly the access refusals this gate hunts.

- **CSDCN-B06** — A status that only appears inside a COMMENT is not an emitted status: the comment describes what the function used to do.

- **CSDCN-B07** — Without the contract section the confrontation is skipped — charging the section's existence belongs to `spec-complete`.

- **CSDCN-B08** — Code that returns no status at all is skipped: a cron or trigger handler has `void` as its contract, and there is nothing to confront.

- **CSDCN-B09** — A status passed as a literal to a locally defined helper counts as emitted — `fail(400, …)` is a 400, wherever the envelope is built.

- **CSDCN-B10** — When the code carries a DYNAMIC status — a helper that takes the code by parameter — the gate stops asserting the phantom side, and keeps charging the literals it did find.

- **CSDCN-B11** — Without a declared `http_status` in the dialect the verdict is Pending, and names the known families — the meter does not fake conformity nor guess the stack.

- **CSDCN-B13** — The API catalog's `Responses` section (`Respostas`, `Respuestas`) is read as the output contract, by its whole title: `Error Responses` (`Respostas de Erro`) is another section.

- **CSDCN-B12** — A project that explicitly waives the `http_status` field is skipped: the opt-out is declared, and it is honoured.

- **CSDCN-I01** — The lexicon that reads the status comes from the project's dialect, never from the gate. Embedding one stack's syntax would make the gate silent on every other one — and silence reads as conformity.

- **CSDCN-I02** — A dialect declared by hand, with no family, teaches the gate its own lexicon. Without this the agnosticism would be a promise limited to the built-in families.

- **CSDCN-I03** — A named constant is worth the number it means: `http.StatusForbidden` is a declared 403, and reading only digits would approve every handler written with constants.

- **CSDCN-X01** — Does not demand the generic ranges, and does not invent their semantics.

- **CSDCN-X02** — Does not judge WHEN each status is right — only whether the number appears on both sides.

- **CSDCN-X03** — Does not charge the phantom side when the code builds the status dynamically.

- **CSDCN-E01** — No map has been built, so the gate receives no graph.

- **CSDCN-E02** — The map links the spec to code, but none of those files can be read (all gone from disk since the last build, or empty).

- **CSDCN-E03** — The project's `dialect.http_status` does not compile as a regular expression.


## CTTST — ContractTested — an API unit is proven against the project's OpenAPI document

The project's OpenAPI is compiled from the specs of its API units (`anchors docs build`), so it says
what the API promises. A contract test is what says the implementation keeps it: each language has
the tool that runs requests against an OpenAPI document and validates the answers — Schemathesis or
Dredd for any stack, kin-openapi in Go, jest-openapi in JavaScript, openapi-core in Python,
swagger-request-validator in Java. Anchors does not choose the tool; it asks for the three things that
make the proof traceable: the contract scenario `{CODE}-CT` in the unit's feature, with the project's
contract regime; a test of the unit that names `{CODE}-CT`; and that test loading the OpenAPI
document, so it validates against the compiled contract and not against a copy written in the test.

It runs on the API unit's main code file — the layers the project tags `interface` — and reads the
spec beside it; a spec with no `Endpoint` section is no API unit.


- **CTTST-B01** — A node that is not code, a code file with no spec beside it, and a spec with no `Endpoint` section or no code leave without a verdict.

- **CTTST-B02** — The feature must carry a scenario line with both `@{CODE}-CT` and the contract regime tag; otherwise the failure names both.

- **CTTST-B03** — A test of the unit — its path names the unit's code, it sits beside the unit under its name, or a folder of its path is named after the unit — must name `{CODE}-CT`; otherwise the failure says no test names it.

- **CTTST-B04** — A test naming `{CODE}-CT` must mention an OpenAPI document; otherwise the failure names the tests that do not load it.

- **CTTST-B05** — The contract regime tag is the one `derived.regimes` maps to a regime naming a contract; without one, `contract-level`.

- **CTTST-B06** — With the scenario and a contract test that loads the document, the gate passes.

- **CTTST-E01** — The spec beside the code file cannot be read, or there is none.

- **CTTST-E02** — The feature or a test file cannot be read.


## CNHNC — CountHonored — a numerical assertion written in a spec must match reality in code

Confronts numerical claims made in specs against the codebase: **does a number that a spec asserts about
the code match the actual count in code?**

The "lying anchor" defect has a numerical variant that ages completely on its own: a spec asserts "the 50
models of the product", an engineer introduces the 51st model, and the sentence instantly becomes a lie
without anyone ever touching or modifying the spec. In real project measurements, a single spec file held 7
numeric assertions, and adding ONE model rendered 10 sentences obsolete at once.

Unlike other anchor defects caused by developer carelessness or skipped steps, numerical drift depends
strictly on the passage of time. Every count written in prose represents debt accumulating compound interest.

The gate does NOT guess what to count. In the measured project, "51 models", "51 authorization clauses", and
"46 indexed models" represented three completely different questions evaluated over the exact same set of
files; no static heuristic can guess all three correctly. Instead, the spec explicitly DECLARES how to count
using the declaration contract:
- A glob pattern alone counts matching FILES (for example, counting model files).
- A glob pattern followed by a regex pattern counts OCCURRENCES matching the regex across those files (for
  example, counting authorization clauses).

Any undeclared number in prose (such as a data retention policy stating "90 days") is deliberately ignored:
the gate verifies what was explicitly contracted, not every arbitrary number appearing in text.

Crucially, the gate also confronts the PROSE situated alongside the count declaration marker, not merely the
marker itself. In real inspection, a marker declared 51 (matching code), yet a prose sentence three lines below
asserted "50 authorization clauses". Passing the marker while ignoring the prose would certify a spec that lies
to any human reader opening the document.


- **CNHNC-B01** — When the confronted node is not a spec, the gate skips confrontation.

- **CNHNC-B02** — When the spec contains no count declarations, the gate skips without asserting green.

- **CNHNC-B03** — When a declared glob expression matches matching files and count equals expected, the gate passes.

- **CNHNC-B04** — When the file count on disk differs from the expected count declared in the marker, the gate fails.

- **CNHNC-B05** — When a regex pattern is provided, the gate counts regex occurrences across files rather than file counts.

- **CNHNC-B06** — When regex occurrence count differs from the expected count declared in the marker, the gate fails.

- **CNHNC-B07** — When a declared glob matches zero files on disk, the gate fails and warns to check the path.

- **CNHNC-B08** — When a glob pattern contains invalid glob syntax, the gate fails reporting the glob error.

- **CNHNC-B09** — When a count pattern contains invalid regex syntax, the gate fails reporting the regex error.

- **CNHNC-B10** — When the marker count matches disk but adjacent prose asserts a conflicting count, the gate fails.

- **CNHNC-B11** — Non-restrictive complements attached to prose labels are confronted as claims about the total count.

- **CNHNC-B12** — Qualifying words following the label indicate subset claims and are excluded from total count confrontation.

- **CNHNC-B13** — Arbitrary numbers in prose without an explicit count declaration are ignored and skip confrontation.

- **CNHNC-B14** — Only FILES are counted: a directory the glob also matches is not a file. Counting it made `models/*` say 2 for one model and one subfolder, while the pattern mode skipped the same directory — the two modes disagreed on the same glob.

- **CNHNC-I01** — The gate only verifies explicit contracts. Numerical assertions without declaration markers never trigger a check.

- **CNHNC-I02** — Both the declaration marker and the adjacent prose must agree with code reality; a correct marker cannot mask lying prose.

- **CNHNC-I03** — Zero file matches are treated as path or glob errors rather than a genuine count of zero files.

- **CNHNC-X01** — Does not guess or infer what to count from arbitrary text in prose.

- **CNHNC-X02** — Does not charge specs that contain no count declaration markers.

- **CNHNC-X03** — Does not confront prose statements that qualify a subset rather than the total.

- **CNHNC-E01** — REF[CNHNC-B08]: a declared glob that does not parse is the configuration failure B08 answers: the gate fails carrying the glob error

- **CNHNC-E02** — REF[CNHNC-B09]: a count pattern that does not compile is the configuration failure B09 answers: the gate fails carrying the regex error

- **CNHNC-E03** — A file the glob matches cannot be read while occurrences of a pattern are counted.


## DEPHN — DependencyHonored — methods promised in the dependency table are consumed in code

Confronts a spec's **Dependency Table** against actual use in the unit's code: **every method
symbol a spec promises to consume must genuinely appear in the non-comment content of the code
it governs.**

This is the **relational gate of the spec→code edge**. It catches the divergence that sibling gates
(such as `feature-test-match`) cannot see: a spec declares `DEP2 → metadataVersioning · resolveVersion, applyEdit`,
but the unit's code calls only `resolveVersion` — the promised `applyEdit` is never invoked.
That exact divergence was the bug that slipped past 65 green tests in the originating project's
end-to-end suite.

The ruler is **static, without execution**:
- **Only symbols are confronted**: identifiers enclosed in `backticks` in the Method field. A prose
  description (such as `"CRUD + queries"` or `"requests"`) is not confrontable and is deliberately ignored;
  prose explains the relationship rather than establishing a verifiable contractual promise.
- Each declared symbol must appear as a token within the non-comment content of the code the spec
  `specifies`. Line comments are stripped before checking so that comments cannot mask unused dependencies.
  A declared symbol absent from code fails.
- When a promised symbol is missing from code, the gate checks for a near rename (e.g. `ping` → `pingHeartbeat`
  or `computeSeatsAmount` → `computeSeatsAmountCents`). If an identifier in code extends the symbol as a prefix
  or suffix and the symbol has at least four characters, the verdict actively suggests the rename instead
  of merely reporting an error.
- Undetermined and unapplicable states are explicit:
  - An artifact that is not a spec skips with reason (`not a spec — only spec has Dependency Table`).
  - When no relational graph is loaded, the verdict is undetermined (`Pending`).
  - When the Dependency Table contains no confrontable symbols, the check skips (`Dependency Table promises no confrontable SYMBOL`).
  - When the spec governs no code (`specifies` edge absent), the verdict is undetermined (`Pending`).


- **DEPHN-B01** — An artifact that is not a spec leaves without a verdict: only specs have a Dependency Table.

- **DEPHN-B02** — Without a relational map the verdict is UNDETERMINED (Pending), because dependency and specification edges cannot be traversed.

- **DEPHN-B03** — A spec declaring no confrontable symbols in its Dependency Table leaves without a verdict (Skip): prose descriptions and empty tables promise no verifiable identifiers.

- **DEPHN-B04** — A spec that specifies no code files leaves the verdict UNDETERMINED (Pending): there is no governed code to confront yet.

- **DEPHN-B05** — When every promised symbol appears in the non-comment content of the governed code, the gate passes.

- **DEPHN-B06** — When a promised symbol is absent from the governed code, the gate fails and names the unused symbol and dependency target file.

- **DEPHN-B07** — When an absent symbol resembles an identifier in code (sharing a prefix or suffix extension), the failure verdict suggests the candidate rename.

- **DEPHN-B08** — Line comments in governed code are stripped before confrontation, so symbols appearing exclusively within comments do not fulfill the promise.

- **DEPHN-I01** — Prose descriptions in dependency methods are never treated as contracts; only backticked identifiers constitute promises to verify.

- **DEPHN-I02** — Symbol presence in code is matched strictly on token word boundaries, never as a substring of a larger identifier name.

- **DEPHN-I03** — Near-symbol rename suggestions are strictly conservative, requiring prefix or suffix containment and a minimum symbol length of four characters.

- **DEPHN-X01** — Static textual confrontation without runtime execution.

- **DEPHN-X02** — Does not interpret dependency semantics, parameter signatures, or method types.

- **DEPHN-E01** — A code file the map says the spec specifies is no longer on disk.


## DCRQD — DocRequired — the aggregated document the unit must feed

Confronts a unit against the documents the project declared as a CONTRACT: **the unit was
delivered, but did it feed what it was supposed to feed?**

The Structure declares which documents are mandatory and WHEN each one must be touched —
a layer triggers a duty. That resolution already existed, and two places consulted it: the
one that informs whoever picks up the card, and the one that lists the duties. **Neither
confronted.**

What that cost, measured: a delivery shipped a new unit — with spec, code, feature and
test — without touching either the interface contract or the data schema. Both were
declared mandatory for that layer, both were untouched, and every check went green. It
only surfaced because two agents collided on the same card and there was something to
compare: one branch carried 267 lines the other did not. **Without the collision nobody
would have noticed** — there is nothing to compare when only one piece of work exists.

**Two questions, and the second is the one that matters:** does the document EXIST, and
does it MENTION this unit? The first alone would approve an empty file created to silence
the gate.


- **DCRQD-B01** — A mandatory document that DOES NOT EXIST fails.

- **DCRQD-B02** — A document that exists and does not MENTION the unit fails too — existence alone would approve an empty file created to silence the gate.

- **DCRQD-B03** — A mention by the unit's identity CODE counts as documented.

- **DCRQD-B04** — A mention by the FILE NAME also counts: the document speaks of the unit either way.

- **DCRQD-B05** — Satisfying one of two declared duties is not enough — each document is charged on its own.

- **DCRQD-B06** — Without a declaration in the Structure nothing is charged: the ruler is what the project committed to, not what one supposes it owes.

- **DCRQD-B07** — A layer with no trigger declared is not charged, even when the document exists.

- **DCRQD-B08** — Aggregated, the verdict is ONE PER DOCUMENT, not one per unit.

- **DCRQD-B09** — In a document that has sections, the unit counts as documented only when a section TITLE names it, by code or by file name: a mention in the body of another section, or in a note, does not count, nor does the title of another unit whose name contains this one's; a document with no section still counts by mention.

- **DCRQD-B10** — Sections and headings are read only in a Markdown document (`.md`, `.markdown`, `.mdx`); in any other — a YAML, an OpenAPI, a script — a line opening with `#` is a comment, the document has no sections, and a mention of the unit in its body counts.

- **DCRQD-I01** — The duty starts from the SPEC, not from the code. The spec is what declares the unit; starting from the code would charge the duty of a file that merely realises it.

- **DCRQD-I02** — The layer used is the UNIT's, not the node's. A spec node always has layer `spec`; reading it would charge every spec of the project the same duty, or none.

- **DCRQD-I03** — Without a map the aggregated verdict is skipped, never approved. Approving without being able to look would stamp what was not measured.

- **DCRQD-X01** — Does not understand the document's CONTENT.

- **DCRQD-X02** — Does not decide WHICH documents are mandatory.

- **DCRQD-E01** — REF[DCRQD-I03]: with no map built the aggregate cannot look, and I03 answers it: skipped, never approved

- **DCRQD-E02** — REF[DCRQD-B01]: a mandatory document that cannot be read is answered as the missing one of B01 — either way the reader is sent to that file

- **DCRQD-E03** — REF[DCRQD-B06]: with no configuration there is no declared duty, and B06 answers that nothing is charged


## DSCDC — DocSelfContained — the spec has to stand on its own

Confronts a spec against the reader who does not have the repository open: **the reference
points somewhere else — does it bring what it points at?**

The spec is read by two audiences, and one of them has no checkout. The compiled `docs/*.md`
inherits the spec's text word for word, and whoever reads the documentation to learn what the
system does has no use for the name of a plan file — navigating to it is exactly what the
documentation mechanism exists to eliminate.

**What the gate does NOT charge.** The reference WITH the text alongside is right and stays:
the reader has the argument in hand. What is left over is the SCAFFOLD — the sentence that
announces the quotation and then does not quote — and it is the scaffold that sends the
person away. Measured in a real project: 48 mentions across 37 specs, almost all of them
followed by the quoted passage.

**No vocabulary, anywhere.** The gate does not look for "plan", "see" or "per": Anchors
governs projects in any language, and a gate that matches words passes in silence over the
project written in the other one — which is worse than not existing, because the spec then
LOOKS protected. What it matches is STRUCTURE, and only what Anchors itself defines: the
PATH of a map node written in the body, and the `{CODE}-R000N` form that is the doctrine's
revision identity.

Unlike `docs-fresh`, this one is informative and offers no command that fixes it: rewriting
a sentence is the work of whoever wrote it, and blocking a commit over a question of form
would stop the flow. The gate marks, and the correction rides along with the card that
already touches the spec.


- **DSCDC-B01** — A reference that comes WITH the passage it announces passes: the reader has the argument in hand and does not leave the page.

- **DSCDC-B02** — A reference that only POINTS is accused, and the finding names the line and shows it, so the correction does not need a hunt.

- **DSCDC-B03** — A quotation counts in any written tradition — straight quotes, typographic ones, guillemets, German low quotes, CJK corner brackets.

- **DSCDC-B04** — A path inside a code fence is not a reference: there it is an EXAMPLE, a command to run, not a sentence sending the reader away.

- **DSCDC-B05** — The spec citing its OWN path is identifying itself, not sending anybody anywhere, and is not accused.

- **DSCDC-B06** — A rule code is not a revision: the `-R000N` form is the one the doctrine reserves, and accusing the other would charge the spec for naming its own subject.

- **DSCDC-B07** — Only the spec is charged; every other kind leaves without a verdict.

- **DSCDC-B08** — A revision cited with an EXPLANATION on the same line passes: the code is the label and the sentence is the content.

- **DSCDC-B09** — With no map the confrontation is skipped, because the list of what counts as a reference comes from the map and nowhere else.

- **DSCDC-I01** — The ruler matches STRUCTURE, never vocabulary. A spec written in English, Spanish, German or Japanese is charged by exactly the same code, because no language changes a path.

- **DSCDC-I02** — The on-line-explanation escape belongs to the REVISION and not to the PATH. A revision code is a label whose sentence is the content; a path is the place the person would have to go, and no amount of surrounding prose says what is there.

- **DSCDC-I03** — The verdict names the LINE and shows what it says. A gate that accuses without pointing hands the diagnostic work to whoever reads it.

- **DSCDC-X01** — Does not judge whether the accompanying content is FAITHFUL to what the reference announces.

- **DSCDC-X02** — Does not offer a command that fixes it, and does not block.

- **DSCDC-X03** — Errs on the side of letting things through when measuring whether a line explains something.

- **DSCDC-E01** — REF[DSCDC-B09]: with no map the list of references cannot be known, and B09 answers that the confrontation is skipped


## DCCVD — DocsCovered — every spec must reach some page of the compiled documentation

The `docs-covered` gate asks whether a spec reaches ANY page of the documentation the project
compiles from its templates. Its sibling, `docs-fresh`, asks whether each page that exists still
reflects the spec it came from; it cannot see a spec that no page asks for.

That is the silent defect this gate exists for. The templates select specs by filter (by layer,
for instance), so a spec that no filter selects compiles to nowhere, and nothing complains: every
page that exists is correct, and the unit has a complete unit and passes every relational gate.
The first spec of a new layer drops out of the documentation without a word, exactly when the
project grows and nobody is watching.

The gate is informative: the fix is a person's decision, to widen a template's filter or to add
the layer's page. Rebuilding the documentation does not fix it, because the build only produces
the pages the templates ask for.

The answer is a property of the whole set of specs, so it is computed once per project root and
map and reused for every spec of the scan. The verdict of each spec, however, is about that spec
only.


- **DCCVD-B01** — An artifact that is not a spec is skipped: only a spec is compiled into documentation.

- **DCCVD-B02** — A project with no templates directory is skipped: demanding documentation from a project that declared none would invent a duty.

- **DCCVD-B03** — A spec that no template reaches fails, and the verdict names the spec.

- **DCCVD-B04** — A spec that some template reaches passes.

- **DCCVD-B05** — When the templates do not compile, the gate skips with no message: the compile failure is the sibling gate's finding, reported there with the compiler's own words, and repeating it here would look like a second defect.

- **DCCVD-B06** — The set of unreached specs is computed once per project root and map, and reused for every spec of the same scan; a different map recomputes it.

- **DCCVD-I01** — The verdict of a spec depends only on whether THAT spec is reached; another spec being an orphan never fails it.

- **DCCVD-X01** — Does not judge whether a page is built or up to date; it only asks whether some template selects the spec.

- **DCCVD-E01** — REF[DCCVD-B05]: the templates failing to compile is the one failure the gate handles, and B05 states its answer: skip, and leave the report to the sibling gate


## DCFRD — DocsFresh — the compiled document has to reflect the spec

Confronts the compiled documentation against the specs that feed it: **the page is there
and it reads well, but is it still saying what the spec says today?**

The documentation is GENERATED. The templates reference passages of the specs, and the
build command produces the pages. The content lives in the spec and only there — the copy
that lives in the compiled output is derived, and derived ages.

Without a gate the failure mode is silent and known: someone changes a rule in the spec,
forgets to recompile, and the documentation goes on asserting the OLD rule. Nobody sees it,
because the page is there, well formed, with real content. That is worse than an empty
page — an obviously incomplete document sends the reader looking for the source; an
out-of-date document convinces.

**INFORMATIVE by explicit decision**: a gate that blocks something a single command fixes
on its own spends a reviewer's attention on machine work. The pipeline runs the build on
merge; this gate exists so the author sees it before.

**Anchored on the SPEC** because the spec is what changes. The compiled document is not a
node of the map — it is generated, and charging a review of a compiler's output would be
charging a review of a compiler's output — so there is nowhere to start but the source.


- **DCFRD-B01** — A compiled document that no longer matches what the templates would produce now FAILS — the spec changed and nobody recompiled.

- **DCFRD-B02** — The verdict NAMES the stale documents and where they live, so the author does not have to diff the whole output directory.

- **DCFRD-B03** — A compiled document that matches what the templates produce passes.

- **DCFRD-B04** — An artifact that is not a spec leaves without a verdict: only the source of the documentation is confronted.

- **DCFRD-B05** — A project with no template directory is not charged — it never opted into the compiled-documentation mechanism.

- **DCFRD-B06** — A declared template whose document was never produced counts as stale by absence, not as satisfied.

- **DCFRD-B07** — A document written by hand, with no generation marker, is NOT charged: the build refuses to overwrite it, and demanding a command that changes nothing is a warning nobody can act on.

- **DCFRD-B08** — A template that cannot be compiled fails, carrying the compiler's own error — a broken template is a defect, not a reason to go quiet.

- **DCFRD-B09** — A project whose specs cannot be read leaves without a verdict, carrying the reason: the gate could not look, and a gate that could not look must not approve.

- **DCFRD-B10** — The comparison happens IN MEMORY. The gate never writes the compiled document.

- **DCFRD-I01** — The gate never repairs what it points at. If it compiled the document, the second run would always pass and the defect would only surface for whoever cloned the repository.

- **DCFRD-I02** — The charge starts from the SPEC and never from the compiled document. The document is a compiler output and not a node of the map; charging a review of it would charge a review of a build artifact.

- **DCFRD-X01** — Does not judge whether the compiled document is GOOD.

- **DCFRD-X02** — Does not run the build, even knowing the fix.

- **DCFRD-X03** — Does not charge documents written by hand.

- **DCFRD-X04** — Does not confront the compiled document as an artifact of its own.

- **DCFRD-E01** — REF[DCFRD-B05]: a missing template directory is the project that never opted in, answered by B05

- **DCFRD-E02** — REF[DCFRD-B08]: a template that fails to compile is answered by B08: it fails with the compiler error

- **DCFRD-E03** — REF[DCFRD-B09]: specs that cannot be read are answered by B09: no verdict, with the reason


## DCTRN — Doctrine — the vertical axis: product doctrine exists, is realized, and is never copied

Co-location ties a unit to one directory, but a business rule that holds for three screens belongs
to none of them. It lives in a product doctrine file, and each spec that concretises it points at the
doctrine rule with a `@realizes` tag. This unit holds the five gates of that vertical axis; each one
catches a silence the others cannot see.

- `plan-doctrine-exists` — a plan that seeds a doctrine promises it will exist. Without this gate the
  promise carries no charge, because the specs the plan also seeds pass every gate with nothing to
  realize.
- `doctrine-realized` — a doctrine rule that no spec realizes is a decision that never reached the
  code, and the doctrine file itself is well-formed, so nothing else looks at that far end.
- `spec-doctrine-exists` — the axis's reference check. A `@realizes` that does not resolve looks
  like traceability while the map gains no edge, typically after a doctrine rule was renamed.
- `doctrine-not-duplicated` — the defect the axis exists to eliminate: a spec that copies the
  doctrine text instead of pointing at it. The ruler is similarity, not equality, because whoever
  copies almost always changes a word.
- `spec-realizes-doctrine` — a layer may declare that every rule of its specs must say which product
  decision it concretises. It is off by default: most rules of a tool like this one are local
  mechanics, while in a product application the proportion inverts, and only the project knows which
  case it is.

The axis has one declared way out, `@TBD: <reason>`, and it means "not yet", never "never": on a
line it turns a charge into debt, which keeps showing up (as Pending) until somebody pays it. A bare
marker, or a marker quoted in backticks by a text that explains it, defers nothing.

An open question of a doctrine (a `-Q` item) names a decision nobody has taken, so it is never a
rule to realize.


- **DCTRN-B01** — A `@TBD` marker defers what its line declares only when it carries a written reason and is not inside backticks.

- **DCTRN-B02** — Without a map, `doctrine-realized`, `spec-doctrine-exists` and `doctrine-not-duplicated` are Pending: they read edges, and there are none to read.

- **DCTRN-B03** — An artifact that is not a plan is skipped.

- **DCTRN-B04** — A plan whose cited doctrines all exist on disk passes.

- **DCTRN-B05** — A plan citing doctrines that do not exist fails, naming each missing doctrine once, in order, with their count.

- **DCTRN-B06** — A missing doctrine cited on a line deferred with `@TBD` is a divergence, naming it, instead of failing.

- **DCTRN-B07** — Only a doctrine path with a directory seeds a doctrine; a bare file name in prose, or a template file, seeds nothing, and a plan that seeds nothing is skipped.

- **DCTRN-B08** — An artifact that is not a product doctrine is skipped.

- **DCTRN-B09** — A doctrine that catalogues no rule is skipped.

- **DCTRN-B10** — A doctrine whose every rule has an incoming realizes edge naming it passes.

- **DCTRN-B11** — A doctrine with unrealized rules fails, naming only the unrealized ones.

- **DCTRN-B12** — When every unrealized rule is deferred with `@TBD` on its own line, the doctrine is a divergence, naming them.

- **DCTRN-B13** — An open question (`-Q`) is not a rule, and is never charged for a realizer.

- **DCTRN-B14** — An artifact that is not a spec is skipped.

- **DCTRN-B15** — A spec that declares no `@realizes` outside the lines deferred with `@TBD` is skipped.

- **DCTRN-B16** — A declared rule passes when a realizes edge of the spec names it, or when a doctrine the spec's edges reach catalogues it.

- **DCTRN-B17** — A declared rule that no reached doctrine catalogues fails, naming only the unresolved rules.

- **DCTRN-B18** — An artifact that is not a spec is skipped.

- **DCTRN-B19** — A spec with no readable realized doctrine text is skipped.

- **DCTRN-B20** — When the rules of both sides together are fewer than four, the gate is Pending: the similarity weights cannot tell shared words apart; from four rules on, the pair is measured.

- **DCTRN-B21** — A spec rule whose text copies the doctrine rule it realizes fails, naming the spec rule, the doctrine rule and the similarity as a percentage.

- **DCTRN-B22** — A near copy, with a word changed or dropped, fails as well.

- **DCTRN-B23** — A spec rule whose text says what is specific to its unit passes.

- **DCTRN-B24** — A copy on a line deferred with `@TBD` is not charged: the wording is still being worked out.

- **DCTRN-B25** — An artifact that is not a spec, or a project with no configuration, is skipped.

- **DCTRN-B26** — A spec whose layer does not require doctrine is skipped.

- **DCTRN-B27** — Where the layer requires doctrine, a rule with no `@realizes` fails, naming it.

- **DCTRN-B28** — A rule that declares what it realizes passes, whether the tag sits on the rule's line or on the lines right below it.

- **DCTRN-B29** — A rule deferred with `@TBD` is a divergence, not failure.

- **DCTRN-B30** — A blank line ends the rule a tag belongs to: a `@realizes` after a blank line declares nothing for the rule above.

- **DCTRN-B31** — The demanding layer is resolved from the target the spec describes: the layers of the code its specifies edges reach, or, before that code exists, the layer the target's path would have.

- **DCTRN-B32** — A catalogued rule and a `@realizes` citation are read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, a 7-character rule and citation are read.

- **DCTRN-B33** — Each rule is charged once, at its first line — a Rule uses row naming it does not repeat it —, and only a rule to realize is charged: an open question (`Q`), a plan phase and a flag scenario (`G`) realize no doctrine, the letters `rule-uses-declared` leaves out.

- **DCTRN-I01** — Only a realizes edge realizes a doctrine rule: an edge of any other type that carries a rule code never counts.

- **DCTRN-X01** — Does not call a pair a copy because the two share a rare word: a pair judged similar must also reach a similarity of one half.

- **DCTRN-E01** — A doctrine that a realizes edge of the spec reaches cannot be read.

- **DCTRN-E02** — REF[DCTRN-B19]: a realized doctrine that cannot be read contributes no text, and B19 answers the spec left with none

- **DCTRN-E03** — REF[DCTRN-B05]: a cited doctrine whose file cannot be found is the missing doctrine B05 charges


## DMDCD — DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid

Confronts a spec against the question the rest of the framework does not ask: **what does this
unit accept as input, and who guarantees that the invalid never arrives?**

The gap it closes was measured: 71% of the specs of a real project had a section of
rules or effects, and only 14% said what the unit accepts. The whole framework is
built on cataloguing EFFECTS — what the unit does —, and every edge defect
found in three rounds of adversarial review lived in what nobody had declared.

The distinction that gives the gate its reason: `## Constraints` says what the unit does NOT do, and pushes
the duty OUTWARD; `## Domain` says what it ACCEPTS, and NAMES WHO is left with it.
Writing more constraints closes nothing — it creates orphans, because every "not mine" needs
someone on the other side.


- **DMDCD-B01** — An artifact that is not a spec leaves the confrontation without a verdict: the gate has no jurisdiction over code, test or guide.

- **DMDCD-B02** — A spec WITHOUT the domain section FAILS. The absence is not silence: it is the assertion not made.

- **DMDCD-B03** — A spec with the section OPENED and EMPTY fails too — opening the title without declaring anything is the same hole wearing the appearance of compliance.

- **DMDCD-B04** — Every declared input must name WHO guarantees it. An input without an owner fails, and the verdict names which ones were left orphaned.

- **DMDCD-B05** — The waiver is DECLARED and with a written reason. Whoever has no external input records that in the spec, and the gate goes quiet — but the trace that someone looked remains.

- **DMDCD-B06** — A line filled only with a pending marker is not a declaration: the untouched mould asserts nothing.

- **DMDCD-B07** — An input whose owner is named passes — it is the other side of the same ruler, and what makes it satisfiable.

- **DMDCD-I01** — The waiver requires a REASON. A bare waiver mark does not silence the gate — silence without a why is what it exists to prevent.

- **DMDCD-I02** — The failing verdict NAMES what is wrong — which input was left without an owner, or that the section is missing. A gate that fails without saying what transfers the diagnostic work to whoever reads it.

- **DMDCD-X01** — Does not judge whether the declared input is RIGHT — only whether it exists and has an owner.

- **DMDCD-X02** — Does not confront the code to check whether the validation in fact exists.

- **DMDCD-E01** — REF[DMDCD-B02]: the one handled path is the absent domain section, and B02 answers it: the spec fails


## DUPLC — Duplication — no code file holds a block copied from somewhere else

The native check behind `no-duplication`. It runs jscpd once per scan and reads its JSON report,
which lists every clone with both files and their lines, instead of reading jscpd's exit code. jscpd
exits 0 with any duplication unless a `threshold` is configured, so a gate that read the exit code
approved what it did not measure: in the reference app, 35 clones (1.17% of the lines), exit 0, and
`check --all` announced the gate "clean, ready to become blocking".

Each file gets its own verdict. A file that holds a copy fails, naming the other side of each clone
and the lines; every other file passes. With `--changed` only the files of the change are judged,
so a commit answers for the clones it touches; jscpd still scans the whole project, because a new
copy only exists relative to its original.

The calibration is the project's `.jscpd.json`, which jscpd reads by itself (`minLines`,
`ignore`, `jscpd:ignore-start` markers). Its `threshold` keeps jscpd's meaning: the percentage of
duplicated lines the project tolerates.


- **DUPLC-B01** — A file that takes part in no clone passes.

- **DUPLC-B02** — A file that takes part in a clone fails, naming for each clone its own lines, the other file and that file's lines, whichever side of the clone it is on.

- **DUPLC-B03** — A clone inside a single file is described by its two line ranges alone.

- **DUPLC-B04** — When `.jscpd.json` declares a `threshold` and the project's duplicated percentage is at or under it, the file's clones are reported as Diverge, not failed; over it they fail.

- **DUPLC-B05** — Without a declared `threshold`, any clone fails, whatever jscpd's exit code.

- **DUPLC-B06** — jscpd runs once per scan: the same root and map reuse the report, and a new map runs it again.

- **DUPLC-B07** — An absolute path in the report is read relative to the project root.

- **DUPLC-B08** — The gate runs a pinned jscpd release, not whatever is newest: a release that cannot run where the project runs does not silently stop the measuring. (`jscpdPackage`, `duplicationCommand`)

- **DUPLC-E01** — jscpd writes no report, or one that is not JSON


## ENVDC — EnvDeclared — the environment variables a unit reads are the ones its spec declares

An environment variable is a contract with whoever deploys. A variable the code reads and no spec names
is found missing in production; a variable the spec names and the code no longer reads is configured
forever for nothing. Like `contract-status-declared` with the status codes, this gate confronts both
sides: the `Environment Variables` section of the spec against the reads in the code it specifies. A
variable declared deprecated may stop being read — it is on its way out.

The code is read with the dialect's `env_read`: the project's own pattern, its language family's
(`os.Getenv`, `process.env`, `os.environ`, `System.getenv`, `ENV[...]`, `getenv`,
`Environment.GetEnvironmentVariable`), or every family's together when the project declares none. A
read by a computed name is out of reach of the text.


- **ENVDC-B01** — A node that is not a spec, and a spec with no variable declared whose code reads none, leave without a verdict.

- **ENVDC-B02** — Each variable the code reads and the spec does not declare is named.

- **ENVDC-B03** — Each variable the spec declares and the code does not read is named — unless it is declared deprecated (`yes`, or `yes: use X`).

- **ENVDC-B04** — A read inside a comment is no read.

- **ENVDC-B05** — With no language family declared, the reads of every family are recognised.

- **ENVDC-E01** — A specified code file cannot be read.


## EVFRV — EvidenceFresh — the score of this test holds against TODAY's code

Every other gate confronts text against text. This one confronts an EXECUTION against the
revisions it measured: **the test went green, and the stamp records the revision of every
node in its closure at that moment. If any of them moved, the score is still written and
has stopped being true.**

It is the failure that makes no noise. A test that breaks screams in the runner; a score
that ages stays green in yesterday's report, and disappears inside the conclusion of
whoever reads it. Measured in the reference app: `utils/login.yaml` is composed into 290
scripts — touching it invalidated 290 measurements, and the signal of not one of them
changed.

**The ruler for when NOT to judge matters as much as the one for when to fail.** A test with
no stamp has no score to expire: that is absence of proof, which is a different debt, with a
different name and a different fix — run it the first time, not revalidate it. Reporting
them together is the defect the edge-level `stale` has: 30 "never validated" drowned in
1419 "revision advanced", and the list stops being read.

It is blocking because the fix is known, local and cheap: run the test. There is no domain
judgement to make — unlike an undeclared letter, where both ways out are legitimate and the
choice belongs to whoever knows the product.


- **EVFRV-B01** — An artifact that is not a test leaves the confrontation without a verdict: only a test carries an execution score.

- **EVFRV-B02** — Without a built map the gate stays quiet: there is no closure to walk, and it will not approve what it could not look at.

- **EVFRV-B03** — A test with NO execution stamp is skipped, never failed — a score that was never written cannot have expired.

- **EVFRV-B04** — A test whose closure is intact PASSES.

- **EVFRV-B05** — The passing verdict says WHAT it checked against: the size of the closure it walked is the difference between "nobody looked" and "I looked and it stands".

- **EVFRV-B06** — A test whose DEPENDENCY advanced a revision FAILS: the score is still written and stopped holding.

- **EVFRV-B07** — The failing verdict NAMES the culprit, so whoever fixes it knows what moved underneath.

- **EVFRV-B08** — The failing verdict states the FIX — run the test again — because a verdict that accuses without saying what to do transfers the work to the reader.

- **EVFRV-B09** — A test whose OWN file changed since the run is reported as such, separately from the closure.

- **EVFRV-B10** — The culprit list is TRUNCATED at five, and the remainder is counted.

- **EVFRV-I01** — Absence of proof and expired proof are NEVER reported as the same finding. They are different debts with different fixes, and merging them is what drowns the list until it stops being read.

- **EVFRV-I02** — A test that never ran is never approved either. The skip is silence, not a stamp: approving would state a freshness nobody measured.

- **EVFRV-I03** — Truncation never hides the size of the problem: what is not listed is COUNTED. Whoever fixes it runs the test once, regardless of how many dependencies moved — but they still learn how far it spread.

- **EVFRV-X01** — Does not charge the ABSENCE of a green test.

- **EVFRV-X02** — Does not RUN the test, nor judge whether the change actually broke it.

- **EVFRV-X03** — Does not read the project's configuration.

- **EVFRV-E01** — REF[EVFRV-B02]: with no map there is no closure to walk, and B02 answers that the gate stays quiet


## EXMCH — ExamplesMatch — every Examples row is a case its tests run

A scenario outline says "for each of these rows"; the parameterised test that proves it runs its
own table, and nothing kept the two together. In the reference app they drifted apart in eight
features: the Examples said `mesada` and the test `allowance`, "Crédito" against "Cartão de
crédito", and a row the code no longer had at all.

The gate reads no test library. Every value of a row must appear in the body of a test that cites
the scenario's code, as a whole token and with its case. The body is read as `test-has-assertion`
reads it. Examples often show what the screen says while the test uses a key: a column whose
header carries `(label)`, in any supported language, is display text and is not looked for.


- **EXMCH-B01** — A row whose values do not all appear in the tests citing its scenario's code fails, naming the code, the row's line and the missing values; a feature whose rows all appear passes. The tests of every citing test join: a row may be in any of them.

- **EXMCH-B02** — A value appears only as a whole token with its case: not glued to a letter, a digit or an underscore on either side. Surrounding quotes or backticks in the cell are not part of the value, and an empty cell asks for nothing. (`tokenIn`)

- **EXMCH-B03** — A column whose header carries `(label)` in any supported language (`(rótulo)`, `(etiqueta)`) is display text and is not looked for. (`isLabelColumn`)

- **EXMCH-B04** — The rows are those under each examples table's header, of the scenario whose coded tag precedes its title, in any Gherkin language; a table ends at the first line that is neither a row, blank nor a comment, a cell past the header's columns is looked for, and a scenario with no code has no rows. (`exampleTables`)

- **EXMCH-B05** — A test's body is read as `test-has-assertion` reads it: with `labels: true` on this gate an empty test stands for the block around it, and a script's `end` bounds it.

- **EXMCH-B06** — A node that is not a feature, a feature with no coded outline with a row, a project with no tests source, and a feature no test of which cites an outline's code are skipped; an outline no test cites is left to `scenario-coverage`.

- **EXMCH-E01** — The tests cannot be listed (a script that fails, output outside the contract)


## EXCMX — ExternalCommand — executes external tools via shell passing targets as positional arguments

Executes an external command (such as jest, eslint, or tsc) defined in `anchors.yaml` against targets without reimplementing the tool, reading its exit code. To prevent command injection, target paths are never interpolated into the shell command string; instead, they are passed as positional arguments (`$1`, `$2`, `"$@"`) via `sh -c`. When target arguments exceed command line length limits (6000 bytes on Windows, 100000 bytes on Unix, or configurable via `ANCHORS_ARGV_MAX`), targets are partitioned into batches and executed sequentially, combining any failures.


- **EXCMX-B01** — An external command exiting with status zero returns `Pass` with empty detail.

- **EXCMX-B02** — An external command exiting with non-zero status returns `Fail` containing the trimmed output.

- **EXCMX-B03** — If an external command fails with empty output, the failure detail reports that the gate produced no output along with the execution error.

- **EXCMX-B04** — A gate's command for a single node runs through `RunExternalArgs`, with the node ID as the single target.

- **EXCMX-B05** — The placeholder `{{file}}` in the command template is rewritten to `"$1"` before shell invocation.

- **EXCMX-B06** — The placeholder `{{files}}` in the command template is rewritten to `"$@"` before shell invocation.

- **EXCMX-B07** — When targets are empty, execution runs exactly once without positional target arguments representing project scope.

- **EXCMX-B08** — Targets fitting within the command line budget run in a single shell execution.

- **EXCMX-B09** — Targets exceeding the budget are partitioned across multiple batches with none dropped or duplicated.

- **EXCMX-B10** — When targets are partitioned, failure in any batch causes the entire gate to fail combining failure outputs.

- **EXCMX-B11** — For a single target, failure output exceeding 500 characters is truncated with a truncation marker.

- **EXCMX-B12** — For batch or project executions, failure output exceeding 4000 characters is truncated with a truncation marker.

- **EXCMX-B13** — Environment variable `ANCHORS_ARGV_MAX` overrides the default argv limit when set to a positive integer.

- **EXCMX-B14** — Under `--index`, a target whose file on disk differs from the index is handed to the command as a copy of what the index holds, and the copy's path is shown back as the project's in what the command printed; a target the index does not have is not handed over; a target the same in both goes as it is. (`indexedTargets`)

- **EXCMX-I01** — Target paths are passed strictly as positional argv arguments and never interpolated directly into the shell script string.

- **EXCMX-I02** — A single target whose path length exceeds the limit is assigned its own batch and never silenced or dropped.

- **EXCMX-I03** — On Windows runtime without environment override the default target limit is 6000 bytes, while Unix defaults to 100000 bytes.

- **EXCMX-X01** — Does not parse or interpret linter or test tool diagnostics beyond reading stdout and stderr.

- **EXCMX-X02** — Does not aggregate cross-file state across partitioned batches.

- **EXCMX-E01** — REF[EXCMX-B02]: a command that exits non-zero is answered by B02: Fail with its trimmed output

- **EXCMX-E02** — REF[EXCMX-B03]: a command that fails with no output is answered by B03: the detail says it produced none, with the execution error


## FLRAI — Failure — the failure a spec declares must be handled, recorded, and every handling declared

A spec catalogues how its unit fails, with its own letter (`-E`) and a table of condition and
result. Until these gates, nothing confronted that catalogue with the code: a declared failure
crossed the whole pipeline without anyone asking whether the code handles it, whether the handling
records it, or whether a handling that exists answers any declared failure at all.

Three gates close that static layer, and none depends on production or on a log format:

- `failure-handled` asks whether the governed code has any path that handles a failure, when the
  spec declares failures;
- `failure-logged` asks whether that handling also records the occurrence;
- `failure-declared` is the inverse of the first, and catches the commonest case: somebody wrote a
  defence and never declared what it prevents.

Handling is not "having a catch". What a check that refuses, a recovered panic and a caught
exception share is the effect, not the syntax, so the project's dialect says what a handling and a
record look like in its language. Without those patterns the gates have measured nothing and say
so. The code is judged as a set: the code does not cite the failure's code, so no rule can be tied
to one specific path, and the gates assert the case that matters, a unit that declares failures
with no handling at all.

A failure can leave the charge only with knowledge written beside it: `@resilient: <reason>` says
"it happens, I know why, and the flow absorbs it". A bare marker exempts nothing. The same rows
carry the conclusions the observation layer reads: `@resilient` and `@observing: <what was ruled
out>`, each with its reason read whole.

A unit whose handling matches are all normal flow (a lazy map initialisation, a pattern that did
not match) closes its failure section with `none — <why>`, and that satisfies `failure-declared`.


- **FLRAI-B01** — Every failure gate skips an artifact that is not a spec.

- **FLRAI-B02** — A failure rule is read in the three catalogued forms — heading, table row, bold bullet; a failure code cited in prose declares nothing.

- **FLRAI-B03** — A spec that declares no failure is skipped by `failure-handled` and `failure-logged`.

- **FLRAI-B04** — Without handling patterns in the dialect every failure gate is Pending, and `failure-logged` is also Pending without recording patterns: a Pass would stamp what was never measured.

- **FLRAI-B05** — Without governed code that can be read — no map, no specifies edge, or no specified file on disk — every failure gate is Pending.

- **FLRAI-B06** — The governed code is read without its comment lines, so a handling written only in a comment does not count.

- **FLRAI-B07** — A spec that declares failures, over code with no handling path at all, fails, naming the charged failures in order.

- **FLRAI-B08** — Any handling path in the governed code passes the gate.

- **FLRAI-B09** — A failure marked `@resilient` with a written reason is not charged by `failure-handled` nor by `failure-logged`; when every declared failure is resilient, both pass.

- **FLRAI-B10** — A bare `@resilient` marker, with no reason, exempts nothing.

- **FLRAI-B11** — Code that handles but records nothing fails, naming the charged failures.

- **FLRAI-B12** — Code that handles and records the occurrence passes.

- **FLRAI-B13** — Code with handling paths and a spec that declares no failure fails.

- **FLRAI-B14** — Code with no handling path is skipped: there is no defence to declare.

- **FLRAI-B15** — Code with handling paths and a spec that declares at least one failure passes.

- **FLRAI-B16** — A spec whose failure section opens with `none` and a written reason passes; a bare `none`, a `none` outside the section, or a section that opens with prose does not close it.

- **FLRAI-B17** — For each declared failure (`FailureConclusions`), the reasons of `@resilient` and `@observing` are read whole; a failure with neither carries no conclusion.

- **FLRAI-B18** — A conclusion's reason ends at its table cell: the next column is never read into it.

- **FLRAI-B19** — A failure rule and its conclusion are read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, a 7-character `-E` rule is a declared failure.

- **FLRAI-B20** — `failure-declared`: each fallible source of a unit — a call of its code that a `dialect.fallible_patterns` entry recognises (comment lines and trailing comments do not count) — is named by a declared failure (`-E`): its row, or its row of the rules' uses, cites the name called as a whole word; a failure that does not name it ("not found") does not answer it. Unanswered, the gate fails naming each source by file and line, unless the Errors section is closed with `none — <reason>` or the spec waives with `@no-failure: <reason>`. (`fallibleSources`, `fallibleCalls`, `FallibleCall`, `FallibleSource`, `uncoveredSources`)

- **FLRAI-B21** — A dependency the spec declares on a file of a layer the project marks `fallible: true` is a fallible source too, named by its `DEPn` and path, and answered by a failure citing its `DEPn` or the file's stem (`useBudget`); a dependency on any other layer is not a source.

- **FLRAI-B22** — `failure-handled`: each fallible call of the unit's code has its handling in its window — the call's statement from its first line (a destructuring above it that reads the error; the climb stops at a blank line or at one ending a statement, `;`, `}` or `)`), and the pattern's `window` after it (`DefaultFallibleWindow` when it declares none) —, matched by the pattern's `handled` or, when it declares none, the project's `handle_patterns`; a call with none fails, named by file, line and text, unless its line or the line above waives it with `@no-handle: <reason>`. (`unhandledCalls`)

- **FLRAI-I01** — `failure-handled` and `failure-logged` never both fail the same spec: with no handling the first charges and the second steps aside, and with handling that records nothing it is the other way round.

- **FLRAI-X01** — Does not tie a declared failure to a specific handling path: one handling path answers for every declared failure.

- **FLRAI-E01** — REF[FLRAI-B05]: a specified file that cannot be read is left out of the governed code, and when none can be read B05 answers with Pending


## FTMFT — FeatureTestMatch — scenarios in feature must be implemented in test by code and description

Confronts the relational edge between a feature file and the test file that executes it: **every scenario
catalogued in the feature must be implemented in the linked test.**

The gate enforces a critical **TWO-TIER VERDICT**:
1. **BY CODE (Defect / Failure)**: Every scenario code (`XXXXX-Y##`) must appear in the non-comment body
   of at least one linked test. A missing scenario code is a DEFECT indicating that an agent or developer
   skipped the scenario entirely or renamed the code without updating the test. When code is missing,
   the gate returns a blocking **Fail**.
2. **BY DESCRIPTION (Signal / Warning)**: When the scenario code exists in the test, the gate checks whether
   the test description corresponds to the scenario title. Free-form text naturally varies between human
   domain language in Gherkin ("when description and value are identical") and code identifiers in tests
   (`classifica`, `expect`). Therefore, descriptive drift is treated as an informative SIGNAL rather than
   a hard failure: the gate issues a non-blocking **Pending** warning so teams can realign descriptions
   without blocking delivery of verified code.

Neighbouring gates govern other dimensions of the unit: `unit-complete` checks the physical presence
of the unit artifacts; `spec-feature-match` ensures requirements defined in specs are covered by feature
scenarios. `feature-test-match` specifically guarantees that documented scenarios are faithfully realized
in automated test suites.


- **FTMFT-B01** — An artifact that is not a feature file skips confrontation without a verdict.

- **FTMFT-B02** — When the dependency graph is nil, the gate returns Pending without approving.

- **FTMFT-B03** — A feature with no scenarios declared skips confrontation.

- **FTMFT-B04** — When no tests are linked by `tested-by` edges, the gate returns Pending.

- **FTMFT-B05** — Scenarios scoped to non-test surfaces (such as `e2e` or `vr`) are skipped by this gate.

- **FTMFT-B06** — A scenario whose code is completely absent from the test body fails.

- **FTMFT-B07** — A scenario code appearing only within comments does not count as implemented and fails.

- **FTMFT-B08** — When scenario codes match and titles are identical, the gate passes.

- **FTMFT-B09** — When scenario codes match but test descriptions drift, the gate issues a divergence.

- **FTMFT-B10** — Test titles containing quotes or nested delimiters are parsed completely without truncation.

- **FTMFT-B11** — Multiple sibling scenario codes cited in a single test title are extracted and attributed cleanly.

- **FTMFT-B12** — A test title shared among sibling scenarios does not require strict single-title equality.

- **FTMFT-B13** — Test comments count towards descriptive coverage while remaining excluded from code presence.

- **FTMFT-B14** — Scenario codes match on exact boundaries so prefix substrings of other codes do not collide.

- **FTMFT-B15** — Go `t.Run` test bindings are recognized alongside standard test runner declarations.

- **FTMFT-B16** — Line comments using `#` and `--` syntax are stripped when evaluating code presence.

- **FTMFT-B17** — The exported function `RootCode` strips scenario sub-indices (`#01`) and returns the root requirement code.

- **FTMFT-B18** — Only a regime tag the project MAPS under `regimes:` (or a canonical regime name) exempts a scenario as belonging to another surface; an unmapped tag that merely looks like a regime (`@nivel-compilacao`) leaves the scenario confronted.

- **FTMFT-B19** — A trailing comment marker counts only OUTSIDE quotes, and `--` or `#` only after whitespace: `"--label"`, `"https://…"` and a `i--` keep the rest of their line. Cut anywhere, the flag argument hid the symbol after it and dependency-honored accused a dependency the code uses.

- **FTMFT-B20** — The similarity verdict named beside each drifting scenario (`FTMFT-B09`) is written in the project's language, through i18n (`divergent` in English, `divergente` in Portuguese).

- **FTMFT-B21** — A test's title is read from the project's tests source — its `dialect.tests` pattern or script, or its family's — so how a test opens (`t.Run`, `it.each(table)`, `.only`) is the project's declaration, and a later test citing the same code is not taken for an earlier one's proof. Without a source no title is read, and the description is confronted with the test's body.

- **FTMFT-B22** — A support file linked to a feature is not among the tests its scenarios are confronted with.

- **FTMFT-B23** — A scenario tag may carry a `#nn` suffix that gives each scenario of one rule its own identity: every code on the tag line keeps its suffix, and a code with no suffix is read as before.

- **FTMFT-B24** — Every test that leads with a scenario's code is compared with the scenario, not only the first: another test whose title DIVERGES from it is named as a warning — it cites the code and talks about something else —, while a merely similar one is not, since several tests of a rule name its variations.

- **FTMFT-B26** — A visual-regression code may carry the state it captures — `BUTTN-VR-S01` — and is read whole: two VR scenarios of different states are two codes, and `scenario-identity` does not take them for one.

- **FTMFT-B25** — A test title matches its scenario when its words begin with the scenario title's words, in order — case and punctuation aside —, optionally followed by detail; the placeholders of a parameterised title — a table-driven test's `%s`, `%d`, `$name`, `${name}`, an outline's `<name>` — are words of neither title. The same words in another order, or fewer, still diverge. (`titleCovers`)

- **FTMFT-E01** — A test the map links by `tested-by` is gone from disk (or cannot be read).

- **FTMFT-E02** — The project's tests source fails, or answers outside its contract.

- **FTMFT-I01** — Absence of scenario code in test implementation is always a Failure.

- **FTMFT-I02** — Descriptive divergence with valid code presence is always a Warning (Diverge), never a Failure.

- **FTMFT-I03** — Comments are strictly separated: excluded when verifying code implementation, included when verifying semantic description.

- **FTMFT-X01** — Does not execute test suites or inspect test execution results.

- **FTMFT-X02** — Does not confront E2E or visual regression scenarios.

- **FTMFT-X03** — Does not fail tests for minor natural language variations in test titles.


## FXIXX — Fix — the self-healer that applies the mechanical, safe repairs of `check --fix`

Some findings have a repair that is mechanical and safe: writing the right date into the
`updated_at` field of a header needs no judgement, only git. `anchors check --fix` applies those
repairs, and this unit is what it calls. It holds a registry of fixers, one per check; a check with
no registered fixer is only reported, never repaired.

Today the registry holds one fixer, for the `updated-at-atual` check. It rewrites only the date
inside an existing field: the right date is today when the file has an uncommitted edit, and the
date of its last commit otherwise. When there is no way to know the right date (no git, no commit
and no edit) the file is left alone, because a guessed date would be worse than a stale one.

The repair is driven by the gate's scope, not by its verdict: every node the gate applies to is
handed to the fixer, and a file whose date is already right comes back unchanged, so it is not
reported. Each repair written, or attempted and failed, is returned for the caller to print.


- **FXIXX-B01** — Only a check with a registered fixer is fixable (`Fixable`); today that is `updated-at-atual` alone.

- **FXIXX-B02** — A stale `updated_at` on a committed file with no pending edit is rewritten to the date of the file's last commit, and the repair is reported as fixed, naming the gate and the file.

- **FXIXX-B03** — A file with an uncommitted edit takes today's date.

- **FXIXX-B04** — A date that already matches is left alone, and nothing is reported for that file.

- **FXIXX-B05** — Only gates with a fixer are run, on every node the gate applies to whatever its verdict (the fixer decides whether there is anything to correct), and only for files present on disk.

- **FXIXX-B06** — Outside a git repository the file is left untouched: there is no right date to find.

- **FXIXX-B07** — A file that was never committed and has no pending edit is left untouched: there is nothing to compare against.

- **FXIXX-B08** — The detail of each repair (fixed, or a write that failed) is written in the project's language, through i18n.

- **FXIXX-B09** — `--fix` writes the header a governed file lacks, at its top after any shebang: the `ref:` of the units the map ties it to (`UnitCodesOf`), or the `layer:` of a guide, a document or a test support file; to a header at the top with no identity it adds that line below `@anchors`, changing nothing written; a file with an identity, with no unit and no such layer, binary, or an executable script is left as it is. (`fixMissingHeader`)

- **FXIXX-B10** — `--fix` gives a file whose header has an identity (or gets one) and no `code:` of its own a code generated from its name and its type, unique among the codes in the map, written below `@anchors` beside the identity; a header with its own code is left as it is. (`fixMissingHeader`)

- **FXIXX-I01** — A repair replaces only the date inside the field; every other byte of the file stays as it was.

- **FXIXX-X01** — Does not create a missing `updated_at` field; it only corrects the value of one that exists.

- **FXIXX-E01** — Writing the repaired file fails.

- **FXIXX-E02** — REF[FXIXX-B05]: a node whose file cannot be read is one of the files B05 leaves out, and nothing is reported for it


## FLSCF — FlagScenarios — the scenarios a feature flag declares are written, complete, cited and tested

A feature flag multiplies the paths of the code without multiplying the spec. A check of a flag's
value creates two behaviours, and the spec usually describes one of them, or both mixed in a sentence
that does not say which holds when. The cost shows up as silences no other gate sees: in review nobody
knows which branch is current and which is dying, the test covers the ON path and leaves the OFF path
unproven, and the "temporary" flag grows old until removing it is archaeology.

A flag file declares its scenarios in a table: a code with the `G` letter, the condition on the value,
and what holds then. Five gates confront that declaration, each answering one question:

- flag-scenario-grammar: is each condition written in the fixed grammar?
- flag-scenarios-complete: does the flag say what happens when it is absent?
- flag-scenario-exists: does every scenario a spec cites exist?
- flag-scenario-governs: is every scenario cited by some rule?
- flag-covered: does every scenario have a test, and did that test pass?

What is deliberately not read is the flag's real value: Anchors has no access to the flag service,
must not have, and the value changes per user and per minute. The gates confront the declared scenarios.


- **FLSCF-B01** — The four gates that confront a flag skip a node that is not a flag, and skip a flag that declares no scenario; the citation gate skips a node that is not a spec.

- **FLSCF-B02** — A flag whose every condition is in the grammar passes.

- **FLSCF-B03** — A condition the grammar refuses fails the gate, and the message names the scenario, its line and the parse error.

- **FLSCF-B04** — A flag that declares a scenario for the absent value passes.

- **FLSCF-B05** — A flag with no absent scenario fails, and the message names the waiver that would close it.

- **FLSCF-B06** — A written waiver for the absent case, with a reason, passes; a bare marker, or one quoted in backticks, does not waive.

- **FLSCF-B07** — A spec with no citation of a flag scenario skips; only a citation of a `G` code counts as one.

- **FLSCF-B08** — A spec whose every cited scenario is declared by some flag of the project passes.

- **FLSCF-B09** — Citations of scenarios no flag declares fail, and the message names each unknown code once, sorted.

- **FLSCF-B10** — Without a map the gate answers Pending, because it cannot know who cites.

- **FLSCF-B11** — Each scenario needs an incoming citation that names it; the scenarios nobody cites fail the gate, named, and the cited ones are not accused.

- **FLSCF-B12** — A scenario whose outcome carries a written waiver of governance, with a reason, is not charged; a bare marker does not waive.

- **FLSCF-B13** — Without a map the gate answers Pending, because it cannot know what was proven.

- **FLSCF-B14** — Coverage is judged per scenario: a scenario among the flag's ingested proven codes is green, and a flag whose every scenario is green passes.

- **FLSCF-B15** — A scenario that no test names fails as having no test; a code that appears only in a comment of a test does not count as written.

- **FLSCF-B16** — A scenario a test names but that is not proven fails with a different message: written but not ingested when no execution was ingested, written and not passing when it was.

- **FLSCF-B17** — A `@gated-by` citation is read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, a citation of a 7-character scenario is confronted.

- **FLSCF-B18** — When the project says how its tests are written, a flag scenario is written in a file the source lists tests in only when a test TITLE cites its code; without a source, or in a file the source lists no test in, a code anywhere in the file outside comments counts.

- **FLSCF-B19** — A flag scenario is green only when a case proved exactly it: a variant `#NN` of it proves the variant, and a case of another variant does not stand for it. (`checkFlagCovered`)

- **FLSCF-I01** — A gate that could not measure never answers Pass: without a map, governance and coverage are Pending.

- **FLSCF-X01** — No waiver is accepted without a written reason.

- **FLSCF-E01** — REF[FLSCF-B03]: a condition the parser refuses is the failure the grammar gate reports, naming the scenario and the parse error

- **FLSCF-E02** — The `flags/` folder, or a flag file in it, cannot be read when a spec cites a flag scenario.

- **FLSCF-E03** — The project's tests source fails, or answers outside its contract


## GTENG — GateEngine — which gates reach which node, and what the run concludes

This is the ENGINE. It answers two questions and nothing else: **which gates apply to
which node**, and **what a run over a set of nodes concludes**. Every ruler belongs to
somebody else — the checker measures, the registry routes the name, the rule unit names
the verification and reads the waiver. What the engine owns is routing and aggregation,
and both have measured failure modes of their own.

**Routing.** A gate reaches a node when the node's kind is in what the gate declares, AND
the node carries none of the excluded labels, AND — when the gate names labels — the node
carries at least one of them, AND — when the gate demands a mark in the content — the
target's text contains it. The ORDER inside that conjunction is the decision, not the
implementation: **exclusion comes BEFORE the positive filter**. Layers carry transversal
labels alongside their own, so a node that matches both a named label and an excluded one
must stay OUT. If the positive filter won, declaring the exception would have no effect
exactly where it matters.

The content filter exists for a measured reason of its own: without it, a judgment gate
declared over specs queues an AI question for EVERY spec in the project, and the pending
counter stops measuring the work and starts measuring the size of the repository. And an
UNREADABLE target does not apply: better to stop charging than to charge blind against a
target whose content nobody knows.

**Aggregation, and the four verdicts that are not two.** A run concludes with a
promotion decision, and the engine's job is to keep four distinct states from collapsing
into pass/fail:

- a gate whose required BINARY is missing steps aside — never fails. The gate did not
  measure, and "I did not measure" is neither "clean" nor "dirty". Failing would say the
  project violated something when what is missing is the tool.
- a JUDGMENT gate does not compute at all. It asks whether somebody already answered, by
  reading the stamp a judgement run left in the map. Without that lookup the gate asked
  forever: the verdict was written, the next run asked again, the pending counter never
  came down, and work already done was invisible.
- a pending item that says "there is a decision STILL TO TAKE" bars promotion; a pending
  item that says "I had nothing to confront" does not. Measured in a real repository:
  treating them alike failed 411 nodes at once, and 410 of them were gates with no signal
  ingested.
- a debt that is ASSUMED — a known duty with a declared moment to be paid — becomes
  recorded work; the other pending items become no issue at all.

**What separates this unit from its neighbours.** The registry knows which function
answers a name; this one never looks inside a checker. The rule unit knows how a
verification is named and whether it was waived; this one only ASKS it, per target, and
turns the answer into a verdict that carries the written reason. And the waiver by target
is applied HERE rather than by filtering the gate out of the list, because a waiver
restricted to codes must leave the gate RUNNING to confront everybody else.


- **GTENG-B01** — A gate reaches a node only when the node's kind is among the kinds the gate declares.

- **GTENG-B02** — A gate that names labels reaches only the nodes carrying at least one of them.

- **GTENG-B03** — A gate that excludes labels never reaches a node carrying one, and ONE excluded label is enough.

- **GTENG-B04** — Exclusion wins over the positive label filter, because layers carry transversal labels and the exception must hold exactly where it matters.

- **GTENG-B05** — A gate demanding a mark in the content reaches only the targets whose text carries it — otherwise the pending counter measures the size of the project, not the work.

- **GTENG-B25** — The same answer — kinds, labels, exclusions, required mark — is offered to the runs that choose what to measure for a gate, so a run never measures a node the gate would not confront. (`Applies`)

- **GTENG-B26** — The gates read the project's files through the source set with `SetFileSource` — the git index in the commit hook —, and from the tree when none is set; the function it returns restores the previous source.

- **GTENG-B27** — Whether a file has changes not yet committed is asked through the source set with `SetChangedSource` — the staged changes in the commit hook —, and of the tree when none is set; `updated-at-current` reads it, so a file edited and not staged is judged by its last commit.

- **GTENG-B28** — A gate that breaks while measuring fails on that target, saying it broke and that the defect is Anchors'; the other gates and targets are still measured. (`runOne`, `runAggregate`, `recoverGate`)

- **GTENG-B29** — A gate whose `presupposes` names a configuration field the project does not declare is Pending on every target, naming the field and the opt-out, and a judgment gate queues no question; when every missing field is in `dialect.opt_out` it is Skip; with all declared it runs as usual. (`presupposedMissing`)

- **GTENG-B06** — An UNREADABLE target does not apply: better to stop charging than to charge blind.

- **GTENG-B07** — A gate with no applicable target does not run at all, so a commit touching one document does not fire a whole-project check.

- **GTENG-B08** — A gate whose required binary is absent steps aside and never fails — the gate did not measure, and the missing piece is the tool.

- **GTENG-B09** — A waiver naming THIS target spares only this node, with the written reason in the verdict, and the gate keeps confronting the others.

- **GTENG-B10** — An aggregate-scope gate runs ONCE and reports a single verdict against the scope itself, not against one of the files.

- **GTENG-B11** — A judgment gate emits a request for judgment when nothing has answered it yet.

- **GTENG-B12** — A judgment gate reads the stamp left by an earlier judgement and turns it into the verdict, so work already done stops being invisible.

- **GTENG-B13** — A judgement recorded as waived becomes Skip and never Pass, because the gate did not measure and Pass would assert an approval nobody gave.

- **GTENG-B14** — A gate declaring neither a command nor a check answers undetermined, naming the omission.

- **GTENG-B15** — Each verdict level does what the gate's `severity` says — `block` bars the promotion, `inform` reports without barring, `ignore` only counts —; unset, a blocking gate blocks its failures and informs its divergences and pending items — an open decision too —, and an informative gate informs. A divergence the project declared (`[declared]`: a debt with its deadline, a rule `@TBD`), or one its own threshold accepts (`[advisory]`: a mutation score above the floor), informs and never bars. "Nothing to confront" is not a level: it is Skip. (`markSeverity`, `Result.Blocks`, `Result.Ignored`)

- **GTENG-B16** — Only the obligations gate produces ASSUMED DEBT, and only that pending item carries a deadline into the record.

- **GTENG-B17** — The engine reconfigures the code grammar from the project's vocabulary before running anything.

- **GTENG-B18** — `Run` is the entry point that confronts the gates with the map alone, delegating with no Structure — the relational checkers then read absence as "no mapping declared".

- **GTENG-B19** — `RunWithConfig` is the same entry point carrying the Structure, which the relational checkers need to read the regimes and the surfaces of the unit.

- **GTENG-B20** — `RunFull` is the one that also knows whether the sweep is the WHOLE project, which is the only thing that lets a gate able to sweep on its own run ONCE instead of receiving thousands of targets in batches.

- **GTENG-B21** — `RunWithWaiver` is the one that honours a waiver BY TARGET, which cannot be served by filtering the gate out of the list.

- **GTENG-B22** — A vendored file — a pipeline Anchors seeded that still carries its template marker — is out of every internal ruler, whose unit, header and identity live upstream; an external command still reaches it, because what the file does in this repository is the project's concern.

- **GTENG-B23** — A gate declared with `run:` executes the custom runner, even when the canonical declaration specifies `check:`.

- **GTENG-B24** — A target the gate's `no_signal` declares is skipped, naming the declared reason, and the gate's check does not run on it.

- **GTENG-I01** — A waiver by target is applied while the gate RUNS, never by removing the gate from the list. Removing it would erase the ruler for the whole repository, and a defect elsewhere would pass along.

- **GTENG-I02** — Skip, Pending and Fail are three different answers and never collapse into two. Each says something the others do not: the gate does not apply, the gate could not measure, the gate measured and the target failed.

- **GTENG-I03** — The verdict of a waived target is Skip WITH the reason written, never silence. It leaves the failure tally without leaving the report — that is the difference between waiving and hiding.

- **GTENG-I04** — The reported target of an aggregate gate is the SCOPE, never one of the files. Blaming one of many files for a verdict about the set would be a statement the engine cannot support.

- **GTENG-I05** — Under an index source (`--index`), no internal gate reads the tree: every project file a gate reads comes through `readFile`, so a commit is judged by what it records whatever the tree holds.

- **GTENG-X01** — Does not decide WHETHER a target is correct.

- **GTENG-X02** — Does not compute the verdict of a judgment gate.

- **GTENG-X03** — Does not invent a map, a configuration or a waiver when it receives none.

- **GTENG-E01** — REF[GTENG-B06]: a target that cannot be read is answered by B06: the gate does not apply to it

- **GTENG-E02** — REF[GTENG-B08]: a required binary missing from the PATH is answered by B08: the gate steps aside


## HDLYD — HeaderLayerDeclared — the layer a header declares is one the Estrutura has

The header's `layer:` decides which templates derive the unit's siblings and where the
documentation files it. A layer the Estrutura does not have is read as nothing: the unit falls back
to the default templates, lands in a page of its own, and shows up — if at all — as a layer
outside every container on the architecture page. In the reference app a section declared
`layer: landing-component` for months while the Estrutura only had `landing-feature`, and it
surfaced by chance.


- **HDLYD-B01** — A header layer the Estrutura declares passes; one it does not fails, naming the layer and the declared layers in order. Only the `@anchors` header is read: a `layer:` line in the body is not a declaration.

- **HDLYD-B02** — A header layer that differs from a declared one only in case fails naming both: the lookup is exact, so the unit is read as unknown.

- **HDLYD-B03** — A file whose header declares no layer, a placeholder layer (`TODO…`), and a project with no layers are skipped.


## IDCND — IdentityConsistent — a unit's spec identity must match its exposed testID and visual baseline

Confronts a unit's declared identity against the external surfaces where it reappears: **does the spec
code agree with the testID prefix exposed in code and the visual regression baseline filename?**

The spec code is the canonical identity of the unit. Other identity gates verify mere PRESENCE
(whether a spec declares a code, or whether a scenario carries a code), never CONCORDANCE across surfaces.
Under presence-only checking, a single unit could declare a spec code of `BGET`, expose `testID=":bdge-screen"`,
and store its visual regression baseline as `BDEDX-VR-*.png` — three conflicting identities for the exact same
unit — while every pipeline gate remains green because each gate inspects its own surface in isolation.

The REAL COST of this discordance is not aesthetic. The project code dictionary consumed by spellcheck is
AUTOMATICALLY GENERATED from the map. When an acronym exists only inside a testID, it is missing from the
generated dictionary. Spellcheck immediately flags the acronym as a typo, and the natural developer
workaround — adding the rogue acronym to the manual spelling dictionary — CRISTALLIZES the divergence.
The symptom becomes permanent vocabulary while the underlying defect disappears from sight. This exact
failure mode occurred in the reference application when `bdgc` and `bdge` entered the manual dictionary,
motivating the creation of this gate.

The gate exercises crucial DISCERNMENT regarding what does NOT constitute divergence: a component may
legitimately use the identity code of ANOTHER unit as its testID prefix when the testID designates WHERE
the element appears. For instance, `SpendingMonthCard` (`SMCD`) exposes `:home-spending-card` because it
lives inside `HomeScreen` (`HOME`), which is how end-to-end flows locate the element starting from the screen.
Demanding strict local identity unification there would break end-to-end navigation flows.

The rule distinguishing the two cases is precise:
1. A prefix that matches the code of ANY unit declared in the map represents a KNOWN IDENTITY (deliberate reuse).
2. A prefix shaped like an identity code (4-5 letters) that belongs to NO unit in the map is an ORPHAN IDENTITY.
   Because it does not exist in the map, it cannot exist in the generated dictionary, making it the exact case
   that forces rogue acronyms into manual spelling dictionaries. Only orphan identities are rejected.

For visual regression baselines (`<Unit>.<CODE>-VR-<variant>.png`), cross-unit delegation is not permitted:
a baseline is the physical proof of THIS specific unit, not a pointer to where it appears.


- **IDCND-B01** — When the confronted node is not a spec, the gate skips confrontation.

- **IDCND-B02** — When the graph is nil, the gate returns Pending without approving.

- **IDCND-B03** — When the spec has no declared code, the gate skips without duplicating findings from the code presence gate.

- **IDCND-B04** — A testID prefix matching the spec's own identity code passes.

- **IDCND-B05** — An orphan testID prefix with code shape that belongs to no unit in the map fails.

- **IDCND-B06** — A testID prefix matching the code of another declared map unit is accepted as legitimate container reuse.

- **IDCND-B07** — A visual regression baseline whose code differs from the spec code fails.

- **IDCND-B08** — Short testID prefixes of three letters or fewer do not have code shape and pass without accusation.

- **IDCND-B09** — The failure verdict cites the conflicting acronyms and their origin files.

- **IDCND-B10** — When no orphan testID prefixes or baseline discrepancies exist, the gate passes.

- **IDCND-B11** — A baseline of one of the unit's rules (`<Unit>.<CODE>-B04-VR-<variant>.png`) is read by its unit code, the part before the first hyphen: it passes when that is the spec's code and fails when it is another unit's.

- **IDCND-B12** — A testID prefix that `anchors.renames.yaml` records as renamed to this unit's code, or to another unit's, is the same identity under its old name and passes — a testID is a contract with the E2E runner and its flows; a baseline passes under an old code only when it was renamed to this unit's.

- **IDCND-I01** — Without a map graph the gate never approves. It returns Pending because concordance cannot be evaluated without the global code inventory.

- **IDCND-I02** — Cross-unit reuse is allowed only for testIDs, never for visual regression baselines. The baseline must prove this specific unit.

- **IDCND-I03** — Absence of a spec code is never double-charged. It skips here so that the dedicated code presence gate reports the single defect.

- **IDCND-X01** — Does not charge testID prefixes of three letters or fewer.

- **IDCND-X02** — Does not charge specs for code absence.

- **IDCND-X03** — Does not forbid components from referencing parent screen codes in testIDs.

- **IDCND-E01** — A file the spec governs (`specifies`) is gone from disk when its testIDs are read.

- **IDCND-E02** — The spec lives under a directory whose name holds glob metacharacters (a Next.js `[slug]` route).


## INCHN — InternalChecks — the registry that routes a declared check name to a function

A project declares its gates in the Structure, and a gate that the CLI answers itself
says only a NAME — the check it wants. Something has to turn that name into a function,
and that something is this unit: the registry of internal checkers, plus the small
checkers whose whole answer is reading text.

**Why it is a registry and not a switch.** The checkers do not all have the same
signature, and the difference is not style. One reads only the content of the target.
One also needs the project ROOT, because it has to invoke the version control system to
answer. One needs the GRAPH, because the question crosses the unit — a feature against
the test linked to it. One needs the declaring GATE itself, because the question is
parameterised: a generic gate only knows what to look for after reading its own
configuration, and a project declares several instances of it. Four registries, and the
routing tries them in that order.

**The measured defect this shape exists to prevent**, and it is the one that costs most:
a name that does not resolve must answer PENDING, never Pass. A checker that is declared
and does not exist has not measured anything, and "I did not measure" is neither "it is
clean" nor "it is dirty". Approving here would stamp green over a verification that never
ran, and the project would read coverage where there is none. The same reasoning drives
the aggregate path: a batch or project scope checker that does not resolve answers
Pending too, naming the check that failed to route.

**Two routing paths, and the difference is not a detail.** The per-node path READS THE
TARGET FILE and fails when the read fails — a checker of content with no content has
nothing to answer. The aggregate path does NOT: its scope is the SET, so there is no one
file to read. It hands the checker an empty node and lets it orient itself by root and
configuration. Trying to read a file there would return a read error, and the gate would
fail over a file that never existed.

**What separates this unit from its neighbours.** The engine decides WHICH gates apply to
which node and what the whole run concludes. The rule unit decides how a verification is
named and whether it was waived. This one decides only WHICH FUNCTION answers, and holds
the small checkers whose entire ruler is the text in front of them: the file is not
empty, the file carries an identity code, the file carries a conformant header, the guide
distils its rules into verifiable points.

**The grammar of codes belongs to the project, not to the engine.** The letters that name
a rule type are declared in the Structure, and every pattern in this package that depends
on them has to be reconfigured together. A pattern added without registering it there
stays frozen on the canonical letters — and a scenario written with a letter the project
declared becomes invisible to that gate, which then reports green over what it never
looked at.


- **INCHN-B01** — A declared check name routes to the function registered under it.

- **INCHN-B02** — A name that does NOT resolve answers undetermined, never approval — a checker that never ran has not measured anything.

- **INCHN-B03** — The routing tries the relational registry first, then the one that needs the root, then the pure content one — so a checker that grows a dependency changes registry without changing name.

- **INCHN-B04** — The per-node path reads the target file, and a read that fails is a failure: a checker of content with no content has nothing to answer.

- **INCHN-B05** — The aggregate path does NOT read any file: its scope is the set, so the checker receives an empty node and orients itself by the root and the configuration.

- **INCHN-B06** — A name that does not resolve in the aggregate path answers undetermined too, and the report names the check that failed to route.

- **INCHN-B07** — An aggregate checker that needs the DECLARING GATE receives it, because a parameterised gate only knows what to look for after reading its own configuration.

- **INCHN-B08** — `SetRuleLetters` reconfigures every pattern of the package that depends on the project's rule-type vocabulary, together.

- **INCHN-B09** — A file that is empty or only whitespace fails the emptiness ruler.

- **INCHN-B10** — A file carrying a scenario code passes the identity ruler, and one carrying none fails it.

- **INCHN-B11** — A governed file with no identity block at all fails the header ruler.

- **INCHN-B12** — A governed file whose header carries ownership OR reference passes; one carrying only a layer does not.

- **INCHN-B13** — A file of a RECOGNIZED layer passes the header ruler with the layer alone, because it has neither an owning spec nor a sibling to reference.

- **INCHN-B14** — A BINARY file steps aside from the header ruler: there is no comment syntax in an image, and charging one would bar every visual baseline commit.

- **INCHN-B33** — A governed file whose `@anchors` block stands below the top without `@fixed-header: <why>` fails the header ruler saying the block is not read as the header, and how to fix it.

- **INCHN-B34** — `scenario-coverage` counts a rule proven only when every scenario the spec's features declare for it is proven — each variant `#NN` a scenario of its own —, and names each variant left unproven: a skipped `#02` is not proven by the green `#01` beside it, nor by a proof recorded for the rule alone. (`checkScenarioCoverage`, `specScenarios`)

- **INCHN-B35** — `header-valid` reads the header's layer as the project declares it: a file whose header names a layer the Structure declares `regime: declarativo` has its identity in `layer:` alone, whatever the layer is called and whether or not the node carries a regime of its own; a layer declared with another regime still asks for `code:` or `ref:`. (`checkHeaderConforms`, `isRecognizedLayerCfg`)

- **INCHN-B36** — `line-coverage` and `coverage-delta` skip a file a coverage report listed with no instrumentable line at its revision — there is nothing to cover —, and answer a divergence for one a whole run of its suite left out of the report — a tool omits a file with no instrumentable line, and also one outside what it collects —; a file never listed nor omitted is pending, never measured. (`coverageAbsence`)

- **INCHN-B37** — `mutation-score` skips a file the mutation tool listed at its current revision with no mutant — an alias, a re-export: it was measured and has nothing to mutate —; a file never listed, or listed at another revision, stays pending.

- **INCHN-B38** — `line-coverage` holds a code file to the floor of the first glob of the gate's `coverage_floors` that matches it, in name order, and otherwise to its `min_coverage`, 70% when undeclared; a file below a glob's floor fails naming the glob and its reason. (`CoverageFloorFor`)

- **INCHN-B39** — A header line is read in every comment dialect the map reads — `//`, `#`, `--`, `<!--` and a block comment's ` * ` (`config.HeaderLinePrefix`): a `-- ref: CODE` header has its identity.

- **INCHN-B40** — `header-valid` reads the header as the map does — the `@anchors` block at the top (`scan.AnchorsHeader`) —, so an `@anchors` further down, in a string or an example, is neither the header nor its identity; and a guide, a document or a test support file, which belong to no unit, have their identity in `layer:` alone.

- **INCHN-B41** — `header-valid` requires, in a header that has its identity, the file's OWN `code:` — a `ref:` alone names the unit, not the file — and fails a code another node of the map carries as its own, naming that file. (`checkOwnCode`)

- **INCHN-B15** — An executable test script steps aside too, by a different path: its format belongs to the runner, and its identity is in the file name.

- **INCHN-B16** — A guide with no compliance-points section, or with the section and no item in it, fails — the AI judgment gate would otherwise fall back on vague heuristics.

- **INCHN-B17** — When the project says how its tests are written, `scenario-coverage` counts a scenario as written in a file the source lists tests in only when a test TITLE cites its code; without a source, or in a file the source lists no test in, a code anywhere in the file outside comments counts.

- **INCHN-B18** — A support file is not judged by `tests-pass` (Skip, saying why), and it does not count as a test that names a scenario for `scenario-coverage`.

- **INCHN-B19** — `non-empty` passes a feature only when it declares a scenario, and a scenario opens with any keyword of the official Gherkin table, in any language and synonyms included; the examples table of an outline is not a scenario.

- **INCHN-B20** — `scenario-coverage` charges only the requirements the spec DEFINES: a code the spec merely cites in its prose is never charged, and a defined requirement with no proven scenario still fails, named.

- **INCHN-B21** — `scenario-coverage` tells a requirement no test names apart from one a test names but no ingested execution proved, even when no execution was ingested at all, and the verdict says which of the two each one is; a requirement an ingested execution proved is not charged.

- **INCHN-B22** — A spec whose layer dispenses `tested-by` is skipped by `scenario-coverage`, saying so, as `unit-complete` does; a layer without that opt-out is still charged.

- **INCHN-B23** — `mutation-score` passes a file whose score reaches the acceptable threshold — the threshold itself included — and fails one below it, naming how many mutants survived and the threshold. A file where no mutant ran does not apply: all ignored says so with the count, and none covered says the line coverage owns it.

- **INCHN-B24** — With no mutation signal ingested, or with one measured at another revision of the file below the floor or under load, `mutation-score` is pending, saying what to ingest or that the signal is stale; a stale score that met the floor, measured without load, does not block — mutation is not remeasured on every change —, and says how to measure it again. The revision is the mutation's own.

- **INCHN-B25** — A score between the acceptable and the desirable threshold is a divergence the project's threshold accepts — it informs and never bars; to make it bar, the project raises the acceptable threshold —, naming both ranges and how far the desirable one is; at or above the desirable it passes clean; with no desirable threshold, or one not above the acceptable, the acceptable threshold alone decides.

- **INCHN-B26** — The mutation thresholds are the ones the ingested report declares; the engine's default of 70% applies only when the report declares none.

- **INCHN-B27** — With the mutation score measured per scope, the verdict follows the ISOLATED score, never the full one, and the report names both and their delta — a large delta read as coupling to dependents, a small one as a missing assertion, never both at once; with no scopes the total score decides.

- **INCHN-B28** — A mutation scope measured at an older revision than the file's decides nothing and is left out of the report; a scope carrying no revision stamp still counts.

- **INCHN-B29** — Outside a git repository `updated-at-atual` skips, naming the missing repository instead of blaming the file as uncommitted; inside one, a new file dated today passes and a wrong date still fails.

- **INCHN-B30** — `spec-sections` fails a section title written in another language of the catalogue than the project's, naming the expected title, unless the gate declares `enforce_section_language: false`; a title outside the catalogue, or a run with no configuration, is not charged for its language.

- **INCHN-B31** — The skeletons `anchors new` emits are born conforming: the spec passes the header and spec-sections rulers, and the feature and the test pass the header ruler.

- **INCHN-B32** — When the share of a file's mutants killed by the time limit is above the gate's `timeout_ceiling` (0.2 by default), its mutation score is Pending as measured under load, whatever the score, and the verdict says to measure again first with fewer workers and with the test cache off — the usual causes are the tool's own parallelism and a time limit taken from a cached run —, and then with no time limit to learn how long a mutant takes; below the ceiling the score decides, and a failure says how many were killed by the time limit.

- **INCHN-I01** — A check name that does not resolve NEVER approves. Approving would stamp green over a verification that never ran, and the project would read coverage where there is none.

- **INCHN-I02** — Every registered name is reachable through exactly one of the routing paths. A name registered in no reachable registry is a gate that is accepted in silence and measures nothing.

- **INCHN-I03** — The compliance ruler is recognised in EVERY language of the catalogue, not only the one the engine was written in. A project seeded in another language would otherwise be born failing a guide its own tooling had just written.

- **INCHN-X01** — Does not decide WHICH checks a project runs.

- **INCHN-X02** — Does not invoke external tooling.

- **INCHN-X03** — Does not judge whether the text it reads is GOOD.

- **INCHN-E01** — A test the map lists is no longer on disk when `scenario-coverage` looks for the tests that name each scenario code.

- **INCHN-E02** — The project's tests source fails, or answers outside its contract, when `scenario-coverage` looks for the tests that name each code


## LYBNL — LayerBoundary — a layer does not reach what is not its own

Confronts a code file against the architecture the project DECLARED: **the layers exist on
paper, but do they respect each other?**

Anchors already declared layers — and did not confront whether they hold. Declaring
`screens/`, `hooks/`, `repositories/` and never verifying that the screen does not talk
straight to the repository is drawing the architecture and not defending it: `layers:`
becomes documentation.

The gap surfaced because every project solved it alone. A real project kept a
**366-line shell script with 15 hand-written architectural rules** — all of the same
shape: "files matching THIS pattern must not contain THAT one". Fifteen instances of a
single mechanism, reimplemented because the framework did not offer it.

The rule is declared in `anchors.yaml` and is language-agnostic — **the project writes the
pattern, the engine knows no dialect**: a `layer` says to whom it applies, `forbid` what
must not appear, `because` the reason that keeps it from being ritual, and `severity`
the maturation of the rule. The honest opt-out is `@allow-boundary: <reason>` on the line:
acknowledged debt stays visible and dated in the code, instead of becoming an exception in
a distant list nobody revisits.


- **LYBNL-B01** — An artifact that is not code leaves without a verdict: there is no import to forbid in a spec or a feature.

- **LYBNL-B02** — Content matching a forbidden pattern FAILS, and the verdict names the LINE and the REASON — a prohibition with no motive turns into ritual.

- **LYBNL-B03** — A rule scoped to a layer charges only that layer: the same import in a hook is legitimate.

- **LYBNL-B04** — A rule with no `layer` holds for ALL code — that is how a global prohibition is declared (raw clock, literal colour, console log).

- **LYBNL-B05** — `severity: warn` records a divergence without failing; the default severity is `error`.

- **LYBNL-B06** — `@allow-boundary: <reason>` on the line waives THAT line — acknowledged debt stays visible and dated where it lives.

- **LYBNL-B07** — The waiver also holds in the comment ON THE LINE ABOVE: in many languages an import has nowhere to carry a readable end-of-line comment, and demanding it inline would push the author not to declare at all.

- **LYBNL-B08** — A BARE marker, with no written reason, does not waive — it would be a silent way to quiet the gate.

- **LYBNL-B09** — With no boundary declared the verdict is Pending, never Pass: Anchors does not know the project's architecture, and pretending it checked is worse than saying what is missing.

- **LYBNL-B10** — An invalid `forbid` pattern FAILS visibly: swallowed in silence it would switch the rule off with nobody knowing.

- **LYBNL-B11** — The pattern is matched against the WHOLE file, so a target spread over several lines is caught — a multi-line import included.

- **LYBNL-I01** — The engine knows no language. The same architectural rule is expressible in the import dialect of TypeScript, Python, Go, Java, Rust or Ruby, because the one who writes the pattern is the project.

- **LYBNL-I02** — Line-to-line matching under-reported. Measured in a reference app: of the 7 screens importing `Modal` from react-native the gate accused ONE — the only one with the import on a single line. The other 6 wrap because they pass 100 columns, and escaped: the screen was accused for FORMATTING, not for being different.

- **LYBNL-I03** — The waiver holds on ANY line of the matched stretch. In an import the formatter wrapped, the natural marking sits on the `from` line; demanding it on the first line of the match would require the author to know where the regex started matching.

- **LYBNL-I04** — A project that WANTS to anchor the pattern to a single line still can: `^` and `$` keep holding per line, because `(?s)` changes the `.`, not the meaning of the anchors.

- **LYBNL-X01** — Does not decide WHICH boundaries exist.

- **LYBNL-X02** — Does not parse the language: it does not read imports, it matches TEXT.

- **LYBNL-X03** — Does not judge whether the forbidden thing is architecturally wrong.

- **LYBNL-E01** — REF[LYBNL-B10]: a forbid pattern that does not compile is answered by B10: it fails visibly


## MRPRM — MarkerParity — the same rule has to appear at BOTH ends that fulfil it

Confronts a rule that lives in TWO places at once: **the promise at one end, and the
fulfilment at the other — are they still the same rule?**

The defect class is the mapping that comes undone on one side only. A rule split across
two ends cannot be checked by any per-file gate: each side, looked at alone, is
impeccable. What breaks is the RELATION, and it disappears without leaving an error.

The case that motivated it, measured (reference app, `EXSC-Q01`): the data-deletion page LISTS
what will be erased in each scope, and the backend ERASES. The two lists were born
together and nothing binds them. If a scope starts erasing more (or less), the page goes
on showing the old version — and the data subject consents on the basis of it. That is
informed consent over the exercise of a right, so the mismatch is not cosmetic.

**This gate is GENERIC on purpose, and that is what separates it from every neighbour.**
It is not canonical: a project may declare SEVERAL instances of it, each with its own
`id`, its own `marker_prefix` and its own `marker_scopes`. What it confronts is a
decision of the project — WHICH rules live at two ends — and that decision does not fit
a universal catalogue. Every other gate in the framework asks a question the framework
itself wrote; this one asks the question the project wrote.

The suffix after the prefix is the NAME of the rule, and it is the name that pairs the
ends. **Two markings on the same side do NOT satisfy the gate:** `marker_scopes` exists
precisely because counting without looking WHERE would let through the very case the
gate exists to catch.


- **MRPRM-B01** — A rule marked at both declared scopes passes: the mapping still has its two ends.

- **MRPRM-B02** — A rule marked at one scope and missing from the other fails, and the verdict NAMES the scope that was left empty.

- **MRPRM-B03** — Two markings on the SAME side do not satisfy the gate — they add up to two and would hide exactly the mismatch it exists to catch.

- **MRPRM-B04** — A rule that survives at one end after being removed from the other fails, and the verdict names the leftover rule.

- **MRPRM-B05** — TOTAL absence of the prefix is not approval: it is almost always a typo in the declaration, and returns Pending asking to check `marker_prefix`.

- **MRPRM-B06** — A declaration with no `marker_prefix` returns Pending: with no prefix there is nothing to confront.

- **MRPRM-B07** — With no scopes declared, the ruler is the COUNT — the weaker mode, which the gate documents as such.

- **MRPRM-B08** — The marking crosses LANGUAGE: the mapping between a `.ts` page and a `.go` handler is the same mapping.

- **MRPRM-B09** — Directories of the ignore list, `node_modules` among them, never count towards parity.

- **MRPRM-E01** — A directory or file under the root cannot be read (permissions) while the markings are walked.

- **MRPRM-I01** — The failing verdict NAMES what is missing — the rule and the scope left empty. A gate that fails without saying what transfers the diagnostic work to whoever reads it, and here the whole point is that neither side looks wrong on its own.

- **MRPRM-I02** — Nothing measured is never approved. Missing prefix, missing count and scopes, and total absence all return Pending, not Pass: approving without having looked would make the gate look vigilant while watching nothing.

- **MRPRM-X01** — Does not read the CONTENT of what each end does.

- **MRPRM-X02** — Does not decide WHICH rules live at two ends.

- **MRPRM-X03** — Does not read files whose extension is not in the text list.


## MKSTP — MockStampGenerator — writes the missing `@contract` stamps, and never rewrites one

A project adopting `mock-stamped` has every double of a governed module unstamped at once.
Measured in the project that asked for this: 1278 `jest.mock`/`vi.mock` calls, none stamped.
Stamping by hand is not viable, and a hand-written stamp is also the one most likely to point at
the wrong snippet.

The generator writes, above each double that has none, the stamp the gate recomputes. It is built
from the gate's own functions, so a stamp it writes is exactly the stamp the gate checks. Measured
on that project, on a copy: 1214 stamps in 278 files, and afterwards `mock-stamped` 278 ✓, 0 ✗.

**It only ADDS.** An existing stamp is never rewritten, even when it diverges. A divergence is the
gate saying the double may have drifted; a tool that refreshed it would let whoever edits the test
regenerate the stamp to match their own mock, and the stamp would certify itself.

Exposed as `anchors stamp [tests...]`, with `--dry-run`.


- **MKSTP-B01** — `GenerateStamps` writes a stamp that passes the gate, and a later change to the stamped snippet fails it.

- **MKSTP-B02** — A specifier matching several files resolves to the one sharing the longest directory prefix with the test; a tie is skipped and reported, never guessed.

- **MKSTP-B03** — Each factory key that names an export of the module gets its own stamp, anchored on that export and covering its block — a change to another member does not make it diverge.

- **MKSTP-B04** — With no factory key naming an export (an automock, for instance), one stamp covers the whole module — which is what an automock replaces.

- **MKSTP-B05** — A double of a module outside the map (a third-party library) is not stamped.

- **MKSTP-B06** — A whole-module stamp starts after the module's `@anchors` header, so an edit that only rewrites the header (`updated_at:`, which `check --fix` bumps on every edit) does not make it diverge.

- **MKSTP-B07** — A key added later to an already-stamped double gets its own stamp when no existing stamp of that module covers its export (same anchor, or a range containing it); the existing stamps are left untouched, and a second run adds nothing.

- **MKSTP-I01** — An existing stamp is never rewritten, even when it diverges.

- **MKSTP-I02** — A line that occurs more than once in the module never anchors a stamp — the gate refuses an ambiguous anchor.

- **MKSTP-X01** — Does not refresh a divergent stamp.

- **MKSTP-E01** — The project's `derived.mock_detect` does not compile as a regular expression, or has no capture group.

- **MKSTP-E02** — The project does not declare `derived.mock_detect`.

- **MKSTP-E03** — The file the map resolves a double to is no longer on disk.


## MCSTM — MockStamped — the double carries the mark of the snippet it replaces, and the gate RECOMPUTES it

Confronts a test against the question no other ruler asks about a test double: **the double
claims to replace a snippet of the real module — is it still the snippet that exists today?**

It is the AGNOSTIC half of the pair that attacks mock drift. `mock-typed` solves the problem
where the language helps — structural types let the compiler check the module's surface — and
only there: most ecosystems have no `Partial<typeof X>`, and not even in TypeScript does it
reach the SHAPE of the returned VALUE. Measured: 202 tests doubled a React Query hook with
two fields where the real one returns about 15, and the type does not tell "deliberate
partial" apart from "out of date".

The stamp does not interpret the code: it READS it. A text hash catches any change —
signature, body, type, constant — and works in Python, Ruby, Go or plain JS, with no
per-language extractor.

**The gate RECOMPUTES instead of validating format**, and that is what makes the mechanism
resistant to whoever writes it. A stamp nobody confronts is theatre: whoever edits the test
would regenerate it to match their own mock, and it would start certifying itself.

**What separates it from its neighbours:** `unit-complete` asks whether the test exists,
`feature-test-match` confronts scenario against case, `tests-green` reads the run — all three
answer "yes" about a double that froze a contract which no longer exists. `mock-typed` demands
the TIE to the real module; this gate demands the recomputable MARK of the snippet.


- **MCSTM-B01** — An artifact that is not a test leaves without a verdict: only a test declares doubles.

- **MCSTM-B02** — A stamp whose recomputed hash matches the module today passes.

- **MCSTM-B03** — A snippet that changed since the stamp was written fails, and the verdict shows the value the contract has today.

- **MCSTM-B04** — The stamp is immune to DISPLACEMENT: the anchor is searched by CONTENT, so editing lines above does not invalidate the stamp of a snippet that did not change.

- **MCSTM-B05** — An anchor that vanished — renamed, removed or rewritten — fails with its own message: it is a finding, not a tool error, because the double is certainly out of date.

- **MCSTM-B06** — An anchor occurring more than once fails as ambiguous: a stamp pointing at "one of the two" proves nothing, and the gate reports rather than choosing.

- **MCSTM-B07** — `StampSnippet` delimits the window by the declared line count: a change beyond it is not reached, and one inside it is. The reach stays in plain sight of whoever reads, and the gate needs no per-language parser to find where the block ends.

- **MCSTM-B08** — Without `derived.mock_detect` declared the gate is Pending, naming the field and the opt-out: it did not check, and saying so is what puts it in the pending items and the doctor; with `mock_detect` in `dialect.opt_out` the project decided, and the gate skips.

- **MCSTM-B09** — A double of a module the project does not govern is not charged.

- **MCSTM-B10** — A stamp whose module no longer exists on disk is a finding with its own message, not a crash.

- **MCSTM-B11** — The ABSENCE of a stamp on a governed double is accused, not skipped — and this is the most important half of the gate.

- **MCSTM-B12** — A stamp that is present and correct satisfies both charges at once: the absence and the correspondence.

- **MCSTM-B13** — A dialect regex that does not compile fails LOUDLY: it is a configuration error, and silencing it would make the gate sweep zero doubles and report green.

- **MCSTM-B14** — A dialect regex with no capture group fails too: without it the gate cannot know WHICH module was doubled.

- **MCSTM-B15** — A different ecosystem's dialect is charged exactly the same way once declared — the stamp is agnostic in fact, not in intention.

- **MCSTM-B16** — A double whose module name contains a dot (`@/src/stores/auth.store`) is matched to the stamp of `auth.store.ts`: the specifier is compared as written and without a final extension, and both forms count. Stripping an "extension" from an import specifier cut part of the NAME.

- **MCSTM-B17** — WHOEVER CHANGES A MODULE SEES THE DOUBLES IT BREAKS. `TestsStamping` resolves the tests whose stamps point at the module, so `check --changed <module>` brings them into the check and this gate runs on them, naming the module file of each stamp. Without it the drift surfaced to whoever next touched the test, far from the change.

- **MCSTM-B18** — `RefreshStamps` (`anchors stamp --refresh <module>`) is how the author of a change updates the stamps: it lists every double stamped against the previous version — test, line, member, old and new hash, and how the stamped block changed from HEAD — and updates those hashes. The list is the work the change created: each double reproduced the old contract, and is adjusted in the same commit. A stamp whose anchor is gone is NOT refreshed, because only a person can say which new line the double now stands for.

- **MCSTM-B19** — The modules a test file's stamps point at are listed each once, in the order they first appear; a file with no stamp lists none. (`StampedModules`)

- **MCSTM-B20** — `StampsHolding` lists the stamps of a test that match their module today, keyed `<module>

- **MCSTM-I01** — The gate RECOMPUTES the hash against the real module; it never validates the stamp's format alone. A stamp nobody confronts would certify itself, because whoever edits the test regenerates it to match their own mock.

- **MCSTM-I02** — The double and the stamp are tied by the module's PATH, matched by suffix without extension. The specifier and the stamped path describe the same file through alias and disk path, and a per-ecosystem alias resolver would be needed otherwise.

- **MCSTM-I03** — A configuration fault fails; a project decision skips; nothing decided is pending. Compiling error and missing capture group are Fail, an undeclared pattern is Pending, a waived one is Skip — the difference is whether somebody CHOSE the silence.

- **MCSTM-X01** — Does not interpret the code of the stamped module.

- **MCSTM-X02** — Does not carry a built-in dialect for detecting doubles.

- **MCSTM-X03** — Does not skip the absence of a stamp to accommodate legacy code.

- **MCSTM-X04** — Does not guarantee cryptographic strength: the hash is truncated.

- **MCSTM-E01** — A stamp declares a line count that is not a positive number (`0`, `-3`, `x`).

- **MCSTM-E02** — A test the map still lists is no longer on disk when `TestsStamping` looks for the doubles of a changed module.


## MCTYM — MockTyped — every test double must DERIVE from the module it replaces

Confronts a test against the hole that is the most treacherous of the whole gate family,
because it is not the absence of proof — it is FALSE proof. **A test that doubles its
neighbour keeps passing after the neighbour changes signature, return or name:** the double
became a frozen copy of a contract that no longer exists. The green certifies the old
version, and nobody goes looking for a defect where there is a green test.

The neighbouring gates do not reach it, and each of them answers "yes" about a test that
lies:

| gate | what it asks | its answer |
| --- | --- | --- |
| `unit-complete` | does the test EXIST? | it does |
| `feature-test-match` | does the scenario match a test case? | it does |
| `tests-green` | did the run pass? | it did |

The defence is not the gate reimplementing type checking: it is demanding that the double be
TIED to the original by a mechanism the language already knows how to check. In TypeScript,
annotating the factory with a partial of the real module's type makes the compiler accuse
both a non-existent method and a divergent return; in Python it is autospec; in Go it is the
interface. **The framework does not know which — the project declares it in
`derived.mock_contract`, and this gate charges the presence of the tie.**

**What separates it from `mock-stamped`:** this one solves the problem where the LANGUAGE
helps, and only there. `mock-stamped` is the agnostic half, which hashes text and works where
there is no structural type to lean on.


- **MCTYM-B01** — A double with NO tie fails, and the verdict names the loose module — this is the false proof the gate exists to catch.

- **MCTYM-B02** — A double whose factory carries the declared tie passes: the annotation is what makes the compiler check name, signature and return against the real module.

- **MCTYM-B03** — The charge is per MODULE, not per file: one loose double among several is enough to fail, and the verdict counts only the loose ones.

- **MCTYM-B04** — A double with no factory is not charged: the automock derives from the real module by construction, so it cannot drift, and charging it would be noise one learns to ignore.

- **MCTYM-B05** — Without `derived.mock_contract` declared the gate is Pending, naming the field and the opt-out; with `mock_contract` in `dialect.opt_out` it skips.

- **MCTYM-B06** — An artifact that is not a test leaves without a verdict.

- **MCTYM-B07** — A test that doubles nobody leaves without a verdict — Skip and not Pass, because nothing was checked and a Pass would inflate the count of greens with nothing.

- **MCTYM-B08** — Another runner of the same ecosystem is recognised the same way: the double is the double regardless of which library spells it.

- **MCTYM-B09** — A third-party library double is not charged.

- **MCTYM-B10** — The third-party exemption is not an escape hatch: one OWN loose double among third-party ones still fails, and only the own one is named.

- **MCTYM-B11** — A relative import that resolves to a node of the map is an own module like any other — the criterion is the graph, not the shape of the specifier.

- **MCTYM-B12** — Another ecosystem's dialect is charged the same way once declared, with its own detection pattern and its own tie shape.

- **MCTYM-B13** — A double whose call is broken over several lines — the module and the factory on the lines after the opening, as a formatter writes it — is read as one written on a single line: an annotated factory passes, and one with no annotation is charged.

- **MCTYM-I01** — What is governed is decided by the GRAPH, never by a prefix list in the config. That asks for no new configuration, assumes no alias convention — which varies per ecosystem — and follows the project on its own: code that is born enters the map and starts being charged.

- **MCTYM-I02** — An undeclared tie shape is PENDING instead of guessing one. Inferring the TypeScript form would assume the ecosystem and report GREEN over what was never checked in any other — the worst possible failure in a measuring device.

- **MCTYM-I03** — The verdict COUNTS the loose doubles, not just names them. Naming one without saying how many leaves the reader unable to tell whether the others are on the list too.

- **MCTYM-X01** — Does not check whether the annotated type actually MATCHES the real module.

- **MCTYM-X02** — Does not cover drift of BEHAVIOUR — only drift of SHAPE: name, signature, return type.

- **MCTYM-X03** — Carries no built-in tie shape and no built-in ecosystem.

- **MCTYM-X04** — Does not charge third-party library doubles.

- **MCTYM-X05** — When the tie is a type on the factory (the form carries `{{module}}`), does not charge a double with no factory — an automock (no second argument) or an options object such as Vitest's `{ spy: true }`. When the tie is an option on the call (`autospec=True`), the bare call is still charged: Python's `patch` without it is a `MagicMock`, not the module.

- **MCTYM-E01** — REF[MCTYM-B05]: with no configuration or no tie shape declared, B05 answers Pending naming what to declare

- **MCTYM-E02** — No map has been built and the test mocks modules.


## OBHNB — ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit

Confronts the CROSS-CUTTING OBLIGATIONS a project declares against reality: **a node whose
header carries the trigger attribute must appear in the files the obligation demands.**

It is the class of defect no per-unit spec catches, because the duty lives OUTSIDE the
unit. The case that motivated it is real: a new data model, carrying free-text annotations
written by the user, was left out of the account-deletion script — personal data that would
never be erased. The project had been bitten by this before, and the script itself carries a
comment calling it a "silent LGPD violation". It repeated, and none of the nine existing
gates saw it, because every one of them looks INSIDE the unit.

**The honest exception, and the third state.** A node can waive itself with
`obligation_waived: <name> — <reason>` in the header, and the reason is MANDATORY: a waiver
without a justification is treated as absent. But waiving is not the only real case — the
most common one is that the duty is REAL and will be paid in another phase, when the handler
that consumes the table does not exist yet. Waiving would be a lie (the duty did not stop
existing), and leaving it red confuses ACKNOWLEDGED DEBT with forgetfulness, which is exactly
the distinction the pillar exists to preserve. `obligation_pending: <name> — <when>` asserts
three things — that the duty is known, that it still holds, and when it will be paid — and
the verdict is Pending, visible in the report and never Pass.

What separates it from its neighbours: `doc-required` charges a document the Structure
declared for a LAYER; this one charges the presence of a TOKEN in the files a named duty
points at, triggered by an attribute the node itself declares.


- **OBHNB-B01** — A node that carries the trigger and does not appear in the demanded file fails, and the verdict carries the declared REASON for the duty.

- **OBHNB-B02** — A node that carries the trigger and does appear passes — the duty is fulfilled.

- **OBHNB-B03** — A node without the trigger attribute contracts no obligation: the duty is charged by what the node declares about itself, not by what it might resemble.

- **OBHNB-B04** — A waiver WITH a written reason exempts the node; the same waiver without one does not, because that is what separates the honest exception from silence.

- **OBHNB-B05** — A project that declares no obligation is skipped: there is nothing to confront, and inventing duties would charge what nobody committed to.

- **OBHNB-B06** — An acknowledged DEBT, with the when written down, yields a divergence — it is a record, visible in the report, never an exemption.

- **OBHNB-B07** — A bare debt marker, with no when, keeps failing: it assumes no debt, it only hides better.

- **OBHNB-B08** — Waiver and debt stay distinct: only the waiver resolves the duty, because the debt is still owed.

- **OBHNB-B09** — The failing verdict OFFERS the three ways out — fulfil, waive with a reason, or acknowledge the debt with a when.

- **OBHNB-I01** — The token searched for is derived through the declared form — `screaming-snake`, `snake`, `kebab`, a free template, or the raw name. Guessing the shape in the engine would put one project's mess inside the framework.

- **OBHNB-I02** — A glob that matches no file produces no violation. Accusing where there was nothing to read would stamp what was never measured.

- **OBHNB-I03** — The node's own `identified_as` wins over the obligation's automatic form. It is the only source that knows the project's real irregularity — one model becomes a plural env var, a sibling becomes a singular one, with no derivable rule; inverting this order accuses 28 correct models (measured).

- **OBHNB-X01** — Does not decide WHICH obligations exist, nor which files satisfy them.

- **OBHNB-X02** — Does not understand what the destination file DOES with the token.

- **OBHNB-X03** — Does not read a declaration written in the body of the document.

- **OBHNB-E01** — A file matched by a `must_appear_in` glob cannot be read (permissions).

- **OBHNB-E02** — A `must_appear_in` glob of a triggered obligation does not parse (`purge[.ts`).


## BLGTN — Obligations — the duties in force, resolved from packs and config, and their status across the project

A cross-cutting obligation is a duty that lives outside the unit: a node whose header carries a
trigger attribute must appear in the files the duty names (the account-deletion script, the audit
log). The obligation-honored gate judges one node against one duty. This unit holds the two pieces
around that gate, which must agree with it.

The first resolves which duties are in force. A project declares duties inline in its configuration
and may adopt packs, sets of duties that come from a norm (a privacy law, for instance). Pack duties
come first and inline ones last, so that a local declaration, the most specific one, is the last word.
A pack duty carries its source into its reason, so the gate's message cites the norm instead of only
asserting the duty. A pack that fails to load is a configuration error and is said on the error
output; it drops only its own duties, never the project's inline ones, because silencing it would turn
a whole set of duties into nothing with a green report. The resolution is read once per project root,
because a full check asks for it for every node.

The second is the compliance report. Its unit is the duty, not the file: "42 nodes are subject, 41
comply" answers an auditor, "file X violates" repeated 42 times does not. For each duty it counts the
subjects, how many fulfil it, how many declared it as acknowledged debt, how many waived it with a
reason, and names the ones that do not comply. It judges each node through the same path as the gate,
because two implementations of one rule diverge, and the wrong one would be the report, where people
trust without checking.


- **BLGTN-B01** — The duties in force, as seen from outside the package (`ObligationsInForce`), are the resolved list the gate uses, pack duties included.

- **BLGTN-B02** — With no configuration there is no duty; with no pack adopted, the duties in force are the inline list as declared.

- **BLGTN-B03** — The pack duties come first and the inline duties last, and each pack duty keeps its trigger, its target files and how it is identified.

- **BLGTN-B04** — A pack duty's reason gets its source appended in parentheses: the authority and the article when both exist, either one alone otherwise; with no reason, the source alone is the reason.

- **BLGTN-B05** — A pack that fails to load is reported on the error output, in the project's language, and the inline duties remain in force.

- **BLGTN-B06** — The pack duties are read once per project root and pack set (the adopted packs, their values and the jurisdictions); later calls with the same root and set are served from that first read.

- **BLGTN-B12** — A call with the same root and another pack set reads that set, instead of being served the list of the first.

- **BLGTN-B07** — The report (`EvaluateObligations`) gives one status per duty, in the order the duties are handed, each with its name and target files.

- **BLGTN-B08** — A node is a subject of a duty only when its header carries the duty's trigger attribute; a duty with no trigger has no subject.

- **BLGTN-B09** — A subject that waives the duty with a reason counts as fulfilled and as waived.

- **BLGTN-B10** — A subject whose judgement is Pending, because it declared the duty as acknowledged debt, counts as debt.

- **BLGTN-B11** — Every other subject that does not pass is named as missing, and the missing list is sorted.

- **BLGTN-I01** — Every subject is counted exactly once: fulfilled, debt, or missing add up to the subjects.

- **BLGTN-X01** — The report evaluates only the duties it is handed: it does not resolve the packs of the configuration it receives.

- **BLGTN-E01** — REF[BLGTN-B05]: a pack that fails to load is the configuration failure B05 reports on the error output while keeping the inline duties

- **BLGTN-E02** — A node of the map whose file cannot be read.


## OPQSP — OpenQuestions — a spec with an open question is not ready to implement

Confronts a spec against the decisions it has NOT yet taken, and keeps it out of "ready"
while there is an open question.

The defect class is UNRESOLVED AMBIGUITY — the cheapest to avoid and the most expensive
to discover late. The path is always the same: the spec does not decide something the code
needs; whoever implements picks a defensible reading and moves on; the choice is never
confronted with the one who had the answer; the product ships with the wrong reading. No other
gate catches it, because all the pieces exist and reference one another — the defect is a decision
nobody took.

What this gate adds to the advice "don't guess, record and report" is a declared
PLACE for the record. Without a place, recording becomes a PR comment that dies in the merge.
With a place, the question is a visible work item, and the spec only becomes implementable when
the section empties.

The intended cycle: whoever writes notices what they do not know and writes it in the section; the gate accuses
while there is an item; the question is taken to whoever decides; the answer BECOMES A RULE, with
a code, and the item leaves the section.


- **OPQSP-B01** — An artifact that is not a spec leaves without a verdict: only the spec has an open decision to demand.

- **OPQSP-B02** — Whoever OPENED the section is confronted by its content: an open item blocks, a closed section releases.

- **OPQSP-B03** — An open item BLOCKS: while there is a question, the spec does not pass as ready.

- **OPQSP-B04** — A section closed honestly — opened under accepted titles such as `## Open Decisions`, `## Open Questions`, or `## Decisões em Aberto` and with no item — releases. Saying "there is no question" is different from not having looked.

- **OPQSP-B05** — An item marked as RESOLVED does not block: the question stays in the trace, and what closed it is the rule that was born from it.

- **OPQSP-B06** — Every question needs a CODE. Without identity it does not become a traceable item nor survive a rewrite of the spec.

- **OPQSP-B07** — `OpenDecisions` COUNTS a spec's pending decisions, for whoever needs the number instead of the verdict — it is what allows reporting the pendency as a systemic lead, in the same standing as an absent signal. The count reads the project's lexicon by the same route as the confrontation: counting zero in a spec whose section is called something else would assert "there is no pending decision" about a spec full of them, which is the silence this unit exists to eliminate.

- **OPQSP-B08** — An item is a question already turned into a rule only when the rule code stands where the answer is written — on a table row, in a cell after the question's text; on a list item, after an arrow or a word of resolution (`→`, `virou`, `became`, `resolved`…); a code the question cites as context (a state, another rule) leaves it open.

- **OPQSP-I01** — Prose is not an item. Explanatory text inside the section does not count as a question — otherwise the author would learn to explain nothing.

- **OPQSP-I02** — The section's boundary is respected: what comes after it is not read as a question. Without that, the whole spec would turn into a decisions section.

- **OPQSP-I03** — The column that says what the question BECOMES is not its identity, and filling it does not close the question. They are two things: the foreseen destination and the answer given.

- **OPQSP-X01** — Does not judge whether the question is GOOD nor whether the answer is right.

- **OPQSP-X02** — Does not FAIL the spec that does not have the section — it records the divergence and says how to close it.

- **OPQSP-E01** — REF[OPQSP-X02]: the one handled path is the absent section, which X02 answers: recorded as a divergence, never failed


## PGNHN — PaginationHonored — what promises a SET does not return the first page in silence

Confronts a function against the promise its name makes: **whoever says "list all" cannot
deliver the first hundred without warning.**

It is the COST/SCALE class of defect — the one that shows up in no test, because the test
runs with three records and production runs with three thousand. Nothing in the code is wrong:
the query is valid, the type is the expected one, the suite passes. The defect is the difference between what
the name promises and what the function delivers when the data grows.

The distinction that gives the ruler: **a CALLER's limit** is not a defect — whoever passed the limit
knows there is more. **A HIDDEN limit** is: a default value the caller does not see makes the
function promise the set and return a slice. The hundred-and-first row is never
processed, and no one is notified.

When the module has sisters that paginate, the proof is by ASYMMETRY: the author knew the
pattern, and the one that does not paginate is forgetfulness, not decision. Without sisters paginating the verdict is
weaker, and it only accuses when the name promises a set unambiguously.


- **PGNHN-B01** — A function with a limit received from the CALLER passes: the page is deliberate, and whoever asked for it knows there is more.

- **PGNHN-B02** — A function with a limit HIDDEN in a default value is accused: the name promises the set and the return is a slice.

- **PGNHN-B03** — The NAME bounds the promise: whoever does not promise a set is not charged, because there is no promise to break.

- **PGNHN-B04** — When there are sisters that paginate in the same module, the ASYMMETRY is the proof: the author knew the pattern.

- **PGNHN-B05** — The waiver is DECLARED and with a written reason, and leaves the report.

- **PGNHN-B06** — The verdict OFFERS the way out: rename to what the function does, expose the limit, return the cursor, or waive with a reason.

- **PGNHN-B07** — Without a declared dialect the verdict is INDETERMINATE, and says so — approving without being able to read the code would be stamping what was not checked.

- **PGNHN-I01** — The ruler is AGNOSTIC: the confronted truth — the name promises a set, the return is partial — belongs to no language. What comes from the project is only how function, loop and cursor are recognized.

- **PGNHN-I02** — Where the unit does not recognize the construct, it stays silent. Silence is better than an invented accusation — a false positive here trains the team to ignore the gate.

- **PGNHN-I03** — A cursor without a loop does not count as pagination: returning the cursor and not walking it leaves the consumer with the same slice, only with the appearance of completeness.

- **PGNHN-I04** — A provider prefix in the name does not hide the promise: what the name says counts, wherever it comes from.

- **PGNHN-X01** — Does not invent a cursor where the provider offers none.

- **PGNHN-X02** — Does not measure PERFORMANCE nor page size.

- **PGNHN-E01** — REF[PGNHN-B07]: an undeclared dialect is answered by B07: indeterminate, and it says so


## PHORP — PhaseOrdered — plan phases and phase dependencies must be ordered and consistent

Confronts internal plan phase ordering and cross-artifact phase dependencies: **phases within a plan must
follow consistent non-cyclic order, and specifications citing phase prerequisites must target existing phases.**

The `needs:` declaration previously resolved execution order between separate plans, but ordering within a single
plan lived purely as informal prose (such as `### Fase 2 — a régua mecânica (depende da Fase 1)`). Natural language
prose cannot be mechanically verified. When task identification pipelines generate work cards for each specification
in a plan, all cards are born into the backlog indistinguishably, and automated claim pipelines deliver any card
without respecting prerequisite readiness.

Measured in the very first real-world usage: an automated agent was assigned the specification for a test harness
(Phase 3) while Phase 1 and Phase 2 remained completely open. Without even a `package.json` created in the repository,
there was nowhere to configure tooling. The assignment only avoided becoming lost work because a human intervened to
read the narrative plan; an agent trusting the task card would have attempted execution and failed.

Phases are therefore promoted to catalogued items with identity codes derived from the plan (such as `<PLAN>-W01`).
A seeded specification declares `needs: <PLAN>-W01` in its header, adopting the exact ordering keyword at phase scope.
Identity codes remain immutable even when authors revise descriptive phase titles, transforming "can this specification
be worked on now?" into an objective mechanical query.

This gate confronts three complementary structural ordering contracts:
1. **Phase order within plans (`phase-ordered`)**: Phases cannot depend on future phases, duplicate codes, or
   themselves. Plans with phase-like sections lacking catalogued codes return **Pending**.
2. **Phase references in specifications (`phase-exists`)**: A specification declaring `needs:` must target catalogued
   phases that exist in project plans, avoiding permanent task blockages.
3. **Parent hierarchy integrity (`parent-valid`)**: Artifacts declaring a `parent:` must target real entities
   (artifacts or phases) without dangling references or circular parent chains.


- **PHORP-B01** — The `PlanPhases` function extracts catalogued phase codes from plan headers in appearance order.

- **PHORP-B02** — When the confronted node is not of kind plan, phase ordering skips confrontation.

- **PHORP-B03** — When a plan contains no phase headings, phase ordering skips confrontation.

- **PHORP-B04** — When a plan contains phase-like sections without catalogued phase codes, phase ordering returns a divergence.

- **PHORP-B05** — When plan phases declare valid backward dependencies on preceding phases, phase ordering passes.

- **PHORP-B06** — When a plan defines duplicate phase codes, phase ordering fails.

- **PHORP-B07** — When a phase declares a dependency on a phase code not catalogued in the plan, phase ordering fails.

- **PHORP-B08** — When a phase declares a dependency on itself, phase ordering fails.

- **PHORP-B09** — When a phase declares a dependency on a future phase defined later in the plan, phase ordering fails.

- **PHORP-B10** — When a specification declares phase dependencies that exist in project plans, phase existence passes.

- **PHORP-B11** — When a specification declares phase dependencies that do not exist in any plan, phase existence fails naming the missing phases.

- **PHORP-B12** — When an artifact declares an existing artifact code or catalogued phase as parent, parent validation passes.

- **PHORP-B13** — When an artifact declares a non-existent parent, self-parenting, or a circular parent chain, parent validation fails.

- **PHORP-B14** — A phase's dependency is read in any supported language — `depends on`, `depende de` and its contractions, as the catalog lists them — whatever the project's `lang`, so an English plan's order is confronted as a Portuguese one's is.

- **PHORP-I01** — Phase dependencies within a plan must be strictly acyclic and backward-directed; a phase cannot depend on itself or subsequent phases.

- **PHORP-I02** — Every declared phase or parent reference must resolve to an existing entity in the map, preventing orphaned items that disappear from dependency hierarchies.

- **PHORP-I03** — Parent chains are strictly cycle-free and bounded to prevent infinite traversal loops during tree assembly.

- **PHORP-I04** — Phase detection identifies structural level-three section boundaries regardless of linguistic naming variations.

- **PHORP-X01** — Does not mandate that small plans define catalogued phases.

- **PHORP-X02** — Does not enforce timing deadlines or calendar durations for phases.

- **PHORP-X03** — Does not restrict parent references to a single hierarchy kind.

- **PHORP-E01** — No map has been built when a spec's `needs:` or an artifact's `parent:` is confronted.

- **PHORP-E02** — A plan the map lists is no longer on disk when the phases it catalogues are collected.


## PLCFL — PlaceholderFilled — the skeleton the generator emits must be FILLED IN

Confronts an artifact against the question no other gate asks: **was this actually
written, or is it still the frame the generator handed over?**

The work prompt promises, textually, that an unfilled placeholder fails. It did not.
Measured: a freshly generated spec — with the layer field, the date, the title and the
first rule all still carrying the generator's marker — crossed EVERY blocking gate with
"can promote", including the header gate, which read a placeholder as a layer and
approved it.

The reason is structural, and it is why this gate cannot be folded into its neighbours:
**header gates validate FORM** — does the field exist, is the shape right — and never ask
whether the value MEANS anything. A marker is a well-formed value. And the completeness
gate counts sections, which the skeleton has all of, empty.

The cost of that silence is the worst kind: an artifact nobody wrote passes as written.
The relational gates that follow find the spec, find the code, confront two things that
reference each other, and the whole pipeline certifies work that does not exist.

Measured before switching it on, against the real repository: **zero findings across 590
specs**. Nothing alive carries a generator marker — which confirms both that whoever
writes, fills in, and that the gate charges only what was left behind.


- **PLCFL-B01** — A raw skeleton is FAILED: the artifact was generated and never written.

- **PLCFL-B02** — A header FIELD whose value is the marker fails — it is the gravest shape, because a placeholder is not a layer and the header gate approves it as well-formed.

- **PLCFL-B03** — A table CELL holding only the marker fails: the rule exists as a code and says nothing.

- **PLCFL-B04** — A TITLE or body line opening with the marker fails.

- **PLCFL-B05** — The verdict NAMES what was left behind, so the reader does not hunt the file for it.

- **PLCFL-B06** — An artifact with every marker replaced passes — the gate charges the generator's leftovers and nothing else.

- **PLCFL-B07** — The marker vocabulary is the project's, declared in `placeholder_markers`; with none declared it is `TODO`, the only word the `anchors new` templates write. A word outside the vocabulary is an ordinary value, and `<…>` is a marker in any vocabulary, because it is a shape and not a word.

- **PLCFL-I01** — A section the author wrote ON PURPOSE to list pending work is legitimate and is never accused. Measured: 77 specs of a real project carry one. Confusing "the author listed what is missing" with "the author wrote nothing" would punish the honesty the framework asks for everywhere else.

- **PLCFL-X01** — Does not judge whether what replaced the marker is GOOD.

- **PLCFL-X02** — Charges only the marker in a VALUE POSITION — header field, table cell, rule title — and never a marker in running prose.

- **PLCFL-Q01** — Should the marker vocabulary come from the project's Structure instead of being fixed in this unit?


## PCJPL — PlanChangeJustified — a modified plan or spec must declare why it changed

Confronts a modified plan or spec against the declaration of its change: **planning makes mistakes,
and whoever implements discovers them — but fixing the plan silently makes the project walk towards
a destination nobody chose.**

Drift is the greatest risk, and it is SILENT by construction: no STATE gate can catch it, because
the corrected plan or spec is perfectly valid — the inconsistency was removed. What exposes the defect
is not the state of the file, but the UNJUSTIFIED CHANGE.

For this reason, this gate inspects the diff rather than just the content: any plan or spec that
appears among the changed files must carry a declared revision (`-R0001`). Without it, the gate blocks.
The ruler relies on the explicit judgment of whoever made the change:
- An INNOCUOUS correction (wording, example, typo, an ambiguity with only one viable interpretation) —
  correct the text and register the sequential revision. The gate verifies that the revision exists.
- A correction that CHANGES DIRECTION, or any doubt about whether it does — do not fix it in-place.
  Escalate via `anchors escalate` (`anchors:needs-user`), halting delivery of the card until a decision is made.

The gate separates silent drift from transparent adaptation: it does not judge whether a correction is
innocuous or directional, but ensures that a human or agent made the conscious choice and recorded it in
the document itself instead of letting drift happen by omission.


- **PCJPL-B01** — An artifact not listed among the changed files is skipped, even if reached by the impact radius.

- **PCJPL-B02** — An artifact that has no identity code is skipped without failing: identity enforcement belongs to another gate.

- **PCJPL-B03** — A changed plan or spec with no declared revision fails, and the verdict explains how to record the revision or escalate.

- **PCJPL-B04** — A changed plan or spec with a valid sequential revision for its own identity code passes.

- **PCJPL-B05** — Citing the revision of another document does not count as justifying this document's change.

- **PCJPL-B06** — Revisions whose numbering has gaps or is not strictly sequential from 1 fail.

- **PCJPL-B07** — Multiple sequential revisions pass, and the verdict cites the most recent revision explaining the change.

- **PCJPL-B08** — Revision markers formatted in bold, blockquotes, raw lines, or GitHub alerts pass.

- **PCJPL-B09** — A change already explained by the plan revision mechanism (`revises:`, `@revised-by`, `@amended-by`) passes without requiring duplicate notation.

- **PCJPL-B10** — A newly created untracked file has no previous state to justify and is skipped.

- **PCJPL-B11** — A newly created staged file has no commit history and is skipped.

- **PCJPL-B12** — An existing committed file that is modified without a revision fails.

- **PCJPL-B13** — When run outside a git repository, the gate trusts the caller's changed list rather than silencing itself.

- **PCJPL-B14** — The exported function `RevisionsOf` extracts all declared revisions in order with code, number, and explanation.

- **PCJPL-B15** — A revision recorded as a section TITLE (`### CODE-R0001 — what changed`) is a revision too, and it is read together with the revisions of the other formats in the same file.

- **PCJPL-B16** — A changed file that differs from its committed version only in its `@anchors` header and in codes `anchors.renames.yaml` records as renamed — both versions read without the header and with each renamed code as its current one are the same text — is a mechanical change and skips; any other change still asks for its revision.

- **PCJPL-I01** — Only files genuinely modified are charged. A file present in the impact radius but not in the changed list is never accused.

- **PCJPL-I02** — Numbering must be strictly sequential starting at 1. Holes or repetitions prevent determining how many times the document changed.

- **PCJPL-I03** — Newly created files are never charged. A file that did not exist in the previous commit has no past state to justify.

- **PCJPL-X01** — Does not distinguish between innocuous corrections and directional changes.

- **PCJPL-X02** — Does not enforce identity codes on files without one.

- **PCJPL-X03** — Does not run or enforce revisions during full-project checks (`--all`).

- **PCJPL-E01** — REF[PCJPL-B13]: a git that cannot report the changed files is the outside-a-repository case B13 answers: the caller's list is trusted


## PLRVP — PlanRevised — mutual revision visibility between superseded and revising plans

Confronts plan nodes against declared revision relationships: **when a plan is revised, the superseded document must explicitly warn readers at the top, and the revising document is reminded until that warning exists.**

A plan is an anchor for architectural intent and execution. When planning errs — as it inevitably does once real implementation begins — editing an existing, already implemented plan destroys the historical record: it rewrites what was actually decided and executed into something that never happened.

For this reason, architectural revisions must be published as a new plan declaring which prior plan it revises. However, this discipline introduces a subtle and dangerous operational defect: out-of-order reading. A developer or agent opening an older, superseded plan has no natural indication that a newer revision exists, and will proceed to implement or follow decisions that have already been overturned. Because the old plan remains internally coherent — it was, after all, the accurate record of its own time — the reader cannot detect the obsolescence from the document alone.

This gate enforces bidirectional visibility across revised plans. When a plan is revised, it must display an explicit warning notice at its very top (within the first 40 lines) naming the revising plan, and must mark each section affected by the revision. Reciprocally, the author of a revising plan is reminded with a non-blocking pending verdict until that top notice is posted on the target file, ensuring that superseded documents are marked at the moment the revision is committed.

This gate operates in distinct territory from neighbouring gates:
- Unlike `plan-change-justified`, which evaluates whether edits inside an existing plan justify their architectural delta, this gate governs the relationship between two distinct plans across the project map.
- Unlike `phase-ordered`, which ensures linear dependency ordering of phases within a single plan, this gate addresses the lifecycle obsolescence of entire plans and phases across revision boundaries.
- Unlike `plan-seeds-valid`, which ensures seed files cited by a plan actually exist in the repository, this gate ensures that revision targets resolve to valid plan nodes and that mutual revision markers exist on disk.


- **PLRVP-B01** — When the confronted node is not a plan, the gate skips confrontation.

- **PLRVP-B02** — When the map graph is nil, confrontation yields a pending verdict because revision relationships cannot be resolved.

- **PLRVP-B03** — When a plan neither revises another plan nor is revised by any plan, the gate skips confrontation.

- **PLRVP-B04** — When a revising plan declares a revision target that does not exist in the map graph, the gate fails.

- **PLRVP-B05** — When a revising plan targets an existing plan whose file lacks a top revision notice, the revising plan receives a divergence reminder.

- **PLRVP-B06** — When a revising plan targets a plan whose file already contains the top revision notice, the divergence reminder clears.

- **PLRVP-B07** — When a revised plan lacks a top revision notice, the gate fails, reporting the revising plan identifier.

- **PLRVP-B08** — When a revised plan places the revision notice past the first 40 lines, the gate fails because the warning arrives too late.

- **PLRVP-B09** — When a revised plan has a top revision notice within the first 40 lines but no section amendment markers, the gate returns a divergence verdict.

- **PLRVP-B10** — When a revised plan contains both a top revision notice in the first 40 lines and section amendment markers, the gate passes.

- **PLRVP-B11** — Both markdown alert blocks and metadata revision directives are accepted as valid notices and section amendments.

- **PLRVP-I01** — The top revision notice must reside within the first 40 lines of the revised document so that readers following top-down reading order encounter the warning before acting on obsolete decisions.

- **PLRVP-I02** — A missing section amendment marker on a revised plan results in a divergence verdict rather than a hard Fail, accommodating whole-plan revisions where individual sections cannot be cleanly partitioned.

- **PLRVP-I03** — The divergence reminder on a revising plan clears as soon as the revised target file contains the required top notice, preventing permanent noise from training teams to disregard gate output.

- **PLRVP-X01** — Does not mandate a specific human language for revision notices, accepting standard markdown alert callouts and language-agnostic directives.

- **PLRVP-X02** — Does not assess the semantic accuracy or completeness of the explanatory prose written beside a revision marker.

- **PLRVP-X03** — Does not enforce section amendment markers on plans that are not targeted by any revision.

- **PLRVP-E01** — REF[PLRVP-B02]: with no map the revision links cannot be resolved, and B02 answers Pending

- **PLRVP-E02** — REF[PLRVP-B05]: a revised plan whose file cannot be read carries no notice, which B05 answers with the divergence reminder


## PSVPL — PlanSeedsValid — specifications seeded in a plan must target valid governed layers

Confronts artifact paths seeded by an implementation plan against the project Structure: **every specification
path that a plan promises to create must belong to a valid governed layer.**

An implementation plan seeds new artifacts to be born during execution. When a plan declares which specifications
will be created, developers and automated agents rely on those paths to execute their tasks. If a plan seeds
specifications in layers where specifications are prohibited or undefined, that structural contradiction is only
discovered downstream during execution.

Observed across three consecutive rounds in project history: an implementation plan repeatedly seeded
`packages/backend/models/metadata.spec.md — **nasce**`, even though `models/` was a recognized declarative layer
(`regime: declarativo`) which does not accept specifications by architectural definition. Each time, a full
execution round was spent rediscovering the contradiction, because downstream defenses (such as creation tools
and gate runners) only act at execution time. The root cause of the defect was the plan itself. Because the plan
is a declared layer in the repository, its structural promises can and must be verified before execution starts.

This gate confronts two structural defects in plan seeds:
1. **Seeding in a declarative layer**: Seeding a specification in a layer declared as `regime: declarativo`
   (which by definition does not carry specifications).
2. **Seeding in an undeclared layer**: Seeding a specification in a concrete repository directory that does not
   match any declared layer in the project Structure (indicating a path typo or an undeclared layer).

What separates this gate from neighbouring gates:
- It deliberately does NOT verify whether the seeded specification already exists on disk (seeds are intended
  to be created during execution).
- It does NOT check whether plan progress is synchronized or complete (which belongs to plan progress gates).
- It does NOT inspect the content or implementation steps of the plan items.
- It restricts evaluation strictly to the structural validity of the specification paths the plan promises to create.

Finally, casual mentions of template specifications (`_TEMPLATE_*.spec.md`), bare file names in prose, and
informal path fragments that do not correspond to top-level repository directories are recognized as narrative
references rather than actionable seeds.


- **PSVPL-B01** — When the confronted node is not of kind plan, the gate skips confrontation.

- **PSVPL-B02** — When project configuration is nil, confrontation returns Pending.

- **PSVPL-B03** — When the plan contains no seeded specification paths, the gate skips confrontation.

- **PSVPL-B04** — Template specification references prefixed with `_TEMPLATE` are ignored as templates.

- **PSVPL-B05** — Bare specification file names without directory paths are ignored as prose references.

- **PSVPL-B06** — Informal path abbreviations whose top-level directory does not exist on disk are ignored as prose references.

- **PSVPL-B07** — When all seeded specification paths target valid governed layers, the gate passes.

- **PSVPL-B08** — Multiple target source file extensions are evaluated when resolving the governed layer.

- **PSVPL-B09** — When a seeded specification targets a declarative layer, the gate fails citing the declarative layer.

- **PSVPL-B10** — When a seeded specification path in a real repository directory matches no declared layer, the gate fails.

- **PSVPL-B11** — Multiple seed defects across declarative and undeclared layers are sorted and aggregated in the failure verdict.

- **PSVPL-I01** — Only plan artifacts are evaluated; all other artifact kinds skip confrontation.

- **PSVPL-I02** — Declarative layers never accept specifications because declarative schema and model files are recognized rather than authored with specifications.

- **PSVPL-I03** — Seed verification is purely structural and never inspects specification implementation or progress status.

- **PSVPL-I04** — Casual prose citations and template references are never confused with actionable artifact creation promises.

- **PSVPL-X01** — Does not verify whether seeded specifications currently exist on disk.

- **PSVPL-X02** — Does not check whether the plan progress file is synchronized with execution state.

- **PSVPL-X03** — Does not enforce specification contents or scenario definitions within seeded files.

- **PSVPL-E01** — REF[PSVPL-B02]: with no configuration there is nothing to confront, and B02 answers Pending

- **PSVPL-E02** — REF[PSVPL-B06]: a path whose top directory is not on disk is answered by B06: read as prose, not charged


## PSDPL — PlanSourceDeclared — a plan that NAMES a source has to declare who builds it

Confronts a plan against the dependency it wrote in PROSE and never declared in `needs:`:
**the plan names a source — who is going to build the adapter for it?**

The real case, measured in the reference app. Plan 0008 (Frontend/Web) said `Fonte: **GA4**, and
it is the architectural exception of the project`, and declared only
`needs: plans/0005-home-e-indice.md`. The GA4 adapter came from plan 0002, and **that
dependency existed only in the prose.**

What happened: 0002 was revised (`PLTFR-R0002`) and `Ga4Adapter` was REMOVED, with a
correct argument — no document of the project sustained it. 0008 stayed intact,
depending on a source nobody was going to build any more. Both plans remained internally
coherent, and the contradiction only surfaced months later, when 0008 was started.

**Why no gate caught it, and what this one separates from its neighbours.**
`dependency-honored` confronts the `needs:` that is DECLARED; here the defect is the
`needs:` that is MISSING. `plan-seeds-valid` looks at what the plan promises to CREATE,
not at what it promises to CONSUME. This is the version, between PLANS, of step 5 of the
review guide ("the prose ages"): nobody touched 0008, and it became wrong anyway —
because another document changed.

**What it measures.** A source line (`Fonte:`/`Fontes:`) names one or more sources in
bold. For each one, the gate looks for the corresponding adapter among the seeds of ALL
plans. If the adapter exists in another plan and this one does not declare it in
`needs:`, it fails.


- **PSDPL-B01** — A source whose adapter another plan seeds, and which this plan does not declare in `needs:`, FAILS — it is the defect that existed only in the prose.

- **PSDPL-B02** — The failing verdict names WHICH source and WHERE its adapter lives, so the fix is the declaration, not an investigation.

- **PSDPL-B03** — With the owning plan declared in `needs:`, the gate passes: the prose and the declaration now say the same thing.

- **PSDPL-B04** — Every source of the line is confronted on its own — one line may name several, and one declared does not cover the rest.

- **PSDPL-B05** — The source name matches the adapter's file regardless of case and punctuation: `GA4` matches `Ga4Adapter.spec.md`.

- **PSDPL-B06** — A source whose adapter NOBODY seeds is not charged: it may be the source of a future plan, and the gate cannot invent a dependency that does not exist.

- **PSDPL-B07** — The plan that seeds the adapter itself is not charged — it does not depend on itself.

- **PSDPL-B08** — A plan with no source line returns Skip: there is nothing to confront, and that is not approval.

- **PSDPL-B09** — An artifact that is not a plan returns Skip: the gate has no jurisdiction over specs, code or features.

- **PSDPL-B10** — Every source line of the plan is read, not only the first: a plan may name its sources on more than one line, and a source on the second line is confronted like one on the first.

- **PSDPL-B11** — Only case and punctuation are ignored in the match: every letter and every digit of the name counts, so `B0`, `B9`, `Ab` and `Zb` do not match `BAdapter.spec.md` — dropping a character would charge a plan for another source's adapter.

- **PSDPL-I01** — What was not measured is never approved. With no graph, and with no adapter seeded anywhere, the verdict is Pending — approving there would stamp a confrontation that never happened.

- **PSDPL-I02** — The ownership convention fails towards the SAFE side. A seeded file whose name is off the `…Adapter.spec.md` pattern owns nothing, and the gate charges nothing for it — a wrong accusation costs more than a missed one, because it teaches the reader to ignore the gate.

- **PSDPL-X01** — Does not confront the ORDER of the phases.

- **PSDPL-X02** — Does not charge a source whose adapter nobody seeds.

- **PSDPL-X03** — Does not read the plan's prose to understand WHAT the source is for.

- **PSDPL-E01** — REF[PSDPL-I01]: with no map nothing was measured, and I01 answers Pending, never approval


## PRSNT — PresentationGates — the presentation validations, confronted

Each row of `Presentation validations` says: this prop or state, under this condition, makes the
unit look like this. Four questions follow from that shape, all answered from the spec's own text —
no language, no rendering: is every value of the prop decided, does one condition lead to one
appearance, is the text shown a message code rather than copy, and can a test point at what
changes. All four are informational by default.


- **PRSNT-B01** — The rows are read from the Presentation validations section by position — rule, what it reads, condition, appearance —, and a spec with none is skipped by all four gates.

- **PRSNT-B02** — `presentation-exhaustive` fails naming, per prop or state, the declared values no row's condition names; a row whose condition says "otherwise" (in any supported language) covers every value; a prop with no declared set of two values or more is not confronted.

- **PRSNT-B03** — `presentation-conflict` fails naming the rules that give one prop and one condition two different appearances.

- **PRSNT-B04** — `presentation-copy-single-source` fails naming the rule whose appearance carries text in double quotes (or “ ” « ») without citing a message code; single quotes are values, not copy.

- **PRSNT-B05** — `presentation-observable` fails naming the rules whose appearance cites, in backticks, no identifier of the spec's Test Identifiers section.


## PRFLO — Profile — the verdicts of a run, gathered per gate and per node

A check run produces one result per gate and per confronted artifact. Quality in Anchors is not a
single number: it is that set of verdicts. The profile gathers the raw results into the three
answers the rest of the tool needs.

Per gate, it counts how many confrontations passed, failed, skipped, stayed pending or await an AI
judgement, and how long they took: the total time and the single most expensive run. The two times
together tell a gate that is expensive by volume (cheap per target, run hundreds of times) from a gate
that is expensive per target, because the fix for each is the opposite.

For the run as a whole, it decides promotion. A failure of a blocking gate blocks it; a failure of a
non-blocking gate is still a failure (it becomes an issue) but does not block. A pending verdict does
not block by itself: in a real repository, letting every pending verdict of a blocking gate block
rejected 411 nodes at once, most of them meaning "there was nothing to confront". Only the gate knows
whether its pending means "a decision is still open", and it says so by marking the result as one that
impedes; a pending that impedes, on a blocking gate, blocks promotion.

Per node, it collapses the results into one verdict the map uses to stamp its edges: which nodes were
actually confronted, and which of them failed a blocking gate.


- **PRFLO-B01** — The aggregation (`Aggregate`) gives each gate a summary that counts its passes, failures, skips, pending verdicts and verdicts awaiting judgement, and the gate names (`GateNames`) come back sorted.

- **PRFLO-B02** — The summary of a gate carries the total time of its confrontations and the single most expensive one.

- **PRFLO-B03** — Every failed result is listed as a failure, whether its gate blocks or not.

- **PRFLO-B04** — A failure of a blocking gate blocks promotion; a failure of a non-blocking gate does not.

- **PRFLO-B05** — A pending verdict blocks promotion only when its gate is blocking and the gate marked the pending as one that impedes.

- **PRFLO-B06** — Every result awaiting an AI judgement is listed as awaiting judgement.

- **PRFLO-B07** — The per-node verdicts (`NodeVerdicts`) include only nodes that were confronted: a node touched only by skips or by verdicts awaiting judgement is left out, and the list is sorted by node.

- **PRFLO-B08** — A node is marked failed when a blocking gate failed on it, or left on it a pending marked as one that impedes (the same results that block promotion, `PRFLO-B05`); a failure of a non-blocking gate, or a pending that does not impede, leaves it confronted but not failed.

- **PRFLO-I01** — Promotion is refused exactly when the list of blocking results is not empty.

- **PRFLO-X01** — Does not decide what a verdict means: whether a pending impedes, and whether a gate blocks, arrive already set on the result.


## PRHNP — ProgressHonest — the progress file tells the truth about the disk

The `-progress.md` is the only Anchors artifact that nothing confronted, and its exclusion
from the map is deliberate: it exists in order to CHANGE, and the gates that demand a
justification for change cannot reach it.

**But "outside the map" turned into "outside any verification", and the two are not the same
thing.** An item that cites a file PATH is trivially confrontable: the file exists, or it
does not.

Measured in the reference project: when the 17 progress files were created, the checkbox
state was carried over from the plans — and the plans were out of date. The progress of
`0002` claimed 6 open items with 7 of the 8 specs already on disk. FIVE items lied, and the
lie was transported faithfully.

**The damage runs in two directions, and the second is worse:** an open box with the file
already there produces rework — somebody redoes what is done. A ticked box with the file
absent declares a plan finished with work still to do, which is the same defect the explicit
`Closes` shut through another door.

**A gate that only looks INSIDE the file never sees what is missing from it.** That is the
third direction, and it is how this gate once passed a plan that had just gained a seeded
spec: the progress said the phase was over, and the next-card command moved on to another
plan with work declared undone.

The gate is anchored on the PLAN — which IS in the map — and confronts its companion. That
is the only way to reach a file that, by design, is not a node.


- **PRHNP-B01** — An artifact that is not a plan leaves the confrontation without a verdict: the gate is anchored on the plan, which is what the map holds.

- **PRHNP-B02** — A plan with NO companion progress file is skipped, not failed — "does the progress exist" is a different question from "is the progress true".

- **PRHNP-B03** — The skip for a missing companion says HOW to create it, so the reader does not have to look the command up.

- **PRHNP-B04** — A TICKED item whose file does not exist FAILS: it declares done what is not.

- **PRHNP-B05** — An OPEN item whose file already exists fails too: it produces rework, because somebody redoes what is done.

- **PRHNP-B06** — A spec the plan SEEDS and the progress does not list fails: a gate that only looks inside the file never sees what is missing from it.

- **PRHNP-B07** — A checkbox item promising no file at all — the untouched template marker — fails: an eternal open box makes the plan look unfinished forever.

- **PRHNP-B08** — The ticked-but-absent finding is reported FIRST, because it is the costliest damage: whoever reads the board decides on it.

- **PRHNP-B09** — An item in prose, citing no path, is not charged — there is nothing to confront.

- **PRHNP-B10** — A progress file that agrees with the disk on every item PASSES.

- **PRHNP-B11** — The verdict NAMES each offending path, so the reader does not have to diff the file against the disk by hand.

- **PRHNP-B12** — The verdict carries only the directions that found something: a direction with nothing to accuse adds no heading, because a "0 item(s)" finding reads as a defect that is not there.

- **PRHNP-I01** — The companion's path has ONE definition, derived from the scanner. The scanner is what must keep the file out of the map; a second constant here would diverge from it in silence, and the gate would confront a file the scanner indexes, or hunt for one that does not exist.

- **PRHNP-I02** — A seed is matched by PATH, never by the item's text. The wording can be shortened, translated or moved between phases without ceasing to be the same item — what identifies it is the file.

- **PRHNP-I03** — A spec mentioned in the plan's PROSE is not a seed. The checkbox is the promise; a revision's prose speaks of what already exists, and charging it accused progress files of not listing specs they did list.

- **PRHNP-I04** — A template file is never a seeded spec. The mould is the shape work is poured into, not work to be done.

- **PRHNP-X01** — Does not charge the EXISTENCE of the progress file.

- **PRHNP-X02** — Does not judge the CONTENT of an item beyond the path it cites.

- **PRHNP-X03** — Does not put the progress file into the map, nor demand a justification for changing it.

- **PRHNP-E01** — REF[PRHNP-B02]: a companion progress file that cannot be read is answered as the missing one of B02: skipped, not failed


## PRJTS — ProjectTests — the gates read the project's tests through the source the project declares

The gates that read a test's title — `feature-test-match`, `test-traceable`, `scenario-coverage` and
`flag-covered` — get the project's tests from here. How a test is written belongs to the project's
test library, not to the engine: the project declares it in `dialect.tests` (a pattern or a script, see
the `testlist` package), or its dialect family does. The gates used to carry Jest's and Go's call
syntax themselves, and a project on any other library had no title read at all.


- **PRJTS-B01** — The source is the project's own `dialect.tests`, or its family's when it declares none; with neither, nothing is declared and the gates fall back to what they can do without titles.

- **PRJTS-B02** — A pattern reads the map's test files, and only those.

- **PRJTS-B03** — The tests are read once per scan: the same root, map and configuration reuse the reading, and a new map reads again.

- **PRJTS-B04** — The tests of a set of files come in the order the files are given, and in each file in the order of their lines; tests of other files are left out.

- **PRJTS-B05** — An error of the source is returned with the declaration, never taken for a project without tests.

- **PRJTS-B06** — Support files are not read as tests, and are dropped from any list of test paths the gates confront.

- **PRJTS-B07** — Under `--index` a test's line is a line of the content the gates read: a pattern scans the files through the current source, and what a script lists of a file whose tree holds something else than the source is left out. (`projectTests`, `sameAsSource`)

- **PRJTS-E01** — REF[PRJTS-B05]: the source fails or answers outside its contract


## PRGTP — PromotableGates — identifies clean informative gates ready for promotion to blocking

Identifies informative gates whose evaluation is completely clean in a given execution profile, making them candidates for promotion to blocking gates (QUALITY section 7). When a gate is newly introduced in a project, it starts as informative so that existing debt does not block work. As the project matures and satisfies the gate's conditions, the gate remains informative unless promoted in `anchors.yaml`, measuring without actively defending the codebase. Promotion remains a human decision; this unit discovers clean gates and surfaces them as reminders directly in workflow commands (`check`, `status`, `next`).


- **PRGTP-B01** — Returns an empty list of `Promotable` gates when the evaluation profile contains no gates.

- **PRGTP-B02** — An informative gate with one or more passes and zero failures is included in the promotable list.

- **PRGTP-B03** — An informative gate with one or more failures is excluded from promotion.

- **PRGTP-B04** — An informative gate with zero passes and zero failures is excluded from promotion as having no measurement data.

- **PRGTP-B05** — A blocking gate is excluded from promotion suggestions because it already actively defends.

- **PRGTP-B06** — Each returned `Promotable` struct populates `Gate` with the gate identifier declared in `anchors.yaml`.

- **PRGTP-B07** — Each returned `Promotable` struct records `Passou` equal to the count of nodes passed by that gate.

- **PRGTP-B08** — `PromotableGates` returns promotable candidates sorted deterministically in alphabetical order of gate names.

- **PRGTP-B09** — Multiple clean informative gates in the profile are all collected into the returned promotable list.

- **PRGTP-I01** — An informative gate is candidate for promotion if and only if it is non-blocking, has zero failures, and has at least one pass.

- **PRGTP-I02** — A gate with zero passes is never classified as clean because lack of measured data must not simulate compliance.

- **PRGTP-X01** — Does not automatically promote gates or modify `anchors.yaml`.

- **PRGTP-X02** — Does not evaluate gate execution results directly from disk or runner outputs.


## PCBPR — ProofCrossesBoundary — when a rule claims a relation, the proof must reach the other side

Confronts a spec against the question the whole unit leaves open: **a rule says it mirrors
another unit — does the governed code actually IMPORT that unit, or is the claim prose while
the proof stays local?**

The measurement comes from an audit of 51 spec×code divergences in a reference app
(2026-08). The three GRAVEST findings had the same shape: two sides defined the same
thing, each side had its own test, and each test confronted its OWN copy.

| finding | the divergence |
| --- | --- |
| RFB codes | `12` = "Terreno" in the app, `12` = "Casa" in the backend — and the document filed with the tax authority carried the wrong heading |
| balance filter | `=== 'statement'` in the app, `!== 'invoice'` in the backend |
| minimum boundary | `<=` on two screens, `<` on the third and in the audit |

In every one of them 53 gates stayed green — correctly, by the rulers they had. The unit
was complete: rule declared, scenario written, test whose title matched. What no gate
asked was whether the PROOF reaches the other side.

The case that motivated the design is ALIVE in the repository and has not yet diverged:
the seat-price rule declares "the price mirrors the backend `orgBilling.ts` — diverging here lies
about the billing", the scenario says "the two values are the same the backend charges",
and the test does `expect(SEAT_PRICE.individual).toBe(15)`. It proves it is 15; it does
not prove it is the same the backend charges. Changing the backend to 18 keeps everything
green.

**What separates it from its neighbours:** `unit-complete` asks whether the test exists,
`feature-test-match` confronts scenario against test case, `rule-implemented` asks whether
the rule reached the code — all three answer "yes" about a rule whose proof never leaves
its own file. And like `dependency-honored`, this gate charges only what the spec DECLARED:
it does not go hunting duplicated concepts across the project.


- **PCBPR-B01** — An artifact that is not a spec leaves without a verdict: the gate has no jurisdiction over code, test or feature.

- **PCBPR-B02** — A rule with the declared single-source mark and a cited file whose governed code does NOT import it fails — the claim is prose and the proof is local.

- **PCBPR-B03** — A citation that lives only in a COMMENT does not satisfy the charge: comments are stripped before the code is read, because counting them would approve exactly the case that motivated the gate.

- **PCBPR-B04** — Governed code that really imports the cited unit passes.

- **PCBPR-B05** — An import written through an alias different from the path in the spec still matches: the comparison is by module base name without extension, because alias and extension vary per project and the module name does not.

- **PCBPR-B06** — A rule carrying the declared waiver with a written reason is not charged: the relation exists and is not importable (network contract, generated file, value living in an external provider).

- **PCBPR-B07** — A rule carrying the OWNER stamp is not charged: the owner does not mirror anybody — it IS the source, and there is nothing to import.

- **PCBPR-B08** — A relation claimed in PROSE, with no declared mark, is reported as a suspicion — it teaches the convention instead of barring the delivery on the first encounter.

- **PCBPR-B09** — The declared mark WITHOUT a target on the line is reported too: the mark says "I mirror something", and with no something there is nothing to confront.

- **PCBPR-B10** — A target declared by RULE CODE is resolved through the map to the files of that unit, so the rule carries stable identity instead of a path that moves.

- **PCBPR-B11** — A rule code that resolves to no unit is reported as unresolved, not silently dropped: a rule pointing at nothing warns nobody.

- **PCBPR-B12** — A rule code of the unit ITSELF is not a target: self-reference is not the other side of a boundary.

- **PCBPR-B13** — A rule with no relation claim at all leaves without a verdict — there was nothing to charge.

- **PCBPR-B14** — Imports in other language shapes (`require`, `from `, `use `, `using `, `#include`, or the project's configured pattern) satisfy the charge the same way.

- **PCBPR-E01** — A spec is confronted with no map built.

- **PCBPR-E02** — A marked rule demands an import, and every file the spec governs (`specifies`) is gone from disk or unreadable.

- **PCBPR-I01** — The claim and the target must be on the SAME rule line. An assertion floating in a surrounding paragraph binds the demand to prose, not to a catalogued rule.

- **PCBPR-I02** — The import is charged on the GOVERNED CODE, never on the test. A test that exercises the real function without importing the unit itself would be a false positive.

- **PCBPR-I03** — A satisfied charge does not swallow a pending suspicion: a file may carry one marked rule that passed and another in prose that nobody confronts.

- **PCBPR-X01** — Does not hunt for duplicated concepts across the project.

- **PCBPR-X02** — Does not compare the VALUES on the two sides.

- **PCBPR-X03** — Does not decide whether an unmarked prose claim blocks.


## RFRSR — RefResolves — the reference points at the spec that REALLY describes the unit

Confronts an artifact against the spec co-located with it: **the reference field is
filled, but does it name the right owner?**

The header gate verifies that the field EXISTS. Nobody verified that it points at the
right place — and a wrong reference is worse than a missing one: it looks like
traceability, the gate goes green, and the whole unit is attributed to the wrong spec.
Every relational gate that depends on that edge starts confronting the wrong pair, in
silence.

The characteristic failure mode is the REFACTORING nobody propagated. Measured on a real
project: 49 model files still referencing the identity from back when all the models lived
in a single file. After the split, each one got its own spec, and not one reference was
updated. The 49 kept pointing at the whole schema, and nothing raised a hand.

**The ruler**: if a SIBLING spec exists — the one the Structure co-locates with this file —
the reference must be that spec's own identity. With no sibling spec the gate goes quiet:
charging the absence of the piece is the unit gate's job, and two gates accusing the same
defect become noise.


- **RFRSR-B01** — A reference that does not match the sibling spec's identity FAILS — the refactoring that nobody propagated.

- **RFRSR-B02** — The failing verdict names BOTH sides: what is written and what the sibling spec declares, so whoever reads it does not have to open two files.

- **RFRSR-B03** — A reference equal to the sibling spec's identity passes.

- **RFRSR-B04** — A spec is not confronted: it OWNS an identity, it does not reference one.

- **RFRSR-B05** — An artifact with no reference declared leaves without a verdict — the absence is the header gate's charge, not this one's.

- **RFRSR-B06** — With no sibling spec on disk the gate goes quiet: the missing piece is the unit gate's charge, and two gates on one defect become noise.

- **RFRSR-B07** — A sibling spec that exists and declares no identity of its own counts as no sibling: there is nothing to compare against.

- **RFRSR-B08** — The sibling is found by NAME convention — same stem, same directory, spec suffix.

- **RFRSR-B09** — The test suffix of each supported language is stripped before the stem is computed, so a test file finds the same sibling its code does.

- **RFRSR-B10** — An intermediate extension is dropped from the stem too, so a file that carries a kind in its name still lands on the unit's spec.

- **RFRSR-B11** — The reference is read from the header whatever the comment syntax of the language — the three families of line marker are accepted.

- **RFRSR-B13** — A leading dot is not a stem separator: a hidden file keeps its whole name, so it is never attributed to a spec named only by the suffix.

- **RFRSR-B12** — The accepted identity length comes from the project's Structure, read at confrontation time and not frozen at process start.

- **RFRSR-B14** — A reference to a code that exists nowhere in the project fails when no sibling spec exists, reporting the unknown code.

- **RFRSR-B15** — Without a sibling spec on disk, an existing declared code in the graph skips.

- **RFRSR-B16** — An inferred identity (`CodeDeclarado: false`) does not satisfy the reference.

- **RFRSR-B17** — Without a graph, absence is not asserted and the gate skips.

- **RFRSR-I01** — The sibling spec on disk is the primary ruler. When a sibling spec exists, the verdict depends on the filesystem alone. When no sibling spec exists, the graph is consulted solely to verify that the cited reference exists as a declared identity.

- **RFRSR-I02** — The gate never writes and never repairs. It reads the artifact and the sibling and returns a verdict; a gate that fixed what it points at would pass on the second run.

- **RFRSR-X01** — Does not charge the ABSENCE of the reference field.

- **RFRSR-X02** — Does not charge the absence of the sibling spec.

- **RFRSR-X03** — Does not consult the map when a sibling spec is present on disk.

- **RFRSR-X04** — Does not judge whether the sibling spec DESCRIBES the unit well.

- **RFRSR-E01** — REF[RFRSR-B06]: a sibling spec that cannot be read is answered as the missing sibling of B06: the gate goes quiet


## RPHRG — RegionPairHonored — every opened source region must close with its own identity code

Confronts source code regions against their pairing markers: **every opened code region must close,
and must close specifying its own exact identity code.**

A region marker pair gives identity a concrete line interval within the source file. Without region
boundaries, one only knows that a given file realizes a requirement, but not WHERE within the file that
implementation resides (TRACEABILITY §3). The region pair is the single fragile component of this mechanism,
and its fragility is particularly perilous because it raises no compiler or syntax errors: closing in the
wrong location produces an interval that is syntactically VALID and completely WRONG.

This gate detects the three critical structural pairing defects that cannot be reliably spotted by human
review when reading large source files top-to-bottom:
1. **Unclosed region (`sem-fecho`)**: A region opened and never closed, causing the requirement interval
   to bleed silently to the end of the file.
2. **Orphan close (`fecho-orfao`)**: An end region marker appearing without a preceding opening marker,
   typically left behind from careless cut-and-paste refactorings.
3. **Mismatched close (`fecho-trocado`)**: An end marker carrying a DIFFERENT identity code than the currently
   open region, indicating inverted or twisted nesting.

The third defect is the exact reason why region close markers must explicitly carry the identity code rather
than using an anonymous end marker. With anonymous end markers, the count of open and close markers would
balance perfectly, the gate would remain green, and the freshness stamp would measure the NEIGHBOR interval:
requirement A would appear to change when neighbor B was edited. A silent, crossed error is the worst kind
of defect: it vanishes during initial measurement and resurfaces weeks later as an erroneous conclusion.

Unlike gates that report non-blocking warnings for domain ambiguities (such as an uncatalogued requirement
letter), this gate issues a definitive blocking **Fail**. A malformed region pair has no two legitimate
interpretations: it is an objective syntax defect in code structure with a single unambiguous, local fix.

Finally, the ABSENCE of region markers is NOT a defect. Granular region delimitation is entirely optional;
when a file does not define regions, the whole file revision applies cleanly.


- **RPHRG-B01** — When the confronted node is neither code nor test, the gate skips confrontation.

- **RPHRG-B02** — When a code or test file contains no region markers, the gate skips without asserting failure.

- **RPHRG-B03** — When all opened regions close with their matching identity code, the gate passes.

- **RPHRG-B04** — An opened region that is never closed fails, reporting the line and missing close marker.

- **RPHRG-B05** — An end region marker without an opening marker fails as an orphan close.

- **RPHRG-B06** — An end region marker closing with a different code than the open region fails, citing both codes and line number.

- **RPHRG-B07** — Multiple pairing errors within a file are ordered sequentially by line number in the verdict.

- **RPHRG-I01** — Region absence is never charged as a failure. Delimitation is optional and files without regions evaluate at whole-file scope.

- **RPHRG-I02** — Pairing defects always result in a blocking Fail verdict, never Pending, because malformed nesting has only one correct repair.

- **RPHRG-I03** — End markers must explicitly match opening codes to prevent silent interval cross-over between neighboring requirements.

- **RPHRG-X01** — Does not mandate region markers in code or test files.

- **RPHRG-X02** — Does not inspect region markers inside specs or markdown files.

- **RPHRG-X03** — Does not enforce semantic validity of the code enclosed within a region.


## RVMTR — ReverseMatch — every scenario still has its rule, and every proven code still has its scenario

The forward gates of the unit walk from the origin to the destination: every rule needs a scenario,
every scenario needs a test. Neither walks back from the destination to ask whether the origin still
exists. The cost was measured: a revert deleted a rule from the spec, the code and the test, while the
feature kept its scenario (a parallel change reintroduced it with no conflict). The scenario went on
asserting a behaviour nobody decided any more, the test stayed green proving a rule nobody declared,
and every gate was green.

This unit holds the return leg of each pair, as two gates of their own, so that whoever reads the
accusation knows it is about the destination:

- feature-spec-match confronts a feature: does every scenario code of this unit correspond to a rule
  that a spec covering this feature still defines?
- test-feature-match confronts a test: does every code the test names correspond to a scenario that a
  feature it exercises still declares?

Both were tuned against false accusations measured in a real project, because a gate that accuses
what is right teaches people to ignore it: numbered variants of one rule, data states defined by name,
the visual baseline, revision codes, and other units' codes cited to build fixtures are not orphans.


- **RVMTR-B01** — The gate skips a node that is not a feature, and a feature that declares no coded scenario.

- **RVMTR-B02** — Without a map, both gates answer Pending.

- **RVMTR-B03** — A feature no spec covers is Pending: the existence of the spec is charged by co-location, not here.

- **RVMTR-B04** — When the covering specs define no requirement and no data state, the gate is Pending.

- **RVMTR-B05** — A scenario code of this unit that no covering spec defines fails the gate, and the message names the code as written in the feature; when every code has its rule, the gate passes.

- **RVMTR-B06** — A numbered variant of a scenario code is the same rule as the code without it; a variant of a rule the spec does not define is still an orphan.

- **RVMTR-B07** — A data state cited in the feature counts as defined when a covering spec defines a state of that name, with or without the unit prefix; a state no spec defines is an orphan.

- **RVMTR-B08** — The unit's visual baseline code is never charged here.

- **RVMTR-B09** — The rules of every spec covering the feature count together: a scenario defined by any of them has an owner.

- **RVMTR-B10** — The gate skips a node that is not a test, and a test no feature exercises, saying the link is `unit-complete`'s to charge: there is nothing to confront.

- **RVMTR-B11** — When the features it exercises declare no coded scenario, the gate is Pending.

- **RVMTR-B12** — A code the test names that no exercised feature declares as a scenario fails the gate, naming the code.

- **RVMTR-B13** — A test that names a rule bare where its feature declares that rule only as numbered variants fails, naming it: each variant is proven on its own, and a proof of the bare rule proves none of them. A feature that declares the bare scenario too takes it.

- **RVMTR-B18** — A test that names a variant its feature does not declare — `CODE-B03#01` under a feature declaring `@CODE-B03` alone — fails, naming it: the proof counts by the scenario it names, and this one no feature declares.

- **RVMTR-B14** — A revision code named by the test is not read as a rule and is not charged.

- **RVMTR-B15** — A data state a spec defines with the unit prefix is read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, `TREXXXX-DS-data-present` defines `DS-data-present`.

- **RVMTR-B16** — A support file is not confronted as a test (Skip, saying why): it proves no scenario, so there is no feature to link it to.

- **RVMTR-B17** — A test with no feature linked whose units under test — the code the project's derivation says the test belongs to — all lie in layers of `regime: declarativo` is skipped, saying why; a test with no unit found, or one of whose units is governed, stays pending.

- **RVMTR-I01** — Neither gate answers Pass when it had nothing to match against: no map, no linked origin, or an origin that declares nothing answers Pending.

- **RVMTR-X01** — Codes of another unit are never charged: only this unit's codes, and in a test only the units the exercised features govern.

- **RVMTR-X02** — A code the test names only in a comment is not a claim of proof.

- **RVMTR-E01** — A linked spec or feature cannot be read.


## RVDUR — ReviewsDue — the targets of a reviewed gate that no review covers at their current revision

A gate that declares `review:` marks its targets to review, apart from how it measures: a judgment gate still asks its question and stamps its verdict, and the same targets are to review until a reviewer looks. A target is to review when the gate applies to it and no review is recorded at its current revision. The list informs and never blocks; the project decides when the reviewer comes.


- **RVDUR-B01** — The targets to review (`ReviewsDue`) are, for every gate that declares `review:`, each node the gate applies to with no review recorded at its current revision, with the gate's review question (`ReviewAsk`); ordered by gate and target.

- **RVDUR-B02** — A gate with no `review:` has nothing to review, and naming a gate lists that gate's alone.

- **RVDUR-B03** — A change to a reviewed target makes it to review again.

- **RVDUR-I01** — Listing what is to review changes no verdict: a judged and reviewed gate's judgment is asked with or without its review.

- **RVDUR-X01** — The list never blocks.


## RVORP — RevisionOrphans — the rules a revision changed the meaning of, without saying so

Confronts a revision declared inside a spec against the OTHER rules of the same unit: **when a revision rewrites a rule, the sibling rules that speak of the same thing must be named — either as revised too, or as read and still valid.**

A rule does not live alone. It shares vocabulary with its siblings, and it is that vocabulary a revision changes — not only the text of the rule it rewrites. A revision that abolishes a concept leaves every other rule still asserting it, and the spec then claims two contradictory things at once.

Measured in the reference app, and it is what produced this gate. The `NTCNN-R0002` changed the notification badge from a COUNT to a DOT, and named the rules it rewrote: `B03`, `B04`, `B07`. The invariant `I02` was not named — and it is titled *"the badge never COUNTS what the list does not show"*, with a body reading *"the NUMBER on the bell matches what appears on opening"*. The invariant governed arithmetic the revision had abolished.

**Nothing accused it.** The `B03` was correct, the `I02` was well-formed, the unit complete, the suite green. The contradiction surfaced MONTHS later, when another agent went to implement and could not tell which of the two to follow — and it became a decision that had to escalate to the user, with nobody left remembering the context. Seven contradictions of this exact shape surfaced in a single batch.

The ruler is ONE shared domain word, and the threshold was measured in both directions. Requiring two failed the very case that produced the gate: the `I02` shares exactly one word with what the revision rewrote — `badge` — and that word carries the whole contradiction. What makes one word enough is CLEANING the title, not counting: with negations and waiver comments in, the same spec accused three rules (`B01` and `B05` entered on "não" alone); with them out, it accuses one — the target.

A ubiquity filter was also tried, discarding terms appearing in more than two thirds of the rules. It discarded `lista` and `badge` — precisely the subject — and made the real case accuse nothing. In a well-written spec the domain vocabulary repeats on purpose: discarding what repeats is discarding the subject.

The ruler is CO-CITATION, not meaning. Asking whether two rules contradict each other requires reading them, and that is judgement — it would make this a judge, not a gate. What a machine decides alone is narrower and sufficient: which rules of this unit share the vocabulary of the rules the revision touched, and were not mentioned.

Leaving the finding is cheap, and deliberately so: `Checked: I02` asserts that somebody read it, not that it is correct. It is cheap to write AFTER reading and impossible to write honestly without reading, and that asymmetry is what makes the ruler work. Demanding that every revision check every rule of the unit would fail always on a nine-rule spec, and satisfying it would become theatre — the whole list pasted in unread.

This gate operates in distinct territory from neighbouring gates:
- Unlike `plan-revised`, which governs the relationship between two distinct PLANS across the map, this gate stays inside one spec, between a revision and the sibling rules of its own unit.
- Unlike `plan-change-justified`, which asks whether an edit carries a written justification, this gate asks whom that justification forgot.
- Unlike `spec-feature-match`, which confronts declared requirements against written scenarios, this gate never leaves the spec: both sides of its confrontation are rules of the same file.


- **RVORP-B01** — When the confronted node is not a spec, the gate skips confrontation.

- **RVORP-B02** — When the spec declares no revision, the gate skips confrontation.

- **RVORP-B03** — When a revision declares no `Revises:`, the gate abstains with a pending verdict: the field is new, and 439 revisions written before it exist in the reference app — accusing all of them at once produces noise, not a queue.

- **RVORP-B04** — When a revision names a rule the spec does not define, the gate fails, reporting the unknown code.

- **RVORP-B05** — When a sibling rule shares significant vocabulary with a revised rule and appears in neither `Revises:` nor `Checked:`, the gate reports it as an orphan, naming the shared terms.

- **RVORP-B06** — When every vocabulary-sharing sibling appears in `Revises:` or `Checked:`, the gate passes.

- **RVORP-B07** — When a rule appears in `Checked:`, it leaves the accusation without asserting that it is correct — only that somebody read it.

- **RVORP-B08** — A rule's title is its heading — or, with none, its first definition —, never a later line that names it: a row of a usage table lists what the rule reads, and taken as the title it made every rule reading the same field a sibling. (`ruleTitles`)

- **RVORP-B09** — A revision whose `Revises:` declares, with its reason, that it revised no rule (`none — <why>`, in any supported language) is an answer: with no code revised and such a declaration, the spec passes; a declaration with no reason declares nothing, and its reason is not read for codes. (`revisedCodes`)

- **RVORP-I01** — A rule never accuses itself: the revised rule is excluded from its own orphan candidates.

- **RVORP-X01** — Negations and waiver comments are stripped from a rule title before comparison: they are not what the rule asserts.

- **RVORP-Q01** — Should the gate block or inform? Decided by the user: it blocks — became `DFGTD-B12` — the default-gates catalogue, where the class is decided.


## RVRNR — RevisionRenumber — the revisions a branch added move to a free number when the base took theirs

A revision is numbered by the file it revises: `PRICX-R0003` is the third recorded change of
`PRICX`. Two open pull requests that revise the same spec each take the next free number — and
both take the SAME one. The second to merge lands a `R0003` that already means something else
on the base: two revisions answering to one code, every citation of it ambiguous, and
`plan-change-justified` seeing the count and the highest number disagree.

The number stays, because it is what says "changed three times". The collision is resolved
where it happens: on the branch that arrives second, at rebase. This unit is the engine behind
`anchors renumber` — text in, text out; the command reads the three versions from git and writes
the result.

**Only what the branch added moves.** A revision the base already has is history other people
have read and cited; renumbering it would point every one of those citations at the wrong
change. The branch's revision is the one nobody outside the branch has seen yet. The same
reasoning decides which citations move: only those on a line the branch added.


- **RVRNR-B01** — `PlanRenumber`: a revision the branch added whose number the base already uses moves to the next free number.

- **RVRNR-B02** — After a rebase, with the base's revision and the branch's sharing a number in the same file, only the branch's moves.

- **RVRNR-B03** — When one added revision collides, every revision the branch added for that code moves, in the order of its old number, after the highest number in use — so the sequence stays in order.

- **RVRNR-B04** — A revision the branch added whose number the base does not use stays as it is.

- **RVRNR-B05** — `RewriteRevisionCitations`: a citation is rewritten only on a line the branch added; a line that existed at the merge base keeps its code.

- **RVRNR-B06** — A branch that adds the same number twice is refused, naming the code: a citation of it could not be told apart.

- **RVRNR-I01** — A revision the base has is never renumbered — including one whose explanation the branch edited.

- **RVRNR-I02** — The rewrite is one pass: a chain of renames (`R0003→R0004`, `R0004→R0005`) never cascades.

- **RVRNR-X01** — Does not rewrite a revision cited without its unit code.


## RTDCL — RouteDeclared — a screen declares how one arrives, and names its neighbours

Confronts a SCREEN spec against the navigation graph it belongs to: **is there a named
route that reaches this screen, and do its edges point at concrete screens?**

It is the navigation traceability ruler. Without the route the screen is a loose node —
something exists that nobody can reach. With generic terms in the navigation tables
("Next screen", "Main menu") the edge points nowhere: it reads like a link and connects to
no node at all, which is worse than an absent edge, because an absent edge is visible and
a vague one passes for a declared one.

**Jurisdiction is the whole design.** Only `layer: screen` is charged. Hooks, business
logic, stores and DAOs have no route, and charging them was the vice of the legacy
validator this gate replaces — it knew only screen and component, so it treated every
non-component as a screen and produced a false positive for every unit that legitimately
has no route. A gate that cries over what cannot be fixed teaches the team to ignore it.

That is also why the Skip carries a stated REASON. A bare indeterminate count leaves the
reader wondering whether the silence is their problem; saying "this is not a screen"
closes the question in one line.


- **RTDCL-B01** — An artifact whose layer is not `screen` leaves without a verdict, and the verdict SAYS why — a bare indeterminate count would leave the reader wondering whether the silence is their problem.

- **RTDCL-B02** — A screen with no named route FAILS: without it the screen is a node nobody can reach.

- **RTDCL-B03** — A screen whose navigation table carries a generic term FAILS — an edge that names no concrete screen points nowhere. Every navigation section is read, one right after another included; a section is navigation only when its heading names it as a whole word (`### In`, not `### Integração`).

- **RTDCL-B04** — A screen with a named route and concrete neighbours passes.

- **RTDCL-B05** — Route and navigation are recognised in either declared language, so a project writing in its own language is measured and not silently approved.

- **RTDCL-I01** — The header's declared layer is the source of truth of identity, and the node's tags are only the fallback when the header is silent. The header is what the author wrote on purpose; a tag can come from inference.

- **RTDCL-I02** — Only TABLE ROWS of the navigation sections are confronted, never surrounding prose. An explanatory sentence mentioning "the next screen" is the author writing, not an edge being declared.

- **RTDCL-X01** — Does not charge a route from any layer other than `screen`.

- **RTDCL-X02** — Does not verify that the declared route EXISTS in the router.


## RTEXR — RouteExists — declared route in specification must exist in application route registry

Confronts the route declared by a specification against actual application code: **the route that the
specification declares must exist where the application registers its routes.**

The neighbouring `route-declared` gate confronts a specification against itself — does the specification
declare a route? This gate confronts the specification against the CODE: does that route actually exist
where the application registers its routes?

The critical distinction emerged in a real end-to-end delivery defect. A specification for a new screen
declared `> **Rota**: MetadataEdit` (or `route: MetadataEdit`), and another specification promised navigation
to it. The blocking `route-declared` gate gave a green checkmark to both specifications because both declared
routes. Yet the route existed nowhere in the application code. Both specifications described a path leading to
an unreachable screen, and the entire verification pipeline remained green.

This represents an anchor that lies in the most elusive way possible: nothing is missing and every document
references each other, but nobody asked whether the destination actually exists in the code.

Measured across 96 screen specifications in a real project prior to enabling this gate: 2 findings, both
true defects (the aforementioned unreachable end-to-end route, and another specification whose screen the
application registered under a different name). Zero false positives.

When a specification declares no route, this gate skips confrontation, as enforcing route declaration is the
exclusive responsibility of `route-declared`. Furthermore, when project configuration does not specify where
routes are registered (`route_registry`), or when zero routes can be extracted, the gate returns **Pending**
rather than falsely approving uninspected routes.


- **RTEXR-B01** — When the confronted node is not of kind spec, the gate skips confrontation.

- **RTEXR-B02** — When the specification does not declare a route, the gate skips confrontation.

- **RTEXR-B03** — When project configuration defines no route registry glob patterns, confrontation returns Pending.

- **RTEXR-B04** — When route registry glob pattern is invalid, the gate returns Pending reporting the glob error.

- **RTEXR-B05** — When route registry files are inspected but zero registered routes are discovered, the gate returns Pending.

- **RTEXR-B06** — When a declared screen name matches a component navigation prop in registered route files, the gate passes.

- **RTEXR-B07** — When a declared screen name matches a navigation stack parameter type entry in registered route files, the gate passes.

- **RTEXR-B08** — When a declared backend route matches an added HTTP resource registration, the gate passes.

- **RTEXR-B09** — When a declared route includes an HTTP method verb prefix, the gate strips the verb and passes if the route path exists.

- **RTEXR-B10** — When a declared route matches a registered route path regardless of leading slash differences, the gate passes.

- **RTEXR-B11** — When project configuration specifies a custom route pattern regex, routes matching that pattern pass.

- **RTEXR-B12** — When a custom route pattern regex is invalid or contains no capture group, the gate returns Pending.

- **RTEXR-B13** — When a declared route is absent from all registered routes, the gate fails citing the missing route, total registered count, and searched globs.

- **RTEXR-I01** — Route presence is validated only for specifications; non-spec nodes skip confrontation.

- **RTEXR-I02** — The gate never approves route existence without inspecting registry files, returning Pending when registry configuration or routes are missing.

- **RTEXR-I03** — Route matching is slash-normalized, treating route paths with and without a leading slash as equivalent representations.

- **RTEXR-I04** — Missing routes always result in a blocking Fail verdict to prevent shipping unreachable destinations.

- **RTEXR-X01** — Does not mandate that every specification declare a route.

- **RTEXR-X02** — Does not validate route parameter schemas, HTTP payload structures, or response codes.

- **RTEXR-X03** — Does not evaluate authentication or access permissions attached to the route.

- **RTEXR-E01** — REF[RTEXR-B04]: a registry glob that does not parse is answered by B04: Pending with the glob error

- **RTEXR-E02** — REF[RTEXR-B12]: a route pattern that does not compile, or has no capture group, is answered by B12: Pending

- **RTEXR-E03** — A file of the route registry cannot be read.


## RLUEX — Rule — the identity of a verification INSIDE a gate, and the waiver that names it

A gate does not verify ONE thing. The one that confronts a spec charges two — that no
placeholder survives, and that at least one rule is catalogued. The one that confronts a
header charges several. Until this unit existed, the verdict said only WHICH GATE failed,
and that had two consequences, both measured:

- whoever waives a deliberate case — the spec that is born before the feature — could
  only waive the ENTIRE gate, and every other verification it did well went down with it;
- two different defects inside the same gate are indistinguishable in a report. The
  reader needs the MESSAGE to know which of the two happened, and a message changes with
  the next rewrite or the next translation.

So the verification gets an identity of its own: `<gate>/<rule>`. It is the same
principle the doctrine already applies to the artifact — a stable identifier is what
lets one speak of a decision without describing it again.

**What separates this unit from its neighbours**: the engine decides WHO runs and
aggregates the answers; the registry decides WHICH function answers. This unit decides
NOTHING about running — it only knows how a verification is named, and whether somebody
declared, in writing and with a reason, that this particular verification should not be
confronted this time. It reads no file, consults no map, and returns no verdict.

**The reason is part of the datum, not a comment beside it.** A waiver without a written
justification is indistinguishable from somebody fleeing a gate that found a defect. The
second measured failure is finer: a waiver restricted to targets that is asked about
WITHOUT a target must answer "not waived" — otherwise the gate is filtered out of the
list entirely, and waiving four new specs would erase the gate for the whole repository.


- **RLUEX-B01** — `NewRuleID` joins gate and rule with a separator, so that short rule names cannot collide between gates.

- **RLUEX-B02** — `NewRuleID` with an empty rule returns the gate name alone: a gate with a single verification gains no separator it does not need.

- **RLUEX-B03** — `Gate` returns the gate half of the identifier, and `Rule` the rule half.

- **RLUEX-B04** — `Rule` is empty when the identifier carries no rule half — the verdict then belongs to the gate as a whole.

- **RLUEX-B05** — `ParseWaiver` refuses an entry with no reason, and an entry whose reason is blank, because the reason is the only thing separating a deliberate waiver from an ignored gate.

- **RLUEX-B06** — `ParseWaiver` refuses an entry with no rule name.

- **RLUEX-B07** — `ParseWaiver` refuses a target that looks like a PATH, and the error names the artifact code as what to use instead.

- **RLUEX-B08** — `ParseWaiver` refuses a target marker with nothing after it, which is a typo that would otherwise produce a waiver that waives nothing.

- **RLUEX-B09** — `Waived` accepts the two granularities: waiving the gate covers every rule inside it, and waiving one rule preserves the rest of the gate.

- **RLUEX-B10** — `Waived` answers "not waived" for a waiver that DECLARES targets, so the caller cannot drop the gate from the list and erase it for the whole repository.

- **RLUEX-B11** — `WaivedTarget` waives the named codes and confronts every code that was not named.

- **RLUEX-B12** — `WaivedTarget` with no declared target waives every code, which is the coarse exit a freshly declared gate still needs.

- **RLUEX-B13** — `WaivedTarget` answers "not waived" for an artifact with no code: without identity there is no specific target to waive.

- **RLUEX-B14** — Each target carries ITS OWN reason, so two waivers of the same rule no longer make the second overwrite the first.

- **RLUEX-B15** — `WaiverFromMessage` reads the waivers declared in the commit message, which is the form that survives in the history beside the why of the change.

- **RLUEX-B16** — `WaiverFromMessage` refuses a marker whose reason is blank, by the same guarantee as the textual form.

- **RLUEX-B17** — `Merge` joins two waivers, so a pipeline hook using the environment variable and an author writing the marker can coexist.

- **RLUEX-I01** — A waiver never reaches what nobody waived. The error that would cost most here is one leaking onto a gate that was never named.

- **RLUEX-I02** — Every accepted waiver carries a written reason, in both input forms. A waiver with no why is indistinguishable from an ignored gate, and the report would show it waived without saying why.

- **RLUEX-I03** — A refusal always produces an error that reaches the caller, never a silent acceptance with an empty reason. A waiver that does not waive would fail the next commit with no visible explanation.

- **RLUEX-X01** — Does not accept a PATH as the waiver target.

- **RLUEX-X02** — Does not decide whether a rule PASSES.

- **RLUEX-X03** — Does not read files, the map, or the project structure.


## RLIMR — RuleImplemented — a spec catalogues rules, and the code shows it realized them

Confronts the spec against the code in the direction that was missing: **did the spec end up
talking to itself?**

It is the inverse of the gate that validates references. That one checks that the codes CITED by
the code exist in the spec; this one checks that the rules DECLARED in the spec got an
implementation. Without it, a spec can declare five new rules and the code gain
not a single line — with all gates green, because the spec exists, the code exists, and the
two reference each other through the header.

Measured: an interface spec gained five rules and 98 lines, and the corresponding file
had ZERO occurrence of the subject. Half of the delivery was dead code declared as
done, and none of the 26 gates asked. The defect only showed up when someone read spec and
code in the same pass.

**The ruler is the declaration, not the guesswork.** Demanding every rule be marked would be false by
construction — measured against 592 units, it would yield 3.121 findings, and not even the well-made units
would pass: among those that mark the code, none marks 100%. The reason is a good one: a constraint ("the
unit does NOT do Y") is satisfied by the ABSENCE of code, and absence has nowhere to receive a
mark. But "at least one" does not serve either — it separates those who implemented from those who did not
implement and says nothing about the other fifteen rules. So whoever writes the spec
DECLARES, rule by rule, whether it has code.


- **RLIMR-B01** — A spec whose rules do not appear in the code is ACCUSED, and the verdict names which ones were left without realization.

- **RLIMR-B02** — A rule waived with a written reason settles the account: the declaration counts as the answer, and what it waives stops being charged.

- **RLIMR-B03** — A unit that predates the practice becomes a DIVERGENCE, not a failure — while the project does not declare that it requires the marking.

- **RLIMR-B04** — Once the requirement is declared in the Structure, the divergence becomes a failure: it is the act of saying "the migration ended here".

- **RLIMR-B05** — A spec with no linked code is not this gate's subject: without the piece on the other side there is no confrontation to make.

- **RLIMR-B06** — The waiver may name the rule it covers, or count for all of them when it names none.

- **RLIMR-B07** — With `data_states.required`, the code of the unit cites each data state its spec defines, as it cites a rule; off, the data states are not asked of the code.

- **RLIMR-I01** — Requiring the marking never punishes whoever already marks. Whoever did the work before the requirement cannot fail for having done it.

- **RLIMR-I02** — Identity survives the rename: code marked with the previous name keeps counting. Losing the mark in a rename would turn identity stability into new debt.

- **RLIMR-X01** — Does not judge whether the implementation is RIGHT — only whether it exists and declares itself.

- **RLIMR-X02** — Does not require a mark on EVERY rule.

- **RLIMR-E01** — The code file the spec describes is on disk but cannot be read (no read permission).


## RLTYR — RuleTypes — the rule VOCABULARY is extensible, but it must be DECLARED

Confronts a spec against the alphabet that makes traceability possible: **each letter of a
rule code (`{CODE}-<letter><NN>`) is the initial of the term that names a section — and a
letter nobody declared is INVISIBLE.**

That is the worst kind of hole, and it is the reason the gate exists. The rule appears in
the spec, in the `.feature` and in the test, and even so `feature-test-match` does not see
it: the code regex does not match the letter. **It looks covered and is not** — a green
that certifies a link nothing ever traversed.

It confronts three things, in this order of gravity: an UNDECLARED LETTER used in the
file; a SECTION cataloguing rules under a title no letter claims; and a CONFLICT in the
vocabulary itself, where two different sections claim the same letter. A fourth finding —
a section the project declared as rule-cataloguing but filled WITHOUT a code — is
Pending rather than Fail, because it is the inverse gap and the ruler is opt-in.

With no vocabulary declared the gate does NOT go quiet: it confronts the CANONICAL
letters. It used to Skip, and the effect was a canonical gate — seeded by `init` in every
project — that measured nothing: it took a row in the `check` table and reported
indeterminate forever. That is the same silence Anchors fights everywhere else, a
declared gate that confronts nothing giving the impression of a defence that does not
exist.


- **RLTYR-B01** — A letter used in the spec and NOT declared in the vocabulary fails, and the verdict names the letter.

- **RLTYR-B02** — A declared letter used under a claimed section passes.

- **RLTYR-B03** — A section that CATALOGUES rules under a title no letter claims fails, and the verdict names the section.

- **RLTYR-B04** — A CONFLICT in the vocabulary itself — the same letter claimed by two different terms — fails before the file is even read: each letter belongs to ONE term.

- **RLTYR-B05** — With no vocabulary declared the gate confronts the CANONICAL letters instead of going quiet.

- **RLTYR-B06** — A heading that IS the rule code itself is not a category section: it is the rule's own header, and is not charged.

- **RLTYR-B07** — A section that merely CITES codes belonging to other sections does not catalogue rules, and claims no letter.

- **RLTYR-B08** — A section that DEFINES a code in the first cell of a table does catalogue rules, and is charged.

- **RLTYR-B09** — A section the project declared as `sections_require_code` that is FILLED and carries no code returns a divergence, naming the section.

- **RLTYR-B10** — A section whose table already carries the code is not charged.

- **RLTYR-B11** — A declared section that is NOT in `sections_require_code` is not charged: it merely enumerates values, and demanding a rule of an index would invent a duty.

- **RLTYR-B12** — A project that does not use `sections_require_code` changes no behaviour — the ruler is born opt-in, or it would accuse an entire existing base at once.

- **RLTYR-I01** — A spec with no rule code at all is not this gate's problem. Demanding a catalogue here would duplicate `spec-complete`, which is who charges the existence of a catalogued rule.

- **RLTYR-I02** — The verdict NAMES the letter AND where to declare it. A gate that fails without saying what transfers the diagnostic work to whoever reads it.

- **RLTYR-I03** — A section without a code borrows a neighbour's. Measured in a real project: 48 specs with "Events / Callbacks" filled in and not one code, and the scenarios proving those events borrowed the code of the neighbouring state — an `-S` governing behaviour.

- **RLTYR-X01** — Does not decide WHICH letters exist.

- **RLTYR-X02** — Without a declared vocabulary it charges only the letter, not sections or terms.

- **RLTYR-X03** — Does not judge whether the letter is the RIGHT one for that rule.

- **RLTYR-X04** — Does not charge format — only traceability.


## RLUSG — RuleUses — each rule says what it uses, and what it uses exists

A spec has the two ends — the rules (behaviours, states, errors…) and the data (the data contract,
the domain, the props) — and nothing between them: what each rule reads lived, when it lived
anywhere, in prose. A field changed in the contract pointed at no rule, and a rule could check a
datum the contract never offered.

Three sections tie them: **Validations** (`V`: a condition on a datum → behaviour), **Presentation
validations** (`P`: a prop or state → appearance) and **Rule uses** (any other rule → what it
uses). Each row starts with the rule's code and goes on with what it uses; the row is read by
position, since the columns' names follow the project's language. Two gates read them:
`rule-uses-declared` asks every rule to say what it uses, and `rule-uses-resolve` asks that what
it uses exists in the spec.


- **RLUSG-B01** — The three sections are found by their title in any supported language, or by the project's own title in `section_titles`; each row whose first cell carries a rule code gives that rule, and its second cell what it uses; a row whose uses are still a TODO says nothing yet.

- **RLUSG-B02** — A uses cell lists its backticked items when it has any, and its comma-separated items otherwise.

- **RLUSG-B03** — `rule-uses-declared` fails naming every rule of the letters asked about that has no row in the three sections; the gate entry's `letters` chooses the letters, and without it every letter but an open question's, a plan phase's and a flag scenario's is asked about; a rule whose line carries `@no-uses: <why>` is waived.

- **RLUSG-B04** — `rule-uses-declared` skips a node that is not a spec, and a spec with no rule of the letters asked about.

- **RLUSG-B05** — `rule-uses-resolve` accepts a field the spec declares — the first cell of a row of any of its other tables, or a backticked name in a heading —, a dotted or indexed name by its first segment, a `DEPn` that is a row of the dependencies table, and leaves codes to `code-reference-valid`; it fails naming each use that resolves to nothing, with its rule.

- **RLUSG-B06** — `rule-uses-resolve` skips a node that is not a spec, and a spec whose rules do not say what they use yet.

- **RLUSG-B07** — `rule-uses-implemented` fails naming each field a rule uses that no code file the spec governs mentions as a whole word — the name, or the first or last segment of a dotted one —; codes and `DEPn` are not fields; a spec with no rule uses, or governing no code, is skipped.


## SCASS — ScenarioAsserts — scenario outcome steps must assert concrete verifiable outcomes

Confronts feature files against tautological step definitions: **the outcome step of a scenario must assert an observable result, rather than merely repeating the requirement code.**

A scenario in Gherkin exists to govern the test: it commits the system to an observable outcome before anyone writes the assertion code. When an outcome step degenerates into a tautology — such as stating that the requirement code is verified — it asserts nothing. Two scenarios for entirely different requirements become indistinguishable, and the definition of what constitutes verification migrates completely into the test implementation, which is the very artifact the scenario was meant to govern.

The operational risk of tautological scenarios is structural: when nobody between the specification and the test is required to write down the expected outcome, the resulting test is written without a discriminating case. A mutation that completely removes the production logic can leave the test suite green, because the test itself was never anchored to a concrete outcome. In a measured reference project, 359 scenarios suffered from this exact defect, having been mechanically copied from templates.

This gate prevents this structural defect by isolating outcome steps and mechanically checking whether the step consists solely of a requirement identity code wrapped in linking words.

This gate operates in distinct territory from neighbouring gates:
- Unlike `feature-test-match`, which verifies scenario-to-test alignment by comparing scenario codes and titles against test titles, this gate inspects the inner substance of outcome steps to ensure they assert real behavior.
- Unlike `spec-complete`, which verifies that requirements are structurally catalogued in specifications, this gate ensures that scenarios realizing those requirements commit to verifiable outcomes.
- Unlike natural language review gates, this gate does not judge prose style; it applies a deterministic filter that flags steps where removing linking words leaves only the requirement identity.


- **SCASS-B01** — When the confronted node kind is not a feature, the gate skips confrontation.

- **SCASS-B02** — When all outcome steps in a feature assert observable results, the gate passes.

- **SCASS-B03** — When an outcome step merely asserts that a requirement code is verified, the gate fails, citing the code.

- **SCASS-B04** — Tautological outcome steps wrapped with linking words and up to two residual words fail.

- **SCASS-B05** — An outcome step citing a requirement code while asserting three or more descriptive content words passes.

- **SCASS-B06** — Multiple tautological outcome steps within a feature are all reported, deduplicated and sorted.

- **SCASS-B07** — Outcome steps written in any recognized dialect keyword are evaluated according to configured project settings.

- **SCASS-B08** — Empty lines and comment lines in feature files are ignored and never evaluated as outcome steps.

- **SCASS-B09** — Non-outcome steps such as setup and action steps are ignored even when citing requirement codes.

- **SCASS-I01** — Up to two residual content words beyond linking words around a code is classified as a tautology, preventing trivial phrasing variations from bypassing the rule.

- **SCASS-I02** — Truly assertive outcome steps are never flagged as tautologies, protecting against false positives that would encourage teams to disable the gate.

- **SCASS-I03** — Language recognition covers all supported dialect alternatives, preventing inherited multi-language features from being silently bypassed.

- **SCASS-X01** — Does not evaluate the semantic accuracy or elegance of prose beyond the mechanical removal of linking words.

- **SCASS-X02** — Does not inspect setup and action steps for requirement code references.

- **SCASS-X03** — Does not enforce the presence of scenarios or tests.


## SCIDS — ScenarioIdentity — two scenarios of the same feature cannot share one code

Confronts a feature against the question that decides whether its scenarios can be paired
at all: **each scenario has a code, but does each code point at ONE scenario?**

A rule legitimately has several scenarios — the happy path and its alternatives. What it
cannot have is two that are INDISTINGUISHABLE. With the same code, nothing links one
scenario to one test, and the relational gates end up comparing N titles against a single
test: at most one matches, and the others become a divergence nobody can resolve.

The way out is the scenario SUFFIX. Numbering keeps the rule legible in the prefix and
gives each case an identity of its own.

Measured on the project that originated the gate: 204 repeated codes, and 65 of them
carrying CONFLICTING kind tags on the same code — one scenario tagged as state, the other
as behaviour. That is the signature of a BORROWED code, not of a rule with two paths. One
of them proved to be a real defect: the spec defined the rule as one thing, and the second
scenario described a different behaviour that had no code of its own.

**PENDING and not FAIL**: numbering scenarios is a migration, and the gate is born over a
base that did not know the notation. Whoever already migrated stays green; whoever did not
sees what is left.


- **SCIDS-B01** — Two scenarios sharing one code are reported: nothing links one of them to one test, and the relational gates compare N titles against a single test.

- **SCIDS-B02** — The report NAMES the repeated code, so the reader does not have to scan the feature to find it.

- **SCIDS-B03** — The report says HOW MANY scenarios share the code, which separates an accidental duplicate from a code borrowed across a whole rule.

- **SCIDS-B04** — The report teaches the way out with the project's OWN repeated code as the example, not a generic one.

- **SCIDS-B05** — The SUFFIX gives each scenario its own identity: numbered, two scenarios of one rule pass.

- **SCIDS-B06** — Distinct codes pass — a rule with several scenarios is the common case and must not be accused.

- **SCIDS-B07** — The verdict is a DIVERGENCE and never a failure: numbering is a migration, and the gate is born over a base that did not know the notation.

- **SCIDS-B08** — An artifact that is not a feature leaves without a verdict: only a feature carries scenarios.

- **SCIDS-B09** — A feature with no coded scenario leaves without a verdict — there is nothing to confront, and that absence is another gate's charge.

- **SCIDS-B10** — Several repeated codes in one feature are reported TOGETHER, in a stable order, so two runs over the same file produce the same message.

- **SCIDS-B11** — Long titles are shortened in the report: the address is the code, and the title only helps recognise which scenario is which.

- **SCIDS-B12** — Two scenarios with different codes whose steps are the same — case and spacing aside — are reported, naming the one that repeats the other: they prove one thing twice, usually a body copied and never rewritten. Scenarios that differ only in their quoted values are variations, and a scenario of a single step only states an outcome many triggers share: neither is reported.

- **SCIDS-I01** — Grouping is by the COMPLETE code, suffix included. Grouping by the prefix alone would accuse exactly the projects that already did the migration the gate asks for.

- **SCIDS-I02** — The message is deterministic. The repeated codes are ordered before being joined, so the same feature always produces the same text and the finding does not churn between runs.

- **SCIDS-X01** — Does not judge whether the two scenarios describe DIFFERENT behaviours.

- **SCIDS-X02** — Does not look across features.

- **SCIDS-X03** — Does not charge the ABSENCE of a code on a scenario.

- **SCIDS-X04** — Does not renumber the scenarios, even knowing the fix.


## SCLTR — ScenarioLetterDeclared — the letter of a scenario code exists in the vocabulary

Confronts a scenario code against the project's own vocabulary: **does the letter it
carries mean anything?**

The project declares which letters it recognises, and a sibling gate charges the spec's
SECTIONS for using the right ones. Nobody charged the same of the code a SCENARIO carries
— and an invented letter passes every other gate. The code matches itself across feature
and test, the unit is complete, the relational gates find both ends of the edge, and
nothing notices that the letter means nothing.

Measured in the project that originated this gate: **18 codes carrying five undeclared
letters**, always accompanied by tags equally outside the vocabulary. The pattern is
recognisable once seen — someone needed a nature the project did not have and invented it
instead of declaring it.

Both repairs are legitimate, and the choice belongs to the project: declare the letter (if
the nature is genuinely missing) or remap the scenario onto a letter that already exists
(if it was only an alias for one). **The gate does not choose — it shows what is outside.**

That is why the verdict is UNDETERMINED and not a failure. Discovering that a nature was
used without registration is information; deciding between adopting it and remapping it is
work for whoever knows the domain.


- **SCLTR-B01** — An artifact that is not a feature leaves without a verdict: scenario codes live in features.

- **SCLTR-B02** — With no declared vocabulary the letters are the canonical ones (`RuleLetters`), as on the spec side (RLTYR-B05): a canonical letter passes and any other is reported.

- **SCLTR-B03** — A feature carrying no scenario code at all leaves without a verdict: there is nothing to judge.

- **SCLTR-B04** — Every letter found inside the declared vocabulary passes.

- **SCLTR-B05** — A letter outside the vocabulary is reported as a DIVERGENCE, never as a failure — the nature may deserve declaring, and that decision is not the gate's.

- **SCLTR-B06** — The verdict NAMES the letters that are outside and the codes that carry them, because that is what the reader needs in order to choose between declaring and remapping.

- **SCLTR-B07** — Codes sharing one unknown letter are grouped into a SINGLE line — eight scenarios with the same invented letter are one finding, not eight.

- **SCLTR-I01** — The scan is over the SHAPE of a code, never over the vocabulary. Reusing the shared scenario parser would make the gate blind to exactly what it exists to find: that parser builds its pattern FROM the declared letters, so a code with an invented letter is precisely what it discards. The vocabulary enters afterwards, when judging each letter found.

- **SCLTR-I02** — The code-length pattern is read at EVERY call, not frozen once. It comes from the project's Structure, which loads after the package globals — a pattern built once at load time would freeze the default and silently ignore the project's declaration.

- **SCLTR-I03** — A valid letter is never named in the verdict, even when the same feature carries invented ones. Reporting what is correct alongside what is not would make the reader hunt for the real finding.

- **SCLTR-X01** — Does not decide whether the invented letter should be declared or remapped.

- **SCLTR-X02** — Does not judge the TAGS accompanying the code, only the letter of the code itself.


## STASC — ScenarioTypeAligned — scenario classification tags must match the code nature letter

Confronts a scenario's classification tag against the nature letter embedded in its identity code:
**the classification tag of a scenario must agree with the rule letter of its code.**

The nature of a rule is declared twice: once in the letter embedded in its identity code (such as `S`
for State or `B` for Behavior) and once in the scenario's classification tag (such as `@state` or
`@behavior`). When the two disagree, one of them lies — and no other gate catches the contradiction
because each inspects only a single facet:
- `rule-types` validates that specification sections conform to the declared rule vocabulary;
- `feature-test-match` pairs scenario identity codes and descriptions against test cases.
Neither gate cross-references the scenario's Gherkin classification tag against the code's embedded letter.

Measured in doctrine: in the project that originated this gate, 46 scenarios exhibited tag and letter
disagreements upon audit. Distinguishing between the two failure modes is essential:
- 36 were wrong tags: for instance, "Intro present is rendered in italics" had the body `When the
  component is rendered / Then I should see the intro text` — pure state without user action, under a
  correct `-S` code, but tagged `@behavior`. The remediation was correcting the tag.
- 10 were borrowed codes: for instance, "Tapping a chip selects priority" was genuine interactive
  behavior, attached to a `-S` code because the specification had not catalogued the event. The
  remediation was creating a dedicated code in the events table.

This gate detects both cases. Deciding whether the tag was misapplied or the code was borrowed requires
reading scenario steps and understanding domain context. The gate guarantees that the divergence does
not remain silent.

The gate issues a **Pending** verdict rather than a blocking Fail: type mismatch represents inherited
technical debt in codebases adopting this gate after scenarios were already written, and resolving each
finding requires case-by-case editorial judgment.

Furthermore, if the project configuration defines no `tags:` mappings in its rule types, the gate skips
silently. Guessing tag meanings across different project languages would introduce false alarms; silence
here respects projects that did not request this check.


- **STASC-B01** — When the confronted node is not of kind feature, the gate skips confrontation.

- **STASC-B02** — When the project configuration is nil or defines no rule types, the gate skips confrontation.

- **STASC-B03** — When configured rule types define no scenario tag mappings, the gate skips confrontation.

- **STASC-B04** — When a feature file contains no scenarios declaring identity codes, the gate skips confrontation.

- **STASC-B05** — Scenarios with identity codes lacking a recognized rule letter are ignored during confrontation.

- **STASC-B06** — Scenario tags not registered in configuration rule types are ignored during confrontation.

- **STASC-B07** — When every recognized scenario tag aligns with its code's rule letter, the gate passes.

- **STASC-B08** — A scenario tag declared under multiple rule letters passes if any of its mapped letters matches the code's letter.

- **STASC-B09** — When a scenario is co-tagged with multiple codes, a tag matching any of those codes passes.

- **STASC-B10** — When a recognized scenario tag disagrees with the code's rule letter, the gate returns a divergence.

- **STASC-B11** — The Diverge message cites the mismatched code, rule letter, tag name, allowed letters, and shortened scenario title.

- **STASC-B12** — Multiple mismatch findings are sorted deterministically and reported together in the divergence verdict.

- **STASC-I01** — Scenario classification alignment is evaluated exclusively on feature files; non-feature nodes skip to keep confrontation scoped to where tags reside.

- **STASC-I02** — Without configured tag mappings, the gate stays silent to prevent false positives across different localization languages.

- **STASC-I03** — Type disagreements return Diverge rather than Fail because inherited divergence requires manual editorial judgment rather than mechanical fixes.

- **STASC-I04** — Co-tagged scenarios with multiple requirement codes accept tags matching any attached code to prevent false defects on valid multi-requirement scenarios.

- **STASC-X01** — Does not enforce presence of classification tags on scenarios.

- **STASC-X02** — Does not decide whether the tag or the code is the erroneous party when disagreement occurs.

- **STASC-X03** — Does not validate or alter specifications or test code files.

- **STASC-X04** — Does not restrict tags from being mapped to more than one rule type letter.


## SBGRD — SiblingGuard — sibling functions treat the same parameter consistently

Confronts a module against its own internal asymmetry: **when two sibling functions guard
a parameter and a third does not, the one that does not is almost always forgetfulness,
not decision.**

The motivating case was real. A versioning module exported three functions over the same
history. The TWO that wrote filtered the history by key before deciding; the one that READ
did not filter — and it was exactly the one receiving the multi-key array straight from
the repository. Asking for one key's version returned another key's, in silence. It passed
11 green gates and 16 tests.

**The asymmetry is the signal**, and it is what makes this detectable without understanding
the domain. The gate does not know what the guard does, only that the siblings apply it and
one does not. That is also why it cannot be a lint rule: nothing in the syntax is wrong.

Deliberately CONSERVATIVE. It accuses only when three conditions hold at once: three or
more exported functions receive the same parameter name, the MAJORITY applies a
recognisable guard over it, and at least one applies none. Below that bar it stays silent,
because a false positive here teaches the team to ignore the gate — and a gate that is
ignored defends nothing.


- **SBGRD-B01** — An artifact that is not code leaves without a verdict: the gate reads function bodies.

- **SBGRD-B02** — Without a declared dialect the verdict is UNDETERMINED, and says so — reading code it cannot recognise would stamp what it never checked.

- **SBGRD-B03** — Fewer than three siblings on the same parameter is left alone: below the bar, asymmetry is not evidence.

- **SBGRD-B04** — When the majority guards and at least one does not, the one that does not is ACCUSED.

- **SBGRD-B05** — When every sibling guards, nothing is accused — there is no asymmetry to report.

- **SBGRD-B06** — When no sibling guards, nothing is accused either: a module that never guards is a decision, not an oversight.

- **SBGRD-B07** — The verdict NAMES the function that fails to guard and the parameter at stake, so the reader does not diff the module by hand.

- **SBGRD-B08** — A waiver declared on the function WITH A WRITTEN REASON silences the accusation — the sibling that legitimately delegates the check says so where whoever reads the function will see it. A bare marker does not waive: a waiver with no why is the silence the gate exists to end.

- **SBGRD-I01** — The gate never judges WHAT the guard does. It reads that the siblings apply one and that one does not — understanding the guard would require understanding the domain, and that is what makes this detectable at all.

- **SBGRD-I02** — The three conservatism conditions hold TOGETHER. Any one of them alone would produce the false positive that costs the gate its credibility.

- **SBGRD-X01** — Does not invent what an exported function or a guard looks like.

- **SBGRD-X02** — Does not accuse a single function in isolation.

- **SBGRD-E01** — REF[SBGRD-B02]: an undeclared exported-function pattern is answered by B02: undetermined, and it says so


## SNGTU — SingleTestPerUnit — a unit has one test file per test layer

Two files testing one unit in one layer split what the unit is proven by, and each looks complete
on its own: in the reference app a hook had `useIapNative.test.ts` and `useIapNative.test.tsx`, and
five models had tests both under `__tests__/unit/models` and under `__tests__/unit/lambdas/`. Which
one to extend, which one a gate reads, which one a mutation run counts — each answer was a guess.

The unit a test tests is the project's own derivation read backwards (`TestedUnits`), so nothing
assumes a language or a layout. A split the project means is declared in one of the files with
`@split-test: <why>`.


- **SNGTU-B01** — A unit tested by more than one file in the same test layer fails, naming the layer and the files; files in different layers are not a split.

- **SNGTU-B02** — A `@split-test: <why>` in any of the files of a layer declares the split, and the unit passes for that layer.

- **SNGTU-B03** — A node that is not code, and a unit no test file tests, are skipped; without a map the gate is pending.


## SFMSP — SpecFeatureMatch — every requirement the spec DEFINES has at least one scenario

Confronts the edge of the unit that had no watcher: **the spec declares a requirement —
is there any scenario that exercises it?**

`feature-test-match` confronts feature→test; `unit-complete` confronts that the PIECES
exist. Nobody confronted spec→feature — and that is where a silent hole lives: the spec
declares a constraint rule, the feature has no scenario carrying its tag, and the requirement
crosses the whole pipeline with nothing verifying it. **Every gate stays green**: the spec
has a code, the feature exists, the feature matches the test. The requirement simply
belongs to nobody.

Measured in a real project: **11 of 287 specs with a feature had a requirement with no
scenario.**

The ruler is the same as the other relational gates — CODE, not prose: every
`{CODE}-{letter}{NN}` the spec DEFINES must appear as a scenario tag in the feature.
**Defining is different from citing:** a spec that mentions another unit's code (in a
Dependency Table, for instance) contracts no obligation. Without that distinction the
gate would be a noise generator and would be switched off.

Honest opt-out (CONCEPT §5.1): the per-requirement waiver marker (`no-scenario`, prefixed with `@`) on the requirement's line, followed by a colon and the reason, waives
that specific requirement, with the reason written. It serves what is genuinely not
observable by scenario — and leaves the trace that it was a decision, not forgetfulness.


- **SFMSP-B01** — A requirement the spec defines and no scenario tags FAILS, and the verdict names it.

- **SFMSP-B02** — A requirement that HAS a scenario is not accused: the verdict carries the uncovered ones only.

- **SFMSP-B03** — With every defined requirement tagged by some scenario, the gate passes.

- **SFMSP-B04** — A code merely CITED — in prose or in the Dependency Table — contracts no obligation, because a spec cites other units' codes all the time.

- **SFMSP-B05** — The per-requirement marker (`no-scenario`) with a written reason waives that one requirement.

- **SFMSP-B06** — A bare per-requirement marker, with nothing after the colon, does not waive — the waiver requires a reason.

- **SFMSP-B07** — The whole-spec marker (`no-feature`) with a reason DRAGS the waiver to every requirement: the spec has no feature, so no requirement of it can have a scenario.

- **SFMSP-B08** — Without that tag the same uncovered requirements keep failing — the drag cannot become a silent way of muting the gate.

- **SFMSP-B09** — A bare whole-spec marker drags nothing, or the marker would be a switch that turns the gate off without accounting for it.

- **SFMSP-B10** — A spec with NO feature returns Skip: that absence is the ruler of `unit-complete`.

- **SFMSP-B11** — A spec covered by SEVERAL features has its requirements looked for across all of them — the requirement only needs to be in some.

- **SFMSP-B12** — An artifact that is not a spec returns Skip: the gate has no jurisdiction over code, test or feature.

- **SFMSP-B13** — A rule written as an ALIAS of another rule of the same spec — `REF[CODE-B05]: <reason>` on its line — needs no scenario of its own: the target's scenario is the proof. It is how a rule is catalogued under one letter (a failure, `-E`) while its behaviour is already stated under another, without writing the decision twice.

- **SFMSP-B14** — An alias that stands for no rule FAILS, naming it: a target this spec does not define, a target that is itself an alias, or no reason after the colon. It is checked before the feature is looked for, because a dangling alias would drop the rule from every scenario check unseen.

- **SFMSP-B15** — A requirement is covered by scenarios whose tag carries it with a `#nn` suffix: the suffix identifies each case of the requirement, not a new requirement.

- **SFMSP-B16** — A code defined in the open-decisions section (`## Open Decisions`, `## Decisões em Aberto`…) is no requirement — for this gate and for every one that reads the spec's requirements (`scenario-coverage`, `reverse-match`, `rule-uses`): an open question gets its scenario once it becomes a rule.

- **SFMSP-B17** — With `data_states.required`, a data state the spec defines — `DS-<field>-<variant>`, bare or with the unit's code, opening a table row, a list item or a heading — is a requirement of the unit, for every gate that reads the spec's requirements (`spec-feature-match`, `scenario-coverage`), with the same no-scenario waiver as a rule; another unit's data state is a citation. Off, the data-state tables only document.

- **SFMSP-B18** — A rule whose line carries the retired marker (`@retired` with a colon, the revision and the reason) is a tombstone: it stays defined, so the old revisions that name it still resolve, and no gate asks anything of it — it is no requirement (`spec-feature-match`, `scenario-coverage`), no rule of the code (`rule-implemented`), uses nothing (`rule-uses-*`) and is no failure to handle. A rule's code is read whole: `CODE-R0001`, a revision, is not the rule `CODE-R00`. (`retiredLine`)

- **SFMSP-I01** — Every waiver requires a written REASON — both the per-requirement marker and the whole-spec one. A bare marker is a switch with no accounting, and silence without a why is what the gate exists to end.

- **SFMSP-I02** — Each gate accuses ONE thing. The missing feature belongs to `unit-complete` and the missing test to `feature-test-match`; accusing them here would print the same defect twice in the report.

- **SFMSP-X01** — Does not judge whether the scenario PROVES the requirement.

- **SFMSP-X02** — Does not charge a code the spec merely CITES.

- **SFMSP-X03** — Does not confront feature→test.

- **SFMSP-E01** — No map has been built, so the gate receives no graph.

- **SFMSP-E02** — A feature the map links to the spec is no longer on disk.


## VTRST — StateTransitions — every change of a visual unit is proven through what the screen shows

A validation changes the screen: the form refuses, a field turns red, a button disables, the flow
takes another path. Each of those is a state, and the states are already proven by their visual
capture. So a validation does not need a capture of its own; it needs to say which state it leads to.
An error, too, shows on the screen — as the message the spec catalogs; the state is the same, what
changes is the text, and the messages are captured.

Two gates tie each change of a visual unit to what the screen shows:

| Gate | Question |
| --- | --- |
| `validation-transitions` | Is every validation the trigger of a State Flow transition, from the state it is checked in to the state it leads to — or does it say `@no-state: <reason>`? |
| `error-message-declared` | Does every error name the message it shows — or say `@no-message: <reason>`? |

They run, like the visual-regression gates, on a visual unit's main code file and read the spec beside
it; sections and columns are read in any language of the catalog.


- **VTRST-B01** — A node that is not code, a code file with no spec beside it and a spec with no code leave both gates without a verdict.

- **VTRST-B02** — `validation-transitions` fails naming each validation — of `Validations` and of `Presentation validations` — that no State Flow row names as its trigger; a spec with no validation leaves without a verdict.

- **VTRST-B03** — A transition counts only when its From and To are states the spec registers; one from or to an unknown state is named with both ends.

- **VTRST-B04** — A validation whose row says `@no-state: <reason>` is not asked; one that says `@no-state` with no reason is still asked, and named as an exemption with no reason.

- **VTRST-B06** — A section is found at any heading level — `### Validações` nested under `## Rules` — and ends at the next heading of its level or above.

- **VTRST-B05** — `error-message-declared` fails naming each error that cites no message code of the unit, and each that cites codes the Messages section does not catalog; `@no-message: <reason>` on its row exempts it; a spec with no error leaves without a verdict.

- **VTRST-E01** — The spec beside the code file cannot be read, or there is none.


## THSAS — TestHasAssertion — every test asserts something in its body

A test that asserts nothing passes whatever the code does. It carries the scenario's code, its
title matches the scenario, the scenario counts as proved — and the body calls the code and checks
nothing. Every gate that reads titles is green over it.

What an assertion looks like is the test library's (`expect(`, `t.Errorf`, a project helper), so the
project declares it in `dialect.tests.assertion`, with a default for the families that have one.
Where a test's body ends is the language's too: a tests script that knows it says it (`end`), and
otherwise the body is read by its layout, which every language keeps whatever its syntax.


- **THSAS-B01** — A test whose body holds no assertion fails, each named by its line and title; a file whose tests all assert passes.

- **THSAS-B02** — A test's body is the block its line opens: a line ending in a bracket or `do` closes on the first later line back at its indentation that starts with a closer; any other line that opens a block ends it before the first later line back at its indentation; a line whose next line is not deeper is a block by itself. (`blockEnd`)

- **THSAS-B03** — A line that starts inside a multi-line literal — opened by a backtick or a triple quote — is text, and never ends a block; a backtick quoted alone is a character, not a delimiter.

- **THSAS-B04** — With `labels: true` on the gate, a test opened and closed on its own line with an empty block is a label for the block around it, and that block is what must assert; without it, the empty test is a test with no assertion.

- **THSAS-B05** — When the tests source says where a test ends (`end`), its body runs from its line to that one.

- **THSAS-B06** — A node that is not a test or is a support file, a project with no tests source or no assertion, and a file where the source lists no test are skipped.

- **THSAS-B07** — A test that carries `@no-assert: <why>` in its body or on the line above it is left out; a waiver with no reason waives nothing.

- **THSAS-B08** — A test listed at a line the content does not have is read as an empty body, never past the content's end. (`checkTestHasAssertion`, `testBody`)

- **THSAS-B09** — A line back at the opener's indentation that starts with a closer and ends opening a block — the table of an `it.each([` closing into the test's body, a `} else {` — goes on with the block; only a closer that opens nothing ends it. (`blockEnd`)

- **THSAS-E01** — The tests cannot be listed (a script that fails, output outside the contract), or the assertion does not compile


## TLVCD — TestLevelCodes — each scenario references only codes its test level accepts

A feature scenario is tagged with its test level (`@unit-level`, `@vr-level`, or the project's own
names) and references rules of the spec by code. The project declares, per level, which codes it
accepts, in the `levels` of this gate's own entry under `gates:`: a list of patterns to `allow`, a list to
`exclude`, or neither. This gate confronts every scenario with the filter of each level it is tagged
with, and names the codes the level does not accept.

It is the general form of a naming rule between level and code, and every rule of that kind is
configuration, not a gate of its own. In the reference app, 9 scenarios tagged with its visual level
had no `-VR` code, so the gate that asks for a baseline image never asked, and one `-VR` code sat
under the integration level. With the visual level allowing only `-VR$` and the other levels
excluding it, both show up here.


- **TLVCD-B01** — A node that is not a feature leaves with `Skip`.

- **TLVCD-B02** — A project that declares no filter per level leaves with `Skip`: every level accepts every code.

- **TLVCD-B03** — A feature with no scenario tagged with a level that declares a filter leaves with `Skip`.

- **TLVCD-B04** — A level that declares `allow` accepts only the codes that match one of its patterns; any other code of its scenarios fails.

- **TLVCD-B05** — A level that declares `exclude` refuses the codes that match one of its patterns, even when `allow` accepts them.

- **TLVCD-B06** — Every code of the scenario is confronted, not only the first, and without the `#NN` scenario suffix.

- **TLVCD-B07** — The failure names each refused code with the level that refused it, sorted.

- **TLVCD-B08** — When every confronted code is accepted, the gate passes.

- **TLVCD-B09** — The levels are read from every gate entry that runs this check; two entries declaring the same level add their lists together, and a `levels` on an entry of another check is not read.


## TSRCH — TestReach — a test reaches the unit it says it tests

A test can carry every code, match every title and assert — over a copy. In the reference app a
test pasted the unit's function into its own file and exercised the paste; another imported a
neighbour; an end-to-end test's `ref:` named one handler and invoked another. The gates that read
titles saw nothing wrong, and the unit could change freely.

Two gates ask it on one routine, apart so the project decides per gate what blocks:

- `test-exercises-unit` — the unit the project's derivation pairs with the test (`TestedUnits`);
- `test-ref-matches-unit` — the unit the test's `ref:` names.

Reaching is read without knowing the language. A test reaches a unit when an import line names
the unit's module, when it names something the unit defines (the dialect's `definition` says what
a definition looks like), or — for the `ref:` — when a call the project declares on the gate
(`invocations`) names it. A test that defines a name its unit defines is exercising a copy.


- **TSRCH-B01** — `test-exercises-unit`: a test that names, as a whole word, something its unit defines passes; one that neither imports the unit nor names anything it defines fails naming the unit. Every unit the derivation pairs with the test is confronted.

- **TSRCH-B02** — An import line of the test — by the dialect's `import_pattern` when declared, or the forms most languages share — that names the unit's file name without extension, or its directory as a path, reaches the unit; the name elsewhere in the test, or a unit at the root with no directory, does not count as an import. (`importsModule`)

- **TSRCH-B03** — A test that defines a name its unit defines fails naming the unit and the names, sorted — even when it reaches the unit otherwise: it is exercising its own copy. A name it copies does not count as reaching.

- **TSRCH-B04** — A defined name shorter than three characters is anyone's, and neither reaches nor copies; a definition's first non-empty capture group is its name. (`definedNames`)

- **TSRCH-B05** — `test-ref-matches-unit`: a test that reaches any of the code files the `ref:`'s spec governs passes; one that reaches none fails naming the `ref:` and the files, sorted.

- **TSRCH-B06** — A declared invocation whose first non-empty capture is the unit's file name without extension, or one of its directories, reaches the unit; another capture does not. (`invokes`)

- **TSRCH-B07** — A node that is not a test or is a support file is skipped by both gates; `test-exercises-unit` skips a test no unit pairs with and a project with no definition; `test-ref-matches-unit` skips a test with no `ref:` and a `ref:` whose spec governs no code.

- **TSRCH-B09** — A definition is identified with its owner when the dialect's `definition` captures one (a Go method by its receiver type, `GormPinger.Ping`), and a definition indented under a type with no owner captured is a member; a test that defines a member of its own type — a fake that implements the unit's interface — copies nothing, and only the same name of the same owner defined again is a copy. A member's bare name still reaches the unit. (`definedNames`)

- **TSRCH-B08** — A test that carries `@no-unit-import: <why>` passes `test-exercises-unit`, copies included; a waiver with no reason waives nothing.


## TSTRT — TestTraceable — a test linked to a feature must declare what scenario it proves

Confronts a test artifact against its linked feature: **a test connected to a feature must state
what scenario it proves.**

The blind spot this gate closes is subtle and occurred twice in a real production project: the test
EXISTS, PASSES, and covers the correct behavior — yet cites no scenario identity code whatsoever.
To every relational gate, such a test is invisible:
- `feature-test-match` checks scenario-to-test alignment by scenario code: without a code in the test,
  it reports the scenarios as unimplemented, mistakenly blaming the feature rather than identifying
  the untraceable test;
- `unit-complete` verifies file presence on disk and marks the test piece satisfied;
- `tests-green` executes the test suite and verifies assertions pass.

The outcome is the worst of both worlds: the engineering work was completed and verified, yet the
pipeline claims it was not. When another developer attempts to fix the reported deficit, they write
a redundant duplicate test for the exact same behavior, having no automated way to discover that the
existing test already exercised it.

Measured in doctrine: in a real project, a shared test exercised across 6 screens used inverted digit
and letter identifiers (such as transposed tokens where specs declared proper codes). Three screens
had tests and none was tracked because code search could not reach them. In the same project, a store
test covered 5 distinct behaviors without citing a single scenario code.

What separates this gate from neighbouring gates:
- `feature-test-match`: enforces granular 1:1 matching between each scenario and a dedicated test case.
  In contrast, `test-traceable` applies the weakest possible threshold: finding at least ONE valid
  scenario code anywhere in the test file satisfies the gate. The question here is simply: "does this
  test declare itself?"
- `unit-complete`: checks only that the four unit files exist on disk, regardless of whether their
  internal contents establish traceability.
- `tests-green`: compiles and executes test suites to ensure zero failures, but remains completely
  agnostic to requirement traceability codes.

Finally, tests without a linked feature (such as unit tests for internal utilities) are not charged:
demanding scenario codes from unlinked tests would require referencing nonexistent features.


- **TSTRT-B01** — When the confronted node is not of kind test, the gate skips confrontation.

- **TSTRT-B02** — When the dependency graph is nil, the gate returns Pending with a missing map verdict.

- **TSTRT-B03** — When a test node has no incoming `tested-by` edge from a feature in the graph, the gate skips confrontation.

- **TSTRT-B04** — When the linked feature file cannot be read from disk, the gate skips confrontation.

- **TSTRT-B05** — When the linked feature file declares no scenario codes, the gate skips confrontation.

- **TSTRT-B06** — When the test content contains at least one scenario code declared by its linked feature, the gate passes.

- **TSTRT-B07** — A single declared scenario code in the test content is sufficient to pass, even when the linked feature defines multiple scenarios.

- **TSTRT-B08** — When the test content contains none of the scenario codes declared by its linked feature, the gate fails.

- **TSTRT-B09** — A test covering behavior with transposed or misspelled codes fails exact substring confrontation.

- **TSTRT-B10** — When failing, the verdict names the linked feature path, lists up to three expected codes, and suggests the first code as a fix hint.

- **TSTRT-B11** — Scenario codes occurring anywhere in the test file content (such as in comments or test titles) satisfy the traceability check.

- **TSTRT-B12** — When the project says how its tests are written (`dialect.tests`, or its family's) and the source lists a test in this file, the test traces to its feature only through a test TITLE that cites one of the feature's codes; a code elsewhere in the file (a fixture, a helper) does not trace it. Without a source, or in a file the source lists no test in — one it does not describe, such as a YAML flow under a Jest pattern — the code counts anywhere in the file.

- **TSTRT-B13** — A support file is not charged with tracing to a scenario (Skip, saying why), even when a feature links to it.

- **TSTRT-I01** — Traceability is strictly a test requirement; non-test artifacts skip confrontation to prevent misattributing test debt to other files.

- **TSTRT-I02** — Tests without a linked feature in the graph are never charged, avoiding demands for references to nothing.

- **TSTRT-I03** — The acceptance threshold requires only one scenario code present in the file, preventing duplicate error reporting with scenario-matching gates.

- **TSTRT-I04** — Confrontation without a dependency graph returns Pending rather than approving an unmeasured relationship.

- **TSTRT-X01** — Does not enforce a 1:1 match between each scenario and a separate test case.

- **TSTRT-X02** — Does not execute test suites or inspect runtime assertion results.

- **TSTRT-X03** — Does not require scenario codes in standalone test files unattached to features.

- **TSTRT-X04** — Does not inspect the semantic validity or quality of test assertions.

- **TSTRT-E01** — REF[TSTRT-B02]: with no map the link to the feature cannot be followed, and B02 answers Pending

- **TSTRT-E02** — REF[TSTRT-B04]: a linked feature that cannot be read is answered by B04: the confrontation is skipped

- **TSTRT-E03** — The project's tests source fails, or answers outside its contract


## TICTS — TestIDContract — a test handle is one contract with four ends: the code exposes it, the spec declares it, a consumer queries it

A test handle (the attribute a project uses to mark an element for tests and automated flows) is one
contract with four ends: the code exposes it, the spec declares it in its test surface inventory, the
feature may describe it, and a test or a flow queries it. The testid-consistent gate confronts the ends
of that contract at once, starting from the spec, and reports per handle, with every end on the same
line. One gate and not several partial ones, because a handle renamed in the code used to produce two
findings in two gates ("exposed without declaring" for the new name, "declared without exposing" for the
old one) when the fact is one: somebody renamed and did not propagate.

Three files hold the contract, and they must agree on what a handle is:

- The gate itself decides when there is something to confront and what fails.
- The parsing half recognises the handle where it is written: in the code (a literal value, a template
  whose suffix is runtime data, a prop that carries a handle to a child, the branches of a conditional)
  and in the spec (the inventory table or list inside the test surface section).
- The consumer half finds who queries a handle: the test linked through the spec's feature, the tests
  beside the spec and in its sibling folders, and the end-to-end flows the project declares.

The attribute is the project's declaration, never a guess: without it the gate skips, because inferring
one would report green over what was never checked. The fourth edge, a consumer querying a handle no
code exposes, is deliberately left out: the flows are the whole project's tree, and charging one spec
for another screen's handle produced 829 findings in a single spec when it was tried. That question
belongs to a project-scope gate.


- **TICTS-B01** — The gate skips a node that is not a spec, and is Pending without a map.

- **TICTS-B02** — Without a declared handle attribute the gate skips.

- **TICTS-B03** — A spec with no readable linked code skips: there is no end that exposes a handle.

- **TICTS-B04** — A unit whose code exposes no handle and whose spec declares none skips: it has no test surface to contract.

- **TICTS-B05** — A handle exposed by the code, declared by the spec and queried by a consumer passes.

- **TICTS-B06** — A handle the code exposes and the spec does not declare fails, named in the report, also when the spec has no inventory at all.

- **TICTS-B07** — A handle the spec declares and the code does not expose fails, named in the report.

- **TICTS-B08** — When some consumer surface exists, a handle no consumer queries fails, named in the report.

- **TICTS-B09** — When there is no consumer surface at all, the queried end is not charged, and the report shows a dash in its place.

- **TICTS-B10** — The handle is read under the attribute the project declares, whatever its ecosystem.

- **TICTS-B11** — A literal handle value counts, also when it reaches the code through a derived prop or an object key whose name contains the attribute.

- **TICTS-B12** — A template handle, whose suffix is runtime data, and a prop that gives a child the head of its handles, both expose the static prefix as a wildcard.

- **TICTS-B13** — In a conditional handle only the literals of the branches are handles; the literals of the condition are not.

- **TICTS-B14** — The spec's inventory is read only inside its test surface section, whose title may carry a qualifier, and the section ends at the next heading.

- **TICTS-B15** — In a table row of the inventory only the first cell is the id; on any other line every quoted id counts.

- **TICTS-B16** — The attribute's own name, quoted in the inventory, is not a declared id.

- **TICTS-B17** — A wildcard at either end covers the concrete ids it opens; a concrete id does not cover a wildcard.

- **TICTS-B18** — A handle is queried when a consumer mentions it, with or without the mark; a wildcard handle, when a consumer mentions its prefix. A mention counts only at an id boundary (`TICTS-B22`).

- **TICTS-B19** — The end-to-end flows are consumers when the project declares the surface with a path: the surface's file template, or else the first override that gives one, whose static prefix is read as a whole; a surface with no path contributes nothing.

- **TICTS-B20** — The test files beside the spec and in the sibling folders of its folder are consumers, even when flows exist, without any edge to them.

- **TICTS-B21** — The static root read for a surface is the directory before the first placeholder or glob wildcard (`*`, `?`, `[`) of its path: `e2e/**/*.yaml` and `e2e/login-*.yaml` read `e2e`, `apps/x-{{module}}/flows` reads `apps`, and a path that opens with a wildcard reads the project root.

- **TICTS-B22** — A mention is a query only when no id character (letter, digit, `.`, `_`, `-`) touches it before, and, for a concrete handle, after: a consumer naming only `abcd-screen-header` or `my-abcd-screen` does not query `abcd-screen`.

- **TICTS-B23** — The test the spec's feature is tested by is a consumer wherever it lives: the link runs spec → feature → test, two hops, and it needs no neighbour folder to reach the test.

- **TICTS-B24** — Each line of the report shows whether the spec's feature describes the handle — a ✓ when the feature mentions it, a ✗ when it does not. It is information, never the reason for the failure: not every handle belongs in a written scenario.

- **TICTS-I01** — One handle is one line of the report, whatever its spelling at each end: the mark is not part of the identity, and the line spells it as the code does.

- **TICTS-X01** — A handle only a consumer mentions, which the spec does not declare and the code does not expose, is not charged to the spec.

- **TICTS-E01** — A linked code file cannot be read.


## TQETS — TestidQueriedExists — every handle queried by an E2E flow must exist in code

Confronts the end-to-end execution surface against the application code: **does an automated flow
query a test handle that no code file exposes?**

This gate closes the fourth edge of the testID contract. The sibling gate `testid-consistent` starts
from a single spec and confronts inventory coherence; however, the E2E surface encompasses the entire
tree of flows across the whole project. Charging an individual screen's spec for an external flow handle
would accuse innocent specs of foreign handles (measured when attempted: 829 findings in a single spec,
all belonging to other screens). Therefore, the question belongs to the PROJECT scope, confronting both
sides collectively.

The measured defect in the reference application (2026-08-25) revealed **13 invented IDs in flows**,
frequently caused by incorrect screen prefixes (`review-*` where the component actually emits `revi-*`).
The assumption that "an invalid ID will simply fail at runtime in the test runner" is FALSE for two
independent reasons:
1. **Flow execution reachability**: The runner must reach the step. In a test suite where step 3 fails,
   phantom IDs in subsequent steps remain completely hidden — surfacing one per run across seven iterations.
2. **Vacuum passing in negative assertions**: In `assertNotVisible`, an invented or non-existent handle
   causes the test to PASS immediately. In `REVI-R03`, asserting `:review-edit-controls` (which never existed)
   proved the exact opposite of reality: green by VACUITY. This is the most dangerous defect, and the test
   runner can never catch it.

The gate enforces a STRICT SINGLE DIRECTION (flow → code). The reverse (code exposing a handle not queried
by any flow) is not an issue: not every marked UI element requires an automated scenario, and `testid-consistent`
governs spec-level declaration.


- **TQETS-B01** — When the project does not declare a test handle attribute, the gate skips without asserting green.

- **TQETS-B02** — When the project does not declare an E2E surface, the gate skips honestly rather than reporting compliance.

- **TQETS-B03** — When no exposed handles are found in project code, the gate skips confrontation.

- **TQETS-B04** — A handle queried by a flow that matches an exposed handle in code passes.

- **TQETS-B05** — A handle queried by a flow that does not exist in any code file fails.

- **TQETS-B06** — The failing verdict names the missing handle and lists all flow files that query it.

- **TQETS-B07** — A non-existent handle inside a negative assertion (`assertNotVisible`) fails, preventing false green by vacuity.

- **TQETS-B08** — A code template handle (`:item-*`) satisfies a concrete instance queried by a flow (`:item-3`).

- **TQETS-B09** — A regex pattern queried by a flow matches against the exposed prefix head in code.

- **TQETS-B10** — A dynamic runtime expression interpolated by the runner (`${...}`) is skipped without failing.

- **TQETS-B11** — A marked handle literal defined in a lookup table or constant object counts as exposed.

- **TQETS-B12** — A template handle defined behind a nullish coalescing fallback (`??`) counts as exposed.

- **TQETS-B13** — A suffix composed in a child component from a parent prop (`*-row-*`) satisfies queries against composite IDs.

- **TQETS-B14** — A terminal suffix composed from a prop (`-toggle`) satisfies composite IDs ending with that segment.

- **TQETS-B15** — Handles referenced only within test files (`.test.tsx`, `.spec.ts`) do not count as exposed in application code.

- **TQETS-I01** — Lack of configuration never produces a false approval. Without an explicit test handle or E2E surface, the gate returns Skip.

- **TQETS-I02** — Findings are grouped by handle ID across flows so that a single missing handle queried multiple times does not multiply reported defects.

- **TQETS-I03** — Negative assertions are confronted with the same rigor as positive queries to eliminate green-by-vacuity passes.

- **TQETS-X01** — Does not charge code for exposing handles that no flow queries.

- **TQETS-X02** — Does not execute flows or evaluate runtime JavaScript expressions.

- **TQETS-X03** — Does not assume a default test handle attribute (such as `testID`).

- **TQETS-E01** — The project declares an E2E surface, but the directory its pattern points to does not exist on disk.

- **TQETS-E02** — A flow of the E2E surface, or a source file of the project, cannot be read.


## TRDCT — TriggerDeclared — cited compliance triggers and obligations must exist in the declared vocabulary

Confronts compliance triggers and obligations cited within specification text against the active
vocabulary declared in project configuration and adopted compliance packs: **every cited obligation
trigger or obligation name must exist in the declared vocabulary.**

Compliance obligations are triggered by declarations in artifact headers (such as `carries: personal-data`).
When a specification instructs an author on what trigger to declare, it teaches the vocabulary. If it
teaches an identifier that no pack declares, the author follows the instruction, writes the header, and
**no obligation ever triggers**.

This produces no structural failure in ordinary gates. The header is syntactically well-formed, header
checks pass, the unit appears covered by regulatory frameworks (such as LGPD or GDPR), yet actual
enforcement is zero. It is the most insidious form of a lying anchor: the instructional text is wrong,
causing every developer and automated agent relying on it to reproduce the defect.

Measured in doctrine: in a real project, 47 model specifications instructed authors to declare `carries: pii`
and cited an obligation named `pii-purgavel`. Neither identifier existed in any pack — the canonical
vocabulary declared across packs was `carries: personal-data`, with the obligations `lgpd-eliminacao` and
`lgpd-portabilidade`. All 47 instances originated from the same boilerplate snippet, copied from spec to spec,
without any check raising an alarm.

What separates this gate from neighbouring gates:
- Header gates (`header-conforms`, `header-valid`): only verify that header syntax is well-formed and fields
  are structured; they do not validate whether the trigger values activate real obligations.
- Obligation gates (`obligation-honored`): verify that active obligations attached to a unit are fulfilled;
  but if a misnamed trigger never fired the obligation in the first place, `obligation-honored` has nothing
  to enforce.
- `trigger-declared`: confronts cited trigger symbols and obligation names mentioned in spec prose against
  the declared universe of packs and configuration.

Finally, citing triggers is entirely optional. Specifications that cite no compliance triggers skip cleanly.
When a project has declared no compliance obligations and loaded no packs, the gate returns **Pending**
rather than approving unverified vocabulary or penalizing projects without compliance.


- **TRDCT-B01** — When the confronted node is not of kind spec, the gate skips confrontation.

- **TRDCT-B02** — When the specification text contains no cited obligation triggers, the gate skips confrontation.

- **TRDCT-B03** — When no compliance vocabulary is declared in configuration or packs, confronting cited triggers returns Pending.

- **TRDCT-B04** — Non-trigger key-value citations (such as layer or code markers) are ignored.

- **TRDCT-B05** — Natural language prose mentioning compliance terms without backtick quotes is ignored.

- **TRDCT-B06** — When every cited trigger value exists in the declared vocabulary, the gate passes.

- **TRDCT-B07** — When cited obligation names exist in the declared vocabulary, the gate passes.

- **TRDCT-B08** — A cited trigger value absent from the declared vocabulary fails, providing a nearest-match suggestion.

- **TRDCT-B09** — When no close match exists for an undeclared trigger, the failure suggests available declared triggers.

- **TRDCT-B10** — A cited obligation name absent from the declared vocabulary fails.

- **TRDCT-B11** — Duplicate citations of the same trigger or obligation within a file are deduplicated in defect reporting.

- **TRDCT-B12** — Multiple vocabulary errors are sorted deterministically and aggregated in the failure verdict.

- **TRDCT-B13** — When no declared trigger is close enough to suggest, the verdict lists the declared ones — capped at four. A list of every trigger a large project declares would bury the advice it exists to give.

- **TRDCT-I01** — Trigger and obligation citations are validated only in specifications; other artifacts skip confrontation.

- **TRDCT-I02** — Missing compliance packs and obligations produce Pending rather than Pass, preventing silent approval of unverified claims.

- **TRDCT-I03** — Trigger keys are restricted to a closed set of recognized compliance predicates to prevent false positives on general key-value metadata.

- **TRDCT-I04** — Undeclared triggers produce a definitive Fail verdict because teaching invalid triggers breaks downstream compliance enforcement.

- **TRDCT-X01** — Does not mandate that specifications cite compliance triggers or obligations.

- **TRDCT-X02** — Does not validate compliance triggers in source code, features, or test files.

- **TRDCT-X03** — Does not enforce implementation of the cited obligations within code.

- **TRDCT-X04** — Does not inspect unquoted natural language mentions of compliance concepts.

- **TRDCT-E01** — REF[TRDCT-B03]: with no vocabulary declared in the configuration or the packs, B03 answers Pending

- **TRDCT-E02** — A pack the project declares cannot be loaded.


## UNTCP — UnitComplete — the pieces that realize a spec EXIST

Confronts a spec of a GOVERNED layer against the simplest question of the unit: **do the pieces
that realize it exist?** The code it specifies, the feature that covers it, and the test
that proves it.

It exists because the relational gates FAIL OPEN by construction. Without a test linked, the
feature↔test confrontation returns "nothing to confront yet" instead of failing; without code,
the dependency one likewise. The side effect is grave: a lone spec, with no
implementation at all, crosses ALL the gates and the pipeline concludes "can promote" — the green
certifying work that does not exist.

This gate closes the hole from the positive side. Instead of asking "do the pieces match?" — which
requires that they exist —, it asks "do the pieces exist?".


- **UNTCP-B01** — An artifact that is not a spec leaves without a verdict: only the spec has a unit to demand.

- **UNTCP-B02** — A RECOGNIZED layer (declarative regime) leaves without a verdict: it has neither spec nor unit by definition.

- **UNTCP-B03** — Without a map the verdict is UNDETERMINED. Approving without being able to look would be asserting what was not measured.

- **UNTCP-B04** — A spec with the three pieces linked passes.

- **UNTCP-B05** — A spec missing some piece fails, and the verdict NAMES which ones are missing and where each one is born.

- **UNTCP-B06** — The layer may waive a piece as a block, declared in the Structure.

- **UNTCP-B07** — The unit may waive a piece in the spec itself, with a written reason — it is the granularity the per-layer waiver does not reach. @realizes SAIDA-R01

- **UNTCP-B08** — A piece declared TO BE DEVELOPED (`@TBD`) leaves the verdict a declared DIVERGENCE, never approved: it is DEBT, and it stays visible until someone writes it. Measured before the fix: `@TBD` shared a bucket with the permanent waiver, and a spec declaring three unwritten pieces came out green — the honest declaration erased the pending work from the radar. @realizes SAIDA-R02 @realizes SAIDA-R03

- **UNTCP-B09** — A waiver written on a line that DEFINES a rule — its table row, its heading (`### CODE-S01: …`) or its bullet — is that rule's, not the unit's: only a marker on a line of its own waives a piece of the unit. Measured before the fix: the reference app had 308 unit failures, 293 of them per-rule `@no-code:` in table rows read as a unit waiver of code, feature and test.

- **UNTCP-I01** — The test is reached in TWO hops — spec → feature → test —, because the one that points at the test is the feature. Checking the test directly on the spec would accuse the whole project of missing tests.

- **UNTCP-I02** — Waiving the test and writing a scenario in the feature is a CONTRADICTION, and fails. The two assertions do not coexist: either the scenario is real and someone must prove it, or it should not exist.

- **UNTCP-I03** — Waiving the test requires saying WHERE the proof is, and the place has to exist. An orphaned reference fails.

- **UNTCP-I04** — The waiver holds only for the declared piece. Waiving one never waives the others.

- **UNTCP-X01** — Does not confront whether the pieces MATCH one another — only whether they exist.

- **UNTCP-X02** — Does not judge the QUALITY of any piece.

- **UNTCP-E01** — A test the map lists is no longer on disk while the gate looks for the test that the `@no-test` reference points at.

- **UNTCP-E02** — A feature the spec is `covered-by` is no longer on disk while the gate counts the scenarios that would contradict `@no-test`.


## VLANV — ValueAnchored — a replicated key is declared where it is used, and every copy carries the same value

A value that lives in more than one place has no address. Changing a colour means changing it
everywhere, and nothing says where "everywhere" is: the copies are plain literals, and the one
that was forgotten keeps compiling.

The gate turns the replicated value into a symbol. Whoever replicates a value declares its KEY
in a comment, and the code line right below carries the value:

    // @code-reference-[COLOR-SUCCESS]-[#1F8A5B]
    success: '#1F8A5B',

Three confrontations follow from the declaration:

- **LOCAL** — the next code line below it (comment and blank lines skipped) contains the
  declared value: the declaration does not lie about its own line.
- **ACROSS** — every declaration of the same key declares the same value. Without this, a
  change reaches one file — value and declaration together — and every other copy stays
  behind while each file is locally consistent.
- **SPEC** — when the key is a rule code whose defining line in the spec declares a value,
  every declaration is confronted with it. The rule is then the source: change it there, and
  every place still carrying the old value is reported.

The key may therefore have a source in the spec, or be only a replicated reference whose
copies are the truth.


- **VLANV-B01** — A declaration whose next code line contains the declared value passes.

- **VLANV-B02** — A declaration whose next code line does not contain the declared value fails, and the verdict shows the key, the declared value and the line.

- **VLANV-B03** — Comment and blank lines between a declaration and the code are skipped, so declarations may be stacked above one line.

- **VLANV-B04** — Declarations of the same key with different values fail, in every file that holds one, listing every place and value.

- **VLANV-B05** — A declaration whose key is a rule declaring a value in the spec, and that disagrees with it, fails, naming the spec and the value it declares.

- **VLANV-B06** — Without a declared pattern the gate skips and names the setting that enables it.

- **VLANV-B07** — A declaration with no code line below it fails: it annotates nothing.

- **VLANV-B08** — Only an anchor that is the whole content of a comment line is a declaration; an anchor inside prose (a comment explaining the syntax) or inside code is a mention, and is not charged.

- **VLANV-I01** — A pattern with fewer than two capture groups is treated as not declared.

- **VLANV-I02** — Every declaration of the project is indexed once per map, not once per file — reading the whole repository per code file is the shape of cost that once made `docs-fresh` 97% of a check.

- **VLANV-I03** — With no built map the verdict is never approval.

- **VLANV-X01** — A rule whose defining line declares no value is not charged against the spec.

- **VLANV-X02** — A literal nobody declared is not charged.

- **VLANV-X03** — Does not decide the anchor's syntax.

- **VLANV-E01** — A spec the map lists is no longer on disk when the project index is built.

- **VLANV-E02** — A code file the map lists is no longer on disk when the project index is built.

- **VLANV-E03** — `derived.value_anchor` is declared with fewer than two capture groups.


## VRBSV — VRBaseline — ensures visual regression scenarios have captured reference baseline images

Verifies that visual regression scenarios declared in feature files have corresponding reference baseline capture images on disk. Visual regression represents a distinct proof surface where verification occurs through visual screen capture rather than unit test assertions. Without a baseline image, a visual scenario exists in the feature and is counted as covered, yet no comparison image exists against which changes can be evaluated, leaving promised visual proofs unverified. In the originating repository audit, out of 105 baselines on disk and 4 features declaring visual regression scenarios, 3 scenarios were declared without any baseline image.


- **VRBSV-B01** — Artifacts that are not feature files leave with verdict `Skip`.

- **VRBSV-B02** — Feature files declaring no visual regression scenarios leave with verdict `Skip`.

- **VRBSV-B03** — Visual regression scenarios having matching baseline images on disk pass with verdict `Pass`.

- **VRBSV-B04** — Visual regression scenarios lacking matching baseline images fail with verdict `Fail` naming the missing scenario codes.

- **VRBSV-B05** — Baseline image matching supports naming variants where variant suffixes are appended after the scenario identifier.

- **VRBSV-B06** — Reads visual regime tag from project configuration `derived.regimes` mapping tag keys to regime names.

- **VRBSV-B07** — Falls back to default tag `vr-level` when project configuration is nil or defines no visual regime.

- **VRBSV-B08** — Multiple missing baseline scenario codes are sorted alphabetically in the failure diagnostic message.

- **VRBSV-B09** — Only scenario codes containing the visual regression suffix `-VR` are matched for baseline existence.

- **VRBSV-I01** — Every visual regression scenario declared in a feature must correspond to at least one baseline image on disk.

- **VRBSV-I02** — The visual regime tag mapping treats the map key as the tag in the feature and the map value as the regime name.

- **VRBSV-X01** — Does not evaluate baseline staleness using disk modification timestamps.

- **VRBSV-X02** — Does not fail commits based on git commit dates of baseline images.


## VRSTC — VRStatesCovered — each state of a visual unit tied to its visual regression, both ways

A state is what a screen or a component looks like under a condition, and a visual capture is what
keeps it looking so. An agent setting up a project with screens wrote specs with states, features and
unit tests, and no visual regression at all — every gate green, and no state of any screen protected
against a visual change. `vr-baseline` only asked a VR scenario that exists for its image.

Four gates ask the four questions that tie a state to its capture, both ways:

| Gate | Question |
| --- | --- |
| `vr-states-covered` | Does every state of the spec have a VR scenario in the feature? |
| `vr-scenarios-tested` | Does every VR scenario have a VR test, and a baseline image? |
| `vr-scenarios-of-states` | Is every VR scenario of a state the spec registers? |
| `vr-tests-of-scenarios` | Is every VR test of a VR scenario the feature declares? |

Every state of a screen or a component is asked for a capture — only visual units are confronted at
all. The exception is written where the state is: `@no-vr: <reason>` on its heading or its row, for a
state with no visual value of its own, such as a transient loading or a state that looks like another.
An exemption with no reason does not exempt.

They run on a visual unit's main code file — the screen or the component, which is what a project tags
as visual — and read the unit's spec, feature and tests from it (`Button.tsx` → `Button.spec.md`,
`Button.feature`). A VR scenario is a scenario tagged with the project's visual regime and the code of
the state it captures (`@{CODE}-S01` or `@{CODE}-VR-S01`). A VR test names `{CODE}-VR-S01` in its path or
its text. A baseline is `<Unit>.{CODE}-VR-S01[-variant].<ext>` beside the unit.


- **VRSTC-B01** — A node that is not code, a code file with no spec beside it (a part of the unit) and a spec with no code leave every gate without a verdict.

- **VRSTC-B02** — `vr-states-covered` fails naming each state of the spec with no VR scenario and no exemption, and the regime tag to use; with no state it leaves without a verdict.

- **VRSTC-B03** — `vr-scenarios-tested` fails naming each VR scenario that no VR test names; with no VR scenario it leaves without a verdict.

- **VRSTC-B04** — `vr-scenarios-tested` also fails naming each VR scenario with no image beside the unit, `<Unit>.{CODE}-VR-<state>` with an optional variant, as png, jpg, jpeg, webp, gif or svg, and the name to save it under.

- **VRSTC-B05** — `vr-scenarios-of-states` fails naming each VR scenario of a state the spec does not register, and each VR scenario of a state the spec exempts with `@no-vr`.

- **VRSTC-B06** — `vr-scenarios-of-states` fails naming each VR scenario that carries no state's code — a scenario capturing the whole unit at once.

- **VRSTC-B07** — `vr-tests-of-scenarios` fails naming each VR state a test names with no VR scenario in the feature, and the tests that name it; with no VR test it leaves without a verdict.

- **VRSTC-B08** — A VR scenario is a feature line carrying the visual-regime tag; its state is a `{CODE}-<state>` or `{CODE}-VR-<state>` code on that line.

- **VRSTC-B09** — A test is of the unit when its path names the unit's code, when it sits beside the unit under the unit's name, or when a folder of its path is named after the unit; an image is never a test.

- **VRSTC-B10** — The State letter is the one of the project's rule type whose term or a section starts with "state" or "estado"; without one, `S`.

- **VRSTC-B11** — `vr-baseline` accepts the same image formats for a VR scenario's baseline.

- **VRSTC-B13** — Every message the spec catalogs (the codes of its User Messages section) is captured like a state: a VR scenario `{CODE}-VR-M01`, a VR test and a baseline image, `@no-vr: <reason>` on its row exempting it — an error shows on the screen as its message, the state is the same.

- **VRSTC-B14** — The states a spec registers are the codes in its States section — whose title may carry a note in parentheses — when it has one; a state code cited elsewhere registers nothing.

- **VRSTC-B15** — A `@no-vr` exempts the code its line declares — the first cell of a table row, or the code before the marker — and never a code its reason mentions: "same frame as M03" exempts the row's own state, not M03.

- **VRSTC-B12** — A state is exempted from visual regression by `@no-vr: <reason>` on a line that declares it — its heading or its row in a states table. An exemption with no reason does not exempt: the state is still asked, and the failure names it as an exemption with no reason.

- **VRSTC-I01** — One state's scenario, test or image never answers for another state.

- **VRSTC-E01** — The spec beside the code file cannot be read, or there is none.

- **VRSTC-E02** — The unit's feature cannot be read, or there is none.

- **VRSTC-E03** — Looking for a baseline image fails on the file system.

- **VRSTC-E04** — A test of the unit cannot be read.



