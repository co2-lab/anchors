<!-- anchors:generated from doct/camadas/gate.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Camada: gate




## CDCTC — CodeCataloged — what the code EXPORTS must be in the spec, or waived in the code









#### CDCTC-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because the ruler starts from the spec that governs the code
```

#### CDCTC-B02 — An exported symbol the spec never names fails, and the verdict names it

```gherkin
    Given a spec cataloguing one rule and a code file exporting three functions
    When the gate confronts it
    Then it returns Fail naming each orphan and its line, because the name alone would
      make the reader hunt for the symbol in the file
```

#### CDCTC-B03 — What the spec already catalogues is never accused

```gherkin
    Given a spec cataloguing one of the exported functions by name
    When the gate confronts it
    Then the verdict does not name that function, because it is already catalogued
```

#### CDCTC-B04 — A no-rule marker with a written reason waives the symbol

```gherkin
    Given an exported function carrying a no-rule marker with a written reason
    And the spec that catalogues every other symbol
    When the gate confronts it
    Then it returns Pass, because not every export deserves a rule and the waiver is what
      makes the noise manageable without lying
```

#### CDCTC-B05 — A bare no-rule marker does not waive

```gherkin
    Given an exported function carrying a no-rule marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a bare marker would be a silent way to quiet the gate
      and the trace that a decision was taken would vanish
```

#### CDCTC-B06 — A spec cataloguing every exported symbol passes

```gherkin
    Given a code file exporting two constants and a spec naming both
    When the gate confronts it
    Then it returns Pass
```

#### CDCTC-B07 — With no code linked the gate leaves without a verdict

```gherkin
    Given a spec with no code file linked to it in the map
    When the gate confronts it
    Then it returns Skip, because the absence belongs to the triad gate and accusing it
      in both places would duplicate the debt
```

#### CDCTC-B08 — Without a declared export pattern the gate skips and says so

```gherkin
    Given a Go file and a project that declared no export pattern
    When the gate confronts the spec that governs it
    Then it does not return Pass, and the verdict names how to enable the pattern,
      because green over what was never read is worse than an honest red
```

#### CDCTC-B09 — With the pattern declared the gate confronts for real in any language

```gherkin
    Given a project declaring an export pattern for Go
    And a Go file exporting a function the spec never names
    When the gate confronts the spec
    Then it returns Fail naming that function
```

#### CDCTC-B10 — The declared dialect family also supplies the pattern

```gherkin
    Given a project that names its dialect family as Go and declares no pattern of its own
    And a Go file exporting a function the spec never names
    When the gate confronts the spec
    Then it returns Fail naming that function, because naming the family is enough
```

#### CDCTC-I01 — The waiver holds in the comment block above the symbol

```gherkin
    Given the no-rule declaration written inline, one line above, in a two-line comment
      block, in a block with paragraphs and in a doc comment
    When the gate reads the context of the symbol in each case
    Then the declaration holds in all of them, because it is documentation and the
      explanation rarely fits on one line
```

#### CDCTC-I02 — The waiver does not leak between symbols

```gherkin
    Given two exported functions where only the first carries a no-rule declaration
    When the gate reads the context of each symbol
    Then only the first carries the declaration, because inheriting would let one marker
      exempt the whole file, which is the opposite of what it is
```

#### CDCTC-I03 — The gate never approves a language it cannot read

```gherkin
    Given a Go file whose exports match no TypeScript syntax
    And a project that declared no export pattern
    When the gate confronts the spec
    Then it does not return Pass, because stamping approval over what was never read is
      the worst possible failure in a measuring instrument
```

#### CDCTC-X01 — The gate does not judge whether the rule describes the symbol well

```gherkin
    Given a spec whose rule names the exported function and describes it wrongly
    When the gate confronts it
    Then it returns Pass, because the ruler is whether the spec NAMES the symbol —
      judging what the rule says about it belongs to another gate
```

#### CDCTC-X02 — The gate does not decide which symbols deserve a rule

```gherkin
    Given an exported function of pure formatting, which a reviewer would exempt
    And no no-rule declaration anywhere near it
    When the gate confronts the spec
    Then it returns Fail, because the exemption is the project's call and the waiver is
      where it records it — deciding here would remove the calibration
```

#### CDCTC-X03 — The gate knows no language, the project declares what is public

```gherkin
    Given two projects whose export patterns recognise different syntaxes
    And the same file, public under one pattern and invisible under the other
    When the gate confronts each
    Then the verdicts differ, because recognising what is public depends on the language
      and Anchors does not presume
```

#### CDCTC-X04 — The gate does not charge the absence of code

```gherkin
    Given a spec cataloguing rules with no code file linked to it
    When the gate confronts it
    Then it returns Skip rather than Fail, because accusing the same debt in two gates
      would duplicate the finding
```


## CDLNG — CodeLanguage — the code does not go back to mixing languages









#### CDLNG-B01 — An identifier in the wrong language is accused, and an English one passes

```gherkin
    Given a file declaring a function named in Portuguese
    And another declaring a function named in English
    When the gate confronts the file
    Then the Portuguese identifier is accused
    And the English one passes without noise
```

#### CDLNG-B02 — The verdict returns the word that accused

```gherkin
    Given a file declaring an identifier in the wrong language
    When the gate confronts it
    Then the verdict names that exact word, so the reader does not hunt the whole file
```

#### CDLNG-B03 — Only a DECLARATION is the subject

```gherkin
    Given a file whose Portuguese words appear outside any declaration
    When the gate confronts it
    Then it accuses nothing, because what declares no identifier is not read
```

#### CDLNG-B04 — The declarations are found in every form the language offers

```gherkin
    Given a file declaring identifiers as function, type, variable and constant
    When the gate confronts it
    Then every one of those declarations is read, not only the most common form
```

#### CDLNG-B05 — Deciding one word is separate from deciding a whole identifier

```gherkin
    Given a compound identifier made of several words
    When the gate decides its language
    Then it breaks the identifier into its words and decides each one
    And that separation is what lets the length floor apply per word instead of to
      the identifier as a whole
```

#### CDLNG-I01 — A short word does not count

```gherkin
    Given a file declaring identifiers below the length floor
    When the gate confronts it
    Then none is accused, because below the floor there is no language to infer and
      accusing there is the noise that costs a gate its credibility
```

#### CDLNG-X01 — The gate does not read comments

```gherkin
    Given a file whose comments are written in the team's own language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because the comment carries the measurement and the why
```

#### CDLNG-X02 — The gate does not read user-facing text

```gherkin
    Given a file whose user-facing strings are written in the project's language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because that text goes through the translation catalog
```

#### CDLNG-X03 — The gate does not use a dictionary to decide the language

```gherkin
    Given a file declaring compound identifiers in English that no common dictionary holds
    When the gate confronts it
    Then none is accused, because the dictionary approach was measured at ninety percent
      false positives — and a gate that wrong is switched off, defending nothing
```


## CSDCN — ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those









#### CSDCN-B01 — A status emitted and not declared is accused by number

```gherkin
    Given an output contract declaring 200, 404 and the 5xx range
    And a handler that also emits 401, 403 and 409 on its refusal branches
    When the gate confronts it
    Then it returns Fail naming 401, 403 and 409, because the client programmed from the
      table does not handle a refusal it was never told about
```

#### CSDCN-B02 — A status declared and emitted by no path is accused as a phantom

```gherkin
    Given an output contract declaring 402 for exceeded quota
    And a handler whose quota branch answers 429 and never 402
    When the gate confronts it
    Then it returns Fail naming 402 as dead code in the client, which disappears with
      nobody noticing
```

#### CSDCN-B03 — A faithful table passes

```gherkin
    Given an output contract whose concrete status codes are exactly the ones the handler emits
    When the gate confronts it
    Then it returns Pass, because a gate that only accuses is a noise generator
```

#### CSDCN-B04 — The 500 of the top-level try/catch is not charged

```gherkin
    Given an output contract declaring 200 and the 5xx range
    And a handler whose only other status is the 500 of its catch block
    When the gate confronts it
    Then it returns Pass, because that 500 is infrastructure every handler carries, not a
      decision of this one
```

#### CSDCN-B05 — The 5xx range covers, the 4xx range does not

```gherkin
    Given an output contract declaring 200, the 4xx range and the 5xx range
    And a handler that emits 503 and 403
    When the gate confronts it
    Then it returns Fail naming 403 and staying silent about 503, because a generic 4xx
      would hide exactly the access refusals this gate hunts
```

#### CSDCN-B06 — A status that lives only in a comment is not emitted

```gherkin
    Given an output contract declaring 200 and 404
    And a handler whose 404 appears only in a comment describing the old behaviour, the
      live branch answering 403
    When the gate confronts it
    Then it returns Fail naming 404 as a phantom, because a comment is not behaviour
```

#### CSDCN-B07 — Without the contract section there is nothing to confront

```gherkin
    Given a spec that catalogues effects and opens no output contract section
    When the gate confronts it
    Then it returns Skip, because charging the section's existence belongs to spec-complete
```

#### CSDCN-B08 — Code that returns no status is skipped

```gherkin
    Given an output contract that declares void, the contract of a cron handler
    And code that persists items and returns no status
    When the gate confronts it
    Then it returns Skip, because there are no numbers on either side to compare
```

#### CSDCN-B09 — A literal status passed to a local helper counts as emitted

```gherkin
    Given an output contract declaring 200 and 400
    And a handler that builds its 400 through a locally defined fail helper called with the literal
    When the gate confronts it
    Then it returns Pass, because a 400 is a 400 wherever the envelope is built
```

#### CSDCN-B10 — With a dynamic status the phantom side goes quiet and the literals still count

```gherkin
    Given an output contract declaring 200 and 404
    And a handler with a helper that takes the status by parameter and one literal 403
    When the gate confronts it
    Then it returns Fail naming 403 and never naming 404, because a declared value may be
      emitted through a call textual reading cannot reach
```

#### CSDCN-B11 — Without a declared dialect the verdict is Pending

```gherkin
    Given a project whose Structure declares no http_status lexicon
    And a handler that emits 200
    When the gate confronts it
    Then it returns Pending, because the meter does not fake conformity nor guess the stack
```

#### CSDCN-B12 — An explicit opt-out of the http_status field is honoured

```gherkin
    Given a project whose Structure waives the http_status field of the dialect
    When the gate confronts a spec with an output contract
    Then it returns Skip, because the waiver is declared and localised, not a silent absence
```

#### CSDCN-I01 — The lexicon comes from the project's dialect, not from the gate

```gherkin
    Given a Go handler that writes its refusal through the net/http writer
    And a Structure declaring the go dialect family
    When the gate confronts it
    Then it returns Fail naming 403, because embedding one stack's syntax would make the
      gate silent on every other one
```

#### CSDCN-I02 — A dialect declared by hand teaches the gate its own lexicon

```gherkin
    Given a Structure declaring the http_status pattern directly, with no family
    And a Ruby handler that renders 422 through that pattern
    When the gate confronts it
    Then it returns Fail naming 422, because the agnosticism cannot stop at the built-in families
```

#### CSDCN-I03 — A named constant is worth the number it means

```gherkin
    Given a handler whose refusal is written as the named forbidden constant and never as digits
    When the gate confronts it
    Then it returns Fail naming 403, because reading only digits would approve every
      handler written with constants
```

#### CSDCN-X01 — The gate does not demand the generic ranges

```gherkin
    Given an output contract that declares no 5xx range at all
    And a handler whose only failure path is the 500 of its catch block
    When the gate confronts it
    Then it returns Pass, because treating a range as a status would mean guessing which
      numbers it covers
```

#### CSDCN-X02 — The gate does not judge when each status is right

```gherkin
    Given an output contract declaring 403 for a branch a reviewer would call a 404
    And a handler that emits exactly that 403
    When the gate confronts it
    Then it returns Pass, because the ruler is the correspondence between two sets of
      numbers, not judgement about the design
```

#### CSDCN-X03 — The gate does not charge the phantom side under a dynamic status

```gherkin
    Given a handler whose envelope helper receives the status by parameter
    And an output contract declaring a status no literal in the code shows
    When the gate confronts it
    Then it never names that status, because mass false positives are what makes a team
      turn the gate off
```


## DCRQD — DocRequired — the aggregated document the unit must feed









#### DCRQD-B01 — A mandatory document that does not exist fails

```gherkin
    Given a project declaring a document as mandatory for this unit's layer
    And that document does not exist on disk
    When the gate confronts the unit
    Then it returns Fail
```

#### DCRQD-B02 — A document that exists and does not mention the unit fails

```gherkin
    Given the mandatory document exists and never names this unit
    When the gate confronts the unit
    Then it returns Fail, because existence alone would approve an empty file created
      to silence the gate
```

#### DCRQD-B03 — A mention by the identity code counts as documented

```gherkin
    Given the mandatory document cites the unit's identity code
    When the gate confronts the unit
    Then it returns Pass
```

#### DCRQD-B04 — A mention by the file name also counts

```gherkin
    Given the mandatory document cites the unit's file name and not its code
    When the gate confronts the unit
    Then it returns Pass, because the document speaks of the unit either way
```

#### DCRQD-B05 — Satisfying one of two duties is not enough

```gherkin
    Given two documents declared mandatory for this unit's layer
    And only one of them mentions the unit
    When the gate confronts the unit
    Then it returns Fail, because each document is charged on its own
```

#### DCRQD-B06 — Without a declaration nothing is charged

```gherkin
    Given a project that declares no mandatory document
    When the gate confronts the unit
    Then it returns Pass, because the ruler is what the project committed to, not what
      one supposes it owes
```

#### DCRQD-B07 — A layer with no trigger is not charged

```gherkin
    Given a mandatory document whose trigger names another layer
    And the document exists and does not mention this unit
    When the gate confronts the unit
    Then it returns Pass, because this layer triggers no duty
```

#### DCRQD-B08 — Aggregated, the verdict is one per document

```gherkin
    Given three units of a triggering layer and one mandatory document
    When the gate runs over the whole project
    Then it reports ONE verdict for that document, not one per unit
```

#### DCRQD-I01 — The duty starts from the spec, not from the code

```gherkin
    Given a unit whose spec and code both exist
    When the gate confronts the project
    Then the charge lands on the spec, because the spec is what declares the unit
```

#### DCRQD-I02 — The layer used is the UNIT's, not the node's

```gherkin
    Given a spec whose node layer is the spec layer
    And the unit it describes belongs to a triggering layer
    When the gate confronts it
    Then the duty is charged, because reading the node's layer would charge every spec
      of the project the same duty, or none
```

#### DCRQD-I03 — Without a map the aggregated verdict is skipped

```gherkin
    Given no graph built
    When the gate runs aggregated
    Then it does not approve, because approving without being able to look would stamp
      what was never measured
```

#### DCRQD-X01 — The gate does not understand the document's content

```gherkin
    Given the mandatory document cites the unit and describes it wrongly
    When the gate confronts the unit
    Then it returns Pass, because the gate separates "not documented" from "documented" —
      judging the quality of the documentation is another ruler
```

#### DCRQD-X02 — The gate does not decide which documents are mandatory

```gherkin
    Given a project whose Structure declares no duty for this layer
    And a document that a reviewer would consider obviously required
    When the gate confronts the unit
    Then it returns Pass, because inventing duties would charge what nobody committed to
```


## DSCDC — DocSelfContained — the spec has to stand on its own









#### DSCDC-B01 — A reference that brings the passage it announces passes

```gherkin
    Given a spec naming the path of a plan node and quoting the passage just below it
    When the gate confronts it
    Then it returns Pass, because the reader has the argument in hand and does not leave the page
```

#### DSCDC-B02 — A reference that only points is accused

```gherkin
    Given a spec whose line names a plan path, or its bare file name, or a revision code,
      and carries nothing of what it announces
    When the gate confronts it
    Then it returns Fail telling the author to bring the text, because the scaffold sends
      the reader to a file they do not have
```

#### DSCDC-B03 — A quotation counts in any written tradition

```gherkin
    Given a spec whose reference is followed by the passage in straight quotes, typographic
      quotes, guillemets, German low quotes or CJK corner brackets
    When the gate confronts each one
    Then every one returns Pass, because demanding Latin quotes would accuse a French or a
      Japanese project unjustly
```

#### DSCDC-B04 — A path inside a code fence is an example, not a reference

```gherkin
    Given a spec showing a command to run inside a fenced block, and that command names a plan path
    When the gate confronts it
    Then it returns Pass, because inside the fence the path is what the reader types, not
      where the reader is being sent
```

#### DSCDC-B05 — The spec citing its own path is identifying itself

```gherkin
    Given a spec whose body names its own file path
    When the gate confronts it
    Then it returns Pass, because it is sending nobody anywhere
```

#### DSCDC-B06 — A rule code is not a revision

```gherkin
    Given a spec whose line names a sibling rule code and explains what that rule proved
    When the gate confronts it
    Then it returns Pass, because the revision form is the one the doctrine reserves and
      accusing the other would charge the spec for naming its own subject
```

#### DSCDC-B07 — Only the spec is charged

```gherkin
    Given the same empty reference written in a plan, a feature, a test and a code node
    When the gate confronts each of them
    Then every one returns Skip, because the plan references sibling plans by function and
      the feature does not become documentation prose
```

#### DSCDC-B08 — A revision cited with an explanation on the same line passes

```gherkin
    Given a spec whose line names a revision code and then says, in substantive prose, what
      that revision changed
    When the gate confronts it
    Then it returns Pass, because the code is the label and the sentence is the content
```

#### DSCDC-B09 — With no map the confrontation is skipped

```gherkin
    Given a spec whose body names what looks like a plan path
    And no graph built
    When the gate confronts it
    Then it returns Skip, because the list of what counts as a reference comes from the map
      and nowhere else
```

#### DSCDC-I01 — The ruler matches structure, never vocabulary

```gherkin
    Given the same bare path reference written in English, Spanish, German and Japanese
    When the gate confronts each spec
    Then every one returns Fail, because a gate that matched words would pass in silence over
      the project written in the other language, and silence is worse than absence
```

#### DSCDC-I02 — The on-line explanation escape belongs to the revision, not to the path

```gherkin
    Given a spec whose line names a plan path and surrounds it with long prose that never says
      what is in the file
    When the gate confronts it
    Then it returns Fail, because a path is the place the person would have to go and no
      amount of surrounding prose says what waits there
```

#### DSCDC-I03 — The verdict names the line and shows what it says

```gherkin
    Given a spec whose fourth line carries a bare plan reference
    When the gate confronts it
    Then the finding carries that line number and an excerpt of the line, so the reader does
      not have to hunt for it
```

#### DSCDC-X01 — The gate does not judge whether the accompanying content is faithful

```gherkin
    Given a spec whose reference is followed by a quotation that says something else entirely
    When the gate confronts it
    Then it returns Pass, because the ruler is whether the reader is left with something to
      read — judging fidelity is another class of gate
```

#### DSCDC-X02 — The gate marks and does not block

```gherkin
    Given a spec with a bare reference
    When the gate confronts it
    Then the finding is recorded as informative and no command is offered to fix it, because
      rewriting a sentence is the work of whoever wrote it
```

#### DSCDC-X03 — Measuring explanation errs on the permissive side

```gherkin
    Given a spec whose line names a revision code followed by prose just past the threshold
      and saying little
    When the gate confronts it
    Then it returns Pass, because the threshold is a coarse ruler and mass false positives
      are what make someone switch an informative gate off
```


## DMDCD — DomainDeclared — the spec declares what the unit ACCEPTS, and who blocks the invalid









#### DMDCD-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec has a domain to declare
```

#### DMDCD-B02 — A spec without the domain section is failed

```gherkin
    Given a spec that catalogues rules and never opens the domain section
    And no declared waiver anywhere in it
    When the gate confronts it
    Then it returns Fail
    And the verdict names the waiver marker, so whoever reads it learns the declared way out
```

#### DMDCD-B03 — A section opened and left empty is failed

```gherkin
    Given a spec whose domain section has only the table header, with no data row
    When the gate confronts it
    Then it returns Fail, because opening the title without declaring anything is the same
      gap wearing the appearance of compliance
```

#### DMDCD-B06 — A row filled only with placeholders is not a declaration

```gherkin
    Given a spec whose domain section carries a single row reading "TODO" in every column
    When the gate confronts it
    Then it returns Fail, because the untouched template asserts nothing
```

#### DMDCD-B04 — An entry with no owner is failed, and the verdict names it

```gherkin
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it is left blank
    When the gate confronts it
    Then it returns Fail
    And the verdict names "chave", so the reader does not have to hunt for the orphan
```

#### DMDCD-B07 — An entry whose owner is named passes

```gherkin
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it reads "the interface, before calling"
    When the gate confronts it
    Then it returns Pass
```

#### DMDCD-B05 — A waiver with a written reason silences the gate

```gherkin
    Given a spec with no domain section
    And a waiver marker followed by the reason "receives only typed values from its own code"
    When the gate confronts it
    Then it returns Skip, and the reason stays in the spec as the record that someone looked
```

#### DMDCD-I01 — A bare waiver, with no reason, does not waive

```gherkin
    Given a spec with no domain section
    And a waiver marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a waiver with no why is the silence the gate exists to end
```

#### DMDCD-I02 — A non-answer in the owner column is not an owner

```gherkin
    Given a spec declaring the entry "chave"
    And the column that says who guarantees it reads "I do not validate (MTVRX-X04)"
    When the gate confronts it
    Then it returns Fail, because carrying a restriction into the owner column names nobody —
      it is the sentence that creates the orphan
```

#### DMDCD-X01 — The gate does not judge whether the declared domain is correct

```gherkin
    Given a spec declaring the entry "month" as accepting any text, which is wider than the real domain
    And the column that says who guarantees it names the caller
    When the gate confronts it
    Then it returns Pass, because the ruler here is the PRESENCE of the declaration —
      whether the accepted set matches reality is judgment, and judgment belongs to another gate
```

#### DMDCD-X02 — The gate does not read the code to check the validation exists

```gherkin
    Given a spec whose domain section is complete and every entry has a named owner
    And the code that realises it performs no validation at all
    When the gate confronts it
    Then it returns Pass, because this layer reads TEXT — crossing the declaration with the
      implementation belongs to the relational gate, which has the map
```


## EVFRV — EvidenceFresh — the score of this test holds against TODAY's code









#### EVFRV-B01 — An artifact that is not a test leaves without a verdict

```gherkin
    Given a node whose kind is code, spec or plan
    When the gate confronts it
    Then it returns Skip, because only a test carries an execution score
```

#### EVFRV-B02 — Without a built map the gate stays quiet

```gherkin
    Given a test node and no graph built
    When the gate confronts it
    Then it returns Skip, because there is no closure to walk and it will not approve
      what it could not look at
```

#### EVFRV-B03 — A test with no execution stamp is skipped, not failed

```gherkin
    Given a test that carries no record of ever having run
    When the gate confronts it
    Then it returns Skip, because a score that was never written cannot have expired
```

#### EVFRV-B04 — A test whose closure is intact passes

```gherkin
    Given a test stamped at the revision it ran on
    And every dependency it recorded still sits at the revision the run measured
    When the gate confronts it
    Then it returns Pass
```

#### EVFRV-B05 — The passing verdict says what it checked against

```gherkin
    Given a test whose intact closure holds one dependency
    When the gate confronts it
    Then the verdict states the size of the closure it walked, because that is the difference
      between "nobody looked" and "I looked and it stands"
```

#### EVFRV-B06 — A test whose dependency advanced a revision fails

```gherkin
    Given a test stamped at the revision it ran on
    And a dependency that has since moved to a newer revision
    When the gate confronts it
    Then it returns Fail, because the score is still written and stopped holding
```

#### EVFRV-B07 — The failing verdict names the culprit

```gherkin
    Given a test whose dependency has since moved to a newer revision
    When the gate confronts it
    Then the verdict names that dependency, so whoever fixes it knows what moved underneath
```

#### EVFRV-B08 — The failing verdict states the fix

```gherkin
    Given a test whose dependency has since moved to a newer revision
    When the gate confronts it
    Then the verdict says to run the test again, because accusing without saying what to do
      transfers the work to the reader
```

#### EVFRV-B09 — A test whose own file changed is reported separately from its closure

```gherkin
    Given a test whose own revision has moved since the run that stamped it
    When the gate confronts it
    Then the verdict reports the test's own file as changed, apart from any closure finding
```

#### EVFRV-B10 — The culprit list is truncated at five and the remainder counted

```gherkin
    Given a test whose stamped closure holds twenty dependencies and all of them moved
    When the gate confronts it
    Then the verdict lists five of them and states how many others there are
```

#### EVFRV-I01 — Absence of proof and expired proof are never the same finding

```gherkin
    Given a test that never ran
    When the gate confronts it
    Then it does not return Fail, because the absent proof is a different debt with a
      different fix — run it the first time, not revalidate it
```

#### EVFRV-I02 — A test that never ran is never approved either

```gherkin
    Given a test that never ran
    When the gate confronts it
    Then it does not return Pass, because approving would state a freshness nobody measured
```

#### EVFRV-I03 — Truncation never hides the size of the problem

```gherkin
    Given a test whose stamped closure holds twenty dependencies and all of them moved
    When the gate confronts it
    Then the verdict counts what it did not list, so the reader still learns how far it spread
```

#### EVFRV-X01 — The gate does not charge the absence of a green test

```gherkin
    Given a test that carries no record of ever having run
    When the gate confronts it
    Then it does not return Fail, because charging the missing test is the coverage gate's ruler
```

#### EVFRV-X02 — The gate does not run the test nor judge whether the change broke it

```gherkin
    Given a test whose dependency moved by a change that could not affect the behaviour
    When the gate confronts it
    Then it returns Fail all the same, because the gate measures whether the evidence still
      covers the current code — the cheap fix settles it for real instead of by opinion
```

#### EVFRV-X03 — The gate does not read the project's configuration

```gherkin
    Given a test with an intact closure and no configuration supplied at all
    When the gate confronts it
    Then it returns Pass, because the confronted truth lives in the map — a ruler depending
      on settings could be turned off by a default nobody chose
```


## LYBNL — LayerBoundary — a layer does not reach what is not its own









#### LYBNL-B01 — An artifact that is not code leaves without a verdict

```gherkin
    Given a node whose kind is spec, test or feature
    When the gate confronts it
    Then it returns Skip, because there is no import to forbid outside code
```

#### LYBNL-B02 — Content matching a forbidden pattern fails, naming line and reason

```gherkin
    Given a boundary forbidding the screens layer to import from the repositories
    And a screen whose second line imports straight from a repository
    When the gate confronts it
    Then it returns Fail naming line 2 and the declared reason, because a prohibition
      with no motive turns into ritual
```

#### LYBNL-B03 — A rule scoped to a layer charges only that layer

```gherkin
    Given the same boundary declared for the screens layer
    And a hook that imports from a repository
    When the gate confronts it
    Then it does not fail, because the hook is exactly who is allowed to reach the data
```

#### LYBNL-B04 — A rule with no layer holds for all code

```gherkin
    Given a boundary with no layer forbidding the raw clock
    And files of the screens, hooks and repositories layers each reading the raw clock
    When the gate confronts each of them
    Then every one fails, because a rule without a layer is how a global prohibition
      is declared
```

#### LYBNL-B05 — Severity warn records without failing, and the default is error

```gherkin
    Given a boundary marked severity warn and a file that violates it
    When the gate confronts it
    Then it returns Pending, and the same boundary with no severity returns Fail,
      because the default is error and warn is the per-rule maturation
```

#### LYBNL-B06 — A waiver with a written reason on the line waives that line

```gherkin
    Given a forbidden import carrying an allow-boundary marker with a written reason
    When the gate confronts it
    Then it returns Pass, and the acknowledged debt stays visible and dated in the code
```

#### LYBNL-B07 — The waiver also holds in the comment on the line above

```gherkin
    Given a forbidden import whose allow-boundary marker with reason sits on the line above
    When the gate confronts it
    Then it returns Pass, because an import has nowhere to carry a readable end-of-line
      comment and demanding it inline would push the author not to declare at all
```

#### LYBNL-B08 — A bare marker with no reason does not waive

```gherkin
    Given a forbidden import carrying an allow-boundary marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a bare marker is a silent way to quiet the gate
```

#### LYBNL-B09 — With no boundary declared the verdict is Pending, never Pass

```gherkin
    Given a project that declares no boundary at all
    When the gate confronts a code file
    Then it returns Pending naming what to declare, because pretending it checked is
      worse than saying what is missing
```

#### LYBNL-B10 — An invalid forbid pattern fails visibly

```gherkin
    Given a boundary whose forbid pattern does not compile
    When the gate confronts a code file
    Then it returns Fail explaining the config problem, because swallowed in silence it
      would switch the rule off with nobody knowing
```

#### LYBNL-B11 — The pattern is matched against the whole file, catching a multi-line import

```gherkin
    Given a boundary forbidding a named import from a package
    And a file whose import of that name is wrapped over several lines
    When the gate confronts it
    Then it returns Fail pointing at the line where the match starts
```

#### LYBNL-I01 — The same rule is expressible in six language dialects

```gherkin
    Given the same architectural rule written in the import dialect of TypeScript,
      Python, Go, Java, Rust and Ruby
    When the gate confronts a violation and a legitimate import in each dialect
    Then the violation fails and the legitimate import passes in all six, because the one
      who writes the pattern is the project and the engine knows no language
```

#### LYBNL-I02 — A single-line import of the same shape is still caught

```gherkin
    Given a boundary forbidding a named import from a package
    And a file whose import of that name fits on one line
    When the gate confronts it
    Then it returns Fail, because accusing only the wrapped form would accuse formatting
      rather than the violation
```

#### LYBNL-I03 — The waiver holds on any line of the matched stretch

```gherkin
    Given a multi-line import whose allow-boundary marker sits on the from line
    When the gate confronts it
    Then it returns Pass, because demanding the marker on the first line of the match
      would require the author to know where the regex started matching
```

#### LYBNL-I04 — A line anchor keeps holding per line

```gherkin
    Given a boundary whose pattern anchors the forbidden text to a whole line
    And a file carrying that text inside a string in the middle of a line
    When the gate confronts it
    Then it returns Pass, because the whole-file match changes the dot, not the meaning
      of the anchors
```

#### LYBNL-X01 — The gate does not decide which boundaries exist

```gherkin
    Given a project whose Structure declares no boundary
    And a screen importing straight from a repository, which a reviewer would forbid
    When the gate confronts it
    Then it does not fail, because inventing boundaries would charge what nobody
      committed to
```

#### LYBNL-X02 — The gate does not parse the language, it matches text

```gherkin
    Given a boundary whose pattern is plain text with no notion of imports
    And a file where the forbidden text appears outside any import statement
    When the gate confronts it
    Then it returns Fail, because the ruler is TEXT — understanding the import graph of
      every language would tie the engine to a set of ecosystems
```

#### LYBNL-X03 — The gate does not judge whether the boundary is the right one to draw

```gherkin
    Given a boundary forbidding something a reviewer would consider harmless
    And a file that matches it
    When the gate confronts it
    Then it returns Fail, because the ruler is the DECLARATION — whether the boundary is
      worth drawing is design judgment, and that belongs to whoever writes the Structure
```


## MRPRM — MarkerParity — the same rule has to appear at BOTH ends that fulfil it









#### MRPRM-B01 — A rule marked at both declared scopes passes

```gherkin
    Given a declaration whose scopes are the page tree and the server tree
    And the same rule name marked once in each of them
    When the gate confronts it
    Then it returns Pass, because the mapping still has its two ends
```

#### MRPRM-B02 — A rule missing from one end fails, and the verdict names the empty scope

```gherkin
    Given the rule marked in the page tree and absent from the server tree
    When the gate confronts it
    Then it returns Fail
    And the verdict names the server scope, because neither side looks wrong on its own
```

#### MRPRM-B03 — Two markings on the same side do not satisfy the gate

```gherkin
    Given the same rule marked twice inside the page tree and never in the server tree
    When the gate confronts it
    Then it returns Fail, because the two add up to the expected count and would hide
      exactly the mismatch the gate exists to catch
```

#### MRPRM-B04 — A rule left over at one end fails and is named

```gherkin
    Given both ends marked for one rule
    And a second rule marked only in the page tree, left behind when the server dropped it
    When the gate confronts it
    Then it returns Fail naming that leftover rule
```

#### MRPRM-B05 — Total absence of the prefix is not approval

```gherkin
    Given a declared prefix that appears nowhere in the tree
    When the gate confronts it
    Then it returns Pending asking to check marker_prefix, because absence is almost
      always a typo in the declaration and Pass would make the gate look vigilant while
      watching nothing
```

#### MRPRM-B06 — A declaration with no prefix returns Pending

```gherkin
    Given a gate declaration whose marker_prefix is empty
    When the gate confronts it
    Then it returns Pending asking for the prefix, because with no prefix there is
      nothing to confront
```

#### MRPRM-B07 — With no scopes declared the ruler is the count

```gherkin
    Given a declaration with no scopes and a required count of two
    And the rule marked in two files anywhere in the tree
    When the gate confronts it
    Then it returns Pass
    And a tree carrying only one of those markings returns Fail
```

#### MRPRM-B08 — The marking crosses language

```gherkin
    Given the rule marked in a TypeScript file of the page tree
    And marked again in a Go file of the server tree
    When the gate confronts it
    Then it returns Pass, because the mapping between the two ends is the same mapping
      whatever language writes each end
```

#### MRPRM-B09 — Ignored directories never count towards parity

```gherkin
    Given both declared ends marked
    And a third copy of the marking inside node_modules
    When the gate confronts it
    Then it returns Pass, and the vendored copy is not counted as an end
```

#### MRPRM-I01 — The failing verdict names the rule and the empty scope

```gherkin
    Given a rule whose server end was never marked
    When the gate confronts it
    Then the verdict carries both the rule name and the scope left empty, so the reader
      does not have to diff the two trees to find the orphan
```

#### MRPRM-I02 — What was not measured is never approved

```gherkin
    Given in turn a declaration with no prefix, one with neither count nor scopes, and
      a prefix that appears nowhere
    When the gate confronts each of them
    Then none of them returns Pass, because approving without having looked would stamp
      what was never measured
```

#### MRPRM-X01 — The gate does not read what each end actually does

```gherkin
    Given both ends marked with the same rule name
    And the page promising a list of items the handler does not erase
    When the gate confronts it
    Then it returns Pass, because presence is deterministic and agreement of meaning is
      not — this gate separates "one end" from "both ends"
```

#### MRPRM-X02 — The gate does not decide which rules live at two ends

```gherkin
    Given a project whose Structure declares no marker-parity gate for a rule a reviewer
      would consider obviously two-ended
    When the gate is asked to run
    Then nothing is charged, because the catalogue of two-ended rules belongs to the
      project and a gate that invented parities would charge what nobody committed to
```

#### MRPRM-X03 — Files outside the text extension list are not read

```gherkin
    Given both declared ends marked
    And a binary asset carrying the same byte sequence
    When the gate confronts it
    Then the asset is not scanned, because erring low costs a marking in an exotic place
      and erring high costs megabytes read on every walk
```


## OBHNB — ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit









#### OBHNB-B01 — A node that carries the trigger and is absent from the demanded file fails

```gherkin
    Given a project declaring an obligation triggered by the personal-data attribute
    And a node whose header carries that attribute and whose token never appears in the purge script
    When the gate confronts it
    Then it returns Fail carrying the declared reason for the duty, so the reader learns
      what the absence costs
```

#### OBHNB-B02 — A node that carries the trigger and does appear passes

```gherkin
    Given the same obligation and a node whose header carries the trigger
    And the purge script naming the token derived from that node
    When the gate confronts it
    Then it returns Pass, because the duty is fulfilled
```

#### OBHNB-B03 — A node without the trigger contracts no obligation

```gherkin
    Given a node whose header does not declare the trigger attribute
    And a purge script that never names it
    When the gate confronts it
    Then it returns Pass, because the duty is charged by what the node declares about
      itself, not by what it might resemble
```

#### OBHNB-B04 — A waiver exempts only when it carries a written reason

```gherkin
    Given a node that carries the trigger and is absent from the purge script
    When the gate confronts it once with the waiver followed by a reason and once with the
      waiver alone
    Then the first returns Pass and the second returns Fail, because the reason is what
      separates the honest exception from silence
```

#### OBHNB-B05 — A project with no declared obligation is skipped

```gherkin
    Given a project whose Structure declares no cross-cutting obligation
    When the gate confronts a node that carries a trigger-looking attribute
    Then it returns Skip, because inventing duties would charge what nobody committed to
```

#### OBHNB-B06 — An acknowledged debt with a written when yields Pending

```gherkin
    Given a node that carries the trigger and is absent from the purge script
    And a debt declaration naming the obligation and the phase in which it will be paid
    When the gate confronts it
    Then it returns Pending carrying that commitment, because the duty still holds and the
      record must stay visible in the report
```

#### OBHNB-B07 — A bare debt marker keeps failing

```gherkin
    Given the same unfulfilled node with a debt marker naming the obligation and nothing more
    When the gate confronts it
    Then it returns Fail, because a marker with no when assumes no debt — it only hides better
```

#### OBHNB-B08 — Waiver and debt stay distinct

```gherkin
    Given one unfulfilled node waiving the obligation with a reason and another acknowledging
      the debt with a when
    When the gate confronts both
    Then only the waiver passes, because the debt is still owed and cannot be stamped as fulfilled
```

#### OBHNB-B09 — The failing verdict offers the three ways out

```gherkin
    Given a node that carries the trigger, is absent from the purge script and declares nothing
    When the gate confronts it
    Then the verdict names fulfilling, waiving with a reason and acknowledging the debt with
      a when, so nobody has to guess what the gate will accept
```

#### OBHNB-I01 — The token is derived through the declared form

```gherkin
    Given a node named MetadataEntry
    When each declared identifier form is applied to it
    Then the raw form yields MetadataEntry, the screaming form METADATA_ENTRY, the snake form
      metadata_entry, the kebab form metadata-entry and a free template the composed token,
      because guessing the shape in the engine would put one project's mess inside the framework
```

#### OBHNB-I02 — A glob that matches no file produces no violation

```gherkin
    Given an obligation whose destination glob matches no file in the project
    And a node that carries the trigger
    When the gate confronts it
    Then it returns Pass, because accusing where there was nothing to read would stamp what
      was never measured
```

#### OBHNB-I03 — The node's own identified_as wins over the automatic form

```gherkin
    Given a node named MetadataEntry whose header declares it is referenced as a plural env var
    And an obligation whose automatic form would derive the singular one
    And the purge script naming only the plural declared by the node
    When the gate confronts it
    Then it returns Pass, because only the node knows the project's real irregularity —
      inverting this order accuses 28 correct models
```

#### OBHNB-X01 — The gate does not decide which obligations exist

```gherkin
    Given a project whose Structure declares no obligation about personal data
    And a node a reviewer would consider obviously purgeable
    When the gate confronts it
    Then it returns Skip, because the duties are the project's decision and a gate that
      invented them would be turned off
```

#### OBHNB-X02 — The gate does not understand what the destination does with the token

```gherkin
    Given a purge script that names the node's token in a dead branch and erases nothing
    When the gate confronts the node that carries the trigger
    Then it returns Pass, because the ruler is presence — separating forgotten from
      remembered is the defect this gate was built for, and judging the implementation is
      another ruler
```

#### OBHNB-X03 — A declaration written in the body is not read

```gherkin
    Given a node that carries the trigger and is absent from the purge script
    And a waiver with a full reason written far down in the document's prose instead of the header
    When the gate confronts it
    Then it returns Fail, because without that cut a quotation in the prose would waive an
      obligation nobody meant to waive
```


## OPQSP — OpenQuestions — a spec with an open question is not ready to implement









#### OPQSP-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has open decisions to demand
```

#### OPQSP-B02 — Whoever OPENED the section is confronted by its content

```gherkin
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because the content is what the ruler reads
```

#### OPQSP-B03 — An open item bars the spec

```gherkin
    Given a spec whose open-decisions section carries one catalogued question
    When the gate confronts it
    Then it returns Fail, because while there is a question the spec does not pass as ready
```

#### OPQSP-B04 — A section closed honestly releases the spec

```gherkin
    Given a spec whose open-decisions section is opened and carries no item
    When the gate confronts it
    Then it returns Pass, because saying "there is no question" differs from not having looked
```

#### OPQSP-B05 — An item marked as resolved does not block

```gherkin
    Given a spec whose question is marked resolved, citing the rule born from it
    When the gate confronts it
    Then it returns Pass, and the question stays in the trail instead of being swept away
```

#### OPQSP-B06 — A question with no code is charged

```gherkin
    Given a spec whose open-decisions section carries a question written without a code
    When the gate confronts it
    Then it returns Fail, because without identity the question is neither a traceable
      item nor survives a rewrite of the spec
```

#### OPQSP-B07 — The count of pending decisions reads the project's own lexicon

```gherkin
    Given a spec whose open-decisions section carries two questions
    And the project names that section with its own wording
    When the number of pending decisions is counted
    Then it answers two, because counting zero over a section named otherwise would
      assert "no pending decision" about a spec full of them
```

#### OPQSP-I01 — Prose is not an item

```gherkin
    Given a spec whose open-decisions section carries explanatory text and no catalogued item
    When the gate confronts it
    Then it returns Pass, because otherwise the author would learn to explain nothing
```

#### OPQSP-I02 — The section boundary is respected

```gherkin
    Given a spec carrying catalogued items in a section that FOLLOWS the open decisions
    And the open-decisions section itself is empty
    When the gate confronts it
    Then it returns Pass, because otherwise the whole spec would read as a section of decisions
```

#### OPQSP-I03 — Filling in what the question BECOMES does not close the question

```gherkin
    Given a spec whose question names the rule it is expected to become
    And the question is not marked resolved
    When the gate confronts it
    Then it returns Fail, because the intended destination and the answer given are two
      different things
```

#### OPQSP-X01 — The gate does not judge whether the question is good

```gherkin
    Given a spec whose only open question is trivial
    When the gate confronts it
    Then it returns Fail all the same, because the ruler is deterministic — an open item
      exists, or it does not; judging the merit of a doubt belongs to another gate
```

#### OPQSP-X02 — A spec with no section is a pending item, and the verdict teaches the way out

```gherkin
    Given a spec with rules catalogued and no open-decisions section at all
    When the gate confronts it
    Then it returns Pending, because the absence does not tell "everything was decided"
      apart from "the section was deleted"
    And the verdict says how to close it: declare that there is no question, or write
      what is not yet decided
```


## PGNHN — PaginationHonored — what promises a SET does not return the first page in silence









#### PGNHN-B01 — A limit received from the caller passes

```gherkin
    Given an exported function whose signature takes the limit as a parameter
    When the gate confronts it
    Then it returns Pass, because the page is deliberate and whoever asked for it
      knows there is more
```

#### PGNHN-B02 — A limit hidden in a default value is accused

```gherkin
    Given an exported function whose name promises the whole set
    And the limit lives in a default value the caller never sees
    When the gate confronts it
    Then it returns Fail, because the hundred-and-first row is never processed and
      nobody is told
```

#### PGNHN-B03 — The NAME bounds the promise

```gherkin
    Given an exported function whose name promises no set at all
    And it returns a partial result
    When the gate confronts it
    Then it returns Pass, because there is no promise to break
```

#### PGNHN-B04 — Sibling functions that paginate are the proof by asymmetry

```gherkin
    Given a module whose sibling functions loop until the cursor is exhausted
    And one exported function returns a single page
    When the gate confronts it
    Then it returns Fail, because the author knew the pattern — the one that does not
      paginate is forgetfulness, not decision
```

#### PGNHN-B05 — A waiver with a written reason leaves the report

```gherkin
    Given an exported function carrying a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, and the function no longer appears in the report
```

#### PGNHN-B06 — The verdict offers the way out

```gherkin
    Given an exported function accused of hiding the limit
    When the gate confronts it
    Then the verdict names the waiver marker, so whoever reads it learns the declared
      way out instead of guessing
```

#### PGNHN-B07 — Without a declared dialect the verdict is undetermined

```gherkin
    Given a project that declares no dialect for its stack
    When the gate confronts any code
    Then it returns Pending, because approving without being able to read the code
      would stamp what was never checked
```

#### PGNHN-I01 — The ruler is agnostic across languages

```gherkin
    Given the same hidden-limit defect written in two different stacks
    And each project declares its own dialect
    When the gate confronts both
    Then both are accused, because the confronted truth — the name promises a set, the
      return is partial — belongs to no language
```

#### PGNHN-I02 — Where the construct is not recognised, the gate stays silent

```gherkin
    Given code whose pattern the declared dialect does not reach
    When the gate confronts it
    Then it accuses nothing, because a false positive here teaches the team to ignore
      the gate
```

#### PGNHN-I03 — A cursor with no loop does not count as pagination

```gherkin
    Given an exported function that returns the cursor and never walks it
    When the gate confronts it
    Then it returns Fail, because the consumer is left with the same slice, now wearing
      the appearance of completeness
```

#### PGNHN-I04 — A provider prefix in the name does not hide the promise

```gherkin
    Given an exported function whose name carries the provider prefix before the promise
    When the gate confronts it
    Then it returns Fail, because what the name says holds wherever it comes from
```

#### PGNHN-X01 — The gate does not invent a cursor the provider does not offer

```gherkin
    Given an exported function calling a dependency with no pagination mechanism
    When the gate confronts it
    Then it returns Pass, because accusing the absence of a mechanism the dependency
      lacks would hand the author a defect that is not theirs
```

#### PGNHN-X02 — The gate does not measure performance or page size

```gherkin
    Given an exported function that exposes its limit and returns a small page
    When the gate confronts it
    Then it returns Pass, because judging whether a hundred is many depends on the
      domain, and that is the project's decision
```


## PSDPL — PlanSourceDeclared — a plan that NAMES a source has to declare who builds it









#### PSDPL-B01 — A source that lives only in the prose is failed

```gherkin
    Given a plan whose prose names a source in bold
    And another plan seeds that source's adapter
    And this plan's needs list does not name that other plan
    When the gate confronts it
    Then it returns Fail, because the dependency existed only in the prose and survived
      the day the other plan dropped the adapter
```

#### PSDPL-B02 — The verdict names which source and where its adapter lives

```gherkin
    Given a plan naming a source whose adapter another plan seeds
    And that other plan absent from the needs list
    When the gate confronts it
    Then the verdict carries the source name and the plan that owns the adapter, so the
      fix is a declaration rather than an investigation
```

#### PSDPL-B03 — With the owning plan declared in needs the gate passes

```gherkin
    Given a plan naming a source whose adapter another plan seeds
    And that other plan listed in this plan's needs
    When the gate confronts it
    Then it returns Pass, because the prose and the declaration now say the same thing
```

#### PSDPL-B04 — Every source of the line is confronted on its own

```gherkin
    Given one source line naming two sources in bold
    And both adapters seeded by a plan this one does not declare
    When the gate confronts it
    Then it returns Fail naming both, because one declared source does not cover the rest
```

#### PSDPL-B05 — The source name matches the adapter regardless of case

```gherkin
    Given a plan naming the source in lower case
    And the seeded adapter file spelling it in mixed case
    When the gate confronts it
    Then the two are matched and the undeclared dependency is charged
```

#### PSDPL-B06 — A source whose adapter nobody seeds is not charged

```gherkin
    Given a plan naming a source
    And no plan in the map seeds an adapter for it
    When the gate confronts it
    Then it does not fail, because the source may belong to a plan that does not exist
      yet and the gate cannot invent a dependency
```

#### PSDPL-B07 — The plan that seeds the adapter is not charged for itself

```gherkin
    Given a plan that names a source and itself seeds that source's adapter
    When the gate confronts it
    Then it does not fail, because a plan does not depend on itself
```

#### PSDPL-B08 — A plan with no source line returns Skip

```gherkin
    Given a plan whose prose names no source at all
    When the gate confronts it
    Then it returns Skip, because there is nothing to confront and that is not approval
```

#### PSDPL-B09 — An artifact that is not a plan returns Skip

```gherkin
    Given a spec whose text carries a source line
    When the gate confronts it
    Then it returns Skip, because the gate has jurisdiction over plans only
```

#### PSDPL-I01 — What was not measured is never approved

```gherkin
    Given in turn a plan confronted with no graph built, and one whose map seeds no
      adapter at all
    When the gate confronts each of them
    Then neither returns Pass, because approving there would stamp a confrontation that
      never happened
```

#### PSDPL-I02 — A seeded file off the naming pattern owns nothing

```gherkin
    Given a plan that seeds a file whose name does not end in the adapter suffix
    And another plan naming that same source in bold
    When the gate confronts the consumer
    Then nothing is charged, because a wrong accusation costs more than a missed one —
      it teaches the reader to ignore the gate
```

#### PSDPL-X01 — The gate does not confront the order of the phases

```gherkin
    Given a plan that declares in needs the plan building its adapter
    And that owning plan scheduled in a later phase than this one
    When the gate confronts it
    Then it returns Pass, because ordering is the ruler of another gate and holding it
      in two places would let the two diverge
```

#### PSDPL-X02 — The gate does not demand a needs pointing at nothing

```gherkin
    Given a plan naming a source no plan in the map builds
    When the gate confronts it
    Then it does not fail, because charging it would demand a declaration pointing at
      nothing — the gate would be asking for a lie instead of catching one
```

#### PSDPL-X03 — The gate does not interpret what the source is for

```gherkin
    Given a plan whose prose names a source in bold only to say it was ruled out
    And another plan seeds that source's adapter
    When the gate confronts it
    Then it still fails, because the ruler is the bold name on the source line —
      deciding whether the plan really consumes it is interpretation, and interpretation
      is not what a blocking gate can hold
```


## PRHNP — ProgressHonest — the progress file tells the truth about the disk









#### PRHNP-B01 — An artifact that is not a plan leaves without a verdict

```gherkin
    Given a node whose kind is spec, code or test
    When the gate confronts it
    Then it returns Skip, because the gate is anchored on the plan, which is what the map holds
```

#### PRHNP-B02 — A plan with no companion progress file is skipped, not failed

```gherkin
    Given a plan with no progress file beside it
    When the gate confronts it
    Then it returns Skip, because "does the progress exist" is a different question
      from "is the progress true"
```

#### PRHNP-B03 — The skip for a missing companion says how to create it

```gherkin
    Given a plan with no progress file beside it
    When the gate confronts it
    Then the verdict names the command that creates the companion, so the reader does not
      have to look it up
```

#### PRHNP-B04 — A ticked item whose file does not exist is failed

```gherkin
    Given a progress item ticked as done and citing a path that is not on disk
    When the gate confronts the plan
    Then it returns Fail, because it declares done what is not
```

#### PRHNP-B05 — An open item whose file already exists is failed

```gherkin
    Given a progress item left open and citing a path that is already on disk
    When the gate confronts the plan
    Then it returns Fail, because it produces rework — somebody redoes what is done
```

#### PRHNP-B06 — A spec the plan seeds and the progress does not list is failed

```gherkin
    Given a plan whose checkbox items seed two specs
    And a progress file that lists only one of them
    When the gate confronts the plan
    Then it returns Fail, because a gate that only looks inside the file never sees
      what is missing from it
```

#### PRHNP-B07 — A checkbox item promising no file at all is failed

```gherkin
    Given a progress item that is an untouched template marker with no path
    When the gate confronts the plan
    Then it returns Fail, because an eternal open box makes the plan look unfinished forever
```

#### PRHNP-B08 — The ticked-but-absent finding is reported first

```gherkin
    Given a progress carrying both a ticked item with no file and an open item whose file exists
    When the gate confronts the plan
    Then the ticked-but-absent finding appears first, because whoever reads the board
      decides on it
```

#### PRHNP-B09 — An item in prose citing no path is not charged

```gherkin
    Given progress items reading "review with the team" and "agreed at the daily"
    When the gate confronts the plan
    Then it returns Pass, because there is nothing to confront
```

#### PRHNP-B10 — A progress that agrees with the disk on every item passes

```gherkin
    Given a ticked item whose file is on disk and an open item whose file is not
    When the gate confronts the plan
    Then it returns Pass
```

#### PRHNP-B11 — The verdict names each offending path

```gherkin
    Given a progress item ticked as done and citing a path that is not on disk
    When the gate confronts the plan
    Then the verdict names that path, so the reader does not have to diff the file
      against the disk by hand
```

#### PRHNP-I01 — The companion's path has one definition, derived from the scanner

```gherkin
    Given several plan paths, with and without an extension
    When the companion of each is derived
    Then each answer matches the scanner's, because a second constant here would diverge
      in silence and the gate would hunt for a file that does not exist
```

#### PRHNP-I02 — A seed is matched by path, never by the item's text

```gherkin
    Given a plan seeding a spec with a long description
    And a progress listing the same path under a shorter wording
    When the gate confronts the plan
    Then the seed is not accused as missing, because what identifies the item is the file
```

#### PRHNP-I03 — A spec mentioned in the plan's prose is not a seed

```gherkin
    Given a plan whose revision paragraph names a spec without a checkbox
    And a progress listing every spec the plan actually seeds
    When the gate confronts the plan
    Then it returns Pass, because the checkbox is the promise and the prose speaks of
      what already exists
```

#### PRHNP-I04 — A template file is never a seeded spec

```gherkin
    Given a plan whose checkbox item cites a template spec path
    And a progress that does not list it
    When the gate confronts the plan
    Then it returns Pass, because the mould is the shape work is poured into, not work to be done
```

#### PRHNP-X01 — The gate does not charge the existence of the progress file

```gherkin
    Given a plan that predates the progress mechanism and has no companion
    When the gate confronts it
    Then it does not return Fail, because merging "does it exist" with "is it true"
      would report two different debts as one finding
```

#### PRHNP-X02 — The gate does not judge the content of an item beyond the path

```gherkin
    Given a progress item whose description contradicts what the cited file contains
    And that file is on disk and the item is ticked
    When the gate confronts the plan
    Then it returns Pass, because the ruler here is the disk, which needs no opinion
```

#### PRHNP-X03 — The gate does not put the progress file into the map

```gherkin
    Given a progress file that changes on every delivery
    When the map is built
    Then the progress file is not a node, because a gate reaching it would charge every
      edit of the artifact that exists in order to change
```


## RLIMR — RuleImplemented — a spec catalogues rules, and the code shows it realized them









#### RLIMR-B01 — A spec whose rules the code ignores is accused, and the verdict names them

```gherkin
    Given a spec cataloguing three rules, of which the code marks only the first
    When the gate confronts it
    Then it returns Fail
    And the verdict names the rule left without realisation, so the reader does not
      have to diff spec against code by hand
```

#### RLIMR-B02 — A rule waived with a written reason closes the account

```gherkin
    Given a spec cataloguing a restriction the code cannot mark, because it is satisfied
      by the ABSENCE of code
    And the rule's row carries a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, because the declaration is the answer the gate asked for
```

#### RLIMR-B03 — A unit that predates the practice is a pending item, not a failure

```gherkin
    Given a spec whose rules carry no mark anywhere in the code
    And the project does not declare that it requires marking
    When the gate confronts it
    Then it returns Pending
    And the verdict NAMES the debt, instead of pretending approval
```

#### RLIMR-B04 — Declaring the requirement turns the pending item into a failure

```gherkin
    Given the same unmarked spec
    And the project declares in its Structure that marking is required
    When the gate confronts it
    Then it returns Fail, because declaring the requirement is the act of saying
      "here the migration is over"
```

#### RLIMR-B05 — A spec with no linked code is not this gate's subject

```gherkin
    Given a spec that no code realises yet
    When the gate confronts it
    Then it returns Skip, because without the piece on the other side there is no
      confrontation to make — and who accuses the absence is the triad gate
```

#### RLIMR-B06 — A waiver with no named rule covers every rule of the spec

```gherkin
    Given a spec whose waiver marker names no rule in particular
    When the gate confronts it
    Then it returns Pass, because the waiver was declared for the unit as a whole
```

#### RLIMR-I01 — Requiring the marking never punishes whoever already marks

```gherkin
    Given a spec whose every rule is marked in the code
    And the project declares that marking is required
    When the gate confronts it
    Then it returns Pass, because whoever did the work before the requirement cannot
      fail for having done it
```

#### RLIMR-I02 — The identity survives a rename

```gherkin
    Given a code marked with the identity the unit carried before being renamed
    When the gate confronts it
    Then it returns Pass, because losing the mark on a rename would turn identity
      stability into new debt
```

#### RLIMR-X01 — The gate does not judge whether the implementation is correct

```gherkin
    Given a spec whose rules are all marked in the code
    And the marked code does something other than what the rule describes
    When the gate confronts it
    Then it returns Pass, because the ruler here is deterministic — the mark exists, or
      the waiver exists with a reason; whether the code honours the rule is judgment
```

#### RLIMR-X02 — The gate does not demand a mark on EVERY rule

```gherkin
    Given a spec whose restrictions are satisfied by the absence of code
    And those rows declare their waiver with a reason
    When the gate confronts it
    Then it returns Pass, because absence has nowhere to receive a comment — demanding
      it would produce thousands of findings and teach the team to ignore the list
```


## RLTYR — RuleTypes — the rule VOCABULARY is extensible, but it must be DECLARED









#### RLTYR-B01 — A letter that is not declared in the vocabulary fails

```gherkin
    Given a vocabulary declaring the letters S, B and E
    And a spec cataloguing a rule under the letter P
    When the gate confronts it
    Then it returns Fail naming the letter P, because a letter the traceability cannot
      see makes the rule look covered when it is not
```

#### RLTYR-B02 — A declared letter under a claimed section passes

```gherkin
    Given a vocabulary declaring the letters S, B and E with their sections
    And a spec cataloguing rules only under those letters and sections
    When the gate confronts it
    Then it returns Pass
```

#### RLTYR-B03 — A section cataloguing rules under a title no letter claims fails

```gherkin
    Given a vocabulary whose declared sections do not include "Regras Inventadas"
    And a spec cataloguing a declared letter under that title
    When the gate confronts it
    Then it returns Fail naming the section, because the letter is claimed by a title and
      the catalogue has to sit where the vocabulary says it does
```

#### RLTYR-B04 — The same letter claimed by two terms is a conflict in the vocabulary

```gherkin
    Given a vocabulary where the letter E is claimed by the terms Error and Estado
    When the gate confronts any spec
    Then it returns Fail signalling the CONFLICT, because each letter belongs to ONE term
```

#### RLTYR-B05 — With no vocabulary declared the gate confronts the canonical letters

```gherkin
    Given a project that declares no vocabulary
    And a spec cataloguing a rule under the letter P, which is outside the canonical set
    When the gate confronts it
    Then it returns Fail saying it confronts the canonical vocabulary, because a gate that
      Skips forever gives the impression of a defence that does not exist
```

#### RLTYR-B06 — A heading that is the rule code itself is not a category section

```gherkin
    Given a spec whose heading is the rule code followed by its title
    When the gate confronts it
    Then it returns Pass, because that heading is the rule's own header and not a
      category that has to claim a letter
```

#### RLTYR-B07 — A section that only cites other sections' codes claims no letter

```gherkin
    Given a spec whose test-id section references codes of other sections in an inner column
    And that section title is claimed by no letter
    When the gate confronts it
    Then it returns Pass, because citing is not cataloguing
```

#### RLTYR-B08 — A section that defines a code in the first table cell is charged

```gherkin
    Given a spec whose unclaimed section carries a table row opening with a rule code
    When the gate confronts it
    Then it returns Fail, because the first cell is where a definition lives
```

#### RLTYR-B09 — A section declared as rule-cataloguing and filled without a code is Pending

```gherkin
    Given a vocabulary declaring "Eventos / Callbacks" as requiring a code
    And a spec whose section of that name carries a filled table and no code at all
    When the gate confronts it
    Then it returns Pending naming the section, because a row that asserts something
      verifiable and carries no code leaves the scenario with nothing to cite
```

#### RLTYR-B10 — A section whose table already carries the code is not charged

```gherkin
    Given a vocabulary declaring "Eventos / Callbacks" as requiring a code
    And a spec whose section of that name carries the rule code in its table
    When the gate confronts it
    Then it does not return Pending, because there is nothing left to charge
```

#### RLTYR-B11 — A declared section outside sections_require_code is not charged

```gherkin
    Given a vocabulary declaring "Variantes" under a letter but not as requiring a code
    And a spec whose section of that name merely enumerates values
    When the gate confronts it
    Then it does not return Pending, because demanding a rule of an index would invent a duty
```

#### RLTYR-B12 — A project that does not use sections_require_code changes no behaviour

```gherkin
    Given a vocabulary declaring "Eventos / Callbacks" with no requires-code marking
    And a spec whose section of that name carries a filled table and no code
    When the gate confronts it
    Then it does not return Pending, because the ruler is born opt-in and would otherwise
      accuse an entire existing base at once
```

#### RLTYR-I01 — A spec with no rule code at all is not this gate's problem

```gherkin
    Given a project with no vocabulary declared
    And a spec of pure prose that catalogues no rule
    When the gate confronts it
    Then it returns Pass, because charging the existence of a catalogued rule belongs to
      the spec-complete gate and doing it here would duplicate the ruler
```

#### RLTYR-I02 — The verdict names the letter and where to declare it

```gherkin
    Given a project with no vocabulary declared
    And a spec using a letter outside the canonical set
    When the gate confronts it
    Then the verdict carries both the letter and the vocabulary key to declare it in,
      because failing without saying what transfers the diagnosis to whoever reads it
```

#### RLTYR-I03 — A filled section with no code is the gap where the scenario loses its anchor

```gherkin
    Given a vocabulary declaring an events section as requiring a code
    And a spec whose events section is filled and carries no code
    When the gate confronts it
    Then it reports the finding, because without a code the scenario borrows the
      neighbouring section's and starts governing what is not its own
```

#### RLTYR-X01 — The gate does not decide which letters exist

```gherkin
    Given a vocabulary declaring a letter the canonical set does not contain
    And a spec cataloguing rules under that letter and its claimed section
    When the gate confronts it
    Then it returns Pass, because the vocabulary is extensible by design and the canonical
      letters are the fallback, not a ceiling
```

#### RLTYR-X02 — Without a declared vocabulary only the letter is charged

```gherkin
    Given a project that declares no vocabulary
    And a spec cataloguing a canonical letter under a title no vocabulary claims
    When the gate confronts it
    Then it returns Pass, because sections and terms only exist once the project declares
      them, and charging them against an implicit vocabulary would invent a rule nobody wrote
```

#### RLTYR-X03 — The gate does not judge whether the letter suits the rule

```gherkin
    Given a vocabulary declaring both S for state and B for behaviour
    And a spec cataloguing a plainly behavioural rule under the state letter and section
    When the gate confronts it
    Then it returns Pass, because which letter a rule deserves is editorial judgment —
      the ruler here is that the traceability can see it
```

#### RLTYR-X04 — The gate charges traceability, not format

```gherkin
    Given a spec whose unclaimed section carries prose and no rule code at all
    When the gate confronts it
    Then it returns Pass, because a section that catalogues nothing has no traceability
      to defend, whatever its shape
```


## SFMSP — SpecFeatureMatch — every requirement the spec DEFINES has at least one scenario









#### SFMSP-B01 — A requirement no scenario tags is failed and named

```gherkin
    Given a spec defining two requirements
    And a linked feature carrying a scenario for only one of them
    When the gate confronts it
    Then it returns Fail naming the uncovered requirement, because otherwise it would
      cross the whole pipeline with nothing verifying it
```

#### SFMSP-B02 — A requirement that has a scenario is not accused

```gherkin
    Given a spec defining one covered requirement and one uncovered requirement
    When the gate confronts it
    Then the verdict carries the uncovered one and never the covered one
```

#### SFMSP-B03 — With every requirement tagged the gate passes

```gherkin
    Given a spec defining two requirements
    And a linked feature carrying one scenario tagged for each
    When the gate confronts it
    Then it returns Pass
```

#### SFMSP-B04 — A code merely cited contracts no obligation

```gherkin
    Given a spec defining one requirement of its own
    And prose and a Dependency Table citing codes of other units
    And a feature covering only its own requirement
    When the gate confronts it
    Then it returns Pass, because a spec cites other units' codes all the time and
      without that distinction the gate would be a noise generator
```

#### SFMSP-B05 — A per-requirement waiver with a written reason waives

```gherkin
    Given a spec whose second requirement carries the waiver marker followed by a reason
    And a feature covering only the first requirement
    When the gate confronts it
    Then it returns Pass, and the reason stays in the spec as the record that it was a
      decision rather than forgetfulness
```

#### SFMSP-B06 — A bare per-requirement waiver does not waive

```gherkin
    Given a spec whose second requirement carries the waiver marker with nothing after it
    And a feature covering only the first requirement
    When the gate confronts it
    Then it returns Fail, because the waiver requires a written reason
```

#### SFMSP-B07 — A whole-spec waiver drags the waiver to every requirement

```gherkin
    Given a spec declaring with a reason that it has no feature
    And two requirements neither of which carries a waiver of its own
    When the gate confronts it
    Then it returns Skip, because with no feature no requirement of it can have a
      scenario — one decision, one place
```

#### SFMSP-B08 — Without the whole-spec waiver the same requirements keep failing

```gherkin
    Given the same two requirements and the same empty feature
    And no whole-spec waiver anywhere in the spec
    When the gate confronts it
    Then it returns Fail, because the drag cannot become a silent way of muting the gate
```

#### SFMSP-B09 — A bare whole-spec waiver drags nothing

```gherkin
    Given a spec whose whole-spec waiver marker has nothing written after it
    And a feature carrying no scenario
    When the gate confronts it
    Then it returns Fail, because otherwise the marker would be a switch that turns the
      gate off without accounting for it
```

#### SFMSP-B10 — A spec with no feature returns Skip

```gherkin
    Given a spec defining a requirement and no feature linked to it
    When the gate confronts it
    Then it returns Skip, because that absence is the ruler of the triad gate and
      accusing it here would print the same defect twice
```

#### SFMSP-B11 — Requirements are looked for across every linked feature

```gherkin
    Given a spec covered by two features
    And each feature carrying the scenario of a different requirement
    When the gate confronts it
    Then it returns Pass, because the requirement only needs to be in some of them
```

#### SFMSP-B12 — An artifact that is not a spec returns Skip

```gherkin
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec defines requirements
```

#### SFMSP-I01 — Every waiver requires a written reason

```gherkin
    Given in turn a bare per-requirement marker and a bare whole-spec marker
    When the gate confronts each of them
    Then both still fail, because a bare marker is a switch with no accounting and
      silence without a why is what the gate exists to end
```

#### SFMSP-I02 — Each gate accuses one thing

```gherkin
    Given a spec with a defined requirement and no feature at all
    When the gate confronts it
    Then it skips rather than failing, because the missing feature belongs to the triad
      gate and reporting it here would print the same defect twice
```

#### SFMSP-X01 — The gate does not judge whether the scenario proves the requirement

```gherkin
    Given a spec defining one requirement
    And a feature whose scenario carries the tag and asserts nothing at all
    When the gate confronts it
    Then it returns Pass, because the tag is deterministic and the shape of the
      assertion is the ruler of another gate
```

#### SFMSP-X02 — A cited code produces no accusation

```gherkin
    Given a spec whose Dependency Table cites three codes of other units
    And a feature covering only the code this spec defines
    When the gate confronts it
    Then none of the cited codes appears in the verdict, because a gate that cries wolf
      gets switched off — which costs more than the defect it was catching
```

#### SFMSP-X03 — The gate does not confront feature against test

```gherkin
    Given a spec whose every requirement has a scenario
    And no test binding any of those scenarios
    When the gate confronts it
    Then it returns Pass, because that edge already has its own watcher and this gate
      exists for the edge before it, which had none
```


## TRCMT — TriadComplete — the pieces that realize a spec EXIST









#### TRCMT-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only a spec has a triad to demand
```

#### TRCMT-B02 — A recognised layer leaves without a verdict

```gherkin
    Given a spec whose layer is declared with the declarative regime
    When the gate confronts it
    Then it returns Skip, because a recognised layer has neither spec nor triad by definition
```

#### TRCMT-B03 — Without a map the verdict is undetermined

```gherkin
    Given a spec and no graph built
    When the gate confronts it
    Then it returns Pending, because approving without being able to look would assert
      what was never measured
```

#### TRCMT-B04 — A spec with the three pieces linked passes

```gherkin
    Given a spec linked to its code, to its feature and to the test that proves it
    When the gate confronts it
    Then it returns Pass
```

#### TRCMT-B05 — A spec missing a piece is failed, and the verdict names which

```gherkin
    Given a spec linked to its code and to nothing else
    When the gate confronts it
    Then it returns Fail
    And the verdict names the feature and the test, and says where each one is born
```

#### TRCMT-B06 — The layer may waive a piece for every spec in it

```gherkin
    Given a layer declaring that the test edge is optional
    And a spec of that layer linked to its code and its feature only
    When the gate confronts it
    Then it returns Pass, because the waiver is declared in the Structure, in plain sight
```

#### TRCMT-B07 — The unit may waive a piece in its own spec, with a written reason

```gherkin
    Given a spec carrying a waiver marker for the test, followed by the reason
    And the spec is linked to its code and its feature
    When the gate confronts it
    Then it returns Pass, because the decision belongs to the unit and is written where
      whoever reads the spec will see it
```

#### TRCMT-I01 — The test is reached in two hops, through the feature

```gherkin
    Given a spec linked to a feature, and that feature linked to the test
    And no edge going straight from the spec to the test
    When the gate confronts it
    Then it returns Pass, because who points at the test is the FEATURE — checking it
      straight on the spec would report a missing test across the whole project
```

#### TRCMT-I02 — Waiving the test while the feature carries a scenario is a contradiction

```gherkin
    Given a spec whose test is waived
    And a linked feature carrying one scenario
    When the gate confronts it
    Then it returns Fail, because either the scenario is real and someone must prove it,
      or it should not exist
```

#### TRCMT-I03 — Waiving the test demands saying where the proof is, and the place must exist

```gherkin
    Given a spec whose test waiver points at a target that no file realises
    When the gate confronts it
    Then it returns Fail, because an orphan reference proves nothing
```

#### TRCMT-I04 — A waiver covers only the piece it declares

```gherkin
    Given a spec that waives the feature and is linked to its code only
    When the gate confronts it
    Then it returns Fail naming the test, because waiving one piece never waives the others
```

#### TRCMT-X01 — The gate does not confront whether the pieces MATCH one another

```gherkin
    Given a spec whose linked feature describes a behaviour the test does not prove
    And the three pieces exist and are linked
    When the gate confronts it
    Then it returns Pass, because matching is the work of the relational gates — this one
      exists precisely because they fail open when the piece is absent
```

#### TRCMT-X02 — The gate does not judge the QUALITY of any piece

```gherkin
    Given a spec linked to a feature with no scenarios and to an empty test
    When the gate confronts it
    Then it returns Pass, because the ruler here is EXISTENCE — confronting the content
      belongs to another gate, and mixing the two would fail by a criterion this one
      cannot measure
```


## VLANV — ValueAnchored — every value of a closed set points at the rule that justifies it, and the anchor carries the value









#### VLANV-B01 — An artifact that is not a spec leaves without a verdict

```gherkin
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because the gate reaches the code through the spec
```

#### VLANV-B02 — A value of a closed set with no anchor is failed

```gherkin
    Given a closed set declaring the values "15m" and "1h" with no anchor comment on any of them
    When the gate confronts it
    Then it returns Fail, because each value is a domain decision that left no address
```

#### VLANV-B03 — The verdict names the unanchored value

```gherkin
    Given a closed set whose value "15m" carries no anchor
    When the gate confronts it
    Then the verdict names "15m", so the reader does not have to hunt for which value it was
```

#### VLANV-B04 — A value whose anchor carries rule key and value passes

```gherkin
    Given a closed set where each value is preceded by an anchor holding the rule key and that same value
    When the gate confronts it
    Then it returns Pass, because the address exists and it checks out
```

#### VLANV-B05 — An anchor that asserts one value while the line says another is failed

```gherkin
    Given an anchor asserting "15m" written above a line whose literal is "5m"
    When the gate confronts it
    Then it returns Fail, because a lying anchor looks like traceability while pointing at the wrong place
```

#### VLANV-B06 — The verdict of a lying anchor shows both sides of the divergence

```gherkin
    Given an anchor asserting "15m" written above a line whose literal is "5m"
    When the gate confronts it
    Then the verdict carries both "15m" and "5m", because one side alone does not show the drift
```

#### VLANV-B07 — Lying anchors are reported before the unanchored ones

```gherkin
    Given a closed set carrying one lying anchor and one value with no anchor at all
    When the gate confronts it
    Then the lying anchor appears first in the verdict, because the absent anchor can be seen
      and the lying one cannot
```

#### VLANV-B08 — Without a declared value anchor pattern the gate skips

```gherkin
    Given a project that declares no value anchor pattern
    When the gate confronts a closed set
    Then it neither approves nor fails, because it cannot read and will not stamp what it did not measure
```

#### VLANV-B09 — The skip names the setting that enables the gate

```gherkin
    Given a project that declares no value anchor pattern
    When the gate confronts a closed set
    Then the verdict names the value_anchor setting, so the reader learns how to turn the gate on
```

#### VLANV-B10 — A declaration that opens no list is not a closed set

```gherkin
    Given a public symbol declared as a single scalar value on one line
    When the gate confronts it
    Then it returns Pass, because a scalar is not a closed set and there is nothing to anchor
```

#### VLANV-I01 — An anchor pattern with a single capture group does not enable the gate

```gherkin
    Given a declared anchor pattern that captures only the rule key
    When the gate confronts a closed set
    Then it neither approves nor fails, because without the second group the anchor asserts
      no value and the confrontation cannot happen
```

#### VLANV-I02 — A line carrying an anchor is never read as the end of the list

```gherkin
    Given a closed set whose second value carries an anchor that lies, after a first anchored value
    When the gate confronts it
    Then it returns Fail, because the brackets inside the anchor must not close the set early
```

#### VLANV-I03 — With no built map the verdict is pending

```gherkin
    Given no graph built
    When the gate confronts a spec
    Then it does not approve, because approving without being able to look would stamp
      what was never measured
```

#### VLANV-X01 — The gate does not judge whether the value is a good one

```gherkin
    Given a closed set whose value "999y" is anchored to a rule key and matches it exactly
    When the gate confronts it
    Then it returns Pass, because the ruler is that the decision has an address and the
      address does not lie — whether the value belongs is judgement
```

#### VLANV-X02 — The gate does not accuse a line whose literal it cannot read

```gherkin
    Given a closed set line carrying a computed expression instead of a quoted literal
    When the gate confronts it
    Then it returns Pass, because a false negative is better than a mass of false positives
      that would train the team to ignore the gate
```

#### VLANV-X03 — The gate does not decide what an anchor or a public symbol looks like

```gherkin
    Given a project whose declared export pattern does not match the way this file writes its set
    When the gate confronts it
    Then it returns Pass, because both shapes are declared by the project — inventing them
      would charge a convention nobody adopted
```



