---
title: "Gate: spec-doctrine-exists"
description: "Garante que a doutrina referenciada pela spec com @realizes realmente existe no projeto."
---

> **Identificador do Gate:** `spec-doctrine-exists` / `doutrina-existe`  
> **Código Interno:** `DCTRN`  
> **Categoria:** [Doutrina e Feature Flags](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Garante que a doutrina referenciada pela spec com @realizes realmente existe no projeto.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Impede apontar para uma regra de produto inexistente ou renomeada.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se o identificador após `@realizes` resolve para uma regra nos arquivos de `product/*.doctrine.md`.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** Todas as doutrinas citadas existem no repositório.
- **`✗ Fail` (Reprovado):** Menção a regra de doutrina que não existe.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Specs sem anotação @realizes.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: spec-doctrine-exists
    on: [spec]
    check: spec-doctrine-exists
    blocking: true
    measures: "a doutrina referenciada pela spec existe no projeto" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante imediato.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Crie a regra no arquivo de doutrina ou corrija o código na anotação @realizes.
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
