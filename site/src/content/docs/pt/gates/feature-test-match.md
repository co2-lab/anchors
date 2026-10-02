---
title: "Gate: feature-test-match"
description: "Garante que cada cenário da feature está implementado no teste por código e descrição."
---

> **Identificador do Gate:** `feature-test-match`  
> **Código Interno:** `FTMFT`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Feature` `Test` |
| **Alvos Avaliados (`on`)** | `feature` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que cada cenário da feature está implementado no teste por código e descrição.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede que o desenvolvedor ou a IA escreva um teste que cita o código do cenário, mas testa outra coisa completamente diferente.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê os cenários da feature e procura nos testes ligados tanto a menção ao código quanto trechos significativos da descrição do cenário.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os cenários da feature estão presentes e descritos nos testes correspondentes.
- **`✗ Fail` (Reprovado):** Cenários da feature não encontrados no teste ou descrição divergente.
- **`~ Indeterminado/Pending`:** O arquivo de teste ainda não foi criado ou ingerido.
- **`Skip` (Dispensado):** Cenários com tag @no-test.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: feature-test-match
    on: [feature]
    check: feature-test-match
    blocking: true
    measures: "cada cenário da feature está no teste ligado por código e descrição" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante assim que a suíte de testes estiver rodando.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Atualize a função de teste ou o describe/it para incluir o código e o título do cenário da feature.
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
