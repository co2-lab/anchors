---
title: "Gate: contract-impact"
description: "Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes."
---

> **Identificador do Gate:** `contract-impact` / `impacto-de-contrato`  
> **Código Interno:** `CTRIM`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Quando um campo de contrato é alterado, identifica as regras que o usam e roda seus testes.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Garante que mudanças de contrato disparem exatamente os testes que dependem do campo alterado.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê o git diff, identifica campos de contrato alterados e cruza com a tabela de usos de regras.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as regras afetadas pela alteração de campo tiveram seus testes executados.
- **`✗ Fail` (Reprovado):** Campos de contrato mudaram e regras dependentes não foram revalidadas.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Sem alteração de contrato no commit.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: contract-impact
    on: [spec]
    check: contract-impact
    blocking: false
    measures: "campo alterado nomeia as regras que o usam e seus testes" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo durante desenvolvimento de APIs e schemas.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Execute os testes das regras impactadas com `anchors test`.
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
