package main

// guideGuide é a régua de como escrever um GUIDE — o meta-guide. Um guide é a régua
// de um tipo de artefato; para ser CONFRONTÁVEL (e não só lido), ele precisa destilar
// suas regras em PONTOS DE CONFORMIDADE verificáveis. Sem eles, o gate de julgamento
// recai em heurística vaga ("respeita o espírito?"); com eles, a IA julga item a item.
const guideGuide = `# Guia de guide (a régua de como escrever uma régua)

Um guide é uma ÂNCORA que rege um tipo de artefato — ele diz como o artefato deve ser,
e é a régua contra a qual o artefato é confrontado. Um guide que ninguém confronta é
só um documento; um guide que rege mas não diz O QUE VERIFICAR só pode ser confrontado
por heurística vaga. Este guia existe para que seus guides sejam CONFRONTÁVEIS.

## A regra central: destile em PONTOS DE CONFORMIDADE

A prosa do guide explica e ensina. Mas o julgamento (o gate de review por IA) precisa
de ALVOS verificáveis, não de prosa. Por isso todo guide de governança tem uma seção
de pontos de conformidade: a lista das coisas que, se um artefato as viola, ele está
FORA da régua. Cada ponto é objetivo o bastante para um veredito pass/fail.

Isto torna o julgamento MENOS heurístico: em vez de "esta tela respeita atomic
design?" (vago, a IA adivinha o critério), vira "verifique CK1, CK2, CK3" (cada um um
veredito, com o critério fixado no guide). A checklist é a prosa destilada — não pode
divergir dela, porque nasce dela.

## A seção obrigatória

Todo guide que rege algo (tem regra 'governs') DEVE ter:

  ## Pontos de conformidade

  - CK1: <um critério verificável, afirmativo> — <como reconhecer a violação>
  - CK2: ...

Regras dos pontos:
- CADA ponto tem um código (CK1, CK2, …) — a identidade, para o laudo referenciar o
  item específico (como os cenários de uma spec).
- AFIRMATIVO e objetivo: "átomos vêm do design system", não "boa organização".
  Alguém deve conseguir dizer pass/fail olhando o alvo, sem adivinhar o que você quis.
- INDEPENDENTE: um ponto testa uma coisa. Se você escreveu "e" no meio, provavelmente
  são dois pontos.
- ANCORADO na prosa: cada ponto destila uma regra que o corpo do guide já explica. Se
  um ponto não tem explicação acima, ou falta a explicação, ou o ponto é invenção.
- HONESTO sobre o não-computável: um ponto pode ser subjetivo ("a separação de
  responsabilidades é clara?") — tudo bem, é justamente o que o julgamento de IA mede.
  Mas prefira o objetivo quando der; quanto mais objetivo o ponto, menos a IA adivinha.

## Pontos POR TARGET (quando o guide rege mais de uma coisa)

Um guide costuma reger MAIS DE UM tipo de alvo — o guide de frontend rege telas E
componentes; o de backend rege handlers E models. Um ponto que vale para telas
("átomos vêm do design system") não vale para um model. Então NÃO faça uma checklist
plana que mistura tudo: AGRUPE os pontos por target, com um subcabeçalho que nomeia a
camada/tag a que se aplicam. Os códigos ganham um prefixo do grupo:

  ## Pontos de conformidade

  ### Para telas (tag: screen)
  - SCR-CK1: a tela injeta componentes, não monta blocos de layout inline
  - SCR-CK2: átomos e moléculas vêm do design system, não são redefinidos aqui

  ### Para componentes (tag: component)
  - CMP-CK1: o componente não busca dados — recebe tudo por props
  - CMP-CK2: um átomo não compõe outros átomos de domínio

  ### Para todos (qualquer alvo regido)
  - GEN-CK1: nomes de arquivo condizem com o que exportam

Como o julgamento usa isto: ao confrontar um alvo, a IA olha a CAMADA/TAG do alvo
(o mapa a informa) e aplica SÓ os pontos do grupo correspondente + os do grupo "Para
todos". Um ponto de tela nunca é cobrado de um model. Assim cada alvo é medido pela
régua certa, e o laudo cita o código do ponto que corresponde (SCR-CK1, CMP-CK2…).

Por que isto importa — torna a análise FOCADA e OBJETIVA:
- FOCO: a IA verifica só os pontos da camada do alvo, sem pesar regras irrelevantes.
- OBJETIVIDADE: cada veredito aponta um CK nomeado, não uma impressão geral —
  reproduzível (o mesmo alvo contra os mesmos CKs converge entre execuções).
- ECONOMIA: a IA consulta o grupo relevante, não a prosa inteira do guide — o que
  ajuda na adoção em lotes (ver 'anchors guide', seção ADOÇÃO).

Se o guide rege um só tipo de alvo, a checklist plana basta — não invente grupos.

## Como o julgamento usa isto

Quando um gate 'measures: judgment' confronta um alvo contra o guide, a IA:
1. lê a seção Pontos de conformidade do guide;
2. avalia o alvo contra CADA ponto;
3. emite um laudo item a item — para cada CK: conforme, ou a violação (o quê, onde,
   por quê, como corrigir). Veja 'anchors guide' (seção JULGAR).

Um guide SEM pontos de conformidade força o julgamento a virar heurística — o
'anchors doctor' avisa isso como débito.

## Estrutura de um guide

1. Propósito — o que este guide rege e por quê (uma ou duas frases).
2. As regras, em prosa — o ensino: o que é certo, exemplos, o raciocínio. É onde
   alguém aprende a fazer, não só a ser verificado.
3. Pontos de conformidade — a seção obrigatória acima; a prosa destilada em alvos.
4. (Opcional) Anti-padrões — os erros comuns e o antídoto.

## O que este guide rege

Um guide precisa declarar (via regra 'governs' no anchors.yaml) QUAL tag ele rege —
senão ele não confronta ninguém e não é guide, é doc. Ao criar um guide, adicione a
regra governs correspondente. Se o guide é transversal (não rege um tipo de artefato),
ele provavelmente é um doc, não um guide.

## Anti-padrões (recuse-os)

- Guide sem pontos de conformidade → só dá para confrontar por heurística; destile.
- Ponto vago ("bem estruturado") → ninguém julga pass/fail; torne objetivo.
- Ponto que junta várias verificações com "e" → separe em pontos independentes.
- Checklist que contradiz a prosa → ela deve NASCER da prosa, não competir com ela.
- Guide que não rege ninguém → declare o governs, ou é doc.
`
