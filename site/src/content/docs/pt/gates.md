---
title: Os gates
description: O catálogo completo, como configurar, e o que ligar em cada tipo de projeto.
---

Um **gate** é uma pergunta que o projeto faz a si mesmo, e cuja resposta o
Anchors registra. A [Qualidade](/docs/qualidade/) explica a doutrina; esta
página é a referência prática: quais gates existem, o que cada um mede, e quais
ligar no seu projeto.

## Como um gate é declarado

```yaml
gates:
    - name: trinca-completa      # como ele aparece na saída
      on: [spec]                 # em que kind de arquivo ele roda
      check: trinca-completa     # o verificador interno que o realiza
      blocking: true             # reprovar BARRA o commit/CI?
      measures: 'a spec tem código, feature e teste que a realizam'
```

| campo | o que decide |
| --- | --- |
| `on` | os `kind` de arquivo em que o gate roda |
| `check` | o verificador **interno** (determinístico) |
| `run` | um comando **externo**, quando o gate é uma ferramenta do projeto |
| `blocking` | `true` barra; `false` informa |
| `measures` | o que ele mede, em uma frase — aparece na saída |
| `requires` | roda só nos alvos cujo conteúdo contém este texto |
| `when` | as fases em que ele roda: `pre-commit`, `pre-push`, `ci` |
| `ask` | a pergunta de um gate de **julgamento por IA** |

### Os três desfechos

| símbolo | significa |
| --- | --- |
| `✓` | passou |
| `✗` | reprovou — vira issue; barra se `blocking: true` |
| `~` | **indeterminado**: o gate não teve o que confrontar |

O `~` não é falha. Um gate de cobertura num arquivo sem teste não reprova — ele
não mede. Confundir os dois faz alguém "consertar" o que não está quebrado.

---

## Comece informativo

**Todo gate novo deve nascer com `blocking: false`.**

Não é timidez: é a única forma de saber o que ele vai acusar antes de ele barrar
o trabalho de alguém. Um gate que nasce bloqueante num projeto que já tem
débito reprova tudo no primeiro dia, e a saída barata vira desligá-lo — o que é
pior que nunca tê-lo ligado.

O caminho é: liga informativo → mede uma semana → conserta o que ele achou →
promove a bloqueante.

---

## O catálogo

### A trinca — a spec tem as peças que a realizam?

| gate | mede |
| --- | --- |
| `trinca-completa` | a spec tem código, feature e teste ligados |
| `spec-tem-codigo` | a spec carrega um código de identidade |
| `spec-completa` | a spec tem ao menos uma regra, sem placeholder |
| `spec-feature-match` | cada regra da spec tem cenário na feature |
| `feature-nao-vazia` | a feature tem cenário de verdade |

O `trinca-completa` é o que impede o defeito mais silencioso do Anchors: uma
spec sozinha atravessa **todos** os gates relacionais — eles falham *aberto*,
sem teste ligado não há o que confrontar — e o pipeline conclui "pode promover"
sobre trabalho que não existe.

### A identidade — dá para achar a regra?

| gate | mede |
| --- | --- |
| `codigo-catalogado` | todo símbolo exportado tem regra na spec, ou dispensa escrita |
| `code-reference-valid` | as referências a códigos apontam para regras que existem |
| `rule-types` | as letras dos códigos estão no vocabulário declarado |
| `teste-rastreavel` | o teste cita o código do cenário que prova |
| `scenario-asserts` | o cenário afirma algo, em vez de só executar |

### O planejamento — o plano ainda descreve a realidade?

| gate | mede |
| --- | --- |
| `fase-existe` | as fases citadas existem no plano |
| `fase-ordenada` | a ordem declarada entre fases é coerente |
| `plan-seeds-valid` | as specs semeadas pelo plano existem |
| `plano-alterado-justificado` | um plano/spec que MUDOU declara por quê |
| `plano-revisado` | a revisão está numerada e explicada |
| `parent-valido` | o `parent:` aponta para uma fase que existe |
| `open-questions-resolved` | as decisões em aberto foram decididas |

O `plano-alterado-justificado` merece nota. Ele olha o **diff**, não o conteúdo,
porque a deriva é silenciosa por construção: um plano corrigido em silêncio fica
perfeitamente válido — a inconsistência foi removida. O que denuncia não é o
estado do arquivo, é a mudança sem justificativa.

### A prova — o teste prova, ou só executa?

| gate | mede |
| --- | --- |
| `tests-green` | a suíte passa |
| `line-coverage` | a linha executou durante o teste |
| `coverage-delta` | a cobertura não caiu com esta mudança |
| `mutation-score` | se a linha mudasse, algum teste quebraria |
| `scenario-coverage` | cada cenário da spec tem teste verde |

**Cobertura e mutação não são a mesma coisa**, e confundi-las é o defeito que o
`mutation-score` existe para pegar. Um projeto medido recentemente tinha **100%
de cobertura de linha** e 47 mutantes sobreviventes — 47 alterações no código
que nenhum teste percebia.

### As fronteiras — quem conhece quem

| gate | mede |
| --- | --- |
| `layer-boundary` | ninguém importa quem a camada proíbe |
| `sibling-guard` | um módulo não alcança o irmão pelo caminho errado |
| `prova-cruza-fronteira` | o teste que cruza fronteira a declara |
| `dependency-honored` | a dependência declarada é a que existe |

### Os dublês — o mock diz a verdade?

| gate | mede |
| --- | --- |
| `mock-carimbado` | todo dublê declara o que ele finge ser |
| `mock-tipado` | o dublê respeita o contrato do que substitui |
| `mock-detect-cobre-o-dialeto` | o regex que reconhece dublê alcança as formas que o projeto usa |

### Segurança e conformidade

| gate | mede |
| --- | --- |
| `secret-nao-vazado` | nenhum segredo entra no repositório |
| `dependencia-vulneravel` | quantas CVEs conhecidas as dependências carregam |
| `sbom-gerado` | o inventário de dependências está publicado |
| `obligation-honored` | os deveres regulatórios declarados são cumpridos |
| `contract-status-declared` | o status de cada contrato externo está dito |

### Julgamento por IA

Três gates não têm resposta determinística — eles **perguntam**:

| gate | pergunta |
| --- | --- |
| `regra-cumprida` | o trecho marcado REALIZA o que a regra descreve? |
| `no-test-prova-real` | a prova apontada pelo `@no-test` exercita mesmo o comportamento? |
| `mock-detect-cobre-o-dialeto` | o padrão declarado alcança todos os dublês? |

O veredito é gravado com `anchors judge`, e tem três valores:

```sh
anchors judge <alvo> --gate <g> --verdict pass|fail|dispensado --reason "..."
```

O **`dispensado`** existe porque os outros dois mentem quando o alvo não existe.
Uma spec que declara `@TBD: code` afirma que o código ainda não foi escrito — e
aí `pass` afirmaria que ele realiza a regra, e `fail` reprovaria trabalho que
ninguém errou.

---

## O que ligar em cada tipo de projeto

Não existe conjunto universal. O que segue são pontos de partida medidos em
projetos reais.

### Projeto novo, começando pela spec

Ligue **tudo como informativo** e não promova nada na primeira semana. Você
precisa ver o que o projeto tem antes de decidir o que barrar.

Os primeiros a promover, quando o projeto tiver 3–4 specs:

```yaml
- name: spec-tem-codigo      # sem identidade, nada é rastreável
- name: spec-completa        # spec com placeholder não decide nada
- name: fase-existe          # o plano aponta para fase que existe
```

### Backend / serverless (o caso do blue-eyes)

O que importa é **fronteira** e **segredo** — as duas coisas que quebram em
produção e não aparecem em teste.

```yaml
- name: layer-boundary            blocking: true   # o contrato não conhece quem o consome
- name: secret-nao-vazado         blocking: true   when: [pre-commit, pre-push, ci]
- name: licenca-compativel        blocking: true   when: [pre-push, ci]
- name: trinca-completa           blocking: true
- name: codigo-catalogado         blocking: true
- name: mutation-score            blocking: false  # informativo até saber quanto o projeto tem
- name: dependencia-vulneravel    blocking: false
```

O `mutation-score` fica informativo **de propósito**: exigir 80% antes de saber
quanto o projeto tem hoje produziria um número escolhido no escuro.

### Aplicação com interface

Acrescente os que confrontam o que a tela promete:

```yaml
- name: contract-status-declared  blocking: true   # a tela declara o estado de cada dado
- name: domain-declared           blocking: true   # o valor que chega errado tem tratamento
- name: pagination-honored        blocking: false  # a lista que pagina, pagina de verdade
- name: count-honored             blocking: false  # a contagem exibida é a contagem real
```

### Projeto com dado regulado (LGPD, saúde, financeiro)

```yaml
- name: obligation-honored        blocking: true
- name: marker-parity             blocking: true   # a mesma regra aparece nas DUAS pontas
- name: secret-nao-vazado         blocking: true
```

O `marker-parity` é o que impede a divergência mais cara desse tipo de projeto:
a tela promete apagar um dado, o backend apaga outro, e nada acusa — porque cada
lado está internamente coerente.

### Biblioteca / CLI

Fronteira importa menos; **prova** importa mais.

```yaml
- name: tests-green               blocking: true
- name: mutation-score            blocking: false → true quando estabilizar
- name: codigo-catalogado         blocking: true
- name: sem-duplicacao            blocking: false
```

---

## Gates externos: quando a ferramenta é do projeto

Nem todo gate é interno. Um comando do projeto vira gate com `run`:

```yaml
- name: licenca-compativel
  on: [code]
  scope: project
  run: 'bash scripts/anchors-licencas.sh'
  needs_tool: go-licenses
  install_hint: 'go install github.com/google/go-licenses@latest'
  blocking: true
  when: [pre-push, ci]
```

O `needs_tool` e o `install_hint` importam: sem eles, um gate que depende de
ferramenta ausente reprova com um erro de shell, e quem o vê não sabe se o
código está errado ou se falta instalar algo.

> **Um gate declarado que nunca roda é pior que gate nenhum**, porque consta na
> configuração como se protegesse. Se o gate declara `when: [ci]` e o seu CI não
> o executa, ele é decoração. Confira com `anchors doctor`.
