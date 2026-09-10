package main

// projectGuide é a régua da fase DESCOBRIR — a passada que faltava: o projeto que
// ainda NÃO existe. O `anchors init` resolve a Estrutura de um projeto que já tem
// arquivos no disco (ele INFERE do que está lá); num diretório vazio ele não tem o
// que inferir, e o agente parte para o plano sem saber em que linguagem escrever.
//
// O conceito vem do Polvo (a versão IDE do Anchors), que conduz o usuário por etapas
// de perguntas e destila um resumo no fim. Aqui ele é PORTADO, não copiado: o Polvo
// tem UI e modelo por trás (nove agentes devolvendo JSON que a interface renderiza em
// formulário); o Anchors não embute IA (ver DECISIONS) — quem tem o modelo é o agente
// que chama o CLI. Então a etapa não é código Go que pergunta: é esta régua, que o
// agente lê e executa na conversa que ele já está tendo com o usuário.
//
// A divisão em DOIS arquivos é o que separa dois públicos: o PROJECT.md é lido a cada
// vez que alguém vai escrever código (curto, técnico, decidido) e o INSIGHTS.md é lido
// quando alguém pergunta "por que isto?" (longo, com as perguntas, as respostas e as
// alternativas descartadas). Misturar os dois produz um documento que ninguém relê.
const projectGuide = `# Guia de projeto (a régua da fase DESCOBRIR)

Esta é a PRIMEIRA passada, e ela só existe quando o projeto AINDA NÃO EXISTE — um
diretório vazio, ou quase. Antes de qualquer plano, spec ou linha de código, o
projeto precisa responder o que ele é tecnicamente: em que linguagem, sob que
paradigma, com que estrutura, com que régua de formatação.

Por que antes do plano: o ` + "`anchors init`" + ` INFERE a Estrutura do que está no disco
(as extensões, os diretórios de código, a co-location). Num projeto vazio não há o
que inferir — ele pergunta "quais diretórios de código tratar como camadas?" e a
resposta honesta é "nenhum ainda". Sem esta passada, o agente escreve a primeira spec
sem saber se o projeto é Go ou TypeScript, e a decisão acaba tomada por acidente no
primeiro arquivo que alguém criar.

## O que esta fase produz

Dois arquivos na RAIZ do projeto, e a divisão entre eles é deliberada:

- ` + "`PROJECT.md`" + ` — o RESUMO TÉCNICO. Curto, decidido, sem alternativa em aberto.
  É o que se lê antes de escrever cada arquivo, então cada linha nele é uma regra que
  o código tem de obedecer. Não carrega justificativa.
- ` + "`INSIGHTS.md`" + ` — a TRANSCRIÇÃO. Cada pergunta que você fez, a resposta que o
  usuário deu, e o que foi descartado no caminho. É o que responde "por que Postgres e
  não Mongo?" seis meses depois, quando ninguém lembra da conversa.

A regra que mantém os dois úteis: **decisão no PROJECT, motivo no INSIGHTS.** Se você
está escrevendo "porque" no PROJECT.md, o texto pertence ao INSIGHTS.md. Se você está
escrevendo uma regra nova no INSIGHTS.md que não está no PROJECT.md, ela vai ser
ignorada — ninguém relê a transcrição para escrever código.

## Quem conduz: VOCÊ, na conversa

O Anchors não pergunta nada aqui — ele não embute modelo. **Quem entrevista é você**,
o agente que está operando o Anchors, na mesma conversa em que o usuário está agora.
O CLI só te deu esta régua.

Isso tem três consequências práticas:

- **A entrevista roda na CONVERSA, nunca num worker de background.** É o usuário quem
  responde; delegar a um subagente que não fala com ele não produz resposta nenhuma.
  Esta é a exceção à regra de não monopolizar a conversa — aqui a conversa É o
  trabalho.
- **Uma etapa por vez, esperando a resposta.** Faça a pergunta, PARE, leia o que o
  usuário respondeu, e só então formule a próxima. Não escreva as cinco etapas numa
  resposta só e peça para ele responder tudo junto: metade do valor está em cada
  pergunta nascer estreitada pela resposta anterior.
- **Os arquivos só são escritos NO FIM**, depois da revisão de inconsistências.
  Durante a entrevista você carrega as conclusões na conversa, não no disco. Escrever
  o PROJECT.md na etapa 2 e ir corrigindo produz um arquivo que já foi lido errado.

## Como conduzir (o formato de cada pergunta)

Você é um arquiteto sênior conversando, não um formulário. As regras que fazem a
diferença entre uma entrevista útil e um questionário que o usuário abandona:

**UMA pergunta por vez.** Nunca despeje as cinco etapas de uma vez. O usuário responde
a primeira, você entende, e a segunda já nasce mais estreita por causa da resposta.

**Toda pergunta tem TRÊS partes:**

1. **O QUE ESTÁ EM JOGO** (2–3 frases). A consequência CONCRETA da resposta na
   arquitetura — não sobre o processo, sobre o sistema. Nomeie coisas: qual banco,
   qual modelo de deploy, que padrão de teste, que teto de escala. Nunca escreva
   "vou te fazer algumas perguntas", "nesta etapa", "isso é importante" — o contexto
   é da arquitetura, não da entrevista.
2. **A PERGUNTA** (1 frase). Direta, sem preâmbulo.
3. **EXEMPLOS CONCRETOS** (3–4 linhas). Comece com "Exemplos para te orientar:" e dê
   um cenário realista por linha, no nível de detalhe que você espera na resposta.
   Exemplo vago produz resposta vaga.

**SEJA OPINATIVO.** Esta é a regra que mais muda o resultado. O usuário raramente é
especialista em tudo que você vai perguntar — se você apresentar um menu neutro de
oito opções, ele escolhe a que reconhece, não a que serve. Use o que você JÁ SABE das
etapas anteriores para eliminar o que não cabe e apresentar no máximo 2–3 opções
viáveis, **liderando com sua recomendação e o motivo dela**. "Dado que é um time de
1 pessoa com prazo de 2 meses, recomendo monólito modular porque X — funciona, ou
você tem um motivo forte para outra coisa?" é uma pergunta melhor do que "monólito,
microsserviços ou serverless?".

**Se o usuário pedir ajuda ou disser que não sabe:** apresente os tradeoffs concretos
no contexto DELE, e faça uma pergunta mais estreita que o leve a decidir. Não pule
para a etapa seguinte com o assunto em aberto.

## As etapas

Cinco etapas, nesta ordem. A ordem importa: cada uma restringe as opções da seguinte,
e inverter produz escolha de ferramenta antes de saber para quê.

### Etapa 1 — Propósito e forma
O que o sistema é, antes de qualquer tecnologia. É API, CLI, app web, mobile,
biblioteca, worker, monorepo com vários destes? Quem consome? Isso decide se existe
camada de interface, se há build de distribuição, se há estado de sessão.
Conclua com: o que é, quem consome, quais artefatos executáveis nascem.

### Etapa 2 — Linguagem e runtime
Só agora, e restrito pela etapa 1. Linguagem, versão, gerenciador de pacotes,
runtime de execução. Seja opinativo com base na forma e na experiência do time —
uma CLI distribuída como binário único empurra para Go/Rust; um app web com um time
que já sabe TypeScript raramente justifica trocar.
Conclua com: linguagem + versão, gerenciador, runtime.

### Etapa 3 — Arquitetura e paradigma
O padrão estrutural e o paradigma: layered, clean/hexagonal, feature-sliced, MVC?
Orientado a objetos, funcional, procedural? Módulos verticais por feature ou camadas
horizontais? Onde mora o domínio?
Isto é o que vira as ` + "`layers:`" + ` do anchors.yaml — a resposta aqui não é acadêmica,
ela decide o glob de cada camada. Cheque os presets: ` + "`anchors init`" + ` oferece
estruturas consagradas por stack, e casar com um preset poupa trabalho e erro.
Conclua com: padrão, paradigma, organização (modular ou em camadas), fronteiras.

### Etapa 4 — Estrutura macro e convenções de arquivo
Os diretórios de primeiro nível e o que mora em cada um. As extensões de cada tipo de
arquivo. A convenção de nome de teste (` + "`*_test.go`" + `, ` + "`*.test.ts`" + `, ` + "`*.spec.ts`" + `) e se
teste/spec ficam AO LADO do código (co-location) ou em árvore separada — o ` + "`init`" + `
pergunta exatamente isso, e a resposta aqui é a que ele vai usar.
Conclua com: árvore de diretórios, extensões, padrão de teste, co-location sim/não.

### Etapa 5 — Ferramental e formatação
A régua mecânica: indentação (tabs ou espaços, quantos), largura de linha, formatador
(gofmt, prettier, black, rustfmt), linter e sua configuração, convenção de nomes
(camelCase, snake_case, PascalCase por tipo de símbolo). Editores e extensões que o
time usa, e se há ` + "`.editorconfig`" + `, ` + "`.vscode/`" + ` ou equivalente a versionar.
Isso parece detalhe e não é: é o que faz o código que VOCÊ escreve ser indistinguível
do que o time escreve. Sem isto declarado, cada arquivo novo negocia o estilo de novo.
Conclua com: indentação, formatador, linter, convenção de nomes, editores/extensões.

### Fechamento — a revisão de inconsistências
Antes de escrever os arquivos, releia as cinco conclusões procurando CONTRADIÇÃO.
Elas são comuns e caras:
  • paradigma funcional declarado + estrutura que só faz sentido com classes
  • co-location "sim" + estrutura macro com ` + "`tests/`" + ` separado no topo
  • linguagem sem tipagem estática + gate que exige tipo
  • padrão modular por feature + diretórios organizados por camada técnica
Apresente cada contradição achada com as DUAS escolhas que colidem, citadas
textualmente, e peça ao usuário para resolver. Uma por vez. Só escreva os arquivos
quando não restar nenhuma.

## Escrevendo o PROJECT.md

Na raiz. Só o que foi DECIDIDO — sem "provavelmente", sem "a definir". Um campo que
não foi decidido não entra: campo vazio num documento que se lê antes de escrever
código é pior do que ausência, porque parece decisão.

` + "```md" + `
# Projeto: <nome>

<uma frase: o que o sistema é e para quem>

## Stack
| item | escolha |
| --- | --- |
| linguagem | Go 1.23 |
| gerenciador | go mod |
| runtime | binário nativo |
| frameworks | cobra (CLI), huh (TUI) |

## Arquitetura
- **Padrão:** layered
- **Paradigma:** procedural com tipos; sem herança
- **Organização:** camadas horizontais (cmd/ → internal/)
- **Fronteiras:** cmd/ só orquestra I/O; a decisão vive em internal/

## Estrutura macro
` + "```" + `
cmd/anchors/     comandos (uma casca fina por comando)
internal/        a lógica, um pacote por domínio
guides/          as réguas dos artefatos
` + "```" + `

## Convenções de arquivo
| item | escolha |
| --- | --- |
| extensões | .go |
| teste | ` + "`*_test.go`" + `, ao lado do código |
| co-location | sim |

## Formatação
| item | escolha |
| --- | --- |
| indentação | tabs |
| formatador | gofmt |
| linter | go vet |
| nomes | PascalCase exportado, camelCase interno |

## Ferramental
- editores: VS Code, Neovim
- extensões: gopls
- versionado: .editorconfig
` + "```" + `

## Escrevendo o INSIGHTS.md

Também na raiz. Uma seção por etapa, e dentro dela cada pergunta com a resposta que
o usuário deu. O que faz este arquivo valer o disco é a terceira parte: **o que foi
descartado e por quê.** Uma decisão sem as alternativas rejeitadas não se revisa —
quem a reabre não sabe se a opção que ele está propondo já foi considerada.

` + "```md" + `
# Insights do projeto

> Transcrição da fase DESCOBRIR. As DECISÕES vivem no PROJECT.md; aqui ficam as
> perguntas, as respostas e o que foi descartado.

## Etapa 2 — Linguagem e runtime

**P:** <a pergunta como você a fez>
**R:** <a resposta do usuário, na palavra dele>

**Decidido:** Go 1.23
**Descartado:**
- Rust — curva de aprendizado não se paga num time que já entrega em Go
- Node/TS — distribuição como binário único era requisito da etapa 1
` + "```" + `

## Depois de escrever

1. ` + "`anchors init`" + ` — agora ele tem o que inferir. As respostas das etapas 3, 4 e 5
   são exatamente o que ele pergunta (preset, camadas, co-location, padrão de teste);
   responda com o que está no PROJECT.md, sem renegociar.
2. ` + "`anchors map build`" + ` — os dois arquivos novos entram no mapa.
3. Siga para o PLANEJAR (` + "`anchors guide plan`" + `): o primeiro plano decide quais specs
   nascem, e agora ele nasce sabendo em que linguagem elas serão implementadas.

## Quando o projeto JÁ existe

Não conduza a entrevista inteira — o disco já respondeu a maior parte. Leia o que
está lá (extensões, diretórios, ` + "`.editorconfig`" + `, config do linter, manifesto de
pacotes), escreva o PROJECT.md com o que você OBSERVOU, e pergunte ao usuário só o
que o código não revela: o paradigma pretendido, as fronteiras que deveriam existir,
e o que no código de hoje é dívida em vez de padrão. O INSIGHTS.md registra o que foi
observado versus o que foi decidido — a diferença entre os dois é a lista de dívidas.

## O que NUNCA fazer

- **Escrever o PROJECT.md sem conversar.** Um resumo técnico que você inventou é um
  palpite com autoridade de documento. Se não houve entrevista, não há PROJECT.md.
- **Deixar campo em aberto.** "A definir" no PROJECT.md vira decisão tomada por
  acidente no primeiro arquivo que alguém escrever.
- **Repetir o INSIGHTS dentro do PROJECT.** O PROJECT.md deixa de ser lido no dia em
  que fica longo, e ele é o arquivo que precisa ser lido sempre.
- **Deixar os dois envelhecerem em silêncio.** Quando uma decisão muda (troca de
  linter, novo módulo, mudança de paradigma), atualize o PROJECT.md e registre no
  INSIGHTS.md o que mudou e por quê. Um PROJECT.md que mente sobre o projeto é pior
  do que não existir — o agente o obedece.
`
