---
title: "Gate: has-code"
description: "Garante que o arquivo possui um código de identidade de cenário e regra."
---

> **Identificador do Gate:** `has-code` / `spec-tem-codigo`  
> **Código Interno:** `SCIDS`  
> **Categoria:** [A Unidade e Estrutura](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `Todas as Camadas` |
| **Alvos Avaliados (`on`)** | `spec` `feature` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o arquivo possui um código de identidade de cenário e regra.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Sem código de identidade (ex: AUTH-B01), os requisitos ficam invisíveis para os verificadores de rastreabilidade.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Examina o conteúdo procurando padrões que casam com a gramática de códigos de regras e cenários do projeto.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** O arquivo contém ao menos um código no formato esperado (ex: PREFIX-B01).
- **`✗ Fail` (Reprovado):** O arquivo não contém nenhum código de regra ou cenário.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Arquivos binários ou da camada declarativa sem regras.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: has-code
    on: [spec, feature]
    check: has-code
    blocking: true
    measures: "o arquivo carrega um código de cenário" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante desde o primeiro dia. Sem código de identidade, nenhuma rastreabilidade funciona.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Adicione um código formal para as regras no arquivo (ex: `### AUTH-B01 — Título da Regra`).
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
