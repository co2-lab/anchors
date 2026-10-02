---
title: "Gate: mock-stamped"
description: "Garante que todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa."
---

> **Identificador do Gate:** `mock-stamped` / `mock-carimbado`  
> **Código Interno:** `MCSTM`  
> **Categoria:** [Dublês e Mocks](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todo dublê de teste carrega a marca @contract do snippet que substitui, e o gate a recomputa.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Quando o código real de um serviço muda, o mock precisa ser atualizado; se o hash divergir, o mock mentiu.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica a presença do carimbo `@contract <hash>` no mock e compara o hash com a assinatura do método real substituído.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** O carimbo bate com o hash da assinatura do método original.
- **`✗ Fail` (Reprovado):** Mock sem carimbo ou hash desatualizado (o método real mudou).
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Mocks declarados em camadas de suporte com dispensa.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: mock-stamped
    on: [test]
    check: mock-stamped
    blocking: true
    measures: "o dublê carrega a marca do snippet que substitui" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para impedir mocks que mentem.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Atualize o mock para refletir o novo contrato e regenere o carimbo com `anchors check --fix`.
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
