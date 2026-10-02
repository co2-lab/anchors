---
title: "Gate: spec-feature-match"
description: "Garante que toda regra catalogada na spec possui ao menos um cenário correspondente na feature."
---

> **Identificador do Gate:** `spec-feature-match`  
> **Código Interno:** `SFMSP`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que toda regra catalogada na spec possui ao menos um cenário correspondente na feature.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita que regras sejam esquecidas no papel sem que ninguém especifique como elas devem ser testadas.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Extrai todos os códigos de regra definidos na spec e confere se cada um é citado como tag (@CODE-B01) na feature ligada.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as regras da spec possuem tags de cenário correspondentes na feature.
- **`✗ Fail` (Reprovado):** Existe alguma regra definida na spec sem cenário na feature.
- **`~ Indeterminado/Pending`:** A feature ligada ainda não foi criada.
- **`Skip` (Dispensado):** Regras marcadas com @no-scenario: motivo ou camadas que dispensam feature.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: spec-feature-match
    on: [spec]
    check: spec-feature-match
    blocking: true
    measures: "cada regra da spec tem cenário na feature" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos onde BDD/Gherkin é parte da governança.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione um cenário na feature com a tag da regra pendente, ou declare @no-scenario na spec com justificativa.
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
