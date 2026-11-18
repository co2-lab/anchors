---
title: Anchors — Documentação Oficial
description: "Um framework spec-first de continuidade para desenvolvimento assistido por IA. Aprenda conceitos, camadas e gates."
template: splash
hero:
  tagline: Um framework spec-first — a spec vem antes do código, e continua sendo a verdade depois dele. Mantém um projeto coerente ao longo do tempo através de âncoras que guiam o trabalho na ida e o confrontam na volta, sem permitir mentiras sobre o código.
  actions:
    - text: Começar pelo Conceito
      link: /pt/docs/concept/
      icon: right-arrow
      variant: primary
    - text: Ver o Catálogo de Gates
      link: /pt/docs/gates/
      icon: document
    - text: Ver no GitHub
      link: https://github.com/co2-lab/anchors
      icon: external
---

## 🌟 O que é o Anchors?

Se você nunca usou o **Anchors**, pense nele como o **sistema imunológico e a bússola do seu repositório**.

Quando você programa com Inteligência Artificial (ChatGPT, Claude, Copilot, etc.), ela frequentemente sofre de **amnésia estrutural**: a cada nova sessão, ela esquece as decisões arquiteturais tomadas ontem, quebra regras de negócio silenciosamente e cria códigos órfãos que ninguém sabe como testar.

O Anchors resolve isso com três princípios inegociáveis:
1. **Spec-First**: A especificação nasce antes do código e continua sendo a dona da verdade depois dele.
2. **[A Unidade (The Unit)](/pt/docs/concepts/unidade/)**: Nenhuma regra de negócio existe sem [Spec](/pt/docs/spec///), [Feature](/pt/docs/concepts/unidade/), [Teste](/pt/docs/quality/) e Código conectados por [códigos de identidade](/pt/docs/concepts/rastreabilidade-e-codigos/).
3. **[Gates Automatizados](/pt/docs/concepts/gates-e-vereditos/)**: Portais de qualidade que seguram a corda do projeto antes do commit ou no CI, impedindo regressões e código desgovernado.

---

## 🧭 Por Onde Começar?

### 1. Se você quer ENTENDER como tudo funciona (A Teoria Descomplicada)
Leia os guias conceituais escritos de forma simples e direta para quem está chegando agora:
- [**O que é uma Âncora?**](/pt/docs/concepts/ancora/) — A metáfora da escalada: como as âncoras apontam, seguram o safepoint e demarcam o caminho.
- [**A Unidade**](/pt/docs/concepts/unidade/) — Por que código sem spec ou teste é débito técnico imediato.
- [**Camadas do Projeto**](/pt/docs/layers/) — A planta da casa: diferença entre Camadas Regidas (negócio) e Reconhecidas (infra/doc).
- [**Rastreabilidade e Códigos**](/pt/docs/concepts/rastreabilidade-e-codigos/) — Como códigos curtos (ex: `AUTH-B01`) conectam requisitos do início ao fim.
- [**Gates e Vereditos**](/pt/docs/concepts/gates-e-vereditos/) — Os 4 desfechos (`✓`, `✗`, `~`, `Skip`) e a regra de ouro: todo gate nasce informativo.
- [**Maturidade e Saúde**](/pt/docs/concepts/maturidade-e-saude/) — O comando `anchors doctor` e como transformar seu projeto em uma fortaleza.

### 2. Se você vai APLICAR o Anchors no seu código hoje (Operação Prática)
Siga o roteiro operacional:
- [**O CLI**](/pt/docs/cli///) — Como instalar o binário e rodar seus primeiros comandos.
- [**O fluxo de trabalho**](/pt/docs/workflow/) — Como trabalhar no dia a dia, da criação de uma spec até o merge seguro.
- [**O anchors.yaml**](/pt/docs/anchors-yaml///) — O arquivo central de configuração do seu projeto.
- [**Catálogo de Gates**](/pt/docs/gates///) — O catálogo com mais de 90 verificadores prontos para ativar.

### 3. A Doutrina Completa dos 7 Pilares
Para quem busca a formulação arquitetural profunda:
1. [**Estrutura de Projeto**](/pt/docs/structure/) — A planta da casa e as fronteiras arquiteturais.
2. [**Planejamento**](/pt/docs/planning/) — Como semear specs e fases ordenadas.
3. [**Spec**](/pt/docs/spec///) — A disciplina de especificação executável.
4. [**Tipos de Spec**](/pt/docs/spec/-types/) — Specs de interface, usecase, serviço e arquitetura.
5. [**Rastreabilidade**](/pt/docs/traceability/) — A fiação contínua entre requisitos e testes.
6. [**Propagação**](/pt/docs/propagation/) — A onda de alterações que mantém o sistema coerente.
7. [**Qualidade**](/pt/docs/quality/) — A teoria da qualidade medida e os limiares de bloqueio.

---

## 🚨 Em Caso de Emergência
- [**Congelar o projeto**](/pt/docs/freeze/) — O botão de pânico: como interromper o trabalho de todos os agentes quando uma inconformidade crítica precisar ser resolvida primeiro.
