---
title: "Gate: code-reference-valid"
description: "Garante que referências cruzadas a códigos de outras regras apontam para regras que realmente existem."
---

> **Identificador do Gate:** `code-reference-valid` / `referencia-de-codigo-valida`  
> **Código Interno:** `CRVCD`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `spec` `feature` `code` `test` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que referências cruzadas a códigos de outras regras apontam para regras que realmente existem.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Quando uma spec cita '@realizes REQ-01' ou '@ref: AUTH-B01', essa regra precisa existir. Caso contrário, é um link quebrado.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Varre o repositório indexando todos os códigos válidos e verifica se todas as menções cruzadas resolvem para um código existente.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as referências a códigos apontam para regras catalogadas no projeto.
- **`✗ Fail` (Reprovado):** Foi encontrada referência a um código que não existe em nenhuma spec.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: code-reference-valid
    blocking: true
    measures: "as referências a códigos apontam para regras que existem" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante obrigatório para manter o grafo íntegro.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Corrija o código digitado ou crie a regra correspondente na spec dona.
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
