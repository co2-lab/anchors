# Guia de cabeçalho — projeto

> O bloco de marcações no topo de CADA arquivo deste projeto. Semeado por
> `anchors init`; é a régua embutida (`anchors guide header`) instanciada para a
> stack aqui. Mandatório: um arquivo sem este cabeçalho é invisível ao que o
> Anchors sabe fazer melhor.

## O bloco, no dialeto desta stack

No CÓDIGO/teste/feature (referencia a unidade da spec):

```
// @anchors
//   ref: LGNN             # referencia a unidade dona (a spec); NÃO é posse
//   updated_at: 2026-08-08 # dia da última alteração (o gate confere vs. git)
//   layer: screen         # camada da Estrutura (normalmente deduzida do caminho)
//   @feature: auth
```

Na SPEC (a DONA da identidade):

```
// @anchors
//   code: LGNN            # a spec POSSUI o código
//   updated_at: 2026-08-08
```

## As marcações

- `code:` — POSSE da identidade (a SPEC é a dona). `ref:` — REFERÊNCIA (code/
  feature/test apontam a unidade da spec; pode ser múltiplo: `ref: A, B`). Todo
  arquivo precisa de um dos dois. Gere o código com `anchors code <nome>`.
- `updated_at:` — o dia da última alteração. Quem altera atualiza; o gate
  `updated-at-atual` confere contra o git (só ano-mês-dia) e `anchors check --fix`
  corrige. NÃO invente a data — deixe bater com o commit.
- `layer:` — a camada; normalmente deduzida do caminho, declare só p/ sobrepor.
- `@feature: <nome>` — o módulo/feature vertical. 
- `@noPropagation`, `@anchors-shared-code` — opt-outs honestos (sempre com o porquê ao lado).

## Regras

- Sempre no TOPO do arquivo.
- `code` é o mínimo obrigatório (gate `header-conforme`).
- `updated_at` bate com o dia do último commit (gate `updated-at-atual`; `--fix` conserta).
- Opt-out sempre com um porquê ao lado.

_(Régua completa e universal: `anchors guide header`.)_
