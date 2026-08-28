package main

// headerGuide é a régua UNIVERSAL e MANDATÓRIA do bloco de cabeçalho de arquivo — o
// padrão de comentário no topo de todo artefato do projeto, que carrega as marcações
// que o Anchors lê: identidade, carimbo de alteração, tags de agrupamento e opt-outs.
// É o guide transversal: rege TODOS os arquivos, não um tipo. O `anchors init` semeia
// um HEADER_GUIDE.md concreto no projeto a partir desta régua.
const headerGuide = `# Guia de cabeçalho (o bloco de marcações no topo de cada arquivo)

Todo arquivo que participa do grafo carrega, no topo, um BLOCO DE CABEÇALHO
padronizado — um comentário que o Anchors lê para saber a identidade do arquivo,
quando mudou, a que grupos pertence, e quais opt-outs ele declara. É MANDATÓRIO:
um arquivo sem cabeçalho conforme é invisível ao que o Anchors sabe fazer melhor.

Este guide é TRANSVERSAL — rege todos os arquivos, de qualquer camada. É a única
régua que não pergunta "que tipo de artefato?"; ela pergunta "todo arquivo tem sua
etiqueta?".

## A forma: chave-valor para dados, @tag para flags

O bloco é um comentário (no dialeto da linguagem do arquivo: '//' em TS/Go, '#' em
Python/shell, '<!-- -->' em markdown). Dentro dele:

- LINHAS 'chave: valor' — os DADOS (identidade, datas, agrupamento com valor).
- LINHAS '@flag' — as FLAGS booleanas e opt-outs (presença = ligado).

Exemplo (TS/Go) — uma tela:

  // @anchors
  //   ref: LGNN            (a TELA realiza a spec LGNN — referencia, não possui)
  //   updated_at: 2026-08-08
  //   layer: screen
  //   @feature: auth
  //   @noPropagation

Exemplo (markdown, na SPEC dona) — a spec POSSUI o código:

  <!-- @anchors
    code: LGNN            (a spec É a dona da identidade)
    updated_at: 2026-08-08
    layer: screen
    @feature: auth
  -->

O bloco abre com '@anchors' e as linhas seguintes são as marcações. O que não se
aplica, omita — mas todo arquivo do grafo precisa de IDENTIDADE: 'code:' se é o dono
(a spec), 'ref:' se referencia (o resto da trinca), OU 'layer:' se é de uma camada
RECONHECIDA sem spec (infra/dao/presentation/vocabulário de domínio — ver abaixo).

## As marcações

### Identidade: 'code:' (posse) vs 'ref:' (referência)

Esta é a distinção que evita a confusão mais comum. Um código de cenário pertence a
UMA unidade; os arquivos que giram em torno dela ou o POSSUEM ou o REFERENCIAM:

- 'code: <CÓDIGO>' — POSSE. O arquivo É o dono canônico da identidade. Só UM valor.
  Quem é o dono: a SPEC (a origem da verdade, SPEC.md). A spec DEFINE o requisito e
  seus cenários; ela é a dona do código. Gere com 'anchors code <nome>' (unicidade).
- 'ref: <CÓDIGO> [, <CÓDIGO>...]' — REFERÊNCIA. O arquivo NÃO é dono; ele realiza,
  cobre ou prova a(s) unidade(s) dona(s). Pode ser MÚLTIPLO. Assim:
    - o CÓDIGO (Divider.tsx) que realiza a spec  → 'ref: DIVI'
    - a FEATURE que cobre os cenários da spec     → 'ref: DIVI'
    - o TESTE que prova os cenários                → 'ref: DIVI'
  E um arquivo pode referenciar VÁRIAS unidades: um util testado por cenários de duas
  telas → 'ref: TXDT, MNDT'; uma tela que compõe componentes → 'ref' dos códigos deles.

Por que importa: pôr 'code:' num teste diria que o teste é DONO do código — mas ele só
o referencia (prova a unidade da spec). Trocar posse por referência cruza a
rastreabilidade e gera colisão de identidade falsa. Na dúvida: a spec tem 'code:'; o
resto da trinca tem 'ref:'.

### Camadas RECONHECIDAS: identidade por 'layer:' (e 'dep:' para dependências)

Nem toda camada tem spec. A Estrutura pode declarar camadas RECONHECIDAS — infra
(helpers puros), acesso a dado (dao), apresentação (mapa domínio→UI), vocabulário de
domínio (enums/catálogos). Elas existem para SAIR DO ESCRUTÍNIO de spec (não há regra
própria a documentar), mas ainda precisam de IDENTIDADE. Como não têm spec dona nem
irmã a referenciar, sua identidade mínima honesta é a própria camada:

- 'layer: <camada>' — para um arquivo de camada reconhecida, o 'layer:' É a identidade
  (satisfaz o gate). Não invente um 'code:' (não é dono de spec) nem um 'ref:' forçado
  (não realiza spec alguma). Declara honestamente "sou desta camada".
- 'dep: <arquivo> [, <arquivo>...]' — como não têm spec, não têm a Tabela de
  Dependências (SPEC_TYPES §5). Então declaram no PRÓPRIO header os ARQUIVOS de que
  dependem (o caminho do arquivo, não um código — o alvo pode ser outra reconhecida sem
  código). Cada um vira uma aresta 'depends-on'. É assim que uma reconhecida referencia
  outra, ou aponta a regida que consome. Ex.: um mapa de apresentação que tira cores de
  'theme/tokens.ts' → 'dep: theme/tokens.ts'.

Camadas REGIDAS (business-logic, validation, hook, store, repository, service, e a spec)
continuam exigindo 'code:'/'ref:' — elas têm regra a documentar, então têm spec e
identidade de código. 'layer:' sozinho NÃO basta para uma regida.

### Outros dados (chave: valor)

- 'updated_at: <AAAA-MM-DD>' — o CARIMBO DE ALTERAÇÃO: quando o conteúdo mudou pela
  última vez. NÃO o mantenha à mão (vai mentir) — o Anchors o preenche/valida contra
  o histórico real (git). Declará-lo aqui é opcional; o valor de verdade vem do grafo.

### Tags de agrupamento (chave: valor OU @tag)

Rótulos transversais que categorizam o arquivo, para consultas e escopo de gates:

- 'layer: <camada>' — a camada da Estrutura a que pertence (screen, component,
  service, model…). Normalmente o Anchors DEDUZ isso do caminho (o pattern da layer);
  declare só se quiser sobrepor.
- '@feature: <nome>' — o módulo/feature vertical (auth, dashboard, budgets…). Agrupa
  arquivos da mesma fatia de domínio, mesmo espalhados por camadas.
- '@<tag>' — qualquer outra tag de agrupamento livre (ex.: @experimental, @legacy).

### Opt-outs (só @flag — presença = ligado)

Dispensas honestas de uma regra do Anchors, registradas no próprio arquivo (datadas
pelo git, localizadas). O Anchors REALMENTE lê estas hoje:

- '@noPropagation' — este filho NÃO depende do pai: a onda de propagação não passa
  por ele. Ver PROPAGATION.
- '@anchors-shared-code' — os códigos de cenário aqui pertencem a OUTRA unidade de
  propósito (um teste de util/handler que prova o cenário da unidade que serve); não
  contam como colisão de identidade. Ver TRACEABILITY.

### Opt-outs COM RAZÃO ('@flag: <razão>')

Estes não são booleanos: exigem a razão escrita depois dos dois-pontos, na MESMA linha.
Um marcador nu ('@no-scenario:' sem texto) NÃO dispensa nada — continua reprovando. A
diferença existe porque estes dispensam uma regra sobre uma linha específica, e sem o
porquê a dispensa vira um buraco com nome bonito: quem ler depois não sabe se foi decisão
ou esquecimento.

Vão na linha do que dispensam (ou no comentário imediatamente acima, onde a linha não
comporta comentário legível):

- '@no-scenario: <razão>' — na linha de um REQUISITO da spec: este requisito não terá
  cenário na feature. Para o que é verdadeiramente não-observável por cenário (uma
  restrição estrutural provada por outro instrumento). Gate 'spec-feature-match'.
- '@no-paginate: <razão>' — numa função que promete um conjunto e não pagina: o limite é
  deliberado e conhecido (ex.: tabela com dezenas de linhas fixadas em migration). Gate
  'pagination-honored'.
- '@allow-boundary: <razão>' — numa linha que cruza uma fronteira de camada declarada em
  'boundaries:'. Para dívida reconhecida, que fica visível e datada no código em vez de
  numa lista de exceções distante. Gate 'layer-boundary'.

Um opt-out DISPENSA a regra, nunca o registro — é a porta legítima, o oposto do
buraco de cobertura silencioso.

## Regras do cabeçalho

- SEMPRE no topo do arquivo (antes do código), para ser o primeiro que se lê.
- code é a única marcação essencial; o resto é conforme a necessidade.
- updated_at NÃO se mantém à mão — deixe o Anchors/git cuidar; escrevê-lo errado é
  pior que omiti-lo (âncora mentindo sobre si).
- Opt-out sempre com um PORQUÊ ao lado (um comentário na mesma linha ou abaixo): quem
  ler depois precisa saber por que a regra foi dispensada.
- O dialeto do comentário é o do arquivo — o Anchors só lê TEXTO e reconhece os
  marcadores em qualquer comentário; não invente sintaxe fora do comentário.

## Anti-padrões (recuse-os)

- Arquivo sem code → órfão invisível à Rastreabilidade; dê identidade.
- updated_at mantido à mão → vai desatualizar e mentir; deixe o git ser a verdade.
- Opt-out sem porquê → ninguém sabe se ainda vale; explique ao lado.
- Marcação fora de comentário (no código executável) → o Anchors lê comentário; e
  poluiria o código.

## O projeto especializa

O 'anchors init' semeia um HEADER_GUIDE.md concreto para o seu projeto — com o
dialeto de comentário da sua stack e as tags/features reais. Este é o piso universal;
o do projeto adapta. Se o seu projeto não tem o guide de header, gere-o pelo init ou
copie este padrão — e avise que a padronização é uma ponta aberta.
`
