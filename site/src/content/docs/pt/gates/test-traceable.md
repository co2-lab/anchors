---
title: "Gate: test-traceable"
description: "Garante que todo teste ligado a uma feature declara no título o código do cenário que prova."
---

> **Identificador do Gate:** `test-traceable` / `teste-rastreavel`  
> **Código Interno:** `TSTRT`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que todo teste ligado a uma feature declara no título o código do cenário que prova.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Sem o código no título do teste, o runner não consegue ligar o resultado da execução ao requisito da spec.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Examina as funções e blocos it/test do arquivo de teste procurando códigos de cenário válidos.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os testes declaram o código do cenário que comprovam.
- **`✗ Fail` (Reprovado):** Teste sem menção a código de cenário ou citando código que não existe.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Arquivos de teste em camadas declarativas ou testes auxiliares.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: test-traceable
    blocking: true
    measures: "o teste cita o código do cenário que prova" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para fechar a fiação de rastreabilidade.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Inclua o código do cenário no nome da função de teste (ex: `TestAuth_AUTH_B01_Bloqueio`).
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
