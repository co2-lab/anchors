# Falhas — declarar, tratar, registrar, e descobrir de onde vêm

> Documento de DESENHO, anterior à implementação.

## O que já existe, e o que nunca foi confrontado

A spec já cataloga como a unidade FALHA. A seção existe, tem letra própria e formato:

    ## Errors / Failures
    | Rule | Condition | Failure |
    | --- | --- | --- |
    | `CRED-E01` | saldo insuficiente | recusa com código 402 |

O propósito declarado: *"como a unidade falha e quem sinaliza. Use quando há entrada
inválida, dependência externa ou estado impossível de tratar."*

**E nenhum gate confronta isso.** Verificado: nenhum checker lê regras `-E`. Elas são
catalogadas, atravessam o pipeline inteiro, e nada pergunta se foram tratadas, se logam, ou
se acontecem.

## As três camadas, e a ordem importa

O pedido tem três partes, e misturá-las foi o erro da primeira versão deste documento.

### Camada 1 — ANTES de rodar: a falha está tratada e REGISTRADA?

A cadeia que precisa estar fechada antes de qualquer coisa:

    a falha está DECLARADA na spec  (`-E` catalogada)
    o código a TRATA                (há um caminho que lida com ela)
    o tratamento LOGA               ← a peça que garante tudo o mais

A terceira é a que sustenta as outras duas camadas, e é a que ninguém cobra: um tratamento
que engole a falha sem registrar nada é o silêncio perfeito — a falha acontece, nada sabe,
e a camada 2 nunca vai ver.

**TRATAR não é "ter um `catch`".** O `catch` é um exemplo, e cravá-lo seria cravar sintaxe
de uma família de linguagens — o oposto do que o Anchors faz. Tratamento é qualquer
caminho que lida com a falha em vez de deixá-la escapar:

    if (x == null) { return recusa() }      trata
    try { ... } catch (e) { log(e) }        trata
    match result { Err(e) => ... }          trata
    if err != nil { return fmt.Errorf(...) } trata

O que todos têm em comum não é a forma, é o EFEITO: a falha vira parte do fluxo, e a
aplicação segue. É isso que a torna resiliente.

Quem sabe reconhecer a forma no dialeto do projeto é o próprio projeto — o mesmo mecanismo
do `GuardPatterns`, que já existe para guarda de parâmetro e não crava sintaxe nenhuma.

| gate | pergunta |
| --- | --- |
| `failure-handled` | a `-E` declarada tem um caminho que a trata? |
| `failure-logged` | esse caminho REGISTRA a ocorrência? |
| `failure-declared` | o tratamento que existe no código corresponde a alguma `-E`? |

O terceiro é o inverso do primeiro, e pega o caso comum: alguém escreveu uma defesa e
nunca declarou o que ela previne.

**É aqui que o Anchors é forte**, porque é confronto estático: tem a spec, tem o código,
tem o dialeto. Não depende de produção nem de log nenhum.

### Camada 2 — RODANDO: a falha aconteceu

Garantida a camada 1, há log quando a falha ocorre. O Anchors não lê esse log — ele
INGERE ocorrências já extraídas dele, como já faz com teste:

> *"Anchors does NOT run the test — you run it and hand over the report."*

O formato do LOG é do projeto. O formato da OCORRÊNCIA é do Anchors, pequeno e sobre o
qual ele tem jurisdição:

    [{"rule": "CRED-E01", "count": 142, "first": "...", "last": "...", "context": {...}}]

### Camada 3 — DEPOIS: de onde a falha vem

É o cerne do pedido, e a parte nova.

Um `try/catch` prevê a falha **sem saber a origem**. A spec declara *"pode falhar ao
consultar o saldo"* e não diz por quê — porque quem escreveu também não sabia. A falha é
mapeada como possível, não como compreendida.

Depois que a aplicação roda, isso muda: **o log tem contexto**. E é o contexto acumulado de
muitas ocorrências que responde o que a spec não sabia responder.

Três desfechos, e todos são progresso:

| desfecho | o que significa | o que se faz |
| --- | --- | --- |
| **CAUSA IDENTIFICADA** | as ocorrências têm um padrão no contexto | a `-E` ganha a causa escrita, e o tratamento pode ser específico em vez de genérico |
| **RESILIENTE** | a falha é compreendida e o fluxo a absorve (o parceiro devolve nulo em migração, o cliente cancelou) | marca-se `@resilient: <razão>`, e ela PARA DE ALARMAR |
| **AINDA NÃO SEI** | ocorre, e o contexto não revelou padrão | fica registrado que está sob observação — e isso é diferente de ninguém ter olhado |

O terceiro desfecho é o que falta em toda ferramenta de observability: a diferença entre
*"ninguém investigou"* e *"investigamos e ainda não sabemos"*. A segunda é conhecimento, e
hoje ela se perde.

**O "parar de alarmar" é o mais valioso no dia a dia.** Uma falha classificada como fluxo
aceitável sai do radar sem sair do registro — e é o oposto de silenciar um alerta, porque a
decisão fica escrita na spec, com quem decidiu e por quê.

### A marcação: `@resilient`

Uma falha compreendida e tratada não precisa continuar alarmando. Ela é marcada, e sai do
radar sem sair do registro:

    | `CRED-E01` | saldo insuficiente | recusa com 402 | @resilient: o parceiro devolve nulo quando a conta está em migração; tratado desde 2026-09 |

É a mesma convenção dos opt-outs que já existem, e uma AFIRMAÇÃO diferente das duas:

    @no-<coisa>: <razão>   "não vai ter"       — dispensa permanente
    @TBD: <razão>          "ainda não tem"     — dívida, continua aparecendo
    @resilient: <razão>    "acontece, eu sei por quê, e está tratada"

A terceira é a que faltava. Não é dispensa (a falha é real e acontece), não é dívida (não
há nada pendente) — é conhecimento adquirido, e a razão escrita é o que a distingue de
simplesmente calar um alerta.

Marcador nu não vale, pela mesma regra dos outros: `@resilient` sem o porquê seria o
silêncio que os gates existem para acabar. A razão é obrigatória, e é ela que responde à
pergunta que alguém vai fazer daqui a seis meses — *"por que ignoramos isso?"*.

## O que o Anchors faz, e o que ele NÃO faz

Ele NÃO analisa o log. Isso é observability, e existe ferramenta boa.

O que ele faz é o que nenhuma delas faz: **fechar o circuito entre a falha DECLARADA e a
falha OBSERVADA**, e guardar o que se aprendeu de volta na spec — onde quem for mexer no
código amanhã vai ler.

A análise de causa é trabalho de quem entende o domínio (ou de uma IA, contra o contexto).
O Anchors dá o alvo, o material e o lugar de registrar a conclusão.

## Decisões em aberto

| Código | Pergunta | Quem decide | Vira |
| --- | --- | --- | --- |
| `Q01` | A camada 1 entra sozinha primeiro? Ela é confronto estático puro — não depende de produção, de log nem de formato nenhum, e é pré-requisito das outras duas: sem `failure-logged`, a camada 2 não tem o que ingerir. | você | o escopo da primeira rodada |
| `Q02` | Como o dialeto reconhece "tratamento" e "log"? O `GuardPatterns` já existe para guarda de parâmetro; o mesmo mecanismo serviria (`catch_patterns`, `log_patterns`), declarado por projeto. | você | dois campos novos em `dialect` |
| `Q03` | A ponte log→regra é o código CARIMBAR o identificador (`{rule: "CRED-E01"}`) ou um padrão declarado por regra? A primeira é exata e exige tocar o código; a segunda adota o legado e erra como toda heurística. | você | se o `anchors.yaml` ganha `errors.match` |
| `Q04` | Onde vive a conclusão da camada 3? Uma coluna nova na tabela `-E` (`causa`, `aceitável`, `sob observação`) mantém tudo num lugar só — mas mistura o que foi DECIDIDO com o que foi DESCOBERTO. | você | o formato da seção `Errors` |
