# Anchors — Fluxo de trabalho: local ou GitHub

> **Estado: desenho.** A configuração (`workflow.mode`) e sua validação existem e têm
> teste. Os comandos do modo `github` **ainda não existem**. Este documento registra as
> decisões já tomadas e as perguntas ainda abertas, para que o refinamento aconteça contra
> um alvo escrito em vez de memória.

## 1. O problema

A fila de trabalho do Anchors sempre foi local: o watcher detecta uma mudança e escreve
uma task em `.anchors/tasks/`; `anchors next` puxa e reivindica; `anchors done` fecha.

Isso funciona para uma pessoa numa máquina. Não funciona para um time — ninguém vê a fila
de ninguém, e não há como dividir trabalho entre pessoas e agentes.

O que já estava pronto para a mudança é o mecanismo de posse. O `queue.go` diz, na abertura:

> *"O claim é ATÔMICO via rename (atômico no POSIX), então dois workers em terminais
> diferentes nunca pegam a mesma task. Isso é o que permite paralelismo manual em QUALQUER
> cliente de IA (dois terminais, dois `anchors next`)."*

Ou seja: **vários trabalhadores puxando de uma fila comum já é a arquitetura.** O modo
`github` não inventa esse modelo — muda onde a fila mora, e ganha o claim atômico de graça
(o `assignee` de uma issue).

## 2. A decisão estruturante: modos EXCLUDENTES

```yaml
workflow:
  mode: github            # ou: local
  repo: owner/nome        # só no github, obrigatório
  labels: [anchors]       # só no github, obrigatório
```

Um modo **ou** o outro. Nunca um com o outro de reserva.

A alternativa — tentar o GitHub e cair no local quando falha — parece robustez e é a pior
escolha disponível. Ela cria a pergunta *"de qual fila veio esta task?"*, que ninguém
consegue responder depois do fato, e permite que dois agentes trabalhem no mesmo item cada
um achando que reivindicou.

O precedente é fresco e custou caro: o `login.yaml` do app de referência era um util que *às vezes*
logava e *às vezes* herdava a sessão existente, dependendo de o app já estar na Home. O
sintoma aparecia três telas adiante ("cartão não visível", "acesso restrito") e uma
investigação inteira foi gasta na camada errada. **Comportamento condicional e implícito
desloca o sintoma da causa.** Declarar o modo faz o arquivo responder.

E torna as validações triviais de escrever:
- modo `github` → `.anchors/tasks/` não deve existir;
- modo `local` → nenhum comando toca a rede.

### O que a validação já cobra (na CARGA, não no uso)

Roda no `config.Load`, então falha no primeiro comando que lê o arquivo. Deixar para o
`anchors next` significaria descobrir a configuração quebrada no meio de um trabalho —
depois de o agente já ter alterado arquivos achando que sabia de onde vinha a task.

| exigência | por quê |
|---|---|
| `repo` obrigatório, **nunca** inferido do remote | num fork, inferir faria a escrita cair no repositório errado, e escrita em lugar errado não se desfaz com revert |
| `labels` obrigatório | sem ela, `anchors next` puxaria qualquer issue — inclusive as de produto, que não têm a forma que o ciclo espera |
| `mode: local` **rejeita** `repo`/`labels` | não são campos inofensivos ali: quem lê o arquivo conclui que a integração está ativa |
| `mode` desconhecido falha | um typo não pode virar `local` em silêncio, senão o agente reivindica no lugar errado |
| bloco ausente = `local` | não transformar toda config existente em erro |

## 3. Sem watcher no modo github

Decisão tomada: **no modo `github` não existe watcher.** Quem cria card é humano ou CI; o
agente apenas consome a fila.

As duas alternativas foram descartadas por motivos distintos:

- **Watcher abrindo issue a cada mudança** — um watcher solto abriria dezenas de cards, e
  ruído numa fila que o time inteiro olha custa mais do que a detecção automática rende.
- **Watcher local + gestão remota** — reintroduz "de qual fila veio?", exatamente o que o
  modo excludente existe para eliminar.

**Consequência assumida:** a detecção automática de drift que o watcher dá hoje não
acontece no modo `github`. O `anchors stale` e o `anchors impact` continuam disponíveis
para quem quiser perguntar — a diferença é que ninguém pergunta por você.

## 4. Transporte: `gh` CLI, não SDK

Existe SDK oficial em Go (`github.com/google/go-github`), e a escolha é deliberada contra
ele.

O fator decisivo é **credencial**. Com SDK, o Anchors passa a precisar de um token, e daí
nascem perguntas sem resposta boa: vem do `anchors.yaml`, que é versionado? De env? O que
acontece quando alguém commita o arquivo? Um framework de governança que passa a carregar
segredo ganha uma superfície de risco que hoje não tem.

Com o `gh`, a autenticação é de quem já a resolveu (a pessoa ou o CI), e o Anchors não vê
token nenhum. Rate limit e paginação também deixam de ser problema dele.

Há precedente na própria filosofia do projeto: **o Anchors não roda teste** — ele ingere o
artefato que a ferramenta do projeto gerou (`anchors ingest --junit`). Mesma régua aqui:
não reimplementar o que a ferramenta especializada já faz.

Custo assumido: dependência de runtime. Se o `gh` não estiver instalado ou autenticado, o
modo `github` deve falhar com mensagem clara — mesma régua da validação de config, e não um
erro de rede genérico no meio do trabalho.

## 5. O de-para

| local | github |
|---|---|
| `.anchors/tasks/*.yaml` | issues com a label declarada |
| `State: pending` | aberta, sem assignee |
| `Claim` (rename atômico) | `assignee` = quem pegou |
| `State: done` | fechada + comentário do que foi feito |
| `anchors queue` | `gh issue list --label <label>` |
| `anchors next` | primeira issue sem assignee + assume |
| o gate que fecha a etapa | check do PR |

O `evidence-fresh` encaixa naturalmente como check de PR: ele já responde *"o placar deste
teste mediu código que mudou"*, que é a pergunta que um review faz.

## 6. Cada etapa guiada por card

Direção dada e **ainda não desenhada**: o ciclo de vida (spec → feature → código → teste)
passa a ter **um card por etapa**, e o card é o que autoriza e registra a etapa.

Isso é mais forte do que "a fila virou GitHub". Muda o que o `anchors next` entrega: não é
"um trabalho qualquer", é *"a próxima etapa do ciclo, com o card que a governa"*. Implica
que o card carrega qual etapa é, sobre qual artefato, e o que fecha ela.

Perguntas abertas, a resolver testando num projeto real:

- **Como o card diz qual etapa é?** Label (`etapa:spec`), campo de Projects, ou convenção
  no título? Label é o mais simples de consultar via `gh`; Projects dá mais estrutura.
- **Quem cria o card da etapa seguinte?** O agente ao terminar a anterior (encadeamento
  automático), ou fica explícito para um humano decidir?
- **Um card por etapa ou um card com checklist de etapas?** O primeiro dá paralelismo real
  entre etapas de unidades diferentes; o segundo mantém a unidade de trabalho junta.
- **O que acontece quando um gate reprova?** Reabre o card, comenta e reatribui, ou abre um
  card novo de correção? (No modo local, uma divergência abre issue em `todo/` — vale
  espelhar isso.)
- **`anchors next` deve inicializar projeto novo?** Levantado e não decidido: o primeiro
  comando útil depende de o projeto já estar iniciado ou não.

## 7. Os pipelines se autoverificam

Um pipeline desatualizado falha em **silêncio**. Ele roda, faz o que a versão dele sabia
fazer, e o que foi corrigido depois simplesmente não acontece — sem erro, sem vermelho, sem
nada. O CI fica verde dizendo que fez um trabalho que ele não faz mais.

Isto foi medido, não previsto: um passo novo do `identify` não rodou por três execuções
porque o pipeline instalado era o antigo. Nada acusou. Só apareceu quando fui ler o log
procurando outra coisa — e até ali eu achava que o mecanismo estava quebrado, e quase
"consertei" o que já funcionava.

Os pipelines que **já têm o Anchors instalado** conferem a si mesmos:

| pipeline | confere | por quê |
|---|---|---|
| `identify` | ✅ | roda a cada PR — o caminho quente |
| `board` | ✅ | roda por schedule, e alcança o repositório **parado** |
| `claim`, `pr-checks`, `stale` | — | não instalam o binário |

Os três de fora ficam de fora de propósito: instalar o binário só para conferir os tornaria
mais lentos a cada execução, e o `claim` é chamado a cada pedido de trabalho.

O `board` importa mais do que parece. Ele roda por `schedule`, então alcança o projeto onde
ninguém abre PR — que é justamente onde o pipeline envelhece sem ninguém ver.

### Avisa por padrão, barra se você pedir

```yaml
workflow:
  stale_pipeline_blocks: true   # padrão: false
```

O padrão **avisa** (`::warning::` no resumo do run) e deixa o CI seguir. A razão é que
barrar troca um problema por outro: o pipeline velho ainda faz o trabalho antigo, e
derrubá-lo passa de *"faz menos do que devia"* para *"não faz nada"*.

Mas há projetos onde o aviso não é lido — CI com dezenas de jobs, equipe que não abre o
resumo do run — e ali o silêncio é pior que a interrupção. Quem declara `true` está dizendo
que prefere o CI vermelho a um pipeline que finge funcionar.

Quem decide é o `anchors.yaml`, nunca o YAML do pipeline. Por isso o passo **não** tem
`continue-on-error` (engoliria a escolha de quem quer barrar) e **não** tem `exit 1` escrito
à mão (barraria quem quer só o aviso) — os dois são proibidos por teste.

O mesmo comando serve na sua máquina:

```console
$ anchors doctor --check-pipelines
✓ os 5 pipelines do fluxo estão no lugar e atualizados (anchors v0.1.6).
```

Ele responde **uma** pergunta e sai com o código certo. O `doctor` completo responde
dezenas, e num pipeline isso é ruído: quem chama de lá quer uma resposta, não um raio-X.

## 7.1 O achado vai para onde alguém o vê

No modo local, o achado de gate vira arquivo em `issues/` e o protocolo é mover a
pasta à mão: `todo/` → `doing/` → `done/`. É assim que se pega trabalho sem GitHub.

No modo github isso é duplicação — e a parte grave não é o arquivo sobrando.
Medido: **onze achados** (`trinca-completa`, `rule-types`, `guide-checklist`)
existiam só em arquivo local, e **nenhum** virou card. O board mostrava o trabalho
planejado e escondia o que os gates tinham encontrado. Quem olhasse o board
concluiria que não havia nada a corrigir.

O ciclo de vida é o mesmo, com outro substrato:

| local | github |
|---|---|
| arquivo em `todo/` | issue aberta, com a label do fluxo e `anchors:to-do` |
| mover para `done/` | issue **fechada** |
| reabrir de `done/` | issue **reaberta**, com o novo laudo em comentário |
| `future/` (dívida) | issue **sem** estado de fluxo — fica fora do board até vencer |

A **deduplicação**, que no local vem de procurar o arquivo pelas pastas, aqui vem de
um marcador estável no corpo:

```html
<!-- anchors-issue-key: gate:alvo:kind -->
```

A busca do GitHub é por texto e devolve aproximações — a confirmação é o marcador
exato. Sem ela, um achado sobre `Foo.spec.md` casaria o card de `FooBar.spec.md`, e
o Anchors fecharia o card errado.

O roteamento fica no pacote, não nos chamadores: são oito pontos entre `check` e
`judge`, e espalhar `if modoGitHub` por eles garantiria que o próximo ponto a nascer
esquecesse. O destino **padrão** continua sendo arquivo — quem não configurou nada
não pode acabar falando com a rede sem pedir.

> **Limitação conhecida:** se o **alvo** do achado é removido, o card não fecha
> sozinho. O `Resolve` só dispara quando um gate que reprovou volta a *passar*, e um
> arquivo apagado não deixa nó a confrontar. Vale para o modo local do mesmo jeito.

### O que o gate NÃO vê: o achado do agente

O roteamento acima cobre o que o **gate** detecta. Falta o outro cenário, e ele é o
mais comum: o agente descobre algo errado enquanto implementa outra coisa.

Uma config que contradiz a doutrina do próprio projeto. Um caminho que ninguém
documentou. Um arquivo no lugar errado. **Nenhum gate vê isso** — e o caminho barato
é consertar na hora e seguir, o que faz o conserto sumir do histórico.

Aconteceu ao implementar a primeira spec de código de um projeto: três achados, e só
um virou card — justamente o que um gate detectou. Os outros dois existiam só na
narração de quem os corrigiu.

```console
$ anchors escalate "o jest mede só packages/*/src/, mas o código vive ao lado da spec" \
    --sobre jest.config.js --card 44
achado registrado: https://github.com/acme/projeto/issues/49
```

O card nasce **sob** o de origem, com a label `anchors:sob-44`:

```console
$ gh issue list --label anchors:sob-44
  #50 Nada indica ONDE o código de uma spec deve nascer
  #49 O `jest.config.js` mede cobertura só de packages/*/src/
```

Os dois se entregam no **mesmo PR**: o achado apareceu fazendo aquele trabalho, e
separá-los faria um dos dois esperar sem razão.

**Por que label, e não texto no corpo.** A primeira versão escrevia *"descoberto
durante o card #44"* na descrição, e isso não se consulta — não dá para listar o que
pende sob um trabalho, nem para saber o que precisa entrar no mesmo PR. Label é
filtrável, aparece na lista e sobrevive a qualquer reescrita do texto.

**Por que não a sub-issue nativa do GitHub.** Ela aceita **um** nível, e a hierarquia
deste fluxo tem mais: plano → fase → spec → achado. A árvore do board se monta pelo
`parent:` do header do artefato (§7 do `PLANNING.md`), e é multinível. A label
`sob-` resolve o caso que aquela não alcança — o achado que **não tem artefato**.

### O julgamento é a exceção, e não vira card

A task de julgamento (`.anchors/tasks/`) **não** vai para o GitHub, e a razão é que
ela é trabalho do **commit atual**: quem mexeu no arquivo é quem tem o contexto para
responder. Mandá-la ao board empurraria para outro alguém uma pergunta que só quem
editou sabe responder — e ela sumiria no meio das outras.

Por isso ela é resolvida **antes** de o trabalho sair da máquina:

```console
$ anchors check --changed <arquivo>
✗ barrado — 1 alvo(s) aguardam julgamento.
$ anchors judge <alvo> --gate <g> --verdict pass --reason "..."
$ anchors check --changed <arquivo>     # destravado
```

Barra em `--changed` (a perspectiva do pre-commit) e não em `--all`: o `--all` é a
foto do projeto inteiro, e derrubá-lo por pendência acumulada tornaria impossível
medir um projeto que já tem.

O `.anchors/` inteiro fica no `.gitignore` — ele guarda o que descreve uma
**execução**, não o projeto. Um julgamento pendente num PR passa a ser sinal de que
alguém contornou o gate.

O que descreve o **projeto** nasce fora dali. O SBOM é o exemplo: ele estava em
`.anchors/sbom.json` prometendo, no próprio gate, ser *"gerado e versionável"* —
enquanto era gravado numa pasta destinada a ser ignorada. Hoje nasce na raiz. A
distinção não é *quem gerou* o arquivo, é **o que ele descreve**.

---

## 8. O que existe hoje

- `config.Workflow` com `Mode`/`Repo`/`Labels` e `Config.ModoGitHub()`
- validação completa no `Load`, com 7 testes que cobrem cada mensagem de erro
- **nada mais**: nenhum comando consulta o GitHub ainda
