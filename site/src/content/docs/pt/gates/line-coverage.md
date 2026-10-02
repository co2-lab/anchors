---
title: "Gate: line-coverage"
description: "Verifica se a cobertura de linhas do arquivo atinge o piso mínimo exigido (ex: >= 70%)."
---

> **Identificador do Gate:** `line-coverage` / `cobertura-de-linha`  
> **Código Interno:** `INCHN`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` `code` |
| **Alvos Avaliados (`on`)** | `code` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Informativo no início (`blocking: false`) |

---

## 🎯 O que este gate mede?

Verifica se a cobertura de linhas do arquivo atinge o piso mínimo exigido (ex: >= 70%).

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Garante que as linhas de código foram pelo menos exercitadas pela suíte durante os testes.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê o relatório LCOV ingerido e calcula a porcentagem de linhas cobertas daquele arquivo.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Cobertura de linha maior ou igual ao limiar configurado.
- **`✗ Fail` (Reprovado):** Cobertura de linha abaixo do mínimo exigido.
- **`~ Indeterminado/Pending`:** Nenhum relatório de cobertura foi ingerido.
- **`Skip` (Dispensado):** Arquivos sem linhas executáveis.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: line-coverage
    on: [code]
    check: line-coverage
    blocking: false
    measures: "a cobertura de linhas atinge o piso mínimo" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Informativo no início; aumente o rigor conforme a base estabiliza.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Escreva testes unitários cobrindo as linhas e ramos que não foram executados.
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
