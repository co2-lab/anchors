package main

// specGuide é a régua universal da fase ESPECIFICAR. Destila a doutrina agnóstica de
// "o que é uma boa spec" — sem amarrar a stack. O projeto especializa no seu próprio
// guide de spec (a camada 'guide' do anchors.yaml); este é o piso comum.
const specGuide = `# Guia de spec (a régua da fase ESPECIFICAR)

A spec é a ORIGEM DA VERDADE de uma unidade. O plano decidiu QUE ela deve existir;
a spec decide COMO ela se comporta — em texto, nunca em código. Da spec nascem o
código, a feature e o teste; se algum deles diverge dela, é a spec que se confronta.

## A regra que não se quebra

DESCREVA COMPORTAMENTO, NÃO IMPLEMENTAÇÃO. A spec é textual e explicativa. Nada de
blocos de código, assinaturas de tipo, árvore de arquivos ou nome de biblioteca. Se
você sentiu vontade de colar código, isso é conteúdo do CÓDIGO, não da spec.

AGNÓSTICA AO FRAMEWORK. O corolário que mais escapa: uma spec não pode depender do
dialeto da ferramenta que hoje a implementa. Trocar de framework não pode invalidar a
spec — se invalidar, ela estava descrevendo a ferramenta, não a decisão.

  ✗ "usa '.identifier(['goalId'])'"      ✓ "a identidade é 'goalId', gerada pelo cliente"
  ✗ "'allow.owner()'"                     ✓ "só o dono lê e escreve o próprio registro"
  ✗ "expõe 'queryField'"                  ✓ "a consulta é exposta ao cliente"
  ✗ "índice secundário do Amplify"        ✓ "responde à pergunta 'X de um usuário'"

O TESTE: se o projeto migrasse de framework amanhã, esta frase continuaria verdadeira?
Se não, ela descreve a ferramenta. A decisão (o que é a identidade, quem enxerga o quê,
que perguntas o dado responde) sobrevive à migração; o nome da API não.

Nomear a ferramenta é legítimo em UM lugar: as Notas de Implementação, onde se registra
como a decisão foi materializada hoje — e fica claro que aquilo é o presente, não a regra.

## Identidade (o fio condutor)

Cada unidade especificável recebe um CÓDIGO curto, único e ESTÁVEL — ele não muda se
o arquivo for renomeado. Esse código prefixa todos os itens internos (regras,
estados, mensagens) e atravessa spec → feature → teste. É o que permite, dado um
código, achar a spec, o cenário e o teste. (O formato exato do código é do projeto;
a EXISTÊNCIA de um ID estável e único é universal.)

## Especialize a spec pela ORIGEM DA VARIAÇÃO

Uma spec não é toda igual — ela se especializa conforme o que faz a unidade variar:
- Unidade de FLUXO/ESTADO — varia por estado interno e navegação, busca dados. A
  spec foca em estados, transições e no contrato de dados.
- Unidade de COMPOSIÇÃO/CONTRATO — varia por entrada (parâmetros/props), não busca
  dados, emite eventos/callbacks. A spec foca no contrato de entrada, na matriz de
  variação e nos eventos emitidos.
(Os nomes "tela"/"componente" são de projeto; a distinção é universal.)

## Seções (adote as que se aplicam à unidade)

1. Cabeçalho com o CÓDIGO e o propósito em uma frase.
2. Visão geral — o que a unidade faz e para quem.
3. Regras — cada regra é um item catalogado, cada um com seu ID (permissões,
   validações, ações, comportamentos). Cada item é uma SEÇÃO com cabeçalho contendo o
   código (não um bullet solto numa lista) — porque a feature e o teste vão referenciar
   o item pelo cabeçalho. O FORMATO EXATO do cabeçalho (nível markdown, pontuação) é
   definido pelo gate de completude do SEU projeto — rode 'anchors check' para
   confirmar que a spec passa; se reprovar, a mensagem do gate diz o formato esperado.
   Tipar as regras ajuda a cobertura depois.
4. Estados — para unidades com estado: cada estado, sua condição de entrada, e o que
   fica visível/oculto/desabilitado nele.
5. Fluxo de estados — as transições entre os estados.
6. Contrato de dados — todo dado dinâmico exibido: origem, obrigatoriedade, formato,
   default. Cuidados que evitam bugs reais: TIMEZONE explícito em datas; enums
   traduzidos (backend↔interface) documentados; campo calculado aponta a função.
7. Estados de dados — um item por ramo condicional / por valor de enum (vazio, erro,
   carregando, cada variante).
8. Mensagens — o texto literal ao usuário é catalogado UMA vez, aqui. É a fonte única
   da copy; regras e ações apenas REFERENCIAM o código da mensagem, nunca duplicam o
   texto.
9. Notas de implementação — TODOs honestos, decisões, ambiguidades assumidas.
10. Histórico de alterações.

## Regras da spec

- Cataloga tudo que varia. Cada estado, validação, ação, mensagem e campo de dados
  vira um item com ID — porque a FEATURE vai exigir um cenário para cada um deles.
- Copy é fonte única. O texto ao usuário mora só nas Mensagens; nunca repita a frase
  em outra seção — referencie o código.
- Ambiguidade vira TODO, nunca chute (herdado do plano).
- Co-localização. A spec vive ao lado do artefato que descreve.

## Anti-padrões (recuse-os)

- Bloco de código na spec → é implementação; descreva o comportamento.
- Estado sem condição de entrada → não dá para testar; declare quando ele acontece.
- Copy duplicada em duas seções → vai divergir; centralize nas Mensagens.
- Requisito sem ID → fica invisível para a feature e o teste; dê um código.

## Especialização do projeto

Este é o piso universal. Seu projeto pode ter um guide de spec mais detalhado (veja a
camada 'guide' no anchors.yaml) com seções extras (acessibilidade, navegação,
referência visual) e o formato exato do código de identidade. Leia-o e siga-o; se ele
não existe, este guia basta — e avise o usuário que a régua específica é uma ponta
aberta.
`
