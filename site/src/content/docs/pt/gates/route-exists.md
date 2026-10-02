---
title: "Gate: route-exists"
description: "Verifica se a rota declarada na especificação existe no registro de rotas da aplicação."
---

> **Identificador do Gate:** `route-exists` / `rota-existe`  
> **Código Interno:** `RTEXR`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `UI` `Telas` `API` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se a rota declarada na especificação existe no registro de rotas da aplicação.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede que uma spec prometa a rota `/checkout/pix` e a aplicação registre `/pagamento/pix`.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Cruza a string da rota na spec com o arquivo de rotas central da aplicação (declarado em router_file).

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A rota existe no arquivo de rotas da aplicação.
- **`✗ Fail` (Reprovado):** A rota não foi encontrada no registro de rotas do código.
- **`~ Indeterminado/Pending`:** Registro de rotas pendente.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: route-exists
    on: [spec]
    check: route-exists
    blocking: true
    measures: "a rota declarada na spec existe no registro da aplicação" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em frontend e backend.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Registre a rota no roteador da aplicação ou corrija a spec.
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
