---
title: "Gate: spec-realizes-doctrine"
description: "Garante que a spec declara e comprova como ela realiza a doutrina."
---

> **Identificador do Gate:** `spec-realizes-doctrine` / `spec-realiza-doutrina`  
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

Garante que a spec declara e comprova como ela realiza a doutrina.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Não basta dizer `@realizes`: a spec precisa descrever a regra local correspondente.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Verifica se a linha que carrega `@realizes` possui contexto funcional associado.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A realização da doutrina é contextualmente válida.
- **`✗ Fail` (Reprovado):** Anotação `@realizes` solta sem regra correspondente.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Não se aplica.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: spec-realizes-doctrine
    on: [spec]
    check: spec-realizes-doctrine
    blocking: true
    measures: "a spec comprova como realiza a doutrina" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em projetos com doutrina ativa.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Associe a anotação `@realizes` ao cabeçalho ou linha da regra que a implementa.
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
