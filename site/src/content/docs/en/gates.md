---
title: "Gates Catalog"
description: "Comprehensive searchable catalog of all 95 Anchors quality gates, organized by category and architectural layers."
---

In **Anchors**, quality is not an afterthought or a subjective discussion in code review. It is an automated, continuous verification process executed by **Gates**.

This catalog lists all **95 quality gates** available in Anchors. Each gate has a dedicated documentation page with verification algorithms, verdict tables, yaml examples, and remediation steps.

---

## Gates by Layer Applicability

Before exploring by category, check which gates govern each [Project Layer](/docs/layers/):

- **[Governed Layers (usecase, service, domain, command)](/docs/layers/):** Require [The Unit](/docs/concepts/unit/) ([`unit-complete`](/docs/gates/unit-complete/), [`spec-feature-match`](/docs/gates/spec-feature-match/), [`feature-test-match`](/docs/gates/feature-test-match/), [`code-cataloged`](/docs/gates/code-cataloged/)).
- **[Recognized Layers (dao, infra, types)](/docs/layers/):** Exempt from business specs, governed by architectural boundaries ([`layer-boundary`](/docs/gates/layer-boundary/), [`circular`](/docs/gates/circular/)).
- **[UI & Presentation Layers (component, page)](/docs/layers/):** Governed by visual baselines ([`vr-baseline`](/docs/gates/vr-baseline/), [`testid-consistent`](/docs/gates/testid-consistent/)).
- **[All Layers](/docs/layers/):** Governed by hygiene and security gates ([`no-secret-leaked`](/docs/gates/no-secret-leaked/), [`deadcode`](/docs/gates/deadcode/)).

---

## AI Judgment & Synthetic

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`regra-cumprida`](/docs/gates/regra-cumprida/) | `RLUEX` | `advisory` | Pergunta ao modelo se o trecho marcado nThe code realmente realiza o que a regra descreve. |
| [`no-test-prova-real`](/docs/gates/no-test-prova-real/) | `RLUEX` | `advisory` | Evaluates whether a justificativa apontada para uma dispensa @no-test é real e legítima. |

## Architectural Boundaries

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`layer-boundary`](/docs/gates/layer-boundary/) | `LYBNL` | `blocking: true` | Enforces architectural import limits and dependency directions between project layers. |
| [`sibling-guard`](/docs/gates/sibling-guard/) | `SBGRD` | `blocking: true` | Prevents módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos. |
| [`proof-crosses-boundary`](/docs/gates/proof-crosses-boundary/) | `PCBPR` | `advisory` | Ensures that quando uma regra afirma relação entre módulos, a prova de teste alcança a outra ponta. |
| [`circular`](/docs/gates/circular/) | `EXCMX` | `blocking: true` | Detects circular dependency loops between units, packages, or architectural layers. |
| [`deadcode`](/docs/gates/deadcode/) | `EXCMX` | `advisory` | Identifica funções, tipos, exports e arquivos órfãos não consumidos. |

## Doctrine & Feature Flags

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`spec-doctrine-exists`](/docs/gates/spec-doctrine-exists/) | `DCTRN` | `blocking: true` | Ensures that a doutrina referenciada pelThe spec com @realizes realmente existe no projeto. |
| [`doctrine-realized`](/docs/gates/doctrine-realized/) | `DCTRN` | `advisory` | Ensures that as regras da doutrina de produto foram concretizadas em código e specs do projeto. |
| [`doctrine-not-duplicated`](/docs/gates/doctrine-not-duplicated/) | `DCTRN` | `blocking: true` | Ensures that regras de produto não foram duplicadas no corpo de specs locais. |
| [`spec-realizes-doctrine`](/docs/gates/spec-realizes-doctrine/) | `DCTRN` | `blocking: true` | Ensures that The spec declara e comprova como ela realiza a doutrina. |
| [`flag-scenario-grammar`](/docs/gates/flag-scenario-grammar/) | `FLSCF` | `blocking: true` | Ensures that os scenarios declarados em arquivos de feature flags seguem a gramática correta. |
| [`flag-scenarios-complete`](/docs/gates/flag-scenarios-complete/) | `FLSCF` | `blocking: true` | Ensures that todos os valores possíveis da feature flag possuem scenarios documentados. |
| [`flag-scenario-exists`](/docs/gates/flag-scenario-exists/) | `FLSCF` | `blocking: true` | Ensures that uma regra que cita @gated-by aponta para uma flag e scenario que realmente existem. |
| [`flag-scenario-governs`](/docs/gates/flag-scenario-governs/) | `FLSCF` | `blocking: true` | Ensures that a flag governa as regras corretas nThe specification. |
| [`flag-covered`](/docs/gates/flag-covered/) | `FLSCF` | `advisory` | Ensures that scenarios governados por flags possuem testes em todos os seus ramos. |

## Failure & Error Handling

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`failure-declared`](/docs/gates/failure-declared/) | `FLRAI` | `blocking: true` | Ensures that possíveis falhas da unidade estão declaradas nThe specification. |
| [`failure-handled`](/docs/gates/failure-handled/) | `FLRAI` | `blocking: true` | Ensures that toda falha declarada nThe spec é tratada nThe code fonte. |
| [`failure-logged`](/docs/gates/failure-logged/) | `FLRAI` | `advisory` | Ensures that falhas tratadas emitem log adequado com contexto. |
| [`region-pair-honored`](/docs/gates/region-pair-honored/) | `RPHRG` | `blocking: true` | Ensures that todo bloco #region nThe code fecha com um #endregion carregando o mesmThe code. |
| [`rule-types`](/docs/gates/rule-types/) | `RLTYR` | `blocking: true` | Ensures that os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto. |
| [`rule-implemented`](/docs/gates/rule-implemented/) | `RLIMR` | `blocking: true` | Ensures that umThe spec cataloga regras e The code mostra que as realizou. |
| [`code-language`](/docs/gates/code-language/) | `CDLNG` | `blocking: true` | Ensures that The code não mistura idiomas de identificadores e comentários. |
| [`marker-parity`](/docs/gates/marker-parity/) | `MRPRM` | `blocking: true` | Ensures that a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend). |
| [`obligation-honored`](/docs/gates/obligation-honored/) | `OBHNB` | `blocking: true` | Ensures that deveres regulatórios declarados fora da unidade são cumpridos. |
| [`header-valid`](/docs/gates/header-valid/) | `INCHN` | `blocking: true` | Ensures that The file possui o bloco de cabeçalho @anchors com identidade mínima. |
| [`updated-at-atual`](/docs/gates/updated-at-atual/) | `INCHN` | `blocking: true` | Ensures that o campo updated_at no cabeçalho reflete a data real da alteração ou commit. |
| [`guide-has-checklist`](/docs/gates/guide-has-checklist/) | `INCHN` | `blocking: true` | Ensures that guias de governança destilam suas regras em uma seção de checklist com itens CK1, CK2... |
| [`docs-fresh`](/docs/gates/docs-fresh/) | `DCFRD` | `blocking: true` | Ensures that a documentação compilada reflete as specifications atuais sem defasagem. |
| [`docs-covered`](/docs/gates/docs-covered/) | `DCCVD` | `advisory` | Ensures that todThe spec alcança alguma página da documentação compilada. |
| [`doc-required`](/docs/gates/doc-required/) | `DCRQD` | `advisory` | Ensures that a unidade alimenta o documento agregado obrigatório do projeto. |
| [`doc-self-contained`](/docs/gates/doc-self-contained/) | `DSCDC` | `blocking: true` | Ensures that The spec é autossuficiente e compreensível sem depender de contextos orais ou implícitos. |

## Planning & Progress

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`phase-exists`](/docs/gates/phase-exists/) | `PHORP` | `blocking: true` | Ensures that as fases citadas nas tarefas e dependencies existem no plano. |
| [`phase-ordered`](/docs/gates/phase-ordered/) | `PHORP` | `blocking: true` | Ensures that a ordem e dependencies entre fases do plano são consistentes e sem ciclos. |
| [`parent-valid`](/docs/gates/parent-valid/) | `PHORP` | `blocking: true` | Ensures that o campo parent aponta para uma fase ou plano pai existente. |
| [`plan-seeds-valid`](/docs/gates/plan-seeds-valid/) | `PSVPL` | `blocking: true` | Ensures that as specs semeadas pelo plano miram camadas governadas válidas. |
| [`plan-source-declared`](/docs/gates/plan-source-declared/) | `PSDPL` | `blocking: true` | Ensures that um plano que cita uma fonte externa declara explicitamente quem a constrói. |
| [`plan-revised`](/docs/gates/plan-revised/) | `PLRVP` | `blocking: true` | Garante visibilidade mútua de revisão entre planos substituídos e planos revisores. |
| [`plan-change-justified`](/docs/gates/plan-change-justified/) | `PCJPL` | `blocking: true` | Ensures that um plano ou spec modificado declare no diff por que a mudança aconteceu. |
| [`open-questions-resolved`](/docs/gates/open-questions-resolved/) | `OPQSP` | `blocking: true` | Ensures that specifications com perguntas ou decisões em aberto não sejam liberadas para implementação. |
| [`progress-honest`](/docs/gates/progress-honest/) | `PRHNP` | `blocking: true` | Ensures that The file de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco. |
| [`revision-orphans`](/docs/gates/revision-orphans/) | `RVORP` | `blocking: true` | Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente. |
| [`revision-renumber`](/docs/gates/revision-renumber/) | `RVRNR` | `blocking: true` | Renumera revisões em branch para evitar conflito quando a branch base já utilizou o mesmo número. |

## Presentation & UI

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`presentation-exhaustive`](/docs/gates/presentation-exhaustive/) | `PRSNT` | `blocking: true` | Ensures that todo valor de prop ou estado lido pela apresentação tem aparência explicitamente decidida. |
| [`presentation-conflict`](/docs/gates/presentation-conflict/) | `PRSNT` | `blocking: true` | Ensures that um prop e uma condição não levam a duas aparências conflitantes. |
| [`presentation-copy-single-source`](/docs/gates/presentation-copy-single-source/) | `PRSNT` | `blocking: true` | Ensures that o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido. |
| [`presentation-observable`](/docs/gates/presentation-observable/) | `PRSNT` | `blocking: true` | Ensures that o elemento alterado pela apresentação possui identificador que testes conseguem apontar. |
| [`identity-consistent`](/docs/gates/identity-consistent/) | `IDCND` | `blocking: true` | Ensures that a identidade dThe spec bate com o testID exposto e a imagem de baseline visual. |
| [`testid-consistent`](/docs/gates/testid-consistent/) | `TICTS` | `blocking: true` | Garante o contrato de testID: The code expõe, The spec declara e The test consome exatamente o mesmo identificador. |
| [`testid-queried-exists`](/docs/gates/testid-queried-exists/) | `TQETS` | `blocking: true` | Ensures that todo testID buscado por um roteiro de fluxo E2E existe de verdade nThe code. |
| [`vr-baseline`](/docs/gates/vr-baseline/) | `VRBSV` | `blocking: true` | Ensures that scenarios de regressão visual (-VR) possuem imagens de baseline capturadas. |

## Proof & Execution

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`tests-pass`](/docs/gates/tests-pass/) | `PRJTS` | `blocking: true` | Verifies whether a suíte de testes passou com zero falhas no relatório ingerido. |
| [`evidence-fresh`](/docs/gates/evidence-fresh/) | `EVFRV` | `blocking: true` | Ensures that o placar dThe test continua fresco e foi medido contra The code atual. |
| [`line-coverage`](/docs/gates/line-coverage/) | `INCHN` | `advisory` | Verifies whether a cobertura de linhas dThe file atinge o piso mínimo exigido (ex: >= 70%). |
| [`coverage-delta`](/docs/gates/coverage-delta/) | `INCHN` | `blocking: true` | Ensures that a alteração atual não reduziu a cobertura de linhas em relação ao baseline. |
| [`mutation-score`](/docs/gates/mutation-score/) | `PRJTS` | `advisory` | Measures test suite defect-catching quality by generating code mutants and checking if tests kill them. |
| [`scenario-coverage`](/docs/gates/scenario-coverage/) | `SFMSP` | `blocking: true` | Ensures that cada scenario declarado nThe spec tem um teste corresponding que rodou e passou. |
| [`single-test-per-unit`](/docs/gates/single-test-per-unit/) | `SNGTU` | `blocking: true` | Ensures that uma unidade tem um únicThe file de teste por camada de teste (ou declara divisão com @split-test). |
| [`test-level-codes`](/docs/gates/test-level-codes/) | `TLVCD` | `blocking: true` | Ensures that cada nível de teste referencia apenas códigos permitidos para seu escopo. |
| [`test-traceable`](/docs/gates/test-traceable/) | `TSTRT` | `blocking: true` | Ensures that todThe test ligado a uma feature declara no título The code do scenario que prova. |

## Rules Usage & Contracts

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`rule-uses-declared`](/docs/gates/rule-uses-declared/) | `RLUSG` | `advisory` | Verifies whether cada regra declara explicitamente quais campos, validações e dependencies utiliza. |
| [`rule-uses-resolve`](/docs/gates/rule-uses-resolve/) | `RLUSG` | `advisory` | Ensures that o que a regra diz usar realmente existe declarado nThe spec. |
| [`rule-uses-implemented`](/docs/gates/rule-uses-implemented/) | `RLIMR` | `advisory` | Verifies whether os campos que a regra diz usar aparecem e são consumidos nThe code governado. |
| [`contract-impact`](/docs/gates/contract-impact/) | `CTRIM` | `advisory` | Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes. |
| [`contract-status-declared`](/docs/gates/contract-status-declared/) | `CSDCN` | `blocking: true` | Ensures that o contrato lista os códigos de status que The code realmente retorna, e apenas esses. |
| [`domain-declared`](/docs/gates/domain-declared/) | `DMDCD` | `blocking: true` | Verifies whether The spec declara o que a unidade aceita e quem bloqueia entradas inválidas. |
| [`count-honored`](/docs/gates/count-honored/) | `CNHNC` | `advisory` | Ensures that asserções numéricas escritas nThe spec batem com os números reais nThe code. |
| [`pagination-honored`](/docs/gates/pagination-honored/) | `PGNHN` | `advisory` | Ensures that funções ou telas que prometem conjuntos paginados não retornam apenas a primeira página em silêncio. |
| [`route-declared`](/docs/gates/route-declared/) | `RTDCL` | `blocking: true` | Ensures that uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação. |
| [`route-exists`](/docs/gates/route-exists/) | `RTEXR` | `blocking: true` | Verifies whether a rota declarada nThe specification existe no registro de rotas da aplicação. |
| [`dependency-honored`](/docs/gates/dependency-honored/) | `DEPHN` | `blocking: true` | Verifies whether os métodos prometidos na tabela de dependencies dThe spec são realmente consumidos nThe code. |
| [`trigger-declared`](/docs/gates/trigger-declared/) | `TRDCT` | `blocking: true` | Ensures that triggers de conformidade e eventos citados existem no vocabulário declarado do projeto. |
| [`value-anchored`](/docs/gates/value-anchored/) | `VLANV` | `blocking: true` | Ensures that constantes ou chaves replicadas em vários arquivos carregam exatamente o mesmo valor. |

## Security & Hygiene

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`no-secret-leaked`](/docs/gates/no-secret-leaked/) | `EXCMX` | `blocking: true` | Ensures that nenhum segredo, token, senha ou chave privada entre no histórico do repositório. |
| [`dependency-vulnerable`](/docs/gates/dependency-vulnerable/) | `EXCMX` | `blocking: true` | Audita dependencies e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs). |
| [`sbom-generated`](/docs/gates/sbom-generated/) | `EXCMX` | `blocking: true` | Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias. |
| [`license-compatible`](/docs/gates/license-compatible/) | `EXCMX` | `blocking: true` | Impede a inclusão de dependencies com licenças incompatíveis ou copyleft forte (como AGPL). |
| [`no-duplication`](/docs/gates/no-duplication/) | `DUPLC` | `advisory` | Detecta blocos de código idênticos ou quase idênticos copiados em múltiplos arquivos. |
| [`spellcheck`](/docs/gates/spellcheck/) | `EXCMX` | `advisory` | Elimina erros de digitação e ortografia em identificadores, comentários e textos. |

## Test Doubles & Mocks

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`mock-stamped`](/docs/gates/mock-stamped/) | `MCSTM` | `blocking: true` | Ensures that todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa. |
| [`mock-typed`](/docs/gates/mock-typed/) | `MCTYM` | `blocking: true` | Ensures that todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui. |
| [`mock-detect-cobre-o-dialeto`](/docs/gates/mock-detect-cobre-o-dialeto/) | `MCSTM` | `advisory` | Evaluates whether a regex declarada para detectar dublês alcança todas as formas que o projeto usa. |

## The Unit & Structure

| Gate | Code | Default Blocking | Measures |
| --- | --- | --- | --- |
| [`unit-complete`](/docs/gates/unit-complete/) | `UNTCP` | `blocking: true` | Verifies that every governed spec has its complete unit: code, feature, and test. |
| [`has-code`](/docs/gates/has-code/) | `SCIDS` | `blocking: true` | Ensures that specification and feature files declare formal rule and scenario identity codes. |
| [`spec-sections`](/docs/gates/spec-sections/) | `SFMSP` | `blocking: true` | Verifies whether The spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto. |
| [`non-empty`](/docs/gates/non-empty/) | `FTMFT` | `blocking: true` | Ensures that The file não é um esqueleto vazio e que a feature possui scenarios de verdade. |
| [`spec-feature-match`](/docs/gates/spec-feature-match/) | `SFMSP` | `blocking: true` | Ensures that toda regra catalogada nThe spec possui ao menos um scenario corresponding na feature. |
| [`feature-test-match`](/docs/gates/feature-test-match/) | `FTMFT` | `blocking: true` | Ensures that cada scenario da feature está implementado nThe test por código e descrição. |
| [`scenario-identity`](/docs/gates/scenario-identity/) | `SCIDS` | `blocking: true` | Ensures that cada scenario é distinguível, com código próprio e passos que não são cópias de outro. |
| [`scenario-asserts`](/docs/gates/scenario-asserts/) | `SCASS` | `blocking: true` | Ensures that o passo de desfecho (Then/Então) afirma um resultado observável concreto. |
| [`scenario-type-aligned`](/docs/gates/scenario-type-aligned/) | `STASC` | `blocking: true` | Ensures that as tags de classificação do scenario batem com a letra dThe code (ex: @unit para regra B). |
| [`scenario-letter-declared`](/docs/gates/scenario-letter-declared/) | `SCLTR` | `blocking: true` | Verifies whether a letra usada nThe code do scenario existe no vocabulário declarado do projeto. |
| [`code-reference-valid`](/docs/gates/code-reference-valid/) | `CRVCD` | `blocking: true` | Ensures that referências cruzadas a códigos de outras regras apontam para regras que realmente existem. |
| [`code-cataloged`](/docs/gates/code-cataloged/) | `CDCTC` | `blocking: true` | Ensures that todo símbolo público (função, tipo, export) exportado pelThe code está catalogado nThe spec. |
| [`placeholder-filled`](/docs/gates/placeholder-filled/) | `PLCFL` | `blocking: true` | Verifies whether os placeholders deixados por geradores ou templates foram preenchidos. |
