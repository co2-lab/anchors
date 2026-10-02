---
title: Maturidade e Saúde do Projeto
description: "O conceito guarda-chuva de maturidade no Anchors, o comando anchors doctor e o roteiro para elevar a confiabilidade do seu software."
---

Construir software não é apenas fazer o código compilar e funcionar em uma demonstração rápida.

Um projeto pode "funcionar" hoje e ainda assim ser **imaturo**:
- Ele quebra a cada refatoração.
- Depende de conhecimentos que só existem na cabeça de uma pessoa.
- Uma IA gera novos arquivos que desrespeitam a arquitetura sem que ninguém perceba.
- A cobertura de testes parece alta, mas se você apagar linhas de validação nenhum teste falha.

No Anchors, o conceito que sintetiza a transição de um projeto frágil para um projeto sólido e perene é a **Maturidade**.

---

## 1. O que é Maturidade no Anchors?

A **Maturidade** é o conceito guarda-chuva do Anchors. Um projeto é maduro quando tem **todos os seus pilares implementados, ativos e rigorosos**:

```
                         ┌─────────────────────────┐
                         │   PROJETO MADURO        │
                         │   (Pronto e Confiável)  │
                         └────────────▲────────────┘
                                      │
           ┌──────────────┬───────────┴───────────┬──────────────┐
           │              │                       │              │
    ┌─────────────┐┌─────────────┐         ┌─────────────┐┌─────────────┐
    │  ESTRUTURA  ││PLANEJAMENTO │  . . .  │ PROPAGAÇÃO  ││  QUALIDADE  │
    │  Definida   ││  Ativo      │         │  Automática ││  Bloqueante │
    └─────────────┘└─────────────┘         └─────────────┘└─────────────┘
```

- **Projeto Imaturo**: Não tem camadas declaradas, não possui specs, seus gates de qualidade são frouxos ou inexistentes, e qualquer alteração pode gerar regressões silenciosas.
- **Projeto Maduro**: Tem sua planta de [Camadas](/pt/docs/layers/) respeitada, suas unidades cobertas pela [Unidade](/pt/docs/concepts/unidade/), testes com alta pontuação de [mutação](/pt/docs/gates//mutation-score/) e gates de segurança bloqueantes.

---

## 2. O Validador de Saúde Global: `anchors doctor`

Maturidade no Anchors não é uma sensação subjetiva ou um selo decorativo. Ela é **medida materialmente** pelo comando:

```sh
anchors doctor
```

O `anchors doctor` atua como um médico especialista examinando a saúde de todo o ecossistema:
1. **Verifica as ferramentas necessárias**: Confere se binários externos declarados no `anchors.yaml` (como `gitleaks`, `syft`, `osv-scanner`) estão instalados no ambiente do desenvolvedor ou no CI.
2. **Avalia a integridade do Grafo**: Verifica se há nós desconectados, specs órfãs ou arestas quebradas.
3. **Mede a vigência dos gates**: Alerta se um gate foi configurado com `when: [ci]`, mas o workflow do GitHub Actions nunca o executa.
4. **Identifica gates prontos para promoção**: Aponta quais gates informativos estão 100% limpos há dias e já podem virar bloqueantes com segurança.

---

## 3. O Roteiro de Amadurecimento: Passo a Passo

Se você está começando agora em um projeto existente, não tente atingir o nível máximo no primeiro dia. Siga o roteiro comprovado de amadurecimento:

### Passo 1: Declare a Planta da Casa (Dia 1)
Rode `anchors init` e declare suas [Camadas](/pt/docs/layers/) no `anchors.yaml`. Diga onde fica o código de negócio (`regime: regra`) e onde fica a infraestrutura (`regime: declarativo`).

### Passo 2: Ligue Gates Básicos como Informativos (Semana 1)
Ative gates de higiene com `blocking: false`:
- [`no-secret-leaked`](/pt/docs/gates//no-secret-leaked/) (este pode nascer bloqueante para evitar vazamentos).
- [`header-valid`](/pt/docs/gates//header-valid/)
- [`layer-boundary`](/pt/docs/gates//layer-boundary/)
Rode `anchors check` e observe o que o projeto já tem.

### Passo 3: Escreva Specs para as Unidades Críticas (Mês 1)
Comece a criar as specs e features para os fluxos mais vitais (ex: pagamentos, autenticação). Ligue os gates de identidade:
- [`has-code`](/pt/docs/gates//has-code/)
- [`spec-sections`](/pt/docs/gates//spec-sections/)
- [`unit-complete`](/pt/docs/gates//unit-complete/)

### Passo 4: Promova a Bloqueante e Endureça a Prova (Mês 2 em diante)
Conforme as violações forem zeradas, altere `blocking: true` nos gates estáveis e introduza medições de qualidade profunda, como o score de mutação do gate [`mutation-score`](/pt/docs/gates//mutation-score/).

---

## 4. Próximos Passos

- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): Entenda os estados de avaliação de cada gate.
- [O CLI](/pt/docs/cli//): Comandos completos do `anchors doctor` e `anchors check`.
- [Qualidade](/pt/docs/quality/): O pilar teórico da qualidade medida.
