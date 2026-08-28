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

## 7. O que existe hoje

- `config.Workflow` com `Mode`/`Repo`/`Labels` e `Config.ModoGitHub()`
- validação completa no `Load`, com 7 testes que cobrem cada mensagem de erro
- **nada mais**: nenhum comando consulta o GitHub ainda
