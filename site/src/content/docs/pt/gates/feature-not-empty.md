---
title: "Gate: non-empty"
description: "Garante que o arquivo não é um esqueleto vazio e que a feature possui cenários de verdade."
---

> **Identificador do Gate:** `non-empty` / `feature-nao-vazia`  
> **Código Interno:** `FTMFT`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Feature` `Spec` |
| **Alvos Avaliados (`on`)** | `feature` `spec` `doc` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o arquivo não é um esqueleto vazio e que a feature possui cenários de verdade.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Um arquivo .feature com 8 linhas de comentários ou cabeçalho Gherkin não é vazio em bytes, mas não declara cenário nenhum. Este gate pega essas cascas vazias.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se o arquivo tem conteúdo não-espaço e, para features, valida a presença de palavras-chave de cenário (Scenario: ou Cenário:).

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** O arquivo possui conteúdo real e, no caso de feature, contém cenários declarados.
- **`✗ Fail` (Reprovado):** Arquivo vazio, só com espaços ou feature sem nenhum cenário declarado.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: feature-not-empty
    on: [feature]
    check: non-empty
    blocking: true
    measures: "a feature tem cenários de verdade" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante imediato.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Escreva ao menos um cenário concreto com Dado/Quando/Então na feature.
2. Reexecute a verificação no terminal:
   ```sh
   anchors check
   ```
3. Se o gate suportar correção automática, você pode tentar o comando:
   ```sh
   anchors check --fix
   ```

---

## 🔗 Conceitos Relacionados

- [Guia Completo de Camadas](/pt/docs/layers/): Entenda quais camadas exigem este gate.
- [A Unidade (The Unit)](/pt/docs/concepts/unidade/): A relação entre Spec, Feature, Teste e Código.
- [Gates e Vereditos](/pt/docs/concepts/gates-e-vereditos/): A mecânica completa de avaliação do Anchors.
- [Catálogo Completo de Gates](/pt/docs/gates//): Retornar ao índice pesquisável de todos os gates.
