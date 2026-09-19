<!-- anchors:generated from doct/camadas/gate.md.tmpl — DO NOT EDIT: run `anchors docs build` -->


# Camada: gate




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



