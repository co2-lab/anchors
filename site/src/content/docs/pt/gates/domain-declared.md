---
title: "Gate: domain-declared"
description: "Verifica se a spec declara o que a unidade aceita e quem bloqueia entradas inválidas."
---

> **Identificador do Gate:** `domain-declared` / `dominio-declarado`  
> **Código Interno:** `DMDCD`  
> **Categoria:** [Uso de Regras e Contratos](/pt/docs/gates//)

| Propriedade | Valor |
| --- | --- |
| **Camadas Aplicáveis** | `Camadas Regidas` |
| **Alvos Avaliados (`on`)** | `spec` |
| **Tipo de Verificação** | `Relacional com Grafo` |
| **Modo Recomendado** | Sim (`blocking: true`) |

---

## 🎯 O que este gate mede?

Verifica se a spec declara o que a unidade aceita e quem bloqueia entradas inválidas.

Em termos simples: este gate garante que o seu software não cometa erros por descuido ou falta de sincronização. Se você ou uma inteligência artificial alterar um arquivo coberto por este gate, ele inspeciona o trabalho imediatamente.

---

## 🛡️ Por que isso é importante?

Garante que valores fora do domínio (ex: idade negativa, string nula) tenham tratamento explícito.

Sem este gate ativo, esse tipo de defeito passa despercebido pelos testes comuns e só estoura em produção ou durante refatorações dolorosas semanas depois.

---

## ⚙️ Como funciona por baixo dos panos?

Lê a seção de Domínio da spec e verifica a presença da coluna 'Quem bloqueia' ou tratamento de erro.

### Condições dos Vereditos:

- **`✓ Pass` (Aprovado):** A tabela de domínio define entradas válidas, entradas inválidas e o componente responsável pelo bloqueio.
- **`✗ Fail` (Reprovado):** Tabela de domínio ausente ou sem declaração de quem bloqueia entradas inválidas.
- **`~ Indeterminado/Pending`:** Não se aplica.
- **`Skip` (Dispensado):** Specs puramente visuais.

---

## 📋 Exemplo de Configuração no `anchors.yaml`

Para ativar este gate no seu projeto, adicione o bloco abaixo na seção `gates:` do seu [`anchors.yaml`](/pt/docs/anchors-yaml//):

```yaml
gates:
  - name: domain-declared
    on: [spec]
    check: domain-declared
    blocking: true
    measures: "a spec declara o que aceita e quem bloqueia o inválido" 
```

---

## 💡 Recomendações de Uso

- **Quando ativar:** Bloqueante em camadas de usecase e domínio.
- **Fase de execução:** Configure em `when: [pre-commit, ci]` para verificações rápidas, ou `when: [pre-push, ci]` para gates que rodam ferramentas mais pesadas.
- **Transição de maturidade:** Comece com `blocking: false` para avaliar o estado atual do repositório com `anchors check`. Quando zerar as ocorrências, altere para `blocking: true`.

---

## 🔧 Como corrigir quando este gate reprovar?

Se o `anchors check` acusar falha (`✗`) neste gate:

1. Preencha a tabela de Domínio na spec especificando entradas aceitas e quem rejeita o inválido.
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
