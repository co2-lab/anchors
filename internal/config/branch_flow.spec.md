<!-- @anchors
  code: BFSBR
  updated_at: 2026-11-06
  layer: config
-->
# BranchFlow — which branch merges into which, and which branches take no commit

> **Code**: `BFSBR`

## Overview

A release flow goes one way: work branch, then review, then production, with hotfixes taken from
production. Nothing said so, and a flow run backwards — commits on the production branch, the review
branch aligned from it — went unnoticed until someone read the history. `branch_flow:` says it in the
configuration: for each base branch, the head branches a pull request may merge into it, and the
branches that take no commit of their own.

```yaml
branch_flow:
  merges:
    review:     [develop, "hotfix/*"]
    production: [review, "hotfix/*"]
    develop:    [review, production, "hotfix/*"]
  protected: [review, production]
```

A branch is named by a glob whose `*` does not cross `/`, or by a regular expression after `re:`. The
commands that confront a merge or a commit read this unit (`anchors branch`).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| `merges` | base pattern → list of head patterns | — | the configuration |
| `protected` | list of branch patterns | — | the configuration |
| a pattern | a glob, or `re:` and a regular expression | anything else | this unit: refuses the configuration (`BFSBR-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `BFSBR-B01` | A merge of a head into a base is allowed when the base's entry names a pattern the head matches; refused, the base's heads are named. The entry naming the base exactly wins over a pattern that matches it. (`MergeAllowed`) |
| `BFSBR-B02` | A base no entry names takes any head, unless it is protected: a protected branch takes only what an entry allows. With no `branch_flow:`, nothing is guarded. |
| `BFSBR-B03` | A branch pattern is a glob whose `*` does not cross `/` — `hotfix/*` names `hotfix/login`, not `hotfix/a/b` —, or a regular expression after `re:`, matched whole. (`BranchMatches`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BFSBR-E01` | A pattern of `merges` or `protected` is neither a glob nor a regular expression. | The configuration does not load, naming `branch_flow` and the pattern. | A pattern matching nothing would leave the flow it guards unguarded. |
