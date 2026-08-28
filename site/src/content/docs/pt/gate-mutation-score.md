---
title: Gate mutation-score
---

> Referência do gate `mutation-score`: o que ele mede, por que a medida importa,
> como configurá-lo e o que fazer com o resultado. Pressupõe o pilar de
> [Qualidade](/docs/qualidade/), que define o que é um gate.

## O que ele responde

Cobertura de linha diz que a linha **executou**. Não diz que alguém verificou o
resultado dela.

O `mutation-score` responde a pergunta que falta: **altere a linha — o teste
percebe?** Se a suíte continua verde depois da alteração, aquele teste não prova
aquela linha. Ele a executa.

A mecânica é essa, e não é do Anchors: uma ferramenta de mutação altera uma
construção do código de cada vez (troca `>` por `>=`, `&&` por `||`, remove uma
chamada) e roda a suíte contra cada alteração. Cada alteração é um **mutante**.
Mutante que a suíte pega está *morto*; mutante que passa despercebido
*sobreviveu*.

    score = mortos / total

O Anchors **não roda mutação**. Ele lê o relatório no formato aberto
[Mutation Testing Elements](https://github.com/stryker-mutator/mutation-testing-elements)
— o mesmo que Stryker (JS), PIT (Java), Infection (PHP) e mutmut (Python)
emitem — e pendura o sinal no nó do código:

```bash
anchors ingest --mutation reports/mutation/mutation.json
```

Sem ingestão o gate fica `Pending`, não `Fail`: exigir a ferramenta seria o
framework decidindo pelo projeto, e nem todo stack tem uma que rode em tempo
viável. Mas `Pending` não é silêncio — aparece no perfil e no `doctor`, com o
que falta e o que se perde por não ter.

## Os dois escopos, e por que o par diz mais que cada número

A pergunta "o teste percebe?" tem duas leituras, e elas divergem na prática:

| escopo | que suíte roda contra os mutantes | pergunta |
|---|---|---|
| `isolated` | só o teste da própria unidade | **este** teste prova o que a unidade faz? |
| `full` | os testes de todos que a importam | **alguém** no sistema prova? |

Ingerir os dois:

```bash
anchors ingest --mutation iso.json  --scope isolated
anchors ingest --mutation full.json --scope full
```

Medido em duas unidades de projetos reais:

| unidade | isolated | full | delta | dependentes |
|---|---|---|---|---|
| átomo de UI (link com estado de carregamento) | **8%** | 77% | **69p** | 17 |
| função pura de negócio (parcelas) | **27%** | 77% | **50p** | 2 |

Nas duas, olhar só o `full` dava a impressão de saúde — 77% acima de qualquer
limiar usual. O par conta outra história: **a maior parte dos mutantes só morre
nos dependentes**. O teste da unidade executa o código e quase não verifica o
resultado; quem prova são os outros arquivos.

Isso importa porque muda o que a suíte protege. Um refactor no átomo não é
guardado pelos testes do átomo — o autor descobre o erro quando quebra uma tela,
longe de onde mexeu.

Note que o segundo caso é uma função **pura, com dois dependentes**. O padrão não
é da UI; é de qualquer unidade cujo teste afirme pouco.

### Como ler o delta

    delta = full − isolated

| delta | leitura | onde está o conserto |
|---|---|---|
| **≥ 40p** | a suíte completa cobre muito mais | o buraco é no teste da unidade: ele não assere o que o código produz |
| **15–40p** | parte da prova vem dos dependentes | acoplamento parcial; vale trazer as asserções para casa |
| **< 15p** com score baixo | os dois escopos concordam | não é acoplamento: falta **asserção**, e nenhum teste do sistema prova aquilo |

A distinção do último caso é o que evita procurar no lugar errado. Delta baixo
com score baixo significa que rodar mais testes não vai ajudar — a linha não tem
prova em lugar nenhum.

### Qual escopo dá o veredito

O **isolado**. Um `full` alto não compensa uma unidade que não se prova sozinha,
porque a pergunta do gate é sobre o teste daquela unidade.

Com só um escopo ingerido (ou nenhum), o gate usa o total da última ingestão — o
comportamento anterior, para quem já ingere hoje.

## Custo, e em que fase cada escopo cabe

Mutação é cara por construção: a suíte roda uma vez por mutante. O custo por
mutante varia por **duas** ordens de grandeza conforme o que a unidade arrasta:

| tipo de unidade | isolated | full |
|---|---|---|
| função pura, poucos dependentes | **~0,5s**/mutante | ~0,7s/mutante |
| componente de UI compartilhado | **~4,4s**/mutante | ~9,1s/mutante |

A diferença não é da ferramenta de mutação — é de **quantos testes tocam a
unidade**. Um átomo importado por 17 arquivos arrasta os testes dos 17; uma
função pura arrasta os seus. E, num stack de UI, cada worker carrega o runtime
do framework antes de rodar o primeiro teste — medido, isso foi 44% do tempo de
uma rodada.

Daí o mapeamento de fase:

| fase | o que roda | por quê |
|---|---|---|
| `pre-commit` | isolado, **só em unidades baratas** (lógica pura) | o hook precisa de segundos |
| `pre-push` | isolado, nas unidades tocadas | minutos são toleráveis antes de publicar |
| `ci` | **completo**, e é onde o delta existe | a máquina espera; o humano não |
| `manual` | ambos, na unidade em que se está trabalhando | quando você quer saber |

Prescrever mutação em toda etapa de teste torna o ciclo inviável. Rodar na
unidade que importa é o uso que paga.

## Limiar: o que significa subir o `break`

A ferramenta de mutação costuma ter um limiar próprio (`thresholds.break` no
Stryker) que faz o comando falhar abaixo de um score. O gate do Anchors tem o
dele. **Manter as duas réguas ativas cria duas fontes discordando sobre o mesmo
número** — escolha onde a decisão vive:

- na ferramenta: falha na hora de rodar, útil em CI isolado;
- no gate: entra no perfil junto com os outros, com o resto da rastreabilidade.

Sobre o valor: cada ponto a mais é uma linha que passou a ter alguém verificando
o resultado. O que se paga é o trabalho de escrever essa asserção — e um risco:
se o corte for alto demais para o que a equipe sustenta, a saída fácil é
escrever teste que mata mutante sem provar comportamento. Aí o gate mede a si
mesmo.

**100% quase nunca é o alvo.** Existe o *mutante equivalente*: alteração que não
muda o comportamento observável, e que por definição nenhum teste pode matar.
Exemplo real — numa função que escolhe o registro de maior valor, trocar `>` por
`>=` muda qual empate vence, e o campo desempatado é o único que o chamador usa:
o resultado é o mesmo. Perseguir esses gasta esforço sem ganho, e a prática usual
fica entre 70% e 85%.

## O valor está nos sobreviventes, não no percentual

O relatório aponta cada mutante que sobreviveu, **com linha e coluna**. É ali que
está o retorno — o número serve para saber se piorou.

Dois sobreviventes reais de uma única função, encontrados numa rodada de 51
segundos:

```
[Survived] StringLiteral   installments.ts:61:78
```
Um `padStart(2, '0')` cujo `'0'` virava `''`. Nenhum teste usava dia menor que
10, então `2026-08-05` viraria `2026-08-5` sem ninguém notar.

```
[Survived] EqualityOperator   installments.ts:137:28
```
Um filtro `>= start && <= end` cujo `>=` virava `>`. O cenário da spec promete
"filtra por intervalo **inclusive**" — e o teste não provava a inclusividade.

Os dois são buraco de teste, não ruído. Nenhum aparece na cobertura de linha:
as linhas executavam em todos os testes.

## Configuração

```yaml
gates:
  - name: mutation-score
    on: [code]
    blocking: false      # comece informativo: a base ainda não tem sinal
    phase: ci            # o completo é caro; o isolado pode ir a pre-push
```

Duas recomendações que vieram da prática:

**Escope a rodada.** Rodar mutação no projeto inteiro é caro demais para
qualquer ciclo. A ferramenta aceita mutar um arquivo por vez (`--mutate`), e é
assim que ela paga.

**Cuide do modo de sandbox.** Ferramentas que mutam **no lugar** (`inPlace`)
restauram o arquivo ao terminar — mas uma interrupção deixa o código
instrumentado no disco. Medido num monorepo real: uma interrupção deixou 848
arquivos com `// @ts-nocheck` e instrumentação. Comite antes de rodar, e confira
`git status` depois.
