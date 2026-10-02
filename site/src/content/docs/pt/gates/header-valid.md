---
title: "Gate: header-valid"
description: "Todo arquivo regido pelo Anchors carrega o cabeçalho @anchors no topo, com a sua identidade."
---

> **Identificador do Gate:** `header-valid`  
> **Código Interno:** `INCHN` | **Categoria:** [A Unidade e Estrutura](/pt/docs/gates/)  
> **Avalia:** `spec`, `feature`, `code`, `test`, `guide`, `doc`, `plan`, `product`, `flag` | **Camadas:** `Todas`  
> **Execução:** `Determinística, interna` | **Reparável:** `anchors check --fix`

---

## 1. O que este gate mede

Todo arquivo regido pelo Anchors carrega o cabeçalho `@anchors` **no topo**, com a sua identidade.

---

## 2. Por que importa

O cabeçalho é como o arquivo entra no mapa: diz a que unidade ele pertence. Um arquivo sem ele fica invisível para os gates que leem a unidade de um arquivo — o alcance de um teste, as referências de um código, a camada que ele declara.

---

## 3. Como funciona

O cabeçalho é o bloco `@anchors` no topo do arquivo, lido exatamente como o mapa lê: antes dele, só linhas em branco, comentários e um shebang. Um `@anchors` mais abaixo — um exemplo num guia, uma string no código — é texto, não o cabeçalho. Um comentário só abre o cabeçalho quando `@anchors` é a primeira palavra dele.

A identidade exigida depende do papel do arquivo:

| Arquivo | Identidade |
| --- | --- |
| spec, plano, doutrina de produto, flag | `code:` — é dono da própria identidade |
| código, teste, feature | `ref:` — o código da spec a que pertence (`ref: A, B` quando compartilhado) |
| guia, documento, arquivo de suporte de teste | `layer:` — não pertence a uma unidade |

O cabeçalho é lido em todo dialeto de comentário: `//`, `#`, `--`, `<!-- -->` e o ` * ` de um comentário de bloco.

---

## 4. Vereditos

| Veredito | Condição | Ação |
| --- | --- | --- |
| **`PASS`** | O cabeçalho está no topo, com a identidade que o papel do arquivo exige. | Nenhuma. |
| **`FAIL`** | Sem cabeçalho no topo, cabeçalho abaixo do topo sem `@fixed-header: <motivo>`, ou cabeçalho sem identidade. | Rode `anchors check --fix`, ou escreva o cabeçalho. |
| **`SKIP`** | Arquivos binários e roteiros de runner externo. | Nenhuma. |

---

## 5. Exemplo de configuração (`anchors.yaml`)

```yaml
gates:
  - name: header-valid
    blocking: true
```

Uma entrada só com o nome herda o `on:` canônico — todos os tipos regidos.

---

## 6. Como corrigir

O `anchors check --all --fix` escreve o cabeçalho que falta, a partir do mapa: o `ref:` da unidade do arquivo — a spec que o especifica, ou, pela feature, a spec de um teste — ou o `layer:` de um guia, documento ou arquivo de suporte de teste. Um cabeçalho sem identidade ganha a linha que falta; nada escrito é mudado. Um shebang continua primeiro, e o `//go:build` do Go continua antes do `package`.

Um arquivo que o mapa não liga a unidade nenhuma fica como está: a identidade dele é uma decisão a tomar.

```go
// @anchors
//   ref: FATUR

package billing
```

---

## 7. Conceitos e gates relacionados

- [A Unidade](/pt/docs/concepts/unit/): a spec e todo irmão que ela rege.
- [Camadas do Projeto](/pt/docs/layers/): fronteiras de arquitetura e regimes de teste.
- [Gates e Vereditos](/pt/docs/concepts/gates-and-verdicts/): como os checks operam.
- [Catálogo de Gates](/pt/docs/gates/): todos os gates de verificação.
