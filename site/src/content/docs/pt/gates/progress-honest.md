---
title: "Gate: progress-honest"
description: "Garante que o arquivo de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco."
---

> **Identificador do Gate:** `progress-honest` / `progresso-honesto`  
> **Código Interno:** `PRHNP`  
> **Categoria:** [Planejamento e Progresso](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Planos` `doc` |
| **Alvos Avaliados (`on`)** | `plan` `doc` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que o arquivo de progresso (ex: checklist de tarefas) diz a verdade sobre os arquivos presentes no disco.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Evita checklists com tarefas marcadas como concluídas [x] quando o arquivo ou teste ainda não existe no repositório.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Cruza itens marcados como concluídos com a existência real dos arquivos e testes no disco.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todos os itens marcados como prontos existem de fato no disco.
- **`✗ Fail` (Reprovado):** Item marcado como concluído no checklist aponta para arquivo inexistente.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: progress-honest
    blocking: true
    measures: "o progresso diz a verdade sobre o disco" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante para manter a gestão transparente.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Crie o arquivo que falta ou desmarque o item [ ] até que o trabalho seja concluído.
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
