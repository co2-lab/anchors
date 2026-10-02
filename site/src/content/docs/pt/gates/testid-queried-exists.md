---
title: "Gate: testid-queried-exists"
description: "Garante que todo testID buscado por um roteiro de fluxo E2E existe de verdade no código."
---

> **Identificador do Gate:** `testid-queried-exists` / `testid-buscado-existe`  
> **Código Interno:** `TQETS`  
> **Categoria:** [Apresentação e Telas](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `E2E` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todo testID buscado por um roteiro de fluxo E2E existe de verdade no código.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita testes E2E falhando por buscar elementos que foram renomeados ou excluídos da tela.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Examina as queries `getByTestId` nos testes E2E e verifica a presença da string no código dos componentes.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os testIDs buscados existem no código fonte.
- **`✗ Fail` (Reprovado):** Teste E2E buscando testID que não existe em nenhum componente.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: testid-queried-exists
    on: [test]
    check: testid-queried-exists
    blocking: true
    measures: "todo testID buscado por fluxo E2E existe no código" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante no CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione o `testID` no componente de destino ou corrija a query no teste.
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
