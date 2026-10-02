---
title: "Gate: docs-fresh"
description: "Garante que a documentação compilada reflete as especificações atuais sem defasagem."
---

> **Identificador do Gate:** `docs-fresh` / `documentos-frescos`  
> **Código Interno:** `DCFRD`  
> **Categoria:** [Falhas e Governança](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `doc` `docs` |
| **Alvos Avaliados (`on`)** | `doc` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que a documentação compilada reflete as especificações atuais sem defasagem.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Se a spec mudou e a documentação compilada não foi regerada, o site de documentação mente para os usuários.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se os arquivos gerados em docs/ são mais novos do que as specs das quais derivam.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A documentação compilada está em dia com as specs.
- **`✗ Fail` (Reprovado):** A documentação compilada está desatualizada (stale).
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Projetos sem compilação de docs.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: docs-fresh
    blocking: true
    measures: "a documentação compilada reflete a spec atual" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em CI.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Recompile a documentação com `anchors docs build`.
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
