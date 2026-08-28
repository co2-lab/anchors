# Guia de spec — como escrever um `.spec.md` neste projeto

> Semeado por `anchors init`. É a régua embutida (`anchors guide spec`)
> instanciada com o dialeto DESTE projeto. Leia antes de escrever qualquer spec.

## Comece pelo comando, não pelo texto

Não escreva a spec do zero. O CLI emite o esqueleto já conforme:

```sh
anchors new spec <Nome> --out <caminho>/<Nome>.spec.md   # gera o esqueleto
anchors new spec --list-sections                          # as seções e QUANDO usar cada
```

O comando resolve o que é mais fácil errar: gera o código de identidade, escreve o
cabeçalho `@anchors` no dialeto certo, e usa o formato exato de regra catalogada.
Depois preencha e confronte com `anchors check --changed <arquivo>`.

## O formato que o gate exige

Uma regra é **catalogada** quando tem código E lugar estruturado. Três formas
valem, e menção solta em prosa NÃO conta:

```md
### ANCH-B01 — descrição da regra          <- cabeçalho (preferido)
| `ANCH-B02` | descrição |                     <- linha de tabela
- **ANCH-B03** descrição                       <- bullet-negrito
```

## Exemplo completo (copie e adapte)

```md
<!-- @anchors
  code: ANCH
  updated_at: 2026-01-15
  layer: screen
-->
# Login — autentica o usuário e o leva ao app

> **Código**: `ANCH`

## Visão Geral

Tela de entrada: recebe e-mail e senha, autentica, e navega para a Home.

## Regras

### ANCH-S01 — Estado inicial
Campos vazios, botão de entrar desabilitado.

### ANCH-A01 — Entrar com credenciais válidas
Autentica e navega para a Home.

### ANCH-V01 — E-mail inválido
O campo mostra a mensagem e o submit não dispara.

### ANCH-R01 — Só anônimo acessa
Sessão ativa é redirecionada para a Home.

## Decisões em aberto

| Pergunta | Quem decide | Vira |
| --- | --- | --- |

nenhuma
```

## A letra do código diz a NATUREZA da regra

Este projeto não declara `rule_types`, então valem as letras canônicas do
framework: `S` estado, `R` permissão, `V` validação, `A` ação, `X` restrição,
`B` comportamento, `N` navegação, `M` mensagem, `D` dado. Declarar as suas em
`rule_types` faz o vocabulário do time valer no lugar do genérico.

## As seções

Três são obrigatórias — cabeçalho, visão geral e regras — mais as decisões em
aberto. As demais entram com `--with <chave>` quando a unidade pede. Rode
`anchors new spec --list-sections` para a lista com o critério de escolha de cada
uma; ela inclui as ALTERNATIVAS mutuamente exclusivas (`contract` ou `signature`,
`rules` ou `effects`), que é onde a escolha errada custa reescrita.

## O que NÃO fazer

- **Regra sem código.** Sem identidade, a feature e o teste não têm o que citar —
  a trinca não fecha e os gates relacionais ficam sem alvo.
- **Descrever implementação.** A spec diz o COMPORTAMENTO; o nome da função e a
  biblioteca mudam sem a regra mudar.
- **Repetir a copy.** O texto ao usuário mora uma vez (na seção de mensagens); as
  outras seções referenciam o código dela.
- **Chutar o que está ambíguo.** Vira linha em *Decisões em aberto* — o gate
  `open-questions-resolved` cobra que alguém decida, e é isso que se quer.

## Especialize este arquivo

Ele nasce genérico. Conforme o projeto firma convenções (perfis de spec por tipo
de unidade, exemplos reais, casos CORRETO/ERRADO tirados do próprio repo), edite
aqui — é a régua DESTE projeto, e o `governs` do `anchors.yaml` liga este guide aos
alvos que ele rege.
