<!-- @anchors
  code: TXSMT
  updated_at: 2026-09-26
  layer: infra
-->
# TextSimilarity — how close two texts that should be equal are, weighted by what each word discriminates

> **Code**: `TXSMT`

## Overview

When a scenario title and the test that should repeat it no longer match, the equality check already
failed; this unit tells the reader WHICH fix applies. Two texts that talk about the same thing with
different words only need one side rewritten; two texts about different things need someone to decide
which side is stale.

Counting shared words is not enough: in the scenarios of one component, words like "component" or
"props" appear almost everywhere and distinguish nothing, while the rare ones ("borderRadius",
"onChange") carry the subject. So each word is weighted by how rare it is in the corpus (inverse
document frequency): a word present in every text weighs nothing, a word present in one weighs most.

Two rulers measure the pair: a weighted Jaccard (the shared fraction of the vocabulary, which penalises
the longer text) and a weighted cosine (which normalises by length). When both say "same subject" the
verdict is firm; when they disagree the pair is borderline, and saying so is more honest than picking a
number. A rare word shared by just these two texts is structural evidence that pulls a low-scoring pair
up to similar. The unit is pure: the same pair and corpus always give the same answer.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the two texts | any text, including empty | — | this unit: an empty text shares nothing and scores 0 |
| the corpus | the texts of the same file the pair comes from | a corpus that does not contain the pair | the caller (the doctrine and feature-test gates) |

## Effects

| Effect | Description |
| --- | --- |
| `TXSMT-B01` | A text is split into upper-cased words on every character that is not a letter or a digit; words of one character and words made only of digits are dropped. (`Tokenize`) |
| `TXSMT-B02` | The weight of a word is the natural logarithm of the corpus size over the number of texts that contain it, so a word in every text weighs 0. (`Weights`) |
| `TXSMT-B03` | When no word of the pair carries weight (a homogeneous corpus, such as a single text), both rulers fall back to the unweighted count, so identical texts score 1 and not 0. (`Score`, `Cosine`) |
| `TXSMT-B04` | Two texts with the same runs of letters and digits in the same order, ignoring case and punctuation, are identical with score 1, before any ruler is applied. Identity reads every run, numbers and one-character words included — unlike the rulers, which read the words of `TXSMT-B01`. (`Classify`) |
| `TXSMT-B05` | When both rulers reach 0.5, the pair is similar. |
| `TXSMT-B06` | When exactly one of the two rulers reaches 0.5, the pair is borderline. |
| `TXSMT-B07` | When neither ruler reaches 0.5 but the two texts share a rare word (weight of at least 1.38, about a quarter of the corpus or less), the pair is similar. |
| `TXSMT-B08` | Otherwise the pair is divergent. |
| `TXSMT-B09` | The score reported with a non-identical verdict is the larger of the two rulers. |
| `TXSMT-B10` | Texts that differ only by a number are not identical: `radius 9999` against `radius 0`, or `2 rows` against `3 rows`, go on to the rulers, while `123` against `123` is identical with score 1. |
| `TXSMT-B11` | Two texts without any run of letters or digits (`""`, `!!`, `??`) are identical with score 1. |
| `TXSMT-B12` | The two rulers fall back to the unweighted count under the same condition (`TXSMT-B03`); when only one text of the pair has no weighted word, both score 0, so a text made only of words present in the whole corpus is divergent, not borderline. (`Score`, `Cosine`) |
| `TXSMT-B13` | A verdict's `String` is its English name (`identical`, `similar`, `borderline`, `divergent`), for logs and tests; a user-facing message translates the verdict itself. |

## Errors

none — the unit is pure computation over text; an empty text or corpus is a value (score 0, empty weights), not a failure.

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
