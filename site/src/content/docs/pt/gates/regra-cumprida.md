---
title: "Gate: regra-cumprida"
description: "Pergunta ao modelo se o trecho marcado no código realmente realiza o que a regra descreve."
---

> **Identificador do Gate:** `regra-cumprida`  
> **Código Interno:** `RLUEX`  
> **Categoria:** [Julgamento por IA](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` `code` |
| **Tipo de Verificação** | `Julgamento por IA (ask)` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Pergunta ao modelo se o trecho marcado no código realmente realiza o que a regra descreve.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Regexes e linters não sabem se uma função faz o que promete semânticamente. A IA avalia a correspondência real de intenção.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Avalia o prompt semântico comparando o texto da regra na spec com o código entre os marcadores no arquivo fonte.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Veredito `pass` gravado com `anchors judge`.
- **`✗ Fail` (Reprovado):** Veredito `fail` acusando que o código diverge da regra.
- **`~ Indeterminado/Pending`:** Nenhum veredito emitido ainda.
- **`Skip` (Dispensado):** Regra marcada como pendente de implementação (@TBD: code) recebe veredito `dispensado`.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: regra-cumprida
    on: [spec]
    ask: "o trecho marcado REALIZA o que a regra descreve?"
    blocking: false 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo; avalie e registre com `anchors judge` em revisões de PR.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Ajuste o código para cumprir a regra da spec e grave o novo laudo com `anchors judge`.
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
