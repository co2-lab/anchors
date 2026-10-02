---
title: "Gate: tests-pass"
description: "Verifica se a suíte de testes passou com zero falhas no relatório ingerido."
---

> **Identificador do Gate:** `tests-pass` / `tests-green`  
> **Código Interno:** `PRJTS`  
> **Categoria:** [Prova e Execução](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Todas as Camadas` `test` |
| **Alvos Avaliados (`on`)** | `test` |
| **Tipo de Verificação** | `Interno Determinístico` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se a suíte de testes passou com zero falhas no relatório ingerido.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Não adianta ter testes se eles estão falhando no pipeline.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê o arquivo JUnit XML gerado pelo runner de testes após a ingestão (`anchors ingest`).

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Zero testes falharam (failed: 0).
- **`✗ Fail` (Reprovado):** Um ou mais testes falharam.
- **`~ Indeterminado/Pending`:** Nenhum resultado de teste foi ingerido ainda.
- **`Skip` (Dispensado):** Arquivos de suporte a testes (helpers).

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: tests-pass
    on: [test]
    check: tests-pass
    blocking: true
    measures: "a suíte de testes passa com zero falhas" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante obrigatório em pre-commit, pre-push e CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Corrija o código que causou a quebra do teste e rode a suíte novamente.
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
