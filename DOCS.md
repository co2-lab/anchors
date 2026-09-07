# Anchors — A Documentação

> A trinca cobre o que está **dentro** de uma unidade. Este documento é sobre o que fica
> **fora** dela: os artefatos agregados que várias unidades alimentam e nenhuma possui.

## 1. Por que a documentação precisa de mecanismo

O Anchors sempre cobrou spec, feature, teste e código — as quatro peças de uma unidade. A
documentação ficava como uma frase no card do agente: *"a documentação evolui junto"*.

Uma frase assim não diz **qual** documentação. Um agente que lê só isso atualiza o que lhe
parece documentação — normalmente um `README` — e o contrato da API segue sem o endpoint
que ele acabou de escrever.

E há um problema anterior, de forma. Uma documentação útil mostra o **conteúdo**. Se ela só
aponta para os arquivos, deixa de ser documentação e vira indexação: quem lê abre link por
link, e desiste. Mas o conteúdo já existe nas specs, e copiá-lo para dentro de um `.md`
escrito à mão cria duplicação que desatualiza em silêncio.

Markdown não resolve isso sozinho. Não há inclusão nativa no CommonMark nem no GFM:
`{% include %}` é Jekyll e só funciona no Pages; `--8<--` é MkDocs; `<iframe>` é sanitizado
pelo GitHub na visualização do repositório. Sem build, a única saída seria duplicar à mão.

## 2. A saída é o build: três camadas, o conteúdo em uma

```
*.spec.md        a fonte — o conteúdo mora aqui, e só aqui
doct/*.md.tmpl   o template — a moldura, que REFERENCIA trechos das specs
docs/*.md        o compilado — conteúdo real, gerado, ninguém edita
```

O agente escreve a **moldura**. O `anchors docs build` produz o documento final, e o
pipeline o roda no merge. A duplicação existe só no compilado — e ali ela não dói, porque
ele é gerado e o gate `docs-fresh` acusa quando envelhece. **A duplicação que dói é a
escrita à mão, que envelhece sem ninguém ver.**

A linguagem é o `text/template` da stdlib. Uma dependência nova só para documentação seria
custo permanente, e o Anchors já usa a mesma linguagem no `derived.files` (`{{dir}}`,
`{{name}}`).

### 2.1 A pasta é `doct`, não `doc[t]`

Colchete no nome de pasta é classe de caracteres em glob de shell e em `.gitignore` —
`doc[t]/` casaria `doct/`, e todo padrão que mencionasse a pasta precisaria escapar.
Alguém esqueceria. Numa URL exige encoding (`doc%5Bt%5D`), que quebra link no GitHub.

### 2.2 O compilador não apaga trabalho de ninguém

Todo arquivo gerado abre com um marcador:

```html
<!-- anchors:generated de doct/regras.md.tmpl — NÃO EDITE: rode `anchors docs build` -->
```

Ele serve a duas coisas, e a segunda é a que protege: quem abre o arquivo sabe que editá-lo
é perder o trabalho; **e o compilador recusa sobrescrever um `.md` que não o tenha.** Uma
doc de produto escrita à mão não é destruída porque alguém criou um template com o mesmo
nome. O comando reporta o arquivo como não gerado, em vez de apagá-lo.

## 3. A matriz: duas visões dos mesmos dados

Uma spec pertence a uma **camada** e tem **seções**. As duas perguntas são legítimas e
nenhuma hierarquia serve às duas:

| a pergunta | o corte | onde |
|---|---|---|
| "todas as regras do sistema" | por **seção**, atravessando camadas | `docs/regras.md` |
| "o que existe em `screens`" | por **camada** | `docs/camadas/screen.md` |

Uma hierarquia só serviria a uma delas, e a outra teria de navegar unidade por unidade —
que é o modo de falha do §1.

## 4. A estrutura organizacional

```
docs/
  produto.md          O QUE o sistema faz e para quem — vem dos planos, escrito à mão
  arquitetura.md      COMO ele é montado — C4 em quatro níveis, com Mermaid
  comportamento.md    ÍNDICE dos cenários → a página da camada
  regras.md           ÍNDICE das regras   → a página da camada
  camadas/*.md        O CONTEÚDO — uma página por camada
  contratos/*.md      as docs específicas do tipo de projeto (OpenAPI, esquema…)
```

### O conteúdo mora num lugar só; os cortes transversais são índices

A primeira versão repetia o texto nas duas visões da matriz, e a medida mostrou por que não
serve: `regras.md` saiu com **9.608 linhas** — o corte transversal de 84 unidades é o
documento inteiro, mais uma vez. E cresce com o projeto, sem limite.

Um índice não tem esse problema: uma linha por regra, e o link leva ao texto. O preço é o
clique, e é justo — quem abre "todas as regras" procura *uma*, e antes rolava por todas.

**A âncora do link é gerada, nunca escrita à mão.** Um índice cujo link não resolve é pior
que não ter índice: ele é clicável, e o navegador fica onde está.

### O formato se ajusta ao tamanho — decidido em um lugar

Uma camada com três unidades cabe numa página com tudo; com trinta, é um documento que
ninguém rola até o fim. O corte é do projeto:

```
anchors docs build --max-units 30 --max-lines 3000
```

Acima do corte, a página traz o **resumo** de cada unidade e diz isso no topo. Os dois
critérios existem porque contar unidades engana: cinco specs longas geram mais página que
vinte curtas.

**O que muda é o que cabe na página, nunca em quantos arquivos a camada se parte.** Dividir
ao cruzar um limiar quebraria todo link externo no dia em que a unidade seguinte entrasse.

E a decisão é tomada **uma vez**, no `docs build`, e distribuída a todas as páginas. Deixar
cada template escolher o seu limiar parecia flexível e foi o defeito: a página resumia por
um critério, o link era montado por outro, e **483 de 812 links saíram quebrados** — todos
clicáveis, todos parando no mesmo lugar.

A ordem é de fora para dentro — produto, arquitetura, comportamento, regras, camadas. É a
ordem em que alguém que chega precisa delas, e não a ordem em que o time as escreveu.

`anchors docs init` escreve esse esqueleto, inclusive uma página por camada existente. Ele
**propõe, não impõe**: os templates são editados pelo time. O que o Anchors garante é que
ninguém precise inventar a organização do zero — inventar a cada projeto é como se chega a
uma documentação onde cada página segue uma lógica diferente, e quem lê tem de descobrir a
lógica antes de achar o que procura.

Rodar de novo **não sobrescreve** sem `--force`: o template já editado carrega a moldura
que o time escreveu, que é justamente a parte não gerada.

### 4.1 Os cenários entram na documentação

A feature é o único artefato da trinca escrito para ser lido por quem não programa —
Gherkin existe para isso. E é ela que responde o que a spec não responde: a spec diz a
**regra** ("dívida aberta bloqueia"), o cenário diz o que **acontece** ("dado uma dívida
aberta, quando confiro a régua, então não libera").

Deixá-lo fora da documentação manteria o comportamento observável do sistema num `.feature`
que só o time de desenvolvimento abre, enquanto a documentação — que é para os outros —
descreve regras em abstrato.

A feature é encontrada pela **aresta do mapa**, não por convenção de nome: um projeto que
organize os arquivos de outro jeito continua funcionando.

## 5. As documentações específicas do tipo de projeto

O que a trinca não alcança são os artefatos **agregados**. Uma API tem um contrato que vive
fora do código — o OpenAPI — e um endpoint que não entra nele é invisível para quem
consome. Um projeto com banco tem o esquema. Um design system tem o catálogo.

Declaram-se no `anchors.yaml`:

```yaml
docs:
    required:
        - kind: openapi
          path: docs/contratos/openapi.yaml
          trigger: [lambdas]        # alterar esta camada OBRIGA tocar esta doc
          why: >-
              é o contrato que diz ao consumidor o que esperar quando a rota falha.
        - kind: c4
          path: docs/arquitetura.md # sem trigger: muda com a ESTRUTURA, não com a unidade
```

**A declaração é o ato.** Um projeto que não declara `docs:` não é cobrado — cobrar OpenAPI
de quem não tem API seria ruído, e ruído no card é o que faz o agente parar de ler o card.

O `trigger` é o que transforma "documente" em algo verificável: mexer numa lambda de rota
muda o contrato da API; mexer num utilitário interno não muda contrato nenhum, e cobrar ali
ensinaria o agente a atualizar o arquivo sem pensar — que é pior que não atualizar.

### 5.1 Nomear o arquivo não basta

"Atualize `docs/openapi.yaml`" leva um agente a acrescentar o endpoint que acabou de
escrever e parar ali — sem os erros, sem os exemplos, sem o esquema do corpo. O resultado
passa em qualquer verificação de existência e é inútil: **um contrato que descreve só o
caminho feliz não é contrato.**

Por isso cada `kind` conhecido carrega o que a documentação daquele tipo precisa
**responder**, e a **armadilha** — o jeito de "cumprir" a tarefa produzindo algo que não
serve. `anchors docs duties [--layer <camada>]` mostra os dois.

Tipos que o Anchors sabe instruir: `openapi`, `c4`, `schema`, `components`, `adr`,
`runbook`. Um tipo desconhecido **não é erro** — o projeto pode ter uma documentação que o
Anchors não conhece, e recusá-la faria o mecanismo servir só ao que já foi previsto. O que
ele não sabe instruir, ele ao menos nomeia e cobra.

### 5.2 A arquitetura em C4

Todo projeto tem arquitetura, e o Anchors já a regula — as camadas, as dependências entre
planos, as fronteiras que o `layer-boundary` cobra. Mas essa regulação vive no
`anchors.yaml` e no mapa, **em formato de máquina**. Quem chega precisa da mesma informação
em formato de gente.

**Cada nível amplia uma caixa do anterior** — essa é a regra central do modelo, e o que ele
existe para impor:

| nível | o que amplia | quantos diagramas |
|---|---|---|
| 1 — Contexto | o sistema como caixa única | um |
| 2 — Contêineres | a caixa "o sistema" | um |
| 3 — Componentes | **um contêiner** | **um por contêiner interno** |
| 4 — Código | um componente | só onde houver algo não óbvio |

Um único diagrama de nível 3 misturando o app, a API e a infraestrutura **não é nível 3 de
coisa nenhuma** — é exatamente a falha que o C4 existe para evitar: um desenho só, com tudo.

### 5.2.1 Contêiner é o que executa **ou armazena** dado

Não é sinônimo de "processo que escrevemos". Banco de dados, fila, cache e sistema de
arquivos são contêineres, e a omissão do banco é o erro mais comum ao aplicar o modelo.

O que roda separado se declara na Estrutura, e é o que dá os níveis 2 e 3:

```yaml
containers:
    - name: app
      description: a interface
      layers: [screen, component, feature-hook]
      talks:
          - to: api
            protocol: HTTPS/JSON
    - name: api
      description: as rotas que o app consulta
      layers: [lambdas, shared]
    - name: banco
      description: o que persiste
      external: true      # é contêiner, e não é nosso por dentro
```

**`external: true`** diz "existe, conversa conosco, e não temos componentes lá dentro". Ele
aparece no nível 2 — a conversa é real e o protocolo importa — e **não ganha nível 3**:
desenhar componentes dentro de um banco de terceiro afirmaria um conhecimento que não temos.

**O protocolo é obrigatório em cada conversa.** Uma seta sem ele diz que os dois se falam e
não diz o que acontece quando a conversa falha — que é a única coisa que um diagrama de
contêineres tem a dizer sobre risco.

### 5.2.2 Camada não é componente

`screen`, `lambdas`, `shared` são agrupamentos de código no repositório. O componente do C4
é a peça com responsabilidade **dentro de um contêiner**. O `containers.layers` é a ponte
entre os dois vocabulários: ele diz que camadas rodam em que contêiner, e é o que permite o
nível 3 mostrar as unidades daquele contêiner em vez da lista de camadas do repositório.

Uma camada que nenhum contêiner declara é **dita** na página, não escondida: um diagrama que
a omite em silêncio afirma, por ausência, que ela não existe.

### 5.2.3 Os diagramas são Mermaid

O GitHub os renderiza nativamente, e o MkDocs e o Starlight também. Um PNG exportado de uma
ferramenta de desenho ficaria fora do controle de versão útil — o diff não diz o que mudou,
e o arquivo-fonte do desenho acaba noutro lugar, ou some. Aqui o diagrama *é* texto,
versionado com o resto.

Os níveis 1 e 2 são escritos à mão no template: descrevem o sistema inteiro, e nenhuma spec
sozinha os conhece. O nível 3 vem dos contêineres declarados.

O C4 **não tem gatilho por camada**, e não é esquecimento: ele não muda quando uma unidade
muda — muda quando a estrutura muda (um contêiner novo, uma fonte externa nova, uma
fronteira que se desloca). Por isso não é trabalho de card de trinca, e sim de issue própria.

## 6. O gate `docs-fresh` é informativo

O compilado envelhece em silêncio: alguém altera uma regra na spec, esquece de recompilar, e
a documentação segue afirmando a versão **antiga** — com conteúdo real, e por isso
convincente. É pior que a página vazia: uma doc obviamente incompleta manda procurar a
fonte; uma doc desatualizada não manda procurar nada.

O gate é **informativo por decisão**, e o critério é geral: *um desvio que um comando
conserta sozinho não deve reprovar o autor.* O `anchors docs build` resolve, o pipeline o
roda no merge, e o gate serve para o autor ver antes. Reprovar aqui gastaria a atenção da
revisão com trabalho de máquina.

O gate **não escreve**. Se ele consertasse o que aponta, a segunda execução sempre passaria
e o defeito só apareceria em quem clonasse o repositório.

## 6.1 A spec tem de bastar por si — `doc-self-contained`

O corpo da spec vira documentação palavra por palavra, e isso muda como ela deve ser
escrita. Uma frase que só **aponta** para um plano manda o leitor da doc a um arquivo que
ele não tem aberto — que é exatamente o que este mecanismo existe para eliminar.

O gate não é sobre citar menos, é sobre **trazer o texto**:

| | |
|---|---|
| acusado | *"O plano `PSHUX` explica por que o canal é assim."* |
| certo | *"O desenho evita expor rota de ingestão: o alarme publica no SNS nativamente."* |

Não é acusado nada que **acompanhe** a referência: a citação, uma tabela ou lista logo
abaixo, ou — só para revisões — a explicação em prosa. Uma revisão é uma etiqueta e a frase
ao lado é o conteúdo (*"a `PLTFR-R0003` corrigiu o escopo: o contador saiu porque a fonte
não o expõe por réplica"*). Um **caminho de arquivo** não ganha esse escape: ele não é
etiqueta de nada, é o lugar aonde a pessoa teria de ir.

### Sem match de idioma

O gate não procura "plano", "ver" nem "conforme". **O Anchors governa projetos em qualquer
língua**, e um gate que casa vocabulário passa em silêncio no projeto escrito na outra — o
que é pior que não existir, porque a spec *parece* protegida.

Ele casa estrutura, e só o que o próprio Anchors define: os **caminhos que o mapa conhece**
(nenhuma língua muda um caminho, e um padrão `plans/*.md` escrito à mão só valeria para
quem chama a pasta assim) e a forma `{CODIGO}-R000N` da doutrina. As aspas reconhecidas são
as de qualquer tradição escrita — `« »`, `„ “`, `「 」` — porque exigir aspas latinas
acusaria injustamente quem escreve em francês, alemão ou japonês.

## 7. O site

O `docs/*.md` compilado é markdown comum, e qualquer gerador estático o consome sem
conversão — Astro Starlight, MkDocs, Docusaurus. O site é do **projeto governado**, não do
Anchors: o Anchors produz o `.md`, e o projeto decide se e como o expõe.

## 8. Resumo

| o problema | o mecanismo |
|---|---|
| doc que só aponta vira indexação | o compilador inclui o **conteúdo** da spec |
| duplicação que desatualiza em silêncio | duplicação só no **gerado**, mais o `docs-fresh` |
| markdown não tem include | build: `doct/*.tmpl` → `docs/*.md` |
| a matriz camada × seção | duas visões geradas dos mesmos dados |
| "documente" não diz o quê | `docs.required` com `trigger` por camada |
| nomear o arquivo não basta | cada `kind` diz o que pede e qual a armadilha |
| a arquitetura só existe em formato de máquina | C4 em Mermaid, um nível 3 por contêiner |
| a spec manda o leitor da doc para fora | `doc-self-contained`, sem match de idioma |
| cada projeto inventa a organização | `anchors docs init` propõe o esqueleto |
