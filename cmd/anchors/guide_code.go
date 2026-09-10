package main

// codeGuide é a régua universal do artefato CÓDIGO: como implementar guiado pela
// spec, com arquitetura que não apodrece. Agnóstico de linguagem/framework — a
// doutrina vem da experiência real (organização por eixo, dependência unidirecional,
// promoção sob demanda). O projeto especializa no seu guide de arquitetura.
const codeGuide = `# Guia de código (implementar guiado pela spec)

O código implementa o comportamento que a spec descreveu. Ele é confrontado contra a
spec (a spec é a régua; o código, o regido). Aqui a doutrina não é sobre sintaxe —
é sobre ONDE cada coisa mora, para o projeto não virar um novelo com o tempo.

## Dois eixos de organização, combinados

Todo código se organiza em duas dimensões ao mesmo tempo:
- POR TIPO/COMPOSIÇÃO (horizontal) — o que é reutilizável e NÃO conhece domínio
  (blocos genéricos, utilidades).
- POR DOMÍNIO/FEATURE (vertical) — fatias isoladas por área de negócio, onde vive a
  lógica daquele domínio.
A lógica de negócio não tem lugar natural no eixo por-tipo — se você só organiza por
tipo, a lógica se espalha. Por isso: domínio organiza o que tem lógica; tipo organiza
o que é reutilizável sem domínio.

## A regra única de classificação

Reutilizável e sem domínio → vive na camada por-tipo (compartilhada).
Acoplado a um domínio → vive na fatia daquele domínio.
Teste decisivo: "este pedaço faria sentido em OUTRO app, sem este domínio?"
Sim → compartilhado. Não → feature.

## Dependência unidirecional

- Um módulo só importa de camadas ABAIXO dele. Nunca para os lados, nunca para cima.
- Uma feature NUNCA importa de outra feature.
- Infra/compartilhado NUNCA importam de uma feature. O compartilhado "sobe"; ele não
  conhece quem o consome.

## Promoção sob demanda (DRY na hora certa)

Algo nasce LOCAL a uma feature. Só SOBE para o compartilhado quando 2+ features o
consomem — não antes. Abstração prematura (subir "por via das dúvidas") é dívida, não
economia. Espere o segundo consumidor.

## Fronteiras de dados

- Só as BORDAS buscam dados (as entradas do fluxo). O resto RECEBE por parâmetro.
  Um bloco que recebe seus dados por parâmetro é testável isolado e reutilizável.
- Separe a LÓGICA PURA (sem framework) num lugar próprio. Funções puras do domínio,
  sem dependência de UI/IO, são as mais fáceis de testar — e as que mais se reusam.

## A trinca co-localizada

Todo artefato de código nasce com sua spec, sua feature e seu teste ao lado. Escrever
código sem os três é deixar a unidade sem régua, sem cobertura e sem prova. (É o
núcleo do Anchors — o mapa espera a trinca.)

## Violação real vs. exceção sancionada

Nem toda quebra de regra é bug. Antes de "corrigir" cegamente uma dependência que
parece torta, pergunte se é uma EXCEÇÃO sancionada (documentada, com motivo) ou uma
VIOLAÇÃO de verdade (sintoma de acoplamento errado). Corrija a violação; respeite a
exceção — e se ela não está documentada, documente o porquê.

## Convenções (em espírito, agnósticas)

- Nome do arquivo condiz com o que ele exporta.
- Evite re-exports/barris que criam ciclos de import.
- Zero valores mágicos — use constantes/tokens nomeados.
- Importe por um alias estável, não por caminhos relativos frágeis.

## Anti-padrões (recuse-os)

- Lógica de negócio na camada por-tipo → mova para a fatia do domínio.
- Feature importando de outra feature → promova o comum para o compartilhado.
- Abstrair no primeiro uso → espere o segundo consumidor.
- Bloco reutilizável que busca os próprios dados → passe os dados por parâmetro.
- Código sem a trinca (spec/feature/teste) → a unidade fica órfã no mapa.

## Especialização do projeto

Os nomes das camadas, a stack, os tokens e a estrutura de pastas concreta são do seu
projeto — veja o guide de arquitetura/código na camada 'guide' do anchors.yaml. Este
guia é a doutrina universal; siga o dialeto do projeto quando existir, e avise se não
existir.
`
