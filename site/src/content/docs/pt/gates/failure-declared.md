---
title: "Gate: failure-declared"
description: "Garante que possíveis falhas da unidade estão declaradas na especificação."
---

> **Identificador do Gate:** `failure-declared` / `falha-declarada`  
> **Código Interno:** `FLRAI`  
> **Categoria:** [Falhas e Governança](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que possíveis falhas da unidade estão declaradas na especificação.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita a ilusão do 'caminho feliz': todo software falha, e as falhas esperadas precisam ser documentadas.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê a seção de Erros/Falhas (`## Errors`) na spec e verifica a catalogação de códigos com prefixo E.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A spec possui catálogo de falhas declaradas.
- **`✗ Fail` (Reprovado):** Spec sem declaração de erros ou falhas possíveis.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Camadas declarativas ou funções puras que não falham.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: failure-declared
    on: [spec]
    check: failure-declared
    blocking: true
    measures: "possíveis falhas da unidade são declaradas na spec" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para usecases e serviços.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione a seção `## Errors` com ao menos uma regra de erro catalogada (ex: `CODE-E01`).
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
