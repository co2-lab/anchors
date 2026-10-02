---
title: "Gate: mock-typed"
description: "Garante que todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui."
---

> **Identificador do Gate:** `mock-typed` / `mock-tipado`  
> **Código Interno:** `MCTYM`  
> **Categoria:** [Dublês e Mocks](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todo dublê de teste implementa ou deriva formalmente do tipo do módulo que substitui.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede a criação de mocks ad-hoc com métodos inventados que a classe real não possui.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se a estrutura do mock implementa a interface ou estende a classe original.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** O mock é formalmente tipado conforme o módulo real.
- **`✗ Fail` (Reprovado):** Dublê sem tipagem estrita ou divergindo da interface original.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Linguagens sem tipagem estática.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: mock-typed
    on: [test]
    check: mock-typed
    blocking: true
    measures: "todo dublê de teste deriva do módulo que substitui" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em Go, TypeScript, Java e Rust.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Faça o mock implementar a interface oficial do serviço substituído.
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
