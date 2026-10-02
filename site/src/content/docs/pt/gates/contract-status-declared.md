---
title: "Gate: contract-status-declared"
description: "Garante que o contrato lista os códigos de status que o código realmente retorna, e apenas esses."
---

> **Identificador do Gate:** `contract-status-declared` / `status-de-contrato-declarado`  
> **Código Interno:** `CSDCN`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `API` `comando` |
| **Alvos Avaliados (`on`)** | `spec` `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o contrato lista os códigos de status que o código realmente retorna, e apenas esses.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita APIs mentirosas que dizem retornar apenas 200 e 400, mas no código lançam 403, 404 e 500 sem documentar.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Extrai retornos HTTP/status do código e confronta com a tabela de status declarada na spec.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Os status retornados no código batem exatamente com os declarados na spec.
- **`✗ Fail` (Reprovado):** Código retorna status não documentado na spec, ou spec documenta status nunca retornado.
- **`~ Indeterminado/Pending`:** Código ainda não implementado.
- **`Skip` (Dispensado):** Camadas sem contratos de saída.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: contract-status-declared
    blocking: true
    measures: "o contrato de saída lista os status reais retornados" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos de API REST / gRPC.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Declare o status na tabela de contrato da spec ou trate o erro no código.
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
