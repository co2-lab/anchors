---
title: Catálogo e Índice de Gates
description: "O índice pesquisável e completo de todos os gates do Anchors, com tags de camadas, tipos de alvo e recomendações de uso."
---

Seja bem-vindo ao **Catálogo Geral de Gates** do Anchors.

Um **gate** é uma pergunta objetiva que o projeto faz a si mesmo em momentos-chave (como antes de commitar ou no CI), e cuja resposta o Anchors confronta de forma rigorosa. Para entender a filosofia por trás dos gates, leia [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/).

> [!TIP]
> **Pesquisa Rápida:** Você pode usar o atalho `Ctrl+K` ou `Cmd+K` para buscar qualquer gate pelo nome (ex: [`unit-complete`](/pt/docs/gates///unit-complete/), [`mutation-score`](/pt/docs/gates///mutation-score/), [`layer-boundary`](/pt/docs/gates///layer-boundary/)) em todo o portal de documentação.

---

## 🏷️ Navegação Rápida por Camada

Descubra quais gates fazem sentido para o seu tipo de trabalho:

- **[Camadas Regidas (usecase, service, domain, comando)](/pt/docs/layers/):** Exigem a [Unidade](/pt/docs/concepts/unidade/) completa ([`unit-complete`](/pt/docs/gates///unit-complete/), [`spec-feature-match`](/pt/docs/gates///spec-feature-match/), [`feature-test-match`](/pt/docs/gates///feature-test-match/), [`code-cataloged`](/pt/docs/gates///code-cataloged/)).
- **[Camadas Reconhecidas (infra, dao, types, doc)](/pt/docs/layers/):** Dispensam spec completa e focam em higiene e fronteiras ([`layer-boundary`](/pt/docs/gates///layer-boundary/), [`header-valid`](/pt/docs/gates///header-valid/)).
- **[Interface e Telas (UI, screen, componentes)](/pt/docs/layers/):** Focam em acessibilidade, visual e rotas ([`presentation-exhaustive`](/pt/docs/gates///presentation-exhaustive/), [`vr-baseline`](/pt/docs/gates///vr-baseline/), [`route-declared`](/pt/docs/gates///route-declared/), [`testid-consistent`](/pt/docs/gates///testid-consistent/)).
- **[Segurança e Governança Externa](/pt/docs/gates///#segurança-e-higiene-externa):** Bloqueiam segredos vazados e dependências com CVEs em todas as camadas ([`no-secret-leaked`](/pt/docs/gates///no-secret-leaked/), [`dependency-vulnerable`](/pt/docs/gates///dependency-vulnerable/)).

---

## 📚 Índice Completo de Gates por Categoria

### A Unidade e Estrutura

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`unit-complete`](/pt/docs/gates///unit-complete/) | `spec` | `Camadas Regidas`, `usecase`... | **Bloqueante** | Verifica se a especificação possui todas as peças que a realizam: código, feature e teste. |
| [`has-code`](/pt/docs/gates///has-code/) | `spec`, `feature` | `Camadas Regidas`, `Todas as Camadas` | **Bloqueante** | Garante que o arquivo possui um código de identidade de cenário e regra. |
| [`spec-sections`](/pt/docs/gates///spec-sections/) | `spec` | `Camadas Regidas`, `Spec` | **Bloqueante** | Verifica se a spec cataloga regras estruturadas (em cabeçalho, tabela ou lista) e usa o idioma correto. |
| [`non-empty`](/pt/docs/gates///non-empty/) | `feature`, `spec`, `doc` | `Camadas Regidas`, `Feature`... | **Bloqueante** | Garante que o arquivo não é um esqueleto vazio e que a feature possui cenários de verdade. |
| [`spec-feature-match`](/pt/docs/gates///spec-feature-match/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que toda regra catalogada na spec possui ao menos um cenário correspondente na feature. |
| [`feature-test-match`](/pt/docs/gates///feature-test-match/) | `feature` | `Camadas Regidas`, `Feature`... | **Bloqueante** | Garante que cada cenário da feature está implementado no teste por código e descrição. |
| [`scenario-identity`](/pt/docs/gates///scenario-identity/) | `feature` | `Camadas Regidas`, `Feature` | **Bloqueante** | Garante que cada cenário é distinguível, com código próprio e passos que não são cópias de outro. |
| [`scenario-asserts`](/pt/docs/gates///scenario-asserts/) | `feature` | `Camadas Regidas`, `Feature` | **Bloqueante** | Garante que o passo de desfecho (Then/Então) afirma um resultado observável concreto. |
| [`scenario-type-aligned`](/pt/docs/gates///scenario-type-aligned/) | `feature` | `Camadas Regidas`, `Feature` | **Bloqueante** | Garante que as tags de classificação do cenário batem com a letra do código (ex: @unit para regra B). |
| [`scenario-letter-declared`](/pt/docs/gates///scenario-letter-declared/) | `feature` | `Camadas Regidas`, `Feature` | **Bloqueante** | Verifica se a letra usada no código do cenário existe no vocabulário declarado do projeto. |
| [`code-reference-valid`](/pt/docs/gates///code-reference-valid/) | `spec`, `feature`, `code`, `test` | `Todas as Camadas` | **Bloqueante** | Garante que referências cruzadas a códigos de outras regras apontam para regras que realmente existem. |
| [`code-cataloged`](/pt/docs/gates///code-cataloged/) | `code` | `Camadas Regidas`, `code` | **Bloqueante** | Garante que todo símbolo público (função, tipo, export) exportado pelo código está catalogado na spec. |
| [`placeholder-filled`](/pt/docs/gates///placeholder-filled/) | `spec`, `plan`, `feature`, `doc` | `Todas as Camadas` | **Bloqueante** | Verifica se os placeholders deixados por geradores ou templates foram preenchidos. |

### Uso de Regras e Contratos

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`rule-uses-declared`](/pt/docs/gates///rule-uses-declared/) | `spec` | `Camadas Regidas` | Informativo | Verifica se cada regra declara explicitamente quais campos, validações e dependências utiliza. |
| [`rule-uses-resolve`](/pt/docs/gates///rule-uses-resolve/) | `spec` | `Camadas Regidas` | Informativo | Garante que o que a regra diz usar realmente existe declarado na spec. |
| [`rule-uses-implemented`](/pt/docs/gates///rule-uses-implemented/) | `spec` | `Camadas Regidas` | Informativo | Verifica se os campos que a regra diz usar aparecem e são consumidos no código governado. |
| [`contract-impact`](/pt/docs/gates///contract-impact/) | `spec` | `Camadas Regidas` | Informativo | Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes. |
| [`contract-status-declared`](/pt/docs/gates///contract-status-declared/) | `spec`, `code` | `Camadas Regidas`, `API`... | **Bloqueante** | Garante que o contrato lista os códigos de status que o código realmente retorna, e apenas esses. |
| [`domain-declared`](/pt/docs/gates///domain-declared/) | `spec` | `Camadas Regidas` | **Bloqueante** | Verifica se a spec declara o que a unidade aceita e quem bloqueia entradas inválidas. |
| [`count-honored`](/pt/docs/gates///count-honored/) | `spec`, `code` | `Camadas Regidas` | Informativo | Garante que asserções numéricas escritas na spec batem com os números reais no código. |
| [`pagination-honored`](/pt/docs/gates///pagination-honored/) | `spec`, `code` | `Camadas Regidas`, `API`... | Informativo | Garante que funções ou telas que prometem conjuntos paginados não retornam apenas a primeira página em silêncio. |
| [`route-declared`](/pt/docs/gates///route-declared/) | `spec` | `UI`, `Telas`... | **Bloqueante** | Garante que uma tela ou endpoint declara como se chega nela e nomeia seus vizinhos de navegação. |
| [`route-exists`](/pt/docs/gates///route-exists/) | `spec` | `UI`, `Telas`... | **Bloqueante** | Verifica se a rota declarada na especificação existe no registro de rotas da aplicação. |
| [`dependency-honored`](/pt/docs/gates///dependency-honored/) | `spec`, `code` | `Camadas Regidas` | **Bloqueante** | Verifica se os métodos prometidos na tabela de dependências da spec são realmente consumidos no código. |
| [`trigger-declared`](/pt/docs/gates///trigger-declared/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que triggers de conformidade e eventos citados existem no vocabulário declarado do projeto. |
| [`value-anchored`](/pt/docs/gates///value-anchored/) | `code`, `spec` | `Todas as Camadas` | **Bloqueante** | Garante que constantes ou chaves replicadas em vários arquivos carregam exatamente o mesmo valor. |

### Fronteiras Arquiteturais

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`layer-boundary`](/pt/docs/gates///layer-boundary/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Garante que nenhuma camada importa ou acessa módulos que a planta da casa proíbe. |
| [`sibling-guard`](/pt/docs/gates///sibling-guard/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Impede que módulos irmãos tratem parâmetros iguais de forma inconsistente ou se alcancem por caminhos proibidos. |
| [`proof-crosses-boundary`](/pt/docs/gates///proof-crosses-boundary/) | `spec`, `test` | `Camadas Regidas` | Informativo | Garante que quando uma regra afirma relação entre módulos, a prova de teste alcança a outra ponta. |
| [`circular`](/pt/docs/gates///circular/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Detecta dependências circulares entre módulos do projeto. |
| [`deadcode`](/pt/docs/gates///deadcode/) | `code` | `Todas as Camadas`, `code` | Informativo | Identifica funções, tipos, exports e arquivos órfãos não consumidos. |

### Prova e Execução

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`tests-pass`](/pt/docs/gates///tests-pass/) | `test` | `Todas as Camadas`, `test` | **Bloqueante** | Verifica se a suíte de testes passou com zero falhas no relatório ingerido. |
| [`evidence-fresh`](/pt/docs/gates///evidence-fresh/) | `test`, `spec` | `Camadas Regidas` | **Bloqueante** | Garante que o placar do teste continua fresco e foi medido contra o código atual. |
| [`line-coverage`](/pt/docs/gates///line-coverage/) | `code` | `Camadas Regidas`, `code` | Informativo | Verifica se a cobertura de linhas do arquivo atinge o piso mínimo exigido (ex: >= 70%). |
| [`coverage-delta`](/pt/docs/gates///coverage-delta/) | `code` | `Camadas Regidas`, `code` | **Bloqueante** | Garante que a alteração atual não reduziu a cobertura de linhas em relação ao baseline. |
| [`mutation-score`](/pt/docs/gates///mutation-score/) | `code` | `Camadas Regidas`, `code` | Informativo | Mede quantos mutantes injetados no código foram eliminados pela suíte de testes. |
| [`scenario-coverage`](/pt/docs/gates///scenario-coverage/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que cada cenário declarado na spec tem um teste correspondente que rodou e passou. |
| [`single-test-per-unit`](/pt/docs/gates///single-test-per-unit/) | `code` | `Camadas Regidas` | **Bloqueante** | Garante que uma unidade tem um único arquivo de teste por camada de teste (ou declara divisão com @split-test). |
| [`test-level-codes`](/pt/docs/gates///test-level-codes/) | `feature` | `Camadas Regidas`, `Feature` | **Bloqueante** | Garante que cada nível de teste referencia apenas códigos permitidos para seu escopo. |
| [`test-traceable`](/pt/docs/gates///test-traceable/) | `test` | `Camadas Regidas`, `test` | **Bloqueante** | Garante que todo teste ligado a uma feature declara no título o código do cenário que prova. |

### Apresentação e Telas

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`presentation-exhaustive`](/pt/docs/gates///presentation-exhaustive/) | `spec` | `UI`, `Telas`... | **Bloqueante** | Garante que todo valor de prop ou estado lido pela apresentação tem aparência explicitamente decidida. |
| [`presentation-conflict`](/pt/docs/gates///presentation-conflict/) | `spec` | `UI`, `Telas`... | **Bloqueante** | Garante que um prop e uma condição não levam a duas aparências conflitantes. |
| [`presentation-copy-single-source`](/pt/docs/gates///presentation-copy-single-source/) | `spec` | `UI`, `Telas` | **Bloqueante** | Garante que o texto exibido na apresentação vem de um código de mensagem central, e não de texto repetido. |
| [`presentation-observable`](/pt/docs/gates///presentation-observable/) | `spec` | `UI`, `Telas`... | **Bloqueante** | Garante que o elemento alterado pela apresentação possui identificador que testes conseguem apontar. |
| [`identity-consistent`](/pt/docs/gates///identity-consistent/) | `spec`, `code` | `UI`, `Telas`... | **Bloqueante** | Garante que a identidade da spec bate com o testID exposto e a imagem de baseline visual. |
| [`testid-consistent`](/pt/docs/gates///testid-consistent/) | `spec`, `code`, `test` | `UI`, `Telas`... | **Bloqueante** | Garante o contrato de testID: o código expõe, a spec declara e o teste consome exatamente o mesmo identificador. |
| [`testid-queried-exists`](/pt/docs/gates///testid-queried-exists/) | `test` | `UI`, `E2E`... | **Bloqueante** | Garante que todo testID buscado por um roteiro de fluxo E2E existe de verdade no código. |
| [`vr-baseline`](/pt/docs/gates///vr-baseline/) | `feature`, `test` | `UI`, `Telas`... | **Bloqueante** | Garante que cenários de regressão visual (-VR) possuem imagens de baseline capturadas. |

### Dublês e Mocks

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`mock-stamped`](/pt/docs/gates///mock-stamped/) | `test` | `Camadas Regidas`, `test` | **Bloqueante** | Garante que todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa. |
| [`mock-typed`](/pt/docs/gates///mock-typed/) | `test` | `Camadas Regidas`, `test` | **Bloqueante** | Garante que todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui. |
| [`mock-detect-cobre-o-dialeto`](/pt/docs/gates///mock-detect-cobre-o-dialeto/) | `test` | `Camadas Regidas`, `test` | Informativo | Avalia se a regex declarada para detectar dublês alcança todas as formas que o projeto usa. |

### Planejamento e Progresso

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`phase-exists`](/pt/docs/gates///phase-exists/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante que as fases citadas nas tarefas e dependências existem no plano. |
| [`phase-ordered`](/pt/docs/gates///phase-ordered/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante que a ordem e dependências entre fases do plano são consistentes e sem ciclos. |
| [`parent-valid`](/pt/docs/gates///parent-valid/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante que o campo parent aponta para uma fase ou plano pai existente. |
| [`plan-seeds-valid`](/pt/docs/gates///plan-seeds-valid/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante que as specs semeadas pelo plano miram camadas governadas válidas. |
| [`plan-source-declared`](/pt/docs/gates///plan-source-declared/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante que um plano que cita uma fonte externa declara explicitamente quem a constrói. |
| [`plan-revised`](/pt/docs/gates///plan-revised/) | `plan` | `Planos`, `doc` | **Bloqueante** | Garante visibilidade mútua de revisão entre planos substituídos e planos revisores. |
| [`plan-change-justified`](/pt/docs/gates///plan-change-justified/) | `plan`, `spec` | `Planos`, `Spec` | **Bloqueante** | Garante que um plano ou spec modificado declare no diff por que a mudança aconteceu. |
| [`open-questions-resolved`](/pt/docs/gates///open-questions-resolved/) | `spec`, `plan` | `Camadas Regidas`, `Planos` | **Bloqueante** | Garante que especificações com perguntas ou decisões em aberto não sejam liberadas para implementação. |
| [`progress-honest`](/pt/docs/gates///progress-honest/) | `plan`, `doc` | `Planos`, `doc` | **Bloqueante** | Garante que o arquivo de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco. |
| [`revision-orphans`](/pt/docs/gates///revision-orphans/) | `spec` | `Camadas Regidas` | **Bloqueante** | Identifica regras cujo significado mudou em uma revisão sem que isso tenha sido declarado formalmente. |
| [`revision-renumber`](/pt/docs/gates///revision-renumber/) | `spec`, `plan` | `Camadas Regidas`, `Planos` | **Bloqueante** | Renumera revisões em branch para evitar conflito quando a branch base já utilizou o mesmo número. |

### Doutrina e Feature Flags

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`spec-doctrine-exists`](/pt/docs/gates///spec-doctrine-exists/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que a doutrina referenciada pela spec com @realizes realmente existe no projeto. |
| [`doctrine-realized`](/pt/docs/gates///doctrine-realized/) | `product` | `produto`, `doutrina` | Informativo | Garante que as regras da doutrina de produto foram concretizadas em código e specs do projeto. |
| [`doctrine-not-duplicated`](/pt/docs/gates///doctrine-not-duplicated/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que regras de produto não foram duplicadas no corpo de specs locais. |
| [`spec-realizes-doctrine`](/pt/docs/gates///spec-realizes-doctrine/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que a spec declara e comprova como ela realiza a doutrina. |
| [`flag-scenario-grammar`](/pt/docs/gates///flag-scenario-grammar/) | `flag` | `flags` | **Bloqueante** | Garante que os cenários declarados em arquivos de feature flags seguem a gramática correta. |
| [`flag-scenarios-complete`](/pt/docs/gates///flag-scenarios-complete/) | `flag` | `flags` | **Bloqueante** | Garante que todos os valores possíveis da feature flag possuem cenários documentados. |
| [`flag-scenario-exists`](/pt/docs/gates///flag-scenario-exists/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que uma regra que cita @gated-by aponta para uma flag e cenário que realmente existem. |
| [`flag-scenario-governs`](/pt/docs/gates///flag-scenario-governs/) | `flag` | `flags` | **Bloqueante** | Garante que a flag governa as regras corretas na especificação. |
| [`flag-covered`](/pt/docs/gates///flag-covered/) | `flag` | `flags` | Informativo | Garante que cenários governados por flags possuem testes em todos os seus ramos. |

### Falhas e Governança

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`failure-declared`](/pt/docs/gates///failure-declared/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que possíveis falhas da unidade estão declaradas na especificação. |
| [`failure-handled`](/pt/docs/gates///failure-handled/) | `spec`, `code` | `Camadas Regidas` | **Bloqueante** | Garante que toda falha declarada na spec é tratada no código fonte. |
| [`failure-logged`](/pt/docs/gates///failure-logged/) | `code` | `Camadas Regidas`, `code` | Informativo | Garante que falhas tratadas emitem log adequado com contexto. |
| [`region-pair-honored`](/pt/docs/gates///region-pair-honored/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Garante que todo bloco #region no código fecha com um #endregion carregando o mesmo código. |
| [`rule-types`](/pt/docs/gates///rule-types/) | `spec`, `feature` | `Todas as Camadas` | **Bloqueante** | Garante que os prefixos e letras de códigos seguem a declaração de tipos de regras do projeto. |
| [`rule-implemented`](/pt/docs/gates///rule-implemented/) | `spec`, `code` | `Camadas Regidas` | **Bloqueante** | Garante que uma spec cataloga regras e o código mostra que as realizou. |
| [`code-language`](/pt/docs/gates///code-language/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Garante que o código não mistura idiomas de identificadores e comentários. |
| [`marker-parity`](/pt/docs/gates///marker-parity/) | `code`, `spec` | `Todas as Camadas` | **Bloqueante** | Garante que a mesma regra aparece em ambas as pontas que a realizam (ex: frontend e backend). |
| [`obligation-honored`](/pt/docs/gates///obligation-honored/) | `spec`, `code` | `Todas as Camadas` | **Bloqueante** | Garante que deveres regulatórios declarados fora da unidade são cumpridos. |
| [`header-valid`](/pt/docs/gates///header-valid/) | `spec`, `feature`, `code`, `test`, `doc` | `Todas as Camadas` | **Bloqueante** | Garante que o arquivo possui o bloco de cabeçalho @anchors com identidade mínima. |
| [`updated-at-atual`](/pt/docs/gates///updated-at-atual/) | `spec`, `feature`, `code`, `test`, `doc` | `Todas as Camadas` | **Bloqueante** | Garante que o campo updated_at no cabeçalho reflete a data real da alteração ou commit. |
| [`guide-has-checklist`](/pt/docs/gates///guide-has-checklist/) | `guide` | `doutrina`, `guide` | **Bloqueante** | Garante que guias de governança destilam suas regras em uma seção de checklist com itens CK1, CK2... |
| [`docs-fresh`](/pt/docs/gates///docs-fresh/) | `doc` | `doc`, `docs` | **Bloqueante** | Garante que a documentação compilada reflete as especificações atuais sem defasagem. |
| [`docs-covered`](/pt/docs/gates///docs-covered/) | `spec` | `Camadas Regidas` | Informativo | Garante que toda spec alcança alguma página da documentação compilada. |
| [`doc-required`](/pt/docs/gates///doc-required/) | `spec` | `Camadas Regidas` | Informativo | Garante que a unidade alimenta o documento agregado obrigatório do projeto. |
| [`doc-self-contained`](/pt/docs/gates///doc-self-contained/) | `spec` | `Camadas Regidas` | **Bloqueante** | Garante que a spec é autossuficiente e compreensível sem depender de contextos orais ou implícitos. |

### Segurança e Higiene Externa

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`no-secret-leaked`](/pt/docs/gates///no-secret-leaked/) | `code`, `test`, `doc`, `spec`, `feature`, `guide`, `plan` | `Todas as Camadas` | **Bloqueante** | Garante que nenhum segredo, token, senha ou chave privada entre no histórico do repositório. |
| [`dependency-vulnerable`](/pt/docs/gates///dependency-vulnerable/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Audita dependências e lockfiles contra bases públicas de vulnerabilidades conhecidas (CVEs). |
| [`sbom-generated`](/pt/docs/gates///sbom-generated/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Gera o inventário de software SBOM (CycloneDX ou SPDX) para compliance e auditorias. |
| [`license-compatible`](/pt/docs/gates///license-compatible/) | `code` | `Todas as Camadas`, `code` | **Bloqueante** | Impede a inclusão de dependências com licenças incompatíveis ou copyleft forte (como AGPL). |
| [`no-duplication`](/pt/docs/gates///no-duplication/) | `code` | `Todas as Camadas`, `code` | Informativo | Detecta blocos de código idênticos ou quase idênticos copiados em múltiplos arquivos. |
| [`spellcheck`](/pt/docs/gates///spellcheck/) | `code`, `doc`, `spec`, `feature` | `Todas as Camadas` | Informativo | Elimina erros de digitação e ortografia em identificadores, comentários e textos. |

### Julgamento por IA

| Gate | Alvo (`on`) | Camadas Compatíveis | Modo Recomendado | O que mede |
| :--- | :---: | :--- | :---: | :--- |
| [`regra-cumprida`](/pt/docs/gates///regra-cumprida/) | `spec`, `code` | `Camadas Regidas` | Informativo | Pergunta ao modelo se o trecho marcado no código realmente realiza o que a regra descreve. |
| [`no-test-prova-real`](/pt/docs/gates///no-test-prova-real/) | `spec` | `Camadas Regidas` | Informativo | Avalia se a justificativa apontada para uma dispensa @no-test é real e legítima. |

---

## 🚀 Como Escolher Gates para o seu Projeto

Não existe um conjunto universal estático. O Anchors recomenda pontos de partida medidos em projetos reais:

### 1. Projeto Novo (Começando pela Spec)
Ligue **tudo como informativo** (`blocking: false`) na primeira semana para observar o que o projeto tem:
- [`has-code`](/pt/docs/gates///has-code/) — Garante identidade nas specs.
- [`spec-sections`](/pt/docs/gates///spec-sections/) — Impede templates vazios.
- [`phase-exists`](/pt/docs/gates///phase-exists/) — Mantém o plano consistente.

### 2. Backend e Microsserviços
O foco principal é **fronteiras arquiteturais** e **segurança**:
- [`layer-boundary`](/pt/docs/gates///layer-boundary/) (`blocking: true`) — Garante isolamento entre camadas.
- [`no-secret-leaked`](/pt/docs/gates///no-secret-leaked/) (`blocking: true`) — Impede vazamento de chaves e senhas.
- [`unit-complete`](/pt/docs/gates///unit-complete/) (`blocking: true`) — Garante a unidade completa em camadas regidas.
- [`code-cataloged`](/pt/docs/gates///code-cataloged/) (`blocking: true`) — Todos os endpoints documentados.
- [`mutation-score`](/pt/docs/gates///mutation-score/) (`blocking: false` até estabilizar).

### 3. Aplicações com Interface Visual (Web e Mobile)
Acrescente gates que confrontam o que a tela promete:
- [`presentation-exhaustive`](/pt/docs/gates///presentation-exhaustive/) (`blocking: true`) — Todos os estados de UI decididos.
- [`testid-consistent`](/pt/docs/gates///testid-consistent/) (`blocking: true`) — Handles consistentes para testes E2E.
- [`vr-baseline`](/pt/docs/gates///vr-baseline/) (`blocking: true`) — Baselines visuais capturados.
- [`route-declared`](/pt/docs/gates///route-declared/) (`blocking: true`) — Telas com rotas mapeadas.

### 4. Sistemas Regulados (LGPD, Saúde, Financeiro)
- [`obligation-honored`](/pt/docs/gates///obligation-honored/) (`blocking: true`) — Cumprimento de deveres legais.
- [`marker-parity`](/pt/docs/gates///marker-parity/) (`blocking: true`) — A mesma regra refletida no frontend e no backend.
- [`license-compatible`](/pt/docs/gates///license-compatible/) (`blocking: true`) — Sem contaminação por licenças copyleft.

Para mais detalhes sobre como configurar cada bloco, veja [O anchors.yaml](/pt/docs/anchors-yaml///).
