---
title: "Gate: evidence-fresh"
description: "Garante que o placar do teste continua fresco e foi medido contra o código atual."
---

> **Identificador do Gate:** `evidence-fresh` / `evidencia-fresca`  
> **Código Interno:** `EVFRV`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `test` `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o placar do teste continua fresco e foi medido contra o código atual.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede que o pipeline aprove commits usando relatórios de testes antigos medidos antes de alterações recentes.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Compara o hash ou data de modificação do teste ingerido com a versão atual dos arquivos no disco.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** O relatório de execução foi gerado na mesma revisão do código atual.
- **`✗ Fail` (Reprovado):** O código ou a spec mudaram após a última execução do teste.
- **`~ Indeterminado/Pending`:** Relatório ausente.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: evidence-fresh
    blocking: true
    measures: "o resultado do teste vale contra o código de hoje" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para impedir testes stale.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Execute a suíte de testes e ingira o novo relatório com `anchors ingest`.
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
