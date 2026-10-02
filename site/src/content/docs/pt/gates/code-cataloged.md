---
title: "Gate: code-cataloged"
description: "Garante que todo símbolo público (função, tipo, export) exportado pelo código está catalogado na spec."
---

> **Identificador do Gate:** `code-cataloged` / `codigo-catalogado`  
> **Código Interno:** `CDCTC`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `code` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todo símbolo público (função, tipo, export) exportado pelo código está catalogado na spec.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede o surgimento de código fantasma ou funções públicas criadas sem especificação que ninguém sabe o que fazem.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Analisa o código procurando funções/tipos exportados e confronta com a spec que governa o arquivo.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os símbolos exportados têm regra correspondente na spec ou dispensa com @no-rule.
- **`✗ Fail` (Reprovado):** Função pública exportada sem nenhuma menção na spec dona.
- **`~ Indeterminado/Pending`:** Spec ainda não criada.
- **`Skip` (Dispensado):** Camadas reconhecidas (infra, dao, types) ou arquivos com regime declarativo.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: code-cataloged
    blocking: true
    measures: "todo símbolo exportado tem regra na spec ou dispensa" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em bibliotecas e serviços de negócio.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Documente a função na spec correspondente ou marque com @no-rule no código explicando o motivo.
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
