# language: en
# @anchors
#   ref: TXSMT
#   updated_at: 2026-09-26
#   layer: feature

@TXSMT
Feature: TextSimilarity — how close two texts that should be equal are, weighted by what each word discriminates

  @TXSMT-B01 @unit-level
  Scenario: Words are upper-cased, and one-character and digit-only words are dropped
    Given the texts "borderRadius 9999" and "a Bc 12 d-ef"
    When they are tokenized
    Then the words are "BORDERRADIUS", and "BC" and "EF"

  @TXSMT-B02 @unit-level
  Scenario: A word in every text weighs nothing and a rarer word weighs the log of its rarity
    Given the corpus "xx yy" and "xx zz"
    When the weights are computed
    Then "XX" weighs 0 and "YY" weighs ln 2

  @TXSMT-B03 @unit-level
  Scenario: A single-text corpus falls back to the unweighted count
    Given the corpus made only of "toque dispara onPress"
    When that text is scored against itself
    Then the score is 1

  @TXSMT-B04 @unit-level
  Scenario: Texts differing only in case and punctuation are identical
    Given the texts "Toque no card dispara onPress" and "toque no card dispara onPress!"
    When they are classified
    Then the verdict is identical

  @TXSMT-B05 @unit-level
  Scenario: Both rulers above the threshold make the pair similar
    Given a corpus of four scenario titles
    And the pair "saving the draft keeps the typed title" and "the draft keeps the typed title after saving it"
    When they are classified
    Then both rulers are at least 0.5 and the verdict is similar

  @TXSMT-B06 @unit-level
  Scenario: Rulers that disagree make the pair borderline
    Given the real corpus of a component's scenario titles
    And the pair "rounded verdadeiro aplica raio de pílula" and "rounded aplica borderRadius 9999"
    When they are classified
    Then exactly one ruler reaches 0.5 and the verdict is borderline

  @TXSMT-B07 @unit-level
  Scenario: A shared rare word pulls a low-scoring pair to similar
    Given a corpus where "onCustomUnit" appears in one title only
    And the pair "Editar o campo personalizado dispara onCustomUnit" and "onCustomUnit recebe o texto digitado"
    When they are classified
    Then the verdict is similar

  @TXSMT-B08 @unit-level
  Scenario: Different subjects are divergent
    Given the real corpus of a component's scenario titles
    And the pair "Estado vazio exibe a mensagem de nenhuma conta cadastrada" and "tocar copiar abre o sheet de cópia seletiva"
    When they are classified
    Then the verdict is divergent

  @TXSMT-B09 @unit-level
  Scenario: The reported score is the larger of the two rulers
    Given the borderline pair of the real corpus
    When it is classified
    Then the score equals the larger of the Jaccard and the cosine
