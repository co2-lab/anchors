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
  arquitetura.md      COMO ele é montado — C4, os quatro níveis
  comportamento.md    O QUE ACONTECE — os cenários das features, todos
  regras.md           AS REGRAS — todas as specs, por seção (visão-matriz A)
  camadas/*.md        POR ONDE — uma página por camada (visão-matriz B)
  contratos/*.md      as docs específicas do tipo de projeto (OpenAPI, esquema…)
```

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

Todo projeto tem arquitetura, e o Anchors já a regula — as camadas, o `needs:` entre
planos, as fronteiras que o `layer-boundary` cobra. Mas essa regulação vive no
`anchors.yaml` e no mapa, **em formato de máquina**. Quem chega precisa da mesma informação
em formato de gente.

O C4 é a notação porque separa em níveis (contexto → contêineres → componentes → código), e
é isso que impede a falha clássica do diagrama de arquitetura: um desenho só, com tudo,
ilegível. **Separar em níveis é o mecanismo, não a decoração** — e o C4 vale tanto pelo que
cada nível mostra quanto pelo que ele omite.

Os níveis 1 e 2 são escritos à mão no template: descrevem o sistema inteiro, e nenhuma spec
sozinha os conhece. Os níveis 3 e 4 vêm das specs e das camadas.

O C4 **não tem gatilho por camada**, e não é esquecimento: ele não muda quando uma unidade
muda — muda quando a estrutura muda (um contêiner novo, uma fonte externa nova, uma
fronteira que se desloca). Por isso ele não é trabalho de card de trinca, e sim de uma
issue própria.

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
| a arquitetura só existe em formato de máquina | C4, com issue própria |
| cada projeto inventa a organização | `anchors docs init` propõe o esqueleto |
