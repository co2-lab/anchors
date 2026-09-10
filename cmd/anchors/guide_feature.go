package main

// featureGuide é a régua universal do artefato FEATURE: os cenários de comportamento
// que traduzem a spec em histórias verificáveis. Agnóstico de Gherkin/ferramenta.
const featureGuide = `# Guia de feature (os cenários de comportamento)

A feature traduz a spec em CENÁRIOS: histórias concretas do que a unidade faz. Ela é
a ponte entre a spec (o que deve ser) e o teste (a prova executável). Cada cenário
reusa um CÓDIGO da spec — nunca inventa um novo.

## A forma de um cenário

Todo cenário tem três tempos, sem misturá-los:
- CONTEXTO (passado) — o estado antes; todas as pré-condições, aqui dentro.
- GATILHO (presente) — o evento único que dispara o comportamento.
- RESULTADO (futuro) — o que se espera observar depois.
Não coloque ação no contexto, nem verificação no gatilho, nem nova ação no resultado.

## Regras da feature

- AUTO-CONTIDO. Cada cenário declara suas próprias pré-condições. Proibido "contexto
  compartilhado" no topo do arquivo — o leitor não deve caçar pré-condições longe do
  cenário. Um cenário é uma história completa.
- NOMES EXPLÍCITOS. "a tela de cadastro", "o botão Salvar" — nunca "a tela", "o botão".
  Linguagem de negócio, não de implementação.
- COBRE O CATÁLOGO DA SPEC. Para CADA item que a spec catalogou — cada estado,
  validação, ação, comportamento, permissão, restrição, mensagem, e cada campo do
  contrato/estado de dados — existe pelo menos um cenário. É assim que a feature
  prova que a spec foi honrada.
- CÓDIGO RASTREÁVEL. Cada cenário carrega o código do item da spec que ele cobre.
  Cenário sem código é invisível ao gate de cobertura.
- COPY LITERAL. Cenário de mensagem verifica o TEXTO EXATO do catálogo da spec
  (acentuação e pontuação reais). Para texto interpolado, verifique a parte estável.
- DATAS E DEFAULTS TÊM CENÁRIO. Data/timezone ganha cenário dedicado ("não devo ver
  o dia anterior"); default ganha cenário de ausência; campo condicional ganha
  cenário presente E ausente.

## Classifique cada cenário em dois eixos

1. NÍVEL DE TESTE — onde ele será provado da forma mais barata e isolada:
   • unidade (lógica pura, sem ambiente)
   • integração (um componente/módulo isolado, com suas bordas mockadas)
   • ponta-a-ponta (o fluxo real atravessando as fronteiras do sistema)
   Os níveis não são exclusivos; escolha o NATURAL para aquele comportamento.
   Sinal anti-E2E: se o observável é estado interno, animação ou scroll, NÃO é E2E —
   forçá-lo gera teste frágil ou vazio.
2. PRIORIDADE — exatamente UMA por cenário (de crítico a baixo), pela regra do MAIOR
   eixo: impacto no negócio, tipo de comportamento, frequência de uso, risco de
   regressão. Regra dura: cenário crítico EXIGE também um teste ponta-a-ponta, além
   do teste barato.

## O contrato feature↔teste

A feature e o teste executável são UM contrato. Cada passo do cenário aparece
espelhado no teste; se você editar um passo da feature sem atualizar o teste, o gate
deve quebrar. Não deixe os dois divergirem em silêncio.

## Anti-padrões (recuse-os)

- Pré-condição no topo do arquivo (contexto compartilhado) → mova para dentro do
  cenário.
- Cenário sem código → invisível ao gate; reuse o código da spec.
- "a tela", "o botão" → nomeie o quê, explicitamente.
- Forçar E2E sobre estado interno → teste frágil; use o nível natural.
- Copy parafraseada → verifique o texto literal do catálogo.

## Especialização do projeto

O formato concreto (a sintaxe dos cenários, os nomes das tags de nível/prioridade/
categoria, a ferramenta de execução) é do seu projeto — veja o guide de feature na
camada 'guide' do anchors.yaml. Este guia é a doutrina; siga o dialeto do projeto se
ele existir, e avise se não existir.
`
